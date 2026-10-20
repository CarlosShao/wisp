#!/bin/bash
# 321-r1 终局尺 v2（v1 报"139/139 全 MISMATCH"＝尺自己坏了：--format 里 %(objecttype) 占第 1 列，
# v1 用 awk '{print $1}' 取到的是 "blob" 这个字、⛔ blob 号 ⇒ 恒不等。两枚单点抽查（winsec.go／判据件
# hash-object ⇄ rev-parse）都对得上 ⇒ 按公共块"两把尺矛盾时先假定新那把坏"修尺重跑。）
# 这把尺证明：成对台面的 Go 字节 == 交件后最终 HEAD 的 Go 字节（本腿跑完后 HEAD 漂过 3 笔，
# 那 3 笔逐枚 name-only 只含 docs/票面/派单存档＝零 Go 字节 ⇒ 台面等价，量出来而不是说）。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
R=/d/tmp/wisp/20261020/321r1
echo "HEAD now = $(git rev-parse --short HEAD)"
bad=0; n=0
while read -r h p; do
  n=$((n+1))
  got=$(git hash-object "$R/treeB/$p" 2>/dev/null)
  [ "$got" = "$h" ] || { echo "MISMATCH $p want=$h got=$got"; bad=$((bad+1)); }
done < <(git ls-tree -r HEAD --format='%(objectname) %(path)' | grep ' cmd/wisp/')
echo "cmd/wisp: checked=$n mismatches=$bad"
bad2=0; n2=0
while read -r h p; do
  n2=$((n2+1))
  got=$(git hash-object "$R/treeM/$p" 2>/dev/null)
  [ "$got" = "$h" ] || { echo "MISMATCH(treeM) $p"; bad2=$((bad2+1)); }
done < <(git ls-tree -r HEAD --format='%(objectname) %(path)' | grep ' cmd/wisp/')
echo "treeM(已还原)：checked=$n2 mismatches=$bad2"
echo "--正控：⛔ 动 treeB（它是成对读数的基线）⇒ 另取一枚 throwaway 副本 treeX 故意改坏，这把尺必须报 MISMATCH--"
mkdir -p "$R/treeX/cmd/wisp"
cp "$R/treeB/cmd/wisp/versioninfo_other.go" "$R/treeX/cmd/wisp/versioninfo_other.go"
printf '\n// 321r1 ruler positive control, deliberately not HEAD bytes\n' >> "$R/treeX/cmd/wisp/versioninfo_other.go"
gx=$(git hash-object "$R/treeX/cmd/wisp/versioninfo_other.go")
gh=$(git rev-parse HEAD:cmd/wisp/versioninfo_other.go)
echo "treeX hash=$gx"; echo "HEAD  hash=$gh"
[ "$gx" = "$gh" ] && echo "CONTROL FAILED (ruler cannot see a坏 byte)" || echo "CONTROL OK：尺看得见坏字节"
echo "--treeB 全程未动的证据：逐枚 hash-object == HEAD 且 sha 与成对读数时同串--"
sha256sum "$R/treeB/cmd/wisp/versioninfo_other.go" "$R/treeB/cmd/wisp/resident_windows.go"
