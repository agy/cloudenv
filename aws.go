package cloudenv

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsProvider struct{}

func newAWSProvider() *cloudprovider {
	return new(awsProvider)
}

func (p *awsProvider) probe(r chan *CloudConfig) {
	s, _ := session.NewSession()
	metadata := ec2metadata.New(s)

	if !metadata.Available() {
		return
	}

	doc, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return
	}

	cfg := new(CloudConfig)

	cfg.Provider = "aws"
	cfg.AZ = doc.AvailabilityZone
	cfg.Region = doc.Region
	cfg.AccountID = doc.AccountID
	cfg.InstanceID = doc.InstanceID
	cfg.Image = doc.ImageID

	r <- cfg

	return
}
