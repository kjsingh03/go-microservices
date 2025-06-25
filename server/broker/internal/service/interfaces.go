package service

import (
	"context"
	"service-broker/types"
)

// Service interfaces
type AuthService interface {
	Authenticate(ctx context.Context, authPayload types.AuthPayload) (*types.AuthResponse, error)
	ValidatePermissions(ctx context.Context, userID string, resource string) error
}

type LogService interface {
	Log(ctx context.Context, level string, message string, data map[string]interface{}) error
}

type MailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendTemplateEmail(ctx context.Context, to, templateID, subject string, data map[string]interface{}) error
}

type RabbitService interface {
	PublishLog(ctx context.Context, payload types.LogPayload) error
	Close() error
}