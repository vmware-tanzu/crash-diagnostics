# Crash Diagnostics Reference

## Running `crashd`
Crash Diagnostics is compiled into a single binary called `crashd`.  The command can be invoked as follows:

```
Usage:
  crashd [command]

Available Commands:
  help        Help about any command
  run         Executes a script file
```

Command `run` executes the specified sript file. Use flag `--help` to get additional help for a given command:

```
> crashd run --help

Usage:
  crashd run [flags] script-file
  ...
```

### Passing script arguments
`crashd` script files can receive parameters from the command-line using the `--args` flag which takes a key/value pair seprated by spaces as shown below:

```
crashd run --args "arg0='value 0' args1='value 1'"
```
These values can be accessed inside a running script using the `args` struct as follows:

```python
ssh_config(username=args.args0, private_key_path=args.args1)
```

### Accessing environment variables
At runtime, `crashd` scripts can also access values stored in environment variables as shown in the following snippet:

```
KUBE_NS=capi-system crashd run file.crsh
```

The running script can access `KUBE_NS` using the `os` struct as shown below:

```python
kube_capture(what="logs", namespaces=[os.getenv("KUBE_NS")])
```

## Starlark: the Crashd Language
Crashd scripts are written in Starlark, a python dialect.  This means that Crashd scripts can have normal programming constructs:
- Variable declarations
- Function definitions
- Simple data types (string, numeric, bool)
- Composite types (dictionary, list, tuple, set, and functions)
- Statements and expressions
- Etc

