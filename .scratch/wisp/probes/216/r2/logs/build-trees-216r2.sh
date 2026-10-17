#!/usr/bin/env bash
# 216-r2 台面建树：票 216 AC#6 两枚仪器的读数一律在仓外导出树上跑
# 用法: bash build-trees.sh <sha1:dirname1> [<sha2:dirname2> ...]
# 配方（派单「台面配方（现行完整版）」逐味）：八件套 + 仓根 docs/ + 手工补三件 + 落 D: 卷
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT=/d/tmp/leg216r2trees
EIGHT="go.mod go.sum internal cmd tools scripts frontend models"
GOTMP=/d/tmp/leg216r2gotmp
mkdir -p "$GOTMP" "$ROOT"
export GOTMPDIR='D:\tmp\leg216r2gotmp'
export GOPATH='D:\work\base\gopath'

build_tree () { # $1=sha $2=目录名
  local d="$ROOT/$2"
  rm -rf "$d"; mkdir -p "$d"
  ( cd "$REPO" && git archive "$1" $EIGHT ) | tar -x -C "$d"
  ( cd "$REPO" && git archive "$1" docs ) | tar -x -C "$d"
  mkdir -p "$d/build" "$d/third_party/sherpa-onnx" "$d/frontend/dist"
  cp "$REPO"/build/*.dll "$d/build/" 2>/dev/null
  cp "$REPO"/third_party/sherpa-onnx/*.dll "$d/third_party/sherpa-onnx/" 2>/dev/null
  cp -r "$REPO"/frontend/dist/. "$d/frontend/dist/" 2>/dev/null
  echo "tree $2 sha=$1 go_files=$(find "$d/internal" "$d/cmd" "$d/tools" -name '*.go' | wc -l) dll_build=$(ls "$d/build"/*.dll 2>/dev/null | wc -l) dll_sherpa=$(ls "$d/third_party/sherpa-onnx"/*.dll 2>/dev/null | wc -l) dist_entries=$(ls "$d/frontend/dist" 2>/dev/null | wc -l) docs_dir=$([ -d "$d/docs" ] && echo yes || echo no)"
}

for pair in "$@"; do
  sha="${pair%%:*}"; name="${pair#*:}"
  echo "########## 建树 $name <- $sha"
  build_tree "$sha" "$name"
  echo "########## build 门 $name"
  ( cd "$ROOT/$name" && go build ./internal/risk/ ; echo "build-$name rc=$?" )
done
echo "########## done  ROOT=$ROOT"
