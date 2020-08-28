// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ReadArgsFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("args file not found: %s", path))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	args := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && len(strings.TrimSpace(line)) != 0 {
			if pair := strings.Split(line, "="); len(pair) == 2 {
				args[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
			} else {
				logrus.Warnf("unknown entry in args file: %s", line)
			}
		}
	}

	return args, nil
}
