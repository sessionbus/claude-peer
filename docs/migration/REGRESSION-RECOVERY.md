# Claude baseline regression recovery

Status: old communication regressions pass; fresh clean wake acceptance is 3/4.
Interactive ACTIVE returned the requested direct reply but omitted the requested
wake-final marker, so its strict test remains failed. Historical failed runs and
retrospective assessments remain unchanged.

## Comparisons that must remain separate

| Comparison | Established result | Outstanding proof |
|---|---|---|
| Installed `8be5802` September 19 acceptance | Ordinary/bypass lanes and ordinary/yolo interactive peers completed native list/send with receiver correlation and cleanup; both deny guards rejected before startup | Fresh four scenarios plus two deny guards pass on the permanent installation |
| Previous installed `a1af0c15` to mandatory-wake `6b4836` | Managed delivery changed from native append to refusal before submission, followed by daemon scheduling of a query run; three append tests were replaced | Reconciliation and missing wrapper/Worker regressions implemented in [MANAGED-DELIVERY-REGRESSIONS.md](MANAGED-DELIVERY-REGRESSIONS.md); same-turn behavior remains changed |
| Extraction baseline `710e5d3` to `ff8471` | Protected runtime preserved with module import relocation | This alone proves nothing about equivalence with the older installed wrapper |
| New wake cells B/D/H/I | Native exchanges independently accepted from retained evidence | All four original drivers failed; fresh CLW923A/B/C drivers pass, while CLW923D lacks its requested wake-final marker |

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

## September 23 implementation and fresh results

Added direct native-frame assertions for idle/active refusal, original-result
retention, wake uncertainty, UUID/session correlation, and the Worker refusal
followed by a separately supplied wake run. Full Go tests and targeted race
checks pass. The Worker test simulates the daemon's second request; daemon queue
and turn-boundary behavior have separate Bus tests. No Claude product runtime
code or permanent installed binary changed during this recovery.

The consolidated local runner is `claude-wake-regression-runner-20260923` under
`/home/antst/sessionbus-evidence`. Independent review and 100 offline tests pass.
It handles optional discovery and final prose, distinct peer carriers and
absorbed mid-turn attachments, exact admission inventory, actual managed
acknowledgements, final native generation and final-to-stop ancestry. Current
native version/hash are observation provenance, with no version allowlist.
The old communication runner separately passes nine offline tests.

### Previous communication contract: four fresh passes and two deny guards

All scenarios used permanent wrapper `ff8471` / SHA256 `4df3ec5d…` and native
Claude 2.1.280. The wrapper, package inventory, static settings and service stayed
exact before/after each test; only native-owned `~/.claude.json` changed.

| Scenario | Fresh cell suffix | Result |
|---|---|---|
| Ordinary managed lane | `lane-default-6f27faea` | Default policy; list({}) → send → direct receipt; completed, acknowledged, closed and forgotten |
| Typed bypass managed lane | `lane-bypass-983539f8` | Explicit bypass policy; same native communication and lifecycle proof |
| Ordinary interactive peer | `peer-default-6bf5399c` | Inherited auto policy; list/send plus successful owned touch, file removed; final/descendant stop and normal exit 0 |
| Interactive yolo peer | `peer-yolo-080c9160` | Explicit bypass policy; list/send, final/descendant stop and normal exit 0 |
| Peer deny guard | `guard-peer-726eb07b` | Exact managed-tool denial before native startup |
| Lane deny guard | `guard-lane-038f3f3e` | Exact managed-tool denial before native startup |

Packets are under `claude-baseline-regression-prep-dev1-20260923/runs`, prefixed
`claude-baseline-`. Root and dev2 reviewed the fresh communication results.
Both interactive cells sent only the counted local exit command after their
native terminal. No permission setting was changed to obtain the ordinary pass.
An unresolved approval prompt would fail this runner without granting it; the
runner does not implement or claim a demonstrated native UI decline.

### Fresh wake contract: three clean passes, one unmet final requirement

