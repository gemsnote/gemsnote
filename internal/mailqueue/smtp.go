package mailqueue

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"
)

func message(id string, p Payload) ([]byte, string, []string, error) {
	from, err := mail.ParseAddress(p.SMTP.Username)
	if err != nil || strings.ContainsAny(p.SMTP.Username+p.Subject+id, "\r\n") || len(p.To) == 0 {
		return nil, "", nil, errors.New("invalidEmail")
	}
	to := make([]string, 0, len(p.To))
	headers := make([]string, 0, len(p.To))
	for _, value := range p.To {
		a, err := mail.ParseAddress(value)
		if err != nil || strings.ContainsAny(value, "\r\n") {
			return nil, "", nil, errors.New("invalidEmail")
		}
		to = append(to, a.Address)
		headers = append(headers, a.String())
	}
	var body bytes.Buffer
	multi := multipart.NewWriter(&body)
	plain := p.Text
	if plain == "" {
		plain = "Please view this message in an HTML-capable email client."
	}
	for _, part := range []struct{ kind, body string }{{"text/plain", plain}, {"text/html", p.Body}} {
		w, err := multi.CreatePart(textproto.MIMEHeader{"Content-Type": {part.kind + "; charset=utf-8"}, "Content-Transfer-Encoding": {"quoted-printable"}})
		if err != nil {
			return nil, "", nil, err
		}
		q := quotedprintable.NewWriter(w)
		if _, err = q.Write([]byte(part.body)); err != nil {
			return nil, "", nil, err
		}
		if err = q.Close(); err != nil {
			return nil, "", nil, err
		}
	}
	if err = multi.Close(); err != nil {
		return nil, "", nil, err
	}
	var result bytes.Buffer
	fmt.Fprintf(&result, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMessage-ID: <%s@gemsnote>\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%q\r\n\r\n", from.String(), strings.Join(headers, ", "), mime.QEncoding.Encode("utf-8", p.Subject), time.Now().UTC().Format(time.RFC1123Z), id, multi.Boundary())
	result.Write(body.Bytes())
	return result.Bytes(), from.Address, to, nil
}

// Send bounds the entire SMTP conversation, including TLS, authentication and DATA.
func Send(ctx context.Context, id string, p Payload) error {
	data, from, to, err := message(id, p)
	if err != nil {
		return errors.New("mailDeliveryFailed")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", net.JoinHostPort(p.SMTP.Host, strconv.Itoa(p.SMTP.Port)))
	if err != nil {
		return errors.New("mailConnectionFailed")
	}
	defer conn.Close()
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	deadline, _ := ctx.Deadline()
	conn.SetDeadline(deadline)
	tlsConfig, err := smtpTLSConfig(p.SMTP.Host)
	if err != nil {
		return errors.New("mailConnectionFailed")
	}
	var transport net.Conn = conn
	if p.SMTP.SSL {
		secured := tls.Client(conn, tlsConfig)
		if err = secured.HandshakeContext(ctx); err != nil {
			return errors.New("mailConnectionFailed")
		}
		transport = secured
	}
	client, err := smtp.NewClient(transport, p.SMTP.Host)
	if err != nil {
		return errors.New("mailConnectionFailed")
	}
	defer client.Close()
	if !p.SMTP.SSL {
		if err = client.StartTLS(tlsConfig); err != nil {
			return errors.New("mailConnectionFailed")
		}
	}
	if err = client.Auth(smtp.PlainAuth("", p.SMTP.Username, p.SMTP.Password, p.SMTP.Host)); err != nil {
		return smtpError(err, "mailAuthenticationFailed")
	}
	if err = client.Mail(from); err != nil {
		return smtpError(err, "mailDeliveryFailed")
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return smtpError(err, "mailRecipientRejected")
		}
	}
	w, err := client.Data()
	if err != nil {
		return smtpError(err, "mailDeliveryFailed")
	}
	if _, err = w.Write(data); err != nil {
		return smtpError(err, "mailDeliveryFailed")
	}
	if err = w.Close(); err != nil {
		return smtpError(err, "mailDeliveryFailed")
	}
	// DATA was accepted. A subsequent QUIT failure must not cause duplicate delivery.
	client.Quit()
	return nil
}
func smtpError(err error, fallback string) error {
	var network net.Error
	if errors.As(err, &network) {
		return errors.New("mailConnectionFailed")
	}
	return errors.New(fallback)
}

// Honor the existing custom-CA setting on all worker platforms, including Windows.
func smtpTLSConfig(host string) (*tls.Config, error) {
	config := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
	if path := os.Getenv("SSL_CERT_FILE"); path != "" {
		pem, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New("mailConnectionFailed")
		}
		roots, err := x509.SystemCertPool()
		if err != nil {
			roots = x509.NewCertPool()
		}
		if !roots.AppendCertsFromPEM(pem) {
			return nil, errors.New("mailConnectionFailed")
		}
		config.RootCAs = roots
	}
	return config, nil
}
