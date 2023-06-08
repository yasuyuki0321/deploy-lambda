package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Parameter struct {
	Arn              string      `json:"ARN"`
	DataType         string      `json:"DataType"`
	LastModifiedDate time.Time   `json:"LastModifiedDate"`
	Name             string      `json:"Name"`
	Selector         interface{} `json:"Selector"`
	SourceResult     interface{} `json:"SourceResult"`
	Type             string      `json:"Type"`
	Value            string      `json:"Value"`
	Version          int         `json:"Version"`
}

type ParameterResult struct {
	Parameter      Parameter `json:"Parameter"`
	ResultMetadata struct{}  `json:"ResultMetadata"`
}

func GetParameterValue(parameterPath string) (ParameterResult, error) {
	urlEncodedPath := url.QueryEscape(parameterPath)
	requestURL := fmt.Sprintf("http://localhost:2773/systemsmanager/parameters/get/?name=%s", urlEncodedPath)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return ParameterResult{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-Aws-Parameters-Secrets-Token", os.Getenv("AWS_SESSION_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ParameterResult{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var result ParameterResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ParameterResult{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	log.Printf("Successfully retrieved parameter %s", parameterPath)
	return result, nil
}
