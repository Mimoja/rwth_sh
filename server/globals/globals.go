package globals

import (
	"crypto/rand"
	"log"
)

func getRandomBytes(length uint) []byte {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

var Secret = getRandomBytes(64)

const Userkey = "user"

var Config AppConfig