> For more on Starlark, see the [language reference](https://github.com/bazelbuild/starlark/blob/master/spec.md).


## The Crashd Script File
A script file is composed Starlark language elements and built-in functions provided by Crashd at runtime. In addition to built-in functions, script authors have the ability to define their own custom functions that can be reused in the script.  The following is an example of a valid script that `crashd` can execute:

```python
def from_hosts():
    hosts = run_local("cat /etc/hosts | grep -E '([0-9]){3}\.' | awk '{print $1}'")
    return hosts.splitlines()

ssh_config(username="username", port=2222, max_retries=10)
resources(hosts=from_hosts())

capture(cmd="sudo crictl info")
copy(path="/var/log/cloud-init-output.log")
copy(path="/var/log/cloud-init.log")
```

The previous example shows the definition of a custom function `from_host` which extracts a list of hosts from the local host file. The script also show the use of several built-in functions including:
* `ssh_config`
* `resources`
* `capture`
* `copy`

These built-in functions are used to configure the script and issue commands against remote compute resources.

## Crashd Built-in Types
Crashd comes with many built-in functions and other types to help you create functioning and useful scripts. Each built-in function falls in to one the following category:
* Configuration functions
* Provider functions
* Resource enumeration function
* Command functions
* Default Values
* OS data and functions
* Argument data

## Configuration Functions
Configuration functions help to declare data structures that are used to store configuration information that can be used in the script.

### `crashd_config()`
This function declares script-wide configuration information that is used to configure the script behavior at runtime.  Values declared here are usually not used directly by the script.  

#### Parameters

| Param | Description | Required |
|:--------|:--------|:--------|
| `workdir` | the working directory used by some functions to store files.| Yes |
| `uid` | User ID used to run local commands|No, defaults to current ID|
| `gid` | Group ID used to run local commands|No, defaults to current ID|
| `default_shell` |The default shell to use to execute commands |No, defaults to no shell|
| `use_ssh_agent` | boolean indicator to start a ssh-agent instance or not |No, defaults to `False`|


#### Output
`crashd_config()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `workdir` | The provided `workdir` |
| `uid` | The current UID set |
| `gid` | The current GID set |
| `default_shell`|The shell set, if any|

#### Example
```python
crashd_config(
    workdir = "{}/crashd".format(os.home)
)
```

#### Internal Crashd ssh-agent
If the crashd operator does not want to rely on the default ssh-agent process, the `crashd_config()` provides the option to start a new instance of the ssh-agent which will be used for all corresponding ssh/scp connections in the script.
```python
# this will force crashd to use a new ssh-agent instance alive for
# the scope of script execution
crashd_config(workdir="/tmp/foo", use_ssh_agent=True)
```

While leveraging the internal crashd agent, any **passphrase protected keys** will pause the script execution and prompt the script operator to enter the passphrase.


### `kube_config()`
This configuration function declares and stores configuration needed to connect to a Kubernetes API server.

#### Parameters
| Param | Description | Required |
| -------- | -------- | ------- |
| `path`  | Path to the local Kubernetes config file. Default: `$HOME/.kube/config`| No |
| `cluster_context`  | The name of a context to use when accessing the cluster. Default: (empty) | No |
| `capi_provider` | A Cluster-API provider (see providers below) to obtain Kubernetes configurations | No |

#### Output
`kube_config()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `path` | The path to the local Kubernetes config that was set |
| `cluster_context` | The name of a context that was set for the cluster |
| `capi_provider`|A provider that was set for Cluster-API usage|

#### Example
```python
kube_config(path=args.kube_conf, cluster_context="my-cluster")
```
### `ssh_config()`
This function creates configuration that can be used to connect via SSH to remote machines.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `username`  | SSH user ID| Yes |
| `private_key_path`| Path for private key | No, default: `$HOME/.ssh/id_rsa` |
| `port` | Port for SSH connection | No,  default `"22"` |
| `jump_user` | Username for an SSH proxy connection | No |
| `jump_host` | Host address for an SSH proxy connection | Yes if `jump_user` is provided |
| `max_retries` | The maximum number of tries to connect to SSH host| No default `5`|

#### Output
`ssh_config()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `username` | The `username` that was set |
| `private_key_path` | The private file that was set |
| `port` | The port value that was set |
| `jump_user`|The proxy user that was set|
| `jump_host`|The proxy host that was set if proxy user was provided|
| `max_retries`|The max number of retries set|

#### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port=args.ssh_port,
    max_retries=5,
)
```

**NOTE**: For passphrase protected keys, add the key to the default ssh-agent prior to running the diagnostics script, to ensure non-interactive execution of the script.

## Provider Functions
A provider function implements the code to cofigure and to enumerate compute resources for a given infrastructure. The result of the provider functions are used by the `resources` function to generate/enumerate the compute resources needed.

### `capa_provider()`
This function configures the Cluster-API provider for AWS (CAPA).  This provider can enumerate management or workload cluster machines in order to execute commands using SSH on those machines. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `ssh_config`|SSH configuration returned by `ssh_config()`|Yes |
| `mgmt_kube_config` |Kubernetes configuration returned by `kube_config`|Yes|
| `workload_cluster`|The name of a workload cluster. When specified the provider will retrieve a cluster's compute nodes for the workload cluster.|No|
| `namespace`|The namespace in which the workload cluster was created, if `workload_cluster` is specified. If no `workload_cluster` is specified, then this should be the namespace of the management cluster.|No|
| `labels`|A list of labels used to filter cluster's compute nodes|No|
| `nodes` |A list of node names that can filter selected cluster nodes|No|

#### Output
`capa_provider()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `kind`| The name of the provider (`capv_provider`)|
|`transport`|The name of the transport to use (i.e. `ssh, http, etc`)|
| `ssh_config` | A struct with SSH configuration |
| `kube_config` | A struct with Kubernetes configuration |
| `workload_cluster` | The name of the  |
| `hosts`|A list of host addresses generated from cluster information|

#### Example
```python

ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port=args.ssh_port,
    max_retries=5,
)

kube=kube_config(path=args.kube_conf)

capa_provider(
    workload_cluster="my-wc-cluster",
    namespace="workloads"
    ssh_config=ssh,
    kube_config=kube
)
```

### `capv_provider()`
This function configures a provider for a Cluster-API managed cluster running on vSphere (CAPV).  By default, this provider will enumerate cluster resources for the management cluster.  However, by specifying the name of a `workload_cluster`, the provider will enumarate cluster compute resources for the workload cluster. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `ssh_config`|SSH configuration returned by `ssh_config()`|Yes |
| `mgmt_kube_config` |Kubernetes configuration returned by `kube_config`|Yes|
| `workload_cluster`|The name of a workload cluster. When specified the provider will retrieve a cluster's compute nodes for the workload cluster.|No|
| `namespace`|The namespace in which the workload cluster was created, if `workload_cluster` is specified. If no `workload_cluster` is specified, then this should be the namespace of the management cluster.|No|
| `labels`|A list of labels used to filter cluster's compute nodes|No|
| `nodes` |A list of node names that can filter selected cluster nodes|No|

#### Output
`capv_provider()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `kind`| The name of the provider (`capv_provider`)|
|`transport`|The name of the transport to use (i.e. `ssh, http, etc`)|
| `ssh_config` | A struct with SSH configuration |
| `kube_config` | A struct with Kubernetes configuration |
| `workload_cluster` | The name of the  |
| `hosts`|A list of host addresses generated from cluster information|

