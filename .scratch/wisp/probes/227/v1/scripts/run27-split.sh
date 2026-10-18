#!/bin/bash
cd '/d/work/workspace/projects plans/Wisp' || exit 2
all=$(git log --format='%H' --grep='票 227 终裁腿 227-v1·')
m=0; o=0; ol=""
for h in $all; do
  case "$h" in
    173796c4*|a13c87fa*|88338d6a*|bc4a3fbb*|65558bb1*|761681de*) m=$((m+1));;
    *) o=$((o+1)); ol="$ol ${h:0:8}";;
  esac
done
echo "prefix-total=$(echo "$all" | grep -c .)  mine=$m  foreign=$o"
echo "foreign:$ol"
