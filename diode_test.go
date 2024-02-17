package main

import (
	"fmt"
	"testing"

	"github.com/acep-uaf/data-diode/utility"
)

func TestCLI(t *testing.T) {
	got := "diode"
	want := "diode"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestConfiguration(t *testing.T) {
	got := Configuration{}
	want := Configuration{}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFileContents(t *testing.T) {
	got := fmt.Sprintf("%x", utility.Checksum())
	want := "ed03bb5d7385010c645c2c72ceabea3b15806db757005071309745c59933586f"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
