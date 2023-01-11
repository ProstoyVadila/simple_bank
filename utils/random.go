package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(rand.Intn(26)) + strconv.Itoa(rand.Intn(999))
}

func RandomBalance() int64 {
	return RandomInt(0, 99999)
}

func RandomCurrency() string {
	c := [3]string{"PHP", "USD", "KZT"}
	return c[rand.Intn(len(c))]
}

func RandomEmail() string {
	return fmt.Sprintf("%v@email.com", RandomOwner())
}
