#!/usr/bin/env bash
# 216-v2 裁决腿（只读）自己重放的尺名册。⛔ 任何 go 命令（无 go test/build/vet/list/env）。
# 全程读**具名提交树**两枚：T1=44f51fc8（216-a2 的锚，验可复现）与 T2=feab45e3（本腿起手锚，验今天）。
# 每把报 0 的尺都现场配正控；stderr 一律并进 stdout 或另存（⛔ 2>/dev/null 吞报错）。
set -u
cd "$(git rev-parse --show-toplevel)"
T1=44f51fc8
T2=feab45e3
PROD=':(exclude)*_test.go'
g() { # g <tree> <flags> <pattern> <paths...> —— 逐命中型
  local tree="$1" fl="$2" pat="$3"; shift 3
  echo "    pat=$pat"; echo "    rng=$*"
  # shellcheck disable=SC2086
  git grep $fl "$pat" "$tree" -- "$@" 2>&1 | sed "s#^$tree:##"
  echo "    rc=${PIPESTATUS[0]}"
}
n() { # n <tree> <pattern> <paths...> —— 计数型（行数；stderr 单独留痕）
  local tree="$1" pat="$2"; shift 2
  echo "    pat=$pat"; echo "    rng=$*"
  git grep -nE "$pat" "$tree" -- "$@" >/dev/null 2>/tmp/e216v2; local rc=$?
  local hits out; out=$(git grep -nE "$pat" "$tree" -- "$@" 2>/dev/null); hits=$(printf '%s' "$out" | grep -c .)
  echo "    count=$hits rc=$rc"
  if [ -s /tmp/e216v2 ]; then echo "    STDERR=$(cat /tmp/e216v2)"; fi
}
mc() { # mc <label> <pattern> <行...> —— 变异正控：把违规形喂给尺本身（⛔ 写任何仓内产码）
  echo "    $1"; shift
  local pat="$1"; shift
  local line
  for line in "$@"; do
    printf '%s\n' "$line" | grep -oE "$pat" >/dev/null 2>&1
    echo "      fired=$?  <<$line>>"
  done
}
echo "=== V0 锚（⚠ 本腿起手记的是 feab45e3；这行是脚本运行那一刻的活 HEAD，尺不依赖它）==="
git rev-parse --short HEAD
echo "T1=$T1 T2=$T2"
A5='strings\.(Contains|Index|HasPrefix|HasSuffix)\([^,]*(Error\(\)|[a-zA-Z_]*[Ee]rr[a-zA-Z_]*\b)|switch .*\.Error\(\)|Error\(\) *=='
A5W='Error\(\)[^;]*(Contains|HasPrefix|HasSuffix|EqualFold|MatchString|Index)|(Contains|HasPrefix|HasSuffix|EqualFold|MatchString|Index|switch|==|!=)[^;]*Error\(\)'
echo "=== N 分母（补 a2 自报不足①：并集、逐行去重）==="
for t in $T1 $T2; do
  for p in 'Error\(\)' 'errors\.Is\(|errors\.As\(' 'Error\(\)|errors\.Is\(|errors\.As\('; do n "$t" "$p" internal cmd tools "$PROD"; done
  n "$t" 'Error\(\)' internal/panel cmd/wisp "$PROD"
