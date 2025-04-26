package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// smtpMailer implements the Mailer interface for SMTP protocol.
type smtpMailer struct {
	host     string
	port     int
	username string
	password string
}

// Send sends an email using SMTP protocol with context support.
func (s *smtpMailer) Send(ctx context.Context, message *EmailMessage) error {
	// Build email headers
	headers := make(map[string]string)
	headers["From"] = message.From
	headers["To"] = strings.Join(message.To, ",")
	if len(message.Cc) > 0 {
		headers["Cc"] = strings.Join(message.Cc, ",")
	}
	headers["Subject"] = message.Subject
	headers["MIME-Version"] = "1.0"

	var contentType string
	if message.IsHTML {
		contentType = "text/html; charset=UTF-8"
	} else {
		contentType = "text/plain; charset=UTF-8"
	}

	boundary := "my-boundary-123456"
	if len(message.Files) > 0 {
		headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)
	} else {
		headers["Content-Type"] = contentType
	}

	// Compose message
	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	if len(message.Files) > 0 {
		// Add body as the first part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString(fmt.Sprintf("Content-Type: %s\r\n\r\n", contentType))
		msg.WriteString(message.Body)
		msg.WriteString("\r\n")

		// Add attachments
		for _, file := range message.Files {
			msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			msg.WriteString(fmt.Sprintf("Content-Type: %s\r\n", file.ContentType))
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", file.Filename))
			msg.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

			b := make([]byte, base64.StdEncoding.EncodedLen(len(file.Content)))
			base64.StdEncoding.Encode(b, file.Content)
			msg.Write(b)
			msg.WriteString("\r\n")
		}

		// Close boundary
		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// No attachment
		msg.WriteString(message.Body)
	}

	// Prepare SMTP authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	serverAddr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Establish TLS connection
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // WARNING: Do proper validation in production
		ServerName:         s.host,
	}

	rawDialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	conn, err := tls.DialWithDialer(rawDialer, "tcp", serverAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Set sender and recipients
	if err = client.Mail(message.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	recipients := append(append(message.To, message.Cc...), message.Bcc...)
	for _, addr := range recipients {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to add recipient %s: %w", addr, err)
		}
	}

	// Start sending data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data command: %w", err)
	}
	_, err = writer.Write(msg.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
