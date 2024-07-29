package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log/slog"
)

// EC2API interface added in order to make mock testing easier if future helper functions require more aws functions they should be added below interface
type EC2API interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// EC2 struct implements ec2.DescribeInstancesAPIClient from aws sdk which allows us to test its functions using the EC2API interface above
type EC2 struct {
	Client ec2.DescribeInstancesAPIClient
}

// NewEC2 initializes new ec2 client to use
func NewEC2(ctx context.Context, logger slog.Logger, accessKey, secretAccessKey, session string) *EC2 {
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, session)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(creds))
	if err != nil {

		logger.Error("error reading AWS config", err)
	}
	return &EC2{Client: ec2.NewFromConfig(cfg)}
}

func (c *EC2) GetInstanceByID(instanceID string) (*types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := c.Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found with ID %s", instanceID)
	}

	return &result.Reservations[0].Instances[0], nil
}

func (c *EC2) GetInstanceByPrivateIP(privateIP string) (*types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: []string{privateIP},
			},
		},
	}

	result, err := c.Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found with private IP %s", privateIP)
	}

	return &result.Reservations[0].Instances[0], nil
}

func (c *EC2) GetInstancesByTag(tagKey, tagValue string) ([]types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []string{tagValue},
			},
		},
	}

	result, err := c.Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("no instances found with tag %s=%s", tagKey, tagValue)
	}

	var instances []types.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}
