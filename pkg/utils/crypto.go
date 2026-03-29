package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateSecureToken(uuid string) (string, error) {

	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	data := append([]byte(uuid), randomBytes...)
	return hex.EncodeToString(data), nil
}

func HashSecureToken(rawToken string) string {

	sum := sha512.Sum512([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func GenerateOTP() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
