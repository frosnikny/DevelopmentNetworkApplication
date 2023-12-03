package app

import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	// Услуги (Разработка)
	r.GET("/api/devs", a.GetAllDevelopmentServices)                                    // Список с поиском
	r.GET("/api/devs/:development_service_id", a.GetDevelopmentService)                // Одна услуга
	r.DELETE("/api/devs/:development_service_id", a.DeleteDevelopmentService)          // Удаление
	r.PUT("/api/devs/:development_service_id", a.ChangeDevelopmentService)             // Изменение
	r.POST("/api/devs", a.AddDevelopmentService)                                       // Добавление
	r.POST("/api/devs/:development_service_id/add_to_request", a.AddToCustomerRequest) // Добавление в заявку

	// Заявки (Заказы)
	r.GET("/api/requests", a.GetAllCustomerRequests)                            // Список (отфильтровать по дате формирования и статусу)
	r.GET("/api/requests/:customer_request_id", a.GetCustomerRequest)           // Одна заявка
	r.PUT("/api/requests/:customer_request_id/update", a.UpdateCustomerRequest) // Изменение (добавление спецификации)
	r.DELETE("/api/requests/:customer_request_id", a.DeleteCustomerRequest)     // Удаление
	r.DELETE("/api/requests/:customer_request_id/delete_development_service/:development_service_id",
		a.DeleteFromCustomerRequest) // Изменение (удаление услуг)
	r.PUT("/api/requests/:customer_request_id/user_confirm", a.UserConfirm)           // Сформировать создателем
	r.PUT("/api/requests/:customer_request_id/moderator_confirm", a.ModeratorConfirm) // Сформировать модератором

	err := r.Run(fmt.Sprintf("%s:%d", a.config.ServiceHost, a.config.ServicePort))
	if err != nil {
		panic(err)
	}

	log.Println("Server down")
}

func New() *Application {
	var err error

	a := &Application{}
	a.config, err = config.NewConfig()
	if err != nil {
		panic(err)
	}

	a.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		panic(err)
	}

	a.minioClient, err = minio.New(a.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	return a
}
