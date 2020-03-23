package internal

import (
	"context"
	"fmt"
	"github.com/parfy-io/mail-notifier/internal/mqtt"
	"github.com/parfy-io/mail-notifier/internal/users"
	"github.com/sirupsen/logrus"
)

type MQTTClient interface {
	Notifications() (<-chan mqtt.Notification, error)
}

type Mailer interface {
	SendPackageArrivedEMail(recipient, name string) error
}

type UsersClient interface {
	ByUserID(clientID, userID string) (users.User, error)
}

type Notifier struct {
	MQTTClient  MQTTClient
	Mailer      Mailer
	UsersClient UsersClient
}

func (n Notifier) Run(ctx context.Context) <-chan error {
	errs := make(chan error)
	go func() {
		defer close(errs)
		notifications, err := n.MQTTClient.Notifications()
		if err != nil {
			errs <- fmt.Errorf("failed to consume mqtt notifications: %w", err)
			return
		}

		for {
			select {
			case notification := <-notifications:
				u, err := n.UsersClient.ByUserID(notification.ClientID, notification.UserID)
				if err != nil {
					logrus.WithError(err).WithFields(map[string]interface{}{
						"clientID":      notification.ClientID,
						"userID":        notification.UserID,
						"correlationID": notification.CorrelationID,
					}).Error("could not get user")
					continue
				}

				err = n.Mailer.SendPackageArrivedEMail(u.EMail, u.FullName)
				if err != nil {
					logrus.WithError(err).WithFields(map[string]interface{}{
						"clientID":      notification.ClientID,
						"userID":        notification.UserID,
						"correlationID": notification.CorrelationID,
					}).Error("could not send email")
					continue
				}
			case _, ok := <-ctx.Done():
				if ok {
					errs <- ctx.Err()
				}

				return
			}
		}
	}()

	return errs
}
