
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
* [ ] Add project badges to README
* [x] Apply tag v0.1.0-alpha.4

* Translate identified TODOs to issues.

# v.0.1.0-alpha.5
* [ ] Uniform directive parameter format (i.e. DIRECTIVE param:value)
* [ ] Introduce support for quoted strings in directives
* [ ] Revamp EVN variable expansion (i.e. use os.Expand)
* [ ] Apply tag v0.1.0-alpha.5

# v0.1.0-alpha.6
* [ ] New directive `RUN` (i.e. RUN shell:"/bin/bash" cmd:"<command string>" )

# v0.1.0-alpha.7
* [ ] New directive `KUBEGET`

# v0.1.0-Beta.8
* [ ] Documentation update (tutorials and how tos)
* [ ] Recipes (i.e. Diagnostics.file files for different well known debg)
* [ ] Cloud API recipes (i.e. recipes to debug CAPV)

# v0.2.0
* Refactor internal executor into a pluggable interface-driven model
  - i.e. possible suppor for different runtime ()
  - default runtime may use ssh and scp while other runtime may choose to use something else
  - default runtime may use ssh/scp all the time regardless of local or remote 
  