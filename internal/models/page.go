package models

// Page represents a Confluence page with all its metadata and content
type Page struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	SpaceKey    string       `json:"spaceKey"`
	Version     int          `json:"version"`
	Content     string       `json:"content"`
	ParentID    string       `json:"parentId,omitempty"`
	URL         string       `json:"url"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
	CreatedBy   string       `json:"createdBy"`
	UpdatedBy   string       `json:"updatedBy"`
	Labels      []Label      `json:"labels,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Label represents a Confluence content label
type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Attachment represents a file attached to a Confluence page
type Attachment struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	FileName    string `json:"fileName"`
	MediaType   string `json:"mediaType"`
	FileSize    int64  `json:"fileSize"`
	DownloadURL string `json:"downloadUrl"`
}

// Space represents a Confluence space
type Space struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
}
