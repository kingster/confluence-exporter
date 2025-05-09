package converter

import (
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// ConvertToMarkdown transforms Confluence HTML content to Markdown
func ConvertToMarkdown(content string) (string, error) {

	markdown, err := htmltomarkdown.ConvertString(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	return markdown, nil

}
