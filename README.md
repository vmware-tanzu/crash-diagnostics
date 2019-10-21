# Crash Recovery and Diagnostics for Kubernetes
A tool to help collect and analyze node information for troubleshooting unresponsive Kubernetes clusters.

This tool is designed for human operators to run against some or all nodes of a troubled or unresponsive cluster to collect logs, configurations, or any other node resource for analysis to help with troubleshooting of Kubernetes.  

## `Crash-Diagnostics`
The tool is compiled into a single binary named `crash-diagnostics`.  Currently, the binary supports two commands:

```
Usage:
  crash-diagnostics [command]

Available Commands:
  help        Help about any command
  run         Executes a diagnostics script file
```

Command `run` uses a diagnostics file to script how and what resources are collected from cluster machines. By default, `crash-diagnostics run` searches for script for `Diagnostics.file` which specifies line-by-line directives and commands that are interpreted into actions to be executed against the nodes in the cluster.

```
> crash-diagnostics run --help

Usage:
  crash-diagnostics run [flags]

Flags:
      --file string     the path to the diagnostics script file to run (default "Diagnostics.file")
      --output string   the path of the generated archive file (default "out.tar.gz")
```


For instance, the following command will execute file `./Diagnostics.file` and store any collected data in file `out.tar.gz`:

```
crash-diagnostics run
```

To run a different script file  or specify a different output archive, use the flags shown below: 

```
crash-diagnostics --file test-cluster.file --output test-cluster.tar.gz
```

## Diagnostics.file Format
`Diagnostics.file` uses a simple line-by-line format (à la Dockerfile) to specify directives on how to collect data from cluster servers:

```
DIRECTIVE [arguments]
```

A directive can either be a `preamble` for runtime configuration or an `action` which can execute a command that runs on each specified host. 

### Example Diagnostics.file
The following is a sample Diagnostics.file that captures command output and copy files from two hosts:
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
commands specified in the file.  The collected information is bundled into archive file `/tmp/crashout/out.tar.gzip` specified by `OUTPUT` (note that 
the output file can also be specified by flag `--output`).

## Diagnostics.file Directives
Currently, `crash-diagnostics` supports the following directives:
```
AS
AUTHCONFIG
CAPTURE
COPY
ENV
FROM
KUBECONFIG
OUTPUT
RUN
WORKDIR
```
Each directive can receive named parameters to pass values to the command it represents.  Each named parameter uses an identifier followed by a colon `:` as shown below:
```
DIRECTIVE name0:<param value 0> name1:<param value 1> ... nameN:<param value N>
```
Optionally, most directives can be declared with a single default unnamed parameter value as shown below:
```
DIRECTIVE <default param value>
```
As an example, directive `WORKDIR` can be declared with its `path` named parameter:
```
WORKDIR path:/some/path
```
Or it can be declared with an unnamed parameter, which internally is assumed to be the `path:` parameter:
```
WORKDIR /some/path
```

### AS
This directive specifies the `userid` and optional `groupid` to  use when `crash-diagnostics` execute its commands against the current machine.
```
AS userid:<userid> [groupid:<groupid>]
```
Example:
```
AS userid:100
```
Or
```
AS userid:vladimir groupid:200
```

### AUTHCONFIG
Configures an authentication for connections to remote node servers.  A `username` must be along with an optional `private-key` which can be used by command backends that support private key/public key certificate such as SSH.

```
AUTHCONFIG username:vladimir private-key:/Users/vladimir/.ssh/ssh_rsa
```

### CAPTURE
This directive captures the output of a command when executed executed on a specified machine (see `FROM` directive).  The output of the executed command is captured and saved in a file that is added to the archive file bundle.

The following shows an example of directive `CAPTURE`:

```
CAPTURE /bin/journalctl -l -u kube-apiserver
```

Or, with its named parameter `cmd:`:
```
CAPTURE cmd:"/bin/journalctl -l -u kube-apiserver"
```

### CAPTURE file names
The captured output will be written to a file whose name is derived from the command string as follows:

```
_bin_journalctl__l__u_kube-api-server.txt
```

### COPY
This directive specifies one or more files (and/or directories) as data sources that are copied
into the arachive bundle as shown in the following example

```
COPY /var/log/kube-proxy.log /var/log/containers

# Or with using its named parameter format with parameter `paths`:

COPY paths:"/var/log/kube-proxy.log /var/log/containers"
```
The previous command will copy file `/var/log/kube-proxy.log` and each file in directory `/var/log/containers` as part of the generated archive bundle.


