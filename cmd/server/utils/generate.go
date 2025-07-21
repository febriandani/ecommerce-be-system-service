package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	mr "math/rand"

	"golang.org/x/crypto/bcrypt"
)

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Hash256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}

func GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(savedPass, incomingPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(incomingPass))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func GenerateOTP() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, 4)
	n, err := io.ReadAtLeast(rand.Reader, b, 4)
	if n != 4 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func CreatePassword(length int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)

	for i := range b {
		b[i] = letterBytes[mr.Intn(len(letterBytes))]
	}
	return string(b)
}
