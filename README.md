![](https://github.com/vmware-tanzu/crash-diagnostics/workflows/Crash%20Diagnostics%20Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/vmware-tanzu/crash-diagnostics)](https://goreportcard.com/report/github.com/vmware-tanzu/crash-diagnostics)

# Crashd - Crash Diagnostics

Crash Diagnostics (Crashd) is a tool that helps human operators to easily interact and collect information from infrastructures running on Kubernetes for tasks such as automated diagnosis and troubleshooting.  

## Crashd Features
* Crashd uses the [Starlark language](https://github.com/google/starlark-go/blob/master/doc/spec.md), a Python dialect, to express and invoke automation functions
* Easily automate interaction with infrastructures running Kubernetes
* Interact and capture information from compute resources such as machines (via SSH)
* Automatically execute commands on compute nodes to capture results 
* Capture object and cluster log from the Kubernetes API server
* Easily extract data from Cluster-API managed clusters 


## How Does it Work?
Crashd executes script files, written in Starlark, that interacts a specified infrastructure along with its cluster resources.  Starlark script files contain predefined Starlark functions that are capable of interacting and collect diagnostics and other information from the servers in the cluster.

For detail on the design of Crashd, see this Google Doc design document [here](https://docs.google.com/document/d/1pqYOdTf6ZIT_GSis-AVzlOTm3kyyg-32-seIfULaYEs/edit?usp=sharing).

## Installation
There are two ways to get started with Crashd. Either download a pre-built binary or pull down the code and build it locally.

### Download binary
1. Dowload the latest [binary relase](https://github.com/vmware-tanzu/crash-diagnostics/releases/) for your platform
2. Extract `tarball` from release
   ```
   tar -xvf <RELEASE_TARBALL_NAME>.tar.gz
   ```
3. Move the binary to your operating system's `PATH`


### Compiling from source
Crashd is written in Go and requires version 1.11 or later.  Clone the source from its repo or download it to your local directory.  From the project's root directory, compile the code with the
following:

```
GO111MODULE=on go build -o crashd .
```

Or, yo can run a versioned build using the `build.go` source code:

```
go run .ci/build/build.go

Build amd64/darwin OK: .build/amd64/darwin/crashd
Build amd64/linux OK: .build/amd64/linux/crashd
```

## Getting Started
A Crashd script consists of a collection of Starlark functions stored in a file.  For instance, the following script (saved as diagnostics.crsh) collects system information from a list of provided hosts using SSH.  The collected data is then bundled as tar.gz file at the end: 

```python
# Crashd global config
crshd = crashd_config(workdir="{0}/crashd".format(os.home))

# Enumerate compute resources 
# Define a host list provider with configured SSH
hosts=resources(
    provider=host_list_provider(
        hosts=["170.10.20.30", "170.40.50.60"], 
        ssh_config=ssh_config(
            username=os.username,
            private_key_path="{0}/.ssh/id_rsa".format(os.home),
        ),
    ),
)

# collect data from hosts
capture(cmd="sudo df -i", resources=hosts)
capture(cmd="sudo crictl info", resources=hosts)
capture(cmd="df -h /var/lib/containerd", resources=hosts)
capture(cmd="sudo systemctl status kubelet", resources=hosts)
capture(cmd="sudo systemctl status containerd", resources=hosts)
capture(cmd="sudo journalctl -xeu kubelet", resources=hosts)

# archive collected data
archive(output_file="diagnostics.tar.gz", source_paths=[crshd.workdir])
```

The previous code snippet connects to two hosts (specified in the `host_list_provider`) and execute commands remotely, over SSH, and `capture` and stores the result.

> See the complete list of supported [functions here](./docs/README.md).

### Running the script
To run the script, do the following:

```
$> crashd run diagnostics.crsh 
```

If you want to output debug information, use the `--debug` flag as shown:

```
$> crashd run --debug diagnostics.crsh

DEBU[0000] creating working directory /home/user/crashd
DEBU[0000] run: executing command on 2 resources
DEBU[0000] run: executing command on localhost using ssh: [sudo df -i]
DEBU[0000] ssh.run: /usr/bin/ssh -q -o StrictHostKeyChecking=no -i /home/user/.ssh/id_rsa -p 22  user@localhost "sudo df -i"
DEBU[0001] run: executing command on 170.10.20.30 using ssh: [sudo df -i]
...
```

## Compute Resource Providers
Crashd utilizes the concept of a provider to enumerate compute resources. Each implementation of a provider is responsible for enumerating compute resources on which Crashd can execute commands using a transport (i.e. SSH). Crashd comes with several providers including

* *Host List Provider* - uses an explicit list of host addresses (see previous example)
* *Kubernetes Nodes Provider* - extracts host information from a Kubernetes API node objects
* *CAPV Provider* - uses Cluster-API to discover machines in vSphere cluster
* *CAPA Provider* - uses Cluster-API to discover machines running on AWS
* More providers coming!


## Accessing script parameters
Crashd scripts can access external values that can be used as script parameters.
### Environment variables
  Crashd scripts can access environment variables at runtime using the `os.getenv` method:
```python
kube_capture(what="logs", namespaces=[os.getenv("KUBE_DEFAULT_NS")])
```

### Command-line arguments
Scripts can also access command-line arguments passed as key/value pairs using the `--args` or `--args-file` flags. For instance, when the following command is used to start a script:

```bash
$ crashd run --args="kube_ns=kube-system, username=$(whoami)" diagnostics.crsh
```

Values from `--args` can be accessed as shown below:

```python
kube_capture(what="logs", namespaces=["default", args.kube_ns])
```

## More Examples
### SSH Connection via a jump host
The SSH configuration function can be configured with a jump user and jump host.  This is useful for providers that requires a host proxy for SSH connection as shown in the following example:
```python
ssh=ssh_config(username=os.username, jump_user=args.jump_user, jump_host=args.jump_host)
hosts=host_list_provider(hosts=["some.host", "172.100.100.20"], ssh_config=ssh)
...
```

### Connecting to Kubernetes nodes with SSH
The following uses the `kube_nodes_provider` to connect to Kubernetes nodes and execute remote commands against those nodes using SSH:

```python
# SSH configuration
ssh=ssh_config(
    username=os.username,
    private_key_path="{0}/.ssh/id_rsa".format(os.home),
    port=args.ssh_port,
    max_retries=5,
)

# enumerate nodes as compute resources
nodes=resources(
    provider=kube_nodes_provider(
        kube_config=kube_config(path=args.kubecfg),
        ssh_config=ssh,
    ),
)

# exec `uptime` command on each node
uptimes = run(cmd="uptime", resources=nodes)

# print `run` result from first node
print(uptimes[0].result)
```


### Retreiving Kubernetes API objects and logs
The`kube_capture` is used, in the folliwng example, to connect to a Kubernetes API server to retrieve Kubernetes API objects and logs.  The retrieved data is then saved to the filesystem as shown below:

```python
nspaces=[
    "capi-kubeadm-bootstrap-system",
    "capi-kubeadm-control-plane-system",
    "capi-system capi-webhook-system",
    "cert-manager tkg-system",
]

conf=kube_config(path=args.kubecfg)

# capture Kubernetes API object and store in files
kube_capture(what="logs", namespaces=nspaces, kube_config=conf)
kube_capture(what="objects", kinds=["services", "pods"], namespaces=nspaces, kube_config=conf)
kube_capture(what="objects", kinds=["deployments", "replicasets"], namespaces=nspaces, kube_config=conf)
```

### Interacting with Cluster-API manged machines running on vSphere (CAPV)
As mentioned, Crashd provides the `capv_provider` which allows scripts to interact with Cluster-API managed clusters running on a vSphere infrastructure (CAPV).  The following shows an abbreviated snippet of a Crashd script that retrieves diagnostics information from the mangement cluster machines managed by a CAPV-initiated cluster:

```python
# enumerates management cluster nodes
nodes = resources(
    provider=capv_provider(
        ssh_config=ssh_config(username="capv", private_key_path=args.private_key),
        kube_config=kube_config(path=args.mc_config)
    )
)

# execute and capture commands output from management nodes
capture(cmd="sudo df -i", resources=nodes)
capture(cmd="sudo crictl info", resources=nodes)
capture(cmd="sudo cat /var/log/cloud-init-output.log", resources=nodes)
capture(cmd="sudo cat /var/log/cloud-init.log", resources=nodes)
...

```

The previous snippet interact with management cluster machines. The provider can be configured to enumerate workload machines (by specifying the name of a workload cluster) as shown in the following example:

```python
# enumerates workload cluster nodes
nodes = resources(
    provider=capv_provider(
        workload_cluster=args.cluster_name,
        ssh_config=ssh_config(username="capv", private_key_path=args.private_key),
        kube_config=kube_config(path=args.mc_config)
    )
)

# execute and capture commands output from workload nodes
capture(cmd="sudo df -i", resources=nodes)
capture(cmd="sudo crictl info", resources=nodes)
...
```

### All Examples
See all script examples in the [./examples](./examples) directory.

## Roadmap
This project has numerous possibilities ahead of it.  Read about our evolving [roadmap here](ROADMAP.md).


## Contributing

New contributors will need to sign a CLA (contributor license agreement). Details are described in our [contributing](CONTRIBUTING.md) documentation.


## License
This project is available under the [Apache License, Version 2.0](LICENSE.txt)