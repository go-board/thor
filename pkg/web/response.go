package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-board/x-go/xnet/xhttp"
)

func WriteJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set(xhttp.HeaderContentType, xhttp.MIMEApplicationJSON)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteError(w http.ResponseWriter, httpStatus int, code int, msg string) error {
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(map[string]interface{}{"code": code, "msg": msg})
}
