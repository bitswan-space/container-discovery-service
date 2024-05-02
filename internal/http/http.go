package http

import (
	"encoding/json"
	"net/http"

	"bitswan.space/container-discovery-service/internal/logger"
	"bitswan.space/container-discovery-service/pkg"
)

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
	// Optional fields. Mostly because of the Bitville API that requires these fields.
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
}

func Response(w http.ResponseWriter, code int, response any) {
	resJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error.Printf("failed to marshal response because: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.WriteHeader(code)

	_, err = w.Write(resJSON)
	if err != nil {
		logger.Error.Printf("failed to write response because %v", err)
	}
}

// Error prints & optionally logs an error message.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := pkg.ErrorCode(err), pkg.ErrorMessage(err)

	// Log & report internal errors.
	if code == pkg.INTERNAL_ERROR {
		pkg.ReportError(r.Context(), err, r)
		LogError(r, err)
	}

	// Print user message to response based on request accept header.
	switch r.Header.Get("Accept") {
	default:
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(ErrorStatusCode(code))
		err = json.NewEncoder(w).Encode(&ErrorResponse{Error: message, Code: code, Status: "failed"})
		if err != nil {
			logger.Error.Printf("failed to write response because %v", err)
		}
	}
}

// LogError logs an error with the HTTP route information.
func LogError(r *http.Request, err error) {
	logger.Error.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	pkg.INVALID_ERROR:         http.StatusBadRequest,
	pkg.NOT_FOUND_ERROR:       http.StatusNotFound,
	pkg.NOT_IMPLEMENTED_ERROR: http.StatusNotImplemented,
	pkg.AUTHENTICATION_ERROR:  http.StatusUnauthorized,
	pkg.INTERNAL_ERROR:        http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a WTF error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
