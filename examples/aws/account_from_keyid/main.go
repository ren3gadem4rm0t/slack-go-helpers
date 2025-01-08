package main

import (
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/aws"
)

func main() {
	// Example AWS Key ID
	awsKeyID := "AKIAIKAZ2VEXAMPLE"

	// Extract the AWS Account ID from the AWS Key ID
	accountID, err := aws.AWSAccountFromAWSKeyID(awsKeyID)
	if err != nil {
		log.Fatalf("Error extracting AWS Account ID: %v", err)
	}

	fmt.Printf("AWS Account ID: %s\n", accountID)
}
