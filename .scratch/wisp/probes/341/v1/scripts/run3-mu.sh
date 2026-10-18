#!/usr/bin/env bash
# 341-v1 AC#1: the four de-gum shots, each re-run by ME from zero, serial, on the
# export tree's PHYSICAL bytes (no -overlay, no tracked file touched).
# ★Pristine copies come from the NAMED ANCHOR BLOB (git show), never from the tree,
#   so an aborted previous shot cannot launder a mutation into a "pristine".
# Content anchors, not line numbers; counts use perl whole-file /s so an anchor across a
# Go string-literal continuation is counted honestly (341-r1's named ruler defect).
W=/d/tmp/wisp/20261018/341v1
T="$W/treeD"
R="/d/work/workspace/projects plans/Wisp"
ANCHOR=6c625a8c23d46c347057c8f5050e8f8a10e33fa3
LOGROOT="$1"
P="$W/pristine"
mkdir -p "$P" "$LOGROOT"
CUR_FILE=""
restore_now() { [ -n "$CUR_FILE" ] && cp -f "$P/$(basename "$CUR_FILE")" "$CUR_FILE" && echo "FORCED_RESTORE $CUR_FILE"; }
trap 'restore_now' EXIT

for rel in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
  ( cd "$R" && git show "$ANCHOR:$rel" ) > "$P/$(basename "$rel")"
  echo "PRISTINE(from blob $ANCHOR:$rel) $(basename "$rel") sha256=$(sha256sum "$P/$(basename "$rel")" | cut -d' ' -f1)"
done
echo "TREE_AT_ENTRY writeguard.go=$(sha256sum "$T/internal/config/writeguard.go" | cut -d' ' -f1)"
echo "TREE_AT_ENTRY loader.go=$(sha256sum "$T/internal/config/loader.go" | cut -d' ' -f1)"

FIVE="TestHandDeletedDefaultValuedLineIsRestoredInvisibly TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed TestHandDeletedProviderBaseUrlComesBackAsThePresetValue TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile"

hits() { perl -0777 -ne 'my $c = () = /'"$2"'/gs; print $c' "$1"; }

run_pkg() {
  ( cd "$T" || exit 9
    export GOCACHE="$(cygpath -w "$W/gocache")"
    export PATH="$(cygpath -w "$T/third_party/sherpa-onnx");$(cygpath -w "$T/build");$PATH"
    go test ./internal/config/ -count=1 -timeout 30m -v > "$1" 2>&1
    echo "GO_RC=$?" )
}

