package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz" // Characters to be used in random string generator

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt return a value n
// min <= n <= max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // Int63n -> 0...(max-min)
}

// RandomString return a random string with length of n
func RandomString(n int) string {
	var sb strings.Builder

	alphabetSize := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetSize)]

		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"BRL", "USD", "EUR"}

	return currencies[rand.Intn(len(currencies))]
}
