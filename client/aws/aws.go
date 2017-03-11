package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// NewEc2Client : Open new Ec2 session
func NewEc2Client() (*session.Session, *ec2.EC2) {
	awsCfg := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	s := session.New(awsCfg)
	ec2client := ec2.New(s)

	return s, ec2client
}
