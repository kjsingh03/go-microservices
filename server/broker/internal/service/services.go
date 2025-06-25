package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-broker/types"
	"time"
)

// Services holds all service implementations
type Services struct {
	AuthService   AuthService
	LogService    LogService
	MailService   MailService
	RabbitService RabbitService
}

// Close closes all services
func (s *Services) Close() error {
	if s.RabbitService != nil {
		return s.RabbitService.Close()
	}
	return nil
}

// Auth Service Implementation
type authService struct {
	baseURL string
	timeout time.Duration
	retries int
	apiKey  string
	client  *http.Client
}

func NewAuthService(baseURL string, timeout time.Duration, retries int, apiKey string) AuthService {
	return &authService{
		baseURL: baseURL,
		timeout: timeout,
		retries: retries,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *authService) Authenticate(ctx context.Context, authPayload types.AuthPayload) (*types.AuthResponse, error) {
	jsonData, err := json.Marshal(authPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid credentials")
	} else if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var authResp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &authResp, nil
}

func (s *authService) ValidatePermissions(ctx context.Context, userID string, resource string) error {
	// Implement permission validation logic based on your requirements
	return nil
}

// Log Service Implementation
type logService struct {
	baseURL string
	timeout time.Duration
	retries int
	client  *http.Client
}

func NewLogService(baseURL string, timeout time.Duration, retries int) LogService {
	return &logService{
		baseURL: baseURL,
		timeout: timeout,
		retries: retries,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *logService) Log(ctx context.Context, level string, message string, data map[string]interface{}) error {
	logEntry := types.LogPayload{
		Name: message,
		Data: fmt.Sprintf("%v", data),
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/logs", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call log service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("log service returned status %d", resp.StatusCode)
	}

	return nil
}

// Mail Service Implementation
type mailService struct {
	baseURL string
	timeout time.Duration
	retries int
	apiKey  string
	client  *http.Client
}

func NewMailService(baseURL string, timeout time.Duration, retries int, apiKey string) MailService {
	return &mailService{
		baseURL: baseURL,
		timeout: timeout,
		retries: retries,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *mailService) SendEmail(ctx context.Context, to, subject, body string) error {
	mailPayload := types.MailPayload{
		To:      to,
		Subject: subject,
		Message: body,
	}

	jsonData, err := json.Marshal(mailPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal mail payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call mail service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mail service returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *mailService) SendTemplateEmail(ctx context.Context, to, templateID, subject string, data map[string]interface{}) error {
	// Implement template email sending logic based on your mail service
	return s.SendEmail(ctx, to, subject, fmt.Sprintf("Template: %s, Data: %v", templateID, data))
}
