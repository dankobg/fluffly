package mailer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/mail"
	"net/textproto"

	gomail "github.com/wneessen/go-mail"
)

var _ Mailer = (*SmtpClient)(nil)

type SmtpAuthType string

const (
	SmtpAuthPlain   SmtpAuthType = "PLAIN"
	SmtpAuthLogin   SmtpAuthType = "LOGIN"
	SmtpAuthCramMD5 SmtpAuthType = "CRAM-MD5"
	SmtpXOauth2     SmtpAuthType = "XOAUTH2"
)

type mailerOpts struct {
	devHost     string
	devPort     int
	devUsername string
	devPassword string
	enabled     bool
	host        string
	port        int
	username    string
	password    string
	tLS         bool
	authMethod  SmtpAuthType
	fromName    string
	fromAddress string
	log         *slog.Logger
}

type SmtpClientOption interface {
	apply(*mailerOpts)
}

type SmtpClientOptions []SmtpClientOption

func (o SmtpClientOptions) apply(s *mailerOpts) {
	for _, opt := range o {
		opt.apply(s)
	}
}

type enabledOpt bool

func (o enabledOpt) apply(c *mailerOpts)     { c.enabled = bool(o) }
func WithEnabled(flag bool) SmtpClientOption { return enabledOpt(flag) }

type devHostOpt string

func (o devHostOpt) apply(c *mailerOpts)       { c.host = string(o) }
func WithDevHost(host string) SmtpClientOption { return devHostOpt(host) }

type devPortOpt int

func (o devPortOpt) apply(c *mailerOpts)    { c.port = int(o) }
func WithDevPort(port int) SmtpClientOption { return devPortOpt(port) }

type devUsernameOpt string

func (o devUsernameOpt) apply(c *mailerOpts)           { c.username = string(o) }
func WithDevUsername(username string) SmtpClientOption { return devUsernameOpt(username) }

type devPasswordOpt string

func (o devPasswordOpt) apply(c *mailerOpts)      { c.password = string(o) }
func WithDevPassword(pwd string) SmtpClientOption { return devPasswordOpt(pwd) }

type hostOpt string

func (o hostOpt) apply(c *mailerOpts)       { c.host = string(o) }
func WithHost(host string) SmtpClientOption { return hostOpt(host) }

type portOpt int

func (o portOpt) apply(c *mailerOpts)    { c.port = int(o) }
func WithPort(port int) SmtpClientOption { return portOpt(port) }

type usernameOpt string

func (o usernameOpt) apply(c *mailerOpts)           { c.username = string(o) }
func WithUsername(username string) SmtpClientOption { return usernameOpt(username) }

type passwordOpt string

func (o passwordOpt) apply(c *mailerOpts)      { c.password = string(o) }
func WithPassword(pwd string) SmtpClientOption { return passwordOpt(pwd) }

type tlsOpt bool

func (o tlsOpt) apply(c *mailerOpts)     { c.tLS = bool(o) }
func WithTLS(flag bool) SmtpClientOption { return tlsOpt(flag) }

type authMethodOpt SmtpAuthType

func (o authMethodOpt) apply(c *mailerOpts)             { c.authMethod = SmtpAuthType(o) }
func WithAuthMethod(auth SmtpAuthType) SmtpClientOption { return authMethodOpt(auth) }

type fromNameOpt string

func (o fromNameOpt) apply(c *mailerOpts)       { c.fromName = string(o) }
func WithFromName(name string) SmtpClientOption { return fromNameOpt(name) }

type fromAddressOpt string

func (o fromAddressOpt) apply(c *mailerOpts)          { c.fromAddress = string(o) }
func WithFromAddress(address string) SmtpClientOption { return fromAddressOpt(address) }

type logOpt struct{ log *slog.Logger }

func (o logOpt) apply(c *mailerOpts)            { c.log = o.log }
func WithLog(log *slog.Logger) SmtpClientOption { return logOpt{log: log} }

type SmtpClient struct {
	Host        string
	Port        int
	Username    string
	Password    string
	TLS         bool
	AuthMethod  SmtpAuthType
	FromName    string
	FromAddress string
	Log         *slog.Logger
}

