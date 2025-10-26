package email

import (
	"bytes"
	"embed"
	"github.com/wneessen/go-mail"
	ht "html/template"
	tt "text/template"
	"time"
)

//go:embed "templates"
var templateFS embed.FS

type MailSender interface {
	Send(recipient, templateFile string, data any) error
}

type mailSender struct {
	client   *mail.Client
	host     string
	port     int
	username string
	password string
	sender   string
}

func (m *mailSender) Send(recipient, templateFile string, data any) error {
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := textTmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	if err := textTmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
		return err
	}

	htmTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	if err := htmTmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
		return err
	}

	msg := mail.NewMsg()
	if err := msg.To(recipient); err != nil {
		return err
	}
	if err := msg.From(m.sender); err != nil {
		return err
	}
	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}

		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}

type Options func(*mailSender)

func WithClient(client *mail.Client) Options {
	return func(s *mailSender) {
		s.client = client
	}
}

func WithHost(host string) Options {
	return func(s *mailSender) {
		s.host = host
	}
}

func WithPort(port int) Options {
	return func(s *mailSender) {
		s.port = port
	}
}

func WithUsername(username string) Options {
	return func(s *mailSender) {
		s.username = username
	}
}

func WithPassword(password string) Options {
	return func(s *mailSender) {
		s.password = password
	}
}

func WithSender(sender string) Options {
	return func(s *mailSender) {
		s.sender = sender
	}
}

func NewMailSender(options ...Options) MailSender {
	s := &mailSender{}
	for _, option := range options {
		option(s)
	}
	return s
}
