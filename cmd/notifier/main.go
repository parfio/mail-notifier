package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/parfio/mail-notifier/internal/mailer"
	"github.com/parfio/mail-notifier/internal/web"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	cfgAsString, err := conf.String(&cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Could not build config string")
	}
	fmt.Print(cfgAsString)
	logrus.Infof("Starting notifier")

	_, err = mailer.New(cfg.Mail.TemplatesFolderPath,
		cfg.Mail.SMTPUsername,
		cfg.Mail.SMTPPassword,
		cfg.Mail.SMTPHost,
		cfg.Mail.SMTPPort,
		cfg.Mail.TLS.InsecureSkipVerify,
		cfg.Mail.TLS.ServerName,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create mailer")
	}

	errs := make(chan error)
	go func() {
		errs <- web.StartAliveEndpoint(cfg.ServerAddress)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signals:
		logrus.WithField("signal", sig).Info("Notifier interrupted")
		break
	case err := <-errs:
		logrus.WithError(err).Error("Notifier could not continue processing")
	}
}
