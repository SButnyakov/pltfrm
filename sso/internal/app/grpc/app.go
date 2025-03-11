package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	slog.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (a *App) Stop() error {
	slog.Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	return nil
}
