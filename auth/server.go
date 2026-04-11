package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func StartOAuth2Server(ctx context.Context, port int) <-chan string {
	codeCh := make(chan string, 1)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      http.HandlerFunc(handleCallbackHTTP(codeCh)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	go server.ListenAndServe()

	return codeCh
}

func handleCallbackHTTP(codeCh chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		slog.Debug("oauth server: got code", "code", code)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(
			[]byte("<html><body>Login successful! You can close this tab.</body></html>"),
		)
		codeCh <- code
	}
}
