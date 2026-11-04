# 370-v1 — 独立裁决（非实现者，v1）

- 票：`.scratch/wisp/issues/370-assembler-maxoutputtokens-is-never-assigned-so-our-ceiling-is-never-sent.md`
- 裁决时 HEAD：`ddbec5ee6814aa6082cc1db2461d967501ec5376`（现跑 `git rev-parse HEAD`）
- 裁决对象：
  - `74f430d86`（AC#2 产码：run.go/loop.go 把 chat 角色的 `RoleSettings.MaxOutputTokens` 接到请求）
  - `e8cd575ec`（AC#0 证据 `370-a1-impl.md` + `internal/agent/maxoutput_370_ac0_test.go`）
  - `ddbec5ee6`（整备：`tools/mockllm/chat.go` 补 `recordRequest`，修 AC#2 读数面）
  - 另：`7d158caf5` 是 `cmd/wisp/task370_ac2_test.go` 的祖先提交（本裁决一并核）
- 方法：Git Bash（`D:\work\soft\Git\bin\bash.exe`），`cmd/wisp` 带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，全部读数由本腿现跑。
- **一句话结论：AC#0 / AC#1 / AC#2 / AC#4 PASS；AC#3 PASS，但须具名一条树级既有假红（未跟踪、被忽略的 `build/mut_rvw.go`，不属本票）。票可以收口。断言完整性 = 是（未削弱任何断言）。**

## 0. 判定表

| AC | 判定 | 一句话依据 |
|---|---|---|
| AC#0 | **PASS** | `370-a1-impl.md` 五小项齐全；父提交 `74f430d86^` 产码零赋值（grep rc=1），HEAD 唯一赋值 `loop.go:261`；正控 `request=1 -> status=truncated -> state=interrupted`；供应商自身上限用 golden `finish_reason=length` 实测仍可达 |
| AC#1 | **PASS** | 落地的确是〔甲′〕：`RoleSettings -> agent.Config -> Assembler -> Request -> max_tokens`；不是〔乙〕（未给 `Endpoint` 加字段、Assembler 不读 endpoint），不是〔丙〕（`StatusTruncated` 与映射原样保留）；与台账 `A1461` 的裁定一致 |
| AC#2 | **PASS** | 正控真跑：屏幕 `truncated` / 库 `interrupted` / exit 3；两向反控 PASS；本腿用 `-overlay` 独立摘掉 `loop.go:261` 那一跳，正控转红在 `task370_ac2_test.go:116`、反控转红在 `:183`，源文件零改动（还原核成立） |
| AC#3 | **PASS（带具名既有红）** | `go vet` rc=0；`internal/agent` ok 7.438s、`internal/llm` ok 40.162s、`internal/config` ok 2.044s、`cmd/wisp -run Test370` ok 4.444s、`tools/mockllm` ok 0.147s；`d22scan` rc=0 clean；长度门 rc=0 GREEN；gofumpt 甲乙两尺均空。`go build ./...` rc=1 **只**因未跟踪且被 `.gitignore:14 build/` 忽略的本机 staging 件 `build/mut_rvw.go`（HEAD 无此文件；`go build ./cmd/... ./internal/... ./tools/...` rc=0） |
| AC#4 | **PASS** | 票面五框仍 `- [ ]`、`Status: **ready-for-agent**`、文件名无 `-done`；三提交未改票面一字；未碰票 365；未碰 `thresholds.go` / golden / `frontend/**` / `design/**`。本腿未 push，只提交本裁决文件 |

## 1. 本腿复跑读数（四把尺 + 格式两尺）

### 1.1 `cmd/wisp` 定向（票面指定 PATH staging）

```console
$ export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
$ go test ./cmd/wisp -run 'Test370' -count=1 -timeout 220s -v
RC=0
--- PASS: Test370TruncatedRunRowsAnInterruptedTask (1.40s)
--- PASS: Test370TheInterruptedRowNeedsTheCeiling (2.92s)
    --- PASS: Test370TheInterruptedRowNeedsTheCeiling/no_ceiling_on_the_wire (1.44s)
    --- PASS: Test370TheInterruptedRowNeedsTheCeiling/ceiling_never_hit (1.48s)
PASS
ok  github.com/CarlosShao/wisp/cmd/wisp  4.444s
```

逐字读数（来自 `-v` 日志）：

