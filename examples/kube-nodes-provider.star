# Copyright (c) 2020 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script shows how to use the kube nodes provider.
# The kube node provider uses the Kubernetes Nodes objects
# to enumerate compute resources that are part of the cluster.
# It uses SSH to execute commands on those on nodes.
#
# This example requires an SSH and a Kubernetes cluster.

# setup and configuration
ssh=ssh_config(
    username=args.username,
    private_key_path=args.key_path,
    port=args.ssh_port,
    max_retries=5,
)

hosts=resources(
    provider=kube_nodes_provider(
        kube_config=kube_config(path=args.kubecfg),
        ssh_config=ssh,
    ),
)

# commands to run on each host
uptimes = run(cmd="uptime", resources=hosts)

# result for resource 0 (localhost)
print(uptimes.result)