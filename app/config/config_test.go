package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const dbName string = "name"

func TestNewConfig(t *testing.T) {
	c := NewConfig(".")

	require.Equal(t, c.Db.Name, dbName)
}
