package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
}

func HashUser(userName string) []byte {
	hash := sha256.New()
	hash.Write([]byte(userName))
	return hash.Sum(nil)
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func CheckErr(err error, text string) {
	if err != nil {
		log.Fatal(text, ": ", err)
	}
}
