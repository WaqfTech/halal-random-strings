import subprocess
import json
import time

def populate_d1_from_output(output_file="output.txt", db_name="my-halal-strings-db"):
    print(f"Reading {output_file}...")
    with open(output_file, 'r') as f:
        lines = f.read().splitlines()

    if not lines:
        print(f"No lines found in {output_file}. Exiting.")
        return

    print(f"Found {len(lines)} strings. Preparing for D1 insertion...")

    # Batch inserts for efficiency
    batch_size = 1000
    for i in range(0, len(lines), batch_size):
        batch = lines[i:i + batch_size]
        values_sql = []
        current_timestamp = int(time.time())
        for line in batch:
            # Escape single quotes within the string for SQL
            escaped_line = line.replace("'", "''")
            values_sql.append(f"('{escaped_line}', NULL, NULL, {current_timestamp})")

        insert_sql = f"INSERT INTO generated_strings (id, used_at, used_by, created_at) VALUES {', '.join(values_sql)};"

        try:
            print(f"Executing batch {i // batch_size + 1}/{(len(lines) + batch_size - 1) // batch_size}...")
            command = [
                "wrangler", "d1", "execute", db_name,
                "--command", insert_sql
            ]
            result = subprocess.run(command, capture_output=True, text=True, check=True)
            print("Stdout:", result.stdout)
            if result.stderr:
                print("Stderr:", result.stderr)
        except subprocess.CalledProcessError as e:
            print(f"Error executing wrangler command: {e}")
            print(f"Command: {' '.join(e.cmd)}")
            print(f"Stdout: {e.stdout}")
            print(f"Stderr: {e.stderr}")
            return
        except FileNotFoundError:
            print("Error: 'wrangler' command not found. Please ensure Cloudflare Wrangler CLI is installed and in your PATH.")
            return

    print("D1 population complete.")

if __name__ == "__main__":
    populate_d1_from_output()
