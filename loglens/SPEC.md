# Loglens: local log analytics and Go learning project

Status: interview-based specification, 2026-10-14. Loglens is a working name.

## Purpose

Build a tool its author uses to understand OpenClaw activity while developing and debugging. Develop fluency in Go context handling, concurrency control, and CPU-cache-aware data processing by implementing the core and measuring its behavior.

The author can reason about goroutines and channels and write Go, but currently looks up concurrency patterns frequently. The learning outcome is being able to implement and explain recurring patterns independently, including their cancellation, ownership, and failure behavior.

## Confirmed decisions

- Start with live OpenClaw logs on this machine.
- Use a command-line pipeline interface; no TUI or browser UI initially.
- Print a summary every 30 seconds. An explicit flag enables individual matching events.
- “Sampling” means periodic summarization, not randomly discarding events.
- Filter by severity, component, text, and time using flags. Add a query language later.
- Filters apply to summaries as well as individual event output.
- Maintain a bounded in-memory window and save raw events to disk before removing them from memory.
- Retain the archive for seven days, with a 1 GB disk limit.
- Preserve every input event during normal operation, accepting delayed output and backpressure.
- Implement core parsing, storage, and query machinery for learning; there is no shipping deadline.
- Use a separate replay benchmark for memory-layout and CPU-cache experiments.
- The user writes the main implementation. The assistant guides, reviews, and supplies incremental hints; it may help implement development setup, build tasks, test layout, and benchmark infrastructure. It must not supply complete core implementations unless requested.

## Proposed defaults

These make the design concrete but were not individually selected in the interview. They can be revised during implementation.

- Summary contents: severity counts, busiest components, frequent exact messages, and up to three error examples. Do not infer semantic importance or add AI summarization in v1.
- Start following files at their current end. Offer an explicit option to read existing content. Stdin processes all incoming content.
- Memory window: at most five minutes of processed events and 64 MiB of accounted retained event data, evicting when either bound is reached. This is not a promise of a 64 MiB process RSS limit.
- Archive cap: 1,000,000,000 bytes, including archive-owned metadata; retention measured from ingestion time.
- On archive failure or insufficient space after deleting expired segments, stop ingestion and report an error. Do not remove younger archive data automatically or silently discard events.
- Text output when interactive; explicit JSON Lines output for scripts.
- Default source in examples is explicit, rather than an installation-specific hardcoded path.
- Editor-specific configuration is deferred; editor preference was not provided.

## Observed input

Read-only inspection found dated files under `/tmp/openclaw/`, including `openclaw-2026-10-14.log`. Sampled lines were JSON objects with `0`, `1`, `2`, `_meta`, `time`, `hostname`, and `message` fields. This is a sample, not a complete schema guarantee. Other files under `~/.openclaw/logs/` included both JSON and plain text.

Before implementing the adapter, inspect field shapes and rotation behavior locally. Do not commit real log contents or private values. Create synthetic fixtures with matching structure. Do not assume OpenClaw is currently running merely because its logs exist.

## First usable workflow

Illustrative CLI, subject to adjustment as the parser is designed:

```sh
loglens watch --file /tmp/openclaw/openclaw-2026-10-14.log
loglens watch --file /tmp/openclaw/openclaw-2026-10-14.log --level error --component gateway
loglens watch --stdin --text timeout --show-events
loglens watch --stdin --since 5m --format json
```

The first useful slice reads JSON lines from stdin and prints a summary. Direct file following comes next, including an explicit directory/pattern mode for OpenClaw's dated files. Do not implement a subprocess supervisor or OpenClaw service management.

Filter semantics:

- Different flags combine with AND.
- Repeated severity or component values combine with OR within that field.
- Severity means exact selected levels, not an implicit minimum level.
- Component matching is exact; text is a case-sensitive literal substring in the normalized message.
- `--since 5m` is a rolling event-time cutoff evaluated when an event enters the query stage. Missing or invalid event timestamps use ingestion time and carry a visible fallback marker.
- Filters select output and query state, not archival: retain all input events in the archive.
- Do not silently guess missing severity/component values. Represent them as unknown unless an adapter has an established mapping.

## Summaries and output

Every 30 seconds of process time, summarize matching events processed since the preceding summary. Intervals do not overlap; an event contributes once. Boundaries use processing time so late input is not silently lost from summaries. Event timestamps remain available for filtering and examples.

Report interval bounds, matching count, severity counts, component counts, frequent messages, and error examples. Include processing totals, malformed counts, and an indication of lag where measurable. A heartbeat summary with zero matching events distinguishes inactivity from a stalled process.

The rolling memory window is distinct from the summary interval: summaries use incremental counters, so evicting a raw event does not remove its contribution to the current summary. Future window queries may only inspect retained events unless explicitly reading the archive.

`--show-events` adds matching individual events alongside summaries. JSON output has an explicit record type (`event`, `summary`) and version. Send diagnostics to stderr; stdout carries data only. Do not use terminal escape sequences in piped output. Flush a final partial summary on graceful shutdown and stdin EOF.

Bound aggregation state too: arbitrary unique messages or component names must not create unlimited maps. Start with capped distinct-key tracking and an explicit overflow count; do not label capped results as exact global top-K. Any later approximate algorithm must disclose its semantics.

## Input integrity, backpressure, and storage

Use bounded queues throughout. When downstream processing or stdout is slow, upstream stages wait with cancellation support. Never create a goroutine per event.

Archive the original raw record before admitting it to eviction-eligible memory state. Preserve malformed records as raw input, count the parse failure, and continue where record boundaries are known. A partial trailing line in a live file waits for completion. Define and test EOF behavior for a final unterminated stdin record.

