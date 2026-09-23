# 136 / AC#1 对抗验收（零样本 fail-closed 那枚钉）— 非实现者

验收方：`acceptor-ticket136-r1`。判据物：票 `.scratch/wisp/issues/136-*.md` 的 **AC#1**。
被验面：`internal/observe/sampler_zerosample_136_test.go`（枚 `4fc65dd` 新增）+ 它所钉的守卫
`internal/observe/sampler.go:332-340`。本程**一字未改生产码与测试码**，所有变异落在仓外纯净树。

---

## §0 锚定与被验面完整性

| 项 | 读数 |
| --- | --- |
| 我锚的 sha（被验版本） | `e8190bfe5b6888244a9ef70bbe2641f17cbe9a20`（短 `e8190bf`） |
| 开工那一刻的 HEAD | `5265c3ad844eb549a05f58011eed0d1e9bd86b8a`（短 `5265c3a`，`2026-09-23T13:51:48Z` 现量） |
| 写本节时的 HEAD | `6bb5a567ab4bfb329323016dbf4b40a8e49110cd`（短 `6bb5a56`，`2026-09-23T14:09:09Z` / 22:09 +08 现量）——树一直在动，本节之后每次取数都重取 |
| 抽取方式 | `git archive e8190bf \| tar -x -C /d/tmp/wisp136-acc-r1` ⇒ `extract rc=0`；**仓内未建 worktree、未 checkout** |
| 变异场 | `/d/tmp/wisp136-acc-r1`（主树）、`/d/tmp/wisp136-acc-r1-tree2`（第二棵树，同 sha 独立抽取）、`/d/tmp/wisp136-acc-r1-tree3`（**装钉前的父树 `a122240`**）；驱动 `/d/tmp/wisp136-acc-r1-runs/{mut.py,battery.sh,whole-repo.sh,tree3-whole.sh}`，一律只建不删 |
| 工具链 | `go version go1.27.1 windows/amd64`；容器 `golang:1.27` = `go1.27.1 linux/amd64`（§4） |
| 仓内 `internal/observe/**` | **只读**，本程未写入一次（`git status --porcelain` 全程只有 §0 末那枚来源未明的未跟踪件） |

被验面完整性（blob 级，不是行数级）：

| 文件 | e8190bf 里的 blob | 我从纯净树 `git hash-object` 算出的 | 行数 |
| --- | --- | --- | --- |
| `internal/observe/sampler_zerosample_136_test.go` | `55a8d25b7b5cb2983735d924b6750a543fd0b72e` | 同左（逐位相同） | 164 |
| `internal/observe/sampler.go` | `6a94434a0a77fa8b5a36e3e70fc9109f86bb5fa2` | 同左 | 503 |

- `git log --all -- internal/observe/sampler_zerosample_136_test.go` ⇒ **只有 `4fc65dd` 一枚**；
  `git show --stat e8190bf` ⇒ 只动票面（`1 file changed, 44 insertions(+)`，删除列 0）。
  即任务书里"枚 `e8190bf` 又追加了票面一节"为真，**被验的那枚钉本身在 `4fc65dd` 与 `e8190bf` 之间一字未动**。
- 当前 HEAD（`6bb5a56`）上该文件 blob 仍是 `55a8d25…` ⇒ 兄弟在飞的只有 `cmd/wisp/**` 与票面/证据，
  被验面与本程读数**不因树在动而失效**。
- 环境坑按派单处理：容器挂载用 `MSYS_NO_PATHCONV=1` + `/d/...`，并在容器里先 `ls -l /src/go.mod`
  （读数：`-rwxrwxrwx 1 root root 883 … /src/go.mod`）证真挂上，不是静默空挂载的假绿。

---

## §1 四发变异三态复现（M1/M2/M3/M4）＋补三发（M6/M7/M8）

基线（主树，未变异）：`go test -count=1 -v ./internal/observe/` ⇒ **rc=0 / RUN=56 / PASS=56 / FAIL=0 / SKIP=0**。
与实现方"装钉后 56 枚"对得上；票 134 那本的 54 枚是装钉前分母（§3 档 C 复算）。

每一发的动作序列固定为：`restore` 到纯净 ⇒ 落一发 ⇒ **落地证明**（`diff -u` 打出被改那一行 + `go build ./...` rc=0）
⇒ 读数 ⇒ `restore` ⇒ 再证（`diff -q` 与纯净逐字相同 + build rc=0 + 复绿四数）。四发之间**不叠加**（驱动每发先无条件 restore）。

