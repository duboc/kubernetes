# ADR-001: Order nil *Version strictly before non-nil *Version across all comparison predicates

- **Status:** Accepted
- **Requirements:** REQ-001
- **Work items:** #8 (commit d12a0b1c)
- **Migrations:** —

## Context

EqualTo(other) already handles nil receiver and nil argument safely, whereas AtLeast, LessThan, and GreaterThan dereference v.components directly in compareInternal and panic on nil.

## Decision

Guard nil receiver and nil argument in compareInternal (or in AtLeast/LessThan/GreaterThan) so nil == nil (0), nil < non-nil (-1), and non-nil > nil (+1), preserving LessThan == !AtLeast.

## Consequences

Callers comparing optional parsed *Version pointers cannot panic; existing non-nil behavior is 100% unchanged.
