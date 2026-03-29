package utils

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordWithHash(pass string) (string, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHashed), err
}

// func CompareHashedPassword(passwordHashed string, pass string) error {
// 	compareHash := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(pass))
// 	return compareHash
// }


func CompareHashedPassword(passwordHashed string, pass string) error {
	start := time.Now()

	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(pass))

	log.Println("bcrypt compare took:", time.Since(start))

	return err
}