#!/usr/bin/env python3
"""
populate_d1.py - Populate Cloudflare D1 Database from generated strings.

Architectural Guidance:
For high-traffic production deployments, consider on-the-fly generation inside
the Cloudflare Worker (e.g. bundling the dictionary or compiling to WebAssembly)
to eliminate D1 read-after-write concurrency races and billing overhead (avoiding
frequent database writes/reads per generated string). Pre-populating D1 is best
suited for pre-allocated claim/reservation token pools.
"""

import argparse
import os
import subprocess
import sys
import tempfile
import time

def run_command(command, check=True, capture_output=True, text=True):
    try:
        result = subprocess.run(command, check=check, capture_output=capture_output, text=text)
        return result
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {' '.join(e.cmd)}")
        if e.stdout:
            print(f"Stdout: {e.stdout}")
        if e.stderr:
            print(f"Stderr: {e.stderr}")
        if check:
            raise
        return None
    except FileNotFoundError:
        print("Error: 'wrangler' command not found. Please ensure Cloudflare Wrangler CLI is installed and in your PATH.")
        raise

def prompt_yes_no(question, default=False, interactive=False):
    if not interactive or not sys.stdin.isatty():
        return default
    try:
        response = input(question).strip().lower()
        return response in ("y", "yes")
    except (EOFError, KeyboardInterrupt):
        return default

def escape_sql_string(val: str) -> str:
    """Sanitize and escape string literal for SQLite/D1."""
    sanitized = val.replace("\x00", "").replace("\r", "").replace("\n", "")
    return sanitized.replace("'", "''")

def format_insert_block(batch, timestamp):
    values = [f"('{escape_sql_string(item)}', NULL, NULL, {timestamp})" for item in batch]
    return (
        f"INSERT OR IGNORE INTO generated_strings (id, used_at, used_by, created_at) VALUES\n"
        f"  {',\n  '.join(values)};\n"
    )

def generate_sql_file(output_file, sql_file_path, chunk_size=500):
    """Stream lines from output_file and write transactional SQL batches."""
    total_strings = 0
    timestamp = int(time.time())

    with open(output_file, 'r', encoding='utf-8') as fin, open(sql_file_path, 'w', encoding='utf-8') as fout:
        fout.write("-- Cloudflare D1 batch insert script\n")
        fout.write("BEGIN TRANSACTION;\n")
        
        batch = []
        for line in fin:
            cleaned = line.strip()
            if not cleaned:
                continue
            batch.append(cleaned)
            total_strings += 1

            if len(batch) >= chunk_size:
                fout.write(format_insert_block(batch, timestamp))
                batch = []
                if total_strings % 5000 == 0:
                    fout.write("COMMIT;\nBEGIN TRANSACTION;\n")

        if batch:
            fout.write(format_insert_block(batch, timestamp))

        fout.write("COMMIT;\n")

    return total_strings

def check_wrangler_setup(db_name, auto_create=False, auto_apply_schema=False, interactive=False):
    print("\n--- Checking Wrangler Setup ---")
    
    # 1. Check if wrangler is installed
    print("Checking if Wrangler CLI is installed...")
    if run_command(["wrangler", "--version"], check=False) is None:
        print("Wrangler not found. Please install it globally: npm install -g wrangler")
        return False
    print("Wrangler CLI is installed.")

    # 2. Check if authenticated
    print("Checking Wrangler authentication status...")
    auth_check = run_command(["wrangler", "whoami"], check=False)
    if auth_check is None or "Logged in" not in auth_check.stdout:
        print("Wrangler is not authenticated. Please run 'wrangler login' and follow the prompts.")
        return False
    print("Wrangler is authenticated.")

    # 3. Check if D1 database exists
    print(f"Checking if D1 database '{db_name}' exists...")
    db_list = run_command(["wrangler", "d1", "list"], check=True)
    if db_list and db_name not in db_list.stdout:
        print(f"D1 database '{db_name}' not found.")
        should_create = auto_create or prompt_yes_no(
            f"Do you want to create '{db_name}' now? (y/N): ",
            default=False,
            interactive=interactive,
        )
        if should_create:
            print(f"Creating D1 database '{db_name}'...")
            create_db_result = run_command(["wrangler", "d1", "create", db_name], check=True)
            if create_db_result:
                print(f"Database '{db_name}' created successfully.")
            else:
                print(f"Failed to create database '{db_name}'. Please check your Cloudflare account and permissions.")
                return False
        else:
            print("Database not created. Exiting.")
            return False
    print(f"D1 database '{db_name}' exists.")

    # 4. Check if schema is applied
    print("Checking if schema is applied...")
    schema_file = "db/schema.sql"
    if not os.path.exists(schema_file):
        print(f"Error: Schema file '{schema_file}' not found. Please ensure it's in the correct location.")
        return False

    table_check = run_command(
        ["wrangler", "d1", "execute", db_name, "--command", "PRAGMA table_info(generated_strings);"],
        check=False,
    )
    if not table_check or "id" not in table_check.stdout:
        print("'generated_strings' table not found or schema not applied.")
        should_apply = auto_apply_schema or prompt_yes_no(
            f"Do you want to apply the schema from '{schema_file}' now? (y/N): ",
            default=False,
            interactive=interactive,
        )
        if should_apply:
            print(f"Applying schema from '{schema_file}'...")
            apply_schema_result = run_command(
                ["wrangler", "d1", "execute", db_name, "--file", schema_file],
                check=True,
            )
            if apply_schema_result:
                print("Schema applied successfully.")
            else:
                print("Failed to apply schema. Exiting.")
                return False
        else:
            print("Schema not applied. Exiting.")
            return False
    print("Schema is applied.")
    
    print("--- Wrangler Setup Complete ---")
    return True

