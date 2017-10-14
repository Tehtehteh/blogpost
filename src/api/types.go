package api


type RtbResponse struct {
	ProvidersInGame []string `json:"providers_in_game"`
	WinBids map[string]Bid `json:"win_bids"`
	Extra map[string]string `json:"ext"`
	BestPrices map[string]float64 `json:"best_prices"`
}

type WinBid struct {
	Provider string `json:"provider"`
	AdId string `json:"ad_id"`
	Adm string `json:"adm"`
	AdDomain []string `json:"adomain"`
	Attributes []string `json:"attr"`
	Categories []string `json:"CAT"`
	Cid string `json:"cid"`
	CreativeId string `json:"crid"`
	Id string `json:"id"`
	ImpressionId string `json:"impid"`
	NUrl string `json:"nurl"`
	Price float64 `json:"price"`
}

type MyDummyResponse struct {
	Id string `json:"id"`
	SeatBids []SeatBid `json:"seatbid"`
	Cur string `json:"cur"`
	BidId string `json:"bidid"`
}

type SeatBid struct {
	Bid []Bid `json:"bid"`
	Seat int `json:"seat"`
}

type Bid struct {
	Provider string `json:"provider"`
	AdId string `json:"ad_id"`
	AdDomain []string `json:"adomain"`
	Adm string `json:"adm"`
	Attributes []string `json:"attr"`
	Categories []string `json:"cat"`
	Cid string `json:"cid"`
	CreativeId string `json:"crid"`
	Id string
	ImpressionId float64 `json:"impid"`
	NUrl string `json:"nurl"`
	Price float64 `json:"price"`
}