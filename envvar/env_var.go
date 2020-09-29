package envvar

import (
	"log"
	"os"
)

func LookupEnv(s string) (r string, err *string) {
	val, ok := os.LookupEnv(s)
	if !ok {
		return "", &s
	}
	return val, nil
}

func Must(s string) string {
	val, ok := os.LookupEnv(s)

	if !ok {
		log.Printf("the required environment variable %s is missing.", s)
		panic("the required environment variable %s is missing.")
	}

	return val
}
