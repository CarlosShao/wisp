#!/bin/bash
# 339-v1 台件⑤：★同族普查尺（票 339 AC#0 后半 ＋ 必答 v）。
# 判据形状＝"设了 Windows-only 环境变量（USERPROFILE / APPDATA / LOCALAPPDATA）而该件没有任何
#   平台档"＝既无 //go:build 行、文件名也不带 _windows／_other／_posix／_unix 后缀。
# 尺甲＝internal/risk/**（票面那一格要的射程）；尺乙＝其余全包（必答 v）；
# 正控＝同一把尺先打在已知在场件上（USERPROFILE 在 334／339 两枚仪器里逐字在场）＝证尺灵。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
LOG="$REPO/.scratch/wisp/probes/339/v1/logs/v1-census-windows-env.txt"
PAT='Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA)"'

verdict_face() { # 余下参数＝pathspec 列表
  local f TAG SFX LINES
  for f in $(git grep -lE "$PAT" HEAD -- "$@" | sed 's/^HEAD://'); do
    TAG=$(grep -c '^//go:build' "$f" || true)
    case "$f" in
      *_windows.go|*_windows_test.go) SFX=windows-suffix ;;
      *_other.go|*_other_test.go|*_posix.go|*_posix_test.go|*_unix.go|*_unix_test.go) SFX=other-suffix ;;
      *) SFX=NO-SUFFIX ;;
    esac
    LINES=$(grep -nE "$PAT" "$f" | cut -d: -f1 | paste -sd, -)
    echo "FILE=$f build_tag_lines=$TAG name_suffix=$SFX setenv_lines=$LINES"
  done
}

{
  echo "# 尺＝git grep -lE '$PAT'（量的、HEAD 具名提交树；工作树与 HEAD 同串已现跑）"
  echo
  echo "== 正控 A（证尺灵）：三枚 Windows-only env 名在全仓 .go 的在场文件数 =="
  echo "files_with_USERPROFILE=$(git grep -l 'USERPROFILE' HEAD -- '*.go' | wc -l)"
  echo "files_with_LOCALAPPDATA=$(git grep -l 'LOCALAPPDATA' HEAD -- '*.go' | wc -l)"
  echo "files_with_APPDATA=$(git grep -l 'APPDATA' HEAD -- '*.go' | wc -l)"
  echo
  echo "== 正控 B：本尺在 internal/risk 的在场命中（应含 334/339 两枚仪器）=="
  verdict_face 'internal/risk/*.go'
  echo
  echo "== 尺甲 internal/risk/**：上面那份即全量；逐枚点名无档件 =="
  echo "risk_no_tag=$(verdict_face 'internal/risk/*.go' | grep -c 'build_tag_lines=0.*NO-SUFFIX' || true)"
  echo
  echo "== 尺乙（必答 v）internal/risk 之外全仓其余包 =="
  verdict_face '*.go' ':!internal/risk'
  echo
  echo "== 尺丙 生产码（非 _test）读 Windows-only env 的落点（只读数、⛔ 本票射程）=="
  git grep -nE 'Getenv\("(USERPROFILE|LOCALAPPDATA|APPDATA)"' HEAD -- '*.go' | grep -v '_test\.go' | sed 's/^HEAD://' || true
  echo
  echo "== 尺丁 全仓 _test 件里 t.Setenv 的环境变量名总表（防判据写窄）=="
  git grep -hoE 'Setenv\("[A-Z0-9_]+"' HEAD -- '*_test.go' | sort | uniq -c | sort -rn
} > "$LOG" 2>&1
cat "$LOG"
