package connect

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"pyrorhythm.dev/libspot"
	"pyrorhythm.dev/libspot/auth/client"
	"pyrorhythm.dev/libspot/auth/session"
	"pyrorhythm.dev/libspot/dealer"
	"resty.dev/v3"
)

const (
	connectionTTL      = 10 * time.Minute
	connectDeviceName  = "libspot"
	connectDeviceModel = "web_player"
)

var clientVersion = libspot.DottedApplicationVer

var (
	errMissingDevice         = errors.New("connect: missing device id")
	errFailedToFetchEndpoint = errors.New("connect: failed to fetch endpoint")
	errNoConnectionID        = errors.New(
		"connect: dealer has no connection id (call Bind after Start?)",
	)
)

// ConnectOptions tweaks connect behaviour.
type ConnectOptions struct {
	// Device selects the playback target by id or name when Play is called
	// without an active device.
	Device string
}

type Connect struct {
	prov libspot.TokenProvider
	rslv libspot.EndpointResolver
	http *resty.Client
	opts ConnectOptions

	dealerMu sync.RWMutex
	dealer   *dealer.Dealer

	hostMu sync.RWMutex
	host   string

	mu              sync.Mutex
	connectDeviceID string
	lastCID         string
	registeredAt    time.Time

	routeMu     sync.RWMutex
	activeID    string
	originID    string
	routeStored bool
}

func New(prov libspot.TokenProvider, rslv libspot.EndpointResolver, opts ConnectOptions) *Connect {
	rc := client.
		NewAuthorizedClient(prov, client.CanRefreshAccessToken(true)).
		Client()
	return &Connect{prov: prov, rslv: rslv, http: rc, opts: opts}
}

func NewFromSession(
	sess session.Session,
	opts ConnectOptions,
) (*Connect, error) {
	rslv, err := sess.Resolver()
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver from session: %w", err)
	}
	return New(sess, rslv, opts), nil
}

// Bind attaches a started dealer as the connection-id source for connect-state.
func (c *Connect) Bind(d *dealer.Dealer) {
	c.dealerMu.Lock()
	c.dealer = d
	c.dealerMu.Unlock()
}

func (c *Connect) connectionID() (string, error) {
	c.dealerMu.RLock()
	d := c.dealer
	c.dealerMu.RUnlock()
	if d == nil {
		return "", errors.New("connect: dealer not bound (call Bind)")
	}
	if id := d.ConnectionID(); id != "" {
		return id, nil
	}
	return "", errNoConnectionID
}

func (c *Connect) spclientHost() (string, error) {
	c.hostMu.RLock()
	if c.host != "" {
		h := c.host
		c.hostMu.RUnlock()
		return h, nil
	}
	c.hostMu.RUnlock()

	eps, ok := c.rslv.Endpoints()
	if !ok || len(eps.Spclient()) == 0 {
		fetched, err := c.rslv.Fetch(libspot.ServiceKindSpclient)
		if err == nil {
			eps = fetched
		}
	}
	if eps == nil || len(eps.Spclient()) == 0 {
		return "", errFailedToFetchEndpoint
	}
	host := stripPort(eps.Spclient()[0])

	c.hostMu.Lock()
	c.host = host
	c.hostMu.Unlock()
	return host, nil
}

func (c *Connect) connectStateBase() (string, error) {
	host, err := c.spclientHost()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s/connect-state/v1", host), nil
}

func (c *Connect) trackPlaybackBase() (string, error) {
	host, err := c.spclientHost()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s/track-playback/v1", host), nil
}

func stripPort(host string) string {
	if i := strings.LastIndex(host, ":"); i > 0 {
		return host[:i]
	}
	return host
}

func connectHeaders() map[string]string {
	return map[string]string{
		"content-type":        "application/json",
		"app-platform":        libspot.AppPlatform().String(),
		"spotify-app-version": clientVersion,
		"accept":              "application/json",
	}
}
