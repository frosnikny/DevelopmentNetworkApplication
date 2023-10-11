package app

import (
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

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

	r.LoadHTMLGlob("templates/html/*")

	r.GET("/development/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil {
			return
		}

		developmentService, err := a.repo.GetDevelopmentServiceByID(uint(number))
		if err != nil {
			log.Printf("cant get developmentService by id %v", err)
			developmentService, err = a.repo.GetDevelopmentServiceByID(0)
			if err != nil {
				c.Error(err)
				return
			}
		}

		c.HTML(http.StatusOK, "full_service_card.gohtml", gin.H{
			"Title":       developmentService.Title,
			"Description": developmentService.Description,
			"ImageName":   developmentService.ImageName,
			"Price":       developmentService.Price,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		searchDevelopmentServiceName := c.Query("developmentServiceName")
		deleteDevelopmentServiceID := c.Query("deleteDevelopmentServiceID")
		if len(searchDevelopmentServiceName) == 0 {
			if len(deleteDevelopmentServiceID) != 0 {
				number, err := strconv.Atoi(deleteDevelopmentServiceID)
				if err != nil {
					return
				}
				err = a.repo.DeleteDevelopmentServiceByID(uint(number))
				if err != nil {
					return
				}
			}
			developmentServices, err := a.repo.GetDevelopmentServices()
			if err != nil {
				c.Error(err)
				return
			}
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"developmentServices": developmentServices,
			})
		} else {
			results, err := a.repo.FindDevelopmentServiceByName(searchDevelopmentServiceName)
			if err != nil {
				c.Error(err)
				return
			}

			c.HTML(200, "index.gohtml", gin.H{
				"developmentServices": results,
				"searchName":          searchDevelopmentServiceName,
			})
		}
	})

	r.GET("/deleted", func(c *gin.Context) {
		deletedDevelopmentServices, err := a.repo.GetDeletedDevelopmentServices()
		if err != nil {
			c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"developmentServices": deletedDevelopmentServices,
		})

	})

	r.Static("/images", "./resources")
	r.Static("/styles", "./templates/css")

	err := r.Run()
	if err != nil {
		return
	}

	log.Println("Server down")
}

func New() *Application {
	a := &Application{}

	log.Println("dsn: " + dsn.FromEnv())
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		panic(err)
	}
	a.repo = repo

	return a
}
