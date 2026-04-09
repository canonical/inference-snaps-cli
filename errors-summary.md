# Errors Summary

All possible error messages that can be printed to the user, with full error chains exactly as they will appear.

## How cobra formats errors

`SilenceUsage: true` is set on the root command; `SilenceErrors` is **not** set. This means:

- **All errors returned from `RunE`** are printed by cobra to **stderr** as:
  ```
  Error: <message>
  ```
- **Argument/flag validation errors** (e.g. wrong number of args, unknown flags — handled before `RunE` is called) are also printed as `Error: <message>` and are additionally followed by:
  ```
  Run '<command> --help' for usage.
  ```
- Errors printed by application code directly via `fmt.Printf`, `fmt.Fprintf(os.Stderr, …)` etc. are **not** prefixed by cobra and appear exactly as written.

All `RunE`-sourced errors in the tables below therefore appear on screen as `Error: <the text shown>`.

---

## `main.go` — startup / setup

These errors are printed directly with `fmt.Printf` — **no** cobra `Error:` prefix.

| Full text printed to user |
|---|
| `Error: retrieving snap services: <snapctl error>` |
| `Error: group with ID "basic" not found` |

---

## `cmd/cli/commands/chat.go` — `chat` command

All rows are `RunE` errors → printed as `Error: <message>` by cobra.

| Full error printed to user |
|---|
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: looking up active engine: <snapctl error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: no active engine` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: engine manifest not found: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: SNAP_COMPONENTS env var not set` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: reading <componentYamlFile>: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: unmarshaling <componentYamlFile>: <yaml error>` |
| `Error: getting OpenAI base URL: getting server endpoints: OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` |
| `Error: getting OpenAI base URL: getting server endpoints: getting server HTTP URL: getting config "http.port": <snapctl error>` |
| `Error: getting OpenAI base URL: getting server endpoints: unsupported protocol "<proto>" for server "<name>" in component "<name>"` |
| `Error: getting OpenAI base URL: "openai" not found in server endpoints` |
| `Error: getting services: <snapctl error>` |
| `Error: server not active`<br><br>`Run "sudo snap start <snap>.server" to start the server.` |
| `Error: initializing readline: <readline error>` |
| `Error: invalid base URL: <url parse error>` |
| `Error: connection refused`<br><br>`Try again when the server is ready.`<br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: <net dial error>` |
| `Error: no models available on server`<br><br>`Try again when the server is ready.`<br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: api: <openai api error>` |
| `Error: <other openai/network error>`<br><br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: server returned no models`<br><br>`Try again when the server is ready.`<br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: expected one but server returned multiple models: <name1>, <name2>, ...` |
| `Error: connection refused`<br><br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: connection closed by server`<br><br>`Run "snap logs <snap>.server" to see the server logs.` |
| `Error: <stream error>`<br><br>`Run "snap logs <snap>.server" to see the server logs.` |

---

## `cmd/cli/commands/get.go` — `get` command

`RunE` errors → `Error: <message>`. The deprecation notice is printed directly to stderr — no cobra prefix.

| Full text printed to user |
|---|
| `Error: getting value of "<key>": <snapctl error>` |
| `Error: no value set for key "<key>"` |
| `Error: serializing value: <yaml marshal error>` |
| `Error: getting values: <snapctl error>` |
| `Error: serializing values: <yaml marshal error>` |
| *(stderr, non-fatal, no cobra prefix)* `Note: "<key>" configuration field is deprecated!` |

---

## `cmd/cli/commands/set.go` — `set` command

`set` uses `cobra.ExactArgs(1)`. Passing the wrong number of args triggers cobra's argument validator before `RunE`, producing:
```
Error: accepts 1 arg(s), received <n>
Run '<snap> set --help' for usage.
```
All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| *(cobra arg validation)* `Error: accepts 1 arg(s), received <n>`<br>`Run '<snap> set --help' for usage.` |
| `Error: permission denied, try again with sudo` |
| `Error: key must not start with an equal sign` |
| `Error: expected key=value, got "<arg>"` |
| `Error: "<key>" is read-only` |
| `Error: setting "<key>" to "<value>": checking existing keys: <snapctl error>` |
| `Error: setting "<key>" to "<value>": unknown key` |
| `Error: setting "<key>" to "<value>": <snapctl error>` |

---

