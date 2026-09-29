# Sessionbus peers development prerelease

This development/beta build publishes all eight peer archives for Linux and macOS, on amd64 and arm64. It adds Pi and OMP to Claude, Codex, Grok, Qwen, OpenCode, and Kilo. DSH is packaged separately.

## Changes

- All eight public peer commands report the peers release and exact source
  revision without starting the native product. `--native-version` is the
  explicit native-version escape outside lane mode. Stable packaging rejects
  drift between the release tag, central release version, Claude manifest and
  Codex plugin base, and each installed product retains its `SOURCE.txt`.
- Claude can publish before its first prompt by observing the exact native parent registry. The existing native hooks take permanent authority at their first usable report; unavailable registries retain hook-based startup. No dummy turn is generated.

- Claude interactive messages restore the structured envelope used by its compact sender display, instead of displaying native peer-origin prose as an ordinary user message. Native model-input guidance remains Claude's own.
- Codex lane results retain every final-answer item from the exact persisted native turn in item order. Multiple items are separated by one blank line; commentary, tool messages, and other turns remain excluded.
- The unified tool exposes parent tracing through `spawn.trace` and the `trace` action. Tracing defaults to `off`, produces at most one ordinary parent copy for a settled Sessionbus send, and adds no durable policy, history, or replay. It requires the forthcoming trace-aware Sessionbus daemon and every involved hub; these peer declarations alone do not enable it.
- Codex, Claude, Grok, Qwen, Pi, and OMP interactive owners reconnect after public daemon loss or initial absence while preserving the native session. Old calls and deliveries are not replayed. Supersession and native exit remain terminal. OMP Task child peers reconnect independently; lane Workers retain their single connection lifetime.
- Shared MCP and Pi/OMP native tool declarations expose the protocol's closed argument fields. In particular, send accepts message and target/targets, not summary. The SDK retains action-specific validation.
- Worker startup and service honor cancellation through the updated Go SDK. A Worker connection remains a single lifetime and does not reconnect.
- OMP graceful shutdown drains its owned native RPC output and report work before the owner finishes. Startup requests wait until the registry owns the native bridge; cancellation releases and joins the wait.
- Installers select latest stable, falling back to the published development prerelease only while no stable release exists. Explicit versions and mirrors remain supported; checksums are verified before installation.
- Pi and OMP include managed native identity, interactive launch, lane Run, staged delivery, and joined process/private bridge ownership.

## Preview status

Pi and OMP remain preview integrations. Installed evidence establishes Pi normal Run, delivery/Forget, interactive replacement, and startup/resume; OMP zero-input Worker/interactive and core normal Run. Some original fixture runs failed on verifier assumptions and were assessed from preserved evidence; those are not reported as full fixture passes. The detailed record is in [the design and acceptance record](https://github.com/sessionbus/pi-omp/blob/710e5d33369cba4fb9468cd24fea0fe844a0219d/docs/designs/pi-omp-0.5.0/DESIGN.md).

The newly integrated reconnect and schema fixes have deterministic and process-fixture coverage, but fresh installed-product validation is still in progress. Pi admitted interrupt/healthy recovery and OMP staged delivery/Forget remain selected work before removing their preview designation; they do not block a stable peers release that explicitly retains that designation. An OMP native extension writing arbitrary stdout can also invalidate the RPC connection.

This prerelease is available now for testing; publication does not mark the remaining stable acceptance work complete. No interrupted action is retried automatically. Native vendor products, credentials, and the Sessionbus daemon must already be installed separately.

The [v0.5.1 notes](v0.5.1.md) describe the coordinated tracing release and its
daemon compatibility requirements.

## Maintenance: native argument grammar

Whenever native Claude's argument grammar is reviewed or changes, re-check the typed-field and managed-tool guards with an isolated, offline probe:
- Use a throwaway `HOME` and `CLAUDE_CONFIG_DIR` and a fake API key.
- Point `ANTHROPIC_BASE_URL` at an unbound local port.
- Set `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`.
- For strict isolation, also run the probe under an enforced network restriction, such as a network namespace. The API URL and these settings alone do not block every outbound connection.

Run each argument list under `timeout 20`, for example `timeout 20 claude -p --append-system-prompt -- --permission-mode bogusmode`:
- If native exits at once with a validation error for `bogusmode`, the `--` after the value-taking option was consumed as its value and parsing continued.
- If `timeout` ends it with exit status 124, the arguments were treated as a prompt. The CLI keeps retrying the unreachable API until then.

Update the option table in `wrappers/claude/interactive/launch.go` and its tests to match. Record the probe commands, the native version and the observed results in the pull request that makes the change. Native versions are provenance, not a supported-version list.
