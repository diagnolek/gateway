package helpers

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	valueInt, err := strconv.Atoi(value)
	if value == "" || err != nil {
		return defaultValue
	}

	return valueInt
}
