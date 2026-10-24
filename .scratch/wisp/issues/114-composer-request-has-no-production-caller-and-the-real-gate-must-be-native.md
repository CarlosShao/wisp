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
- [x] **AC#4** `R-92-1` 的覆盖面：要么把门二的扩展名集扩到与 d22scan 的文本类一致 〔10-14 2026-10-14 16:25 编排者翻勾：`114-v1` 判**成立**——自造单删 `.json` ⇒ 对齐尺红点名；全量回退旧 7 类 ⇒ 点名 9 枚；旧集下 `.json` 种形同时逃过结构钉与门二＝"改前逃过两枚仪器" old-state 复读成立。裁决表 `docs/evidence/s1/ticket114-114-v1-accept-v1.md`（`387e2d2b`）；账 `A937`〕
      （`.js/.mjs/.jsx/.html/.json`），要么把这族判定改挂在**结构**上（全树计数，不是单文件字面量）。
      **不许**为了通过而删判定或缩范围；验收方自己种的越界形状必须红。
- [ ] **AC#5** 不越界：面板侧**只许显示 + 发起请求**（R20 的明写）。工作区/档位两个输入口不得变成授权口；
      `frontend/` 里出现 `approval.decide` 即 ban #6 红；零 emoji（ban #8）含测试文件与注释。
- [ ] **AC#6** 可见性：`R-92-5` 欠的那一笔——真机**差分截屏**（改档位前后各一张，能看出输入框上的档位文字变化）。
      fixture HTML 属渲染级证据，**不许升格成像素证据**。owner 不在签收窗口时先停在"已生成、待签收"。
- [x] **AC#7** 门禁：`bash -n` / `sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降； 〔10-14 2026-10-14 16:25 编排者翻勾：`114-v1` 判**成立（抽核口径）**——panel 整包 live 223/143/5/0 vs 对照件 220/141/5/0，**红集 diff 空**＋抽 2 名逐字同句。账 `A937`〕
      `gofumpt -l` 对本轮改动文件**必须真跑并贴原始输出**（票 92 因为"没跑却说没二进制"被打回，本票不得重犯）；
      POSIX 三合包数字按 `-v` 口径报（分包/顶层说清楚，参考 `R-92-4`）。

## 编排者追加（21:5x，来源=`acceptor-ticket92b` 的 `docs/evidence/s1/92b-adversarial-acceptance.md`）

验收方原本建议"另开票 123 做 composer 接线"。**我并进来而不另开**——本票标题逐字就是这件事，
再开一张只会得到两张各自半对的票。它的实测形状里**两形全绿**，这两形是本票的硬判据来源：

- [x] **AC#8**（`R-92b-1` 的正解）：**把"两边同改"从承诺变成门**——往 `knownComposerMethod`（原生侧方法白名单） 〔10-14 2026-10-14 16:25 编排者翻勾：`114-v1` 判**成立**——仓外树真种"页侧单边漂移" ⇒ 新门 FAIL 逐字点名（"adding a route on one side only is the drift…"）；同一棵树上 253 尺与 N7 **PASS（盲）**＝两把尺覆盖不同方向，现读引其自陈行 `inbound_roster_253_test.go:19-20`。账 `A937`〕
      加一条时，**必须**同时让 `frontend/` 出现对应字面量，否则用例红。
      现在的洞是：`sendRequest("panel"+"."+noun+"."+verb)` 在运行时等于 `panel.approval.decide`，
      而**三枚判据各自只认字面量** ⇒ 运行时拼接**全绿穿过**（含 `sh scripts/d22scan.sh` rc=0 clean）。
