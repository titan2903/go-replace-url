package handler

import (
	"context"
	"replace-url-gin/config"
	"replace-url-gin/usecase"
)

type handler struct {
	ctx     context.Context
	config  *config.Config
	usecase usecase.Usecase
}

func NewHandler(ctx context.Context, conf *config.Config) Handler {
	return handler{
		ctx:     ctx,
		config:  conf,
		usecase: usecase.NewUsecase(ctx, conf),
	}
}
