# Kubernetes Version Utility Factory knowledge plane

Read this index before you plan, implement, review, or deploy any work item. Open files through the paths in this
index; never guess a file name (L4). Cite requirement IDs (`REQ-NNN`) and decision IDs (`ADR-NNN`) in work items,
commits, reviews, and evaluation results.

## Where to look

| When you need to | Read |
|---|---|
| Know what to build and how it's accepted | [product/requirements.md](product/requirements.md) |
| Know why the system is built this way | [Decisions](#decisions) |
| Check an API or data contract | [contracts/README.md](contracts/README.md) |
| Review a change | [invariants.md](invariants.md) |
| Understand a domain term | [glossary.md](glossary.md) |
| Follow a project-specific procedure | [playbooks/README.md](playbooks/README.md) |
| Know what is live right now | [state/LIVE.md](state/LIVE.md) |
| Find the backlog | GitHub issues in `duboc/kubernetes` (`status:*` labels) |

## Learned sources

These files existed in the repository before the factory started. They stay where they are; this index links to them.

- Contracts: [`staging/src/k8s.io/apimachinery/pkg/util/version/version.go`](../../../../../../../../staging/src/k8s.io/apimachinery/pkg/util/version/version.go)
- Guides: [`README.md`](../../../../../../../../README.md)
- Guides: [`CONTRIBUTING.md`](../../../../../../../../CONTRIBUTING.md)
- Guides: [`AGENTS.md`](../../../../../../../../AGENTS.md)

## Decisions

| ADR | Title | Status | Requirements |
|---|---|---|---|
| [ADR-001](adr/ADR-001-order-nil-version-strictly-before-non-nil-version-across-all.md) | Order nil *Version strictly before non-nil *Version across all comparison predicates | Accepted | REQ-001 |
| [ADR-002](adr/ADR-002-scope-go-fuzzing-to-deterministic-in-package-round-trip-and-.md) | Scope Go fuzzing to deterministic in-package round-trip and ordering laws | Accepted | REQ-002 |

## Keep this knowledge plane healthy

- Run `python3 knowledge-plane/kp-check.py` after you change any file here. It fails on broken links and unknown
  `REQ-`/`ADR-` IDs.
- Add a decision with [adr/_template.md](adr/_template.md), and add a row to the table above.
- Don't put live state in this file. The doc-writer keeps [state/LIVE.md](state/LIVE.md) current (nothing is deployed).
- Don't put evaluation golden datasets here. They belong to the evaluation team.