def populate_d1_from_output(output_file, db_name, batch_size=500, dump_sql=None):
    print(f"\n--- Populating D1 Database '{db_name}' ---")
    if not os.path.exists(output_file):
        print(f"Error: Output file '{output_file}' not found. Please run 'make generate' first.")
        return False

    target_sql_path = dump_sql
    is_temp = False
    if not target_sql_path:
        temp_fd, target_sql_path = tempfile.mkstemp(prefix="d1_batch_", suffix=".sql")
        os.close(temp_fd)
        is_temp = True

    try:
        print(f"Preparing SQL statements from '{output_file}'...")
        total = generate_sql_file(output_file, target_sql_path, chunk_size=batch_size)
        if total == 0:
            print(f"No valid strings found in {output_file}. Nothing to populate.")
            return True
        print(f"Prepared {total} strings in '{target_sql_path}'.")

        if dump_sql:
            print(f"SQL dump completed successfully. File saved to '{dump_sql}'.")
            return True

        print(f"Executing batch SQL file against D1 database '{db_name}' via wrangler...")
        result = run_command(
            ["wrangler", "d1", "execute", db_name, "--file", target_sql_path],
            check=False,
        )
        if result is None or result.returncode != 0:
            print("Failed to execute batch SQL file on D1.")
            return False

        print("D1 population complete.")
        return True
    finally:
        if is_temp and os.path.exists(target_sql_path):
            try:
                os.remove(target_sql_path)
            except OSError:
                pass

def main():
    parser = argparse.ArgumentParser(
        description="Populate Cloudflare D1 database with generated halal strings via consolidated SQL file.",
    )
    parser.add_argument(
        "output_file",
        nargs="?",
        default="output.txt",
        help="Path to output file containing generated strings (default: output.txt)",
    )
    parser.add_argument(
        "db_name",
        nargs="?",
        default="my-halal-strings-db",
        help="Target Cloudflare D1 database name (default: my-halal-strings-db)",
    )
    parser.add_argument(
        "--batch-size",
        type=int,
        default=500,
        help="Batch size for D1 insert statements per transaction block (default: 500)",
    )
    parser.add_argument(
        "--dump-sql",
        type=str,
        default=None,
        help="Output path to dump batch SQL file without executing wrangler",
    )
    parser.add_argument(
        "--auto-create",
        action="store_true",
        help="Automatically create the D1 database if it does not exist (non-interactive)",
    )
    parser.add_argument(
        "--apply-schema",
        action="store_true",
        help="Automatically apply schema if table is missing (non-interactive)",
    )
    parser.add_argument(
        "--interactive",
        action="store_true",
        help="Enable interactive confirmation prompts if running in a TTY",
    )

    args = parser.parse_args()

    # If user just wants to dump SQL to file, bypass Wrangler checks
    if args.dump_sql:
        success = populate_d1_from_output(
            args.output_file,
            args.db_name,
            batch_size=args.batch_size,
            dump_sql=args.dump_sql,
        )
        sys.exit(0 if success else 1)

    if check_wrangler_setup(
        args.db_name,
        auto_create=args.auto_create,
        auto_apply_schema=args.apply_schema,
        interactive=args.interactive,
    ):
        success = populate_d1_from_output(args.output_file, args.db_name, batch_size=args.batch_size)
        sys.exit(0 if success else 1)
    else:
        print("Wrangler setup failed. Please resolve the issues and try again.")
        sys.exit(1)

if __name__ == "__main__":
    main()