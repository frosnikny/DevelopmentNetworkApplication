package ds

type DevelopmentServiceStatus struct {
	ID         uint `gorm:"primarykey"`
	StatusName string
}

type DevelopmentService struct {
	ID           uint `gorm:"primarykey"`
	Title        string
	Description  string
	ImageName    string
	Price        uint
	RecordStatus uint
}

func GetDevelopmentServiceStatuses() []DevelopmentServiceStatus {
	return []DevelopmentServiceStatus{
		{0, "Works"},
		{1, "Deleted"},
	}
}
