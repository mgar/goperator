package formatter

import (
	"os"

	"github.com/mgar/goperator/aws/ec2"
	"github.com/olekukonko/tablewriter"
)

// Ec2ToTable : converts ec2.DescribeInstancesOutput to table
func Ec2ToTable(instances []ec2.Ec2Instance) {

	data := [][]string{}

	for _, inst := range instances {
		data = append(data, []string{
			inst.InstanceID,
			inst.PublicIP,
			inst.PrivateIP,
			inst.InstanceType,
			inst.WorkingVersion,
			inst.Service,
			inst.AutoScalingGroup,
			inst.State,
		})
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
