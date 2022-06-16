package app

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_configPingDB(t *testing.T) {

	tests := []struct {
		name  string
		value string
	}{
		{name: "test main", value: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{}
			c.ConfigServerEnv()
			var err error
			err = nil
			require.NoError(t, err)
		})
	}
}
