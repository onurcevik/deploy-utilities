package aws

import (
	"context"
	"github.com/onurcevik/deploy-utilities/common/config"
	"log/slog"
)

// AWSClient holds fields for logger and AWS services
type AWSClient struct {
	Logger slog.Logger
	EC2    *EC2
}

// NewAWSClient initializes AWSClient and its service fields
func NewAWSClient(ctx context.Context, logger slog.Logger, conf config.Config) *AWSClient {
	aws := new(AWSClient)
	aws.Logger = logger
	aws.EC2 = NewEC2(ctx, logger, conf.AWS.AWSAccessKey, conf.AWS.AWSSecretAccessKey, conf.AWS.Session)
	return aws
}