| 运行 | 屏幕逐字 | 库 `task_log.state` | exit |
|---|---|---|---|
| `Test370TruncatedRunRowsAnInterruptedTask`（`max_output_tokens = 1`） | `"truncated"`（`task370_ac2_test.go:110`） | `"interrupted"` | 3 |
| `.../no_ceiling_on_the_wire`（不设上限） | `"completed"`（`:146`） | `"done"` | 0 |
| `.../ceiling_never_hit`（`max_output_tokens = 5000`，上线但撞不到） | `"completed"`（`:170`） | `"done"` | 0 |

### 1.2 `internal/agent -run 370`

```console
$ go test ./internal/agent -run '370' -count=1 -v
RC=0
--- PASS: Test370AC0HandSetCeilingRowsTruncated (0.13s)
    maxoutput_370_ac0_test.go:130: positive control: request.MaxOutputTokens=1; Loop.Run status="truncated"; task_log.state="interrupted"
--- PASS: Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling (0.01s)
    maxoutput_370_ac0_test.go:160: no ceiling on the wire (body has no max_tokens); recorded finish_reason=length -> status="truncated"
--- PASS: Test370RoleMaxOutputReachesTheAssemblerAndTheWire
--- PASS: Test370UnsetRoleMaxOutputSendsNoCeiling
--- PASS: Test370MaxOutputIsPerLoopNotGlobal
PASS
ok  github.com/CarlosShao/wisp/internal/agent  0.192s
```

### 1.3 `tools/mockllm`（独立 module）

```console
$ (cd tools/mockllm && go test ./... -count=1)
RC=0
ok  github.com/CarlosShao/wisp/tools/mockllm  0.147s
```

### 1.4 `d22scan`

```console
$ sh scripts/d22scan.sh
RC=0
d22scan: examined 408 production Go files under internal/ and cmd/ ...
d22scan: clean - no D22 ban violations; ...
```

### 1.5 gofumpt（甲 = 入库 blob，乙 = 工作树）

```console
$ /d/work/base/gopath/bin/gofumpt.exe -l <入库 blob 临时副本 7 件>
（无输出）BLOB_GOFUMPT_RC=0
$ /d/work/base/gopath/bin/gofumpt.exe -l <工作树 7 件>
（无输出）WT_GOFUMPT_RC=0
```

覆盖文件：`cmd/wisp/run.go`、`internal/agent/loop.go`、`internal/agent/maxoutput_370_test.go`、
`internal/agent/maxoutput_370_ac0_test.go`、`tools/mockllm/chat.go`、`tools/mockllm/server.go`、
`cmd/wisp/task370_ac2_test.go`。

### 1.6 AC#3 其余读数（本腿加跑）

```console
$ GOFLAGS= go build ./...
BUILD_RC=1   # 只红在未跟踪、被忽略的 build/mut_rvw.go（见 §6 残余 R1）
$ GOFLAGS= go build ./cmd/... ./internal/... ./tools/...
BUILD_TRACKED_RC=0
$ GOFLAGS= go vet ./cmd/wisp ./internal/agent ./internal/llm ./internal/config
VET_RC=0
$ go test -count=1 ./internal/agent ./internal/llm ./internal/config
ok  github.com/CarlosShao/wisp/internal/agent   7.438s
ok  github.com/CarlosShao/wisp/internal/llm     40.162s
ok  github.com/CarlosShao/wisp/internal/config  2.044s
$ sh scripts/check-path-length-budget.sh
PATHLEN_RC=0
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
```

## 2. AC#0 判定：PASS

`docs/evidence/s1/370-a1-impl.md`（337 行，`e8cd575ec`）+ `internal/agent/maxoutput_370_ac0_test.go`（166 行）覆盖五小项：

