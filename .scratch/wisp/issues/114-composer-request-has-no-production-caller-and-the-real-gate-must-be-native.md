# 114 — `ParseComposerRequest` 在**生产里零调用者**：面板那条"档位/附件/工作区"的输入通路今天还没接上；接上的那一票**必须把门放在原生侧**（票 92 验收 `R-92-1` / `R-92-2`）

> **⛔ 2026-09-23 09:22 编排者加注（owner 把前端整块交给别的 agent，台账 `A102` / README 规则 7）**：
> 本票**不整体冻结**——它的正解主体是 Go 侧权限接线（`internal/panel/`、`cmd/wisp/` 装配根），仍该被干。
> 但 **AC#6（`R-92-5` 真机差分截屏）与任何要写进 `frontend/src/` 的改动不属于本票的派发范围**：
> 领这票的代理**只做非前端腿**，AC#6 那格**保持未勾**并在 Progress log 写 `skipped=frontend(owner-delegated)`，
> **不许替它勾、也不许"只改一角"继续写**。判据：**先把要写的文件列出来，有一枚在 `frontend/` 下就整段停手上报**。

**Status:** in-progress（2026-09-23 10:4x `agent-ticket114-ac2` 交 AC#2 原生侧门 + AC#1 现状表入库；
              **AC#3/#4/#5/#7/#8/#9/#10/#11 未做**，AC#6 按票头 09-23 横幅**保持未勾**（`skipped=frontend(owner-delegated)`）；
              2026-09-21 20:2x 编排者建；来源=`acceptor-ticket92` 的 `docs/evidence/s1/92-adversarial-acceptance.md`
              总判一节的两条记账：`R-92-1`（文本门只覆盖三种扩展名的字面量）、`R-92-2`（真边界在原生侧））
**Type:** 能力已实现但没接线（memory 第 8 条那个形状的**第五次**）+ **一条前置条件**（接线方式决定门是否成立）
**Blocks:** owner 要的"在主面板输入框里看到并请求档位"（R20 的最后一段路）· **Blocked by:** 票 92 结案（同一批文件）
**Packages:** `internal/panel/`（请求解析的原生侧落点）、`cmd/wisp/`（装配根：把 composer 请求接到 `perm.Store`/审批链）、
              `internal/perm/`（存储语义**别改**）、`frontend/src/`（发起侧，只许发请求，不许出现 `approval.decide` = ban #6）。
              **禁改**：`internal/risk/**`、`docs/PLAN.md`、`docs/specs/*.md`、`tools/d22scan/**`、`allowlist.txt`、任何阈值/断言/golden。

## 为什么现在立案（不是"以后再说"）

`acceptor-ticket92` 实测：`ParseComposerRequest` 的**生产调用者 = 0**，`perm.Store.Set` 的**生产调用者 = 0**
⇒ 它因此判"面板改不了档位"**实质成立，但成立的原因不是票 92 声称的那两道门**——
两道门只钉住了两种字面写法，而真正的防线是**这条通路根本不存在**。

这就是危险所在：**"通路不存在"是一个会被下一次接线动作消灭的属性。**
接上之后，唯一挡着面板改档位的就只剩 `composerModeWriteRe` 那类**源码文本门**——
而验收代理已经证明：被攻陷的渲染进程在运行时调 `window.chrome.webview.postMessage(...)`
**不需要源码里出现任何被禁字样**（它造的 `.js` + 运行时拼名两种形状全绿）。
⇒ 所以这条账**不能只登记在票 92 的残留里**：它必须变成接线那一票的**硬 AC**，否则 92 结案那天
    我们等于把"绿的原因"换成了"绿的样子"。

## AC（1:1，裁决表 `docs/evidence/s1/114-*.md` 由验收方出）

- [x] **AC#1** 先出**现状表**：`ParseComposerRequest` / `perm.Store.Set` / 面板三块输入（档位、附件、工作区）
      各自的**生产调用者**逐条点名（文件:行），并明确"接通之后哪一条从 0 变成 N"。不许用"应该已经接了"结框。
      ⇒ 交付：`docs/evidence/s1/114-ac1-status-table.md`（只读代理 `audit-ticket114-ac1`，锚定 `7ad6eb4`，本代理入库）。
      ⚠ 该表的 ①.1 结论是**"WebView2 事件 → `ParseComposerRequest` 那一跳今天根本不存在"**，故 AC#3 的反向对照在本段**做不到**（见 Progress log 的"没做到的"）。
