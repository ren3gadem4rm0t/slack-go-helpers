package main

import (
	"fmt"
	"log"

	"github.com/ren3gadem4rm0t/slack-go-helpers/aws"
)

func main() {
	// Example AWS Key IDs with different prefixes
	awsKeyIDs := []string{
		"AKIAEXAMPLEKEYID1234",
		"AIDAEXAMPLEUSERID5678",
		"ASIAEXAMPLETEMPKEY9012",
		"UNKNOWNPREFIX1234",
	}

	for _, keyID := range awsKeyIDs {
		resourceType, err := aws.AWSResourceTypeFromPrefix(keyID)
		if err != nil {
			log.Printf("Error determining resource type for Key ID '%s': %v\n", keyID, err)
			continue
		}
		fmt.Printf("AWS Key ID: %s => Resource Type: %s\n", keyID, resourceType)
	}
}
