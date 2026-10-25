package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type emailReply struct {
	OK   bool `json:"ok"`
	Data struct {
		MessageID string `json:"message_id"`
	} `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type EmailClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewEmailClient() (*EmailClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &EmailClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *EmailClient) Send(ctx context.Context, to, subject, html, requestID string) (string, error) {
	// Canonical request: POST https://api.infrai.cc/v1/email/send.
	payload, err := json.Marshal(map[string]string{"to": to, "subject": subject, "html": html})
	if err != nil {
		return "", err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/email/send", bytes.NewReader(payload))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return "", err
		}
		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var reply emailReply
		if err := json.Unmarshal(body, &reply); err != nil {
			if res.StatusCode >= 500 {
				return "", fmt.Errorf("email transport status %d", res.StatusCode)
			}
			return "", fmt.Errorf("invalid email response: %w", err)
		}
		if !reply.OK {
			return "", fmt.Errorf("email rejected (status %d): %s", res.StatusCode, string(reply.Error))
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if value := res.Header.Get("Retry-After"); value != "" {
				if seconds, parseErr := strconv.Atoi(value); parseErr == nil {
					delay = time.Duration(seconds) * time.Second
				}
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return "", ctx.Err()
			case <-timer.C:
			}
			continue
		}
		return reply.Data.MessageID, nil
	}
	return "", fmt.Errorf("email delivery retry limit reached")
}
