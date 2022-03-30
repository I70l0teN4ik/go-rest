package rest

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func AssertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("bad response status, got %d, want %d", response.Code, want)
	}
}

func AssertContentType(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	got := response.Header().Get(ContentType)
	if got != JsonContent {
		t.Errorf("bad response content type, got %s, want: %s", got, JsonContent)
	}
}

func ParseJSON(t *testing.T, res *httptest.ResponseRecorder, got any) {
	t.Helper()
	err := json.NewDecoder(res.Body).Decode(got)
	if err != nil {
		t.Fatalf("failed to parse res %q to valid %T, %v", res.Body, got, err)
	}
	t.Log(got)
}