## `cmd/cli/commands/list-engines.go` — `list-engines` command

`list-engines` uses `cobra.NoArgs`. Passing any argument triggers:
```
Error: unknown command "<arg>" for "<snap> list-engines"
Run '<snap> list-engines --help' for usage.
```
All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: checking engines: loading engines: <enginesDir>: <os error>` |
| `Error: checking engines: loading engines: <filepath>/engine.yaml: <os error>` |
| `Error: checking engines: loading engines: <enginesDir>: <yaml error>` |
| `Error: checking engines: getting machine info: getting memory info: querying host meminfo: reading /proc/meminfo: <os error>` |
| `Error: checking engines: getting machine info: getting memory info: parsing meminfo: parsing MemTotal: parsing kB value: <strconv error>` |
| `Error: checking engines: getting machine info: getting memory info: parsing meminfo: parsing MemTotal: parsing byte value: <strconv error>` |
| `Error: checking engines: getting machine info: getting memory info: parsing meminfo: parsing SwapTotal: parsing kB value: <strconv error>` |
| `Error: checking engines: getting machine info: getting memory info: parsing meminfo: parsing SwapTotal: parsing byte value: <strconv error>` |
| `Error: checking engines: getting machine info: getting cpu info: querying host cpuinfo: reading /proc/cpuinfo: <os error>` |
| `Error: checking engines: getting machine info: getting cpu info: getting host uname: <uname exec error>` |
| `Error: checking engines: getting machine info: getting cpu info: parsing cpu data: parsing cpuinfo: unsupported architecture: <arch>` |
| `Error: checking engines: getting machine info: getting cpu info: parsing cpu data: parsing cpuinfo: amd64: <strconv error>` |
| `Error: checking engines: getting machine info: getting cpu info: parsing cpu data: parsing cpuinfo: arm64: <strconv error>` |
| `Error: checking engines: getting machine info: getting cpu info: parsing cpu data: filtering cpu info: converting cpu info: unsupported architecture: <arch>` |
| `Error: checking engines: getting machine info: getting disk info: getting directory info for /var/lib/snapd/snaps: statfs: <syscall error>` |
| `Error: checking engines: getting machine info: getting pci devices: querying host lspci: <lspci exec error>` |
| `Error: checking engines: getting machine info: getting pci devices: parsing lspci response: unexpected format for pci slot: <slot>` |
| `Error: checking engines: getting machine info: getting pci devices: parsing lspci response: cannot parse pci bus number: <part>` |
| `Error: checking engines: scoring engines: parsing required memory: <strconv error>` |
| `Error: checking engines: scoring engines: total memory not reported by host system` |
| `Error: checking engines: scoring engines: parsing required disk space: <strconv error>` |
| `Error: checking engines: scoring engines: disk space not reported by host system` |
| `Error: looking up active engine: <snapctl error>` |
| `Error: printing table: adding data: <tablewriter error>` |
| `Error: printing table: rendering: <tablewriter error>` |
| `Error: printing json: marshalling engines: <json marshal error>` |
| `Error: unknown format "<format>"` |
| *(stderr, non-fatal, no cobra prefix)* `No engines found.` |

---

## `cmd/cli/commands/show-engine.go` — `show-engine` command

`show-engine` uses `cobra.MaximumNArgs(1)`. All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: invalid number of arguments` |
| `Error: looking up active engine: <snapctl error>` |
| `Error: no active engine` |
| `Error: checking engines: …` *(same tree as list-engines above)* |
| `Error: engine "<name>" does not exist` |
| `Error: printing engine manifest: json: <json marshal error>` |
| `Error: printing engine manifest: yaml: <yaml marshal error>` |
| `Error: printing engine manifest: unknown format "<format>"` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error: loading engines: <enginesDir>: <os error>` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error: loading engines: <filepath>/engine.yaml: <os error>` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error: loading engines: <enginesDir>: <yaml error>` |

---

## `cmd/cli/commands/show-machine.go` — `show-machine` command

`show-machine` uses `cobra.NoArgs`. All `RunE` errors → `Error: <message>`. The `fmt.Fprintln(os.Stderr, …)` messages are printed directly by application code — no cobra prefix.

| Full text printed to user |
|---|
| `Error: getting machine info: getting memory info: …` *(same subtree as list-engines above)* |
| `Error: getting machine info: getting cpu info: …` *(same subtree as list-engines above)* |
| `Error: getting machine info: getting disk info: getting directory info for /var/lib/snapd/snaps: statfs: <syscall error>` |
| `Error: getting machine info: getting pci devices: querying host lspci: <lspci exec error>` |
| `Error: getting machine info: getting pci devices: parsing lspci response: unexpected format for pci slot: <slot>` |
| `Error: getting machine info: getting pci devices: parsing lspci response: cannot parse pci bus number: <part>` |
| `Error: json: <json marshal error>` |
| `Error: yaml: <yaml marshal error>` |
| `Error: unknown format "<format>"` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get friendly name for pci device: opening pci database: <pcidb error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: AMD: getting gpu properties: looking up vram: <os error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: AMD: getting gpu properties: looking up gfx architecture: <os error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: AMD: getting gpu properties: looking up gfx architecture: gfx_target_version not found for device with pci slot <slot>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: NVIDIA: getting gpu properties: looking up vram: querying nvidia-smi: <exec error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: NVIDIA: getting gpu properties: looking up vram: querying nvidia-smi: <exec error>: <nvidia-smi stdout>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: NVIDIA: getting gpu properties: looking up vram: <strconv error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: NVIDIA: getting gpu properties: looking up compute capability: querying nvidia-smi: <exec error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: Intel: getting gpu properties: looking up vram: querying clinfo: <exec error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: Intel: getting gpu properties: looking up vram: parsing clinfo response: <json unmarshal error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: Intel: getting gpu properties: looking up vram: clinfo: no devices found` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to get additional properties for pci device: Intel: getting gpu properties: looking up vram: clinfo: no online devices found` |

---

## `cmd/cli/commands/status.go` — `status` command

`status` uses `cobra.NoArgs`. All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: waiting for component: looking up active engine: <snapctl error>` |
| `Error: waiting for component: no active engine` |
| `Error: waiting for component: loading engine manifest: engine manifest not found: <os error>` |
| `Error: waiting for component: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: waiting for component: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: waiting for component: SNAP_COMPONENTS env var not set` |
| `Error: waiting for component: timeout after waiting 3600s for required components: <comp1>, <comp2>, ...` |
| `Error: getting json status: getting status: looking up active engine: <snapctl error>` |
| `Error: getting json status: getting status: no active engine` |
| `Error: getting json status: getting status: getting list of services: <snapctl error>` |
| `Error: getting json status: getting status: unexpected service name format: "<name>"` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: looking up active engine: <snapctl error>` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: no active engine` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: loading engine manifest: engine manifest not found: <os error>` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: SNAP_COMPONENTS env var not set` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: reading <componentYamlFile>: <os error>` |
| `Error: getting json status: getting status: getting server api endpoints: loading engine environment: unmarshaling <componentYamlFile>: <yaml error>` |
| `Error: getting json status: getting status: getting server api endpoints: OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` |
| `Error: getting json status: getting status: getting server api endpoints: getting server HTTP URL: getting config "http.port": <snapctl error>` |
| `Error: getting json status: getting status: getting server api endpoints: unsupported protocol "<proto>" for server "<name>" in component "<name>"` |
| `Error: getting json status: marshalling json: <json marshal error>` |
| `Error: getting yaml status: getting status: …` *(same subtree as getting json status above)* |
| `Error: getting yaml status: marshalling yaml: <yaml marshal error>` |
| `Error: unknown format "<format>"` |

---

## `cmd/cli/commands/use-engine.go` — `use-engine` command

`use-engine` uses `cobra.MaximumNArgs(1)`. All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: permission denied, try again with sudo` |
| `Error: cannot specify both engine name and --auto flag` |
| `Error: cannot specify both engine name and --fix flag` |
| `Error: engine name not specified` |
| **`--auto` path:** |
| `Error: checking engines: …` *(same tree as list-engines above)* |
| `Error: finding top engine: no compatible engines found` |
| `Error: use engine: "<engine>" not found` |
| `Error: use engine: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: use engine: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: use engine: installing missing components: checking installed components: checking component directory "<comp>": <os error>` |
| `Error: use engine: installing missing components: checking installed components: component "<comp>" exists but is not a directory` |
| `Error: use engine: installing missing components: installing components: timed out while installing "<comp>":`<br>`Monitor the installation progress with "snap changes"`<br><br>`Rerun this command once the installation is complete` |
| `Error: use engine: installing missing components: installing components: snap not known to the store:`<br>`Rerun this command after manually installing "<comp>"` |
| `Error: use engine: installing missing components: installing components: installing "<comp>": <snapctl error>` |
| `Error: use engine: looking up active engine: <snapctl error>` |
| `Error: use engine: un-setting engine configurations: un-setting engine configurations: <snapctl error>` |
| `Error: use engine: un-setting engine configurations: loading engine manifest: engine manifest not found: <os error>` |
| `Error: use engine: un-setting engine configurations: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: use engine: un-setting engine configurations: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: use engine: un-setting engine configurations: un-setting configuration "<key>": <snapctl error>` |
| `Error: use engine: setting active engine: engine name cannot be empty` |
| `Error: use engine: setting active engine: <snapctl error>` |
| `Error: use engine: setting new engine configurations: setting engine configuration "<key>": <snapctl/json marshal error>` |
| **`--fix` path:** |
| `Error: looking up active engine: <snapctl error>` |
| `Error: no active engine` |
| `Error: loading active engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: loading active engine manifest: <enginesDir>: <yaml error>` |
| `Error: installing missing components: …` *(same subtree as --auto above)* |
| `Error: un-setting engine configurations: …` *(same subtree as --auto above)* |
| `Error: setting engine configurations: setting engine configuration "<key>": <snapctl/json marshal error>` |
| **Named engine path:** |
| `Error: "<engine>" not found` |
| `Error: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: installing missing components: …` *(same subtree as --auto above)* |
| `Error: looking up active engine: <snapctl error>` |
| `Error: un-setting engine configurations: …` *(same subtree as --auto above)* |
| `Error: setting active engine: engine name cannot be empty` |
| `Error: setting active engine: <snapctl error>` |
| `Error: setting new engine configurations: setting engine configuration "<key>": <snapctl/json marshal error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: SNAP_NAME is not set. Likely not inside a snap` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: SNAP_REVISION is not set` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: not installed from store` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: parsing snap revision: <strconv error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting snap info: creating http request: <http error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting snap info: making http request: <http error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting snap info: http status not OK: <status code>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting snap info: decoding json: <json error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: fetching refresh data from store: marshalling request: <json error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: fetching refresh data from store: creating http request: <http error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: fetching refresh data from store: making http request: <http error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: fetching refresh data from store: http status not OK: <status code>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: fetching refresh data from store: decoding json: <json error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: store returned no refresh results` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: finding snap components: getting components: no refresh results found for snap id <id>` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error loading engines: <enginesDir>: <os error>` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error loading engines: <filepath>/engine.yaml: <os error>` |
| *(tab-completion, stdout, direct `fmt.Printf`, no cobra prefix)* `Error loading engines: <enginesDir>: <yaml error>` |
| *(stderr, non-fatal, no cobra prefix)* `Warning: previously active engine "<name>" not found; skipping user configuration cleanup.` |

