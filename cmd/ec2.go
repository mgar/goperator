// Copyright © 2017 Miguel Ángel García <mgarcia.inf@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	ec2_service "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mgar/goperator/aws/ec2"
	"github.com/mgar/goperator/client"
	"github.com/mgar/goperator/formatter"
	"github.com/spf13/cobra"
)

var ec2Client = client.NewEc2Client()

// listCmd represents the list command
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

var cmdTerminateInstance = &cobra.Command{
	Use:   "terminate [instance-id]",
	Long:  "Terminate one or more EC2 instances given [instance-id ...]",
	Short: "Terminate one or many EC2 instances",
	Run:   terminateInstance,
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

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Manage Ec2 resources",
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(cmdListInstances)
	ec2Cmd.AddCommand(cmdSSHInstance)
	ec2Cmd.AddCommand(cmdTerminateInstance)
	ec2Cmd.AddCommand(cmdStopInstance)
	ec2Cmd.AddCommand(cmdStartInstance)
	ec2Cmd.AddCommand(execCommandInInstance)

	// Here you will define your flags and configuration settings.
	var serviceName string
	cmdListInstances.Flags().StringVarP(&serviceName, "service", "s", "", "Filter by service name")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func listInstances(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		fmt.Println("You need to specify [environment] and [component]")
		os.Exit(1)
	}

	service := cmd.Flag("service").Value.String()
	environment, component := args[0], args[1]
	instances := []ec2.Ec2Instance{}
	filter := []*ec2_service.Filter{}

	if service != "" {
		filter = []*ec2_service.Filter{
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
				Name: aws.String("tag:service"),
				Values: []*string{
					aws.String(service),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("stopped"), aws.String("shutting-down"),
					aws.String("stopping"), aws.String("pending")},
			},
		}
	} else {
		filter = []*ec2_service.Filter{
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
		}
	}

	params := &ec2_service.DescribeInstancesInput{
		DryRun:  aws.Bool(false),
		Filters: filter,
	}

	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp.Reservations != nil {
		for _, res := range resp.Reservations {
			instance, err := ec2.NewEc2Instance(res.Instances[0])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			instances = append(instances, *instance)
		}
	}

	formatter.Ec2ToTable(instances)
}

func sshInstance(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id]")
		os.Exit(1)
	}
	instanceID := args[0]

	params := &ec2_service.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2_service.Filter{
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

	instance, err := ec2.NewEc2Instance(resp.Reservations[0].Instances[0])
	if err != nil {

	}
	err = connectSSH(instance)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func terminateInstance(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id ...]")
		os.Exit(1)
	}
	instances := stringsSliceToStringPointersSlice(args)

	params := &ec2_service.TerminateInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: instances,
	}
	_, err := ec2Client.TerminateInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func stopInstance(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("You need to specify [instance-id ...]")
		os.Exit(1)
	}

	instances := stringsSliceToStringPointersSlice(args)

	params := &ec2_service.StopInstancesInput{
		DryRun:      aws.Bool(false),
		InstanceIds: instances,
	}

	resp, err := ec2Client.StopInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	for idx := range resp.StoppingInstances {
		if *resp.StoppingInstances[idx].PreviousState.Name != "running" {
			fmt.Printf("Instance %s is not running. It can't be stopped.\n", *resp.StoppingInstances[idx].InstanceId)
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

	instances := stringsSliceToStringPointersSlice(args)

	params := &ec2_service.StartInstancesInput{
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

	command := args[len(args)-1]
	args = args[:len(args)-1]

	instancesList := stringsSliceToStringPointersSlice(args)

	params := &ec2_service.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2_service.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
			{
				Name:   aws.String("instance-id"),
				Values: instancesList,
			},
		},
	}

	resp, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	instances := []ec2.Ec2Instance{}

	if resp.Reservations != nil {
		for _, res := range resp.Reservations {
			instance, err := ec2.NewEc2Instance(res.Instances[0])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			instances = append(instances, *instance)
		}
	}

	fmt.Printf("COMMAND: [%s] ********************************\n", command)

	for _, instance := range instances {
		err := runCommand(&instance, command)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func runCommand(inst *ec2.Ec2Instance, command string) error {

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
		fmt.Sprintf("%s/%s.pem", path, inst.Key),
		"-l",
		"ec2-user",
		inst.PublicIP,
		command)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command: %s", err)
		return err
	}

	fmt.Printf("%s:\n %s\n", inst.InstanceID, out)

	return nil
}

func connectSSH(instance *ec2.Ec2Instance) error {

	path, err := filepath.Abs(filepath.Join("./ssh-keys"))
	if err != nil {
		os.Exit(1)
	}

	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		log.Fatalf("failed to find ssh executable: %s", err)
		return lookErr
	}

	args := []string{"ssh", "-i", fmt.Sprintf("%s/%s.pem", path, instance.Key), "-l", "ec2-user", instance.PublicIP}
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		log.Fatalf("Failed to execute syscall: %s", err)
		return err
	}

	return nil
}

func stringsSliceToStringPointersSlice(stringsSlice []string) (stringsPointersSlice []*string) {

	stringsPointersSlice = []*string{}

	for item := range stringsSlice {
		stringsPointersSlice = append(stringsPointersSlice, &stringsSlice[item])
	}

	return stringsPointersSlice
}
