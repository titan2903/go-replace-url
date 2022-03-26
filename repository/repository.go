package repository

import (
	"context"
	"replace-url-gin/config"
	"replace-url-gin/database"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	ctx    context.Context
	config *config.Config
	db     *sqlx.DB
}

func NewRepository(ctx context.Context, conf *config.Config) Repository {
	return repository{
		ctx:    ctx,
		config: conf,
		db:     database.Postgres(),
	}
}
