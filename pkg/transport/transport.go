package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Options control how every libspot component reaches the network. The zero
// value behaves like net/http's defaults, including honouring the HTTP_PROXY,
// HTTPS_PROXY and NO_PROXY environment variables.
type Options struct {
	// ProxyURL routes all traffic through an HTTP(S) proxy. It overrides the
	// environment proxy settings. Accesspoint connections, which are not HTTP,
	// are tunnelled with CONNECT.
	ProxyURL *url.URL

	// Dialer is the base dialer used for every TCP connection, including the
	// connection to the proxy itself.
	Dialer *net.Dialer

	// TLSClientConfig applies to outgoing TLS connections.
	TLSClientConfig *tls.Config
}

var (
	mu   sync.RWMutex
	opts Options
	rt   *http.Transport
)

// Configure replaces the network options. Call it before creating a session;
// connections already established keep their old settings.
func Configure(o Options) {
	mu.Lock()
	defer mu.Unlock()
	opts = o
	if rt != nil {
		rt.CloseIdleConnections()
		rt = nil
	}
}

// dynamic resolves the underlying transport per request, so round trippers
// handed out before Configure still pick up the new options.
type dynamic struct{}

func (dynamic) RoundTrip(r *http.Request) (*http.Response, error) {
	return httpTransport().RoundTrip(r)
}

// RoundTripper returns the shared round tripper. It always reflects the most
// recent Configure call, whenever it was handed out.
func RoundTripper() http.RoundTripper { return dynamic{} }

func httpTransport() *http.Transport {
	mu.Lock()
	defer mu.Unlock()
	if rt != nil {
		return rt
	}
	proxy := http.ProxyFromEnvironment
	if opts.ProxyURL != nil {
		proxy = http.ProxyURL(opts.ProxyURL)
	}
	rt = &http.Transport{
		Proxy:                 proxy,
		DialContext:           baseDialer(opts).DialContext,
		TLSClientConfig:       opts.TLSClientConfig,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
	return rt
}

// HTTPClient returns a client using the configured transport. A zero timeout
// means no timeout.
func HTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Transport: RoundTripper(), Timeout: timeout}
}

// DialContext opens a raw TCP connection to addr, tunnelling through the
// configured HTTP proxy with CONNECT when one applies. It is used for the
// accesspoint connection, which does not speak HTTP.
func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	mu.RLock()
	o := opts
	mu.RUnlock()

	proxy := o.ProxyURL
	if proxy == nil {
		var err error
		proxy, err = http.ProxyFromEnvironment(&http.Request{URL: &url.URL{Scheme: "https", Host: addr}})
		if err != nil {
			return nil, fmt.Errorf("resolve proxy for %s: %w", addr, err)
		}
	}

	d := baseDialer(o)
	if proxy == nil {
		return d.DialContext(ctx, network, addr)
	}
	return dialConnect(ctx, o, d, proxy, addr)
}

func baseDialer(o Options) *net.Dialer {
	if o.Dialer != nil {
		return o.Dialer
	}
	return &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
}

func dialConnect(
	ctx context.Context,
	o Options,
	d *net.Dialer,
	proxy *url.URL,
	addr string,
) (net.Conn, error) {
	conn, err := d.DialContext(ctx, "tcp", proxyAddr(proxy))
	if err != nil {
		return nil, fmt.Errorf("dial proxy %s: %w", proxy.Host, err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	if proxy.Scheme == "https" {
		cfg := &tls.Config{}
		if o.TLSClientConfig != nil {
			cfg = o.TLSClientConfig.Clone()
		}
		if cfg.ServerName == "" {
			cfg.ServerName = proxy.Hostname()
		}
		tc := tls.Client(conn, cfg)
		if err = tc.HandshakeContext(ctx); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("proxy tls handshake: %w", err)
		}
		conn = tc
	}

	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: make(http.Header),
	}
	if u := proxy.User; u != nil {
		pw, _ := u.Password()
		req.Header.Set("Proxy-Authorization",
			"Basic "+base64.StdEncoding.EncodeToString([]byte(u.Username()+":"+pw)))
	}
	if err = req.Write(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("write CONNECT %s: %w", addr, err)
	}

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("read CONNECT response for %s: %w", addr, err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_ = conn.Close()
		return nil, fmt.Errorf("proxy refused CONNECT %s: %s", addr, resp.Status)
	}
	// The tunnel must be byte-clean: anything the proxy sent early would be
	// stranded in br and lost by the caller.
	if br.Buffered() > 0 {
		_ = conn.Close()
		return nil, fmt.Errorf("proxy sent %d unexpected bytes after CONNECT %s", br.Buffered(), addr)
	}

	_ = conn.SetDeadline(time.Time{})
	return conn, nil
}

func proxyAddr(p *url.URL) string {
	if p.Port() != "" {
		return p.Host
	}
	if p.Scheme == "https" {
		return net.JoinHostPort(p.Hostname(), "443")
	}
	return net.JoinHostPort(p.Hostname(), "80")
}
