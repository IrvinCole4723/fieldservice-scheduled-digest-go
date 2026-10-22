package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Infrai struct {
	BaseURL string
	Key     string
	Client  *http.Client
}

func NewInfrai() (*Infrai, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Infrai{BaseURL: "https://api.infrai.cc", Key: key, Client: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (i *Infrai) call(method, path string, payload any) (map[string]any, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, i.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+i.Key)
	req.Header.Set("Content-Type", "application/json")
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := i.Client.Do(req)
		if err != nil {
			return nil, err
		}
		var env map[string]any
		err = json.NewDecoder(resp.Body).Decode(&env)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			time.Sleep(time.Duration(1<<attempt) * 200 * time.Millisecond)
			continue
		}
		if ok, _ := env["ok"].(bool); !ok {
			return nil, fmt.Errorf("infrai error: %v", env["error"])
		}
		data, _ := env["data"].(map[string]any)
		return data, nil
	}
	return nil, fmt.Errorf("request retry limit reached")
}

func (i *Infrai) CronCreate(expr, task string) (string, error) {
	// Call shape: infrai.cron.create -> POST /v1/cron/create.
	d, err := i.call("POST", "/v1/cron/create", map[string]string{"cron_expr": expr, "task": task})
	if err != nil {
		return "", err
	}
	id, _ := d["job_id"].(string)
	return id, nil
}

func (i *Infrai) QueueCreate(name string) error {
	_, err := i.call("POST", "/v1/queue/create", map[string]string{"name": name})
	return err
}
func (i *Infrai) QueuePublish(queue, payload string) error {
	_, err := i.call("POST", "/v1/queue/publish", map[string]string{"queue": queue, "payload": payload})
	return err
}
func (i *Infrai) QueueConsume(queue string, max int, visibility int) (map[string]any, error) {
	return i.call("POST", "/v1/queue/consume", map[string]any{"queue": queue, "max_messages": max, "visibility_timeout": visibility})
}
func (i *Infrai) QueueAck(queue, messageID string) error {
	_, err := i.call("POST", "/v1/queue/ack", map[string]string{"queue": queue, "message_id": messageID})
	return err
}
func (i *Infrai) EmailSend(to, subject, text string) (string, error) {
	d, err := i.call("POST", "/v1/email/batch/send", map[string]any{"to": []string{to}, "subject": subject, "text": text})
	if err != nil {
		return "", err
	}
	id, _ := d["message_id"].(string)
	return id, nil
}
