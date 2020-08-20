# Copyright (c) 2020 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

conf=crashd_config(workdir=args.workdir)
kube_capture(what="logs", namespaces=["default", "cert-manager", "tkg-system"], kube_config = kube_config(path=args.kubecfg))

# bundle files stored in working dir
archive(output_file=args.output, source_paths=[args.workdir])