1. **现读清点**：父提交 `74f430d86^` 上 `git grep -n -E '\.MaxOutputTokens[[:space:]]*=' -- '*.go'` 现跑 rc=1（零命中）；HEAD 上产码唯一赋值 `internal/agent/loop.go:261: asm.MaxOutputTokens = opt.Config.MaxOutputTokens`。票面点名的 `prompt.go:132` 字段与 `llm.Request` 的 `provider.go:74` 字段都在；父提交 `prompt.go:253` 只读、抄的是零值。与 `370-a1-impl.md` §1 的 19 处逐条清单一致。
2. **resolver 落点**：现读 `internal/llm/resolver.go:186` 签名是 `func (r *Resolver) ResolveRole(name RoleName) (Endpoint, RoleSettings, error)`；`:202-206` 逐字 `MaxOutputTokens: role.MaxOutputTokens` 落进 **`llm.RoleSettings`**，不是 `llm.Endpoint`（`Endpoint` 只有模型目录语义的 `MaxOutput: spec.MaxOutputTokens`，`resolver.go:142`）。**票面 AC#0② 原句「已经落进 llm.Endpoint」措辞不准，`370-a1` §2 已逐字更正**；AC#0 的实质问题（配置值有没有出 resolver）答案是「有，在第二个返回值里」。这不是缺陷，是证据纠正。
3. **整条链表**：`370-a1-impl.md` §3.1 给出 HEAD 八跳全通的表；§3.2 给出父提交断点表，精确到「`agent.Config` 无字段 + `run.go:503 ep, _, err` 丢掉第二返回值」两处叠加。本腿现读核对：`run.go:516-525` 现为 `ep, roleSettings, err := res.ResolveRole(llm.RoleChat)` 且 `rt.roleSettings = roleSettings`；HEAD `Config.MaxOutputTokens` 在 `run.go:1560` 赋值；`loop.go:261` 抄进 Assembler；`prompt.go:253` 抄进 `Request`；`openaichat/request.go:75` 渲染 `MaxTokens`。
4. **正控**：`Test370AC0HandSetCeilingRowsTruncated` 手设 `Config.MaxOutputTokens = 1`，C5 由假 provider 返回 `StopMaxTokens`，`Loop.Run` 走真状态机，`Journal` 接真 SQLite；**没有构造 `Result`**（`res` 是 `Loop.Run` 返回值，行是 loop 自己经 journal 写入）。现跑读数逐字：`request.MaxOutputTokens=1; Loop.Run status="truncated"; task_log.state="interrupted"`。证明状态机本身好，父提交「到不了」确是配置没传。
5. **供应商自身上限的射程**：`Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling` 用 golden `internal/agent/testdata/golden/max-tokens-toolcall.sse`（第 17 行逐字 `"finish_reason":"length"`），配置不设上限、断言请求体无 `max_tokens`，仍得到 `status="truncated"`。三适配器读数：openai-chat/openai-responses 省略字段仍可回 `length`；anthropic 适配器把 0 填成本地 4096。故「我们设的天花板到不了」在父提交为真，「`StatusTruncated` 这个状态到不了」为假。与票面两条射程一致。

AC#0 只读性：`e8cd575ec` 只增 `docs/evidence/s1/370-a1-impl.md` 与 `internal/agent/maxoutput_370_ac0_test.go`，产码零字节。

## 3. AC#1 判定：PASS（落地的就是甲′）

**甲′ 实证（HEAD）**：

```text
run.go:516  ep, roleSettings, err := res.ResolveRole(llm.RoleChat)   // 保留第二返回值
run.go:525  rt.roleSettings = roleSettings                            // boot 副本
run.go:1498 _, roleSettings, rerr := llm.NewResolver(cfg, nil).ResolveRole(llm.RoleChat)
run.go:1503 roleSettings = rt.roleSettings                            // 失败回退 boot 副本
run.go:1560 MaxOutputTokens: roleSettings.MaxOutputTokens,            // 每任务 live 配置
loop.go:261 asm.MaxOutputTokens = opt.Config.MaxOutputTokens          // Config -> Assembler
prompt.go:253 MaxOutputTokens: a.MaxOutputTokens,                     // Assembler -> Request
openaichat/request.go:75 MaxTokens: req.MaxOutputTokens               // Request -> wire max_tokens
```

- **不是〔乙〕**：`internal/llm/resolver.go` 的 `Endpoint` 没有加 `MaxOutputTokens` 字段；`Assembler` 也没有从 endpoint 读。`Endpoint.MaxOutput` 仍是 `ModelSpec` 的模型目录上限（`:142`），语义未被挪用。
- **不是〔丙〕**：`StatusTruncated` 仍定义于 `internal/agent/loop.go:52`，仍在 `:62` 的 `Valid` 名册，截断分支 `:543` 仍返回它，映射 `:1129-1130 StatusTruncated -> memory.TaskStateInterrupted` 原样，`cmd/wisp/run.go:83-84` 仍给 exit 3。没有删状态、没有删映射。
- **与编排者裁定一致**：台账 `docs/reports/pending-and-issues.md:30623` 起 `A1461` 逐字「**370〔甲′〕落地**（`74f430d86`，3 文件）……**机主批准〔甲′〕，已兑现 `max_output_tokens` 真流到请求。**」；`370-a1-impl.md` §6.2 的三形代价表也推荐维持甲′。
- 契约面：`agent.Config` 是包内 Go 结构体，不在 `C1-C32` / `D1-D47` 清单；本票未改任何 C/D 契约文字，也未改 `llm.Request`（`provider.go:74` 父提交就有）。

