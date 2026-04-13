package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/pyrorhythm/libspot/auth"
	"github.com/pyrorhythm/libspot/auth/session"
	"github.com/pyrorhythm/libspot/dealer"
	"github.com/pyrorhythm/libspot/pkg/keychain"
)

func main() {
	ctx, cl := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cl()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	redirectPort := 4382
	sess := session.New(
		session.WithRedirectPort(redirectPort),
		session.WithGracefulContext(ctx),
	)
	err := sess.Load()
	if err != nil && errors.Is(err, keychain.ErrItemNotFound) {
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

	at, err := sess.AccessToken(ctx, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("access token: %s\n", at)

	testDealer(ctx, sess)
}

func testDealer(ctx context.Context, sess session.Session) {
	d, err := dealer.NewFromSession(sess)
	if err != nil {
		debug.Stack()
		panic(err)
	}

	if err = d.Start(ctx); err != nil {
		debug.Stack()
		panic(err)
	}

	timer := time.NewTimer(time.Hour)

	select {
	case <-timer.C:
	case <-ctx.Done():
	}

	_ = d.Stop()
}
