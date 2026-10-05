# 票 267 落地腿 `267-r1` — 证据件

对象票＝`.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish.md`
裁定形＝**ⓐ 值域进门**（编排者 2026-10-05 09:1x，锚点 HEAD `21bec8a1`）。
写面＝**只 `internal/config/**`**。

---

## §0 起手锚与并发

（本节答：起手时刻／HEAD／branch／并发腿／编排者 09:13 预取基线的引用与复核）

- 起手时刻 **`2026-10-05T09:15:10+0800`**，起手 HEAD＝**`f761a017`**，branch＝`dev`（尺＝`date "+%Y-%m-%dT%H:%M:%S%z"` + `git rev-parse --short HEAD`）。起手时工作树**已有别人的脏**（`git status --porcelain` 首行 `M .gitignore`，另有 `.scratch/wisp/probes/152|161/**` 若干 `M`）＝⛔ 不是本腿造的，本腿一律不动。
- 并发：**`167-a2`**（只读普查，不跑 go）、**`ci-phantom-1`**（只读，不跑 go）⇒ 本腿＝当机唯一跑 `go` 的腿；按派单⛔ **未跑** `go test ./cmd/wisp/`、⛔ **未跑** `wisp slo`。
- 编排者 09:13 预取基线（引用不重跑）：`go test ./internal/config/` ＝ **ok 1.008s**（09:13:02，全绿）／`gofumpt -l internal/config/` ＝ 空（09:13）／`sh scripts/d22scan.sh` ＝ rc=0 `clean - no D22 ban violations`（08:56:34→08:57:05，分母 `ban #8 internal/=507`、`cmd/=101`）／`sh scripts/check-path-length-budget.sh --with-self-test` ＝ VERDICT GREEN（08:57，over-budget 57＝roster 57）。
- **复核一次**（本腿自己取数，见 §6 带时刻）：`gofumpt -l internal/config/` 于 `09:24:16` 复跑＝仍空；`go test -count=1 ./internal/config/` 于 `09:24:16` 复跑＝`ok 1.001s` ⇒ 与 09:13 基线同形（全绿），红名册＝空。
- 本腿第一次 commit＝**`42115f65`**（证据件骨架，pathspec 只有 `evidence.md`）。


---

## §1 写面与形状

（本节答：validate 那一段逐字设计 + 同源守卫那一格怎么落的 + 依赖方向实测尺）

**写面＝5 枚，全在 `internal/config/**`**（名册见 §8）：`validate.go`（带＋检查＋推导注释）、`validate_test.go`（两枚新钉）、`validate_267_test.go`（**新文件**，外部测试包：同源守卫＋加载路径带＋AC#2 读数）、`schema.go`（仅注释面）、`unwired.go`（仅 `:121` 一行释文，见 §5）。

### 1.1 落进 `validateRisk` 的那一段（逐字）

```go
	if sec := c.Risk.ConfirmTimeoutSec; sec < confirmTimeoutSecMin || sec > confirmTimeoutSecMax {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]",
			sec, confirmTimeoutSecMin, confirmTimeoutSecMax))
	}
```

常量 `confirmTimeoutSecMin = 31` / `confirmTimeoutSecMax = 3600`；注释把"为什么是 31"写成**不变式**（`gate.go:528` 只在 `Timeout()-WarningLead() > 0` 时武装提示；`queue.go:88-89` 在生产者没传 lead 时恒兜底成 30s；生产里**没有**生产者传 lead＝`267-a1` §1 实测）。错误句形状照抄既有先例 `ball.size %d out of range [44, 72]`＝**同层同形状**（裁定理由②"复用现成形状"），语义＝拒载（`observe.ClassConfig`），⛔ 不是"程序自己决定无视配置"。`permission_mode` 那枚枚举检查一字未动，新检查排在它后面。

### 1.2 ★ 那枚"我自己先判的依赖问题"——**先跑尺，没猜，落成了**

