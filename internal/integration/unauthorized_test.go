package integration

import (
	"net/http"
	"testing"
)

func TestUnauthorizedCreateProjectWithoutToken(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/projects",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			res.StatusCode,
		)
	}
}

func TestUnauthorizedCreateProjectWithInvalidToken(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/projects",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(
		"Authorization",
		"Bearer abcxyz123",
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			res.StatusCode,
		)
	}
}
