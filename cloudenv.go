package cloudenv

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultTimeout = 5 // seconds
)

// CloudConfig stores the cloud provider's config
type CloudConfig struct {
	Provider   string
	Region     string
	AccountID  string
	AZ         string
	InstanceID string
	Image      string
}

type cloudprovider interface {
	probe(ctx context.Context, r chan *CloudConfig)
}

type cloudproviders struct {
	providers []cloudprovider
}

// Provider is a function which adds to the list of cloud providers
type Provider func(*cloudproviders)

func (c CloudConfig) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", c.Provider, c.Region, c.AccountID, c.AZ, c.InstanceID, c.Image)
}

// AWS appends AWS to the list of providers
func AWS() Provider {
	return func(c *cloudproviders) {
		c.providers = append(c.providers, newAWSProvider())
		return
	}
}

// GCP appends GCP to the list of providers
func GCP() Provider {
	return func(c *cloudproviders) {
		c.providers = append(c.providers, newGCPProvider())
		return
	}
}

// Discover returns a CloudConfig for the cloud provider that is being
// used. Timeout is the timeout value when attempting to probe each cloud
// provider. Providers are chosen from the passed in list of providers and
// defaults to AWS and GCP if none are supplied.
func Discover(timeout time.Duration, providers ...Provider) *CloudConfig {
	if timeout < (1 * time.Second) {
		timeout = defaultTimeout * time.Second
	}

	if len(providers) == 0 {
		providers = append(providers, []Provider{AWS(), GCP()}...)
	}

	cp := new(cloudproviders)

	for _, p := range providers {
		if p == nil {
			continue
		}
		p(cp)
	}

	r := make(chan *CloudConfig, 1)

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	cfg := new(CloudConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, p := range cp.providers {
		go p.probe(ctx, r)
	}

	select {
	case cfg = <-r:
		return cfg
	case <-ticker.C:
		return cfg
	}
}
