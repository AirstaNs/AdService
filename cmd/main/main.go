package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	grpcImpl "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"homework10/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPort      = ":50051"
	httpPort      = ":18080"
	server        = "localhost"
	serverNetwork = "tcp"
)

func main() {
	repo := adrepo.New()
	uRep := userrepo.New()
	formatter := util.NewDateTimeFormatter(time.RFC3339)
	newApp := app.NewApp(repo, uRep, formatter)
	signals := append([]os.Signal{}, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	httpLogger := log.New(os.Stdout, "[HTTP] ", 0)
	rpcLogger := log.New(os.Stdout, "[gRPC] ", 0)
	sysLogger := log.New(os.Stdout, "[SYSTEM] ", log.Ldate|log.Ltime)

	g, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, signals...)

	hServer := httpgin.NewHTTPServer(server+httpPort, newApp, httpLogger)
	gServer := grpcImpl.NewServer(rpcLogger, newApp)

	g.Go(func() error {
		select {
		case sig := <-sigQuit:
			msg := fmt.Sprintf("received %s signal, starting graceful shutdown", sig.String())
			sysLogger.Println(msg)
			return fmt.Errorf("captured signal: %v", sig)
		case <-ctx.Done():
			return nil
		}
	})

	g.Go(func() error {
		errCh := make(chan error)
		defer func() {
			if err := hServer.Stop(); err != nil {
				sysLogger.Printf("error stopping HTTP server: %v\n", err)
			}
			close(errCh)
		}()

		go func() {
			if err := hServer.Start(server, httpPort); err != nil {
				errCh <- err
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	g.Go(func() error {
		errCh := make(chan error)
		defer func() {
			if err := gServer.Stop(); err != nil {
				sysLogger.Printf("error stopping HTTP server: %v\n", err)
			}
			close(errCh)
		}()

		go func() {
			if err := gServer.Start(serverNetwork, server+grpcPort); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err2 := g.Wait(); err2 != nil {
		sysLogger.Printf("gracefully shutting down the servers: %v\n", err2)
	}
}
