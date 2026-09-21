# 102 — C26 展开改写：独立对抗验收裁决表

**验收代理**：`acceptor-ticket102`（独立于 `agent-ticket102`；未参与实现）
**被验收物**：票面 `.scratch/wisp/issues/102-c26-expansion-rewrites-onto-another-tree.md`，AC#1–AC#5 五框
**三枚 commit**：`0117459`（AC#1 红用例 + 判 (B) 原句）→ `a66aadf`（实现 + 四条腿 + AC#5 判据）→ `a1613d9`（票面读数，HEAD）
**本表由验收方出（票面 Rules 明写"不是自裁"）**

## 0. 复现环境（所有读数来自哪棵树）

四枚快照一律 `git archive <rev> | tar -x -C <dir>`，**全在 `/tmp`，仓内零 worktree/checkout**（A38④），目录名带本会话后缀 `ac2`；
每枚都拷入 `third_party/sherpa-onnx`（票 98 的账：缺 dll 会让 `cmd/wisp` 加载期 `0xc0000135`）。

| 快照 | rev | 用途 |
| --- | --- | --- |
| `/tmp/wisp-102-pre-ac2` | `0117459`（修前） | AC#1 的"修前必须红" |
| `/tmp/wisp-102-head-ac2` | `HEAD` = `a1613d9` | AC#1 修后绿 + AC#4 全部门禁 |
| `/tmp/wisp-102-mut1-ac2` | `a66aadf` | AC#3 双向变异（M1、M2 同一快照串行，M1 先还原再打 M2） |
| `/tmp/wisp-102-instr-ac2` | `HEAD` | AC#5 静态判据的反证（探针 A/B） |

**共享工作树 `D:\work\workspace\projects plans\Wisp` 只被读过（`grep`/`sed`/`git show`），没在里面跑任何门禁或改码。**
`internal/winsec/`、`internal/perm/`、`cmd/wisp/` 此刻有别的代理在写 —— 本表全部绿读数**来自 `a1613d9` 的仓外快照**，不是来自共树。
⚠ 本代理交件时共树 HEAD 已前进到 `3eea08b`（票 101 结案 + 台账 A76），并且 `internal/winsec/seam_guard_windows_test.go` 以**未跟踪**状态躺在共树里 ——
这些**不在本表射程内**：AC#4 的四包绿**只属于 `a1613d9` 那一棵树**，任何把这段读数当 HEAD 现状引用的做法都要重跑。〔独立复现（HEAD 移动由 `git log` 亲测）〕

**档位图例**：〔独立复现〕= 本代理亲自敲命令、读到真实 exit code；〔日志＋归档，我抽验〕= 数字来自实现方，本代理核对文件/字段对得上；〔仅自述，不背书〕。

---

## AC#1 —— 先复现再修（修前必须红，红名 + 断言原文）

> 票面原句：「把验收代理的 PROBE B2 做成**包内可重跑的用例**（不是临时探针）：输入含 `%VAR%`/前导 `~` ⇒ 断言"结果树与请求树**必须同一棵**，否则必须报错"。⚠ 这条用例在**修之前必须红**（红名+断言原文进票面）。」

| # | 证据 | 命令 | 读数 | 档位 |
| --- | --- | --- | --- | --- |
| 1 | **修前红（本代理自己复现，不是抄实现方）** | `/tmp/wisp-102-pre-ac2$`<br>`go test -count=1 -v -run TestC26ExpansionMustNotRewriteOntoAnotherTree ./internal/risk/` | **rc=1**；`=== RUN` **3**（1 父 + 2 子）；`--- FAIL` **3**（父 1 + `percent_env_var_spelling`、`leading_tilde_spelling`）；断言原文里 **`FAIL-OPEN:` 字样出现 9 次**；`--- SKIP` 0 | 〔独立复现〕 |
| 2 | 三条腿逐条点名（修前） | 同上，读日志 | `risk.c26Pipeline.Resolve` / `winsec.ResolvePath (sealing caller's view)` / `risk.syncSet.resolveTarget (fs.write sync judgement)` 各红；`%VAR%` 腿 `tree acted upon = …\aelsewhereb\artifacts` 而 `caller's spelling = …\a%WISP102PROBE_TARGET%b\artifacts`；`~` 腿 `tree acted upon = C:\Users\swq\wisp102-tilde-probe\artifacts` —— **全部 `err = nil`**，即 PROBE B2 的形状原样落在包内用例里 | 〔独立复现〕 |
| 3 | **修后绿，同一命令** | `/tmp/wisp-102-head-ac2$`<br>`go test -count=1 -v -run TestC26… ./internal/risk/` | **rc=0**；`=== RUN` **3**（>0，满足"不匹配的用例名也会打印 ok"这条要求）；`--- PASS` **3**；`--- FAIL` 0；`FAIL-OPEN` **0**；打印的是 `refused as required: risk: expansion moved this path onto a tree the caller did not name …` | 〔独立复现〕 |

