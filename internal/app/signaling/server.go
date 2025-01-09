package signaling

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

type Config struct {
	Addr string `split_words:"true" default:":8080"`
}

func SetupServer(ctx context.Context, cfg Config, h http.Handler) http.Server {
	return http.Server{
		Addr:    cfg.Addr,
		Handler: h,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
}

func SetupHandler(log *slog.Logger) http.Handler {
	r := http.NewServeMux()

	h := NewHandler(log)
	r.HandleFunc("/ws", h.Handle)

	return r
}
