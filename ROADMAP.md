# Roadmap
This project has just started and is going through a steady set of iterative changes to create a tool that will be useful for Kubernetes human operators.  The release cadance is designed to allow the implemented features to mature overtime and lessen technical debts. Each release series will consist of alpha and beta releases before each major release to allow time for the code to be properly exercized by the community.

This roadmap has a short and medium term views of the type of design and functionalities that the tool should support prior to a `1.0` release.

## v0.1.x-Releases
These releases are meant to introduce new features and introduce fundamental designs that will allow the tool to feature-scale. This will change often and may break backward compactivity until a GA version is reached. Starting with `v0.1.0-alpha.0`, the main themes for `v0.1.x` release series are:

* Feedback - continued requirement gathering from early adopters. 
* Documentation - solidify the documentation early for easy usage.
* Standardization - ensure that script directives are consistent and predictable for improved usability.
* Feature growth - continue to improve on current features and add new ones.

### Features
* Go API - ensure a clear API surface for code reuse and embedding.
* Pluggable executors - make executors (the code that executes the translated script directives) work using pluggable API (i.e. interface) 
* `KUBEGET` - new directive to collect objects from a running Kubernetes API server
* Develop troubleshooting recipes

## v0.2.x-Releases
These releases will provide stability of features introduced in previous release series (v0.1.0).  The main theme in this release series is optimization. `Crash-diagnostics` has plenty of opportunities to improve operations when collecting cluster information.  

### Features
Some features may include:
* Parsing and execution optimization (i.e. parallel execution)
* A Uniform retry strategies (smart enough to requeue actions when failed)
* Automated diagnostics (would be nice)
* And more...

## v0.3.x-Releases
TBD