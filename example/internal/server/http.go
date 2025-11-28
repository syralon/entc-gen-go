package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"

    "github.com/syralon/entc-gen-go/example/internal/conf"
    "github.com/syralon/entc-gen-go/example/internal/service"
    pb "github.com/syralon/entc-gen-go/example/proto/example"
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

	pb.RegisterGroupServiceHTTPServer(srv, services.GroupService)
	pb.RegisterUserServiceHTTPServer(srv, services.UserService)
	
	return srv
}
