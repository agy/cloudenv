package cloudenv

import (
	"context"
	"reflect"
	"testing"
	"time"
)

const (
	fakeProviderName = "fake"
	fakeAZ           = "az"
	fakeRegion       = "region"
	fakeAccountID    = "accountid"
	fakeInstanceID   = "instanceid"
	fakeImage        = "imageid"
)

func addFakeProvider() Provider {
	return func(c *cloudproviders) {
		c.providers = append(c.providers, fakeProvider{})
		return
	}
}

type fakeProvider struct{}

func newFakeProvider() cloudprovider {
	return fakeProvider{}
}

func (p fakeProvider) probe(ctx context.Context, r chan *CloudConfig) {
	select {
	case <-ctx.Done():
		return
	default:
		cfg := new(CloudConfig)

		cfg.Provider = fakeProviderName
		cfg.AZ = fakeAZ
		cfg.Region = fakeRegion
		cfg.AccountID = fakeAccountID
		cfg.InstanceID = fakeInstanceID
		cfg.Image = fakeImage

		r <- cfg

		return
	}

	return
}

func TestAWS(t *testing.T) {
	cp := new(cloudproviders)
	p := AWS()
	p(cp)

	if len(cp.providers) != 1 {
		t.Errorf("expected aws provider to be added to provider list")
	}

	switch kind := cp.providers[0].(type) {
	case awsProvider:
	default:
		t.Errorf("invalid provider type, got %v, expected %s", kind, "awsProvider")
	}
}

func TestGCP(t *testing.T) {
	cp := new(cloudproviders)
	p := GCP()
	p(cp)

	if len(cp.providers) != 1 {
		t.Errorf("expected gcp provider to be added to provider list")
	}

	switch kind := cp.providers[0].(type) {
	case gcpProvider:
	default:
		t.Errorf("invalid provider type, got %v, expected %s", kind, "gcpProvider")
	}
}

func TestDiscover(t *testing.T) {
	cases := []struct {
		Name     string
		Timeout  time.Duration
		Provider Provider
		Expected *CloudConfig
	}{
		{
			Name:     "HappyPath",
			Provider: addFakeProvider(),
			Timeout:  1 * time.Second,
			Expected: &CloudConfig{
				Provider:   fakeProviderName,
				AZ:         fakeAZ,
				Region:     fakeRegion,
				AccountID:  fakeAccountID,
				InstanceID: fakeInstanceID,
				Image:      fakeImage,
			},
		},
		{
			Name:     "TooShortTimeout",
			Provider: addFakeProvider(),
			Timeout:  1,
			Expected: &CloudConfig{
				Provider:   fakeProviderName,
				AZ:         fakeAZ,
				Region:     fakeRegion,
				AccountID:  fakeAccountID,
				InstanceID: fakeInstanceID,
				Image:      fakeImage,
			},
		},
		{
			Name:     "ZeroTimeout",
			Provider: addFakeProvider(),
			Timeout:  0,
			Expected: &CloudConfig{
				Provider:   fakeProviderName,
				AZ:         fakeAZ,
				Region:     fakeRegion,
				AccountID:  fakeAccountID,
				InstanceID: fakeInstanceID,
				Image:      fakeImage,
			},
		},
		{
			Name:     "NegativeTimeout",
			Provider: addFakeProvider(),
			Timeout:  -1,
			Expected: &CloudConfig{
				Provider:   fakeProviderName,
				AZ:         fakeAZ,
				Region:     fakeRegion,
				AccountID:  fakeAccountID,
				InstanceID: fakeInstanceID,
				Image:      fakeImage,
			},
		},
		{
			Name:     "NilProvider",
			Provider: nil,
			Timeout:  1 * time.Second,
			Expected: &CloudConfig{},
		},
	}

	for _, c := range cases {
		result := Discover(c.Timeout, c.Provider)

		if !reflect.DeepEqual(result, c.Expected) {
			t.Errorf("%s: got %v, expected %v\n", c.Name, result, c.Expected)
		}
	}

	{
		expected := &CloudConfig{}
		result := Discover(1 * time.Second)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("NoProvider: got %v, expected %v\n", result, expected)
		}
	}
}
