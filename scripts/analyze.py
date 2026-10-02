import collections
import sys

def analyze_uniqueness(filename):
  with open(filename, 'r') as f:
    lines = f.read().splitlines()

  total_lines = len(lines)
  unique_lines = len(set(lines))

  print(f"Total lines: {total_lines}")
  print(f"Unique lines: {unique_lines}")

  if total_lines == 0:
    print("Uniqueness percentage: 0.00% (file is empty)")
    print("\nNo lines to analyze.")
    return

  uniqueness_percentage = (unique_lines / total_lines) * 100
  print(f"Uniqueness percentage: {uniqueness_percentage:.2f}%")

  duplicates = [line for line, count in collections.Counter(lines).items() if count > 1]
  if duplicates:
    print("\nDuplicate lines:")
    for line in duplicates:
      print(line)
  else:
    print("\nNo duplicate lines found.")

if __name__ == "__main__":
  if len(sys.argv) < 2:
    print("Usage: python3 analyze.py <output_file>")
    sys.exit(1)
  analyze_uniqueness(sys.argv[1])