# Copyright (c) 2020 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

conf=crashd_config(workdir="/tmp/crashobjs")
nspaces=[
    "capi-kubeadm-bootstrap-system",
    "capi-kubeadm-control-plane-system",
    "capi-system capi-webhook-system",
    "capv-system capa-system",
    "cert-manager tkg-system",
]

set_as_default(kube_config = kube_config(path=args.kubecfg))

# capture Kubernetes API object and store in files (under working dir)
kube_capture(what="objects", kinds=["services", "pods"], namespaces=nspaces)
kube_capture(what="objects", kinds=["deployments", "replicasets"], namespaces=nspaces)
kube_capture(what="objects", kinds=["clusters", "machines", "machinesets", "machinedeployments"], namespaces="tkg-system")

# bundle files stored in working dir
archive(output_file="/tmp/crashobjs.tar.gz", source_paths=[conf.workdir])