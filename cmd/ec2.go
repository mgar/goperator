package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	client "github.com/mgar/goperator/client/aws"
	"github.com/mgar/goperator/formatter"
	"github.com/spf13/cobra"
)

var session, ec2Client = client.NewEc2Client()

var cmdListInstances = &cobra.Command{
	Use:   "list [environment] [component]",
	Long:  "List EC2 instances based on [environment] and [component] tags",
	Short: "List EC2 instances",
	Run:   listInstances,
}

var cmdSSHInstance = &cobra.Command{
	Use:   "ssh [instance-id]",
	Long:  "SSH into an EC2 Instance given the its [instance-id]",
	Short: "SSH into an EC2 instance",
	Run:   sshInstance,
}

var cmdStopInstance = &cobra.Command{
	Use:   "stop [instance-id]",
	Long:  "Stop one or more EC2 instances given [instance-id ...]",
	Short: "Stop one or many EC2 instances",
	Run:   stopInstance,
}

var cmdStartInstance = &cobra.Command{
	Use:   "start [instance-id]",
	Long:  "Start one or more EC2 instances given [instance-id ...]",
	Short: "Start one or many EC2 instances",
	Run:   startInstance,
}

var execCommandInInstance = &cobra.Command{
	Use:   "exec [instance-id ...]  [command]",
	Long:  "Execute a command on one or a given number of EC2 instances given [instance-id ...]",
	Short: "Execute a command on one or many EC2 instances",
	Run:   executeCommand,
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
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("stopped"), aws.String("shutting-down"),
					aws.String("stopping"), aws.String("pending")},
			},
		},
	}
	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp.Reservations != nil {
		formatter.Ec2ToTable(resp)
	}
}

func sshInstance(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id]")
		os.Exit(1)
	}
	instanceID := args[0]
	var component, environment, IP string

	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("stopped")},
			},
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(instanceID)},
			},
		},
	}
	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	inst := resp.Reservations[0].Instances
	for _, inst := range inst {
		for _, keys := range inst.Tags {
			switch tag := *keys.Key; tag {
			case "component":
				component = *keys.Value
			case "environment":
				environment = *keys.Value
			}
		}
		if inst.PublicIpAddress != nil {
			IP = *inst.PublicIpAddress
		} else {
			IP = *inst.PrivateIpAddress
		}
	}
	err = connectSSH(environment, component, IP)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func stopInstance(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id ...]")
		os.Exit(1)
	}
	instances := []*string{}
	for inst := range args {
		instances = append(instances, &args[inst])
	}

	params := &ec2.StopInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: instances,
	}

	resp, err := ec2Client.StopInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	for idx := range resp.StoppingInstances {
		if *resp.StoppingInstances[idx].PreviousState.Name != "running" {
			fmt.Printf("Instance %s in not running. It can't be stopped.\n", *resp.StoppingInstances[idx].InstanceId)
		} else {
			fmt.Printf("Stopping instance: %s\n", *resp.StoppingInstances[idx].InstanceId)
		}
	}
}

func startInstance(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id ...]")
		os.Exit(1)
	}

	instances := []*string{}
	for inst := range args {
		instances = append(instances, &args[inst])
	}

	params := &ec2.StartInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: instances,
	}

	resp, err := ec2Client.StartInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	for idx := range resp.StartingInstances {
		if *resp.StartingInstances[idx].CurrentState.Name == *resp.StartingInstances[idx].PreviousState.Name {
			fmt.Printf("Instance %s is already running. Skipping...\n", *resp.StartingInstances[idx].InstanceId)
		} else {
			fmt.Printf("Starting instance: %s\n", *resp.StartingInstances[idx].InstanceId)
		}
	}
}

func executeCommand(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		fmt.Println("You need to specify [instance-id ...] [command]")
		os.Exit(1)
	}

	var component, environment, IP string

	command := args[len(args)-1]
	args = args[:len(args)-1]
	instancesIDs := []*string{}

	for inst := range args {
		instancesIDs = append(instancesIDs, &args[inst])
	}

	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
			{
				Name:   aws.String("instance-id"),
				Values: instancesIDs,
			},
		},
	}

	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("COMMAND: [%s] ********************************\n", command)

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, keys := range inst.Tags {
				switch tag := *keys.Key; tag {
				case "component":
					component = *keys.Value
				case "environment":
					environment = *keys.Value
				}
			}
			if inst.PublicIpAddress != nil {
				IP = *inst.PublicIpAddress
			} else {
				IP = *inst.PrivateIpAddress
			}

			err := runCommand(*inst.InstanceId, component, environment, IP, command)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
}

func runCommand(instanceID, component, environment, IP, command string) error {

	_, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		log.Fatalf("failed to find ssh executable: %s", lookErr)
		return lookErr
	}

	path, err := filepath.Abs(filepath.Join("./ssh-keys"))
	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command(
		"ssh",
		"-i",
		fmt.Sprintf("%s/%s-%s.pem", path, environment, component),
		"-l",
		"ec2-user",
		IP,
		command)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command: %s", err)
		return err
	}

	fmt.Printf("%s:\n %s\n", instanceID, out)

	return nil
}

func connectSSH(environment, component, IP string) error {

	path, err := filepath.Abs(filepath.Join("./ssh-keys"))
	if err != nil {
		os.Exit(1)
	}

	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		log.Fatalf("failed to find ssh executable: %s", err)
		return lookErr
	}

	args := []string{"ssh", "-i", fmt.Sprintf("%s/%s-%s.pem", path, environment, component), "-l", "ec2-user", IP}
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		log.Fatalf("Failed to execute syscall: %s", err)
		return err
	}

	return nil
}
