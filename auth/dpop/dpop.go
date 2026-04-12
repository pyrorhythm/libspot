package dpop

import (
	"bytes"
	"context"
	"crypto/ecdsa"
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
	"github.com/pyrorhythm/libspot"
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

type provider interface {
	libspot.TokenProvider
	GetNonce() (string, bool)
	SetNonce(string)
}

type Transport struct {
	base        http.RoundTripper
	clientToken string
	prov        provider

	key *ecdsa.PrivateKey
	mu  sync.RWMutex
}

func NewClientWithKey(
	clientToken string,
	provider provider,
	key *ecdsa.PrivateKey,
) *http.Client {
	return &http.Client{
		Transport: &Transport{
			base:        &loggingTransport{},
			clientToken: clientToken,
			prov:        provider,
			key:         key,
		},
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	firstReq := req.Clone(req.Context())

	resp, err := t.roundTripWithNonce(firstReq)
	if err != nil {
		return nil, err
	}

	hasUseDPoPNonce := t.respHasUseDPoPNonce(resp)

	if !hasUseDPoPNonce {
		return resp, nil
	}

	if nonce, ok := parseDPoPNonceHeader(resp); !ok {
		return nil, fmt.Errorf("server requests dpop with nonce but could not find it in header")
	} else {
		t.prov.SetNonce(nonce)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	slog.Debug("dpop: retrying with server nonce")

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

	retryResp, err := t.roundTripWithNonce(retryReq)
	if err != nil {
		return nil, err
	}

	if nonce, ok := parseDPoPNonceHeader(retryResp); ok {
		t.prov.SetNonce(nonce)
	}

	return retryResp, nil
}

func (t *Transport) respHasUseDPoPNonce(resp *http.Response) (retry bool) {
	return parseUseDPoPNonceBody(resp) || parseUseDPoPNonceWWW(resp)
}

func parseUseDPoPNonceWWW(resp *http.Response) (found bool) {
	auth := resp.Header.Get("WWW-Authenticate")
	if auth == "" {
		return false
	}

	return strings.Contains(auth, "error=\"use_dpop_nonce\"") ||
		strings.Contains(auth, "error=use_dpop_nonce")
}

func parseUseDPoPNonceBody(resp *http.Response) (found bool) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return false
	}

	var errResp struct {
		Error           string `json:"error"`
		ErrorDescrption string `json:"error_descrption"`
	}
	if err := sonic.Unmarshal(bodyBytes, &errResp); err != nil {
		return false
	}
	return errResp.Error == "use_dpop_nonce"
}

func parseDPoPNonceHeader(resp *http.Response) (string, bool) {
	nonce := resp.Header.Get("Dpop-Nonce")

	return nonce, nonce != ""
}

func (t *Transport) roundTripWithNonce(req *http.Request) (*http.Response, error) {
	proof, err := t.createProof(req.Context(), req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("DPoP", proof)
	if t.clientToken != "" {
		req.Header.Set("Client-Token", t.clientToken)
	}

	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func (t *Transport) createProof(
	ctx context.Context,
	req *http.Request,
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

	if nonce, ok := t.prov.GetNonce(); ok {
		claims["nonce"] = nonce
	}

	if accessToken, err := t.prov.AccessToken(ctx); err == nil && accessToken != "" {
		hash := sha256.Sum256([]byte(accessToken))
		claims["ath"] = base64.RawURLEncoding.EncodeToString(hash[:])
	} else if err != nil {
		slog.Debug("dpop: no access token available, omitting ath claim", "error", err)
	}

	slog.Debug("dpop: proof claims", "claims", claims)

	return jwt.Signed(signer).Claims(claims).Serialize()
}
