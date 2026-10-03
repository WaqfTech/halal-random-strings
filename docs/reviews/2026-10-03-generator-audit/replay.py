#!/usr/bin/env python3
"""Replay local audit probes. Results are observations, not safety certification."""

import argparse
import collections
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import time


def run(command, root, env, timeout=60):
    started = time.monotonic()
    completed = subprocess.run(
        command, cwd=root, env=env, capture_output=True, text=True, timeout=timeout
    )
    return {
        "command": list(map(str, command)),
        "exit": completed.returncode,
        "seconds": time.monotonic() - started,
        "stdout": completed.stdout,
        "stderr": completed.stderr,
    }


def corpus_metrics(text, categories):
    normalize = lambda s: s.lower().replace("'", "").replace(" ", "-")
    entity_categories = [
        "sahaba", "muslim_names_male", "muslim_names_female", "asma_allah",
        "holy_sanctuaries",
    ]
    food_categories = [
        "animals", "birds", "arabic_food", "jordanian_food", "saudi_food",
        "yemeni_food", "vegetables", "fruits", "spices",
    ]
    entities = {normalize(w) for c in entity_categories for w in categories[c]}
    food = {normalize(w) for c in food_categories for w in categories[c]}
    markers = {
        "shrine-holy", "zawiya-shrine", "khanqah-monastery", "mausoleum-tomb",
        "ablution-area-wash", "mihrab", "minbar", "prayer-niche", "qibla-wall",
        "adhan-tower", "minbar-stairs", "caliphate",
    }
    counts = collections.Counter(text.splitlines())
    findings = {"entity_food": [], "religious_marker_food": []}
    totals = collections.Counter()
    for line_number, line in enumerate(text.splitlines(), 1):
        tokens = line.split("-")
        phrases = {
            "-".join(tokens[i:j])
            for i in range(len(tokens)) for j in range(i + 1, len(tokens) + 1)
        }
        if not food & phrases:
            continue
        for key, vocabulary in (("entity_food", entities), ("religious_marker_food", markers)):
            if vocabulary & phrases:
                totals[key] += 1
                if len(findings[key]) < 12:
                    findings[key].append({"line": line_number, "text": line})
    return {
        "lines": sum(counts.values()), "unique": len(counts),
        "duplicate_excess": sum(counts.values()) - len(counts),
        "sha256": hashlib.sha256(text.encode()).hexdigest(),
        "candidate_counts": dict(totals), "candidate_samples": findings,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--go", default="go", help="Path to an installed Go executable")
    parser.add_argument("--out", type=Path, help="Receipt directory (default: a new temporary directory)")
    parser.add_argument("--samples", type=int, default=100000, help="Default corpus size per seed, 1..100000")
    parser.add_argument("--race", action="store_true", help="Also run the adversarial caller-mutation race probe")
    args = parser.parse_args()
    if not 1 <= args.samples <= 100000:
        parser.error("--samples must be between 1 and 100000")
    here = Path(__file__).resolve().parent
    root = here.parents[2]
    out = args.out.resolve() if args.out else Path(tempfile.mkdtemp(prefix="hrs-audit-replay-"))
    out.mkdir(parents=True, exist_ok=True)
    go = shutil.which(args.go)
    if go is None:
        parser.error("Go executable not found; specify --go")
    env = dict(os.environ)
    env["PATH"] = str(Path(go).parent) + os.pathsep + env.get("PATH", "")
    env["GOCACHE"] = str(out / "go-cache")
    cli = out / "generator"
    receipts = {"build": run([go, "build", "-o", str(cli), "./cmd/halal-random-strings"], root, env)}
    if receipts["build"]["exit"]:
        (out / "receipts.json").write_text(json.dumps(receipts, indent=2))
        raise SystemExit(receipts["build"]["stderr"])
    categories = json.loads((root / "words.json").read_text())["categories"]
    cases = [
        ("default-42", ["-r", str(args.samples), "--seed", "42", "--no-random-number"]),
        ("default-2026", ["-r", str(args.samples), "--seed", "2026", "--no-random-number"]),
        ("animals-places", ["--categories", "animals,nouns_places", "--seed", "42", "-r", "10000", "--no-random-number"]),
        ("animals-architecture", ["--categories", "animals,architectural_elements", "--seed", "42", "-r", "10000", "--no-random-number"]),
        ("asma-alone", ["--categories", "asma_allah", "--min-words", "2", "--max-words", "2", "--seed", "42", "-r", "20", "--no-random-number"]),
        ("rejected-name-food", ["--categories", "muslim_names_male,animals,arabic_food"]),
        ("donkey-separator", ["--categories", "muslim_names_male", "--sep", "donkey", "--min-words", "2", "--max-words", "2", "--seed", "42", "-r", "10", "--no-random-number"]),
        ("newline-separator", ["--categories", "muslim_names_male", "--sep", "\n", "--min-words", "2", "--max-words", "2", "--seed", "42", "--no-random-number"]),
        ("ignored-category", ["--categories", "colors_arabic,muslim_names_female", "--seed", "42", "-r", "10", "--no-random-number"]),
        ("inverted-bounds", ["--min-words", "9", "--max-words", "2", "--seed", "42"]),
        ("impossible-length", ["--categories", "servant_prefixes,asma_allah", "--min-words", "1", "--max-words", "1", "--seed", "42"]),
    ]
    receipts["cli"] = {}
    for label, options in cases:
        result = run([str(cli), *options], root, env)
        content = result.pop("stdout")
        (out / (label + ".txt")).write_text(content)
        result["samples"] = content.splitlines()[:20]
        if label.startswith("default-") or label.startswith("animals-"):
            result["corpus_metrics"] = corpus_metrics(content, categories)
        receipts["cli"][label] = result
    api_source = out / "api-probe.go"
    shutil.copyfile(here / "evidence" / "api-probe.go.txt", api_source)
    receipts["api"] = run([go, "run", str(api_source)], root, env)
    receipts["analyzer"] = []
    for value in ["ahmad-donkey-food", "ahmad-cow-rice", "ahmad_cow_rice", "abu-ubaidah-ibn-al-jarrah-cow", "sad-ibn-abi-waqqas-cow"]:
        fixture = out / "analyzer-input.txt"
        fixture.write_text(value + "\n")
        result = run(["python3", "scripts/analyze.py", str(fixture)], root, env)
        result["input"] = value
        receipts["analyzer"].append(result)
    if args.race:
        race_source = out / "race-probe.go"
        shutil.copyfile(here / "evidence" / "race-probe.go.txt", race_source)
        receipts["race"] = run([go, "run", "-race", str(race_source)], root, env)
    (out / "receipts.json").write_text(json.dumps(receipts, indent=2) + "\n")
    print("Audit receipts:", out / "receipts.json")
    print("Inspect the observations; this replay does not certify sensitivity compliance.")


if __name__ == "__main__":
    main()
