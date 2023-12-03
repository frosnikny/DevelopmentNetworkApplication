package app

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

func (a *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension
	log.Println(imageName)
	_, err = a.minioClient.PutObject(c, a.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", a.config.MinioEndpoint, a.config.BucketName, imageName)
	return &imageURL, nil
}

func (a *Application) GetDevelopmentService(c *gin.Context) {
	var request schemes.DevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentService, err := a.repo.GetDevelopmentServiceByID(request.DevelopmentServiceId)
	if err != nil {
		log.Println("here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("услуга по разработке не найдена"))
		return
	}
	c.JSON(http.StatusOK, developmentService)
}

func (a *Application) GetAllDevelopmentServices(c *gin.Context) {
	log.Println(c.Request.MultipartForm)

	var request schemes.GetAllDevelopmentServicesReq
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentServices, err := a.repo.GetDevelopmentServicesByName(request.DevelopmentServiceName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	draftCustomerRequest, err := a.repo.GetDraftCustomerRequest(a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllDevelopmentServicesResponse{DraftCustomerRequest: nil, DevelopmentServices: developmentServices}
	if draftCustomerRequest != nil {
		response.DraftCustomerRequest = &schemes.CustomerRequests{UUID: draftCustomerRequest.UUID}
		developmentServices, err := a.repo.GetServiceRequests(draftCustomerRequest.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftCustomerRequest.DevelopmentServicesCount = len(developmentServices)
	}
	c.JSON(http.StatusOK, response)
}

func (a *Application) DeleteDevelopmentService(c *gin.Context) {
	var request schemes.DevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentService, err := a.repo.GetDevelopmentServiceByID(request.DevelopmentServiceId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("услуга по разработке не найдена"))
		return
	}
	developmentService.RecordStatus = ds.DSDeleted
	if err := a.repo.SaveDevelopmentService(developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (a *Application) ChangeDevelopmentService(c *gin.Context) {
	var request schemes.AddDevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentService, err := a.repo.GetDevelopmentServiceByID(request.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("контейнер не найден"))
		return
	}
	if request.Title != "" {
		developmentService.Title = request.Title
	}
	if request.Description != "" {
		developmentService.Description = request.Description
	}
	if request.Price != 0 {
		developmentService.Price = request.Price
	}
	if request.Technology != "" {
		developmentService.Technology = request.Technology
	}
	if request.DetailedCost != "" {
		developmentService.DetailedCost = request.DetailedCost
	}
	if request.Image != nil {
		imageUrl, err := a.uploadImage(c, request.Image, developmentService.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		developmentService.ImageUrl = imageUrl
	}

	if err := a.repo.SaveDevelopmentService(developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, developmentService)
}

func (a *Application) AddDevelopmentService(c *gin.Context) {
	var request schemes.AddDevelopmentServiceReq
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Println(c.Request.MultipartForm)

	developmentService := request.DevelopmentService
	if err := a.repo.AddDevelopmentService(&developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageUrl, err := a.uploadImage(c, request.Image, developmentService.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Println(imageUrl)
		developmentService.ImageUrl = imageUrl
	}

	if err := a.repo.SaveDevelopmentService(&developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (a *Application) AddToCustomerRequest(c *gin.Context) {
	var request schemes.DevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Println(request.DevelopmentServiceId)
	var err error

	// Проверить существует ли услуга по разработке
	developmentService, err := a.repo.GetDevelopmentServiceByID(request.DevelopmentServiceId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("контейнер не найден"))
		return
	}

	// Получить черновую заявку
	var customerRequest *ds.CustomerRequest
	customerRequest, err = a.repo.GetDraftCustomerRequest(a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		//log.Println("HEHEHEHE")
		customerRequest, err = a.repo.CreateDraftCustomerRequest(a.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	//log.Println("HEHEHE")

	// Создать связь между перевозкой и контейнером
	if err = a.repo.AddToCustomerRequest(customerRequest.UUID, request.DevelopmentServiceId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Вернуть список всех контейнеров в перевозке
	var developmentServices []ds.DevelopmentService
	developmentServices, err = a.repo.GetServiceRequests(customerRequest.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllDevelopmentServicesResponse{DevelopmentServices: developmentServices})
}

func (a *Application) getCustomer() string {
	return "1d6b2213-f5e5-4eb8-939d-3ab21f60108f"
}

func (a *Application) getModerator() *string {
	moderatorId := "01a6cec6-954d-4ce9-aeb1-3850d00162b4"
	return &moderatorId
}

func (a *Application) GetAllCustomerRequests(c *gin.Context) {
	var request schemes.GetAllCustomerRequestReq
	request.RecordStatus = 10
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//log.Println(request.RecordStatus)

	customerRequests, err := a.repo.GetAllCustomerRequests(request.FormationDateStart, request.FormationDateEnd, request.RecordStatus)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputCustomerRequests := make([]schemes.CustomerRequestOutputResponse, len(customerRequests))
	for i, customerRequest := range customerRequests {
		outputCustomerRequests[i] = schemes.ConvertCustomerRequestResponse(&customerRequest)
	}
	c.JSON(http.StatusOK, schemes.AllCustomerRequestsResponse{CustomerRequests: outputCustomerRequests})
}

func (a *Application) GetCustomerRequest(c *gin.Context) {
	var request schemes.CustomerRequestReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	customerRequest, err := a.repo.GetCustomerRequestById(request.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заявка не найдена"))
		return
	}

	developmentServices, err := a.repo.GetServiceRequests(request.CustomerRequestId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.CustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest), DevelopmentServices: developmentServices})
}

func (a *Application) UpdateCustomerRequest(c *gin.Context) {
	var request schemes.UpdateCustomerRequestReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("перевозка не найдена"))
		return
	}
	customerRequest.WorkSpecification = request.WorkSpecification
	if err = a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateCustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest)})
}

func (a *Application) DeleteCustomerRequest(c *gin.Context) {
	var request schemes.CustomerRequestReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	customerRequest, err := a.repo.GetCustomerRequestById(request.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("перевозка не найдена"))
		return
	}
	customerRequest.RecordStatus = ds.CRDeleted

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Application) DeleteFromCustomerRequest(c *gin.Context) {
	var request schemes.DeleteFromCustomerRequestReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	customerRequest, err := a.repo.GetCustomerRequestById(request.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заказ не найден"))
		return
	}
	if customerRequest.RecordStatus != ds.CRDraft {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать заказ со статусом: %d", customerRequest.RecordStatus))
		return
	}

	if err := a.repo.DeleteFromCustomerRequest(request.CustomerRequestId, request.DevelopmentServiceId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	developmentServices, err := a.repo.GetServiceRequests(request.CustomerRequestId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllDevelopmentServicesResponse{DevelopmentServices: developmentServices})
}

func (a *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заказ не найден"))
		return
	}
	if customerRequest.RecordStatus != ds.CRDraft {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать заказ со статусом %d", customerRequest.RecordStatus))
		return
	}
	customerRequest.RecordStatus = ds.CRWorks
	customerRequest.FormationDate = time.Now()

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.RecordStatus != ds.CRCompleted && request.RecordStatus != ds.CRDeclined {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("статус %d запрещен", request.RecordStatus))
		return
	}

	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("перевозка не найдена"))
		return
	}
	if customerRequest.RecordStatus != ds.CRWorks {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%d\" на \"%d\"", customerRequest.RecordStatus, request.RecordStatus))
		return
	}
	customerRequest.RecordStatus = request.RecordStatus
	customerRequest.ModeratorId = a.getModerator()
	if request.RecordStatus == ds.CRCompleted {
		customerRequest.CompletionDate = time.Now()
	}

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
