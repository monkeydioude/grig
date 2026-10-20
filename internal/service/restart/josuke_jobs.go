package restart

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"monkeydioude/grig/internal/model"
	"net/http"
	"slices"
	"strings"
)

var ErrNoJobsAPI = errors.New("josuke config has no jobs_api token: grig cannot check its jobs")

// JobsClient reads Josuke's jobs API (GET /jobs).
type JobsClient struct {
	URL    string
	Token  string
	Client *http.Client
}

// NewJobsClient builds the client from Josuke's own config: its address
// and its jobs_api section (token, optional route).
func NewJobsClient(jk model.Josuke, client *http.Client) (JobsClient, error) {
	api, _ := jk.Extra["jobs_api"].(map[string]any)
	token, _ := api["token"].(string)
	if token == "" {
		return JobsClient{}, ErrNoJobsAPI
	}
	route, _ := api["route"].(string)
	route = cmp.Or(strings.TrimSuffix(route, "/"), "/jobs")
	host := cmp.Or(jk.Host, "localhost")
	return JobsClient{
		URL:    fmt.Sprintf("http://%s:%d%s", host, jk.Port, route),
		Token:  token,
		Client: client,
	}, nil
}

type job struct {
	Status string `json:"status"` // running, stopping or waiting
}

// Busy tells whether a job is running or being stopped.
func (c JobsClient) Busy(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	res, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("GET %s: %d", c.URL, res.StatusCode)
	}
	jobs := []job{}
	if err := json.NewDecoder(res.Body).Decode(&jobs); err != nil {
		return false, err
	}
	return slices.ContainsFunc(jobs, func(j job) bool { return j.Status != "waiting" }), nil
}
