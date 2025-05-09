package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"confluence-exporter/internal/models"
)

// ConfluenceClient handles all interactions with the Confluence API
type ConfluenceClient struct {
	BaseURL    string
	Username   string
	APIToken   string
	HTTPClient *http.Client
}

// NewConfluenceClient creates a new client for interacting with Confluence
func NewConfluenceClient(baseURL, username, apiToken string) *ConfluenceClient {
	return &ConfluenceClient{
		BaseURL:  baseURL,
		Username: username,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetSpaces retrieves all spaces the user has access to
func (c *ConfluenceClient) GetSpaces() ([]models.Space, error) {
	endpoint := "/rest/api/space"

	var allSpaces []models.Space
	start := 0
	limit := 500

	for {
		params := url.Values{}
		params.Add("type", "global")
		params.Add("start", strconv.Itoa(start))
		params.Add("limit", strconv.Itoa(limit))

		resp, err := c.sendRequest("GET", endpoint, params, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result struct {
			Results []models.Space `json:"results"`
			Size    int            `json:"size"`
			Limit   int            `json:"limit"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		allSpaces = append(allSpaces, result.Results...)

		if len(result.Results) < limit {
			break
		}

		start += limit
	}

	return allSpaces, nil
}

// GetPages retrieves all pages in a space
func (c *ConfluenceClient) GetPages(spaceKey string) ([]models.Page, error) {
	endpoint := "/rest/api/content"

	var allPages []models.Page
	start := 0
	limit := 100

	for {
		params := url.Values{}
		params.Add("spaceKey", spaceKey)
		params.Add("type", "page")
		params.Add("expand", "body.storage,version,space")
		params.Add("start", strconv.Itoa(start))
		params.Add("limit", strconv.Itoa(limit))

		resp, err := c.sendRequest("GET", endpoint, params, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result struct {
			Results []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				Type  string `json:"type"`
				Space struct {
					Key string `json:"key"`
				} `json:"space"`
				Body struct {
					Storage struct {
						Value string `json:"value"`
					} `json:"storage"`
				} `json:"body"`
				Version struct {
					Number int `json:"number"`
				} `json:"version"`
				Links struct {
					WebUI string `json:"webui"`
				} `json:"_links"`
			} `json:"results"`
			Size  int `json:"size"`
			Limit int `json:"limit"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		for _, p := range result.Results {
			page := models.Page{
				ID:       p.ID,
				Title:    p.Title,
				SpaceKey: p.Space.Key,
				Version:  p.Version.Number,
				Content:  p.Body.Storage.Value,
				URL:      p.Links.WebUI,
			}
			allPages = append(allPages, page)
		}

		if len(result.Results) < limit {
			break
		}

		start += limit
	}

	return allPages, nil
}

// GetAttachments retrieves all attachments for a page
func (c *ConfluenceClient) GetAttachments(pageID string) ([]models.Attachment, error) {
	endpoint := fmt.Sprintf("/rest/api/content/%s/child/attachment", pageID)

	params := url.Values{}
	params.Add("expand", "version")

	resp, err := c.sendRequest("GET", endpoint, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Metadata struct {
				MediaType string `json:"mediaType"`
				Size      int64  `json:"size"`
			} `json:"metadata"`
			Links struct {
				Download string `json:"download"`
			} `json:"_links"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var attachments []models.Attachment
	for _, a := range result.Results {
		attachment := models.Attachment{
			ID:          a.ID,
			Title:       a.Title,
			FileName:    a.Title,
			MediaType:   a.Metadata.MediaType,
			FileSize:    a.Metadata.Size,
			DownloadURL: a.Links.Download,
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// sendRequest sends an HTTP request to the Confluence API
func (c *ConfluenceClient) sendRequest(method, endpoint string, params url.Values, body io.Reader) (*http.Response, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	// Join the base URL path with the endpoint
	apiURL := baseURL.ResolveReference(&url.URL{
		Path: path.Join(baseURL.Path, endpoint),
	})

	// Add query parameters
	if params != nil {
		apiURL.RawQuery = params.Encode()
	}

	// Create request
	req, err := http.NewRequest(method, apiURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Add basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.APIToken))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	return c.HTTPClient.Do(req)
}

// GetAttachmentContent downloads the content of an attachment
func (c *ConfluenceClient) GetAttachmentContent(downloadURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, err
	}

	// Add basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.APIToken))
	req.Header.Add("Authorization", "Basic "+auth)

	return c.HTTPClient.Do(req)
}

// GetBaseURL returns the base URL of the Confluence instance
func (c *ConfluenceClient) GetBaseURL() string {
	return c.BaseURL
}