1. **尺①** `grep -rn "internal/config" internal/agent/approval/` ⇒ **rc=1，0 命中**＝`approval` 不**直接**指 `config`。**只看这一步会误判成"可以随便 import"**。
2. **尺②（真判据）** `go list -deps ./internal/agent/approval | grep wisp/internal/config` ⇒ **命中**。直插 `config` 的两位＝`internal/llm`、`internal/agent`；链条＝`approval -> internal/agent -> internal/config`（名册＝`.scratch/wisp/probes/267/r1/approval-deps.txt`）。⇒ **反向依赖存在** ⇒ 在 `package config`（内部测试包）里 import `approval` **就是 import cycle**：这不是推理，本腿第一版就是这么写的，`go vet` 当场报死。
3. **落法＝派单 ② 的可行变体**：守卫落在新文件 `internal/config/validate_267_test.go`，包名 **`package config_test`（外部测试包）**。Go 允许外部测试包 import 那些依赖被测包的包——而"`config` 自己不依赖 `approval`"这件事**正是本票要守的不变式**，文件头注释写明了这层关系。⇒ **没有常量值复制**：`warningLeadSec()` 返回 `int(approval.DefaultApprovalWarning / time.Second)`，守卫断的是「`lead` 值必被拒、`lead+1` 值必能载」两枚**比较**，且走真 `config.LoadFile`。⛔ 不需要交回你判。
4. **反向自证**（09:34:42，`final-verify.txt` (1)(2)(3)）：`go list -deps ./internal/config | grep wisp/internal/agent` ＝ **0** ⇒ 产码面没新出指向 `internal/agent/approval` 的边；`go list -test -deps ./internal/config | grep .../approval` ＝ **1** ⇒ 该边**只在测试二进制里**。`internal/config/*.go` 中 "internal/agent" 的 3 处命中**全是注释**（`validate.go:109/117/118`），`validate.go` 的 import 块逐字未变（`fmt/slices/strings` + `observe/risk/secret`）。

### 1.3 一枚自纠（已写进产码注释）

第一版 `validate_test.go` 把种子写成 `confirmTimeoutSecMax + 1`／`confirmTimeoutSecMin - 1`。**突变 M2 当场把它打回原形**：抬上界时种子跟着搬家（3601 -> 360001，仍然"合法地被拒"），于是上界那一支**测不到**＝正是票面 AC#1 警告的"随便换个越界值就不响"。⇒ 现在两枚钉的种子与期望句**全是字面量**（`30/10/3601/99999/0/-1` 与 `"out of range [31, 3600]"`），注释具名记了这次教训（commit `ec6a47a8`）。符号派生的那一支没丢，它搬去由符号说话的地方＝`validate_267_test.go` 的同源守卫。


---

## §2 AC#1 四种子读数

（本节答：30／10／3601／99999 四发越界种子各自的指名用例与红句逐字）

**指名用例（本腿新增；绿色＝带在位时确实逐枚拒载，带摘掉时逐枚转红＝§4）**

| 层 | 用例 | 越界种子 | 合法种子 |
|---|---|---|---|
| 语义层 `validate()` | `internal/config/validate_test.go TestValidateRiskConfirmTimeoutRange` | `reject_30`／`reject_10`／`reject_3601`／`reject_99999`／`reject_0`／`reject_-1` | `accept_31`／`accept_300`／`accept_3600` |
| 语义层（带边语义） | 同文件 `TestValidateRiskConfirmTimeoutBandRejectsWarningLeadPlusOne` | 30（＝lead，必须拒） | 31（＝lead+1，必须载） |
| **加载层**（真 `LoadFile`＝拒载语义） | `internal/config/validate_267_test.go TestConfirmTimeoutIsBandedOnTheLoadPath` | 30／10／**1／2／20**／3601／99999／0／−1 | 31／300／3600 |
| 加载层（同源） | 同文件 `TestConfirmTimeoutFloorSitsAboveApprovalWarningLead` | `lead`（符号现算＝30） | `lead+1`（＝31） |

**程序实际吐出的拒载句**：形状钉在 `validate.go` 的格式串
`"config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]"`
⇒ 种 99999 那句＝`config.toml: risk.confirm_timeout_sec 99999 out of range [31, 3600]`。**"哪枚键／超出哪个范围／越界值本身"三件事都由断言钉着**（不是本件手抄）：`reject_99999` 为绿即代表 `err.Error()` 同时含 `risk.confirm_timeout_sec`、`out of range [31, 3600]`、`99999 out of range`。

**加载层逐枚拒载读数**（09:33:39，`.scratch/wisp/probes/267/r1/seeds-cmd-proof.txt`）：
`reject_30 PASS`｜`reject_10 PASS`｜`reject_1 PASS`｜`reject_2 PASS`｜`reject_20 PASS`｜`reject_3601 PASS`｜`reject_99999 PASS`｜`reject_0 PASS`｜`reject_-1 PASS`｜`accept_31 PASS`｜`accept_300 PASS`｜`accept_3600 PASS`。
（`1`／`2`／`20` 不是编的数＝`cmd/wisp` 今天真的在种的种子，后果见 §7.1。）

