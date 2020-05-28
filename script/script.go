// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"strings"
)

type baseDirective struct {
	index int
	name  string
	raw   string
}

func (d *baseDirective) Index() int {
	return d.index
}

func (d *baseDirective) Name() string {
	return d.name
}

func (d *baseDirective) Raw() string {
	return d.raw
}

// Script is the internal structure of a script
type Script struct {
	directives []Directive
	Preambles  map[string][]Directive // directive commands in the script
	Actions    []Directive            // action commands
}

func New() *Script {
	return &Script{}
}

func (s *Script) AddConfigDirective(index int, name, raw string) *Script {
	dir := s.makeDirective(index, name, raw)
	s.directives = append(s.directives, ConfigDirective(dir))
	return s
}

func (s *Script) AddExecDirective(index int, name, raw string) *Script {
	dir := s.makeDirective(index, name, raw)
	s.directives = append(s.directives, ExecDirective(dir))
	return s
}

// Directives return stored directives in script
func (s *Script) Directives(names ...string) []Directive {
	if len(names) == 0 {
		return s.directives
	}
	var directives []Directive
	for _, dir := range s.directives {
		for _, name := range names {
			if strings.EqualFold(name, dir.Name()) {
				directives = append(directives, dir)
			}
		}
	}
	return directives
}

func (s *Script) ConfigDirectives(names ...string) []ConfigDirective {
	var configs []ConfigDirective
	for _, dir := range s.Directives(names...) {
		switch dir.(type){
		case ConfigDirective:
			configs = append(configs, dir)
		}
	}

	return configs
}


func (s *Script) ExecDirectives(names ...string) []ExecDirective {
	var execs []ExecDirective
	for _, dir := range s.Directives(names...) {
		switch dir.(type){
		case ExecDirective:
			execs = append(execs, dir)
		}
	}

	return execs
}
func (s *Script) makeDirective(index int, name, raw string) Directive {
	return &baseDirective{
		index: index,
		name:  name,
		raw:   raw,
	}
}
