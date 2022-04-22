package repository

import (
	"context"
	"fmt"
	"replace-url-gin/config"
	"replace-url-gin/entity"

	"github.com/jmoiron/sqlx"
)

type BaseRepository interface {
	ReplaceImage() (entity.UploadFileModels, error)
	BulkUpdateImage(payloads entity.ModifyUploadFileModels) error
	UpdateUrlImage(payloads entity.ModifyUploadFileModelUrls) error
}

type baseRepository struct {
	ctx    context.Context
	config *config.Config
	db     *sqlx.DB
}

func (r repository) BaseRepository() BaseRepository {
	return &baseRepository{
		ctx:    r.ctx,
		config: r.config,
		db:     r.db,
	}
}

var (
	getUploadFile = `
    SELECT
      up.id,
      up.formats,
      up.url
    FROM upload_file up
  `

	updateUploadFile = `
    UPDATE upload_file
    SET formats = :formats
    WHERE id = :id
  `

	updateUrl = `
	UPDATE upload_file
	SET url = :url
	WHERE id = :id
  `
)

func (r baseRepository) ReplaceImage() (entity.UploadFileModels, error) {
	dest := entity.UploadFileModels{}
	err := r.db.SelectContext(r.ctx, &dest, getUploadFile)
	if err != nil {
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository ReplaceImage: %v\n", err.Error())
		return nil, err
	}

	return dest, nil
}

func (r baseRepository) BulkUpdateImage(payloads entity.ModifyUploadFileModels) error {
	tx, err := r.db.Beginx()
	if err != nil {
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage Begin: %v\n", err.Error())
		return err
	}

	for _, payload := range payloads {
		_, err = tx.NamedExecContext(r.ctx, updateUploadFile, payload)
		if err != nil {
			fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage NamedExec: %v\n", err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage Commit: %v\n", err.Error())
		return err
	}

	return nil
}

func (r baseRepository) UpdateUrlImage(payloads entity.ModifyUploadFileModelUrls) error {
	tx, err := r.db.Beginx()
	if err != nil {
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage Begin: %v\n", err.Error())
		return err
	}

	for _, payload := range payloads {
		_, err = tx.NamedExecContext(r.ctx, updateUrl, payload)
		if err != nil {
			fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage NamedExec: %v\n", err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkUpdateImage Commit: %v\n", err.Error())
		return err
	}

	return nil
}
