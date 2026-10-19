#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "WINSEC_GO_BLOB_NOW=$(git rev-parse HEAD:internal/winsec/winsec.go)"
echo "T342_342v2_BULLET_MENTIONS=$(grep -c 'agent=`342-v2`' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md || true)"
echo "T348_342v2_BULLET_MENTIONS=$(grep -c 'agent=`342-v2`' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md || true)"
echo '--- my bullet head line in each ticket (must be the LAST bullet) ---'
tail -1 .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md | cut -c1-150
tail -1 .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | cut -c1-150
echo '--- boxes ---'
echo "T342 x=$(grep -c '^- \[x\]' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md) open=$(grep -c '^- \[ \]' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md)"
echo "T348 x=$(grep -c '^- \[x\]' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md) open=$(grep -c '^- \[ \]' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md)"
