package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

func TestRun(t *testing.T) {
	instanceId := "i-033631851a5922563"
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Errorf("error while creating AWS session err=%s", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)

	result, err := Run(ssmClient, instanceId, "eu-west-1", "sudo id\n")
	if err != nil {
		t.Errorf("error not expected: err=%s", err)
	}

	if result != "blah" {
		t.Errorf("wrong result: result=%s", result)
	}
}