**结论：通过。** 三点补充（都不翻转判定）：

1. 票面 AC#1 段写"**6 条断言红**"，本代理实测 **9 条**（`~` 那条腿在同一个子用例里跑了 `\` 与 `/` 两种拼写 ⇒ 3 腿 × 2 拼写 = 6，加 `%VAR%` 的 3 条 = 9）。方向是**实现方少报了自己的红量**，属记账误差、不是假绿；`=== RUN`=3 与红名两条与票面逐字一致。→ 记 **R-102-1**。
2. 用例是**包内可重跑**（`internal/risk/pathresolver_expansion_test.go`，`t.TempDir()` + `t.Setenv`），不是临时探针；无网络、无特权。
3. "请求树"取自包内 `lexCanonical`（`Abs`+`Clean`），没有新增第二套规范化 —— 本代理核对该文件未新增 `filepath.Abs`/裸 `Clean` 调用点，`sh scripts/d22scan.sh` 纯净快照 rc=0（见 AC#4），静态扫描也没多一处表面。

---

## AC#2 —— 选定 (A)/(B) + 每个生产调用点的处置表（不许只改 winsec 那条腿）

> 票面原句：「选定 (A) 还是 (B)，并给**每个生产调用点**的处置表（`grep` 出所有 `Resolve(` 非测试命中，逐点写"要不要消费改写标记/是否拒绝"）。**不许只改 winsec 那条腿。**」

### 2.1 命中集合：本代理自己跑，与实现方的表逐条对照

```
grep -rn "risk\.Resolve(" --include=*.go . | grep -v _test.go
./internal/tools/paths.go:81:	res, err := risk.Resolve(raw, p.exceptions)     ← 1 命中
grep -rn "Resolve(" --include=*.go internal/risk | grep -v _test.go
pathresolver.go:105(定义) syncdirs.go:140 syncdirs.go:198 syncdirs.go:222 winsec_c26.go:40   ← 包内 4 处消费
```

**非测试命中集合 = 上表 5 条腿，与实现方处置表完全一致：多一条没有，少一条也没有。**〔独立复现〕
⚠ 行号漂移：票面表第 5 行写 `internal/tools/paths.go:56`，HEAD 实测是 `:81`；第 1 行写 `winsec_c26.go:34`，实测 `:39-44`。指向的函数没错，是行号没跟着自己的改动重算。→ 记 **R-102-2**。
（同类：票面 AC#2 表第 2/3/4 行的 `syncdirs.go:140/197/211` 实测 `140/198/222`，其中 `res.Rewritten` 的 case 落在 `:143`。）

### 2.2 逐腿：读**测试体**而不是读测试名，判"钉的是行为结果还是字段等于某值"

