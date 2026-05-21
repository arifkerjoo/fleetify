package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Payload struct {
	Event     string      `json:"event"`
	ReportID  string      `json:"report_id"`
	Status    string      `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func SendWebhook(url string, payload Payload) {
	go func() {
		if url == "" {
			return
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[webhook] marshal error: %v", err)
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			log.Printf("[webhook] build request error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[webhook] send error: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("[webhook] event=%s report=%s status=%d",
			payload.Event, payload.ReportID, resp.StatusCode)
	}()
}
