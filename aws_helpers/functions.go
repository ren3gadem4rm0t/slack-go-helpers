package aws_helpers

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"math/big"
)

// AWSAccountFromAWSKeyID decodes an AWS Key ID to extract the account ID.
//
// Some code referenced from: https://medium.com/@TalBeerySec/a-short-note-on-aws-key-id-f88cc4317489
func AWSAccountFromAWSKeyID(awsKeyID string) (string, error) {
	if len(awsKeyID) <= 4 {
		return "", fmt.Errorf("AWSKeyID is too short")
	}
	raw := awsKeyID[4:]
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("failed to decode base32: %w", err)
	}
	if len(decoded) < 6 {
		return "", fmt.Errorf("decoded AWSKeyID is too short")
	}
	z := new(big.Int).SetBytes(decoded[:6])

	maskBytes, err := hex.DecodeString("7fffffffff80")
	if err != nil {
		return "", fmt.Errorf("failed to decode mask: %w", err)
	}
	mask := new(big.Int).SetBytes(maskBytes)

	z.And(z, mask)
	z.Rsh(z, 7)

	return fmt.Sprintf("%012d", z), nil
}

// AWSResourceTypeFromPrefix returns the AWS resource type based on the Key ID prefix.
//
// Prefixes referenced from: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-unique-ids
func AWSResourceTypeFromPrefix(awsKeyID string) (string, error) {
	if len(awsKeyID) < 4 {
		return "", fmt.Errorf("AWSKeyID is too short to determine prefix")
	}
	prefix := awsKeyID[:4]
	resourceTypes := map[string]string{
		"ABIA": "AWS STS service bearer token",
		"ACCA": "Context-specific credential",
		"AGPA": "User group",
		"AIDA": "IAM user",
		"AIPA": "Amazon EC2 instance profile",
		"AKIA": "Access key",
		"ANPA": "Managed policy",
		"ANVA": "Version in a managed policy",
		"APKA": "Public key",
		"AROA": "Role",
		"ASCA": "Certificate",
		"ASIA": "Temporary (AWS STS) access key ID",
	}
	if t, ok := resourceTypes[prefix]; ok {
		return t, nil
	}
	return "Unknown resource type", nil
}
