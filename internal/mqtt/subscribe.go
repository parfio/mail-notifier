package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt/client"
	"strings"
)

//topic: notify/mail/<clientID>/<correlationID>
const notificationTopic string = "notify/mail/+/+"

func (c Client) Notifications() (<-chan Notification, error) {
	notifications := make(chan Notification)
	err := c.c.Subscribe(&client.SubscribeOptions{SubReqs: []*client.SubReq{
		{
			QoS:         2,
			TopicFilter: []byte(notificationTopic),
			Handler: func(topicName, message []byte) {
				topicLevels := strings.Split(string(topicName), "/")
				if len(topicLevels) != 4 {
					logrus.WithField("topic", string(topicName)).Error("Could not parse clientID and correlationID out of topic")
				}

				rawMsg := struct {
					UserID string `json:"user-id"`
				}{}

				err := json.Unmarshal(message, &rawMsg)
				if err != nil {
					logrus.WithError(err).WithField("message", string(message)).Error("Could not unmarshal json")
				}

				notifications <- Notification{
					ClientID:      topicLevels[2],
					UserID:        rawMsg.UserID,
					CorrelationID: topicLevels[3],
				}
			},
		},
	}})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to notifications-topic: %w", err)
	}

	return notifications, nil
}

type Notification struct {
	ClientID      string
	UserID        string
	CorrelationID string
}
