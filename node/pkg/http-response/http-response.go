package httpresponse

import (
	"encoding/json"
	"net/http"
)

type HttpResponsePayload struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpResponseErrorPayload struct {
	Message string `json:"message"`
	ErrCode string `json:"err_code"`
	Data    any    `json:"data"`
}

func WriteSuccessResponse(w http.ResponseWriter, r *http.Request, code int, msg string, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	payload := HttpResponsePayload{
		Message: msg,
		Data:    data,
	}
	json.NewEncoder(w).Encode(&payload)
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, code int, msg, errCode string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	payload := HttpResponseErrorPayload{
		Message: msg,
		ErrCode: errCode,
		Data:    nil,
	}
	json.NewEncoder(w).Encode(&payload)
}
