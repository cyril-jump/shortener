package users

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
)

//DBUsers struct
type DBUsers struct {
	randNum []byte
	ctx     context.Context
}

//New DBUsers constructor
func New(ctx context.Context) *DBUsers {
	key := []byte("secret")

	return &DBUsers{
		randNum: key,
		ctx:     ctx,
	}
}

//CreateToken Create token for user
func (MU *DBUsers) CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": userID})
	tokenString, _ := token.SignedString(MU.randNum)
	return tokenString, nil
}

//CheckToken Check user token
func (MU *DBUsers) CheckToken(tokenString string) (string, bool) {

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
