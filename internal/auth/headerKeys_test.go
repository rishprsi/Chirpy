package auth

import (
	"net/http"
	"testing"
)

func TestReturnToken(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer sampletoken")
	token, err := GetBearerToken(headers)
	if err != nil || token != "sampletoken" {
		t.Fatal(err)
	}
}
