package web

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartAliveEndpoint(serverAddress string) error {
	return http.ListenAndServe(serverAddress, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/internal/alive" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		_, err := w.Write([]byte(`{"alive":true}`))
		if err != nil {
			logrus.WithError(err).Error("Failed to write alive response body")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}
