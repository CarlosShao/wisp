#!/usr/bin/env bash
# 328-v1 台面建树（仓外导出树，两枚成对）：规矩块八件套 + 仓根 docs/ + 手工三件。
# 树放 D: 卷、目录名 <=8 字符无空格（8.3 别名＝名字本身）；GOTMPDIR 先 mkdir -p。
# ⛔ 在共享工作树跑 go test / go build。
set -u
REPO='D:/work/workspace/projects plans/Wisp'
EIGHT='go.mod go.sum internal cmd tools scripts frontend models'
BEFORE_ANCHOR='648eda04'   # 腿的起手锚（现跑：cat-file -t = commit；merge-base --is-ancestor = YES）
AFTER_ANCHOR='78d3b46a'    # 我的入场锚（git log -1 现认）
GB='D:\work\base\gopath'

build() { # <dir> <anchor>
  local T="$1" A="$2"
  if [ -e "$T" ]; then echo "TREE_EXISTS $T (临时件只建不删 ⇒ 换名，⛔ rm)"; return 1; fi
  mkdir -p "$T" || { echo "MKDIR_FAIL $T"; return 1; }
  git -C "$REPO" archive "$A" $EIGHT | tar -x -C "$T" || { echo "ARCHIVE8_FAIL $A"; return 1; }
  git -C "$REPO" archive "$A" docs | tar -x -C "$T" || { echo "ARCHIVE_DOCS_FAIL $A"; return 1; }
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx"
  cp "$REPO"/build/*.dll "$T/build/" || echo "DLL_BUILD_FAIL"
  cp "$REPO"/third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" || echo "DLL_SHERPA_FAIL"
  cp -r "$REPO/frontend/dist" "$T/frontend/dist" || echo "DIST_FAIL"
  local n
  n=$(find "$T" -name '*.go' | wc -l)
  echo "TREE $T anchor=$A go_files=$n"
  echo "  docsPLAN=$([ -f "$T/docs/PLAN.md" ] && echo yes || echo NO)"
  echo "  dll=$(ls "$T"/build/*.dll 2>/dev/null | wc -l)/$(ls "$T"/third_party/sherpa-onnx/*.dll 2>/dev/null | wc -l) dist=$([ -d "$T/frontend/dist" ] && echo yes || echo NO)"
  echo "  gitblob=$(git -C "$REPO" rev-parse "$A:internal/panel/subagent_roster_197.go" 2>/dev/null) ondiskblob=$(cd "$T" && git hash-object internal/panel/subagent_roster_197.go 2>/dev/null)"
}

export GOPATH="$GB"
export GOTMPDIR='D:\tmp\328v1-gotmp'
mkdir -p "$GOTMPDIR" || echo "GOTMPDIR_FAIL"
echo "GOTMPDIR=$GOTMPDIR exists=$([ -d "$GOTMPDIR" ] && echo yes || echo NO)"
build 'D:/work/workspace/w328v1b' "$BEFORE_ANCHOR"
build 'D:/work/workspace/w328v1a' "$AFTER_ANCHOR"
