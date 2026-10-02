# 票 256 — 常驻那条腿吃不到 `[risk]` 那两项配置：ⓘ 那一支（把 `approval.New` 移进 `assembleRuntime`）今天没裁，因为**没人量过"移动会不会破 `Confirming` 那一维"**

**立票时刻**：2026-10-02 12:24:00+0800，锚点 HEAD `4a9851d6`（`dev`）
**来路**：非实现者裁决腿 `248-v1c` 的 `docs/evidence/s1/248-settings-write-path-v1.md` §1（AC#10 那格）＋ §5 第 8 条＋ §6 第 4 条；台账 `A530`。母票＝**票 248 AC#10**（我裁了 ⓑ＝回执文案说真话，⛔ **ⓘ 这一支不在今天**）。

## 现量（编排者本机自跑，2026-10-02 12:2x，⚠ 引用前先重跑，别把这几行当常量）

1. 常驻那枚门是**在会话账本之前**建的：`cmd/wisp/resident_approval_windows.go:108-113` 逐字 `ra.gate = approval.New(approval.Options{` 里只有 `UI`／`Channels`／`Logf` 三项——**没有 `Grants`、没有超时、没有 L1 窗口**。
2. 常量在那两处：`internal/agent/approval/queue.go` 的 `DefaultApprovalTimeout = 300s`／`DefaultL1Window = 3s`／`MaxL1Window = 3s`；钳位在 `internal/agent/approval/gate.go:139-145` 的 `win <= 0 → DefaultL1Window`（行号来自母票 `246-v2` 现量，**属待验断言**）。
3. ⇒ 真实代价两条：① 常驻腿里「本会话内允许」因为 `Options.Grants` 为 nil **落不下一行**；② `confirm_timeout_sec` 那类改动**在常驻腿不生效**，而设置页今天不说这句话（这句"不说"由 ⓑ 那一支先补上，归 `255-r2`）。

## 要建什么（ⓘ 那一支，⛔ 一动手就撞上"改票 246 AC#1 裁过的乙形次序"）

- [ ] **AC#0（本票第一格，且是闸门）＝先把"移动会不会破 `Confirming` 那一维"量出来**：现读 `approval.New` 现在吃哪三项、挪进 `assembleRuntime` 之后谁在挪之前需要那枚门（常驻卡片 UI／热键答复／托盘「允许一次」各一处），逐处带 `file:line` 与"挪之后还拿不拿得到"。⛔ **本格不许直接改产码**；交件判据＝量不到的那一处要具名说"量不到"，⛔ 不许用"应该没问题"填空。
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
