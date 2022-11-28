package common

import (
	"log"
	"strings"

	"github.com/johnaoss/htpasswd/apr1"
)

func CheckPasswordHash(password, hash string) bool {
	substrings := strings.Split(hash, "$")

	if len(substrings) != 4 || substrings[0] != "" {
		log.Panic("Invalid hash provided, please use salt")
	}

	if substrings[1] != "apr1" {
		log.Panic("Only apr1 is supported for password hashing")
	}

	mhash, err := apr1.Hash(password, substrings[2])
	return err == nil && mhash == hash
}
