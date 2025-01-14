package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
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

func HashText(text string) string {
	// Create a new SHA-256 hash
	hash := sha256.New()
	// Write the input text to the hash
	hash.Write([]byte(text))
	// Get the finalized hash as a byte slice
	hashedBytes := hash.Sum(nil)
	// Convert the byte slice to a hexadecimal string
	return hex.EncodeToString(hashedBytes)
}

// isAlphanumeric checks if a string contains only alphanumeric characters
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// The following code is for checking if the subevent Names are of the correct naming syntax
// validateString checks if a string ends with _YES or _NO and has valid text before it
func ValidateSubEventName(input string) error {
	// Ensure input ends with _YES or _NO
	if strings.HasSuffix(input, "_YES") || strings.HasSuffix(input, "_NO") {
		// Extract the prefix (the part before _YES or _NO)
		prefix := input[:len(input)-4]

		// Ensure the prefix is non-empty and consists of valid characters
		if len(prefix) > 0 && isAlphanumeric(prefix) {
			return nil
		}
	}
	return errors.New("invalid subevent name : doesnt SUFFIX with _YES or _NO")
}
