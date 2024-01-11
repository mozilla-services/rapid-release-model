package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// AnnotationsFilter holds parameters for filtering Grafana Annotations
type AnnotationsFilter struct {
	App   string
	From  string
	To    string
	Limit int
}

// CreateURLValues creates URL values based on the given AnnotationsFilter
func CreateURLValues(f *AnnotationsFilter) url.Values {
	values := make(url.Values)

	// Fixed parameters for deployments
	values.Add("type", "annotation")
	values.Add("tags", "event_type:deployment")
	values.Add("tags", "event_status:complete")

	// Variable parameters specified by the user
	values.Add("tags", fmt.Sprintf("app:%s", f.App))
	values.Add("from", f.From)
	values.Add("to", f.To)
	values.Add("limit", strconv.Itoa(f.Limit))

	return values
}

// Annotation represents an annotation object as returned by the Grafana REST API
type Annotation struct {
	Text      string   `json:"text"`
	CreatedAt int64    `json:"created"`
	UpdatedAt int64    `json:"updated"`
	Time      int64    `json:"time"`
	TimeEnd   int64    `json:"timeEnd"`
	Tags      []string `json:"tags"`
}

// DockerImage used for a Deployment
type DockerImage struct {
	Repo string `json:"repo"`
	Tag  string `json:"tag"`
}

// Deployment event referenced by a Grafana annotation
type Deployment struct {
	DockerImage *DockerImage
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Env         string
	Canary      bool
}

// QueryDeployments fetches information about Deployments from the Grafana REST API
func QueryDeployments(ctx context.Context, httpClient HTTPClient, filter *AnnotationsFilter) ([]Deployment, error) {
	// Create HTTP request query parameters
	urlValues := CreateURLValues(filter)

	respData, err := httpClient.Get(ctx, "api/annotations", urlValues)
	if err != nil {
		return nil, err
	}

	var annotations []Annotation

	if err := json.Unmarshal(respData, &annotations); err != nil {
		return nil, err
	}

	var deployments []Deployment

	for _, a := range annotations {
		// TODO: Parse Docker image, environment, and canary from the annotation text
		deployment := &Deployment{
			CreatedAt: time.UnixMilli(a.CreatedAt).UTC(),
			UpdatedAt: time.UnixMilli(a.UpdatedAt).UTC(),
		}
		deployments = append(deployments, *deployment)
	}

	return deployments, nil
}
