import subprocess
import json
import time
import os

def run_command(command, check=True, capture_output=True, text=True):
    try:
        result = subprocess.run(command, check=check, capture_output=capture_output, text=text)
        return result
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {' '.join(e.cmd)}")
        print(f"Stdout: {e.stdout}")
        print(f"Stderr: {e.stderr}")
        if check:
            raise
        return None
    except FileNotFoundError:
        print("Error: 'wrangler' command not found. Please ensure Cloudflare Wrangler CLI is installed and in your PATH.")
        raise

def check_wrangler_setup(db_name):
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
        response = input(f"Do you want to create '{db_name}' now? (y/N): ").lower()
        if response == 'y':
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

    # 4. Check if schema is applied (basic check by trying to insert a dummy row and delete it)
    print("Checking if schema is applied...")
    schema_file = "db/schema.sql"
    if not os.path.exists(schema_file):
        print(f"Error: Schema file '{schema_file}' not found. Please ensure it's in the correct location.")
        return False

    # Try to get table info to see if table exists
    table_check = run_command(["wrangler", "d1", "execute", db_name, "--command", "PRAGMA table_info(generated_strings);"], check=False)
    if not table_check or "id" not in table_check.stdout:
        print("'generated_strings' table not found or schema not applied.")
        response = input(f"Do you want to apply the schema from '{schema_file}' now? (y/N): ").lower()
        if response == 'y':
            print(f"Applying schema from '{schema_file}'...")
            apply_schema_result = run_command(["wrangler", "d1", "execute", db_name, "--file", schema_file], check=True)
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

def populate_d1_from_output(output_file, db_name):
    print(f"\n--- Populating D1 Database '{db_name}' ---")
    print(f"Reading strings from {output_file}...")
    
    try:
        with open(output_file, 'r') as f:
            lines = f.read().splitlines()
    except FileNotFoundError:
        print(f"Error: Output file '{output_file}' not found. Please run 'make generate' first.")
        return

    if not lines:
        print(f"No strings found in {output_file}. Nothing to populate.")
        return

    print(f"Found {len(lines)} strings to insert.")

    # Batch inserts for efficiency
    batch_size = 1000
    total_batches = (len(lines) + batch_size - 1) // batch_size
    
    for i in range(0, len(lines), batch_size):
        batch = lines[i:i + batch_size]
        values_sql = []
        current_timestamp = int(time.time())
        for line in batch:
            # Escape single quotes within the string for SQL
            escaped_line = line.replace("'", "''")
            values_sql.append(f"('{escaped_line}', NULL, NULL, {current_timestamp})")

        insert_sql = f"INSERT INTO generated_strings (id, used_at, used_by, created_at) VALUES {', '.join(values_sql)};"

        print(f"Executing batch {i // batch_size + 1}/{total_batches}...")
        result = run_command(["wrangler", "d1", "execute", db_name, "--command", insert_sql])
        if result is None:
            print(f"Failed to execute batch {i // batch_size + 1}. Aborting population.")
            return

    print("D1 population complete.")

if __name__ == "__main__":
    import sys

    output_file = "output.txt"
    db_name = "my-halal-strings-db"

    if len(sys.argv) > 1:
        output_file = sys.argv[1]
    if len(sys.argv) > 2:
        db_name = sys.argv[2]

    if check_wrangler_setup(db_name):
        populate_d1_from_output(output_file, db_name)
    else:
        print("Wrangler setup failed. Please resolve the issues and try again.")