# Roadmap
This project has just started and is going through a steady set of iterative changes to create a tool that will be useful for Kubernetes human operators.  The release cadance is designed to allow the implemented features to mature overtime and lessen technical debts. Each release series will consist of alpha and beta releases before each major release to allow time for the code to be properly exercized by the community.

This roadmap has a short and medium term views of the type of design and functionalities that the tool should support prior to a `1.0` release.

## v0.1.x-Releases
These releases are meant to introduce new features and introduce fundamental designs that will allow the tool to feature-scale. This will change often and may break backward compactivity until a GA version is reached. Starting with `v0.1.0-alpha.0`, the main themes for `v0.1.x` release series are:

* Feedback - continued requirement gathering from early adopters. 
* Documentation - solidify the documentation early for easy usage.
* Standardization - ensure that script directives are consistent and predictable for improved usability.
* Feature growth - continue to improve on current features and add new ones.

## v0.2.x-Releases
The 0.2.x releases will continue to provide stability of features introduced in previous release series (v0.1.0).  
There will be two main themes in this release:
* Introduction of new `KUBEGET` directive to query objects from the API server
* Start a collection of Diagnostics files for troubleshooting recipes
* Redesign the execution backend into a pluggable system allowing different execution runtime (i.e. SSH, HTTP, gRPC, cloud provider, etc)

### Features
The following additional features are also planned for this series.
* Go API - ensure a clear API surface for code reuse and embedding.
* Pluggable executors - make executors (the code that executes the translated script directives) work using pluggable API (i.e. interface) 


## v0.3.x-Releases
This series of release will start introduce optimization features.  Some features may include:
* Parsing and execution optimization (i.e. parallel execution)
* A Uniform retry strategies (smart enough to requeue actions when failed)

## v0.4.x-Releases
Exploring other interesting ideas: 
* Automated diagnostics (would be nice)
* And more...

TBD