package ssm

import (
	"errors"

	"github.com/letsgodevops/ssm-sync/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const defaultParameterType = "SecureString"

// PutObject writes a given secret value on SSM
// it uses PutParameter API call
// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_PutParameter.html
func (c *Client) PutObject(object *types.PutObjectInput) error {
	if object.Key == "" {
		return errors.New("key name is not valid")
	}
	// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#PutParameterInput
	putParameterInput := &ssm.PutParameterInput{
		KeyId:     aws.String(buildKeyAliasPath(object.KmsKeyAlias)),
		Name:      aws.String("/" + removeSlashPrefix(object.Key)),
		Type:      aws.String(defaultParameterType),
		Value:     aws.String(object.Value),
		Overwrite: aws.Bool(true),
	}

	// Ignore PutParameter returned Version
	_, err := c.svc.PutParameter(putParameterInput)
	return err
}

// GetObject returns a secret for given key
func (c *Client) GetObject(object *types.GetObjectInput) (*types.GetObjectOutput, error) {
	params := &ssm.GetParameterInput{
		// we decided to use path based keys without `/` at the begining
		// so we need to add it here
		Name: aws.String("/" + removeSlashPrefix(object.Key)),
		// Retrieve all parameters in a hierarchy with their value decrypted.
		WithDecryption: aws.Bool(true),
	}

	resp, err := c.svc.GetParameter(params)
	if err != nil {
		return nil, err
	}

	return &types.GetObjectOutput{Value: *resp.Parameter.Value}, nil
}
