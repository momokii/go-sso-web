package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateUUIDV4() (string, error) {
	// create 16 byte array for UUID
	uuid := make([]byte, 16)

	// fill array with random data
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// config UUID to v4
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4 (0100)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 1 (10)

	// Format uuid to string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
