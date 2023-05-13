package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework10/internal/util"
	"log"
	"time"
)

func LoggerInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		formatter := util.NewDateTimeFormatter(time.DateTime)
		startTime, _ := formatter.ToTime(time.Now().UTC())

		resp, err := handler(ctx, req)
		var errStr string
		if err == nil {
			errStr = "good request"
		} else {
			errStr = err.Error()
		}

		latencyTime := time.Since(startTime)
		logger.Printf("[%s] | %s | %s | %s",
			formatter.ToString(time.Now().UTC()),
			info.FullMethod,
			latencyTime.String(),
			errStr,
		)
		return resp, err
	}
}

func RecoveryInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("PANIC ERROR: %v\n", r)
				err = status.Errorf(codes.Internal, "Internal Server Error")
			}
		}()
		return handler(ctx, req)
	}
}
