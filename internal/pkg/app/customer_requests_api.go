package app

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (a *Application) GetAllCustomerRequests(c *gin.Context) {
	var request schemes.GetAllCustomerRequestReq
	request.RecordStatus = 10
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	customerRequests, err := a.repo.GetAllCustomerRequests(request.FormationDateStart, request.FormationDateEnd, request.RecordStatus)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputCustomerRequests := make([]schemes.CustomerRequestOutputResponse, len(customerRequests))
	for i, customerRequest := range customerRequests {
		serviceRequests, err := a.repo.GetServiceRequestsByCustId(customerRequest.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		outputCustomerRequests[i] = schemes.ConvertCustomerRequestResponse(&customerRequest, serviceRequests)
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

	serviceRequests, err := a.repo.GetServiceRequestsByCustId(request.CustomerRequestId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.CustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest, serviceRequests), DevelopmentServices: developmentServices})
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

	serviceRequests, err := a.repo.GetServiceRequestsByCustId(customerRequest.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateCustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest, serviceRequests)})
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

	if err := paymentRequest(customerRequest.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`payment service is unavailable: {%s}`, err))
		return
	}
	paymentStatus := ds.PaymentStarted
	customerRequest.PaymentStatus = &paymentStatus
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

func (a *Application) Payment(c *gin.Context) {
	var request schemes.PaymentReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Token != a.config.Token {
		c.AbortWithStatus(http.StatusForbidden)
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
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	var paymentStatus string
	if request.PaymentStatus {
		paymentStatus = "1"
	} else {
		paymentStatus = "0"
	}
	customerRequest.PaymentStatus = &paymentStatus

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
