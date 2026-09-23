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
| 仓内 `internal/observe/**` | **只读**，本程未写入一次。仓树在我跑动期间被别人动过的是别的路径（末次现量 `22:24 +08`：`git status --porcelain` ⇒ ` M .scratch/wisp/issues/122-*.md`、` M docs/evidence/s1/133-ac2-fix.md`、`?? docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`），**都不是我改的、我也没提交它们**；本程三枚 commit 的暂存清单每次只有本文件一枚 |

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

另两发 **M9**（把 `:297` 的 `append` 换成 `_ = sample`，丢弃分支一字未动）与 **M10**（把 `CheckSettle` 的
`sampler.go:477` `> 0` 换成 `>= 0`）不在 AC#1 的判据物上，各归 §5.1（证 panic 的原因）与 §7 `R-136-1`（settle 那侧无钉）。

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
| `:66` 要求 `rep.SampleErrors != 0` | 丢弃计数 | `sampler.go:288` / `:293` |
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
我**没有**在 linux 上跑全仓 `./...` 两态对照；linux 腿我只测了"这枚钉进不进 core scope、有没有被 skip 收走、能不能真红"（§4）。

---

## §4 CI 落点：我自己量的读数

### 4.1 名单与行号（可重算）

| 主张 | 我的实测（锚 `e8190bf` 纯净树） |
| --- | --- |
| observe 在 core 名单里 | `scripts/portable-tests.sh:140` 逐字 `github.com/CarlosShao/wisp/internal/observe`（`core_pin` 块从 `:124` 起）；同文件 `:175` 的 core scope glob 含 `./internal/observe/...` |
| CI 哪一步跑它 | `.github/workflows/ci.yml` job `test-core`（`:224`）、`runs-on: ubuntu-latest`（`:225`）、step 名在 `:266`、**`:288` 就是 `run: bash scripts/portable-tests.sh --scope=core`** |
| 编排者那句"行号引偏" | **不成立**：`:267-287` 才是注释块（讲的是 GUARD A/B/C 与"哪些包没有分母"），`:288` 是 run 命令行本身。实现方引对了。重算：`grep -n '' .github/workflows/ci.yml \| sed -n '266,289p'` |
| Windows 腿 | `--scope=windows` 的 scope 在 `:184-190`、`win_pin` 在 `:151-160`，**两处的包名单里都没有 observe** ⇒ 实现方"Windows 那条腿本来不含 observe"为真 |

### 4.2 有没有被 skip／ledger 收走

- `scripts/portable-tests.sh` 的 `ledger`（`:319-331`）共 **11 条**，容器里实跑打印的是
  `portable-tests.sh: 11 ledger entries, 8 accounted on this platform`；
  生效于 linux 的 8 条逐名：`TestDefaultDeadlineWallClockMeasurement`、`TestSubprocessCrashWriter`、
  `TestRealDownloadVadThroughPipeline`、`TestRealDownloadPuncArchiveThroughPipeline`、
  `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`、`TestWorkspaceSwitchRefusesAJunctionToOutside`、
  `TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`。
- 实跑打出的完整 pattern（原文）：
  `^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestC26RewrittenSyncRootDoesNotDisarmSuspectNet|TestWorkspaceSwitchRefusesAJunctionToOutside|TestD34WriteMatrix|TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop)$`
  ⇒ 本票两枚用例名**都不在其中**；且 pattern 是 `^(…)$` 全锚，`-skip` 不可能顺带吞掉别的名字。
- 反向那道保险也在：ledger 条目若在该包里 `-list` 不到（改名／删掉／被 build tag 挡掉），脚本走 `stale` 分支 **exit 1**（`:396-404`）。
  也就是说"往名单里塞一条来收走新钉"这条路本身就是会红的。

### 4.3 ubuntu 那腿到底有没有分母（容器原生 `golang:1.27` = `go1.27.1 linux/amd64`）

挂载按派单给的形做，先证真挂上再读数：`MSYS_NO_PATHCONV=1` + `/d/...` + 容器内
`ls -l /src/go.mod` ⇒ `-rwxrwxrwx 1 root root 883 … /src/go.mod`；M1 那一发用**单文件 bind-mount**
（`-v …/sampler.go.M1:/src/internal/observe/sampler.go:ro`）落，宿主那棵树全程不被改（落地证明：
容器内 `grep -c "len(rep.Samples) == 0"` ⇒ RUN A 打 `1`、RUN B 打 `0`，`ls -l` 显示 18636 字节 vs 19098）。

