# Changelog

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