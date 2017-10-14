package api

import (
	format "fmt"
	"database/sql"

	"github.com/valyala/fasthttp"

	"db"
	"time"
	"encryption"
	. "logger"
	"settings"
	"request"
	"reflect"
	"crypto/md5"
	"encoding/hex"
	_ "auction"
	"encoding/json"
)

func PanicHandler(ctx *fasthttp.RequestCtx, reason interface{}){
	ctx.Error("", 500)
	Log.Errorf("500 Internal server error: %s", reason)
}

func BuyerHandler(ctx *fasthttp.RequestCtx){
	if !settings.ExchangeHandler.Ready() {
		ctx.Error("RTB not ready", 500)
		Log.Error("RTB not ready")
	}
	startExecutionTime := time.Now()
	//publisherId := encryption.Encryption.DecodeId(ctx.UserValue("buyer_id").(string))
	publisherId := encryption.Encryption.DecodeId("14567721834")

	var (
		ParsedRequest *request.Request
		allowed bool = false
	)
	ParsedRequest, allowed = request.Parse(ctx.PostBody())

	if !allowed || !settings.ExchangeHandler.AllowedByWebsiteRestrictions(int(publisherId), ParsedRequest.Site["domain"]){
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(204)
		format.Fprint(ctx, "{}")
		return
	}

	if v, hasKey := settings.ExchangeHandler.SourceIds[ParsedRequest.Extra.WebsiteId]; hasKey{
		//ParsedRequest.Extra.NetworkId = v["network_id"]  // Buyer should be "imonomy"
		ParsedRequest.Extra.PublisherId = v["publisher_id"]
	}

	if reflect.DeepEqual(new(request.Extra), ParsedRequest.Extra) {
		ParsedRequest.Extra.WebsiteId = float64(publisherId)
		ParsedRequest.Extra.UserMatching = make(map[string]string)
		ParsedRequest.Extra.Location = "demand"
		ParsedRequest.Extra.SegmentPrice = true
		ParsedRequest.Extra.PublisherId = float64(publisherId)
		ParsedRequest.Extra.NetworkId = "imonomy"
		ParsedRequest.Extra.Source = "buyer"
	}

	if reflect.DeepEqual(make(map[string]string), ParsedRequest.User) {
		uId := md5.New()
		uId.Write([]byte(ParsedRequest.Device.Ip))
		uId.Write([]byte(ParsedRequest.Device.UserAgent))
		ParsedRequest.User["id"] = hex.EncodeToString(uId.Sum(nil))
	}

	//responses := make(chan map[string][]byte)

	responses := ParsedRequest.Send()

	var DummyResponse = new(RtbResponse)

	DummyResponse.ProvidersInGame = make([]string, len(responses))
	for k := range responses {
		DummyResponse.ProvidersInGame = append(DummyResponse.ProvidersInGame, k)
	}
	var DD = new(MyDummyResponse)
	json.Unmarshal(responses["dummy"].Body(), &DD)
	Log.Info(DD)

	DummyResponse.Extra = map[string]string{}
	DummyResponse.BestPrices = map[string]float64{"1": 2.3}
	DummyResponse.WinBids = make(map[string]Bid)
	DummyResponse.WinBids["1"] = DD.SeatBids[0].Bid[0]

	//log.Println(string(responses["dummy"].Body()))
	//json.Unmarshal(responses["dummy"].Body(), &DD)
	//DummyResponse.WinBids = make(map[string]WinBid, 1)
	//DummyResponse.WinBids[ParsedRequest.Id] = res
	//resp := []byte{}
	resp, _ := json.Marshal(DummyResponse)


	//format.Fprint(ctx, string(resp))
	//log.Println(DD.SeatBids[0].Bid[0].Adm)
	format.Fprint(ctx, string(resp))

	//log.Println("Parsed Request: ", ParsedRequest)

	Log.Printf("POST execution time %s; request %s; publisherId %s", time.Since(startExecutionTime).Seconds(),
		ParsedRequest, publisherId)
}


func HomepageHandler(ctx *fasthttp.RequestCtx){
	var query string = "SELECT * FROM ad_units WHERE idad_units = ?"
	id := ctx.QueryArgs().GetUintOrZero("qeq")
	args := make([]interface{}, 1)
	args[0] = id
	channel := make (chan *sql.Row)
	go db.FetchOne(query, channel, args...)
	post := <-channel
	format.Fprintf(ctx, "Blog post: %s", post)
}