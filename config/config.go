package config

import (
	"log"
	"os"
)

var DecryptionKey string

func getEnvironmentVariable(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalln("missing environment variable:", key)
	}
	return value
}

func init() {
	DecryptionKey = getEnvironmentVariable("DECRYPTION_KEY")
}
