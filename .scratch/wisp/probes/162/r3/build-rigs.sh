#!/usr/bin/env bash
# 162-r3 mutation rigs. Writes ONLY into .scratch/wisp/probes/162/r3/**.
# Production files are never edited here: every mutation reaches the compiler
# through `go test -overlay`. Two rigs:
#   base  = fs_edit.go at this cell's anchor (AC#3b's two checks ABSENT)
#           -> the 恒真 question: do the new refusal cases stay quiet without the fix?
#   m3    = M-3 re-applied ON TOP OF r3's code (force the bytes about to land
#           into CRLF) -> the reverse ruler must go red here.
# Every copy carries the repo's own anti-no-op discipline: `proof` fails the run
# unless the mutation marker is present and the symbol it replaced is GONE.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 1
ABS="D:/work/workspace/projects plans/Wisp"
V=".scratch/wisp/probes/162/r3"
ANCHOR=ff0007840ac9489c6c5ef017ca5092cdd2ed6094
mkdir -p "$V/mut/base" "$V/mut/m3" "$V/logs"

sub() { # sub <file> <old> <new> -> exactly one replacement, else FATAL
  local f="$1"
  OLD="$2" NEW="$3" perl -0777 -i -pe 'our $c; $c += s/\Q$ENV{OLD}\E/$ENV{NEW}/g; END{print STDERR "COUNT=".($c//0)."\n"}' "$f" 2>"$V/logs/.cnt"
  local c; c=$(cat "$V/logs/.cnt")
  echo "  $f : $c"
  [ "$c" = "COUNT=1" ] || { echo "FATAL: expected exactly 1 replacement in $f ($c)"; return 1; }
}

proof() { # proof <file> <must-be-present> <must-be-absent>
  local f="$1" pres="$2" ab="$3"
  local nh oh
  nh=$(grep -c -F "$pres" "$f"); oh=$(grep -c -F "$ab" "$f")
  echo "  PROOF $f present-hits=$nh absent-hits=$oh"
  [ "$nh" -ge 1 ] && [ "$oh" -eq 0 ] || { echo "FATAL: overlay text not applied in $f"; return 1; }
}

ovl() { # ovl <name> <repo-relative-file> <mut-file>
  printf '{ "Replace": {\n    "%s/%s": "%s/%s"\n  } }\n' "$ABS" "$2" "$ABS" "$3" > "$V/overlay-$1.json"
  echo "  wrote $V/overlay-$1.json"
}

echo "== BASE: fs_edit.go at anchor $ANCHOR (both AC#3b checks absent) =="
git show "$ANCHOR:internal/tools/fs_edit.go" > "$V/mut/base/fs_edit.go" || exit 1
proof "$V/mut/base/fs_edit.go" 'Deliberately NOT here (ticket 162 r2): CRLF/BOM handling' 'leadingBOMLen' || exit 1
proof "$V/mut/base/fs_edit.go" 'Deliberately NOT here' 'lineEndingConvention' || exit 1
ovl base "internal/tools/fs_edit.go" "$V/mut/base/fs_edit.go"

echo "== M-3 on top of r3: force the staged bytes to CRLF =="
cp internal/tools/fs_edit.go "$V/mut/m3/fs_edit.go"
sub "$V/mut/m3/fs_edit.go" '	updated := b.String()' '	updated := strings.ReplaceAll(strings.ReplaceAll(b.String(), "\r\n", "\n"), "\n", "\r\n") // MUT-R3 same shape as 162-v1 MUT-3' || exit 1
proof "$V/mut/m3/fs_edit.go" 'MUT-R3 same shape as 162-v1 MUT-3' '	updated := b.String()' || exit 1
ovl m3 "internal/tools/fs_edit.go" "$V/mut/m3/fs_edit.go"

echo "== runs =="
go test -count=1 -overlay="$V/overlay-base.json" -v -run 'TestFSEdit' ./internal/tools/ > "$V/logs/base-on-unfixed-code.log" 2>&1
echo "  base rc=$? (expected NONZERO: the refusal cases must be red without the fix)"
grep -E "(^|[[:space:]])--- FAIL" "$V/logs/base-on-unfixed-code.log" | sed 's/^/  BASE-RED  /'
go test -count=1 -overlay="$V/overlay-m3.json" -v -run 'TestFSEdit' ./internal/tools/ > "$V/logs/m3-on-r3.log" 2>&1
echo "  m3 rc=$? (expected NONZERO: the reverse ruler must be red under forced CRLF)"
grep -E "(^|[[:space:]])--- FAIL" "$V/logs/m3-on-r3.log" | sed 's/^/  M3-RED    /'
echo "== rigs built =="

# ---------------------------------------------------------------- AC#4b ruler
# M-4, re-derived HERE (162-v1's own mut/m4 is read-only evidence): take the one
# staging call out of fs.edit and write the target directly with production's own
# chunked writer, same WriteChunk, same kill seam. That is the falsification of
# the byte-level assertion in the real-taskkill case - if the rig is a no-op the
# ruler is decorative.
echo "== M-4 on top of r3: temp+rename removed from fs.edit (AC#4b ruler) =="
mkdir -p "$V/mut/m4"; cp internal/tools/fs_edit.go "$V/mut/m4/fs_edit.go"
sub "$V/mut/m4/fs_edit.go" '	res := t.d.stageAndRename(ctx, target, strings.NewReader(updated), onUpdate)
	res.Origin = target' '	se := t.d.ledger()
	f, oerr := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // MUT-R3-M4 no staging file, direct write (same kill seam)
	if oerr != nil {
		return Result{Text: "打开目标失败：" + oerr.Error(), IsError: true, Origin: target}, nil
	}
	nw, werr := se.writeAll(ctx, f, strings.NewReader(updated), onUpdate, limit)
	_ = f.Close()
	res := Result{Origin: target, AppliedSteps: se.snapshot()}
	if werr != nil {
		res.Text = "写入未完成：" + werr.Error()
		res.IsError = true
	} else {
		res.Text = fmt.Sprintf("fs.edit 已改写 %s：%d 枚编辑全部生效，%d 字节 / %d 行 → %d 字节 / %d 行（直接写，无临时文件）",
			target, len(matches), len(content), countNL(content), nw, countNL(updated))
	}' || exit 1
proof "$V/mut/m4/fs_edit.go" 'MUT-R3-M4 no staging file' 't.d.stageAndRename' || exit 1
ovl m4 "internal/tools/fs_edit.go" "$V/mut/m4/fs_edit.go"
go test -count=1 -overlay="$V/overlay-m4.json" -v -run '^TestFSEditRealTaskkill' ./internal/tools/ > "$V/logs/m4-on-r3-ac4b.log" 2>&1
echo "  m4 rc=$? (expected NONZERO: the target must stop being the original bytes)"
grep -E "(^|[[:space:]])--- FAIL" "$V/logs/m4-on-r3-ac4b.log" | sed 's/^/  M4-RED    /'
grep -E 'violated' "$V/logs/m4-on-r3-ac4b.log" | sed 's/^/  M4-CAUSE  /'
echo "== all rigs built =="
