package app

import (
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// GetAllCustomerRequests @Summary		Получить все заказы
// @Tags		Заказы
// @Description	Возвращает все заказы с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус заказа"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирования"
// @Success		200 {object} schemes.AllCustomerRequestsResponse
// @Router		/api/requests [get]
func (a *Application) GetAllCustomerRequests(c *gin.Context) {
	var request schemes.GetAllCustomerRequestReq
	var err error
	request.RecordStatus = 10
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	log.Println(userId, userRole)
	var customerRequests []ds.CustomerRequest
	if userRole == role.Customer {
		customerRequests, err = a.repo.GetAllCustomerRequests(&userId, request.FormationDateStart, request.FormationDateEnd, request.RecordStatus)
	} else {
		customerRequests, err = a.repo.GetAllCustomerRequests(nil, request.FormationDateStart, request.FormationDateEnd, request.RecordStatus)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputCustomerRequests := make([]schemes.AllCustomerRequestOutputResponse, len(customerRequests))
	for i, customerRequest := range customerRequests {
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		outputCustomerRequests[i] = schemes.ConvertAllCustomerRequestResponse(&customerRequest)
	}
	c.JSON(http.StatusOK, schemes.AllCustomerRequestsResponse{CustomerRequests: outputCustomerRequests})
}

// GetCustomerRequest @Summary		Получить один заказ
// @Tags		Заказы
// @Description	Возвращает подробную информацию о заказе и его составе
// @Produce		json
// @Param		id path string true "id заказа"
// @Success		200 {object} schemes.CustomerRequestResponse
// @Router		/api/requests/{customer_request_id} [get]
func (a *Application) GetCustomerRequest(c *gin.Context) {
	var request schemes.CustomerRequestReq
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var customerRequest *ds.CustomerRequest
	if userRole == role.Moderator {
		customerRequest, err = a.repo.GetCustomerRequestById(request.CustomerRequestId, "")
	} else {
		customerRequest, err = a.repo.GetCustomerRequestById(request.CustomerRequestId, userId)
	}
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
	log.Println(developmentServices)

	serviceRequests, err := a.repo.GetServiceRequestsByCustId(request.CustomerRequestId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.CustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest, serviceRequests, developmentServices)})
}

type SwaggerUpdateTransportationRequest struct {
	WorkSpecification string `json:"work_specification"`
}

// UpdateCustomerRequest @Summary		Указать спецификацию заказа
// @Tags		Заказы
// @Description	Позволяет изменить спецификацию заказа и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		transport body SwaggerUpdateTransportationRequest true "Спецификация"
// @Success		200
// @Router		/api/requests [put]
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

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заявка не найдена"))
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

	developmentServices, err := a.repo.GetServiceRequests(customerRequest.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateCustomerRequestResponse{CustomerRequest: schemes.ConvertCustomerRequestResponse(customerRequest, serviceRequests, developmentServices)})
}

// DeleteCustomerRequest @Summary		Удалить заказ
// @Tags		Заказы
// @Description	Удаляет заказ
// @Success		200
// @Router		/api/requests [delete]
func (a *Application) DeleteCustomerRequest(c *gin.Context) {
	var request schemes.CustomerRequestReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.CustomerRequestId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заявка не найдена"))
		return
	}
	customerRequest.RecordStatus = ds.CRDeleted

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// DeleteFromCustomerRequest @Summary		Удалить услугу по разработке из чернового заказа
// @Tags		Заказы
// @Description	Удалить услугу по разработке из чернового заказа
// @Produce		json
// @Param		id path string true "id заказа"
// @Success		200
// @Router		/api/requests/{customer_request_id}/delete_development_service/{development_service_id} [delete]
func (a *Application) DeleteFromCustomerRequest(c *gin.Context) {
	var request schemes.DeleteFromCustomerRequestReq
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.CustomerRequestId, userId)
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

// UserConfirm @Summary		Сформировать заказ
// @Tags		Заказы
// @Description	Сформировать заказ пользователем
// @Success		200
// @Router		/api/requests/user_confirm [put]
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

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, userId)
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

	//if err := paymentRequest(customerRequest.UUID); err != nil {
	//	c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`payment service is unavailable: {%s}`, err))
	//	return
	//}
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

// ModeratorConfirm @Summary		Подтвердить заказ
// @Tags		Заказы
// @Description	Подтвердить или отменить заказ модератором
// @Param		id path string true "id заказа"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/requests/{id}/moderator_confirm [put]
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

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, userId)
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
	customerRequest.ModeratorId = &userId
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

	userId := getUserId(c)
	customerRequest, err := a.repo.GetCustomerRequestById(request.URI.CustomerRequestId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if customerRequest == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("заявка не найдена"))
		return
	}
	if customerRequest.RecordStatus != ds.CRWorks {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	var paymentStatus string
	if *request.PaymentStatus {
		paymentStatus = "1"
	} else {
		paymentStatus = "0"
	}
	//paymentStatus = request.PaymentStatus
	customerRequest.PaymentStatus = &paymentStatus

	if err := a.repo.SaveCustomerRequest(customerRequest); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
