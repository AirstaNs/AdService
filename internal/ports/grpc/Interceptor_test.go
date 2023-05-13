package grpc

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"testing"
)

func TestLoggerInterceptor(t *testing.T) {
	logOutput := bytes.Buffer{}
	logger := log.New(&logOutput, "", 0)
	interceptor := LoggerInterceptor(logger)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/service.Service/Method",
	}

	handlerFunc := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}

	_, _ = interceptor(context.Background(), nil, info, handlerFunc)

	assert.NotEqual(t, logOutput.String(), 0, "log output should be empty")
}

func TestLoggerInterceptorErr(t *testing.T) {
	logOutput := bytes.Buffer{}
	logger := log.New(&logOutput, "", 0)
	interceptor := LoggerInterceptor(logger)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/service.Service/Method",
	}

	handlerFunc := func(ctx context.Context, req any) (any, error) {
		return nil, assert.AnError
	}

	_, _ = interceptor(context.Background(), nil, info, handlerFunc)

	assert.NotEqual(t, logOutput.String(), 0, "log output should be empty")
}

func TestRecoveryInterceptor(t *testing.T) {
	logOutput := bytes.Buffer{}
	logger := log.New(&logOutput, "", 0)
	interceptor := RecoveryInterceptor(logger)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/service.Service/Method",
	}

	handlerFunc := func(ctx context.Context, req any) (any, error) {
		panic("panic error!")
	}

	_, err := interceptor(context.Background(), nil, info, handlerFunc)

	assert.NotEqual(t, logOutput.String(), 0, "log output should be empty")
	assert.Equal(t, codes.Internal, status.Code(err))
}