| 读数 | RUN A（pristine） | RUN B（M1） |
| --- | --- | --- |
| `bash scripts/portable-tests.sh --scope=core` 对 observe 那一行 | `ok (own line)  github.com/CarlosShao/wisp/internal/observe` | **`FAIL (own line) github.com/CarlosShao/wisp/internal/observe`** |
| core 步里 `=== RUN`／`--- PASS` 两条新用例 | `1471:=== RUN TestSampleStateZeroSampleWindowFailsClosed` / `1472:--- PASS (0.12s)` / `1473:=== RUN …Trustworthy…` / `1474:--- PASS (0.06s)` | `1471:=== RUN …ZeroSample…` / **`1473:--- FAIL (0.05s)`** / `1474:=== RUN …Trustworthy…` / `1475:--- PASS` |
| 整条 core 步里的红名 | `TestComposerRenderFixtureTellsTheTruth`（`internal/panel`，与守卫无关） | `TestSampleStateZeroSampleWindowFailsClosed` ＋ 上面那枚 |
| 定点：`sh tools/d22scan/runtests.sh ./internal/observe/... -count=1 -skip "<上面那串原文>"` | rc=0，**RUN=56 / PASS=56 / FAIL=0 / SKIP=0**，`runtests.sh: OK … top-level: PASS=56 FAIL=0 SKIP=0, === RUN=56, '[no tests to run]'=0` | rc=1，**RUN=56 / PASS=55 / FAIL=1 / SKIP=0**，红名＝那枚钉 |

⇒ **"这枚钉会真进 CI"为真，且不是靠读 yml 读出来的**：ubuntu 腿有分母（56 枚、逐名可点），
把守卫拆掉它在 CI 的形状里就是红的；`-skip` 名单没收到它；Windows 腿不含本包这条如实保留。

顺带量到、**不归本票**的一枚（写清免得下一个人重新发现）：`e8190bf` 的 pristine 树在 linux 上
`--scope=core` **本来就 rc=1**（`internal/panel` 的 `TestComposerRenderFixtureTellsTheTruth` 在 POSIX 上红）。
门禁处在这种"已经红着"的状态时，新增真伤不改变步的颜色，只多一行红名；
归因还能给是因为 `portable-tests.sh` 逐包打 `ok/FAIL (own line)`（GUARD B）。这条属票 114／门禁那族账，不写进 AC#1 判据。

---

## §5 AC#8 那本邻居账：定性 + 被它拖走的枚数（不修）

AC#8 的文本不在我锚的 `e8190bf` 票面上（该 sha 的票面是 7 格 AC；AC#8 是编排者 22:1x 追加的，
`git log -- .scratch/wisp/issues/136-*.md` ⇒ 追加发生在 `6da6710` 之后）。我按**当前票面**上 AC#8 的结案判据①来定这一格。

### 5.1 panic 为真：两发"空 Samples"，走两条互不相干的路

| 发 | 怎么把 `Samples` 弄空 | 落地证明 | 全包四数 + panic |
| --- | --- | --- | --- |
| M3（＝§1 那发） | `sampler.go:290` `<= 0` ⇒ `>= 0`（逢读数都丢） | `diff` 单行 + `go build ./...` rc=0 | rc=1 / RUN=**52** / PASS=**47** / FAIL=**5** / SKIP=0 + **1 枚 panic** |
| M9（补，独立路） | `sampler.go:297` `rep.Samples = append(rep.Samples, sample)` ⇒ `_ = sample`（**丢弃分支一字未动**） | `diff` 单行（`-rep.Samples = append…` / `+_ = sample`）+ `go build ./...` rc=0 | rc=1 / RUN=**52** / PASS=**47** / FAIL=**5** / SKIP=0 + **同一枚 panic** |

两发的 panic 原文同形：

