package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestFullAPIFlow(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	email := fmt.Sprintf("duytest_%d@gmail.com", time.Now().UnixNano())
	password := "123456"

	// 1. Register
	doJSON(t, http.MethodPost, baseURL+"/auth/register", "", map[string]any{
		"email":     email,
		"password":  password,
		"full_name": "Duy Test",
	}, http.StatusCreated)

	// 2. Login
	loginRes := doJSON(t, http.MethodPost, baseURL+"/auth/login", "", map[string]any{
		"email":    email,
		"password": password,
	}, http.StatusOK)

	t.Logf("login response: %+v", loginRes)

	token := findString(loginRes, "token", "access_token")
	if token == "" {
		t.Fatal("expected login token")
	}

	// 3. Create Project
	projectRes := doJSON(t, http.MethodPost, baseURL+"/projects", token, map[string]any{
		"name":        "Integration Test Project",
		"description": "Project created by integration test",
	}, http.StatusCreated)

	projectID := findString(projectRes, "id")
	if projectID == "" {
		t.Fatal("expected project id")
	}

	// 4. Create Task
	taskRes := doJSON(t, http.MethodPost, baseURL+"/tasks", token, map[string]any{
		"project_id":  projectID,
		"title":       "Integration Test Task",
		"description": "Task created by integration test",
		"status":      "TODO",
	}, http.StatusCreated)

	taskID := findString(taskRes, "id")
	if taskID == "" {
		t.Fatal("expected task id")
	}

	// 5. Create Comment
	commentRes := doJSON(t, http.MethodPost, baseURL+"/tasks/"+taskID+"/comments", token, map[string]any{
		"content": "Integration test comment",
	}, http.StatusCreated)

	commentID := findString(commentRes, "id")
	if commentID == "" {
		t.Fatal("expected comment id")
	}
}

func doJSON(
	t *testing.T,
	method string,
	url string,
	token string,
	body map[string]any,
	expectedStatus int,
) map[string]any {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != expectedStatus {
		t.Fatalf(
			"expected status %d, got %d, response: %s",
			expectedStatus,
			res.StatusCode,
			string(rawBody),
		)
	}

	var result map[string]any
	if len(rawBody) > 0 {
		if err := json.Unmarshal(rawBody, &result); err != nil {
			t.Fatalf("failed to decode response as JSON object: %s", string(rawBody))
		}
	}

	return result
}

func findString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := data[key].(string); ok {
			return value
		}
	}

	for _, parentKey := range []string{"data", "Data"} {
		if nested, ok := data[parentKey].(map[string]any); ok {
			for _, key := range keys {
				if value, ok := nested[key].(string); ok {
					return value
				}
			}
		}
	}

	return ""
}
