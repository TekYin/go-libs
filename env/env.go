package env

import "os"

func GetEnvOrEmpty(key string) string {
	return os.Getenv(key)
}

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Missing env variable: " + key)
	}
	return value
}
