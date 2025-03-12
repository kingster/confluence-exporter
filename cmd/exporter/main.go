package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"confluence-exporter/internal/api"
	"confluence-exporter/internal/config"
	"confluence-exporter/internal/converter"
	"confluence-exporter/internal/models"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Confluence client
	client := api.NewConfluenceClient(
		cfg.Confluence.BaseURL,
		cfg.Confluence.Username,
		cfg.Confluence.APIToken,
	)

	// Get all pages from specified space
	pages, err := client.GetPages(cfg.Export.SpaceKey)
	if err != nil {
		log.Fatalf("Failed to fetch pages: %v", err)
	}

	fmt.Printf("Found %d pages to export\n", len(pages))

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(cfg.Export.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Process each page
	for _, page := range pages {
		fmt.Printf("Exporting '%s'...\n", page.Title)

		// Convert HTML content to Markdown
		markdown, err := converter.ConvertToMarkdown(page.Content)
		if err != nil {
			log.Printf("Failed to convert page %s: %v", page.Title, err)
			continue
		}

		// Create safe filename
		safeFilename := getSafeFilename(page.Title)
		outputPath := filepath.Join(cfg.Export.OutputDir, safeFilename+".md")

		// Write markdown content to file
		if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
			log.Printf("Failed to write page %s: %v", page.Title, err)
			continue
		}

		// Get and save attachments if enabled
		if cfg.Export.IncludeAttachments {
			attachments, err := client.GetAttachments(page.ID)
			if err != nil {
				log.Printf("Failed to get attachments for page %s: %v", page.Title, err)
				continue
			}

			if len(attachments) > 0 {
				// Create attachments directory
				attachmentsDir := filepath.Join(cfg.Export.OutputDir, "attachments", safeFilename)
				if err := os.MkdirAll(attachmentsDir, 0755); err != nil {
					log.Printf("Failed to create attachments directory for page %s: %v", page.Title, err)
					continue
				}

				fmt.Printf("  Saving %d attachments\n", len(attachments))
				for _, attachment := range attachments {
					outputPath := filepath.Join(attachmentsDir, attachment.FileName)
					fmt.Printf("    Downloading: %s\n", attachment.FileName)

					if err := downloadAttachment(client, attachment, outputPath); err != nil {
						log.Printf("Failed to download attachment %s: %v", attachment.FileName, err)
						continue
					}
				}
			}
		}
	}

	fmt.Println("Export completed successfully!")
}

// getSafeFilename converts a string to a safe filename
func getSafeFilename(name string) string {
	// Replace characters that are not allowed in filenames
	// This is a simplified version, you might need to handle more cases
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
		" ", "_",
	)
	return replacer.Replace(name)
}

// downloadAttachment downloads and saves an attachment to disk
func downloadAttachment(client *api.ConfluenceClient, attachment models.Attachment, outputPath string) error {
	// Construct the full download URL
	downloadURL := client.GetBaseURL() + attachment.DownloadURL

	// Get the file
	resp, err := client.GetAttachmentContent(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the content to the file
	_, err = io.Copy(out, resp.Body)
	return err
}