| 发 | 改哪一行（改成什么） | 落地 | build rc | 全包四数 | 红名 | 红点 | 另一腿同发 | 还原后 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| M1 删守卫 | 删 `sampler.go:332-340` 整块 + 随之未用的 `"fmt"` import | diff 两块都在，`grep -c 'len(rep.Samples) == 0'`⇒0 | 0 | rc=1 / 56 / 55 / 1 / 0 | `TestSampleStateZeroSampleWindowFailsClosed` | `sampler_zerosample_136_test.go:75` | 正向对照腿 `--- PASS` | rc=0 / 56 / 56 / 0 / 0 |
| M2 零足迹永不误丢 | `:290` `<= 0` ⇒ `< 0` | diff 单行 | 0 | rc=1 / 56 / 55 / 1 / 0 | 同上 | **`:61` 前提腿**：`precondition broken: unmeasurable window produced 5 samples` | `--- PASS` | rc=0 / 56 / 56 / 0 / 0 |
| M3 逢读数都丢 | `:290` `<= 0` ⇒ `>= 0` | diff 单行 | 0 | rc=1 / **52** / **47** / **5** / 0 + **panic 1 枚** | 见 §5（钉根本没跑到） | 定点读数 ⇒ **`:151` 正向对照腿** | 定点读数里腿 A `--- PASS` | rc=0 / 56 / 56 / 0 / 0 |
| M4 守卫在、牙没了 | `:337` `Pass: false` ⇒ `Pass: true` | diff 单行 | 0 | rc=1 / 56 / 55 / 1 / 0 | 同 M1 | `:75`（verdicts 里 `sampling … Pass:true Gate:true` 全绿那串） | `--- PASS` | rc=0 / 56 / 56 / 0 / 0 |
| M6 降级成非门（补） | `:337` `Gate: true` ⇒ `Gate: false` | diff 单行 | 0 | rc=1 / 56 / 55 / 1 / 0 | 同 M1 | `:75` | `--- PASS` | rc=0 |
| M7 换掉那行的名字（补） | `:336` `Metric: "sampling"` ⇒ `"sampling_neutered"` | diff 单行 | 0 | rc=1 / 56 / 55 / 1 / 0 | 同 M1 | **`:84`**：`carries no \`sampling\` verdict row` | `--- PASS` | rc=0 |
| M8 守卫条件反写（补） | `:332` `== 0` ⇒ `!= 0` | diff 单行 | 0 | rc=1 / 56 / **52** / **4** / 0 | 两腿**都**红 | `:75`（零样本窗口反而绿）＋ `:158`（`the sampling guard fired on a measurable window: {Metric:sampling Measured:6 valid / 0 errors …}`） | — | rc=0 |

**四发是否都响**：都响。M1/M4 响"钉本身"，M2 响"前提腿"，M3 响"正向对照腿"。
**每发是不是只响该响的那条腿**：M1/M2/M4/M6/M7 五发都是"包内只有 1 枚 FAIL、另一腿同发 `--- PASS`"；
M3 那一发**全包读数取不到两条腿**——既有的 `TestSamplerGoroutineAccountingFollowsRegistry` 先 panic 杀了测试二进制
（读数 `47 PASS / 5 FAIL + panic`，与实现方报备的事故①**逐位相同**），只能用定点读数，取到的是
"腿 B 红在 `:151`、腿 A 同发 `--- PASS`"。所以"自证腿不哑"这一条我这边成立，
但**M3 的"另一腿仍绿"是定点读数给的，不是全包关门读数给的**，这条边界要写在这。

补的三发是往"守卫被无声改弱而不被删除"那一族形状上打的：M6（把门降级成记录行）、M7（把那一行改名让找不到的循环落空）、
M8（把判据反写）——三发都响，且 M6/M7 只响该响的那条腿。M8 让两腿同时红是**正确**行为：条件反写同时破坏了
"零样本要红"和"可信窗口不许长门"两面，两条腿各抓一面。

---

## §2 断言链是不是落在生产码上（含"是不是自己算自己"）

### 2.1 逐跳（测试行 ⇒ 生产行）

生产符号定位：`NewSampler` `internal/observe/sampler.go:208`、`SampleState` `:252`、
`StateReport` `:158`（`Samples` 标签 `json:"samples"` 在 `:163`、`Pass` 标签 `json:"pass"` 在 `:190`）、
`Verdict` `:124`（`metric` `:125`、`pass` `:128`、`gate` `:133`）。

