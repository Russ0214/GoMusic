package api

import (
	"net/http"
	"GoMusic/handler"
)

var router = handler.NewRouter()

// Handler handles all requests
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
} 