```
--- FAIL: TestSamplerGoroutineAccountingFollowsRegistry (0.03s)
panic: runtime error: index out of range [0] with length 0 [recovered, repanicked]
github.com/CarlosShao/wisp/internal/observe.TestSamplerGoroutineAccountingFollowsRegistry(...)
	D:/tmp/wisp136-acc-r1/…/internal/observe/sampler_test.go:308 +0x23b
```

M9 的意义：**丢弃分支没动、`sampling` 守卫还在**，只把 `Samples` 弄空 ⇒ 照样 panic ⇒
所以吞读数的原因确实是"`rep.Samples` 为空 ＋ `:308` 无长度守卫直取 `[0]`"，不是 M3 那发的语义。
另把该用例单拎出来跑（`-run TestSamplerGoroutineAccountingFollowsRegistry`）在同发下也是 panic ⇒ 归因到枚。

### 5.2 被它拖走的枚数（逐名，roster 差集，不是估计）

以未变异基线的 `=== RUN` 名册（56 枚）为分母，两发各跑一次后取差集，**两发拖走的是同一批 4 枚**：

| # | 被拖走（该发里既没红也没绿，压根没跑到） | 它本来在名册的第几位 |
| --- | --- | --- |
| 1 | `TestLiveRegistryBaselineWithinSleepingGate` | 53 |
| 2 | `TestThresholdTableCoversAllStates` | 54 |
| 3 | **`TestSampleStateZeroSampleWindowFailsClosed`**（本票 AC#1 那枚钉） | 55 |
| 4 | **`TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`**（本票 AC#1 正向对照腿） | 56 |

（`TestSamplerGoroutineAccountingFollowsRegistry` 自己是第 52 位，它**有**读数——`--- FAIL` ＋ panic，
不计入"被拖走"；差集里也没有别的名字，`56 − 52 = 4` 与逐名列表一致。）

### 5.3 定性

- **既有件，不是装钉造出来的**：`git blame -L 305,310` ⇒ `739bb15f (CarlosShao 2026-09-20)`，早于 `4fc65dd`（本票钉）。
- **同包 `[0]` 直取普查**：`grep -rn 'Samples\[\|Verdicts\[' internal/observe/*_test.go` ⇒ 只有
  `sampler_test.go:308` 一处是无守卫直取；本票那两处在 `sampler_zerosample_136_test.go:79-80`，走的是
  `for i := range rep.Verdicts` 的下标，长度天然安全。⇒ AC#8 里"别一把改完"那句眼下没有第二枚要改。
- **对本票 AC#1 的实际影响**：只影响"读数的形状"，不影响判据。M3 那一发的**全包关门读数**取不到两条腿
  （被吞了），我用定点读数补齐（腿 B 红在 `:151`、腿 A 同发绿）。而 AC#1 的**结案判据那一发（M1）不受它影响**：
  M1 落地时全包 56 枚跑完、`panic=0`、红名逐名可点（§1）。
- **AC#8 结案判据①＝我已复算成立**（造得出空 `Samples`、panic 为真、拖走 4 枚可逐名）。判据②（加守卫后同发只红这一枚、
  其余照常跑完）与③（M3 因此可以走全包读数）**没做**——那是修的人的账，本程按派单**一字未修**。
- 归谁：**归 AC#8 那一格**（修法只许动 `sampler_test.go` 那枚用例，不许改成 Skip——跳过＝把"没测"洗成"通过"）。

---

## §6 AC#1 终判，以及哪几项该另立一格

### 6.1 票面 AC#1 的每一条判据，逐条对表

