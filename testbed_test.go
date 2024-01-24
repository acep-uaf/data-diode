package main

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	name := "Input"
	t.Run(name, func(t *testing.T) {
		t.Parallel()
	})
}

func TestNewServer(t *testing.T) {
	name := "Output"
	t.Run(name, func(t *testing.T) {
		t.Parallel()
	})
}