- [x] **AC#9**（`R-92b-2`）：结构钉的**扫描根比 `ban #6` 与"渲染器实际加载的文件集合"都窄**—— 〔10-14 2026-10-14 16:25 编排者翻勾：`114-v1` 判**成立（钉半支）**——仓外树真种 `public/…-door.js`：旧根 overlay **PASS**（"绿与事实无关"复现）／新根 **FAIL** 点名到 `:1` 行与 `panel.mode.set`；真树 113 files/2 call sites；精确加载集（J3）留档。账 `A937`〕
      第二形实测：把一条通道放进 `frontend/public/*.js` + 在 `index.html` 里引用它，
      那枚钉报的是 `21 files scanned, 2 host call sites`，**一个数都没动** ⇒ 全绿。
      ⇒ 本票要回答"**渲染器到底会加载哪些文件**"，并把可达性钉挂在**那个集合**上，不是挂在 `frontend/src` 上。
      ⚠ 已知它不是扫描器的问题：同一位置直接写字面 `approval.decide` 时 `d22scan` 会红并点名该文件 ⇒ **漏在钉的扫描根**。
- [x] **AC#10**（`R-92-2` 的另一半，验收方新补）：`panel.ts` 注释自陈宿主有**两种装法** 〔10-14 2026-10-14 16:25 编排者重裁（`114-v1` 顶回）：**面＝本队可写仪器格，⛔ 非 frontend(owner-delegated)**——a2 普查 `:198` 逐字列为仪器格（"另有 1 枚（AC#6）永不归本编队"）；a5 行① 落点全在 Go；真门＝`panel_host_windows.go:80` const＋`:842 w.Bind`（a5 引的 `:82/:319` 系旧锚已腐烂）；现钉 `composer_test.go:441` 只认 `.postMessage(` ⇒ **补钉方向＝让可达性钉认第三装法**（归 `114-r6` 落地）；③"本波不落"站得住（普查自陈＋J3）。panel.ts 注释那半（若非本格必需）归前端。账 `A937`〕
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
- 2026-09-23 10:5x +08 `agent-ticket114-ac2` did=交装配与自证：`cmd/wisp/run.go:378` 把 `rt.modeWrites` 装配出来
  （`Confirm` 交的是 `perm.New` 上面那**同一枚** `confirm`，字段 `run.go:229`），三发变异逐发还原
  （M-1 门退回不检查 nil ⇒ `composer_handlers_test.go:137` 红；M-2 只记日志不拒但照样写 ⇒ 红在 **`:145` 的
  "injected ModeWriter was called"**，即只断字符串的那一发放过去、本票的判据接住；还原后 `grep -c MUTATION` = 0）；
  四数复绿 windows `go test -count=2 -v ./internal/panel/` = `152/92/0/0`，**Docker 真跑 linux = 同一读数**，
  `bash scripts/wisp-cli-tests.sh` = `98/51/0/0`；`sh scripts/d22scan.sh` clean rc=0、`gofumpt -l` 三枚改动文件零输出。
  证据：`docs/evidence/s1/114-ac2-native-gate.md`（含"没做到的"六条，第一条＝**通路仍不存在、`rt.modeWrites` 生产零调用者**）。
  next=票 33/35 落宿主那一跳时必须调 `rt.modeWrites`（并解冻 `cmd/wisp/*_test.go` 补"Confirm 非 nil"的装配用例）；AC#3 在那之后。
