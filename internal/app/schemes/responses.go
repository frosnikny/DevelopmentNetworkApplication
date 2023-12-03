package schemes

import (
	"awesomeProject/internal/app/ds"
)

type AllDevelopmentServicesResponse struct {
	DevelopmentServices []ds.DevelopmentService `json:"development_services"`
}

type CustomerRequests struct {
	UUID                     string `json:"uuid"`
	DevelopmentServicesCount int    `json:"development_services_count"`
}

type GetAllDevelopmentServicesResponse struct {
	DraftCustomerRequest *CustomerRequests       `json:"draft_customer_request"`
	DevelopmentServices  []ds.DevelopmentService `json:"development_services"`
}

type AllCustomerRequestsResponse struct {
	CustomerRequests []CustomerRequestOutputResponse `json:"customer_requests"`
}

type CustomerRequestResponse struct {
	CustomerRequest     CustomerRequestOutputResponse `json:"customer_request"`
	DevelopmentServices []ds.DevelopmentService       `json:"development_services"`
}

type UpdateCustomerRequestResponse struct {
	CustomerRequest CustomerRequestOutputResponse `json:"customer_request"`
}

type CustomerRequestOutputResponse struct {
	UUID              string  `json:"uuid"`
	RecordStatus      uint    `json:"record_status"`
	CreationDate      string  `json:"creation_date"`
	FormationDate     *string `json:"formation_date"`
	CompletionDate    *string `json:"completion_date"`
	WorkSpecification string  `json:"work_specification"`
	Moderator         *string `json:"moderator"`
	Creator           string  `json:"creator"`
}

func ConvertCustomerRequestResponse(customerRequest *ds.CustomerRequest) CustomerRequestOutputResponse {
	output := CustomerRequestOutputResponse{
		UUID:              customerRequest.UUID,
		RecordStatus:      customerRequest.RecordStatus,
		CreationDate:      customerRequest.CreationDate.Format("2006-01-02 15:04:05"),
		WorkSpecification: customerRequest.WorkSpecification,
		Creator:           customerRequest.Creator.Name,
	}

	if !customerRequest.FormationDate.IsZero() { // != nil
		formationDate := customerRequest.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if !customerRequest.CompletionDate.IsZero() { // != nil
		completionDate := customerRequest.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if customerRequest.Moderator != nil {
		output.Moderator = &customerRequest.Moderator.Name
	}

	return output
}
