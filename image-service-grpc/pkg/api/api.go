package api

import (
	"io"
	"net/http"
)

func ImageEndpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: Call the image service from here
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