| 票面 AC#1 的要求（原文拆条） | 我的裁定 | 依据 |
| --- | --- | --- |
| "把 `sampler.go:332-340` 那枚守卫退回旧实现或删掉 ⇒ 本用例必须转红，红名点到本用例" | **成立** | §1 M1：`rc=1 / 55 PASS / 1 FAIL`，红名 `TestSampleStateZeroSampleWindowFailsClosed`，红点 `sampler_zerosample_136_test.go:75`。"退回旧实现"这一支我另量了一句：`git log -L 332,340:internal/observe/sampler.go` ⇒ 那段区间历史上**只有 `adbdfa2` 一枚 commit 碰过**（一次写成、从未改过），所以"旧实现"＝没有旧实现，删掉就是它唯一能落的形式 |
| "变异先证落地（`grep -n` 出被改后那一行 ＋ `go build` rc=0）再读数；三态原文都要贴" | **成立** | §1 表：四发 + 补五发每发都留了 `diff -u` 被改行原文与 `build rc=0`；还原后 `diff -q` 逐字相同 + build rc=0 + 56/56 复绿 |
| "不许用放宽任何断言或阈值换绿" | **成立** | 装钉那一枚 commit（`4fc65dd`）净面只有新增的 164 行测试文件（`git show --stat` ⇒ `1 file changed, 164 insertions(+)`），`thresholds.go`／golden 一字未动；fixture 落在冻结阈值之下（M1 打印的 verdicts 串可自证：`<=25MB`、`<=0.5%` 那些行本来就是绿的） |
| "也不许把这枚钉写成恒真（自证腿不许是哑的：往它自己的 sanity 腿上也拆一发看会不会红）" | **成立** | §2：没有本地重算、`:99-106` 隔离环堵掉"靠别的红门撑绿"这条路；§1：M2 单独打死前提腿（`:61`）、M3 单独打死正向对照腿（`:151`），另有 M6/M7/M8 三发打"守卫被改弱"族 |
| 票面开头那句："新增用例，令'零样本 ⇒ 不出报告 / 出报告即带 `samples=0` 标记 ⇒ 判红'这条行为被钉住" | **成立（钉的是后半句那支）** | 生产在零样本时**确实出报告**（`SampleState` 返回非 nil + `err==nil`），所以能钉的就是"出报告即带标记 ⇒ 判红"：`:74` 判红、`:79-94` 要求那枚 `sampling` 门行在、写着计数、且是 gate；`:110-131` 要求出线 JSON 带 `samples` 空 + `pass:false` + `"gate":true` 那行 |
| **本格唯一结案判据**："若这枚守卫被无声删掉，本仓哪一枚仪器会红？"——要点名到用例 | **答出来了，且两侧都是读数** | 点名：`internal/observe` 包的 `TestSampleStateZeroSampleWindowFailsClosed`。装钉**前**（父树 `a122240`，全仓 `./...` 两态）拆守卫 ⇒ 红包集合与红名集合 `diff` 无输出、`observe` 仍 `ok`（§3 档 C）＝"当时确实谁都不红"；装钉**后**同一发 ⇒ 全仓只多红这一枚（§3 档 B），在 CI 的 linux 腿形状里也是同一名（§4.3） |

### 6.2 终判

**PASS，无附条件。** 判据物要的那句"点名那枚会红的仪器"我这边不是抄来的：我自己造了 M1、M4、M6、M7 四发"守卫还在但已经不咬人"与"守卫被删"两种形状，
四发都在**同一枚**用例上红，红点分别是 `:75` / `:75` / `:75` / `:84`；我造不出"守卫被拆掉而本仓没有任何仪器状态变化"那一形（在装钉后的树上）。
派单里那句"若造得出就退回"没被触发。

三件我没否决它、但要写清的边界（都不构成退回理由）：
① M3 那一发的**全包关门读数**取不到两条腿（被既有件 panic 吞了，§5），它的"另一腿仍绿"是定点读数给的；
② 恒真攻击里有一处**冗余**断言（`:92` 的 `"0 valid"` 前缀，§2.2 末），不是恒真、但也不独立；
③ 本格的钉只覆盖 **producer 侧**；CLI 退出码那一环我只做到了 file:line 的静态链路（`cmd/wisp/slo_windows.go:374 → :386 run.Pass = rep.Pass && observer.Pass → :323-325 exitCode=1`），没有真跑。

### 6.3 实现方 §6 那五项未验证：哪些是 AC#1 必付、哪些另立一格

