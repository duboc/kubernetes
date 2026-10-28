# Requirements

Mission: Harden k8s.io/apimachinery/pkg/util/version with nil-safe ordering predicates, a Core() semantic version helper, and round-trip/antisymmetry fuzz tests without breaking existing Kubernetes callers.

Non-goals:

- Editing generated go.mod, go.work, vendor, or zz_generated.* files
- Changing packages outside staging/src/k8s.io/apimachinery/pkg/util/version
- Deploying any HTTP server or Kubernetes cluster workload as an application service

Each requirement has a stable ID. Never renumber an ID; mark it `Withdrawn` instead.

## REQ-001: Nil-safe ordering methods and Core() semantic release helper



### Acceptance criteria

- [x] (*Version).AtLeast, (*Version).LessThan, and (*Version).GreaterThan never panic when receiver or argument is nil, ordering nil strictly before any non-nil *Version and equal to nil
- [x] (*Version).Core() returns a semantic Major.Minor.Patch *Version copy stripped of preRelease and buildMetadata (or nil when receiver is nil)
- [x] Table-driven unit tests in version_test.go cover nil receivers, nil arguments, generic versions, and pre-release/build-metadata stripping

### Work items

- Wave 1.1: Make AtLeast, LessThan, and GreaterThan nil-safe and add (*Version).Core() helper

## REQ-002: Round-trip and antisymmetry Go fuzz test suite



### Acceptance criteria

- [x] version_fuzz_test.go carries the Kubernetes Apache-2.0 header and defines FuzzParseAndCompare(f *testing.F)
- [x] go test -fuzz=FuzzParseAndCompare -fuzztime=3s . passes with zero panics and verifies ParseSemantic(v.String()).EqualTo(v) plus Compare antisymmetry

### Work items

- Wave 2.1: Add FuzzParseAndCompare round-trip and ordering antisymmetry fuzz suite
