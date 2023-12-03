package app

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (a *Application) GetAllDevelopmentServices(c *gin.Context) {
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
	customerRequest, err = a.repo.GetDraftCustomerRequest(a.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		customerRequest, err = a.repo.CreateDraftCustomerRequest(a.getCustomer())
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
	developmentServices, err = a.repo.GetServiceRequests(customerRequest.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllDevelopmentServicesResponse{DevelopmentServices: developmentServices})
}
