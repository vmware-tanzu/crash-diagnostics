package ssm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"net/url"
	"strings"
	"testing"
)

func TestRunDocument(t *testing.T) {
	instanceId := "i-033631851a5922563"
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Errorf("error while creating AWS session err=%s", err)
	}

	originalSSMClient := ssm.NewFromConfig(cfg)
	ssmClientStruct := &SSMClient{
		Client: originalSSMClient,
	}

	result, err := Run(ctx, ssmClientStruct, instanceId, "sudo crictl images")
	if err != nil {
		t.Errorf("error not expected: err=%s", err)
	}

	if result != "blah" {
		t.Errorf("wrong result: result=%s", result)
	}
}

func TestURL(t *testing.T) {
	testUrl := "aws:///eu-west-1a/i-075404133d12be157"

	u, err := url.Parse(testUrl)
	if err != nil {
		t.Fatal(err)
	}

	idList := strings.Split(u.Path, "/")

	t.Log(idList[len(idList)-1])
}