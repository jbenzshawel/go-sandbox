package infrastructure

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jbenzshawel/go-sandbox/common/rest"
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
	username string
	password string
}

type EmailConfig struct {
	Addr     string
	Host     string
	From     string
	Username string
	Password string
}

func NewEmailClient(cfg EmailConfig) *EmailClient {
	return &EmailClient{
		smtpAddr: cfg.Addr,
		host:     cfg.Host,
		from:     cfg.Password,
		username: cfg.Username,
		password: cfg.Password,
	}
}

const htmlMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func (c *EmailClient) SendMail(email *Email) error {
	auth := c.plainAuth()

	to := []string{email.To}
	msg := []byte(fmt.Sprintf("To: %s\n", email.To) +
		fmt.Sprintf("Subject: %s\n", email.Subject))
	if email.IsHtml {
		msg = append(msg, []byte(htmlMIME)...)
	} else {
		msg = append(msg, []byte("\r\n")...)
	}
	msg = append(msg, []byte(email.Body)...)
	return smtp.SendMail(c.smtpAddr, auth, c.from, to, msg)
}

func (c *EmailClient) plainAuth() smtp.Auth {
	return PlainAuth("", c.username, c.password, c.host)
}

func (c *EmailClient) HealthCheck() rest.HealthCheckTask {
	return func() (bool, string, error) {
		healthCheckName := "smtp"
		ec, err := smtp.Dial(c.smtpAddr)
		if err != nil {
			return false, healthCheckName, err
		}

		err = ec.Auth(c.plainAuth())
		if err != nil {
			return false, healthCheckName, err
		}

		err = ec.Noop()
		if err != nil {
			return false, healthCheckName, err
		}

		return true, healthCheckName, nil
	}
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
