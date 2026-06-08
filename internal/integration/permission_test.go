package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestPermissionUserCannotUpdateOtherUserProject(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// =====================
	// User A register/login
	// =====================

	emailA := fmt.Sprintf("user_a_%d@gmail.com", time.Now().UnixNano())
	password := "123456"

	doJSON(t, http.MethodPost, baseURL+"/auth/register", "", map[string]any{
		"email":     emailA,
		"password":  password,
		"full_name": "User A",
	}, http.StatusCreated)

	loginA := doJSON(t, http.MethodPost, baseURL+"/auth/login", "", map[string]any{
		"email":    emailA,
		"password": password,
	}, http.StatusOK)

	tokenA := findString(loginA, "token", "access_token")
	if tokenA == "" {
		t.Fatal("expected token A")
	}

	// =====================
	// User A create project
	// =====================

	projectRes := doJSON(t, http.MethodPost, baseURL+"/projects", tokenA, map[string]any{
		"name":        "User A Project",
		"description": "Created by User A",
	}, http.StatusCreated)

	projectID := findString(projectRes, "id")
	if projectID == "" {
		t.Fatal("expected project id")
	}

	// =====================
	// User B register/login
	// =====================

	emailB := fmt.Sprintf("user_b_%d@gmail.com", time.Now().UnixNano())

	doJSON(t, http.MethodPost, baseURL+"/auth/register", "", map[string]any{
		"email":     emailB,
		"password":  password,
		"full_name": "User B",
	}, http.StatusCreated)

	loginB := doJSON(t, http.MethodPost, baseURL+"/auth/login", "", map[string]any{
		"email":    emailB,
		"password": password,
	}, http.StatusOK)

	tokenB := findString(loginB, "token", "access_token")
	if tokenB == "" {
		t.Fatal("expected token B")
	}

	// ============================
	// User B update User A project
	// Expected: 403 Forbidden
	// ============================

	body, err := json.Marshal(map[string]any{
		"name":        "Hacked Project",
		"description": "User B should not update this project",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		baseURL+"/projects/"+projectID,
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenB)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", res.StatusCode)
	}
}

func TestPermissionUserCannotDeleteOtherUserProject(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	emailA := fmt.Sprintf("delete_a_%d@gmail.com", time.Now().UnixNano())
	emailB := fmt.Sprintf("delete_b_%d@gmail.com", time.Now().UnixNano())
	password := "123456"

	doJSON(t, http.MethodPost, baseURL+"/auth/register", "", map[string]any{
		"email":     emailA,
		"password":  password,
		"full_name": "Delete User A",
	}, http.StatusCreated)

	loginA := doJSON(t, http.MethodPost, baseURL+"/auth/login", "", map[string]any{
		"email":    emailA,
		"password": password,
	}, http.StatusOK)

	tokenA := findString(loginA, "token", "access_token")
	if tokenA == "" {
		t.Fatal("expected token A")
	}

	projectRes := doJSON(t, http.MethodPost, baseURL+"/projects", tokenA, map[string]any{
		"name":        "Delete Permission Project",
		"description": "Created by User A",
	}, http.StatusCreated)

	projectID := findString(projectRes, "id")
	if projectID == "" {
		t.Fatal("expected project id")
	}

	doJSON(t, http.MethodPost, baseURL+"/auth/register", "", map[string]any{
		"email":     emailB,
		"password":  password,
		"full_name": "Delete User B",
	}, http.StatusCreated)

	loginB := doJSON(t, http.MethodPost, baseURL+"/auth/login", "", map[string]any{
		"email":    emailB,
		"password": password,
	}, http.StatusOK)

	tokenB := findString(loginB, "token", "access_token")
	if tokenB == "" {
		t.Fatal("expected token B")
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		baseURL+"/projects/"+projectID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenB)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", res.StatusCode)
	}
}
