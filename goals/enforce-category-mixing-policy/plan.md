# Plan: Enforce reviewed category mixing policy

## Execution Steps

- [x] **Policy and data**: Define explicit groups, isolate prophet spellings and titles, disable standalone prefixes, and curate cross-group aliases and religious general-list entries.
- [x] **Implementation**: Disable custom dictionaries, use immutable indexes and feasible category/length plans, enforce safe separators and bounds, compile blocked phrases into a trie, and validate audit/export input.
- [x] **Verification**: Run the full race suite, vet, build, Python tests, all category pairs and CLI corpus replay; preserve receipts and policy/API migration documentation.
- [ ] **Attribution commit**: Commit with `(goals/enforce-category-mixing-policy/goal.md)`, run `sila goals` to verify attribution, and save the handoff.