| 测试里的断言 | 读的是生产哪个值 | 生产出处 |
| --- | --- | --- |
| `:60 len(rep.Samples) != 0` | 采样窗口里的有效样本数 | `sampler.go:296-298`（`derive` + `append`），丢弃支 `:290-294` |
| `:63 reads == 0` | fixture 自己的读计数器（证明"树真被读过"） | `sampler.go:285` 的 `s.tree.ReadTree()` 循环 |
| `:66 rep.SampleErrors == 0` 为假 | 丢弃计数 | `sampler.go:288` / `:293` |
| `:69` 含 `"zero private working set"` | 生产那句错误原文 | `sampler.go:294` 字面量 `"read returned a zero private working set for a live tree"` |
| `:74 if rep.Pass` **（本格的钉）** | `rep.Pass` 由生产那段 gate 归约算出 | `sampler.go:341-346`（`rep.Pass = true` ⇒ 任一 `Gate && !Pass` 即 false） |
| `:79-84` 找 `Metric == "sampling"` | 那一行 verdict 只能由守卫产出 | `sampler.go:335-339`；全仓 `"sampling"` 这个 Metric 只有这一处生产者（`grep -rn '"sampling"' --include=*.go` ⇒ `sampler.go:336` + 本测试 4 处引用） |
| `:86 guard.Gate` / `:89 !guard.Pass` | 同一行的 `Gate:true` / `Pass:false` | `sampler.go:337` |
| `:92 HasPrefix(guard.Measured, "0 valid")` | 同一行的 `Measured` | `sampler.go:336` 的 `fmt.Sprintf("%d valid / %d errors", …)` |
| `:99-106` 隔离环：除 `sampling` 外不许有别门红 | `buildVerdicts` 产出的全部行 | `thresholds.go:84-114`（mem `:116`、cpu `:133`、gdi `:88`、handles `:100`、goroutine `:156`、sleeping 两行 `:190-206`） |
| `:110-131` 出线 JSON 的 `samples` / `pass` / `"metric":"sampling"` / `"gate":true` | 生产结构体的 json 标签 | `sampler.go:163` / `:190` / `:125` / `:133` |

### 2.2 是不是"自己算自己的期望值"，或把生产逻辑抄一遍从而恒真

结论：**不是**。三条依据：

1. **没有本地重算**。测试文件里既不调 `buildVerdicts`，也不自己遍历 verdicts 去算一个 `expectedPass` 再和 `rep.Pass` 比。
   `:74` 直接读生产产物 `rep.Pass`。唯一带"计算味"的是 `:92` 的 `"0 valid"` 前缀和 `:79` 的 `"sampling"` 字面量——
   它们是对生产输出字符串的**匹配**，不是对生产逻辑的**复写**（生产那段拼串在 `sampler.go:336`，测试里没有第二份 `Sprintf`）。
2. **恒真最可能的入口已被 :99-106 那圈堵死**。若没有这圈，`rep.Pass=false` 可能来自任何一枚别的红门，
   那么删掉守卫后别的门仍红 ⇒ 钉仍绿 ⇒ 恒真。实测：M1 落地后该钉在 `:75` 红，说明"零样本窗口反而全绿"这一形
   确实是删守卫后的真实生产输出（M1 读数里打印的 verdicts 串正是 `tree_private_bytes 0.0MB Pass:true`、
   `cpu_percent_all_core 0.000% Pass:true` … 九行全绿）；隔离环在同发未触发（没有出现"别的门红了"那条消息）。
