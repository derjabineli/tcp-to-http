package response

import (
	"fmt"
)

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	statusLine := []byte(fmt.Sprintf("HTTP/1.1 %v %v \r\n", statusCode, reasonPhrase))
	return statusLine
}