#!/bin/bash
# 339-v1 台件⑤b：同族普查尺的★收窄判据形（必答 v／票 339 AC#0 后半）。
# 本票那枚缺陷的真形不是"设了 Windows-only env"，而是：
#   设了 Windows-only env（USERPROFILE/APPDATA/LOCALAPPDATA）
#   ＋ 同件不设 POSIX 对应量（HOME）
#   ＋ 该件没有任何平台档（无 //go:build 行、文件名无 _windows/_other/_posix/_unix 后缀）
#   ＋ 件内无 runtime.GOOS 现场守卫
# ⇒ 才叫"假主目录只有 Windows 有世界"。逐枚给 file:line，0 命中必须带正控。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
LOG="$REPO/.scratch/wisp/probes/339/v1/logs/v1-census-narrowed-shape.txt"
{
  echo "== 尺＝逐枚验四条件（量的、HEAD 具名树）=="
  echo "正控＝票 339 落码前的 anchor_miss_reason_334_test.go（在 e6fd0f50 上＝四条件全中的原形）"
  echo
  echo "--- 正控（旧锚 e6fd0f50 的 334 件）：设 env 逐行 / 有 HOME 吗 / 有档吗 / 有 GOOS 守卫吗"
  git show 'e6fd0f50:internal/risk/anchor_miss_reason_334_test.go' > /d/tmp/wisp/20261018/339v1/ctrl-334-old.go
  C=/d/tmp/wisp/20261018/339v1/ctrl-334-old.go
  echo "  setenv_lines=$(grep -nE 'Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA|HOME)"' "$C" | cut -d: -f1 | paste -sd, -)"
  echo "  sets_HOME=$(grep -cE 'Setenv\("HOME"' "$C" || true) (正控＝1 ⇒ 旧锚那件本来就在设 HOME ⇒ 票面那句 USERPROFILE 诊断作废)"
  echo "  build_tag_lines=$(grep -c '^//go:build' "$C" || true) (正控＝0 ⇒ 四条件全中 ⇒ 票 339 那两枚红就是这么来的)"
  echo
  echo "--- 全仓 .go 逐枚筛（HEAD）：设 Windows-only env 且件内不设 HOME、无档、无 GOOS 守卫"
  for f in $(git grep -lE 'Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA)"' HEAD -- '*.go' | sed 's/^HEAD://'); do
    W=$(grep -cE 'Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA)"' "$f" || true)
    H=$(grep -cE 'Setenv\("HOME"' "$f" || true)
    T=$(grep -c '^//go:build' "$f" || true)
    S=$(case "$f" in *_windows.go|*_windows_test.go|*_other.go|*_other_test.go|*_posix.go|*_posix_test.go|*_unix.go|*_unix_test.go) echo 1;; *) echo 0;; esac)
    G=$(grep -c 'runtime\.GOOS' "$f" || true)
    if [ "$H" -eq 0 ] && [ "$T" -eq 0 ] && [ "$S" -eq 0 ] && [ "$G" -eq 0 ]; then
      echo "HIT-SAME-SHAPE file=$f winenv_lines=$W lines=$(grep -nE 'Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA)"' "$f" | cut -d: -f1 | paste -sd, -)"
    else
      echo "no  file=$f winenv=$W setsHOME=$H tag=$T suffix=$S goos=$G"
    fi
  done
  echo
  echo "--- 同族另一形（只读参考）：设 HOME 不设 Windows-only env 的无档件（POSIX-only 世界）"
  for f in $(git grep -lE 'Setenv\("HOME"' HEAD -- '*.go' | sed 's/^HEAD://'); do
    W=$(grep -cE 'Setenv\("(USERPROFILE|APPDATA|LOCALAPPDATA)"' "$f" || true)
    T=$(grep -c '^//go:build' "$f" || true)
    S=$(case "$f" in *_windows.go|*_windows_test.go|*_other.go|*_other_test.go) echo 1;; *) echo 0;; esac)
    [ "$W" -eq 0 ] && [ "$T" -eq 0 ] && [ "$S" -eq 0 ] && echo "POSIX-side-only file=$f lines=$(grep -nE 'Setenv\("HOME"' "$f" | cut -d: -f1 | paste -sd, -)"
  done
} > "$LOG" 2>&1
cat "$LOG"
