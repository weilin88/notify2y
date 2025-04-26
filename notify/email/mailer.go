package mailer

import "context"

// Mailer defines the interface for sending emails.
type Mailer interface {
	Send(ctx context.Context, message *EmailMessage) error
}

// EmailMessage represents an email message structure.
type EmailMessage struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
	IsHTML  bool
	Files   []Attachment
}

// Attachment represents a file attachment.
type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

// NewSMTPMailer creates a new SMTP mailer instance.
func NewSMTPMailer(host string, port int, username, password string) Mailer {
	return &smtpMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}
