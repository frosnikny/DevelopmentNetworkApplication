package api

import (
	"awesomeProject/internal/api/ds"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//type Person struct {
//	Index       int
//	Title       string
//	Description string
//	ImageName   string
//	Price       int
//}

func search(query string, persons []ds.Person) []ds.Person {
	// Создаем результирующий массив
	results := []ds.Person{}

	// Проходим по массиву продуктов
	for _, person := range persons {
		// Проверяем, содержит ли продукт указанный запрос
		if strings.Contains(person.Title, query) {
			// Добавляем продукт в результирующий массив
			results = append(results, person)
		}
	}

	// Возвращаем результирующий массив
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
		// Получить переменную запроса page
		page := c.Param("page")

		// Считать значение переменной запроса
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

	// Добавляем обработчик событий для кнопки
	r.GET("/search", func(c *gin.Context) {
		// Получаем значение запроса из формы
		query := c.Query("query")

		// Выполняем поиск
		results := search(query, pipe)

		// Отображаем результаты поиска
		c.HTML(200, "index.gohtml", gin.H{
			"pipeline": results,
		})
	})

	r.Static("/images", "./resources")
	r.Static("/styles", "./templates/css")

	r.Run()

	log.Println("Sever down")
}
