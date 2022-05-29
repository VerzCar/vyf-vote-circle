package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRawJsonError_Error(t *testing.T) {
	expectedErrorMsg := "Error"
	jsonError := RawJsonError{
		RawMessage: []byte(expectedErrorMsg),
	}

	require.Equal(t, expectedErrorMsg, jsonError.Error())
}
