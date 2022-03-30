package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I70l0teN4ik/go-rest/src/rest"
)

func TestServer(t *testing.T) {
	type reqArgs struct {
		method string
		route  string
	}
	tests := []struct {
		name string
		reqArgs
		wantCode int
	}{
		{"default", reqArgs{http.MethodGet, rest.DefaultEndpoint}, http.StatusOK},
		{"health", reqArgs{http.MethodGet, rest.HealthEndpoint}, http.StatusOK},
		{"not found", reqArgs{http.MethodHead, "/nonsense"}, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := rest.NewServer("testing", http.NewServeMux(), nil)
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.route, nil)
			srv.ServeHTTP(res, req)
			rest.AssertContentType(t, res)
			rest.AssertStatus(t, res, tt.wantCode)
		})
	}
}

func TestStubHealthHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, rest.HealthEndpoint, nil)
	handler := http.HandlerFunc(rest.StubHealthHandler)

	handler.ServeHTTP(res, req)

	rest.AssertStatus(t, res, http.StatusOK)
	rest.AssertContentType(t, res)

	var got rest.HealthResponse
	rest.ParseJSON(t, res, &got)
	if got.Status != rest.OK {
		t.Fatalf("bad health check status, got %s, want %s", got.Status, rest.OK)
	}
}

func TestHandleResponse(t *testing.T) {
	tests := []struct {
		name   string
		data   interface{}
		status int
	}{
		{"ok", map[string]string{"test": "ok"}, http.StatusOK},
		{"created", map[string]string{"name": "new"}, http.StatusCreated},
		{"bad req", nil, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			rest.HandleResponse(res, tt.data, tt.status)
			rest.AssertContentType(t, res)
			rest.AssertStatus(t, res, tt.status)
			var data interface{}
			rest.ParseJSON(t, res, &data)
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		data   rest.CodeError
		status int
	}{
		{&rest.BadRequestError{}, http.StatusBadRequest},
		{&rest.MethodNotAllowedError{}, http.StatusMethodNotAllowed},
		{rest.NewError(errors.New("teapot"), http.StatusTeapot), http.StatusTeapot},
		{rest.NewInternalError(errors.New("internal")), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.data.Error(), func(t *testing.T) {
			res := httptest.NewRecorder()
			rest.HandleError(res, tt.data)
			rest.AssertContentType(t, res)
			rest.AssertStatus(t, res, tt.status)
			var data rest.ErrorResponse
			rest.ParseJSON(t, res, &data)
		})
	}
}
