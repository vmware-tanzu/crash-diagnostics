# Flare
A tool that collects server data for troubleshooting.

A human operator would run flare directly on a running server to take a snapshot of the current state of the server for analysis.  Flare can collect data from different data sources specified at runtime.  The collected data is packaged and compressed in a tar bundle.


## Pre-requisites
 * Assumes Linux/Unix

## Flare CLI
Flare is designed to support multiple commands in the future to help diagnose and troubleshoot cluster machines.

### flare out
This command can gather information from multiple data sources running on a machine. Flare out uses the flare.file (see below)
to specify what and how to collect machine information such as log files and output from commands.  The result is then bundled into an archive 
file.

The following example uses default `./flare.file` if one exists in the current directory and output the result to file `flare.tar.gz`:

```
flare out --output flare.tar.gz
```

A flare file can be specified using `--file` flag:

```
flare out -file flare.file -output out.tar.gz
```

## Flare.file Format
Some flare commands, like `flare out` above, use a file (by default named `flare.file`) that contains a simple dialect to declartively specify commands executed by the flare program.  

Flare uses a simple format for the command file (similar to Dockerfile).  Each command uses a single line:

```
COMMAND [arguments]
```

### Example flare.file
```
FROM local
WORKDIR /tmp/flareout

COPY /var/log/system.log
CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i
```

### Flare.file Commands
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


## Compile and Running Flare
Flare  is written and Go 1.11 or later.  Clone or download the source to your local directory.  From the project's root directory, compile the code with the
following:

```
GO111MODULE="on" go install .
```

This should place the compiled flare binary in `$(go env GOPATH)/bin`.  You can test this with:
```
flare --help
```
If this does not work properly, ensure that your Go environment is setup properly.  Next setup a sample flare.file to test with your environment:

```
FROM local
WORKDIR /tmp/flareout

CAPTURE df -h
CAPTURE df -i
CAPTURE netstat -an
CAPTURE ps -ef
CAPTURE lsof -i
```

In the directory where the flare.file is located, run flare to generate output bundle:
```
flare out --loglevel debug
```
You should see log messages on the screen showing flare working:
```
DEBU[0000] Parsing flare script
DEBU[0000] Parsing [1: FROM local]
DEBU[0000] FROM parsed OK
DEBU[0000] Parsing [2: WORKDIR /tmp/flareout]
...
DEBU[0002] Archiving [/tmp/flareout] in out.tar.gz
DEBU[0002] Archived /tmp/flareout/df_-h.txt
DEBU[0002] Archived /tmp/flareout/df_-i.txt
DEBU[0002] Archived /tmp/flareout/lsof_-i.txt
DEBU[0002] Archived /tmp/flareout/netstat_-an.txt
DEBU[0002] Archived /tmp/flareout/ps_-ef.txt
DEBU[0002] Archived /tmp/flareout/system.log
INFO[0002] Created archive out.tar.gz
INFO[0002] Output done
```