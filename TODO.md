
# v.0.1.0-alpha.0
This tag/version reflects migration to github
* [x] Reset tag/version to alpha.0
* [x] Add ROADMAP.md
* [x] Start CHANGELOG notes
* [x] Apply tag v0.1.0-alpha.0 
* [x] Update release notes in GitHub

# v0.1.0-alpha.4
* [x] Add automated build (Makefile etc)
* [x] Git describe version reporting (i.e. `crash-diagnostics --version`)
* [x] Add project badges to README
* [x] Apply tag v0.1.0-alpha.4

* Translate identified TODOs to issues.

# v.0.1.0-alpha.5
* [x] Uniform directive parameter format (i.e. DIRECTIVE param:value)

# v.0.1.0-alpha.6
* [x] Introduce support for quoted strings in directives

# v0.1.0-alpha.7
* [x] Revamp EVN variable expansion (i.e. use os.Expand)
* [x] Suppport for shell variable expansion format (i.e. ${Var})
* [x] Deprecate support for Go style templating ( i.e. {{.Var}} )
* [x] Apply tag v0.1.0-alpha.7

# v0.1.0-alpha.8
* [x] New directive `RUN` (i.e. RUN shell:"/bin/bash" cmd:"<command string>" )
* [x] Grammatical/word corrections in crash-diagnostics --help output

# v0.1.0-alpha.9
* [x] Embedded quotes bug fix for `RUN` and `CAPTURE`

# v0.1.0 Release
* [x] Doc udpdate
* [x] Todo and Roadmap updates

# v0.1.1
* [x] Fix parameter expansion clash between tool parser and shell
* [x] Introduce ability to escape variable expansion
* [x] Update docs
* [x] Update changelog doc

# v0.2.0-alpha.0
* [x] New directive `KUBEGET`
* [x] Update doc


# v0.2.1-alpha.0
* [x] Introduce support for command echo parameter
* [x] Documentation update

# v0.2.1-alpha.1
* [x] Remove support for local execution model
* [x] The default executor will use SSH/SCP even when targeting local machine
* [x] Update test for new executor backend
* [x] Update CI/CD to automate end-to-end tests using SSH server
* [x] Documentation update

# v0.2.2-alpha.0
* [ ] Initial CloudAPI support

# Other Tasks
* [ ] Documentation update (tutorials and how tos)
* [ ] Recipes (i.e. Diagnostics.file files for different well known debg)
* [ ] Cloud API recipes (i.e. recipes to debug CAPV)

# v0.3.0
* Refactor internal executor into a pluggable interface-driven model
  - i.e. possible suppor for different runtime ()
  - default runtime may use ssh and scp while other runtime may choose to use something else
  - default runtime may use ssh/scp all the time regardless of local or remote 
  