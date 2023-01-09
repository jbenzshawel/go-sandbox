package infrastructure

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

type Email struct {
	To      string
	Subject string
	Body    string
	IsHtml  bool
}

type EmailSender interface {
	SendMail(email *Email) error
}

type EmailClient struct {
	smtpAddr string
	host     string
	from     string
}

func NewEmailClient(smtpAddr, host, from string) *EmailClient {
	return &EmailClient{smtpAddr: smtpAddr, host: host, from: from}
}

const htmlMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func (c *EmailClient) SendMail(email *Email) error {
	auth := PlainAuth("", "", "", c.host)

	to := []string{email.To}
	msg := []byte(fmt.Sprintf("To: %s\n", email.To) +
		fmt.Sprintf("Subject: %s\n", email.Subject))
	if email.IsHtml {
		msg = append(msg, []byte(htmlMIME)...)
	} else {
		msg = append(msg, []byte("\r\n")...)
	}
	msg = append(msg, []byte(email.Body)...)
	err := smtp.SendMail(c.smtpAddr, auth, c.from, to, msg)
	return err
}

type plainAuth struct {
	identity, username, password string
	host                         string
}

// PlainAuth is the same as smtp.PlainAuth in the standard library, however it has
// been updated to treat the hostname used with docker-compose as localhost to allow
// testing locally
func PlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuth{identity, username, password, host}
}

func isLocalhost(name string) bool {
	return name == "localhost" ||
		name == "127.0.0.1" ||
		name == "::1" ||
		(name == os.Getenv("SMTP_HOST") && os.Getenv("SERVER_ENVIRONMENT") == "docker-compose")
}

func (a *plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises PLAIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}
