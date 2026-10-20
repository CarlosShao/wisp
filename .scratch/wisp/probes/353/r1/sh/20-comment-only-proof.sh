#!/bin/sh
# 353-r1 判据形状尺（公共块 判据形状★那条）＝"只动注释"怎么证＝
#   ① 工作树 ⇄ HEAD 的 diff -U0 里每条 ^[+-] 行都必须以 // 开头或为注释续行（可执行行变化＝0）
#   ② 用例名集 symdiff 0 ③ 字符串字面量排序后 sha256 两侧同串 ④ 行尾注释 symdiff 0
#   ＋ 第五把＝取 HEAD blob、仓外只跑格式化工具到不动点、所得哈希 == 入库 blob（＝交件字节≡工具输出）
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
mkdir -p "$L"
{
  echo "== @ $(date '+%F %T %z') HEAD=$(git log -1 --format=%h) =="
  for f in internal/winsec/pinnedleaf_342_windows.go internal/winsec/pinnedleaf_342_unix.go \
           internal/winsec/nametopup_353_windows_test.go; do
    echo "-- $f --"
    TOT=$(git diff -U0 -- "$f" | grep -cE '^[+-]' || true)
    HEADER=$(git diff -U0 -- "$f" | grep -cE '^(\+\+\+|---)' || true)
    NONCOMMENT=$(git diff -U0 -- "$f" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' \
      | grep -vcE '^[+-][[:space:]]*(//|$)' || true)
    echo "   DIFF_CHANGED_LINES=$TOT HEADER_LINES=$HEADER NON_COMMENT_HUNK_LINES=$NONCOMMENT（★期望 0＝⛔ 可执行行变化）"
    git diff -U0 -- "$f" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vE '^[+-][[:space:]]*(//|$)' | sed 's/^/     BAD| /' || true
  done
  echo "== 用例名集 symdiff（HEAD blob ⇄ 工作树）=="
  for f in internal/winsec/nametopup_353_windows_test.go internal/winsec/ancestorwindow_342r2_windows_test.go; do
    git cat-file blob "HEAD:$f" | grep -oE '^func Test[A-Za-z0-9_]+' | sort > "$L/names-head.txt"
    grep -oE '^func Test[A-Za-z0-9_]+' "$f" | sort > "$L/names-wt.txt"
    echo "   $f SYMDIFF=$(comm -3 "$L/names-head.txt" "$L/names-wt.txt" | wc -l)"
  done
  echo "== 字符串字面量集合 sha256（HEAD blob ⇄ 工作树，排序后）=="
  for f in internal/winsec/nametopup_353_windows_test.go; do
    a=$(git cat-file blob "HEAD:$f" | grep -oE '"[^"]*"' | sort | sha256sum | cut -d' ' -f1)
    b=$(grep -oE '"[^"]*"' "$f" | sort | sha256sum | cut -d' ' -f1)
    echo "   $f HEAD_STRINGS_SHA=$a WT_STRINGS_SHA=$b SAME=$([ "$a" = "$b" ] && echo YES || echo NO)"
  done
  echo "== 格式化工具到不动点（仓外，⛔ 改工作树）=="
  git cat-file blob "HEAD:internal/winsec/pinnedleaf_342_windows.go" > "$L/fmtcheck-in.go"
  cp "$L/fmtcheck-in.go" "$L/fmtcheck-run.go"
  /d/work/base/gopath/bin/gofumpt.exe -w "$L/fmtcheck-run.go"
  echo "   GOOFMPT_IDEMPOTENT_AT_HEAD=$([ "$(sha256sum < "$L/fmtcheck-in.go")" = "$(sha256sum < "$L/fmtcheck-run.go")" ] && echo YES || echo NO)"
} > "$L/comment-only-proof.txt" 2>&1
cat "$L/comment-only-proof.txt"