# args: name rel cut new-re disc-re replacement OB OA NB NA DB DA
shot() {
  name="$1"; rel="$2"; cut="$3"; newre="$4"; discre="$5"; REPL="$6"
  OB="$7"; OA="$8"; NB="$9"; NA="${10}"; DB="${11}"; DA="${12}"
  f="$T/$rel"; CUR_FILE="$f"; base=$(basename "$f")
  d="$W/log-$name"; mkdir -p "$d"; out="$d/test.log"
  L="$LOGROOT/10-mu-$name.txt"
  {
  echo "############ SHOT $name :: $rel (anchor $ANCHOR) ############"
  echo "anchor_old(cut)   =$cut      want before=$OB after=$OA"
  echo "anchor_new        =$newre    want before=$NB after=$NA"
  echo "anchor_discrim    =$discre   want before=$DB after=$DA"
  echo "-- hash at entry (want == pristine-from-blob) --"; sha256sum "$f" | cut -d' ' -f1
  diff -q "$f" "$P/$base" >/dev/null || { echo "SHOT ABORTED: tree file != pristine blob"; exit 1; }
  OLDB=$(hits "$f" "$cut"); NEWB=$(hits "$f" "$newre"); DISCB=$(hits "$f" "$discre")
  echo "-- BEFORE: old=$OLDB new=$NEWB disc=$DISCB --"
  { [ "$OLDB" = "$OB" ] && [ "$NEWB" = "$NB" ] && [ "$DISCB" = "$DB" ]; } ||
    { echo "SHOT ABORTED: pre-cut anchor bookkeeping wrong"; exit 1; }
  REPVAL="$REPL" perl -e '
    my ($f, $re) = @ARGV;
    my $rp = $ENV{REPVAL};
    local $/;
    open my $fh, "<", $f or die "read: $!"; my $t = <$fh>; close $fh;
    my $n = ($t =~ s/$re/$rp/g);
    die "cut matched $n times, want 1\n" unless $n == 1;
    open my $o, ">", $f or die "write: $!"; print $o $t; close $o or die "close: $!";
    print "PERL_CUT_OK (1 substitution)\n";
  ' "$f" "$cut" || { echo "SHOT ABORTED: perl cut refused"; exit 1; }
  OLDA=$(hits "$f" "$cut"); NEWA=$(hits "$f" "$newre"); DISCA=$(hits "$f" "$discre")
  echo "-- AFTER: old=$OLDA new=$NEWA disc=$DISCA --"
  { [ "$OLDA" = "$OA" ] && [ "$NEWA" = "$NA" ] && [ "$DISCA" = "$DA" ]; } ||
    { echo "SHOT ABORTED: post-cut anchor bookkeeping wrong"; cp -f "$P/$base" "$f"; exit 1; }
  echo "-- mutated region --"; grep -nE -B2 -A2 "$newre" "$f" | head -18
  echo "-- hash after cut (want != pristine) --"; sha256sum "$f" | cut -d' ' -f1
  echo "-- full package run (mutated) --"
  run_pkg "$out"
  echo "RUN=$(grep -c '^=== RUN' $out || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' $out || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' $out || true)"
  echo "SKIP=$(grep -c '^--- SKIP' $out || true)"
  echo "-- red roster (top-level, named) --"
  grep '^--- FAIL' "$out" | sed 's/^--- FAIL: //' | cut -d' ' -f1 | sort | tee "$d/rednames.txt"
  echo "-- per-case verdict for the five 341 cases --"
  for t in $FIVE; do
    v=$(grep -E "^ *--- (PASS|FAIL): $t( |\()" "$out" | sed 's/^ *//' | awk '{print $2}' | head -1)
    printf "%s = %s\n" "$t" "${v:-ABSENT}"
  done
  echo "-- verbatim red lines in the 341 carrier --"
  grep -E "^ +writeguard_unnailed_341_test\.go:[0-9]+: " "$out" | sed 's/^ *//'
  echo "-- failing 341 case :: its first red line --"
  awk '/^--- FAIL: /{n=$3} /^ +writeguard_unnailed_341_test\.go:[0-9]+: /{ if (!(n in seen)) { print n" :: "$0; seen[n]=1 } }' "$out"
  echo "-- restore --"
  cp -f "$P/$base" "$f"; CUR_FILE=""
  echo "-- restored hash (want == pristine) --"; sha256sum "$f" | cut -d' ' -f1
  echo "-- old anchor hits after restore (want $OB) --"; hits "$f" "$cut"; echo
  } > "$L" 2>&1
  rc=$?
  echo "LOGGED $L rc=$rc bytes=$(wc -c < "$L") lines=$(wc -l < "$L")"
  [ $rc -eq 0 ] || { echo "SERIAL STOP after $name (rc=$rc)"; exit $rc; }
}

# P2: cut ONLY the INFO-branch delete-shape clause; the WARN branch keeps its copy.
#     disc "a key line the" 2 -> 1 = the cut was branch-local (proves the two branches
#     are different lines, and that P2's own clause is the one removed).
shot p2 internal/config/writeguard.go \
  'are not preserved; a key line the\s*"\+\s*"operator deleted by hand comes back at its schema default value' \
  'hand-chosen key order are not preserved\)"' \
  'a key line the' \
  'are not preserved' \
  1 0 1 2 2 1

# P3: merge base退回内存快照 (ticket 226's original defect shape).
#     new anchor form only exists after the cut (the declaration at :133 uses ":=").
shot p3 internal/config/writeguard.go \
  'base = deepCopyConfig\(disk\)' \
  'base = deepCopyConfig\(m\.cur\)' \
  'disk != nil' \
  'base = deepCopyConfig(m.cur)' \
  1 0 0 1 2 2

# P5: SaveFile stops stamping the current schema version (pins to 1).
shot p5 internal/config/loader.go \
  'cp\.SchemaVersion = SchemaVersionCurrent' \
  'cp\.SchemaVersion = 1' \
  'func SaveFile' \
  'cp.SchemaVersion = 1' \
  1 0 0 1 1 1

# P6: the preset fill for base_url becomes a no-op.
#     disc "pre.Protocol" stays 1 = the protocol fill survived = preset path only.
shot p6 internal/config/loader.go \
  'p\.BaseURL = pre\.BaseURL' \
  'p\.BaseURL = ""' \
  'pre\.Protocol' \
  'p.BaseURL = ""' \
  1 0 0 1 1 1

echo "ALL_SHOTS_DONE"