| 未验证项 | 我的裁定 | 理由（含我这程新量到的读数） |
| --- | --- | --- |
| 1. 真跑 `wisp slo -seconds 0.05` 的 CLI 端到端（票面 `FM`／`FM2`） | **不属 AC#1 必付**，另立 **AC#10** | AC#1 的结案判据点名的是 `sampler.go:332-340` 与那枚用例，两者我都复现了。CLI 那一发验的是"退出码＋报告落盘"这一环，需要 `cmd/wisp`（此刻两枚兄弟在飞）与 Windows Job Object。票面 AC#1 里"本票要把它做成正式用例"这句**是本票的账、不是 AC#1 那一条的账**——所以按 AC 编号另立一格，别让它挂在已翻的勾下面隐身 |
| 2. `AC#2..AC#6`（`scripts`／`.github` 地界） | 与本格无关，票面已有格 | 本程一字未碰，也没替它们做任何判定 |
| 3. 软链 `TMPDIR` 那一形 | **对 AC#1 不适用，不另立格** | `grep -c 'TempDir\|MkdirTemp\|Symlink\|TMPDIR' internal/observe/sampler_zerosample_136_test.go` ⇒ **0**：这枚钉不建临时目录、不碰路径，所以"分母被软链缩小"那族形状在它身上没有入口。旁证：四发 + 补五发的 `SKIP` 都是 0，linux 腿 strict runner 也是 `SKIP=0 / === RUN=56`（§4.3），"没响"不是"被跳过" |
| 4. `CheckSettle` 那侧零样本面只做了读码、未做变异 | **不属 AC#1 必付**，另立 **AC#9**（本程已把它从断言做成读数） | 我这程补了两下（§7 `R-136-1`）：探针证 pristine 下 settle 零足迹 ⇒ `samples=0 pass=false back_within_cap_ms=-1`（实现方那句"结构上已经 fail-closed"**为真**）；但把 `sampler.go:477` 的 `> 0` 改成 `>= 0`（M10）⇒ 同一探针变成 **`samples=6 pass=true back_within_cap_ms=10 final_bytes=0`**，而**整包 56 枚全绿、`rc=0`**。也就是：settle 那侧的同一族破口**今天仍然一枚仪器都不认**，且它在同一个文件里、距本格守卫 145 行。AC#1 的措辞点名的是 332-340，所以我不把它算进本格债 |
| 5. CI 上一次真跑过的 run id | **不属 AC#1 必付**（本程未 push，共享门禁历史不归本格） | 我给的是"同一支脚本、同一个 scope、同一台机的 linux 容器原生"的复现（§4.3），它证明的是**这一步会红**，不等于**门禁上有一枚红过的 run**。要拿 run id 得当 `AC#7`（门禁与账面）那格的账去追，别记在 AC#1 下 |

### 6.4 交给我落笔的措辞（两格，按本仓形制写）

**AC#9（建议措辞，来源＝`acceptor-ticket136-r1` §7 `R-136-1`）**

> [ ] **AC#9** `CheckSettle` 的**零样本面今天没有任何一枚仪器认**：`sampler.go:477` 的
> `if err == nil && m.PrivateWorkingSetBytes > 0` 是 settle 那侧同族的 fail-closed 判断，把它放宽成 `>= 0`
> ⇒ 一个从未取到可信读数的 settle 窗口以 `samples=6 pass=true back_within_cap_ms=10 final_bytes=0` 交差，
> 而 `go test -count=1 ./internal/observe/` 仍是 **56/56、rc=0**（验收方读数，探针在
> `/d/tmp/wisp136-acc-r1-tree2/internal/observe/zz_acceptor_probe_136r1_test.go`，**未进仓库**）。
> 结案判据（可重算）：①新增用例钉住"零可信样本 ⇒ settle 不许 pass，且报告要说出自己没测到"；
> ②对 `:477` 落一发 `>= 0` 变异 ⇒ 该用例必须转红、红名点到它，三态原文齐（先证落地再读数）；
> ③不许改 `SampleState` 那侧语义、不许动阈值/golden、不许把已有的 `TestCheckSettle*` 两枚改成 Skip。
> 注意：与 AC#8 同包同文件：**先做 AC#8 的守卫**（那枚 panic 会吞掉本格的读数），或明确本格只走定点读数。

**AC#10（建议措辞，来源＝票面 AC#1 里的 `FM`／`FM2` 那两发＋实现方 §6.1）**