### ENV
This directive is used to inject environment variables that are made available to other commands at runtime:
```
ENV key0=val0 key1=val1 ... keyN=valN
```
Multiple variables can be declared for each `ENV`.  A Diagnostics file can have one or more `ENV` declarations.

#### ENV Variable Expansion
`Crash-Diagnostics` supports a simple version of Unix-style variable expansion using `$VarName` and `${varName}` formats.  The following example shows how this works:

```
# environment vars
ENV logroot=/var/log kubefile=kube-proxy.log
ENV containerlogs=/var/log/containers

# references vars above
COPY $logroot/${kubefile} 
COPY ${containerlogs}
```
The `ENV` command can optionally used parameter name `vars:` as shown below:
```
ENV vars:"Foo=bar Blat=bat"
```

### FROM
`FROM` specifies a space-separated list of nodes from which data is collected.  Each 
host is specified by an address endpoint consisting of `<host-address>:<port>` as shown in the following example:

```
FROM 10.10.100.2:22 10.10.100.3:22 10.10.100.4:22

# Or using its named parameter `hosts`

FROM hosts:"10.10.100.2:22 10.10.100.3:22 10.10.100.4:22"
```

By default the `crash-diagnostics` will use SSH as a runtime to interact with the specified remote hosts.

### KUBECONFIG
This directive specifies the fully qualified path of the Kubernetes client configuration file. The
tool will use this value to load the Kubernetes configuration file to communicate with the API server to retrieve vital cluster information if available.  

`KUBECONFIG` is declared as shown below:

```
KUBECONFIG /path/to/kube/config

# Or using its named parameter `path`

KUBECONFIG path:"/path/to/kube/config"
```

By default, the following resources will be retrieved from the API server:

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
This directive configures the location and file name of the generated archive file as shown in the following example:

```
OUTPUT /tmp/crashout/out.tar.gz

# Or with its named parameter path

OUTPUT path:"/tmp/crashout/out.tar.gz"
```

If `OUTPUT` is not specified in the `Diagnostics.file`, the tool will apply the value of flag `--output` if provided.

### RUN
This directive executes the specified command on each machine in the `FROM` list. Unlike `CAPTURE` however, the output of the command is not written to the archive file bundle.

The following shows an example of `RUN`:

```
RUN /bin/journalctl -l -u kube-apiserver
 
# Or with its named parameter `cmd`

RUN cmd:"/bin/journalctl -l -u kube-apiserver"
```

`RUN` is useful and helps to execute commands to interact with the remote node for tasks such as data preparation or gathering before aggregation.

The following shows how `RUN` can be used (see [originating issue](https://github.com/vmware-tanzu/crash-diagnostics/issues/4#issuecomment-540926598))

```
# prepare needed data
RUN mkdir -p /tmp/containers
RUN /bin/bash -c 'for file in $(ls /var/log/containers/); do sudo cat /var/log/containers/$file > /tmp/containers/$file; done'
COPY /tmp/containers

# clean up
RUN /usr/bin/rm -rf /tmp/containers
```

### WORKDIR
In a Diagnostics.file, `WORKDIR` specifies the working directory used when building the archive bundle as shown in the following example:

```
WORKDIR /tmp/crashdir

# Or using its named parameter path

WORKDIR path:"/tmp/crashdir"
```

The directory is  used as a temporary location to store data from all data sources specified in the file.  When the tar is built, the content of that directory is removed.

### Example File

```
FROM local 162.164.10.1:2222 162.164.10.2:2222
KUBECONFIG ${USER}/.kube/kind-config-kind
AUTHCONFIG username:test private-key:${USER}/.ssh/id_rsa
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
KUBECONFIG ${USER}/.kube/kind-config-kind
AUTHCONFIG username:test private-key:${USER}/.ssh/id_rsa
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


## Compile and Run
`crash-diagnostics` is written in Go and requires version 1.11 or later.  Clone the source from its repo or download it to your local directory.  From the project's root directory, compile the code with the
following:

```
GO111MODULE="on" go install .
```

This should place the compiled `crash-diagnostics` binary in `$(go env GOPATH)/bin`.  You can test this with:
```
crash-diagnostics --help
```
If this does not work properly, ensure that your Go environment is setup properly.

Next run `crash-diagnostics` using the sample Diagnostics.file in this directory. Ensure to update it to reflect your
current environment:

```
crash-diagnostics run --output crashd.tar.gzip --debug
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