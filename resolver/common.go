package resolver

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pyrorhythm/fn"
	"github.com/pyrorhythm/libspot"
	"github.com/valyala/fasthttp"
)

const resolveEndpoint = "https://apresolve.spotify.com/"

type endpoints struct {
	SpclientEps    []string `json:"spclient"`
	DealerEps      []string `json:"dealer"`
	DealerG2Eps    []string `json:"dealer-g2"`
	AccesspointEps []string `json:"accesspoint"`
}

func (e *endpoints) merge(oth endpoints) {
	if len(oth.SpclientEps) != 0 {
		e.SpclientEps = oth.SpclientEps
	}

	if len(oth.DealerEps) != 0 {
		e.DealerEps = oth.DealerEps
	}

	if len(oth.DealerG2Eps) != 0 {
		e.DealerG2Eps = oth.DealerG2Eps
	}

	if len(oth.AccesspointEps) != 0 {
		e.AccesspointEps = oth.AccesspointEps
	}
}

func (e endpoints) Spclient() []string {
	return e.SpclientEps
}

func (e endpoints) Dealer() []string {
	return e.DealerEps
}

func (e endpoints) DealerG2() []string {
	return e.DealerG2Eps
}

func (e endpoints) Accesspoint() []string {
	return e.AccesspointEps
}

func tostring(kind libspot.ServiceKind) string {
	return string(kind)
}

type fetcher struct {
	endpoints   *endpoints
	clientToken string
}

func (t *fetcher) Endpoints() (libspot.Endpoints, bool) {
	return t.endpoints, t.endpoints != nil
}

func (t *fetcher) Fetch(kinds ...libspot.ServiceKind) (libspot.Endpoints, error) {
	v := url.Values{
		"type": fn.Map(kinds, tostring),
	}

	if slices.Contains(kinds, libspot.ServiceKindAccesspoint) {
		v["time"] = []string{strconv.FormatInt(time.Now().Unix(), 10)}
	}

	var (
		ep         endpoints
		err        error
		retries    int
		retriesCap = 5
	)

	u := fasthttp.AcquireURI()
	_ = u.Parse(nil, []byte(resolveEndpoint))
	u.SetQueryString(v.Encode())

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetURI(u)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Client-Token", t.clientToken)
	req.Header.SetUserAgent("Spotify/128600502 (43; 0; 2)")

retryPoint:
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = fasthttp.Do(req, resp)
	if err != nil {
		// Error is not a bad status but a transport error
		return nil, fmt.Errorf(
			"internal transport failure: %w",
			err,
		)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		if retries < retriesCap {
			fasthttp.ReleaseResponse(resp)
			retries++
			goto retryPoint
		}

		return nil, fmt.Errorf("max retries reached: %d", retries)
	}

	if err := sonic.Unmarshal(resp.Body(), &ep); err != nil {
		return nil, fmt.Errorf("failed to unmarshal endpoints: %w", err)
	}

	t.endpoints.merge(ep)

	return t.endpoints, nil
}

func New(clientToken string) libspot.EndpointResolver {
	return &fetcher{clientToken: clientToken, endpoints: &endpoints{}}
}
