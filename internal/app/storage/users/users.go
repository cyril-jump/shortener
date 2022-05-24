package users

import (
	"fmt"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/golang-jwt/jwt"
	"log"
)

type DBUsers struct {
	randNum []byte
}

func New() *DBUsers {
	key, err := utils.GenerateRandom(16)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return &DBUsers{
		randNum: key,
	}
}

func (MU *DBUsers) CreateCookie(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": userID})
	tokenString, _ := token.SignedString(MU.randNum)
	return tokenString, nil
}

func (MU *DBUsers) CheckCookie(tokenString string) (string, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexected signing method: %v", token.Header["alg"])
		}
		return MU.randNum, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%s", claims["user"]), ok
	}
	return "", false
}
