package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pyrorhythm/zlog"
	"pyrorhythm.dev/libspot/auth/session"
	"pyrorhythm.dev/libspot/auth/store"
	"pyrorhythm.dev/libspot/connect"
	"pyrorhythm.dev/libspot/dealer"
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

	sess := session.New(
		session.RedirectPort(9292),
		session.GracefulContext(ctx),
		session.Keychainer(store.Zalando),
		session.Interactive(func(url string) {
			fmt.Println("open this URL to log in:", url)
		}),
	)
	if err := sess.Load(); err != nil {
		panic(err)
	}

	d, err := startDealer(ctx, sess)
	if err != nil {
		panic(err)
	}
	defer d.Stop()

	time.Sleep(5 * time.Second)

	conn, err := connect.NewFromSession(sess, connect.ConnectOptions{})
	if err != nil {
		panic(err)
	}

	conn.Bind(d)

	pb, err := conn.Playback(ctx)
	if err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(pb, "", "\t")

	fmt.Println(string(bs))

	conn.Play(ctx, "")

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
