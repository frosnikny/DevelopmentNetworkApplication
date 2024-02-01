package app

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetDevelopmentService @Summary		Получить услугу по разработке
// @Tags		Услуги по разработке
// @Description	Возвращает более подробную информацию об одной из услуг по разработке
// @Produce		json
// @Param		type query string false "тип услуги для фильтрации"
// @Success		200 {object} schemes.GetAllDevelopmentServicesResponse
// @Router		/api/devs/{development_service_id} [get]
func (a *Application) GetDevelopmentService(c *gin.Context) {
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
	c.JSON(http.StatusOK, developmentService)
}

// GetAllDevelopmentServices @Summary		Получить все услуги по разработке
// @Tags		Услуги по разработке
// @Description	Возвращает все доступные услуги по разработке с опциональной фильтрацией по типу
// @Produce		json
// @Param		id path string true "id услуги"
// @Success		200 {object} ds.DevelopmentService
// @Router		/api/devs [get]
func (a *Application) GetAllDevelopmentServices(c *gin.Context) {
	var request schemes.GetAllDevelopmentServicesReq

	log.Println("aaaaaa")

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentServices, err := a.repo.GetDevelopmentServicesByName(request.DevelopmentServiceName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userId := getUserId(c)
	log.Println("userId: ", userId)
	draftCustomerRequest, err := a.repo.GetDraftCustomerRequest(userId)
	if err != nil {
		log.Println("error draft")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllDevelopmentServicesResponse{DraftCustomerRequest: nil, DevelopmentServices: developmentServices}
	if draftCustomerRequest != nil {
		response.DraftCustomerRequest = &schemes.CustomerRequests{UUID: draftCustomerRequest.UUID}
		developmentServices, err := a.repo.GetDevServices(draftCustomerRequest.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftCustomerRequest.DevelopmentServicesCount = len(developmentServices)
	}
	c.JSON(http.StatusOK, response)
}

// DeleteDevelopmentService @Summary		Удалить услугу по разработке
// @Tags		Услуги по разработке
// @Description	Удаляет услугу по разработке по id
// @Param		id path string true "id услуги"
// @Success		200
// @Router		/api/devs/{development_service_id} [delete]
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
	if err := a.deleteImage(c, developmentService.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	developmentService.ImageUrl = nil
	developmentService.RecordStatus = ds.DSDeleted
	if err := a.repo.SaveDevelopmentService(developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// ChangeDevelopmentService @Summary		Изменить услугу по разработке
// @Tags		Услуги по разработке
// @Description	Изменить данные полей об услуге по разработке
// @Accept		mpfd
// @Param		id path string true "Идентификатор услуги" format:"uuid"
// @Param		image formData file false "Изображение услуги по разработке"
// @Param		title formData string true "Название" format:"string" maxLength:100
// @Param		description formData string true "Описание" format:"string" maxLength:500
// @Param		price formData int true "Цена" format:"int"
// @Param		technology formData string true "Технологии" format:"string" maxLength:100
// @Param		detailed_price formData float32 true "Цена за день" format:"real"
// @Success		200
// @Router		/api/devs/{development_service_id} [put]
func (a *Application) ChangeDevelopmentService(c *gin.Context) {
	var request schemes.ChangeDevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentService, err := a.repo.GetDevelopmentServiceByID(request.DevelopmentServiceId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Println(developmentService)

	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("услуга по разработке не найдена"))
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
	if request.DetailedPrice != 0 {
		developmentService.DetailedPrice = request.DetailedPrice
	}
	if request.Image != nil {
		if developmentService.ImageUrl != nil {
			if err := a.deleteImage(c, developmentService.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
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

// AddDevelopmentService @Summary		Добавить услугу по разработке
// @Tags		Услуги по разработке
// @Description	Добавить новую услугу по разработке
// @Accept		mpfd
// @Param		image formData file false "Изображение услуги по разработке"
// @Param		title formData string true "Название" format:"string" maxLength:100
// @Param		description formData string true "Описание" format:"string" maxLength:500
// @Param		price formData int true "Цена" format:"int"
// @Param		technology formData string true "Технологии" format:"string" maxLength:100
// @Param		detailed_price formData float32 true "Цена за день" format:"real"
// @Success		200
// @Router		/api/devs [post]
func (a *Application) AddDevelopmentService(c *gin.Context) {
	var request schemes.AddDevelopmentServiceReq
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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
		developmentService.ImageUrl = imageUrl
	}

	if err := a.repo.SaveDevelopmentService(&developmentService); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// AddToCustomerRequest @Summary		Добавить в заказ
// @Tags		Услуги по разработке
// @Description	Добавить выбранную услугу по разработке в черновик заказа
// @Param		id path string true "id услуги"
// @Success		200
// @Router		/api/devs/{development_service_id}/add_to_request [post]
func (a *Application) AddToCustomerRequest(c *gin.Context) {
	var request schemes.DevelopmentServiceReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
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
	userId := getUserId(c)
	customerRequest, err = a.repo.GetDraftCustomerRequest(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		customerRequest, err = a.repo.CreateDraftCustomerRequest(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Создать связь между заказом и разработкой
	if err = a.repo.AddToCustomerRequest(customerRequest.UUID, request.DevelopmentServiceId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Вернуть список всех разработок в заказе
	var developmentServices []ds.DevelopmentService
	developmentServices, err = a.repo.GetDevServices(customerRequest.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllDevelopmentServicesResponse{DevelopmentServices: developmentServices})
}
