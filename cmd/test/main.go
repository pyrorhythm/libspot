package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot/auth/server"
	"github.com/pyrorhythm/libspot/auth/session"
	"github.com/pyrorhythm/libspot/auth/store"
	"github.com/pyrorhythm/libspot/dealer"
	"github.com/pyrorhythm/libspot/pathfinder"
	pfq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	"github.com/pyrorhythm/libspot/pkg/keychain"
	"github.com/pyrorhythm/zlog"
)

func main() {
	ctx, cl := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cl()

	logger := zlog.New(zlog.Options{
		Sinks: []zlog.Sink{{
			Writer: os.Stdout,
			Level:  zlog.LevelTrace,
			Style:  zlog.DefaultStyle,
		}},
	})

	slog.SetDefault(logger)

	redirectPort := 4382
	sess := session.New(
		session.RedirectPort(redirectPort),
		session.GracefulContext(ctx),
		session.Keychainer(store.Zalando))
	err := sess.Load()
	if err != nil && errors.Is(err, keychain.ErrItemNotFound) {
		srvctx, cancel := context.WithCancel(ctx)
		codeCh := server.StartOAuth2Server(srvctx, redirectPort)
		url, pkce := sess.AuthUrl("")
		_ = browser.OpenURL(url)
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

	d, err := startDealer(ctx, sess)
	if err != nil {
		panic(err)
	}
	defer d.Stop()

	pf := pathfinder.New(sess)
	sugg, err := pf.Query(ctx,
		pfq.Suggestions().
			WithQuery("neonate").
			WithCommons(pfq.Commons().WithLimit(1).WithTopResults(1)))
	if err != nil {
		slog.Log(ctx, zlog.LevelPanic, "failed to query suggestions", "error", err)
	}

	bs, _ := json.MarshalIndent(sugg, "", "\t")
	fmt.Println(string(bs))

	top, err := pf.Query(ctx,
		pfq.Top().
			WithQuery("neonate").
			WithCommons(pfq.Commons().WithLimit(1).WithTopResults(1)))
	if err != nil {
		slog.Log(ctx, zlog.LevelPanic, "failed to query suggestions", "error", err)
	}

	bs, _ = json.MarshalIndent(top, "", "\t")
	fmt.Println(string(bs))

	<-ctx.Done()
}

func startDealer(ctx context.Context, sess session.Session) (*dealer.Dealer, error) {
	d, err := dealer.NewFromSession(sess)
	if err != nil {
		return nil, err
	}

	if err = d.Start(ctx); err != nil {
		return nil, err
	}

	return d, nil
}
