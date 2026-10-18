#!/bin/bash
# 227-v1 AC#6 + 必答(iv) gauges (committed bytes only; no go test in cmd/wisp = lane-held)
cd '/d/work/workspace/projects plans/Wisp'
out='.scratch/wisp/probes/227/v1/logs/08-ac6-fix-and-iv-gauges.txt'
{
echo "tree gauge = HEAD=$(git rev-parse --short HEAD) (2038c45d..HEAD touched 0 files under internal/config cmd/wisp)"
echo
echo "== FIX: my earlier tr-s gauge was mismeasured by me: tr -s ' \t' does not join the newline,"
echo "   so the sentence split over :43/:44 can never match on ONE line. Correct gauges below,"
echo "   both content-anchored, each with a negative control that must return 0 (rc=1). =="
echo "A) two-line window join (sed):"
git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | sed -n '43,44p' | tr -d '\n' | tr -s ' \t/' ' ' | grep -c 'running and nothing polls: \[app\] theme, a key'
a=$(git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | sed -n '43,44p' | tr -d '\n' | tr -s ' \t/' ' ' | grep -c 'ZZZ not a phrase' || true)
echo "A-negctrl (must be 0): hits=$a (grep rc=1 when 0)"
echo "B) whole file joined on space:"
git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | tr '\n' ' ' | tr -s ' \t/' ' ' | grep -c 'running and nothing polls'
b=$(git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | tr '\n' ' ' | tr -s ' \t/' ' ' | grep -c 'polls nothing whatsoever' || true)
echo "B-negctrl (must be 0): hits=$b"
echo "C) per-part hits with line numbers (the 09-29 leg reading this was 'nothing' twice):"
git grep -n 'nothing polls\|and nothing$' HEAD -- cmd/wisp/always_write_no_clobber_226_test.go | cut -c1-120
echo
echo "=== 必答(iv): do the three log literals tell the truth today? gauges: ==="
echo "[iv-a] the narrowed nothing-written sentence: does the CONDITIONAL half hold in production code?"
echo "  the branch returns nil WITHOUT SaveFile when written==0:"
git show HEAD:internal/config/writeguard.go | sed -n '149,159p'
echo "[iv-b] does the schema-migration really re-lay inside readConfigFile before that sentence logs?"
echo "  mergeWrite step (1) readConfigFile at writeguard.go:118 -> loader.go applyMigrations -> migrate.go PrivateFile+atomicWrite:"
git grep -n 'readConfigFile(m.path)' HEAD -- internal/config/writeguard.go
git grep -n 'applyMigrations(path, raw, ver)' HEAD -- internal/config/loader.go
git grep -n 'winsec.PrivateFile(backup, raw, 0o600)\|atomicWrite(path, data)' HEAD -- internal/config/migrate.go
echo "[iv-c] the delete-shape sentence (两句): does a deleted scalar line really come back at the SCHEMA DEFAULT,"
echo "  not at memory's value? readout = my MU-1 run + the leg's own c1 + MY named PASS list above"
echo "  (01-deliv: TestHandDeletedScalarLineComesBackAtSchemaDefault PASS; under MU-2b that case reds ONLY on the 文案 line :64"
echo "   => the restoration half is true without my sentence, the sentence names a fact production already has)."
echo "[iv-d] residue check: does the Info branch still EVER fire with a schema-stale file that DID move bytes?"
echo "  (yes: matrixBase hand case with no schema_version -> leg P5' + my MU-1 red 404/407 prove bytes+mtime moved while"
echo "   the file 'already carried the key'. The narrowed sentence is the honest one.)"
} > "$out" 2>&1
echo "written bytes=$(wc -c < "$out")"