#### Example
```python

ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port=args.ssh_port,
    max_retries=5,
)

kube=kube_config(path=args.kube_conf)

capv_provider(
    workload_cluster="my-wc-cluster",
    ssh_config=ssh,
    kube_config=kube
)
```

### `host_list_provider()`
As its name suggests, this provider is used to explicitly specify a list of host addresses directly. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `hosts` | A list of IP addresses or machine names | Yes |
| `ssh_config` | An SSH configuration as returned by ssh_config() | Yes |

#### Output
`host_list_provider()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `hosts` | The list of hosts that was set |
| `ssh_config` | The SSH configuration that was set|

#### Output
`capv_provider()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `kind`| The name of the provider (`host_list_provider`)|
| `transport`|The name of the transport to use (i.e. `ssh, http, etc`)|
| `ssh_config` | A struct with SSH configuration |
| `hosts`|The list of host addresses|

#### Example

```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

host_list_provider(hosts=["172.100.10.20", "ctlplane.local"], ssh_config=ssh)
```

### `kube_nodes_provider()`
This provider captures configuration information to enumerate a Kubernetes cluster nodes. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `kube_config` | Kubernetes config returned by `kube_config()` | Yes |
| `ssh_config` | An SSH configuration as returned by ssh_config() | Yes |
| `names`|A list of names used to filter nodes |No|
| `labels`|A list of labels used to filter nodes|No|

#### Output
`kube_nodes_provider()` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
| `kind`| The name of the provider (`kube_nodes_provider`)|
| `transport`|The name of the transport to use (i.e. `ssh, http, etc`)|
| `ssh_config` | A struct with SSH configuration |
| `kube_config` | The Kubernetes configuration that was set |
| `hosts`|A list of host addresses generated from cluster information|

#### Example

```python
ssh=ssh_config(
    username=args.username,
    private_key_path=args.key_path,
    port=args.ssh_port,
    max_retries=5,
)

kube_nodes_provider(
    kube_config=kube_config(path=args.kubecfg),
    ssh_config=ssh,
)
```

## Resource Enumeration
Crashd uses the notion of a compute resource to which the running script can connect and possibly execute commands (see Command Functions). 

### `resrouces()`
The Crashd script uses the `resources` function along with a provider (see providers above) to properly enumerate compute resources. Each provider implements its own logic which determines how resources are enumerated.

#### Parameter
| Param | Description | Required |
| -------- | -------- | -------- |
|`provider`|Species the provider to use for resource enumeration|Yes|

#### Output
`resources` returns a list of structs based on the type of provider that is used.

For `host_list_provider`, `kube_nodes_provider`, and `capv_provider`, each struct has the following fields.

| Field | Description |
| --------| --------- |
| `kind` | The kind for the resources (`host_resource`) |
| `provider` | The name of the provider that generated the resource |
| `host` | Host address |
| `transport`|transport to use|
| `ssh_config`|SSH configuration|

#### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

hosts=resources(
    provider=host_list_provider(
        hosts=["localhost", "127.0.0.1"], 
        ssh_config=ssh,
    ),
)

run(cmd="uptime", resources=hosts)
```
In the previous example, `hosts` contains the a list of informatation about hosts that can be used in command functions such as `run`.

## Command Functions
Command functions can execute commands on all specified enumerated compute resources automatically or be used in a custom function (`def`) for more control.

### `archive()`
The archive function bundles the specified directories into a single archive file (format tar.gz).

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
|`source_paths`|A list of directories to be archived|Yes|
|`output_file`|The name of the generated archive file|No, default `archive.tar.gz`|

#### Output
`archive` returns the full path of the created bundled file.


### `capture()`
This function runs its command all provided compute resources automatically. The output of the executed command is captured and saved in a file for each execution.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `cmd`|The command string to execute|Yes|
| `resources`|The value returned by `resources()`|Yes|
| `workdir`|A parent directory where captured files will be saved|No, defaults to `crashd_config.workdir`|
| `file_name`|The path/name of the generated file|No, auto-generated based on command string, if omitted|
| `desc`|A short description added at the start of the file|No|

#### Output
`capture()` returns a list `[]` of command result struct for each compute resource where the command was executed. Each struct contains the following fields.

| Field | Description |
| --------| --------- |
| `resource` | The address or name of the compute resource |
| `result` | the path of the file created |
| `err` | An error message if one was encountered |

#### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

hosts=resources(
    provider=host_list_provider(
        hosts=["localhost", "127.0.0.1"], 
        ssh_config=ssh,
    ),
)

capture(cmd="sudo df -i", resources=hosts)
capture(cmd="sudo crictl info", resources=hosts)
capture(cmd="df -h /var/lib/containerd", resources=hosts)
capture(cmd="sudo systemctl status kubelet", resources=hosts)

```

