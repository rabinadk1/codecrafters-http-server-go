package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
)

func getHeaderValue(headerName string, stringBuffer string) (string, error) {
	contentTypeIndex := strings.Index(stringBuffer, headerName)
	if contentTypeIndex == -1 {
		return "", fmt.Errorf("Could not find %s in request: %s\n", headerName, stringBuffer)
	}

	remString := stringBuffer[contentTypeIndex+len(headerName):]

	rnIndex := strings.Index(remString, "\r\n")
	if rnIndex == -1 {
		return "", fmt.Errorf("Could not find \\r\\n in request: %s\n", stringBuffer)
	}

	return remString[:rnIndex], nil
}

func ifGZIPAccepted(stringBuffer string) bool {
	const acceptEncodingPrefix = "Accept-Encoding: "
	acceptEncodingHeader, err := getHeaderValue(acceptEncodingPrefix, stringBuffer)
	if err != nil {
		return false
	}

	for _, acceptEncoder := range strings.Split(acceptEncodingHeader, ", ") {
		if acceptEncoder == "gzip" {
			return true
		}
	}

	return false
}

func compressGZIP(content []byte) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)

	_, err := w.Write(content)
	if err != nil {
		return buffer, err
	}

	// Close early because the buffer may not have written
	err = w.Close()

	return buffer, err
}
