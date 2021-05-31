package goutils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
)

// GenerateRandomKey generates a random secure key
func GenerateRandomKey(l int) []byte {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

// IsVaildAESkeylen returns true if len is a valid aes keylength
func IsVaildAESkeylen(len int) bool {
	switch len {
	case 16, 24, 32:
		return true
	}
	return false
}

// EncodeBase64 encodes the given bytes in Base64
func EncodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

// DecodeBase64 decodes given bytes from Base64
func DecodeBase64(b []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		fmt.Println("Error: Bad Key!")
		os.Exit(1)
	}
	return data
}
