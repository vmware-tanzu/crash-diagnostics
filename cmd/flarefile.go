// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

const flarefile = `FROM local
WORKDIR /tmp/flareout
COPY /var/log/messages
COPY /var/log/syslog
COPY /var/log/system.log
`
