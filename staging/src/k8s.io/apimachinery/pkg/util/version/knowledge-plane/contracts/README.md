# Contracts

Contracts are the source of truth for the L2 contract checks that the evaluation team runs. Change a contract in the
same work item as the code that implements it.

## Contract files

- [`staging/src/k8s.io/apimachinery/pkg/util/version/version.go`](../../../../../../../../../staging/src/k8s.io/apimachinery/pkg/util/version/version.go)

## Endpoints under contract

- None recorded.

## Contract invariants

- Existing exported symbols in k8s.io/apimachinery/pkg/util/version remain backward-compatible
- LessThan(other) == !AtLeast(other) holds across all nil and non-nil *Version pairs
