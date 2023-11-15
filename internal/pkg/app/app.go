package app

import (
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/repository"
	"github.com/gin-gonic/gin"
	"log"
)

type Application struct {
	//config *config.Config
	//router *http.ServeMux
	repo *repository.Repository
}

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/api/containers/:container_id", a.GetDevelopmentService)

	log.Println("Server down")
}

func New() *Application {
	a := &Application{}

	//log.Println("dsn: " + dsn.FromEnv())
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		panic(err)
	}
	a.repo = repo

	return a
}
