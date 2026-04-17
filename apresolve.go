package libspot

type Endpoints interface {
	Spclient() []string
	Dealer() []string
	DealerG2() []string
	Accesspoint() []string
}

// EndpointResolver keeps the historical two-method shape so existing callers
// (notably the dealer package) keep compiling. New code should prefer the
// Option/IOResult-returning methods when available on concrete types.
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
