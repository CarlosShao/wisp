# 票 256 — 常驻那条腿吃不到 `[risk]` 那两项配置：ⓘ 那一支（把 `approval.New` 移进 `assembleRuntime`）今天没裁，因为**没人量过"移动会不会破 `Confirming` 那一维"**

**立票时刻**：2026-10-02 12:24:00+0800，锚点 HEAD `4a9851d6`（`dev`）
**来路**：非实现者裁决腿 `248-v1c` 的 `docs/evidence/s1/248-settings-write-path-v1.md` §1（AC#10 那格）＋ §5 第 8 条＋ §6 第 4 条；台账 `A530`。母票＝**票 248 AC#10**（我裁了 ⓑ＝回执文案说真话，⛔ **ⓘ 这一支不在今天**）。

## 现量（编排者本机自跑，2026-10-02 12:2x，⚠ 引用前先重跑，别把这几行当常量）

1. 常驻那枚门是**在会话账本之前**建的：`cmd/wisp/resident_approval_windows.go:108-113` 逐字 `ra.gate = approval.New(approval.Options{` 里只有 `UI`／`Channels`／`Logf` 三项——**没有 `Grants`、没有超时、没有 L1 窗口**。
2. 常量在那两处：`internal/agent/approval/queue.go` 的 `DefaultApprovalTimeout = 300s`／`DefaultL1Window = 3s`／`MaxL1Window = 3s`；钳位在 `internal/agent/approval/gate.go:139-145` 的 `win <= 0 → DefaultL1Window`（行号来自母票 `246-v2` 现量，**属待验断言**）。
3. ⇒ 真实代价两条：① 常驻腿里「本会话内允许」因为 `Options.Grants` 为 nil **落不下一行**；② `confirm_timeout_sec` 那类改动**在常驻腿不生效**，而设置页今天不说这句话（这句"不说"由 ⓑ 那一支先补上，归 `255-r2`）。

## 要建什么（ⓘ 那一支，⛔ 一动手就撞上"改票 246 AC#1 裁过的乙形次序"）

- [x] **AC#0（本票第一格，且是闸门）＝先把"移动会不会破 `Confirming` 那一维"量出来**：现读 `approval.New` 现在吃哪三项、挪进 `assembleRuntime` 之后谁在挪之前需要那枚门（常驻卡片 UI／热键答复／托盘「允许一次」各一处），逐处带 `file:line` 与"挪之后还拿不拿得到"。⛔ **本格不许直接改产码**；交件判据＝量不到的那一处要具名说"量不到"，⛔ 不许用"应该没问题"填空。**〔10-02 14:3x 编排者现读复认后翻勾；凭据与"我票面三句被推翻"在下面第 7 节；第三处并列不成立〕**
- [ ] **AC#1 只在 AC#0 交完并由编排者落一枚具名 `A##` 批准之后才许动**：`approval.New` 的 `Options` 里补上 `Grants` ＋ 来自 `[risk]` 的两项（超时／L1 窗口），常驻腿那次运行里「本会话内允许」**真落一行**。凭据形状**复用票 224 r2 那套授权仪器**，⛔ 不许新造一台。
- [ ] **AC#2 反向判据＋正控**：种一发 `confirm_timeout_sec` 改动 ⇒ 常驻腿那枚窗口的实际钳位值随它变（正控）；⛔ 且 `GRANT-DROPPED` 那行**必须不再出现**（母票 AC#10 写的就是这个形状）。
- [ ] **AC#3 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。⛔ **D43 转移表一个字不许动**——如果量出来"必须新增一个状态才挪得动"，那**不是本票的活**，停下来上报。

## 禁区（本票全程）

