package api

import (
	"awesomeProject/internal/api/ds"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func search(query string, persons []ds.Person) []ds.Person {
	var results []ds.Person

	for _, person := range persons {
		if strings.Contains(person.Title, query) {
			results = append(results, person)
		}
	}

	return results
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/html/*")

	pipe := ds.GetPipeline()

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"pipeline": pipe,
		})
	})

	r.GET("/full/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil || number > len(pipe) {
			number = 0
		}

		c.HTML(http.StatusOK, "full-card.gohtml", gin.H{
			"Title":       pipe[number].Title,
			"Description": pipe[number].Description,
			"ImageName":   pipe[number].ImageName,
		})
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("query")

		results := search(query, pipe)

		c.HTML(200, "index.gohtml", gin.H{
			"pipeline": results,
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
