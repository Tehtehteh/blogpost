package settings

import (
	"db"
	"fmt"
	. "logger"
	"database/sql"
	"strings"
	"math/rand"
)

var ExchangeHandler = New()

func New() *ExchangeHandlerClass{
	return &ExchangeHandlerClass{
		TMTBlacklistLastChanged:"1970-01-01",
		CapsLastChanged: "1970-01-01",
		ProvidersLastChanged: "1970-01-01",
		WebsiteRestrictionsLastChanged: "1970-01-01",
	}
}

func (ExchangeHandler *ExchangeHandlerClass) AllowedByWebsiteRestrictions(websiteId int, domain string) bool{
	if domain == "" {
		return true
	}
	domain = strings.Replace(domain, "www", "", 1)
	if _, isInRestrictionList := ExchangeHandler.WebsiteRestrictions[websiteId]; isInRestrictionList{
		isBlocked := ExchangeHandler.WebsiteRestrictions[websiteId].IsBlocked
		if _, in := ExchangeHandler.WebsiteRestrictions[websiteId].Data[domain]; in{
			if isBlocked {
				return true
			} else {
				if _, in := ExchangeHandler.WebsiteRestrictions[websiteId].BlockedData[domain]; in {
					return true
				} else {
					return false
				}
			}
		} else if _, in := ExchangeHandler.WebsiteRestrictions[-1].Data[domain]; ExchangeHandler.WebsiteRestrictions[websiteId].PassPercentage > 0 && rand.Intn(100) <= ExchangeHandler.WebsiteRestrictions[websiteId].PassPercentage && in{
			if _, in := ExchangeHandler.WebsiteRestrictions[websiteId].BlockedData[domain]; in {
				return true
			} else {
				return false
			}
		} else {
			return !isBlocked
		}
	}
	return true
}


func (ExchangeHandler *ExchangeHandlerClass) Init(){
	Log.Info("ExchangeHandler initializing")
	go ExchangeHandler.loadSourceIds()
	go ExchangeHandler.loadProviders()
	go ExchangeHandler.loadTMTBlacklist(ExchangeHandler.TMTBlacklistLastChanged)
	go ExchangeHandler.loadWebsiteRestrictions(ExchangeHandler.WebsiteRestrictionsLastChanged)
	//go ExchangeHandler.loadCaps()
}

//func (ExchangeHandler *ExchangeHandler) loadHooks(){
//	hooks := make(Hooks)
//	Hooks.
//}

func (ExchangeHandler *ExchangeHandlerClass) Ready() bool {
	if len(ExchangeHandler.SourceIds) > 0 {
		return true
	} else {
		return false
	}
}

func (ExchangeHandler *ExchangeHandlerClass) loadWebsiteRestrictions(date string){
	Log.Info("Loading website restrictions")
	channel := make(chan *sql.Rows)

	query := fmt.Sprintf(`CALL get_rtb_website_restrictions('%s')`, date)

	args := make([]interface{}, 0)

	go db.FetchMany(query, channel, args...)

	rows := <- channel
	defer rows.Close()

	var (
		hasRows bool = false
		WebsiteRestrictions WebsiteRestriction
		WebsiteRestrictionsData map[int]WebsiteRestriction = make(map[int]WebsiteRestriction)
	)

	for rows.Next(){
		var resultSet interface{}

		rows.Scan(&resultSet)

		if fmt.Sprintf("%s", resultSet) == "empty" {
			break
		}

		var (
			websiteId int
			typeId int
			isBlocked bool
			value string
			blockedValue string
			passPercentage int
			updatedDate string
		)
		hasRows = true
		rows.Scan(&websiteId, &typeId, &isBlocked, &value, &blockedValue, &passPercentage, &updatedDate)
		ExchangeHandler.WebsiteRestrictionsLastChanged = updatedDate
		WebsiteRestrictions.Data = make(map[string]bool)
		WebsiteRestrictions.BlockedData = make(map[string]bool)
		for _, domain := range strings.Split(value, ",") {
			domain = strings.ToLower(strings.Replace(strings.Replace(domain, "www", "", 1), " ", "", 1))
			WebsiteRestrictions.Data[domain] = true
		}
		for _, domain := range strings.Split(value, ",") {
			domain = strings.ToLower(strings.Replace(strings.Replace(domain, "www", "", 1), " ", "", 1))
			WebsiteRestrictions.BlockedData[domain] = true
		}
		WebsiteRestrictions.IsBlocked = isBlocked
		WebsiteRestrictions.PassPercentage = passPercentage
		WebsiteRestrictionsData[websiteId] = WebsiteRestrictions
	}
	if hasRows {
		ExchangeHandler.WebsiteRestrictions = WebsiteRestrictionsData
		Log.Info("Website Restrictions reloaded.")
	} else {
		Log.Info("No website restrictions in DB")
	}
}

