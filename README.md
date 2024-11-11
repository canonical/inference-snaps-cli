# Hardware Detect

This program detects system hardware and provides a summary in JSON format.

## Required dependencies

External dependencies are kept to a minimum.
The only extra dependency is the new standard Unix Syscall interface library.
This can be installed by running:

```
go get golang.org/x/sys/unix
```

## Build

To build the CLI for hardware-info, run the following command in the root of this repository:

```bash
go build -o build/hardware-info ./cmd/cli
```

## Usage

A help message is printed out when providing the `-h` or `--help` flags.

```bash
$ build/hardware-info -h
Usage of build/hardware-info:
  -file string
        Output json to this file. Default output is to stdout.
  -pretty
        Output pretty json. Default is compact json.
```

By default, the `hardware-info` application will print out a summary of the host system to STDOUT in compact JSON format.
By specifying the `--pretty` flag, the JSON will be formatted for easier readability.
The `--file` argument allows writing the JSON data to a file, rather than to STDOUT.

Errors and warnings are printed to STDERR.
