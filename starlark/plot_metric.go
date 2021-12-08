// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/metric"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"go.starlark.net/starlark"
)

// plotMetric is a built-in starlark function that plots a metric returned from an OpenTelemetry endpoint.
// It takes an array of metric names (and optionally, client certificate, server key, and endpoint), checks
// a list of known metrics and uses capture() to ssh onto the provided list of nodes curl the endpoints
// and plot the metrics into pngs. Starlark format: plot_metric(prog=<prog_name>, resources=nodes)
func plotMetric(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var workdir, fileName, desc, clientKey, clientCert, endpoint, cmdStr string
	var metricNames, resources *starlark.List

	if err := starlark.UnpackArgs(
		identifiers.capture, args, kwargs,
		"resources?", &resources,
		"workdir?", &workdir,
		"file_name?", &fileName,
		"desc?", &desc,
		"metrics", &metricNames,
		"clientCert?", &clientCert,
		"clientKey?", &clientKey,
		"endpoint?", &endpoint,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.capture, err)
	}
	if len(workdir) == 0 {
		if dir, err := getWorkdirFromThread(thread); err == nil {
			workdir = dir
		}
	}
	if len(workdir) == 0 {
		workdir = defaults.workdir
	}

	if resources == nil {
		res, err := getResourcesFromThread(thread)
		if err != nil {
			return starlark.None, fmt.Errorf("%s: %s", identifiers.copyFrom, err)
		}
		resources = res
	}

	for _, m := range toSlice(metricNames) {
		var agent ssh.Agent
		var ok bool
		if agentVal := thread.Local(identifiers.sshAgent); agentVal != nil {
			agent, ok = agentVal.(ssh.Agent)
			if !ok {
				return starlark.None, errors.New("unable to fetch ssh-agent")
			}
		}
		switch m {
		case "etcd":
			client := metric.GetClient(m, clientKey, clientCert, endpoint, workdir)
			cmdStr = client.GetCommandOutput()
		default:
			logrus.Errorf("Unknown Metric %s", m)
		}

		results, err := execCapture(cmdStr, workdir, fileName, desc, agent, resources)
		if err != nil {
			return starlark.None, fmt.Errorf("%s: %s", identifiers.capture, err)
		}

		// build list of struct as result
		var resultList []starlark.Value
		for _, result := range results {
			if len(results) == 1 {
				return result.toStarlarkStruct(), nil
			}
			resultList = append(resultList, result.toStarlarkStruct())
		}

	}

	//results ,err := execPlotMetrics(metricNames, clientCert, clientKey, endpoint, resource)
	//
	//if err != nil {
	//	return starlark.None, err
	//}

	// build list of struct as result
	//var resultList []starlark.Value
	//	if len(results) == 1 {
	//		return result.toStarlarkStruct(), nil
	//	}
	//	resultList = append(resultList, result.toStarlarkStruct())
	//}
	return nil, nil
	//return starlark.NewList(resultList), nil
}

//func execPlotMetrics(metricNames, clientCert, clientKey, endpoint, resources, workdir) ([]commandResult,error){
//	for _, metricNames {
//
//
//	}
//
//}
