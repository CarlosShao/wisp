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

<!-- 待填：AC#1/2/3/4/6/7/8/9/10 十枚框原文，sed 行号已核对与否要具名 -->

---

## §2 AC#1 第④列名册（问 2）

<!-- 待填：包 × (有无测试文件 / 在不在 CI 某一步 / 那一步真给过结论的 run+job+step)，含 census 步打印的四数逐作业 -->

---

## §3 GUARD D 与空分母（问 3）

<!-- 待填：session/watchdog 今日跟踪的 _test.go 枚数、在册与否、census success 的成因（步级证据 + portable-tests.sh 退出码路径） -->

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

## §5 红名作差（问 5）

<!-- 待填：test-core step7 / lint step8+step12 staticcheck / test-windows step7+step9 红句逐名，与 37396530365 同口径作差 -->

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
