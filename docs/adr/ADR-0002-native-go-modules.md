# ADR-0002 — Use native Go Modules without a vendored dependency tree

- **Recorded:** 2026-11-19
- **Status:** Accepted for dependency resolution; Protobuf generation replacement remains under review
- **Scope:** `orqesta/moby`

## Context

The fork inherited checked-in `vendor/` directories and custom dependency synchronization routines. That duplicated the dependency state already maintained in `go.mod` and `go.sum` and required special handling for builds.

## Decision

Use the standard Go module graph and checksum verification as the sole source of dependency versions. Do not check in `vendor/`, introduce dependency-installation scripts, or add fallback dependency mechanisms. Normal `go build` and `go test` resolve module dependencies through Go.

## Implementation

Commit `e2579f4` removed the root and manual-page vendored directories, obsolete vendor synchronization/validation integration and vendoring workflow. Commit `f8c656f` removed transitional dependency helper scripts. Commit `3d5e782` restored executable bits accidentally lost on existing validators. The original Protobuf generator still required separate remediation; commits `124d8c7`, `a020e0e`, and `9e6c35c` concern that distinct, not-yet-accepted work.

## Consequences

The Moby repository no longer tracks the dependency vendor tree. The module graph remains auditable and reproducible using Go. This decision does **not** assert that SwarmKit's separate repository is already vendor-free.

## Validation

On commit `a020e0e`, the [CI](https://github.com/orqesta/moby/actions/runs/44242193232), [test](https://github.com/orqesta/moby/actions/runs/44242193209), [Buildkit](https://github.com/orqesta/moby/actions/runs/44242193321), [Windows](https://github.com/orqesta/moby/actions/runs/44242193210), [VM](https://github.com/orqesta/moby/actions/runs/44242195808) and [binary image](https://github.com/orqesta/moby/actions/runs/44242193914) workflows passed. The new compiler/generator implementation in `9e6c35c` has separate CI that must not be represented as passed until confirmed.

## References

- [Moby commit e2579f4](https://github.com/orqesta/moby/commit/e2579f4ff7277d5e688cabf9e4024cf7904070d4)
- [Moby commit f8c656f](https://github.com/orqesta/moby/commit/f8c656f56b041ddcc9fa780e1f1a4a85827512ca)
- [Moby commit 3d5e782](https://github.com/orqesta/moby/commit/3d5e78296860f7117e8831c1c73ba5e24d8758e0)
- [ADR-0001 — SwarmKit dependency](ADR-0001-orqesta-swarmkit-dependency.md)