**正控方向**（"摘掉那枚界限检查必须让指名用例红"）＝§4 的 M1／M2／M3，红句逐字在那一节，⛔ 不在这里复述。


---

## §3 AC#2 存活读数

（本节答：合法带内最小值那一发的 `Timeout() - WarningLead() > 0` 读数，或写面外那半格的具名归口）

**这一发在写面里做到了，而且不是近似**——它构造的是**真的 `approval.Queue`**（⛔ 不是 mock），喂法逐字照装配根（`cmd/wisp/run.go:616` 与 `cmd/wisp/resident_approval_windows.go:460` 都是 `time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second` 且**不传** warning lead）。

用例＝`internal/config/validate_267_test.go TestMinimumLegalConfirmTimeoutKeepsC18Warning`，读数逐字（09:24:04）：

```
    validate_267_test.go:127: smallest legal timeout 31s, warning lead 30s, arming lead 1s > 0
--- PASS: TestMinimumLegalConfirmTimeoutKeepsC18Warning (0.00s)
```

即带内最小值 31s -> `q.Timeout()=31s`、`q.WarningLead()=30s`（`queue.go:88-89` 的兜底，因为生产者不传）、`gate.go:528` 分支所依据的那个表达式 `Timeout()-WarningLead() = 1s > 0` ⇒ **C18 那句提示仍武装**。配上 §2 的"种 30 必被拒"，两半合起来才是 AC#2 要的结论：**"静默把保护撤了"这一终态在配置层不可达**。

⚠ **剩下那一半我确实够不着，具名归口＝§7.2**：真让 `warn` channel 响一次、看到面板那句提示到达，要构造 gate 对象＋clock＋prompt，落点在 `internal/agent/approval/**` 与 `cmd/wisp/**`（两枚都是本腿禁改面）。⛔ 本腿没越界写。


---

## §4 突变自证

（本节答：摘掉界限检查 ⇒ 指名用例必红的每一发：原样命令 + 红句逐字 + 还原出处 + 三枚 md5）

**统一还原出处**＝`git cat-file blob HEAD:internal/config/validate.go > internal/config/validate.go`（⛔ 全程没用 `checkout`／`stash`／`restore`／`reset`）。
**md5 锚**＝起手＝HEAD blob＝还原后，三者皆 **`fecfefcb1c4240801fb2f9c47bcd66f9`**（`validate.go` 在 `37f8e5c6`→`ec6a47a8`→`873c3063` 之间是同一枚 blob，只有测试文件在动，所以三发突变共用这一枚锚）。原始输出＝同目录 `m1-red.log`／`m2-red.log`／`m2b-red.log`／`m3-red.log`（已 commit）。

### M1 — 把那段界限检查整个摘掉（票面点名的正控）

- 改法：删掉 `validateRisk` 里那 5 行 `if sec := ...` 块，换成一行 `// MUTATION-M1: band check removed on purpose (ticket 267 r1 self-proof).`
- 三枚 md5：起手 `fecfefcb1c4240801fb2f9c47bcd66f9` → 突变 **`2a39e278756683ef97db414ccfecf6f0`** → 还原 **`fecfefcb1c4240801fb2f9c47bcd66f9`**（逐字节＝起手；另用 `grep -c MUTATION` ＝ 0 复核残留）
- 原样命令：`go test ./internal/config/ -count=1 -v -run 'TestValidateRiskConfirmTimeoutRange|TestConfirmTimeoutIsBandedOnTheLoadPath'`
- 红句逐字（语义层六枚**全红**，`--- FAIL: TestValidateRiskConfirmTimeoutRange/reject_30 … reject_-1` 六行子测试名册在 `m1-red.log`）：
  ```
  validate_test.go:196: risk.confirm_timeout_sec=30 must be rejected (out of range [31, 3600])
  validate_test.go:196: risk.confirm_timeout_sec=10 must be rejected (out of range [31, 3600])
  validate_test.go:196: risk.confirm_timeout_sec=3601 must be rejected (out of range [31, 3600])
  validate_test.go:196: risk.confirm_timeout_sec=99999 must be rejected (out of range [31, 3600])
  validate_test.go:196: risk.confirm_timeout_sec=0 must be rejected (out of range [31, 3600])
  validate_test.go:196: risk.confirm_timeout_sec=-1 must be rejected (out of range [31, 3600])
  ```
