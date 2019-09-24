# Crash Recovery and Diagnostics for Kubernetes
A tool to help collect and analyze node information for troubleshooting unresponsive Kubernetes clusters.

This tool is designed for human operators to run against some or all nodes of a troubled or unresponsive cluster to collect logs, configurations, or any other node resource for analysis to help with troubleshotting of Kubernetes.  

## Pre-requisites
 * Assumes the Kubernetes cluster is running on Linux

## `Crash-Diagnostics`
The tool is compiled into a a single binary known named `crash-dianostics`.  It uses a configuration/command file to script how and what resource to be collected from one or more machines. By default, `crash-diagnostics` searches for script `diagnostics.file` which specifies line-by-line directives and commands that are interpreted into actions to be executed against the nodes in the cluster.  

To collect and bundle information resources from cluster machines into a gizipped tar file, use the `--output` flag to specify the tar file name:

```
crash-diagnostics --output small-cluster.tar.gz
```

By default, when `crash-diagnostics` runs, it will automatically look for script file `./Diagnostics.file`.  The script file can be specified, however, using the `--file` flag:

```
crash-diagnostics --file small-cluster.file -output small-cluster.tar.gz
```

## Diagnotics.file Format
`Diagnostics.file` uses a simple line-by-line directive format (à la Dockerfile) to specify how to collect data from cluster servers:

```
DIRECTIVE [arguments]
```

A directive can either be a `preamble` or a `action`.  Preambles provide runtime configuration settings for the tool while actions can execute a command against the clsuter nodes speficied. 

### Example Diagnostics.file
```
FROM 127.0.0.1:22 192.168.99.7:22
WORKDIR /tmp/crashdir

COPY /var/log/system.log
CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i
```

## Diagnostics.file Directives
Currently, `crash-dianostics` supports the following directives:
```
AS
AUTHCONFIG
CAPTURE
COPY
ENV
FROM
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
if available. If theAPI server is not available or unreponsive, retrieval of cluster information is skipped.

```
KUBECONFIG /path/to/kube/config
```


### WORKDIR
Specifies the working directory used when building the archive bundle.  The
directory is  used as temporary location to store data from all data sources
specified in the file.  When the tar is built, the content of that directory
is removed.

```
WORKDIR <relative or absolute path>
```

## Example File

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
```

## Diagnostics.file Templating
The Diagnostics.file supports templated content to dynamically insert values when the tool loads the scipt file. The file
uses Go-style templates where attributes are wrapped in double curly braces `{{ <template-attribute> }}`.  Currently, the following
attributes are supported:

* `{{.Home}}` - emits the home directory of the user running the `crash-diagnostics` binary
* `{{.Username}}` - emits the current username that runs the `crash-diagnostics` binary

The following shows the previous example using templated values:

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
crash-diagnostics --output crashd.tar.gzip --loglevel debug
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