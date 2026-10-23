# Requirements

Mission: Keep the kubectl library in staging/src/k8s.io/kubectl correct and well tested: fix bugs and add focused unit tests without changing CLI output.

Non-goals:

- Changes outside staging/src/k8s.io/kubectl unless the lead approves
- API or feature-gate changes
- Release or vendoring work

Each requirement has a stable ID. Never renumber an ID; mark it `Withdrawn` instead.

## REQ-001: Focused unit tests for staging/src/k8s.io/kubectl



### Acceptance criteria

- [ ] New behavior has a table-driven unit test
- [ ] Package tests pass

### Work items

- Add tests for slice.ContainsString with and without a modifier
