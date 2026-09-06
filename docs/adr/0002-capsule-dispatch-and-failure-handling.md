# 0002 — Capsule dispatch and failure handling

**Status:** accepted

## Context

Capsules are delivered by a cron job, not by the request that created them, so
delivery has to survive restarts and must never send the same capsule twice.
Duplicate delivery is the worse failure here: the recipient sees a message the
sender wrote once, arriving twice, years later.

## Decision

`capsules.status` is a three-state enum: `waiting → in progress → done`.

Each tick claims work in a single statement (`repository.Capsule.ClaimDue`): a
CTE selects due rows `FOR UPDATE SKIP LOCKED` and the enclosing `UPDATE` flips
them to `in progress` and returns them. `SKIP LOCKED` means several dispatcher
instances can run concurrently without ever handing the same row to two of them,
and the claim is atomic, so a crash between select and update is impossible.
A partial index on `send_at WHERE status = 'waiting'` keeps the scan cheap as
the `done` rows accumulate.

Failure handling is deliberately asymmetric:

- A failure **before** the mail leaves — the object store is unreachable, SMTP
  refuses — reverts the row to `waiting`, and the next tick retries it.
- A failure **after** SMTP accepted the message, i.e. the `done` update fails,
  leaves the row `in progress` and logs it. It is never retried automatically.

Delivery fan-out is bounded (`capsuleSendConcurrency`) so a large batch cannot
open one SMTP connection per capsule.

## Consequences

- At-most-once delivery: the system prefers a stuck row to a duplicate email.
- A row left `in progress` needs a human, or a sweeper that has some way to
  confirm whether the message actually went out. There is no such sweeper.
- A capsule whose send fails repeatedly retries on every tick forever; there is
  no attempt counter and no dead-letter state.
- Cron schedules are evaluated in UTC, and `send_at` is compared against `NOW()`,
  so delivery does not shift with the host's timezone.
