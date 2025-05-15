package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type WAResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func SendWACustomMessage(targets []string, message string) (WAResponse, error) {

	var response WAResponse

	wa_url := os.Getenv("WA_NOTIFIER_URL")
	wa_url += "/api/wa/messages"

	if len(targets) <= 0 {
		return response, errors.New("targets cannot be empty")
	}

	if message == "" {
		return response, errors.New("message cannot be empty")
	}

	body_request := map[string]interface{}{
		"messages":         message,
		"whatsapp_numbers": targets,
	}

	body_request_json, err := json.Marshal(body_request)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", wa_url, bytes.NewBuffer(body_request_json))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	// decode response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}
