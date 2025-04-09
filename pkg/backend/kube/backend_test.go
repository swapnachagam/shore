package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	testLog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var logger *logrus.Logger

func TestSavePipeline(t *testing.T) {
	// Arrange
	// Arrange
	logger, _ = testLog.NewNullLogger()

	client := &SpinClient{log: logger}

	pipelineJSON := `{
   "iamPolicy": {
      "apiVersion": "iam.services.k8s.aws/v1alpha1",
      "kind": "Policy",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1-policy",
         "namespace": "default"
      },
      "spec": {
         "name": "log-forwarder-1-policy",
         "policyDocument": {
            "Statement": [
               {
                  "Action": [
                     "logs:PutSubscriptionFilter",
                     "logs:DeleteSubscriptionFilter",
                     "s3:PutObject",
                     "s3:GetObject"
                  ],
                  "Effect": "Allow",
                  "Resource": "*"
               }
            ],
            "Version": "2012-10-17"
         }
      }
   },
   "logForwarder1": {
      "apiVersion": "logforwarding.ucp.adskeng.net/v1alpha1",
      "kind": "CloudWatchLogForwarder",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1",
         "namespace": "default"
      },
      "spec": {
         "account": "123456789012",
         "logGroupNames": [
            "log-group-1",
            "log-group-2"
         ],
         "moniker": "log-forwarder-1",
         "region": "us-west-2"
      }
   },
   "s3Bucket": {
      "apiVersion": "s3.services.k8s.aws/v1alpha1",
      "kind": "Bucket",
      "metadata": {
         "labels": {
            "app": "log-forwarding"
         },
         "name": "log-forwarder-1-failed-delivery"
      },
      "spec": {
         "name": "log-forwarder-1-failed-delivery"
      }
   }
}`
	folderPath := "./generated"

	// Act
	_, err := client.SavePipeline(pipelineJSON)

	// Assert
	assert.NoError(t, err)

	// Verify the files exist
	expectedFiles := []string{"iamPolicy.yaml", "logForwarder1.yaml", "s3Bucket.yaml"}
	for _, fileName := range expectedFiles {
		filePath := filepath.Join(folderPath, fileName)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, fmt.Sprintf("Expected the file %s to be created", fileName))
	}

	// Verify the content of one of the files (e.g., iamPolicy.yaml)
	iamPolicyPath := filepath.Join(folderPath, "iamPolicy.yaml")
	content, err := os.ReadFile(iamPolicyPath)
	assert.NoError(t, err, "Expected the YAML file to be readable")

	expectedYAML := `apiVersion: iam.services.k8s.aws/v1alpha1
kind: Policy
metadata:
  labels:
    app: log-forwarding
  name: log-forwarder-1-policy
  namespace: default
spec:
  name: log-forwarder-1-policy
  policyDocument:
    Statement:
    - Action:
      - logs:PutSubscriptionFilter
      - logs:DeleteSubscriptionFilter
      - s3:PutObject
      - s3:GetObject
      Effect: Allow
      Resource: '*'
    Version: "2012-10-17"
`
	assert.Equal(t, expectedYAML, string(content), "Expected the YAML content to match")

	// Clean up
	os.RemoveAll(folderPath)
}
