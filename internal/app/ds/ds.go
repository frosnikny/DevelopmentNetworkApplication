package ds

import "time"

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string `gorm:"size:30;not null"`
	Password  string `gorm:"size:30;not null"`
	Name      string `gorm:"size:50;not null"`
	Moderator bool   `gorm:"not null"`
}

type DevelopmentService struct {
	UUID         string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid"`
	Title        string  `gorm:"size:100"`
	Description  string  `gorm:"type:text"`
	ImageUrl     *string `gorm:"size:100" json:"image_url"`
	Price        uint    `gorm:"type:integer"`
	RecordStatus uint    `gorm:"type:integer"`
	Technology   string  `gorm:"type:text"`
	DetailedCost string  `gorm:"type:text"`
}

type CustomerRequest struct {
	UUID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RecordStatus      uint      `gorm:"type:integer"`
	CreationDate      time.Time `gorm:"type:timestamp"`
	FormationDate     time.Time `gorm:"type:timestamp"`
	CompletionDate    time.Time `gorm:"type:timestamp"`
	WorkSpecification string    `gorm:"type:text"`
	CreatorId         string    `gorm:"not null"`
	ModeratorId       *string   `json:"-"`

	Creator   User
	Moderator *User
}

type ServiceRequest struct {
	DevelopmentServiceId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"development_service_id"`
	CustomerRequestId    string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"customer_request_id"`
	WorkScope            string `gorm:"type:text"`
	WorkingDays          uint   `gorm:"type:integer"`

	DevelopmentService *DevelopmentService `gorm:"foreignKey:DevelopmentServiceId"`
	CustomerRequest    *CustomerRequest    `gorm:"foreignKey:CustomerRequestId"`
}

const (
	DSWorks = iota
	DSDeleted
)

const (
	CRDraft = iota
	CRWorks
	CRCompleted
	CRDeclined
	CRDeleted
)

const ZeroUUID = "00000000-0000-0000-0000-000000000000"
