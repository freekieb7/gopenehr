package util

import "math/rand"

func RandomLowerAlphaString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {

		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
