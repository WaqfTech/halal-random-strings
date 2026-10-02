# Goals Status & Pragmatic Execution Roadmap

This index tracks all goal packages under `goals/`, their execution mode (Independent ⚡ vs Dependent 🔗), dependencies, and progress status.

## 🤖 Agent Guide: How to Create, Claim & Execute Goals

1. **Create a Goal**: Scaffold a new goal package conforming to standards:
   ```bash
   sila goals create <slug> --title="..." [--depends-on="..."] [--independent]
   ```
2. **Claim a Goal**: Lock the goal so no other agent duplicates work:
   ```bash
   sila goals claim <slug> --agent=<name> [--note="evaluating..."]
   ```
3. **Launch & Implement**: Follow `facts.md` and `plan.md` until all tests pass:
   ```bash
   /goal goals/<slug>/goal.md
   ```
4. **Progress Updates**: Report status updates during execution:
   ```bash
   sila goals report <slug> "Running e2e test suite"
   ```
5. **Commit with Goal Reference**: Include the goal path in your commit message:
   ```bash
   git commit -m "feat(scope): implement description (goals/<slug>/goal.md)"
   ```
6. **Auto-Reconciliation**: Run `sila goals` to scan commits, mark 🟢 **Implemented**, and clear the claim lock.

## Summary
- **Total Goals**: 14
- 🟢 **Implemented & Verified**: 14
- 🟡 **Ready to Execute (Pending)**: 0 (0 independent ⚡, 0 unblocked 🔗)
- 🎯 **By Tier**: 0 immediate (Tier 1), 14 roadmap (Tier 2), 0 nice-to-have (Tier 3), 0 wont-fix (Tier 4)
- 📋 **Execution Tasks Progress**: 65/65 completed (100%)

---

## 🟢 Implemented & Verified

| Goal Package | Mode | Shape | Task / Ref | Commit |
| :--- | :---: | :--- | :--- | :--- |
| [`fix-words-embedding-and-init`](fix-words-embedding-and-init/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `e38e9e0` |
| [`clean-repo-artifacts-and-clutter`](clean-repo-artifacts-and-clutter/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0bb7f9a` |
| [`fix-scunthorpe-and-blocked-list`](fix-scunthorpe-and-blocked-list/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `7fc7218` |
| [`fix-historical-scholar-attributions`](fix-historical-scholar-attributions/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `ba5c235` |
| [`fix-sanctuaries-and-animal-rules`](fix-sanctuaries-and-animal-rules/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `7431c7e` |
| [`harden-d1-scripts-and-architecture`](harden-d1-scripts-and-architecture/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `dfd58b1` |
| [`harmonize-licensing-and-security-policy`](harmonize-licensing-and-security-policy/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `9626920` |
| [`sanitize-asma-allah-and-virtues`](sanitize-asma-allah-and-virtues/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `fddf653` |
| [`sila-onboarding`](sila-onboarding/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0fd3c0d` |
| [`fix-separator-and-word-count-logic`](fix-separator-and-word-count-logic/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `e732521` |
| [`fix-category-filtering-and-fallback`](fix-category-filtering-and-fallback/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `e664af4` |
| [`fix-makefile-and-scripts-stability`](fix-makefile-and-scripts-stability/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `317035e` |
| [`fix-test-suite-and-panics`](fix-test-suite-and-panics/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `509943b` |
| [`optimize-generator-performance-and-allocations`](optimize-generator-performance-and-allocations/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `65c7117` |
