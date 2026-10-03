#!/usr/bin/env python3
"""Audit retained million-line runs with a separate trie/bitmask parser.

This evidence helper does not import the production Python NamingPolicy. Its
category table is an explicit copy of the user's approved contract. Vocabulary
still comes from words.json; this cannot certify missing cultural knowledge.
"""

import argparse
import collections
import hashlib
import json
from pathlib import Path
import random
import re
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[3]
GROUPS = {
    "divine": ["asma_allah"],
    "prophets": ["prophets"],
    "names": ["muslim_names_male", "muslim_names_female"],
    "islamic": ["adjectives", "sahaba", "islamic_art_forms", "islamic_events",
                "islamic_golden_age_scholars", "islamic_inventions", "islamic_months",
                "islamic_virtues", "scholarly_terms", "adab_terms", "holy_sanctuaries",
                "muslim_empires", "nouns_concepts"],
    "general": ["animals", "arabic_food", "architectural_elements", "birds",
                "colors_arabic", "days_of_week_arabic", "flowers", "fruits",
                "gems_minerals", "geographic_features", "jordanian_food",
                "nouns_objects", "nouns_places", "saudi_food", "spices", "suffixes",
                "trees", "vegetables", "yemeni_food"],
}
PROBES = ["ahmad-donkey-food", "ahmad-cow-olive", "aisha-cow-olive",
          "aisha-ramadan", "rahman-ramadan", "mihrab-cow",
          "sad-ibn-abi-waqqas-cow", "hamza-ibn-abd-al-muttalib-olive",
          "kaffir", "kaffer", "kafir", "caffre"]
PROBE_PATTERN = re.compile(r"(?<![a-z])(?:" + "|".join(map(re.escape, PROBES)) + r")(?![a-z])")


def canonical(value):
    return re.sub(r"[^a-z0-9]+", "-", value.lower().strip()
                  .replace("'", "").replace("\u2019", "")).strip("-")


def insert(trie, phrase, mask):
    node = trie
    for token in phrase.split("-"):
        node = node.setdefault(token, {})
    node[None] = node.get(None, 0) | mask


class IndependentParser:
    def __init__(self):
        data = json.loads((ROOT / "words.json").read_text())
        expected = {category: group for group, categories in GROUPS.items() for category in categories}
        assert expected == data["policy"]["category_groups"], "Contract table drift"
        assert set(expected) == set(data["categories"]), "Unexpected categories"
        self.names = list(GROUPS)
        self.trie = {}
        self.blocked = {}
        self.expressions = {}
        for index, (group, categories) in enumerate(GROUPS.items()):
            mask = 1 << index
            for category in categories:
                for raw in data["categories"][category]:
                    phrase = canonical(raw)
                    insert(self.trie, phrase, mask)
                    self.expressions[phrase] = group
        for raw in data["blocked"]:
            insert(self.blocked, canonical(raw), 1)

    def blocked_matches(self, tokens):
        for start in range(len(tokens)):
            node = self.blocked
            for end in range(start, len(tokens)):
                node = node.get(tokens[end])
                if node is None:
                    break
                if None in node:
                    yield "-".join(tokens[start:end + 1])

    def group_mask(self, tokens):
        # Each position carries all group hypotheses that reached it. A complete
        # expression can extend only hypotheses with the same group bit.
        reach = [0] * (len(tokens) + 1)
        reach[0] = (1 << len(self.names)) - 1
        for start in range(len(tokens)):
            if not reach[start]:
                continue
            node = self.trie
            for end in range(start, len(tokens)):
                node = node.get(tokens[end])
                if node is None:
                    break
                reach[end + 1] |= reach[start] & node.get(None, 0)
        return reach[-1]


def raw_lines(path):
    with path.open("rb") as source:
        yield from source