- [x] **AC#2** 原生侧门（`R-92-2` 的正身）：**mode 写入处理器在注入的 `Confirm`（L2 确认腿）为 nil 时必须拒**，
      并且**加一条用例钉住"没有 L2 确认腿就不许写档位"**。判据要能被变异证伪：
      把 nil-检查删掉 ⇒ 该用例必须红；把它改成"只记日志不拒" ⇒ 也必须红。
      ⇒ 交付：`internal/panel/composer_handlers.go`（处理器 + `Confirm == nil` 的 fail-closed）+
      `internal/panel/composer_handlers_test.go`（判据 = **注入的 `ModeWriter` 未被调用**，不是"返回了 error"）+
      `cmd/wisp/run.go` 的装配那一格。覆盖面按编排者 P1 裁定：**只挡"把档位往更宽切"，变严不挡**；审计三档全写未被改变。
      ⚠ 门已立，**通路仍不存在**（P0 地界：宿主那一跳留票 33/35）——见 `docs/evidence/s1/114-ac2-native-gate.md`。
- [ ] **AC#3** 反向对照：请求**自动档/更宽档**在真机上必须走 L2 原生卡片（C18），超时 300s = 拒；
      并钉住**审计三档全写**（R20）这一条不因接线而漏。
- [ ] **AC#4** `R-92-1` 的覆盖面：要么把门二的扩展名集扩到与 d22scan 的文本类一致
      （`.js/.mjs/.jsx/.html/.json`），要么把这族判定改挂在**结构**上（全树计数，不是单文件字面量）。
      **不许**为了通过而删判定或缩范围；验收方自己种的越界形状必须红。
- [ ] **AC#5** 不越界：面板侧**只许显示 + 发起请求**（R20 的明写）。工作区/档位两个输入口不得变成授权口；
      `frontend/` 里出现 `approval.decide` 即 ban #6 红；零 emoji（ban #8）含测试文件与注释。
- [ ] **AC#6** 可见性：`R-92-5` 欠的那一笔——真机**差分截屏**（改档位前后各一张，能看出输入框上的档位文字变化）。
      fixture HTML 属渲染级证据，**不许升格成像素证据**。owner 不在签收窗口时先停在"已生成、待签收"。
- [ ] **AC#7** 门禁：`bash -n` / `sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
      `gofumpt -l` 对本轮改动文件**必须真跑并贴原始输出**（票 92 因为"没跑却说没二进制"被打回，本票不得重犯）；
      POSIX 三合包数字按 `-v` 口径报（分包/顶层说清楚，参考 `R-92-4`）。

## 编排者追加（21:5x，来源=`acceptor-ticket92b` 的 `docs/evidence/s1/92b-adversarial-acceptance.md`）

验收方原本建议"另开票 123 做 composer 接线"。**我并进来而不另开**——本票标题逐字就是这件事，
再开一张只会得到两张各自半对的票。它的实测形状里**两形全绿**，这两形是本票的硬判据来源：

- [ ] **AC#8**（`R-92b-1` 的正解）：**把"两边同改"从承诺变成门**——往 `knownComposerMethod`（原生侧方法白名单）
      加一条时，**必须**同时让 `frontend/` 出现对应字面量，否则用例红。
      现在的洞是：`sendRequest("panel"+"."+noun+"."+verb)` 在运行时等于 `panel.approval.decide`，
      而**三枚判据各自只认字面量** ⇒ 运行时拼接**全绿穿过**（含 `sh scripts/d22scan.sh` rc=0 clean）。
- [ ] **AC#9**（`R-92b-2`）：结构钉的**扫描根比 `ban #6` 与"渲染器实际加载的文件集合"都窄**——
      第二形实测：把一条通道放进 `frontend/public/*.js` + 在 `index.html` 里引用它，
      那枚钉报的是 `21 files scanned, 2 host call sites`，**一个数都没动** ⇒ 全绿。
      ⇒ 本票要回答"**渲染器到底会加载哪些文件**"，并把可达性钉挂在**那个集合**上，不是挂在 `frontend/src` 上。
      ⚠ 已知它不是扫描器的问题：同一位置直接写字面 `approval.decide` 时 `d22scan` 会红并点名该文件 ⇒ **漏在钉的扫描根**。