- 红句逐字（加载层同一命令的另一半）：
  ```
  validate_267_test.go:94: confirm_timeout_sec = 30 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 10 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 3601 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 99999 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 0 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = -1 must be refused at load
  ```
  ⚠ 行号口径：M1 跑在"字面量改造之前"的那一版测试文件上，所以句里的行号是当时的（`validate_test.go:196`／`validate_267_test.go:94`）。同一发在改造后又跑了一次，语义层句变成 `validate_test.go:200`（见 M2b/M3 的引用）。**行号只对写下那一刻的树负责**＝票面 09:13 那条定式。
- 补跑两枚守卫钉（`go test ./internal/config/ -count=1 -run 'TestValidateRiskConfirmTimeoutBandRejectsWarningLeadPlusOne|TestConfirmTimeoutFloorSitsAboveApprovalWarningLead|TestMinimumLegal' -v`）红句逐字：
  ```
  validate_test.go:236: a timeout equal to the C18 warning lead (30s) must be rejected: it arms no warning at all
  validate_267_test.go:68: confirm_timeout_sec = 30 must be refused: it equals approval.DefaultApprovalWarning (30s), so gate.go's `lead := Timeout() - WarningLead()` is <= 0 and the C18 warning never arms
  ```
- ⚠ **如实登记一处不敏感**：M1 下 `TestMinimumLegalConfirmTimeoutKeepsC18Warning` **仍 PASS**（照样打出 `smallest legal timeout 31s, warning lead 30s, arming lead 1s > 0`）。原因不是它坏了：它断的是**正向存活**（31 这一发提示还在），摘不摘带都不影响那步算术。它的敏感面由上面那些**拒载钉**承担。把这条写出来是因为"摘掉检查还全绿＝不敏感"是本票明令不许蒙的形状，我不替自己圆。

### M2 — 只松上界（证明上界**单独**可打红，不是搭 M1 的车）

- 改法：`confirmTimeoutSecMax = 3600` -> `= 360000`（一行）
- 三枚 md5：起手 `fecfefcb1c4240801fb2f9c47bcd66f9` → 突变 **`8636660c0e4bd3fcfd344b0c42b8049d`** → 还原 `fecfefcb1c4240801fb2f9c47bcd66f9`
- 原样命令：同 M1
- 红句逐字（上界两支，两层各一）：
  ```
  validate_test.go:200: risk.confirm_timeout_sec=3601 must be rejected (out of range [31, 3600])
  validate_test.go:200: risk.confirm_timeout_sec=99999 must be rejected (out of range [31, 3600])
  validate_267_test.go:94: confirm_timeout_sec = 3601 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 99999 must be refused at load
  ```
- **干净分离读数在加载层**：那一枚只断"含键名"＋"含 out of range"、不断具体数，所以 M2 下 `reject_30`／`reject_10`／`reject_0`／`reject_-1` 在加载层**保持 PASS**（`m2b-red.log` 逐字 `--- PASS: TestConfirmTimeoutIsBandedOnTheLoadPath/reject_30` 等）＝**松上界只打红上界**。语义层另有 30/10/0/-1 转红，那是期望句 `[31, 3600]` 与实吐 `[31, 360000]` 不符＝带文本漂移也算打红，属预期的过敏，不是误判。
- ⚠ **第一发 M2 作废并如实登记**（`m2-red.log`，突变 md5 `7e3703fa61560290d7685df7acbc7b8e`）：那一发跑在字面量改造**之前**，旧种子 `confirmTimeoutSecMax+1` 跟着搬到 360001、仍然"合法地被拒"，于是上界那一支**假绿**——正是 §1.3 那枚自纠的来源。旧日志只建不删，留在原处。

### M3 — 只松下界（证明下界单独可打红，且**同源守卫真的响**）

- 改法：`confirmTimeoutSecMin = 31` -> `= 1`（一行）
- 三枚 md5：起手 `fecfefcb1c4240801fb2f9c47bcd66f9` → 突变 **`b502f6b507a9db5c768e3c7847ac3a34`** → 还原 `fecfefcb1c4240801fb2f9c47bcd66f9`
- 原样命令：`go test ./internal/config/ -count=1 -v -run 'TestValidateRiskConfirmTimeoutRange|TestConfirmTimeout|TestMinimumLegal'`
- 红句逐字：
  ```
  validate_test.go:200: risk.confirm_timeout_sec=30 must be rejected (out of range [31, 3600])
  validate_test.go:200: risk.confirm_timeout_sec=10 must be rejected (out of range [31, 3600])
  validate_267_test.go:94: confirm_timeout_sec = 30 must be refused at load
  validate_267_test.go:94: confirm_timeout_sec = 10 must be refused at load
  validate_267_test.go:68: confirm_timeout_sec = 30 must be refused: it equals approval.DefaultApprovalWarning (30s), so gate.go's `lead := Timeout() - WarningLead()` is <= 0 and the C18 warning never arms
  ```
