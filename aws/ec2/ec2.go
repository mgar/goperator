package ec2

import ec2_service "github.com/aws/aws-sdk-go/service/ec2"

type Ec2Instance struct {
	InstanceID       string
	InstanceType     string
	PublicIP         string
	PrivateIP        string
	Component        string
	Environment      string
	Service          string
	AutoScalingGroup string
	State            string
	Key              string
	WorkingVersion   string
}

func NewEc2Instance(res *ec2_service.Instance) (*Ec2Instance, error) {

	instance := &Ec2Instance{
		InstanceID:   *res.InstanceId,
		State:        *res.State.Name,
		Key:          *res.KeyName,
		InstanceType: *res.InstanceType,
	}

	if res.PrivateIpAddress != nil {
		instance.PrivateIP = *res.PrivateIpAddress
	}

	if res.PublicIpAddress != nil {
		instance.PublicIP = *res.PublicIpAddress
	}
	for _, keys := range res.Tags {
		switch tag := *keys.Key; tag {
		case "working_version":
			instance.WorkingVersion = *keys.Value
		case "aws:autoscaling:groupName":
			instance.AutoScalingGroup = *keys.Value
		case "component":
			instance.Component = *keys.Value
		case "environment":
			instance.Environment = *keys.Value
		case "service":
			instance.Service = *keys.Value
		}
	}

	return instance, nil
}
