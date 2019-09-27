# Crash Recovery and Diagnostics for Kubernetes
A tool to help collect and analyze node information for troubleshooting unresponsive Kubernetes clusters.

This tool is designed for human operators to run against some or all nodes of a troubled or unresponsive cluster to collect logs, configurations, or any other node resource for analysis to help with troubleshotting of Kubernetes.  

## Pre-requisites
 * Assumes the Kubernetes cluster is running on Linux

## `Crash-Diagnostics`
The tool is compiled into a a single binary named `crash-dianostics`.  Currently, the binary supports two commands:

```
Usage:
  crash-diagnotics [command]

Available Commands:
  help        Help about any command
  run         Executes a diagnostics script file
```

Command `run` uses a diagnostics file to script how and what resources to are collected from cluster machines. By default, `crash-diagnostics run` searches for script for `Diagnostics.file` which specifies line-by-line directives and commands that are interpreted into actions to be executed against the nodes in the cluster.

```
> crash-diagnostics run --help

Usage:
  crash-diagnotics run [flags]

Flags:
      --file string     the path to the dianostics script file to run (default "Diagnostics.file")
      --output string   the path of the generated archive file (default "out.tar.gz")
```


For instance, the following command will execte file `./Diagnostics.file` and store any collected data in file `out.tar.gz`:

```
crash-diagnostics run
```

To run a different script file  or specify a different output archive, use the flags shown below: 

```
crash-diagnostics --file test-cluster.file -output test-cluster.tar.gz
```

## Diagnotics.file Format
`Diagnostics.file` uses a simple line-by-line format (à la Dockerfile) to specify directives on how to collect data from cluster servers:

```
DIRECTIVE [arguments]
```

A directive can either be a `preamble` for runtime configuration or an `action` which can execute a command that runs on each specified node in a cluster. 

### Example Diagnostics.file
```
FROM 127.0.0.1:22 192.168.99.7:22
WORKDIR /tmp/crashdir

COPY /var/log/kube-apiserver.log
CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i
CAPTURE journalctl -l -u kube-apiserver
COPY /var/log/kubelet.log
COPY /var/log/kube-proxy.log

OUTPUT path:/tmp/crashout/out.tar.gzip

```
In the previous example, the tool will collect information from servers `127.0.0.1:22` and `192.168.99.7:22` by executing the COPY and the CAPTURE 
commands specified in the file.  The colleted information is bundled into archive file `/tmp/crashout/out.tar.gzip` specified by `OUTPUT` (note that 
the output file can also be specified by flag `--output`).

## Diagnostics.file Directives
Currently, `crash-dianostics` supports the following directives:
```
AS
AUTHCONFIG
CAPTURE
COPY
ENV
FROM
KUBECONFIG
OUTPUT
WORKDIR
```

### AS
This directive specifies the userid id and optional group id to  use when `crash-diagnostics` execute its commands against the current machine.
```
AS <userid>[:<groupid>]
```

### AUTHCONFIG
Configures the authentication for connections to remote node servers with a username and a private key that is used
in a keypair configuration for tools such as SSH.

```
AUTHCONFIG username:<name> private-key:</path/to/private/key>
```

### CAPTURE
Executes a shell command on the specified machines (see FROM) and captures the output in the archived bundle for analysis.

```
CAPTURE [<shell command>]
```

### COPY
This action specifies one or more files (and/or directories) as data sources that are copied
into the arhive bundle.

```
COPY <space-separated files or directories>
```

### ENV
This directive is used to inject environment variables that to be processed by shell commands executed by the CAPTURE action at runtime.

```
ENV key0=value0
ENV key1=value1
ENV keyN=valueN
```

### FROM
This specifies a space-separated list of nodes from which data can be collected.  Each 
machine is specified by an address and a service port as `<host-address>:<port>`. By default
the tool will use SSH as at runtime to interact with the specified remote hosts.

An address of `local` indicates that the current machine, where the `crash-diagnostics` binary, 
is running will be used as the source allowing the tool to directly access and execute commands.

```
FROM local 10.10.100.2:22
```

### KUBECONFIG
This directive specifies the fully qualified path of the Kubernetes client configuration file. The
tool will use this value to communicate with the API server to retrieve vital cluster information 
if available.  

```
KUBECONFIG /path/to/kube/config
```

