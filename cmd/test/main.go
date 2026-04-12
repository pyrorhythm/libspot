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
	"github.com/pyrorhythm/libspot/auth/store"
	"github.com/pyrorhythm/libspot/dealer"
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

	at, err := sess.GetOrRefreshToken(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("access token: %s\n", at)

	testDealer(sess)
}

func testDealer(sess session.Session) {
	d, err := dealer.NewSession(sess, dealer.ExponentialJitter2Delay(time.Second))
	if err != nil {
		debug.Stack()
		panic(err)
	}

	if err = d.Start(); err != nil {
		debug.Stack()
		panic(err)
	}

	time.Sleep(1 * time.Hour) // lol
}