func (ExchangeHandler *ExchangeHandlerClass) loadTMTBlacklist(date string){
	Log.Info("Loading TMT Blacklist")
	channel := make(chan *sql.Rows)

	query := fmt.Sprintf(`CALL sp_ab_blacklist_last_days_get_list('%s')`, date)

	args := make([]interface{}, 0)

	go db.FetchMany(query, channel, args...)

	var (
		hasRows bool = false
		urls []string
		creativeIds = make(map[string]bool)
	)

	rows :=<- channel
	defer rows.Close()

	for rows.Next(){
		var resultSet interface{}

		rows.Scan(&resultSet)

		if fmt.Sprintf("%s", resultSet) == "empty" {
			break
		}
		var (
			id int
			url string
			_type int
			blockedCount int
			isFixed bool
			isActive bool
			createdDate string
			updatedDate string
		)
		hasRows = true
		rows.Scan(&id, &url, &_type, &blockedCount, &isFixed, &isActive, &createdDate, &updatedDate)
		ExchangeHandler.TMTBlacklistLastChanged = updatedDate
		if _type == 1{
			urls = append(urls, url)
		} else {
			creativeIds[url] = true
		}
	}

	if hasRows {
		ExchangeHandler.TMTBlacklistUrls = urls
		ExchangeHandler.TMTBlacklistCreativeIds = creativeIds
		Log.Info("TMT Blacklist reloaded. Urls:", urls, "; Creative ids:", creativeIds)
	} else {
		Log.Info("No new rows for TMT Blacklist")
	}
}

func (ExchangeHandler *ExchangeHandlerClass) loadCaps(){
	Log.Info("Loading caps from database")
	channel := make(chan *sql.Rows)

	query := `SELECT provider_id, country FROM rtb_country_blacklist`


	args := make([]interface{}, 0)

	go db.FetchMany(query, channel, args...)

	rows := <- channel

	var blacklistDict = make(map[int]map[string]bool)
	//var hasRows bool = false

	for rows.Next(){

		//hasRows := true

		var (
			publisherId int
			country string
		)

		rows.Scan(&publisherId, &country)

		blacklistDict[publisherId] = make(map[string]bool)
		blacklistDict[publisherId][country] = true

	}

	query = `SELECT p.id AS 'provider_id', p.blacklist_country AS 'list_type',
			 		UNIX_TIMESTAMP(p.updated) AS 'updated', cc.country AS 'country',
					cc.requests AS 'requests', cc.period AS 'period'
			 FROM rtb_provider AS p
			 LEFT JOIN rtb_cap_country AS cc ON (p.id = cc.provider_id)
			`

	rows.Close()
	go db.FetchMany(query, channel, args...)

	rows =<- channel

	type NestedOne struct{
		updated int
		country struct{
			listType int
			blackWhiteList map[string]bool
			countries map[string]struct{
				requests float64
				period string
			}
		}
	}
	var result = make(map[int]*NestedOne)


	for rows.Next(){

		var (
			providerId, blackListCountryType int
			updated int
			country string
			requests float64
			period string
		)


		rows.Scan(&providerId, &blackListCountryType, &updated, &country,
			&requests, &period)

		result[providerId] = &NestedOne{}
		result[providerId].updated = updated
		//result[providerId].country = make(map[])

	}


}

func (ExchangeHandler *ExchangeHandlerClass) loadProviders(){
	Log.Info("Loading providers from database")
	channel := make(chan *sql.Rows)

	query := `SELECT provider, end_point, method, at, pricing_model, use_native,
			  		 support_https, is_cached
			  FROM rtb_provider
			  WHERE end_point IS NOT NULL AND end_point != "" AND status = 1`

	args := make([]interface{}, 0)

	go db.FetchMany(query, channel, args...)

	var (
		result = make(map[string]ProviderParameters)
		hasRows bool
	)
	rows :=<- channel
	defer rows.Close()

	for rows.Next(){

		var resultSet interface{}

		rows.Scan(&resultSet)

		if fmt.Sprintf("%s", resultSet) == "empty" {
			break
		}

		var (
			providerParams ProviderParameters
			providerName string
		)

		hasRows = true
		rows.Scan(&providerName, &providerParams.Endpoint, &providerParams.Method, &providerParams.At,
				  &providerParams.PricingModel, &providerParams.UseNative, &providerParams.SupportHttps,
			      &providerParams.IsCached)
		result[providerName] = providerParams
	}

	if hasRows {
		ExchangeHandler.Providers = result
		Log.Info("Providers reloaded: ", result)
	} else {
		Log.Info("No new providers in database")
	}
}

func (ExchangeHandler *ExchangeHandlerClass) loadSourceIds(){
	Log.Info("Loading source ids")
	channel := make(chan *sql.Rows)

	query := `SELECT w.id as website_id, w.publisher_id AS publisher_id, p.network_id AS network_id
			  FROM websites  AS w
	 		  LEFT JOIN publishers AS p ON (w.publisher_id = p.id)`

	args := make([]interface{}, 0)

	go db.FetchMany(query, channel, args...)
	var (
		result = make(map[float64]map[string]float64)
		hasRows bool = false
	)

	rows := <-channel
	defer rows.Close()

	for rows.Next() {
		hasRows = true
		var (website_id float64
			 publisher_id float64
			 network_id float64)
		rows.Scan(&website_id, &publisher_id, &network_id)
		result[website_id] = make(map[string]float64)
		result[website_id]["publisher_id"] = publisher_id
		result[website_id]["website_id"] = website_id
	}

	if hasRows {
		ExchangeHandler.SourceIds = result
		Log.Info("Source ids reloaded: ", result)
	} else {
		Log.Info("No new source ids from database")
	}

}