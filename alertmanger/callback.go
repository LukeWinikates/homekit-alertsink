package alertmanger

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewHandlerFunc(callback func(payload *WebhookPayload) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var webhookPayload *WebhookPayload
		err := json.NewDecoder(request.Body).Decode(webhookPayload)
		if err != nil {
			log.Printf("error decoding request: %s", err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
		err = callback(webhookPayload)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}
