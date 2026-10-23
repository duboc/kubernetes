# Golden data (evaluation team only)

This directory belongs to the evaluation project (`kubectl-eval`). Copy it into the evaluation
workspace, never into the build workspace or the build knowledge plane, so the build team can't tune the product to
the test set.

- **Metric:** `package_test_pass_rate` with a target of `1.0` or higher.
- **Probe:** go test k8s.io/kubectl/pkg/util/slice

Store each dataset as a versioned file (for example, `golden-v1.jsonl`), and record the version in every evaluation
run. The evaluation team reads the build team's knowledge plane read-only at `/scion-volumes/knowledge-plane`.
