package oteljsonlogflattenerprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	require.NoError(t, err)
}