| 腿 | 处置（实现方口径） | 钉住它的用例（本代理读测试体确认） | 断的是行为还是字段 | 档位 |
| --- | --- | --- | --- | --- |
| 1 `winsec_c26.go:39-44` `c26Pipeline.Resolve` | **拒**：`return res.Actable()` | `pathresolver_expansion_test.go` 的 `requireSameTreeOrRefusal`，对三条腿分别要求 `err != nil`；winsec 侧断 `winsec.ResolvePath` 报错（`%w` 传播，本代理核对该错误串出现在 HEAD 输出里） | **行为结果**（"报成功"这件事本身被否掉） | 〔独立复现〕 |
| 2 `syncdirs.go:140` `syncSet.add` | **消费不拒**：`case err == nil && res.Rewritten:` → `canonical=false` + `logf` | `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`：断 `rewritten.roots[0].canonical == false` **且** `rewritten.complete == false`，并**自带对照腿**（干净根必须 `canonical=true` 且 `complete=true`，否则 `t.Fatalf`）；本代理另外核对它**先 `Skipf` 于"平台给不出 handle-resolved 形"的机器上**，所以 Windows 上这条是真断言 | **行为结果**：`complete` 是兜底网的开关（`finalize()` 读它），不是记账字段 | 〔独立复现〕 |
| 3 `syncdirs.go:198` `resolveTarget` | **拒**（fail-closed，上游按 sync-suspect） | AC#1 用例的第三条断言腿（`risk.syncSet.resolveTarget … refused as required`），且 `match()` 的 `s.resolveTarget` 错误路径就是 sync-suspect（`syncdirs.go:407`） | **行为结果** | 〔独立复现〕 |
| 4 `syncdirs.go:222` 祖先重解析 | **消费**：`ares.Actable()` | 没有**独立**用例把"祖先被改写"这一支单独钉住 —— 全链 `anceCanon` 替换 + AC#5 文本判据要求 `syncdirs.go` 出现 `Actable(`；实现方自己写明"祖先形是从 C26 自己出来的，今天恒不触发" | **半行为半结构**：分支今天不可达（本代理认同其理由），所以只能是结构判据。→ 记 **R-102-3** | 〔日志＋归档，我抽验〕（可达性是本代理读 `resolveTarget` 体推出的，非实测） |
| 5 `paths.go:81`（`resolve`；`risk.Resolve(` 唯一非测试命中） | 根=**记账 + 不可确认则不授权**；target=**不拒** | `TestPathCanonicalizerAccountsForRewrittenRoots`：断 `RewrittenRoots()` 恰 1 条且含拼写/展开树/构造名；断"改写且未确认"的根 `Roots()` **为空**（=什么也不授权）；断"确认过的改写根照样授权"（`InAllowlist(c)==true`，即**不许变成 (A)**） | **混合**：`RewrittenRoots()` 那半是**字段/ getter 形状**；`Roots()` 空与 `InAllowlist` 真半是**行为结果** | 〔独立复现〕（本代理在 M2 快照里看到它红，说明断言不是装饰） |

### 2.3 "不许只改 winsec 那条腿"是否成立

成立：`internal/risk/` 三处 + `internal/tools/` 一处都动了，HEAD 的 `git show --name-status a66aadf` 六个文件 = `pathresolver.go` `syncdirs.go` `winsec_c26.go` `paths.go` + 两个新测试文件。
`internal/winsec/`、`internal/secret/`、`cmd/wisp/`、`internal/perm/` **确实一字未动**（diff 文件清单里没有它们）。〔独立复现〕

**AC#2 结论：通过（附两条记账/覆盖残留 R-102-2、R-102-3）。**

---

## AC#3 —— 双向变异（本代理自己做，不照抄）

> 票面原句：「变异：把修法退回"照旧展开且不记账" ⇒ AC#1 红；另做一发**反向**：把 `Result` 的标记字段变成恒 `false` ⇒ **也必须有用例红**（否则 (B) 是纸面承诺）。」

快照 `/tmp/wisp-102-mut1-ac2`（`git archive a66aadf`）。每发都在**同一条 `&&` 链**里 `grep -n` 打印被改后的整行，且 `go build ./internal/risk/` rc=0 之后才跑测试。

**M1「`Actable()` 不拒」** —— `pathresolver.go:94`，落地证据：`94:	if false && r.Rewritten {`；`M1 BUILD rc=0`；`go test -count=1 -v ./internal/risk/ ./internal/tools/` **rc=1**（`=== RUN` 263，`--- SKIP` 1）：

```
--- FAIL: TestC26ExpansionMustNotRewriteOntoAnotherTree            (2 子用例红)
--- FAIL: TestC26RewriteAccountIsRecorded                          (3 子用例红: percent_env_var / dollar_env_var / leading_home_tilde)
```

**M2「`Rewritten` 恒 false」** —— 先还原 M1（`grep -n` 打印回 `94:	if r.Rewritten {`），再打 `pathresolver.go:116`，落地证据：`116:	res := Result{Spelling: input, Rewritten: false, Rewrites: kinds}`；`M2 BUILD rc=0`；同两包 **rc=1**（`=== RUN` 263）：

```
--- FAIL: TestC26ExpansionMustNotRewriteOntoAnotherTree            (2 子)
--- FAIL: TestC26RewriteAccountIsRecorded                          (3 子)
--- FAIL: TestC26RewrittenSyncRootDoesNotDisarmSuspectNet          (+2 条断言: canonical=true / complete=true)
--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots           (internal/tools)
```