- 2026-10-14 16:0x +08 `agent-ticket114-a5` did=九格拆分＋只落干净子集（写腿；实现者；⛔ 未翻框/未动 `Status:`、票名/零 push；窗口内 HEAD 由他腿推进 `037bfb7f`→`88c7e752`）。交件＝`.scratch/wisp/probes/114/a5/00-split-and-landing.md`（拆分表全文＋逐格"改前/改后"读数，78 行）＋`logs/01…06`（6 枚，每枚含原命令行）。**①拆分表**：九格逐格给"要写的 `file:line`／面／现读三态"，零枚被已收票整格做掉。**②落了 3 格仪器**（均非前端非契约）：**AC#4**＝`internal/panel/composer_test.go` 的 `rendererSourceFile` 扩为 d22scan `isTextFile` 的**超集**（+`.json .md .txt .yaml .yml .toml .go .svg .vue`、留 `.mjs`，只宽不缩）＋新对齐尺 `TestDoorTwoReadsEveryTextTypeTheScannerWalks`（现抽 `main.go:1312` 15 类＋`.go/.json/.html` 三地板名＋合成 switch 子控）＋`.json` 种形；改前（overlay 只回退扩展名集）尺点名 **9 枚**未读类、`.json` 门同时逃过两枚仪器，改后全绿。**AC#9**＝结构钉根 `frontend/src`→**整棵 `frontend`**（与门二同根、任何加载集的超集；精确加载集留界面侧 J3）＋新种形 `frontend/public/ac92-public-door.js`（R-92b-2 第二形）；改前（overlay 只回退根）种形两断言红、真树不种形仍绿（"绿与事实无关"现场复现），改后真树 `109 files scanned, 2 host call sites`（窗口内 `frontend/` 被界面侧写到 111 枚，两发都绿）、种形点名。**AC#8**＝**253-r5 未覆盖此方向**（其尺自陈只读 `internal/panel` 生产源；它 gates declared/guard/router 三处 Go 侧登记脊柱）；新增派生门 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage`（名字经 253 的 `parseFileInbound` 从 `bridge.go` 现取，逐一要求 `frontend/src/lib/panel.ts` 字面量；自带盘上种形正控）；盘上种读数＝仓外模块拷贝把 `config.unmirrored114` 种进 `bridge.go` ⇒ 旧套件（HEAD 版）**PASS**（盲）、新门 **FAIL** 逐字点名（A713：overlay 换不到 `os.ReadFile`，故用盘上种，⛔ 未拿 overlay 绿抵账）。**③AC#7 读数**＝build rc=0／vet rc=0／gofumpt -l（`.exe` 全路径）零行／`bash -n` 两脚本 rc=0／`sh scripts/d22scan.sh` rc=0 clean（scopes 逐数在 log 05）／长度门 VERDICT GREEN（60=60）／`internal/panel` 整包 `223/143/5/0`，5 红与 HEAD 版对照件**完全同集**（3 枚前端耦合既有红＋2 枚 35-r8，均非本波）／**POSIX 三合包** `621/441/5/5 rc=1`（FAIL 5 名与 windows 逐名同；SKIP 5 名逐条在 log 06c；首跑因网络搁浅 RUN=0 未当读数）。**④本波⛔没做（具名）**：**AC#3**（`cmd/wisp/panel_inbound.go:254/:264` 两枚 `Confirm: nil` 未接、`rt.modeWrites` 死码未动、未开真窗——J1"常驻腿拿什么举 C18 卡"无一处被具名 ⇒ 未定义即停）；**AC#6、AC#10＝`skipped=frontend(owner-delegated)`**（不写不勾；AC#10 补量：`AddHostObjectToScript` 在本仓零调用点、真门是第三种 `w.Bind("wispDispatch")`，若改判要落 Go 侧钉另派）；**AC#11 停手上报**（`.github/**` 契约邻接＋111/85a 让渡未查）。**⑤纪律偏差自报**：产码 commit `d8dd4b38`（两枚测试件 +219/−87）**意外带上编排者先前已暂存未提交的 `.scratch/wisp/issues/253-…holes.md` 旧名删除**（253 改名 `-done` 的未完成半边；内容无损：`-done` 收口版在 `88c7e752`、旧名 blob `02d5632…` 在历史；⛔ 未 amend；教训＝下一腿用 `git commit -- <显式 pathspec>`）；交付件（本票尾／probes/114/a5）与产码分两笔提交。next=编排者裁 AC#3 的 J1 形＋AC#10/AC#11 去向；AC#5 的 frontend 半支按仪器读数已在件 §2.5。
- 2026-10-14 16:4x +08 `agent-ticket114-v1` did=验收腿收件（非实现者；只读＋overlay／仓外真种盘＋重跑；锚 `9b51859b`，窗口内 `20-r2` 在飞）：**四格判语**＝**AC#4 成立**（自造"单删 `.json`"变异⇒新对齐尺红且逐字点名；全量回退旧 7 类⇒点名 9 枚；旧集下 `.json` 种形同时逃过结构钉 unknown 2≠3 与门二少点名＝"改前逃过两枚仪器"复读）／**AC#7 成立（抽核口径）**（panel 整包 `223/143/5/0` 与 `88c7e752` 对照件 `220/141/5/0` 红集逐名同集；vet rc=0／gofumpt -l 零行／`bash -n`×2 rc=0 亲跑；d22scan 全跑／长度门／POSIX 三合包未复跑具名）／**AC#8 成立**（仓外树真种盘页侧 `config.set`⇒`config.sets`⇒新门 FAIL 逐字点名 `config.set`，而同一棵树 253 尺 `6/6/6` 与 N7 PASS＝旧盲新有牙；253 盲据＝其尺 `:19-20`"reads nothing outside internal/panel's production sources"）／**AC#9 成立（钉半支**；仓外树旧根 overlay 对种门树仍绿 `89 files scanned, 2 host call sites`＝"绿与事实无关"复现、新根同树 FAIL 点名 `public/ac92-v1-door.js`；真树 `113 files scanned, 2 host call sites` 绿；junction 只登记不演练）。**两停手格复核**＝**AC#3 停得对**（`composer_handlers.go:133` `ErrNoL2Confirm`／`panel_inbound.go:254/:264` 两 `Confirm: nil`／`run.go:340/:677` `rt.modeWrites` 零读者死码；J1 无具名）／**AC#11 停得对**（`ci.yml:1026` 后无 `git diff`、`grep -c exit-code`=0；票 111/85 双 open 未让渡；措辞更正＝"契约邻接"→"他票占用未让渡＋CI 治理面"，`.github` 不在明列冻结）。**拆分表抽核**（`shuf -e AC#5 AC#6 AC#10 -n 2`⇒**AC#10、AC#5**）＝AC#5 判对（一处更正：`config_handlers.go:108` 拒 7 族＋2 叶名，非"两族"）；**AC#10 ②面判不成立⇒停手上报**（普查 `:198` 列其为仪器格、本行① 落点全 Go、真门 `w.Bind`@`panel_host_windows.go:842`＋const `:80` 而钉正则 `composer_test.go:441` 只认 `.postMessage(`；③"本波不落"站得住）——请编排者按已挂的"AC#10 去向"重裁面签。交件＝`docs/evidence/s1/ticket114-114-v1-accept-v1.md`＋`probes/114/v1/`（锚＋logs 5）；仓内零突变（两测试件 sha256 前后同值）、框未翻、零 push。next=编排者裁 AC#10 面签＋AC#7 未复跑三处的接受度（d22scan 全跑／长度门／POSIX）。
- 2026-10-14 17:0x +08 `agent-ticket114-r6` did=写腿落地 AC#10「可达性钉认第三装法」（实现者；⛔ 未翻框/未动 `Status:`、票名/零 push；起手锚 `6c5515e8`，窗口内他腿 `201-a1` 落件、`frontend/` 被界面侧续写）。**①落点**＝新件 `internal/panel/bound_dispatch_door_114_test.go`（530 行；产码 commit `3b6f170f`）／`TestTheBoundDispatchDoorIsReachableFromThePage`，能力形两侧同名：host 侧 AST 走 `cmd/wisp` 产码取每个 `w.Bind(<名>,回调)`，`<名>` 须为包级字符串常量（`:80` `panelDispatchBinding="wispDispatch"`→`:842` `w.Bind` 即此；单参 `ws.Bind(run.paths)` 记 skipped；literal 的 `wispProbeRT` 不入 reviewed），page 侧要求该常量值＝`frontend/src/lib/panel.ts` 以 `window.<名>` 读到、且整棵 `frontend/` 该门读只许在 panel.ts（`//`／`/* */`／`<!-- -->` 注释剥掉，`:329` 的 prose 顶不了数）；floor＝reviewed 空即红。⛔ 未动 `cmd/wisp/**` 产码一字节、未动 `frontend/**`。**②三向读数（`probes/114/r6/logs/` 逐字）**＝基线绿 `bound door held: 36 cmd/wisp production file(s), 2 w.Bind door site(s) […:842 binds "wispDispatch" via panelDispatchBinding; …:892 binds "wispProbeRT"…], 27 window.* read site(s)… panel.ts:350 window.wispDispatch`；仓外真种盘（`git archive HEAD` 树，A713）三突变各自 FAIL 且点名：m1 摘 `w.Bind`〔span `[42508,42619)` 换 `error(nil)`〕⇒"bind no door through a w.Bind(<package-level string constant>…"＋剩余站点清单；m2 常量值换名⇒`binds "wispDispatchX" through w.Bind, and … never reads window.wispDispatchX`；m3 页读换名（宿主不动）⇒`… binds "wispDispatch" … never reads window.wispDispatch - … panel.ts:350 window.wispDispatchX`；同树上旧 `.postMessage(` 尺三发全 PASS（m3 系页面侧改动仍在其射程外）＝"旧盲新牙"现场；③既有牙不回退＝`red as required: 3 outside call sites, 1 assembled routes, 3 unknown literals…` 逐字不动。**③抽核**＝`internal/panel` 整包同树 A/B `225/142/7` vs `223/141/7` **红集 diff SAME SET**（树 7 红＝live 既有 5＋两枚树工件：树无 `.git`、导出件 CRLF 而 fixture 尺分隔符写字面 `" -->\n"`）；live 整包 `225/144/5/0 rc=1`，5 红与 v1/a5 记录**同集**。**④门禁**＝build rc=0／vet rc=0／gofumpt -l 零行／`sh scripts/d22scan.sh` clean rc=0／长度门 VERDICT GREEN（60=60）／`-run` rc=0。**⑤具名没做**＝未跑 `go test ./cmd/wisp`（314 死区，派单禁）／POSIX 三合包／真窗；`panel.ts:146/329` 注释自陈"两种装法"那半按派单归前端；AC#10 判语与翻勾留非实现者。交件＝`.scratch/wisp/probes/114/r6/00-anchor-and-gates.md`＋logs 4；产码/交付两笔提交、显式 pathspec、仓内零突变（`cmd/wisp/panel_host_windows.go` `45c1059b…`／`panel.ts` `2b7a9fd0…`／`composer_test.go` `d2b2e848…` 起手＝收手）。next=编排者派非实现者复核本钉（两侧同名能力形＋三突变点名）后裁 AC#10 勾。
- 2026-10-14 17:1x +08 `agent-ticket114-v2` did=验收腿收件（非实现者；只读＋仓外真种盘跨树突变＋重跑；起手锚 `7e8ac3fa`，窗口内 `318-r1` 在飞）：**AC#10 判 成立**（编排者重裁口径"让可达性钉认第三装法"）——第三装法被两侧同名读盘钉住：host 36 枚产码取 `:842 w.Bind(panelDispatchBinding,…)`（常量 `:80`＝`"wispDispatch"`）、页 `panel.ts:350 window.wispDispatch` 且全 `frontend/` 只许 panel.ts 读。**我自己的 9 发种形**（仓外树 `D:\tmp\wisp114v2\tree` 真盘＋sha256 逐发；⛔ 不用 overlay）＝8 红 1 绿：t1 摘常量绑 `w.Bind` 红（floor 逐字、剩余站点列 `:886 wispProbeRT`）／t2 常量声明搬家 `run.go` **绿**（跨文件解析成立）／t2b 声明与调用拼写不一致红（两半各点名 `names … no package-level string constant…`＋floor `"<unreadable>"`）／t3 页真读改注释红（`never reads window.wispDispatch`，prose 在 :329 同样剥掉）／t4 第二文件读门红（点名 `config-view.ts:168`）／t5 reviewed 空形 A（实参换 literal、两 Bind 都在）红／t6 reviewed 空形 B（binds 全空）红／t7 单参 skip 翻双参红（`:311` 边界 fail-closed）／t8 `cmd/wisp` 根消失红（`not walkable`）；**旧 `.postMessage(` 尺同树两发 PASS**（t1/t3，`113 files scanned, 2 host call sites`）＝旧盲新牙现场；**射程四条逐条独立核**（skipped 单参／literal 不入 reviewed／注释剥离／only-panel.ts，各抽一行现读＋对应种形）。**基线**＝live 新钉 rc=0（绿句逐字）＋live 整包 `225/144/5/0 rc=1` 5 红与 r6/v1/a5 **同集**（抽 2 名逐字）；同树 A/B（有新钉 `225/142/7` vs 移走 `223/141/7`）**FAIL 名集 diff rc=0＝SAME SET**。**门禁**＝build rc=0／两枚 `-run` rc=0／`sh scripts/d22scan.sh` clean rc=0（ban #8 internal/=556）／长度门 `VERDICT GREEN`（60=60，tracked=12876）。**留档句**①常量窄口径（literal 不入 reviewed，可预见唯一摩擦）②AC#10 原信 `AddHostObjectToScript` 半支：Go 生产零命中为真、本钉仅在红句散文陈述、未成门（注释半归前端，要门须另派）。**具名没做**＝`go test ./cmd/wisp`（派单禁）／POSIX／真窗／r6 的 span 数值未复核（其结论由我自有九发独立复核）。交件＝`docs/evidence/s1/ticket114-114-v2-accept-v1.md`＋`probes/114/v2/`（锚＋logs 4）；仓内零突变（5 枚关键件 sha256 起手＝收手 `45c1059b…`／`2b7a9fd0…`／`d2b2e848…`／`6924b8e1…`／`016b6943…`）、框未翻、零 push。next=编排者按判语＋留档句①②裁 AC#10 翻勾。

