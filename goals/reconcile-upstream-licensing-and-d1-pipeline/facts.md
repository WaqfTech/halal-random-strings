# Facts: Reconcile Upstream Licensing Attribution and Harden D1 Pipeline

## Architectural Invariants & Constraints
- Legal attribution compliance: Derivative works of MIT-licensed software must preserve the copyright notice of the original authors.
- Honest branding: A fork maintained under a separate organization cannot use upstream badges or claims of official inclusion without authorization.
- Cloudflare D1 operational limits: Wrangler CLI commands run as distinct Node.js runtimes; batching through `--file` with SQLite transactions is mandatory for large-scale ingestion.

## File & Interface Contracts
- `LICENSE` / `NOTICE`: Upstream attribution text.
- `README.md`: Upstream credit and branding sections.
- `scripts/populate_d1.py`: `populate_d1_from_output`, batch execution.
- `scripts/analyze.py`: `analyze_uniqueness`.
