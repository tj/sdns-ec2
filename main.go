package main

import "github.com/awslabs/aws-sdk-go/aws/credentials"
import "github.com/awslabs/aws-sdk-go/service/ec2"
import "github.com/awslabs/aws-sdk-go/aws"
import "github.com/tj/docopt"
import "github.com/tj/sdns"
import "strconv"
import "log"
import "os"

var Version = "0.0.1"

const Usage = `
  Usage:
    sdns-ec2 [--region name] [--ttl n]
    sdns-ec2 -h | --help
    sdns-ec2 --version

  Options:
    --ttl n          record ttl [default: 300]
    --region name    aws region [default: us-west-2]
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	log := log.New(os.Stderr, "", 0)

	args, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	question, err := sdns.Read(os.Stdin)
	if err != nil {
		log.Fatalf("error parsing stdin: %s", err)
	}

	ttl, err := strconv.Atoi(args["--ttl"].(string))
	if err != nil {
		log.Fatalf("error parsing --ttl: %s", err)
	}

	client := ec2.New(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region:      args["--region"].(string),
	})

	res, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(question.Name),
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if len(res.Reservations) == 0 {
		log.Fatalf("error: failed to find a matching host")
	}

	if len(res.Reservations[0].Instances) == 0 {
		log.Fatalf("error: failed to find a matching host")
	}

	ip := res.Reservations[0].Instances[0].PrivateIPAddress

	answer := sdns.Answers{
		&sdns.Answer{
			Type:  "A",
			Value: *ip,
			TTL:   uint32(ttl),
		},
	}

	sdns.Write(answer, os.Stdout)
}
