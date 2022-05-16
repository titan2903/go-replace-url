package usecase

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"replace-url-gin/config"
	"replace-url-gin/entity"
	"replace-url-gin/repository"
	"strings"
	"sync"
)

type BaseUsecase interface {
	ReplaceImage() error
	ReplaceImageUrl() error
	BulkInsertNumber() error
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

func (u baseUsecase) ReplaceImageUrl() error {
	imageData, err := u.repository.ReplaceImage()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	modifyChan := make(chan entity.ModifyUploadFileModelUrl)
	for _, val := range imageData {
		wg.Add(1)

		go u.readBufferUrl(wg, val, modifyChan)
	}

	go func() {
		wg.Wait()
		close(modifyChan)
	}()

	imagePayload := entity.ModifyUploadFileModelUrls{}
	for data := range modifyChan {
		imagePayload = append(imagePayload, data)
	}

	err = u.repository.UpdateUrlImage(imagePayload)
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

func (u baseUsecase) readBufferUrl(wg *sync.WaitGroup, data entity.UploadFileModel, modify chan entity.ModifyUploadFileModelUrl) {
	defer wg.Done()

	if data.Url == "" {
		return
	}

	if data.Url != "" {
		data.Url = u.replaceSingleImage(data.Url, "https://dlo75jjjtr3ck.cloudfront.net/emi-banner.jpg")
	}

	dest := entity.ModifyUploadFileModelUrl{
		Id:  data.Id,
		Url: data.Url,
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

func (u baseUsecase) replaceSingleImage(imageUrl, replaceUrl string) string {
	uat := u.extractSubDomain(replaceUrl)
	dev := u.extractSubDomain(imageUrl)

	if dev == nil || uat == nil {
		return imageUrl
	}

	if !strings.EqualFold(*dev, *uat) {
		return strings.Replace(imageUrl, *dev, *uat, -1)
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

func (u baseUsecase) BulkInsertNumber() error {
	file, err := os.Open("data/number.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = u.insertNumberData(file)
	if err != nil {
		panic(err)
	}

	return nil
}

func (u baseUsecase) insertNumberData(file *os.File) error {
	reader := csv.NewReader(file)
	fmt.Printf("reader: %+v \n", reader)

	var arrayNumbers entity.PhoneNumbers

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if len(line) == 0 {
			continue
		}

		if line[0] != "phone_number" {
			fmt.Printf("line 0: %+v \n", line[0])

			arrayNumbers = append(arrayNumbers, entity.PhoneNumber{
				PhoneNumber: line[0],
			})

		}
	}

	fmt.Printf("numbers: %+v", arrayNumbers)
	err := u.repository.BulkInsertNumber(arrayNumbers)
	if err != nil {
		return err
	}

	return nil
}
