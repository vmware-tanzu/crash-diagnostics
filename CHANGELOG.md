# Changelog
## v0.2.0-alpha.0
This release introduces new directive `KUBEGET` to retrieve API objects and pod logs from the API server.  When an API connection is configured properly using `KUBECONFIG`, KUBEGET can be used to retrieve any accessible arbitrary API objects or access logs for running pods as shown in the following example:

```
KUBEGET objects groups:"core" kinds:"pods" namespaces:"kube-system default" containers:"kindnet-cni etcd"
```

The previous command would retrieve all pods from namespace kube-system or default having containers named kindnet-cni or etcd.

See [README](https://github.com/vmware-tanzu/crash-diagnostics/tree/master/docs#kubeget) for detail.

## v0.1.1
This release fixes variable expansion clash whereby Crash Diagnostics would interpret variables early that are intended to be sent to remote servers.  For instance, the following `RUN` would not work properly since `$f` would be interpreted by Crash Diagnostics as empty:

```
RUN /bin/bash -c 'for f in $(find /var/logs/containers -type f); do cat $f; done'
```
This PR introduces the ability to escape the variable expansion, allowing it to be sent to the server as intended as shown below:

```
RUN /bin/bash -c 'for f in \$(find /var/logs/containers -type f); do cat \$f; done'
``` 

See [README](https://github.com/vmware-tanzu/crash-diagnostics/blob/master/docs/README.md#escaping-variable-expansion) for detail.

## Changelog
6fca28a Merge pull request #31 from vladimirvivien/variable-expansion-clash-fix
79783de Documentation update for variable expansion escape
b6ffd97 Fix for variable expansion clash with expansion escape

## v0.1.0
This is the first release of the project.  It marks the end of the 0.1.0-alpha.x release series designed to get the project to a stable and usable place.  It includes the following high level features:

* Support for diagnostics script file
* Support for several directive commands including RUN, COPY, CAPTURE, etc
* Flexible script directive format with support fro nested quotes for complex shell commands
* Script supports environment variable declaration with variable expansion
* Ability to execute diagnostics script commands on remote machines
* Automatically collect and collate script result as a tar file
* Etc

See the [README](https://github.com/vmware-tanzu/crash-diagnostics/blob/master/README.md) for feature detail.

#### Changelog

79322fa Merge pull request #28 from vladimirvivien/v0.1.0-release
8a0d801 v0.1.0 Release and doc update
32e599b Merge pull request #27 from vladimirvivien/release-v0.1.0
7124db5 v0.1.0 release doc updates

## v0.1.0-alpha.9
This is a fix release that corrects the way directives such as `RUN` and `CAPTURE` handle embedded quotes (see issue #23 for detail).  This release also adds the `shell:` named param to `RUN` and `CAPTURE` to specify a shell program to use to wrap a given command.  So now the followings are supported for these two actions (note the way quotes can now be embedded):

* `{RUN | CAPTURE} echo "Hello World!"`
* `{RUN | CAPTURE} "echo 'Hello World!'"`
* `{RUN | CAPTURE} 'echo "Hello World"'`
* `{RUN | CAPTURE} "/bin/sh -c 'date -u'"`
* `{RUN | CAPTURE} /bin/bash -c 'echo "Hello World"'`
* `{RUN | CAPTURE} cmd:"echo 'HELLO WORLD'"`
* `{RUN | CAPTURE} shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'"`

#### Changelog
264a134 Merge pull request #25 from vladimirvivien/quote-bug-fix
88071f2 Changes and tests for quoted SSH commands
bf355da Executor test updates to handle embedded quotes
e9b93ec Script parser updates to handle embeded quotes
dac401b Merge pull request #22 from YanzhaoLi/topic/fix-run-error-info
fd7e470 Change Run command's error info


## v0.1.0-alpha.8
This release introduces the new `RUN` directive to execute commands on remote machines without capturing the output in the generated archive file like `CAPTURE`.  This is useful help interact with the remote nodes for tasks such as data preparation such as the following example:

```
# gather logs
RUN mkdir -p /tmp/all-logs
RUN cp /var/log/containers/* /tmp/all-logs 
RUN cp /var/log/kube-apiserver.log /tmp/all-logs
RUN cp /var/log/kubelet.log /tmp/all-logs
RUN cp /var/log/kube-proxy.log /tmp/all-logs

# copy all logs
COPY /tmp/all-logs

# clean up
RUN /usr/bin/rm -rf /tmp/all-logs
```

#### Changelog
e2a9055 Merge pull request #19 from vladimirvivien/run-command-support
1cd6097 Documentation updates for RUN command
f8442b3 Executor updates for RUN command
64d62a5 Parser updates for RUN command


## v0.1.0-alhpa.7
This release introduces Unix-style variable expansion in Diagnostics.file script.  The previous `Go-style` templating has been `deprecated`  in favor of the more traditional variable expansion.

## Variable Expansion
Now all directives that appear in the Diagnostic script file can support variable expansion. This provides a more familiar feel to those with experience writing shell script.  For instance, the following shows how variable expansion can be used in a Diagnostics script:

```
ENV clogs=/var/log/containers
ENV klogroot=/var/log
COPY ${clogs}
COPY $klogroot/kubelet.log
COPY ${klogroot}/kube-proxy.log
COPY ${klogroot}/kube-apiserver.log
```

Diagnostics script files can also access any environment variables made available to the process at runtime:
```
FROM 127.0.0.1:22
AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
CAPTURE /bin/echo "HELLO WORLD"
```

#### Changelog
e64ccca Merge pull request #18 from vladimirvivien/var-expansion-redo
43b2371 Documentation updates
0621e2c Executor updates with var expansion support
bfd2430 Parser update  with variable expansion support


## v0.1.0-alpha.6
This release introduces better more robust directive string parsing with support for quoted values.  With this release, the followings are possible:
* `DIRECTIVE paramName:ParamValue` example: `WORKDIR path:/tmp/mypath`
* `DIRECTIVE paramName:"paramValue"` example `OUTPUT path:"/tmp/my dir"`
* `DIRECTIVE paramName:'value0 "value1"'`  example `CAPTURE cmd:'/bin/echo "Hello World!"'`

Also most directive also support now a default parameter than can be specified without a parameter name which can also support quoted values:

* `WORKDIR "/tmp/my dir"`
* `ENV "key0=val0 key1=val1"`
* Etc. (see documentation)

#### Changelog
647c71b Merge pull request #14 from vladimirvivien/quoted-cmd-parsing
4c22fef Better command parsing, quoted value support


## v0.1.0-alpha.5
This release introduces a uniform way of specifying parameters for DIRECTIVEs in the Diagnostics.file.  Now, named parameters may be used to pass values into all directives.  See the README.

#### Changelog
5710199 Merge pull request #13 from vladimirvivien/uniform-directive-params
3cfaf0c Update README doc
ba3e873 exec package supports named params
370d43a Parser updated to support named params
22ac1fc Start refactoring for all commands


## v0.1.0-alpha.4
This release sets up automated build and release with goreleaser and travis
No new features introduced.

#### Changelog
4b6c1b7 Merge pull request #8 from vladimirvivien/fix-travis-encrypt-val-parsing
ba8fe45 Fix for encrypted env in travis
b7eab26 Merge pull request #7 from vladimirvivien/travis-fix
a239215 (origin/travis-fix) Fixing encrypted global var
8ca1f57 Merge pull request #6 from vladimirvivien/release-automation-fix
f76f5ff (origin/release-automation-fix) Fix to travis encrypted values
81b4c23 Merge pull request #3 from vladimirvivien/release-automation
8aa0fd0 (origin/release-automation) Build and release automation with travis, gorelaser
0f6b3f8 Merge pull request #2 from vladimirvivien/automated-build
9a0ef20 (origin/automated-build) Initial CI setup for automatic build/test


## v0.1.0-alpha.0 (10/2/2019)
Initial release.

The initial set of features of the tool allow cluster operators to investigate unresponsive or dead clusters. The code has been in development for a few weeks prior to this initial release and includes the following features:

* Gather specified resources (log files, etc) from one or more cluster machines
* Collected resources are automatically bundled into an archive (tar) file for further analysis
* Use of user-defined declarative script (i.e. `Diagnostics.file`) to declare resources to gather
* Script include directives:
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
* Support for parameterized script using Go-style templated values
* See [README.md](README.md) for detail

## Previous changes

#### Changelog
e805951 (HEAD -> alpha.0-migration-doc, upstream/master, upstream/HEAD, update-author-info, master) Migration to GitHub/VMware-Tanzu
bd4846e OSS compliance docs and updates
2f70416 Naming refactor, fixes, and tests
c4a9f59 Refactor SSHCONFIG to AUTHCONFIG, add OUTPUT command for bundle config, fixes
e6965a8 Rebranding, doc update, code update
c586945 Remove textual reference to old project name
ff664b4 Adds source and package documentation
44a4167 Adds source and package documentation
8e5c98f Adding copyright and license to source
b7ff4d9 Prevent program on command failure, update script file
2df345d New tests for all remote functions, doc upste
2edd2e1 Executes CAPTURE command remotely with tests
db5bbdb Code refactor for remote exec support, test for remote exec
f0d68a9 Preparing for remote cli exec with ssh
342bcdf Abtraction of command execution
326c0d3 Merge branch 'rewrite-copy-command' into 'master'
f4e0117 Refactor COPY command to use cli command cp
1188aec Merge branch 'impl-templating' into 'master'
e07a95e Adds support for Go style templating in script file
56d4e5e Merge changes from branch impl-kubecfg-support
69c04cf Supports KUBECONFIG directive
c278548 Merge branch 'impl-as-command'
87bdfa6 Merge branch 'impl-as-command' into 'master'
d9e0f54 Internal refactor to implement all commands
5e76f64 Refactor all commands and added new tests
f2d675c Adding more exec tests for each command
bb1fe24 Ads test for parser
babd246 Refactor commands and additional tests
ded9ebf Refactor parser/executor to type-centric
20ab146 Added ENV instruction including tests
7a05184 ENV impl WIP
03d509a Tests for AS instruction
2cd6c88 Starting impl of AS instruction
6a47cf1 Update README.md with more instructions
62294f5 Added archiver to tar file
2ae42bc Able to execute a flare file
4eade80 Adding support for flare.file
7bbcd54 creating the tar utility
64d1608 More project name updates.
a39c067 Update with project name flare
e6aec70 Added preliminary README detail.
9226859 Initial commit