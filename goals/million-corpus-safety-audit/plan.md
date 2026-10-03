# Plan: Generate and audit at least one million real identifiers

## Execution Steps
- [x] Build the current CLI and generate two fresh corpora of 1,000,000 actual identifiers each: default numbered and unnumbered.
- [x] Run the production analyzer against every physical line and retain hashes, exact counts, commands, exits, distribution and duplicate metrics.
- [x] Independently parse every identifier against the explicit approved group table, test requested sensitivity probes, and verify injected violations at the start, middle and end of a million-line corpus.
- [x] Inspect actual group samples, record findings and limitations, verify repository checks, and commit the evidence with `(goals/million-corpus-safety-audit/goal.md)` attribution.

## Result

Two million actual generated identifiers passed both complete mechanical audits.
The million-line control correctly failed with exactly three inserted violations.
Semantic review requests changes: actual `amr-ibn-al-as-jinn`,
`abdur-rahman-ibn-abu-bakr-makruh`, and `prophet-said-ibn-zayd-illumination-qadr`
are allowed by the current Islamic grouping. A full follow-up scan found 3,267
candidate lines for review. Dictionary annotations/classification and repeated
identifiers also need attention. Audit delivery is complete; it does not certify
the broader sensitivity goal. See the durable report and receipts in
`docs/reviews/2026-10-03-million-corpus/`.
