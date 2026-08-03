package transport

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"
)

// fakeProxy accepts one CONNECT request and splices the tunnel to target.
// It reports the CONNECT target and the Proxy-Authorization header it saw.
func fakeProxy(t *testing.T, target string) (addr string, seen chan [2]string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	seen = make(chan [2]string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		br := bufio.NewReader(conn)
		req, err := http.ReadRequest(br)
		if err != nil {
			_ = conn.Close()
			return
		}
		seen <- [2]string{req.Host, req.Header.Get("Proxy-Authorization")}
		if req.Method != http.MethodConnect {
			// A forwarded (non-tunnelled) request; answering it is enough to
			// prove it reached the proxy.
			_, _ = io.WriteString(conn, "HTTP/1.1 204 No Content\r\n\r\n")
			_ = conn.Close()
			return
		}
		up, err := net.Dial("tcp", target)
		if err != nil {
			_, _ = io.WriteString(conn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			_ = conn.Close()
			return
		}
		_, _ = io.WriteString(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		go func() { _, _ = io.Copy(up, conn) }()
		go func() { _, _ = io.Copy(conn, up); _ = conn.Close() }()
	}()
	return ln.Addr().String(), seen
}

// echoServer stands in for the accesspoint: raw TCP, no HTTP.
func echoServer(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func() { _, _ = io.Copy(conn, conn); _ = conn.Close() }()
		}
	}()
	return ln.Addr().String()
}

// The accesspoint speaks a raw binary protocol, so a configured proxy must be
// reached with CONNECT and then get out of the way entirely.
func TestDialContextTunnelsRawTCPThroughProxy(t *testing.T) {
	origin := echoServer(t)
	proxyAddr, seen := fakeProxy(t, origin)

	proxyURL, err := url.Parse("http://user:pass@" + proxyAddr)
	if err != nil {
		t.Fatal(err)
	}
	Configure(Options{ProxyURL: proxyURL})
	t.Cleanup(func() { Configure(Options{}) })

	conn, err := DialContext(context.Background(), "tcp", origin)
	if err != nil {
		t.Fatalf("dial through proxy: %v", err)
	}
	defer func() { _ = conn.Close() }()

	got := <-seen
	if got[0] != origin {
		t.Errorf("CONNECT target = %q, want %q", got[0], origin)
	}
	// dXNlcjpwYXNz is base64("user:pass"); without it an authenticating proxy
	// would reject every accesspoint connection.
	if want := "Basic dXNlcjpwYXNz"; got[1] != want {
		t.Errorf("Proxy-Authorization = %q, want %q", got[1], want)
	}

	if _, err = conn.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 4)
	if _, err = io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "ping" {
		t.Errorf("tunnelled payload = %q, want %q", buf, "ping")
	}
}

// Round trippers are handed out during package init, long before main parses
// PROXY_URL, so they must resolve the proxy per request rather than at capture.
func TestRoundTripperPicksUpLaterConfigure(t *testing.T) {
	rt := RoundTripper()

	proxyAddr, seen := fakeProxy(t, "")

	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		t.Fatal(err)
	}
	Configure(Options{ProxyURL: proxyURL})
	t.Cleanup(func() { Configure(Options{}) })

	req, err := http.NewRequest(http.MethodGet, "http://spclient.wg.spotify.com/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case got := <-seen:
		if got[0] != "spclient.wg.spotify.com" {
			t.Errorf("proxied request host = %q", got[0])
		}
	default:
		t.Fatal("request bypassed the proxy configured after RoundTripper() was called")
	}
}
