# 票 321 — 托盘/球/热键那一跳的执行器接线判据（落地腿 `321-r1` 交付件）

**腿**：`321-r1`（落地腿，只装判据）· **票**：`.scratch/wisp/issues/321-*.md`（AC#1/AC#2 框**未翻**，⛔ 本腿不自裁）
**锚**：`dev` @ `78ea61e2f8395aac66e80991b762ad3ce526f356`（开工前 `git log -1` 现读，与派单锚同枚）
**时刻**：开工 2026-10-15 19:35 +08，交件 2026-10-15 19:59 +08（全部 `date` 插值，⛔ 手打）
**裁决表归属**：AC#1/AC#2 的成立与否**必须由非实现者另出**（`AGENTS.md` §0.3 / `SPEC-12 §4.3`#1#3）。本件只是实现方读数。

---

## 1. 现量锚逐枚复核（派单 §2 要求"引用前自己再跑一遍"）

| 锚 | 派单给的 | 本腿现跑 | 判 |
|---|---|---|---|
| 装配处 | `cmd/wisp/resident_windows.go:268-269` | `:268` `...withPanelHost(func(via string) bool {` ／ `:269` `return holder.toggle(via)` | 对 |
| holder 构建 | — | `:201` `holder := newPanelHolder(rt.Registry, rp)`；构造器 `:501` | 新增引用 |
| 两个执行器 | — | `toggle` `:536` ／ `show` `:595`（同文件） | 新增引用 |
| hook 定义 | — | `cmd/wisp/resident_ball_windows.go:107` `func withPanelHost(showPanel func(via string) bool) ballHostHook` | 新增引用 |
| 反面尺 | `cmd/wisp/panel_resident_windows_test.go:1009` | `:1009 func TestBallPanelGesturesReachThePanelThread`，体内 `hs := panelHostHooks{showPanel: rp.RequestToggle}`＝**自备夹具、不读产码接线** | 对 |
| AST 先例 | `resident_ball_228_test.go:261` 族 | 该文件 `:64 productionFiles228`／`:132 funcBody228`；`SelectorExpr` 计数钉在 `:257-264` | 对（已复用） |
| 装配体抽取器 | — | `cmd/wisp/resident_hotkey_258_test.go:25 runResidentBody258`（票 294 明令"不要长第二枚 walker"） | 对（已复用） |

`withPanelHost` 全仓调用点现读＝**1 枚**（`resident_windows.go:268`），无第二处可冒充。

## 2. 装了什么（判据形状，派单 §5）

新件：`cmd/wisp/exec_holder_toggle_321_test.go`（**200 行整**，sha256 前缀 `717b2d2d9b76d8ae`，8671 字节，无 build tag ⇒ 非 Windows 也能跑这枚尺，与 228/294 同形）。

- **case A `TestAC321PanelHopDrivesTheHolderToggle`＝接线尺**（AST，三支）：
  A1 `runResident` 内 `startResidentBall` 调用**恰 1 枚**且其参数列里 `withPanelHost` **恰 1 枚**（"hook 存在但到不了球宿主"不算过）；
  A2 该实参必须是函数字面量，且体内**确有**对 holder 的方法调用（holder 变量由 `newPanelHolder(...)` 绑定而来——尺认**构造器**不认 `holder` 这个拼法）；
  A3 被调方法名集合**恰为 `{toggle}`**。
- **case B `TestAC321NamedExecutorIsTheToggleRoad`＝能力尺**：把 case A 读出的方法名拿去全产码找方法声明（枚数钉死为 1，防"造第二枚同名 decoy"），打开它的体，断言**发出 `RequestToggle`**。
  ⇒ 这一支是**反词面尺**的那一半：派单禁的是"源码里出现 toggle 这个词"，本腿防的更进一步——**把 `show` 的体重命名为 `toggle`** 这种改名也红。
- 零行号锚（只读 AST 形状）；零 mock；零产码改动（见 §5 面核）。

## 3. 三发正控／突变（派单 §5"缺一不作凭据"）

日志全在 `.scratch/wisp/probes/321/r1/logs/`。