Use append-only archive segments and explicit metadata; core archive logic is user-written Go. Flush buffered writes at summary boundaries and shutdown. Document that a successful buffered append is not a guarantee against power loss. Stronger fsync and crash-recovery guarantees are a later decision, not an implicit v1 promise.

Delete only expired segments; if mixed-age segments prevent reclaiming space, stop rather than exceed the cap. Account for metadata and avoid deleting a segment still in use. Archive failures must be visible and produce a nonzero exit status.

“Preserve every event” means no intentional ingestion drops and no eviction of an unarchived admitted event. It does not mean infinite retention or guaranteed recovery after arbitrary process/host failure. Saving events before eviction changes the earlier memory-only preference: this product now has both a bounded memory window and a raw disk archive.

For stdin, blocking can propagate to a cooperative producer. For file following, blocking only slows our reader; OpenClaw can continue writing. Track file identity and offsets, handle rename-based rotation, and detect truncation/disappearance where possible. Warn clearly on detected gaps. Copy-truncate and source deletion can lose unread data, so do not claim exactly-once delivery or guaranteed gap detection. Verify OpenClaw's actual rotation behavior before promising the adapter's coverage.

Set record-size bounds explicitly during the input milestone. An oversized record must either stream safely into archival handling or halt with a clear error; never silently skip it or allocate unbounded memory.

## Architecture and ownership

Logical flow:

```text
stdin / file follower
  -> bounded raw-record reader
  -> raw archive append
  -> parser and OpenClaw normalization
  -> filter + bounded memory window + summary accumulator
  -> terminal / JSON output
```

These are logical stages, not a requirement to give every stage a goroutine. Start sequentially. Add parallel parsing or batching only through explicit experiments. Preserve input sequence identifiers if parallel processing would otherwise reorder output.

One command owns the root lifecycle. Pass contexts through blocking operations. Define which component owns each buffer, channel, file, and mutable data structure. Prefer a single owner for aggregation initially; compare alternatives only after establishing a baseline.

Shutdown has two phases: stop accepting new input, then drain admitted work within a bounded shutdown deadline and flush the final summary/archive. Cancellation of intake must not accidentally cancel all draining work. Fatal archive/output errors stop the pipeline and return an error. An interrupted drain reports incomplete work rather than claiming success.

Use the standard library freely. “Implement the core” means own parsing/normalization flow, archive structure, window management, filters, and aggregation; it does not require writing a JSON decoder or replacing OS file APIs. Implement unfamiliar concurrency patterns in small exercises before deciding whether a standard helper is appropriate in the product.

## Learning and delivery milestones

1. **Input contract and setup.** Inspect representative shapes without publishing sensitive content; create synthetic fixtures; establish a Go module, build/test commands, and test directories. User writes a sequential parser and normalizer. Exit: valid, malformed, missing-field, and partial-line cases have explicit behavior.
2. **First useful stream.** User implements stdin ingestion, filters, 30-second summaries, optional event output, and an injectable time source for tests. Exit: a synthetic stream yields predictable summaries without real-time test sleeps.
3. **Bounded state and archive.** User implements raw segments, retention/cap handling, a bounded window, and bounded aggregate state. Exit: stress input cannot grow these structures without limit, and archive failure cannot silently discard data.
4. **Live OpenClaw following.** Implement file follow and dated-file discovery. Exit: append, rotation, truncation, and partial writes have tested behavior and documented limits.
5. **Context and concurrency.** Add cancellation, controlled parallelism, backpressure, and graceful drain. Exit: cancellation releases blocked producers/consumers, ordering is defined, and race checks pass.
6. **Replay and cache experiments.** Benchmark equivalent implementations on fixed representative data. Exit: report speed, allocations, working-set size, and correctness results; retain only improvements supported by evidence.

For each milestone, the assistant supplies a small design task, questions about ownership/failure behavior, and acceptance checks. The user implements, then receives review and hints. Periodically rebuild a worker pool, cancellable pipeline, or bounded window from memory and explain when it terminates. Progress is demonstrated by independent implementation and explanation, not memorizing APIs.

## Verification and performance experiments

Use synthetic fixtures and deterministic clocks. Necessary tests include filter combinations, interval boundaries, timestamp fallback, malformed input, retention/cap edges, archive write failure, slow consumers, cancellation while queues are full, EOF, rotation, and shutdown with pending work. Run race detection for concurrency changes. Fuzz record parsing once its contract stabilizes.

Separate three measurements: parsing, in-memory filtering/aggregation, and full replay including archive I/O. Keep workload and output equivalent across variants; setup and fixture generation stay outside timed sections.

Experiments progress from a readable baseline to:

- Slice of records versus separate columns for frequently scanned fields.
- Full string fields versus compact IDs for severity and components.
- Individual records versus bounded batches of different sizes.
- Shared aggregation versus per-worker aggregation followed by reduction.
- Worker-count sweeps and working sets below and above relevant cache capacities.
- A controlled false-sharing microbenchmark, separate from claims about the product.

Record Go version, CPU, dataset shape/size, worker count, batch size, throughput, latency, bytes allocated, and allocations per operation. Measure repeated runs and variability. Use CPU profiles and hardware performance counters where available and permitted; document limitations. A faster benchmark alone does not establish fewer L1 misses. Cache sizes and behavior are machine-specific, and storage throughput must not obscure an in-memory experiment.

Set numerical performance targets after measuring the baseline and real source rate. Initial success means keeping up with observed OpenClaw input, bounded accounted state, correct summaries, and clean cancellation. Report archive growth and estimated time to the disk limit from measured input volume.

## Deferred scope

No TUI/browser interface, remote collection service, distributed storage, custom query language, AI summaries, automatic message-template clustering, or CPU-specific assembly in v1. Archive search/replay as a user-facing query command follows the live workflow; benchmark replay is included earlier for learning. No core implementation is part of this specification task.
