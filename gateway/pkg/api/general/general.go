package general

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct{}

type GrpcClient struct {
	Conn                     *grpc.ClientConn
	CustomClientInterceptors []grpc.UnaryClientInterceptor
	CustomDialOptions        []grpc.DialOption
}

type BaseConfig struct {
	LogLevel        string `env:"LOG_LEVEL,required" envDefault:"DEBUG"`
	Address         string `env:"HTTP_ADDRESS" envDefault:"0.0.0.0"`
	Port            string `env:"HTTP_PORT" envDefault:"8081"`
	ReadTimeout     int    `env:"READ_TIMEOUT" envDefault:"10"`
	WriteTimeout    int    `env:"WRITE_TIMEOUT" envDefault:"60"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT" envDefault:"30"`
	Listener        net.Listener
	MacHeaderBytes  int
	Router          *gin.Engine
}

type Resources struct {
	Log         *zap.SugaredLogger
	BaseCfg     *BaseConfig
	Engine      *gin.Engine
	GrpcClients map[string]*GrpcClient
}

type ServiceGeneral struct {
	*Resources
}

type BaseService struct {
	Config *BaseConfig
}

func Start(
	svc Service,
) {

	cfg := &BaseConfig{}
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Errorf("failed to read base config, %+v", err))
	}

	svc.SetBaseConfig(cfg)

	srv := http.Server{
		Addr:           net.JoinHostPort(cfg.Address, cfg.Port),
		Handler:        svc.GetEngine(),
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MacHeaderBytes,
	}

	run := func() {
		fmt.Printf("server listening on %s...", srv.Addr)
		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			_ = fmt.Errorf("failed to setup net.Listener, %s", err)
			panic(err)
		}

		defer func(ln net.Listener) {
			err := ln.Close()
			if err != nil {
				_ = fmt.Errorf("failed to close net.Listener, %s", err)
			}
		}(ln)

		if err := srv.Serve(ln); err != nil {
			_ = fmt.Errorf("failed to serve on %s, %s", srv.Addr, err)
			panic(err)
		}
	}

	teardown := func() {

		for service, grpcClient := range svc.GetGrpcClients() {
			if grpcClient.Conn != nil {
				err := grpcClient.Conn.Close()
				if err != nil {
					_ = fmt.Errorf("failed to close connection to service %s, %s", service, err)
				}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			_ = fmt.Errorf("failed to shutdown server, %s", err)
			panic(err)
		}
	}

	go run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	teardown()
}

type Service interface {
	SetBaseConfig(cfg *BaseConfig)
	GetGrpcClients() map[string]*GrpcClient
	GetEngine() *gin.Engine
}
