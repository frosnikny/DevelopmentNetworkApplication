package api

import (
	"awesomeProject/internal/api/ds"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func search(query string, developmentServices []ds.DevelopmentService) []ds.DevelopmentService {
	var results []ds.DevelopmentService

	query = strings.ToLower(query)

	for _, developmentService := range developmentServices {
		if strings.Contains(strings.ToLower(developmentService.Title), query) {
			results = append(results, developmentService)
		}
	}

	return results
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/html/*")

	developmentServices := ds.GetDevelopmentServices()

	r.GET("/service/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil || number > len(developmentServices) {
			number = 0
		}

		c.HTML(http.StatusOK, "full_service_card.gohtml", gin.H{
			"Title":       developmentServices[number].Title,
			"Description": developmentServices[number].Description,
			"ImageName":   developmentServices[number].ImageName,
			"Price":       developmentServices[number].Price,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		searchDevelopmentServiceName := c.Query("developmentServiceName")

		results := search(searchDevelopmentServiceName, developmentServices)

		if len(results) == 0 {
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"developmentServices": developmentServices,
			})
		} else {

			c.HTML(200, "index.gohtml", gin.H{
				"developmentServices": results,
			})
		}
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
