package server

import (
	"fmt"
	"net/http"

	h "simple-repo/pkg/handlers"
)

func New(port int) (*http.Server, error) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(h.Handler),
	}

	return srv, nil
}
