package httpentity

import (
	"fmt"
	"strconv"
)

type Response struct {
	ResponseCode uint16
	Version      string
	Body         []byte
}

func (response *Response) Encode() []byte {
	data := make([]byte, 0, 1024)
	antent := buildResponseCodeText(response.Version, response.ResponseCode)
	data = append(data, []byte(antent)...)
	data = append(data, response.Body...)
	return data
}

func buildResponseCodeText(version string, responseCode uint16) string {
	return fmt.Sprintf("%s %s %s\r\n\r\n",
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