- **这一发的分量在最后一行**：把下界压到符号以下，红的不是"我抄的那个 31"，而是**由 `approval.DefaultApprovalWarning` 现算出来的那条断言**＝同源守卫确实能看见两边脱钩（＝派单 ② 那一格要的东西，不是装饰）。
- 三发跑完的最终树：`go test -count=1 ./internal/config/` ＝ **`ok 2.443s`**（09:34:42）、`gofumpt -l internal/config/` ＝ 空、`md5sum internal/config/validate.go` ＝ `fecfefcb1c4240801fb2f9c47bcd66f9`（＝HEAD blob）、`git diff --name-only internal/config/` 里 `validate.go` 不出现 ⇒ **零突变残留**。


---

## §5 `unwired.go` 释文改写前后逐字

（本节答：过期指认的原文、新文、理由出处）

改前（`internal/config/unwired.go:121`，逐字）：

```go
	"risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go builds the approval timeout from it",
```

改后（同处，逐字）：

```go
	"risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go and cmd/wisp/resident_approval_windows.go both build the approval timeout from it (the resident leg since ticket 256); ticket 267 bands it [31, 3600] at load",
```

- 改这条的出处＝票 267 §排程 那一行（"只说了 `run.go`——票 256 落地后**常驻腿也消费它**了 ⇒ 落地腿同批改释，⛔ 不单独立票，见 `A603`"）＋派单写面里点名的 `unwired.go:121-122`。
- 现跑的尺（按票面 09:13 那条自纠定式：**用被指的字符去 grep，不抄行号**）：`grep -n 'ConfirmTimeoutSec) \* time.Second' cmd/wisp/resident_approval_windows.go cmd/wisp/run.go` ⇒ **`resident_approval_windows.go:460`** ＋ **`run.go:616`** 两枚消费点。⇒ 新句把两位都点了名并标出"常驻腿自票 256 起"的来历，顺带把本票新加的带写进指认（"bands it [31, 3600] at load"），这样下一个人不会再把这枚键当成"能配什么随你"。
- `:122` 的 `l1_window_sec` 那行**故意没改**——它的过期方式与 `:121` 一字不差（常驻腿也读它，见 `resident_approval_windows.go:382` 的 `window_sec_read`），但票面与派单都写着 `l1_window_sec` 一字不动 ⇒ ⛔ 我没把它读成"顺带改注释"，具名交回＝**§7.3**。
- 这行属 `lockedKeyDisposition` 表，归 `unwired_test.go TestEveryLockedSectionKeyIsAccountedFor` 那枚完整性钉管：值若非 `unwired:<path>` 形状就是"今天允许静默加载"的说明文字。本腿改的是**说明**，键名与分类前缀都没动 ⇒ 该钉与全包同绿（§6）。


---

## §6 门禁四数

（本节答：go test / gofumpt / d22scan / path-length 四门读数，各带取数时刻）

取数时刻 **09:31:18–09:31:49 +08**，当时 HEAD＝**`d7a10152`**（⚠ 起手是 `f761a017`，期间别的腿提交了 commit＝共享树正常现象）。完整原样输出＝`.scratch/wisp/probes/267/r1/gates-final.txt`。

| 门 | 读数（带时刻） | 与 09:13/08:57 基线比 |
|---|---|---|
| `go test -count=1 ./internal/config/` | **`ok 1.057s`**（09:31:19），全绿 | 基线 09:13:02 `ok 1.008s` 全绿 ⇒ **红名册＝空**，未扩大 |
| `gofumpt -l internal/config/` | **空**（09:31:19；另 09:24:16 也空） | 基线 09:13 空 ⇒ 未扩大 |
| `sh scripts/d22scan.sh` | **rc=0 `clean - no D22 ban violations`**（09:31:21） | 违规＝0 不变。分母 **`ban #8 internal/=508`（基线 507，+1）**＝本腿新增 `internal/config/validate_267_test.go` 这一枚 `internal/` 下的 `.go` 文件被 ban #8 计入（该尺连注释与 `_test.go` 一起扫），**不是新增违规**；`cmd/=101` 与基线逐字相同 ⇒ 本腿未碰 `cmd/` |
| `sh scripts/check-path-length-budget.sh --with-self-test` | **`VERDICT GREEN`**（09:31:44），`tracked paths=5701`、`over-budget=57`、`covered by roster=57`、`not in roster=0`、`in the wall interval=0`、`past the old debt line(180)=0` | 基线 08:57 over-budget 57＝roster 57、wall 区间 0 ⇒ **一枚没多**（本腿最长新路径 `internal/config/validate_267_test.go`＝37 字符，远未过 hat） |

