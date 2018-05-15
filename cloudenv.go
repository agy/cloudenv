package cloudenv

import (
	"time"
)

const (
	defaultTimeout = 5 // seconds
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

type cloudproviders struct {
	providers []cloudprovider
}

type provider func(*cloudproviders)

func AWS() provider {
	return func(c *cloudproviders) {
		c.providers = append(c.providers, newAWSProvider())
		return
	}
}

func GCP() provider {
	return func(c *cloudproviders) {
		c.providers = append(c.providers, newGCPProvider())
		return
	}
}

func Discover(timeout time.Duration, providers ...provider) *CloudConfig {
	if timeout <= 0 {
		timeout = defaultTimeout * time.Second
	}

	cp := new(cloudproviders)

	if len(providers) == 0 {
		providers = []provider{
			AWS(),
			GCP(),
		}
	}

	for _, p := range providers {
		p(cp)
	}

	r := make(chan *CloudConfig, 1)
	defer close(r)

	for _, p := range cp.providers {
		p.probe(r)
	}

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	cfg := new(CloudConfig)

	select {
	case cfg = <-r:
	case <-ticker.C:
	}

	return cfg
}
