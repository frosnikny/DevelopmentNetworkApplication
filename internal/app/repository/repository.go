package repository

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"

	"awesomeProject/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetDevelopmentServices() (*[]ds.DevelopmentService, error) {
	developmentServices := &[]ds.DevelopmentService{}

	err := r.db.Where("record_status = ?", 0).Find(developmentServices).Error
	if err != nil {
		return nil, err
	}

	return developmentServices, nil
}

func (r *Repository) GetDeletedDevelopmentServices() (*[]ds.DevelopmentService, error) {
	deletedDevelopmentServices := &[]ds.DevelopmentService{}

	err := r.db.Where("record_status = ?", 1).Find(deletedDevelopmentServices).Error
	if err != nil {
		return nil, err
	}

	return deletedDevelopmentServices, nil
}

func (r *Repository) GetDevelopmentServiceByID(id uint) (*ds.DevelopmentService, error) {
	developmentService := &ds.DevelopmentService{}

	err := r.db.First(developmentService, "id = ?", strconv.Itoa(int(id))).Error
	if err != nil {
		return nil, err
	}

	return developmentService, nil
}

func (r *Repository) FindDevelopmentServiceByName(name string) (*[]ds.DevelopmentService, error) {
	developmentServices := &[]ds.DevelopmentService{}

	err := r.db.Where("lower(title) ilike ?", "%"+name+"%").Find(developmentServices).Error
	if err != nil {
		return nil, err
	}

	return developmentServices, nil
}

func (r *Repository) DeleteDevelopmentServiceByID(id uint) error {
	developmentService := &ds.DevelopmentService{}

	r.db.Exec("UPDATE development_services SET record_status = 1 WHERE id = @id", sql.Named("id", id))
	err := r.db.First(developmentService, "id = ?", strconv.Itoa(int(id))).Error
	return err
}

func (r *Repository) CreateDevelopmentService(developmentService ds.DevelopmentService) error {
	return r.db.Create(developmentService).Error
}
