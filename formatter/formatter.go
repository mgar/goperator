package formatter

import (
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

// Ec2ToTable : converts ec2.DescribeInstancesOutput to table
func Ec2ToTable(resp *ec2.DescribeInstancesOutput) {
	data := [][]string{}
	workingVersion := ""

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, keys := range inst.Tags {
				if *keys.Key == "working_version" {
					workingVersion = *keys.Value
				}
			}
			data = append(data, []string{
				*inst.InstanceId,
				*inst.PublicIpAddress,
				*inst.PrivateIpAddress,
				*inst.InstanceType,
				workingVersion,
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
	})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
