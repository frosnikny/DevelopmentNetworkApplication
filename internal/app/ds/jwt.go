package ds

import (
	"awesomeProject/internal/app/role"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID string    `json:"user_uuid"`
	Role     role.Role `json:"role"`
	Login    string    `json:"login"`
}
