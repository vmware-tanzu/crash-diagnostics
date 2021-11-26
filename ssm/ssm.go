// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"strings"
	"time"
)

const (
	DocumentAWSRunShellScript = "AWS-RunShellScript"
)

// ssmClientAPI is a proxy to the AWS SDK. Tests should mock this.
type ssmClientAPI interface {
	SendCommand(context.Context, *ssm.SendCommandInput, ...func(*ssm.Options)) (*ssm.SendCommandOutput, error)
	GetClient() *ssm.Client
}

type SSMClient struct {
	client *ssm.Client
}

func (s *SSMClient) SendCommand(ctx context.Context, input *ssm.SendCommandInput, opts ...func(*ssm.Options)) (*ssm.SendCommandOutput, error)  {
	return s.client.SendCommand(ctx, input, opts...)
}

func (s *SSMClient) GetClient() *ssm.Client {
	return s.client
}

func Run(ctx context.Context, ssmClient ssmClientAPI, instanceId string, cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	input := &ssm.SendCommandInput{
		DocumentName:           aws.String(DocumentAWSRunShellScript),
		Comment:                aws.String(fmt.Sprintf("crash-diagnostic run at %s", time.Now().String())),
		DocumentHash:           nil,
		DocumentVersion:        nil,
		InstanceIds:            []string{instanceId},
		MaxConcurrency:         nil,
		MaxErrors:              nil,
		Parameters: map[string][]string{
			"commands": {cmd},
		},
		ServiceRoleArn:         nil,
		TimeoutSeconds:         60,
	}

	command, err := ssmClient.SendCommand(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error sending command err=%s", err)
	}

	waiter := ssm.NewCommandExecutedWaiter(ssmClient.GetClient())
	waiterInput := &ssm.GetCommandInvocationInput{
		CommandId:  command.Command.CommandId,
		InstanceId: aws.String(instanceId),
	}
	waiterOutput, err := waiter.WaitForOutput(ctx, waiterInput, 60*time.Second)
	if err != nil {
		return "", fmt.Errorf("error wating for output of the command: cmd=%s, err=%s", cmd, err)
	}

	// TODO: get output from S3 as this is onle first 24k characters of the output.
	return *waiterOutput.StandardOutputContent, nil
}