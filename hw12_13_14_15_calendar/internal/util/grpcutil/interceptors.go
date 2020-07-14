package grpcutil

import (
	"context"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := logging.FromContext(ctx)
		start := time.Now()

		remoteAddress := "UNKNOWN"
		p, ok := peer.FromContext(ctx)
		if ok {
			remoteAddress = p.Addr.String()
		}

		resp, err := handler(ctx, req)

		duration := time.Since(start).String()
		logLine := remoteAddress + " " + "[" + start.Format(time.RFC3339) + "]" + " " + info.FullMethod + " " + status.Code(err).String() + " " + duration
		logger.Info(logLine)

		return resp, err
	}
}
