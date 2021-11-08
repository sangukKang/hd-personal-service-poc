package middleware

import (
	"context"
	"os"

	"api-server-poc/logger"
	"api-server-poc/utils"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggerUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UNKNOW"
	}
	requestID := utils.RandToken(5)
	entry := logger.WithFields(logrus.Fields{
		logger.LFhostname:  hostname,
		logger.LFRequestID: requestID,
	})
	entry.Info(info.FullMethod)

	return handler(ctx, req)
}
