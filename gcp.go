package cloudenv

import (
	"strings"

	gcpmetadata "cloud.google.com/go/compute/metadata"
)

type gcpProvider struct{}

func newGCPProvider() *cloudprovider {
	return new(gcpProvider)
}

func (p *gcpProvider) probe(r chan *CloudConfig) {
	if !gcpmetadata.OnGCE() {
		return
	}

	cfg := new(CloudConfig)

	cfg.Provider = "gcp"

	zone, _ := gcpmetadata.Zone()
	cfg.AZ = zone

	cfg.Region = regionFromZone(zone)

	projectID, _ := gcpmetadata.ProjectID()
	cfg.AccountID = projectID

	instanceID, _ := gcpmetadata.InstanceID()
	cfg.InstanceID = instanceID

	image, _ := gcpmetadata.Get("instance/image")
	cfg.Image = image

	r <- cfg

	return
}

func regionFromZone(z string) string {
	parts := strings.Split(z, "-")

	if len(parts) < 3 {
		return ""
	}

	return strings.Join(parts[:len(parts)-1], "-")
}