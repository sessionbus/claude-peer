# Stable functionality checklist — Claude separation

Baseline: original peers main `710e5d33369cba4fb9468cd24fea0fe844a0219d`.
The F01–F20 requirement IDs are shared with the migration checklist. Product
semantics and known limitations stay explicit; a failed check does not remove a
requirement. Extraction source preservation and fresh installed behavior are
separate evidence. Current status: local source checks pass; independent review,
permanent installation and fresh installed acceptance are pending.

| ID | Preserved functionality | Existing regression coverage | Installed evidence / limit |
|---|---|---|---|
| F01 | Complete archive/install/update/remove, native plugin/hooks/skill and private aliases | interactive/install_test.go; claude/install; release download tests | Install and reinstall the exact archive in the real home; only owned paths removed |
| F02 | One public binary, wrapper/native version and token-selected lane dispatch | cmd/claude-peer/main_test.go; version_test.go; common peerversion | Exact installed source/version and aliases |
| F03 | Native argv ownership/order, literal --, groups/name, aliases, resume and yolo | interactive/launch_test.go, alias_values_test.go, permissions_test.go; claude_test.go | Preserve native selection and permission semantics, including positional and variadic values |
| F04 | Default native policy plus scoped managed communication grant; reject defeating deny rules | interactive/permissions_test.go; claude_test.go | Normal-policy native communication; do not broaden unrelated grants |
| F05 | Ordinary native launch remains ordinary | interactive/launch_test.go; per-launch plugin contract | No global plugin/MCP registration introduced |
| F06 | Native identity/title, initial publication, rename/withdrawal and resume | interactive/startup_test.go, owner_test.go; lane_endpoint_test.go | Exact native/public identity join; preserve native history |
| F07 | Discovery, single/multiple/group sends, public schema/error semantics | interactive/tools_test.go; common MCP tests and public SDK | Native list/send with authenticated receiver correlation |
| F08 | Zero-input lane open/readiness, spawn/resume | claude_test.go; lane_endpoint_test.go | Ready only on actual native confirmation; no model turn before run |
| F09 | Run/start/status/wait/ack, result cursor, exact outcome and interruption | stream_test.go; claude_test.go; common lane tests | Native result retained and collected before acknowledgement |
| F10 | Parent lifecycle, completion notification and direct-child tracing | interactive/tools_test.go; common trace and lane policy tests | Authority and scheduling remain daemon-owned; preserve prior qualifications |
| F11 | Interactive idle inbound autonomously wakes | interactive/delivery_test.go, owner_test.go | One inbound -> exact native reply/final with no later harness input |
| F12 | Interactive active delivery and processing | interactive/delivery_test.go, startup_test.go | Original active turn/process evidence -> actual receipt -> reply/final |
| F13 | Managed idle inbound autonomously wakes | wake_worker_test.go; wake_test.go; stream_test.go | One inbound -> actual automatic native turn and managed result |
| F14 | Managed active admission, queued/submitted distinction and original turn guard | wake_test.go; wake_worker_test.go; stream_test.go | Follow Claude's actual admission semantics; never substitute another product's contract |
| F15 | Resident reconnect, latest identity, no replay or worker resurrection | interactive/owner_test.go, transport_test.go | Preserve source-tested reconnect and no-replay behavior; report native limits |
| F16 | Cancellation, overlapping reports, bounds and protocol fidelity | interactive/mcp_test.go, transport_test.go; common tests | Controlled race evidence is distinct from native observations |
| F17 | Normal/failed startup, close, native death and owned cleanup | claude_test.go; interactive/startup_test.go, owner_test.go | Owned rows/processes absent; generic forced-death descendant containment remains qualified |
| F18 | Native history, persistence and independent auto-close policy | claude_test.go; common lane tests; retained shared-lifecycle notes | Close retains native history; no local policy scheduler added |
| F19 | Independent module/archive/CI/release and platform builds | architecture_test.go; interactive/install_test.go; release tests | Exact extracted archive and source; Linux/macOS amd64/arm64 builds |
| F20 | Every original product/common runtime, asset, fixture and test preserved | PRESERVED-FILES.json and baseline | 64 protected files and 133 original test functions resolve locally or in exact peer-common |

Coverage paths without a prefix are under `wrappers/claude/`. Shared support is
pinned to peer-common `eb655f686e4456a4c1121054763318e3d27e89b0`; its tests run in
that repository. All Claude-specific tests remain here.

Native clients are expected to update routinely. No exact native version
allowlist is introduced. Record the actual native version/hash for each test;
check concrete capabilities and behavior rather than rejecting a new version.
The build dependency pin on peer-common is separate from native client support.

## Evidence rules and retained limitations

- Never relabel historical acceptance as acceptance of the extracted binary.
- Use the permanent real-home/config/PATH installation for both interactive and
  lane tests. Install/reinstall and test that same integration.
- Each fresh cell gets one send, no replay, and no post-inbound harness prompt,
  PTY input or lifecycle turn. Preserve first failures and actual receipts.
- Presence `running:false` alone is not native-idle evidence. `no_receipt` is
  uncertain and never permission to resend.
- The historical Node implementation and held lane skills are retained under
  `docs/designs/claude-0.5.0` for knowledge preservation, not shipped or activated.
- Native hooks, SessionStart/first-report behavior, title confirmation and
  native delivery framing remain Claude-specific; Codex observations do not
  establish Claude behavior.
- Historical cleanup gaps remain explicit in Claude facts and design notes.
  Builds on macOS do not imply a fresh native macOS test.

## Completion rule

Account for every protected file/test and dependency; review the concrete diff;
run retained tests/race/vet/lint/packaging checks; install the exact artifact;
verify interactive and lane behavior with retained native evidence; report any
remaining limitation before merge. No runtime redesign, version bump or release
publication is part of the extraction.
