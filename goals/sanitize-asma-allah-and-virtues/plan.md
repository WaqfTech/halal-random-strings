# Plan: Segregate Asma' Allah al-Husna from Generic Adjectives

## Execution Steps
- [x] **Phase 1 — Extract Divine Names**: Extract all exclusive Names and Attributes of Allah (`rahman`, `khaliq`, `bari`, `musawwir`, `ghaffar`, `qahhar`, `wahhab`, `razzaq`, `fattah`, `alim`, `ahad`, `samad`, `qadir`, `hayy`, `qayyum`, `malik-al-mulk`, `dhul-jalali-wal-ikram`, `mumit`, `muhyi`, etc.) out of the `adjectives` array in `words.json`.
- [x] **Phase 2 — Isolated Divine Names Category**: Place them into a dedicated category (e.g. `divine_attributes` or `asma_allah`).
- [x] **Phase 3 — Audit Combination Rules**: Ensure Rules 42 (`{animals}-{adjectives}`), 53 (`{vegetables}-{adjectives}`), 55 (`{arabic_food}-{adjectives}`), and 56 (`{jordanian_food}-{adjectives}`) only draw from mundane descriptive adjectives and cannot pair with Asma' Allah.
- [x] **Phase 4 — Respectful Formulas**: Define dedicated rules for Divine Names, such as `abd-{divine_name}` or within pious compound phrases.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/sanitize-asma-allah-and-virtues/goal.md)` using a `fix:` subject, then run `sila goals`.
