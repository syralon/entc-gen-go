package server

import (
    "log/slog"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"

    "{{.module}}/internal/conf"
    "{{.module}}/internal/service"
    pb "{{.proto_package}}"
)

func NewHTTPServer(c *conf.Server, services *service.Services) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			//logging.Server(logger),
			//validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	{{ range .services }}pb.Register{{.}}HTTPServer(srv, services.{{.}})
	{{ end }}
	slog.Info("http server started", "addr", c.Http.Addr)

	return srv
}
