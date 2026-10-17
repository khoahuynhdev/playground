# Loglens

A local OpenClaw log analytics CLI built as a Go learning project. The implementation has not started. See [SPEC.md](SPEC.md) for decisions and milestones.

## Development

Go 1.25 or later and Make are sufficient for the initial scaffold. The local toolchain used during setup was Go 1.25.8. Race checks need a working C toolchain on supported platforms.

```sh
cd loglens
make check
make fmt
make bench
```

The parser currently contains only a package declaration. Checks validate the scaffold; they do not establish application correctness. Benchmarks and executable build targets become useful once their implementations exist.

## First exercise

Read [exercises/01-parser.md](exercises/01-parser.md). Write your proposed event type and parser function signature first, then discuss them before implementation if helpful.

## Test layout

- Unit tests live beside their package in `*_test.go` files.
- Synthetic parser inputs live in `internal/parser/testdata/`.
- Add command integration tests when a CLI exists.
- Add benchmarks beside the code they measure; keep fixture setup outside timing.

Create additional packages only when a milestone needs them. Raw archives belong in an explicit data directory, never in committed fixtures.
