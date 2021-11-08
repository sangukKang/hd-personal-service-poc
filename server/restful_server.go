package server

import (
	"context"
	"fmt"
	"net/http"

	"api-server-poc/config"
	"api-server-poc/logger"
	"api-server-poc/proto/generated"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func APIServer(lifecycle fx.Lifecycle, conf *config.Config) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := generated.RegisterS3ManagerServiceHandlerFromEndpoint(ctx, mux, "localhost:"+conf.RPCPort, opts)
	if err != nil {
		logger.Error(err)
	}
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				var ip string
				if conf.Host == "" {
					ip = "localhost"
				} else {
					ip = conf.Host
				}
				host := fmt.Sprintf("%s:%s", ip, conf.Port)
				logger.Info("run http restful service", host)
				go http.ListenAndServe(host, mux)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				cancel()
				return nil
			},
		},
	)
}
