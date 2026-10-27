# Plan: non-blocking component installation with progress reporting

Status: ready to implement (investigation complete, no code written yet)

## Repository context

- Repo: `inference-snaps-cli`, module `github.com/canonical/inference-snaps-cli/v2`, Go 1.26.1.
- The CLI binary is `cmd/modelctl`; shared logic lives in `cmd/modelctl/common`, reusable packages
  in `pkg/`.
- Snap packaging: `snap/snapcraft.yaml` (snap name `stack-utils`, `base: core24`, `grade: devel`).
- Build: `go build ./cmd/modelctl` — Test: `go test -count 1 -failfast ./...` — Snap: `snapcraft -v`.
- Related upstream issue this change addresses:
  [inference-snaps-cli#122](https://github.com/canonical/inference-snaps-cli/issues/122)
  (component install times out / blocks with no feedback). It is referenced in a comment in
  `cmd/modelctl/common/component.go` — remove that comment once fixed.

## Goal

Replace the blocking `snapctl install <snap>+<component>` call with the newer non-blocking flow:

```
snapctl install +<component> --no-wait   # prints change id on stdout
snapctl tasks --format=json <change-id>  # progress polling
```

Progress is printed as `label: done/total` for now. A real progress bar is out of scope.

`snapctl is-ready <change-id>` also exists, with exit codes:

```
0: change completed successfully (Done)
1: fatal errors (invalid change ID, permissions error, too many args)
2: change is ready but did not complete successfully (Undone, Error, Hold)
3: change is not ready
stdout: always empty, the exit code conveys readiness
stderr: empty for exit codes 0 and 3, contains relevant errors for exit codes 1 and 2
```

**Decided: we do not use `is-ready`.** The `tasks --format=json` output already carries `ready` and
`status`, so a single poll gives both completion and failure detection without a second exec.

Reference output:

```
root@host:~# snapctl install gemma4+llamacpp-cuda --no-wait
3998
root@host:~# snapctl tasks --format=json 3998
{"id":"3998","kind":"snapctl-install","summary":"Installing components [llamacpp-cuda] for snap gemma4","status":"Doing",
 "tasks":[{"id":"35298","kind":"download-component","summary":"Download component \"llamacpp-cuda\" (311)","status":"Doing",
 "progress":{"label":"gemma4+llamacpp-cuda","done":118582886,"total":795006091},
 "spawn-time":"2026-08-24T15:48:08.656741269+02:00","ready-time":"0001-01-01T00:00:00Z"}, ...],
 "ready":false,"spawn-time":"...","ready-time":"0001-01-01T00:00:00Z"}
root@host:~# snapctl is-ready 3998 ; echo $?
3
```

Once complete, the same change ID keeps working — `status`/task statuses become `Done`, `ready` is
`true`, `ready-time` is set, and the download task's `done` equals `total`:

```
root@host:~# snapctl is-ready 3998 ; echo $?
0
root@host:~# snapctl tasks --format=json 3998
{"id":"3998","kind":"snapctl-install","summary":"Installing components [llamacpp-cuda] for snap gemma4","status":"Done",
 "tasks":[{"id":"35298","kind":"download-component","summary":"Download component \"llamacpp-cuda\" (311)","status":"Done",
 "progress":{"label":"gemma4+llamacpp-cuda","done":795006091,"total":795006091},
 "spawn-time":"2026-08-24T15:48:08.656741269+02:00","ready-time":"2026-08-24T15:54:36.712923562+02:00"}, ...],
 "ready":true,"spawn-time":"...","ready-time":"2026-08-24T15:54:41.789954346+02:00"}
```

An unknown change ID fails on stderr with exit code 1:

```
root@host:~# snapctl tasks --format=json 88888
error: snapctl: change "88888" not found
```

## Findings

**`go-snapctl v1.0.0-beta.6` cannot do this.** Its `InstallComponents(...).Run()` just execs
`snapctl install +<comp>` and returns `CombinedOutput()` merged into an error. There is:

- no `--no-wait` option on the builder,
- no `tasks` / `is-ready` commands,
- no exported escape hatch (`run()` is unexported, output is combined so a change ID on stdout
  can't be isolated cleanly),
- `v1.0.0-beta.6` is the newest published version — no upgrade path.

So the non-blocking flow must be implemented in our own `pkg/snap` by exec'ing `snapctl` directly
(upstreaming to go-snapctl later is optional).

**Current call chain:** `InstallMissingComponents` → `InstallComponents(ctx, comps)` →
`installComponents(ctx, comps, 60m, 10s)` → per component: `StartProgressSpinner("Installing X")` +
blocking `ctx.Snap.InstallComponent(name)` with a retry loop matching 4 error substrings
(`already installed`, `unknown to the store`, `timeout exceeded while waiting for response`,
`change in progress`). Tests in `cmd/modelctl/common/component_test.go` (`TestInstallComponents`)
drive it via `snap.MockWithInstall(func(name string) error)`.

Environment confirms availability: local snapd is 2.77, `snapctl` at `/usr/bin/snapctl`; base is `core24`.

**Impact inventory** (verified by grepping the whole repo):

| symbol | where | impact |
| --- | --- | --- |
| `Snap.InstallComponent` | `pkg/snap/snap.go` (interface + impl), `pkg/snap/mock_snap.go` | signature changes |
| `ctx.Snap.InstallComponent(...)` | `cmd/modelctl/common/component.go` lines ~266, ~292, ~297 — **the only production call sites** | rewritten |
| `snap.MockWithInstall(...)` | `cmd/modelctl/common/component_test.go` only | tests updated |
| `snap.Mock()` | `cmd/modelctl/commands/{use-engine,unset,set}_test.go`, `cmd/modelctl/common/status_test.go` | **unaffected**, as long as `Mock()` keeps satisfying the interface |
| `snap.MockWithServiceStatuses(...)` | `cmd/modelctl/common/status_test.go` | unaffected |

`dev/smoke-test.sh` does not assert on component-install output, so no change needed there.

## Plan

**1. `pkg/snap/change.go` (new)** — JSON model + parsing for `snapctl tasks --format=json <id>`.
The payload is snapd's `client.Change` (see
[`client/change.go`](https://github.com/canonical/snapd/blob/master/client/change.go)), so mirror it:

```go
type Change   struct { ID, Kind, Summary, Status string; Tasks []Task; Ready bool; Err string; SpawnTime, ReadyTime time.Time }
type Task     struct { ID, Kind, Summary, Status string; Log []string; Progress Progress; SpawnTime, ReadyTime time.Time }
type Progress struct { Label string; Done, Total int64 }
```

- json tags: `id`, `kind`, `summary`, `status`, `tasks`, `ready`, `err`, `log`, `progress`,
  `spawn-time`, `ready-time`.
- `err` (change level) and `log` (task level) carry the failure detail — use them for error messages.
- Upstream tags these times `omitzero`, and observed output contains
  `"ready-time":"0001-01-01T00:00:00Z"`. Parsing into `time.Time` tolerates both present-zero and
  absent, so no custom unmarshalling is needed.

**2. `pkg/snap/snap.go`** — add a private runner that execs `snapctl` capturing **stdout and stderr
separately** (the change ID arrives on stdout; snapd's error text on stderr must stay in the
returned error, because `component.go` classifies errors by substring):

```go
func runSnapctl(args ...string) (stdout string, err error) // exec.Command("snapctl", args...)
```

Extend the `Snap` interface (keep the remaining methods untouched; `RemoveComponents`,
`Restart`, `Services` etc. continue to use `go-snapctl`):

- `InstallComponent(name string) (changeID string, err error)` → `snapctl install --no-wait +<name>`,
  returns trimmed stdout as change ID. **Breaking signature change** — the only production caller is
  `installComponents`. **An empty change ID with a nil error is a valid, successful result** meaning
  the component was already installed (see verified behaviour below); document this on the method.
- `Change(changeID string) (*Change, error)` → `snapctl tasks --format=json <id>`, unmarshalled via
  `ParseChange`.
- No `is-ready` wrapper (decided): `Change.Ready` and `Change.Status` cover it.

**3. `cmd/modelctl/common/component.go` → `installComponents`**

Split the current single loop into a submit phase and a poll phase, e.g.:

```go
func InstallComponents(ctx *Context, components []string) error {
    return installComponents(ctx, components, 60*time.Minute, 10*time.Second, 1*time.Second)
}
func installComponents(ctx *Context, components []string, installTimeout, retryDelay, pollInterval time.Duration) error
func submitComponentInstall(ctx *Context, component string, start time.Time, installTimeout, retryDelay time.Duration) (changeID string, err error)
func waitForChange(ctx *Context, component, changeID string, start time.Time, installTimeout, pollInterval time.Duration) error
```

- Phase A (submit): call `InstallComponent` → changeID.
  - **Empty change ID + nil error ⇒ the component was already installed.** Print
    `Installed <component>` and continue with the next component without polling. This replaces the
    old `already installed` error branch, which snapd no longer produces on this path.
  - Keep the remaining retry classification loop, using the error substrings already defined as
    constants at the top of `installComponents`: `cannot install components for a snap that is
    unknown to the store` → hard error with the existing "snap not known to the store" message;
    `timeout exceeded while waiting for response` and `change in progress` → sleep `retryDelay` and
    retry, honouring `installTimeout`; anything else → return immediately. Keep the
    `already installed` substring branch as a cheap defensive fallback.
- Phase B (poll): loop every `pollInterval` calling `ctx.Snap.Change(changeID)`:
  - pick the first task with status `Doing` and render it:
    - download tasks (non-empty `progress.label`, `total > 1`):
      `gemma4+llamacpp-cuda: 113M/758M` using `utils.FmtBytesShort`,
    - all other tasks (empty label, `0/1` progress): print the task summary, e.g.
      `Mount component "gemma4+llamacpp-cuda" (311)`,
  - tty (`utils.IsTerminalOutput()`): redraw in place with `\r`, padding to clear the previous,
    longer line; non-tty: print a new line **only when the rendered line changes**, to avoid log spam,
  - exit the loop when `change.Ready`; if the change is then in `Error`/`Undone`/`Hold`, return an
    error built from `change.Err`, falling back to the first failing task's summary and its `log`
    entries,
  - enforce the same overall `installTimeout`, keeping the existing user-facing message
    ("Monitor the installation progress with `snap changes`" / "Rerun this command once the
    installation is complete"),
  - keep the final `fmt.Println("Installed " + component)` per component.
- Spinner: `StartProgressSpinner` (`cmd/modelctl/common/spinner.go`) writes from its own goroutine,
  so it must be stopped before any `\r` progress line is printed. Use it only for the submit phase
  and while no task with real progress is available yet.
- Suggested new helper `cmd/modelctl/common/progress.go` holding the "single updating line" renderer
  (`Update(text string)` / `Clear()`), so it is unit-testable against a `bytes.Buffer` and reusable.

**4. `pkg/snap/mock_snap.go`** — update the mock to the new interface:

- `MockWithInstall(fn func(name string) error)` keeps its existing signature (returning a synthetic
  incrementing change ID) so the existing `component_test.go` error-classification subtests need
  minimal changes.
- Add a way to return an **empty** change ID, to cover the already-installed path.
- Add a way to drive `Change`, e.g. `MockWithChanges(changes ...Change)` returning the changes in
  order and repeating the last, plus `MockWithChangeError(err error)`.
- `Change` default (used by plain `Mock()`) must return a ready, successful change, otherwise every
  existing test using `snap.Mock()` would hang in the poll loop.

**5. Tests**

- `pkg/snap/change_test.go` (new): unmarshal both sample payloads above (in-progress and completed)
  and assert id/status/`ready`/task progress/`ready-time`; cover `Failed()`, `FailureReason()`
  (including `err` and task `log`), `CurrentTask()` (nil when everything is `Done`) and the
  download vs. non-download rendering of `Task.Description()`.
- `cmd/modelctl/common/component_test.go` (`TestInstallComponents`): adapt the existing subtests to
  the extra `pollInterval` argument, and add:
  - **empty change ID (already installed) → success, no polling, `Installed <component>` printed**,
  - progress polling over several `Doing` snapshots until `Done`,
  - a change that ends in `Error` → error mentions the failing task summary,
  - `Change` returning an error → surfaced, not silently retried forever,
  - polling exceeding `installTimeout` → the existing "timed out while installing" message.
- Progress renderer test against a `bytes.Buffer`: non-tty emits one line per *changed* text only.

**6. snapd compatibility — `assumes:` (decided)**

**Decided: no fallback code path.** Declare the requirement as a top-level key in
`snap/snapcraft.yaml` (alongside `base:`/`grade:`/`confinement:`, before the `slots:` block):

```yaml
assumes:
  - snapd2.77
```

so snapd refuses installation/refresh on older systems with a clear message, and the Go code keeps a
single clean path. Agreed with the maintainer: snapd 2.77 will be publicly released by the time this
change lands, so no transitional handling is needed.

Version research (snapd upstream, verified 2026-08-24):

| feature | PR | merged | first release |
| --- | --- | --- | --- |
| `snapctl install/remove --no-wait` | [#16852](https://github.com/canonical/snapd/pull/16852) | 2026-05-04 | flag is **commented out** in the 2.76.x branch (`install.go` line 47), enabled on master |
| `snapctl tasks` / `change` (+ `--format=json`) | [#16894](https://github.com/canonical/snapd/pull/16894) | 2026-06-03 | `overlord/hookstate/ctlcmd/tasks.go` absent from tag 2.76.3, present on master |
| `snapctl is-ready` exit codes | [#17114](https://github.com/canonical/snapd/pull/17114) | 2026-06-01 | master (`is_ready.go` exists in 2.76.3, exit codes fixed after) |

⇒ **snapd 2.77 is the minimum**. As of 2026-08-24 the latest *stable* snapd is 2.76.3 (released
2026-08-20) and 2.77 is only in edge/beta (the dev machine runs `2.77+g75.bee548f`). Per decision 4
below, 2.77 is expected to be publicly released before this change lands, so this is accepted.

Other details confirmed from upstream `tasks.go`:

- `tasks` and `change` are aliases for the same command; `--format` only accepts `json`.
- Exit codes: `0` = change info reported successfully regardless of change state, `1` = any error
  (invalid change ID, permissions). So JSON parse + `status`/`ready` is the only thing we need.
- JSON is written to stdout with `json.NewEncoder` (trailing newline).
- The change must belong to the calling snap (`getAssociatedChange`), which is exactly our case.

Note: this repo builds `stack-utils`, but `modelctl` is also staged into the consumer inference
snaps — the same `assumes:` entry must be added to those snaps' `snapcraft.yaml`, and mentioned in
the README/integration docs.

## Decisions (settled)

1. **`assumes: [snapd2.77]` in `snapcraft.yaml`, no old-snapd fallback.**
2. **No `is-ready` wrapper** — poll `tasks --format=json` and use `ready` + `status`.
   (For the record, `is-ready` exit codes are: `0` Done, `1` fatal error, `2` ready but not
   successful — Undone/Error/Hold, `3` not ready; stdout always empty.)
3. **Output format** — downloads: `label: done/total` via `utils.FmtBytesShort`; non-download tasks:
   task summary; non-tty: emit a line only when the rendered text changes.
4. **snapd 2.77 will be publicly released before this change lands**, so `assumes: snapd2.77` is
   acceptable and no transitional/fallback handling is required.

## Suggested implementation order

1. `pkg/snap/change.go` + `pkg/snap/change_test.go` (pure, no dependencies).
2. `pkg/snap/snap.go`: `runSnapctl`, new `InstallComponent`, new `Change`, updated interface.
3. `pkg/snap/mock_snap.go` to match the interface — at this point `go build ./...` must pass.
4. `cmd/modelctl/common/progress.go` (+ test).
5. Rewrite `installComponents` in `cmd/modelctl/common/component.go`.
6. Update/extend `cmd/modelctl/common/component_test.go`.
7. `snap/snapcraft.yaml`: add `assumes:`.
8. Docs: note the snapd 2.77 requirement in `README.md`.

## Verification

```bash
gofmt -l ./cmd ./pkg          # must print nothing
go build ./...
go test -count 1 -failfast ./...
```

Manual check on a machine with snapd ≥ 2.77 (unit tests cannot cover the real snapctl calls):

```bash
snapcraft -v
sudo snap install --dangerous ./stack-utils_*.snap
sudo snap connect stack-utils:hardware-observe
# trigger a component install that actually downloads, e.g.
stack-utils use-engine <engine-needing-an-uninstalled-component>
```

Confirm: a single updating progress line with sensible byte counts, correct behaviour when piping
to a file (`| cat` / redirect: one line per change, no `\r` spam), a clean error when the change
fails, and that Ctrl-C leaves the terminal in a sane state.

## Out of scope

- A graphical progress bar / percentage (explicitly deferred by the maintainer).
- Using `--no-wait` for `RemoveComponents` (same opportunity exists; separate change).
- Installing multiple components concurrently — snapd allows only one `snapctl-install` change per
  snap at a time (see the verified "change in progress" behaviour above), so the sequential loop
  stays.
- Upstreaming these commands to `go-snapctl`.

### Possible refinement (needs a decision before implementing)

`snapctl install` accepts several components in one invocation, so all missing components could be
submitted as a **single change** (`snapctl install --no-wait +a +b +c`) and tracked with one poll
loop. Benefits: one change instead of N, no chance of tripping over our own "change in progress",
and snapd downloads/installs them as one unit. Costs: the per-component
"Installed &lt;component&gt;" progression and the per-component "already installed" handling would
need reworking, since the error/retry classification then applies to the whole batch.
The plan above keeps the current **one component per change** behaviour, which is the smaller,
lower-risk change.

## Verified snapd behaviour (snapd 2.77, manually tested 2026-08-24)

- ⚠️ **"Already installed" is NOT an error under `--no-wait`.** Confirmed both by experiment and by
  reading [`overlord/hookstate/ctlcmd/install.go`](https://github.com/canonical/snapd/blob/master/overlord/hookstate/ctlcmd/install.go)
  on master:

  ```go
  id, affectedComponents, err := runSnapManagementCommand(...)
  if err != nil {
      if _, ok := err.(*snap.AlreadyInstalledError); !ok {
          return err          // AlreadyInstalledError is swallowed
      }
  }
  if len(affectedComponents) < len(comps) {
      // "snapctl: component %q is already installed" -> c.stderr, per unaffected component
  }
  if c.NoWait {
      fmt.Fprintf(c.stdout, "%s\n", id)   // id is EMPTY when nothing had to be done
  }
  return nil                              // exit code 0
  ```

  ```
  root@host:~# snapctl install gemma4+llamacpp --no-wait
                                             <- empty change ID on stdout
  snapctl: component "llamacpp" is already installed   <- stderr, exit code 0
  ```

  **Consequence for the design:** the existing `already installed` error-substring check never fires
  on this path. `InstallComponent` must return `("", nil)` and the caller must treat an **empty
  change ID as "nothing to do"** — print `Installed <component>` and skip polling. Do *not* treat an
  empty change ID as a failure. Keep the substring check as a cheap defensive fallback.
  (Note this also means the already-installed check wins over "change in progress": running
  `install +llamacpp` while another change was running still returned "already installed".)
- ✅ "Change in progress" **is** a real error, reported synchronously by the submit call while
  another snapctl change for the same snap is running:

  ```
  root@host:~# snapctl install gemma4+llamacpp-rocm --no-wait
  4000
  root@host:~# snapctl install gemma4+llamacpp-rocm --no-wait
  error: snapctl: snap "gemma4" has "snapctl-install" change in progress
  ```

  The existing `change in progress` substring check therefore still works, and retrying with
  `retryDelay` remains the correct handling. Note the poll phase already waits for each component's
  change to finish before the next component is submitted, so this should now only be triggered by
  changes started outside our process.
- ✅ A completed change stays queryable: `snapctl tasks --format=json <id>` keeps returning the full
  task list with `"status":"Done"`, `"ready":true` and `ready-time` set, so the poll loop can rely on
  a final successful poll and does not need to treat a missing change as success.
- ✅ An unknown/expired change ID is a hard error (`error: snapctl: change "88888" not found`,
  exit 1) — surface it as an error, no special-casing.
- ✅ On completion the download task reports `done == total`, so the last rendered progress line is
  naturally 100%.

Still unverified (low risk): whether `timeout exceeded while waiting for response` can still be
returned by the submit call — with `--no-wait` returning immediately it should no longer occur, but
the retry handling is kept as a safety net.

## Definition of done

- [ ] `snapctl install --no-wait` + `snapctl tasks --format=json` polling replaces the blocking call.
- [ ] Progress is visible and updates while a component downloads, on tty and in logs.
- [ ] Existing error classification and user-facing messages are preserved.
- [ ] `assumes: snapd2.77` added; README mentions the requirement.
- [ ] `go build ./...` and `go test -count 1 -failfast ./...` pass; `gofmt -l` clean.
- [ ] Manually verified against a real component download on snapd ≥ 2.77.
- [ ] The `// This is blocking, but there is a timeout bug: .../issues/122` comment in
      `component.go` is removed and the issue referenced in the PR.
- [ ] This plan file is deleted (or moved into the PR description) before merging.