## 4. AC#2 判定：PASS

### 4.1 正控（真组装根，非构造 Result）

`Test370TruncatedRunRowsAnInterruptedTask` 用 `probe365` 的真 `runTextTask` 入口（`main()` 调的同一入口）、真 DPAPI secret store、真 SQLite、真 openaichat 适配器；配置写 `[llm.roles.chat] max_output_tokens = 1`，mockllm 在 `max_tokens=1` 处切断回 `finish_reason:"length"`，适配器归一化为 `StopMaxTokens`，行从 memory DAO 读回。现跑 PASS，逐字：

```text
屏幕逐字 = "truncated"；库 task_log.state 逐字 = "interrupted"；exit = 3
```

### 4.2 反向判据 + 本腿独立牙齿证明

仓内 `Test370TheInterruptedRowNeedsTheCeiling` 有两向控制，均 PASS：不设上限 -> `completed`/`done`/exit 0；设 5000 但撞不到 -> 请求体确带 `max_tokens`、`completed`/`done`/exit 0。两者都证明「`interrupted` 那一行不是随便什么运行都有的」。

但票面 AC#2 的字面要求是「**把新加的那一跳摘掉 ⇒ 用例必须转红并逐名报红行**」。仓内证据没有留下这次源码变异的日志；本腿用 **Go `-overlay`** 独立补做（不改任何源文件、不碰测试）：

```console
$ sed 's/asm.MaxOutputTokens = opt.Config.MaxOutputTokens/asm.MaxOutputTokens = 0/' \
    internal/agent/loop.go > .scratch/370-v1-overlay/loop.go
$ OVERLAY="<repo>/.scratch/370-v1-overlay/overlay.json"
$ go test ./cmd/wisp -run 'Test370TruncatedRunRowsAnInterruptedTask$' -count=1 -overlay "$OVERLAY"
--- FAIL: Test370TruncatedRunRowsAnInterruptedTask (1.39s)
    task370_ac2_test.go:103: 屏幕逐字: ...（completed，1 轮，...）
    task370_ac2_test.go:110: AC#2(2) 屏幕逐字 = "completed"；库 task_log.state 逐字 = "done"
    task370_ac2_test.go:116: this run did not end at the output ceiling (screen word "completed", exit 0), ...
FAIL
$ go test ./cmd/wisp -run 'Test370TheInterruptedRowNeedsTheCeiling$' -count=1 -overlay "$OVERLAY"
--- FAIL: Test370TheInterruptedRowNeedsTheCeiling/ceiling_never_hit
    task370_ac2_test.go:183: the request body carries no max_tokens field, so this control never exercised the hop it is reversing: ...
FAIL
```

**还原核**：`git status --porcelain` 对 `internal/agent/loop.go` 与全部相关源文件为空；`sha256sum internal/agent/loop.go` = `b6600a0977370b15c34ae5cd3965bd539255d2caf8f5fffcd4630a7a12929f17`。整个变异只发生在 overlay 临时副本里，工作树从未被改。⇒ 两枚用例都有牙：摘掉那一跳，正控红在 `:116`、反控红在 `:183`。

### 4.3 AC#2(2) 两词记录与判决

| 场景 | 屏幕逐字 | 库 `task_log.state` 逐字 | 记录位置 |
|---|---|---|---|
| 撞我们设的天花板（`max_output_tokens=1`） | `truncated` | `interrupted` | `task370_ac2_test.go:110` |
| 完全不上限（`n=0`） | `completed` | `done` | `:146` |
| 上限 5000、上线但撞不到 | `completed` | `done` | `:170` |

**判决：本票只要求「逐字记下来 + 判要不要一致」，不需要在本票统一。** 理由：屏幕词来自 agent 状态机的 `StatusTruncated = "truncated"`（`loop.go:52`），库词来自 memory 名册的 `TaskStateInterrupted = "interrupted"`（`loop.go:1129-1130` 的刻意映射）；票 365/196 把「loop 词汇」与「memory 名册」分开是有意设计，且硬提醒 1 明确禁止把 `truncated` 加进名册。若机主要屏幕也说 `interrupted`，那是改屏幕文案/名册的独立任务，属 365 的射程。本票已把这一对词钉进测试日志，可测可复跑。

