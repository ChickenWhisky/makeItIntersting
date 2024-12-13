package helpers

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"strconv"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
func ConvertStringToInt(s string) int {
	lengthFromEnv, err := strconv.Atoi(os.Getenv("CONTRACT_ID_LENGTH"))
	if err != nil {
		log.Println("Error converting CONTRACT_ID_LENGTH to int")
	}
	return lengthFromEnv
}
