package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type API struct {
	URL      string // URL of Gitea server
	Username string
	Password string
	Token    string
}

func (api *API) UpdateCommitState(repository string, commitID string, status CreateStatusOption) error {
	u, err := url.Parse(api.URL)
	if err != nil {
		return err
	}
	u.Path = path.Join("api", "v1", "repos", repository, "statuses", commitID)

	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(status)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	if api.Token != "" {
		r.Header.Add("Authorization", "token "+api.Token)
	} else {
		r.SetBasicAuth(api.Username, api.Password)
	}
	client := http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 &&
		resp.StatusCode != 201 {
		return fmt.Errorf("invalid status code returned: %d", resp.StatusCode)
	}
	return nil
}
