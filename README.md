# CRE Bundle Tool

Tool to collect machine environment data from a Kubernetes cluster.

## Pre-requisites
 * Linux
 * Access to cluster nodes

## Usage

```shell
bundle help

Usage: bundle <command>

Commands:
  collect    Collects state and environment data from local machine.
  show       Display summary view of state and environment data from local machine
```

### bundle collect
```shell
bundle help collect 

Usage: bundle collect <options>
Collects state and environment data from local machine.

options:
  --all        Collect all state and environment data from machine
  --k8s-info   Collects Kubernetes related data including infrastructure and logs
  --node-info  Collects machine related infrastructure and logs.
  --output     Specifies the location of generated tarball.
```

### bundle show
```shell
bundle help view

Usage: bundle view
Displays a summary of Kubnernetes and infrastructure information from local machine.
```

### Example

```shell
$> bundle collect --all --output /tmp/node01.gzip
```