done
echo "=== F 禁止那半ⓐ：按错误串子串/正则分支 ==="
echo "### F1 ＝a2 的 A5 逐字重放（两枚锚都该 0）"
for t in $T1 $T2; do g "$t" -nE "$A5" internal cmd tools "$PROD"; done
echo "### F2 变异正控：6 枚候选违规形喂 A5（fired=0＝尺抓到了；fired=1＝尺看不见那一形）"
mc "[A5]" "$A5" 'if strings.Contains(err.Error(), "413") {' 'if strings.Contains(strings.ToLower(oe.Error()), "too long") {' 'switch err.Error() {' 'if msg := err.Error(); strings.Contains(msg, "413") {' 'if regexp.MustCompile("4[0-9][0-9]").MatchString(err.Error()) {' 'if strings.Contains(e.Message, "413") {'
echo "### F2b 同 6 枚喂本腿的宽形尺 A5W"
mc "[A5W]" "$A5W" 'if strings.Contains(err.Error(), "413") {' 'if strings.Contains(strings.ToLower(oe.Error()), "too long") {' 'switch err.Error() {' 'if msg := err.Error(); strings.Contains(msg, "413") {' 'if regexp.MustCompile("4[0-9][0-9]").MatchString(err.Error()) {' 'if strings.Contains(e.Message, "413") {'
echo "### F3 宽形尺 A5W 的生产命中（本腿自加）"
g "$T2" -nE "$A5W" internal cmd tools "$PROD"
echo "### F4 正控·同族子串尺在判级包出数（a2 的 A6）"
n "$T2" 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)' internal/risk "$PROD"
echo "### F5 别名链绑定数（a2 的 A7 字面）"
g "$T2" -nE '(:?|\+)= *[a-zA-Z_.]*[Ee]rr(or)?\.Error\(\)' internal cmd tools "$PROD"
echo "### F6 ＝a2 的 A8 字面（射程＝它枚举的 6 枚文件）"
g "$T2" -nE 'strings\.(Contains|HasPrefix|HasSuffix|Index)\((line|detail|lastErr|reason|msg|text)' internal/agent/loop.go internal/audio/audio.go internal/llm/probe.go internal/panel/config_handlers.go cmd/wisp/panel_config_store.go cmd/wisp/resident_audio_windows.go
echo "### F7 本腿补的别名链尺（射程＝全生产树，不枚举文件）"
g "$T2" -nE 'strings\.(Contains|HasPrefix|HasSuffix|EqualFold|Index)\((line|detail|lastErr|reason|msg|text|message|errMsg|errText)\b' internal cmd tools "$PROD"
echo "### F9 测试面拿带数字错误串做断言的枚数（a2 报 16）"
n "$T2" 'Contains\([^,]*(Error\(\)|[a-zA-Z_]*[Ee]rr[a-zA-Z_]*\b), *"[^"]*[0-9][^"]*"' internal cmd tools
echo "### F10 有没有任何尺/钉已把'⛔ 按错误串分支'写成 CI 禁令或脚本断言（该 0）"
g "$T2" -nE 'error-text|errortext|错误串|错误文案' tools/d22scan scripts .github "$PROD"
echo "### F11 别名字段的子串判定（a2 的尺只认裸变量名，不认 res.Detail 这类字段）"
g "$T2" -nE 'strings\.(Contains|HasPrefix|HasSuffix|EqualFold|Index|ContainsAny)[(][a-zA-Z_]+\.(Detail|Message|CredentialError|Reason|Text|Error)[,)]' internal cmd tools "$PROD"
echo "### F12 变异正控：三枚别名/字段违规形逐尺喂（哪把尺看得见哪形，逐行标尺名）"
ALIAS='strings\.(Contains|HasPrefix|HasSuffix|EqualFold|Index|ContainsAny)[(][a-zA-Z_]+\.(Detail|Message|CredentialError|Reason|Text|Error)[,)]'
for line in 'if strings.Contains(res.Detail, "413") {' 'if strings.Contains(row.CredentialError, "auth") {' 'if strings.Contains(strings.ToLower(msg), "413") {'; do
  printf '%s\n' "$line" | grep -oE "$A5" >/dev/null 2>&1; echo "      [A5]  fired=$? <<$line>>"
  printf '%s\n' "$line" | grep -oE "$A5W" >/dev/null 2>&1; echo "      [A5W] fired=$? <<$line>>"
  printf '%s\n' "$line" | grep -oE "$ALIAS" >/dev/null 2>&1; echo "      [F11] fired=$? <<$line>>"
