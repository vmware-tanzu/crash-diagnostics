// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bufio"
	"regexp"
	"strings"
	"testing"

	"github.com/vladimirvivien/echo"
)

func TestParseAndValidateAgentInfo(t *testing.T) {
	tests := []struct {
		name      string
		info      string
		shouldErr bool
	}{
		{
			name:      "valid info",
			shouldErr: false,
			info: `SSH_AUTH_SOCK=/foo/bar.1234; export SSH_AUTH_SOCK;
SSH_AGENT_PID=4567; export SSH_AGENT_PID;
echo Agent pid 4567;`,
		},
		{
			name:      "invalid info",
			shouldErr: true,
			info: `FOO=/foo/bar.1234; export BAR;
BLAH=4567; export BLOOP;
echo lorem ipsum 4567;`,
		},
		{
			name:      "invalid info",
			shouldErr: true,
			info: `SSH_AUTH_SOCK=/foo/bar.1234; export SSH_AUTH_SOCK;
BLAH=4567; export BLOOP;
echo lorem ipsum 4567;`,
		},
		{
			name:      "invalid info",
			shouldErr: true,
			info: `FOO=/foo/bar.1234; export BAR;
SSH_AGENT_PID=4567; export SSH_AGENT_PID;
echo lorem ipsum 4567;`,
		},
		{
			name:      "invalid info",
			shouldErr: true,
			info: `lorem ipsum 1;
lorem ipsum 2.`,
		},
		{
			name:      "invalid info",
			shouldErr: true,
			info:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			agentInfo, err := parseAgentInfo(strings.NewReader(test.info))
			if err != nil {
				t.Fail()
			}
			err = validateAgentInfo(agentInfo)
			if err != nil && !test.shouldErr {
				// unexpected failures
				t.Fail()
			} else if !test.shouldErr {
				if _, ok := agentInfo[AgentPidIdentifier]; !ok {
					t.Fail()
				}
				if _, ok := agentInfo[AuthSockIdentifier]; !ok {
					t.Fail()
				}
			} else {
				// asserting error scenarios
				if err == nil {
					t.Fail()
				}
			}
		})
	}
}

func TestStartAgent(t *testing.T) {
	a, err := StartAgent()
	if err != nil || a == nil {
		t.Fatalf("error should be nil and agent should not be nil: %v", err)
	}
	out := echo.New().Run("ps -ax")
	if !strings.Contains(out, "ssh-agent") {
		t.Fatal("no ssh-agent process found")
	}

	failed := true
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ssh-agent") {
			pid := strings.Split(strings.TrimSpace(line), " ")[0]
			// set failed to false if correct ssh-agent process is found
			agentStruct, _ := a.(*agent)
			if pid == agentStruct.Pid {
				failed = false
			}
		}
	}
	if failed {
		t.Fatal("could not find agent with correct Pid")
	}

	t.Cleanup(func() {
		_ = a.Stop()
	})
}

func TestAgent(t *testing.T) {
	a, err := StartAgent()
	if err != nil {
		t.Fatalf("failed to start agent: %v", err)
	}

	tests := []struct {
		name   string
		assert func(*testing.T, Agent)
	}{
		{
			name: "GetEnvVariables",
			assert: func(t *testing.T, agent Agent) {
				vars := agent.GetEnvVariables()
				if len(strings.Split(vars, " ")) != 2 {
					t.Fatalf("not enough variables")
				}

				match, err := regexp.MatchString(`SSH_AGENT_PID=[0-9]+ SSH_AUTH_SOCK=\S*`, vars)
				if err != nil || !match {
					t.Fatalf("format does not match")
				}
			},
		},
		{
			name: "Stop",
			assert: func(t *testing.T, agent Agent) {
				if err := agent.Stop(); err != nil {
					t.Errorf("failed to stop agent: %s", err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.assert(t, a)
		})
	}
}
