package libspot

type Endpoints interface {
	Spclient() []string
	Dealer() []string
	DealerG2() []string
	Accesspoint() []string
}

type EndpointResolver interface {
	Endpoints() (endpoints Endpoints, ok bool)
	Fetch(...ServiceKind) (Endpoints, error)
}

type ServiceKind string

const (
	ServiceKindSpclient    ServiceKind = "spclient"
	ServiceKindDealer      ServiceKind = "dealer"
	ServiceKindDealerG2    ServiceKind = "dealer-g2"
	ServiceKindAccesspoint ServiceKind = "accesspoint"
)
