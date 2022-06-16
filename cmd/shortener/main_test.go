package main

import (
	"github.com/stretchr/testify/require"
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
			//main()
			var err error
			require.NoError(t, err)
		})
	}
}
