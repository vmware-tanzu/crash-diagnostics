# Copyright (c) 2020 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

conf=crashd_config(workdir="/tmp/crashlogs")
kube_config(path="{0}/.kube/config".format(os.home))
kube_capture(what="logs", namespaces=["default", "cert-manager", "tkg-system"])

# bundle files stored in working dir
archive(output_file="/tmp/craslogs.tar.gz", source_paths=[conf.workdir])