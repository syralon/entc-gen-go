package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/syralon/entc-gen-go/example/ent"
	"github.com/syralon/entc-gen-go/example/internal/service"
	pb "github.com/syralon/entc-gen-go/example/proto/example"
	"net/http"
)

func main() {
	ctx := context.Background()
	client := ent.NewClient()

	mux := runtime.NewServeMux()

	_ = pb.RegisterGroupServiceHandlerServer(ctx, mux, service.NewGroupService(client))
	_ = pb.RegisterUserServiceHandlerServer(ctx, mux, service.NewUserService(client))

	if err := http.ListenAndServe("0.0.0.0:0", mux); err != nil {
		fmt.Println(err)
	}
}
