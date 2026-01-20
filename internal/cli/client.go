package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *APIError       `json:"error,omitempty"`
	Meta    *APIMeta        `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type APIMeta struct {
	Total  int64 `json:"total,omitempty"`
	Limit  int   `json:"limit,omitempty"`
	Offset int   `json:"offset,omitempty"`
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	if result != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}
	}

	return nil
}

type Link struct {
	ID          int64    `json:"id"`
	Slug        string   `json:"slug"`
	OriginalURL string   `json:"original_url"`
	Domain      string   `json:"domain,omitempty"`
	HasPassword bool     `json:"has_password"`
	ExpiresAt   *string  `json:"expires_at,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ClickCount  int64    `json:"click_count"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type CreateLinkRequest struct {
	URL      string   `json:"url"`
	Slug     string   `json:"slug,omitempty"`
	Domain   string   `json:"domain,omitempty"`
	Password string   `json:"password,omitempty"`
	TTLHours int      `json:"ttl_hours,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func (c *Client) CreateLink(req CreateLinkRequest) (*Link, error) {
	var link Link
	if err := c.do("POST", "/api/v1/links", req, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (c *Client) GetLink(slug string) (*Link, error) {
	var link Link
	if err := c.do("GET", "/api/v1/links/"+slug, nil, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

type ListLinksOptions struct {
	Search   string
	Tags     []string
	FolderID *int64
	Limit    int
	Offset   int
}

func (c *Client) ListLinks(opts ListLinksOptions) ([]Link, error) {
	path := "/api/v1/links"
	params := url.Values{}

	if opts.Search != "" {
		params.Set("search", opts.Search)
	}
	for _, tag := range opts.Tags {
		params.Add("tags", tag)
	}
	if opts.FolderID != nil {
		params.Set("folder_id", fmt.Sprintf("%d", *opts.FolderID))
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var links []Link
	if err := c.do("GET", path, nil, &links); err != nil {
		return nil, err
	}
	return links, nil
}

func (c *Client) DeleteLink(slug string, permanent bool) error {
	path := "/api/v1/links/" + slug
	if permanent {
		path += "?permanent=true"
	}
	return c.do("DELETE", path, nil, nil)
}

type BulkDeleteRequest struct {
	Slugs     []string `json:"slugs"`
	Permanent bool     `json:"permanent"`
}

type BulkDeleteResponse struct {
	Deleted []string `json:"deleted"`
	Failed  []string `json:"failed"`
}

func (c *Client) BulkDeleteLinks(slugs []string, permanent bool) (*BulkDeleteResponse, error) {
	req := BulkDeleteRequest{
		Slugs:     slugs,
		Permanent: permanent,
	}
	var result BulkDeleteResponse
	if err := c.do("DELETE", "/api/v1/links", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ClickStats struct {
	TotalClicks  int64           `json:"total_clicks"`
	ClicksByDay  []DayStats      `json:"clicks_by_day,omitempty"`
	TopReferrers []ReferrerStats `json:"top_referrers,omitempty"`
}

type DayStats struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

func (c *Client) GetStats(slug string) (*ClickStats, error) {
	var stats ClickStats
	if err := c.do("GET", "/api/v1/stats/"+slug, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
