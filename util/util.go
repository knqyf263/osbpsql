package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var rsLowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")

// RandLowerString returns random string
func RandLowerString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rsLowerLetters[rand.Intn(len(rsLowerLetters))]
	}
	return string(b)
}
