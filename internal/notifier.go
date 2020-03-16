package internal

//go:generate moq -out mailer_moq_test.go . Mailer
type Mailer interface {
	SendPackageArrivedEMail(recipient, name string) error
}

type Notifier struct {
	Mailer Mailer
}

func Run() error {

}
