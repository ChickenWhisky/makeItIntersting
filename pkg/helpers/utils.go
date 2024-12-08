package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strconv"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
func ConvertStringToInt(s string) int {
	lengthFromEnv, err := strconv.Atoi(os.Getenv("CONTRACT_ID_LENGTH"))
	if err != nil {
		log.Println("Error converting CONTRACT_ID_LENGTH to int")
	}
	return lengthFromEnv
}
