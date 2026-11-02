# 票 337 注入式宿主半格验证腿 337-v1c · 裁决表

- 腿：`337-v1c`＝**cmd/wisp 半格真验腿**（两枚 TestTicket337R2*经验红；反向判据 -overlay 摘线→红；还原核）。
- 派单存档：`.scratch/wisp/probes/orch/20261102/txt/337v1c-dispatch.txt`＝**未落盘**（起手 `glob` 该目录 14 枚 txt，无此名；最接近的两枚＝同目录 `337v1-dispatch.txt`／`337v1b-dispatch.txt`）。本腿据**腿单正文**（＝本派消息）＋票面执行，⛔ 代编排者补建存档件。
- 票面：`.scratch/wisp/issues/337-fail-closed-reason-has-no-reader-wiring-census.md`（未动一字）。
- 交付核验：`7fd2f37b5`（337-r2 ＝ cmd/wisp 侧注入式宿主装配点补接 `risk.Logf = rt.auditf`）在 HEAD 祖先。
- 边界：**只读复核**＋**证据件**；⛔ 改码／⛔ 改票面／⛔ 台账／⛔ push；临时件 `D:\tmp\wisp\20261102\337v1c\`（只建不删）。

---

## 0. 起手三核与环境 staging 准备

| 核 | 尺（逐字可重跑） | 读数 | 判 |
|---|---|---|---|
| ① HEAD ≥ `059011810` | `git rev-parse HEAD` / `git merge-base --is-ancestor 059011810 HEAD` | `HEAD=2ebee4fa3`（短 9 位）；`rc=0`（相等） | **PASS** |
| ② 索引干净 | `git status --porcelain --untracked-files=no` 行数 | **0**（跟踪面零脏；余下全是 .scratch 等别腿未跟踪影件） | **PASS** |
| ③ 交付在 HEAD 树里 | `git merge-base --is-ancestor 7fd2f37b5 HEAD` | `rc=0` | **PASS** |

**staging 环境**（按票面 sherpa 三层陷阱）：
- 测试二进制导出至 `D:\tmp\wisp\20261102\337v1c\bin\wisp337v1c.test.exe`（`go test -c ./cmd/wisp`）。
- 三 DLL 拷入同目录：`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`（`third_party/sherpa-onnx/`）。
- 环境变量指向主仓模型资产：`WISP_SILERO_MODEL=third_party/spike-models/silero_vad.onnx`；`WISP_SPIKE_MODELS=third_party/spike-models`。
- PATH 前缀设为 staging bin：`$env:PATH = "$tmp\bin;" + $env:PATH`，确保 System32 旧 DLL 被遮蔽（红控），不给＝`0xc0000135`。
- `-timeout 30m` 整包必加（票面边界要求）。
- gofumpt 路径为绝对 `/d/work/base/gopath/bin/gofumpt.exe`（本腿未触发 gofumpt，仅核迹）。

---

## 1. Grid ① — 真跑注入支两枚测试（逐名报 PASS/FAIL＋字数）

**名称校注**：派单 `-run 'TestRiskSinkInjected337r2'` 的 literal 匹配**磁盘 0 命中**（名字漂移），实际测试名为 `^TestTicket337R2InjectedGateArmsRiskSink(OnRuntimeAudit|WhenHostNeverArmed)$`。本腿按实际名复跑。

| # | 测试名 | 命令片段（关键参数已齐） | rc | 判定行数｜断言行号 | PASS/FAIL | 一句话读数 |
|---|--------|--------------------------|-----|------------------|-----------|------------|
| 1 | `TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit` | `./wisp337v1c.test.exe -test.run='^TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit$' -test.v -test.count=1 -test.timeout=30m` | 0 | --- PASS ... `(2.72s)`；审计句落在 **:208** | **PASS** | audit sink 收到 fail-closed 实话，foreign sink 收行为 0（装配根重授同一个 rt.auditf）。 |
| 2 | `TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed` | `... -test.run='^TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed$' -test.v -test.count=1 -test.timeout=30m` | 0 | --- PASS ... `(2.08s)` | **PASS** | 宿主从未武装 risk.Logf 的分支，组装后 risk.Logf 仍被装配根接管且落 sink。 |

**log 归档**：  
- `D:\tmp\wisp\20261102\337v1c\logs\run1-TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit.log`（7993 字节）。  
- `D:\tmp\wisp\20261102\337v1c\logs\run2-TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed.log`（4719 字节）。

**结论**：两枚注入支测试**同时绿**；rc=0。⛔ "cmd/wisp 整包 rc=1" 之既有真窗族＋进程级截断为其他 leg 的归因，本腿对这两枚专名测试的计数＝**全绿**（未跑整包）。

---

## 2. Grid ② — （可选反向判据）仓外 -overlay 摘 run.go:758 ⇒ 该测试红；还原核

**突变**：仓外 `D:\tmp\wisp\20261102\337v1c\overlay\run.go` 从 HEAD blob 剥离 `risk.Logf = rt.auditf` 一行（原 1926 行 → 1925 行），并生成 `-overlay.json` 指向替换。

| 格 | 尺 | 读数（红点/断言行号） | rc | 判 |
|---|-----|----------------------|-----|-----|
| A | overlay build | `go test -c -overlay <json> -o wisp337v1c-mut.test.exe ./cmd/wisp` | 0 | 编译成功（无语法错） |
| B | 突变‑run1 | `... wisp337v1c-mut.test.exe -test.run=...OnRuntimeAudit ...` | 1 | `risk_sink_injected_337r2_windows_test.go:184: an A deny on the injected-host assembly did not reach the runtime audit sink; got:` |
| C | 突变‑run2 | `... wisp337v1c-mut.test.exe -test.run=...WhenHostNeverArmed ...` | 1 | `risk_sink_injected_337r2_windows_test.go:244: the injected assembly did not arm risk.Logf: the fail-closed reason is still dropped` |
| D | 还原核 | `git status --porcelain --untracked-files=no`（ tracked dirty）＋`git hash-object` HEAD vs WT 逐枚 | **0** 行；blob 均**IDENTICAL** | **PASS** |
| E | 还原重跑 | `wisp337v1c.test.exe`（绿色）再跑两枚 | 0 | 同 §1 之 PASS 读数；**restore**成立 |

**红点详述**（**证明尺有牙**）：
- mut‑run1 **`:184` Fatalf**：测试期望 DENY 落 stderr 审计尾，但删除 `rt.auditf` 赋值的注入支不重写风险日志 sink，于是 A 级拒绝仍发向 host 先前的 foreign sink，导致 `errs.Tail(mark)` 尾部不命中 `[audit] risk: A-list DENY`，触发 `t.Fatalf`。
- mut‑run2 **`:244` Fatalf**：宿主从未武装过 risk.Logf，删除本行的注入支导致运行时 risk.Logf 仍为 nil，`assembleRuntime` 返回即被起手的 `if risk.Logf == nil { t.Fatal(...) }` 捕获。

**说明**：header 文档 booked case 1 红在 "not re-arm"（line ~188 Errorf），但因 `Fatalf(:184)` 提前终止执行，真实可见红点为 `:184`；case 2 确实红在 `:244` Fatalf。两者均为“拆线即死”的证据。

**日志归档**：
- `D:\tmp\wisp\20261102\337v1c\logs\mut-run1-OnRuntimeAudit.log`（4469 字节）。
- `D:\tmp\wisp\20261102\337v1c\logs\mut-run2-WhenHostNeverArmed.log`（4428 字节）。
- `D:\tmp\wisp\20261102\337v1c\overlay\run.go`（突变版本）。
- `D:\tmp\wisp\20261102\337v1c\overlay.json`（overlay 映射）。

**结论**：仓外 `-overlay` 摘下 `run.go:758` 赋值，测试**双红**；还原后工作树跟踪面零脏、blob 同串、重跑双绿。**反向判据通过**。

---

## 3. Grid ③ — 核证 run.go:758 注入支行在（读码逐字）

| 核 | 尺 | 读数 | 判 |
|---|-----|------|-----|
| HEAD 逐字 | `git show HEAD:cmd/wisp/run.go \| Select-Object -Skip 756 -First 3` | 第 758 行（1-based）逐字：<br>`		risk.Logf = rt.auditf` | **存在** |
| 工作树一致 | `git hash-object cmd/wisp/run.go` vs `git rev-parse HEAD:cmd/wisp/run.go` | `82e052e03606a6a0ea09dfe5e8a108633c2ed078` vs 同上 | **IDENTICAL** |
| 生产侧 census（剔 `_test.go` 及 `.scratch` 影） | `git grep -nE 'risk\.Logf[[:space:]]*=' HEAD -- '*.go' \| grep -v '_test\.go:' \| grep -v '\.scratch/'` | 恰 **2 枚**：<br>`internal/agent/approval/gate.go:181`<br>`cmd/wisp/run.go:758` | **符合 337-v1 §1 census** |
| 总 census（含 .scratch 探针影） | `git grep -nE 'risk\.Logf[[:space:]]*=' HEAD -- '*.go'`（剔 _test.go） | 3 枚（+ `.scratch/wisp/probes/175/r1/main.go:82`） | 影件按普查纪律剔除 |

**read code verbatim**（§3 交付件之关键行，逐字）：
```go
                // local audit log only - never the model text, never the panel bytes - and it
                // cannot allow or deny anything. The folded (normPath) paths it carries are
                // exactly what an auditor needs to attribute a deny, which is the surface a
                // local audit line is for (337-r1's own note).
                risk.Logf = rt.auditf
        } else {
                rt.liveCards = newNativeCards()
```

**结论**：注入支行 `risk.Logf = rt.auditf` 在 `cmd/wisp/run.go:758` **实存且与 HEAD 同串**；生产侧风险日志汇入口 **恰 2 枚**（控制台门 `gate.go:181` ＋ 注入式宿主门 `run.go:758`）。Grid ③ **通过**。

---

## 4. 逐格总表

| 格 | 内容概要 | 结果 |
|---|-----------------------------------------------|------|
| ① 真跑注入支两测 | TestOnRuntimeAudit（pass）、TestWhenHostNeverArmed（pass）；rc=0；审计句落 sink | **PASS** |
| ② -overlay 摘线反证 | mut 红于 :184/:244；还原核：跟踪面 0 行，blob 同串，重跑绿 | **PASS** |
| ③ 读码逐字确认 | run.go:758 逐字 `risk.Logf = rt.auditf`；HEAD/WT 同串；非测产 censure＝2 | **PASS** |
| 纪律检查 | 临时件只建不删；commit 清单＝本证据件；⛔ add -A／amend/reset/rebase/stash/checkout；⛔ push | **PASS** |

---

## 5. 是否同意翻框（337 可否结清）

**同意结清票 337 cmd/wisp 半格**。

理由：交付的注入式宿主安全汇接点已在 `cmd/wisp/run.go:758` 实装并通过经验红线检验。两枚注入支测试在 staging（DLL 贴身＋模型 env 指向）中绿过；移除该赋值则测试双红（明确断言行号点名：`:184`、`:244`），还原则恢复通行；生产侧 `risk.Logf` 汇接口的数量与位置与 337-r1 预期一致（gate.go:181 + run.go:758 = 2）。反向判据证实接线非死码，正控证实句子进入本地审计日志。

**保留提示**：本腿对 `TestTicket337R2*` 专名的复跑通过，派单中的 dispatch regex (`TestRiskSinkInjected337r2`) 与磁盘测试名存在漂移，属文档口径问题，不影响数值判据。若需要严格对齐派单命名，应在票面或注释中校正正则。

*验证腿：`337-v1c`（非实现者，cmd/wisp 半格唯一验腿）。本腿产码字节＝0；只 commit 本件，⛔ push。*

— 

*证据附件：临时件目录 `D:\tmp\wisp\20261102\337v1c\` 保存所有二进制、日志、overlay 与映射文件；本裁决表提交到仓库 `docs/evidence/s1/337-v1c-verdict.md`。*
