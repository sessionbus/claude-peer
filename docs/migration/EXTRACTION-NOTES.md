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

Local protected-file normalization passed for all 64 files and 133 tests.
Independent source review and hosted Linux/macOS/scan/workflow checks passed at
`ff8471b6d9965a39c934e7a8fe3a0a4f706b3e4e`. PR #1 merged as
`99709ed7a5af1241a1a2f508886bdfd3a35c6369` with the same tree. Four platform
archives built; the Linux amd64 archive was installed and reinstalled in UMKA's
real home with identical installed inventories. The common pin is unchanged.

Fresh managed idle/active and interactive idle/active exchanges are accepted
through separately reviewed retrospective assessments. The original test-driver
failures and their limitations remain recorded; see the
[installed acceptance record](FUNCTIONALITY-CHECKLIST.md#installed-acceptance).
The assessment process made no model calls or product changes. Native versions
2.1.278 and 2.1.280 identify those observations, not a compatibility allowlist.
No release, tag or version bump accompanies this extraction.

The previous UMKA a1af0c wrapper predated the merged 6b4836 mandatory-wake fix.
This extraction preserves 6b4836 from the 710e baseline and validates the newly
installed ff8471 artifact. Earlier explicit communication passes were not reused
as fresh automatic-wake acceptance. The obsolete b6fb0e6 Node portability branch
remains archived on the original Codex repository; it was not cherry-picked into
the Go implementation.

Local extraction checks (test, race, vet, module verify/tidy, golangci-lint,
actionlint and protected source/test normalization) passed. Native macOS behavior
and the lifecycle qualifications in the checklist remain outside these four
UMKA wake observations.
