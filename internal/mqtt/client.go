package mqtt

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt/client"
)

type Client struct {
	c *client.Client
	n chan Notification
}

func NewClient(brokerAddress, clientID, username, passowrd string) (*Client, error, <-chan error) {
	errs := make(chan error)
	mqttClient := client.New(&client.Options{ErrorHandler: func(err error) {
		errs <- fmt.Errorf("mqtt client failed: %w", err)
	}})

	err := mqttClient.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      brokerAddress,
		ClientID:     []byte(clientID),
		UserName:     []byte(username),
		Password:     []byte(passowrd),
		CleanSession: false,
	})
	if err != nil {
		close(errs)
		return nil, fmt.Errorf("failed to connect to broker: %w", err), nil
	}

	return &Client{c: mqttClient, n: make(chan Notification)}, nil, errs
}

func (c Client) Stop() error {
	close(c.n)
	err := c.c.Unsubscribe(&client.UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte(notificationTopic),
		},
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to unsubscribe from notifications-topic")
	}

	err = c.c.Disconnect()
	c.c.Terminate()

	if err != nil {
		return fmt.Errorf("failed to disconnect from broker: %w", err)
	}

	return nil
}
