package utils

import (
	"authentication/types"
	"encoding/json"
	"errors"
	"net/http"
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
