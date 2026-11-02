# 385-v1 裁决表：mockllm 400 字节截断修复·非实现者验收

**票题**: 385 — mockllm echo 400-byte clamp cut multi-byte UTF-8 mid-sequence  
**落地提交**: `fad2f5e4f9945c2b4919258d47b938d06f05b9eb`「票 385 落地腿 385-r1」(HEAD ancestor: `2806bb9f0`)  
**裁决者**: 非实现者（本腿 v1，仅读复核与对抗验证）  
**日期**: 2026-11-02  

---

## AC#1 修法形状｜字符边界安全（PASS）

| 检查点 | 读数/事实 | 结论 |
|--------|----------|------|
| echo 路径原 `len(answer) > 400 => answer[:400]` 是否移除 | `chat.go:159-165`, `messages.go:107-109`, `responses.go:89-94` —— 三方言均删除字节帽；注释说明“防超长意图的护栏为 8MB LimitReader 与 max_tokens 请求语义” | ✅ PASS |
| 现存的长度控制位点（max_tokens）是否逐码点安全 | `chat.go:217` → `answer[:runeBoundaryWithin(req.MaxTokens*4, answer)]`；函数定义见 `chat.go:365-373`，从 n 起**回退**到最近的 `utf8.RuneStart` 边界再截断（ASCII 时恰为 min(n,len(s))） | ✅ PASS |
| streaming chunking 是否逐码点安全 | `splitChunks()` (chat.go:376-384) 以 `[]rune(s)` 分块，永不按字节切 | ✅ PASS |

**读数**: 3 处旧 clamp 全部摘除；剩余 truncation 通过 `runeBoundaryWithin` + rune 分块实现，多字节安全。

---

## AC#2 多字节回归一发｜402 汉字案通不过即红（PASS）

| 测试项 | 读数 | 结论 |
|--------|------|------|
| 包内运行 | 正确尺＝模块内跑（tools/mockllm 有独立 go.mod，从根跑 `./tools/mockllm` 必 setup failed 假红，定式见编排者 A1304补）：`cd tools/mockllm; go test . -count=1 -timeout 10m` | **ok 0.154s**, rc=0 |
| 顶层计数 | PASS=20, FAIL=0, SKIP=0, RUN=36 (含子测试) | ✅ PASS |
| 定向案 | `Test385EchoDeliversLongMultibyteInFull` (ticket30Sample: 402 hanzi + summary = 1243 bytes / 419 runes) | **PASS**, assembled text length matches exactly in all three dialects (stream SSE, non-stream JSON, messages/responses) |
| 输出一致 | 剥离 "echo: " 签名后，与 C5 seam 输出 byte-identical（测试内断言 `text != want` / `content != want`） | ✅ PASS |

**读数**: 流式 SSE、非流 JSON、messages/responses 三模式各返回完整 1243 字节；尾部无 U+FFFD，419 runes 完整。

---

## AC#3 反向判据｜还原核（PASS）

**方法**: 出仓 overlay 恢复三文件 400-byte hard clamp → 复跑定向案 → 预期 RED；撤 overlay 重跑 → 预期 GREEN；sha256 核对零漂移。

| 步骤 | 操作 | 结果 | 读数/失败消息 | 结论 |
|------|------|------|--------------|------|
| 1 | 备份当前三文件 (git sha256 baseline) | chat.go=5C96…09F; messages.go=874B…16E; responses.go=017F…CB61 | — | ✅ PASS |
| 2 | overlay.json 替换为重新插入 `if len(answer) > 400 { answer = answer[:400] }` 的版本 | `D:\tmp\wisp\20261102\385v1\overlay.json` 映射三文件 | — | ✅ PASS |
| 3 | `go test . -run Test385EchoDeliversLongMultibyteInFull -overlay overlay.json` | **FAIL** rc=1 | 三方言子案全红：“assembled text = **402 bytes / 138 runes**, want 1243 bytes / 419 runes”，尾部残迹含 U+FFFD（与落地腿记录 402/138 逐字吻合） | ✅ PASS (red as expected) |
| 4 | 撤 overlay 重跑同一定向案 | **PASS** rc=0 | ok | ✅ PASS |
| 5 | 撤 overlay 后核对三文件 sha256 | ZERO-DRIFT | chat/messages/responses 三 hash 与步骤 1 基线逐一相同；`git status -- tools/mockllm` 空 | ✅ PASS |

