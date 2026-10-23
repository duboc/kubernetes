# Requirements

Mission: Keep the kube-scheduler code in pkg/scheduler correct and well tested: fix bugs and add focused unit tests without changing scheduling decisions.

Non-goals:

- Changes outside pkg/scheduler unless the lead approves
- API or feature-gate changes
- Release or vendoring work

Each requirement has a stable ID. Never renumber an ID; mark it `Withdrawn` instead.

## REQ-001: Focused unit tests for pkg/scheduler



### Acceptance criteria

- [ ] New behavior has a table-driven unit test
- [ ] Package tests pass

### Work items

- Add unit tests for util.IsScalarResourceName covering extended, hugepages, and native resources