- ⛔ 未定义即停：票 246 AC#1 裁过的那条**次序**属既有裁定，改它＝**人工批准**（`SPEC-12 §4.1`），编排者单方面派腿＝跑歪模式。
- ⛔ 不动 `internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件；不许为变绿放宽断言；不许 `t.Skip`；不许把 SKIP 读成通过。
- ⛔ 凭据值绝不进对话／日志／表（只写变量名）。
- Git：只 commit 不 push；显式 pathspec；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。
- `frontend/**`／`design/**` 两层禁令；`grep`/`find` 显式根（`cmd internal tools docs scripts .scratch`）。

## 排程

写面＝`cmd/wisp`（＋`internal/agent/approval`）⇒ ⛔ 与按住中的 `197-r3`／`167-r1`／`255-r2`／`174 AC#2b` 同面**串行**；AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与整包测量同机即互洗）。
**默认不排**：本票排在 ⓑ 那一支（`255-r2`）交完、且 owner 对"常驻腿该不该吃配置"没有相反意思之后。撤销口令「256 撤」。

---

## 7. `256-a1` 收档 ⇒ AC#0 翻勾（编排者现读复认）＋母票 248 AC#10 那枚具名裁定落地（2026-10-02 14:3x，现量 HEAD `1d8106d8`，`git status --porcelain -- cmd internal`＝0 ⇒ Go 面与 `bf498283` 同形）

**1. 交件只认盘上**：`.scratch/wisp/probes/256/a1/census.md` **184 行／40,694 字节**（14:51:49 现量；我 14:3x 读到的是 181 行，其后又落两发）、八节全填、`grep -cE '（待填）|（填写中）|TODO'`＝0；流水六发（14:5x 现量 `git log --format=%h -- <那枚 md>`＝`f4299a8d`→`3245656c`→`bf498283`→`b40a8e72`→`fcdf0f96`→`1abf7f2a`），**逐名 `git show --numstat` 数过去：除第一发＝4 枚之外，其余五发各 1 枚**（都在它自己的写点内）。⚠ 归属事故记在 `A532`：`f4299a8d` 那一发**同时带走了我这轮暂存的三枚文件**——尺＝`git show --numstat --format= f4299a8d`（四行：`census.md 49/0`＋票 255 `2/0`＋票 257 `28/0`＋台账 `17/0`）；那枚腿自己在 `fcdf0f96` 里把这件事具名登记，并当场把写法改成"add 与 commit 同发、commit 必带 pathspec"。⛔ 不 revert、⛔ 不改写已成史那发。

**2. 编排者现读复认四把尺（14:3x，只读 `grep`／`sed`，⛔ 我没跑任何 Go 命令）**：
- `cmd/wisp/resident_approval_windows.go:109-113` 逐字只有 `UI`／`Channels`／`Logf` 三枚字段 ⇒ **本票 §现量.1 内容成立、区间起点偏一行**（真身 `:109`；我票面写 `:108-113`）。
- `internal/agent/approval/gate.go` 里 `g.grants` **只有两枚命中**（`:667` 判 nil、`:678` `Record`），同在一支 `allowSession` ⇒ "Grants 为 nil 不改判定、只改落不落盘＋写哪一句审计"成立。
- `cmd/wisp/resident_task_source_windows.go`：`:218 console := interactiveStdin()` → `:230 return nil` → `:265 run, code := assembleRuntime(spec)` ⇒ **那枚提前 return 确实在装配之前**。
- `internal/ball/tray_windows.go`：菜单项只有 `:86 打开面板`／`:88 静音`／`:89 暂停唤醒`／`:91 退出`，`Allow`／`允许` **零命中** ⇒ 托盘「允许一次」今天不存在。

**3. AC#0 翻勾**：凭据＝上四把尺＋普查件 §1 那 11 行逐处 `file:line`＋判语只用"拿得到／拿不到／量不到"三词；⛔ 它没有用"应该没问题"填空（第 9 处举卡入口具名"量不到"＝`AskOnTaskRoot`／`askConfirmation` 产码调用者 **0 枚**）。⇒ **最小爆炸半径读数＝两行**：`:116` 把门当**指针值**递给账本（`Attach` 只有整包覆盖语义、没有"稍后补门"的入口）、`:157` boot 那一刻加载 Esc。撤销口令「256 重判」。

**4. 它推翻我票面三句（都成立，原话不改、就地记，三笔都归我）**：
- ★ ①AC#0 题面那句"常驻卡片 UI／热键答复／托盘「允许一次」各一处"**并列不成立**——第三处今天没有（它的归口＝票 244 J9，属**待建**不属**待挪**）⇒ AC#0 能交的现状只有两处。
- ②§现量.2 的钳位真身是 `gate.go:137-145`（我引的是母票 `246-v2` 的 `:139-145`；**内容一致、行号窄了一格**），常量三枚复认在 `queue.go:116`（3s）／`:121`（2s）／`:122`（3s）、`DefaultApprovalTimeout` 在 `:107`（300s）。
- ③★ **它把我那个问法换了**（本票今天最值钱的一枚读数）："`Confirming` 会不会破"——**表那一半没有可破的东西**（常驻腿不持有 Machine：`statemachine.New` 全仓产码只一枚调用者 `cmd/wisp/models.go:303`；球不查表：`internal/ball/ball_windows.go:311-313` 的 `SetState` 只 `PostTask(applyStateLocked)`，没有 `Fire`／没有转移检查；`Confirming` 这个名字只吃卡片自己的 `Level` 字符串：`resident_approval_windows.go:501-506`）⇒ **⛔ 本票不必上报"要改 D43"**；真变的是**门在哪一刻存在**。我要是按原问法裁"不破就可以挪"，就会裁在一枚没被那一问问到的前置条件上。**这一格它拒绝替我补，补是我的活。**

**5. ★ 母票 248 AC#10 定案＝ⓑ（账 `A534`，落在我这条裁定里）**：**不移动 `approval.New`**；改在设置那几项的回执里逐字写明"这两项只作用于跑任务的进程，常驻腿今天用常量 300s／3s"，且那句话**必须由登记表同源产出**（归 `255-r2`，⛔ 不许硬写死一句漂亮话）。
- **理由（现量，不是我的判断）**：ⓘ 若做，门的存在时刻从"进程启动即存在"改成"仅当这发常驻进程有交互控制台、且装配根没先失败"（`resident_task_source_windows.go:218→:230→:265`）⇒ **双击／Explorer 拉起那一支（正是 D2／票 228 那个"用户真正启动的进程"形状）永远走不到装配**＝没有门、没有 `PendingWindow`、球永不进 `Confirming`、Esc 永不借、`:157` 永不加载。⛔ 这不是"那一维变弱"，是**载体不存在**。
- **边界**：本票 **AC#1／AC#2／AC#3 仍不排腿**（§排程那句"默认不排"继续有效）；D43 转移表一个字没动、也不需要动。撤销口令**「248 AC#10 改 ⓘ」**。
- ⚠ 谁要重开 ⓘ，**必须先答两格**（都已有现量、无人给修法）：① 无控制台那一支的门从哪儿来；② §7-6 丙——"本会话内允许"在常驻腿真实的前置门槛是**四道**，`Grants` 为 nil 只是第四道，而第三道（得有人说得出具 `session` 那个词，链在 `cmd/wisp/approval_reply.go:566 case "session"` → `:567` → `:257`，产码里**只有这一条**）在无控制台形状里**根本没有入口** ⇒ **本票 AC#1 那句"真落一行"在无控制台形状里不可判**；我把它当待验断言带着，⛔ 不改票面原句。

**6. AC#2 判据字面有洞（我票面自己写的"必须不再出现"，就地定形，原句不改）**：普查件 §3④★★ 量到 `GRANT-DROPPED` 有**两形**——第一形 `gate.go:667-670`（`g.grants == nil`，且发生在放行**成功之后**），第二形 `gate.go:673-675`（"卡片主题在答复前已离开队列"，**与 `Grants` 无关、挪门挪不掉它**）；而既有仪器 `internal/agent/approval/ticket224_reply_grant_test.go:283` 只查前缀 `approval: GRANT-DROPPED`、对两形**不加区分**。⇒ AC#2 真到那天定形＝**只禁第一形**，第二形必须仍能出现；⛔ 不许为了"前缀零命中"去动第二形那句。
**另记一枚无钉的格（§3⑤）**：常驻腿"会话档不落盘"这一格今天**零用例守**（`GRANT-DROPPED` 的三枚命中里没有一枚在常驻腿：`cmd/wisp/ticket224_assembly_test.go:194` 是反向钉、`ticket224_reply_grant_test.go:128`／`:283` 是包内仪器），只被两处注释与三处文档记着 ⇒ 这一格归本票 AC#2 那天的凭据，⛔ 现在不派腿。
