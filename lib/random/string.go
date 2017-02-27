package random

import (
	"crypto/rand"
	"log"
)

const availableCharBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// secureRandomBytes returns the requested number of bytes using crypto/rand
func secureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

// SecureString returns a string of the requested length,
// made from the byte characters provided (only ASCII allowed).
// Uses crypto/rand for security. Will panic if len(availableCharBytes) > 256.
//
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func SecureString(length int) string {

	// Compute bitMask
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 || availableCharLength > 256 {
		panic("availableCharBytes length must be greater than 0 and less than or equal to 256")
	}
	var bitLength byte
	var bitMask byte
	for bits := availableCharLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1

	// Compute bufferSize
	bufferSize := length + length/3

	// Create random string
	result := make([]byte, length)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			// Random byte buffer is empty, get a new one
			randomBytes = secureRandomBytes(bufferSize)
		}
		// Mask bytes to get an index into the character slice
		if idx := int(randomBytes[j%length] & bitMask); idx < availableCharLength {
			result[i] = availableCharBytes[idx]
			i++
		}
	}

	return string(result)
}
