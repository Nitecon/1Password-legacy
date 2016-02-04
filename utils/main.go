package utils

import (
	"encoding/json"
	"net/http"
)

// AppData contains a standardized response object for the application.
type AppData struct {
	IsError    bool        `json:"isError"`
	Error      string      `json:"errorMessage"`
	StatusCode int         `json:"statusCode"`
	Content    interface{} `json:"content"`
}

// WriteRestResponse writes a standardized rest response.
func (a *AppData) WriteRestResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(a.StatusCode)
		w.Write([]byte(`{"isError":true, "errorMessage":"Critical error", "statusCode":500,"content":{}}`))
		return
	}
	w.WriteHeader(a.StatusCode)
	w.Write(data)
	return
}
