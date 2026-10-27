# Invariants

Every change MUST keep these rules. The reviewer checks each one and cites the rule number in a FAIL verdict.

1. Every new or modified .go file in staging/src/k8s.io/apimachinery/pkg/util/version carries the Apache-2.0 license header.
2. Zero edits to go.mod, go.sum, go.work, vendor/, or files outside staging/src/k8s.io/apimachinery/pkg/util/version (except staging/src/k8s.io/apimachinery/pkg/util/version/knowledge-plane/).
3. `go build . && go vet . && go test -race -count=1 .` exits 0 inside staging/src/k8s.io/apimachinery/pkg/util/version in under 30 seconds.
4. Existing exported signatures in k8s.io/apimachinery/pkg/util/version remain backward-compatible.
