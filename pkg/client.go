// Copyright 2025 Lincoln Institute of Land Policy
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
)

// A client for the OGC API Processes Spec
type ProcessesClient struct {
	// base url to the ogc compliant api
	BaseUrl string
	// http client to use for requests
	httpClient *http.Client
}

func NewProcessesClientFromHttpClient(baseUrl string, httpClient *http.Client) (*ProcessesClient, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("base url cannot be empty")
	}
	return &ProcessesClient{
		BaseUrl:    baseUrl,
		httpClient: httpClient,
	}, nil
}

func NewProcessesClient(baseUrl string) (*ProcessesClient, error) {
	return NewProcessesClientFromHttpClient(baseUrl, http.DefaultClient)
}

type Link struct {
	Type_    string `json:"type"`
	Rel      string `json:"rel"`
	Href     string `json:"href"`
	Title    string `json:"title"`
	Hreflang string `json:"hreflang"`
}

type JobControlOption string

const (
	SyncSupport  JobControlOption = "sync-execute"
	AsyncSupport JobControlOption = "async-execute"
)

type ProcessInfoResponse struct {
	ProcessInfo []struct {
		Version            string             `json:"version"`
		Id                 string             `json:"id"`
		Title              string             `json:"title"`
		Description        string             `json:"description"`
		JobControlOptions  []JobControlOption `json:"jobControlOptions"`
		Keywords           []string           `json:"keywords"`
		OutputTransmission []string           `json:"outputTransmission"`
		Links              []Link             `json:"links"`
	} `json:"processes"`
	Links []Link `json:"links"`
}

// List all processes available on the server
func (c *ProcessesClient) ListProcesses() (ProcessInfoResponse, error) {
	response, err := c.httpClient.Get(c.BaseUrl + "/processes?f=json")
	if err != nil {
		return ProcessInfoResponse{}, err
	}

	defer func() { _ = response.Body.Close() }()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ProcessInfoResponse{}, err
	}

	if response.StatusCode >= 400 {
		return ProcessInfoResponse{}, fmt.Errorf("error in process list: %s", string(bodyBytes))
	}

	var processInfoResponse ProcessInfoResponse
	err = json.Unmarshal(bodyBytes, &processInfoResponse)
	return processInfoResponse, err
}

type ProcessInfo struct {
	Version            string             `json:"version"`
	Id                 string             `json:"id"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	Links              []Link             `json:"links"`
	JobControlOptions  []JobControlOption `json:"jobControlOptions"`
	Keywords           []string           `json:"keywords"`
	Inputs             map[string]IOInfo  `json:"inputs"`
	Outputs            map[string]IOInfo  `json:"outputs"`
	Example            Example            `json:"example"`
	OutputTransmission []string           `json:"outputTransmission"`
}

// Represents either an input or an output
type IOInfo struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"schema"`
	MinOccurs   *int           `json:"minOccurs"`
	MaxOccurs   *int           `json:"maxOccurs"`
	Keywords    []string       `json:"keywords"`
}

// Example input section
type Example struct {
	Inputs map[string]any `json:"inputs"`
}

func (c *ProcessesClient) GetProcessInfo(processId string) (ProcessInfo, error) {
	if processId == "" {
		return ProcessInfo{}, fmt.Errorf("process name cannot be empty")
	}

	url := c.BaseUrl + "/processes/" + processId + "?f=json"

	response, err := c.httpClient.Get(url)
	if err != nil {
		return ProcessInfo{}, err
	}

	defer func() { _ = response.Body.Close() }()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ProcessInfo{}, err
	}

	if response.StatusCode >= 400 {
		return ProcessInfo{}, fmt.Errorf("error in process info: %s", string(bodyBytes))
	}

	var processInfo ProcessInfo
	err = json.Unmarshal(bodyBytes, &processInfo)
	return processInfo, err
}

type ProcessExecutionMode string

const (
	Sync  ProcessExecutionMode = "sync"
	Async ProcessExecutionMode = "async"
)

func (c *ProcessesClient) ExecuteAsync(processID string, inputs ...map[string]any) (ExecuteResponse, error) {
	return c.execute(processID, Async, inputs...)
}

func (c *ProcessesClient) ExecuteSync(processID string, inputs ...map[string]any) (ExecuteResponse, error) {
	return c.execute(processID, Sync, inputs...)
}

type ExecuteResponse struct {
	ID     string `json:"id"`
	Value  string `json:"value"`
	JobUrl string
}

func (c *ProcessesClient) execute(processID string, mode ProcessExecutionMode, inputs ...map[string]any) (ExecuteResponse, error) {
	if processID == "" {
		return ExecuteResponse{}, fmt.Errorf("process name cannot be empty")
	}
	if mode != Sync && mode != Async {
		return ExecuteResponse{}, fmt.Errorf("invalid mode: %s; mode must be sync or async", mode)
	}

	url := c.BaseUrl + "/processes/" + processID + "/execution?f=json"

	mergedInputs := make(map[string]any)
	for _, input := range inputs {
		maps.Copy(mergedInputs, input)
	}

	payload := map[string]any{
		"mode":   mode,
		"inputs": mergedInputs,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return ExecuteResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return ExecuteResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ExecuteResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExecuteResponse{}, err
	}

	if resp.StatusCode >= 400 {
		bodyAsString := string(bodyBytes)
		return ExecuteResponse{}, fmt.Errorf("error in process execution: %s", bodyAsString)
	}

	var jobUrl string
	headers := resp.Header
	for k, v := range headers {
		if k == "Location" {
			jobUrl = v[0]
		}
	}

	var executeResponse ExecuteResponse
	err = json.Unmarshal(bodyBytes, &executeResponse)
	if err != nil {
		return ExecuteResponse{}, err
	}

	executeResponse.JobUrl = jobUrl
	return executeResponse, nil
}
