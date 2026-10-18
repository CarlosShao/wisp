#!/bin/bash
# 341-r1 AC#1 de-gum shots: four shapes, FOUR DIFFERENT landing spots (dispatch: 两发落点⛔ 相同).
#   P2 -> writeguard.go  the INFO branch's delete-shape clause ONLY (WARN branch must survive untouched)
#   P3 -> writeguard.go  the merge base taken from the FILE (line 141)          [different line than P2]
#   P5 -> loader.go      SaveFile's schema-version stamp (line 257)
#   P6 -> loader.go      applyPresets' base_url fill (line 213)
# Discipline (dispatch: 落刀前后都验命中数 + 逐发当场还原并报还原哈希):
#   - shots land on physical bytes of an EXPORTED tree outside the repo, never on tracked files;
#   - content anchors, never absolute line numbers;
#   - killed anchor: count BEFORE (must equal the expected number) and AFTER (must be 0);
#   - replacement anchor: count AFTER (must be >= 1);
#   - survivor anchor: a nearby string that must keep the SAME count, which proves the shot hit only
#     the one branch/line this shape is about (for P2 that is the whole difference between INFO and WARN);
#   - pristine hash verified before the shot, restore hash reported after, compared to the delivery copy;
#   - shots are SERIAL, one fresh numbered tree each; nothing is deleted (temp pieces are built, never removed).
set -u
repo='/d/work/workspace/projects plans/Wisp'
root='D:/tmp/wisp/20261018/341r1'
out="$repo/.scratch/wisp/probes/341/r1/logs"
deliv="$root/deliv"
export GOTMPDIR="$root/gotmp"
unset TMPDIR

n() { grep -cF -- "$1" "$2" || true; }   # count without killing the pipeline on 0 hits

shot() { # $1 tag $2 rel-file $3 perl-expr $4 killed-anchor $5 want $6 new-anchor $7 survivor-anchor $8 survivor-want $9 case-regex
  local tag="$1" rel="$2" px="$3" old_a="$4" want="$5" new_a="$6" sur_a="$7" sur_want="$8" casere="$9"
  local tree="$root/mu_$tag"
  local src="$deliv/internal/config/$rel"
  local f="$tree/internal/config/$rel"
  local lg="$out/05-mu-$tag.txt" full="$out/05-mu-$tag-full.txt"
  cp -r "$deliv" "$tree" || { echo "$tag: copy failed"; return 1; }
  local h0=$(sha256sum "$f" | cut -d' ' -f1) hs=$(sha256sum "$src" | cut -d' ' -f1)
  local c_before=$(n "$old_a" "$f") c_sur_before=$(n "$sur_a" "$f")
  {
    echo "=== SHOT $tag  internal/config/$rel ==="
    echo "tree=$tree (exported, outside the repo; tracked files untouched)"
    echo "tree copy hash  = $h0"
    echo "delivery hash   = $hs    tree==delivery ? $([ "$h0" = "$hs" ] && echo YES || echo NO)"
    echo "killed anchor   [$old_a]  BEFORE=$c_before  want=$want"
    echo "survivor anchor [$sur_a]  BEFORE=$c_sur_before  want=$sur_want  (proves the shot does not eat the other branch)"
  } > "$lg" 2>&1
  if [ "$c_before" != "$want" ] || [ "$h0" != "$hs" ] || [ "$c_sur_before" != "$sur_want" ]; then
    echo "ABORTED: pre-flight gate failed (anchor hit count or pristine hash mismatch) - no shot landed" >> "$lg"
    echo "$tag ABORTED" ; return 1
  fi
  perl -0pi -e "$px" "$f"
  local c_after=$(n "$old_a" "$f") c_new=$(n "$new_a" "$f") c_sur_after=$(n "$sur_a" "$f")
  {
    echo "after shot: killed anchor=$c_after (want 0)  new anchor [$new_a]=$c_new (want >=1)  survivor=$c_sur_after (want $sur_want)"
    echo "landed hash = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")"
  } >> "$lg" 2>&1
  cd "$tree" || exit 2
  go test ./internal/config/ -count=1 -timeout 30m -v > "$full" 2>&1
  local r=$?
  {
    echo "go test rc=$r"
    echo "RUN=$(grep -c '^=== RUN' "$full" || true) PASS_TOP=$(grep -c '^--- PASS' "$full" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$full" || true) SKIP=$(grep -c '^--- SKIP' "$full" || true)"
    echo "---- full red roster (⛔ mine only; other legs' cases reddening is a legitimate reading of a dead guard) ----"
    grep '^--- FAIL' "$full" || echo "(no top-level FAIL at all = THE SHOT DID NOT FIRE, this nail has no teeth)"
    echo "---- the ticket-341 case this shot targets, named ----"
    grep -E "^--- (PASS|FAIL): ($casere)" "$full" || echo "(case absent from the run = bench failure, not a reading)"
    echo "---- verbatim red sentences in writeguard_unnailed_341_test.go ----"
    grep -E 'writeguard_unnailed_341_test\.go:[0-9]+:' "$full" | cut -c1-200
    cp "$src" "$f"
    local hr=$(sha256sum "$f" | cut -d' ' -f1)
    echo "RESTORED hash = $hr   restored==pristine ? $([ "$hr" = "$h0" ] && echo YES || echo NO)"
    echo "re-check killed anchor after restore = $(n "$old_a" "$f") (want $want)"
  } >> "$lg" 2>&1
  echo "$tag done rc=$r -> $lg"
}

# P2: remove ONLY the INFO branch's delete-shape clause. INFO ends `... a key line the "+\n "operator deleted ... value)"`;
# WARN ends `... a key line the operator deleted by "+\n "hand comes back ... value",`  -> different byte shapes, so the
# kill can be branch-exact, and the WARN clause is the survivor anchor that must stay at 1.
shot p2 writeguard.go \
  's/"in canonical form, so comments and hand-chosen key order are not preserved; a key line the "\+\n\s*"operator deleted by hand comes back at its schema default value\)"/"in canonical form, so comments and hand-chosen key order are not preserved)"/' \
  'are not preserved; a key line the "' 1 \
  'are not preserved)"' \
  'a key line the operator deleted by "' 1 \
  'TestHandDeletedDefaultValuedLineIsRestoredInvisibly'

# P3: the merge base from memory instead of the file (this is ticket 226's original bug shape).
shot p3 writeguard.go \
  's/base = deepCopyConfig\(disk\)/base = deepCopyConfig(m.cur)/' \
  'base = deepCopyConfig(disk)' 1 \
  'base = deepCopyConfig(m.cur)' \
  'if disk != nil {' 2 \
  'TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory'

# P5: SaveFile stops forcing the current schema version.
shot p5 loader.go \
  's/cp\.SchemaVersion = SchemaVersionCurrent/cp.SchemaVersion = 1/' \
  'cp.SchemaVersion = SchemaVersionCurrent' 1 \
  'cp.SchemaVersion = 1' \
  'data, err := MarshalCanonical(cp)' 1 \
  'TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed'

# P6: applyPresets stops supplying the endpoint from the preset table.
shot p6 loader.go \
  's/p\.BaseURL = pre\.BaseURL/p.BaseURL = ""/' \
  'p.BaseURL = pre.BaseURL' 1 \
  'p.BaseURL = ""' \
  'p.Protocol = pre.Protocol' 1 \
  'TestHandDeletedProviderBaseUrlComesBackAsThePresetValue'

echo "ALL FOUR SHOTS DONE"
