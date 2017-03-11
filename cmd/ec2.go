package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	client "github.com/mgar/goperator/client/aws"
	"github.com/mgar/goperator/formatter"
	"github.com/spf13/cobra"
)

var session, ec2Client = client.NewEc2Client()

var cmdListInstances = &cobra.Command{
	Use:  "list [environment] [component]",
	Long: "List Ec2 instances",
	Run:  listInstances,
}

var cmdSSHInstance = &cobra.Command{
	Use:  "ssh [instance-id]",
	Long: "SSH into an Ec2 Instance",
	Run:  sshInstance,
}

func listInstances(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("You need to specify [environment] and [component]")
		os.Exit(1)
	}
	environment, component := args[0], args[1]

	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:component"),
				Values: []*string{
					aws.String(component),
				},
			},
			{
				Name: aws.String("tag:environment"),
				Values: []*string{
					aws.String(environment),
				},
			},
		},
	}
	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	formatter.Ec2ToTable(resp)
}

func sshInstance(cmd *cobra.Command, args []string) {
	fmt.Print("SSHing...")
}
