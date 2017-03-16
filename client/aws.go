package client

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// NewEc2Client : Open new Ec2 session
func NewEc2Client() *ec2.EC2 {
	region := os.Getenv("AWS_REGION")

	// Set default region if none has been selected
	if region == "" {
		region = "us-east-1"
	}

	awsCfg := &aws.Config{
		Region: aws.String(region),
	}
	s := session.New(awsCfg)
	ec2client := ec2.New(s)

	return ec2client
}
