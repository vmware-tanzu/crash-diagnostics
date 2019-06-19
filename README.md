# Flare

A too to collect machine environment data to investigate issues with a Kubernetes cluster.

## Pre-requisites
 * Linux
 * Access to cluster nodes

## Usage

```shell
flare help

Usage: flare <command>

Commands:
  collect    Collects state and environment data from local machine.
  show       Display summary view of state and environment data from local machine
```

### bundle collect
```shell
flare help collect 

Usage: flare collect <options>
Collects state and environment data from local machine.

options:
  --all        Collect all state and environment data from machine
  --k8s-info   Collects Kubernetes related data including infrastructure and logs
  --node-info  Collects machine related infrastructure and logs.
  --output     Specifies the location of generated tarball.
```

### bundle show
```shell
flare help view

Usage: flare view
Displays a summary of Kubnernetes and infrastructure information from local machine.
```

### Example

```shell
$> flare collect --all --output /tmp/node01.gzip
```