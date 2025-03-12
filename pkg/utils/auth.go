package utils

import (
    "fmt"
    "os"
)

// GetAuthToken retrieves an authentication token for accessing the Confluence API.
func GetAuthToken() (string, error) {
    token := os.Getenv("CONFLUENCE_API_TOKEN")
    if token == "" {
        return "", fmt.Errorf("authentication token not found in environment variables")
    }
    return token, nil
}