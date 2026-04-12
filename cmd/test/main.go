package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/pyrorhythm/libspot/auth"
	"github.com/pyrorhythm/libspot/auth/session"
	"github.com/pyrorhythm/libspot/auth/store"
)

func main() {
	ctx, cl := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cl()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	redirectPort := 4382
	sess := session.New(
		auth.NewDefaultOAuthConfig(redirectPort),
		store.NewZalandoKeychainer,
	)
	err := sess.Load()
	if err != nil && errors.Is(err, store.ErrItemNotFound) {
		srvctx, cancel := context.WithCancel(ctx)
		codeCh := auth.StartOAuth2Server(srvctx, redirectPort)
		url, pkce := sess.AuthUrl("")
		println(url)
		code := <-codeCh
		cancel()
		err = sess.AuthCode(ctx, code, pkce)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	at, err := sess.AccessToken(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("access token: %s\n", at)
}
