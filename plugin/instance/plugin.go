package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/infrakit/pkg/spi/instance"
	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	ZONE_OPERATION_STATUS_DONE   = "DONE"
	INSTANCE_STATUS_RUNNING      = "RUNNING"
	INSTANCE_STATUS_PROVISIONING = "PROVISIONING"
)

type gceInstancePlugin struct {
	service *compute.Service
	project string
	zone    string
}

// NewInstancePlugin creates a new plugin that creates instances in GCE.
func NewInstancePlugin(service *compute.Service, project, zone string) instance.Plugin {
	return &gceInstancePlugin{
		service: service,
		project: project,
		zone:    zone,
	}
}

// Validate performs local checks to determine if the request is valid.
func (p *gceInstancePlugin) Validate(req json.RawMessage) error {
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

	var instanceName string

	if spec.LogicalID != nil {
		if props.Instance.Name != "" {
			instanceName = fmt.Sprintf("%s-%s", props.Instance.Name, spec.LogicalID)
		} else {
			instanceName = fmt.Sprintf("%s", spec.LogicalID)
		}
	} else {
		if props.Instance.Name != "" {
			instanceName = fmt.Sprintf("%s-%d", props.Instance.Name, rand.Int31())
		} else {
			instanceName = fmt.Sprintf("%s-%d", spec.Tags["infrakit.group"], rand.Int31())
		}
	}

	props.Instance.Name = instanceName

	if props.Instance.Metadata == nil {
		props.Instance.Metadata = &compute.Metadata{
			Items: []*compute.MetadataItems{},
		}
	}

	for k, v := range spec.Tags {
		key := ensureToMetadataKey(k)
		value := v
		props.Instance.Metadata.Items = append(props.Instance.Metadata.Items, &compute.MetadataItems{Key: key, Value: &value})
	}

	if spec.Init != "" {
		props.Instance.Metadata.Items = append(props.Instance.Metadata.Items, &compute.MetadataItems{Key: "startup-script", Value: &spec.Init})
	}

	resp, err := p.service.Instances.Insert(p.project, p.zone, props.Instance).Do()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	instanceId := instance.ID(instanceName)
	log.Debugf("Instance id is %v The operation of it is %v", instanceId, resp.Name)

	//Checking provision processing
	{
		timer := time.NewTimer(time.Second * 20)
		ticker := time.NewTicker(time.Second * 1)
	C:
		for _ = range ticker.C {

			select {
			case <-timer.C:
				return nil, fmt.Errorf("Instance %v provisioning operation %v has timed out", instanceName, resp.Name)
			default:
				resp, err := p.service.ZoneOperations.Get(p.project, p.zone, resp.Name).Do()
				if err != nil {
					log.Error(err)
					return nil, err
				}
				log.Debugf("Provision instance %v has %v status", instanceName, resp.Status)
				if resp.Status == ZONE_OPERATION_STATUS_DONE {
					break C
				}
			}

		}
	}

	return &instanceId, nil
}

// Destroy terminates an existing instance.
func (p gceInstancePlugin) Destroy(instanceId instance.ID) error {
	_, err := p.service.Instances.Delete(p.project, p.zone, string(instanceId)).Do()
	if err != nil {
		return err
	}

	// TODO(anarcher): Implement

	return nil
}

// DescribeInstances implements instance.Provisioner.DescribeInstances.
func (p gceInstancePlugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {
	descriptions := []instance.Description{}

	call := p.service.Instances.List(p.project, p.zone)
	//todo(anarcher): Currently,GCP compute instances list API filter doesn't support metadata filtering.
	//So getting all instances information and then filtering by tags for now.
	/*
		for k, _ := range tags {
			call = call.Filter(fmt.Sprintf("metadata.items.key=%s", k))
		}
	*/
	//call = call.Filter("status=RUNNING")

	call = call.Fields(googleapi.Field("items(id,metadata/items,name,status)"))

	ctx := context.Background() //todo(anarcher)
	if err := call.Pages(ctx, func(page *compute.InstanceList) error {
		for _, v := range page.Items {

			metadataJson, _ := v.Metadata.MarshalJSON()
			log.Debugf("Instance Name:%v,Status:%v,Metadata:%s", v.Name, v.Status, metadataJson)

			if v.Status == INSTANCE_STATUS_RUNNING || v.Status == INSTANCE_STATUS_PROVISIONING {
			} else {
				continue
			}

			found := false
			iTags := instanceTags(v.Metadata)
			for k, v := range iTags {
				for _k, _v := range tags {
					if k == _k && v == _v {
						found = true
					}
				}
			}

			if !found {
				continue
			}

			logicalID := instance.LogicalID(v.Name)
			descriptions = append(descriptions, instance.Description{
				ID:        instance.ID(v.Name),
				LogicalID: &logicalID,
				Tags:      iTags,
			})
		}
		return nil
	}); err != nil {
		log.Error(err)
		return descriptions, err
	}

	log.Debugf("There is %v related instances", len(descriptions))
	return descriptions, nil
}