3. **五发独立攻击都落空**（不是只测"删除"这一形）：把守卫**改弱**（M6 降级为非门、M7 改名、M4 翻 Pass）
   或**反写**（M8）都红，且红点分别落在 `:75` / `:84` / `:75` / `:75+`:158`。
   特别是 M7：那一行还在、还判红、`rep.Pass` 仍 false，只有"名字被换掉"——
   钉靠 `:79-84` 的"必须找得到那枚 `sampling` 行"抓住了它。这说明该钉不只钉"pass 值"，
   还钉"报告要自己说清为什么不许过"这半句（票面 AC#1 措辞里的"出报告即带标记 ⇒ 判红"）。

一处**轻度冗余**（不构成恒真，登记不判缺陷）：`:60` 已断 `len(rep.Samples)==0`，`:92` 的 `"0 valid"` 前缀因此近乎必然；
它多要的是"那行得把计数写在脸上"这层表达，价值在 M7 之外没被单独打死过（见 §7 `R-136-3`，低）。

---

## §3 全仓对照：我复算到哪一档

派单只要求"至少两棵树各跑一次 `go test -count=1 ./internal/observe/...`"。我做了三档，**其中两档是全仓**：

### 档 A｜本包，两棵 e8190bf 独立树

| 树 | 状态 | 四数 |
| --- | --- | --- |
| `/d/tmp/wisp136-acc-r1` | pristine | rc=0 / 56 / 56 / 0 / 0 |
| 同上 | M1 | rc=1 / 56 / 55 / 1 / 0，红名 `TestSampleStateZeroSampleWindowFailsClosed` |
| `/d/tmp/wisp136-acc-r1-tree2`（同 sha 第二次 `git archive` 独立抽取） | pristine | rc=0（`go test -count=1 ./internal/observe/...`） |
| 同上 | M1（**从主树逐字拷过去的同一份变异文件**，落地后 `grep -c` 两处都为 0、`build rc=0`） | rc=1，红名同上 |

⇒ 两棵树"红名集合逐名相同"在我这边成立：都只有 `TestSampleStateZeroSampleWindowFailsClosed` 一枚。

### 档 B｜全仓 `go test -count=1 ./...`，e8190bf 两态（windows/amd64 宿主）

| 树态 | 顶层 rc | 红包集合 | `internal/observe` 那一行 | 全包 `--- FAIL` 名 |
| --- | --- | --- | --- | --- |
| pristine | 1 | `cmd/wisp`、`internal/panel`、`internal/risk` | `ok … internal/observe 1.576s` | `TestComposerRenderFixtureTellsTheTruth`、`TestResolvePerCallBudget` |
| M1 | 1 | 上述三枚 **+ `internal/observe`** | `FAIL … internal/observe 1.613s` | 上述两枚 **+ `TestSampleStateZeroSampleWindowFailsClosed`** |

⇒ 装钉后，拆守卫在**全仓**只多红这一枚；其余红包/红名逐名不变。

### 档 C｜决定性那档：装钉**前**的父树 `a122240`（`4fc65dd` 的父 commit），全仓两态

`a122240` 里 `internal/observe/sampler.go` 与 `e8190bf` **逐字节相同**（`diff -q` 无输出，503 行），
`ls internal/observe/ | grep -ci 136` ⇒ **0**（钉还不存在）。

| 树态 | 顶层 rc | 红包集合 | `internal/observe` | 全包 `--- FAIL` 名 |
| --- | --- | --- | --- | --- |
| pristine | 1 | `cmd/wisp`、`internal/panel`、`internal/risk` | `ok … 2.853s`，包内 **54/54** | 同档 B 两枚 |
| M1（守卫删掉，钉不存在） | 1 | **同上，`diff` 无输出** | `ok … 1.580s`，包内 **54 PASS / 0 FAIL / rc=0** | **同上，`diff` 无输出** |

⇒ **票 134 那句"我答不出"被两侧都算出来了**：装钉前拆掉守卫，全仓没有任何一枚仪器的状态发生变化（连本包都仍 `rc=0`、54 枚全过）；
装钉后同一发变异，全仓只多红一枚，名字可点。这一形我**造不出**"守卫被拆掉而本仓没有任何仪器状态变化"（在装钉后的树上）。

### 读数分歧（两边都留，不改实现方一行）

| 项 | 实现方读数（`docs/evidence/s1/136-ac1-zero-sample-nail.md` §1.3，锚 `09edf02`） | 我的读数（锚 `e8190bf`，windows/amd64） | 以哪份为准 / 为什么 |
| --- | --- | --- | --- |
| 全仓既有红包集合 | 4 枚：`cmd/wisp`、`internal/agent/approval`、`internal/panel`、`internal/risk` | 3 枚：**无** `internal/agent/approval`（该包打 `ok 0.720s`） | **各自锚点下各为准**。分歧与本格判据无关：两边要的是"两棵树之间逐名相同"，这一点两边都成立；`git diff --name-only 09edf02..e8190bf -- internal/agent/approval internal/panel internal/risk cmd/wisp` 输出为空 ⇒ 不是锚点漂移解释得通的，剩下能说的是负载/环境差异，我未追（§8）。**不把"3 枚"写成对实现方的纠正**，只登记。 |
| `ci.yml:288` 指什么 | 实现方引 `ci.yml:288` 作 observe 进 CI 的落点；派单里编排者记为"行号引偏（那是'哪些包不在本 scope'的注释行）" | `grep -n` 出：`266` 是 step 名，`267-287` 是注释，**`288: run: bash scripts/portable-tests.sh --scope=core`** | **以我的实测为准**（纯净树字节级 grep，可重算）：`ci.yml:288` 就是那条 run 命令本身，实现方引对了；编排者登记在票面上的那句"行号引偏"不成立。详见 §4。 |

### 我复算到哪一档（不许把本包写成全仓）

我在 **windows/amd64 宿主**上把**全仓 `./...` 两态**复算到了（档 B、档 C），本包两棵树也各跑了一次（档 A）。
我**没有**在 linux 上跑全仓 `./...` 两态对照；linux 腿我只测了"这枚钉进不进 core scope、有没有被 skip 收走"（§4）。
