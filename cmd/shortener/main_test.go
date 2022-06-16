package main

import (
	"testing"
)

func TestMain_ServerMain(t *testing.T) {

	tests := []struct {
		name  string
		value string
	}{
		{name: "test main", value: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
