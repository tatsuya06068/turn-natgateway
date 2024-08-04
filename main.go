package main

import (
    "context"
    "log"
    "os"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

var (
    subnetID       = os.Getenv("SUBNET_ID")
    eipAllocationID = os.Getenv("EIP_ALLOCATION_ID")
)

func handleRequest(ctx context.Context) (string, error) {
    sess := session.Must(session.NewSession())
    svc := ec2.New(sess)

    // Describe NAT Gateways
    input := &ec2.DescribeNatGatewaysInput{}
    result, err := svc.DescribeNatGateways(input)
    if err != nil {
        log.Fatalf("Unable to describe NAT Gateways, %v", err)
    }

    for _, natGateway := range result.NatGateways {
        if *natGateway.State == "available" {
            // Delete NAT Gateway
            delInput := &ec2.DeleteNatGatewayInput{
                NatGatewayId: natGateway.NatGatewayId,
            }
            _, err := svc.DeleteNatGateway(delInput)
            if err != nil {
                log.Fatalf("Unable to delete NAT Gateway, %v", err)
            }

            // Create NAT Gateway
            createInput := &ec2.CreateNatGatewayInput{
                SubnetId:     aws.String(subnetID),
                AllocationId: aws.String(eipAllocationID),
            }
            _, err = svc.CreateNatGateway(createInput)
            if err != nil {
                log.Fatalf("Unable to create NAT Gateway, %v", err)
            }
        }
    }

    return "NAT Gateways restarted successfully", nil
}

func main() {
    lambda.Start(handleRequest)
}