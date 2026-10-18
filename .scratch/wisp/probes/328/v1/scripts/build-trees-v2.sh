#!/usr/bin/env bash
# 328-v1 台面建树 第二枚（落点改判：机主 2026-10-18 10:1x 亲口新铁律＝盘符根目录零污染，
# 覆盖派单台面配方那两行 ⇒ 新建两棵仓外导出树只许落 D:\tmp\wisp\20261018\328v1\ 下）。
# 规矩块八件套 go.mod go.sum internal cmd tools scripts frontend models ＋仓根 docs/ ＋手工三件。
# ⛔ 复用上一程的 wisp-328r1-*-tree／w328v1a／w328v1b（只建不删＝机主的决定）。
# 每段目录名 <=8 字符 ASCII 无空格 ⇒ 8.3 别名＝名字本身（尺＝cmd /X dir 现量，逐字入库）。
set -u
REPO='D:/work/workspace/projects plans/Wisp'
EIGHT='go.mod go.sum internal cmd tools scripts frontend models'
BASE='D:/tmp/wisp/20261018/328v1'
BEFORE_ANCHOR='648eda04'   # 落地腿自己的起手锚（派单指名＝改前侧）
AFTER_ANCHOR="$(git -C "$REPO" log -1 --format=%H)"   # 我的入场锚＝现跑，⛔ 抄号
export GOPATH='D:\work\base\gopath'
export GOTMPDIR="$BASE/gotmp"

echo "date=$(date '+%F %T %z')"
echo "after_anchor=$AFTER_ANCHOR"
git -C "$REPO" log -1 --format='after_full=%H %ad %s' --date=iso "$AFTER_ANCHOR"
git -C "$REPO" log -1 --format='before_full=%H %ad %s' --date=iso "$BEFORE_ANCHOR"

mkdir -p "$GOTMPDIR" || echo "GOTMPDIR_FAIL"
echo "GOTMPDIR=$GOTMPDIR exists=$([ -d "$GOTMPDIR" ] && echo yes || echo NO)"

build() { # <dir> <anchor>
  local T="$1" A="$2"
  if [ -e "$T" ]; then echo "TREE_EXISTS $T ⇒ 换名，⛔ rm"; return 1; fi
  mkdir -p "$T" || { echo "MKDIR_FAIL $T"; return 1; }
  git -C "$REPO" archive "$A" $EIGHT | tar -x -C "$T" || { echo "ARCHIVE8_FAIL $A"; return 1; }
  git -C "$REPO" archive "$A" docs | tar -x -C "$T" || { echo "ARCHIVE_DOCS_FAIL $A"; return 1; }
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx"
  cp "$REPO"/build/*.dll "$T/build/" || echo "DLL_BUILD_FAIL"
  cp "$REPO"/third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" || echo "DLL_SHERPA_FAIL"
  cp -r "$REPO/frontend/dist" "$T/frontend/dist" || echo "DIST_FAIL"
  echo "TREE $T anchor=$A go_files=$(find "$T" -name '*.go' | wc -l)"
  echo "  docsPLAN=$([ -f "$T/docs/PLAN.md" ] && echo yes || echo NO)"
  echo "  dll=$(ls "$T"/build/*.dll 2>/dev/null | wc -l)/$(ls "$T"/third_party/sherpa-onnx/*.dll 2>/dev/null | wc -l) dist=$([ -d "$T/frontend/dist" ] && echo yes || echo NO)"
  echo "  gitblob=$(git -C "$REPO" rev-parse "$A:internal/panel/subagent_roster_197.go") ondiskblob=$(cd "$T" && git hash-object internal/panel/subagent_roster_197.go)"
  echo "  shortpath=$(/c/Windows/System32/cmd.exe //c 'for %I in ('"$(echo "$T" | sed 's|/|\\|g')"') do @echo %~sI' 2>/dev/null)"
  echo "  alias83=$(/c/Windows/System32/cmd.exe //c 'dir /X /AD D:\tmp\wisp\20261018\328v1' 2>/dev/null | grep -i "$(basename "$T")" | head -1)"
}

build "$BASE/before" "$BEFORE_ANCHOR"
build "$BASE/after" "$AFTER_ANCHOR"
echo "build_trees_v2_rc=0"
