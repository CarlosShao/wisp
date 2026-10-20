#!/bin/sh
# gw-1 sh/10 那把名册尺的**修正发**。
# ★具名自抓坏尺＝我上一发（sh/10 第 3⇄4 段）的 comm 作差⛔ 成立＝gofumpt 打的是 `D:/tmp/...` 正斜杠盘符形，
#   而我册子是 `.scratch/...` 仓相对形⇒两串不同源⇒comm 把 35 枚全判成"只在本地脏"＝假读数；
#   随后我在 pwsh 内联里补的那把尺又因 sed 表达式被 shell 吃掉而**整条失败**（stderr 两行 `unterminated s command`
#   ⇒ 名册文件被写成 0 行＝BLOB_DIRTY=0 是**尺坏**、⛔ 是读数＝公共块"报 0 命中前先跑正控＋看 stderr"那一味）。
# 本发＝在仓外用脚本文件跑（⛔ pwsh 内联），并把两侧路径**归一到同一个仓相对形**再作差。
set -u
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
B="$R/blobcheck"
cd "$B" || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
echo "== 归一化＝在 blobcheck 根里跑（相对路径）＋正控先行 =="
/d/work/base/gopath/bin/gofumpt.exe -l . > "$LOGS/blobform-relative-gofumpt.txt" 2>&1
echo "GOFUMPT_RC=$?  BYTES=$(wc -c < "$LOGS/blobform-relative-gofumpt.txt")"
grep -c '\.go' "$LOGS/blobform-relative-gofumpt.txt" | sed 's|^|GOFUMPT_LISTED=|' || true
gofmt -l . > "$LOGS/blobform-relative-gofmt.txt" 2>&1
echo "GOFMT_RC=$?  BYTES=$(wc -c < "$LOGS/blobform-relative-gofmt.txt")"
grep -c '\.go' "$LOGS/blobform-relative-gofmt.txt" | sed 's|^|GOFMT_LISTED=|' || true
echo
echo "== ★正控（证明这把尺在 blobcheck 上看得见东西）＝当场造一枚脏件 =="
mkdir -p "$B/posctrl"
printf 'package p\n\nfunc  F(){ }\n' > "$B/posctrl/posctl.go"
/d/work/base/gopath/bin/gofumpt.exe -l . > "$LOGS/blobform-posctrl.txt" 2>&1
echo "POSCTRL_RC=$?"
echo "POSCTRL_LISTED=$(grep -c '\.go' "$LOGS/blobform-posctrl.txt" || true)"
grep '\.go' "$LOGS/blobform-posctrl.txt" | sed 's|^|  POS |'
rm -f "$B/posctrl/posctl.go" 2>/dev/null || true
/d/work/base/gopath/bin/gofumpt.exe -l . > "$LOGS/blobform-posctrl-after.txt" 2>&1
echo "AFTER_REMOVE_LISTED=$(grep -c '\.go' "$LOGS/blobform-posctrl-after.txt" || true)（须＝POSCTRL_LISTED - 1＝尺跟着还原走）"
echo
echo "== 名册（blob 形仍脏＝CI 真红的那批；⛔ 拿枚数顶名册）=="
grep '\.go' "$LOGS/blobform-relative-gofumpt.txt" | tr '\\' '/' | sed 's|^\./||' | sort > "$LOGS/blob-dirty.names3"
echo "STILL_DIRTY_IN_BLOB=$(grep -c . "$LOGS/blob-dirty.names3" || true)"
grep . "$LOGS/blob-dirty.names3" | sed 's|^|  BDIRTY |'
echo
echo "== 作差（同形对同形：册子⇄blob 名册，都归一到仓相对）=="
sort "$LOGS/g4-35-paths.txt" > "$LOGS/paths35.norm"
echo "CLEAN_IN_BLOB(＝只在本地工作树脏＝CRLF 那一形)=$(comm -23 "$LOGS/paths35.norm" "$LOGS/blob-dirty.names3" | grep -c . || true)"
comm -23 "$LOGS/paths35.norm" "$LOGS/blob-dirty.names3" | sed 's|^|  LOCALONLY |'
echo "ODD_NAMES(尺冒出的不在 35 册里的名＝须 0)=$(comm -13 "$LOGS/paths35.norm" "$LOGS/blob-dirty.names3" | grep -c . || true)"
comm -13 "$LOGS/paths35.norm" "$LOGS/blob-dirty.names3" | sed 's|^|  ODD |'
echo
echo "== 生产面（非 .scratch/**）在 blob 形仍脏的枚数＝CI 判据那一问 =="
grep -v '^\.scratch/' "$LOGS/blob-dirty.names3" > "$LOGS/blob-dirty.prod" || true
echo "PRODUCTION_BLOB_DIRTY=$(grep -c . "$LOGS/blob-dirty.prod" || true)"
grep . "$LOGS/blob-dirty.prod" | sed 's|^|  PRODIRTY |'
echo
echo "== 两把独立尺在 blob 形上对撞 =="
grep '\.go' "$LOGS/blobform-relative-gofmt.txt" | tr '\\' '/' | sed 's|^\./||' | sort > "$LOGS/blob-gofmt.names3"
echo "ONLY_GOFUMPT=$(comm -23 "$LOGS/blob-dirty.names3" "$LOGS/blob-gofmt.names3" | grep -c . || true)"
comm -23 "$LOGS/blob-dirty.names3" "$LOGS/blob-gofmt.names3" | sed 's|^|  GONLY |'
echo "ONLY_GOFMT=$(comm -13 "$LOGS/blob-dirty.names3" "$LOGS/blob-gofmt.names3" | grep -c . || true)"
comm -13 "$LOGS/blob-dirty.names3" "$LOGS/blob-gofmt.names3" | sed 's|^|  MONLY |'
echo "BOTH=$(comm -12 "$LOGS/blob-dirty.names3" "$LOGS/blob-gofmt.names3" | grep -c . || true)"
