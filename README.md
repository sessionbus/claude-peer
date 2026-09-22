# Sessionbus Claude peer

Connect native Claude sessions through [Sessionbus](https://github.com/sessionbus/sessionbus).
`claude-peer` provides interactive peers and managed lanes with the same permanent
Go installation, native history and permission controls.

## Install

Install native Claude and Sessionbus first using your normal home, login and PATH.

```sh
curl -fsSL https://raw.githubusercontent.com/sessionbus/claude-peer/main/scripts/install-claude.sh | sh
```

This repository is being separated from the original peers tree. No independent
release is published yet; use a reviewed archive built from source until release.
The installer retains checksum verification and latest stable/development selection;
it does not silently fetch another product. Older published installers remain at
[the original repository](https://github.com/sessionbus/codex-peer).

See [the Claude guide](claude/README.md) for full install/update/remove instructions,
private aliases, native flags, permissions, resume, identity and lifecycle behavior.
Ordinary `claude` stays ordinary. The archive includes the Go executable and native
plugin/hooks/skill; no Node.js/npm runtime is required for this integration.

## Build and test

```sh
git clone https://github.com/sessionbus/claude-peer.git
cd claude-peer
GOWORK=off go mod download
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
scripts/package-claude ./dist
```

Builds require Go 1.24 or newer; target installations do not need Go. Packaging
supports Linux/macOS amd64/arm64. RELEASE_VERSION and the Claude manifest agree;
this extraction does not bump either. Publication remains held during validation.

Shared support uses the exact peer-common version/checksum in go.mod/go.sum.
Native Claude versions are not pinned: users routinely update native clients.
Observed native versions/hashes identify test evidence, not a runtime allowlist.

The [stable functionality checklist](docs/migration/FUNCTIONALITY-CHECKLIST.md)
and [preservation inventory](docs/migration/PRESERVED-FILES.json) track separation.
Historical behavior and limitations remain in [Claude facts](docs/products/claude.md)
and [implementation records](docs/designs/claude-0.5.0/GO-MIGRATION.md).
The historical Node reference and held skills stay documentation only and are not
packaged or activated. Fresh extracted-build validation remains pending.

## Delivery and presence


A delivery reported as `rejected` with reason `no_receipt` means that no usable
receipt was obtained. It does not prove the message was never submitted or
consumed, including when the recipient disconnects. Preserve the delivery ID,
reason and any run reference; report the uncertainty without automatically
resending. A later connected or idle-looking row does not make replay safe.

In a `list` row, `connected` describes the Sessionbus attachment and `running`
describes a daemon-managed Run. An interactive peer's `running:false` does not
prove its native model is idle. Do not use these flags to predict delivery
admission; follow the actual receipt.

The orchestrator's installed tool declaration governs the tool identifier and
`{action, arguments}` envelope. After selecting a lane product, that product's
`describe` response and product skill/README govern its open fields, native
permission values, delivery receipts and lifecycle. `describe` lists available
fields; product documentation supplies their meaning and allowed native values.
Do not apply the orchestrator product's native options to a different lane product.

### Parent tracing

The unified tool's `trace` action can enable `events` or `content` tracing for
one live direct child, and `spawn` accepts the same initial `trace` setting. The
default is `off`. The daemon sends at most one ordinary message copy to the
eligible parent after the original Sessionbus send settles; `content` includes
the body, while `events` contains message and delivery metadata. Run lifecycle,
native prompts, and native results are outside this initial scope.

Tracing adds no history, durable policy, replay, catch-up, or separate event
transport. It requires Sessionbus v0.5.4 or later on the daemon and every involved hub,
and updated peer tools. Upgrade all of them before requesting tracing across
hosts. See the [v0.5.1 notes](docs/releases/v0.5.1.md) and the
[Sessionbus communication-trace contract](https://github.com/sessionbus/sessionbus/blob/main/docs/designs/COMMUNICATION-TRACE.md)
for the complete authority and delivery rules.

The tracing relationship ends with the parent's live lifetime, even for a
persistent child. Reconnecting with an old parent ID does not recover that
relationship or replay copies. Older daemons reject the new `trace` action or
spawn field; a trace-aware daemon returns `unsupported_trace` when an involved
federation link cannot enforce the requested tracing controls.