---

## `cmd/cli/commands/prune-cache.go` — `prune-cache` command

All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: permission denied, try again with sudo` |
| `Error: looking up active engine: <snapctl error>` |
| `Error: no active engine` |
| `Error: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: cannot prune the active engine "<engine>"` |
| `Error: "<engine>" not found` |
| `Error: loading manifests: <enginesDir>: <os error>` |
| `Error: loading manifests: <filepath>/engine.yaml: <os error>` |
| `Error: loading manifests: <enginesDir>: <yaml error>` |
| `Error: checking component directory "<comp>": <os error>` |
| `Error: component "<comp>" exists but is not a directory` |
| `Error: un-setting engine configurations: <snapctl error>` |
| `Error: un-setting engine configurations: loading engine manifest: engine manifest not found: <os error>` |
| `Error: un-setting engine configurations: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: un-setting engine configurations: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: un-setting engine configurations: un-setting configuration "<key>": <snapctl error>` |
| `Error: removing components: <snapctl error>` |
| `Error: getting list of inactive engines: looking up active engine: <snapctl error>` |
| `Error: getting list of inactive engines: <enginesDir>: <os error>` |
| `Error: getting list of inactive engines: <filepath>/engine.yaml: <os error>` |
| `Error: getting list of inactive engines: <enginesDir>: <yaml error>` |
| `Error: confirming component: <snapctl error>` |
| *(stdout, non-fatal, no cobra prefix)* `Warning: unable to query component sizes: …` *(same subtree as use-engine above)* |
| *(stderr, non-fatal, no cobra prefix)* `Warning: previously active engine "<name>" not found; skipping user configuration cleanup.` |

---

## `cmd/cli/commands/serve-ui.go` — `serve-ui` command (hidden)

`serve-ui` uses `cobra.ExactArgs(1)`. All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: looking up active engine: <snapctl error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: no active engine` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: engine manifest not found: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: SNAP_COMPONENTS env var not set` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: reading <componentYamlFile>: <os error>` |
| `Error: getting OpenAI base URL: getting server endpoints: loading engine environment: unmarshaling <componentYamlFile>: <yaml error>` |
| `Error: getting OpenAI base URL: getting server endpoints: OPENAI_BASE_PATH env in component "<name>" is deprecated; set server settings in "servers".` |
| `Error: getting OpenAI base URL: getting server endpoints: getting server HTTP URL: getting config "http.port": <snapctl error>` |
| `Error: getting OpenAI base URL: getting server endpoints: unsupported protocol "<proto>" for server "<name>" in component "<name>"` |
| `Error: getting OpenAI base URL: "openai" not found in server endpoints` |
| `Error: getting active engine: <snapctl error>` |
| `Error: no engine is active` |
| `Error: invalid configuration: invalid OpenAI base URL: <url parse error>` |
| `Error: invalid configuration: unsupported capability: "<cap>"` |
| *(from `ui.Serve` → `http.ListenAndServe`)* `Error: <net listen error>` |
| *(from `ui.Serve` → `verifyStaticContent`)* `Error: unexpected static files: checking "<indexFile>": <os error>` |

