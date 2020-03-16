package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/parfio/mail-notifier/internal/mailer"
	"github.com/sirupsen/logrus"
	"log"
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

	m, err := mailer.New(cfg.Mail.TemplatesFolderPath,
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

	err = m.SendPackageArrivedEMail("max.marche@live.de", "Max MMARRCHHCE")
	if err != nil {
		log.Fatal(err)
	}
}
