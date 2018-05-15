package cloudenv

import (
	"time"
)

const (
	probeTimeout = 15 // seconds
)

type CloudConfig struct {
	Provider   string
	Region     string
	AccountID  string
	AZ         string
	InstanceID string
	Image      string
}

type cloudprovider interface {
	probe(r chan *CloudConfig)
}

func Discover() *CloudConfig {
	r := make(chan *CloudConfig, 1)
	defer close(r)

	providers := []cloudprovider{
		newAWSProvider(),
		newGCPProvider(),
	}

	for _, provider := range providers {
		provider.probe(r)
	}

	ticker := time.NewTicker(probeTimeout * time.Second)
	defer ticker.Stop()

	cfg := new(CloudConfig)

	select {
	case cfg = <-r:
	case <-ticker.C:
	}

	return cfg
}
