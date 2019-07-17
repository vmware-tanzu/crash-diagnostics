# Flare
A tool that collects server data for troubleshooting.

A human operator would run flare directly on a running server to take a snapshot of the current state of the server for analysis.  Flare can collect data from different data sources specified at runtime.  The collected data is packaged and compressed in a tar bundle.


## Pre-requisites
 * Assumes Linux/Unix


## Flare.file
Flare uses a file (usually named `flare.file`) that contains a simple dialect to declartively specify where and how data is collected from the running server.  The following example uses default `./flare.file` if one exist on the current directory. If one is not found, flare uses sensible default to gather machine state data from known data sources:

 ```
 flare -o flare.tar.gz
 ```

A flare file can be specified using `-f`:

```
flare -f flare.file.name -o out.tar.gz
```

### Flare.file Format
Flare uses a simple format for the command file
(similar to Dockerfile) as shown below:

```
COMMAND [arguments]
```

Currently, flare supports the following commands:
```
AS
CAPTURE
COPY
ENV
FROM
WORKDIR
```

### AS
Specifies the group id and the user id to  use when running flare commands.
```
AS <userid>[:<groupid>]
```

### CAPTURE
Executes a shell command and captures the output as a data source that is copied
into the built archive bundle.

```
CAPTURE [<shell command>]
```

### COPY
Specifies one or more files (or directories) as data sources that are copied
into the arhive bundle.

```
COPY <file or directory list>
```

### ENV
Can be used to set up environment variables that are then exposed to commands
executed by the CAPTURE command.

```
ENV key=value
```

### FROM
Specifies the machine to use as the source of the data collected.  Currently
only value `local` is supported.

```
FROM local
```

### WORKDIR
Specifies the working directory used when building the archive bundle.  The
directory is  used as temporary location to store data from all data sources
specified in the file.  When the tar is built, the content of that directory
is removed.

```
WORKDIR <relative or absolute path>
```