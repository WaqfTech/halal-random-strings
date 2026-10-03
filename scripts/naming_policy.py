"""Validation of identifiers against the reviewed words.json mixing policy."""

import json
from pathlib import Path
import re


def normalize(value):
    value = value.lower().strip().replace("'", "").replace("’", "")
    return re.sub(r"[^a-z0-9]+", "-", value).strip("-")


class NamingPolicy:
    def __init__(self, words_path=None):
        path = Path(words_path) if words_path else Path(__file__).resolve().parent.parent / "words.json"
        data = json.loads(path.read_text(encoding="utf-8"))
        policy = data["policy"]
        categories = data["categories"]
        groups = policy["category_groups"]
        if (policy["version"] != 1 or policy["max_words"] != 64
                or policy["max_repeat"] != 1000000
                or policy["allowed_separators"] != ["-", "_", ".", "/"]
                or set(groups) != set(categories) or not data["blocked"]):
            raise ValueError("Missing or unsupported mixing policy")
        protected = {"asma_allah": "divine", "prophets": "prophets",
                     "muslim_names_male": "names", "muslim_names_female": "names"}
        if any(groups.get(c) != g for c, g in protected.items()):
            raise ValueError("Protected category assigned to wrong group")
        self.separators = policy["allowed_separators"]
        self.max_words = policy["max_words"]
        self.blocked = {normalize(s) for s in data["blocked"]}
        if "" in self.blocked:
            raise ValueError("Empty blocked expression")
        self.max_blocked = max(len(s.split("-")) for s in self.blocked)
        self.phrases = {g: set() for g in ("divine", "prophets", "names", "islamic", "general")}
        owners = {}
        for category, entries in categories.items():
            group = groups[category]
            if (group not in self.phrases or not entries
                    or category in policy["disabled_categories"]
                    or group in ("divine", "prophets", "names") and protected.get(category) != group):
                raise ValueError(f"Invalid category {category}")
            seen = set()
            for word in entries:
                canonical = normalize(word)
                if (not word.isascii() or not re.fullmatch(r"[a-z]+(?:-[a-z]+)*", canonical)
                        or canonical in seen or len(canonical.split("-")) > self.max_words
                        or self.has_blocked(canonical.split("-"))):
                    raise ValueError(f"Invalid, duplicate or blocked entry {word} in {category}")
                if canonical in owners and owners[canonical] != group:
                    raise ValueError(f"Ambiguous expression {canonical}")
                seen.add(canonical)
                owners[canonical] = group
                self.phrases[group].add(canonical)
        if any(not phrases for phrases in self.phrases.values()):
            raise ValueError("Missing mixing group")
        self.max_phrase = {g: max(len(s.split("-")) for s in phrases)
                           for g, phrases in self.phrases.items()}
        for rule in data["rules"]:
            pattern = rule["pattern"]
            if (not pattern or any(c not in groups for c in pattern)
                    or len({groups[c] for c in pattern}) != 1
                    or rule["template"] != "-".join("{" + c + "}" for c in pattern)):
                raise ValueError(f"Invalid rule {rule}")

    def has_blocked(self, tokens):
        return any("-".join(tokens[i:j]) in self.blocked
                   for i in range(len(tokens))
                   for j in range(i + 1, min(len(tokens), i + self.max_blocked) + 1))

    def validate(self, value, sep="-"):
        """Return (group, word count), or raise ValueError. Never guess a separator."""
        if sep not in self.separators:
            raise ValueError(f"Unsupported separator {sep!r}")
        tokens = value.split(sep)
        if len(tokens) > 1 and re.fullmatch(r"[1-9][0-9]{3}", tokens[-1]):
            tokens.pop()
        if not 1 <= len(tokens) <= self.max_words or any(not re.fullmatch(r"[a-z]+", s) for s in tokens):
            raise ValueError("Invalid identifier format or word count")
        if self.has_blocked(tokens):
            raise ValueError("Blocked vocabulary")
        # Match complete expressions, including long names and compounds. Words
        # inside 'abdul-hadi' or the plant 'bird-of-paradise' are not separate entries.
        for group, phrases in self.phrases.items():
            reachable = [False] * (len(tokens) + 1)
            reachable[0] = True
            for i in range(len(tokens)):
                if not reachable[i]:
                    continue
                for j in range(i + 1, min(len(tokens), i + self.max_phrase[group]) + 1):
                    if "-".join(tokens[i:j]) in phrases:
                        reachable[j] = True
            if reachable[-1]:
                return group, len(tokens)
        raise ValueError("Unknown words or incompatible mixing groups")


def corpus_lines(filename):
    """Preserve whitespace as invalid input and keep physical line numbers."""
    with open(filename, encoding="utf-8") as source:
        for number, line in enumerate(source, 1):
            yield number, line.rstrip("\n").removesuffix("\r")
