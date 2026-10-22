# 164-r1（票 164 · AC#2b）实施表 —— 把"未修码读数"改署名成"未知工具拒绝形状＋D37 分类的守卫"

- 格：**只做 AC#2b**（AC#4／AC#5 未碰；框归编排者）。锚点 `ce84d53d`；写面＝`internal/tools/task_output_ac2_before_test.go`＋本件。
- 起手 `git status --porcelain -- internal/tools docs/evidence/s1`＝**空**（别人的改动全在 `.scratch/**`，未提交未动）。

## 0 结论（先给判语）

★**那枚用例是"守卫"，不是"恒真装饰品"**——两枚互相独立的变异各杀它一次（下面 §1 逐字）：
- **M1 改坏桥的拒绝文案** ⇒ 红在**形状**断言（`Text` 不再含"未知工具"）；
- **M2 把该支的 D37 类改成 `internal`** ⇒ 红在**分类**断言；
- **两发都还原**（`bridge.go` sha256 三处逐字相同）⇒ **绿**。
⇒ 它今天钉的是**活的行为**（不是树状态），所以"可以留"这个前提**成立**，且它的署名必须是**守卫**。
★**它仍然不是、也不许再被引用为"AC#2 的未修码读数"**：桥下面挂的注册表是 `BuiltinFSEntries`（fs 家族），
生产在它之上另挂 task 家族（`cmd/wisp/run.go`）⇒ 这枚用例对"后台任务今天能不能读"这一问**不可证伪**。

## 1 两向量到（同一枚**交付后**文件；每发先落地、后读数、读完立刻还原）

落地证据（`grep -n`）＋`go build ./internal/tools/` rc=0（**两发都是 0**）：

| 发 | 改动 | 落地行 | 读数 |
|---|---|---|---|
| 原样 | — | — | `--- PASS` |
| M1 | `bridge.go` 拒绝文案 `"未知工具 "`→`"无法识别工具 "` | `276: "无法识别工具 "+req.Name+"，可用工具见 list_tools",` | **FAIL** |
| M2 | 同支 `observe.ClassTool`→`observe.ClassInternal` | `277: string(observe.ClassInternal), OutcomeKindNotFound)` | **FAIL** |

```
=== [原样绿]
=== RUN   TestTaskOutputAC2BeforeLegIsUnreachable
    task_output_ac2_before_test.go:66: unknown-tool refusal: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.01s)

=== [M1 红]
    task_output_ac2_before_test.go:66: unknown-tool refusal: IsError=true ErrorClass="tool" Truncated=false Text="无法识别工具 task.output，可用工具见 list_tools"
    task_output_ac2_before_test.go:73: expected the unknown-tool refusal (Bridge.Execute, Lookup-failure branch), got: "无法识别工具 task.output，可用工具见 list_tools"
--- FAIL: TestTaskOutputAC2BeforeLegIsUnreachable (0.02s)

=== [M2 红]
    task_output_ac2_before_test.go:66: unknown-tool refusal: IsError=true ErrorClass="internal" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
    task_output_ac2_before_test.go:76: error_class = "internal", want tool (an unknown tool is self-correctable, D37)
--- FAIL: TestTaskOutputAC2BeforeLegIsUnreachable (0.01s)
```

★**三处 sha256 逐字相同**（pristine／M1 还原／M2 还原）：
`a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd  internal/tools/bridge.go`
变异中两枚（只为证明"读数确实来自变异"）：M1 `89210b0035b6b67e68e66aead76a3ceffeecaea7a830e138a5748d7344c09c7a`／M2 `589b537e5bb4d9b1e773d959234ee832ed1911e84b0f4ed37b7714017727e51c`。
还原后 `git status --porcelain -- internal/tools/bridge.go`＝**空**、`git diff --stat` ＝**空**。

## 2 改后的署名（逐字＝文件第 11–49 行，⛔"只说改了"）

