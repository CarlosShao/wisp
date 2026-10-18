#!/usr/bin/env bash
# 328-v1 交件自证三把尺（逐枚现跑入库）。★第二稿：两把尺的写法自纠，具名登记在下面。
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 1
TICK='.scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md'
VERD='docs/evidence/s1/328-v1-verdict.md'
RULER='([T]ODO|[F]IXME|[T]BD|待[补]|占[位]|X{3})'
PC='internal/agent/approval/ticket146_liveapprovals_backing_test.go'
ANCHOR='5db0b538dad3bf2841adec834a130e4a6b28d8df'
OUT=.scratch/wisp/probes/328/v1/logs/self-evidence-three-rulers.txt

{
echo "date=$(date '+%F %T %z') leg=328-v1 裁决锚=$ANCHOR 第二稿（两把尺写法自纠，逐字在下）"
echo
echo '## 尺① 逐枚提交号＋逐枚 git show --name-only（严格前缀＝编排者收件的尺）'
CS=$(git log --grep='^票 328 终裁腿 328-v1·' --format='%h')
ALL=''
for c in $CS; do
  echo "--- $c | $(git log -1 --format='%s' "$c" | cut -c1-40)…"
  files=$(git show --name-only --format='' "$c" | grep -v '^$')
  echo "$files" | sed 's/^/      /'
  echo "      枚数=$(echo "$files" | wc -l)"
  ALL="${ALL}${files}"$'\n'
done
UNIQ=$(printf '%s' "$ALL" | grep -v '^$' | sort -u)
echo
echo '## 尺①补 全集越格自查（⊄ 三处写面者枚数，应 0）'
echo '  ★自纠①具名：第一稿这把尺写成 grep -vE "^(目录前缀|B|C)$" ⇒ 行尾锚把"目录前缀"那一支钉死，'
echo '    58 枚里 56 枚被我自己的错尺误报成越格＝尺的缺陷，⛔ 写面缺陷。改用逐枚 case 判：'
echo "  去重总枚数=$(printf '%s\n' "$UNIQ" | grep -c .)"
n_in=0; : > /tmp/outsofscope.txt
while IFS= read -r f; do
  [ -z "$f" ] && continue
  case "$f" in
    .scratch/wisp/probes/328/v1/*) n_in=$((n_in+1)) ;;
    docs/evidence/s1/328-v1-verdict.md) n_in=$((n_in+1)) ;;
    "$TICK") n_in=$((n_in+1)) ;;
    *) printf '%s\n' "$f" >> /tmp/outsofscope.txt ;;
  esac
done <<EOF2
$UNIQ
EOF2
echo "  落在三处写面内=$(printf '%s\n' "$UNIQ" | grep -c . | awk -v o="$(grep -c . /tmp/outsofscope.txt || echo 0)" '{print $1-o}')（越格 0 时＝总数−越格）"
echo "  越格枚数=$(grep -c . /tmp/outsofscope.txt || true)"
grep . /tmp/outsofscope.txt | sed 's/^/      ✗ /'
echo "  ★.go 命中（尺＝名册里以 .go 结尾者）=$(printf '%s\n' "$UNIQ" | grep -c '\.go$' || true)"
echo "    正控＝这把尺认得 .go 这个词尾：同一名册里 *.go.snap 命中 $(printf '%s\n' "$UNIQ" | grep -c '\.go\.snap$') 枚（第一程存的只读快照，在 probes/** 允许面内、⛔ 产码路径）"
echo "    逐枚对撞尺（快照字节 ⇄ 具名锚的真实 blob，证它们＝只读镜像、非产码）："
for s in $(printf '%s\n' "$UNIQ" | grep '\.go\.snap$'); do
  b=$(basename "$s" .snap)
  case "$b" in
    subagent_roster_197.go|subagent_blocked_220_test.go|subagent_corr_task_prefix_328_test.go) src="internal/panel/$b" ;;
    corr_roster_join_328_test.go) src="internal/agent/$b" ;;
    *) src="?"; ;;
  esac
  printf '      %s hash=%s  %s@锚=%s same=%s\n' "$b" "$(git hash-object "$s")" "$src" "$(git rev-parse "$ANCHOR:$src")" "$([ "$(git hash-object "$s")" = "$(git rev-parse "$ANCHOR:$src")" ] && echo YES || echo 'NO（快照存于锚之前/之后，字节差已具名）')"
done
echo
echo '## 尺② 证据件绝对路径＋wc -c＋行数＋空槽尺（逐字命令＋正控）'
echo "  绝对路径＝$(pwd)/$VERD"
printf '  %s bytes=%s lines=%s 帽=130 over=%s\n' "$VERD" "$(wc -c < "$VERD")" "$(wc -l < "$VERD")" "$([ "$(wc -l < "$VERD")" -le 130 ] && echo NO || echo YES)"
echo "\$ grep -nE '$RULER' $VERD"
echo "  命中=$(grep -cE "$RULER" "$VERD" || true)（0＝空槽尺；正控在下面一枚）"
echo "\$ grep -nE '$RULER' $PC   ← 正控：同一串尺打到已知存在的夹具件，先证这把尺会响"
echo "  正控命中行号=$(grep -nE "$RULER" "$PC" | cut -d: -f1 | tr '\n' ' ')行数=$(grep -cE "$RULER" "$PC" || true)"
echo "  写面 0 字节枚数（尺＝find -size 0，射程＝probes/328/v1/** ＋判语件）=$(find .scratch/wisp/probes/328/v1 "$VERD" -type f -size 0 | wc -l)"
echo '  ★自纠②具名：第一稿这把空槽尺写成不断形的 (TODO|FIXME|…|XXX)，把我自己那一行命令匹配进去（命中 1），'
echo '    且正控射程（我自己的 logs）当时真 0 命中⇒拿不到"尺会响"的证据。改成自避分词形＋到 internal/ 找到真正控。'
echo '  逐枚 wc 名册＝logs/deliverable-sizes.txt（入库笔 5615f14e）'
echo
echo '## 尺③ 票面（⛔ 翻框＋只追加一枚）'
echo "  未勾框=$(grep -c '^[[:space:]]*- \[ \]' "$TICK") 带勾框=$(grep -c '^[[:space:]]*- \[x\]' "$TICK") 票面总行=$(wc -l < "$TICK")"
echo "  ★AC 标题行尺：全档命中=$(grep -cE '\*\*AC#[0-4] ' "$TICK")；判据区射程（sed -n '29,36p'）命中=$(sed -n '29,36p' "$TICK" | grep -cE '\*\*AC#[0-4] ')"
echo '    ⇒ 那把全档尺被我自己的票尾 bullet 污染（正文里引用了"AC#2"字样）＝⛔ 用它数 AC 枚数；AC 枚数以判据区射程＝5 为准'
echo '    （起手指 §2 那把 - [ 形状尺另具名：五枚 AC 框＋两枚 dated bullet＝7，我那枚追加后＝8，dated bullet 同形所致）'
echo '  我那枚提交（f255d53c）numstat（追加行＝1，删除＝0）：'
git show --numstat --format='' f255d53c | sed 's/^/      /'
echo
echo '## 尺③补 ★编排者代提那枚（7f805583）：字节同原稿三串对撞，⛔ 判"件丢了"、⛔ 补第二枚同内容件'
SP='.scratch/wisp/probes/328/v1/logs/start-rulers-v2.txt'
echo "  git rev-parse 7f805583^:$SP = $(git rev-parse "7f805583^:$SP")"
echo "  git rev-parse 7f805583:$SP  = $(git rev-parse "7f805583:$SP")"
echo "  git rev-parse HEAD:$SP      = $(git rev-parse "HEAD:$SP")"
echo "  git hash-object <工作树>    = $(git hash-object "$SP")"
echo "  numstat（7f805583 对该件）＝$(git show --numstat --format='' 7f805583 -- "$SP")"
echo
echo '## 尺④ 决定性那把尺逐字（起手侧 ⇄ 交件侧，带 file:line＋rc=）'
echo '--- 去牙一发（派单 AC#2③ 那一发）红句逐字（截 packets 前）---'
cut -c1-160 .scratch/wisp/probes/328/v1/logs/mut-V1-redlines.txt | head -2 | sed 's/^/      /'
echo '--- 改前改后四数＋还原哈希 ---'
grep -E 'pristine_before|mut-V1-boundary:|restored:|restored_same_as_pristine|pristine_equals_anchor_blob|after-panel-220table-restored' .scratch/wisp/probes/328/v1/logs/color-driver-raw.log | cut -c1-150 | sed 's/^/      /'
echo '--- 成对第二发红句逐字（落点不同枚）---'
cut -c1-160 .scratch/wisp/probes/328/v1/logs/mut-V2-redlines.txt | head -3 | sed 's/^/      /'
echo '--- 门禁四数 rc（重打那发，GOTMPDIR unset）---'
grep -E 'gofumpt_rc|d22scan_rc|length_rc|vet_rc' .scratch/wisp/probes/328/v1/logs/gates-recheck-nogotmp.log | sed 's/^/      /'
echo '--- 红名册逐名作差（复跑两侧）---'
diff .scratch/wisp/probes/328/v1/logs/red-before-panel-rerun2.txt .scratch/wisp/probes/328/v1/logs/red-after-panel-rerun2.txt && echo '      REDROSTER_DIFF_RC=0 逐名同集（各 7 枚）'
echo '--- 成对整包四数逐枚点名 ---'
grep -E '^(before|after)-panel-full|^(before|after)-agent-full' .scratch/wisp/probes/328/v1/logs/color-driver-raw.log .scratch/wisp/probes/328/v1/logs/attribute-driver-raw.log 2>/dev/null | grep -E 'RUN=' | cut -c1-140 | sed 's/^/      /'
} > "$OUT" 2>&1
echo "written $OUT bytes=$(wc -c < "$OUT") lines=$(wc -l < "$OUT")"