⚠ 同一次 `git diff --name-only`（09:31:18）里出现 `.gitignore`、`.scratch/wisp/probes/152|161/**`、`design/**` 共 30+ 枚 —— 那**不是本腿的改动**，是起手 09:15 `git status` 就在工作树里的别人的脏（§0 记了首屏）。本腿自己的改动一律以 §8 的 `git show --name-only` 逐枚列，只有 `internal/config/**` 与 `.scratch/wisp/probes/267/r1/**`。

**终值复跑**（09:34:42–09:35:20，HEAD＝`873c3063`，⛔ 不是我 §6 上表那次读数之后树没动——期间别的腿又提交了，HEAD 见文末；完整输出＝`.scratch/wisp/probes/267/r1/final-verify.txt`）：`go test -count=1 ./internal/config/` ＝ **`ok 2.443s`** 全绿／`gofumpt -l internal/config/` ＝ **空**／`sh scripts/d22scan.sh` ＝ **`clean - no D22 ban violations`**（分母仍 `internal/=508`、`cmd/=101`，与上表逐字同）／`check-path-length-budget.sh --with-self-test` ＝ **`VERDICT GREEN`**、`over-budget=57`、`roster=57`、`not in roster=0`、`wall interval=0`。⇒ 四门**取的是最后一次落地后的树**，与交件树一致（§8 的 `git status --porcelain internal/config/` 为空可复核）。


---

## §7 判不动／量不到

（本节答：具名空格 + 归口，⛔ 不许"应该没问题"）

### 7.1 ★ 阻塞级：**[31, 3600] 会打红 `cmd/wisp` 里现成的绿用例**（⛔ 本腿无权修，且**任何**满足票面不变式的下界都躲不开）

裁定要求下界**严格大于** `DefaultApprovalWarning`＝30s ⇒ 下界 ≥31。**不是"31 挑坏了"**：**任何** ≥3 的下界都会拒掉下面这些种子（要不动它们只能把下界压到 ≤2，那直接违反不变式）。尺（`grep -rn "confirm_timeout_sec\|ConfirmTimeoutSec"` 全仓，排除 `.scratch/`）拉出的种子名册：

