package headers 
 
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderPar(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

  // Test: Valid single header with hyphen
  headers = NewHeaders()
  data = []byte("Content-Type: application/json       \r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "application/json", headers["Content-Type"])
  assert.False(t, done)
 
  // Test: Valid single header with underscore 
  headers = NewHeaders()
  data = []byte("Accept_Language: de   \r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "de", headers["Accept_Language"])
  assert.False(t, done)

  // Test: Valid multiple headers
  headers = NewHeaders()
  data = []byte("host: localhost:41209\r\n Content-Type: application/json\r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "localhost:41209", headers["host"])

  // Test: Valid done
   headers = NewHeaders()
  data = []byte("\r\n")  
  n, done, err = headers.Parse(data)
  assert.Equal(t, done, true)

  // Test: Valid multiple header read
  headers = NewHeaders()
  data = []byte("host: localhost:41209\r\n Content-Type: application/json\r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "localhost:41209", headers["host"])
  
  data = data[n:]
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  assert.Equal(t, "application/json", headers["Content-Type"])
  
  data = data[n:]
  _, done, _ = headers.Parse(data)
  require.NoError(t, err)
  assert.Equal(t, done, true)
}