```go
// Guard, not a historical reading: the shape of the unknown-tool refusal and
// its D37 class.
//
// This case is NOT evidence for ticket 164 AC#2 ("start something in the
// background, then read what it printed - can it be done today"). That reading
// was taken on the real tree at anchor f206e9f and lives in the ticket's AC#2
// note; this file cannot speak for it, because the bridge below is composed
// from BuiltinFSEntries, which is the fs family only, while the production
// composition root appends the task family separately (cmd/wisp/run.go,
// BuiltinTaskEntries). What this file pins is the bridge's own contract at the
// choke point every tool call crosses (Bridge.Execute):
//
//  1. REFUSAL SHAPE - an unknown name falls into the Lookup-failure branch and
//     comes back as DATA, not as a Go error: err == nil, because D37 says a
//     rejection is not a fault (SPEC-07 2 states the C3 rule as a refusal, not
//     as an error). The text names the requested tool as unknown and points the
//     model at list_tools, and it does not echo anything it never looked up.
//  2. D37 CLASS - that refusal carries class "tool" (observe.ClassTool:
//     self-correctable, so the task keeps running), never "internal".
//  3. FS ENTRY POINT - BuiltinFSEntries does not smuggle a task.* name into the
//     directory: the fs family stays single-purpose. This clause is about the
//     fs registration path only. It is NOT a claim about the production
//     registry, which does carry task.output today.
//
// The teeth were measured in both directions on 2026-10-22 (tree ce84d53d):
// rewording the refusal opener in bridge.go makes this red on the shape
// assertion; switching that branch's class to internal makes it red on the
// class assertion; restoring bridge.go (same sha256) makes it green. It pins
// the live behaviour of the fs-only bridge, so it is green on a tree where the
// feature has already landed (the measurement above was taken on one).
//
// Reproduce:
//
//	go test -count=1 -v -run '^TestTaskOutputAC2BeforeLegIsUnreachable$' ./internal/tools/
//
// The file name and the test name are historical keys - ticket 164 AC#2b, the
// archived instrument .scratch/wisp/probes/164/v1/ac2-anchor-reading.sh and the
// CI logs address this case by them - so they are kept unchanged even though
// the case no longer signs itself as a before/after reading.
```

同批改的三处**说明文字**（不是断言）：`t.Logf` 标签 `AC#2 BEFORE verbatim:`→`unknown-tool refusal:`；两处 `Fatalf` 措辞去掉了过期行号 `bridge.go:247`；`a task.* tool is registered`→`BuiltinFSEntries carried a task.* name`。

## 3 断言逻辑一字未动的证明（可复跑）

- 去注释后逐行 `diff`（`diff <(git show HEAD:$F | grep -vE '^[[:space:]]*//') <(grep -vE '^[[:space:]]*//' $F)`）⇒ **只剩 4 行**，全是 `Logf`／`Fatalf` 的**措辞**。
- 计数与 HEAD 同：`t.Fatalf`4／`t.Errorf`2／`if `7／`strings.Contains`2／`ErrorClass`4。
- 文件形状：`CR=0`、`LF=95`、4207 字节（`tr -cd "\r" | wc -c` 量的，⛔ 不用 `grep` 数 CR）。

## 4 门禁（自己现跑，四数分开报，⛔ 相减）

| 门 | 命令 | 读数 |
|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | **rc=0** |
| 静态 | `go vet ./internal/tools/` | **rc=0** |
| 单包 | `go test -count=1 -v ./internal/tools/`（主检出） | **rc=0**；`RUN=329`／顶层 `PASS=246`／`FAIL=0`／`SKIP=0` |
| 族 | `go test -count=1 -run TaskOutput ./internal/tools/` | **rc=0**；RUN=4／PASS=4／FAIL=0／SKIP=0 |
| D22 | `sh scripts/d22scan.sh` | **rc=0**（正控 `runtests.sh`：`PASS=35 FAIL=0 SKIP=0`，`RUN=77`；扫描 clean） |
| 路径 | `sh scripts/check-path-length-budget.sh` | **rc=0**（`VERDICT GREEN`，over-budget=61 全在名册） |

**红名册逐名（取主检出，⛔ 取导出树）**：`./internal/tools/` 顶层 FAIL＝**0 枚**、子测试 FAIL＝**0 枚** ⇒ **逐名＝（无）**，新增红＝**0**。

## 5 格式门禁甲乙两尺（★先问"量的是哪个对象"）

- **甲（判据）**：`git archive HEAD` → 仓外 `export-<sha>/` → `/d/work/base/gopath/bin/gofumpt.exe -l .`（v0.12.0 go1.27.1）。
  **读数：总 20 行，全部 `.scratch\…`（被跟踪的探针件）；非 `.scratch` 产品面＝`grep -vcE '^\.scratch[\\/]'`＝0** ⇒ 我改的那枚**不在名单里**。
  **正控**（同尺打在另一份导出的 `zz_probe_unformatted.go` 上）：产品面读数 **1**（`zz_probe_unformatted.go`）⇒ 这把尺**能**看见坏 `.go`，不是空转。
