package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	helpers "mailer/internal/helper"
	"mailer/internal/mailer"
	"mailer/types"
)

type Handler struct {
	mailerService *mailer.Service
}

func NewHandler(mailerService *mailer.Service) *Handler {
	return &Handler{
		mailerService: mailerService,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	payload := types.JsonResponse{
		Success: true,
		Message: "Welcome to Mailer Service API",
		Data: map[string]any{
			"version":   "1.0.0",
			"status":    "healthy",
			"endpoints": []string{"/api/v1/send","/api/v1/send/batch"},
		},
	}

	if err := helpers.WriteJSON(w, http.StatusOK, payload); err != nil {
		log.Printf("Error writing JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailRequest struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var req mailRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		log.Printf("Error reading JSON: %v", err)
		helpers.ErrorJSON(w, fmt.Errorf("invalid JSON payload"), http.StatusBadRequest)
		return
	}

	msg := types.Message{
		From:    strings.TrimSpace(req.From),
		To:      strings.TrimSpace(req.To),
		Subject: strings.TrimSpace(req.Subject),
		Data:    req.Message,
	}

	if err := h.mailerService.SendSMTPMessage(msg); err != nil {
		log.Printf("Error sending email: %v", err)
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := types.JsonResponse{
		Success: true,
		Message: fmt.Sprintf("Email sent successfully to %s", req.To),
		Data: map[string]string{
			"recipient": req.To,
			"status":    "sent",
		},
	}

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *Handler) SendBatchMail(w http.ResponseWriter, r *http.Request) {
	type batchMailRequest struct {
		From      string   `json:"from"`
		To        []string `json:"to"`
		Subject   string   `json:"subject"`
		Message   string   `json:"message"`
		BatchSize int      `json:"batch_size,omitempty"`
	}

	var req batchMailRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid JSON payload"), http.StatusBadRequest)
		return
	}

	if len(req.To) == 0 {
		helpers.ErrorJSON(w, fmt.Errorf("at least one recipient is required"), http.StatusBadRequest)
		return
	}

	if req.BatchSize == 0 {
		req.BatchSize = 10
	}

	// Need to implment RabbitMQ
	var failed []string
	var sent []string

	for _, recipient := range req.To {
		msg := types.Message{
			From:    req.From,
			To:      strings.TrimSpace(recipient),
			Subject: req.Subject,
			Data:    req.Message,
		}

		if err := h.mailerService.SendSMTPMessage(msg); err != nil {
			log.Printf("Failed to send to %s: %v", recipient, err)
			failed = append(failed, recipient)
		} else {
			sent = append(sent, recipient)
		}
	}

	payload := types.JsonResponse{
		Success: len(failed) == 0,
		Message: fmt.Sprintf("Sent to %d recipients, %d failed", len(sent), len(failed)),
		Data: map[string]interface{}{
			"sent":        sent,
			"failed":      failed,
			"total_sent":  len(sent),
			"total_failed": len(failed),
		},
	}

	status := http.StatusOK
	if len(failed) > 0 {
		status = http.StatusPartialContent
	}

	helpers.WriteJSON(w, status, payload)
}