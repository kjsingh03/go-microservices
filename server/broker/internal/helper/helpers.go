package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"service-broker/types"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1024 * 1024 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("only one JSON object allowed")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	
	out, err := json.Marshal(data)
	
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	return err
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {

	statuscode := http.StatusBadRequest

	payload := types.JsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return WriteJSON(w, statuscode, payload)
}

func ErrorJSONWithExample(w http.ResponseWriter, err error, examples interface{}, status ...int) error {
	statusCode := http.StatusBadRequest
	
	if len(status) > 0 {
		statusCode = status[0]
	}
	
	var payload types.JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	payload.Data = examples
	
	return WriteJSON(w, statusCode, payload)
}

func GetRequestFormatExample() map[string]interface{} {
	return map[string]interface{}{
		"error": "Request must be valid JSON",
		"format": "Content-Type: application/json",
		"example": map[string]interface{}{
			"action": "auth",
			"auth": map[string]string{
				"email":    "user@example.com",
				"password": "your-password",
			},
		},
	}
}

func GetValidActionExamples() map[string]interface{} {
	return map[string]interface{}{
		"examples": map[string]interface{}{
			"auth": map[string]interface{}{
				"action": "auth",
				"auth": map[string]string{
					"email":    "user@example.com",
					"password": "your-password",
				},
			},
			"log": map[string]interface{}{
				"action": "log",
				"log": map[string]string{
					"name": "user-action",
					"data": "User logged in successfully",
				},
			},
			"logdirect": map[string]interface{}{
				"action": "logdirect",
				"log": map[string]string{
					"name": "system-event",
					"data": "Direct log to service",
				},
			},
			"mail": map[string]interface{}{
				"action": "mail",
				"mail": map[string]string{
					"from":    "noreply@example.com",
					"to":      "user@example.com",
					"subject": "Welcome!",
					"message": "Welcome to our service",
				},
			},
		},
	}
}