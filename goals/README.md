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
- 🟢 **Implemented & Verified**: 3
- 🟡 **Ready to Execute (Pending)**: 9 (6 independent ⚡, 3 unblocked 🔗)
- ⛔ **Blocked on Prerequisites**: 2
- 🎯 **By Tier**: 0 immediate (Tier 1), 14 roadmap (Tier 2), 0 nice-to-have (Tier 3), 0 wont-fix (Tier 4)
- 📋 **Execution Tasks Progress**: 13/65 completed (20%)

---

## 🟡 Ready to Execute (Pending Goals)

| Goal Package | Tier | Mode & Sequence | Dependencies | Focus & Description | Launch Command |
| :--- | :---: | :---: | :--- | :--- | :--- |
| [`fix-scunthorpe-and-blocked-list`](fix-scunthorpe-and-blocked-list/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Fix Scunthorpe Problem and Sanitize Blocked Word Filter — Replace naive substring search in isSafe() with token-boundary checks to prevent... | `/goal goals/fix-scunthorpe-and-blocked-list/goal.md` |
| [`fix-historical-scholar-attributions`](fix-historical-scholar-attributions/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Correct Historical and Religious Attribution of Non-Muslim Scholars — In words.json, rename muslim_scientists to islamic_golden_age_scholars or catego... | `/goal goals/fix-historical-scholar-attributions/goal.md` |
| [`fix-sanctuaries-and-animal-rules`](fix-sanctuaries-and-animal-rules/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Segregate Holy Sanctuaries and Fix Disrespectful Animal Pairings — Segregate holy Islamic sanctuaries (Kaaba, Masjid, Madina) from general places i... | `/goal goals/fix-sanctuaries-and-animal-rules/goal.md` |
| [`harden-d1-scripts-and-architecture`](harden-d1-scripts-and-architecture/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Harden D1 Population Script and Eliminate CI Blocking Prompts — Remove interactive input() prompts from populate_d1.py to fix CI/CD execution, a... | `/goal goals/harden-d1-scripts-and-architecture/goal.md` |
| [`harmonize-licensing-and-security-policy`](harmonize-licensing-and-security-policy/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Harmonize License Inconsistency, Update Security Policy and Docs — Clarify MIT vs Waqf GPL copyleft conflict, attribute WaqfTech in the license, re... | `/goal goals/harmonize-licensing-and-security-policy/goal.md` |
| [`sanitize-asma-allah-and-virtues`](sanitize-asma-allah-and-virtues/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⚡ Independent` | - | Segregate Asma' Allah al-Husna from Generic Adjectives — Segregate the 99 Names and Attributes of Allah from generic adjectives in words.... | `/goal goals/sanitize-asma-allah-and-virtues/goal.md` |
| [`fix-separator-and-word-count-logic`](fix-separator-and-word-count-logic/goal.md) | `🗺️ Roadmap` | `🚢 SHIP 🔗 Ready` | Deps met: fix-words-embedding-and-init | Fix Custom Separator Formatting and Word Count Boundaries — Ensure custom --sep applies to template joints and within compound words without... | `/goal goals/fix-separator-and-word-count-logic/goal.md` |
| [`fix-category-filtering-and-fallback`](fix-category-filtering-and-fallback/goal.md) | `🗺️ Roadmap` | `🚢 SHIP 🔗 Ready` | Deps met: fix-words-embedding-and-init | Fix Category Filtering and Add Single-Category Fallback Logic — Fix the rule matching logic in halal-random-strings.go so specifying any of the ... | `/goal goals/fix-category-filtering-and-fallback/goal.md` |
| [`fix-makefile-and-scripts-stability`](fix-makefile-and-scripts-stability/goal.md) | `🗺️ Roadmap` | `🚢 SHIP 🔗 Ready` | Deps met: clean-repo-artifacts-and-clutter | Fix Makefile Typo, Process Fork Bomb, and Script Zero-Division — Fix %(GO_APP_NAME) typo in Makefile, replace the 1,000-process bash loop with a ... | `/goal goals/fix-makefile-and-scripts-stability/goal.md` |
| [`fix-test-suite-and-panics`](fix-test-suite-and-panics/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⛔ Blocked` | Prereqs: fix-scunthorpe-and-blocked-list | Fix Test Suite Panics and Assertions — Fix the index out of range panic in TestIncludeRandomNumber, remove pass-by-valu... | `/goal goals/fix-test-suite-and-panics/goal.md` |
| [`optimize-generator-performance-and-allocations`](optimize-generator-performance-and-allocations/goal.md) | `🗺️ Roadmap` | `🚢 SHIP ⛔ Blocked` | Prereqs: fix-separator-and-word-count-logic | Eliminate Redundant Rule Filtering and Hot-Path Heap Allocations — Move rule filtering outside generation loops to eliminate 1M redundant slice all... | `/goal goals/optimize-generator-performance-and-allocations/goal.md` |

---

## 🟢 Implemented & Verified

| Goal Package | Mode | Shape | Task / Ref | Commit |
| :--- | :---: | :--- | :--- | :--- |
| [`fix-words-embedding-and-init`](fix-words-embedding-and-init/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `e38e9e0` |
| [`clean-repo-artifacts-and-clutter`](clean-repo-artifacts-and-clutter/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0bb7f9a` |
| [`sila-onboarding`](sila-onboarding/goal.md) | `⚡ Indep` | 🚢 SHIP | Git Commit | `0fd3c0d` |