def scan(parser, item):
    started = time.monotonic()
    path = Path(item["path"])
    counts = collections.Counter()
    groups = collections.Counter()
    word_counts = collections.Counter()
    single_expressions = collections.Counter()
    probes = collections.Counter()
    samples = {group: [] for group in GROUPS}
    extremes = {}
    rng = random.Random(item["seed"])
    digest = hashlib.sha256()
    failures = []
    violation_count = ambiguous_count = blocked_count = total_bytes = total = 0
    for number, raw in enumerate(raw_lines(path), 1):
        digest.update(raw)
        total_bytes += len(raw)
        total += 1
        value = raw.decode("ascii").removesuffix("\n")
        counts[value] += 1
        tokens = value.split("-")
        reason = None
        if not raw.endswith(b"\n") or "\r" in value:
            reason = "Noncanonical physical line"
        if item["numbered"]:
            if not re.fullmatch(r"[1-9][0-9]{3}", tokens[-1]):
                reason = "Missing/invalid four-digit suffix"
            else:
                tokens.pop()
        if not 5 <= len(tokens) <= 8 or any(not re.fullmatch(r"[a-z]+", t) for t in tokens):
            reason = "Wrong format or requested 5-8 token bounds"
        blocked = list(parser.blocked_matches(tokens))
        if blocked:
            blocked_count += 1
            reason = "Blocked phrase: " + ",".join(blocked)
        mask = parser.group_mask(tokens)
        if mask == 0:
            reason = "Cannot cover whole identifier with one approved group"
        elif mask & (mask - 1):
            ambiguous_count += 1
            reason = "Multiple group readings"
        for match in PROBE_PATTERN.finditer(value):
            probes[match.group()] += 1
            reason = "Explicit prohibited sensitivity probe"
        if reason:
            violation_count += 1
            if len(failures) < 10:
                failures.append({"line": number, "value": value, "reason": reason})
            continue
        group = parser.names[mask.bit_length() - 1]
        groups[group] += 1
        word_counts[len(tokens)] += 1
        if "-".join(tokens) in parser.expressions:
            single_expressions[group] += 1
        sample = {"line": number, "value": value}
        if len(samples[group]) < 20:
            samples[group].append(sample)
        else:
            slot = rng.randrange(groups[group])
            if slot < 20:
                samples[group][slot] = sample
        if group not in extremes:
            extremes[group] = {"shortest": sample, "longest": sample}
        else:
            for key, better in [("shortest", len(value) < len(extremes[group]["shortest"]["value"])),
                                ("longest", len(value) > len(extremes[group]["longest"]["value"]))]:
                if better:
                    extremes[group][key] = sample
    result = {
        "name": item["name"], "physical_lines": total, "bytes": total_bytes,
        "sha256": digest.hexdigest(), "unique_lines": len(counts),
        "duplicate_excess": total - len(counts),
        "duplicate_identifiers": sum(count > 1 for count in counts.values()),
        "max_occurrences": max(counts.values()),
        "top_repeated": [{"value": value, "count": count} for value, count in counts.most_common(20) if count > 1],
        "group_counts": dict(sorted(groups.items())), "word_counts": dict(sorted(word_counts.items())),
        "matches_one_complete_expression": dict(sorted(single_expressions.items())),
        "blocked_lines": blocked_count, "ambiguous_group_lines": ambiguous_count,
        "sensitivity_probe_hits": dict(probes), "violations": violation_count,
        "violation_examples": failures,
        "reservoir_samples": {group: sorted(items, key=lambda s: s["line"]) for group, items in samples.items()},
        "length_extremes": extremes, "elapsed_seconds": round(time.monotonic() - started, 6),
    }
    assert total == item["physical_lines"] == 1000000, result
    assert total_bytes == item["bytes"] and result["sha256"] == item["sha256"], "Corpus changed"
    assert not violation_count, result
    assert set(groups) == set(GROUPS), "Missing group coverage"
    return result


def negative_control(directory, source):
    path = directory / "injected-negative-control.txt"
    injections = {1: "ahmad-donkey-food-1234", 500000: "aisha-cow-olive-2345",
                  1000000: "rahman-ramadan-3456"}
    digest = hashlib.sha256()
    with Path(source).open("rb") as original, path.open("wb") as output:
        for number, raw in enumerate(original, 1):
            raw = (injections[number] + "\n").encode() if number in injections else raw
            digest.update(raw)
            output.write(raw)
    assert number == 1000000
    arguments = [sys.executable, "scripts/analyze.py", str(path), "--sep", "-"]
    started = time.monotonic()
    result = subprocess.run(arguments, cwd=ROOT, capture_output=True, text=True, timeout=120)
    receipt = {"path": str(path), "sha256": digest.hexdigest(), "physical_lines": number,
               "injections": injections, "arguments": arguments, "exit": result.returncode,
               "stdout": result.stdout, "stderr": result.stderr,
               "elapsed_seconds": round(time.monotonic() - started, 6)}
    assert result.returncode == 1 and "[FAIL] Policy violations: 3" in result.stdout, receipt
    for line, value in injections.items():
        assert f"Line {line}: {value!r}" in result.stdout, receipt
    (directory / "negative-control-audit.txt").write_text(result.stdout + result.stderr)
    return receipt


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--corpus-dir", type=Path, required=True)
    ap.add_argument("--receipts", type=Path, default=Path(__file__).with_name("runtime.json"))
    args = ap.parse_args()
    receipts = json.loads((args.corpus_dir / "audited.json").read_text())
    assert receipts["dictionary_sha256"] == hashlib.sha256((ROOT / "words.json").read_bytes()).hexdigest()
    parser = IndependentParser()
    receipts["independent_parser"] = "explicit category table, shared token trie, group hypothesis bitmasks"
    receipts["independent_scans"] = []
    for item in receipts["corpora"]:
        assert item["generation_exit"] == item["audit"]["exit"] == 0
        result = scan(parser, item)
        receipts["independent_scans"].append(result)
        print(f"PASS independent {item['name']}: {result['physical_lines']:,} lines, "
              f"{result['violations']} violations; {result['duplicate_excess']} duplicate excess", flush=True)
    # These controls check the independent parser as well as the CLI auditor.
    receipts["independent_negative_probes"] = []
    for value in PROBES:
        tokens = value.split("-")
        blocked = list(parser.blocked_matches(tokens))
        mask = parser.group_mask(tokens)
        assert blocked or mask == 0, value
        receipts["independent_negative_probes"].append({"value": value, "group_mask": mask, "blocked": blocked})
    receipts["negative_control"] = negative_control(args.corpus_dir, receipts["corpora"][0]["path"])
    print("PASS control: all 3 injected violations detected at lines 1, 500000, 1000000; exit 1", flush=True)
    args.receipts.write_text(json.dumps(receipts, indent=2) + "\n")


if __name__ == "__main__":
    main()