## ★前提更正（2026-10-24，凭 `114-r1`/`114-r2` 现跑）

> **本节只更正前提**：⛔ 不改上面任何一句原文、⛔ 不勾任何框、⛔ 不动 `Status:` 与票名（那些是编排者的动作）。
> 原文留着，是为了让后来人看见"当时为什么那么写"；**引用时应引本节这一份读数**。
> 下表每一行的读数都由写腿 `114-r3` **自己现跑**（尺与原始输出见下面「现跑的尺」），⛔ 不是抄的。

**原文当时为什么为真**：2026-09-21 立案时，面板→Go 的入向那一跳**整层不存在**，所以
`ParseComposerRequest`／`HandleModeRequest`／`perm.Store.Set` 的生产调用者确实都是 **0**，票 33（C27 宿主）
与票 35（C17 桥）也确实是 ready-for-agent（`docs/evidence/s1/114-ac1-status-table.md` ①.1 量的正是那个状态）。
**那一跳后来真的落地了**（宿主 `w.Bind` 门 + `wisp panel-inbound` 管两条入口），三处说法因此过期。

### 三处过期句 ↔ 现跑读数（原文逐字，⛔ 未改）

| 位置 | 原文（逐字） | 现跑读数 |
|---|---|---|
| **标题**（`:1`） | 「`ParseComposerRequest` 在**生产里零调用者**」 | **≠ 0**：非 test 调用点 `internal/panel/composer_dispatch.go:178`（`req, err := ParseComposerRequest(raw)`；定义在 `internal/panel/bridge.go:134`） |
| **`:21`** | 「`ParseComposerRequest` 的**生产调用者 = 0**，`perm.Store.Set` 的**生产调用者 = 0**」 | **两半都不是 0**：①`ParseComposerRequest` ⇒ `composer_dispatch.go:178`；②`perm.Store.Set`（`internal/perm/store.go:189`）⇒ 全仓非 test **恰 1 处**：`internal/panel/composer_handlers.go:141`（经本地 `ModeWriter` 接口；`internal/tools/mode.go:20-22` 自陈写侧不在 tools，读侧只有 `PermissionMode`） |
| **`:44`** | 「⚠ 门已立，**通路仍不存在**（P0 地界：宿主那一跳留票 33/35）」 | **通路已在**，且两条产线入口都到同一枚 `ComposerDispatch.Handle`（`internal/panel/composer_dispatch.go:176`）：① **页面门**＝`cmd/wisp/panel_host_windows.go:929` `w.Bind(panelDispatchBinding, …)` → `:930 m.dispatchRaw` → `:948 m.disp.Handle(ctx, raw)`；② **CLI 管**＝`cmd/wisp/panel_inbound.go:173` `reply, err := disp.Handle(ctx, raw)`（`wisp panel-inbound`）。router 的 `dispatch` 在 `composer_dispatch.go:204` 把 `MethodModeRequest` 这一 case 交给 `HandleModeRequest`（实现 `composer_handlers.go:109`） |

