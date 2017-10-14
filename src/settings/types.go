package settings


type ExchangeHandlerClass struct {
	SourceIds map[float64]map[string]float64
	AuctionType float64

	Providers map[string]ProviderParameters
	ProvidersLastChanged string
	UrlToProvider []UrlToProvider

	Caps Caps
	CapsLastChanged string
	WebsiteBlackList map[string][]int

	TMTBlacklistUrls []string
	TMTBlacklistCreativeIds map[string]bool
	TMTBlacklistLastChanged string

	WebsiteRestrictions map[int]WebsiteRestriction
	WebsiteRestrictionsLastChanged string
}

type Caps map[string]int

type WebsiteRestriction struct{
	PassPercentage int
	IsBlocked bool
	BlockedData map[string]bool
	Data map[string]bool
}

type UrlToProvider struct {
	Provider map[string][]string
}

type ProviderParameters struct {
	Hooks map[string]interface{}
	Endpoint string
	Method string
	At float64
	PricingModel string
	UseNative bool
	SupportHttps bool
	IsCached bool
}