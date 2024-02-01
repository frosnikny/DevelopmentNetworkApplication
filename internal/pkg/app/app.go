package app

import (
	"awesomeProject/docs"
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/redis"
	"awesomeProject/internal/app/repository"
	"awesomeProject/internal/app/role"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	docs.SwaggerInfo.Title = "Development Services"
	docs.SwaggerInfo.Description = "API SERVER"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1:8080"
	docs.SwaggerInfo.BasePath = "/"

	r.Use(ErrorHandler())

	api := r.Group("/api")
	{
		// Услуги (Разработка)
		d := api.Group("/devs")
		{
			d.GET("", a.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), a.GetAllDevelopmentServices)                     // Список с поиском
			d.GET("/:development_service_id", a.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), a.GetDevelopmentService) // Одна услуга
			d.DELETE("/:development_service_id", a.WithAuthCheck(role.Moderator), a.DeleteDevelopmentService)                              // Удаление
			d.PUT("/:development_service_id", a.WithAuthCheck(role.Moderator), a.ChangeDevelopmentService)                                 // Изменение
			d.POST("", a.WithAuthCheck(role.Moderator), a.AddDevelopmentService)                                                           // Добавление
			d.POST("/:development_service_id/add_to_request", a.WithAuthCheck(role.Customer, role.Moderator), a.AddToCustomerRequest)      // Добавление в заявку
		}

		// Заявки (Заказы)
		req := api.Group("/requests")
		{

			req.GET("", a.WithAuthCheck(role.Customer, role.Moderator), a.GetAllCustomerRequests)                                                         // Список (отфильтровать по дате формирования и статусу)
			req.GET("/:customer_request_id", a.WithAuthCheck(role.Customer, role.Moderator), a.GetCustomerRequest)                                        // Одна заявка
			req.PUT("/:customer_request_id/update", a.WithAuthCheck(role.Customer, role.Moderator), a.UpdateCustomerRequest)                              // Изменение (добавление спецификации)
			req.PUT("/:customer_request_id/change_scope/:development_service_id", a.WithAuthCheck(role.Customer, role.Moderator), a.UpdateServiceRequest) // Изменение (добавление спецификации)
			req.DELETE("/:customer_request_id", a.WithAuthCheck(role.Customer, role.Moderator), a.DeleteCustomerRequest)                                  // Удаление
			req.DELETE("/:customer_request_id/delete_development_service/:development_service_id", a.WithAuthCheck(role.Customer, role.Moderator),
				a.DeleteFromCustomerRequest) // Изменение (удаление услуг)
			req.PUT("/:customer_request_id/user_confirm", a.WithAuthCheck(role.Customer, role.Moderator), a.UserConfirm) // Сформировать создателем
			req.PUT("/:customer_request_id/moderator_confirm", a.WithAuthCheck(role.Moderator), a.ModeratorConfirm)      // Сформировать модератором
			req.PUT("/:customer_request_id/payment", a.Payment)
		}

		// Пользователи (авторизация)
		u := api.Group("/user")
		{
			u.POST("/sign_up", a.Register)
			u.POST("/login", a.Login)
			u.GET("/logout", a.Logout)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	a.minioClient, err = minio.New(a.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		log.Println("point: ", a.config.Minio.Endpoint)
		panic(err)
	}

	a.redisClient, err = redis.New(a.config.Redis)
	if err != nil {
		panic(err)
	}

	return a
}
