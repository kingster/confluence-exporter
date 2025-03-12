package converter

import (
	"bytes"
	"fmt"
	"strings"
)

// ConvertToMarkdown converts Confluence content to Markdown format.
func ConvertToMarkdown(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	var markdownBuffer bytes.Buffer

	// Simple conversion logic (this can be expanded)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "h1. ") {
			markdownBuffer.WriteString("# " + strings.TrimPrefix(line, "h1. ") + "\n")
		} else if strings.HasPrefix(line, "h2. ") {
			markdownBuffer.WriteString("## " + strings.TrimPrefix(line, "h2. ") + "\n")
		} else if strings.HasPrefix(line, "* ") {
			markdownBuffer.WriteString("- " + strings.TrimPrefix(line, "* ") + "\n")
		} else {
			markdownBuffer.WriteString(line + "\n")
		}
	}

	return markdownBuffer.String(), nil
}