### `capture_local()`
This function runs a command locally on the machine running the script.  It then captures its output in a specified file. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `cmd`|The command string to execute|Yes|
| `workdir`|A parent directory where captured files will be saved|No, defaults to `crashd_config.workdir`|
| `file_name`|The path/name of the generated file|No, auto-generated based on command string, if omitted|
| `desc`|A short description added at the start of the file|No|
| `append` | boolean indicator to append to a file if it already exists or not |No, defaults to `False`|

#### Output
`capture_local()` returns the full path of the capured output file.


### `copy_from()`
This command specifies a list of files that are copied from a remote location to the local machine running the script.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `path`|The path of the remote file|Yes|
| `resources`|The value returned by `resources()`|Yes|
| `workdir`|A parent directory where files are copied to|No, defaults to `crashd_config.workdir`|

#### Output
`copy()` returns a list `[]` of command result struct for each compute resource where the command was executed. Each struct contains the following fields.

| Field | Description |
| --------| --------- |
| `resource` | The address or name of the compute resource |
| `result` | the path of the file copied |
| `err` | An error message if one was encountered |

#### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

hosts=resources(
    provider=host_list_provider(
        hosts=["localhost", "127.0.0.1"], 
        ssh_config=ssh,
    ),
)

copy_from(path="/var/log/kube*.log", resources=hosts)
```

### `log()`
This function prints a log message on the terminal.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `msg`|The message to print on the screen.|Yes|
| `prefix`|An optional prefix that is printed prior to the message.|No|

#### Output
None

#### Example
```python
log(msg="Hello World!")
log(msg="Failed to reach server", prefix="ERROR")
```

### `run()`
This function executes its specified command string on all provided compute resources automatically.  It then returns a list of result objects containing information about the remote compute resource, where the command was executed, and the result of the command. 

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `cmd`|The command string to execute on each compute resource|Yes|
| `resources`|A collection of compute resources returned by `resources()`|Yes|

#### Output
`run()` returns a list `[]` of command result structs for each compute resource where the command was executed. 
Each struct contains the following fields.

| Field | Description |
| --------| --------- |
| `resource` | The address or name of the compute resource where the command was executed |
| `result` | The result of the command on the resource |
| `err` | An error message if one was encountered |

#### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

hosts=resources(
    provider=host_list_provider(
        hosts=["ctrlplane.local", "172.10.20.30"], 
        ssh_config=ssh,
    ),
)

# run uptime command on all hosts
uptimes = run(cmd="uptime", resources=hosts)

#print result for each host
print(uptimes[0].result)
print(uptimes[1].result)
```
### `run_local()`
This function executes a command locally on the machine running the script and returns the result as a string.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
| `cmd`|The command string to execute|Yes|

#### Output
`run_local` returns the result of the command as a string value.

#### Example

```python
# run_local to parse local /etc/hosts file
def from_hosts():
    hosts = run_local("""cat /etc/hosts | grep -E '([0-9]){3}\\.' | awk '{print $1}'""")
    return hosts.splitlines()

ssh_config(username=os.user, port=2222, retries=10)
hosts=resources(provider=host_list_provivider(hosts=from_hosts()))

# run on hosts
uptimes = run(cmd="uptime", resources=hosts)
print(uptimes[0].result)
print(uptimes[1].result)
```
## Kubernetes Functions
These are functions used to execute API requests against a running Kubernetes cluster using a Kubernetes configuration (either explicitly defined or from predeclared default). 

### `kube_capture()`
The `kube_capture` function retrieves Kubernetes API objects and container logs.  The captured information is stored in local files with directory structure similar to that of `kubectl cluster-info dump`.

#### Parameters
| Param | Description | Required |
| -------- | -------- | -------- |
|`what`|Specifies what to get inclusing `objects` or `logs`|Yes|
|`groups`|A list of API groups from which to retrieve API objects.  The core group is named `core`|No|
|`kinds`|A list of object kinds to select|No|
|`namespaces`|A list of namespaces from which to select objects|No|
|`versions`|A list of API versions used to select objects|No|
|`names`|A list used to filter retrieved object by names|No|
|`labels`|A list of label selector expressions used to filter objects|No|
|`containers`|A list of container names used to filter when selecting pod objects|No|
|`kube_config`|The Kubernetes configuration used for this call|No, uses default if omitted|

