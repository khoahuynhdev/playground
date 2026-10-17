# Exercise 1: define and normalize one event

Goal: turn one raw JSON record into an event usable by filters and summaries. Keep this first step synchronous. No channels, file following, timers, or archive code yet.

## Before coding

Write down your answers and propose an event type and function signature:

1. What does the query stage need from an event? Which raw information must remain available to the archive instead?
2. How will callers distinguish malformed JSON from valid JSON with missing optional metadata?
3. Who supplies ingestion time so timestamp fallback can be tested deterministically?
4. Will the event own its bytes/strings, or borrow a caller's buffer? What happens when that buffer is reused?
5. Does parsing a single bounded record need a context? Explain where cancellation becomes useful in the larger reader pipeline.

Use `encoding/json` and `time`; implementing the core does not mean rewriting JSON decoding.

## Proposed first contract

- Normalize timestamp, message, severity, and component from the synthetic valid records.
- Missing/invalid time falls back to caller-supplied ingestion time, with a marker.
- Missing/unrecognized severity and missing component become unknown.
- Malformed JSON and non-object JSON return a descriptive error.
- Explicitly decide how wrong-typed fields behave. Keep this distinct from a missing optional field.
- Do not print from the parser or read the clock internally.

These rules narrow the first exercise; review real schema variation before expanding the adapter.

## Done when

- You wrote the implementation and table-driven tests covering valid records and the chosen edge-case behavior.
- Tests assert observable fields/errors, rather than duplicating parser logic.
- `make check` passes.
- You can explain ownership and fallback behavior without opening the implementation.

Ask for a conceptual hint, a review of your signature, or a failing-test diagnosis whenever useful. The assistant should not supply the completed parser unless you request it.
