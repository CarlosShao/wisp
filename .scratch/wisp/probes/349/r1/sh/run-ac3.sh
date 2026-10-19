#!/bin/sh
# 349-r1 AC#3：'注释点名的符号不存在' 该不该上通用仪器＝我自己独立跑一把（⛔ 抄 343-r1/343-v1 的数）。
# 做法：把我自己那把三形尺（SA 驼峰／SB 调用形／SC 声明形）在【三枚产码文件】上的非 GREEN 行逐枚
# 人工读一遍，标 TRUE（注释确实在撒"包内有这符号"的谎）／FP（撒不了谎的那一类），再算 precision。
# 分类表写在这里＝判据摊开给人复核，⛔ 藏在脚本外。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
C=.scratch/wisp/probes/349/r1/logs/census
L=.scratch/wisp/probes/349/r1/logs/ac3
mkdir -p "$L" /d/tmp/wisp/20261019/349r1/logs/ac3
{
echo "### AC#3 misfire surface, measured by 349-r1 on its own gauge"
echo "== raw non-GREEN rows, the three production files, both sides =="
for side in base del; do
  echo "-- side=$side --"
  awk '/^-- /{f=$2} /^RED/{print f"  " $0}' "$C/ac0-$side-prodonly.txt"
done
echo ""
echo "== 手工逐枚分类（TRUE＝注释真在断言'包内有此符号'且⛔ 有；FP＝词面与假话不可分／根本不是断言）=="
cat <<'TBL'
side=START(base)
  blacklist.go:310  SB fires        TRUE   票 349 名册第①枚（真身在别包 config/unwired.go:54 的字段）
  provenance.go:54  SA taskScopeID  TRUE   名册第②枚（全仓⛔ 定义＝散文占位）
  provenance.go:427 SA disposalScope TRUE  名册第③枚（真身＝别包类型 internal/plugin.DisposalScope）
  provenance.go:427 SB Defer        FP     Defer 真是 internal/plugin/disposal.go:186 的方法；谎在 disposalScope 不在 Defer
  provenance.go:58  SB DEFERRED     FP     DEFERRED(C25-loop-wiring)＝本仓标记约定（AGENTS.md 1.1 节要求它与 SPEC-12 第 5 节 1:1 对得上）
  provenance.go:365 SB CloseScope   FP     句子自己写 "Before it, ... was TWO functions — OpenScope(id) and CloseScope(id)"＝票 160 故意删的历史名，词面与假话不可分
  provenance.go:534 SC Test         FP     散文里嵌的一条命令 `grep -n "^func Test"`＝被声明形尺当成"func Test 的定义声明"
  provenance.go:1077 SC escape      FP     括注 "(B-2 type escape)"＝名词短语，不是声明
  syncdirs.go:453   SB Complete     TRUE   名册第④枚（真名 SyncDetectionComplete:472＝纯小写调用形）
  syncdirs.go:249   SA macOS        FP     平台名
  syncdirs.go:15    SB Nutstore     FP     厂商名（后紧跟全角左括号"坚果云"才被切成调用形）
  syncdirs.go:59    SC the          FP     散文冠词（"env var the sync client sets"）
side=DELIVERY(del)
  provenance.go:58  SB DEFERRED     FP     同上
  provenance.go:365 SB CloseScope   FP     同上
  provenance.go:534 SC Test         FP     同上
  provenance.go:1077 SC escape      FP     同上
  syncdirs.go:249   SA macOS        FP     同上
  syncdirs.go:15    SB Nutstore     FP     同上
  syncdirs.go:59    SC the          FP     同上
TBL
echo ""
echo "== precision（START：TRUE=4 / 非GREEN行=12；DELIVERY：TRUE=0 / 非GREEN行=7）=="
echo "-- 逐枚复核这两把算术（⛔ 我手算，尺自己数一遍）--"
awk '/^RED/{n++} END{printf "start_non_green_rows=%s\n", n}' "$C/ac0-base-prodonly.txt"
awk '/^RED/{n++} END{printf "del_non_green_rows=%s\n", n}' "$C/ac0-del-prodonly.txt"
echo ""
echo "== 三形各自的命中构成（全 package 两侧，含 _test.go＝另一维读数）=="
grep '^shape=' "$C/ac0-base-full.txt" | sed 's/^/  base /'
grep '^shape=' "$C/ac0-del-full.txt"  | sed 's/^/  del  /'
echo ""
echo "== SC（声明形）单独摊：它是'放宽切词'的那一味，看它抓到什么 =="
grep ' SC ' "$C/ac0-base-full.txt" | sed 's/^/  /'
echo "-- SC 全 package：14 枚命中里 GREEN-pkg＝0 ⇒ 这一形结构性地只会捞散文（本腿独立复现 343-v1 的放宽即不稳）--"
echo ""
echo "== 有界词表形（票 350 那一族 ^a[A-Z]）另开一把尺＝见 sh/run-ac3-wordlist.sh 与 logs/ac3/ac3-wordlist.log =="
} > "$L/ac3-misfire.log" 2>&1
echo "wrote=$L/ac3-misfire.log rc=$?"