#### Output
Function `kube_capture` returns a struct with the following fields.

| Field | Description |
| --------| --------- |
|`file`|The root directory where the captured files are saved|
|`error`|An error message, if any was encountered|

#### Example
```python

kube = kube_config(path=args.kube_cfg)

pod_ns=["default", "kube-system"]

kube_capture(what="logs", namespaces=pod_ns, kube_config=kube)
kube_capture(what="objects", kinds=["pods", "services"], namespaces=pod_ns, kube_config=kube)
kube_capture(what="objects", kinds=["deployments", "replicasets"], groups=["apps"], namespaces=pod_ns, kube_config=kube)
```

## Default Values
Some value types can be saved as default values during the execution of a
script.  When the following values are saved as default, Crashd will automatically use
the last known default value for that type when appropriate:
* `kube_config` - the struct created by calling `kube_config`
* `ssh_config` - the struct created by calling `ssh_config()`
* `resources` - the list of struct created by calling `resources()`

### Setting Default Values
Default values are set using the `set_defaults()` function.  Each time this function 
is called, it will save the last instance of a given type (overwriting the previous)
value. 

For instance, consider the following script:
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

hosts=resources(
    provider=host_list_provider(
        hosts=["localhost", "127.0.0.1"], 
        ssh_config=ssh,
    ),
)

capture(cmd="sudo df -i", resources=hosts)
capture(cmd="sudo crictl info", resources=hosts)
capture(cmd="df -h /var/lib/containerd", resources=hosts)
capture(cmd="sudo systemctl status kubelet", resources=hosts)
```

The previous script can be simplified using default values:
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5,
)

set_defaults(hosts=resources(
    provider=host_list_provider(
        hosts=["localhost", "127.0.0.1"], 
        ssh_config=ssh,
    ),
))

capture(cmd="sudo df -i")
capture(cmd="sudo crictl info")
capture(cmd="df -h /var/lib/containerd")
capture(cmd="sudo systemctl status kubelet")
```
The previous can be further simplified by setting the `ssh_config` as a default
value as follows:

```python
set_defaults(ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port="2222",
    max_retries=5
))

# host_list_provider shortcut for resources
set_defaults(resources(hosts=["localhost", "127.0.0.1"]))

capture(cmd="sudo df -i")
capture(cmd="sudo crictl info")
capture(cmd="df -h /var/lib/containerd")
capture(cmd="sudo systemctl status kubelet")
```
The previous snippet sets the values for both the `ssh_config` and `resources`
as default. Notice also that `resources()` supports a shortcut to specify
host lists directly as a parameter.  Internally, `resources()` creates an
instance of the `host_list_provider` when this shortcut is used.

## OS Struct
At runtime, executing scripts are able to access OS information via a global OS struct.

| Field | Description |
| ------- | ---------- | 
|`os.name`| Returns the name of the OS running the script |
|`os.username`|The current username running the script|
|`os.home`|The home directory associated with the user running the script|
| `os.getenv()` | A function which returns the value of the provided environment variable name|

### Example
```python
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    max_retries=5,
)
```

## Argument Struct
A running script can receive argument values from the command that invoked
the script using the `--args` flag which takes a space-separated key/value pair
as shown:

```
crashd run --args "ssh_user='capv' ssh_port='2121' kube_cfg='~/my/cfg' file.crsh
```
In the script, the args can be used as follows:
```python
ssh=ssh_config(
    username=args.ssh_user,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port=args.ssh_port
    max_retries=5,
)

kube_config(path=args.kube_cfg)
```

### Arguments file
In the case, when the script requires mutliple values to be provided by the user, the `--args` flag becomes difficult to use. The `run` command exposes the `--args-file` flag which takes a file path as input.

The supplied args file should follow the format:
* A line contains a single key-value pair separated by `=` sign (eg: foo=bar|foo =bar|foo= bar|foo = bar)
* A line can either contain a key-value pair in the above format or a comment statement starting with #
* Blank lines are allowed

Any line not adhering to the said format will result in a warning message to appear on the screen, and would be ignored.

```bash
$ cat /tmp/script.args
foo=bar
bloop blah

# this will result in a warning message with foo=bar as the only pair pairs to be passed to the .crsh file
$ crash run diagnsotics.crsh --args-file /tmp/script.args
WARN[0000] unknown entry in args file: blooop blah
```
