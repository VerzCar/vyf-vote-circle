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

	// Open and validate file
	contentFile, err := c.openAndValidateFile(multiPartFile)
	if err != nil {
		return "", err
	}
	defer contentFile.Close()

	// Resize image if needed
	imageForUpload, err := c.resizeImageIfNeeded(contentFile)
	if err != nil {
		return "", err
	}

	// Load the image to the temporary file
	tempImageFile, err := c.loadTempImageFile(imageForUpload)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempImageFile.Name())

	imageEndpoint, err := c.uploadFile(ctx, imageForUpload, circleId, tempImageFile)
	if err != nil {
		return "", err
	}

	return imageEndpoint, nil
}

func (c *circleUploadService) openAndValidateFile(multiPartFile *multipart.FileHeader) (
	contentFile multipart.File,
	err error,
) {
	contentFile, err = multiPartFile.Open()
	if err != nil {
		c.log.Errorf("error opening multipart file: %s", err)
		return
	}

	err = c.detectMimeType(contentFile)
	return
}

func (c *circleUploadService) resizeImageIfNeeded(contentFile multipart.File) (resizableImage image.Image, err error) {
	resizableImage, _, err = image.Decode(contentFile)
	if err != nil {
		return
	}

	size, calculated := utils.CalculatedImageSize(resizableImage, image.Point{X: 600, Y: 400})
	if calculated {
		resizableImage = utils.ResizeImage(resizableImage, size.Max)
	}
	return
}

func (c *circleUploadService) loadTempImageFile(circleImage image.Image) (*os.File, error) {
	tempImageFile, err := os.CreateTemp("", "circle-image")
	if err != nil {
		c.log.Errorf("error creating temp file: %s", err)
		return nil, err
	}

	err = png.Encode(tempImageFile, circleImage)
	if err != nil {
		c.log.Errorf("error encoding image file to png: %s", err)
		return nil, err
	}

	_, err = tempImageFile.Seek(0, 0)
	return tempImageFile, err
}

func (c *circleUploadService) uploadFile(
	ctx context.Context,
	circleImage image.Image,
	circleId int64,
	tempImageFile *os.File,
) (string, error) {
	filePath := fmt.Sprintf("circle/image/%d/%s", circleId, "main.png")
	_, err := c.extStorageService.Upload(ctx, filePath, tempImageFile)
	if err != nil {
		c.log.Errorf("error uploading file to external storage service: %s", err)
		return "", err
	}

	imageEndpoint := fmt.Sprintf("%s/%s", c.extStorageService.ObjectEndpoint(), filePath)
	updateCircleReq := &model.CircleUpdateRequest{
		ImageSrc: &imageEndpoint,
	}
	_, err = c.circleService.UpdateCircle(ctx, circleId, updateCircleReq)
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
