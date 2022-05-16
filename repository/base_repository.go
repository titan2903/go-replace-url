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
	BulkInsertNumber(payloads entity.PhoneNumbers) error
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

	bulkInsertNumber = `
		INSERT INTO customer_numbers (phone_number)
		VALUES (:phone_number)
		ON CONFLICT (phone_number)
		DO UPDATE
		SET phone_number = :phone_number;
	`

	getCustomerNumbers = `
		SELECT
			phone_number
		FROM customer_numbers
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

func (r baseRepository) GetCustomerNumbers() (entity.PhoneNumbers, error) {
	dest := entity.PhoneNumbers{}
	err := r.db.SelectContext(r.ctx, &dest, getCustomerNumbers)
	if err != nil {
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository ReplaceImage: %v\n", err.Error())
		return nil, err
	}

	return dest, nil
}

func (r baseRepository) BulkInsertNumber(payloads entity.PhoneNumbers) error {
	tx, err := r.db.Beginx()
	if err != nil {
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkInsertNumber Begin: %v\n", err.Error())
		return err

	}

	for _, payload := range payloads {
		fmt.Printf("payload: %+v", payload)
		_, err = tx.NamedExecContext(r.ctx, bulkInsertNumber, payload)
		if err != nil {
			fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkInsertNumber NamedExec: %v\n", err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("\033[1;31m [ERROR] \033[0m Repository BulkInsertNumber Commit: %v\n", err.Error())
		return err
	}

	return nil
}