func NewSmtpClient(opts ...SmtpClientOption) *SmtpClient {
	mopts := &mailerOpts{
		authMethod:  SmtpAuthPlain,
		fromName:    "Fluffly",
		devHost:     "mailpit",
		devPort:     1025,
		devUsername: "test",
		devPassword: "test",
	}

	for _, o := range opts {
		o.apply(mopts)
	}

	client := &SmtpClient{
		Host:        mopts.devHost,
		Port:        mopts.devPort,
		Username:    mopts.devUsername,
		Password:    mopts.devPassword,
		TLS:         mopts.tLS,
		AuthMethod:  mopts.authMethod,
		FromName:    mopts.fromName,
		FromAddress: mopts.fromAddress,
		Log:         mopts.log,
	}

	if mopts.enabled {
		client.Host = mopts.host
		client.Port = mopts.port
		client.Username = mopts.username
		client.Password = mopts.password
	}

	return client
}

func getAuthType(authType SmtpAuthType) gomail.SMTPAuthType {
	switch authType {
	case SmtpAuthPlain:
		return gomail.SMTPAuthPlain
	case SmtpAuthLogin:
		return gomail.SMTPAuthLogin
	case SmtpAuthCramMD5:
		return gomail.SMTPAuthCramMD5
	case SmtpXOauth2:
		return gomail.SMTPAuthXOAUTH2
	default:
		return gomail.SMTPAuthPlain
	}
}

func formatAddresses(addresses []mail.Address, withName bool) []string {
	result := make([]string, len(addresses))

	for i, addr := range addresses {
		if withName && addr.Name != "" {
			result[i] = addr.String()
		} else {
			result[i] = addr.Address
		}
	}

	return result
}

func (client *SmtpClient) Send(ctx context.Context, msg *Message) error {
	if client.Log == nil {
		client.Log = slog.Default()
	}

	opts := []gomail.Option{
		gomail.WithPort(client.Port),
		gomail.WithUsername(client.Username),
		gomail.WithPassword(client.Password),
		gomail.WithSMTPAuth(getAuthType(client.AuthMethod)),
		gomail.WithTLSPortPolicy(gomail.TLSMandatory),
		gomail.WithSSLPort(true),
	}

	c, err := gomail.NewClient(client.Host, opts...)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	m := gomail.NewMsg()

	if err := m.From(msg.From.Address); err != nil {
		client.Log.Error("invalid `from` header: %w", slog.Any("error", err))
		return err
	}
	if err := m.To(formatAddresses(msg.To, true)...); err != nil {
		client.Log.Error("invalid `to` header: %w", slog.Any("error", err))
		return err
	}
	if err := m.Cc(formatAddresses(msg.Cc, true)...); err != nil {
		client.Log.Error("invalid `Cc` header: %w", slog.Any("error", err))
		return err
	}
	if err := m.Bcc(formatAddresses(msg.Bcc, true)...); err != nil {
		client.Log.Error("invalid `Bcc` header: %w", slog.Any("error", err))
		return err
	}

	m.Subject(msg.Subject)
	m.SetBodyString(gomail.TypeTextHTML, msg.HTML)
	m.AddAlternativeString(gomail.TypeTextPlain, msg.Text)
	m.SetGenHeader(gomail.HeaderReplyTo, msg.ReplyTo)

	for k, v := range msg.Headers {
		m.SetGenHeader(gomail.Header(k), v)
	}

	attachments := make([]*gomail.File, len(msg.Attachments))
	for i, att := range msg.Attachments {
		attachments[i] = &gomail.File{
			Name:        att.FileName,
			Desc:        att.Description,
			ContentType: gomail.ContentType(att.ContentType),
			Header:      textproto.MIMEHeader{},
			Enc:         gomail.EncodingQP,
			Writer: func(w io.Writer) (int64, error) {
				return io.Copy(w, att.Content)
			},
		}
	}
	if len(attachments) > 0 {
		m.SetAttachements(attachments)
	}

	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
