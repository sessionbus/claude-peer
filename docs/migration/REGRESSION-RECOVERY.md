# Claude baseline regression recovery

Status: in progress. This record supersedes a blanket completion claim; the
sealed historical runs and retrospective assessments remain unchanged.

## Comparisons that must remain separate

| Comparison | Established result | Outstanding proof |
|---|---|---|
| Installed `8be5802` September 19 acceptance | Ordinary/bypass lanes and ordinary/yolo interactive peers completed native list/send with receiver correlation and cleanup; both deny guards rejected before startup | Repeat those supported communication scenarios on the current permanent installation |
| Previous installed `a1af0c15` to mandatory-wake `6b4836` | Managed delivery changed from native append to refusal before submission, followed by daemon scheduling of a query run; three append tests were replaced | Reconciliation and missing wrapper/Worker regressions implemented in [MANAGED-DELIVERY-REGRESSIONS.md](MANAGED-DELIVERY-REGRESSIONS.md); same-turn behavior remains changed |
| Extraction baseline `710e5d3` to `ff8471` | Protected runtime preserved with module import relocation | This alone proves nothing about equivalence with the older installed wrapper |
| New wake cells B/D/H/I | Native exchanges independently accepted from retained evidence | All four original drivers failed; clean driver completion remains pending |

Historical reference: `claude-comms-grant-installed-20260919/PLAN.md`,
`OBSERVATIONS.md` and `verify_histories.py` under the development evidence root.
They record explicit tool communication, not four-surface automatic wake.

## Required regressions

Retain the historical ordinary and explicit bypass policy variants; never add a
bypass to make an ordinary-policy test pass. Both managed-tool deny guards must
still reject before native startup. Communication requires the native exact
Sessionbus tool, authenticated identity/list result, settled send and direct
receiver observation. Preserve the ordinary interactive unrelated-operation
policy observation without changing native policy to force its outcome.

Earlier installed evidence also demonstrated same-turn managed delivery initiated during a
held tool and admitted after its release (`lfc-active-allowed` and `lfc-active-admission`). The
mandatory-wake fallback does not preserve that scheduling behavior: an active
lane message waits for a subsequent managed run. Record this change explicitly;
a passing explicit list/send test cannot prove same-turn compatibility. The
retained mandatory-wake design explains the terminal-boundary and blocked-tool
liveness tradeoff, which must not be hidden by extraction preservation claims.

Source tests must preserve correlation, uncertainty after attempted submission,
no premature receipt, original-result retention and cancellation across the
new managed-delivery contract. Removed internal append APIs are not restored
merely to rerun obsolete unit tests. Any supported external behavior loss is a
product defect to fix and review separately from the runner.

## Clean wake completion

Use one consolidated runner candidate, tested offline against retained positive
histories and mutation negatives before live use. Known valid native forms
include optional discovery, final prose with an exact standalone marker,
interactive peer carrier rows and active absorbed-mid-turn attachments.
Counting only already-matching rows must not hide unrelated input or tools.

Each of the four fresh surfaces must complete its own driver and phase checks:
actual idle/active witness, one authenticated inbound, correlated reply and
native terminal, no later model/keyboard input, collected and acknowledged
managed terminals, and owned cleanup. Retain final process evidence before
cleanup where required. Collect roster before an expected unknown-session
error on a controller that exits on errors. An offline assessment does not
convert a failed original driver into a clean pass.

Run the previous communication scenarios first, then managed idle, managed
active, interactive idle and interactive active serially. Preserve the first
failure and investigate; do not resend an uncertain message or repeatedly run
models to tune assertions.

## Environment and reporting

Use UMKA's actual permanent installation and real home/config/login/service.
Dev1 owns host writes and live tests; source work and independent review happen
off-host. Record current native version and hashes as provenance, with no
product version allowlist. Keep current installed bytes unless a concrete
product defect requires a reviewed rebuild. No version bump or release.

Report source checks, historical evidence, retrospective assessments and fresh
clean driver results separately. Completion requires both the old communication
regressions and four clean wake surfaces; until then this record stays open.

## September 23 implementation and first live results

- Added direct native-frame assertions for idle/active refusal, original-result
  retention, wake uncertainty, UUID/session correlation, and the Worker refusal
  followed by a separately seeded wake run. Full Go tests and targeted race
  checks pass; no product runtime code changed.
- Consolidated the existing managed and interactive wake runner in the local
  `claude-wake-regression-runner-20260923` evidence directory. It handles the
  retained native discovery/final forms, distinct peer carrier and absorbed
  mid-turn attachment, actual managed acknowledgements, final live generation,
  and cleanup ordering. Independent review and 99 offline Claude tests pass.
  This is source validation, not fresh installed wake acceptance.
- Both original managed-tool deny guards pass against the permanent UMKA
  installation, with strict installed pre/post observations and no model turn.
- The ordinary managed communication run started, reached its exact Sessionbus
  `list({})` call, and hit the controller's bounded transport timeout before a
  reply. Further model cells remain pending. The original failure is retained;
  it is not converted into a pass.

The all-host list is independently stalled outside Claude too. At the running
pdev daemon source `326bc81`, `collectFederatedList` waits for every directed
host reply without a per-host or aggregate deadline. A directed `mbp` query
also did not return during the observed interval. Host-local identity listing
works and is used for the test controller, but does not replace the historical
communication regression's all-host query. No daemon or authentication change
was made to force a pass. The detailed local record is
`claude-baseline-regression-prep-dev1-20260923/FEDERATION-BLOCKER.md`.