### 现跑的尺（四条，逐条可重跑）

```
git grep -n "ParseComposerRequest(" -- '*.go' ':!*_test.go'
# ⇒ internal/panel/bridge.go:134(定义) ＋ internal/panel/composer_dispatch.go:178(调用)
git grep -n "HandleModeRequest(" -- '*.go' ':!*_test.go'
# ⇒ composer_dispatch.go:73(接口) ＋ :204(调用) ＋ composer_handlers.go:109(实现)
git grep -n "\.Set(ctx" -- '*.go' ':!*_test.go'
# ⇒ composer_handlers.go:141（全仓唯一一处；锚的是调用形 ".Set(ctx"，⛔ 裸符号名）
git grep -n "disp.Handle(ctx, raw)" -- 'cmd/wisp/*.go'
# ⇒ panel_inbound.go:173 ＋ panel_host_windows.go:948
```

### 仍未变的两条（⛔ 把本节读成"票 114 的活变小了"）

1. **`Confirm` 的 nil 语义一字未变**：处理器读的仍是 `h.Confirm` 的 nil-ness（`composer_handlers.go:131`）；
   `wisp panel-inbound` 那条入口**仍递 nil**（`panel_inbound.go:236` → `newComposerDispatchChain(…, nil)`，
   理由自陈在 `:33-40`），所以**管这一路的变宽仍被拒**。常驻面板那条入口递的是
   `panelModeConfirm(ra.confirmModeSwitch)`（`cmd/wisp/panel_resident_windows.go:233`），与 `wisp run` 同一条腿。
2. **AC#3 仍未做**：本节只更正"有没有人调"，⛔ 不构成反向对照（真机 L2 卡片 + 300s 超时＝拒）的证据；那一格仍空。

- ★**同一批前提的另外两处已由 `114-r3` 就地改准**（写腿；在本票面之外的两枚产码注释／字符串）：
  `internal/panel/composer_handlers.go:41-45` 的「Nothing calls this handler yet … tickets 33/35 are still
  ready-for-agent」⇒ 改成点名两条产线入口；`:132` 括注「R20/M4：只有变宽要确认」⇒ 改回机主原话
  「**只有全自动要确认**」（出处＝`docs/reports/pending-and-issues.md:3019` 的裁定 `R20` 逐字；旧措辞在
  `PLAN.md`／`specs`／台账 **0 命中**）。逐处改前／改后与现跑读数见 `docs/evidence/s1/114-r3-impl.md`。
