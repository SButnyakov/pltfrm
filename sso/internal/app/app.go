package app

import (
	"sso/internal/app/grpc"
	"sso/internal/config"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(cfg *config.Config) *App {
	grpcApp := grpcapp.New(cfg.GRPC.Port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
