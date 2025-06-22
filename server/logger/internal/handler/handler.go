// internal/handler/handler.go
package handlers

import (
	"encoding/json"
	"logger/internal/helpers"
	"logger/types"
	"net/http"
	"context"
	"time"

	"github.com/gorilla/mux"
)

type LogHandler struct {
	logService types.LogServiceInterface
}

func NewLogHandler(logService types.LogServiceInterface) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

func (h *LogHandler) Health(w http.ResponseWriter, r *http.Request) {
	payload := types.JsonResponse{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		},
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) Home(w http.ResponseWriter, r *http.Request) {
	payload := types.JsonResponse{
		Success: true,
		Message: "Welcome to Logger Service API v1",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"endpoints": map[string]string{
				"health": "/health",
				"logs":   "/api/v1/logs",
				"stats":  "/api/v1/logs/stats",
			},
		},
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) GetAllLogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	logs, err := h.logService.GetAllLogs(ctx)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to fetch logs",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Logs retrieved successfully",
		Data:    logs,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) GetLogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	logEntry, err := h.logService.GetLogByID(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "log not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "log ID is required" {
			statusCode = http.StatusBadRequest
		}

		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, statusCode, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log retrieved successfully",
		Data:    logEntry,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var req types.CreateLogRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Invalid JSON format",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	createdLog, err := h.logService.CreateLog(ctx, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "name field is required" || err.Error() == "data field is required" {
			statusCode = http.StatusBadRequest
		}

		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, statusCode, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry created successfully",
		Data:    createdLog,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req types.UpdateLogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Invalid JSON format",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	updatedLog, err := h.logService.UpdateLog(ctx, id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "log not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "log ID is required" {
			statusCode = http.StatusBadRequest
		}

		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, statusCode, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry updated successfully",
		Data:    updatedLog,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	err := h.logService.DeleteLog(ctx, id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "log not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "log ID is required" {
			statusCode = http.StatusBadRequest
		}

		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, statusCode, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry deleted successfully",
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) DropAllLogs(w http.ResponseWriter, r *http.Request) {
	// Check for confirmation query parameter
	if r.URL.Query().Get("confirm") != "true" {
		payload := types.JsonResponse{
			Success: false,
			Message: "This operation requires confirmation. Add query parameter 'confirm=true'",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	err := h.logService.DropAllLogs(ctx)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "All logs have been deleted successfully",
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}

func (h *LogHandler) GetLogsStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	stats, err := h.logService.GetLogStats(ctx)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log statistics retrieved successfully",
		Data:    stats,
	}

	helpers.WriteJSON(w, http.StatusOK, payload)
}