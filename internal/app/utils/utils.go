package utils

import (
	"crypto/md5"
	"fmt"
)

func Hash(url []byte, hostName string) string {
	hash := md5.Sum(url)
	return fmt.Sprintf("%s%s%x", hostName, "/", hash)
}
