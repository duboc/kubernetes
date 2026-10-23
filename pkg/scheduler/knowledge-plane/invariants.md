# Invariants

Every change MUST keep these rules. The reviewer checks each one and cites the rule number in a FAIL verdict.

1. Work stays inside pkg/scheduler unless the lead approves.
2. Every change runs go test for the touched package before review.