- [ ] **AC#10**（`R-92-2` 的另一半，验收方新补）：`panel.ts` 注释自陈宿主有**两种装法**
      （`AddHostObjectToScript` / `postMessage`），而可达性钉**只认后者**
      ⇒ 第一种要么**显式禁用**、要么进同一枚钉。
- [ ] **AC#11**（`R-92b-3`）：把 `npm run render:composer && git diff --exit-code` 做成 CI 一步
      （关掉"fixture 可以手写/可以腐坏"这一格）。⚠ 动 `ci.yml` 之前先看票 111/85a 是否已让出该文件。

⚠ **两条我原样带过来的诚实边界**：① 验收方复算出 `gofumpt` 那格有一条新仪器坑——
票面那句 `export PATH="$PATH:$(go env GOPATH)/bin"` 在 Git Bash 下会被 MSYS 吃成 `\work\base\gopath/...`
⇒ **rc=127"命令不存在"是假读数**，必须写 `.exe` 全路径；② `R-92b-7`：交件里"本轮 `frontend/` 一行未改"
在**树级别为假**（差 1 行，实为注释 ⇒ 结论仍成立，但这种话以后要按树说）。

## Rules（本仓固定）

- 只 commit 不 push；共树禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）。
- `git add` 只用显式路径；每次 commit 前 `git diff --cached --name-only`。
- 登记"未跑"必须引**命令原文 + 错误原文**（本轮新增，来源=票 92 的 `R-92-3`）。
- 验收期若发现探针达成了本票 AC 声称要防的结局 ⇒ 判 FAIL，不写"通过附条件"。

## Progress log（append-only）

- 2026-09-21 20:2x（编排者）：建票。来源是票 92 验收的 `R-92-1`/`R-92-2`。
  建票前已确认**无同族票**：`grep -rln "ParseComposerRequest" .scratch/wisp/issues/` 只命中
  票 92 本身与两张已结案票（101/110）的引用；"composer 原生侧接线"没有归属。
  与票 92 的关系写死：本票 **Blocked by 票 92 结案**（同文件地界），不在 92b 修窗口期内动手。
- 2026-09-23 10:4x +08 `agent-ticket114-ac2` did=落 AC#2 原生侧门：新建 `internal/panel/composer_handlers.go`
  （`ModeWriteHandler.HandleModeRequest`，第一句之后的 fail-closed = `modeIsWidening(from,to) && h.Confirm == nil` ⇒ 拒，
  **位置在任何 `ModeWriter.Set` 调用之前**）+ `internal/panel/composer_handlers_test.go`（6 枚顶层用例，判据 =
  "注入的 `ModeWriter` 一次都没被调用"）；覆盖面按编排者 P1 = **只挡变宽、变严不挡**，审计三档全写未动
  （拒的一支自己写 audit 行，过的一支交给 `perm.Store.Set` 记）。`internal/panel` 零新增包依赖（`ModeWriter`/`ModeConfirm`
  都是本地形状，未 import `internal/perm`）。门禁：`go test -count=2 -v ./internal/panel/` = `=== RUN 152 / --- PASS(顶层) 92 / FAIL 0 / SKIP 0`。
  next=AC#3 的反向对照要等票 33/35 的宿主那一跳（本票 P0 地界外）；AC#4/#8/#9 的门二覆盖面与扫描根仍未动。
- 2026-09-23 10:5x +08 `agent-ticket114-ac2` did=`skipped=frontend(owner-delegated)` —— AC#6（`R-92-5` 真机差分截屏）
  与任何 `frontend/src/` 改动**未做、不勾**，按票头 09-23 横幅与 `A102③` 记在本行；本轮写码文件清单里**零枚**在 `frontend/` 下
  （`internal/panel/composer_handlers.go`、`internal/panel/composer_handlers_test.go`、`cmd/wisp/run.go`、本票面、`docs/evidence/s1/114-ac2-native-gate.md`）。
  next=编排者把 AC#6 交给 owner 指派的外部 agent；本票剩余格见上一条。
