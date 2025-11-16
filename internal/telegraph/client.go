package telegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// DefaultBaseURL is the default base URL for the Telegraph API.
	DefaultBaseURL = "https://api.telegra.ph"
)

// Config holds the configuration for the Telegraph client.
type Config struct {
	// BaseURL is the base URL for the Telegraph API.
	// If not set, DefaultBaseURL will be used.
	BaseURL string

	// HTTPClient is the HTTP client to use for making requests.
	// If not set, http.DefaultClient will be used.
	HTTPClient *http.Client
}

// client implements the Client interface.
type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Telegraph client with the given configuration.
func NewClient(cfg Config) Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// CreateAccount creates a new Telegraph account.
func (c *client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	var result Account
	if err := c.callMethod(ctx, "createAccount", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EditAccountInfo updates information about a Telegraph account.
func (c *client) EditAccountInfo(ctx context.Context, req EditAccountInfoRequest) (*Account, error) {
	var result Account
	if err := c.callMethod(ctx, "editAccountInfo", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccountInfo gets information about a Telegraph account.
func (c *client) GetAccountInfo(ctx context.Context, req GetAccountInfoRequest) (*Account, error) {
	var result Account
	if err := c.callMethod(ctx, "getAccountInfo", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeAccessToken revokes access_token and generates a new one.
func (c *client) RevokeAccessToken(ctx context.Context, req RevokeAccessTokenRequest) (*Account, error) {
	var result Account
	if err := c.callMethod(ctx, "revokeAccessToken", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePage creates a new Telegraph page.
func (c *client) CreatePage(ctx context.Context, req CreatePageRequest) (*Page, error) {
	var result Page
	if err := c.callMethod(ctx, "createPage", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EditPage edits an existing Telegraph page.
func (c *client) EditPage(ctx context.Context, req EditPageRequest) (*Page, error) {
	var result Page
	method := fmt.Sprintf("editPage/%s", url.PathEscape(req.Path))
	if err := c.callMethod(ctx, method, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPage gets a Telegraph page.
func (c *client) GetPage(ctx context.Context, req GetPageRequest) (*Page, error) {
	var result Page
	method := fmt.Sprintf("getPage/%s", url.PathEscape(req.Path))
	if err := c.callMethod(ctx, method, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPageList gets a list of pages belonging to a Telegraph account.
func (c *client) GetPageList(ctx context.Context, req GetPageListRequest) (*PageList, error) {
	var result PageList
	if err := c.callMethod(ctx, "getPageList", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetViews gets the number of views for a Telegraph article.
func (c *client) GetViews(ctx context.Context, req GetViewsRequest) (*PageViews, error) {
	var result PageViews
	method := fmt.Sprintf("getViews/%s", url.PathEscape(req.Path))
	if err := c.callMethod(ctx, method, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// callMethod makes an API call to the specified method.
// The Telegraph API supports both GET and POST. We use POST with form-encoded data
// for better compatibility with complex data types like content arrays.
func (c *client) callMethod(ctx context.Context, method string, req any, result any) error {
	// Build URL
	apiURL := fmt.Sprintf("%s/%s", c.baseURL, method)

	// Convert request to form values
	reqData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	var reqMap map[string]any
	if err := json.Unmarshal(reqData, &reqMap); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// Build form values
	values := url.Values{}
	for k, v := range reqMap {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			values.Set(k, val)
		case []any:
			// Handle arrays (like fields parameter) - marshal to JSON string
			if len(val) > 0 {
				jsonBytes, _ := json.Marshal(val)
				values.Set(k, string(jsonBytes))
			}
		case []string:
			// Handle string arrays
			jsonBytes, _ := json.Marshal(val)
			values.Set(k, string(jsonBytes))
		case bool:
			if val {
				values.Set(k, "true")
			} else {
				values.Set(k, "false")
			}
		case *bool:
			if val != nil {
				if *val {
					values.Set(k, "true")
				} else {
					values.Set(k, "false")
				}
			}
		case int:
			values.Set(k, fmt.Sprintf("%d", val))
		case *int:
			if val != nil {
				values.Set(k, fmt.Sprintf("%d", *val))
			}
		default:
			// For complex types, marshal to JSON
			jsonBytes, _ := json.Marshal(val)
			values.Set(k, string(jsonBytes))
		}
	}

	// Create POST request with form-encoded body
	body := bytes.NewBufferString(values.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !apiResp.OK {
		return fmt.Errorf("telegraph API error: %s", apiResp.Error)
	}

	// Unmarshal result
	resultBytes, err := json.Marshal(apiResp.Result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := json.Unmarshal(resultBytes, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return nil
}
