package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
	"time"

	"mailer/internal/config"
	"mailer/internal/templates"
	"mailer/types"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Service struct {
	config *config.MailerConfig
}

func NewService(cfg *config.MailerConfig) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) SendSMTPMessage(msg types.Message) error {
	if msg.From == "" {
		msg.From = s.config.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = s.config.FromName
	}

	if err := s.validateMessage(msg); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	data := map[string]any{
		"message": msg.Data,
	}
	msg.DataMap = data

	formattedMessage, err := s.buildHTMLMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to build HTML message: %w", err)
	}

	plainMessage, err := s.buildPlainTextMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to build plain text message: %w", err)
	}

	server := mail.NewSMTPClient()
	server.Host = s.config.Host
	server.Port = s.config.Port
	server.Username = s.config.Username
	server.Password = s.config.Password
	server.Encryption = s.getEncryption(s.config.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			if err := s.validateAttachment(attachment); err != nil {
				log.Printf("Skipping invalid attachment %s: %v", attachment, err)
				continue
			}
			email.AddAttachment(attachment)
		}
	}

	if err := email.Send(smtpClient); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", msg.To)
	return nil
}

func (s *Service) validateMessage(msg types.Message) error {
	var errors []string

	if msg.To == "" {
		errors = append(errors, "recipient email is required")
	}
	if msg.Subject == "" {
		errors = append(errors, "subject is required")
	}
	if !s.isValidEmail(msg.To) {
		errors = append(errors, "invalid recipient email format")
	}
	if msg.From != "" && !s.isValidEmail(msg.From) {
		errors = append(errors, "invalid sender email format",msg.From)
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (s *Service) isValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	return len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (s *Service) validateAttachment(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("attachment path must be absolute")
	}
	
	return nil
}

func (s *Service) buildHTMLMessage(msg types.Message) (string, error) {
	
	t, err := template.New("email-html").ParseFS(templates.EmailTemplates, "mail.html.gohtml")
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	formattedMessage := tpl.String()
	
	formattedMessage, err = s.inlineCSS(formattedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to inline CSS: %w", err)
	}

	return formattedMessage, nil
}

func (s *Service) buildPlainTextMessage(msg types.Message) (string, error) {

	t, err := template.New("email-html").ParseFS(templates.EmailPlainTemplates, "mail.plain.gohtml")
	if err != nil {
		return "", fmt.Errorf("failed to parse plain text template: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", fmt.Errorf("failed to execute plain text template: %w", err)
	}

	return tpl.String(), nil
}

func (s *Service) inlineCSS(htmlContent string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(htmlContent, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (s *Service) getEncryption(encType string) mail.Encryption {
	switch strings.ToLower(encType) {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}