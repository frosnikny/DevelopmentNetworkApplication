package api

import (
	"awesomeProject/internal/api/ds"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func search(query string, services []ds.Service) []ds.Service {
	var results []ds.Service

	for _, service := range services {
		if strings.Contains(service.Title, query) {
			results = append(results, service)
		}
	}

	return results
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/html/*")

	services := ds.GetServices()

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"services": services,
		})
	})

	r.GET("/service/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil || number > len(services) {
			number = 0
		}

		c.HTML(http.StatusOK, "full_service_card.gohtml", gin.H{
			"Title":       services[number].Title,
			"Description": services[number].Description,
			"ImageName":   services[number].ImageName,
		})
	})

	r.GET("/search", func(c *gin.Context) {
		searchServiceName := c.Query("serviceName")

		results := search(searchServiceName, services)

		c.HTML(200, "index.gohtml", gin.H{
			"services": results,
		})
	})

	r.Static("/images", "./resources")
	r.Static("/styles", "./templates/css")

	err := r.Run()
	if err != nil {
		log.Println("Server start up error", err)
		return
	}

	log.Println("Sever down")
}
