#!/usr/bin/env bash
# ban #8 分母两枚（改前/改后）：同一把仪器、两枚【完整】导出树，只差本腿两枚新件。
# ⛔ 部分导出树：八件套缺 design/ 会让 ban #8 那枚射程被 GUARD 拒（exit 2＝射程拒绝，不是红）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
HERE="$(cd "$(dirname "$0")" && pwd)"
ROOT=/d/tmp/leg216r2trees
export GOPATH='D:\work\base\gopath'
mkdir -p "$ROOT"
for pair in "72b65c61:wisp216r2-full-before" "671d3585:wisp216r2-full-after"; do
  sha="${pair%%:*}"; name="${pair#*:}"; d="$ROOT/$name"
  rm -rf "$d"; mkdir -p "$d"
  ( cd "$REPO" && git archive "$sha" ) | tar -x -C "$d"
  echo "== $name sha=$sha tracked-export go_files_total=$(find "$d" -name '*.go' | wc -l) design=$([ -d "$d/design" ] && ls "$d/design" | wc -l || echo missing) docs=$([ -d "$d/docs" ] && echo yes || echo no)"
  ( cd "$d/tools/d22scan" && go run . -root "$d" 2>&1 | grep -E '^d22scan: scope|^d22scan: examined|clean|VIOLATION|GUARD' ; echo "scan-$name rc=${PIPESTATUS[0]}" )
done
echo "== 两树文件清单作差（只差本腿两枚新件，逐名）"
diff <( cd "$ROOT/wisp216r2-full-before" && find internal cmd -name '*.go' | sort ) \
     <( cd "$ROOT/wisp216r2-full-after"  && find internal cmd -name '*.go' | sort )
