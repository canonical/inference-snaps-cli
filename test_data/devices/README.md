# Machine hardware info

Each subdirectory represents a single machine.
This directory contains files with raw data from the respective machine.
It optionally contains a `hardware-info.json` file which is the output of the `cmd/hardware-info` application.

## cpuinfo.txt

```
cat /proc/cpuinfo
```

## lspci.txt

```
lspci -vmmnD
```

## uname-m.txt

```
uname -m
```

## disk.json

Normally obtained with `hardware-info`, the total and available disk space can also be looked up with `df`.

```
{
  "/": {
    "total": <bytes>,
    "avail": <bytes>
  },
  "/var/lib/snapd/snaps": {
    "total": <bytes>,
    "avail": <bytes>
  }
}
```

## memory.txt

Memory and swap sizes can be looked up with `cat /proc/meminfo`.

```
{
  "total_ram": <bytes>,
  "total_swap": <bytes>
}
```

## additional-properties.json (optional)

Normally the additional properties are looked up on the host, using vendor specific tools.
This is not possible during testing when we do not have access to the host.
If any additional properties need to be added to PCI devices, add it to this file, based on the device slot address.

```
{
  "0000:00:02.0": {
    "vram": "14482374656",
    "compute_capability": "12.4"
  }
}
```

## hardware-info.json (optional)

If we have access to the machine, and we can run the hardware-info application on it, add the output to this file.
It is used for cross-validation of the output between the actual hardware-info vs the hardware info from test data.

```
stack-utils.hardware-info --pretty --friendly
```
