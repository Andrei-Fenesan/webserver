package httpentity

import (
	"fmt"
	"strconv"
)

type Response struct {
	ResponseCode uint16
	Version      string
	Body         []byte
	Headers      map[string]string
}

func (response *Response) Encode() []byte {
	data := make([]byte, 0, 1024)
	data = appendResponseCodeText(data, response.Version, response.ResponseCode)
	data = append(data, encodeHeaders(response.Headers)...)
	data = append(data, response.Body...)
	return data
}

func appendResponseCodeText(data []byte, version string, responseCode uint16) []byte {
	return fmt.Appendf(data, "%s %s %s\r\n",
		version,
		strconv.FormatUint(uint64(responseCode), 10),
		getResponseMessage(responseCode))
}

func getResponseMessage(responseCode uint16) string {
	switch responseCode {
	case 200:
		return "OK"
	case 400:
		return "BAD REQUEST"
	case 404:
		return "NOT FOUND"
	}
	return ""
}

func encodeHeaders(headers map[string]string) []byte {
	if headers == nil {
		return []byte{'\r', '\n'}
	}
	encodedHeaders := make([]byte, 0, 64)
	for headerName, headerVal := range headers {
		encodedHeaders = fmt.Appendf(encodedHeaders, "%s: %s\r\n", headerName, headerVal)
	}
	encodedHeaders = append(encodedHeaders, '\r', '\n')
	return encodedHeaders
}
