package api

import (
	"context"
	"fmt"
	awsx "github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type CircleUploadService interface {
	UploadImage(
		ctx context.Context,
		multiPartFile *multipart.FileHeader,
		circleId int64,
	) (string, error)
}

type circleUploadService struct {
	circleService     CircleService
	extStorageService awsx.S3Service
	config            *config.Config
	log               logger.Logger
}

func NewCircleUploadService(
	circleService CircleService,
	extStorageService awsx.S3Service,
	config *config.Config,
	log logger.Logger,
) CircleUploadService {
	return &circleUploadService{
		circleService:     circleService,
		extStorageService: extStorageService,
		config:            config,
		log:               log,
	}
}

func (c *circleUploadService) UploadImage(
	ctx context.Context,
	multiPartFile *multipart.FileHeader,
	circleId int64,
) (string, error) {
	contentFile, err := multiPartFile.Open()

	if err != nil {
		c.log.Errorf("error opening multipart file: %s", err)
		return "", err
	}

	defer contentFile.Close()

	err = c.detectMimeType(contentFile)

	if err != nil {
		return "", err
	}

	decodedImage, _, err := image.Decode(contentFile)

	size, calculated := utils.CalculatedImageSize(decodedImage, image.Point{X: 600, Y: 400})

	circleImage := decodedImage

	if calculated {
		circleImage = utils.ResizeImage(decodedImage, size.Max)
	}

	tempImageFile, err := os.CreateTemp("", "circle-image")

	if err != nil {
		c.log.Errorf("error creating temp file: %s", err)
		return "", err
	}

	defer os.Remove(tempImageFile.Name())

	err = png.Encode(tempImageFile, circleImage)

	if err != nil {
		c.log.Errorf("error encoding image file to png: %s", err)
		return "", err
	}

	_, _ = tempImageFile.Seek(0, 0)

	filePath := fmt.Sprintf("circle/image/%d/%s", circleId, "main.png")

	_, err = c.extStorageService.Upload(
		ctx,
		filePath,
		tempImageFile,
	)

	if err != nil {
		c.log.Errorf("error uploading file to external storage service: %s", err)
		return "", err
	}

	imageEndpoint := fmt.Sprintf("%s/%s", c.extStorageService.ObjectEndpoint(), filePath)

	updateCircleReq := &model.CircleUpdateRequest{
		ImageSrc: &imageEndpoint,
	}

	_, err = c.circleService.UpdateCircle(ctx, updateCircleReq)

	if err != nil {
		return "", err
	}

	return imageEndpoint, nil
}

func (c *circleUploadService) detectMimeType(contentFile multipart.File) error {
	bytes, err := io.ReadAll(contentFile)

	if err != nil {
		c.log.Errorf("error reading content of file: %s", err)
		return err
	}

	mimeType := http.DetectContentType(bytes)

	if !utils.IsImageMimeType(mimeType) {
		c.log.Infof("content type is not of mime type image")
		return fmt.Errorf("image is of wrong type")
	}

	_, err = contentFile.Seek(0, 0)

	if err != nil {
		c.log.Errorf("error resetting file to start position: %s", err)
		return err
	}

	return nil
}
