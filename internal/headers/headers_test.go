package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	v, ok := headers.Get("host")
	assert.Equal(t, "localhost:42069", v)
	assert.True(t, ok)
	assert.Equal(t, len(data), n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("       Host:localhost:42069       \r\nHost:localhost:42019\r\n\r\n")
	n, done, err = headers.Parse(data)
	v, ok = headers.Get("host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069, localhost:42019", v)
	require.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.True(t, done)

	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
