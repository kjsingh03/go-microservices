package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-broker/internal/helper"
	"service-broker/internal/service"
	"service-broker/types"
	"time"
)

type Handler struct {
	services *service.Services
}

func New(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Service Broker API",
		"status":  "running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload types.RequestPayload
	
	if err := helper.ReadJSON(w, r, &requestPayload); err != nil {
		helper.ErrorJSONWithExample(w, fmt.Errorf("invalid JSON format: %v", err), helper.GetRequestFormatExample(), http.StatusBadRequest)
		return
	}
	
	// Validate the payload
	if err := requestPayload.Validate(); err != nil {
		helper.ErrorJSONWithExample(w, err, helper.GetValidActionExamples(), http.StatusBadRequest)
		return
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	switch requestPayload.Action {
	case "auth":
		if requestPayload.Auth == nil {
			helper.ErrorJSON(w, fmt.Errorf("auth payload is required"), http.StatusBadRequest)
			return
		}
		h.authenticate(ctx, w, *requestPayload.Auth)
	case "log":
		if requestPayload.Log == nil {
			helper.ErrorJSON(w, fmt.Errorf("log payload is required"), http.StatusBadRequest)
			return
		}
		h.logEventViaRabbit(ctx, w, *requestPayload.Log)
	case "logdirect":
		if requestPayload.Log == nil {
			helper.ErrorJSON(w, fmt.Errorf("log payload is required"), http.StatusBadRequest)
			return
		}
		h.logItem(ctx, w, *requestPayload.Log)
	case "mail":
		if requestPayload.Mail == nil {
			helper.ErrorJSON(w, fmt.Errorf("mail payload is required"), http.StatusBadRequest)
			return
		}
		h.sendMail(ctx, w, *requestPayload.Mail)
	default:
		helper.ErrorJSONWithExample(w, fmt.Errorf("unknown action '%s'", requestPayload.Action), helper.GetValidActionExamples(), http.StatusBadRequest)
	}
}

func (h *Handler) authenticate(ctx context.Context, w http.ResponseWriter, authPayload types.AuthPayload) {
	authResp, err := h.services.AuthService.Authenticate(ctx, authPayload)
	if err != nil {
		helper.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	
	response := types.JsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    authResp,
	}
	
	helper.WriteJSON(w, http.StatusAccepted, response)
}

func (h *Handler) logItem(ctx context.Context, w http.ResponseWriter, logPayload types.LogPayload) {
	err := h.services.LogService.Log(ctx, "INFO", logPayload.Name, map[string]interface{}{
		"data": logPayload.Data,
	})
	if err != nil {
		helper.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	response := types.JsonResponse{
		Error:   false,
		Message: "logged",
	}
	
	helper.WriteJSON(w, http.StatusAccepted, response)
}

func (h *Handler) logEventViaRabbit(ctx context.Context, w http.ResponseWriter, logPayload types.LogPayload) {
	err := h.services.RabbitService.PublishLog(ctx, logPayload)
	if err != nil {
		helper.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	response := types.JsonResponse{
		Error:   false,
		Message: "logged via RabbitMQ",
	}
	
	helper.WriteJSON(w, http.StatusAccepted, response)
}

func (h *Handler) sendMail(ctx context.Context, w http.ResponseWriter, mailPayload types.MailPayload) {
	err := h.services.MailService.SendEmail(ctx, mailPayload.To, mailPayload.Subject, mailPayload.Message)
	if err != nil {
		helper.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	response := types.JsonResponse{
		Error:   false,
		Message: "Message sent to " + mailPayload.To,
	}
	
	helper.WriteJSON(w, http.StatusAccepted, response)
}