done
echo "### F8 三个数字 token 在生产里的命中（a2 甲②）"
for t in '413' 'too long' 'exit status'; do g "$T2" -n "$t" internal cmd tools "$PROD"; done
echo "=== M 必做那半ⓑ：按码选一句用户可见文案 ==="
echo "### M1 MessageKey 生产消费者（两枚锚都该只剩声明＋注释）"
for t in $T1 $T2; do g "$t" -n 'MessageKey' internal cmd tools "$PROD"; done
echo "### M2 正控·MessageKey 含测试全量"
g "$T2" -n 'MessageKey' internal cmd tools
echo "### M3 正控·MessageKey 扩到 frontend docs（看有没有别处消费者）"
g "$T2" -n 'MessageKey' internal cmd tools frontend docs
echo "### M4 生产里按 Class 分支的行（分的是句、状态还是退出码，逐枚眼看）"
g "$T2" -nE 'switch [^;]*\.Class|\.Class *==|observe\.ErrorClass\(' internal cmd tools "$PROD"
echo "### M5 声明行逐字现场（feab45e3:internal/observe/errors.go 158-165）"
git show "$T2:internal/observe/errors.go" | sed -n '158,165p' | cat -n
echo "=== C 票面反面教材 ChatInput.tsx ==="
echo "### C1 提交树两枚锚的名单里找 chatinput（都该 rc=1）"
for t in $T1 $T2; do echo "    tree=$t"; git ls-tree -r --name-only "$t" 2>&1 | grep -i chatinput; echo "      rc=$?"; done
echo "### C2 正控·同尺 *.tsx 枚数：提交树两枚锚 vs 工作树（a2 报的 90 是工作树口径）"
for t in $T1 $T2; do echo "      commit-tree $t count=$(git ls-tree -r --name-only "$t" | grep -c '\.tsx$')"; done
echo "      working-tree find count=$(find . -name '*.tsx' -not -path './.git/*' | wc -l)"
echo "### C3 正控·同尺找 chat-composer"
git ls-tree -r --name-only "$T2" 2>&1 | grep -i chat-composer; echo "      rc=$?"
echo "### C5 openchamber 在本仓文档面的出现处（射程＝那枚调研件）"
g "$T2" -n 'openchamber' docs/reports/survey-2026-09-28-composer-plus-menu.md
echo "### C7 票面那形的在仓逐字抄件（survey :376-382）"
git show "$T2:docs/reports/survey-2026-09-28-composer-plus-menu.md" | sed -n '376,382p'
echo "### C8 本仓 212-a3 假引用尺对这枚件名的归类"
g "$T2" -n 'ChatInput' .scratch/wisp/probes/212/a3/work/roster2.tsv
echo "### C9 d22scan phantom-citation 禁令的射程"
g "$T2" -n 'phantom-citation' tools/d22scan/main.go
echo "### C6 前端侧 413 与按文案选句族（射程＝提交树）＋正控 includes 枚数"
g "$T2" -n '413' frontend
g "$T2" -nE '\.(includes|indexOf|startsWith|match|search|test)\(.{0,3}(rror|too long|exit status|413|429)' frontend/src frontend/scripts
n "$T2" '[.]includes[(]' frontend/src
echo "=== G 推迟登记（裁ⓑ形状要用：这一层是不是已登记推迟）==="
echo "### G1 产码里的 DEFERRED 标记逐枚（⛔ 管道接 head，rc 才是 git grep 的）"
g "$T2" -n 'DEFERRED(' internal cmd tools "$PROD"
echo "### G2 自陈推迟那三行有没有 DEFERRED 码（errors.go 155-165 全窗）"
g "$T2" -n 'DEFERRED|deferred' internal/observe/errors.go
echo "### G3 SPEC-12 §5 登记表里那一行逐字"
g "$T2" -n '完整错误文案体系' docs/specs/SPEC-12-roadmap-governance.md
echo "### G4 PLAN.md 里 D23/S6 与错误文案的授权句（枚数＋前几枚）"
n "$T2" '错误文案' docs/PLAN.md docs/specs
g "$T2" -n '错误文案' docs/PLAN.md
echo "### G5 工单池里有没有别的票已经认领这层（按码选句/文案）"
g "$T2" -nE '错误文案|按码选|MessageKey' .scratch/wisp/issues
echo "=== D AC#3 没被打的那半：名字→判定的可达性 ==="
echo "### D1 判级包生产字符串匹配总数（a2 的 D1 字面）"
n "$T2" 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)' internal/risk internal/perm "$PROD"
echo "### D2 ALLOW 侧两枚件里的字符串比较逐枚"
g "$T2" -nE 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)' internal/risk/rules_shell.go internal/risk/rules_network.go "$PROD"
echo "### D3 变异正控（算子尺）：今天那行 vs 改成 Contains 那行"
mc "[今日 EqualFold 行]" 'strings\.(Contains|HasPrefix|HasSuffix)' 'return strings.EqualFold(argv0, entry) || strings.EqualFold(base, entry)'
mc "[变异 Contains 行]" 'strings\.(Contains|HasPrefix|HasSuffix)' 'return strings.Contains(argv0, entry) || strings.Contains(base, entry)'
echo "### D14 ALLOW 侧三处函数体行窗里的子串算子枚数（今天该 0）"
for spec in 'internal/risk/rules_shell.go:91,104' 'internal/risk/rules_network.go:114,128' 'internal/risk/blacklist.go:115,128'; do
  f=${spec%%:*}; w=${spec##*:}
  git show "$T2:$f" | sed -n "${w}p" > /tmp/win216.txt
  echo "      win=$spec substring_ops=$(grep -cE 'strings\.(Contains|ContainsAny|ContainsRune|Index|IndexAny)[(]' /tmp/win216.txt)"
done
echo "      变异窗（喂那枚算子尺）ops=$(printf '%s\n' 'if strings.Contains(argv0, entry) {' | grep -cE 'strings\.(Contains|ContainsAny|ContainsRune|Index|IndexAny)[(]')"
echo "### D15 名字真进判定的那枚 Contains 在哪一侧（逐字）"
git show "$T2:internal/risk/blacklist.go" | sed -n '252p'
echo "### D10 现有代码/测试里有没有钉已盯住 ALLOW 侧那三枚函数"
g "$T2" -nE 'shellAllowlisted|hostAllowlisted|overrideApplies' internal/risk
echo "### D11 ALLOW 侧三处比较算子的逐字现场（行窗）"
git show "$T2:internal/risk/rules_shell.go" | sed -n '91,104p'
git show "$T2:internal/risk/rules_network.go" | sed -n '114,128p'
git show "$T2:internal/risk/blacklist.go" | sed -n '115,128p'
echo "### D12 基名图案六枚的逐字现场（deny 侧）"
git show "$T2:internal/risk/blacklist.go" | sed -n '236,254p'
echo "### D13 显示名入口与洗函数（含本腿自抓的假 0 一例）"
g "$T2" -nE 'ContainsAny\(name' internal/panel/attachments.go
g "$T2" -nE 'scrubForDisplay|scrubComposerForDisplay' internal/panel
echo "      ↓ 同一枚尺去掉 -E ＝假 0（BRE 里竖线是字面量）"
g "$T2" -n 'scrubForDisplay|scrubComposerForDisplay' internal/panel
echo "### D4 '允许'类字面当判据（该 0）＋正控（名字子串判定确有）"
g "$T2" -nE 'strings\.[A-Za-z]+\([^,]*(name|Name|base|baseName)[^,]*, *"(允许|已允许|approved|allowed)"' internal cmd tools "$PROD"
g "$T2" -nE 'strings\.[A-Za-z]+\(base' internal/risk/blacklist.go
echo "### D5 附件显示名到判级包可达性（该 0）＋正控 D6＋非空目录 D7"
g "$T2" -ci 'attachment' internal/risk internal/perm internal/tools internal/agent "$PROD"
g "$T2" -ci 'params' internal/risk "$PROD"
git ls-tree --name-only "$T2:internal/perm"; echo "      rc=$?"
echo "### D8 a2 引的 attachments.go:69/:260 在两枚锚上是否仍对得上"
for t in $T1 $T2; do g "$t" -n 'never used to build a path' internal/panel/attachments.go; g "$t" -nE 'name := "attachment-' internal/panel/attachments.go; done