| 对照 | 实现方登记 | 本代理实测 | 判定 |
| --- | --- | --- | --- |
| M1 红几条 | 5 条（AC#1 两条子 + IsRecorded 三条子）；`IsConsumedAtEverySecurityLeg` 与 `internal/tools` 那条**不红** | 5 条子用例（2 父），不红的两条**确实不红** | **对得上** 〔独立复现〕 |
| M2 另红 | 2 条（sync 腿 + tools 腿） | 恰好这 2 个父用例新增转红，且 sync 腿两条断言原文一致 | **对得上** 〔独立复现〕 |

**判语（按简报要求分辨"结果被钉住"还是"某行被钉住"）**：M1 单独关 `Actable` 时 `internal/tools` 全绿，**这不是测试无效** —— 授权腿按设计读 `Rewritten` 而不是 `Actable()`（两条腿的语义不同：一个动手并报成功，一个只是判定树==执行树），属**分层冗余缺失的那一格由 M2 兜住**：M2 把账擦掉时 tools 腿立刻红。也就是说"账"与"闸"各有其钉住的用例，两发变异合起来把 (B) 的两半都钉住了 ⇒ 不是纸面承诺。**M1 的 `IsRecorded` 3 条子用例也红**，说明该用例是**经过闸**的（`requireRewritten` 里断 `Actable()` 返回 `ErrRewrittenPath`），不只是看字段。

**AC#3 结论：通过。** 变异后快照已弃用（未合回共树）。

---

## AC#4 —— 影响面回归 + 门禁

> 票面原句：「`go test -count=2 ./internal/risk/ ./internal/winsec/ ./internal/tools/ ./internal/memory/` 全绿，并逐条点名 SKIP/FAIL；`sh scripts/d22scan.sh` 纯净树 rc=0、`ban #6`/`ban #8` 对 `frontend/` 的数**不许降**。」

### 4.1 四包 `-count=2 -v`（HEAD 快照，本代理自己跑，`-json` 未用，靠 `tac`+`awk` 分块归包）

| 包 | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` |
| --- | --- | --- | --- | --- |
| `internal/risk` | 324 | 322 | **0** | 2 |
| `internal/winsec` | 58 | 58 | **0** | 0 |
| `internal/tools` | 202 | 202 | **0** | 0 |
| `internal/memory` | 134 | 132 | **0** | 2 |

命令：`cd /tmp/wisp-102-head-ac2 && go test -count=2 -v ./internal/risk/ ./internal/winsec/ ./internal/tools/ ./internal/memory/` ⇒ **rc=0**（后台任务 exit code 0）。
`=== RUN` == 2 × 不同名 四包全部成立（324/58/202/134）。**FAIL = 0 条，无需点名。**〔独立复现〕

SKIP **逐条点名**（都是 `-v` 才打印得出来，非 `-v` 会静默）：

- `internal/risk/TestSyncRegistryProbeLive`（`syncdirs_test.go:347`：`no registry-grade sync record on this machine (HKCU Accounts without UserFolder …)`）
- `internal/memory/TestSubprocessCrashWriter`
- **是否本票引入**：本代理在 **`0117459`（修前）快照**里跑 `go test -count=2 -v -run TestSyncRegistryProbeLive ./internal/risk/` ⇒ 两轮**同样 2 条 SKIP**，理由串逐字相同 ⇒ 环境闸，非回归。`memory` 那条本代理只核对"名字与实现对上 + 修前快照里同包无新增 SKIP"，未单独复跑它的子进程机制。〔独立复现〕／memory 那条〔日志＋归档，我抽验〕
- **`TestResolvePerCallBudget`**：四包 FAIL=0 里含它，本代理**没量到红**，也未拿它当借口。

### 4.2 其余门禁（HEAD 快照）

| 门禁 | 命令 | 读数 | 档位 |
| --- | --- | --- | --- |
| 格式 | `gofmt -l .` | **空**，rc=0 | 〔独立复现〕 |
| 格式 | `"$(go env GOPATH)/bin/gofumpt" -l .` | **空**，rc=0（第一次跑 `gofumpt` 裸命令是 `command not found` rc=127，改绝对路径才算数——**PATH 里的形状与实现方不同，但读数一致**） | 〔独立复现〕 |
| vet | `go vet ./internal/{risk,tools,winsec,memory}/` 逐包 | **4 × rc=0** | 〔独立复现〕 |
| vet | `GOOS=linux go vet ./internal/{risk,tools,winsec,memory}/` 逐包 | **4 × rc=0** | 〔独立复现〕 |
| vet | `GOOS=linux go vet ./...`（仓根） | **rc=1**，串是 `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files …` ⇒ **既有仪器坑，非回归** | 〔独立复现〕 |
| cmd/wisp | `go test -count=1 -v ./cmd/wisp/` | **rc=1**、`exit status 0xc0000135`、`=== RUN` **0 条** ⇒ 与票 98 的加载期缺 DLL 同形，与实现方登记逐字一致；本代理判"仪器不可跑"，**不判本票回归**（本票未动 `cmd/wisp`，diff 清单可证） | 〔独立复现〕 |

**核对实现方"按包 vs 整仓"这句话有没有说谎**：它写的是"按包跑四包 rc=0（仓根 `./...` 的已知坑不在此列）"。本代理两个读数**分别**跑过：按包 8 × rc=0、整仓 rc=1 ⇒ **没有把整仓读数说成按包**，也没有把按包绿说成整仓绿。〔独立复现〕

### 4.3 D22 静态扫描（**逐字同形** `sh scripts/d22scan.sh`，不是 `go run ./tools/d22scan`）

`cd /tmp/wisp-102-head-ac2 && sh scripts/d22scan.sh` ⇒ **rc=0**，18.6s 量级（含 `go test ./...` 正对照 `ok tools/d22scan 3.089s`）。台账 8 数：

```
bans #1-5 internal/=197  bans #1-5 cmd/=20  ban #6 frontend/=37  ban #7 internal/tools/=17
ban #8 design/=16  ban #8 frontend/=37  ban #8 internal/=340  ban #8 cmd/=26
```

与实现方登记**逐字相同**；`ban #6 frontend/` **37**、`ban #8 frontend/` **37** 相对票 94 参照 **37/37** 未降；`#8 internal/` 340、`#8 cmd/` 26 相对参照 335/25 只增不减。〔独立复现〕

