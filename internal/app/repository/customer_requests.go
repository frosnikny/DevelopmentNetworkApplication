package repository

import (
	"awesomeProject/internal/app/ds"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

func (r *Repository) GetAllCustomerRequests(customerId *string, formationDateStart, formationDateEnd *time.Time, recordStatus uint) ([]ds.CustomerRequest, error) {
	var customerRequests []ds.CustomerRequest
	log.Println(recordStatus)
	query := r.db.Preload("Creator").Preload("Moderator")
	if recordStatus != 10 {
		query = query.Where("record_status = ?", recordStatus)
	}
	query = query.Where("record_status != ?", ds.CRDeleted)
	if customerId != nil {
		log.Println("cus: ", customerId)
		query = query.Where("creator_id = ?", *customerId)
	}
	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}

	if err := query.Find(&customerRequests).Error; err != nil {
		return nil, err
	}
	return customerRequests, nil
}

func (r *Repository) GetDraftCustomerRequest(customerId string) (*ds.CustomerRequest, error) {
	customerRequest := &ds.CustomerRequest{}
	if customerId == "" {
		return nil, nil
	}
	query := r.db.Table("customer_requests").
		Where("record_status = ?", ds.CRDraft)
	query = query.Where("creator_id", customerId)
	err := query.First(customerRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return customerRequest, nil
}

func (r *Repository) CreateDraftCustomerRequest(customerId string) (*ds.CustomerRequest, error) {
	customerRequest := &ds.CustomerRequest{CreationDate: time.Now(), CreatorId: customerId, RecordStatus: ds.CRDraft}
	err := r.db.Create(customerRequest).Error
	if err != nil {
		return nil, err
	}
	return customerRequest, nil
}

func (r *Repository) AddToCustomerRequest(customerRequestId, developmentServiceId string) error {
	serviceRequest := ds.ServiceRequest{CustomerRequestId: customerRequestId, DevelopmentServiceId: developmentServiceId}
	err := r.db.Create(&serviceRequest).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetServiceRequests(customerRequestId string) ([]ds.DevelopmentService, error) {
	var developmentServices []ds.DevelopmentService

	err := r.db.Table("service_requests").
		Select("development_services.*").
		Joins("JOIN development_services ON service_requests.development_service_id = development_services.uuid").
		Where(ds.ServiceRequest{CustomerRequestId: customerRequestId}).
		Scan(&developmentServices).Error

	if err != nil {
		return nil, err
	}
	return developmentServices, nil
}

func (r *Repository) GetServiceRequestsByCustId(customerRequestId string) ([]ds.ServiceRequest, error) {
	var serviceRequests []ds.ServiceRequest

	err := r.db.Table("service_requests").
		Where(ds.ServiceRequest{CustomerRequestId: customerRequestId}).
		Scan(&serviceRequests).Error

	if err != nil {
		return nil, err
	}
	return serviceRequests, nil
}

func (r *Repository) GetCustomerRequestById(customerRequestId, userId string) (*ds.CustomerRequest, error) {
	customerRequest := &ds.CustomerRequest{}
	query := r.db.Preload("Moderator").Preload("Creator").
		Where("record_status != ?", ds.CRDeleted)
	//if userId != "" {
	//	query = query.Where("creator_id = ?", userId)
	//}
	query = query.First(customerRequest, ds.CustomerRequest{UUID: customerRequestId})
	err := query.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return customerRequest, nil
}

func (r *Repository) SaveCustomerRequest(customerRequest *ds.CustomerRequest) error {
	err := r.db.Save(customerRequest).Error
	if err != nil {
		log.Println("e")
		return err
	}
	return nil
}

func (r *Repository) DeleteFromCustomerRequest(customerRequestId, DevelopmentServiceId string) error {
	err := r.db.Delete(&ds.ServiceRequest{CustomerRequestId: customerRequestId, DevelopmentServiceId: DevelopmentServiceId}).Error
	if err != nil {
		return err
	}
	return nil
}
