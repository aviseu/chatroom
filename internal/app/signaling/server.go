package signaling

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Addr            string        `default:":8080"`
	ShutdownTimeout time.Duration `split_words:"true" default:"5s"`
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
