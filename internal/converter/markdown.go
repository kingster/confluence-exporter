package converter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ConvertToMarkdown transforms Confluence HTML content to Markdown
func ConvertToMarkdown(content string) (string, error) {
	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Process the document
	var markdown strings.Builder
	processNode(doc.Selection, &markdown, 0)

	// Clean up the result
	result := cleanupMarkdown(markdown.String())
	return result, nil
}

// processNode recursively converts HTML nodes to Markdown
func processNode(s *goquery.Selection, markdown *strings.Builder, depth int) {
	s.Each(func(i int, node *goquery.Selection) {
		// Check for Confluence specific macros and elements
		if processConfluenceMacro(node, markdown) {
			return // Skip regular processing if this was a special macro
		}

		// Process based on HTML element
		switch goquery.NodeName(node) {
		case "h1":
			markdown.WriteString("# " + node.Text() + "\n\n")
		case "h2":
			markdown.WriteString("## " + node.Text() + "\n\n")
		case "h3":
			markdown.WriteString("### " + node.Text() + "\n\n")
		case "h4":
			markdown.WriteString("#### " + node.Text() + "\n\n")
		case "h5":
			markdown.WriteString("##### " + node.Text() + "\n\n")
		case "h6":
			markdown.WriteString("###### " + node.Text() + "\n\n")
		case "p":
			text := trimText(node.Text())
			if text != "" {
				markdown.WriteString(text + "\n\n")
			}
		case "ul":
			processChildren(node, markdown, depth)
			markdown.WriteString("\n")
		case "ol":
			processChildren(node, markdown, depth)
			markdown.WriteString("\n")
		case "li":
			indent := strings.Repeat("  ", depth)
			if node.Parent().Is("ol") {
				index := node.Index() + 1
				markdown.WriteString(indent + fmt.Sprintf("%d. ", index))
			} else {
				markdown.WriteString(indent + "- ")
			}
			processChildren(node, markdown, depth+1)
			markdown.WriteString("\n")
		case "a":
			href, exists := node.Attr("href")
			if exists {
				markdown.WriteString("[" + node.Text() + "](" + href + ")")
			} else {
				markdown.WriteString(node.Text())
			}
		case "strong", "b":
			markdown.WriteString("**" + node.Text() + "**")
		case "em", "i":
			markdown.WriteString("*" + node.Text() + "*")
		case "code":
			markdown.WriteString("`" + node.Text() + "`")
		case "pre":
			language := ""
			if node.Find("code").Length() > 0 {
				codeClass, _ := node.Find("code").Attr("class")
				if strings.Contains(codeClass, "language-") {
					language = strings.TrimPrefix(codeClass, "language-")
				}
			}
			markdown.WriteString("```" + language + "\n" + node.Text() + "\n```\n\n")
		case "br":
			markdown.WriteString("\n")
		case "hr":
			markdown.WriteString("---\n\n")
		case "img":
			alt, _ := node.Attr("alt")
			src, exists := node.Attr("src")
			if exists {
				markdown.WriteString("![" + alt + "](" + src + ")")
			}
		case "table":
			processTable(node, markdown)
		case "div", "span":
			// For general containers, just process their children
			processChildren(node, markdown, depth)
		case "#text":
			text := trimText(node.Text())
			if text != "" {
				markdown.WriteString(text)
			}
		default:
			// For other elements, just process their children
			processChildren(node, markdown, depth)
		}
	})
}

