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
- 🟢 **Implemented & Verified**: 7
- 🟡 **Ready to Execute (Pending)**: 7 (5 independent ⚡, 2 unblocked 🔗)
- 🎯 **By Tier**: 0 immediate (Tier 1), 14 roadmap (Tier 2), 0 nice-to-have (Tier 3), 0 wont-fix (Tier 4)
- 📋 **Execution Tasks Progress**: 33/65 completed (50%)

---

## 🟡 Ready to Execute (Pending Goals)

| Goal Package | Tier | Mode & Sequence | Dependencies | Focus & Description | Launch Command |
| :--- | :---: | :---: | :--- | :--- | :--- |
| [`fix-historical-scholar-attributions`](fix-historical-scholar-attributions/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Correct Historical and Religious Attribution of Non-Muslim Scholars — In words.json, rename muslim_scientists to islamic_golden_age_scholars or catego... | `/goal goals/fix-historical-scholar-attributions/goal.md` |
| [`fix-sanctuaries-and-animal-rules`](fix-sanctuaries-and-animal-rules/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Segregate Holy Sanctuaries and Fix Disrespectful Animal Pairings — Segregate holy Islamic sanctuaries (Kaaba, Masjid, Madina) from general places i... | `/goal goals/fix-sanctuaries-and-animal-rules/goal.md` |
| [`harden-d1-scripts-and-architecture`](harden-d1-scripts-and-architecture/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Harden D1 Population Script and Eliminate CI Blocking Prompts — Remove interactive input() prompts from populate_d1.py to fix CI/CD execution, a... | `/goal goals/harden-d1-scripts-and-architecture/goal.md` |
| [`harmonize-licensing-and-security-policy`](harmonize-licensing-and-security-policy/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Harmonize License Inconsistency, Update Security Policy and Docs — Clarify MIT vs Waqf GPL copyleft conflict, attribute WaqfTech in the license, re... | `/goal goals/harmonize-licensing-and-security-policy/goal.md` |
| [`sanitize-asma-allah-and-virtues`](sanitize-asma-allah-and-virtues/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Segregate Asma' Allah al-Husna from Generic Adjectives — Segregate the 99 Names and Attributes of Allah from generic adjectives in words.... | `/goal goals/sanitize-asma-allah-and-virtues/goal.md` |
| [`fix-makefile-and-scripts-stability`](fix-makefile-and-scripts-stability/goal.md) | `🗺️ Roadmap` | `🚢 SHIP 🔗 Ready` | Deps met: clean-repo-artifacts-and-clutter | Fix Makefile Typo, Process Fork Bomb, and Script Zero-Division — Fix %(GO_APP_NAME) typo in Makefile, replace the 1,000-process bash loop with a ... | `/goal goals/fix-makefile-and-scripts-stability/goal.md` |
| [`optimize-generator-performance-and-allocations`](optimize-generator-performance-and-allocations/goal.md) | `🗺️ Roadmap` | `🚢 SHIP 🔗 Ready` | Deps met: fix-words-embedding-and-init, fix-separator-and-word-count-logic | Eliminate Redundant Rule Filtering and Hot-Path Heap Allocations — Move rule filtering outside generation loops to eliminate 1M redundant slice all... | `/goal goals/optimize-generator-performance-and-allocations/goal.md` |

---

## 🟢 Implemented & Verified

| Goal Package | Mode | Shape | Task / Ref | Commit |
| :--- | :---: | :--- | :--- | :--- |
| [`fix-words-embedding-and-init`](fix-words-embedding-and-init/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `e38e9e0` |
| [`clean-repo-artifacts-and-clutter`](clean-repo-artifacts-and-clutter/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0bb7f9a` |
| [`fix-scunthorpe-and-blocked-list`](fix-scunthorpe-and-blocked-list/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `7fc7218` |
| [`sila-onboarding`](sila-onboarding/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0fd3c0d` |
| [`fix-separator-and-word-count-logic`](fix-separator-and-word-count-logic/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `e732521` |
| [`fix-category-filtering-and-fallback`](fix-category-filtering-and-fallback/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `e664af4` |
| [`fix-test-suite-and-panics`](fix-test-suite-and-panics/goal.md) | `🔗 Dep` | 🚢 SHIP | Git Commit | `509943b` |
