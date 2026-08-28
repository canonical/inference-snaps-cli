# Platform Accelerator Contract

Inference snaps should depend on capabilities, not on the device-node names
of individual SoCs.

## Application Contracts

An inference snap that needs a directly accessible NPU declares one stable
device plug and a runtime-content plug for each incompatible runtime ABI it
supports:

```yaml
plugs:
  inference-npu:
    interface: custom-device
    custom-device: inference-npu
  inference-npu-runtime-qcom-htp:
    interface: content
    content: inference-npu-runtime-qcom-htp-v1
    target: $SNAP/inference-npu-runtime-qcom-htp
```

The inference server app must list both plugs. The engine manifest lists the
same names under `snap-connections`, so engine selection rejects an engine
whose required platform contract is not connected.

The application snap does not declare `/dev/fastrpc-*`, `/dev/drpai*`,
`/dev/apusys*`, or any other vendor-specific path. It also does not read
driver libraries from the host `/usr` tree.

## Platform Provider Responsibilities

Each platform gadget implements the stable device contract using its own
hardware details. Runtime contents are ABI- and layout-specific, so providers
use versioned content identifiers rather than publishing incompatible files
under one generic identifier. A Qualcomm Hexagon provider may use
`inference-npu-runtime-qcom-htp-v1`, while another backend uses its own
content contract.

| Platform | Device contract implementation | Runtime contract implementation |
| --- | --- | --- |
| Qualcomm Hexagon HTP | FastRPC CDSP and DMA-heap nodes | Versioned Qualcomm HTP runtime content |
| MediaTek APUSYS | APUSYS/MDLA device nodes | Versioned MediaTek runtime content |
| Renesas DRP-AI | DRP-AI device nodes | Versioned Renesas runtime content |

For Qualcomm, the gadget slot is expected to expose at least
`/dev/fastrpc-cdsp`, `/dev/fastrpc-cdsp-secure`, and
`/dev/dma_heap/system`. The content provider must make the architecture- and
firmware-matched FastRPC library available below `lib/` and the QCOM tree
available below `usr/share/qcom/` in the content mount.

The gadget remains responsible for ordinary Unix ownership and permissions,
udev tagging, firmware, and any system daemon required by the platform. A
custom-device connection grants AppArmor/device-cgroup access, but does not
replace those host-level requirements.

## Why This Scales

The number of platform implementations can grow without multiplying the
device-access surface of application snaps. New platforms add a
gadget/provider mapping and an engine/runtime implementation only when their
inference backend differs. Existing model snaps keep the same `inference-npu`
device plug; runtime content plugs are added per incompatible ABI family, not
per board or SoC.

This deliberately does not create a new snapd interface for Qualcomm,
MediaTek, Renesas, or individual boards. The existing `custom-device` and
`content` interfaces carry the access policy. A new content plug is justified
only when a runtime ABI or file-layout contract cannot be shared safely.

If two platforms cannot share a runtime ABI or layout, they should still use
the same device capability plug and expose a backend-specific runtime through
the provider contract selected by the engine. A new snapd interface is only
needed when the access policy cannot be expressed by `custom-device`,
`content`, or another existing interface.