1. **正控（还原态绿）**：`restored-astpin-shot1/2/3-final-disk.txt` — 三发全 `--- PASS` ×2、`ok`、`rc=0`（第 3 发是本轮**第二次物理变异并还原之后**的盘上态，防"我只验过一次还原"）。
2. **突变甲**（票面 AC#1 指定形，`:269` 同一行 `holder.toggle(via)` → `holder.show(via)`）＝ **两枚尺全红**，`rc=1`：
   - case A（`mutant-A-astpin.txt`，逐字）：
     `AC#1 RED: this hop is NOT toggle semantics - runResident (resident_windows.go) hands the ball host the executor set [show] on newPanelHolder (receiver "holder", argument a function literal), want exactly [toggle]. A panel hop that shows instead of flipping is the rewiring ticket 319-r3 split the two executors to forbid: it hid the operator's window on a double click (probes/319/v/03-ac2-live-adjudication.txt shape A).`
   - case B（同文件，逐字）：
     `AC#1 RED: this hop is NOT toggle semantics by capability either - the executor show the panel hop drives posts [Lock Unlock isFinished spawn Info RequestShow RequestRaise], which lacks RequestToggle: the tray item, the ball and the panel hot key no longer flip the panel's state`
   ⇒ 红句**点到"这一跳不是 toggle 语义"**，正是 AC#1 要的那句；且 case B 把 show 到底发了什么摊开（`RequestShow RequestRaise`，无 `RequestToggle`）。
3. **突变乙**（反形：闭包体不调任何 holder 方法，`return false`）＝ **同一两枚尺全红**，`rc=1`：
   - case A：`AC#1 RED: the panel hop in runResident (resident_windows.go) is a function literal that drives NO executor on newPanelHolder (receiver "") - the tray item, the ball and the panel hot key now execute nothing, whatever this file is named`
   - case B：`AC#1 RED: the panel hop (resident_windows.go) drives [] executor(s), want exactly one toggle - no single method here has a capability that could be checked`
   ⇒ 本仓定式"判据换成反形也要红"成立；**两枚尺都有牙，没交出一枚没牙的尺**。

**反面用例的颜色（缺口属实，不是我编的）**：`mutant-A-negative-ruler.txt` ＝ 在**突变甲在飞期间** `TestBallPanelGesturesReachThePanelThread` `--- PASS (0.74s) rc=0`；`restored-negative-ruler.txt` ＝ 还原态同一枚 `--- PASS (0.73s) rc=0`。
⇒ 产码被改成 `show` 而既有那枚尺**照样绿**＝票 321 立案的那句话被现量复认；本格只能由新两枚尺接住（派单 §5 禁的"拿它冒充本格"因此有凭有据）。

**变异台件形态与一枚方法学发现**：本腿用**物理改写＋sha256 还原**，⛔ `-overlay`——理由是现量的：这两枚尺在测试运行时用 `parser.ParseFile`／`os.ReadDir` **直接读盘上的产码**，而 `-overlay` 只换 go 命令编译期的输入，不改测试进程看到的文件内容 ⇒ 用 overlay 变异被测量那一边，尺会**恒绿**，那是假读数不是判据（319-v2 的 overlay 之所以成立，是因为它驱动的是**编译后的**函数）。
还原链现读：基线 sha `816017b68c495b6ca3770b08e01b0cf502d91822cb31dec0305eedafb5f9d90d`；突变甲态 sha `248d5baebb9a651f850ad58876fbfe0503ceb7f8e0122504428588602c90d57e` 且 `git diff --numstat`＝`1 1`（**同一行**，非搬家）；两轮还原后均回到基线 sha，且 `git status --porcelain -- cmd/wisp/resident_windows.go` 空。备份件留在 `probes/321/r1/resident_windows.go.bak`（⛔ 删）。

## 4. 门禁（票面 AC#2；逐把自落 rc=，`logs/gates.txt`）

| 尺 | rc | 读数 |
|---|---|---|
| `go build ./...` | 0 | 无输出 |
| `gofumpt -l <我的件>` | 0 | 空列表＝已净 |
| `go vet -tags winlive ./cmd/wisp/` | 0 | 无告警 |
| `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations`（cmd/ 136 枚 Go 件、注释与 _test.go 全在内 ⇒ 我的注释/字符串都过） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 0 | `VERDICT GREEN`，roster==tree 60 枚，双单位交叉核对同集 |
| 测试两发取交集 | 0 | `TestAC321*` shot1/shot2 **同名两绿**，无"只在一发里出现的色" |

## 5. 面核（改动逐枚）

`git diff --name-only -- cmd/wisp internal docs scripts` ＝ **0 行**（提交前现跑）；`git status --porcelain -- cmd/wisp` 只 `?? cmd/wisp/exec_holder_toggle_321_test.go`。
本腿路径只有四类：新测试件 1 枚／本证据件 1 枚／`probes/321/r1/**` 读数与备份／票 321 尾部 **1** 条 bullet。
⛔ 产码一字（`resident_windows.go` 连注释都没加——它的 sha 与基线逐字节相等就是凭据）；`internal/panel/**`、`internal/winsec/**`、`internal/tools/**`、`frontend/**`、`design/**`、阈值/golden/`allowlist.txt`/`.github/**` 零命中；AC 框、Status、票名、票面既有文本一律未碰。

