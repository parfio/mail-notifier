package mail

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/mail.v2"
)

type Mailer struct {
	dialer    *mail.Dialer
	templates map[string]Template
}

func New(templatesFolderPath, username, password, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) (*Mailer, error) {
	d := mail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify, ServerName: tlsServerName}

	//check connection and auth
	rd, err := d.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	defer rd.Close()

	pwRestTmpl, err := Load(templatesFolderPath, PackageArrivedTemplateName)
	if err != nil {
		return nil, fmt.Errorf("failed to load password-reset tempate: %w", err)
	}

	return &Mailer{
		dialer: d,
		templates: map[string]Template{
			PackageArrivedTemplateName: pwRestTmpl,
		},
	}, nil
}

func (m Mailer) SendPackageArrivedEMail(recipient, name string) error {
	mailData := struct {
		EMail string
		Name  string
	}{
		EMail: recipient,
		Name:  name,
	}

	msg, err := m.templates[PackageArrivedTemplateName].Render(mailData)
	if err != nil {
		return fmt.Errorf("failed to render mail template: %w", err)
	}

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err)
	}

	return nil
}
