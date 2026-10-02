#!/usr/bin/env python3
"""
analyze.py - Stream-based uniqueness & semantic sensitivity analyzer for generated string outputs.
"""

import collections
import json
import os
import re
import sys

def stream_lines(filename):
    """Yield non-empty stripped lines from file."""
    with open(filename, 'r', encoding='utf-8') as f:
        for line in f:
            stripped = line.strip()
            if stripped:
                yield stripped

def load_word_categories():
    """Load and classify categories from words.json if present."""
    base_dir = os.path.dirname(os.path.abspath(__file__))
    words_path = os.path.join(base_dir, '..', 'words.json')
    if not os.path.exists(words_path):
        return None

    try:
        with open(words_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
    except Exception:
        return None

    cats = data.get('categories', {})

    sacred_cat_names = {
        'sahaba', 'muslim_names_male', 'muslim_names_female',
        'asma_allah', 'servant_prefixes', 'holy_sanctuaries',
        'scholarly_terms', 'islamic_events', 'adab_terms',
        'islamic_virtues', 'islamic_art_forms', 'islamic_golden_age_scholars',
        'days_of_week_arabic', 'islamic_months', 'nouns_concepts', 'adjectives',
        'muslim_empires'
    }

    mundane_cat_names = {
        'animals', 'birds', 'arabic_food', 'jordanian_food',
        'saudi_food', 'yemeni_food', 'vegetables', 'fruits',
        'spices', 'trees', 'nouns_objects'
    }

    sacred_phrases = set()
    for sc in sacred_cat_names:
        for w in cats.get(sc, []):
            norm = re.sub(r'[^a-z0-9]+', '-', w.lower()).strip('-')
            if norm:
                sacred_phrases.add(norm)

    mundane_phrases = set()
    for mc in mundane_cat_names:
        for w in cats.get(mc, []):
            norm = re.sub(r'[^a-z0-9]+', '-', w.lower()).strip('-')
            if norm:
                mundane_phrases.add(norm)

    # Specific highly sensitive entities for explicit checks
    holy_entities = set()
    for w in (cats.get('holy_sanctuaries', []) + cats.get('asma_allah', []) + 
              cats.get('sahaba', []) + cats.get('muslim_names_male', []) + 
              cats.get('muslim_names_female', [])):
        norm = re.sub(r'[^a-z0-9]+', '-', w.lower()).strip('-')
        if norm:
            holy_entities.add(norm)

    mundane_derogatory = set()
    for w in (cats.get('animals', []) + cats.get('birds', []) + 
              cats.get('arabic_food', []) + cats.get('jordanian_food', []) + 
              cats.get('saudi_food', []) + cats.get('yemeni_food', []) + 
              cats.get('vegetables', []) + cats.get('fruits', [])):
        norm = re.sub(r'[^a-z0-9]+', '-', w.lower()).strip('-')
        if norm:
            mundane_derogatory.add(norm)

    return {
        'sacred': sacred_phrases,
        'mundane': mundane_phrases,
        'holy_entities': holy_entities,
        'mundane_derogatory': mundane_derogatory,
    }

def analyze_corpus(filename):
    if not os.path.exists(filename):
        print(f"Error: File '{filename}' not found.")
        sys.exit(1)

    counts = collections.Counter()
    total_lines = 0
    word_counts = collections.Counter()
    char_lengths = []

    cat_data = load_word_categories()
    sacred_mundane_violations = []
    holy_mundane_violations = []
    scraped_artifacts = []

    artifact_patterns = [
        re.compile(r'\b(judgment-day-concept|badr-event|caliphate-system|rahman-attribute|fajr-dawn|arafat-day|ruh-cleanliness)\b'),
        re.compile(r'\b(haram-area-sanctuary|prayer-room-retreat|mussalla-prayer|jameh-mosque-great)\b'),
        re.compile(r'\b(peace-salam|justice-adl|favor-ni-ma|courage-shaja-a)\b'),
        re.compile(r'\bchrist-s-thorn\b'),
    ]

    for line_num, line in enumerate(stream_lines(filename), 1):
        counts[line] += 1
        total_lines += 1

        tokens = line.split('-')
        # Exclude trailing 4-digit number for word count if present
        actual_tokens = tokens[:-1] if (len(tokens) > 1 and tokens[-1].isdigit()) else tokens
        word_counts[len(actual_tokens)] += 1
        char_lengths.append(len(line))

        # Check for scraped glossary artifacts
        for pat in artifact_patterns:
            if pat.search(line):
                if len(scraped_artifacts) < 10:
                    scraped_artifacts.append((line_num, line))
                break

        # Semantic domain checks
        if cat_data:
            ngrams = set()
            for i in range(len(actual_tokens)):
                for j in range(i + 1, min(len(actual_tokens) + 1, i + 5)):
                    ngrams.add("-".join(actual_tokens[i:j]))

            has_holy = ngrams & cat_data['holy_entities']
            has_mundane_derog = ngrams & cat_data['mundane_derogatory']

            if has_holy and has_mundane_derog:
                if len(holy_mundane_violations) < 10:
                    holy_mundane_violations.append((line_num, line, list(has_holy), list(has_mundane_derog)))

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
        # Show top duplicates (up to 20)
        for line, count in duplicates[:20]:
            print(f"{line} (x{count})")
        if len(duplicates) > 20:
            print(f"... and {len(duplicates) - 20} more duplicate entries.")
    else:
        print("\nNo duplicate lines found.")

    # Semantic & Metric Report
    print("\n" + "=" * 60)
    print("           SEMANTIC SENSITIVITY & AUDIT REPORT")
    print("=" * 60)
    print(f"Length metrics: min={min(char_lengths)} chars, max={max(char_lengths)} chars, avg={sum(char_lengths)/len(char_lengths):.1f} chars")
    print("Word count distribution:")
    for wc in sorted(word_counts.keys()):
        print(f"  {wc} words: {word_counts[wc]} lines ({word_counts[wc]/total_lines*100:.2f}%)")

    violations_found = False

    if holy_mundane_violations:
        violations_found = True
        print(f"\n[FAIL] Found {len(holy_mundane_violations)} Holy/Sacred entity + Mundane pairings:")
        for lnum, lstr, holy, mund in holy_mundane_violations:
            print(f"  Line {lnum}: {lstr} (Holy: {holy}, Mundane: {mund})")
    else:
        print("\n[PASS] Zero Holy entity (Qur'an, Prophets, Sahaba, Asma Allah) paired with food/animals.")

    if scraped_artifacts:
        violations_found = True
        print(f"\n[FAIL] Found {len(scraped_artifacts)} scraped glossary artifacts:")
        for lnum, lstr in scraped_artifacts:
            print(f"  Line {lnum}: {lstr}")
    else:
        print("[PASS] Zero scraped glossary artifacts detected.")

    print("=" * 60)

    if violations_found:
        print("\nWarning: Sensitivity violations detected in output corpus.")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 analyze.py <output_file>")
        sys.exit(1)
    analyze_corpus(sys.argv[1])