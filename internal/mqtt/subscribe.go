package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt/client"
)

//topic: notify/mail/<clientID>/<correlationID>
const notificationTopic string = "notify/mail/+/+"

func (c Client) Notifications(ctx context.Context) (chan<- Notification, error) {
	notifications := make(chan<- Notification)
	err := c.c.Subscribe(&client.SubscribeOptions{SubReqs: []*client.SubReq{
		{
			QoS:         2,
			TopicFilter: []byte(notificationTopic),
			Handler: func(topicName, message []byte) {
				var clientID, correlationID string
				_, err := fmt.Sscanf(string(topicName), "notify/mail/%s/%s", &clientID, &correlationID)
				if err != nil {
					logrus.WithError(err).WithField("topic", topicName).Error("Could not parse clientID, userID and correlationID out of topic")
					return
				}

				rawMsg := struct {
					UserID string `json:"userID"`
				}{}

				err = json.Unmarshal(message, &rawMsg)
				if err != nil {
					logrus.WithError(err).WithField("message", string(message)).Error("Could not unmarshal json")
				}

				notifications <- Notification{
					ClientID:      clientID,
					UserID:        rawMsg.UserID,
					CorrelationID: correlationID,
				}
			},
		},
	}})
	if err != nil {
		close(notifications)
		return nil, fmt.Errorf("failed to subscribe to notifications-topic: %w", err)
	}

	go func() {
		<-ctx.Done()
		close(notifications)
		err := c.c.Unsubscribe(&client.UnsubscribeOptions{
			TopicFilters: [][]byte{
				[]byte(notificationTopic),
			},
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to unsubscribe from notifications-topic")
		}
	}()

	return notifications, nil
}

type Notification struct {
	ClientID      string
	UserID        string
	CorrelationID string
}
