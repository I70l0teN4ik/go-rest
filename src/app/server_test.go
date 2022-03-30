package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I70l0teN4ik/go-rest/src/app"
	"github.com/I70l0teN4ik/go-rest/src/rest"
)

func Test_Server(t *testing.T) {
	srv := app.NewServer()

	type reqArgs struct {
		method string
		route  string
		body   io.Reader
	}
	tests := []struct {
		reqArgs
		wantCode int
	}{
		{reqArgs{http.MethodGet, rest.DefaultEndpoint, nil}, http.StatusOK},
		{reqArgs{http.MethodGet, "/hello/", nil}, http.StatusBadRequest},
		{reqArgs{http.MethodGet, "/hello/Rick", nil}, http.StatusOK},
		{reqArgs{http.MethodPost, "/hello/Morty", nil}, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.reqArgs.method+" "+tt.reqArgs.route, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.route, tt.body)
			srv.ServeHTTP(res, req)
			rest.AssertContentType(t, res)
			rest.AssertStatus(t, res, tt.wantCode)

			switch tt.reqArgs.route {
			case rest.DefaultEndpoint:
				var got rest.DefaultResponse
				rest.ParseJSON(t, res, &got)
				if got.App != app.AppName {
					t.Fatalf("bad app name, got %v, want %v", got.App, app.AppName)
				}
			}
		})
	}
}
