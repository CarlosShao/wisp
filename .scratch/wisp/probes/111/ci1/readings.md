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

<!-- 先写满：真写内容；拿不到的读数要具名到 run+job+step，不许用本地或旧归档顶替 -->
