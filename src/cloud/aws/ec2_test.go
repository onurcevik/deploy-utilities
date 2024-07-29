package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2Client "github.com/onurcevik/deploy-service/src/cloud/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2API is a mock type for the EC2API interface
type MockEC2API struct {
	mock.Mock
}

// DescribeInstances provides a mock function with given fields: ctx, params, optFns
func (m *MockEC2API) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, params, optFns)
	output, ok := args.Get(0).(*ec2.DescribeInstancesOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return output, args.Error(1)
}

func TestGetInstanceByID(t *testing.T) {
	mockEC2 := new(MockEC2API)
	client := &ec2Client.EC2{Client: mockEC2}

	instanceID := "i-123456"
	mockEC2.On("DescribeInstances", mock.Anything, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}, mock.Anything).Return(&ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{InstanceId: &instanceID},
				},
			},
		},
	}, nil)

	instance, err := client.GetInstanceByID(instanceID)
	assert.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, instanceID, *instance.InstanceId)

	mockEC2.AssertExpectations(t)
}

func TestGetInstanceByPrivateIP(t *testing.T) {
	mockEC2 := new(MockEC2API)
	client := &ec2Client.EC2{Client: mockEC2}

	privateIP := "10.0.0.1"
	instanceID := "i-123456"
	mockEC2.On("DescribeInstances", mock.Anything, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: []string{privateIP},
			},
		},
	}, mock.Anything).Return(&ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{InstanceId: &instanceID, PrivateIpAddress: &privateIP},
				},
			},
		},
	}, nil)

	instance, err := client.GetInstanceByPrivateIP(privateIP)
	assert.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, privateIP, *instance.PrivateIpAddress)

	mockEC2.AssertExpectations(t)
}

func TestGetInstancesByTag(t *testing.T) {
	mockEC2 := new(MockEC2API)
	client := &ec2Client.EC2{Client: mockEC2}

	tagKey := "Name"
	tagValue := "test-instance"
	instanceID := "i-123456"
	mockEC2.On("DescribeInstances", mock.Anything, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []string{tagValue},
			},
		},
	}, mock.Anything).Return(&ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{InstanceId: &instanceID},
				},
			},
		},
	}, nil)

	instances, err := client.GetInstancesByTag(tagKey, tagValue)
	assert.NoError(t, err)
	assert.NotNil(t, instances)
	assert.Len(t, instances, 1)
	assert.Equal(t, instanceID, *instances[0].InstanceId)

	mockEC2.AssertExpectations(t)
}
