package helpers

import (
	"encoding/json"
	"net/http"

)

// APIResponse ...
type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// APIResponseError ...
type APIResponseError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// Response handler
func Response(w http.ResponseWriter, httpStatus int, data interface{}) {
	apiResponse := new(APIResponse)
	apiResponse.Status = httpStatus
	apiResponse.Data = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiResponse)
}

// ResponseError handler
func ResponseError(w http.ResponseWriter, httpStatus int, err error) {
	apiResponse := new(APIResponseError)
	apiResponse.Error = err.Error()
	apiResponse.Status = httpStatus

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiResponse)
}
