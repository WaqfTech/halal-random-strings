#!/usr/bin/env python3
"""Report uniqueness and enforce the reviewed identifier policy; FAIL exits 1."""

import argparse
import collections
import sys

from naming_policy import NamingPolicy, corpus_lines


def analyze_corpus(filename, sep="-"):
    policy = NamingPolicy()  # Missing/invalid dictionary must fail closed.
    counts = collections.Counter()
    word_counts = collections.Counter()
    groups = collections.Counter()
    examples = []
    total = violations = char_total = 0
    shortest = longest = None
    for number, line in corpus_lines(filename):
        counts[line] += 1
        total += 1
        char_total += len(line)
        shortest = len(line) if shortest is None else min(shortest, len(line))
        longest = len(line) if longest is None else max(longest, len(line))
        try:
            group, words = policy.validate(line, sep)
            groups[group] += 1
            word_counts[words] += 1
        except ValueError as error:
            violations += 1
            if len(examples) < 10:
                examples.append((number, line, str(error)))
    print(f"Total lines: {total}")
    print(f"Unique lines: {len(counts)}")
    if not total:
        print("Uniqueness percentage: 0.00% (file is empty)")
        print("[FAIL] Empty corpus cannot prove safe generation.")
        return False
    print(f"Uniqueness percentage: {len(counts) / total * 100:.2f}%")
    duplicates = sorted(((line, count) for line, count in counts.items() if count > 1),
                        key=lambda item: (-item[1], item[0]))
    if duplicates:
        print("\nDuplicate lines:")
        for line, count in duplicates[:20]:
            print(f"{line} (x{count})")
        if len(duplicates) > 20:
            print(f"... and {len(duplicates) - 20} more duplicate entries.")
    else:
        print("\nNo duplicate lines found.")
    print(f"Length metrics: min={shortest} chars, max={longest} chars, avg={char_total / total:.1f} chars")
    print("Word count distribution (valid identifiers):")
    for words, count in sorted(word_counts.items()):
        print(f"  {words} words: {count} lines")
    print("Mixing groups:")
    for group, count in sorted(groups.items()):
        print(f"  {group}: {count}")
    print(f"[{ 'FAIL' if violations else 'PASS' }] Policy violations: {violations}")
    for number, line, reason in examples:
        print(f"  Line {number}: {line!r} ({reason})")
    return violations == 0


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("output_file")
    parser.add_argument("--sep", default="-", choices=["-", "_", ".", "/"])
    args = parser.parse_args()
    try:
        return 0 if analyze_corpus(args.output_file, args.sep) else 1
    except (OSError, ValueError, KeyError, TypeError) as error:
        print(f"[FAIL] Cannot audit corpus: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
