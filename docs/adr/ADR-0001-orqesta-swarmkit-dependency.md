# ADR-0001 — Pin the Orqesta SwarmKit fork

- **Recorded:** 2026-11-19
- **Status:** Accepted; integration at the pinned revision verified, latest SwarmKit not yet integrated
- **Scope:** `orqesta/moby`

## Context

The Moby fork must consume the intended Orqesta SwarmKit implementation, including the stop-first grace-period correction, without introducing a second scheduler or changing Moby's import paths.

## Decision

Use one Go module replacement: `github.com/moby/swarmkit/v2 => github.com/orqesta/swarmkit/v2 v2.0.0-20261112170507-cf6d063b76df`. The pinned revision is explicit and reproducible. Updating the fork version is a separate change subject to validation.

## Implementation

The dependency, checksums, and the initial vendored content were synchronized in commits `2e066c9`, `6ffafa0`, and `afa79db`. Subsequent ADR-0002 removed the vendored dependency tree from Moby while preserving the module replacement in `go.mod`.

## Consequences

Moby builds against a single pinned SwarmKit fork via native module resolution. This does not mean Moby consumes the current SwarmKit head (`0d604ab`) or that the full Orqesta engine has been validated in a running cluster.

## Validation

Orqesta integration CI run [43660896447](https://github.com/orqesta/orqesta/actions/runs/43660896447) passed its dependency and compile checks. Those checks are not equivalent to operating a Docker Engine built from the fork.

## References

- [Moby commit 2e066c9](https://github.com/orqesta/moby/commit/2e066c917df44fe63d10fb6a3e97e0f30ea18a11)
- [Moby commit afa79db](https://github.com/orqesta/moby/commit/afa79db7e712694344618a057d5bf591d7e43723)
- [Orqesta ADR-0008](https://docs.google.com/document/d/1SIs2f6wokcoUJzZtvK8kHGlbvvwlawr7-hhdkAspE8w/edit)
- [ADR-0002 — native Go Modules](ADR-0002-native-go-modules.md)
