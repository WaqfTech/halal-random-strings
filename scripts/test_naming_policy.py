"""Executable policy and offline import regression tests (no Cloudflare calls)."""

import contextlib
import io
import json
from pathlib import Path
import sqlite3
import subprocess
import sys
import tempfile
import unittest

from analyze import analyze_corpus
from naming_policy import NamingPolicy, normalize
from populate_d1 import generate_sql_file


class NamingPolicyTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.policy = NamingPolicy()

    def test_reviewed_examples_and_separators(self):
        good = {"aisha-maher": "names", "abdul-hadi-maher": "names",
                "ahmad-musa": "prophets", "rahman-rahim": "divine",
                "mihrab-ramadan": "islamic", "sad-ibn-abi-waqqas-ramadan": "islamic",
                "olive-apple": "general", "makrut-lime-leaves-olive": "general", "bird-of-paradise-olive": "general"}
        bad = ["taha-ramadan", "muzammil-ramadan", "kaffir-lime-leaves", "ahmad-donkey-food", "ahmad-cow-rice", "aisha-ramadan",
               "muhammad-aisha", "rahman-ramadan", "abd-rahman", "mihrab-cow",
               "sad-ibn-abi-waqqas-cow", "hamza-ibn-abd-al-muttalib-rice",
               "shrine-holy-dwelling-green-peppercorn-falafel", "s-mores", "smores",
               "olive-123", "olive-0123", "olive-12345", "olive--apple",
               " olive", "Olive", "olive\napple", "unknownword", ""]
        for sep in self.policy.separators:
            for value, group in good.items():
                with self.subTest(value=value, sep=sep):
                    self.assertEqual(self.policy.validate(value.replace("-", sep), sep)[0], group)
                    self.assertEqual(self.policy.validate(value.replace("-", sep) + sep + "1234", sep)[0], group)
            for value in bad:
                with self.subTest(value=value, sep=sep):
                    with self.assertRaises(ValueError):
                        self.policy.validate(value.replace("-", sep), sep)
        for sep in ("donkey", "wine", "x", "\n", "--"):
            with self.assertRaises(ValueError):
                self.policy.validate("ahmad" + sep + "cow", sep)

    def test_normalization_and_long_blocked_phrases(self):
        self.assertEqual(normalize("Sa'd ibn Abi Waqqas"), "sad-ibn-abi-waqqas")
        self.assertEqual(normalize("s'mores"), "smores")
        policy = NamingPolicy()
        policy.blocked.add("one-two-three-four-five")
        policy.max_blocked = 5
        self.assertTrue(policy.has_blocked("prefix-one-two-three-four-five-suffix".split("-")))

    def test_missing_or_invalid_dictionary_fails_closed(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "words.json"
            with self.assertRaises(OSError):
                NamingPolicy(path)
            path.write_text("{}")
            with self.assertRaises(KeyError):
                NamingPolicy(path)
            original = json.loads((Path(__file__).resolve().parent.parent / "words.json").read_text())
            original["policy"]["category_groups"]["asma_allah"] = "general"
            path.write_text(json.dumps(original))
            with self.assertRaises(ValueError):
                NamingPolicy(path)

    def test_analyzer_exit_counts_and_empty_input(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "corpus.txt"
            path.write_text("olive-1234\n" + "ahmad-donkey-food\n" * 13 + "\n")
            output = io.StringIO()
            with contextlib.redirect_stdout(output):
                self.assertFalse(analyze_corpus(path))
            self.assertIn("Policy violations: 14", output.getvalue())
            self.assertEqual(output.getvalue().count("  Line "), 10)
            result = subprocess.run([sys.executable, str(Path(__file__).with_name("analyze.py")), str(path)], capture_output=True)
            self.assertEqual(result.returncode, 1)
            path.write_text("")
            with contextlib.redirect_stdout(io.StringIO()):
                self.assertFalse(analyze_corpus(path))
            path.write_text("olive_apple_1234\n")
            with contextlib.redirect_stdout(io.StringIO()):
                self.assertFalse(analyze_corpus(path))
                self.assertTrue(analyze_corpus(path, "_"))

    def test_import_rejection_preserves_existing_dump(self):
        with tempfile.TemporaryDirectory() as directory:
            source, dump = Path(directory) / "corpus.txt", Path(directory) / "dump.sql"
            source.write_text("olive-1234\nahmad-donkey-food\n")
            dump.write_text("existing valid dump")
            with self.assertRaises(ValueError):
                generate_sql_file(source, dump, chunk_size=1)
            self.assertEqual(dump.read_text(), "existing valid dump")
            self.assertEqual(sorted(p.name for p in Path(directory).iterdir()), ["corpus.txt", "dump.sql"])
            result = subprocess.run([sys.executable, str(Path(__file__).with_name("populate_d1.py")), str(source), "--dump-sql", str(dump)], capture_output=True)
            self.assertEqual(result.returncode, 1)
            self.assertEqual(dump.read_text(), "existing valid dump")
            source.write_text("olive-1234\n\n")
            with self.assertRaises(ValueError):
                generate_sql_file(source, dump)
            for size in (0, -1, 501):
                with self.assertRaises(ValueError):
                    generate_sql_file(source, dump, size)

    def test_valid_sql_import_and_deduplication(self):
        with tempfile.TemporaryDirectory() as directory:
            source, dump = Path(directory) / "corpus.txt", Path(directory) / "dump.sql"
            source.write_text("olive_1234\nolive_1234\naisha_maher_5678\n")
            self.assertEqual(generate_sql_file(source, dump, 1, "_"), 3)
            sql = dump.read_text()
            self.assertNotIn("BEGIN TRANSACTION", sql)
            self.assertNotIn("COMMIT", sql)
            with sqlite3.connect(":memory:") as db:
                schema = Path(__file__).resolve().parent.parent / "db" / "schema.sql"
                db.executescript(schema.read_text())
                db.executescript(sql)
                self.assertEqual(db.execute("SELECT id FROM generated_strings ORDER BY id").fetchall(), [("aisha_maher_5678",), ("olive_1234",)])


if __name__ == "__main__":
    unittest.main()
