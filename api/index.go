package handler

import (
	"go-auth-backend-api/app"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	h, err := app.Handler()
	if err != nil {
		http.Error(w, "bootstrap failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.ServeHTTP(w, r)
}
