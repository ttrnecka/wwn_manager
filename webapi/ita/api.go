package ita

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

var httpClient = &http.Client{
	Timeout: 3 * time.Minute,
	Transport: &http.Transport{
		MaxIdleConns:    5,
		IdleConnTimeout: 90 * time.Second,
	},
}

type Client struct {
	httpClient *http.Client
	logger     *zerolog.Logger
	token      string
	baseURL    string
}

func NewITAClient(logger *zerolog.Logger) (*Client, error) {
	baseURL := os.Getenv("ITA_API_URI")
	token := os.Getenv("ITA_TOKEN")
	if baseURL == "" || token == "" {
		return nil, errors.New("ITA environment variables are not set up")
	}
	return &Client{
		httpClient: httpClient,
		logger:     logger,
		token:      token,
		baseURL:    baseURL,
	}, nil
}

func (c *Client) GenerateReportTemplate(templateID string, page int, pageSize int) ([]byte, error) {
	urlStr := fmt.Sprintf("%sreport-templates/%s", c.baseURL, templateID)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	q.Set("currentPage", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))

	u.RawQuery = q.Encode()

	// Marshal the empty struct -> this becomes "{}"
	payload, err := json.Marshal(EmptyPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := c.httpClient

	// Create POST request (no body)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Authorization header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	// Perform the request
	c.logger.Info().Msg(fmt.Sprintf("Calling %+v", u))

	// #nosec G704
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
