package instance

import (
	"google.golang.org/api/compute/v1"
)

// Return instance tags from compute.Metadata
func instanceTags(metadata *compute.Metadata) map[string]string {
	tags := map[string]string{}
	for _, item := range metadata.Items {
		tags[item.Key] = *item.Value
	}
	return tags
}
