package main

import (
	"math/rand"
)

const randomWord = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) (value string) {
	valueBytes := make([]byte, length)
	wordLength := len(randomWord)

	for i := 0; i < length; i++ {
		valueBytes[i] = randomWord[rand.Intn(wordLength)]
	}

	value = string(valueBytes)
	return
}
