package instance

import (
	"encoding/json"
	"net/http"

	"github.com/docker/infrakit/pkg/spi/instance"
	"google.golang.org/api/compute/v1"
)

type gceInstancePlugin struct {
	client  *http.Client
	service *compute.Service
}

// NewInstancePlugin creates a new plugin that creates instances in GCE.
func NewInstancePlugin(client *http.Client) instance.Plugin {
	return &gceInstancePlugin{client: client}
}

// Validate performs local checks to determine if the request is valid.
func (p *gceInstancePlugin) Validate(req json.RawMessage) error {
	service, err := compute.New(p.client)
	if err != nil {
		return err
	}
	p.service = service

	// TODO(anarcher): Implement
	return nil
}

// Provision creates a new instance.
func (p gceInstancePlugin) Provision(spec instance.Spec) (*instance.ID, error) {
	// TODO(anarcher): Implement
	return nil, nil
}

// Destroy terminates an existing instance.
func (p gceInstancePlugin) Destroy(id instance.ID) error {
	// TODO(anarcher): Implement
	return nil
}

// DescribeInstances implements instance.Provisioner.DescribeInstances.
func (p gceInstancePlugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {
	return nil, nil
}
