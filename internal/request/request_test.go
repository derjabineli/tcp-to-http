package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
data            string
numBytesPerRead int
pos             int
}

// Read implementaion that reads chunked data to simulate reading from network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
  if cr.pos >= len(cr.data) {
    return 0, io.EOF
  }
  endIndex := cr.pos + cr.numBytesPerRead
  
  if endIndex > len(cr.data) {
    endIndex = len(cr.data)
  }
  
  n = copy(p, cr.data[cr.pos:endIndex])
  cr.pos += n
  if n > cr.numBytesPerRead {
    n = cr.numBytesPerRead
    cr.pos -= n - cr.numBytesPerRead
  }
  
  return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

  // Test: Good GET Request line with path
	reader = &chunkReader{
	data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 1024,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = &chunkReader{
	data: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 1024,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// Test: Too many parts in Request line
	reader = &chunkReader{
  data: "GET /coffee HTTP/1.1 POST\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 1024,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestHeadersParse(t *testing.T) {
	// Test: Out of order Method in Request line
	reader := &chunkReader{
		data: "/coffee HTTP/1.1 GET\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1024,
		}
		_, err := RequestFromReader(reader)
		require.Error(t, err)
	
		// Test: Invalid version in Request line
		reader = &chunkReader{
		data: "GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 40,
		}
		_, err = RequestFromReader(reader)
		require.Error(t, err)
	
		// Test: Good GET Request line with larger read size
		reader = &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 8,
		}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	
	  // Test: Standard Headers
	  reader = &chunkReader{
		  data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		  numBytesPerRead: 1024,
	  }
	  r, err = RequestFromReader(reader)
	  require.NoError(t, err)
	  require.NotNil(t, r)
	  assert.Equal(t, "localhost:42069", r.Headers["host"])
	  assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	  assert.Equal(t, "*/*", r.Headers["accept"])
	
	  // Test: Standard Headers with slow read speed
	  reader = &chunkReader{
		  data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		  numBytesPerRead: 1,
	  }
	  r, err = RequestFromReader(reader)
	  require.NoError(t, err)
	  require.NotNil(t, r)
	  assert.Equal(t, "localhost:42069", r.Headers["host"])
	  assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	  assert.Equal(t, "*/*", r.Headers["accept"])
	
	  // Test: Empty Headers
	  reader = &chunkReader{
		  data: "GET / HTTP/1.1\r\n\r\n",
		  numBytesPerRead: 8,
	  }
	  r, err = RequestFromReader(reader)
	  require.NoError(t, err)
	  require.NotNil(t, r)
	  require.NotNil(t, r.RequestLine)
	  assert.Empty(t, r.Headers)
	
	  // Test: Malformed Headers
	  reader = &chunkReader{
		  data: "GET / HTTP/1.1\r\nHost: localhost:42068\r\nUser-Agent : curl/7.81.0\r\nAccept: */*\r\n\r\n",
		  numBytesPerRead: 1,
	  }
	  _, err = RequestFromReader(reader)
	  require.Error(t, err)
	
	  // Test: Duplicate Header request
	  reader = &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nHost: anotherHost:8080\r\nAccept: */*\r\n\r\n",
		  numBytesPerRead: 8,
	  }
	  r, err = RequestFromReader(reader)
	  require.NoError(t, err)
	  require.Equal(t, "localhost:42069, anotherHost:8080", r.Headers["host"])
	  require.Equal(t, "*/*", r.Headers["accept"])
	
	  // Test: Case insensitive headers 
	  reader = &chunkReader{
		data: "GET / HTTP/1.1\r\nSet-Person: eli\r\nset-PeRsoN: vika\r\naCcEpT: */*\r\n\r\n",
		  numBytesPerRead: 8,
	  }
	  r, err = RequestFromReader(reader)
	  require.NoError(t, err)
	  require.Equal(t, "eli, vika", r.Headers["set-person"])
	  require.Equal(t, "*/*", r.Headers["accept"])
	
	  // Test: Missing crlf at end of Headers
	  reader = &chunkReader{
		data: "GET / HTTP/1.1\r\nSet-Person: eli\r\nset-PeRsoN: vika\r\naCcEpT: */*\r\n",
		  numBytesPerRead: 8,
	  }
	  _, err = RequestFromReader(reader)
	  require.Error(t, err)
	  require.Equal(t, "incomplete request", err.Error())
}