- **乙（只作诊断）**：工作树 `gofumpt -l internal/` ⇒ 2 行，**都不是我的**：`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`。
  根因现量：这两枚工作树副本是 **CRLF**（各 225 个 CR），而 `text eol=lf` 让 `git status` 看不见它、`git archive` 出去的是 LF ⇒ **甲干净、乙脏**，正是"甲判据／乙诊断"的理由。我那一枚 `gofumpt -l` **零输出**（乙对单文件 rc=0）。

## 6 提交（逐枚；`git show --name-only` 逐枚核 `.go` 与越界）

- ① `62a1521cb411421c069b7dcf960fe1c2490b95a3` —— 本格**唯一 `.go`**：`git show --name-only` 列 1 枚
  `internal/tools/task_output_ac2_before_test.go`（`grep -c '\.go$'`＝1；写面外＝0）；1 file changed，+41/−26，全在注释／措辞。
- ② `ec8786e18…`（`ec8786e1`）—— 承载**本件**：`--name-only` 列 1 枚 `docs/evidence/s1/164-r1-impl.md`（`.go`＝0；写面外＝0）。
- ③ 把 ② 的提交号写进本节的那笔**也是 docs-only**（只列本件、`.go`＝0）——本行不写死它自己的号，命令可核。
- **正控（证明这把 `--name-only` 尺不是瞎子）**：同尺打在上面那笔**含 `.go` 的真提交**上 ⇒
  `git show --name-only --format="" 62a1521cb` 输出 **1 枚**＝`internal/tools/task_output_ac2_before_test.go`（`grep -c '\.go$'`＝**1**）⇒ 尺能看见 `.go`，"没列 `.go`"才是真读数。
- 越界：`internal/session/**`／`internal/memory/**`／`internal/risk/**`／`internal/panel/**`／`internal/config/**`／`cmd/**`／`scripts/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件＝**零字节**（写面只有上列两枚）。⛔ push；⛔ 台账；⛔ 票名／`Status:`／`-done`。

## 7 票面前提新鲜度（逐条现跑，★本仓票面腐坏已 ≥9 次）

- ✅ `grep -c '^- \[x\]'`＝**3**、`grep -c '^- \[ \]'`＝**3**、缩进不敏感总数＝**6**（票面 `:59` 的 3 勾／3 未勾**仍成立**；未勾＝AC#2b／AC#4／AC#5）。
- ✅ `internal/tools/task_output_ac2_before_test.go` 在，且旧署名确是"AC#2 的 BEFORE 读数"（§2 引的是改后文本）。
- ❌ **过期三处（不阻塞本格，报给编排者）**：① 票面 `:4`"`task.output` 全仓零命中"已不成立（现 17 枚 `.go`）；② 票面 `:56`"`RunAsync` 生产零调用点"已不成立（`cmd/wisp/run.go:1176`、`internal/tools/subagent_197.go:372`）；③ 票面 `:52`／`:59` 引的拒绝分支行号 `bridge.go:251`／`:247-253` 已漂到 **`:272-278`**（文案本体 `:276`）——**旧测试注释里那枚 `bridge.go:247` 就是这处腐坏**，本格已在措辞里去掉行号，改引函数与分支名。

## 8 残余与风险（★不许被"AC#2b 已交"压缩掉）

1. **名字没改**（文件名＋函数名）：它们是**历史键**——票面 AC#2b 按文件名点它、`.scratch/wisp/probes/164/v1/ac2-anchor-reading.sh:41,44` 与 CI 日志按函数名跑它。改名会让那支归档仪器的 HEAD 那一发变成 `no tests to run` 的**假绿**、也会让别腿的名册差集读出一枚"被吞用例"。⇒ 署名改在**注释与 Logf 标签**上，文件头第 36–39 行明写这一取舍。**若编排者要真名，另起一枚改名提交即可**（代价＝上面两条，先把归档脚本的 `-run` 改成新名）。
2. **第三节那颗牙射程窄**：它只钉 `BuiltinFSEntries` 不夹带 `task.*`，**不是**关于生产注册表的断言（生产在它之上另挂 task 家族）。旧注释那句"模型今天看到的目录没有 `task.*` 条目，所以模型连问都学不会"**对生产是假的**，已按射程改写（属于"改说明文字"，未动断言）。
3. **一把盲尺的教训（本轮自己踩的）**：`grep -c $(printf "\r")` 对一枚 **225 个 CR** 的文件读出 **0** ⇒ 判 CRLF 必须用 `tr -cd "\r" | wc -c`。§5 乙那两枚 `internal/agent/**` 的 CRLF 是**别人留下的工作树状态**（git 因 `eol=lf` 看不见），不属本格、未动。
4. **AC#4／AC#5 未碰**：取消"还在写的尾巴"仍无被测对象；契约面零字节。
