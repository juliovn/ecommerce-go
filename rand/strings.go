package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// RememberToken is a helper function designed to generate remember tokens of predetermined size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// String will generate a byte slice of size nBytes and
// return a string that is the base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Bytes will help to generate n random bytes, or will return an error if there was one
// This uses crypto/rand package so safe to use with things like remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

