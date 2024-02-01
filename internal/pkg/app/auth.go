package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/role"
	"awesomeProject/internal/app/schemes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Register @Summary		Регистрация
// @Tags		Авторизация
// @Description	Регистрация нового пользователя
// @Accept		json
// @Param		user_credentials body schemes.RegisterReq true "login and password"
// @Success		200
// @Router		/api/user/sign_up [post]
func (a *Application) Register(c *gin.Context) {
	request := &schemes.RegisterReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingUser, err := a.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if existingUser != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := ds.User{
		Role:     role.Customer,
		Login:    request.Login,
		Password: generateHashString(request.Password),
	}
	if err := a.repo.AddUser(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// Login @Summary		Авторизация
// @Tags		Авторизация
// @Description	Авторизует пользователя по логиню, паролю и отдаёт jwt токен для дальнейших запросов
// @Accept		json
// @Produce		json
// @Param		user_credentials body schemes.LoginReq true "login and password"
// @Success		200 {object} schemes.AuthResp
// @Router		/api/user/login [post]
// @Consumes     json
func (a *Application) Login(c *gin.Context) {
	JWTConfig := a.config.JWT
	request := &schemes.LoginReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := a.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if user.Password != generateHashString(request.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	token := jwt.NewWithClaims(JWTConfig.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTConfig.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserUUID: user.UUID,
		Role:     user.Role,
		Login:    user.Login,
	})
	if token == nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
		return
	}

	strToken, err := token.SignedString([]byte(JWTConfig.Token))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
		return
	}

	userId, _ := c.Get("userId")
	log.Println(userId)

	log.Println(user.Role)

	c.JSON(http.StatusOK, schemes.AuthResp{
		AccessToken: strToken,
		TokenType:   "Bearer",
		Role:        int(user.Role),
	})
}

// Logout @Summary		Выйти из аккаунта
// @Tags		Авторизация
// @Description	Выход из аккаунта
// @Accept		json
// @Success		200
// @Router		/api/user/loguot [get]
func (a *Application) Logout(c *gin.Context) {
	jwtStr := c.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf(jwtStr))
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JWT.Token), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	err = a.redisClient.WriteJWTToBlacklist(c.Request.Context(), jwtStr, a.config.JWT.ExpiresIn)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
