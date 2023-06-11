package util

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var currencies = [3]string{"USD", "INR", "EUR"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomFloat generates a random floating value between min and max
func RandomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min+1)
}

// RandomUUID generates a random UUID
func RandomUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Print("[ERROR] error generating random UUID, ", err)
		return ""
	}
	return uuid.String()
}

func RandomCurrency() string {
	return currencies[RandomInt(0, int64(len(currencies)-1))]
}
