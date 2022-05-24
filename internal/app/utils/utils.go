package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
)

func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
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
