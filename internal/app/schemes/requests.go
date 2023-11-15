package schemes

type DevelopmentServiceRequest struct {
	DevelopmentServiceId string `uri:"development_service_id" binding:"required,uuid"`
}
