package session

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/pyrorhythm/fn"
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth"
	"github.com/pyrorhythm/libspot/auth/store"
	"github.com/pyrorhythm/libspot/pkg/keychain"
	"github.com/pyrorhythm/libspot/resolver"
	"golang.org/x/oauth2"
)

const sessionKey = "authorizationData"

type storedCredentials struct {
	*oauth2.Token `json:"token,omitempty"`

	ClientToken string `json:"clientToken,omitempty"`
	UserName    string `json:"userName,omitempty"`
	DeviceId    string `json:"deviceId,omitempty"`

	Proof struct {
		Privkey   string `json:"privKey,omitempty"`
		LastNonce string `json:"lastNonce,omitempty"`
	} `json:"proof,omitempty"`
}

type session struct {
	mu      sync.RWMutex
	nonceMu sync.RWMutex

	kcer keychain.Keychainer[storedCredentials]

	creds *storedCredentials
	conf  *oauth2.Config

	dpopClient *http.Client

	resolver    libspot.EndpointResolver
	gracefulCtx context.Context
}

func (s *session) Resolver() (libspot.EndpointResolver, error) {
	if s.resolver != nil {
		return s.resolver, nil
	}

	ctok, err := s.safeClientToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid client-token: %w", err)
	}

	s.resolver = resolver.New(ctok)
	return s.resolver, nil
}

func (s *session) ClientToken() (string, error) {
	return s.safeClientToken()
}

func (s *session) GetNonce() (string, bool) {
	s.nonceMu.RLock()
	defer s.nonceMu.RUnlock()
	return s.creds.Proof.LastNonce, s.creds.Proof.LastNonce != ""
}

func (s *session) SetNonce(nonce string) {
	s.nonceMu.Lock()
	defer s.nonceMu.Unlock()
	s.creds.Proof.LastNonce = nonce
}

type Option func(*session)

func WithRedirectPort(c int) Option {
	return func(s *session) {
		s.conf = auth.NewDefaultOAuthConfig(c)
	}
}

func WithKeychainerFunc(kcfn func(key string) keychain.Keychainer[storedCredentials]) Option {
	return func(s *session) {
		s.kcer = kcfn(sessionKey)
	}
}

func WithGracefulContext(ctx context.Context) Option {
	return func(s *session) {
		s.gracefulCtx = ctx
	}
}

func applyDefaults(s *session) {
	s.conf = auth.NewDefaultOAuthConfig(9292)
	s.kcer = store.Zalando[storedCredentials](sessionKey)
	s.gracefulCtx = context.Background()
}

func New(
	opts ...Option,
) Session {
	s := &session{}
	
	for _, opt := range opts {
		opt(s)
	}
	
	applyDefaults(s)
	
	return s
}

func generateDeviceId() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func validateDeviceId(deviceId string) bool {
	if len(deviceId) != 40 {
		return false
	}
	_, err := hex.DecodeString(deviceId)
	return err == nil
}

func (s *session) getOrCreateDPoPKey() (*ecdsa.PrivateKey, error) {
	if s.creds.Proof.Privkey != "" {
		block, _ := pem.Decode([]byte(s.creds.Proof.Privkey))
		if block == nil {
			return nil, errors.New("invalid DPoP key PEM")
		}
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse DPoP key: %w", err)
		}
		return key, nil
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate DPoP key: %w", err)
	}

	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshal DPoP key: %w", err)
	}
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	})
	s.creds.Proof.Privkey = string(pemBlock)

	return key, nil
}

// safeClientToken returns the Client‑Token, fetching it if necessary.
func (s *session) safeClientToken() (string, error) {
	clientToken := s.creds.ClientToken
	deviceId := s.creds.DeviceId

	if clientToken != "" {
		return clientToken, nil
	}

	if !validateDeviceId(deviceId) {
		deviceId = generateDeviceId()
		s.creds.DeviceId = deviceId
	}

	token, err := libspot.RetrieveClientToken(http.DefaultClient, deviceId)
	if err != nil {
		return "", fmt.Errorf("fetch client token: %w", err)
	}

	s.creds.ClientToken = token
	return token, nil
}

func (s *session) injectDPoPClient(ctx context.Context) (context.Context, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.dpopClient == nil {
		clientToken, err := s.safeClientToken()
		if err != nil {
			return nil, err
		}
		key, err := s.getOrCreateDPoPKey()
		if err != nil {
			return nil, err
		}
		s.dpopClient = newSessionDpopClient(clientToken, s, key)
	}
	return context.WithValue(ctx, oauth2.HTTPClient, s.dpopClient), nil
}

func (s *session) AuthUrl(state string) (url, pkce string) {
	pkce = oauth2.GenerateVerifier()
	return s.conf.AuthCodeURL(state, oauth2.S256ChallengeOption(pkce)), pkce
}

func (s *session) AuthCode(ctx context.Context, code, pkce string) error {
	ctx, err := s.injectDPoPClient(ctx)
	if err != nil {
		return fmt.Errorf("prepare DPoP client: %w", err)
	}

	slog.Debug("processing code", "code", code, "code_verifier", pkce)

	tok, err := s.conf.Exchange(ctx, code, oauth2.VerifierOption(pkce))
	if err != nil {
		return fmt.Errorf("code exchange failed: %w", err)
	}

	return s.SaveToken(tok)
}

func (s *session) GetToken() (string, bool) {
	s.mu.RLock()
	valid := s.creds.Valid()
	accessToken := s.creds.AccessToken
	s.mu.RUnlock()

	return fn.If(valid, accessToken, ""), valid
}

func (s *session) GetOrRefreshToken(ctx context.Context) (string, error) {
	s.mu.RLock()
	valid := s.creds.Valid()
	accessToken := s.creds.AccessToken
	refreshToken := s.creds.RefreshToken
	s.mu.RUnlock()

	if valid {
		return accessToken, nil
	}
	if refreshToken == "" {
		return "", errors.New("no refresh token available")
	}

	// refresh required
	ctx, err := s.injectDPoPClient(ctx)
	if err != nil {
		return "", err
	}

	ts := s.conf.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})
	newTok, err := ts.Token()
	if err != nil {
		return "", fmt.Errorf("refresh token: %w", err)
	}

	if err := s.SaveToken(newTok); err != nil {
		return "", err
	}
	return newTok.AccessToken, nil
}

// SaveToken updates the stored token and persists it to the keychain.
func (s *session) SaveToken(tok *oauth2.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.creds.Token = tok

	if username, ok := tok.Extra("username").(string); ok {
		s.creds.UserName = username
	}

	return s.kcer.Save(s.creds)
}

func (s *session) RefreshToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.creds.RefreshToken
}

func (s *session) User() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.creds.UserName
}

func (s *session) DeviceId() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.creds.DeviceId
}

func (s *session) Load() error {
	creds, err := s.kcer.Load(false)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.creds = creds
	return nil
}

func (s *session) Valid() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.creds.Valid()
}

func (s *session) Clear(clearKeychain bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if clearKeychain {
		if err := s.kcer.Delete(); err != nil {
			return err
		}
	}

	s.creds = &storedCredentials{Token: &oauth2.Token{}}
	s.dpopClient = nil

	return nil
}
