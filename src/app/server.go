package app

import (
	"net/http"
	"strings"

	"github.com/I70l0teN4ik/go-rest/src/rest"
)

const AppName = "Hello world"

func NewServer() *rest.Server {
	router := http.NewServeMux()
	router.HandleFunc("/hello/", HelloHandler)
	return rest.NewServer(AppName, router, nil)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/hello/")
	if "" == name {
		rest.HandleError(w, &rest.BadRequestError{})
		return
	}
	rest.HandleResponse(w, map[string]string{"hello": name}, http.StatusOK)
}
