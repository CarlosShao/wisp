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

未判

---

## §2 AC#1 四种子读数

（本节答：30／10／3601／99999 四发越界种子各自的指名用例与红句逐字）

未判

---

## §3 AC#2 存活读数

（本节答：合法带内最小值那一发的 `Timeout() - WarningLead() > 0` 读数，或写面外那半格的具名归口）

未判

---

## §4 突变自证

（本节答：摘掉界限检查 ⇒ 指名用例必红的每一发：原样命令 + 红句逐字 + 还原出处 + 三枚 md5）

未判

---

## §5 `unwired.go` 释文改写前后逐字

（本节答：过期指认的原文、新文、理由出处）

未判

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

未判