By default, the following resourcess will be retrieved from the API server:

 * Namespaces
 * Nodes
 * Events
 * Replication Controllers
 * Services
 * DaemonSets
 * Deployments
 * ReplicaSets
 * Pods

If `KUBECONFIG` is not specified, the tool will attempt to search for:
 * Environment variable `KUBECONFIG`
 * If the `KUBECONFIG` env variable is not set, path $HOME/.kube/config will be used

If a Kubernetes configuration file is not found or the API server is unresponsive, cluster information will be skipped.

### OUTPUT
This preamble configures the the location and file name of the generated archive file that contains the collected information
from the specified servers.

```
OUTPUT path:<path of archive file>
```

If `OUTPUT` is not specified, the tool applies the value of flag `--output` as specified on the command line at runtime.

### WORKDIR
Specifies the working directory used when building the archive bundle.  The
directory is  used as temporary location to store data from all data sources
specified in the file.  When the tar is built, the content of that directory
is removed.

```
WORKDIR <relative or absolute path>
```

### Example File

```
FROM local 162.164.10.1:2222 162.164.10.2:2222
KUBECONFIG /home/username/.kube/kind-config-kind
AUTHCONFIG username:test private-key:/home/testuser/.ssh/id_rsa
WORKDIR /tmp/output

CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i

OUTPUT path:/tmp/crashout/out.tar.gz
```

### Comments
Each line that starts with with `#` is considered to be a comment and is ignored at runtime as shown in 
the following example:

```
# This shows how to comment your script
FROM local 162.164.10.1:2222 162.164.10.2:2222
KUBECONFIG /home/username/.kube/kind-config-kind
AUTHCONFIG username:test private-key:/home/testuser/.ssh/id_rsa
WORKDIR /tmp/output

# Capture the following commands
CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i

# send output here
OUTPUT path:/tmp/crashout/out.tar.gz
```

## Templating
The script also supports templated content to dynamically insert values when the tool runs. The file
uses the Go programming language's style of text template where attributes are wrapped in double curly braces `{{ .<template-attribute> }}`.  Currently, the following
attributes are supported:

* `{{.Home}}` - emits the home directory of the user running the `crash-diagnostics` binary
* `{{.Username}}` - emits the current username that runs the `crash-diagnostics` binary

The following script uses templated values that outputs the current `HOME` directory and
`username`:

```
FROM local 162.164.10.1:2222 162.164.10.2:2222
KUBECONFIG {{.Home}}/.kube/kind-config-kind
AUTHCONFIG username:{{.Username} private-key:{{.Home}}/.ssh/id_rsa
WORKDIR /tmp/output

CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i 

OUTPUT path:{{.Home}}/.crashout/out.tar.gz
```


## Compile and Run
`crash-diagnostics` is written in Go and requires version 1.11 or later.  Clone the source from its repo or download it to your local directory.  From the project's root directory, compile the code with the
following:

```
GO111MODULE="on" go install .
```

This should place the compiled `crash-diagnostics` binary in `$(go env GOPATH)/bin`.  You can test this with:
```
crash-dianostics --help
```
If this does not work properly, ensure that your Go environment is setup properly.

Next run crash-diagonostics using the sample Diagnostics.file in this directory. Ensure to update it to reflect your
current environment:

```
crash-diagnostics run --output crashd.tar.gzip --loglevel debug
```

You should see log messages on the screen similar to the following:
```
DEBU[0000] Parsing script file
DEBU[0000] Parsing [1: FROM local]
DEBU[0000] FROM parsed OK
DEBU[0000] Parsing [2: WORKDIR /tmp/crasdir]
...
DEBU[0000] Archiving [/tmp/crashdir] in out.tar.gz
DEBU[0000] Archived /tmp/crashdir/local/df_-i.txt
DEBU[0000] Archived /tmp/crashdir/local/lsof_-i.txt
DEBU[0000] Archived /tmp/crashdir/local/netstat_-an.txt
DEBU[0000] Archived /tmp/crashdir/local/ps_-ef.txt
DEBU[0000] Archived /tmp/crashdir/local/var/log/syslog
INFO[0000] Created archive out.tar.gz
INFO[0002] Created archive out.tar.gz
INFO[0002] Output done
```

## Contributing

New contributors will need to sign a CLA (contributor license agreement). Details are described in our [contributing](CONTRIBUTING.md) documentation.


## License
This project is available under the [Apache License, Version 2.0](LICENSE.txt)