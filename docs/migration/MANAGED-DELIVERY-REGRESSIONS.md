# Claude managed lane delivery: old contract to current contract

This compares two source points for managed-lane delivery:

- **Old:** `8be5802`, pinned to Bus SDK 69c1024f and last installed on
  2026-09-19; and `a1af0c1`, pinned to SDK v0.5.5 and the last pre-split
  installed wrapper.
- **Current:** mandatory wake `6b4836`, pinned to SDK d765d80, and the split
  commits `ff8471b`/`536ae75`, pinned to SDK v0.5.7.

Only `6b4836` changed Claude lane delivery. It left interactive delivery through
the native messaging socket, the explicit list/send path and the managed-tool
deny guards untouched. Their tests are unchanged:

- `TestOpenRejectsManagedToolDenyBeforeLifetime`
- `TestNativeArgumentsKeepGrantWithBypass`
- `interactive/permissions_test.go`
- `interactive/delivery_test.go`

Native client versions below are test provenance, not a runtime allowlist.

## Behavior map

| # | Old observable behavior | Current behavior | Classification |
|---|---|---|---|
| 1 | **Idle:** Deliver writes a `shouldQuery:false` append. After exact replay it reports `queued_for_next_turn`. The text sits in native context until the next explicit run (installed `lfc-core`). | SDK v0.5.7 refuses idle delivery before calling the product. Daemon v0.5.7 forces `idle_message` to "run" and seeds a fresh managed run through the native query path, reported as `injected` at exact replay (installed CLW922B). | Removed by the platform mandatory-wake decision. Daemon `policy.go` says passive staging is no longer selectable. Not a Claude defect. |
| 2 | **Active:** the same append is `injected` at exact replay during the confirmed active run, and the text affects that turn's final result. This was verified on 2.1.260 and installed as `lfc-active-allowed`/`lfc-active-admission`. | Claude refuses with NotRunning before any write. The daemon reports `queued_for_next_turn` and starts a distinct automatic wake run after the terminal (installed CLW922D). | **Observable same-turn capability loss, Claude only.** Codex steer and Grok interject stay same-turn. The reason is recorded in the NATIVE-BOUNDARIES Claude row: no demonstrated turn-safe `shouldQuery:false` handoff at the terminal boundary, and no replay independent of a blocked tool batch. Restoring it needs an explicit owner contract decision and proof. It is not reintroduced by assumption here. |
| 3 | Append crossing a run boundary is uncertain (-32603). | Not applicable: there is no append, and each wake is its own run. | Removed along with row 2. |
| 4 | An attempted write without admission (stream loss) is uncertain (-32603); an unattempted one is `rejected/not_submitted`. | The wake `execute` path applies the same rule. | Preserved. Stream-loss coverage is added below. |
| 5 | Admission requires an exact UUID and session replay. | The wake path requires the same. | Preserved. Wrong-UUID coverage is added below. |
| 6 | An unavailable lane gives `rejected/native_unavailable`. | Unchanged. | Preserved. Coverage is added below. |

## Removed tests and their dispositions

- **`TestAppendWaitsForNativeReplay`:** the passive append behavior (row 1) was
  removed by design. Its correlation check lives on in
  `TestWakeReplayFixesReceiptWhileTerminalDrains` (session) and
  `TestWakeAdmissionRequiresExactUUIDAndSession` (UUID).
- **`TestAppendClassifiesAtReplayAcrossRunBoundaries`:** the behavior in rows 2
  and 3 was removed. The replacement contract is covered by
  `TestActualWorkerActiveDeliveryDefersToDistinctWakeRun`: with the real SDK
  worker and an active run, Claude refuses without writing, the active run
  completes unchanged, and a later daemon-seeded run uses the query path with the
  source envelope and a new native UUID.
- **`TestAttemptedAppendWithoutAdmissionIsUncertain`:** row 4. The stream-stop
  and native-EOF variants are covered by `TestWakeNativeLossAfterWriteIsUncertain`.
  The cancellation variant was already covered by
  `TestWakeWithoutReplayPreservesSubmissionUncertainty`.

## Weakness in the replacement test

`TestLaneDeliveryDefersBeforeAnyNativeWrite` has three weaknesses:

- It asserts only -32004.
- It covers only the idle state, which SDK v0.5.7 never passes to Deliver.
- It reads no native input. A write before the refusal blocks its pipe and fails
  only by timeout, and an asynchronous write would pass.

It is kept unchanged. `TestManagedDeliveryNeverWritesIdleOrDuringActiveRun`
uses a non-blocking frame recorder:

- zero frames while idle;
- exactly the original frame during an active run, both before and after its
  replay;
- the original run's result unchanged.

`TestManagedDeliveryUnavailableOrCancelledNeverWrites` covers row 6 and
cancellation.

## Mutation check

| Mutation | Result |
|---|---|
| Write an append before refusing | The idle/active and worker tests fail. The old replacement test hangs. |
| Report loss after write as `not_submitted` | The native-loss wake test fails. |
| Match replay by session only | The exact-correlation test fails. |

## Outside this change

- Daemon queue and retention, tested in the bus repository.
- Live native behavior, covered by the old-contract runner.
- Any decision to restore same-turn delivery (row 2).
