package main

import (
	"fmt"
	"net/http"
)

func Readiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)

	status, err := writer.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("failed to write return: %s\n", err)
	}
	fmt.Printf("Body is: %v\n", status)
}
