package usecase

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"replace-url-gin/config"
	"replace-url-gin/entity"
	"replace-url-gin/repository"
	"strings"
	"sync"
)

type BaseUsecase interface {
	ReplaceImage() error
}

type baseUsecase struct {
	ctx        context.Context
	config     *config.Config
	repository repository.BaseRepository
}

func (u usecase) BaseHandler() BaseUsecase {
	return &baseUsecase{
		ctx:        u.ctx,
		config:     u.config,
		repository: u.repository.BaseRepository(),
	}
}

func (u baseUsecase) ReplaceImage() error {
	imageData, err := u.repository.ReplaceImage()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	modifyChan := make(chan entity.ModifyUploadFileModel)
	for _, val := range imageData {
		wg.Add(1)

		go u.readBuffer(wg, val, modifyChan)
	}

	go func() {
		wg.Wait()
		close(modifyChan)
	}()

	imagePayload := entity.ModifyUploadFileModels{}
	for data := range modifyChan {
		imagePayload = append(imagePayload, data)
	}

	err = u.repository.BulkUpdateImage(imagePayload)
	if err != nil {
		return err
	}

	return nil
}

func (u baseUsecase) readBuffer(wg *sync.WaitGroup, data entity.UploadFileModel, modify chan entity.ModifyUploadFileModel) {
	defer wg.Done()

	if data.Format == nil {
		return
	}

	var toStruct entity.UploadFileFormat
	err := json.Unmarshal([]byte(*data.Format), &toStruct)
	if err != nil {
		log.Fatal("Error json.Unmarshal buffer", err)
	}

	if toStruct.Medium != nil {
		toStruct.Medium.Url = u.replaceImage(toStruct.Medium.Url, data.Url)
	}

	if toStruct.Large != nil {
		toStruct.Large.Url = u.replaceImage(toStruct.Large.Url, data.Url)
	}

	if toStruct.Small != nil {
		toStruct.Small.Url = u.replaceImage(toStruct.Small.Url, data.Url)
	}

	if toStruct.Thumbnail != nil {
		toStruct.Thumbnail.Url = u.replaceImage(toStruct.Thumbnail.Url, data.Url)
	}

	toJSON, err := json.Marshal(toStruct)
	if err != nil {
		log.Fatal("Error json.Marshal buffer", err)
	}

	dest := entity.ModifyUploadFileModel{
		Id:     data.Id,
		Format: string(toJSON),
	}

	modify <- dest
}

func (u baseUsecase) replaceImage(imageUrl, originalUrl string) string {
	original := u.extractSubDomain(originalUrl)
	image := u.extractSubDomain(imageUrl)

	if image == nil || original == nil {
		return imageUrl
	}

	if !strings.EqualFold(*image, *original) {
		return strings.Replace(imageUrl, *image, *original, -1)
	}

	return imageUrl
}

func (u baseUsecase) extractSubDomain(imageUrl string) *string {
	url, err := url.Parse(imageUrl)
	if err != nil {
		return nil
	}
	arrDomain := strings.Split(url.Host, ".")
	if len(arrDomain) == 0 {
		return nil
	}
	return &arrDomain[0]
}
