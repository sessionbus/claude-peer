# Stable functionality checklist — Claude separation

Baseline: original peers main `710e5d33369cba4fb9468cd24fea0fe844a0219d`.
The F01–F20 requirement IDs are shared with the migration checklist. Product
semantics and known limitations stay explicit; a failed check does not remove a
requirement. Extraction source preservation and fresh installed behavior are
separate evidence. Source preservation, independent review, hosted CI and permanent
installation/reinstallation passed. Four wake exchanges were accepted through
retrospective assessments; none of their original automated drivers completed
successfully. Fresh September 23 baseline communication passes 4/4 plus both
deny guards; clean automated wake passes 3/4. Interactive ACTIVE replies but
omits its requested wake-final marker. The
[regression recovery record](REGRESSION-RECOVERY.md) records the separate results.
The installed evidence below remains valid with its original limitations.

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

## Installed acceptance

Tested source: `ff8471b6d9965a39c934e7a8fe3a0a4f706b3e4e`; merge
`99709ed7a5af1241a1a2f508886bdfd3a35c6369` has the same tree. Installed wrapper
SHA256: `4df3ec5d644a0e5650126999edb89dbcc09e2b6c3f0967701a14d6292aad0d16`.
Install/reinstall packet seal:
`746d006040c6ae63cb08e763aeb3aab5f01681d455a8f9fdfb83447f3d963f0c`.
The common pin above, package/plugin files and private aliases were verified in
the permanent real-home installation. No alternate prefix or configuration was
used. B and D ran on native 2.1.278; H and I ran on 2.1.280 after the native
client updated during this work.

Root and dev2 independently reviewed the following four exchanges. Each original
run has a preserved failed driver outcome and null acceptance; the separate
assessment validates retained native history and cleanup without replaying the
message or running the model again. These are retrospective functional passes,
not claims that the original driver or original phase check passed.

| Requirement | Cell | Observed native behavior | Assessment seal (SHA256 of SHA256SUMS) |
|---|---|---|---|
| F13 managed idle | CLW922B | Injected delivery; distinct automatic wake run, reply and final | `6fb95574c71e483f7dd71701b2f92f81b67c124c4d331900d1748e9c7fc76662` |
| F14 managed active | CLW922D | Queued during original Bash; distinct automatic wake run after it | `5f5e2d87e061a28d1eed18cca6bd2a89b078ee57396f73b975321fd43a4c0c68` |
| F11 interactive idle | CLW922H | Written delivery; native peer carrier, reply, final and stop | `d57420481cdc59e77ed0c64702cbfc98bda1910feadb8d8a9d0205c3eb74e8f1` |
| F12 interactive active | CLW922I | Written while Bash active; native queue absorbed into the original turn as an attachment, reply and combined final/stop | `bc68d7075c6657464caf67e8487ed53fe1235d6e83c27b21f94b80dd82d75162` |

Original packet seals, in the same order:

- B: `b488d4ccec84f0aae23b191ee140d90407644c209459b72a070bad97f77d8f45`.
- D: `80b7d5f6104e91f45568f881e06992dd35185fb7b15e0fe2bd309d2bb69f2fdf`.
- H: `79fce9dd58c247c07d1fbe0bb5abaf361407eb3b672dc103091c7641200b39d1`.
- I: `f3fb395ef1d71c5795ffd882544088db4617ab5f54301e8f3175386b8423e805`.

Packets are retained under the development-host evidence directory
`/home/antst/sessionbus-evidence/claude-split-wake-installed-dev1-20260922/`;
they are not published repository artifacts. Seals identify exact retained
packets, and do not by themselves give public access to their contents.

B's driver failed when it queried the expected unknown-session error before
collecting the roster. D's strict text check rejected valid final prose first.
Its projector would also have rejected native ToolSearch `max_results: 3`, but
that check never ran in the original driver. D's original run was acknowledged;
its wake result was collected but
**not acknowledged**. No synthetic acknowledgement is claimed.

H's checker omitted native peer framing. I's checker expected a second ordinary
user row, while native history recorded `absorbed_mid_turn` and a
`queued_command` attachment on the original turn. The assessments verify the
entire carrier and independently reconstructed authenticated envelope. They do
not invent a second ordinary user or separate wake turn for I.

H and I retain the original timeout, missing runtime terminal artifact and
`cleanup_before_wake_terminal` status. Their final/stop evidence comes from
persisted native history; process generations come from the historical
ready/send-window witnesses (also post-receipt for I), with **no final process
sample**. No original close/forget event is claimed. Owned interruption left no
survivors; independently retained roster/process observations were empty and
target queries returned unknown-session. Static configuration and installed
inventory checks passed. Native-owned state updates were retained separately.

Every accepted cell used one inbound send and no later model, keyboard or
lifecycle input. Interactive setup used the native positional prompt with zero
PTY writes. Active fixtures granted only the exact sleep command per launch.
Earlier A/C/E/F/G diagnostics remain preserved; none is represented as a wake
pass. Root independently confirmed that G's settings serialization change
preserved all JSON values; the writer remains unknown. G's separate presence
correction preserves the observed public session without claiming a completed
setup.

These retrospective assessments establish the four observed wake exchanges on
the installed Linux artifact. They do not close clean automated acceptance of
F11–F14, or replace the earlier explicit communication regression scenarios.
Source preservation is measured against the extraction baseline, which already
contained changes absent from the prior UMKA installation. Fresh coverage of
every lifecycle path, platform, permission mode or update is not claimed. No
native version allowlist, product redesign or release is introduced.
