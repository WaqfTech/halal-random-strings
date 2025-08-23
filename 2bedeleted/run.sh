#!/bin/bash

for i in {1..1000000}
do
  ./halal -r 1000 >> output.txt
done
