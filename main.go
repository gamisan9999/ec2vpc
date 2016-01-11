package main

/***************************************************************************************************
ec2vpc
機能: 指定したEC2インスタンスIDから、EC2がLaunchされているVPC-IDを表示する

$ ./ec2vpc_linux_amd64 --instance-id `curl -s http://169.254.169.254/latest/meta-data/instance-id`
vpc-XXXXXXXX
***************************************************************************************************/

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
)

func getVpcIDFromInstanceID(svc *ec2.EC2, s string) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(s)},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	if resp.Reservations[0].Instances[0].VpcId != nil {
		fmt.Println(*resp.Reservations[0].Instances[0].VpcId)
	}
}

func getRegionFromInstanceMetaData() (region string) {
	metadata := ec2metadata.New(session.New())
	region, err := metadata.Region()
	if err != nil {
		panic(err)
	}
	return region
}

func main() {
	var instanceID, region, profile string
	app := cli.NewApp()
	app.Name = "ec2vpc"
	app.Usage = "VPC ID from EC2 Instance ID"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "instance-id, i",
			Value:       "",
			Usage:       "--instance-id `http://169.254.169.254/meta-data/latest/instance-id`",
			Destination: &instanceID,
		},
		cli.StringFlag{
			Name:        "profile, p",
			Value:       "",
			Usage:       "--profile default",
			Destination: &profile,
		},
		cli.StringFlag{
			Name:        "region, r",
			Value:       "",
			Usage:       "--region ap-northeast-1",
			Destination: &region,
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("region") == "" {
			region = getRegionFromInstanceMetaData()
		}
		config := &aws.Config{
			// Credentials指定したらIAM Role読まなくなっちゃったよ
			//			Credentials: credentials.NewSharedCredentials("", profile),
			Region: aws.String(region),
		}
		if c.String("instance-id") != "" {
			getVpcIDFromInstanceID(ec2.New(session.New(), config), instanceID)
			os.Exit(0)
		}
	}
	app.Run(os.Args)
}
