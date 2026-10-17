#!/usr/bin/env bash
# 216-v3 交件自证三把尺：整把跑一遍，输出即交件套件（⛔ 手抄读数）。
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
T=.scratch/wisp/issues/216-menu-display-base-scrub-control-chars-and-two-sentences.md
V=docs/evidence/s1/216-v3-verdict.md
L=.scratch/wisp/probes/216/v3/logs
echo "########## 尺一：逐枚提交号＋逐枚 git show --name-only（写面只许三处）"
for c in 6feaa41d b0deec82 c3995bbc 18436906 2a3b9ff8; do
  echo "--- $c $(git log -1 --format=%s $c | cut -c1-90)"
  git show --name-only --format= $c | sed 's/^/    /'
  echo "    numstat=$(git show --numstat --format= $c | awk '{i+=$1; d+=$2} END{print i" 增/"d" 删"}')"
done
echo "########## 尺二：判语件绝对路径＋wc＋占位符尺（含正控）"
echo "abs=$(cygpath -w "$V")"
echo "worktree  wc=$(wc -l < "$V") bytes=$(wc -c < "$V")"
echo "committed-at-18436906 blob wc=$(git show 18436906:"$V" | wc -l) bytes=$(git show 18436906:"$V" | wc -c)"
echo "  ★自抓并已更正：那一笔的 commit-msg 写的是 10913 字节＝我量完正文后又改了一处引用行号（winsec 那半句的 file:line），落笔前未复量；盘上与 2a3b9ff8 票尾 bullet 用的都是上面这两行现量值。历史不改写，只在这里具名。"
for w in TODO 待填 占位 TBD FIXME; do printf 'placeholder[%s]=%s ' "$w" "$(grep -c "$w" "$V" || true)"; done; echo
echo "wide尺=$(grep -cE 'TODO|待填|占位|TBD|FIXME|XXX' "$V" || true)  正控A(同尺读本件在场词 成立)=$(grep -c '成立' "$V")  正控B(同尺族在 s1 目录看得见 TODO 的件数)=$(grep -rl 'TODO' docs/evidence/s1 | wc -l)"
echo "本腿五枚交件套件 0 字节核对："; find "$L" -maxdepth 1 -type f -size 0 | sed 's/^/  ZERO: /'; echo "  zero_count=$(find "$L" -maxdepth 1 -type f -size 0 | wc -l)"
echo "########## 尺三：决定性那把尺，起手侧 ⇄ 交件侧，逐字带 file:line"
echo "--- 起手侧＝编排者那对台面（本腿自己复跑，日志 surface-C-default-216only.log / surface-D-redirect-216only.log）"
grep -n 'name_claim_two_arm_216_test.go:105' "$L/surface-C-default-216only.log" | head -1
grep -n 'name_claim_two_arm_216_test.go:110' "$L/surface-C-default-216only.log" | head -1
echo "C 默认台面四数：RUN=$(grep -c '^=== RUN' $L/surface-C-default-216only.log || true) PASS_TOP=$(grep -c '^--- PASS' $L/surface-C-default-216only.log || true) FAIL_TOP=$(grep -c '^--- FAIL' $L/surface-C-default-216only.log || true) SKIP=$(grep -c '^--- SKIP' $L/surface-C-default-216only.log || true)"
echo "D 台面六枚仪器四数：RUN=$(grep -c '^=== RUN' $L/surface-D-redirect-216only.log || true) PASS_TOP=$(grep -c '^--- PASS' $L/surface-D-redirect-216only.log || true) FAIL_TOP=$(grep -c '^--- FAIL' $L/surface-D-redirect-216only.log || true) SKIP=$(grep -c '^--- SKIP' $L/surface-D-redirect-216only.log || true)"
echo "C 整包四数：$(grep -c '^=== RUN' $L/surface-C-default-full.log) / $(grep -c '^--- PASS' $L/surface-C-default-full.log) / $(grep -c '^--- FAIL' $L/surface-C-default-full.log) / $(grep -c '^--- SKIP' $L/surface-C-default-full.log)（多出来的那一枚红＝TestC26RewriteAccountIsConsumedAtEverySecurityLeg，两把台面都在册）"
echo "--- 交件侧＝突变重放三发＋还原哈希（日志 mutations-216v3-runner.log）"
grep -n 'LANDS\|restored \|tree=\|blob=' "$L/mutations-216v3-runner.log" | sed 's/^/  /'
echo "  M1 红句：$(grep -o 'allow_substring_pin_216_test.go:83: [^"]*' $L/mut-M1-D.log | head -1 | cut -c1-150)"
echo "  M1 附带 ⓑ：$(grep -o 'name_claim_two_arm_216_test.go:135: .*' $L/mut-M1-D.log | head -1)"
echo "  M2 红句（D 台面）：$(grep -o 'name_claim_two_arm_216_test.go:105: class changed.*= none' $L/mut-M2-D.log | head -1)"
echo "  M2 打在 C 台面（hostile 读 = A ⇒ 冒充牙证）：$(grep -o 'name_claim_two_arm_216_test.go:105: class changed.*= A' $L/mut-M2-Cdefault.log | head -1)"
echo "  R1 必答题红句：$(grep -o 'allow_substring_pin_216_test.go:68: [^"]*' $L/mut-R1-D.log | head -1 | cut -c1-160)"
echo "--- 交件侧＝假主目录两发（归因）"
echo "  发一 reason：$(grep -o 'non-overridable): [^\\"]*' $L/fakehome-aws-absent.log | sort -u | paste -sd'; ' -)"
echo "  发二 reason（补造 .aws 之后）：$(grep -o 'non-overridable): [^\\"]*' $L/fakehome-aws-present.log | sort -u | paste -sd'; ' -)"
echo "--- 交件侧＝CI 形台面（GOTMPDIR 落 %LOCALAPPDATA%）：$(grep -c '^=== RUN' $L/surface-CIshaped-gotmpdir.log) / $(grep -c '^--- PASS' $L/surface-CIshaped-gotmpdir.log) / $(grep -c '^--- FAIL' $L/surface-CIshaped-gotmpdir.log) / $(grep -c '^--- SKIP' $L/surface-CIshaped-gotmpdir.log)"
echo "########## 门禁四数 rc"
grep -E 'gofumpt_rc=|d22scan_rc=|d22scan_retry_rc=|pathlen_rc=|vet_rc=' "$L/gates-216v3.log"
echo "########## 票面：框数一枚没动＋本腿那一笔 numstat"
echo "unchecked=$(grep -c '^[[:space:]]*- \[ \]' $T || true) 逐名："
grep -n '^[[:space:]]*- \[ \]' "$T" | cut -c1-24
echo "checked=$(grep -c '^[[:space:]]*- \[x\]' $T || true)"
echo "2a3b9ff8 票那一笔：$(git show --numstat --format= 2a3b9ff8 -- $T)"
echo "########## 产码零字节：internal/risk 工作树 ⇄ HEAD"
git status --porcelain -- internal/risk; echo "porcelain_rc=$? (空＝零字节)"
echo "diff 7a94c3bd..HEAD on internal/risk：$(git diff --numstat 7a94c3bd HEAD -- internal/risk | paste -sd'; ' -)"
echo "########## 注：logs/logs/ 里是第一发（台件 HERE 路径写错那发）的副本，按规矩只建不删、未入库"
ls "$L/logs" | paste -sd' ' -
