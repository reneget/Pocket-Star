package kuma

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
		return fmt.Errorf("kuma login: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
		Ok    bool   `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("kuma login decode: %w", err)
	}
	if !result.Ok {
		return fmt.Errorf("kuma login failed")
	}
	c.token = result.Token
	return nil
}

type MonitorConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Retry    int    `json:"retry"`
	Maxretry int    `json:"maxretry"`
}

func (c *Client) AddPingMonitor(name, ip string) error {
	cfg := MonitorConfig{
		Name:     name,
		Type:     "ping",
		URL:      ip,
		Interval: 60,
		Retry:    3,
		Maxretry: 5,
	}
	data, _ := json.Marshal(cfg)

	req, err := http.NewRequest("POST", c.BaseURL+"/api/monitor", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("kuma request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("kuma add monitor: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("kuma returned status %d", resp.StatusCode)
	}
	return nil
}

const ComposeYAML = `services:
  uptime-kuma:
    image: louislam/uptime-kuma:latest
    container_name: uptime-kuma
    ports:
      - "3001:3001"
    volumes:
      - ./pstar-data/uptime:/app/data
    restart: unless-stopped
`
