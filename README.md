# ML Snap Utils

This repo contains utilities used in snapping machine learning (AI) workloads.

## Build

The CLIs included in this repo can be built using the following commands.

Hardware Info:

```bash
go build github.com/canonical/hardware-info/cmd/hardware-info
```

Select Stack:

```bash
go build github.com/canonical/hardware-info/cmd/select-stack
```

To build a snap for these applications, run:

```bash
snapcraft -v
```

Then install the snap and connect the required interfaces:

```bash
sudo snap install --dangerous ./ml-snap-utils_*.snap
sudo snap connect ml-snap-utils:hardware-observe 
```

## Usage

### Hardware Info

See [readme](cmd/hardware-info/README.md).

### Select Stack

The output from `hardware-info` can be piped into `select-stack`.
You need to provide the location of the stack definitions from which the selection should be made.

The result is written as json to STDOUT, while any other log messages are available on STDERR.

Example:

```bash
$ ml-snap-utils.hardware-info | ml-snap-utils.select-stack --stacks=test_data/stacks/
2024/12/10 11:28:03 Vendor specific info for Intel GPU not implemented
2024/12/10 11:28:03 Stack cpu-f32 not selected: not enough memory
2024/12/10 11:28:03 Stack fallback-cpu matches. Score = 4.000000
2024/12/10 11:28:03 Stack fallback-gpu not selected: any: could not find a required device
2024/12/10 11:28:03 Stack llamacpp-avx2 matches. Score = 3.200000
2024/12/10 11:28:03 Stack llamacpp-avx512 not selected: any: could not find a required device
{"name":"fallback-cpu","components":["llamacpp","model-q4-k-m-gguf"],"score":4}
```