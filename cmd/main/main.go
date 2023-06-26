package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework10/internal/adapters/repository/adrepo"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/app"
	"homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"homework10/internal/util"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	//grpcPort = ":50051"
	//httpPort      = ":8080"
	server        = "localhost"
	serverNetwork = "tcp"
)

var PORT_REST string
var PORT_gRPC string

//func main() {
//	r := gin.New()
//	r.GET("/get", func(context *gin.Context) {
//		context.JSON(http.StatusOK, gin.H{"HERLLO HORLD": "DGOGOGOGO"})
//	})
//	r.Run()
//}

// "example.crt", "example.key"
func main() {
	PORT_REST = setPortEnv(3333, "PORT_REST", ":")
	PORT_gRPC = setPortEnv(5051, "PORT_gRPC", ":")

	cert := flag.String("pCertFile", "not_found", "get path to Cert File for TLS")
	key := flag.String("pKeyFile", "not_found", "get path to Key File for TLS")

	flag.Parse()
	fmt.Println(PORT_REST)
	repo := adrepo.New()
	uRep := userrepo.New()
	formatter := util.NewDateTimeFormatter(time.RFC3339)
	newApp := app.NewApp(repo, uRep, formatter)
	signals := append([]os.Signal{}, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	httpLogger := log.New(os.Stdout, "[HTTP] ", 0)
	rpcLogger := log.New(os.Stdout, "[gRPC] ", 0)
	sysLogger := log.New(os.Stdout, "[SYSTEM] ", log.Ldate|log.Ltime)

	g, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, signals...)

	hServer := httpgin.NewHTTPServer(PORT_REST, newApp, httpLogger, *cert, *key)
	gServer := grpc.NewServer(rpcLogger, newApp)

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
			if err := hServer.Start("", PORT_REST); err != nil {
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
			if err := gServer.Start(serverNetwork, PORT_gRPC); err != nil {
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

func setPortEnv(dPort int, name, sep string) (port string) {
	port, ok := os.LookupEnv(name)
	if ok {
		port = sep + port
	} else {
		port = sep + strconv.Itoa(dPort)
	}
	return port
}
