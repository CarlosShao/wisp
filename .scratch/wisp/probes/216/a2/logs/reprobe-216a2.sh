#!/usr/bin/env bash
# 216-a2 只读普查的尺名册（⛔ 产码判定、⛔ go test/build/vet）。
# 一律读具名提交树 44f51fc8（⛔ 混在飞工作树）；每把尺都带正控。
set -u
ANCHOR=44f51fc8
cd "$(git rev-parse --show-toplevel)"
r() { echo "### $1"; shift; git grep "$@" 2>&1 | sed "s#^44f51fc8:##"; echo "rc=${PIPESTATUS[0]}"; }

echo "=== 甲① 分母 ==="
printf '### A1 Error() 站点·含测试 计数＝'; git grep -n 'Error()' "$ANCHOR" -- internal/panel cmd/wisp | wc -l
r "A2 Error() 站点·生产(去 _test)" -n 'Error()' "$ANCHOR" -- internal/panel cmd/wisp ':(exclude)*_test.go'
printf '### A3 errors.Is/As 含测试 计数＝'; git grep -nE 'errors\.Is\(|errors\.As\(' "$ANCHOR" -- internal cmd tools | wc -l
printf '### A4 errors.Is/As 生产 计数＝'; git grep -nE 'errors\.Is\(|errors\.As\(' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go' | wc -l
r "A5 生产里按错误串子串分支(类b)" -n -E 'strings\.(Contains|Index|HasPrefix|HasSuffix)\([^,]*(Error\(\)|[a-zA-Z_]*[Ee]rr[a-zA-Z_]*\b)|switch .*\.Error\(\)|Error\(\) *==' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
printf '### A6 正控·同族尺(strings.Contains)在生产 internal/risk 计数＝'; git grep -nE 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)' "$ANCHOR" -- internal/risk ':(exclude)*_test.go' | wc -l
r "A7 别名链：err.Error() 绑进变量(生产)" -n -E '(:?|\+)= *[a-zA-Z_.]*[Ee]rr(or)?\.Error\(\)' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
r "A8 那些绑定变量被子串判定过吗(应 0)" -n -E 'strings\.(Contains|HasPrefix|HasSuffix|Index)\((line|detail|lastErr|reason|msg|text)' "$ANCHOR" -- internal/agent/loop.go internal/audio/audio.go internal/llm/probe.go internal/panel/config_handlers.go cmd/wisp/panel_config_store.go cmd/wisp/resident_audio_windows.go
r "A9 选句按码：switch on class / switch err" -n -E 'switch *(observe\.ErrorClass\(|err$|err *\{)' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
r "A10 MessageKey 的生产消费者(应 0)" -n 'MessageKey' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
r "A11 正控·MessageKey 消费者在测试里存在(全量)" -n 'MessageKey' "$ANCHOR" -- internal cmd tools

echo "=== 甲② 数字串被下游子串判定 ==="
for t in 413 'too long' 'exit status'; do
  r "B1[$t] 生产命中" -n "$t" "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
done
r "B2 正控·observe.ClassProvider 生产命中" -c 'observe.ClassProvider' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
r "B3 类型化退出码处理(类a 证据)" -n -E 'exec\.ExitError|ExitCode\(\)' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'

echo "=== 甲③ 反面教材件在本仓存在吗 ==="
echo "### C1 工作树 find"; find . -name 'ChatInput.tsx' -not -path './.git/*' 2>/dev/null; echo "rc=$?"
echo "### C2 正控·同尺找 *.tsx(应非 0)"; find . -name '*.tsx' -not -path './.git/*' 2>/dev/null | wc -l
echo "### C3 提交树 44f51fc8 名单"; git ls-tree -r --name-only "$ANCHOR" | grep -i chatinput; echo "rc=$?"
echo "### C4 正控·同尺已知存在件"; git ls-tree -r --name-only "$ANCHOR" | grep -i chat-composer

echo "=== 乙 名字里的字当判级输入 ==="
r "D1 判级包内所有字符串匹配(分母)" -n -E 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold|Index)' "$ANCHOR" -- internal/risk internal/perm ':(exclude)*_test.go'
r "D2 判级包内 allow 语义行的子串匹配(逐枚看方向)" -n -E 'strings\.(Contains|HasPrefix|HasSuffix)' "$ANCHOR" -- internal/risk/rules_shell.go internal/risk/rules_network.go ':(exclude)*_test.go'
r "D3 生产里对名字做'允许'类字面判定(应 0)" -n -E 'strings\.[A-Za-z]+\([^,]*(name|Name|base|baseName)[^,]*, *"(允许|已允许|approved|allowed)"' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
r "D4 正控·同尺去掉字面(名字子串判定确有)" -n -E 'strings\.[A-Za-z]+\(base' "$ANCHOR" -- internal/risk/blacklist.go
r "D5 Attachment 一词进判级包了吗(应 0)" -ci 'attachment' "$ANCHOR" -- internal/risk internal/perm internal/tools internal/agent ':(exclude)*_test.go'
r "D6 正控·同尺 params 确有" -ci 'params' "$ANCHOR" -- internal/risk ':(exclude)*_test.go'
r "D7 显示名→路径 的唯一出口(生成名=内容哈希)" -n -E 'name := "attachment-' "$ANCHOR" -- internal/panel/attachments.go
r "D8 门控调用点(判定入口的生产 caller)" -n -E 'risk\.(Gate|Classify)\(|classifier\.Classify\(|assessor\.Assess\(|\.Inspect\(|CheckText\(' "$ANCHOR" -- internal cmd tools ':(exclude)*_test.go'
echo "### D9 internal/perm 非空目录(D5 的 0 不是空转)"
git ls-tree --name-only "$ANCHOR":internal/perm

echo "=== 甲②′ 界面侧（只读，射程＝提交树 frontend/**；我零写面）==="
echo "### F1 前端里 413 这个数字串"; git grep -n '413' "$ANCHOR" -- frontend; echo "rc=$?"
F2P='\.(includes|indexOf|startsWith|match|search|test)\(.{0,3}(rror|too long|exit status|413|429|4[0-9][0-9])'
echo "### F2 前端按错误文案子串判定族(尺见上一行 F2P)"; git grep -nE "$F2P" "$ANCHOR" -- frontend/src frontend/scripts; echo "rc=$?"
echo "### F3 正控·同族尺 includes 调用在 frontend/src 的枚数"; git grep -nE '[.]includes[(]' "$ANCHOR" -- frontend/src | wc -l