**读数**: overlay 时字节帽在第 400 字节硬切 ⇒ 402 bytes / 138 runes（应为 1243/419），SSE 装配尾部出现残 rune；还原后立即绿且文件 hash 零漂移、仓面 git 干净。

---

## AC#4 零回归｜现有消费者全绿（PASS，本腿亲测除 cmd/wisp）

| 消费方 | 本腿复跑 | 读数 | 状态 |
|--------|----------|------|------|
| tools/mockllm 包内 | ✅ | PASS=20, FAIL=0（AC#2 已记） | ✅ PASS |
| internal/llm/... + internal/present（票 30 E2E 真子进程） | ✅ 本腿现跑 `go test ./internal/llm/... ./internal/present -count=1 -timeout 10m` | 7 包全 ok、rc=0（llm 60.4s / adaptertest 4.4s / anthropic 21.9s / golden 0.1s / openaichat 32.3s / openairesponses 29.3s / present 4.6s） | ✅ PASS |
| cmd/wisp 定向消费（票 12 E2E） | ⛔ 本腿边界禁触 cmd/wisp | 沿用落地腿 385-r1 记录（46/36/0/0 rc=0，编排者 A1304/A1304补 已独立复跑认可） | 📎 在案 |

**读数**: 包内＋票 30 E2E 消费者两路本腿亲测全绿；cmd/wisp 支属本腿禁区、取落地记录（且该记录已被编排者复跑过）。

---

## AC#5 纪律格｜范围门禁（PASS）

| 项目 | 检查 | 读数 | 结论 |
|------|------|------|------|
| HEAD ≥ 2806bb9f0? | `git log --oneline -5` | HEAD=`2806bb9f0` | ✅ PASS |
| 索引干净？ | `git status --porcelain=v1` | 0 tracked modifications (untracked logs only) | ✅ PASS |
| fad2f5e4f 写面与题面一致？ | `git show --name-only fad2f5e4f` | `docs/evidence/s1/385r1-impl.md`, `tools/mockllm/chat.go`, `chat_echo_385_test.go`, `messages.go`, `responses.go` | ✅ PASS (证据件允许范围) |
| tools/mockllm/** 外是否改码？ | 对比 stat | 否 | ✅ PASS |
| gofmt/gofumpt 清洁度（盘现状） | `gofmt -l .` → 0 行；绝对路径 `$env:GOPATH\bin\gofumpt.exe -l .` (tools/mockllm 内) → 0 行 exit=0 | ✅ PASS |
| 台账/票面/push 是否改写？ | — | 本腿不改动任何台账/commit-msg | ✅ PASS |

**读数**: 写面严格限制在 `tools/mockllm/**` + 其自身证据文件；无其他修改。

---

## 综合判定

| 格子 | 结论 |
|------|------|
| AC#1 修法形状 | PASS |
| AC#2 多字节回归 | PASS |
| AC#3 反向判据 | PASS |
| AC#4 零回归 | PASS（cmd/wisp 支取落地记录＋编排者复跑在案，本腿 ⛔ 禁区） |
| AC#5 纪律格 | PASS |

**是否同意结 385**: ✅ **是**

理由: 修复已落盘且多字节安全；定向案通过并完整还原性验证；范围纪律合规。

---

**附件**:  
- 原始测试日志: `D:\tmp\wisp\20261102\385v1\go-test-mockllm.txt`  
- 定向案 verbose: `D:\tmp\wisp\20261102\385v1\go-test-mockllm-v.txt`  
- Overlay red run: `D:\tmp\wisp\20261102\385v1\reverse-overlay.txt`  
- Restore green run: `D:\tmp\wisp\20261102\385v1\restore-rerun.txt`  
- 票 30 E2E 消费者 7 包：`D:\tmp\wisp\20261102\385v1\consumers-internal.txt`  

---

签署: 385-v1 非实现者验收 | 2026-11-02
