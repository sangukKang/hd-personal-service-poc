package server

import (
	"context"
	"fmt"
	"net"

	"api-server-poc/config"
	"api-server-poc/logger"
	"api-server-poc/proto/generated"
	"api-server-poc/server/handlers"
	"api-server-poc/server/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func GRPCServer(lifecycle fx.Lifecycle, conf *config.Config) {
	var ip string
	if conf.Host == "" {
		ip = "localhost"
	} else {
		ip = conf.Host
	}

	host := fmt.Sprintf("%s:%s", ip, conf.RPCPort)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.Errorf("failed to listen: %s", err.Error())
		return
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(), //서버 리커버리
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(
					grpc_ctxtags.CodeGenRequestFieldExtractor,
				),
			),
			middleware.LoggerUnaryInterceptor,
			middleware.AuthUnaryInterceptor,
			grpc_prometheus.UnaryServerInterceptor,
		),
	)

	srv := handlers.NewS3ManagerServiceServer()

	generated.RegisterS3ManagerServiceServer(s, srv)

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("run grpc service", host)
				go s.Serve(lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.Stop()
				return nil
			},
		},
	)
}