## 6. 255 名册那发＝**存量红，非本腿分母**（起手/终态各一发）

- 起手 `255-roster-start.txt`：`rc=1`，`--- FAIL: TestTicket255RosterStillMatchesTheActualReadSites`，红两节（逐字）：
  `section "llm" (field LLM) is read by [cmd/wisp/panel_config_store.go cmd/wisp/panel_models_320.go cmd/wisp/providers.go cmd/wisp/run.go internal/llm/resolver.go] in production, the roster names [...四条...]`
  `section "models" (field Models) is read by [cmd/wisp/models.go internal/panel/composer_dispatch.go] in production, the roster names [cmd/wisp/models.go]`
- 终态 `255-roster-end.txt`：`rc=1`，红名集合与起手**逐字同集**（`diff 255-start-names.txt 255-end-names.txt` 空）⇒ 本腿没动分母。
- **归因（现跑，非猜）**：多出的两枚读者都是 `272405f5`（10-15 16:27，320-r2 编排者亲做）引入的——`cmd/wisp/panel_models_320.go` 与 `internal/panel/composer_dispatch.go` 都 tracked、worktree 干净；而 `cmd/wisp/config_readers_255.go` 的 `sectionReadSites` 末次改动是 `3e486af1`（10-14 12:44）。
  ⇒ 这是**票 320 落地后名册未回判**的存量红，落在 255 那族（`LLM`／`Models` 两行），**不是** 321 造成的，也不是本腿能修的（改 `config_readers_255.go` 是别人的面与本票未授权的写面）。**上报编排者处置。**

## 7. 整包分母（派单 §6）

**整包分母不可达（314 死区），四数归编排者另取。** 本腿 ⛔ 跑整包，只跑了定向 `-run`：`TestAC321*`（5 发）、`TestBallPanelGesturesReachThePanelThread`（2 发，两态各 1，均 PASS）、`TestTicket255RosterStillMatchesTheActualReadSites`（2 发，同为存量红）。没编数。

## 8. 没做完的格子（逐枚具名，派单 §9）

1. **票面 AC#1/AC#2 两枚框未翻、Status 未改**——⛔ 自翻；需非实现者裁决表（`docs/evidence/s1/ticket321-321-v*-accept-v1.md`）。需要谁：编排者派 `321-v1`。
2. **本票没写裁决表**：写者＝判者违规，故只交实现方读数。需要谁：另一枚腿。
3. **CI 色未取**：本腿只 commit、⛔ push，`test-windows` 那步的色拿不到。需要谁：编排者推后取（派单铁律）。
4. **255 存量红未修**（§6）。需要谁：编排者裁（归票 320 尾部或 255-r2）。
5. **三枚可选宽容形没装（预算换掉的，具名）**：派单每件 ≤200 行，我为压线砍掉了 (a) `withPanelHost(holder.toggle)` 这种**方法值**拼法的耐受、(b) `via` 实参必须是闭包形参那枚**票面未要求**的自加支、(c) "toggle 不得同时发出 RequestShow/RequestRaise"那枚额外支。后果：这三形今天会**红**（是"过严"不是"没牙"），已在件头 `射程` 段落文里写明"改动者复判本文件、⛔ 把断言改松"。需要谁：非实现者裁要不要补。
6. **一枚纪律违例（自抓，如实入账）**：我为把突变甲的 sha 与 `1 1` numstat 记全，临时 `cp` 出了一枚快照 `logs/mutant-A-subject.txt`，随后 `rm` 掉了它——违反 `issues/README` 规则 8／派单 §9"临时件只建不删"。内容零丢失（同字节仍在 `resident_windows.go.bak`，且基线 sha 已复现两次），但**动作本身违规**，记在我这腿。需要谁：编排者入账（台账由编排者写，本腿不动台账）。

## 9. 你需要做什么（编排者）

1. 收件三把尺现跑：`git log -1 <本腿提交>`、`wc -l cmd/wisp/exec_holder_toggle_321_test.go`（应 200）、`sha256sum cmd/wisp/resident_windows.go`（应 `816017b6…`）＋ `go test ./cmd/wisp -run 'TestAC321' -count=1`。
2. 派 **`321-v1`（非实现者）** 裁 AC#1/AC#2，并重点复判 §8.5 那三枚我砍掉的耐受形与 §3 的"overlay 对读源码的尺无效"这条方法学（它影响本仓所有 AST 钉的突变台件选择）。
3. 处置 §6 那枚 255 存量红（归口 320 或 255-r2），并把 §8.6 的违例入账。
4. 推送由你做；本腿 ⛔ push。
