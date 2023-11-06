package ds

import "time"

type DevelopmentServiceStatus struct {
	ID         uint
	StatusName string
}

type CustomerRequestStatus struct {
	ID         uint
	StatusName string
}

type User struct {
	UserId    uint   `gorm:"primaryKey;not null"`
	Name      string `gorm:"type:varchar(30)"`
	Username  string `gorm:"type:varchar(30)"`
	Password  string `gorm:"type:varchar(30)"`
	Moderator bool   `gorm:"type:bool"`
}

type DevelopmentService struct {
	DevelopmentServiceId uint   `gorm:"primaryKey;not null"`
	Title                string `gorm:"type:varchar(100)"`
	Description          string `gorm:"type:text"`
	ImageName            string `gorm:"type:varchar(100)"`
	Price                uint   `gorm:"type:integer"`
	RecordStatus         uint   `gorm:"type:integer"`
	Technology           string `gorm:"type:text"`
	DetailedCost         string `gorm:"type:text"`
}

type CustomerRequest struct {
	CustomerRequestId uint      `gorm:"primaryKey;not null"`
	RecordStatus      uint      `gorm:"type:integer"`
	CreationDate      time.Time `gorm:"type:date"`
	FormationDate     time.Time `gorm:"type:date"`
	CompletionDate    time.Time `gorm:"type:date"`
	WorkSpecification string    `gorm:"type:text"`
	CreatorId         uint      `gorm:"not null"`
	ModeratorId       uint      `gorm:"not null"`
	Creator           *User     `gorm:"foreignKey:CreatorId"`
	Moderator         *User     `gorm:"foreignKey:ModeratorId"`
}

type ServiceRequest struct {
	DevelopmentServiceId uint                `gorm:"primaryKey;not null;autoIncrement:false"`
	CustomerRequestId    uint                `gorm:"primaryKey;not null;autoIncrement:false"`
	DevelopmentService   *DevelopmentService `gorm:"foreignKey:DevelopmentServiceId"`
	CustomerRequest      *CustomerRequest    `gorm:"foreignKey:CustomerRequestId"`
	WorkScope            string              `gorm:"type:text"`
	WorkingDays          uint                `gorm:"type:integer"`
}

func GetDevelopmentServiceStatuses() []DevelopmentServiceStatus {
	return []DevelopmentServiceStatus{
		{0, "Works"},
		{1, "Deleted"},
	}
}

func GetCustomerRequestStatus() []CustomerRequestStatus {
	return []CustomerRequestStatus{
		{0, "Draft"},
		{1, "Works"},
		{2, "Completed"},
		{3, "Declined"},
		{4, "Deleted"},
	}
}