---

## `cmd/cli/commands/run.go` — `run` command (hidden)

`run` uses `cobra.MaximumNArgs(1)`. All `RunE` errors → `Error: <message>`. The subprocess error is returned directly from `execCmd.Run()` and passed through `RunE` unchanged, so it also gets the `Error:` prefix.

| Full text printed to user |
|---|
| `Error: unexpected number of arguments, expected 1 got <n>` |
| `Error: waiting for component: looking up active engine: <snapctl error>` |
| `Error: waiting for component: no active engine` |
| `Error: waiting for component: loading engine manifest: engine manifest not found: <os error>` |
| `Error: waiting for component: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: waiting for component: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: waiting for component: SNAP_COMPONENTS env var not set` |
| `Error: waiting for component: timeout after waiting 3600s for required components: <comp1>, <comp2>, ...` |
| `Error: loading engine environment: loading engine component settings: looking up active engine: <snapctl error>` |
| `Error: loading engine environment: loading engine component settings: no active engine` |
| `Error: loading engine environment: loading engine component settings: loading engine manifest: engine manifest not found: <os error>` |
| `Error: loading engine environment: loading engine component settings: loading engine manifest: <filepath>/engine.yaml: <os error>` |
| `Error: loading engine environment: loading engine component settings: loading engine manifest: <enginesDir>: <yaml error>` |
| `Error: loading engine environment: loading engine component settings: SNAP_COMPONENTS env var not set` |
| `Error: loading engine environment: loading engine component settings: reading <componentYamlFile>: <os error>` |
| `Error: loading engine environment: loading engine component settings: unmarshaling <componentYamlFile>: <yaml error>` |
| `Error: loading engine environment: SNAP_COMPONENTS env var not set` |
| `Error: loading engine environment: invalid env var "<kv>"` |
| `Error: loading engine environment: setting env var "COMPONENT": <os error>` |
| `Error: loading engine environment: unsetting env var "COMPONENT": <os error>` |
| `Error: loading engine environment: setting env var "<key>": <os error>` |
| `Error: <subprocess exec error>` *(passed through directly from `execCmd.Run()`)* |

