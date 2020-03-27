package webhook

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// onSuccess will be called if a request to /webhook has been successfully validated
	onSuccess := func(typ EventType, ev Event, w http.ResponseWriter, r *http.Request) {
		// do nothing
	}

	secretKey := "123456"
	http.HandleFunc("/webhook", Handler(secretKey, onSuccess))
	ts := httptest.NewServer(Handler(secretKey, onSuccess))
	defer ts.Close()

	type Tst struct {
		Method         string
		ExpectedResult int
		Headers        http.Header
		Body           string
	}

	tests := []Tst{
		// Wrong method
		{Method: http.MethodGet, ExpectedResult: http.StatusBadRequest},

		// No content type
		{Method: http.MethodPost, ExpectedResult: http.StatusBadRequest},

		// No event type
		{Method: http.MethodPost, ExpectedResult: http.StatusBadRequest, Headers: http.Header{"Content-type": []string{"application/json"}}},

		// No signature
		{Method: http.MethodPost, ExpectedResult: http.StatusBadRequest, Headers: http.Header{
			"Content-type":  []string{"application/json"},
			"X-Gitea-Event": []string{"push"},
		}},

		// Invalid signature
		{Method: http.MethodPost, ExpectedResult: http.StatusBadRequest, Headers: http.Header{
			"Content-type":      []string{"application/json"},
			"X-Gitea-Event":     []string{"push"},
			"X-Gitea-Signature": []string{"blah"},
		}},

		// Valid
		{Method: http.MethodPost, ExpectedResult: http.StatusOK, Headers: http.Header{
			"Content-type":      []string{"application/json"},
			"X-Gitea-Event":     []string{"push"},
			"X-Gitea-Signature": []string{"cd2f9b218db846d088a6ed5d7cb0ec0ee8f6da141dab90c3fd826d3e7e7918fd"},
		},
			Body: `{"secret": "123456", "number": 23}`,
		},
	}

	for _, tst := range tests {
		body := bytes.NewBufferString(tst.Body)
		req, err := http.NewRequest(tst.Method, ts.URL, body)
		if err != nil {
			t.Errorf("got error: %s", err)
			return
		}

		req.Header = tst.Headers
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("got error: %s", err)
			return
		}

		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("ioutil.ReadAll returned error: %s", err)
			return
		}

		if res.StatusCode != tst.ExpectedResult {
			t.Errorf("unexpected status - expected %d, got %d. body: %s", tst.ExpectedResult, res.StatusCode, string(resp))
			return
		}
		res.Body.Close()
	}
}