| 种子出处 | 值 | 落带 | 机制（为什么会被拒） |
|---|---|---|---|
| `cmd/wisp/run_mode101_test.go:110` | `confirm_timeout_sec = 1` | **带外** | 同文件 `:172` 是 `config.NewManager(h.configPath(), nil)`，走 `LoadFile`→`readConfigFile`（`internal/config/loader.go:66-67` 逐字："runs the sec 4.1 pipeline up to (and including) semantic validation"）⇒ 加载即拒 |
| `cmd/wisp/panel_pump_test.go:136` | `confirm_timeout_sec = 2` | **带外** | 写进 fixture 的 `config.toml`，由 run 腿加载 |
| `cmd/wisp/approval_reply_201_test.go:145` | `int(h.l2Wait.Seconds())` ⇒ **20**（`:202`/`:272`/`:350`/`:504`/`:578`）与 **2**（`:439`） | **带外** | 同一 NewManager/LoadFile 链 |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:202` 等 | 45 / 90 | 带内 ✓ | 票面"为什么不是 60s"那枚钉，不受影响 |
| `internal/config/unwired_test.go:66`／`:188`、`boundary_test.go:142` | 300 / 60 | 带内 ✓ | 本包自证，§6 全绿 |

⇒ **归口＝编排者**：这三枚文件（`cmd/wisp/**`）在票 256/201/145 那批钉的地界里，本腿⛔一字未动、也按派单**没跑** `go test ./cmd/wisp/`（那一包的基线你说另有安排）。所以这一条给的是**机制读数**（种子值 + 它走的加载链都在表里逐字可查），不是"我跑过它红"。要落 ⓐ 就得配套解冻这 8 枚调用点的种子（45s 那枚你已具名，这 3 枚文件＝新的具名项）。

⚠ 顺带把**根因**说清，别让它被当成"测试写死了烂值"：这些用例种 1/2/20 秒是**故意的**——它们要在一次测试运行里亲眼看着一张 L2 卡开合并超时（`panel_pump_test.go:124` 逐字："so the case can watch an L2 card open and close"）。值域进门之后，**测试没有合法途径表达"两秒的卡"**：`confirm_timeout_sec` 是整秒、下界 31。这一格是设计后果，不是 bug，⛔ 不由本腿裁（可选出口：给测试留一条构造 `approval.NewQueue` 的窄缝、或把 C18 提示提前量做成可配、或接受 cmd 腿改用 fake clock——三选一归你）。

### 7.2 AC#2 的**另一半**在写面外（具名归口）

本腿能做到的是**真对象读数**（见 §3）：真 `approval.Queue` + 真 `Timeout()/WarningLead()`。做不到的是驱动 `gate.go:528` 那枚 `warn` channel 真的响一次——那是 gate 对象＋clock＋prompt 的构造，落点在 `internal/agent/approval/**`（⛔ 禁改）或 `cmd/wisp`（⛔ 禁改）。⇒ **归口＝编排者排给哪枚腿**：建议下一枚 `cmd/wisp` 写腿顺手一发"种 31 ⇒ 提示到达"的端到端读数（它那批文件本来就要动，见 7.1）。

### 7.3 `unwired.go:122` 那枚**同样过期**的指认（本腿故意没改）

`"risk.l1_window_sec": "consumed: cmd/wisp/run.go builds the L1 auto-approve window from it"` —— 票 256 之后常驻腿也读它（`resident_approval_windows.go:382` 的 `window_sec_read` 字段），过期方式与 `:121` 一字不差。**但票面与派单都写着 `l1_window_sec` 一字不动**，本腿不敢把它读成"顺带把 :122 的注释也改了"⇒ **只改 :121**（§5），:122 交回你：要么并进 7.1 那枚 cmd 腿，要么具名一枚释文票。

### 7.4 裁定里列成残余的格子，本腿一律没做（逐条确认）

① 下界抬到 60s＝要动 256 的种子钉 ⇒ 未做；② `confirm_timeout_sec` 拒载 vs `l1_window_sec` 钳位的行为不一致 ⇒ 未碰；③ `WarningLead` 生产里恒零、靠 `queue.go:88-89` 兜底 ⇒ 未碰（**本腿的下界正是**依赖这个事实才成立的）；④ `llm.timeout_ms` 等 6 枚无上界数值键 ⇒ 未碰。Q-77（配置值 vs C18 写死 300s 谁优先）⇒ 未表态，schema `default:"300"` 与 `DefaultApprovalTimeout` 一字未动（AC#3，见 §8 自证）。

### 7.5 量不到（说清为什么量不到，不写成"应该没问题"）

- `go test ./cmd/wisp/` ＝派单明令禁跑 ⇒ 7.1 的后果**未被本腿实测**，只给了机制链。
- `wisp slo`／D32 采样 ＝派单明令禁跑（self-hosted runner 刚被用过）⇒ 未取。
- CI 端 phantom-citation `cmd/wisp/panel_host_windows.go:343` ＝别的只读腿的地界 ⇒ 未读未跑。


---

## §8 污染面与提交名册

（本节答：本腿全部 commit + `git diff --name-only` 名单 + 票面 AC 框零改动自证）

### 8.1 提交名册（只 commit，⛔ 全程没 push；每枚都带显式 pathspec，⛔ 没有 `git add -A`/`.`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`restore`）

| commit | 内容 | 带的 pathspec |
|---|---|---|
| `42115f65` | 证据件骨架（§0–§8，八节正文＝占位符，本件现已全部填实） | `evidence.md` |
| `37f8e5c6` | 落地：带＋检查＋注释、两枚语义层钉、新外部测试文件、schema/unwired 注释面 | `internal/config/{validate.go,validate_test.go,validate_267_test.go,schema.go,unwired.go}` ＋ 本件 3 枚临时件 |
| `ec6a47a8` | 自纠：钉改成字面量种子（M2 暴露的问题） | `internal/config/validate_test.go` |
| `22fc968a` | 证据 §0/§6/§7 ＋ 四发突变日志 ＋ 门禁原样输出 | 本件 ＋ `gates-final.txt`、`m1/m2/m2b/m3-red.log` |
| `873c3063` | 外部守卫补种 `cmd/wisp` 今天真在种的 1/2/20 ⇒ 拒载在加载层逐枚实测 | `internal/config/validate_267_test.go`、`seeds-cmd-proof.txt` |
| （本件最后一枚＝§1–§5/§6 复跑/§8 落字 ＋ `final-verify.txt`） | 证据件收口 | `evidence.md`、`final-verify.txt` |

### 8.2 污染面自证（AC#4）

- **本腿 5 枚 commit 的全部文件名册**（`git show --name-only` 去重后，`final-verify.txt` (10)）＝
  `internal/config/schema.go`、`internal/config/unwired.go`、`internal/config/validate.go`、`internal/config/validate_267_test.go`、`internal/config/validate_test.go` ＋ `.scratch/wisp/probes/267/r1/**`（证据件、四份突变日志、门禁输出、种子名册、两枚 commit-msg 临时件）。
- **尺**：`for c in <上述 5 枚>; do git show --name-only --pretty=format: $c; done | sort -u | grep -v -e '^internal/config/' -e '^\.scratch/wisp/probes/267/r1/' | wc -l` ＝ **0** ⇒ 本腿的 commit 里**没有一枚**文件落在写面之外。
- **禁路径逐名比对**（`final-verify.txt` (11)，尺＝`git diff --name-only f761a017..HEAD -- docs/PLAN.md docs/specs internal/observe/thresholds.go tools/d22scan/allowlist.txt frontend design internal/agent cmd | wc -l`）＝ **0**。⚠ 这一把的范围**还包含了别的腿**（`f761a017..HEAD` 里有 `d7a10152`、`8aa9ff6f` 等别人提交）⇒ 整个区间都没人动过禁路径，比"只有我干净"更强。
- **AC#3**：`DefaultApprovalTimeout`／`DefaultApprovalWarning` 两个常量的定义行、`queue.go`、`gate.go` **一行未动**（上面的 `internal/agent` diff＝0 就是这一条的尺）；`schema.go` 的 `default:"300"` 逐字仍在（我只在它**上面**加了注释）；`docs/PLAN.md` 的 C18 面＝禁路径且 diff＝0。⛔ 没把 300s 改小、没把默认挪走。
- **`l1_window_sec`**：值域／钳位／`unwired.go:122` 释文**三处全部一字未动**（§7.3 具名交回）。
- **工作树里别人的脏**（`.gitignore`、`design/**`、`.scratch/wisp/probes/152|161/**`）＝起手就在（§0），本腿**一枚没碰、也没顺手修**。

### 8.3 票面 AC 框零改动自证

- **尺①（ scoped 到本票）**：`git status --porcelain -- .scratch/wisp/issues/267*` ＝ **空**（09:41:21）⇒ 票 267 文件在工作树里没被我改过一个字（progress log 归编排者）。
- **尺②（谁最后碰过它）**：`git log --oneline -3 -- .scratch/wisp/issues/267*` ＝ **`f761a017`**，那是**编排者自己**翻 AC#0＋裁形的那枚 commit ⇒ 本腿之后没有第二枚。
- ⚠ 同一次 `git status --porcelain -- .scratch/wisp/issues/` 里另有 `M .scratch/wisp/issues/167-...md` ＝**别人的**（`167-a2` 那枚普查腿在飞的票），本腿没碰、也不去"顺手清理"。
- 票面 AC 框现状（`grep -c "^- \[ \] \*\*AC" .scratch/wisp/issues/267*.md`）＝ **4** ＝ AC#1–AC#4 四枚**全部还是未勾**，AC#0 保持编排者翻好的 `[x]` ⇒ 本腿没有自翻勾。
- 证据件本身：`grep -c "未判|待填|填写中"` ＝ **0**（09:41:21），290 行，八节全部填实。


### 8.4 交件形状（一句话）

只落了裁定的 **ⓐ 值域进门**：`[risk].confirm_timeout_sec` 在 `internal/config/validate.go` 获得带 `[31, 3600]`，越界＝`observe.ClassConfig` 拒载并逐字念出键名、越界值与范围；下界的不变式（严格大于 `approval.DefaultApprovalWarning`）由 `package config_test` 里的**同源守卫**从符号现算钉住，产码面**没有**新出 `internal/config -> internal/agent/approval` 的依赖边。**唯一需要你处理的成块后果＝§7.1**（带会拒掉 `cmd/wisp` 三枚文件里 8 个调用点现种的 1/2/20/30 秒；这不是 31 选错了，是任何满足不变式的下界都会撞上，需要具名解冻或给测试留一条合法的"短卡"表达口）。

