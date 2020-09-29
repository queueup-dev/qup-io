package envvar

import (
	"os"
	"testing"
)

func TestMustWithExistingValue(t *testing.T) {
	err := os.Setenv("TEST", "Hello World!")
	if err != nil {
		t.Fail()
	}

	if Must("TEST") != "Hello World!" {
		t.Fail()
	}
}

func TestMustPanicWithNonExistingValue(t *testing.T) {
	_ = os.Unsetenv("TEST")
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Fail()
	}()

	Must("TEST")
}

func TestLookupEnvWithExistingValue(t *testing.T) {
	err := os.Setenv("TEST", "Hello World!")

	if err != nil {
		t.Fail()
	}

	val, _ := LookupEnv("TEST")

	if val != "Hello World!" {
		t.Fail()
	}
}

func TestLookupEnvWithNonExistingValue(t *testing.T) {
	_ = os.Unsetenv("TEST")
	val, errString := LookupEnv("TEST")

	if val != "" {
		t.Fail()
	}

	if *errString != "TEST" {
		t.Fail()
	}
}
