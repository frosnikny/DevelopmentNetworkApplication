package app

import (
	"awesomeProject/internal/app/schemes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) GetDevelopmentService(c *gin.Context) {
	var request schemes.DevelopmentServiceRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	developmentService, err := a.repo.GetDevelopmentServiceByID(request.DevelopmentServiceId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if developmentService == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("услуга по разработке не найдена"))
		return
	}
	c.JSON(http.StatusOK, developmentService)
}
