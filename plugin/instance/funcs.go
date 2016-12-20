package instance

import (
	"math/rand"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Return instance tags from compute.Metadata
func instanceTags(metadata *compute.Metadata) map[string]string {
	tags := map[string]string{}
	if metadata == nil || metadata.Items == nil {
		return tags
	}
	for _, item := range metadata.Items {
		key := ensureToTagKey(item.Key)
		tags[key] = *item.Value
	}
	return tags
}

// Replace dot(.) in value to dash(-) in key value (GCP Metadata key must be a match of regex '[a-zA-Z0-9-_]{1,128})
func ensureToMetadataKey(v string) string {
	return strings.Replace(v, ".", "---", -1)
}

func ensureToTagKey(v string) string {
	return strings.Replace(v, "---", ".", -1)
}
