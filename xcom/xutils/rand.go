package xutils

import "math/rand"

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(bytes)
}