**AC#4 结论：通过。**

---

## AC#5 —— "谁消费了这个字段"的静态可检查判据（防"扫空=绿"）

> 票面原句：「若走 (B)：给"谁消费了这个字段"一条**静态可检查**的判据……而不是"注释里提醒后人记得读"。」

判据实体：`internal/risk/pathresolver_rewrite_account_test.go:TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（两半：逐腿点名 needle + 全仓扫 `.Canonical` 字段读取者，除生产者外必须出现 `Actable(` 或 `.Rewritten`，一条没扫到就 `t.Fatalf`）。
`readsCanonicalField` 用**后随字符**把 `Canonicalize` 方法排掉 —— 本代理认同这正是"按常量名 grep 对内联裸字面量失明"那一族的正确解法：**它匹配的是形状的读法而不是符号名表**。

### 5.1 两枚反证（本代理自己在 `/tmp/wisp-102-instr-ac2` 里做，实现方没做过）

| 探针 | 做法 | 读数 | 说明 |
| --- | --- | --- | --- |
| **A：新第六条腿直接抓 `Canonical`** | 新增 `internal/risk/zz_probe_ticket102.go`：`res, _ := Resolve(...); return res.Canonical`（不读账），`go build ./internal/risk/` rc=0 | **rc=1**：`pathresolver_rewrite_account_test.go:200: …\zz_probe_ticket102.go:9 reads Result.Canonical without reading the rewrite account …` ⇒ **判据真能咬住新增消费者，且报到行号** | 〔独立复现〕 |
| **B：扫空** | 把 `canonicalReaders` 的 `repo` 指向不存在的目录（`no-such-dir-102`），使 `.Canonical` 命中数为 0 | **rc=1**：`pathresolver_rewrite_account_test.go:192: the scan found no .Canonical reader at all: the instrument is broken, not the code` ⇒ **"扫空=绿"这条路被堵死** | 〔独立复现〕 |

两枚探针跑完即删/弃，未进共树。

### 5.2 判据的射程边界（如实登记，不当它没写过的话）

- **needle 半只挡"符号被删/搬走"**：M1（只把 `if r.Rewritten` 改成 `if false && …`）**不红** —— 实现方自己在票面 AC#5 段末承认了这点，本代理复现成立。
- **needle 半不覆盖祖先重解析那一支的可达性**（见 R-102-3）。
- **扫描范围有 6 个 `SkipDir`（`.git/.scratch/node_modules/dist/build/testdata/design/docs`）+ 整个 `tools/` 前缀被排除**。`testdata/`/`docs/` 里不该有生产判腿，可**"新增一条藏在 `testdata/` 或 `tools/` 下的第六条腿"是这条判据扫不到的**（当前不存在，故不翻转判定）。→ 记 **R-102-4**。

**AC#5 结论：通过（有条件：R-102-4）。** 这条不是"注释里提醒后人"——它是一枚被本代理双向攻过的不变式。

---

## 6. 契约引句逐字核对（硬项）

| 引用 | 本代理读到的当前真实文本 | 判定 |
| --- | --- | --- |
| `docs/specs/SPEC-06-security-gatekeeping.md:50` | 第 50 行逐字：`展开(env / ~) → 绝对化 → Clean → 打开句柄取 GetFinalPathNameByHandle(VOLUME_NAME_DOS) 真实路径` | **存在且行号精确**（`grep -n` 命中 50） |
| `docs/PLAN.md:2375` | 2374 行 `**修复 = C26 PathResolver，唯一入口**：`，**2375 行起** `展开(env/~)` → `绝对化` → `Clean` → **打开句柄取 … 得到真实路径** | **存在，行号指到 pipeline 那一行** |
| `docs/PLAN.md:1376` | 1376 行 `| **C26** | **PathResolver** | **唯一**的路径规范化入口：展开 → 绝对化 → 取句柄真实路径（Win: …）→ **拒绝 reparse point…**。白名单与黑名单**都必须**经它 |` | **存在且逐字** |
| 支撑句 `SPEC-06:63-66`（A 档黑名单本身用 `~/.git-credentials`、`%APPDATA%\wisp\config.toml` 拼写） | 63-66 行确有 `~/.git-credentials` · `%APPDATA%\wisp\config.toml` 等拼写 | 存在 |

⇒ **"展开是契约 mandates 的 pipeline 第一步"这个结论逐字站得住**；(A)「入口直接拒绝 `%VAR%`/`~`」确实等于删掉契约明写的一步，并且会让"黑名单经 C26"落不了地。
⇒ **"本票是修实现、未触发 D22 停手"成立**：`git show --name-status a66aadf` 与 `a1613d9` 里 `docs/PLAN.md`、`docs/specs/*` **一字未动**（本代理核对 diff 文件清单）。本代理亦未改任何 `docs/PLAN.md`/`docs/specs/*.md`。

**档位：〔独立复现〕**（三处文本均由本代理 `sed -n` 直读当前工作树）。

---

## 7. 争议裁决：`internal/tools/paths.go` 该不该动、改动形状对不对

### 7.1 自述 vs diff 逐行核对（`git show a66aadf -- internal/tools/paths.go`）

| 实现方自述 | diff 实测 | 结论 |
| --- | --- | --- |
| 只加 `rewritten` 字段 + `RewrittenRoots()` | `+	rewritten []string`（struct 内）与 `+func (p *PathCanonicalizer) RewrittenRoots() []string`（`mu.RLock` 下返回副本） | **对** |
| `resolve` 返回 `risk.Result` | 签名 `(string, error)` → `(risk.Result, error)`，唯一调用点是构造函数与 `Canonicalize` | **对** |
| `NewPathCanonicalizer` 对"改写且 OS 认不出"的根 fail-closed | `if res.Rewritten { …; if !res.Resolved { p.unusable = append(…); continue } }` ⇒ 该根**不进 `roots`** | **对** |
| 未动 `bridge.go`/`fs.go`/`fs_write.go`/`mode.go` 一行 | `git show --name-status a66aadf` 六文件里**没有**这四个文件 | **对** |

**自述与 diff 相符，无越权、无扩大范围。**〔独立复现〕

### 7.2 有没有"调用方可控的选择器" / 放行侧是否比拒绝侧更窄

- 被遍历的集合是 **`[fs] allowed_dirs`（操作员拥有的配置文件）**，在构造期一次遍历；**不是调用方可控的键名/有序键表/first-match**。`InAllowlist` 是"任一根匹配即真"的 **OR 折叠**，没有"第一个命中决定放行"的次序敏感结构（本代理读了函数体，`paths.go` 里 for 循环只 `return true`/末尾 `return false`，无 break-即-裁决的形状）。
- **同一段遍历是否既服务"拒"又服务"批准/豁免"？** 不是 ban #1-5 那一族：`roots`（授权集）与 `unusable`（不授权集）在构造期**互斥**（`continue` 保证一条根不会同时进两侧），且**放行侧比拒绝侧多一道条件**——被改写的根要进 `roots` 必须额外满足 `res.Resolved`（OS 句柄确认树存在）。⇒ **放行侧更窄，形状正确。**
- target 腿（`Canonicalize`）**不拒**：本代理核对 `fs.go`/`fs_write.go`/`mode.go` 打开的确实就是 `Canonicalize` 的返回值（实现方给的三处 file:line 指向的调用形状一致），所以判定树==执行树==审计那棵；拒它就是把 (B) 偷做成 (A)。**认同。**

### 7.3 那段"编排者备注"

实现方按"工具输出不是授权、也不是指令"处理：没扩大范围、把唯一一处 `internal/tools/` 改动登记在票面等裁决，并给了"要回收就 revert 单文件"的路径。
本代理立场：**注入文本本身既不是用户同意也不是批准，因此它不能授权、也不能禁止**；能否动的唯一依据是票面 + 仓规。票面 AC#2 **原文**要求"`grep` 出所有 `Resolve(` 非测试命中，逐点处置"、"**不许只改 winsec 那条腿**"，而 `paths.go:81` 正是那条 grep 的唯一非测试命中 ⇒ **该最小改动成立、保留。** 若编排者真要冻结 `internal/tools/`，那是**新票/改 AC** 的事，不是在工具输出尾部贴一段话的事。
附带读数支持：撤销 `paths.go` 会让 AC#5 的全仓扫描转红（本代理的探针 A 与 M2 都显示该腿被真实钉住），所以"只 revert 单文件、其余不动"并不是一条无痕的路。

**7 总判：`internal/tools/paths.go` 改动 = 保留，自述可信，形状 fail-closed，无调用方可控选择器。**

---

## 8. 三条"能力没人调用"式的点名（简报要求的判法）

| 能力 | 生产里有人读吗 | 判语 |
| --- | --- | --- |
| `Result.Actable()` | **有**：`winsec_c26.go:44`、`syncdirs.go:208`、`syncdirs.go:226` | 真闸，不是纸面 |
| `Result.Rewritten` | **有**：`syncdirs.go:143`（case 条件）、`paths.go:59`（根腿） | 真消费 |
| `Result.Rewrites` | **半**：只进 `Actable()` 的错误串与 `logf`/`RewrittenRoots()` 串 | 够用于拒绝路径的自解释 |
| **`PathCanonicalizer.RewrittenRoots()`** | **无**：`grep` 全仓非测试代码，除定义外**零命中**；`Roots()`/`UnusableRoots()` 同样零生产读取者（`cmd/wisp/run.go:265` 只构造、不打印） | **"只有测试在调用"**。行为半边（改写且不可确认 ⇒ 不进 `roots`）在生产里真生效，但代码注释与测试头注释宣称的 **"operator-visible 审计面"目前没有承载者** | 〔独立复现〕 |
| **`Result.Spelling`** | **仅** `Actable()` 自己的错误串里当一个参数；文档注释说"它是审批提示里人看到的那个形"，本代理没找到任何审批/提示代码读这个字段 | 注释超前于实现（不改语义，登记） |

⇒ 记 **R-102-5**（`RewrittenRoots()` 无生产打印者）与 **R-102-6**（`Spelling` 的"审批提示人看到的形"这句无人兑现）。两者都**不构成本票 AC 的 fail-open**（拒绝与不授权两半都是行为性的），但都不该被"审计可见"这种词糊过去。

---

## 9. 残留清单（验收期不改码，逐条该有自己的票）

- **R-102-1**（记账误差，轻）：票面 AC#1 段"6 条断言红"实为 **9 条**（`~` 子用例跑 `\`/`/` 两种拼写 ×3 腿）。方向是实现方少报。补救：下次追加时把该行订正为实测值并注明"由 `acceptor-ticket102` 复算"。
- **R-102-2**（引用漂移，轻）：AC#2 处置表与"需要协调"段的 file:line 未随自己的改动重算（`paths.go:56`→`:81`，`winsec_c26.go:34`→`:39-44`，`syncdirs.go:197/211`→`198/222`）。补救：按 HEAD 重算一次，或改成"函数名 + 相对锚点"。
- **R-102-3**（覆盖缺口，中）：腿 4（祖先重解析走 `Actable()`）**只有结构判据、没有行为用例**——该分支今天不可达，所以只能钉符号；但一旦哪天祖先形真被改写，没有用例描述"应该拒"。补救：给它一条构造性用例（造一个"目标不存在、最近存在祖先形含 `%VAR%`"的拼写），或在票面明写"本支不可达、由 needle 兜"。
- **R-102-4**（判据射程，中）：`canonicalReaders` 的 `SkipDir` 名单排除 `testdata/`、`docs/` 以及整个 `tools/` 前缀 ⇒ 藏在这些目录里的第六条腿扫不到。补救：把排除表本身变成一条断言（"排除目录里不得出现 `risk.Resolve(`"），或至少把 `tools/` 纳入扫描。
- **R-102-5**（能力无承载者，中）：`RewrittenRoots()` 零生产读取者 ⇒ 票面/注释里的"operator-visible 审计面"未兑现。补救：在 `cmd/wisp/run.go` 的启动审计行里连 `Roots()/UnusableRoots()/RewrittenRoots()` 一起打（那是 `cmd/wisp` 的地界，需另票）。
- **R-102-6**（注释超前，轻）：`Result.Spelling` 的文档说"是审批提示里人看到的形"，无任何审批面读取。补救：改注释或接入 L2 提示。
- **R-102-7**（环境，非本票）：`go test ./cmd/wisp/` 本机 `0xc0000135`、`GOOS=linux go vet ./...` 整仓 rc=1 两条仪器坑仍在；票 98 的注入命令本代理未持有。

---

## 10. 总判

| AC | 判定 |
| --- | --- |
| AC#1 | **通过**（修前红 rc=1 / 修后绿 rc=0，两半均由本代理亲自复现，含 `=== RUN` 计数与 `FAIL-OPEN:` 字样） |
| AC#2 | **通过**（命中集合与处置表 1:1；逐腿读的是测试体；残留 R-102-2/R-102-3） |
| AC#3 | **通过**（M1 红 5 子、M2 另红 2 父，与实现方登记逐字对得上；两发都 build rc=0 + 落地行打印） |
| AC#4 | **通过**（四包 rc=0、FAIL=0、SKIP 4 次 = 2 名 × 2 轮且修前同样 SKIP；d22scan rc=0 台账 8 数逐字对得上，`#6/#8 frontend/` 未降） |
| AC#5 | **通过但有条件**（探针 A/B 双向攻过，判据非 vacuous、非"扫空=绿"；条件 = R-102-4 的排除表缺口） |
| 契约引句 | **逐字成立** ⇒ "本票是修实现、未触发 D22 停手"、"本票不需人工批准"的结论**成立** |
| `tools/paths.go` 争议 | **改动保留、自述与 diff 相符、形状 fail-closed、放行侧不比拒绝侧宽** |

# **总判：PASS WITH CONDITIONS**

五框全部勾得住，没有一条是照着一个并不存在的 API 写的（故本票没有"判据本身有毛病"格）。条件是**残留性质的**：R-102-3/R-102-4/R-102-5 三条应在票 103+ 各立一票，不阻塞本票 `-done`；R-102-1/R-102-2/R-102-6 是文档记账卫生，随下一笔追加一起清。
**本代理未行使 `-done` 命名权、未改 Status 归类、未改禁区任何文件。**

## 11. 没验完的（如实留断点）

1. `internal/memory/TestSubprocessCrashWriter` 的 SKIP 本代理只核对名字与理由串，**未复跑其子进程机制**（下一条命令：`cd /tmp/wisp-102-pre-ac2 && go test -count=2 -v -run 'TestSubprocessCrashWriter' ./internal/memory/`）。
2. 腿 4（祖先重解析）的不可达性是本代理**读码推定**，没构造实验。要证伪：在 `/tmp/wisp-102-instr-ac2` 里构造"目标不存在 + 最近存在祖先含 `%VAR%`"的拼写，看 `resolveTarget` 是否返回 `ErrRewrittenPath`。
3. `ban #7 internal/tools/=17` 的**分母变化**没对基线 `a66aadf^` 复算（只与本票参照比未降）。要补：`git archive 0117459 | tar -x -C /tmp/wisp-102-ban7-ac2` 后跑 `sh scripts/d22scan.sh` 对台账 8 数。
4. winsec `%w` 传播那半（`internal/winsec/resolve.go:113-119`）本代理**只从 AC#1 的输出串反推**（`winsec: refusing to seal … : risk: expansion moved …`），没读那段码。
