#!/usr/bin/env python3
"""
analyze.py - Stream-based uniqueness analyzer for generated string outputs.
"""

import collections
import os
import sys

def stream_lines(filename):
    """Yield non-empty stripped lines from file."""
    with open(filename, 'r', encoding='utf-8') as f:
        for line in f:
            stripped = line.strip()
            if stripped:
                yield stripped

def analyze_uniqueness(filename):
    if not os.path.exists(filename):
        print(f"Error: File '{filename}' not found.")
        sys.exit(1)

    counts = collections.Counter()
    total_lines = 0

    for line in stream_lines(filename):
        counts[line] += 1
        total_lines += 1

    unique_lines = len(counts)

    print(f"Total lines: {total_lines}")
    print(f"Unique lines: {unique_lines}")

    if total_lines == 0:
        print("Uniqueness percentage: 0.00% (file is empty)")
        print("\nNo lines to analyze.")
        return

    uniqueness_percentage = (unique_lines / total_lines) * 100
    print(f"Uniqueness percentage: {uniqueness_percentage:.2f}%")

    duplicates = [(line, count) for line, count in counts.items() if count > 1]
    if duplicates:
        print("\nDuplicate lines:")
        duplicates.sort(key=lambda x: (-x[1], x[0]))
        for line, count in duplicates:
            print(f"{line} (x{count})")
    else:
        print("\nNo duplicate lines found.")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 analyze.py <output_file>")
        sys.exit(1)
    analyze_uniqueness(sys.argv[1])