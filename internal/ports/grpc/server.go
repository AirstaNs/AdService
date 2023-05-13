package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"homework10/internal/app"
	"log"
	"net"
)

type grpcServer struct {
	gServer *grpc.Server
	log     *log.Logger
}

type Server interface {
	Start(network, address string) error
	Stop() error
}

func NewServer(loggerRPC *log.Logger, newApp app.App) Server {

	loggerInterceptor := LoggerInterceptor(loggerRPC)
	recoveryInterceptor := RecoveryInterceptor(loggerRPC)

	server := grpc.NewServer(
		grpc.Creds(nil),
		grpc.ChainUnaryInterceptor(loggerInterceptor, recoveryInterceptor),
	)
	RegisterAdServiceServer(server, GServer{App: newApp})

	return &grpcServer{gServer: server, log: loggerRPC}
}

func (s *grpcServer) Start(network, address string) error {
	lis, err := net.Listen(network, address)

	if err != nil {
		msg := fmt.Sprintf("failed to listen: %v", err)
		s.log.Println(msg)
		return fmt.Errorf(msg)
	}
	s.log.Printf("starting gRPC server on %s", lis.Addr())
	if err1 := s.gServer.Serve(lis); err1 != nil && err1 != grpc.ErrServerStopped {
		s.log.Fatalf("failed to start gRPC server: %v", err1)
		return err1
	}
	return nil
}

func (s *grpcServer) Stop() error {
	if s.gServer != nil {
		s.gServer.GracefulStop()
	} else {
		msg := "gServer is nil"
		s.log.Printf(msg)
		return fmt.Errorf(msg)
	}
	return nil
}
