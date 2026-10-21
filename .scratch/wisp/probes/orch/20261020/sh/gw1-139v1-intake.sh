#!/usr/bin/env bash
# 编排者收件尺 v5：对 gw-1（量色腿）与 139-v1（终裁腿）的交件做盘上复量。
cd "/d/work/workspace/projects plans/Wisp" || exit 9

echo "=== HEAD 与锚 ==="
git log -1 --format="HEAD=%H %ct"

echo
echo "=== ① 名册三把尺（gw-1 / 139-v1 / 353-v1）==="
for P in "gw-1 量色腿 gw-1" "票 139 终裁腿 139-v1" "353-v1"; do
  a=$(git log --format='%s' | grep -c "^$P")
  b=$(git log --oneline --grep="^$P" | wc -l)
  c=$(git log --format='%H%x09%s' | awk -F'\t' -v P="$P" 'index($2,P)==1{n++} END{print n+0}')
  d=$(git log --format='%H%x09%s' | awk -F'\t' -v P="$P" 'index($2,P)>0{n++} END{print n+0}')
  echo "[$P] 锚定=$a grep=$b 内存前缀=$c 任意位置=$d"
done

echo
echo "=== ② 逐枚 numstat（只增不减的验证：删除列必须为 0 或具名）==="
git log --format='COMMIT %h %s' --numstat a19d5bf2..HEAD 2>/dev/null | grep -E '^COMMIT|^[0-9]' | cut -c1-120 | head -60

echo
echo "=== ③ 139-v1 第 5 件三把尺（独立复跑，⛔ 采信它的数）==="
echo "-- (a) cmd/wisp 对 Compression 的命中（含测试）"
echo -n "  grep -n 'Compression' HEAD -- cmd/wisp  => "; git grep -n 'Compression' HEAD -- cmd/wisp | wc -l
echo -n "  非注释命中（去 ^..*// 形）            => "; git grep -n 'Compression' HEAD -- cmd/wisp | grep -vcE ':[0-9]+:[[:space:]]*(//|\*)'
echo "-- (b) BuildDiagnosticsBundle 的产码调用者"
echo -n "  全部命中 => "; git grep -n 'BuildDiagnosticsBundle' HEAD -- 'internal' 'cmd' | wc -l
git grep -n 'BuildDiagnosticsBundle' HEAD -- 'internal' 'cmd' | cut -c1-110
echo "-- (c) .logDir() 的调用者"
echo -n "  全部命中 => "; git grep -n '\.logDir()' HEAD -- 'internal' 'cmd' | wc -l
git grep -n '\.logDir()' HEAD -- 'internal' 'cmd' | cut -c1-110

echo
echo "=== ④ ★本 HEAD 对 ec826c63 的 .go 变化清单（139-v1 声称 5 枚全在 winsec/models）==="
git diff --name-only ec826c63..HEAD -- '*.go'
echo -n "  .go 枚数="; git diff --name-only ec826c63..HEAD -- '*.go' | wc -l

echo
echo "=== ⑤ agent 三件套 blob 是否未漂 ==="
for f in internal/agent/compress.go internal/agent/loop.go internal/agent/result.go; do
  echo "$(git rev-parse --short=8 HEAD:$f 2>/dev/null) $f"
done
echo "  139-v1 声称：4fcd9a32 / ae8c25c7 / 2357dca8"

echo
echo "=== ⑥ 钉死的锚 blob 复算（139-v1 声称 16e153d8…2f50）==="
git cat-file blob a18f801d:docs/evidence/s1/139-adversarial-acceptance.md | sha256sum | cut -c1-24

echo
echo "=== ⑦ 判语件首尾与关键行的盘上证据（痕＝compress.go:240）==="
sed -n '240p' internal/agent/compress.go | cut -c1-140
echo -n "  Info(history compressed) 两把独立形："
echo -n "    grep 命中="; grep -rn 'history compressed' --include='*.go' internal | wc -l
echo -n "    git grep HEAD 命中="; git grep -c 'history compressed' HEAD -- 'internal/agent/compress.go' | wc -l

echo
echo "=== ⑧ 四枚票面的未勾/已勾现量（翻框前基线）==="
for N in 139 115 352 321 353; do
  f=$(ls .scratch/wisp/issues/${N}-*.md 2>/dev/null | head -1)
  [ -z "$f" ] && continue
  echo "  票 ${N} $(basename "$f") 未勾=$(grep -c '^- \[ \]' "$f") 已勾=$(grep -c '^- \[x\]' "$f")"
done

echo
echo "=== ⑨ ★115 那枚漂着的票面（未入库的别腿写面）到底改了什么 ==="
git diff --numstat -- '.scratch/wisp/issues/115-*.md'
git diff -- '.scratch/wisp/issues/115-*.md' | sed -n '1,40p' | cut -c1-140

echo
echo "=== ⑩ tracked 脏（逐枚名，⛔ 计数）==="
git status --porcelain | grep -E '^ ?M' | cut -c1-100