### 4.4 `ddbec5ee6` 的角色

`tools/mockllm/chat.go` 补 `s.recordRequest("chat", body, r.Header)`，并把 `takeFault()` 移到读体/记录之后，与 `messages.go:161-170` 的 `read -> record -> fault` 同序。这是**读数面的仪器修复**：`/__control/last_request?route=chat` 此前永远 `present:false`，导致 `ceiling_never_hit` 的 `max_tokens` 断言取不到 body（台账 `:30997` 已记「本环境红（mockllm captured no request）」）。修后 `cmd/wisp` 与 `tools/mockllm` 两处都绿。它没有放宽任何断言，也没有改产码语义。

## 5. AC#3 判定：PASS（带一条具名既有红）

| 门 | 读数 | 判定 |
|---|---|---|
| `GOFLAGS= go build ./...` | rc=1；10 行 `undefined:` 全在 `build/mut_rvw.go` | **非本票**：`git ls-files build/mut_rvw.go` 为空，`git cat-file -e HEAD:build/mut_rvw.go` 报 "exists on disk, but not in HEAD"，`git check-ignore -v` 指到 `.gitignore:14:build/`。`go build ./cmd/... ./internal/... ./tools/...` rc=0 |
| `go vet` 动过的包 | `./cmd/wisp ./internal/agent ./internal/llm ./internal/config` rc=0 | PASS |
| `go test -count=1`（分开报） | `internal/agent` ok 7.438s；`internal/llm` ok 40.162s；`internal/config` ok 2.044s；`cmd/wisp -run Test370` ok 4.444s；`tools/mockllm` ok 0.147s | PASS |
| `sh scripts/d22scan.sh` | rc=0，clean，408 枚产码 | PASS |
| `sh scripts/check-path-length-budget.sh` | rc=0，VERDICT GREEN | PASS |
| gofumpt 甲（入库 blob） | 空 | PASS |
| gofumpt 乙（工作树） | 空 | PASS |

`cmd/wisp` 整包按票面**不作判据**（MARK299 硬终止），只跑 `-run 'Test370'`，且带 sherpa PATH，未出现 `0xc0000135`。四数分开记录、未相减。

## 6. AC#4 判定：PASS

- 票面五格仍 `- [ ]`（行 38/44/50/56/60），`Status: **ready-for-agent**`（行 7），文件名无 `-done`。**没有翻框、没有改名、没有改 Status。**
- `git diff 74f430d86^ HEAD -- .scratch/wisp/issues/370-...md` 为空；三提交的 `--name-only` 没有票面文件。
- 票 365 未被碰：三提交无 `issues/365*` 路径。
- 禁面未被碰：三提交无 `thresholds` / `golden` / `frontend/` / `design/`。
- 本腿未 push；本腿只新建并提交本文件。
- 观察（不判红）：`git log --decorate` 显示 `origin/dev` 已在 `ddbec5ee6`，即三提交已发布。推送是编排者的授权动作；本腿无法也无需归因是谁推的，仅记录现状。

## 7. 断言完整性：是（未削弱）

`git show --stat ddbec5ee6`：

```text
 tools/mockllm/chat.go   | 16 +++++++++++-----
 tools/mockllm/server.go |  6 ++++--
 2 files changed, 15 insertions(+), 7 deletions(-)
```

- **没有碰 `cmd/wisp`**：`git show --name-only --format="" ddbec5ee6` = 上述两文件；`grep -c "cmd/wisp"` = 0。
- **测试文件逐字节未变**：`git rev-parse HEAD:cmd/wisp/task370_ac2_test.go` = `55c6744b1488e63bf80d1532c827f868bbaea90d`，`ddbec5ee6^:` 同一 blob；`git diff ddbec5ee6^ HEAD -- cmd/wisp/task370_ac2_test.go` 为空；工作树 sha256 = `4dc03c68bedb6a477610efb3e9afa3f7bd3fe03b460e64eba75d0b48b1e3565b`（与 `ddbec5ee6^` 取出的 blob 同串）。
- **断言原文仍在**（现 grep）：
  - 前提：`:115 if screenWord != string(agent.StatusTruncated)`；`:120 if code != 3`。
  - AC#2(1)：`:125 if !t370Interrupted(rowState)`，`:126-128 ... want memory.TaskStateInterrupted`。
  - 反控 1：`:149 if t370Interrupted(rowState)` 必须为假；`:154 if rowState != string(memory.TaskStateDone)`。
  - 反控 2：`:173 if t370Interrupted(rowState)` 必须为假；`:182 if !strings.Contains(body, "\"max_tokens\"")`；`:186 if rowState != string(memory.TaskStateDone)`。
  - `t370Interrupted`（`:92-94`）仍是 `rowState == string(memory.TaskStateInterrupted)`，正控与两反控用**同一谓词**。
