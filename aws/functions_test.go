package aws_test

import (
	"testing"

	aws_helpers "github.com/ren3gadem4rm0t/slack-go-helpers/aws"
)

func TestAWSAccountFromAWSKeyID(t *testing.T) {
	tests := []struct {
		awsKeyID    string
		want        string
		expectError bool
	}{
		{"AKIAIKAZ2VEXAMPLE", "571284826414", false},
		{"AKI", "", true},
	}
	for _, tt := range tests {
		got, err := aws_helpers.AWSAccountFromAWSKeyID(tt.awsKeyID)
		if (err != nil) != tt.expectError {
			t.Errorf("expected error=%v, got=%v", tt.expectError, err != nil)
			continue
		}
		if !tt.expectError && got != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}

func TestAWSResourceTypeFromPrefix(t *testing.T) {
	tests := []struct {
		awsKeyID    string
		want        string
		expectError bool
	}{
		{"AKIA", "Access key", false},
		{"AIDA", "IAM user", false},
		{"ABIA", "AWS STS service bearer token", false},
		{"XXXX", "Unknown resource type", false},
		{"A", "", true},
	}
	for _, tt := range tests {
		got, err := aws_helpers.AWSResourceTypeFromPrefix(tt.awsKeyID)
		if (err != nil) != tt.expectError {
			t.Errorf("expected error=%v, got=%v", tt.expectError, err != nil)
			continue
		}
		if !tt.expectError && got != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}
