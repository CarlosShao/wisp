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

