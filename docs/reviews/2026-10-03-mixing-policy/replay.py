#!/usr/bin/env python3
"""Replay local CLI policy checks; save receipts without retaining large corpora."""

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import time

ROOT = Path(__file__).resolve().parents[3]
sys.path.insert(0, str(ROOT / "scripts"))
from naming_policy import NamingPolicy


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--cli", type=Path, default=ROOT / "halal-random-strings")
    parser.add_argument("--receipts", type=Path, default=Path(__file__).with_name("runtime.json"))
    args = parser.parse_args()
    cli = str(args.cli.resolve())
    policy = NamingPolicy()
    receipts = {"dictionary_sha256": hashlib.sha256((ROOT / "words.json").read_bytes()).hexdigest(),
                "cli_sha256": hashlib.sha256(Path(cli).read_bytes()).hexdigest(),
                "cli_cases": [], "corpora": [], "analyzer_regressions": []}
    def invoke(arguments):
        started = time.monotonic()
        result = subprocess.run(arguments, cwd=ROOT, capture_output=True, text=True, timeout=60)
        return {"arguments": arguments, "exit": result.returncode,
                "elapsed_seconds": round(time.monotonic()-started, 6),
                "stdout": result.stdout, "stderr": result.stderr}
    positives = ["asma_allah", "prophets", "muslim_names_male,muslim_names_female",
                 "sahaba,islamic_months", "islamic_golden_age_scholars,islamic_inventions",
                 "animals,architectural_elements", "colors_arabic,nouns_places",
                 "arabic_food,spices", "birds,flowers"]
    for categories in positives:
        receipt = invoke([cli, "--seed", "42", "-r", "8", "--categories", categories,
                          "--min-words", "2", "--max-words", "12", "--no-random-number"])
        assert receipt["exit"] == 0, receipt
        values = receipt["stdout"].splitlines()
        assert len(values) == 8
        for value in values:
            policy.validate(value)
        receipts["cli_cases"].append(receipt)
    negatives = ["asma_allah,prophets", "asma_allah,islamic_virtues", "asma_allah,muslim_names_male",
                 "servant_prefixes,asma_allah", "prophets,sahaba", "prophets,animals,arabic_food",
                 "muslim_names_male,animals", "muslim_names_female,islamic_virtues",
                 "sahaba,nouns_places", "islamic_inventions,animals", "adjectives,nouns_places",
                 "nouns_concepts,fruits", "muslim_empires,animals"]
    for categories in negatives:
        receipt = invoke([cli, "--categories", categories])
        assert receipt["exit"] == 1 and not receipt["stdout"], receipt
        receipts["cli_cases"].append(receipt)
    for sep in ["donkey", "wine", "x", "\n", " ", "--"]:
        receipt = invoke([cli, "--sep", sep, "--categories", "prophets"])
        assert receipt["exit"] == 1 and not receipt["stdout"], receipt
        receipts["cli_cases"].append(receipt)
    receipt = invoke([cli, "-r", "1000000", "--categories", "sahaba", "--min-words", "1", "--max-words", "1"])
    assert receipt["exit"] == 1 and not receipt["stdout"], receipt
    receipts["cli_cases"].append(receipt)
    with tempfile.TemporaryDirectory(prefix="hrs-policy-replay-") as directory:
        cases = [("default_numbered", 200000, "-", [], True),
                 ("default_no_number", 100000, "-", [], False),
                 ("animals_architecture", 10000, "-", ["animals","architectural_elements"], False)]
        cases += [("default_sep_" + str(i), 10000, sep, [], True) for i, sep in enumerate(["_", ".", "/"])]
        for name, count, sep, categories, numbered in cases:
            path = Path(directory) / (name + ".txt")
            command = [cli, "--seed", "424242", "-r", str(count), "--sep", sep]
            if categories:
                command += ["--categories", ",".join(categories)]
            if not numbered:
                command += ["--no-random-number"]
            started = time.monotonic()
            with path.open("w") as output:
                result = subprocess.run(command, cwd=ROOT, stdout=output, stderr=subprocess.PIPE, text=True, timeout=60)
            assert result.returncode == 0, result.stderr
            generation_seconds = round(time.monotonic()-started, 6)
            data = path.read_bytes()
            assert len(data.splitlines()) == count
            audit = invoke([sys.executable, "scripts/analyze.py", str(path), "--sep", sep])
            assert audit["exit"] == 0, audit
            receipts["corpora"].append({"name":name, "generation_arguments":command,
                                       "generation_seconds":generation_seconds, "count":count,
                                       "sha256":hashlib.sha256(data).hexdigest(),
                                       "samples":data.decode().splitlines()[:12], "audit":audit})
        bad_cases = ["ahmad-donkey-food", "ahmad-cow-olive", "aisha-ramadan", "rahman-ramadan",
                     "mihrab-cow", "sad-ibn-abi-waqqas-cow", "hamza-ibn-abd-al-muttalib-olive",
                     "shrine-holy-dwelling-green-peppercorn-falafel", "kaffir-lime-leaves", "smores", ""]
        for index, value in enumerate(bad_cases):
            path = Path(directory) / f"bad-{index}.txt"
            path.write_text(value + "\n")
            receipt = invoke([sys.executable, "scripts/analyze.py", str(path)])
            assert receipt["exit"] == 1, receipt
            receipts["analyzer_regressions"].append({"input":value, **receipt})
        path = Path(directory) / "bad-import.txt"
        path.write_text("olive-1234\nahmad-donkey-food\n")
        dump = Path(directory) / "dump.sql"
        dump.write_text("sentinel")
        receipt = invoke([sys.executable, "scripts/populate_d1.py", str(path), "--dump-sql", str(dump)])
        assert receipt["exit"] == 1 and dump.read_text() == "sentinel", receipt
        receipts["import_rejection"] = receipt
    args.receipts.write_text(json.dumps(receipts, indent=2) + "\n")
    print(f"PASS: {sum(c['count'] for c in receipts['corpora']):,} generated identifiers audited; "
          f"{len(receipts['cli_cases'])} CLI cases; {len(bad_cases)} failing analyzer fixtures; unsafe import rejected.")


if __name__ == "__main__":
    main()
