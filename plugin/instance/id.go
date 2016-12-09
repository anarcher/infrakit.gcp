package instance

import (
	"fmt"
	"strings"
)

type ID struct {
	project string
	zone    string
	name    string
}

func NewID(project, zone, name string) *ID {
	id := &ID{
		project: project,
		zone:    zone,
		name:    name,
	}

	return id
}

func GetID(id string) (*ID, error) {
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("id isn't good %v", id)
	}
	project, zone, name := parts[0], parts[1], parts[2]

	return NewID(project, zone, name), nil
}

func (id ID) String() string {
	return fmt.Sprintf("%s-%s-%s", id.project, id.zone, id.name)
}
