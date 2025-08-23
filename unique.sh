awk '
{ count[$0]++ } 
END {
  total = 0; unique_lines = 0; repeated_lines = 0
  for (line in count) {
    total += count[line]
    if (count[line] == 1) unique_lines += 1
    else repeated_lines += count[line]
  }
  printf "Total lines: %d\n", total
  printf "Unique lines: %d (%.2f%%)\n", unique_lines, unique_lines/total*100
  printf "Repeated lines: %d (%.2f%%)\n", repeated_lines, repeated_lines/total*100
}' output.txt

