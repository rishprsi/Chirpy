package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("no authorization token")
	}

	return strings.Replace(tokenString, "Bearer ", "", 1), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no api key found")
	}

	return strings.Replace(apiKey, "ApiKey ", "", 1), nil
}
