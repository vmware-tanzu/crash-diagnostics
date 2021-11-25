// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/sirupsen/logrus"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	expect "github.com/Netflix/go-expect"
)

func expectNoError() expect.ConsoleOpt {
	return expect.WithExpectObserver(
		func(matchers []expect.Matcher, buf string, err error) {
			if err == nil {
				return
			}

			if len(matchers) == 0 {
				logrus.Fatalf("Error occurred while matching %q: %s\n%s", buf, err, string(debug.Stack()))
			} else {
				var criteria []string
				for _, matcher := range matchers {
					criteria = append(criteria, fmt.Sprintf("%q", matcher.Criteria()))
				}
				logrus.Fatalf("Failed to find [%s] in %q: %s\n%s", strings.Join(criteria, ", "), buf, err, string(debug.Stack()))
			}

		})
}

func sendNoError() expect.ConsoleOpt {
	return expect.WithSendObserver(
		func(msg string, num int, err error) {
			if err != nil {
				logrus.Fatalf("Failed to send %q: %s\n%s", msg, err, string(debug.Stack()))
			}
			if len(msg) != num {
				logrus.Fatalf("Only sent %d of %d bytes for %q\n%s", num, len(msg), msg, string(debug.Stack()))
			}
		})
}

func Run(ssmClient *ssm.Client, instanceId string, region string, cmd string) (string, error) {
	input := &ssm.StartSessionInput{
		Target:       aws.String(instanceId),
	}
	sess, err := ssmClient.StartSession(context.TODO(), input)
	if err != nil {
		return "", err
	}


	inputJson, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("error marshaling the session input err=%s", err)
	}
	sessionToken, err := json.Marshal(sess)
	if err != nil {
		return "", fmt.Errorf("cannot marshal session: err=%s", err)
	}

	defer func() {
		if _, err := ssmClient.TerminateSession(context.TODO(), &ssm.TerminateSessionInput{SessionId: sess.SessionId}); err != nil {
			fmt.Printf("unable to terminate session: err=%s", err)
		}
	}()

//	cmdLine := fmt.Sprintf(`/Users/dzyla/bin/session-manager-plugin '%s' '%s' StartSession '%s' '%s'`, sessionToken, region, string(inputJson), "https://ssm.eu-west-1.amazonaws.com")
	mainCommand := "session-manager-plugin"
	cmdLine := []string{
		fmt.Sprintf("%s", sessionToken),
		fmt.Sprintf("%s", region),
		"StartSession",
		"",
		fmt.Sprintf("%s", string(inputJson)),
		"https://ssm.eu-west-1.amazonaws.com",
	}

	var buf bytes.Buffer
	consoleOpts := []expect.ConsoleOpt{
		expectNoError(),
		sendNoError(),
		expect.WithDefaultTimeout(5*time.Second),
		expect.WithStdout(&buf),
	}
	c, err := expect.NewConsole(consoleOpts...)
	if err != nil {
		logrus.Errorf("cannot create an expect console: err=%s", err)
	}

	command := exec.Command(mainCommand, cmdLine...)
	command.Stdin = c.Tty()
	command.Stdout = c.Tty()
	command.Stderr = c.Tty()

	shellStart := regexp.MustCompile(`\n\$`)

	//var out string
	//var i int

	go func() {
		c.Expect(expect.Regexp(shellStart))
		c.Send(fmt.Sprintf("%s\n", cmd))
		c.Expect(expect.String(""))
		c.Send("exit\n")
		c.ExpectEOF()
	}()

	if err := command.Start(); err != nil {
		return "", fmt.Errorf("cannot start the command: err=%s, cmd=%s", err, cmdLine)
	}

	if err := command.Wait(); err != nil {
		return "", fmt.Errorf("error waiting for the command: err=%s", err)
	}

	return strings.TrimSpace(buf.String()), nil
}