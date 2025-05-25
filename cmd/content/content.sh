#!/bin/bash

CONTENT_DIR="$HOME/work/src/archive/app/content"

#SECTION=""
#SECTION="islam_kamil"
SECTION="islam_kamil/part_1"
HADITH_SECTION="islam_kamil_hadiths/part_1"
#SECTION="islam_kamil_hadiths/part_1"
#SECTION="islam_kamil/part_3/chapter_2"

if [ $2 == "hadith" ]; then
  SECTION=$HADITH_SECTION
  INDEX=$3
else
  INDEX=$2
fi


SECTION_DIR="$CONTENT_DIR/$SECTION"
ITEM_FILE="$SECTION_DIR/$INDEX.json"

FIELDS_FILE="tmp.json"
TEXT_FILE="tmp.md"

#./content --unpack --item-file=$ITEM_FILE
#./content --pack --item-file="$ITEM_FILE"

#./content --repack=id --item-dir="$SECTION_DIR"
#./content --repack=index --item-dir="$SECTION_DIR" --skip=203000
#./content --repack=dir_index --item-dir="$SECTION_DIR"

#./content --repack=meta --item-dir="$SECTION_DIR"

# main command

./content "--$1" --item-file="$ITEM_FILE" --fields-file="$FIELDS_FILE" --text-file="$TEXT_FILE"

# main command hadith

#./content --repack=title --item-dir="$SECTION_DIR" --title-template='Хадис {{index}}'

# main command loop

#SRC_SECTION_DIR="$CONTENT_DIR/islam_kamil/_part_1"
#
#for INDEX in $(seq 1 107)
#do
#  ITEM_FILE="$SECTION_DIR/$INDEX.json"
#  FIELDS_FILE="$SRC_SECTION_DIR/$INDEX.json"
#  TEXT_FILE="$SRC_SECTION_DIR/$INDEX.md"
#  ./content --pack --item-file="$ITEM_FILE" --fields-file="$FIELDS_FILE" --text-file="$TEXT_FILE"
#done