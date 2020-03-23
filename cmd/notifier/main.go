package main

import (
	"context"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/parfy-io/mail-notifier/internal"
	"github.com/parfy-io/mail-notifier/internal/mail"
	"github.com/parfy-io/mail-notifier/internal/mqtt"
	"github.com/parfy-io/mail-notifier/internal/users"
	"github.com/parfy-io/mail-notifier/internal/web"
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

	mailer, err := mail.New(cfg.Mail.TemplatesFolderPath,
		cfg.Mail.SMTPUsername,
		cfg.Mail.SMTPPassword,
		cfg.Mail.SMTPHost,
		cfg.Mail.SMTPPort,
		cfg.Mail.TLS.InsecureSkipVerify,
		cfg.Mail.TLS.ServerName,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create mail")
	}

	mqttClient, err, mqttErrs := mqtt.NewClient(cfg.MQTT.BrokerAddress, cfg.MQTT.ClientID, cfg.MQTT.Username, cfg.MQTT.Passowrd)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create mqtt-client")
	}

	usersClient := users.NewClient(cfg.UsersServiceBaseURL)

	webErrs := make(chan error)
	go func() {
		webErrs <- web.StartAliveEndpoint(cfg.ServerAddress)
	}()

	ctx, cancel := context.WithCancel(context.Background())

	errs := internal.Notifier{
		Mailer:      mailer,
		MQTTClient:  mqttClient,
		UsersClient: usersClient,
	}.Run(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	hasErr := false
	select {
	case sig := <-signals:
		logrus.WithField("signal", sig).Info("Notifier interrupted")
	case err := <-webErrs:
		hasErr = true
		logrus.WithError(err).Error("Web-Server could not continue processing")
	case err := <-mqttErrs:
		hasErr = true
		logrus.WithError(err).Error("MQTT-client could not continue processing")
	case err := <-errs:
		hasErr = true
		logrus.WithError(err).Error("Notifier could not continue processing")
	}

	cancel()

	err = mqttClient.Stop()
	if err != nil {
		logrus.WithError(err).Error("Could not stop mqtt-client")
	}

	if hasErr {
		os.Exit(1)
	}
}
