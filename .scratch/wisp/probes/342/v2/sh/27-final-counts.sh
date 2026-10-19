#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "FILES_342_v2=$(find .scratch/wisp/probes/342/v2 -type f | wc -l)"
echo "FILES_348_v1=$(find .scratch/wisp/probes/348/v1 -type f | wc -l)"
echo "ZERO_BYTE=$(find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 -type f -size 0 | wc -l)"
echo "DOT_OUT=$(find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 -name '*.out' | wc -l)"
echo "MY_ROSTER_6=$(git log --all --format='%H%x09%s' | awk -F'\t' 'index($2,"\xE7\xA5\xA8 342\xEF\xBC\x8B\xE7\xA5\xA8 348 \xE7\xBB\x88\xE8\xA3\x81\xE8\x85\xBF 342-v2")==1' | wc -l)"
echo "MY_COMMITS:"
git log --all --format='%H%x09%s' | awk -F'\t' 'index($2,"\xE7\xA5\xA8 342\xEF\xBC\x8B\xE7\xA5\xA8 348 \xE7\xBB\x88\xE8\xA3\x81\xE8\x85\xBF 342-v2")==1{print substr($1,1,8)}'
echo "DIRTY_TRACKED=$(git status --porcelain -uno | wc -l)"
echo "GO_BYTES_IN_ANY_OF_MY_COMMITS:"
git log --all --format='%H' | awk -F'\t' '{print}' > /dev/null
for h in $(git log --all --format='%H%x09%s' | awk -F'\t' 'index($2,"\xE7\xA5\xA8 342\xEF\xBC\x8B\xE7\xA5\xA8 348 \xE7\xBB\x88\xE8\xA3\x81\xE8\x85\xBF 342-v2")==1{print $1}'); do
  echo "  $h go=$(git show --name-only --format= $h | grep -c '\.go$' || true) files=$(git show --name-only --format= $h | grep -c .)"
done
