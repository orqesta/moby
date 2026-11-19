# Orqesta integration and build conventions

This document describes how the Orqesta fork of Moby is currently assembled. It is not a guide to the upstream Docker Engine release process.

## Dependency resolution

The root Moby module and its `api`, `client` and `man` modules use their own `go.mod` and `go.sum` files. Dependencies are resolved by the Go toolchain using those module definitions and checksum verification. No dependency `vendor/` tree is tracked in this fork.

From the root module:

```sh
go build ./...
go test ./...
```

The project's existing CI also performs platform-specific builds and validation. Ordinary compilation and tests do not require running dependency synchronization scripts or regenerating generated files.

## SwarmKit integration

The root `go.mod` declares `github.com/moby/swarmkit/v2` as a dependency and uses a single module replacement:

```go
replace github.com/moby/swarmkit/v2 => github.com/orqesta/swarmkit/v2 v2.0.0-20261112170507-cf6d063b76df
```

This preserves Moby's existing SwarmKit import path while selecting a specific, immutable source revision from the Orqesta fork. It does not yet select the latest SwarmKit fork revision. Moby also still depends directly on `github.com/docker/go-events` for its own consumers; SwarmKit's migration of its internal event transport does not eliminate those Moby usages.

## Protocol Buffers

The Protobuf schemas are contracts. Generated `.pb.go` files are checked into the repository and ordinary builds consume them without invoking generators.

The current base branch retains the existing `protoc`-based regeneration procedure. Replacement of that procedure by Go-managed tooling is being reviewed independently in [the Protobuf integration PR](https://github.com/orqesta/moby/pull/1); it is not part of the verified mainline architecture yet.

## Validation boundary

A successful build or CI run demonstrates the corresponding compile and test scope. It does not establish that an Orqesta-built Docker Engine has passed end-to-end single-node or multi-node deployment tests. The newer SwarmKit event-transport changes must be integrated and validated separately.
