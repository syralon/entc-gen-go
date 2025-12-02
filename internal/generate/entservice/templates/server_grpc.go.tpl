package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"{{.module}}/internal/conf"
	"{{.module}}/internal/service"
	pb "{{.proto_package}}"
)

func NewGRPCServer(c *conf.Server, services *service.Services) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			//logging.Server(logger),
			//validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

    {{ range .services }}pb.Register{{.}}Server(srv, services.{{.}})
    {{ end }}
	return srv
}
