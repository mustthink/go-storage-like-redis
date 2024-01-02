package handlers

import (
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

func Handler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	request := RequestByMethod(r.Method)
	request.readRequestBody(r.Body)
	data, code := request.getResponse(s)

	w.WriteHeader(code)
	w.Write(data)
}
