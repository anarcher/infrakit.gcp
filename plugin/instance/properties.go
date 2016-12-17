package instance

import (
	"google.golang.org/api/compute/v1"
)

type Properties struct {
	Instance *compute.Instance `json:"instance"`
}
