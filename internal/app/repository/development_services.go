package repository

import (
	"awesomeProject/internal/app/ds"
	"errors"
	"gorm.io/gorm"
	"strings"
)

func (r *Repository) GetDevelopmentServiceByID(id string) (*ds.DevelopmentService, error) {
	developmentService := &ds.DevelopmentService{UUID: id}
	query := r.db.Table("development_services")
	query = query.Where("record_status = ?", ds.DSWorks)
	query = query.Where("uuid = ?", id)
	err := query.First(developmentService).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return developmentService, nil
}

func (r *Repository) AddDevelopmentService(developmentService *ds.DevelopmentService) error {
	err := r.db.Create(&developmentService).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveDevelopmentService(developmentService *ds.DevelopmentService) error {
	err := r.db.Save(developmentService).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetDevelopmentServicesByName(developmentServiceName string) ([]ds.DevelopmentService, error) {
	var developmentServices []ds.DevelopmentService

	err := r.db.Where("lower(title) ilike ? AND record_status = 0", "%"+strings.ToLower(developmentServiceName)+"%").Find(&developmentServices).Error
	if err != nil {
		return nil, err
	}

	return developmentServices, nil
}