- **没有 skip**：三个 370 测试文件 `grep -n "t.Skip\|SkipNow\|testing.Short"` 无命中。
- `ddbec5ee6` 的 diff 只做两件事：`recordRequest` 前移、`takeFault` 后移；fault 分支仍 `serveFault`，未删未弱化。`server.go` 只改注释。`messages.go:161-170` 是同一顺序的既有实现，chat 现在是向它对齐，不是反向放宽。

**结论：断言完整性 = 是。**

## 8. 残余（不阻塞收口）

- **R1（环境，非本票）**：本机工作树里有一个未跟踪、被 `.gitignore` 忽略的 `build/mut_rvw.go`（本机 staging 目录里的一枚副本，引用 `voiceTtsFiles` 等 `cmd/wisp` 的包内符号，故单独编不过）。它不在 HEAD，干净检出/CI 不会有它。`go build ./cmd/... ./internal/... ./tools/...` rc=0。票面 AC#3 的 `go build ./...` 在本机因此无法显示 rc=0；此红与 370 无关，具名交编排者。既有先例：`364-ac-verdict.md:7`、`378-ac-verdict.md:25` 记录同一假红。
- **R2（证据形态，已由本腿补足）**：仓内 AC#2 的「反向判据」是配置面负控 + 断言，不是实现者留下的「摘源码一跳 -> 转红 -> 还原」日志。本腿用 `-overlay` 补做了字面变异（§4.2），红行、还原核都在本文件里；若编排者要求实现者自己的变异记录，这一条是本腿的独立替代，不改变 AC#2 的实质成立。
- **R3（轻微，非削弱）**：`ceiling_never_hit` 的线上断言只查 `"max_tokens"` 子串，不查值 `5000`；正控的结局（`truncated`/exit 3）已把值 `1` 钉死。该文件自 `7d158caf5` 起就是这样，`ddbec5ee6` 未动。
- **R4（设计，非缺陷）**：屏幕 `truncated` 与库 `interrupted` 仍不一致。按票面 AC#2(2)，本票只记录并判「要不要一致」；判：**本票不统一**，统一属 365/屏幕文案的独立任务（硬提醒 1 禁止加名册词）。
- **R5（票面措辞，已纠正）**：AC#0② 原句说值落进 `llm.Endpoint`，实际落进 `llm.RoleSettings`；`370-a1-impl.md` §2 已逐字更正。不改票，只记录。
- **R6（观察）**：`origin/dev` 已在 `ddbec5ee6`（三提交已发布）。本腿未 push。

## 9. 复跑清单

```console
$ cd "/d/work/workspace/projects plans/Wisp"
$ git rev-parse HEAD                                                   # ddbec5ee6
$ export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
$ go test ./cmd/wisp -run 'Test370' -count=1 -timeout 220s -v          # rc=0, ok 4.444s
$ go test ./internal/agent -run '370' -count=1 -v                      # rc=0, 5 PASS, ok 0.192s
$ (cd tools/mockllm && go test ./... -count=1)                         # rc=0, ok 0.147s
$ sh scripts/d22scan.sh                                                # rc=0, clean
$ sh scripts/check-path-length-budget.sh                               # rc=0, GREEN
$ /d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/run.go \
    internal/agent/loop.go internal/agent/maxoutput_370_test.go \
    internal/agent/maxoutput_370_ac0_test.go tools/mockllm/chat.go \
    tools/mockllm/server.go cmd/wisp/task370_ac2_test.go                # 空
$ GOFLAGS= go build ./cmd/... ./internal/... ./tools/...                # rc=0
$ GOFLAGS= go vet ./cmd/wisp ./internal/agent ./internal/llm ./internal/config   # rc=0
```

## 10. 裁决

**AC#0、AC#1、AC#2、AC#4 PASS；AC#3 PASS（R1 具名、非本票）；断言完整性 = 是。票 370 可以收口。**
