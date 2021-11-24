// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
)

func Run(ssmClient *ssm.Client, region string, cmd string) (string, error) {
	input := &ssm.StartSessionInput{
		Target:       aws.String("i-033631851a5922563"),
	}
	sess, err := ssmClient.StartSession(context.TODO(), input)
	if err != nil {
		return "", err
	}

	defer func() {
		if _, err := ssmClient.TerminateSession(context.TODO(), &ssm.TerminateSessionInput{SessionId: sess.SessionId}); err != nil {
			fmt.Printf("unable to terminate session: err=%s", err)
		}
	}()

	sessionToken := *sess.TokenValue

	customEndpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == ssm.ServiceID && region == region {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           fmt.Sprintf("https://%s.%s.amazonaws.com", service, region),
				SigningRegion: region,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	endpoint, err := customEndpointResolver.ResolveEndpoint("ssm", region)
	if err != nil {
		return "", err
	}


	cmdLine := fmt.Sprintf("session-manager-plugin %s %s StartSession %s", sessionToken, region, endpoint.URL)
	e, _, err := expect.Spawn(cmdLine, 60)
	if err != nil {
		return "", err
	}

	shellStart := regexp.MustCompile(`\n\$`)
	if result, _, err := e.Expect(shellStart, 10*time.Second); err != nil {
		return "", fmt.Errorf("did not find shell: err=%s, output=%s", err, result)
	}

	if err := e.Send(cmd); err != nil {
		return "", fmt.Errorf("unable to send command: err=%s", err)
	}

	result, _, err := e.Expect(shellStart, 20*time.Second)
	if err != nil {
		return "", fmt.Errorf("shell did not respond: err=%s", err)
	}

	return result, nil
}