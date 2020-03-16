package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
)

type config struct {
	ServerAddress       string `conf:"help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080',default:0.0.0.0:80"`
	UsersServiceBaseURL string `conf:"help:Base URL to users service,required"`
	Mail                struct {
		TemplatesFolderPath string `conf:"help:Path to mail-templates folder,default:/mail-templates"`
		SMTPUsername        string `conf:"env:MAIL_SMTP_USERNAME,help:SMTP username to authorize with,required"`
		SMTPPassword        string `conf:"env:MAIL_SMTP_PASSWORD,help:SMTP password to authorize with,required,noprint"`
		SMTPHost            string `conf:"env:MAIL_SMTP_HOST,help:SMTP host to connect to,required"`
		SMTPPort            int    `conf:"env:MAIL_SMTP_PORT,help:SMTP port to connect to,default:587"`
		TLS                 struct {
			InsecureSkipVerify bool   `conf:"help:true if certificates should not be verified,default:false"`
			ServerName         string `conf:"help:name of the server who expose the certificate"`
		}
	}
	MQTT struct {
		BrokerAddress string `conf:"help:MQTT Broker address,required"`
		ClientID      string `conf:"help:MQTT ClientID,required"`
		Username      string `conf:"help:MQTT Username to authorize"`
		Passowrd      string `conf:"help:MQTT Password to authorize,noprint"`
	}
}

func newConfig() (config, error) {
	cfg := config{}

	if origErr := conf.Parse(os.Environ(), "MMN", &cfg); origErr != nil {
		usage, err := conf.Usage("MMN", &cfg)
		if err != nil {
			return cfg, err
		}
		fmt.Println(usage)
		return cfg, origErr
	}

	return cfg, nil
}
