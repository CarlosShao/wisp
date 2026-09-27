#!/usr/bin/env bash
# 162-v1 acceptance rig builder + runner. Writes ONLY into
# .scratch/wisp/probes/162/v1/** ; production files are never edited - every
# mutation reaches the compiler through `go test -overlay`.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 1
ABS="D:/work/workspace/projects plans/Wisp"
V=".scratch/wisp/probes/162/v1"
mkdir -p "$V/mut" "$V/logs"

sub() { # sub <file> <old> <new>  -> exactly one replacement, else FATAL
  local f="$1"
  OLD="$2" NEW="$3" perl -0777 -i -pe 'our $c; $c += s/\Q$ENV{OLD}\E/$ENV{NEW}/g; END{print STDERR "COUNT=".($c//0)."\n"}' "$f" 2>"$V/logs/.cnt"
  local c; c=$(cat "$V/logs/.cnt")
  echo "  $f : $c"
  [ "$c" = "COUNT=1" ] || { echo "FATAL: expected exactly 1 replacement in $f ($c)"; return 1; }
}

proof() { # proof <file> <must-be-present> <must-be-absent> - the overlay bytes that reach the compiler
  local f="$1" pres="$2" ab="$3"
  local nh oh
  nh=$(grep -c -F "$pres" "$f"); oh=$(grep -c -F "$ab" "$f")
  echo "  PROOF $f present-hits=$nh absent-hits=$oh"
  [ "$nh" -ge 1 ] && [ "$oh" -eq 0 ] || { echo "FATAL: overlay text not applied in $f"; return 1; }
}

ovl() { # ovl <name> <entries...> ; entries are "virtual=real" repo-relative paths
  local name="$1"; shift
  { echo '{ "Replace": {'; local first=1 e k v
    for e in "$@"; do
      k="${e%%=*}"; v="${e#*=}"
      [ $first -eq 1 ] || echo ','
      first=0
      printf '    "%s/%s": "%s/%s"' "$ABS" "$k" "$ABS" "$v"
    done
    echo ''; echo '  } }'
  } > "$V/overlay-$name.json"
  echo "  wrote $V/overlay-$name.json"
}

# ---------------------------------------------------------------- mutations
echo "== MUT-1 ambiguity refusal removed (AC#2 判据3) =="
mkdir -p "$V/mut/m1"; cp internal/tools/fs_edit.go "$V/mut/m1/fs_edit.go"
sub "$V/mut/m1/fs_edit.go" '		case count > 1:
			return refusal(
				fmt.Sprintf("第 %d 枚编辑的 old 在文件里命中 %d 处（多于 1 处即为歧义，不取第一处）", i+1, count),
				refusalUniqueWording), nil' '		case count > 1:
			// MUT-1: ambiguity refusal removed, first occurrence wins.' || exit 1
proof "$V/mut/m1/fs_edit.go" '// MUT-1: ambiguity refusal removed' 'refusalUniqueWording), nil' || exit 1
ovl m1 "internal/tools/fs_edit.go=$V/mut/m1/fs_edit.go"

echo "== MUT-2 overlap check neutered (AC#2 判据4 all-or-nothing) =="
mkdir -p "$V/mut/m2"; cp internal/tools/fs_edit.go "$V/mut/m2/fs_edit.go"
sub "$V/mut/m2/fs_edit.go" '			if matches[i].start < matches[j].stop && matches[j].start < matches[i].stop {' '			if false && matches[i].start < matches[j].stop && matches[j].start < matches[i].stop { // MUT-2 overlap check disabled' || exit 1
proof "$V/mut/m2/fs_edit.go" 'if false && matches[i].start' 'if matches[i].start < matches[j].stop &&' || exit 1
ovl m2 "internal/tools/fs_edit.go=$V/mut/m2/fs_edit.go"

echo "== MUT-3 output forced to CRLF (AC#3 endings are load-bearing) =="
mkdir -p "$V/mut/m3"; cp internal/tools/fs_edit.go "$V/mut/m3/fs_edit.go"
sub "$V/mut/m3/fs_edit.go" '	updated := b.String()' '	updated := strings.ReplaceAll(strings.ReplaceAll(b.String(), "\r\n", "\n"), "\n", "\r\n") // MUT-3 unify to CRLF' || exit 1
proof "$V/mut/m3/fs_edit.go" 'MUT-3 unify to CRLF' '	updated := b.String()' || exit 1
ovl m3 "internal/tools/fs_edit.go=$V/mut/m3/fs_edit.go"

echo "== MUT-4 temp+rename removed from fs.edit (AC#4 摘掉那一味) =="
mkdir -p "$V/mut/m4"; cp internal/tools/fs_edit.go "$V/mut/m4/fs_edit.go"
sub "$V/mut/m4/fs_edit.go" '	res := t.d.stageAndRename(ctx, target, strings.NewReader(updated), onUpdate)
	res.Origin = target' '	se := t.d.ledger()
	f, oerr := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // MUT-4 no staging file, direct write
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
proof "$V/mut/m4/fs_edit.go" 'MUT-4 no staging file' 't.d.stageAndRename' || exit 1
ovl m4 "internal/tools/fs_edit.go=$V/mut/m4/fs_edit.go"

echo "== MUT-5 fs.write starts reporting the delta (AC#1 falsifiability) =="
mkdir -p "$V/mut/m5"; cp internal/tools/fs_write.go "$V/mut/m5/fs_write.go"
sub "$V/mut/m5/fs_write.go" '	if err := os.Rename(tmpName, target); err != nil {' '	prev := 0
	if fi, serr := os.Stat(target); serr == nil {
		prev = int(fi.Size()) // MUT-5 learn the pre-write size BEFORE the rename
	}
	if err := os.Rename(tmpName, target); err != nil {' || exit 1
sub "$V/mut/m5/fs_write.go" '		Text: fmt.Sprintf("已写入 %s（%d 字节，临时文件+原子重命名）",
			target, written),' '		Text: fmt.Sprintf("已写入 %s（%d 字节，比原文件少 %d 字节，临时文件+原子重命名）",
			target, written, prev-written), // MUT-5 the missing instrument, added' || exit 1
proof "$V/mut/m5/fs_write.go" 'MUT-5 the missing instrument' '已写入 %s（%d 字节，临时文件+原子重命名）' || exit 1
ovl m5 "internal/tools/fs_write.go=$V/mut/m5/fs_write.go"

echo "== overlay for my own AC#1/AC#3 cells =="
ovl mine "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go"
ovl m5mine "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go" \
            "internal/tools/fs_write.go=$V/mut/m5/fs_write.go"

# combined: mutation + my acceptance cells, so one run reads both rulers
ovl m1c "internal/tools/fs_edit.go=$V/mut/m1/fs_edit.go" \
        "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go"
ovl m2c "internal/tools/fs_edit.go=$V/mut/m2/fs_edit.go" \
        "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go"
ovl m3c "internal/tools/fs_edit.go=$V/mut/m3/fs_edit.go" \
        "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go"
ovl m4c "internal/tools/fs_edit.go=$V/mut/m4/fs_edit.go" \
        "internal/tools/probe162v1_ac13_test.go=$V/probe162v1_ac13_test.go"
echo "== rig built =="
