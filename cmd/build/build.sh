#!/bin/bash

INPUT_DIR="$HOME/var/archive"
OUTPUT_DIR="$HOME/src/archive/app/content/"

#./build \
#    --input-dir="$INPUT_DIR" \
#    --output-dir="$OUTPUT_DIR"

# ./build --input-dir=$HOME/var/archive --output-dir=$HOME/src/archive/app/content/ --category=/islam_kamil/part_1 --index=4

CATEGORY=/islam_kamil/part_1
INDEX=1

pandoc ~/var/archive/md/$CATEGORY/$INDEX.md -f markdown -t html -o ~/var/archive/pages/$CATEGORY/$INDEX.html
./build \
    --input-dir="$INPUT_DIR" \
    --output-dir="$OUTPUT_DIR" \
    --category=$CATEGORY --index=$INDEX