package schemes

import (
	"awesomeProject/internal/app/ds"
	"mime/multipart"
	"time"
)

type DevelopmentServiceReq struct {
	DevelopmentServiceId string `uri:"development_service_id" binding:"required,uuid"`
}

type GetAllDevelopmentServicesReq struct {
	DevelopmentServiceName string `form:"name" json:"name"`
}

type AddDevelopmentServiceReq struct {
	ds.DevelopmentService
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type GetAllCustomerRequestReq struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04:05"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04:05"`
	RecordStatus       uint       `form:"status" json:"status" binding:"required"`
}

type CustomerRequestReq struct {
	CustomerRequestId string `uri:"customer_request_id" binding:"required,uuid"`
}

type UpdateCustomerRequestReq struct {
	URI struct {
		CustomerRequestId string `uri:"customer_request_id" binding:"required,uuid"`
	}
	WorkSpecification string `form:"work_specification" json:"work_specification" binding:"required,max=50"`
}

type DeleteFromCustomerRequestReq struct {
	CustomerRequestId    string `uri:"customer_request_id" binding:"required,uuid"`
	DevelopmentServiceId string `uri:"development_service_id" binding:"required,uuid"`
}

type UserConfirmReq struct {
	URI struct {
		CustomerRequestId string `uri:"customer_request_id" binding:"required,uuid"`
	}
	Confirm bool `form:"confirm" binding:"required"`
}

type ModeratorConfirmReq struct {
	URI struct {
		CustomerRequestId string `uri:"customer_request_id" binding:"required,uuid"`
	}
	RecordStatus uint `form:"status" binding:"required"`
	Confirm      bool `form:"confirm" binding:"required"`
}

type PaymentReq struct {
	URI struct {
		CustomerRequestId string `uri:"customer_request_id" binding:"required,uuid"`
	}
	PaymentStatus bool   `form:"payment_status" binding:"required"`
	Token         string `form:"token" binding:"required"`
}
