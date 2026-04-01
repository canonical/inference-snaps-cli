# inference-snaps-cli — Complete Error Reference

All errors produced by first-party code, with their **full literal output** as seen by the user.
Errors from external packages are the boundary; they are listed as the leaf cause but not traced further.

## How cobra formats errors

`SilenceUsage: true` is set on the root command; `SilenceErrors` is **not** set.
Every error returned by a `RunE` (or `PersistentPreRunE`) function is printed to **stderr** by cobra as:

```
Error: <err.Error()>
```

followed, only when a command/flag parse fails (not a `RunE` error), by:

```
Run '<command> --help' for usage.
```

Errors printed directly with `fmt.Printf`/`fmt.Fprintf` are noted explicitly. All other entries in
this document represent errors that cobra prefixes with `Error: `.

---

## Table of Contents
1. [cmd/cli/main.go](#1-cmdclimaingo)
2. [cmd/cli/commands/](#2-cmdclicommands)
    - [use-engine.go](#21-use-enginego)
    - [list-engines.go](#22-list-enginesgo)
    - [show-engine.go](#23-show-enginego)
    - [status.go](#24-statusgo)
    - [chat.go](#25-chatgo)
    - [get.go](#26-getgo)
    - [set.go](#27-setgo)
    - [prune-cache.go](#28-prune-cachego)
    - [run.go](#29-rungo)
    - [show-machine.go](#210-show-machinego)
    - [version.go](#211-versiongo)
    - [debug/validate.go](#212-debugvalidatego)
    - [debug/select.go](#213-debugselectgo)
    - [debug/chat.go](#214-debugchatgo)
3. [cmd/cli/common/](#3-cmdclicommon)
    - [engine.go](#31-enginego)
    - [endpoints.go](#32-endpointsgo)
    - [component.go](#33-componentgo)
    - [errors.go](#34-errorsgo)
4. [pkg/engines/](#4-pkgengines)
    - [load.go](#41-loadgo)
    - [validate.go](#42-validatego)
    - [devices.go](#43-devicesgo)
    - [busses.go](#44-bussesgo)
    - [cpu.go](#45-cpugo)
5. [pkg/selector/](#5-pkgselector)
    - [select_stack.go](#51-select_stackgo)
    - [pci/pci.go](#52-pcipci-go)
    - [pci/properties.go](#53-pcipropertiesgo)
6. [pkg/hardware_info/](#6-pkghardware_info)
    - [hardware-info.go](#61-hardware-infogo)
    - [memory/](#62-memory)
    - [cpu/](#63-cpu)
    - [disk/](#64-disk)
    - [pci/pci_devices.go](#65-pcipci_devicesgo)
    - [pci/lspci.go](#66-pcilspcigo)
    - [pci/nvidia/](#67-pcinvidia)
    - [pci/amd/](#68-pciamd)
    - [pci/intel/](#69-pciintel)
7. [pkg/snap_store/snap_store.go](#7-pkgsnap_storesnap_storego)
8. [pkg/storage/](#8-pkgstorage)
    - [cache.go](#81-cachego)
    - [config.go](#82-configgo)
    - [snapctl_storage.go](#83-snapctl_storagego)
9. [pkg/utils/utils.go](#9-pkgutilsutilsgo)

---

## 1. `cmd/cli/main.go`

These are printed **directly** with `fmt.Printf` (not through cobra), so there is no cobra `Error:` prefix added — but the strings themselves include `"Error: "`.

| Full stdout output | Condition |
|---|---|
| `Error: could not retrieve snap services: <snapctl error>` | `snapctl.Services().Run()` fails at startup; printed to stdout then `return` |
| `Error: <err>` | `appendCommandToGroup` fails (group ID `"basic"` not found — programming invariant); printed to stdout then `return` |

`PersistentPreRunE` can return one error (cobra then prints it to stderr):

| Full stderr output | Condition |
|---|---|
| `Error: <os.Setenv error>` | `os.Setenv("VERBOSE", "true")` fails when `--verbose` is passed |

---

## 2. `cmd/cli/commands/`

### 2.1 `use-engine.go`

> Shell-completion helper `validateArgs` prints directly to **stdout** (not through cobra RunE):
> `Error: loading engines: <LoadManifests error>`

#### `run()` — returned to cobra → stderr

| Full stderr output | Condition |
|---|---|
| `Error: permission denied, try again with sudo` | Not root user |
| `Error: cannot specify both engine name and --auto flag` | `--auto` used with a positional argument |
| `Error: cannot specify both engine name and --fix flag` | `--fix` used with a positional argument |
| `Error: engine name not specified` | No `--auto`, no `--fix`, and no positional argument |

#### `autoSelectEngine()` → propagated to `run()` → stderr

| Full stderr output | Condition |
|---|---|
| `Error: loading engines: <LoadManifests error>` | `engines.LoadManifests` fails — see [§4.1](#41-loadgo) |
| `Error: machine info: <hardware_info.Get error>` | `hardware_info.Get` fails — see [§6.1](#61-hardware-infogo) |
| `Error: scoring engines: <ScoreEngines error>` | `selector.ScoreEngines` fails — see [§5.1](#51-select_stackgo) |
| `Error: finding top engine: no compatible engines found` | No stable, compatible engine exists |
| `Error: engine "<name>" not found` | Named engine manifest missing (ErrManifestNotFound) |
| `Error: loading engine manifest: <OS/YAML error>` | Other error reading the engine manifest |
| `Error: checking installed components: checking installed component directory "<comp>": <os error>` | `os.Stat` on a component directory fails with non-ErrNotExist |
| `Error: checking installed components: component "<comp>" exists but is not a directory` | Component path exists but is a regular file |
| `Error: timed out while installing "<comp>":\nMonitor the installation progress with "snap changes"\n\nRerun this command once the installation is complete` | Component not fully installed within 60-minute timeout |
| `Error: snap not known to the store:\nRerun this command after manually installing "<comp>"` | snapctl reports the snap is unknown to the store |
| `Error: installing "<comp>": <snapctl error>` | Any other unhandled `snapctl.InstallComponents` error |
| `Error: getting active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails |
| `Error: un-setting engine configurations: <snapctl error>` | `Config.Unset(".", EngineConfig)` fails |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load error inside `UnsetEngineConfig` (when `unsetUserOverrides=true`) |
| `Error: un-setting configuration "<key>": <snapctl error>` | `Config.Unset(k, UserConfig)` fails for a key |
| `Error: setting active engine: <snapctl error>` | `Cache.SetActiveEngine` fails |
| `Error: setting engine configuration "<key>": <snapctl error>` | `Config.SetDocument` fails |

#### `fixActiveEngine()` → propagated to `run()` → stderr

| Full stderr output | Condition |
|---|---|
| `Error: getting active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails |
| `Error: no active engine to fix` | Active engine name is empty |
| `Error: loading active engine manifest: <OS/YAML error>` | Other manifest load error (non-ErrManifestNotFound) |
| *(all autoSelectEngine or switchEngine errors above)* | ErrManifestNotFound → autoSelectEngine; other errors from installMissingComponents / UnsetEngineConfig / SetEngineConfig |

#### Non-fatal, printed directly to stderr (no cobra prefix)

| Full stderr output | Condition |
|---|---|
| `Warning: unable to get component sizes: <snap_store error>` | `snap_store.ComponentSizes()` fails **and** `--verbose` is set |
| `Warning: unable to get component sizes` | `snap_store.ComponentSizes()` fails **and** `--verbose` is **not** set |

---

### 2.2 `list-engines.go`

| Full stderr output | Condition |
|---|---|
| `Error: loading engines: <LoadManifests error>` | `engines.LoadManifests` fails — see [§4.1](#41-loadgo) |
| `Error: machine info: <hardware_info.Get error>` | `hardware_info.Get` fails — see [§6.1](#61-hardware-infogo) |
| `Error: scoring engines: <ScoreEngines error>` | `selector.ScoreEngines` fails — see [§5.1](#51-select_stackgo) |
| `Error: get active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails |
| `Error: unknown format "<value>"` | `--format` value is not `"table"` or `"json"` |
| `Error: json: <json error>` | `json.MarshalIndent` fails (internal invariant) |
| `Error: adding data to table: <tablewriter error>` | `tablewriter.Bulk` fails (external library) |
| `Error: rendering table: <tablewriter error>` | `tablewriter.Render` fails (external library) |

Non-fatal, printed directly to stderr:

| Full stderr output | Condition |
|---|---|
| `No engines found.` | Engine list is empty |

---

### 2.3 `show-engine.go`

> Shell-completion helper `validateArgs` prints directly to **stdout** (not through cobra RunE):
> `Error: loading engines: <LoadManifests error>`

| Full stderr output | Condition |
|---|---|
| `Error: get active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails (no-arg path) |
| `Error: no active engine` | Active engine name is empty |
| `Error: loading engines: <LoadManifests error>` | `engines.LoadManifests` fails — see [§4.1](#41-loadgo) |
| `Error: machine info: <hardware_info.Get error>` | `hardware_info.Get` fails — see [§6.1](#61-hardware-infogo) |
| `Error: scoring engines: <ScoreEngines error>` | `selector.ScoreEngines` fails — see [§5.1](#51-select_stackgo) |
| `Error: engine "<name>" does not exist` | Named engine not found in scored list |
| `Error: json: <json error>` | `json.MarshalIndent` fails |
| `Error: yaml: <yaml error>` | `yaml.Marshal` fails |
| `Error: unknown format "<value>"` | `--format` is not `"json"` or `"yaml"` |

---

### 2.4 `status.go`

| Full stderr output | Condition |
|---|---|
| `Error: unknown format "<value>"` | `--format` is not `"json"` or `"yaml"` |
| `Error: getting active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails |
| `Error: no engine is active` | Active engine name is empty |
| `Error: getting services: <snapctl error>` | `snapctl.Services().Run()` fails |
| `Error: unexpected service name format: "<name>"` | A service name has no `.` separator |
| `Error: looking up active engine: <snapctl error>` | `Cache.GetActiveEngine()` inside `EngineComponentSettings` fails |
| `Error: no active engine` | Active engine name empty inside `EngineComponentSettings` |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load fails inside `EngineComponentSettings` |
| `Error: SNAP_COMPONENTS env var not set` | `SNAP_COMPONENTS` not in environment |
| `Error: reading <path>: <os error>` | `os.ReadFile` fails on a `component.yaml` |
| `Error: unmarshaling <path>: <yaml error>` | `yaml.Unmarshal` fails on a `component.yaml` |
| `Error: OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` | Deprecated env key found in component environment |
| `Error: unsupported protocol "<proto>" for server "<server>" in component "<comp>"` | Protocol is not `"http"` or `"https"` |
| `Error: getting "http.port": <snapctl error>` | `Config.Get("http.port")` fails |
| `Error: yaml: <yaml error>` | `yaml.Marshal` of the status struct fails |
| `Error: json: <json error>` | `json.MarshalIndent` of the status struct fails |

---

### 2.5 `chat.go`

| Full stderr output | Condition |
|---|---|
| `Error: looking up active engine: <snapctl error>` | `Cache.GetActiveEngine()` inside `EngineComponentSettings` fails |
| `Error: no active engine` | Active engine name empty inside `EngineComponentSettings` |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load fails inside `EngineComponentSettings` |
| `Error: SNAP_COMPONENTS env var not set` | `SNAP_COMPONENTS` not in environment |
| `Error: reading <path>: <os error>` | `os.ReadFile` fails on a `component.yaml` |
| `Error: unmarshaling <path>: <yaml error>` | `yaml.Unmarshal` fails on a `component.yaml` |
| `Error: OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` | Deprecated env key |
| `Error: unsupported protocol "<proto>" for server "<server>" in component "<comp>"` | Unsupported protocol |
| `Error: getting "http.port": <snapctl error>` | `Config.Get("http.port")` fails |
| `Error: "openai" not found in server endpoints` | No `openai` key in the endpoints map |
| `Error: services: <snapctl error>` | `snapctl.Services(serviceName).Run()` fails |
| `Error: server not active\n\n<suggestion string>` | The `.server` service `Current` state is `"inactive"` |
| `Error: initializing readline: <error>` | `readline.NewEx` fails |
| `Error: invalid base URL: <url.Parse error>` | `url.Parse(baseUrl)` fails in `handshake` |
| `Error: connection refused\n\n<suggest-start>\n<suggest-logs>` | TCP dial to server is refused in `handshake` |
| `Error: no models available on server\n\n<suggest-start>\n<suggest-logs>` | Server returns 503 `unavailable_error` for >60 s in `checkServerReady` |
| `Error: api: <api error>` | Server returns a non-503 API error in `checkServerReady` or `lookupModelName` |
| `Error: <other error>\n\n<suggest-logs>` | Any other error from `checkServerReady` or `lookupModelName` |
| `Error: no models available on server\n\n<suggest-start>\n<suggest-logs>` | Server returns 503 `unavailable_error` for >60 s in `lookupModelName` |
| `Error: server returned no models\n\n<suggest-start>\n<suggest-logs>` | Model list is empty for >60 s in `lookupModelName` |
| `Error: server returned multiple models; expected one: <name1>, <name2>` | Server lists more than one model and `--model` not specified |
| `Error: connection refused\n\n<suggest-logs>` | TCP connection refused mid-stream in `processStream` |
| `Error: connection closed by server\n\n<suggest-logs>` | Unexpected EOF mid-stream in `processStream` |
| `Error: <stream error>\n\n<suggest-logs>` | Any other streaming error in `processStream` |

---

### 2.6 `get.go`

| Full stderr output | Condition |
|---|---|
| `Error: getting value of "<key>": <snapctl error>` | `Config.Get(key)` fails |
| `Error: no value set for key "<key>"` | Key exists in schema but has no stored value |
| `Error: serializing value: <yaml error>` | `yaml.Marshal` fails on the single-key result |
| `Error: getting values: <snapctl error>` | `Config.GetAll()` fails |
| `Error: serializing values: <yaml error>` | `yaml.Marshal` fails on the full config map |

Non-fatal, printed directly to stderr:

| Full stderr output | Condition |
|---|---|
| `Warning: "<key>" configuration field is deprecated!` | Requested key is in the deprecated-config list and output is a terminal |

---

### 2.7 `set.go`

| Full stderr output | Condition |
|---|---|
| `Error: permission denied, try again with sudo` | Not root user |
| `Error: key must not start with an equal sign` | First character of argument is `=` |
| `Error: expected key=value, got "<arg>"` | No `=` in the argument |
| `Error: "<key>" is read-only` | User tries to set a deprecated/internal config key |
| `Error: setting "<key>" to "<value>": checking key: <snapctl error>` | `Config.Get(key)` fails during UserConfig key validation |
| `Error: setting "<key>" to "<value>": unknown key` | Key has no existing value (not in schema) |
| `Error: setting "<key>" to "<value>": <snapctl error>` | `storage.Set` (snapctl) fails |

---

### 2.8 `prune-cache.go`

| Full stderr output | Condition |
|---|---|
| `Error: permission denied, try again with sudo` | Not root user |
| `Error: <snapctl error>` | `Cache.GetActiveEngine()` fails (error returned unwrapped) |
| `Error: no active engine found` | `engines.ErrManifestNotFound` for the active engine |
| `Error: loading engine manifest: <OS/YAML error>` | Any other manifest load error for active engine |
| `Error: cannot prune the active engine "<name>"` | `--engine` value matches the active engine |
| `Error: "<name>" not found` | `engines.ErrManifestNotFound` for the named `--engine` |
| `Error: loading engine manifest: <OS/YAML error>` | Other load error for named engine |
| `Error: getting list of inactive engines: <LoadManifests or snapctl error>` | `engines.LoadManifests` or `Cache.GetActiveEngine` fails inside `inactiveEngines()` |
| `Error: loading manifests: <LoadManifests error>` | `engines.LoadManifests` fails inside `getAllComponentsToRemove` or `pruneAllInactiveEngines` |
| `Error: checking installed component directory "<comp>": <os error>` | `os.Stat` fails with non-not-exist error inside `calculateRemovableComponents` |
| `Error: component "<comp>" exists but is not a directory` | Component path is a file inside `calculateRemovableComponents` |
| `Error: un-setting engine configurations: <snapctl error>` | `Config.Unset(".", EngineConfig)` fails inside `pruneEngine` |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load error inside `UnsetEngineConfig` |
| `Error: un-setting configuration "<key>": <snapctl error>` | `Config.Unset(k, UserConfig)` fails |
| `Error: removing components: <snapctl error>` | `snapctl.RemoveComponents` fails |

Non-fatal, printed directly to **stdout** (no cobra prefix):

| Full stdout output | Condition |
|---|---|
| `Warning: unable to get component sizes: <snap_store error>` | `snap_store.ComponentSizes()` fails **and** `--verbose` is set |
| `Warning: unable to get component sizes` | `snap_store.ComponentSizes()` fails **and** `--verbose` is **not** set |

---

### 2.9 `run.go`

| Full stderr output | Condition |
|---|---|
| `Error: unexpected number of arguments, expected 1 got <n>` | Not exactly 1 positional argument |
| `Error: getting active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails inside `waitForComponents` |
| `Error: no active engine` | Active engine name is empty inside `waitForComponents` |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load error inside `waitForComponents` |
| `Error: SNAP_COMPONENTS env var not set` | `SNAP_COMPONENTS` not in environment inside `checkMissingComponents` |
| `Error: <n>s timeout waiting for required components: <comp1>, <comp2>` | Components not mounted after 3600 s |
| `Error: looking up active engine: <snapctl error>` | `Cache.GetActiveEngine()` inside `EngineComponentSettings` fails |
| `Error: no active engine` | Active engine name empty inside `EngineComponentSettings` |
| `Error: loading engine manifest: <OS/YAML error>` | Manifest load fails inside `EngineComponentSettings` |
| `Error: SNAP_COMPONENTS env var not set` | `SNAP_COMPONENTS` not in environment inside `LoadEngineEnvironment` |
| `Error: reading <path>: <os error>` | `os.ReadFile` fails on a `component.yaml` |
| `Error: unmarshaling <path>: <yaml error>` | `yaml.Unmarshal` fails on `component.yaml` |
| `Error: invalid env var "<entry>"` | An env entry in the component YAML has no `=` separator |
| `Error: setting "COMPONENT": <os error>` | `os.Setenv(COMPONENT, ...)` fails |
| `Error: unsetting "COMPONENT": <os error>` | `os.Unsetenv(COMPONENT)` fails |
| `Error: setting "<key>": <os error>` | `os.Setenv(k, v)` fails for an engine env var |
| `Error: <exec error>` | The subprocess exits non-zero (exec error returned directly) |

---

### 2.10 `show-machine.go`

| Full stderr output | Condition |
|---|---|
| `Error: memory info: error reading /proc/meminfo: <os error>` | `/proc/meminfo` unreadable |
| `Error: memory info: error parsing MemTotal: error parsing kB value: <strconv error>` | MemTotal kB parse fails |
| `Error: memory info: error parsing MemTotal: error parsing byte value: <strconv error>` | MemTotal byte parse fails |
| `Error: memory info: error parsing SwapTotal: error parsing kB value: <strconv error>` | SwapTotal kB parse fails |
| `Error: memory info: error parsing SwapTotal: error parsing byte value: <strconv error>` | SwapTotal byte parse fails |
| `Error: cpu info: reading /proc/cpuinfo: <os error>` | `/proc/cpuinfo` unreadable |
| `Error: cpu info: getting host uname: <exec error>` | `uname --machine` binary missing or fails |
| `Error: cpu info: unsupported architecture: <arch>` | `uname -m` output not in lookup table |
| `Error: cpu info: unsupported architecture: <arch>` | Architecture from `/proc/cpuinfo` not `amd64`/`arm64` |
| `Error: cpu info: <strconv error>` | Integer field parsing in `/proc/cpuinfo` fails |
| `Error: disk info: statfs failed: <os error>` | `unix.Statfs` on the snap storage path fails |
| `Error: pci devices: getting lspci data: <exec error>` | `lspci -vmmnD` fails |
| `Error: pci devices: parsing lspci data: unexpected format for pci slot: <value>` | PCI slot field has wrong format |
| `Error: pci devices: parsing lspci data: cannot parse pci bus number: <value>` | Bus-number hex parse fails |
| `Error: json: <json error>` | `json.MarshalIndent` fails |
| `Error: yaml: <yaml error>` | `yaml.Marshal` fails |
| `Error: unknown format "<value>"` | `--format` is not `"json"` or `"yaml"` |

---

### 2.11 `version.go`

| Full stderr output | Condition |
|---|---|
| `Error: json: <json error>` | `json.MarshalIndent` fails |
| `Error: yaml: <yaml error>` | `yaml.Marshal` fails |
| `Error: unknown format "<value>"` | `--format` is not `"json"` or `"yaml"` |

---

### 2.12 `debug/validate.go`

| Full stderr output | Condition |
|---|---|
| `Error: no engine manifest specified` | Zero positional arguments |
| `Error: not all manifests are valid` | At least one `engines.Validate` call returned an error |

Per-file errors are printed **directly to stdout** (not returned, no cobra prefix):

| Full stdout output | Condition |
|---|---|
| `❌ <path>: manifest file must be called engine.yaml: <path>` | File name doesn't end with `engine.yaml` |
| `❌ <path>: manifest file does not exist: <path>` | `os.Stat` → ErrNotExist |
| `❌ <path>: getting file info: <os error>` | `os.Stat` fails for any other reason |
| `❌ <path>: reading file: <os error>` | `os.ReadFile` fails |
| `❌ <path>: empty yaml data` | File is empty after trimming |
| `❌ <path>: yaml: <yaml error>` | YAML parse error or unknown field |
| `❌ <path>: required field is not set: name` | `name` field missing |
| `❌ <path>: engine directory name should match name in manifest: <dir> != <name>` | Directory name ≠ manifest name |
| `❌ <path>: required field is not set: description` | `description` field missing |
| `❌ <path>: required field is not set: vendor` | `vendor` field missing |
| `❌ <path>: required field is not set: grade` | `grade` field missing |
| `❌ <path>: grade should be 'stable' or 'devel'` | `grade` is set but invalid |
| `❌ <path>: parsing memory: <strconv error>` | `memory` value not a valid size string |
| `❌ <path>: parsing disk space: <strconv error>` | `disk-space` value not a valid size string |
| `❌ <path>: configuration field <key> is not a primitive value: <value>` | A configurations entry is non-primitive |
| `❌ <path>: invalid device: allof <n>/<total>: cpu: architecture field required` | CPU device in `allof` missing `architecture` |
| `❌ <path>: invalid device: allof <n>/<total>: cpu: invalid architecture: <value>` | Unknown architecture |
| `❌ <path>: invalid device: allof <n>/<total>: cpu: invalid field for amd64: <field>` | Invalid field for amd64 CPU |
| `❌ <path>: invalid device: allof <n>/<total>: cpu: invalid field for arm64: <field>` | Invalid field for arm64 CPU |
| `❌ <path>: invalid device: allof <n>/<total>: gpu: invalid bus: <value>` | Unknown bus type for GPU |
| `❌ <path>: invalid device: allof <n>/<total>: gpu: usb device validation not implemented` | GPU on USB bus |
| `❌ <path>: invalid device: allof <n>/<total>: gpu: pci device: invalid field: <field>` | Invalid field for PCI GPU |
| `❌ <path>: invalid device: allof <n>/<total>: npu: <bus/field error>` | NPU device errors (same bus/field patterns) |
| `❌ <path>: invalid device: allof <n>/<total>: typeless: <bus/field error>` | Typeless device errors |
| `❌ <path>: invalid device: allof <n>/<total>: invalid device type: <value>` | `type` not `cpu`, `gpu`, `npu`, or `""` |
| *(same set with `anyof` instead of `allof`)* | Same errors for `anyof` devices |
| `✅ <path>` | Manifest is valid |

---

### 2.13 `debug/select.go`

| Full stderr output | Condition |
|---|---|
| `Error: decoding hardware info: <yaml error>` | `yaml.Decode` on stdin fails |
| `Error: loading engines from directory: <LoadManifests error>` | `engines.LoadManifests` fails — see [§4.1](#41-loadgo) |
| `Error: checking engines: parsing required memory: <strconv error>` | `utils.StringToBytes` fails for a manifest memory field |
| `Error: checking engines: total memory not reported by host system` | `hardwareInfo.Memory.TotalRam == 0` |
| `Error: checking engines: parsing required disk space: <strconv error>` | `utils.StringToBytes` fails for a manifest disk-space field |
| `Error: checking engines: disk space not reported by host system` | Snap storage path not in disk info map |
| `Error: finding top engine: no compatible engines found` | No stable compatible engine |
| `Error: json: <json error>` | `json.MarshalIndent` fails |
| `Error: yaml: <yaml error>` | `yaml.Marshal` fails |
| `Error: unknown format "<value>"` | `--format` is not `"json"` or `"yaml"` |

---

### 2.14 `debug/chat.go`

| Full stderr output | Condition |
|---|---|
| `Error: the --base-url parameter is required` | `--base-url` flag not provided |
| `Error: initializing readline: <error>` | `readline.NewEx` fails |
| `Error: invalid base URL: <url.Parse error>` | `url.Parse(baseUrl)` fails in `handshake` |
| `Error: connection refused\n\n<suggest-start>\n<suggest-logs>` | TCP dial to server is refused |
| `Error: no models available on server\n\n<suggest-start>\n<suggest-logs>` | Server 503 `unavailable_error` for >60 s in `checkServerReady` or `lookupModelName` |
| `Error: api: <api error>` | Server returns a non-503 API error |
| `Error: <other error>\n\n<suggest-logs>` | Any other error from `checkServerReady` or `lookupModelName` |
| `Error: server returned no models\n\n<suggest-start>\n<suggest-logs>` | Model list empty for >60 s |
| `Error: server returned multiple models; expected one: <name1>, <name2>` | Server lists >1 model and `--model` not specified |
| `Error: connection refused\n\n<suggest-logs>` | TCP connection refused mid-stream |
| `Error: connection closed by server\n\n<suggest-logs>` | Unexpected EOF mid-stream |
| `Error: <stream error>\n\n<suggest-logs>` | Any other streaming error |

---

## 3. `cmd/cli/common/`

### 3.1 `engine.go`

These errors are propagated up through callers and ultimately reach cobra, which prefixes them with `Error: `.
The tables below show the message fragment that each function contributes; callers add their own prefix before it.

#### `EngineComponentSettings(ctx)` errors

| Error message fragment | Condition |
|---|---|
| `looking up active engine: <snapctl error>` | `Cache.GetActiveEngine()` fails |
| `no active engine` | Active engine name is empty |
| `loading engine manifest: <OS/YAML error>` | Manifest load error — see [§4.1](#41-loadgo) |
| `SNAP_COMPONENTS env var not set` | `SNAP_COMPONENTS` not in environment |
| `reading <path>: <os error>` | `os.ReadFile` fails on a component's `component.yaml` |
| `unmarshaling <path>: <yaml error>` | `yaml.Unmarshal` fails on `component.yaml` |

#### `LoadEngineEnvironment(ctx)` errors

| Error message fragment | Condition |
|---|---|
| *(all EngineComponentSettings errors above, propagated as-is)* | |
| `SNAP_COMPONENTS env var not set` | Second check, inside the env-expansion loop |
| `invalid env var "<entry>"` | Env entry in component YAML has no `=` |
| `setting "COMPONENT": <os error>` | `os.Setenv(COMPONENT, ...)` fails |
| `unsetting "COMPONENT": <os error>` | `os.Unsetenv(COMPONENT)` fails |
| `setting "<key>": <os error>` | `os.Setenv(k, v)` fails |

#### `SetEngineConfig(engine, ctx)` errors

| Error message fragment | Condition |
|---|---|
| `setting engine configuration "<key>": <snapctl error>` | `Config.SetDocument` fails |

#### `UnsetEngineConfig(engineName, unsetUserOverrides, ctx)` errors

| Error message fragment | Condition |
|---|---|
| `un-setting engine configurations: <snapctl error>` | `Config.Unset(".", EngineConfig)` fails |
| `loading engine manifest: <OS/YAML error>` | Manifest load error when `unsetUserOverrides=true` |
| `un-setting configuration "<key>": <snapctl error>` | `Config.Unset(k, UserConfig)` fails |

Non-fatal, printed directly to stderr (no cobra prefix):

| Full stderr output | Condition |
|---|---|
| `Warning: previously active engine "<name>" not found; skipping user configuration cleanup.` | `engines.ErrManifestNotFound` when `unsetUserOverrides=true` |

#### `ScoreEngines(ctx)` errors

| Error message fragment | Condition |
|---|---|
| `loading engines: <LoadManifests error>` | `engines.LoadManifests` fails — see [§4.1](#41-loadgo) |
| `machine info: <hardware_info.Get error>` | `hardware_info.Get(false)` fails — see [§6.1](#61-hardware-infogo) |
| `scoring engines: <ScoreEngines error>` | `selector.ScoreEngines` fails — see [§5.1](#51-select_stackgo) |

---

### 3.2 `endpoints.go`

#### `ServerEndpoints(ctx)` / `serverEndpoints()` / `serverHttpUrl()` errors

| Error message fragment | Condition |
|---|---|
| *(all EngineComponentSettings errors, propagated as-is)* | `EngineComponentSettings` fails |
| `OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` | Deprecated env key found |
| `unsupported protocol "<proto>" for server "<server>" in component "<comp>"` | Protocol not `"http"` or `"https"` |
| `getting "http.port": <snapctl error>` | `Config.Get("http.port")` fails |

---

### 3.3 `component.go`

#### `ComponentInstalled(component)` errors

| Error message fragment | Condition |
|---|---|
| `checking installed component directory "<comp>": <os error>` | `os.Stat` fails with non-ErrNotExist |
| `component "<comp>" exists but is not a directory` | Path exists but is a regular file |

---

### 3.4 `errors.go`

| Sentinel | Full value |
|---|---|
| `common.ErrPermissionDenied` | `permission denied, try again with sudo` |

---

## 4. `pkg/engines/`

### 4.1 `load.go`

#### `LoadManifests(manifestsDir)` errors

| Error message fragment | Condition |
|---|---|
| `<manifestsDir>: <os error>` | `os.ReadDir(manifestsDir)` fails |
| `<file path>: <os error>` | `os.ReadFile` fails on an engine YAML |
| `<manifestsDir>: <yaml error>` | `yaml.Unmarshal` fails on an engine YAML |

#### `LoadManifest(manifestsDir, engineName)` errors

| Error message fragment | Condition |
|---|---|
| `engine manifest not found: <os ErrNotExist>` | File does not exist; wraps `ErrManifestNotFound` |
| `<file path>: <os error>` | `os.ReadFile` fails for any other reason |
| `<manifestsDir>: <yaml error>` | `yaml.Unmarshal` fails |

---

### 4.2 `validate.go`

Errors from `Validate` are used by `debug/validate` (printed per-file, see [§2.12](#212-debugvalidatego)) and by `LoadManifest`; they are not surfaced directly by cobra.

| Error message fragment | Condition |
|---|---|
| `manifest file must be called engine.yaml: <path>` | File name doesn't end with `engine.yaml` |
| `manifest file does not exist: <path>` | `os.Stat` → ErrNotExist |
| `getting file info: <os error>` | `os.Stat` fails for any other reason |
| `reading file: <os error>` | `os.ReadFile` fails |
| `empty yaml data` | File is empty after trimming |
| `yaml: <yaml error>` | YAML parse error or unknown field |
| `required field is not set: name` | `name` field missing |
| `engine directory name should match name in manifest: <dir> != <name>` | Directory name ≠ manifest name |
| `required field is not set: description` | `description` field missing |
| `required field is not set: vendor` | `vendor` field missing |
| `required field is not set: grade` | `grade` field missing |
| `grade should be 'stable' or 'devel'` | `grade` is set but invalid |
| `parsing memory: <strconv error>` | `memory` value is not a valid size string |
| `parsing disk space: <strconv error>` | `disk-space` value is not a valid size string |
| `configuration field <key> is not a primitive value: <value>` | A configurations entry is non-primitive |
| `invalid device: allof <n>/<total>: <device error>` | A device in `allof` fails validation |
| `invalid device: anyof <n>/<total>: <device error>` | A device in `anyof` fails validation |

---

### 4.3 `devices.go`

| Error message fragment | Condition |
|---|---|
| `invalid device: allof <n>/<total>: <device error>` | A device in `allof` fails its own `validate()` |
| `invalid device: anyof <n>/<total>: <device error>` | A device in `anyof` fails its own `validate()` |
| `cpu: <cpu error>` | `validateCpu()` fails — see [§4.5](#45-cpugo) |
| `gpu: <bus/field error>` | `validateGpu()` fails — see [§4.4](#44-bussesgo) |
| `npu: <bus/field error>` | `validateNpu()` fails — see [§4.4](#44-bussesgo) |
| `typeless: <bus/field error>` | `validateTypelessDevice()` fails — see [§4.4](#44-bussesgo) |
| `invalid device type: <value>` | `type` field is not `cpu`, `gpu`, `npu`, or `""` |

---

### 4.4 `busses.go`

| Error message fragment | Condition |
|---|---|
| `invalid bus: <value>` | `bus` field is not `"pci"`, `"usb"`, or `""` |
| `usb device validation not implemented` | `bus == "usb"` |
| `pci device: invalid field: <field>` | A non-zero struct field is not in the PCI allow-list |

---

### 4.5 `cpu.go`

| Error message fragment | Condition |
|---|---|
| `architecture field required` | `device.Architecture == nil` |
| `invalid architecture: <value>` | Architecture is set but not `"amd64"` or `"arm64"` |
| `invalid field for amd64: <field>` | A non-zero struct field is not in the amd64 allow-list |
| `invalid field for arm64: <field>` | A non-zero struct field is not in the arm64 allow-list |

---

## 5. `pkg/selector/`

### 5.1 `select_stack.go`

#### `TopEngine(scoredEngines)` — sentinel

| Error message fragment | Condition |
|---|---|
| `no compatible engines found` | No engine with `Score > 0` and `Grade == "stable"` |

#### `checkEngine(hardwareInfo, manifest)` — returned through `ScoreEngines`

| Error message fragment | Condition |
|---|---|
| `parsing required memory: <strconv error>` | `utils.StringToBytes(*manifest.Memory)` fails |
| `total memory not reported by host system` | `hardwareInfo.Memory.TotalRam == 0` |
| `parsing required disk space: <strconv error>` | `utils.StringToBytes(*manifest.DiskSpace)` fails |
| `disk space not reported by host system` | Snap storage path not in disk info map |

*Device-level compatibility failures (wrong vendor, wrong type, etc.) are recorded as compatibility issues and reflected only in the score — they are never returned as errors.*

---

### 5.2 `pci/pci.go`

`Match` never returns an `error`. Incompatibility reasons are returned as `[]string` device issues, which surface in `CompatibilityIssues` fields and verbose output:

| Issue string | Condition |
|---|---|
| `no pci devices on host system` | `hostPciDevices` is empty |
| `device not found` | No device survives vendor/device ID filtering |
| `pci <slot>: vendor id mismatch: 0x<hex>` | Vendor ID doesn't match |
| `pci <slot>: device id mismatch: 0x<hex>` | Device ID doesn't match |
| `pci <slot>: device class 0x<hex> not of required type <type>` | Device class doesn't match required type |
| `pci <slot>: <err from checkProperties>` | Property check failure — see [§5.3](#53-pcipropertiesgo) |
| `pci <slot>: error checking snap connection "<conn>": <snapctl error>` | `snapctl.IsConnected` fails |
| `pci <slot>: "<conn>" is not connected` | Required snap interface not connected |
| `pci <slot>: no criteria met` | Device found but no scoring criterion matched |

---

### 5.3 `pci/properties.go`

Errors returned from `checkProperties` become issue strings in `pci.go` (prefixed with `pci <slot>: `).

| Error message fragment | Condition |
|---|---|
| `<strconv error>` | `utils.StringToBytes(*device.VRam)` fails — required vRAM is not a valid size string |
| `error parsing vRAM: <strconv error>` | `utils.StringToBytes(vram)` fails on the device-reported vRAM value |
| `not enough vRAM: <available bytes>` | Available vRAM < required |
| `unable to detect vRAM` | `"vram"` key absent from `AdditionalProperties` |
| `microarchitecture does not match: <actual>` | Microarch doesn't match required value |
| `unable to detect microarchitecture` | `"microarchitecture"` key absent from `AdditionalProperties` |

---

## 6. `pkg/hardware_info/`

### 6.1 `hardware-info.go`

#### `Get(friendlyNames)` errors — propagated up (prefixed by callers)

| Error message fragment | Condition |
|---|---|
| `memory info: <memory error>` | `memory.Info()` fails — see [§6.2](#62-memory) |
| `cpu info: <cpu error>` | `cpu.Info()` fails — see [§6.3](#63-cpu) |
| `disk info: <disk error>` | `disk.Info()` fails — see [§6.4](#64-disk) |
| `pci devices: <pci error>` | `pci.Devices(friendlyNames)` fails — see [§6.5](#65-pcipci_devicesgo) |

`show-machine` uses `fmt.Errorf("%s", err)` (bare, no extra prefix), so the above fragments appear directly after `Error: ` in the terminal. All other callers add their own prefix first (e.g. `machine info: ...` from `common/engine.go`).

---

### 6.2 `memory/`

#### `memory.Info()` / `parseProcMemInfo()` errors

| Error message fragment | Condition |
|---|---|
| `error reading /proc/meminfo: <os error>` | `os.ReadFile("/proc/meminfo")` fails |
| `error parsing MemTotal: error parsing kB value: <strconv error>` | `strconv.ParseInt` on the kB part of MemTotal fails |
| `error parsing MemTotal: error parsing byte value: <strconv error>` | `strconv.ParseInt` on a bare byte value of MemTotal fails |
| `error parsing SwapTotal: error parsing kB value: <strconv error>` | Same for SwapTotal |
| `error parsing SwapTotal: error parsing byte value: <strconv error>` | Same for SwapTotal |

---

### 6.3 `cpu/`

#### `cpu.Info()` / `InfoFromRawData()` / `parseProcCpuInfo()` / `debianArchitecture()` errors

| Error message fragment | Condition |
|---|---|
| `reading /proc/cpuinfo: <os error>` | `os.ReadFile("/proc/cpuinfo")` fails |
| `getting host uname: <exec error>` | `exec.Command("uname", "--machine")` fails |
| `unsupported architecture: <arch>` | `uname -m` output not in the debian-arch lookup table |
| `unsupported architecture: <arch>` | Architecture string from `/proc/cpuinfo` not `"amd64"` or `"arm64"` |
| `<strconv error>` | Parsing integer fields (processor index, implementer-id, part-number, variant, revision, cpu-khz) |

---

### 6.4 `disk/`

#### `disk.Info()` / `disk.InfoFromRawData()` errors

| Error message fragment | Condition |
|---|---|
| `statfs failed: <os error>` | `unix.Statfs(path, &fs)` fails |
| `error parsing df: not 6 columns` | A `df` output line has wrong field count |
| `error parsing df: error parsing 'total blocks' field: <strconv error>` | Total-blocks column parse fails |
| `error parsing df: error parsing 'available blocks' field: <strconv error>` | Available-blocks column parse fails |
| `df did not return info for all dirs` | Fewer parsed entries than expected directories |

---

### 6.5 `pci/pci_devices.go`

#### `pci.Devices(friendlyNames)` errors

| Error message fragment | Condition |
|---|---|
| `getting lspci data: <exec error>` | `exec.Command("lspci", "-vmmnD").Output()` fails |
| `parsing lspci data: unexpected format for pci slot: <value>` | PCI slot has wrong format |
| `parsing lspci data: cannot parse pci bus number: <value>` | Bus-number hex parse fails |

Additional-properties errors are non-fatal, printed directly to stderr (no cobra prefix):

| Full stderr output | Condition |
|---|---|
| `Warning: failed to get additional properties: AMD: vRAM: <os error>` | AMD `/sys/bus/pci/devices/<slot>/mem_info_vram_total` unreadable |
| `Warning: failed to get additional properties: AMD: vRAM: <strconv error>` | AMD vRAM file content not a valid integer |
| `Warning: failed to get additional properties: AMD: gfx architecture: <os error>` | `/sys/class/kfd/kfd/topology/nodes` unreadable |
| `Warning: failed to get additional properties: AMD: gfx architecture: gfx_target_version not found for device with pci slot <slot>` | No KFD node matched the device |
| `Warning: failed to get additional properties: AMD: gfx architecture: gfx_target_version is invalid for this device` | Value is `"0"` |
| `Warning: failed to get additional properties: AMD: gfx architecture: gfx_target_version has an unexpected format: <value>` | Hex string shorter than 6 chars |
| `Warning: failed to get additional properties: AMD: gfx architecture: parsing major version from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 0–1 |
| `Warning: failed to get additional properties: AMD: gfx architecture: parsing minor version from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 2–3 |
| `Warning: failed to get additional properties: AMD: gfx architecture: parsing revision from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 4–5 |
| `Warning: failed to get additional properties: AMD: gfx architecture: unexpected format for gfx_target_version: <line>` | Line not 2 space-separated parts |
| `Warning: failed to get additional properties: NVIDIA: vRAM: nvidia-smi: <exec error>` | `nvidia-smi` fails, no stdout output |
| `Warning: failed to get additional properties: NVIDIA: vRAM: nvidia-smi: <exec error>: <stdout>` | `nvidia-smi` fails with stdout error message |
| `Warning: failed to get additional properties: NVIDIA: vRAM: <strconv error>` | vRAM value string not a valid integer |
| `Warning: failed to get additional properties: NVIDIA: compute capability: nvidia-smi: <exec error>` | Same patterns for compute-capability query |
| `Warning: failed to get additional properties: NVIDIA: compute capability: nvidia-smi: <exec error>: <stdout>` | Same |
| `Warning: failed to get additional properties: Intel: vRAM: <exec error>` | `clinfo --json` fails |
| `Warning: failed to get additional properties: Intel: vRAM: parsing clinfo json: <json error>` | JSON parse of clinfo output fails |
| `Warning: failed to get additional properties: Intel: vRAM: clinfo: no devices found` | `clinfo.Devices` is empty |
| `Warning: failed to get additional properties: Intel: vRAM: clinfo: no online devices found` | `clinfo.Devices[0].Online` is empty |

---

### 6.6 `pci/lspci.go`

#### `ParseLsPci()` errors

| Error message fragment | Condition |
|---|---|
| `unexpected format for pci slot: <value>` | Slot field not two `:` separators |
| `cannot parse pci bus number: <value>` | Bus-number hex string can't be parsed |

Friendly-name lookup failures are non-fatal, printed directly to stderr:

| Full stderr output | Condition |
|---|---|
| `Error looking up friendly name: opening pci database: <pcidb error>` | `pcidb.New()` fails |

---

### 6.7 `pci/nvidia/`

All nvidia errors surface as `Warning: failed to get additional properties: NVIDIA: ...` via [§6.5](#65-pcipci_devicesgo).

| Error chain (after `NVIDIA: `) | Condition |
|---|---|
| `vRAM: nvidia-smi: <exec error>` | `nvidia-smi` not found / fails, no stdout output |
| `vRAM: nvidia-smi: <exec error>: <stdout>` | `nvidia-smi` fails with stdout error message |
| `vRAM: <strconv error>` | vRAM value string not a valid integer |
| `compute capability: nvidia-smi: <exec error>` | Same patterns for compute-capability query |
| `compute capability: nvidia-smi: <exec error>: <stdout>` | Same |

---

### 6.8 `pci/amd/`

All AMD errors surface as `Warning: failed to get additional properties: AMD: ...` via [§6.5](#65-pcipci_devicesgo).

| Error chain (after `AMD: `) | Condition |
|---|---|
| `vRAM: <os error>` | `/sys/bus/pci/devices/<slot>/mem_info_vram_total` unreadable |
| `vRAM: <strconv error>` | File content not a valid integer |
| `gfx architecture: <os error>` | `/sys/class/kfd/kfd/topology/nodes` unreadable |
| `gfx architecture: gfx_target_version not found for device with pci slot <slot>` | No KFD node matched the device |
| `gfx architecture: gfx_target_version is invalid for this device` | Value is `"0"` |
| `gfx architecture: gfx_target_version has an unexpected format: <value>` | Hex string shorter than 6 chars |
| `gfx architecture: parsing major version from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 0–1 |
| `gfx architecture: parsing minor version from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 2–3 |
| `gfx architecture: parsing revision from gfx_target_version: <strconv error>` | `strconv.Atoi` fails on chars 4–5 |
| `gfx architecture: unexpected format for gfx_target_version: <line>` | Line not 2 space-separated parts |

---

### 6.9 `pci/intel/`

All Intel errors surface as `Warning: failed to get additional properties: Intel: ...` via [§6.5](#65-pcipci_devicesgo).

| Error chain (after `Intel: `) | Condition |
|---|---|
| `vRAM: <exec error>` | `clinfo --json` binary missing or command fails |
| `vRAM: parsing clinfo json: <json error>` | `json.Unmarshal` on clinfo output fails |
| `vRAM: clinfo: no devices found` | `clinfo.Devices` is empty |
| `vRAM: clinfo: no online devices found` | `clinfo.Devices[0].Online` is empty |

*(When vRAM lookup succeeds but no device PCI address matches `device.Slot`, `nil` is returned without error.)*

---

## 7. `pkg/snap_store/snap_store.go`

All errors from this package surface as one of:
- `Warning: unable to get component sizes: components sizes: <error>` — `use-engine` with `--verbose` (to stderr)
- `Warning: unable to get component sizes` — `use-engine` without `--verbose`, error detail suppressed (to stderr)
- `Warning: unable to get component sizes: <error>` — `prune-cache` with `--verbose`, wraps one level less (to stdout)
- `Warning: unable to get component sizes` — `prune-cache` without `--verbose` (to stdout)

Full error chains from `ComponentSizes()`:

| Full error chain from `ComponentSizes` | Condition |
|---|---|
| `components sizes: SNAP_NAME is not set. Likely not inside a snap` | `$SNAP_NAME` env-var is empty |
| `components sizes: SNAP_REVISION is not set` | `$SNAP_REVISION` env-var is empty |
| `components sizes: not installed from store` | Revision starts with `"x"` (sideloaded) |
| `components sizes: parsing snap revision: <strconv error>` | `strconv.ParseInt` on revision string fails |
| `components sizes: getting snap info: creating new http request: <net error>` | `http.NewRequest` for snap info fails |
| `components sizes: getting snap info: HTTP request: <net error>` | Network error fetching snap info |
| `components sizes: getting snap info: HTTP status not OK: <code>` | Snap info response not 200 |
| `components sizes: getting snap info: json: <json error>` | JSON decode of snap info response fails |
| `components sizes: getting components: fetching refresh data from store: json: <json error>` | `json.Marshal` of refresh request body fails |
| `components sizes: getting components: fetching refresh data from store: new http request: <net error>` | `http.NewRequest` for refresh fails |
| `components sizes: getting components: fetching refresh data from store: HTTP request: <net error>` | Network error on refresh call |
| `components sizes: getting components: fetching refresh data from store: HTTP status not OK: <code>` | Refresh response not 200 |
| `components sizes: getting components: fetching refresh data from store: json: <json error>` | JSON decode of refresh response fails |
| `components sizes: getting components: store returned no refresh results` | `Results` slice is empty |
| `components sizes: getting components: no refresh results found for snap id <id>` | No result matches the snap ID |

---

## 8. `pkg/storage/`

### 8.1 `cache.go`

#### `SetActiveEngine(engine)` errors

| Error message fragment | Condition |
|---|---|
| `engine name cannot be empty` | `engine == ""` |
| `<snapctl.Set error>` | External snapctl error (unwrapped) |

#### `GetActiveEngine()` errors

| Error message fragment | Condition |
|---|---|
| `<snapctl.Get or json.Unmarshal error>` | Any non-ErrorNotFound storage error (unwrapped) |

---

### 8.2 `config.go`

#### `Config.Set(key, value, confType)` errors

| Error message fragment | Condition |
|---|---|
| `checking key: <snapctl error>` | `Config.Get(key)` fails during UserConfig validation |
| `unknown key` | Key not in schema (no existing value) |
| `<snapctl.Set error>` | `storage.Set` fails (unwrapped) |

#### `Config.SetDocument(key, value, confType)` errors

| Error message fragment | Condition |
|---|---|
| `<json.Marshal error>` | Value cannot be marshalled (unwrapped) |
| `<snapctl.Set error>` | `snapctl.Set(...).Document().Run()` fails (unwrapped) |

#### `Config.Get(key)` / `Config.GetAll()` errors

| Error message fragment | Condition |
|---|---|
| `<snapctl.Get or json.Unmarshal error>` | `loadConfigs` → `storage.Get` fails (unwrapped) |

#### `Config.Unset(key, confType)` errors

| Error message fragment | Condition |
|---|---|
| `<snapctl.Unset error>` | `snapctl.Unset(...).Run()` fails (unwrapped) |

---

### 8.3 `snapctl_storage.go`

The external boundary. All errors are raw errors from `github.com/canonical/go-snapctl` or stdlib:

| Method | External error source |
|---|---|
| `Set` | `snapctl.Set(...).Run()` |
| `SetDocument` | `json.Marshal` (stdlib) or `snapctl.Set(...).Document().Run()` |
| `Get` | `snapctl.Get(...).Run()` or `json.Unmarshal` (stdlib) |
| `Unset` | `snapctl.Unset(...).Run()` |

`ErrorNotFound` (`"not found"`) is returned (and immediately handled in `cache.go`) when `snapctl.Get` returns an empty string.

---

## 9. `pkg/utils/utils.go`

#### `StringToBytes(sizeString)` errors

| Error message fragment | Condition |
|---|---|
| `<strconv.ParseUint error>` | Numeric part of the size string is not a valid integer (unwrapped) |

#### `SubDirectories(dirPath)` errors

| Error message fragment | Condition |
|---|---|
| `failed to read directory: <os error>` | `os.ReadDir(dirPath)` fails |

---

## Cross-cutting Notes

### Non-fatal outputs (never cobra-prefixed, never returned as errors)

| Full output | Stream | Source | Condition |
|---|---|---|---|
| `Warning: unable to get component sizes: <error>` | stderr | `use-engine.go` | `snap_store.ComponentSizes()` fails with `--verbose` |
| `Warning: unable to get component sizes` | stderr | `use-engine.go` | Same, without `--verbose` |
| `Warning: unable to get component sizes: <error>` | stdout | `prune-cache.go` | `snap_store.ComponentSizes()` fails with `--verbose` |
| `Warning: unable to get component sizes` | stdout | `prune-cache.go` | Same, without `--verbose` |
| `Warning: previously active engine "<name>" not found; skipping user configuration cleanup.` | stderr | `engine.go` | `ErrManifestNotFound` during unset with `unsetUserOverrides=true` |
| `Warning: failed to get additional properties: <vendor>: <error>` | stderr | `pci_devices.go` | Per-device, non-fatal; all AMD/NVIDIA/Intel tool errors |
| `Error looking up friendly name: <error>` | stderr | `lspci.go` | Per-device, non-fatal PCI database lookup failure |
| `No engines found.` | stderr | `list-engines.go` | Engine list is empty |
| `Warning: "<key>" configuration field is deprecated!` | stderr | `get.go` | Deprecated key requested on a terminal |

### Sentinel errors (used with `errors.Is`)

| Sentinel | Package | Checked in |
|---|---|---|
| `engines.ErrManifestNotFound` | `pkg/engines` | `use-engine`, `prune-cache`, `run`, `engine.go` (`UnsetEngineConfig`) |
| `selector.ErrorNoCompatibleEngine` | `pkg/selector` | `use-engine` (`autoSelectEngine`), `debug/select` |
| `storage.ErrorNotFound` | `pkg/storage` | `cache.go` (`GetActiveEngine`) — handled internally, never surfaces to the user |
| `pci.ErrorVendorNotSupported` | `pkg/hardware_info/pci` | `pci_devices.go` (`addAdditionalProperties`) — silently ignored |

### Root-only commands (return `Error: permission denied, try again with sudo` immediately if not root)

`use-engine`, `set`, `prune-cache`
