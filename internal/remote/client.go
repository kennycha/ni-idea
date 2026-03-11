package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is an API client for ni-idea server
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NoteMeta represents note metadata returned from API
type NoteMeta struct {
	Path    string   `json:"path"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Tags    []string `json:"tags"`
	Private bool     `json:"private"`
	Created string   `json:"created"`
	Updated string   `json:"updated"`
}

// Note represents a full note with content
type Note struct {
	NoteMeta
	Content string `json:"content"`
}

// SearchResult represents a search result
type SearchResult struct {
	Path    string   `json:"path"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Tags    []string `json:"tags"`
	Matches []string `json:"matches"`
	Score   float64  `json:"score"`
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListNotes returns all notes from the remote
func (c *Client) ListNotes() ([]*NoteMeta, error) {
	req, err := c.newRequest("GET", "/api/notes", nil)
	if err != nil {
		return nil, err
	}

	var notes []*NoteMeta
	if err := c.do(req, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

// GetNote returns a specific note by path
func (c *Client) GetNote(path string) (*Note, error) {
	req, err := c.newRequest("GET", "/api/notes/"+url.PathEscape(path), nil)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := c.do(req, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// PushNote uploads or updates a note
func (c *Client) PushNote(note *Note) error {
	body, err := json.Marshal(note)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PUT", "/api/notes/"+url.PathEscape(note.Path), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, nil)
}

// DeleteNote removes a note from the remote
func (c *Client) DeleteNote(path string) error {
	req, err := c.newRequest("DELETE", "/api/notes/"+url.PathEscape(path), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Search searches notes on the remote
func (c *Client) Search(query string) ([]*SearchResult, error) {
	req, err := c.newRequest("GET", "/api/search?q="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}

	var results []*SearchResult
	if err := c.do(req, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Ping checks if the server is reachable and auth is valid
func (c *Client) Ping() error {
	req, err := c.newRequest("GET", "/api/ping", nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// GetNoteMeta returns metadata for a specific note (without content)
func (c *Client) GetNoteMeta(path string) (*NoteMeta, error) {
	// We use GetNote and extract metadata since server returns full note
	note, err := c.GetNote(path)
	if err != nil {
		return nil, err
	}
	return &note.NoteMeta, nil
}

// ListNotesMap returns a map of path -> NoteMeta for efficient lookup
func (c *Client) ListNotesMap() (map[string]*NoteMeta, error) {
	notes, err := c.ListNotes()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*NoteMeta, len(notes))
	for _, note := range notes {
		result[note.Path] = note
	}
	return result, nil
}

// IsNotFound returns true if the error indicates a not found response
func IsNotFound(err error) bool {
	return err != nil && err.Error() == "not found"
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized: invalid or missing token")
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
