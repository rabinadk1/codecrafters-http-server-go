package main

import (
	"fmt"
	"strings"
)

func getHeaderValue(headerName string, stringBuffer string) (string, error) {
	contentTypeIndex := strings.Index(stringBuffer, headerName)
	if contentTypeIndex == -1 {
		return "", fmt.Errorf("Could not find Content-Type in request: %s\n", stringBuffer)
	}

	remString := stringBuffer[contentTypeIndex+len(headerName):]

	rnIndex := strings.Index(remString, "\r\n")
	if rnIndex == -1 {
		return "", fmt.Errorf("Could not find \\r\\n in request: %s\n", stringBuffer)
	}

	return remString[:rnIndex], nil
}
