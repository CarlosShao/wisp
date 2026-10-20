#!/bin/bash
cd '/d/work/workspace/projects plans/Wisp' || exit 9
S7=$(printf '\xc2\xa7')
grep -n "$S7" docs/evidence/s1/353-v1-verdict.md | cut -c1-240
echo "--- count: $(grep -o "$S7" docs/evidence/s1/353-v1-verdict.md | wc -l)"
echo "--- 正控：票 353 面同尺命中 $(git cat-file blob 'HEAD:.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md' | grep -o "$S7" | wc -l)"
