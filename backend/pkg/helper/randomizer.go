package helper

import (
	"math/rand"
	"time"
)

func RandString(length int) (res string) {
	if length < 1 {
		length = 1
	}
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"
	for i := 0; i < length; i++ {
		res += string(charset[randSeed.Intn(len(charset))])
	}
	return
}
