package types

import (
	"fmt"
	"strings"
)

var ValidActions = []string{"auth", "log", "logdirect", "mail"}

func (r *RequestPayload) Validate() error {
	// Check if action is provided
	if r.Action == "" {
		return fmt.Errorf("action is required. Valid actions: %s", strings.Join(ValidActions, ", "))
	}

	// Check if action is valid
	validAction := false
	for _, action := range ValidActions {
		if r.Action == action {
			validAction = true
			break
		}
	}

	if !validAction {
		return fmt.Errorf("invalid action '%s'. Valid actions: %s", r.Action, strings.Join(ValidActions, ", "))
	}

	// Validate payload based on action
	switch r.Action {
	case "auth":
		if r.Auth == nil {
			return fmt.Errorf("auth payload is required for action 'auth'")
		}
		return r.Auth.Validate()
	case "log", "logdirect":
		if r.Log == nil {
			return fmt.Errorf("log payload is required for action '%s'", r.Action)
		}
		return r.Log.Validate()
	case "mail":
		if r.Mail == nil {
			return fmt.Errorf("mail payload is required for action 'mail'")
		}
		return r.Mail.Validate()
	}

	return nil
}

func (a *AuthPayload) Validate() error {
	if a.Email == "" {
		return fmt.Errorf("email is required for authentication")
	}
	if a.Password == "" {
		return fmt.Errorf("password is required for authentication")
	}
	if !strings.Contains(a.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (l *LogPayload) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("log name is required")
	}
	if l.Data == "" {
		return fmt.Errorf("log data is required")
	}
	return nil
}

func (m *MailPayload) Validate() error {
	if m.To == "" {
		return fmt.Errorf("recipient email (to) is required")
	}
	if m.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if m.Message == "" {
		return fmt.Errorf("email message is required")
	}
	if !strings.Contains(m.To, "@") {
		return fmt.Errorf("invalid recipient email format")
	}
	if m.From != "" && !strings.Contains(m.From, "@") {
		return fmt.Errorf("invalid sender email format")
	}
	return nil
}