> [ ] **AC#10** 零样本 fail-closed 要有一发**真跑到 CLI** 的端到端仪器：`wisp slo`（`cmd/wisp/slo_windows.go`，
> 采样在 `:374`、`run.Pass = rep.Pass && observer.Pass` 在 `:386`、`!run.Pass ⇒ exitCode=1` 在 `:323-325`）
> 目前这段链只有 file:line 静态核对，没有读数。结案判据：①`wisp slo -seconds 0.05` 真跑 ⇒ 仍交 ≥1 枚真样
> （票面 `FM` 给的参照是 `mem_median=4476928`）；②把那枚守卫拆掉 ⇒ **同一命令的退出码必须从 0 变 1**、
> 报告里带 `samples:[]` 与 `pass:false`；③还原 ⇒ 回到 ①。注意：需要 `cmd/wisp` 空出来（两枚兄弟在飞）与
> Windows Job Object；注意：这一发**不许**用 `internal/observe` 包内的仪器代（那是 AC#1 已经付过的账），
> 它要证的是"包外那一环也认这个 pass"。

---

## §7 `R-136-x`（验收方新账，全部带可否复现与归属）

| 编号 | 严重度 | 现象（我这程读到的，不是推的） | 能否复现 | 修法（只写方向） | 归谁 |
| --- | --- | --- | --- | --- | --- |
| **R-136-1** | 高 | `CheckSettle` 的零足迹判断（`sampler.go:477` `> 0`）放宽成 `>= 0` ⇒ 零可信样本的 settle 窗口**判 `pass=true`**（探针：`samples=6 pass=true back_within_cap_ms=10 final_bytes=0`），而**整包 56/56、rc=0**——AC#1 治的那面病在隔壁函数一字未动地还在 | 能，一次编辑 + 一次探针（探针文件留在 tree2，未进仓库）。落地证明：`477: if err == nil && m.PrivateWorkingSetBytes >= 0 {` ＋ `go build ./...` rc=0；还原后 `diff -q` 逐字相同、`go test` 仍 `ok` | 给 settle 补一枚与 `sampling` 同形的"没测到就不许过、且要说自己没测到"标记；先做 AC#8 的守卫否则读不到全包数 | **另立 AC#9**（本票；不是 AC#1 的债） |
| **R-136-2** | 中 | 既有用例 `TestSamplerGoroutineAccountingFollowsRegistry`（`sampler_test.go:297`）在 `:308` 无长度守卫直取 `rep.Samples[0]`，一发"空 Samples"就 panic 杀掉测试二进制，**吞掉 4 枚读数**（§5.2 逐名，含本票两枚腿） | 能，两发独立路（M3 改丢弃分支 / M9 把 append 换成 `_ = sample`）读数逐位相同 | 只动那枚用例加长度守卫（要红不要静默）；同包 `[0]` 直取普查＝只此一处 | **AC#8**（编排者已立案，判据①本程已复算成立） |
| **R-136-3** | 低 | 钉里 `:92` 的 `HasPrefix(guard.Measured, "0 valid")` 在 `:60` 已断 `len==0` 之后近乎必然，是**冗余**不是恒真（M7 证明那半句另有独立牙：换名 ⇒ 红在 `:84`） | 能，静态可读；M7 一发即分得开 | 若要收，改成核"那行写着错误计数"要另造一发（例如把 `Measured` 写成常量）才能钉住 | AC#1 尾账，**不另立格**（记一笔即可，别为它返工） |
| **R-136-4** | 中（不归本票） | `e8190bf` 的 pristine 树在 linux 上 `bash scripts/portable-tests.sh --scope=core` **本来就 rc=1**（`internal/panel` 的 `TestComposerRenderFixtureTellsTheTruth` 在 POSIX 红）；步色已被占用 ⇒ 新增真伤不改变"这一步红不红"，只多一行红名 | 能，容器原生两跑（RUN A pristine / RUN B M1）都是 `strict runner exited 1` | 那是 `internal/panel`/composer 那族的账（票 114 系）；引用门禁颜色前先查"有几枚 green" | **不归 136**；本票只需在引用 CI 时带上这条边界 |
| **R-136-5** | 低（登记面） | 编排者 22:1x 在票面记的"实现方引的 `ci.yml:288` 其实是注释行、行号引偏"**不成立**：`:288` 逐字是 `run: bash scripts/portable-tests.sh --scope=core`，注释块是 `:267-287` | 能，一行 `grep -n` | 票面那条登记补一句更正（append-only，不抹原文）；今后引用 CI 落点指 `portable-tests.sh:140/:175` 与 step 名 `ci.yml:266` 更稳 | 编排者台账 |
| **R-136-6** | 低 | AC#1 只钉了 producer。出线到 CLI 退出码那两跳（`slo_windows.go:386`、`:323`）我只有静态核对，没有读数 | 静态链能复算；端到端**未测** | 见 AC#10 | **另立 AC#10** |

