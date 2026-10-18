#!/bin/bash
cd '/d/work/workspace/projects plans/Wisp' || exit 2
t='.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md'
b=/d/tmp/wisp/20261018/227v1/ticket-before.txt
echo "before-file CR count = $(tr -dc '\r' < "$b" | wc -c)"
echo "after-file  CR count = $(tr -dc '\r' < "$t" | wc -c)"
echo "my bullet source CR  = $(tr -dc '\r' < /d/tmp/wisp/20261018/227v1/msgs/bullet227.txt | wc -c)"
echo "committed blob CR    = $(git show HEAD:"$t" | tr -dc '\r' | wc -c)"
echo "new lines introduced = $(( $(wc -c < "$t") - $(wc -c < "$b") )) byte delta, CR delta = $(( $(tr -dc '\r' < "$t" | wc -c) - $(tr -dc '\r' < "$b" | wc -c) ))"
