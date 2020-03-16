package internal

//go:generate moq -out mailer_moq_test.go . Mailer
type Mailer interface {
	SendPasswordResetRequestEMail(recipient, passwordResetLink string) error
}

type Notifier struct {
	Mailer Mailer
}
