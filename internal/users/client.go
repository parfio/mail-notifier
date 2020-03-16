package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var UserNotFoundErr = errors.New("users not found")

type Client struct {
	usersServiceBaseURL string
}

func NewClient(userServiceBaseURL string) Client {
	if !strings.HasSuffix(userServiceBaseURL, "/") {
		userServiceBaseURL += "/"
	}

	return Client{usersServiceBaseURL: userServiceBaseURL}
}

func (c Client) ByUserID(clientID, userID string) (User, error) {
	resp, err := http.Get(fmt.Sprintf("%sv1/users/%s?clientID=%s", c.usersServiceBaseURL, userID, clientID))
	if err != nil {
		return User{}, fmt.Errorf("failed to request users-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return User{}, UserNotFoundErr
	}

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("failed to request users-service: expected status code 200, given %d", resp.StatusCode)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, fmt.Errorf("failed to decode users-service response body: %w", err)
	}

	return user, nil
}

type User struct {
	EMail string `json:"email"`
}
