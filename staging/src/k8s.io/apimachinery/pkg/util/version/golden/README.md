# Golden data (evaluation team only)

This directory belongs to the evaluation project (`k8s-version-eval`). Copy it into the evaluation
workspace, never into the build workspace or the build knowledge plane, so the build team can't tune the product to
the test set.

- **Metric:** `version_statement_coverage` with a target of `0.92` or higher.
- **Probe:** Measure `go test -cover .` statement ratio (>= 0.92) and run clean-room ordering/Core()/fuzz property checks (1.0 pass rate) against origin/factory/gke-main.

Store each dataset as a versioned file (for example, `golden-v1.jsonl`), and record the version in every evaluation
run. The evaluation team reads the build team's knowledge plane read-only at `/scion-volumes/knowledge-plane`.
