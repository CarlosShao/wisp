# 票 259 落地腿 `259-r1` — AC#2 拒因可指名 + AC#3 三枚能力尺

本件＝写码腿 `259-r1` 的自证表。派单口径：**只两格**（AC#2 拒因可指名／AC#3 三枚能力尺）；
AC#1（ⓐ具名降级／ⓑ做出牙 二选一）与 AC#4（`cmd/wisp` 载具前置）**本轮不落**，理由与归口见 §6。
写面＝`internal/agent/approval/**` ＋本证据件目录。`cmd/wisp`／`internal/panel`／`internal/tools`
正被 `257-r2`／`181-r3` 取整包红名册，那三枚包本腿连跑都没跑（§5 逐枚列实跑命令）。

---

## §0 起手锚（同发取，2026-10-05 11:3x +08）

| 项 | 读数 |
|---|---|
| `date` | `2026-10-05 11:35:48 +0800` |
| `git log -1 --format=%h` | `d7236abc` |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal tools` | 见下（**起手态，逐字**） |

```
 M internal/panel/git.go
 M internal/panel/git_test.go
 M internal/panel/instructions_200.go
 M internal/panel/workspace.go
 M internal/panel/workspace_test.go
 M internal/tools/paths.go
 M internal/tools/paths_workspace.go
 M internal/tools/paths_workspace_test.go
 M internal/tools/task_pointer_authority_ac3_174r4_windows_test.go
