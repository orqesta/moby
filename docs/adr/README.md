# Architecture decision records

This directory records architectural decisions implemented in the Orqesta Moby fork. Each decision is immutable once accepted; later changes receive a new ADR.

- [ADR-0001 — Pin the Orqesta SwarmKit fork](ADR-0001-orqesta-swarmkit-dependency.md)
- [ADR-0002 — Use native Go Modules without a vendored dependency tree](ADR-0002-native-go-modules.md)

These records describe verified changes. The replacement of the Protobuf compiler/generator (commit `9e6c35c`) is under validation and is not an accepted decision here. The Moby fork still pins an older SwarmKit revision; updating that revision is a separate pending integration.

The cross-repository decision log is maintained in [Orqesta's Google Drive ADR register](https://drive.google.com/drive/folders/1Y-ErdX024KOy4zicZ1G35YSvYLUvTyno).
