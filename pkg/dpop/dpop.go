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
)


type (
	getter func() (string, bool)
	setter func(string)
)

type Provider struct {
	GetAccessToken getter
	GetNonce       getter
	SetNonce       setter
}

// Transport is an http.RoundTripper that attaches a DPoP proof to every request
// and handles server-issued nonce challenges (RFC 9449 §4.3).
type Transport struct {
	Base http.RoundTripper
	Prov *Provider
	Key  *ecdsa.PrivateKey
	mu   sync.RWMutex
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	firstReq := req.Clone(req.Context())

	resp, err := t.roundTripWithProof(firstReq)
	if err != nil {
		return nil, err
	}

	if !t.respHasUseDPoPNonce(resp) {
		return resp, nil
	}

	nonce, ok := parseDPoPNonceHeader(resp)
	if !ok {
		return nil, fmt.Errorf("dpop: server requested nonce but DPoP-Nonce header missing")
	}
	t.Prov.SetNonce(nonce)

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	slog.Debug("dpop: retrying with server nonce")

	retryReq := req.Clone(req.Context())
	if req.Body != nil {
		if req.GetBody == nil {
			return nil, errors.New("dpop: request body cannot be re-read for retry")
		}
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("dpop: get body for retry: %w", err)
		}
		retryReq.Body = body
	}

	retryResp, err := t.roundTripWithProof(retryReq)
	if err != nil {
		return nil, err
	}

	if nonce, ok := parseDPoPNonceHeader(retryResp); ok {
		t.Prov.SetNonce(nonce)
	}

	return retryResp, nil
}

func (t *Transport) respHasUseDPoPNonce(resp *http.Response) bool {
	return parseUseDPoPNonceBody(resp) || parseUseDPoPNonceWWW(resp)
}

func parseUseDPoPNonceWWW(resp *http.Response) bool {
	auth := resp.Header.Get("WWW-Authenticate")
	return strings.Contains(auth, "error=\"use_dpop_nonce\"") ||
		strings.Contains(auth, "error=use_dpop_nonce")
}

func parseUseDPoPNonceBody(resp *http.Response) bool {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return false
	}

	if nod, err := sonic.Get(bodyBytes, "error"); err != nil {
		return false
	} else if jserr, err := nod.String(); err != nil {
		return false
	} else {
		return jserr == "use_dpop_nonce"
	}
}

func parseDPoPNonceHeader(resp *http.Response) (string, bool) {
	nonce := resp.Header.Get("DPoP-Nonce")
	return nonce, nonce != ""
}

func (t *Transport) roundTripWithProof(req *http.Request) (*http.Response, error) {
	proof, err := t.createProof(req.Context(), req)
	if err != nil {
		return nil, err
	}
	req.Header.Set("DPoP", proof)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

func (t *Transport) createProof(ctx context.Context, req *http.Request) (string, error) {
	pub := jose.JSONWebKey{
		Key:       &t.Key.PublicKey,
		Algorithm: string(jose.ES256),
		Use:       "sig",
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: t.Key},
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

	if nonce, ok := t.Prov.GetNonce(); ok {
		claims["nonce"] = nonce
	}

	if accessToken, ok := t.Prov.GetAccessToken(); ok && accessToken != "" {
		hash := sha256.Sum256([]byte(accessToken))
		claims["ath"] = base64.RawURLEncoding.EncodeToString(hash[:])
	} else {
		slog.Debug("dpop: no access token, omitting ath claim")
	}

	slog.Debug("dpop: proof claims", "claims", claims)

	return jwt.Signed(signer).Claims(claims).Serialize()
}
