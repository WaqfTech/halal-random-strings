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

def populate_d1_from_output(output_file, db_name, batch_size=50):
    print(f"\n--- Populating D1 Database '{db_name}' ---")
    print(f"Reading strings from {output_file}...")
    
    try:
        with open(output_file, 'r', encoding='utf-8') as f:
            lines = [line.strip() for line in f if line.strip()]
    except FileNotFoundError:
        print(f"Error: Output file '{output_file}' not found. Please run 'make generate' first.")
        return False

    if not lines:
        print(f"No strings found in {output_file}. Nothing to populate.")
        return True

    print(f"Found {len(lines)} strings to insert.")

    # Safe batch inserts to avoid ARG_MAX and CLI length overflows
    batch_size = max(1, min(batch_size, 100))
    total_batches = (len(lines) + batch_size - 1) // batch_size
    
    for i in range(0, len(lines), batch_size):
        batch = lines[i:i + batch_size]
        values_sql = []
        current_timestamp = int(time.time())
        for line in batch:
            # Escape single quotes within the string for SQL
            escaped_line = line.replace("'", "''")
            values_sql.append(f"('{escaped_line}', NULL, NULL, {current_timestamp})")

        # INSERT OR IGNORE avoids aborting on duplicate collisions
        insert_sql = (
            f"INSERT OR IGNORE INTO generated_strings (id, used_at, used_by, created_at) "
            f"VALUES {', '.join(values_sql)};"
        )

        batch_num = i // batch_size + 1
        print(f"Executing batch {batch_num}/{total_batches} ({len(batch)} items)...")
        result = run_command(["wrangler", "d1", "execute", db_name, "--command", insert_sql], check=False)
        if result is None:
            print(f"Failed to execute batch {batch_num}. Aborting population.")
            return False

    print("D1 population complete.")
    return True

def main():
    parser = argparse.ArgumentParser(
        description="Populate Cloudflare D1 database with generated halal strings.",
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
        default=50,
        help="Batch size for D1 insert statements (default: 50, safe for CLI limits)",
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