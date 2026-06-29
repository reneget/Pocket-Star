package pulse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Login(username, password string) error {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	data, _ := json.Marshal(body)
	resp, err := c.http.Post(c.BaseURL+"/api/login", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pulse login: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
		Ok    bool   `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("pulse login decode: %w", err)
	}
	if !result.Ok {
		return fmt.Errorf("pulse login failed")
	}
	c.token = result.Token
	return nil
}

type HostInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Platform string `json:"platform"`
	Status   string `json:"status"`
}

func (c *Client) GetHosts() ([]HostInfo, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/api/hosts", nil)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pulse hosts: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Hosts []HostInfo `json:"hosts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("pulse hosts decode: %w", err)
	}
	return result.Hosts, nil
}

type MetricsSnapshot struct {
	CPU     float64 `json:"cpu"`
	RAM     float64 `json:"ram"`
	Disk    float64 `json:"disk"`
	Uptime  int64   `json:"uptime"`
}

func (c *Client) GetHostMetrics(hostID string) (*MetricsSnapshot, error) {
	req, _ := http.NewRequest("GET", c.BaseURL+"/api/hosts/"+hostID+"/metrics", nil)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pulse metrics: %w", err)
	}
	defer resp.Body.Close()

	var ms MetricsSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&ms); err != nil {
		return nil, fmt.Errorf("pulse metrics decode: %w", err)
	}
	return &ms, nil
}

func (c *Client) Health() error {
	resp, err := c.http.Get(c.BaseURL + "/api/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("pulse health returned %d", resp.StatusCode)
	}
	return nil
}
