# Claude extraction record

Baseline: `710e5d33369cba4fb9468cd24fea0fe844a0219d` from the complete organisation
copy. The original repository history and sibling copies retain removed product
paths. Claude implementation, plugin hooks/manifest/MCP declaration, skill, native
argument and permission handling, install layout and private aliases remain.

The frozen inventory accounts for 64 protected files and 133 original Go test
functions. Forty-seven files remain byte-identical; sixteen Go files change only
module imports and gofmt, and `claude/README.md` changes its bootstrap URL.
Product runtime contains no behavior change. Shared host/MCP/version/socket
support resolves to `github.com/sessionbus/peer-common` at the reviewed immutable
version `v0.0.0-20260922143100-eb655f686e44`, with no replace or alternate workspace.
The exact common code/tests preserve their original behavior. The linker stamps
the Claude product release/revision into the common version helper.

Root boundary/packaging/version/download tests and workflows are scoped to
Claude. Claude-specific boundary and manifest reachability checks remain;
other-product test cases are removed with their product paths. The version guard
retains Claude manifest/tag agreement and drops only the removed Codex manifest.
Package/installer implementation is unchanged apart from canonical download URL
and linker import path. Third-party notices now include the common module.
The package remains one Go integration with no Node/npm runtime dependency.

All Claude design notes, the reviewed historical Node reference, held skills,
mandatory-wake notes and product facts remain. Historical external-artifact links
that were already unresolved are not represented as live repository targets.
Three deletion-induced Pi/OMP links now point to immutable original source in
the family repo; the root README install anchor is updated.

Local protected-file normalization passed for all 64 files and 133 tests. Fresh
installation/acceptance and independent extraction review remain pending. A
historical pass does not validate this new artifact. Release publication remains
held and no version bump is made. Native client version updates are expected;
exact native versions in evidence are provenance, not a compatibility allowlist.

Retained native installation inventory identifies UMKA's installed a1af0c wrapper
as predating the repository's merged 6b4836 mandatory-wake fix. The extraction
preserves 6b4836 because it is already in the 710e baseline; fresh installed wake
acceptance for Claude is still outstanding. Earlier explicit communication
passes do not establish automatic wake. The obsolete b6fb0e6 Node portability
branch remains archived on the original Codex repository and is not blindly
cherry-picked into the Go implementation.

Local extraction checks: go test, race, vet, module verify/tidy, golangci-lint,
actionlint and protected source/test normalization pass. There is no new model
call or UMKA mutation in this source preparation.
