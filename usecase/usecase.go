package usecase

import (
	"context"
	"replace-url-gin/config"
	"replace-url-gin/repository"
)

type usecase struct {
	ctx        context.Context
	config     *config.Config
	repository repository.Repository
}

func NewUsecase(ctx context.Context, conf *config.Config) Usecase {
	return usecase{
		ctx:        ctx,
		config:     conf,
		repository: repository.NewRepository(ctx, conf),
	}
}
