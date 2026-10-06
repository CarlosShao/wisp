# 111-ci1 — 推送后两发 CI run 的步级读数（只取数，不裁框）

> 本件是**只读取证腿** `111-ci1` 的交件。范围＝把票 111 里逐字要求"CI 侧真读数"的那几格
> （AC#1 第④列 / AC#3 前提今判 / AC#6 skipped 全名册 / 常红归因）喂到**步级**。
> **本腿不翻任何 `- [ ]` 框、不改票面、不裁 AC 的生死**；每张表末尾的〔建议位〕只是给编排者的理由，
> 判语一律留给编排者。

---

## §0 起手锚（三把尺原文）

| 尺 | 命令原文 | 读数 |
|---|---|---|
| HEAD | `git rev-parse --short HEAD` | `c98fc5db` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` |
| 时刻 | `date '+%m-%d %H:%M'` | `10-06 11:46` |
| 自证零 go 命令 | 见下 | 本腿全程未调用 `go test` / `go build` / `go vet` / `go run` 任何一枚；枚数与行数一律走 `git ls-files`／`wc`／`sed`，CI 侧读数一律走 `gh`。go 使用权整轮留给写腿 `232-r2`。 |

**本腿的尺原文清单（可复量）**
- `gh run list --branch dev -L 8`
- `gh run view 37405698188 --json jobs -q '...'`（步级 conclusion）
- `gh run view 37406757402 --json jobs -q '...'`
- 计数／名单：`git ls-files` 系列、`wc -l`、`sed -n`、`grep -n`
- ⛔ 无 `go`／无 `go test`／无 `go vet`

**run 号存活确认**（`gh run list --branch dev -L 8` 原文摘录，10-06 11:4x）
```
completed  failure  ## A634                                              ci  dev  push  37406757402  9m36s  2026-10-06T02:58:39Z
completed  failure  docs(ticket-269): 把已撤销就地标到 Status 行         ci  dev  push  37406422380  9m29s  2026-10-06T02:54:30Z
completed  failure  ## 5. 编排者翻勾记录（2026-10-06 10:3x +08…          ci  dev  push  37405698188  8m34s  2026-10-06T02:45:36Z
completed  failure  ci                                                 ci  dev  sched   37396530365  9m6s   2026-10-06T00:55:36Z
completed  success  slo-fresh                                 slo-fresh dev  sched   37392522160  15s    2026-10-06T00:09:51Z
```
- 主口径两发（任务给的）：`37405698188`（推送后新配置第一发）、`37406757402`（headSha `cc315261`，推送后第二发）。
- **推送后其实有第三发** `37406422380`（10-06 02:54:30Z，夹在两发之间）＝任务前文未列；本件主表仍按任务口径只用那两发，
  第三发只在它能让某一问的读数更硬时才补，补了会在对应小节具名。
- 推前旧配置对照发：`37396530365`（schedule，10-06 00:55:36Z）＝§5 作差的左端。

---

## §1 票面逐字（问 1）

尺原文：`sed -n '20,35p;70,90p' .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md`
**行号复核结论：没漂**。本腿另用 `grep -n '^- \[ \]\|^- \[x\]\|^Status:' .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md`
现量 ⇒ 十枚框分别在 `:22 / :24 / :27 / :29 / :31 / :71 / :75 / :77 / :81 / :86`，
与编排者给的 `20,35p`＋`70,90p` 两段**完全覆盖**；票面总行数 `wc -l` = **306**，**零枚 `- [x]`**（全 10 枚仍是 `- [ ]`）。

```
22| - [ ] **AC#1** 复算并出一张**全仓对账表**：每个包 ×（有无测试文件 / 在不在 CI 某一步 / 那一步真给过结论的 run id + step 号）。
23|       表格必须能自证完整（`go list ./...` 的 33 行都在，不许只列零覆盖那几个）。
24| - [ ] **AC#2** 逐包接入，**先易后难**，每包一次可核对的步级读数。
25|       ⚠ **加严可以直接做**；**不许**为了让某包变绿而放宽它的断言、调它的阈值、或给它加 `//go:build`/`t.Skip` 挡掉。
26|       接入第一天就红 ⇒ **那是发现**：红名逐条登记进本票面并**开票**，不许撤步骤。
27| - [ ] **AC#3** `session`/`watchdog` 这种"空分母"要**响亮**：scope 校验加上
28|       "**声明在范围内但该平台没有任何测试文件 ⇒ 直接失败**"的守卫（与票 93 的"条目腐坏即红"同族），并人为抽掉一个包证明它会红。
29| - [ ] **AC#4** `cmd/wisp` 那一格：给出它在 CI 上**能不能跑**的实测结论（能 ⇒ 接入；不能 ⇒ 写清缺什么、归票 98 还是新票），
30|       **不许默默留在零覆盖列**。
31| - [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
32|       ⚠ 新步若排在"会失败的步骤"之后 ⇒ 必须放前面或 `if: always()`（本仓实测过这道门因此从未执行）。
```

```
71| - [ ] **AC#6** 让**新增的门不再吃掉后面的步骤**：把 step4 之后的每一步都还能跑
72|       （`if: always()` 或把新步挪到该 job 最后，二选一并说明为什么）。
73|       ⚠ 加 `always()` 属**加严**可以直接做；**不许**反过来把新步删掉或挪到 `continue-on-error`（那等于把门拆了）。
74|       判据：**同一枚 run 里 step4 与 step5–8 同时有结论**（允许 step4 红），并给一次这样的**步级**读数（run id + 各步 conclusion）。
75| - [ ] **AC#7** `R-110-4`：票 110 承诺过"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿（step7）"，
76|       验收复算发现**至今 0 次** ⇒ 要么落地并给 step7 读数，要么在票 110/111 面把它**当众改口径**（不许留在原地当已做）。
77| - [ ] **AC#8** `R-110-3`：包匹配式缺前缀锚定（`"winsec"` 宽松串曾让我把 0 命中说成 18 次）
78|       ⇒ 匹配式要能区分"**被测包**"与"日志里出现过这个词"，并用一次阳性自证（种一个只在字符串里出现的包名 ⇒ 不许计入分母）。
81| - [ ] **AC#9（编排者 21:0x 追加，来源=run `35599458439` 的真实读数）** **`internal/winsec` 的 POSIX 半边今天零覆盖**：
82|       `test-core`（ubuntu 腿）step7 的逐字 scope 是那 16 个包、**不含** `./internal/winsec/`，全日志里 `internal/winsec` 出现 **0 次**
83|       ⇒ 票 113 刚交的链接腿（`winsec_other.go`）**没有任何 CI 回归保护**，只有编排者本机 Docker 跑过。
84|       判据：ubuntu 腿里出现一步真跑 winsec 的 `!windows` 半边，并给出**该步的 run id + job id + step 号 + 结论**。
85|       ⚠ 不许用"本机 Docker 跑过"替代；也不许把它接成"只编译不执行"（`GOOS=linux go vet` 那一类）就算数。
86| - [ ] **AC#10（编排者 21:0x 追加）** `test-core` step7 现在报 **PASS=578 FAIL=0 SKIP=1**，而那枚 skip
87|       （`TestWorkspaceSwitchRefusesAJunctionToOutside`，`paths_workspace_test.go:198`，理由是"C26 reparse 检测是 Windows-only"）
88|       **不在任何台账里** ⇒ 与票 93 同族（**步不再把 SKIP 记成 ok**），但这一枚要的是：**未记账的 skip 必须响亮**——
89|       要么进"已知双平台跳过"清单并写明归谁，要么在 POSIX 上给出等价判据。**不许**为消掉数字而 `Skip` 掉它。
```

**一处要具名的口径差（⛔ 不是本腿改票面）**：任务前文把 AC#6 的生死判据引成
「只要还有一步是 `skipped` 就判 FAIL，不接受'通过附条件'」并说它在"约 `:71-74`"。
**票面 `:71-74` 那四行逐字如上，里面没有这句话。**
这句话真身有两处、都在 Progress log 里而非 AC 段里：
- `:255`（next= 1 的判据本体，点名 test-windows 的五枚非 post 步）：「只要还有一步是 skipped ⇒ **没修好，AC#6 当场判 FAIL**（不许写"通过附条件"）」
- `:291`：「是本票唯一的生死判据：**只要还有一步是 skipped，就判 AC#6 FAIL**，不接受"通过附条件"」

⇒ 两处的**射程不一样**（一处限定五枚、一处未限定），这正是 §4.3 把两读法摆开的依据。本腿按硬边界**一枚框、一字票面都没动**，只把行号与原文差异登记在这里。

**AC 段之外、与问 1 同一段落里的一条过期陈述**（同号具名，不改）：票头 `:1` 与 `:23` 的 **33** 这个分母，
今天三发 census 步逐字打的是 `go list ./... = 35 packages`（见 §2.1、§6.8）。

---

## §2 AC#1 第④列名册（问 2）

### 2.0 尺与分母（先自证完整）

- 包全集（第①～③列的前两块）：`git ls-files '*.go' | grep -v '^\.scratch/'` 按目录聚合现量（枚数＝跟踪的 `_test.go` / 非测试 `.go` 各几枚），
  core 名单来自 `scripts/portable-tests.sh`（`core_pin` `:164-192`、`scope(core)` `:234-247`），
  windows 腿名单同文件（`win_pin` `:193-204`、`scope(windows)` `:248-255`），`cli` `:256-262`、`winsec` `:263-270`，`tiers='core windows cli winsec census'` `:220`。
- 第④列**只从这两发 run 的步级 log 里取**：`gh run view --job <id> --log`（全量，⛔ 不用 `--log-failed`），落盘 9 枚 gz 后逐枚 grep。
- **分母今值＝35，不是票面那 33**（三发 census 步逐字 `go list ./... = 35 packages`；差集口径见 §6.8）。
  下表**35 行全在**，满足票面 `:23` 那句"不许只列零覆盖那几个"。
- ⚠ `tools/d22scan`／`tools/mockllm`／`scripts/spike` 是**另三枚 `go.mod`**（`git ls-files` 现量 4 枚：`go.mod`、`tools/d22scan/go.mod`、`tools/mockllm/go.mod`、`scripts/spike/go.mod`），
  不在主 module 的 `go list` 里 ⇒ 不在 35 行内；它们**有** CI 读数（`lint` step4 / step13），单列在 §2.3。

### 2.1 census 步打印的四数（逐作业，第④列的"表头"）

| run | job | 步 | census totals 原文（四个数） | 结论 |
|---|---|---|---|---|
| `37405698188` | `112082660385`（test-windows） | **step8** | `census GOOS=windows - go list ./... = 35 packages` ／ `census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` | success |
| `37406757402` | `112085927937`（test-windows） | **step8** | 同上，**四数逐字相同** | success |
| `37406422380`（第三发，参照） | `112084901765` | **step8** | 同上，四数亦相同 | success |
| 推前 `37396530365` | `112053739011`（test-windows） | — | **无 census 步**：该发 step8 是 `Portable windows tests`，step9 是 `PathResolver junction`；`gh run view 37396530365 … test("census")` 返回**空** | 不存在 |

（⛔ 这三行不是"core 那发 1518/1028/4/1 一形"。那四枚是 `runtests.sh` 的 **RUN/PASS/FAIL/SKIP 四数**，
不是 census 的 totals；两套四数分别在 §2.4 与本节。census 的四个数是 `packages/with-zero-compiled-tests/claimed-by-no-scope/unclaimed-with-tests`。）

### 2.2 三十五行全名册

列义：①`git ls-files` 跟踪的 `_test.go` 枚数 ②census 在 windows 上打印的 `TESTS(t/x)`（＝真编进该平台 test 二进制的那几枚，`!windows` 的会少）
③在不在 CI 某一步（哪个 tier→哪个作业第几步）④那一步**真给过结论**的 run＋job＋步＋该枚包自己的 own-line 判语。

| # | 包（`github.com/CarlosShao/wisp/` 之后） | ①`_test.go` | ②t/x | ③在册（tier→步） | ④步级真读数 |
|---|---|---|---|---|---|
| 1 | `cmd/balldebug` | 0 | `0/0` | **NO-SCOPE**（`<-NO-TESTS`） | 无任何步 ⇒ 无读数 |
| 2 | `cmd/llmrecord` | 1 | `1/0` | core＋windows | core：`37406757402`/`112085927688`/step7 `ok (own line)`；windows：`112085927937`/step9 `ok` （run1 `112082660423`/step7、`112082660385`/step9 同形） |
| 3 | `cmd/wisp` | 65 | `58/0` | cli→`test-windows` step7 | 三发都**给了结论**：run1 `112082660385`/step7 `FAIL (own line)`（`RUN=335 PASS=231 FAIL=6 SKIP=1`）；run2 `112085927937`/step7 `FAIL`（`335/229/8/1`）；run3 `112084901765`/step7 `FAIL`（`335/231/6/1`）。推前 `112053739011`/step7 亦 FAIL（`323/220/6/1`） |
| 4 | `frontend` | 0 | `0/0` | NO-SCOPE | 无读数 |
| 5 | `internal/agent` | 19 | `19/0` | core | run1 `112082660423`/step7 `ok`；run2 `112085927688`/step7 `ok`（该步整体 failure，但本枚 own-line 是 ok） |
| 6 | `internal/agent/approval` | 19 | `8/11` | core | 同上两发 step7 `ok (own line)` |
| 7 | `internal/agent/scheduler` | 0 | `0/0` | NO-SCOPE | 无读数（`doc.go` 型空壳，`portable-tests.sh:107-111` 具名它为何不能进 core） |
| 8 | `internal/audio` | 5 | `5/0` | core | 两发 step7 `ok` |
| 9 | `internal/ball` | 16 | `11/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 10 | `internal/buildinfo` | 1 | `1/0` | core | 两发 step7 `ok` |
| 11 | `internal/config` | 18 | `16/1` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 12 | `internal/llm` | 10 | `2/8` | core | 两发 step7 `ok` |
| 13 | `internal/llm/adaptertest` | 1 | `1/0` | core | 两发 step7 `ok` |
| 14 | `internal/llm/anthropic` | 3 | `3/0` | core | 两发 step7 `ok` |
| 15 | `internal/llm/golden` | 1 | `1/0` | core | 两发 step7 `ok` |
| 16 | `internal/llm/openaichat` | 4 | `4/0` | core | 两发 step7 `ok` |
| 17 | `internal/llm/openairesponses` | 2 | `2/0` | core | 两发 step7 `ok` |
| 18 | `internal/memory` | 11 | `11/0` | core | run1 `112082660423`/step7 **`ok`**；run2 `112085927688`/step7 **`FAIL (own line)`**（`--- FAIL: TestConcurrentWritersReaders`，见 §6.6） |
| 19 | `internal/models` | 17 | `16/0` | core | 两发 step7 `ok` |
| 20 | `internal/observe` | 15 | `15/0` | core | 两发 step7 `ok` |
| 21 | `internal/panel` | 20 | `20/0` | core | 两发 step7 **`FAIL (own line)`**（4 枚 C21/theme 契约用例，见 §5.3） |
| 22 | `internal/perm` | 3 | `3/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 23 | `internal/plugin` | 1 | `1/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 24 | `internal/proc` | 11 | `11/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 25 | `internal/projctx` | 1 | `0/1` | core＋windows（**推送才在册**） | core 两发 step7 `ok (own line)`＋`ok …projctx 0.006s`；windows 两发 step9 `ok`。推前 `37396530365` **无此枚读数**（pre 的 own-line 名册 10 枚里没有 projctx/session，见 §5.4） |
| 26 | `internal/risk` | 22 | `20/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 **`FAIL (own line)`**（12 枚用例，推前也红，见 §5.5） |
| 27 | `internal/secret` | 4 | `4/0` | core＋windows | core 两发 step7 `ok`；windows 两发 step9 `ok` |
| 28 | `internal/session` | **2** | `2/0` | core＋windows（**推送才在册**） | core：run1 `112082660423`/step7 `ok (own line)`＋`ok …session 0.127s`；run2 `112085927688`/step7 `ok`＋`0.098s`。windows：`112082660385`/`112085927937` step9 均 `ok (own line)`。推前**零枚**（§5.4）⇒ 这是 AC#1 第④列今天第一次有值的那格（口径分歧见 §6.5） |
| 29 | `internal/speech` | 0 | `0/0` | NO-SCOPE | 无读数 |
| 30 | `internal/statemachine` | 2 | `2/0` | core | 两发 step7 `ok` |
| 31 | `internal/streamkey` | 0 | `0/0` | NO-SCOPE | 无读数 |
| 32 | `internal/tools` | 49 | `46/3` | core | 两发 step7 `ok` |
| 33 | `internal/watchdog` | **0** | `0/0` | **NO-SCOPE**（`<-NO-TESTS`） | 无读数；跟踪件只有 `internal/watchdog/doc.go`（`git ls-files internal/watchdog` = 1 枚）。§3 的主角 |
| 34 | `internal/winsec` | 27 | `12/8` | core＋winsec tier | winsec 独立步：三发 `test-windows` **step4 success**（`RUN=101 PASS=58 FAIL=0 SKIP=0`，own-line `ok`）；ubuntu 半边 core step7 `ok (own line)`＋`ok …winsec 0.025s` 级读数（run1/run2 都有） |
| 35 | `tools/signmodels` | 0 | `0/0` | NO-SCOPE | 无读数 |

**自证完整**：35 行＝census 步打印的 35 枚（尺：`gzip -dc job-112085927937.log.gz | grep -cE 'portable-tests.sh: +github'` = **35**；run1 与 run3 各测一遍亦 35，三发名册逐字节相同——
尺原文 `diff /tmp/job-112082660385.census /tmp/job-112085927937.census` → 无输出）。
**②列与①列的差集只有六枚**（尺原文＝`awk` 对 `/tmp/gocounts.txt` 与 census 的 `t+x` 求差，逐枚可对上账）：

| 包 | ①跟踪 `_test.go` | ②windows 上编译进去 `t/x` | 差的账（HEAD 尺逐枚数出来的 tag） |
|---|---|---|---|
| `cmd/wisp` | 65 | `58/0` | 1 枚 `!windows` ＋ 6 枚 `windows && winlive`（CI 不设 `winlive`）⇒ 65−7=58 |
| `internal/ball` | 16 | `11/0` | 5 枚 `windows && winlive` ⇒ 16−5=11 |
| `internal/config` | 18 | `16/1` | 1 枚 `!windows` ⇒ 18−1=17（=16+1） |
| `internal/models` | 17 | `16/0` | 1 枚 `!windows` ⇒ 16 |
| `internal/risk` | 22 | `20/0` | 2 枚 `!windows` ⇒ 20 |
| `internal/winsec` | 27 | `12/8` | 7 枚 `!windows` ⇒ 20（=12+8）；tag 串里"含 windows"的 24 枚有 7 枚其实是 `!windows`（子串陷阱，正是 AC#8 那枚病，本腿按首行整条 tag 数的） |

其余 29 枚包 ①＝②（`internal/agent/approval 19↔8/11`、`internal/llm 10↔2/8` 是**内部拆 xtest**、总数不变，⛔ 不是 build-tag 缺口）。
⇒ 票面 AC#9 关心的"`internal/winsec` 的 POSIX 半边"**在 windows census 上必然显示为 20/27**，
那 7 枚 `!windows` 只能由 ubuntu 腿证明——本件 §5.6 给的就是那一步的读数。
尺原文：`git ls-tree -r --name-only HEAD -- <pkg> | grep '_test.go$'` ＋逐枚 `git show HEAD:<file> | grep -m1 '^//go:build'`（**一律按 HEAD，不碰工作树**）。
本腿**没跑 `go list`**（零 go 命令），所以②列只转述 census 自己打印的数，①列与②列的账靠上面这把 tag 尺闭合。

### 2.3 35 行之外、但有 CI 读数的两个 module（防"20/33"式漏计）

| 包 | 在不在 35 行 | CI 步 | 今天读数 |
|---|---|---|---|
| `tools/d22scan`（`*_test.go` 2 枚） | **不在**（独立 `go.mod`） | `lint` step4 `D22 scanner positive control` | run1 `112082660317`／run2 `112085927883` 均 **success**，`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77` |
| `tools/mockllm`（`*_test.go` 3 枚） | 不在 | `lint` step13 `mockllm module vet`（`go vet ./...`） | 三发 **success**（但那是 vet，不是 test；**没有任何一步跑 mockllm 的测试**） |
| `scripts/spike`（11 枚 `.go`） | 不在 | 无 | 无读数 |

### 2.4 各步的四数（第④列的"这一步有没有真给结论"的分母）

| run | job | 步 | 四数原文（`portable-tests.sh: four numbers …`） | 步结论 |
|---|---|---|---|---|
| `37405698188` | `112082660423` | test-core step7 | `=== RUN=1587 --- PASS=1071 --- FAIL=4 --- SKIP=1` | failure |
| `37406757402` | `112085927688` | test-core step7 | `=== RUN=1587 --- PASS=1070 --- FAIL=5 --- SKIP=1` | failure |
| `37406422380` | `112084901657` | test-core step7 | `=== RUN=1587 --- PASS=1071 --- FAIL=4 --- SKIP=1` | failure |
| 推前 `37396530365` | `112053739109` | test-core step7 | `=== RUN=1518 --- PASS=1028 --- FAIL=4 --- SKIP=1` | failure |
| `37405698188` | `112082660385` | test-windows step4 | `RUN=101 PASS=58 FAIL=0 SKIP=0` | success |
| `37405698188` | `112082660385` | test-windows step7 | `RUN=335 PASS=231 FAIL=6 SKIP=1` | failure |
| `37405698188` | `112082660385` | test-windows step8 | **census**（四数见 §2.1） | success |
| `37405698188` | `112082660385` | test-windows step9 | `RUN=577 PASS=386 FAIL=12 SKIP=1` | failure |
| `37406757402` | `112085927937` | test-windows step4/7/8/9/10 | `101/58/0/0` ／ `335/229/8/1` ／ census ／ `577/386/12/1` ／ junction `PASS=1 FAIL=0 SKIP=0` | success／failure／**success**／failure／success |
| `37406422380` | `112084901765` | test-windows step4/7/8/9/10 | `101/58/0/0` ／ `335/231/6/1` ／ census ／ `577/386/12/1` | 同上形 |
| 推前 `37396530365` | `112053739011` | test-windows step4/7/8/9 | `101/58/0/0` ／ `323/220/6/1` ／ `536/362/12/1` ／ junction success | **无 step8 census**（步号整体前移一格） |

〔建议位——仍归编排者裁〕第④列今天**第一次能整列填上 CI 侧读数**：35 枚里 28 枚有真步级结论、
7 枚 NO-SCOPE（`cmd/balldebug`／`frontend`／`internal/agent/scheduler`／`internal/speech`／`internal/streamkey`／`internal/watchdog`／`tools/signmodels`，全部 `0/0` ⇒ 零测试文件，不是"有测试没在册"）。
⇒ "CI 只测 33 个包里的 20 个"这句**今天数字已翻**（28/35 有结论，且 census 的 `unclaimed-with-tests=0` 说明**再没有"带测试却零覆盖"那一形**）；
但分母 33→35 属口径变更（§6.8），本腿不据此翻框。

---

## §3 GUARD D 与空分母（问 3，AC#3 前提今判）

### 3.1 两枚主角今天的跟踪测试枚数（尺原文，全部按 HEAD／不碰工作树）

```
git ls-files 'internal/session/*_test.go'
  internal/session/grants_test.go
  internal/session/ticket224_pattern_dialect_test.go            ⇒ 2 枚
git ls-files 'internal/watchdog/*_test.go'   ⇒ 空（0 枚）
git ls-files 'internal/watchdog/'            ⇒ internal/watchdog/doc.go，仅此 1 枚
```
票面 `:27-28`（AC#3）的前提今天**两半不同**：
- `session` 的"在清单里却一个测试文件都没有"这一支**已翻**（2 枚）；
  且**推前那一发也已翻**——`git ls-tree -r --name-only c6cf66e6 -- internal/session | grep -c '_test.go'` = **2**
  ⇒ 测试文件不是推送才出现的，推送带来的是**在册**（见 3.2）。
- `watchdog` 那一支**今天仍成立**：0 枚、只有 `doc.go`；且它**不在任何名册里**（下条），
  所以它现在属"零测试＋不在册"＝GUARD A 管不到、GUARD D 也不算洞（`0/0` 走 `case 0/0) ;;` 那一支）。

### 3.2 在不在 scope 名册／GUARD D 的名册里（HEAD 尺，行号现量）

| 包 | `core_pin` | `win_pin` | `scope(core)` | `scope(windows)` | GUARD D 眼里的形状 |
|---|---|---|---|---|---|
| `internal/session` | **在**，`:188` | **在**，`:203` | **在**，`:244` | **在**，`:252` | 已在册 ⇒ 不进 `unclaimed` 计数 |
| `internal/watchdog` | 不在 | 不在 | 不在 | 不在 | **NO-SCOPE ＋ `0/0`** ⇒ census 打 ` NO-SCOPE <-NO-TESTS`，计入 `claimed-by-no-scope=7` 但**不计** `unclaimed-with-tests` |

推前同尺对照（`git show c6cf66e6:scripts/portable-tests.sh`）：`grep -c "GUARD D"` = **0**、
`grep -c "guardd\|UNCLAIMED"` = **0**、`git show c6cf66e6:.github/workflows/ci.yml | grep -c "scope=census"` = **0**
⇒ **GUARD D 这支守卫与它的调用步在推前那个 head 上根本不存在**（`tiers=`/`census)` 分支当时已有，但没有 GUARD D 的判红）。
`internal/session` 在 `c6cf66e6` 的 `core_pin`/`win_pin` 里也**查无**（`git show c6cf66e6:scripts/portable-tests.sh | grep -n "wisp/internal/session\|wisp/internal/projctx"` = 无输出）
⇒ "测试文件早就有、名册里一直没有"这件事，是**这次推送才补上的**，且和日志读数对得上（§5.4）。

### 3.3 "有测试却没在册"这一形，census 步在 CI 上**响没响**

尺原文（三发 census 步各测一遍）：
```
for j in job-112082660385 job-112085927937 run3-job-112084901765; do
  gzip -dc $j.log.gz | grep -c 'UNCLAIMED-HAS-TESTS'   # 三发都 = 0
  gzip -dc $j.log.gz | grep -c 'GUARD D -'             # 三发都 = 0
  gzip -dc $j.log.gz | grep -c '<-NO-TESTS'            # 三发都 = 7
  gzip -dc $j.log.gz | grep -cE 'portable-tests.sh: +github'   # 三发都 = 35
done
```
⇒ **三发都没响**（`unclaimed-with-tests=0`，GUARD D 的判红句一次都没打印）。

**success 的成因判定：是"真没漏"，不是"这一支不产红"。** 三条步级证据（⛔ 全部靠日志与行号，不靠注释）：
1. **七个 NO-SCOPE 行逐枚都带 `0/0`**：`cmd/balldebug 0/0`／`frontend 0/0`／`internal/agent/scheduler 0/0`／
   `internal/speech 0/0`／`internal/streamkey 0/0`／`internal/watchdog 0/0`／`tools/signmodels 0/0`，
   且行尾统一挂 `<-NO-TESTS`。这正对应 `portable-tests.sh:392-393` 的 `case $counts in` / `0/0) ;;` ⇒
   这七枚**根本不该**进 `guardd` 计数。35 枚里其余 28 枚的 `t/x` 至少一边非零，而它们**全部**有 tier 认领（`where` 非空 ⇒ 走不到 `:382` 那个 `[ -z "$where" ]`）。
2. **退出码路径逐行号可追**：`:360 guardd=0` → `:396 guardd=$((guardd + 1))`（唯一的自增处，在无认领＋非 `0/0` 那一支里）
   → `:406` 打印 totals（今天打出的是 `unclaimed-with-tests=0`）→ `:407 if [ "$guardd" -ne 0 ]` → `:422 exit 1`／`:424 exit 0`。
   今天的 totals 就是 `0` ⇒ 走 `:424`。**这条 `exit 1` 是真在的**，同一步里还有另外三处 `exit 1`（`:305` 读不到自身字节、`:324` tiers 与分支不等集、`:350` `go list` 非零），今天也都没走。
3. **default-deny 的反面没有漏**：`:392-399` 那支把"计数读不出来"（`?/?` 或空）也算洞。三发 35 行里**没有一行是 `?/?`**（逐枚都是 `数字/数字`）
   ⇒ 今天不是"读不出所以没计"，是"读得出且确实为 `0/0`"。

**但必须摆的另一半（这才是"响亮"与"能红"的差别）**：CI 上**至今没有任何一发**给过 GUARD D 真的红一次的读数，
而票面 AC#3 逐字要的第二半"并人为抽掉一个包证明它会红"的载体 `scripts/portable-tests-selftest.sh`
在 `.github/workflows/ci.yml` 里被引用 **0 次**（`grep -c portable-tests-selftest .github/workflows/ci.yml` = 0；`slo-fresh.yml` 亦 0）。
⇒ **本腿给的结论是**："success＝真没漏"这一条有步级证据（上面 1/2/3）；
"这一支能不能产红"这一条**今天仍无 CI 证据**，属 §6.3，⛔ 不许拿"逻辑上行号在那儿"当"CI 上验过"。
`watchdog` 的"响亮"今天也只是**打印**（` NO-SCOPE <-NO-TESTS` 那一行，与有覆盖行同音量），
不是 AC#3 那句"直接失败"——因为它压根不在 scope 里，GUARD A（`:570-583`，`exit 1`）碰不到它。
GUARD A 今天同样没响，且不是只测一发：**八枚作业全量日志**（四发 × core/windows 两条腿：
`pre-112053739109`／`pre-112053739011`／`job-112082660423`／`job-112085927688`／`job-112082660385`／`job-112085927937`／
`run3-job-112084901657`／`run3-job-112084901765`）里 `grep -c 'GUARD A -'` **逐枚 = 0**，`grep -c 'GUARD C -'` 亦逐枚 = 0。

---

## §4 skipped 全名册与两分（问 4，自我对抗节 A）

尺原文：`gh run view <run> --json jobs -q '.jobs[] | .name as $n | .steps[] | select(.conclusion=="skipped") | "\($n) :: step\(.number) [\(.name)]"'`
（四发各跑一遍，输出逐字贴在下面；⛔ 不是 `--log-failed`，那把尺会少算。）

### 4.1 全名册（逐枚，含推前与第三发做参照）

| run | 作业 | 步号 | 步名（原文） | 为什么被跳过（能不能从日志判） |
|---|---|---|---|---|
| `37405698188` | `lint` | 9 | `gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)` | step8 `gofmt (gofumpt)` failure 吃掉；该步在 `ci.yml:184`，**无 `if:`** ⇒ 默认 `success()` |
| `37405698188` | `lint` | 10 | `go vet (module)` | 同上（`ci.yml:206`，无 `if:`） |
| `37405698188` | `lint` | 11 | `go vet (tools/d22scan module)` | 同上（`ci.yml:209`，无 `if:`） |
| `37405698188` | `lint` | 25 | `Post Run actions/setup-go@v5` | action 自带的 post；`setup-go@v5.0.0/action.yml:30-31` = `post: dist/cache-save/index.js` / **`post-if: success()`**，作业已红 ⇒ skipped |
| `37405698188` | `test-windows` | 18 | `Post Cache third_party (deps.toml key)` | `actions/cache@v4.0.2/action.yml:39-40` = **`post-if: "success() \|\| github.event.inputs.save-always"`**；`ci.yml` 里 `save-always` 出现 **0 次**（`grep -c save-always .github/workflows/ci.yml` = 0），作业红 ⇒ skipped。⚠ 该步自身的 `if: ${{ !cancelled() }}`（`ci.yml:474`）只管**恢复半（step5）**，管不到 post 半 |
| `37405698188` | `test-windows` | 19 | `Post Run actions/setup-go@v5` | 同 `setup-go` `post-if: success()` |
| `37405698188` | `test-core` | 15 | `Post Run actions/setup-go@v5` | 同上 |
| `37406757402` | `test-core` | 15 | `Post Run actions/setup-go@v5` | 同上（枚枚同因，只是作业顺序变了） |
| `37406757402` | `lint` | 9 | `gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)` | step8 failure 吃掉（同一枚红句，见 §5.2） |
| `37406757402` | `lint` | 10 | `go vet (module)` | 同上 |
| `37406757402` | `lint` | 11 | `go vet (tools/d22scan module)` | 同上 |
| `37406757402` | `lint` | 25 | `Post Run actions/setup-go@v5` | `post-if: success()` |
| `37406757402` | `test-windows` | 18 | `Post Cache third_party (deps.toml key)` | `post-if: success() \|\| save-always` |
| `37406757402` | `test-windows` | 19 | `Post Run actions/setup-go@v5` | `post-if: success()` |

**两发主口径 run 的 skipped 枚数＝7 : 7，且名册逐枚同构**（只是 `test-core` 那枚在两发里都排第 15）。
作差用的参照两发（同尺现量，⛔ 二手）：
- 第三发推送后 `37406422380`：skipped 名册与 `37406757402` **逐枚相同**（7 枚：test-core 15 / lint 9,10,11,25 / test-windows 18,19）。
- 推前 `37396530365`：也是 7 枚，但 test-windows 那两枚编号是 **16、17**（不是 18、19）——
  因为那一发**没有 census 步**（见 §3.4），插一步就把 post 编号往后推一格。
  ⇒ 结论性事实：**skipped 的形状不是推送引入的，推前推后同构**；差别只在步号。

**反证尺（证明"不是整个 post 阶段被跳过"）**：每一发的每一个作业里
`Post Run actions/checkout@v4` **都是 success**（`checkout@v4/action.yml:104` 只声明 `post:`、**不声明 `post-if`** ⇒ 走 runner 默认，红作业也跑）。
所以 skipped 是**逐 action 的 `post-if` 属性**，不是"作业红 ⇒ post 全灭"。这条差别只能由 `Post Run actions/checkout@v4` 承载，本腿已把它从全名册里单列出来。

### 4.2 两分：谁承载读数，谁不承载

**① 被前一步失败吃掉的实质测试步（这一类承载读数，且读数永久采不到）**
- `lint` step9 / step10 / step11 三枚。
- 它们各自承载什么：step9 承载 `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`（161 AC#7 form A 的"受跟踪集当分母"那道判据）；
  step10 承载 `go vet ./...`（主模块，`ci.yml:207`）；step11 承载 `go vet ./...`（`tools/d22scan` 模块，`ci.yml:210`）。
- **一步红会吃掉后续步 ⇒ 被 skip 的步没有日志**：这三枚在 `37405698188`／`37406757402`／`37406422380`／`37396530365`
  四发里**全部**没有日志行（`gh run view --log` 拿不到它们的任何输出）。
  ⇒ **⛔ 不许把"没日志"读成"没红"**：`go vet` 这两枚今天是红是绿**未知**，不是"应该绿"。
  这三枚的形状正是票 111 AC#6 想治的那枚病（`ci.yml` 的注释 `:180` 附近也这么写），
  同作业里带 `if: ${{ !cancelled() }}` 的 step12 `staticcheck`（`ci.yml:263`）与 step13 `mockllm module vet`（`:300`）
  **照样给了结论**（step12 failure／step13 success），这道对照就是"加了 `!cancelled()` 的步能穿透前一步的红"的现成读数。

**② setup 的 post 清理步（这一类不承载任何读数）**
- `Post Cache third_party (deps.toml key)`、`Post Run actions/setup-go@v5`、`Post Run actions/checkout@v4`。
- 它们是谁写的：**不是 `ci.yml` 写的**，是 `uses:` 动作自带的 post 钩子由 runner 注入的（所以名册里它们的步名一律是 `Post …` 前缀、
  编号和主步之间有空洞：test-windows 主步到 step10，post 从 step18 起）。
- 它们承载什么：`setup-go` 的 post 写 Go 模块缓存（`dist/cache-save/index.js`）；`cache` 的 post 保存/落盘 `third_party` 缓存；
  `checkout` 的 post 删临时 git config。**四数（RUN/PASS/FAIL/SKIP）、own-line 名册、census 表——一枚都不在这些步的输出里。**
- 现量证据：`37406757402`/`test-windows` 全日志里 `Post Cache`/`Post Run` 相关内容只有 `Post Run actions/checkout@v4` 打了
  git 版本与 config 那几行（`:5596-5609`），另两枚 skipped ⇒ **零字节输出**，没有任何测试结论可丢。

### 4.3 AC#6 的两读法（本腿不裁，只把两面摆平）

票面那把生死尺的**原文有两个射程**（逐字抄自票面，行号现量）：
- **窄读法（`:250-255`，next= 1 的判据本体）**：点名五枚**非 post 的 test-windows 实质步**——
  `Cache third_party` / `cgo build smoke` / `cmd/wisp CLI tests` / `Portable windows tests` / `PathResolver junction placeholder`
  "各自必须有 conclusion（success 或 failure，不能是 skipped）；只要还有一步是 skipped ⇒ 没修好，AC#6 当场判 FAIL"。
  两发主口径 run 的读数：`step5 success／step6 success／step7 failure／step9 failure／step10 success` ⇒ **这五枚无一 skipped**（第三发同形）。
  ⚠ 同一枚 run 里 census（step8）也 success ⇒ `step4 与 step5–8 同时有结论` 这半句（`:74`）今天**第一次**成立。
- **宽读法（`:291`）**：「**只要还有一步是 skipped，就判 AC#6 FAIL**，不接受"通过附条件"」——没限定步集合。
  按这半句逐字执行 ⇒ 两发各有 **7 枚 skipped**（上表），**包含 lint 的三枚实质步**，无论 post 类算不算，
  `lint` step9/10/11 都是被吃掉的**测试步** ⇒ 宽读法下**光凭 lint 那一作业就翻不了**。

〔建议位——留给编排者裁，本腿不自翻〕
- 事实层面可以给死的三点：①`lint` 那道病**推前推后一模一样**（四发同形），本票 AC#6 若要的是"新门不再吃掉后面的步"，
  那 `lint` 作业里的 step9/10/11 至今仍是"被吃掉且无日志"，**没有任何一处出现它们今天的结论**；
  ②`test-windows` 的五枚点名步**今天首次全数给结论**（且 census success 是历史第一次）；
  ③post 类三枚**不承载读数**，把它们算进 AC#6 的"还有一步是 skipped"里，等于把 runner 的 `post-if` 属性记在本票头上。
- ⇒ 建议：AC#6 的"翻/不翻"取决于裁的是 `:250-255` 的五枚点名集还是 `:291` 的全集；
  若裁全集，**必须同时说明 `lint` step9/10/11 的处置**（这三枚归 `ci.yml` 面，票面 `:306` 已具名"归 ci.yml 面，与票 85 地界同处"）。
  本腿一枚框都不动。

---

## §5 红名作差（问 5，推前推后同口径）

同口径做法：四发都用 `gh run view --job <id> --log`（**全量日志，⛔ 不用 `--log-failed`**）＋同一把尺
`grep -oE "\-\-\- FAIL: [A-Za-z0-9_]+" | sort -u`；
推前那发 `37396530365`（head `c6cf66e6`）与推送后三发（`ec84cff1`／`b948bcb8`／`cc315261`）。
⚠ 推前那发的四枚 job 日志第二列同样全是 `UNKNOWN STEP`（§6.7 同病）⇒ 推前的**步号归属**由 `ci.yml` 当时的步序＋包名归属推，
步级 conclusion 则来自 `gh run view --json jobs`（API，与日志无关，可靠）。

### 5.1 `test-core` step7（Portable package tests，ubuntu 腿）

| run | job | 步 | 四数 | 红名（逐字） |
|---|---|---|---|---|
| 推前 `37396530365` | `112053739109` | step7 failure | `RUN=1518 PASS=1028 FAIL=4 SKIP=1` | `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`（own-line FAIL 只有 1 枚＝`internal/panel`） |
| `37405698188` | `112082660423` | step7 failure | `RUN=1587 PASS=1071 FAIL=4 SKIP=1` | **同上 4 枚，逐名相同**（own-line FAIL＝`internal/panel` 1 枚） |
| `37406757402` | `112085927688` | step7 failure | `RUN=1587 PASS=1070 FAIL=5 SKIP=1` | 上面 4 枚 **＋** `TestConcurrentWritersReaders`（own-line FAIL＝`internal/memory`＋`internal/panel` 2 枚） |
| `37406422380`（参照） | `112084901657` | step7 failure | `RUN=1587 PASS=1071 FAIL=4 SKIP=1` | 回到 4 枚（`TestConcurrentWritersReaders` 这发 **PASS**） |

**作差结论**：那 4 枚 `internal/panel` 的 C21/theme 契约用例**推前就红**，一字未变；
今天**唯一新增的红名＝`TestConcurrentWritersReaders`**，且只在 `37406757402` 出现（三发 Go 源码逐字节相同，见 §6.6 ⇒ 写"同码不同果"，不写归因）。
`RUN=1518 → 1587`（＋69）与 `PASS=1028 → 1071`（＋43）＝推送把 `internal/session`＋`internal/projctx` 两枚包的测试**第一次带进 ubuntu 腿**的账（5.4）。
今天那**一枚 SKIP** 的名字：`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`
（尺：`grep -oE "\-\-\- SKIP: [A-Za-z0-9_]+" | sort | uniq -c` ⇒ **推前与三发推送后的 core 日志都恰好是 2 行同名**，逐发相同；
原因句逐字 `task_output_pointer_notice_test.go:428: 前置条件缺失：造不出真 NTFS junction（mklink /J /tmp/TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3…/001/…）`，
文件在 `internal/tools/task_output_pointer_notice_test.go`，`git ls-files` 尺确认）。
⇒ 票面 `:86-89`（AC#10）当年点名的 `TestWorkspaceSwitchRefusesAJunctionToOutside` **今天不在 SKIP 里、在 `-skip` 正则里**（读数见 §6 末条），
但今天这枚 SKIP **在 `runtests.sh` 打印的 `[fixture]/[reexec]/[opt-in]` 名册里查无**（ledger-hits=0，阳性对照=4/1/1，尺原文见 §6）。

**一条必须补进来的硬发现（与 AC#10 的"响亮"同族；⛔ 但它不是今天才有的，见末段）：**
```
$ gzip -dc job-112085927688.log.gz | grep "unaccounted SKIP"
test-core  Portable package tests (core scope; …)  portable-tests.sh: unaccounted SKIP lines, each with the file:line and reason it printed:
test-core  Portable package tests (core scope; …)      task_output_pointer_notice_test.go:428: 前置条件缺失：造不出真 NTFS junction（mklink /J …）
test-core  Portable package tests (core scope; …)  --- SKIP: TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3 (0.00s)
test-core  Portable package tests (core scope; …)  portable-tests.sh: strict runner exited 1 for the core scope
```
同一枚 core 步里 `=== RUN TestCanonicalize…` 只 **1 次**、`--- SKIP:` 却 **2 行**（`grep -c` 现量），
第 2 行紧跟在 `unaccounted SKIP lines` 那句之后 ⇒ 那第 2 行是**守卫自己的复述**，不是重复执行；
顶层四数里的 `SKIP=1` 因此**没有**被记成 ok（AC#10 括号里那句"步不再把 SKIP 记成 ok"这一形，日志里看得见）。

**但这道守卫不是推送带来的新行为，本腿不许把它写成"今天第一次响"**（尺原文＝逐枚日志 `grep -c 'unaccounted SKIP'`）：
```
pre-112053739109（推前 core） = 1     job-112082660423（run1 core） = 1
job-112085927688（run2 core） = 1     run3-job-112084901657（run3 core） = 1
pre-112053739011（推前 windows） = 2  job-112085927937（run2 windows） = 2
```
⇒ **推前推后一样在响**（每枚"有未记账 skip 的步"打印一次，windows 腿两枚步各一次）。
所以 AC#10 那句"**不在任何台账里**"指的应是名册（`[fixture]/[reexec]/[opt-in]` 名册里查无，ledger-hits=0 已量），
⛔ 不是"日志里静悄悄"。这一格今天到底是"守卫已在做它该做的（⇒ 与票 93 同族那半已成立）"
还是"仍缺一枚把未记账 skip 写进名册的落地（⇒ 票面第二半没做）"——**归编排者裁，本腿不翻框**，
只登记：三枚未记账 skip 的名字与所在步已在 §5.5／§6 末条列全。

### 5.2 `lint` step8 `gofmt (gofumpt)`

三发推送后 ＋ 推前那一发，**红句逐字节同一枚**（尺：`awk -F'\t' '$2 ~ /^gofmt \(gofumpt\)$/'` ＋ `grep -o '##\[error\]…'`）：
```
##[error].scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations
##[error]Process completed with exit code 2.
```
⇒ **推前就红**（推前 `112053738745` 同一枚红句，`grep -oE "##\[error\]\.scratch[^ ]*"` 命中）。
红的是**一枚台账探针用的坏文件**（`.scratch/wisp/probes/185/c1/mut/`），不是产码；
它的直接后果就是 §4 那三枚 skipped（step9/10/11）。⚠ 本腿**没有**判"该不该把这枚 mut 文件排出口径"的权限——那是口径变更，交编排者。

### 5.3 `lint` step12 `staticcheck`（带 `if: ${{ !cancelled() }}`，`ci.yml:263`）

| run | job | findings / 名数（按 `file:line:col` 去重） | 步结论 |
|---|---|---|---|
| 推前 `37396530365` | `112053738745` | `modules=3 packages=36 findings=53`；unique 名行 **54**（含 §5.2 那枚 fs_broken，同日志混排） | failure |
| `37405698188` | `112082660317` | `modules=3 packages=36 findings=52`；unique 名行 **53** | failure |
| `37406757402` | `112085927883` | 同上 `findings=52`；`--- module .: exit=1 packages=34 findings=48`／`tools/d22scan: exit=1 packages=1 findings=3`／`tools/mockllm: exit=1 packages=1 findings=1` | failure |

**作差**（尺：先 `sed -E 's/:[0-9]+:[0-9]+:/:POS:/'` 把行号列号归一，再 `diff`，最后再按原名 `comm` 双向）：
- **今天消失的只有 1 枚**：`internal/agent/approval/queue.go:546:17: func (*Queue).pendingCount is unused (U1000)`（推前有、三发推送后全无）。
- **纯位移、不是新账的 2 枚**：`internal/agent/approval/ticket242_binding_test.go` 的 `sameDigest`（`:111`→`:118`）
  与 `internal/panel/git_test.go` 的 `composerMethodWhitelist`（`:448`→`:450`）——归一化 diff 里**不出现**，
  ⇒ ⛔ 别把行号漂移当新增红。
- 其余 **51 枚逐名逐位相同**（`cmd/wisp` 11 枚、`internal/audio` 4 枚、`internal/llm/*` 若干、`frontend/embed.go:4:1 SA9009`、
  `tools/mockllm/chat.go:219:6`、`tools/d22scan/gitignore.go:13:1` 等）⇒ **推前就红，今天照红**。
- 两发推送后之间：**名册逐字节相同**（`diff /tmp/job-112082660317.names /tmp/job-112085927883.names` → `IDENTICAL`）。

### 5.4 推送第一次带进 CI 的两枚包（"推前没有读数"的那一格）

推前 `37396530365` 的 own-line 名册尺原文（`awk '/own line/ {print $NF}' | sort -u`）：
- core 腿 **25 枚**，里面**没有** `internal/session`、**没有** `internal/projctx`（当时 `core_pin` 里也查无，§3.2）。
- windows 腿 **10 枚**，同样没有这两枚。
今天：core 腿 **27 枚**（多出这两枚），windows 腿 step9 **10 枚**里含 `internal/session`＋`internal/projctx`，
ubuntu 上各打 `ok (own line)` ＋ 顶层 `ok …session 0.127s/0.098s`、`ok …projctx 0.006s`（`37405698188`／`37406757402` 两发都有）。
⇒ 这两枚的第④列**今天第一次有 CI 侧读数**；同一步整体 failure 的折扣见 §6.5。

### 5.5 `test-windows` step7（CLI）与 step9（Portable windows tests）

step9 `internal/risk` 的 **12 枚红名推前就红、今天逐名相同**
（尺：`grep -oE "\-\-\- FAIL: [A-Za-z0-9_]+" | sort -u`，把推前 18 枚红名与今天 step9 的 12 枚逐一 `grep -qx` ⇒ **ALL 12 PRESENT PRE-PUSH**；
四发的该步四数 `RUN=536/577 PASS=362/386 FAIL=12 SKIP=1`，`FAIL=12` 推前推后一样）。
step7 `cmd/wisp`：推前 6 枚、`37405698188` 6 枚（**逐名相同**）、`37406757402` **8 枚**、`37406422380` 回到 6 枚。
今天**只在 `37406757402` 冒出来的 2 枚**＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`、`TestAC14GoSideEvalPushReachesThePage`；
两枚在推前 `37396530365`、`37405698188`、`37406422380` **都是 `--- PASS`**（八枚日志逐名跟踪表见本节末）⇒ 与 `TestConcurrentWritersReaders` 同形，**同码不同果**（§6.6），不是"接入第一天就红"那种发现（AC#2 `:26`）。
`SKIP=1` 那枚今天叫 `TestSyncRedTeamRealOneDrive`（step9）与 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（step7），两枚在 `[fixture]` 名册里同样查无（§6 末条）。

### 5.6 顺手量到、与 AC#7 / AC#9 直接相关的两枚读数（⛔ 不翻框）

- **AC#7（`:75-76`）**：`TestSyncRegistryProbeLive` 在两发 windows 腿上**一次都没被 `-run` 跑过**
  （尺：`grep -c -- "-run TestSyncRegistryProbeLive"` ⇒ `37405698188`/`112082660385` = **0**、`37406757402`/`112085927937` = **0**；
  它出现的 9 行全部落在 `…Pipeline|TestSyncRegistryProbeLive)$` 那串 **`-skip` 正则**里＋`[fixture]` 名册行，
  ⇒ 正是 AC#8（`:77-78`）点的那枚"词出现在日志里 ≠ 被测包跑过"）。
  推前 `37396530365`/`112053739011` 亦 0；ubuntu 腿三发全 0（该用例在 ubuntu 由平台出 scope）。
  ⇒ **"至今 0 次"这一条今天仍成立**，本腿只登记，不裁。
- **AC#9（`:81-85`）**：ubuntu 腿今天**真跑了** `internal/winsec` 的 `!windows` 半边，四发都有顶层 own 结果行：
  推前 `112053739109`/step7 `ok github.com/CarlosShao/wisp/internal/winsec 0.025s`；
  `112082660423`/step7 `0.011s`；`112085927688`/step7 `0.017s`；`112084901657`/step7 `0.012s`；
  再各自带 `ok (own line) github.com/CarlosShao/wisp/internal/winsec`（GUARD B 的锚定行）。
  ⚠ 票面 `:82` 说"全日志里 `internal/winsec` 出现 **0 次**"——那句的来源是旧 run `35599458439`；
  **在今天这四发里该句字面不成立**（尺：`grep -c "internal/winsec"` 全 >0），本腿只把新读数摆出来，AC#9 翻不翻归编排者。

〔建议位〕§5 的净结论只有一条是硬的：**四枚常红里 `test-core step7`（panel 那 4 枚）、`lint step8`、`lint step12 staticcheck`（51 枚原样＋1 枚消失）、
`test-windows step9`（risk 12 枚）全部推前就红；今天新增的红名一共只有 3 枚，
且三枚都在另外两发 PASS（同码不同果）**。⇒ 没有一枚是"接入新包把门撑红"的形状；
`session`／`projctx` 带进来的是 **＋69 RUN／＋43 PASS 的绿读数**，不是红。AC#2 `:26` 那句"接入第一天就红 ⇒ 那是发现"今天**没有对应实例**。

---

## §6 判不动的地方（问 6，自我对抗节 B）

这一节**不是空的**，八枚都真写着东西。凡本腿拿不到的读数，一律具名到 `run + job + step`，
⛔ 没有一处用本地读数或旧归档读数顶替（本腿零 go 命令，压根产不出本地绿）。

| # | 判不动的东西 | 缺的是哪一发的哪一个 job 的哪一步 | 为什么本腿补不了 |
|---|---|---|---|
| 6.1 | `go vet`（主模块）与 `go vet`（`tools/d22scan` 模块）**今天到底是红还是绿** | `37405698188`/job `112082660317`/step10+step11；`37406757402`/job `112085927883`/step10+step11；第三发 `37406422380`/job `112084901703`/step10+step11 全同 | 三枚都是 skipped ⇒ **零字节日志**。上游 step8 `gofmt (gofumpt)` 在三发里都 failure（同一枚红句，见 §5.2），红一天不退，这三枚的读数一天采不到。**"没日志"不等于"没红"**，本腿拒绝把它读成绿。 |
| 6.2 | `gofmt (gofumpt) - the tracked set is the denominator`（161 AC#7 form A）今天的结论 | 同上三发的 **step9**（skipped） | 同一枚病：step8 吃掉 step9。该步承载的是"受跟踪集当分母"那道判据，⛔ 没有替代读数。 |
| 6.3 | **GUARD D 的正向自证**（票 AC#3 逐字要求"人为抽掉一个包证明它会红"） | 不属于任何一发：**CI 上从来没跑过这道反证**。`scripts/portable-tests-selftest.sh`（票 250 AC#3 造的载体）在 `.github/workflows/ci.yml` 里被引用 **0 次**（`grep -c portable-tests-selftest .github/workflows/ci.yml` = 0，两个 workflow 文件各测一遍都是 0），全仓非自身的提及处只有 `scripts/portable-tests.sh:481` 与 `scripts/testdata/portable-tests/go:3,8` 三处**注释** | 本腿只读：造反证要改代码＋推一发 run，越过"零 go 命令／一票面零改动"两条边界。§3 能给到的上限是**退出码路径**＋**三发步级的 marker 计数**（见 §3.3），那是"这一支今天有没有响"，不是"这一支能不能响"。 |
| 6.4 | **linux 侧的 census 形状** | `--scope=census` 的唯一调用点是 `ci.yml:596`，落在 `test-windows`（`ci.yml:419` 起）⇒ 三发打印全是 `census GOOS=windows`；`test-core`（job `112082660423`／`112085927688`／`112084901657`）**没有任何一步跑 census** | 于是"某包在 linux 上被 build tag 清成 0/0"这一形**今天无 CI 读数**。本腿能确证的只有：仓里带 `//go:build !windows` 测试文件的包（`cmd/wisp`／`internal/config`／`internal/models`／`internal/risk`／`internal/winsec`；尺原文＝`git grep -l "go:build !windows" HEAD -- '*_test.go'`，**按 HEAD 取而非工作树**，因为写腿 `232-r2` 正在改工作树）全部**已在 scope 内**，若真在 linux 空掉会撞 GUARD A（`portable-tests.sh:570-583`，`exit 1`），而 linux 腿今天的 step7 红名里**没有** GUARD A 那句（见 §5.3）。**这条推理靠日志措辞，不靠注释**，但它仍不是"linux census 读数"。 |
| 6.5 | `internal/session` 的**第一枚 ubuntu 读数到底算不算真给过结论** | `37405698188`/job `112082660423`/step7 `ok (own line) github.com/CarlosShao/wisp/internal/session`（+ 顶层 `ok github.com/CarlosShao/wisp/internal/session 0.127s`）；`37406757402`/job `112085927688`/step7 同形（`0.098s`） | 读数在，**但该步整体 failure**：run1 那一步有 **1 枚** own-line FAIL（`internal/panel`），run2 那一步有 **2 枚**（`internal/memory`＋`internal/panel`）。⇒ "session 所在的那一步真给过结论"这半句成立；"session 那一步是绿的"不成立。这一格到底按哪一档记，属编排者的口径，本腿只把两半摆开。（另：`scripts/portable-tests.sh:110-127` 那段注释自陈"the ubuntu reading has never been taken by anyone"（逐字在 `:124-126`）——这句**已被上面两发推翻**，属产码注释过期，具名交回，本腿不改一字。） |
| 6.6 | `TestConcurrentWritersReaders`（`internal/memory`）为什么只在 `37406757402` 红 | `37406757402`/job `112085927688`/step7 `--- FAIL`（10.05s）对照 `37405698188`/`112082660423`/step7 `--- PASS`（10.03s）与第三发 `37406422380`/`112084901657`/step7 `--- PASS` | **三发的 Go 源码逐字节相同**：`git diff --name-only ec84cff1..cc315261 \| grep -c '\.go$'` = **0**，`b948bcb8..cc315261` 也只差 1 枚纯台账提交（`cc315261` 改 `msg-a634.txt`＋`pending-and-issues.md`，`git show --stat` 现量）。⇒ 只能写成"同码不同果"，**根因判不了**（时序／负载／runner 差异三种都没尺区分）。⛔ 不许据此放宽任何断言——那是 AC#2 明令禁止的方向。 |
| 6.7 | `37405698188` 的 `lint`／`test-core` 两枚 job 的全量日志**没有步名列** | 尺原文：`gh run view --job 112082660317 --log` 与 `--job 112082660423 --log`，第二列逐行都是 `UNKNOWN STEP`（对照 `112085927883`／`112085927688`／两发 `test-windows` 都带真步名） | ⇒ 那一发的这两枚只能**按内容归步**（如 `##[error]...fs_broken.go` 落在 gofmt 的 `Run` 块里），归到步号的映射靠 `ci.yml` 的步序推。凡本件里由 run1 这两枚给出的"第 N 步"，都要按这个折扣读；run2 的读数不受影响（步名列齐全）。这是本腿尺子的缺陷，不隐瞒。 |
| 6.8 | AC#1 的分母：票标题写 **33**，census 今天打 **35** | 三发 census 步（`test-windows` step8）逐字 `go list ./... = 35 packages` | 差 2 枚是谁挤进来的，**本腿判不了**：`tools/d22scan`／`tools/mockllm`／`scripts/spike` 是**独立 module**（`git ls-files` 现量 4 枚 `go.mod`：`go.mod`、`tools/d22scan/go.mod`、`tools/mockllm/go.mod`、`scripts/spike/go.mod`），主 module 的 `go list` 里**没有**它们；`internal/session`／`internal/projctx` 又明明早已在 35 行里。⇒ "20 个 / 33 个"那句今天**过期**，但"20"该重算成几、33 该改成 35 还是把另两个 module 的包并入分母，属**口径变更**（碰 D／C 契约级），交编排者，本腿不自裁。 |

**另外三处"看似能判、本腿故意不判"**
1. **AC#6 翻不翻**——§4.3 两读法已摆平，判语留白。
2. **`Post Cache third_party` skipped 是否真的丢了那次缓存写入**（即下一发是否付冷缓存代价）——那两枚 post 步零输出，无从读；三发 census 步都在 step7 之后仍 success，只说明"缓存缺失没弄死测试步"，**不等于** post 的副作用没丢。
3. **今天 3 枚未登记的 skip 该按 AC#10 翻哪一格**——现量事实：三发里被 `runtests.sh` 打印的 `[fixture]/[reexec]/[opt-in]` 名册中**查无**这三枚（`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`／`TestSyncRedTeamRealOneDrive`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 的 ledger-hits 全 = 0，尺原文见本节的 `grep -cE "portable-tests.sh: {1,4}$n"`），而票面 `:86-90`（AC#10）点名的 `TestWorkspaceSwitchRefusesAJunctionToOutside` 今天**已在 `-skip` 正则里**（属"已接住"，与票写它时"未记账"的形状不同）。要不要据此翻 AC#10、还是把这三枚当**新发现的未记账 skip** 另开票 ⇒ 编排者裁。本腿只登记名册。
