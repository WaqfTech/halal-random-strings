#!/usr/bin/env python3
"""Generate two million real identifiers, then run mechanical and semantic audits.

Mechanical passes check the reviewed list policy. Semantic candidate output still
requires review; this command's success is not a cultural-safety certification.
"""

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[3]
HERE = Path(__file__).resolve().parent


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--cli", type=Path, required=True)
    ap.add_argument("--corpus-dir", type=Path, required=True)
    ap.add_argument("--receipts", type=Path, default=HERE / "runtime.json")
    ap.add_argument("--semantic-receipts", type=Path, default=HERE / "semantic-candidates.json")
    args = ap.parse_args()
    cli = str(args.cli.resolve())
    directory = args.corpus_dir.resolve()
    directory.mkdir(parents=True, exist_ok=True)
    receipts = {"started_utc": datetime.now(timezone.utc).isoformat(),
                "source_commit": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
                "dictionary_sha256": hashlib.sha256((ROOT / "words.json").read_bytes()).hexdigest(),
                "cli_sha256": hashlib.sha256(Path(cli).read_bytes()).hexdigest(), "corpora": []}
    # Complete generation of both batches before beginning any audit.
    for name, seed, numbered in [("default-numbered", 1003202601, True),
                                 ("default-unnumbered", 1003202602, False)]:
        path = directory / (name + ".txt")
        arguments = [cli, "--seed", str(seed), "-r", "1000000"]
        if not numbered:
            arguments.append("--no-random-number")
        started = time.monotonic()
        with path.open("w") as output:
            result = subprocess.run(arguments, cwd=ROOT, stdout=output, stderr=subprocess.PIPE, text=True, timeout=60)
        assert result.returncode == 0, result.stderr
        elapsed = time.monotonic() - started
        digest = hashlib.sha256()
        count = size = 0
        with path.open("rb") as source:
            for raw in source:
                digest.update(raw)
                count += 1
                size += len(raw)
        assert count == 1000000, count
        receipts["corpora"].append({"name": name, "path": str(path), "numbered": numbered, "seed": seed,
                                    "generation_arguments": arguments, "generation_exit": result.returncode,
                                    "generation_stderr": result.stderr, "generation_seconds": round(elapsed, 6),
                                    "physical_lines": count, "bytes": size, "sha256": digest.hexdigest()})
        print(f"GENERATED {name}: {count:,} physical lines, sha256 {digest.hexdigest()}", flush=True)
    receipts["finished_utc"] = datetime.now(timezone.utc).isoformat()
    (directory / "generation.json").write_text(json.dumps(receipts, indent=2) + "\n")
    for item in receipts["corpora"]:
        arguments = [sys.executable, "scripts/analyze.py", item["path"], "--sep", "-"]
        started = time.monotonic()
        result = subprocess.run(arguments, cwd=ROOT, capture_output=True, text=True, timeout=120)
        item["audit"] = {"arguments": arguments, "exit": result.returncode,
                         "elapsed_seconds": round(time.monotonic() - started, 6),
                         "stdout": result.stdout, "stderr": result.stderr}
        (directory / (item["name"] + "-audit.txt")).write_text(result.stdout + result.stderr)
        assert result.returncode == 0, item["audit"]
        print(f"MECHANICAL PASS {item['name']}: every line audited", flush=True)
    (directory / "audited.json").write_text(json.dumps(receipts, indent=2) + "\n")
    subprocess.run([sys.executable, str(HERE / "independent_scan.py"), "--corpus-dir", str(directory),
                    "--receipts", str(args.receipts.resolve())], cwd=ROOT, check=True, timeout=120)
    subprocess.run([sys.executable, str(HERE / "semantic_scan.py"), "--corpus-dir", str(directory),
                    "--receipts", str(args.semantic_receipts.resolve())], cwd=ROOT, check=True, timeout=120)
    print("REVIEW REQUIRED: semantic-candidates.json contains potential sensitive associations.", flush=True)


if __name__ == "__main__":
    main()