?? internal/panel/workspace_account_181r3_test.go
?? internal/tools/paths_workspace_account_181r3_test.go
```

⚠ 登记不归因：以上全部落在 `internal/panel`／`internal/tools`＝派单点名的**别的腿的脏面**
（`181-r3`）。`internal/agent/approval/` 起手零命中；`cmd/wisp` 起手零命中（`257-r2` 尚未落笔或未进该面）。
本腿⛔ 未 `stash`／未 `checkout`／未 `clean`，一个字节都没清。

起手备份锚（`git cat-file blob HEAD:<path> > <path>.orig`，本腿突变还原的**唯一**比对基准）：见 §4 表末列的逐枚 md5，
每枚文件起手值＝`git cat-file blob HEAD:` 抽出的那份，终值＝交付前的 `md5sum`。

---

## §1 起手名册 ＋ 本腿会撞的既有断言原文（逐枚读过，不是 grep 新符号）

### §1.1 起手逐名读数（⛔ 只跑 `internal/agent/approval/` 这一枚包）

| 命令 | 起／止时刻 | 读数 |
|---|---|---|
| `go test -count=1 ./internal/agent/approval/` | 11:36:10 → 11:36:18 | `ok github.com/CarlosShao/wisp/internal/agent/approval 0.618s`（rc=0） |
| `go test -count=1 -v ./internal/agent/approval/` | 11:36:23 → 11:36:29 | **77 枚 `--- PASS`／0 枚 `--- FAIL`**，`ok 0.611s`；顶层名 **58 枚**（含子测试共 77 行 PASS） |

名册落盘：`.scratch/wisp/probes/259/r1/base-v.txt`（逐字输出）＋`base-names.txt`（58 枚顶层名排序）。
⛔ 起手态**零红名册**——本包全绿，所以本腿交付后如果出现红，只能是本腿种出来的。

对照前人读数：`docs/evidence/s1/242-grant-binding-v1.md` §0 写的是 **66 PASS**（10-03 锚）。
今天 58 枚顶层名／77 行 PASS＝242-v1 之后又落了票 220／260 系那几批用例，**不是同一名册**；
引用时按今天的读数，不要拿 66 当常量。

### §1.2 会撞我新返回形状的既有断言（逐枚原文，逐枚判怎么办）

`spend` 的返回形状从 `bool` 改成带拒因的形状，仓里**只有 6 枚调用者**（尺＝`grep -rn "\.spend("`）：

1. `internal/agent/approval/queue.go:376` 逐字 `spent := it.grants.spend(nonce, it.bind)`——**产码唯一调用点**，本腿要改的就是这一跳。
2. `ticket242_binding_test.go:23` 逐字 `if s.spend("nonce-live", "digest-of-item-two") {` → 红句 `AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding`
3. `ticket242_binding_test.go:29` 逐字 `if s.spend("nonce-live", "digest-of-item-one") {` → 红句 `AC#1 RED: the rejected nonce was still spendable afterwards - single-use was broken by the forged attempt`
4. `ticket242_binding_test.go:40` 逐字 `if s.spend("nonce-e", "") {` → 红句 `AC#1 RED: empty binding spent a live grant`
5. `ticket242_binding_test.go:50` 逐字 `if !happy.spend("nonce-a", want) {` → 红句 `AC#1 RED: the exact minted digest failed to spend its own grant - the mint/spend round trip is broken`
6. `ticket242_binding_test.go:104` 逐字 `if itA.grants.spend("leaked-or-guessed-nonce", itB.bind) {` → 红句 `AC#1 RED: item B's binding spent item A's grant - cross-item binding is broken`

**处置**：2–6 这五枚按**逐字等价**改判据形状（`spend` 返 `denialNone` ⇔ 旧 `true` ⇔ 旧"允许放行"那一支），
红句文本一字不改、失败条件不放宽。等价式写死在本件 §2。⛔ 不删任何一枚、⛔ 不把 `t.Fatal` 降级成 `t.Log`。

### §1.3 会撞我新审计行的既有钉（读到的原文，逐枚）

- `cmd/wisp/approval_reply_201_test.go:418` 钉逐字 `"approval: FORGED-OR-STALE allow rejected corr=t201-panel-corr"`，
  匹配方式＝`strings.Contains(audit, want) || strings.Contains(sentences, want)`（`:424-427` 那圈）。
- `cmd/wisp/subagent_selfapproval_197_test.go:336` 钉逐字 `"approval: FORGED-OR-STALE allow rejected corr=" + rec.corr1`，
  同包同一枚 `strings.Contains(audit, want)`（`:340-343`），同发还钉 `"approval: ANSWER-ALLOW corr=" + rec.corr1 + " tool=fs.write route=native decision=allow"`。
  ⇒ 结论：`queue.go:379` 那枚 `FORGED-OR-STALE` 行**必须一字不动**（含那句 `missing/spent/misbound` 合并口径），
  新增的指名行只能**另起一行**。枚数/计数式断言在两枚钉里都不存在（是 Contains），本腿另在 §5 复核"没有一枚断言数审计行的条数"。
- 对外文案钉：`cmd/wisp/approval_reply.go:227`／`:272` 逐字 `原生令牌无效（缺失/已用/与本次请求不绑定）`，
  `cmd/wisp/approval_reply_201_test.go:387` 钉 `"原生令牌无效"`；`ui.go:130` 是这三份副本里的本包那一份。
  ⇒ 编排者已裁「对外维持合并」：`ui.go:127-130`（含 `// The message never says which, so the API cannot be used to probe` 那句）
  **一字不改**，`ErrBadGrant` 的文案与错误值**一字不改**。撤销口令＝「259 对外也区分」，本腿没碰。
- `internal/agent/approval/queue_test.go:132`／`:166` 用 `errors.Is(err, approval.ErrBadGrant)` 钉拒答的错误值；
  `ticket87_veto_l2_test.go:119`、`ticket97_alias_direction_test.go:87` 同。⇒ 四种拒因对外**仍然**折成 `ErrBadGrant` 才不红，这正是裁定的形状。
- `internal/agent/approval/fakes_test.go:276-280`：`newGate` 的 `t.Cleanup` 断言 `g.Queue().Depth() != 0` 即报红
  ⇒ 本腿所有走路由的新用例**必须把卡收干净**（答复或拒绝），否则会被这枚既有钉打死。
- `internal/agent/approval/ticket242_panelface_test.go:21-28` 的 `panelItemReadFace` 是 M-D2 那发绕过的那份名单
  （`"Permitted"` 一加就绿）。⇒ 本腿⛔ 不改这枚既有尺的名单与判据（改了就是"把被绕过的钉拔下来换新形"，
  对抗验收会读成洗读数）；新尺另立一枚文件，**按能力**判，与它并存。

### §1.4 已知正控 `197-v1` 的 D2（本腿新尺不许把它做成常红或常绿）

出处＝`docs/evidence/s1/242-grant-binding-v1.md` §2 M-A2 行（该腿把票面 D2 落实的形）：
`approval.go` 的 `grantStore.spend` 改成不比摘要。⚠ 派单点名的是 **197-v1 的 D2**，本腿在
`docs/evidence/s1/` 里检索到的**同形实跑读数**在 242-v1 表 M-A2（`-- 不比摘要 --` 那一发）：
`TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce` 与 `TestTicket242SpendRequiresTheExactBinding` 红。
本腿的处置＝§4 的 **PC-D2** 一发：种下同形 ⇒ 指名"绑定不对"那一枚新用例**必须红**（非常绿），
不种时全绿（非常红）。两向都量，读数在 §4。
