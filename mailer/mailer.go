package mailer

import (
	"context"
	"io"

	"net/mail"
)

type Mailer interface {
	Send(ctx context.Context, msg *Message) error
}

type Message struct {
	From        mail.Address      `json:"from"`
	To          []mail.Address    `json:"to"`
	Bcc         []mail.Address    `json:"bcc"`
	Cc          []mail.Address    `json:"cc"`
	Subject     string            `json:"subject"`
	ReplyTo     string            `json:"reply_to"`
	HTML        string            `json:"html"`
	Text        string            `json:"text"`
	Headers     map[string]string `json:"headers"`
	Attachments []Attachment      `json:"attachments"`
}

type Attachment struct {
	FileName    string
	Description string
	Content     io.Reader
	ContentType string
	Inline      bool
}
