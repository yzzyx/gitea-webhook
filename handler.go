package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Handler returns a http handler function that validates a gitea webhook request,
// and, if successful, passes the event information to the 'onSuccess'-function.
func Handler(secretKey string, onSuccess func(typ EventType, ev Event, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		contentType := strings.ToLower(r.Header.Get("Content-type"))
		if contentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid content type")
			return
		}

		eventTypeStr := r.Header.Get("X-Gitea-Event")
		if eventTypeStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "No event header specified")
			return
		}

		eventType, ok := eventTypeTrans[eventTypeStr]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unsupported eventtype %s", eventTypeStr)
			return
		}

		signature := r.Header.Get("X-Gitea-Signature")
		if signature == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "No signature header specified")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not ready body: %s", err)
			return
		}

		mac := hmac.New(sha256.New, []byte(secretKey))
		mac.Write(body)
		actualSignature := fmt.Sprintf("%x", mac.Sum(nil))
		if actualSignature != signature {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not validate signature of body")
			return
		}

		var event Event
		err = json.Unmarshal(body, &event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not unmarshal body: %s", err)
			return
		}

		onSuccess(eventType, event, w, r)
	}
}
