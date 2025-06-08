package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("contraseña_secreta_para_este_ejemplo")

func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": userID,
        "exp":    time.Now().Add(time.Minute * 5).Unix(), 
    })

    return token.SignedString(jwtSecret)
}