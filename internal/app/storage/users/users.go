package users

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	"log"
)

type DBUsers struct {
	storageUsers map[string]ModelUser
	randNum      []byte
	CookieWord   string
}

type ModelUser struct {
	userID string
	cookie string
}

func New() *DBUsers {
	key, err := utils.GenerateRandom(16)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	return &DBUsers{
		storageUsers: map[string]ModelUser{},
		randNum:      key,
		CookieWord:   "cookie",
	}
}

func (MU *DBUsers) GetUserID(userName string) (string, error) {
	if _, ok := MU.storageUsers[userName]; !ok {
		return "", errs.ErrNoContent
	}
	return MU.storageUsers[userName].userID, nil
}

func (MU *DBUsers) SetUserID(userName string) {
	var model ModelUser
	hash := utils.HashUser(userName)
	model.userID = hex.EncodeToString(hash)
	MU.storageUsers[userName] = model
}

func (MU *DBUsers) CreateCookie(userID string) (string, error) {
	code := MU.randNum
	id, err := hex.DecodeString(userID)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, code)
	h.Write(id)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (MU *DBUsers) CheckCookie(cookieOld, userID string) bool {
	oldCookie, _ := hex.DecodeString(cookieOld)
	code := MU.randNum
	mac := hmac.New(sha256.New, code)
	id, _ := hex.DecodeString(userID)
	mac.Write(id)
	newCookie := mac.Sum(nil)
	return hmac.Equal(newCookie, oldCookie)

}
