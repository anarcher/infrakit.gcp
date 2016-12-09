package instance

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if spec.Properties == nil {
		return nil, errors.New("Properties must be set")
	}

	props := Properties{}
	err := json.Unmarshal(*spec.Properties, &props)
	if err != nil {
		return nil, fmt.Errorf("Invalid input formatting: %s", err)
	}

	resp, err := p.service.Instances.Insert(props.Project, props.Zone, props.Instance).Do()
	if err != nil {
		return nil, err
	}

	instanceId := instance.ID(NewID(props.Project, props.Zone, resp.Name).String())

	return &instanceId, nil
}

// Destroy terminates an existing instance.
func (p gceInstancePlugin) Destroy(instanceId instance.ID) error {
	id, err := GetID(string(instanceId))
	if err != nil {
		return err
	}

	project, zone, name := id.project, id.zone, id.name

	_, err = p.service.Instances.Delete(project, zone, name).Do()
	if err != nil {
		return err
	}

	// TODO(anarcher): Implement

	return nil
}

// DescribeInstances implements instance.Provisioner.DescribeInstances.
func (p gceInstancePlugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {
	return nil, nil
}
