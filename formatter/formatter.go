package formatter

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

// Ec2ToTable : converts ec2.DescribeInstancesOutput to table
func Ec2ToTable(resp *ec2.DescribeInstancesOutput) {

	data := [][]string{}

	var (
		workingVersion string
		asg            string
		publicIP       string
		service        string
		privateIP      string
	)

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, keys := range inst.Tags {
				switch tag := *keys.Key; tag {
				case "working_version":
					workingVersion = *keys.Value
				case "aws:autoscaling:groupName":
					asg = *keys.Value
				case "service":
					service = *keys.Value
				}
			}
			if inst.PublicIpAddress != nil {
				publicIP = *inst.PublicIpAddress
			}
			if inst.PrivateIpAddress != nil {
				privateIP = *inst.PrivateIpAddress
			}

			data = append(data, []string{
				*inst.InstanceId,
				publicIP,
				privateIP,
				*inst.InstanceType,
				workingVersion,
				service,
				asg,
				*inst.State.Name,
			})
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Instance ID",
		"Public IP",
		"Internal IP",
		"Instance Type",
		"Working Version",
		"Service",
		"Auto Scaling group",
		"Instance State",
	})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
