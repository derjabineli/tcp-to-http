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
	assert.Equal(t, "localhost:42069", headers["host"])
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
  assert.Equal(t, "application/json", headers["content-type"])
  assert.False(t, done)
 
  // Test: Valid single header with underscore 
  headers = NewHeaders()
  data = []byte("Accept_Language: de   \r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "de", headers["accept_language"])
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
  assert.Equal(t, "application/json", headers["content-type"])
  
  data = data[n:]
  _, done, _ = headers.Parse(data)
  require.NoError(t, err)
  assert.Equal(t, done, true)

  // Test: Invalid header field name
  headers = NewHeaders()
  data = []byte("h™ost: localhost:41209\r\n\r\n")
  _, _, err = headers.Parse(data)
  require.Error(t, err)

  // Test: Valid and invalid field names
  headers = NewHeaders()
  data = []byte("host: localhost:41209\r\n Content-T¢pe: application/json\r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  assert.Equal(t, "localhost:41209", headers["host"])

  data = data[n:]
  n, done, err = headers.Parse(data)
  require.Error(t, err)

  // Test: Multiple valid values for one field name
  headers = NewHeaders()
  data = []byte("Set-Person: Eli\r\n Set-Person: Vika\r\n\r\n")
  
  for {
    n, done, err := headers.Parse(data)
    require.NoError(t, err)

    data = data[n:]
    if done {
      break
      }
    }

  assert.Equal(t, "Eli, Vika", headers["set-person"])

  // Test: Valid header with mutiple field names and combined values
  headers = NewHeaders()
  data = []byte("Host: localhost:41209\r\n Set-Person: Eli\r\n Set-Person: Vika\r\n\r\n")
  
  for {
    n, done, err := headers.Parse(data)
    require.NoError(t, err)

    data = data[n:]
    if done {
      break 
    }
  }

  assert.Equal(t, "localhost:41209", headers["host"])
  assert.Equal(t, "Eli, Vika", headers["set-person"])
}
