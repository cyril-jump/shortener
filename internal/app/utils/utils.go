package utils

import (
	"crypto/md5"
	"fmt"
	"log"
)

func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
}

func CheckErr(err error, text string) {
	if err != nil {
		log.Fatal(text, ": ", err)
	}
}
