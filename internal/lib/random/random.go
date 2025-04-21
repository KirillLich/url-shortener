package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
