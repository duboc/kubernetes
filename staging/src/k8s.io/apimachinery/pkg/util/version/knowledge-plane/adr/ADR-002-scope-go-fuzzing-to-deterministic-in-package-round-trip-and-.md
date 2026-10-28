# ADR-002: Scope Go fuzzing to deterministic in-package round-trip and ordering laws

- **Status:** Accepted
- **Requirements:** REQ-002
- **Work items:** #9 (commit c0bc927d)
- **Migrations:** —

## Context

Version parsing uses regular expressions and numeric splitting that must never panic or violate total-order antisymmetry on arbitrary inputs.

## Decision

Implement native Go testing.F fuzzing in version_fuzz_test.go seeded from TestSemanticVersions and TestGenericVersions corpuses, bounded to 3s in automated CI/eval checks.

## Consequences

Zero external dependencies required; runs offline against the vendored tree.
