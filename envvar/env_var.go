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
		log.Fatalf("the required environment variable %s is missing.", s)
	}

	return val
}

func FilterValuesAndDereference(arr ...*string) ([]string, int) {
	n := 0
	b := make([]string, len(arr))
	for _, val := range arr {
		if val != nil {
			b[n] = *val
			n++
		}
	}

	return b[:n], n
}