// processConfluenceMacro handles specific Confluence macros and elements
func processConfluenceMacro(node *goquery.Selection, markdown *strings.Builder) bool {
	// Check for ac: namespaced elements
	if namespace, _ := node.Attr("xmlns:ac"); namespace != "" ||
		strings.HasPrefix(goquery.NodeName(node), "ac:") {

		// Handle structured macros (info panels, notes, code blocks, etc)
		if node.Is("ac\\:structured-macro") || node.HasClass("confluence-structured-macro") {
			macroType, _ := node.Attr("ac:name")
			switch macroType {
			case "code":
				// Code block macro
				language := node.Find("ac\\:parameter[ac\\:name='language']").Text()
				code := node.Find("ac\\:plain-text-body").Text()
				markdown.WriteString("```" + language + "\n" + code + "\n```\n\n")
			case "info", "note", "warning", "tip":
				// Info/Note/Warning/Tip macro
				content := node.Find(".confluence-information-macro-body").Text()
				markdown.WriteString("> **" + strings.ToUpper(macroType) + ":** " + content + "\n\n")
			default:
				// Other macros - extract any text content
				markdown.WriteString(node.Text() + "\n\n")
			}
			return true
		}

		// Handle task lists
		if node.Is("ac\\:task-list") || node.HasClass("task-list") {
			node.Find("ac\\:task").Each(func(i int, task *goquery.Selection) {
				completed, _ := task.Find("ac\\:task-status").Attr("ac:name")
				taskText := task.Find("ac\\:task-body").Text()

				if completed == "complete" {
					markdown.WriteString("- [x] " + taskText + "\n")
				} else {
					markdown.WriteString("- [ ] " + taskText + "\n")
				}
			})
			markdown.WriteString("\n")
			return true
		}

		// Handle Confluence images
		if node.Is("ac\\:image") || node.HasClass("confluence-embedded-image") {
			imgSrc, _ := node.Find("ri\\:url").Attr("ri:value")
			alt := node.Find("ac\\:alt").Text()
			markdown.WriteString("![" + alt + "](" + imgSrc + ")\n\n")
			return true
		}
	}

	return false
}

// processTable converts HTML tables to Markdown tables
func processTable(table *goquery.Selection, markdown *strings.Builder) {
	var rows [][]string

	// Process header row
	var headers []string
	table.Find("thead tr th").Each(func(i int, th *goquery.Selection) {
		headers = append(headers, th.Text())
	})
	if len(headers) > 0 {
		rows = append(rows, headers)
	}

	// Process body rows
	table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		var row []string
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			row = append(row, td.Text())
		})
		rows = append(rows, row)
	})

	// Build the markdown table
	if len(rows) > 0 {
		// Determine column count from the row with most cells
		maxCols := 0
		for _, row := range rows {
			if len(row) > maxCols {
				maxCols = len(row)
			}
		}

		// Print header row
		if len(headers) > 0 {
			for i := 0; i < len(headers); i++ {
				if i > 0 {
					markdown.WriteString(" | ")
				}
				markdown.WriteString(headers[i])
			}
			markdown.WriteString("\n")

			// Print separator row
			for i := 0; i < len(headers); i++ {
				if i > 0 {
					markdown.WriteString(" | ")
				}
				markdown.WriteString("---")
			}
			markdown.WriteString("\n")
		}

		// Print data rows
		rowStart := 0
		if len(headers) > 0 {
			rowStart = 1
		}

		for i := rowStart; i < len(rows); i++ {
			for j := 0; j < len(rows[i]); j++ {
				if j > 0 {
					markdown.WriteString(" | ")
				}
				markdown.WriteString(rows[i][j])
			}
			markdown.WriteString("\n")
		}

		markdown.WriteString("\n")
	}
}

// processChildren processes child nodes
func processChildren(s *goquery.Selection, markdown *strings.Builder, depth int) {
	s.Contents().Each(func(i int, child *goquery.Selection) {
		processNode(child, markdown, depth)
	})
}

// trimText removes extra whitespace
func trimText(text string) string {
	// Remove extra whitespace
	text = strings.TrimSpace(text)
	return text
}

// cleanupMarkdown performs final cleanup on the markdown text
func cleanupMarkdown(markdown string) string {
	// Remove excessive newlines
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")

	return strings.TrimSpace(markdown)
}
