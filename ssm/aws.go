package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Client with AWS services
type Client struct {
	svc ssmiface.SSMAPI
}

// New returns clients for AWS services
func New(region string, awsConfig aws.Config) (*Client, error) {
	if region != "" {
		awsConfig = *awsConfig.WithRegion(region)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: awsConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %s", err)
	}

	// Create a SSM client with additional configuration
	svc := ssm.New(sess)

	return &Client{
		svc: svc,
	}, nil
}