Packets are under `claude-wake-regression-live-dev1-20260923/cells-clw923*/`.
These are new runs, separate from CLW922B/D/H/I and their retrospective records.

| Surface | Cell | Original driver result | Observed behavior |
|---|---|---|---|
| Managed idle | CLW923A | PASS, exit 0 and original phase pass | One injected inbound; exact envelope/reply/final; distinct setup/wake runs, both collected and acknowledged |
| Managed active | CLW923B | PASS, exit 0 and original phase pass | One queued inbound with live original sleep witnesses; successful Bash; distinct automatic wake, reply/final; both runs collected and acknowledged |
| Interactive idle | CLW923C | PASS, exit 0 and original phase pass | One written inbound; distinct peer carrier; exact reply/final/descendant stop and final native generation |
| Interactive active | CLW923D | FAIL, exit 1; no original phase pass | One written inbound and exact direct reply; absorbed attachment on original turn; sole substantive final contains only the setup marker, with no requested wake-final marker |

Dev2 independently reproduced the first three projections and original phase
checks. Managed cleanup includes successful close/forget and unknown target;
interactive idle cleanup followed the observed terminal with no survivors.
For A–C, retained scoped public rosters and owned process scans were empty.
No wake test received
a later model prompt or keyboard input, and no inbound was resent.

For D, the retained `native-history-during-wait.jsonl` (84 rows), SHA256
`b68cb8c6b1c92242420d67ec73a6b05e24549950313ce0aa1f6a8cc62cc68d87`
has valid absorbed admission, successful Bash, one correlated Sessionbus reply,
and a substantive final with a descendant stop. That final is exactly
`SETUP_CLAUDE_PEER_ACTIVE_CLW923D`; the requested
`WAKE_FINAL_CLAUDE_PEER_ACTIVE_CLW923D` is absent. The live predicate correctly
rejects it. This is an unmet test instruction, not another carrier-schema miss
or evidence that delivery/reply failed. No assertion is relaxed and no automatic
retry is authorized by this record. The original driver timed out waiting for
its terminal artifact; its SSH cleanup wait also timed out. The retained remote
cleanup reports `cleanup_before_wake_terminal`, SIGINT and no survivors. The
86-row `native-history-after-cleanup.jsonl` extends the earlier copy only with
last-prompt/cost-state rows and retains the same missing-marker verdict. The
terminal capture is retained; final installed
observations show no owned Claude processes and unchanged product/static
settings/service. A separate root native-tool observation after cleanup records
unknown target and an empty scoped roster in
`root-public-cleanup-observation.json`; it is not an original driver artifact.
The owned temporary root and local driver/SSH/controller processes were removed.
Cleanup does not turn the missing final into a pass.
Full four-surface clean acceptance remains
open; same-turn managed compatibility remains separately un-restored.

### Preserved prerequisite failures

The earlier ordinary cell `lane-default-64fe2ff6` reached `list({})` and then
hit the controller's 180-second deadline. Native history records a user-rejected
tool result at that deadline; it does not prove a normal permission-policy
rejection. No send occurred, and owned cleanup completed. The earlier uncertain
spawn `lane-default-e33e26f0` did create a disconnected lane despite its RPC
error; that row was recovered, closed and forgotten. Both failures remain
separate from the fresh passing cells.

All-host listing also stalled independently outside Claude while MBP slept.
At running pdev daemon source `326bc81`, `collectFederatedList` has no per-host
or aggregate deadline. After the owner woke MBP, the unchanged all-host query
returned promptly and the fresh baseline passed. The owner reported the sleep
and wake; each fresh baseline cell also retains a successful all-host query,
including MBP rows. No MBP session is required by
the tests, and no daemon/authentication change was made. The missing federation
deadline remains a separate Bus defect; these observations do not prove which
host leg the earlier native call reached. The detailed local record is
`claude-baseline-regression-prep-dev1-20260923/FEDERATION-BLOCKER.md`.
