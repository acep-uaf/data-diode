package example

import (
	"regexp"
	"testing"
)

// Name + Valid Return Value

func TestExampleName(t *testing.T) {
	name := "Solomon"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Communication("Solomon")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Communication("Solomon) = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// Empty String -- Error Checking
func TestExampleEmpty(t *testing.T) {
	msg, err := Communication("")
	if msg != "" || err == nil {
		t.Fatalf(`Communication("") = %q, %v, want "", error`, msg, err)
	}
}
