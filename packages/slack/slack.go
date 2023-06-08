package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lambda-deploy/packages/utils"
	"net/http"
)

type MessageData struct {
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

func SendMessage(m MessageData) error {
	webhookResult, err := utils.GetParameterValue("LAMBDA_DEPLOY_SLACK_URL")
	if err != nil {
		return fmt.Errorf("error getting webhook URL: %v", err)
	}

	payload, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	resp, err := http.Post(webhookResult.Parameter.Value, "application/json; charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error posting message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: HTTP status code is not 200 (HTTP status code: %v)", resp.StatusCode)
	}

	return nil
}
