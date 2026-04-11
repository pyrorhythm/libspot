package dpop

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type loggingTransport struct{}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	id := uuid.Must(uuid.NewV7())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	slog.Debug(
		"request",
		"id", id.String(),
		//"headers", r.Header,
		"dest", r.URL.String(),
		"body", string(body),
	)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	slog.Debug(
		"response",
		"id", id.String(),
		//"headers", resp.Header,
		"code", resp.StatusCode,
		"body", string(respBody),
	)

	return resp, err
}

type Transport struct {
	Base           http.RoundTripper
	ClientToken    string
	GetAccessToken func(context.Context) (string, error)

	key   *ecdsa.PrivateKey
	nonce string
	mu    sync.RWMutex
}

func NewClient(clientToken string, getToken func(context.Context) (string, error)) *http.Client {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return NewClientWithKey(clientToken, getToken, key)
}

func NewClientWithKey(
	clientToken string,
	getToken func(context.Context) (string, error),
	key *ecdsa.PrivateKey,
) *http.Client {
	return &http.Client{
		Transport: &Transport{
			Base:           &loggingTransport{},
			ClientToken:    clientToken,
			GetAccessToken: getToken,
			key:            key,
		},
	}
}

func (t *Transport) SetNonce(nonce string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.nonce = nonce
}

func (t *Transport) Nonce() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.nonce
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	firstReq := req.Clone(req.Context())

	resp, err := t.roundTripWithNonce(firstReq, t.Nonce())
	if err != nil {
		return nil, err
	}

	serverNonce, needRetry := t.extractNonceFromResponse(resp)

	if !needRetry {
		return resp, nil
	}

	t.SetNonce(serverNonce)

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	slog.Debug("dpop: retrying with server nonce", "nonce", serverNonce)

	retryReq := req.Clone(req.Context())

	if req.Body != nil {
		if req.GetBody == nil {
			return nil, errors.New("dpop: request body cannot be re‑read for retry")
		}
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("dpop: get body for retry: %w", err)
		}
		retryReq.Body = body
	}

	retryResp, err := t.roundTripWithNonce(retryReq, serverNonce)
	if err != nil {
		return nil, err
	}

	if nonce, ok := t.extractNonceFromResponse(retryResp); ok {
		t.SetNonce(nonce)
	}

	return retryResp, nil
}

func (t *Transport) extractNonceFromResponse(resp *http.Response) (nonce string, retry bool) {
	if n, ok := parseUseDPoPNonceWWW(resp); ok {
		return n, true
	}

	if n, ok := parseDPoPNonceHeader(resp); ok {
		return n, true
	}

	if n, ok := parseUseDPoPNonceBody(resp); ok {
		return n, true
	}

	return "", false
}

func parseUseDPoPNonceWWW(resp *http.Response) (nonce string, found bool) {
	auth := resp.Header.Get("WWW-Authenticate")
	if auth == "" {
		return "", false
	}
	if !strings.Contains(auth, "error=\"use_dpop_nonce\"") &&
		!strings.Contains(auth, "error=use_dpop_nonce") {
		return "", false
	}

	parts := strings.Split(auth, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "nonce=") {
			val := strings.TrimPrefix(part, "nonce=")
			return strings.Trim(val, "\""), true
		}
	}
	return "", false
}

func parseUseDPoPNonceBody(resp *http.Response) (nonce string, found bool) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return "", false
	}

	var errResp struct {
		Error string `json:"error"`
		Nonce string `json:"nonce"`
	}
	if err := sonic.Unmarshal(bodyBytes, &errResp); err != nil {
		return "", false
	}
	if errResp.Error == "use_dpop_nonce" && errResp.Nonce != "" {
		return errResp.Nonce, true
	}
	return "", false
}

func parseDPoPNonceHeader(resp *http.Response) (string, bool) {
	nonce := resp.Header.Get("Dpop-Nonce")

	return nonce, nonce != ""
}

func (t *Transport) roundTripWithNonce(req *http.Request, nonce string) (*http.Response, error) {
	proof, err := t.createProof(req.Context(), req, nonce)
	if err != nil {
		return nil, err
	}

	req.Header.Set("DPoP", proof)
	if t.ClientToken != "" {
		req.Header.Set("Client-Token", t.ClientToken)
	}

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func (t *Transport) createProof(
	ctx context.Context,
	req *http.Request,
	nonce string,
) (string, error) {
	pub := jose.JSONWebKey{
		Key:       &t.key.PublicKey,
		Algorithm: string(jose.ES256),
		Use:       "sig",
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: t.key},
		(&jose.SignerOptions{}).WithType("dpop+jwt").WithHeader("jwk", pub),
	)
	if err != nil {
		return "", fmt.Errorf("dpop: create signer: %w", err)
	}

	htu := &url.URL{
		Scheme:   req.URL.Scheme,
		Host:     req.URL.Host,
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}
	if htu.Scheme == "" {
		htu.Scheme = "https"
	}

	claims := map[string]any{
		"jti": uuid.New().String(),
		"iat": time.Now().Unix(),
		"htm": req.Method,
		"htu": htu.String(),
	}

	if nonce != "" {
		claims["nonce"] = nonce
	}

	if accessToken, err := t.GetAccessToken(ctx); err == nil && accessToken != "" {
		hash := sha256.Sum256([]byte(accessToken))
		claims["ath"] = base64.RawURLEncoding.EncodeToString(hash[:])
	} else if err != nil {
		slog.Debug("dpop: no access token available, omitting ath claim", "error", err)
	}

	slog.Debug("dpop: proof claims", "claims", claims)

	return jwt.Signed(signer).Claims(claims).Serialize()
}
