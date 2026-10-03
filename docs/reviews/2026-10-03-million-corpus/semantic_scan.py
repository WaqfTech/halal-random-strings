#!/usr/bin/env python3
"""Locate named-person/concept associations for manual review, not certification."""

import argparse
import collections
import hashlib
import json
from pathlib import Path
import re

from independent_scan import ROOT, canonical

CUES = ["jinn", "creator", "prophet", "messenger", "prophethood", "makruh"]


def boundary_pattern(values):
    return re.compile(r"(?<![a-z])(?:" + "|".join(re.escape(s) for s in
                      sorted(values, key=len, reverse=True)) + r")(?![a-z])")


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--corpus-dir", type=Path, required=True)
    ap.add_argument("--receipts", type=Path, default=Path(__file__).with_name("semantic-candidates.json"))
    args = ap.parse_args()
    manifest = json.loads((args.corpus_dir / "audited.json").read_text())
    data = (ROOT / "words.json").read_bytes()
    assert hashlib.sha256(data).hexdigest() == manifest["dictionary_sha256"], "Dictionary changed"
    words = json.loads(data)
    people = {canonical(value): category for category in ("sahaba", "islamic_golden_age_scholars")
              for value in words["categories"][category]}
    people_re = boundary_pattern(people)
    cue_re = boundary_pattern(CUES)
    # A token inside an atomic event such as 'Birth of Prophet Muhammad' is
    # not a standalone 'prophet' descriptor. Exclude those readings here.
    containers = {canonical(value) for entries in words["categories"].values() for value in entries
                  if "-" in canonical(value) and cue_re.search(canonical(value))}
    container_re = boundary_pattern(containers) if containers else None
    results = []
    for corpus in manifest["corpora"]:
        counts = collections.Counter()
        samples = {cue: [] for cue in CUES}
        name_annotation_count = unique_candidate_lines = total = 0
        digest = hashlib.sha256()
        with Path(corpus["path"]).open("rb") as source:
            for line, raw in enumerate(source, 1):
                digest.update(raw)
                total += 1
                value = raw.decode("ascii").rstrip("\n")
                if "-yumna-success-" in "-" + value + "-":
                    name_annotation_count += 1
                spans = [match.span() for match in container_re.finditer(value)] if container_re else []
                found = {match.group() for match in cue_re.finditer(value)
                         if not any(start <= match.start() and match.end() <= end for start, end in spans)}
                if not found:
                    continue
                entities = [{"value": match.group(), "category": people[match.group()]}
                            for match in people_re.finditer(value)]
                if not entities:
                    continue
                unique_candidate_lines += 1
                for cue in found:
                    counts[cue] += 1
                    if len(samples[cue]) < 8:
                        samples[cue].append({"line": line, "value": value, "people": entities})
        assert total == corpus["physical_lines"] and digest.hexdigest() == corpus["sha256"], "Corpus changed"
        result = {"name": corpus["name"], "sha256": digest.hexdigest(), "physical_lines": total,
                  "unique_candidate_lines": unique_candidate_lines, "cue_person_lines": dict(sorted(counts.items())),
                  "examples": samples, "yumna_success_lines": name_annotation_count,
                  "excluded_atomic_cue_containers": sorted(containers),
                  "interpretation": "Heuristic candidates, not automatic policy violations or religious judgments. "
                  "Whole-token dictionary person expressions occur with standalone concept tokens. "
                  "Different cue counts can overlap on the same line."}
        results.append(result)
        print(f"{corpus['name']}: {unique_candidate_lines} distinct candidate lines; "
              f"cue counts {dict(sorted(counts.items()))}; Yumna-Success {name_annotation_count}", flush=True)
    args.receipts.write_text(json.dumps(results, indent=2) + "\n")


if __name__ == "__main__":
    main()
