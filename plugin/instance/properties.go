package instance

import (
	"google.golang.org/api/compute/v1"
)

type Properties struct {
	Project  string            `json:"project"`
	Zone     string            `json:"zone"`
	Instance *compute.Instance `json:"instance"`
}
