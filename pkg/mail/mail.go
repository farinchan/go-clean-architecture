package mail

import (
	"crypto/tls"

	"github.com/your-username/go-clean-architecture/config"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"gopkg.in/gomail.v2"
)

// Mailer handles email sending
type Mailer struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
}

// NewMailer creates a new mailer instance
func NewMailer(cfg *config.SMTPConfig) *Mailer {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Mailer{
		dialer:   dialer,
		from:     cfg.From,
		fromName: cfg.FromName,
	}
}

// EmailData holds email data
type EmailData struct {
	To          []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string
	CC          []string
	BCC         []string
}

// Send sends an email
func (m *Mailer) Send(data EmailData) error {
	msg := gomail.NewMessage()

	// Set sender
	msg.SetAddressHeader("From", m.from, m.fromName)

	// Set recipients
	msg.SetHeader("To", data.To...)

	// Set CC if provided
	if len(data.CC) > 0 {
		msg.SetHeader("Cc", data.CC...)
	}

	// Set BCC if provided
	if len(data.BCC) > 0 {
		msg.SetHeader("Bcc", data.BCC...)
	}

	// Set subject
	msg.SetHeader("Subject", data.Subject)

	// Set body
	if data.IsHTML {
		msg.SetBody("text/html", data.Body)
	} else {
		msg.SetBody("text/plain", data.Body)
	}

	// Add attachments
	for _, attachment := range data.Attachments {
		msg.Attach(attachment)
	}

	// Send email
	if err := m.dialer.DialAndSend(msg); err != nil {
		logger.Errorf("Failed to send email: %v", err)
		return err
	}

	logger.Infof("Email sent successfully to: %v", data.To)
	return nil
}

// SendSimple sends a simple text email
func (m *Mailer) SendSimple(to, subject, body string) error {
	return m.Send(EmailData{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}

// SendHTML sends an HTML email
func (m *Mailer) SendHTML(to, subject, htmlBody string) error {
	return m.Send(EmailData{
		To:      []string{to},
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	})
}
