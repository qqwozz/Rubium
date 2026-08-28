package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxRetries     = 2
	retryBaseDelay = 200 * time.Millisecond
)

type Client struct {
	baseURL    string
	anonKey    string
	serviceKey string
	http       *http.Client
}

func NewClient(baseURL, anonKey, serviceKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		anonKey:    anonKey,
		serviceKey: serviceKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{
		"apikey":        c.anonKey,
		"Authorization": "Bearer " + c.anonKey,
		"Content-Type":  "application/json",
	}
}

func (c *Client) headersService() map[string]string {
	return map[string]string{
		"apikey":        c.serviceKey,
		"Authorization": "Bearer " + c.serviceKey,
		"Content-Type":  "application/json",
	}
}

// do — выполнение запроса с retry и отменой по context
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBaseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("executing request: %w", err)
			continue
		}

		if resp.StatusCode < 400 {
			return resp, nil
		}

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()

		if resp.StatusCode < 500 {
			return nil, fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
		}

		lastErr = fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return nil, lastErr
}

func (c *Client) RawQuery(ctx context.Context, endpoint string, useServiceRole bool) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return body, nil
}

func (c *Client) Query(ctx context.Context, endpoint string, useServiceRole bool, result interface{}) error {
	body, err := c.RawQuery(ctx, endpoint, useServiceRole)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

func (c *Client) Patch(ctx context.Context, endpoint string, useServiceRole bool, payload interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=minimal"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Post(ctx context.Context, endpoint string, useServiceRole bool, payload interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=representation"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	if result != nil && len(body) > 0 {
		return json.Unmarshal(body, result)
	}

	return nil
}

func (c *Client) RPC(ctx context.Context, function string, useServiceRole bool, params interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/rpc/%s", c.baseURL, function)

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshaling rpc params: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating rpc request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading rpc response: %w", err)
	}

	if result != nil && len(body) > 0 {
		return json.Unmarshal(body, result)
	}
	return nil
}

func (c *Client) AuthUser(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/auth/v1/user", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("apikey", c.anonKey)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.do(ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("невалидный токен")
	}

	var u struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &u); err != nil {
		return "", err
	}

	if u.ID == "" {
		return "", fmt.Errorf("не удалось получить user_id")
	}
	return u.ID, nil
}

func (c *Client) Delete(ctx context.Context, endpoint string, useServiceRole bool) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=minimal"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}