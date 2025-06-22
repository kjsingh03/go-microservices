package handler

import (
	"encoding/json"
	"log"
	"logger/internal/helpers"
	"logger/internal/models"
	"logger/types"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Home handler - welcome message
func Home(w http.ResponseWriter, r *http.Request) {
	payload := types.JsonResponse{
		Success: true,
		Message: "Welcome to Logger Service",
	}

	err := helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
	}
}

// GetAllLogs - Retrieve all log entries
func GetAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := models.AppModels.LogEntry.All()
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to fetch logs",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		log.Printf("Error fetching logs: %v", err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Logs retrieved successfully",
		Data:    logs,
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// GetLogByID - Retrieve a single log entry by ID
func GetLogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if strings.TrimSpace(id) == "" {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log ID is required",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	logEntry, err := models.AppModels.LogEntry.GetOne(id)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log not found",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusNotFound, payload)
		log.Printf("Error fetching log with ID %s: %v", id, err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log retrieved successfully",
		Data:    logEntry,
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// CreateLog - Create a new log entry
func CreateLog(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Invalid JSON format",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// Validate required fields
	if strings.TrimSpace(requestBody.Name) == "" {
		payload := types.JsonResponse{
			Success: false,
			Message: "Name field is required",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	if strings.TrimSpace(requestBody.Data) == "" {
		payload := types.JsonResponse{
			Success: false,
			Message: "Data field is required",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Create new log entry
	newLog := models.LogEntry{
		Name: strings.TrimSpace(requestBody.Name),
		Data: strings.TrimSpace(requestBody.Data),
	}

	err = models.AppModels.LogEntry.Insert(newLog)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to create log entry",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		log.Printf("Error creating log: %v", err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry created successfully",
		Data: map[string]interface{}{
			"name": newLog.Name,
			"data": newLog.Data,
		},
	}

	err = helpers.WriteJSON(w, http.StatusCreated, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// UpdateLog - Update an existing log entry
func UpdateLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if strings.TrimSpace(id) == "" {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log ID is required",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// First, check if the log exists
	existingLog, err := models.AppModels.LogEntry.GetOne(id)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log not found",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusNotFound, payload)
		log.Printf("Error finding log with ID %s: %v", id, err)
		return
	}

	var requestBody struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Invalid JSON format",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// Update fields if provided
	if strings.TrimSpace(requestBody.Name) != "" {
		existingLog.Name = strings.TrimSpace(requestBody.Name)
	}
	if strings.TrimSpace(requestBody.Data) != "" {
		existingLog.Data = strings.TrimSpace(requestBody.Data)
	}

	result, err := models.AppModels.LogEntry.Update(existingLog)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to update log entry",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		log.Printf("Error updating log: %v", err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry updated successfully",
		Data: map[string]interface{}{
			"matched_count":  result.MatchedCount,
			"modified_count": result.ModifiedCount,
			"updated_log":    existingLog,
		},
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// DeleteLog - Delete a log entry by ID
func DeleteLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if strings.TrimSpace(id) == "" {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log ID is required",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Validate ObjectID format
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Invalid log ID format",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Check if log exists before deleting
	_, err = models.AppModels.LogEntry.GetOne(id)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Log not found",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusNotFound, payload)
		return
	}

	// Delete the log entry
	result, err := models.AppModels.LogEntry.Delete(objectID)
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to delete log entry",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		log.Printf("Error deleting log: %v", err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log entry deleted successfully",
		Data: map[string]interface{}{
			"deleted_count": result.DeletedCount,
		},
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// DropAllLogs - Drop the entire logs collection (use with caution)
func DropAllLogs(w http.ResponseWriter, r *http.Request) {
	// Add confirmation header for safety
	confirmation := r.Header.Get("X-Confirm-Drop")
	if confirmation != "yes" {
		payload := types.JsonResponse{
			Success: false,
			Message: "This operation requires confirmation. Add header 'X-Confirm-Drop: yes'",
		}
		helpers.WriteJSON(w, http.StatusBadRequest, payload)
		return
	}

	err := models.AppModels.LogEntry.DropCollection()
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to drop logs collection",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		log.Printf("Error dropping collection: %v", err)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "All logs have been deleted successfully",
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// GetLogsStats - Get statistics about logs
func GetLogsStats(w http.ResponseWriter, r *http.Request) {
	logs, err := models.AppModels.LogEntry.All()
	if err != nil {
		payload := types.JsonResponse{
			Success: false,
			Message: "Failed to fetch logs for statistics",
			Error:   err.Error(),
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	stats := map[string]interface{}{
		"total_logs": len(logs),
		"timestamp":  "Generated statistics",
	}

	if len(logs) > 0 {
		stats["oldest_log"] = logs[len(logs)-1].CreatedAt
		stats["newest_log"] = logs[0].CreatedAt
	}

	payload := types.JsonResponse{
		Success: true,
		Message: "Log statistics retrieved successfully",
		Data:    stats,
	}

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}