---

## `cmd/cli/commands/version.go` — `version` command

All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: marshalling json: <json marshal error>` |
| `Error: marshalling yaml: <yaml marshal error>` |
| `Error: unknown format "<format>"` |

---

## `cmd/cli/commands/debug/chat.go` — `debug chat` command (hidden)

All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: the --base-url parameter is required` |
| *(then any error from `chatClient.Start` — see `chat` command section above, same `Error:` prefix applies)* |

---

## `cmd/cli/commands/debug/validate.go` — `debug validate-engines` command (hidden)

`validate-engines` uses `cobra.MinimumNArgs(1)`. Passing no args triggers cobra's validator:
```
Error: requires at least 1 arg(s), only received 0
Run '<snap> debug validate-engines --help' for usage.
```
The per-file `❌` lines are printed directly with `fmt.Printf` — **no** cobra prefix. The final summary error is a `RunE` return → `Error: <message>`.

| Full text printed to user |
|---|
| *(cobra arg validation)* `Error: requires at least 1 arg(s), only received 0`<br>`Run '<snap> debug validate-engines --help' for usage.` |
| `Error: not all manifests are valid` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: manifest file must be called engine.yaml: <path>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: manifest file does not exist: <path>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: getting file info: <os error>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: reading file: <os error>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: empty yaml data` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: decoding manifest: <yaml error>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: required field is not set: name` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: engine directory name should match name in manifest: <dir> != <manifest-name>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: required field is not set: description` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: required field is not set: vendor` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: required field is not set: grade` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: grade should be 'stable' or 'devel'` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: parsing memory: <strconv error>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: parsing disk space: <strconv error>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: configuration field <key> is not a primitive value: <value>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: cpu: architecture field required` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: cpu: invalid architecture: <arch>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: cpu: amd64: invalid field: <fieldname>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: cpu: arm64: invalid field: <fieldname>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: gpu: bus: invalid bus: <bus>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: gpu: bus: usb: device validation not implemented` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: gpu: bus: pci: invalid field: <fieldname>` |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: npu: bus: …` *(same as gpu subtree)* |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: allof <i>/<n>: typeless: bus: …` *(same as gpu subtree)* |
| *(stdout, direct `fmt.Printf`, no cobra prefix)* `❌ <path>: invalid device: anyof <i>/<n>: …` *(same subtree as allof above)* |

---

## `cmd/cli/commands/debug/select.go` — `debug select-engine` command (hidden)

All `RunE` errors → `Error: <message>`.

| Full text printed to user |
|---|
| `Error: decoding hardware info: <yaml decode error>` |
| `Error: loading engines from directory: <enginesDir>: <os error>` |
| `Error: loading engines from directory: <filepath>/engine.yaml: <os error>` |
| `Error: loading engines from directory: <enginesDir>: <yaml error>` |
| `Error: checking engines: parsing required memory: <strconv error>` |
| `Error: checking engines: total memory not reported by host system` |
| `Error: checking engines: parsing required disk space: <strconv error>` |
| `Error: checking engines: disk space not reported by host system` |
| `Error: finding top engine: no compatible engines found` |
| `Error: marshalling json: <json marshal error>` |
| `Error: marshalling yaml: <yaml marshal error>` |
| `Error: unknown format "<format>"` |

---

## Standalone sentinel / leaf errors (building blocks in chains above)

| Sentinel | Declared in |
|---|---|
| `permission denied, try again with sudo` | `cmd/cli/common/errors.go` |
| `looking up active engine` | `cmd/cli/common/errors.go` |
| `no active engine` | `cmd/cli/common/errors.go` |
| `not found` | `pkg/storage/storage.go` |
| `engine manifest not found` (wrapped with underlying OS error) | `pkg/engines/load.go` |
| `no compatible engines found` | `pkg/selector/select_stack.go` |
| `unknown key` | `pkg/storage/config.go` |

---

## Notes

- `<snapctl error>` is any error returned by the `go-snapctl` library (e.g. snapd not running, permission issues).
- `<os error>`, `<yaml error>`, `<json marshal error>`, `<strconv error>`, `<http error>`, `<exec error>` are standard library errors whose text is set by the Go standard library and the underlying OS.
- The `show-machine` command calls `hardware_info.Get(true)` (with `friendlyNames=true`) so it additionally triggers the PCI database and vendor-specific tool calls that `list-engines` / `use-engine` do not (those use `friendlyNames=false`).
- Additional-properties errors from AMD/NVIDIA/Intel sub-packages are silently dropped unless they pass through `addAdditionalProperties → deviceAdditionalProperties`, at which point they are written to stderr as `Warning: unable to get additional properties for pci device: …` without the cobra `Error:` prefix.