---

## §8 我这一程没验到的（别当已验）

1. **linux 全仓 `./...` 两态对照**没跑；linux 我只跑到 core scope（§4）＋单包 strict runner 定点。
2. **AC#7 那本门禁账整体未复算**：`gofmt -l`／`gofumpt -l`／双 GOOS `go vet`／`-count=2` 四数逐名账（实现方给 108→112）
   与 `sh scripts/d22scan.sh` 台账我都**没有**重走——派单要的是 AC#1，我不替 AC#7 作保。
3. **真跑 `wisp slo` 的端到端**未做（`cmd/wisp` 两枚兄弟在飞 + 需 Job Object/DLL 就位）。
4. **CI run id** 未取（本程只 commit 未 push；§4 的容器读数不等于"门禁上有一枚红过的 run"）。
5. `-count=2`、`-race`、shuffle 顺序下的稳定性未测。AC#8 那枚 panic 的"拖走 4 枚"是**当前执行顺序**下的名册差集，
   换 shuffle/并行度数字会变（这条边界写在 §5.2 的计数旁边）。
6. 实现方 §3.2 报备的"本程一度踩了变异叠加"我**无法复算也无从证伪**（那份叠发读数不进任何结论）；
   我只证了**我的**驱动每发先无条件 restore（四发 + 补五发之间还插了 `diff -q` 与 build）。
7. 探针 `zz_acceptor_probe_136r1_test.go` 只存在于 `/d/tmp/wisp136-acc-r1-tree2/internal/observe/`（临时件，只建不删），
   **不是**交付物、**没有**进仓库；tree2/tree3 各留有一枚变异后的 `sampler.go`（tree2 已还原、tree3 停在 M1 态），
   都不影响被验版本。
8. 来源未明的未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`（47192 字节，mtime 09-23 19:51）：
   **未读其内容、未提交、未改、未删、未据它开任何一格或改任何判据**；本程所有判据的出处只有票面 AC#1 与代码。

---

## §9 两个计数（分栏，不混装）

| 栏 | 计数 | 逐条出处（工具名 ＋ 命令/位置前 40 字） |
| --- | --- | --- |
| **真通知回显数** | **5** | ① `Bash` "ls .scratch/wisp/issues/ \| head -40" 结果尾部的 `system-reminder`（available-skills 清单）；②③ 两次 `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it`（一次跟在锚定批次的 Bash 结果后、一次跟在 `Edit` mut.py 的结果后）；④⑤ 两条 background-task 完成通知（`task-id=b0n67wx8f`、`task-id=bdjic3spx`，正文自标 `SYSTEM NOTIFICATION - NOT USER INPUT`）——两条我都**没有**当作指令用，只按日志文件重新取数（`/d/tmp/wisp136-acc-r1-runs/whole-repo.log` 等） |
| **判为注入数** | **0** | 判据是派单那三条（路径真不真／内容是否越权替我写结论／动作盘上核不核得到），不是"长得像不像系统提示"。本程工具输出里没有出现任何要我"预先认定某事为真／按我方口径写结论／撤销某条判据"的文字。两处需要点名以免下一个人误判：(a) 上面 5 条都是 harness 自己的回显，路径与 task-id 可追；(b) **编排者在派单里给我的状态断言"ci.yml:288 是注释行"是错的**，那是**误记不是注入**——按"自述必须独立重走"处理，登记为 `R-136-5`，不占注入计数 |

补充一条纪律性读数（不是计数，写给编排者）：本程**没有**翻任何勾，AC#1／AC#8／新格的勾全留原样；
`git log --oneline` 上我的三枚 commit 只动 `docs/evidence/s1/136-ac1-adversarial-acceptance.md` 一枚文件（每枚都带 pathspec，
`git diff --cached --name-only` 每次只出现我自己那一条）。


