package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
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

// readDockerImageFromText reads the Docker image repo and tag from the
// Annotation's text. If no complete info is found, return an error.
func readDockerImageFromText(t string) (*DockerImage, error) {
	pattern := `<b>Docker Image:</b> (?P<repo>\w+):(?P<tag>[a-zA-Z0-9\.]+)<br>`
	re := regexp.MustCompile(pattern)

	var repo string
	var tag string

	for _, match := range re.FindAllStringSubmatch(t, 1) {
		for i, name := range re.SubexpNames() {
			if name == "repo" {
				repo = match[i]
			}
			if name == "tag" {
				tag = match[i]
			}
		}
	}

	if repo == "" || tag == "" {
		return nil, fmt.Errorf("no Docker image tag in text")
	}

	return &DockerImage{Repo: repo, Tag: tag}, nil
}

// readEnvironmentFromTags
func readEnvironmentFromTags(tags []string) (string, error) {
	pattern := `env:(\w+)`
	re := regexp.MustCompile(pattern)
	for _, t := range tags {
		if submatches := re.FindStringSubmatch(t); submatches != nil {
			return submatches[1], nil
		}
	}
	return "", fmt.Errorf("cannot find env in tags")
}

// readCanaryFromText returns if it's a canary deployment
func readCanaryFromText(t string) bool {
	pattern := `Canary Deployment:`
	re := regexp.MustCompile(pattern)
	return re.MatchString(t)
}

// newDeploymentFromAnnotation creates a new Deployment from the given Annotation
func newDeploymentFromAnnotation(a *Annotation) (*Deployment, error) {
	dockerImage, err := readDockerImageFromText(a.Text)
	if err != nil {
		log.Printf("unable to read Docker Image from annotation: %v", err)
	}

	env, err := readEnvironmentFromTags(a.Tags)
	if err != nil {
		log.Printf("unable to read Environment from annotation: %v", err)
	}

	deployment := &Deployment{
		CreatedAt:   time.UnixMilli(a.CreatedAt).UTC(),
		UpdatedAt:   time.UnixMilli(a.UpdatedAt).UTC(),
		Env:         env,
		DockerImage: dockerImage,
		Canary:      readCanaryFromText(a.Text),
	}

	return deployment, nil
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
		deployment, err := newDeploymentFromAnnotation(&a)
		if err != nil {
			log.Printf("error creating Deployment from Annotation %v", err)
			continue
		}
		deployments = append(deployments, *deployment)
	}

	return deployments, nil
}
