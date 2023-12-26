package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func (a *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" {
		return nil, fmt.Errorf("разрешены только jpg изображения")
	}
	imageName := UUID + extension
	//log.Println(imageName)
	_, err = a.minioClient.PutObject(c, a.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("http://%s/%s/%s", a.config.MinioEndpoint, a.config.BucketName, imageName)
	return &imageURL, nil
}

func (a *Application) deleteImage(c *gin.Context, UUID string) error {
	imageName := UUID + ".jpg"
	//fmt.Println(imageName)
	err := a.minioClient.RemoveObject(c, a.config.BucketName, imageName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *Application) getCustomer() string {
	return "1d6b2213-f5e5-4eb8-939d-3ab21f60108f"
}

func (a *Application) getModerator() *string {
	moderatorId := "01a6cec6-954d-4ce9-aeb1-3850d00162b4"
	return &moderatorId
}

func paymentRequest(customerRequestId string) error {
	url := "http://localhost:8000/"
	payload := fmt.Sprintf(`{"customer_request_id": "%s"}`, customerRequestId)

	resp, err := http.Post(url, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf(`delivery failed with status: {%s}`, resp.Status)
	}
	return nil
}
