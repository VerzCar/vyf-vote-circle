package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const environment string = "development"

func TestNewConfig(t *testing.T) {
	c := NewConfig(".")

	require.Equal(t, c.Environment, EnvironmentDev)
}
