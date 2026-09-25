# 77 — Frontend scaffold: React + TS + Tailwind + shadcn, vendored react-bits / beautifului components, sharing C21 tokens

> **⛔ 2026-09-23 09:22 编排者按 owner 指令冻结本票（不改文件名，只改状态位）**：
> 原话「我是想让别的 agent 干前端部分的」⇒ 本票**余下四格 AC#1 / AC#3 / AC#4 / AC#6 的正解全在 `frontend/` 下**，
> 我编队**不再派发、不再代写**。已入库的三格（AC#2 / AC#5 / AC#7）**保留不回滚**。
> 台账依据 `A102`，规则依据 `issues/README.md` 规则 7。
> ⚠ 唯一例外（不属前端、可继续）：AC#4 欠在 `tools/d22scan` 里的 ban #8 `frontend/` 作用域那一行是**编排者地界**，
> 它是扫描器配置不是界面代码。

**Status:** **blocked（09-23，owner 把前端收走了）** / 上一状态 **in progress**（AC#2 / AC#5 / AC#7 已交 `63ef895`；AC#1 PARTIAL）/ 原 claimed（2026-09-21 14:0x，`agent-ticket77` 已认领并开工，
**2026-09-21 16:5x 由 `agent-ticket77b` 接续**（前任撞 150 轮上限，断点见编排者 16:2x 那条与我的 16:5x 那条），
按 owner 拍的板：**基座 = Beautiful UI**，一期 **React Bits 零代码进树**）。
> **2026-09-21 15:1x 第三任 `agent-ticket77d` 的 Status 快照（不覆盖上面那行，只追加）**：
> AC#1 / AC#3 / AC#4 的**判据数字全部到位**（见 Progress log 14:5x 与 15:1x 两条），
> 四框**仍不勾**，且每一框只剩同一段"最后 mile"：**真 WebView2 host（票 33）与 bridge 推送（票 35）不在这棵树里**
> ——AC#4 另欠 `tools/d22scan` 里 ban #8 的 `frontend/` 作用域那一行（编排者地界，见 15:1x 条末）。
> AC#6 的 `lint-frontend` job **已写进 `ci.yml` 但无 run id**（本票不许 push）⇒ 框不勾，
> **明写"欠 push 后复跑"**，不许用"本地 6 步 rc=0"替代（`npm ci` 那步本地**故意没跑**，理由见 15:1x 条）。
上一状态：**unblocked**（2026-09-21 13:5x，owner 第二次指令 ⇒ 裁定 **R19**）——
他挑完了动画组件（12 条清单在 R19 表里），并且给了**排序**：
**第一版只用 Beautiful UI 那套 agent 组件为基座（+ shadcn 基础件 + Tailwind），React Bits 二期再叠。**
⚠ 一期因此**不引 React Bits 的任何代码** ⇒ 它的 **MIT + Commons Clause** 许可审查推到二期
（但**引入前必须 owner 复核**这条留在 AC#7 的台账判据里，不许因为"二期"就忘掉）。
"等 owner 逐屏给参考图"那条前提**早已作废**（A47），这次也不再有任何等他答的视觉项。

**一期范围 = R19 的排序 + Q-24 的判据**：只做 **L2 确认卡 + 面板骨架**，其余屏后面按票排。
**两条一期的硬约束（先于任何视觉效果）**：
1. **桌面那颗球不许换实现**：Win32 分层窗口 + Direct2D 画的，D32 门 = 休眠 **CPU ≤0.5% / RSS ≤25MB**
   （`PLAN.md:527`、`:1032`）。R19 表里的 `Agentic Ball` **只能**当"面板里的次要指示"，
   **绝不能**变成原生球的替代品。
2. **全屏 WebGL 氛围层同屏最多 1 个**（`Glass Flow`/`Aura Blob`/`Neural Float`/`Fog Sphere` 四条互相排斥），
   且**面板一隐藏就必须销毁**（不是"暂停渲染循环"就算）——WebView2 的 CPU 是要计进 D32 预算的。
   二期落地时这条要有用例钉住"隐藏后 CPU 回落"。

**Type:** feature / infrastructure (first real frontend code in this repo)
**Blocks:** 面板与 L2 确认卡的所有后续票 · **Blocked by:** ~~owner 的挑选~~ **已解除**；
票 78（Linux `go vet`）已由 **A54③** 查清"那条命令在本机永远 rc=1"的仪器问题、CI 侧样本仍挂票 70；
Q-23/Q-24（我的推荐已生效，见 R19 第 1 条）
**Owner ruling:** **R18**（栈与 `design/` 降级）+ **R19**（基座=Beautiful UI、React Bits 二期、12 条清单与我给的取舍建议）

**Spec refs:** `PLAN.md:981`（栈本就是 React+TS+Tailwind+shadcn/ui，宿主 `jchv/go-webview2`）、
`PLAN.md:1032-1036`（界面分层）、`PLAN.md:1038-1040`（**共享 token 是硬要求**）、`PLAN.md:1044`（**前端必须无状态**）、
C21、D23/§17 零 emoji 图标、D29 人工视觉签收、ban #6 / ban #8

## 前提事实（已核，不必重查）
- 仓内**没有任何前端工程**：无 `package.json`、无 `frontend/`；`design/` 只有 11 屏静态原型 HTML + `assets/tokens.css`。
- **悬浮球与 L1 提示条保持原生**（Direct2D）；本票只做 **WebView 侧**：面板（命令/结果/配置/历史）+ **L2 强确认卡**。
- 许可：`beautifului.dev` → 上游 `TurboKach/ai-native-react-components`（**MIT**，19 个组件，含
  approval flows / thinking states / tool traces）；`react-bits`（47.7k★，**`MIT + Commons Clause`**）；
  `shadcn/ui`（MIT）。三家都是 **copy-paste** 交付 ⇒ 按 **Q-20** 的推荐**vendored 进仓**，每文件头注明来源与许可。

## What to build
1. **工程骨架**：`frontend/`（Vite + React + TS + Tailwind），构建产物由 Go 侧 `go:embed` 打进 `wisp.exe`，
   **不引入运行时 CDN**。宿主用已冻结选型 `jchv/go-webview2`（`PLAN.md:1008`）。
2. **token 单一来源**（这条是本票的**验收核心**，不是附带工作）：Tailwind theme 必须**由 C21 那份定义生成**
   （`docs/contracts/` 的 token → 生成 `tailwind` config / CSS variables），**禁止在 `frontend/` 里手写第二套颜色/圆角/字号**。
   并把票 74 交付的三方对账 `TestC21TableColourRowsMatchTokensCSS`（`tokens.css` = 表 = `tokens.go`，实测 78/78）**扩成四方**，
   第四方就是生成出来的前端 theme。**变异检验**：只改前端一处颜色 → 那条用例必须红。
   理由：`PLAN.md:1038-1040` 把"共享 token"列为硬要求，破了它的表现是"看起来像两个东西拼的"——
   那正是**没有机器检查时人会以为没问题**的那种缺陷。
3. **L2 强确认卡**：用 beautifului 的 approval 组件族，卡片必须显示
   **完整命令/参数 + 风险等级 + 调用链 + 拒绝或放行的原因**（Q-17/Q-23 的落点；
   原因枚举依赖票 20 第 7 框那条 D22 项，**若尚未落地就把原因字段做成可选并显式标注"信息不足"**，
   绝不允许把"未知原因"显示成"无风险"）。
4. **无状态**（`PLAN.md:1044`）：WebView 销毁即全丢；一切从 Go 侧读（历史 SQLite、配置 TOML）。
   ⇒ 不得引入前端本地持久化（localStorage 等）。
5. **同批武装两条门禁**（今天的 A40④ 判据在此生效，**不许分两次做**）：
   - `ban #6 panel-approval` 现在扫描面是 `frontend/=0 [NOT COVERED]`，票 71 已把"声明作用域走到 0 文件"变成
     **rc=2 致命** ⇒ 第一批 `frontend/` 文件落地时，`ban #6` 必须**真的扫到文件**并通过
     （它禁的是前端出现 `approval.decide`——放行只能在原生侧）。
   - `ban #8` 零 emoji：vendored 组件里若带字形，**要么 ASCII 化、要么显式不把它们纳入扫描面并登记理由**，
     两种都行，**但必须选一种并写进 commit**。参考今天的实测：一次覆盖面扩展没同批改存量，
     HEAD 的 lint 红了 13 分钟。
6. **CI 新 job**：`lint-frontend`（`npm ci && npm run typecheck && npm run lint && npm run build`）。
   ⚠ D22 的"每个 job 都不许被跳过"（`PLAN.md` D22 mode-6）同样适用于新 job；且**必须排在票 78 修好之后**，
   否则会复制 A44① 那个病症——**步骤存在但从未产出结论**。判据：交回一个**真实 run id 与 job id**，
   显示这一步跑过并给出过红或绿。

## AC（1:1 裁决表，与框逐条对应；能力类框必须回答"生产里谁调用它"）
- [ ] **AC#1** `frontend/` 可构建，产物被 `go:embed` 打进二进制，`wisp.exe` 在无 node 环境的机器上能拉起面板。
      判据：**剥光 PATH 里的 node/npm** 后跑一次真拉起（这条来自我们自己踩过的"签收命令只在开发机能跑"）。
- [x] **AC#2** token 四方对账用例落地 + 一次"只改前端一处颜色 ⇒ 用例红"的变异（grep 证落盘 → 跑 → 还原 → 证还原）。
- [ ] **AC#3** L2 确认卡渲染真数据：由 `internal/risk` 的真实决策对象驱动，**不是**组件库 demo 的假 props。
      ⚠ 能力类判据：**说清生产调用者是谁**（哪个装配根、什么事件流），否则按 A33 只能记 PARTIAL。
- [ ] **AC#4** `ban #6` 与 `ban #8` 对真实 `frontend/` 的作用域台账：贴出扫描器自报的**逐作用域文件数**（>0），
      并给出一次故意在前端写 `approval.decide` ⇒ **rc=1 且点名文件** 的真红。
- [x] **AC#5** 前端无任何本地持久化；重启 WebView 后一切从 Go 侧恢复到同一状态（一条真机或半真机用例证明）。
- [ ] **AC#6** CI 新 job 有真实 run id + 结论（见前提 6）；**不允许用"本地跑过了"替代**。
      > **2026-09-25 11:4x 前端会话就地标注（只加这一段，原句不改，框仍不勾）**：这一格的**唯一缺口是
      > 一次推送**，不是"做不了"、也不是界面侧欠工。编排者 11:0x 明示他按住了未推的积压（本会话现量
      > `git rev-list --count origin/dev..HEAD` = **152**，其中含本会话今天的 `d6c52ef` / `a21336e` /
      > `53a1359` / `35c3724`），并给了口径：**标「待编排者推送」，别标「做不了」**。
      > 推送之后这一步要取的是 `gh run view <id> --json jobs` 里 `lint-frontend` 的**步级**结论
      > （六步逐名可引，见 `:17` 与前提 6），不是 job 级的绿。
      > ⚠ 另有一枚会让这一格**红着到手**的在先：今天 07:2x 的 `3b59512` 把 `<=` 写进了 JSX 文本，
      > `npm run typecheck` 从那时起红到 11:1x（`composer.tsx(170,62) TS1003`），已由 `a21336e` 修好。
      > 所以**推送落在 `a21336e` 之后**时第 2 步应当是绿的；若有人单独推更早的 sha，那枚红是本会话造的，
      > 别记成接线问题。
- [x] **AC#7** vendored 清单：`frontend/VENDORED.md` 逐文件列 来源仓库 / 组件名 / 许可（含 Commons Clause 标注），
      且 `git grep` 能在每个 vendored 文件头找到出处。
      ⇒ **R19 之后的一期口径**：一期**只会出现 Beautiful UI（MIT）与 shadcn/ui（MIT）两类来源**；
      React Bits 一条代码都不许进树。台账里请**明写一行"React Bits = 二期，许可为 MIT + Commons Clause，
      引入前需 owner 复核"**，让它不会因为"这期不用"而从账上消失。

## 硬规矩（本项目 24 小时内的事故换来的）
- 只许你动 `frontend/**`、`docs/contracts/**`、`.github/workflows/ci.yml`（新增 job 那一段）、
  以及 `internal/panel/**` 的 embed 接线。**不许动 `internal/ball`、`internal/risk`、`internal/tools`**（他人领地）。
- `design/**` **只读**：降级为参考不等于可以删它（票 65/68 的历史证据还指着它）。
- 提交一律 `git commit -q -F - -- <显式路径> <<'MSGEOF'`，**heredoc 必须加引号**；禁 `git add -A`；
  提交前 `git diff --cached --name-only` 核清单（A31/A36）；禁 `--amend`/`reset`/`rebase`/`stash`（A34）。
- 不 push（编排者串行推双远程）。第一个 checkpoint commit 必须落在前 15 次工具调用内；
  每次 commit 同步票面 Status + 勾框 + 一行 `next=`（A30/A32：票面落后 HEAD 会让整票被重做）。
- 未证明的框**留着不勾并写明原因**；`t.Skip`/放宽阈值/加豁免换绿，一律视为未完成。

## Progress log（append-only）

- **14:0x `agent-ticket77` 认领 + 第一枚 checkpoint**（本 commit 只改票面）。
  现场事实核对：`frontend/` 与 `docs/contracts/` **均不存在**（票面"前提事实"仍成立）；
  token 现状 = `design/assets/tokens.css` + `internal/ball/tokens.go` + `internal/ball/tokens_table_test.go`
  的三方对账 `TestC21TableColourRowsMatchTokensCSS`；`internal/panel/` 只有 `doc.go`（embed 接线是空地）；
  `cmd/wisp/` 最近 5 手 = 票 78/67/70 的 vet・ASCII・build-tag 收敛，**本票对它只加不改**。
  领地声明：本票动 `frontend/**`、`docs/contracts/**`、`internal/panel/**`、`cmd/wisp/**`（仅新增文件）、
  `.github/workflows/ci.yml`（仅新增 job 行）；**不碰 `internal/ball/**`、`internal/risk/**`、`tools/d22scan/**`**。
  `next=` AC#1 —— 立 `frontend/` 骨架（Vite+React+TS+Tailwind）+ `internal/panel` 的 `go:embed` 静态服务
  + `cmd/wisp` 新增接线文件，跑真构建与真拉起拿 exit code。

- **14:5x `agent-ticket77` AC#1 落地（构建 + go:embed + 无 node 真拉起）**。
   `npm run build`（Vite 8 + React 19 + TS 6 + Tailwind v4）rc=0：产物 index.html 1043 B、
  assets/index-BQ3J3QdO.css 37955 B、assets/index-DIrih6dy.js 264929 B（清掉 create-vite 模板遗留的
  两个 svg 之后重跑的最终数字；css 哈希变、js 未变）。
  `go build -o build/wisp.exe ./cmd/wisp` rc=0（gofmt -l 空、go vet ./internal/panel/... ./cmd/wisp/... rc=0）。
  **剥光 node/npm 的真拉起**（`PATH=/c/Windows/System32:/c/Windows:/usr/bin`，`command -v node` 为空）：
  `wisp.exe panel-assets` rc=0 = "panel assets embedded: 4 files, entry=index.html built=true"；
  `-manifest` rc=0 逐文件 sha256 前缀；`-render /` rc=0 输出 1043 B（引用 ./assets/...，相对路径）；
  `-render /assets/index-DIrih6dy.js` rc=0 输出 264929 B，Content-Type text/javascript；
  `-render ../go.mod` **rc=1 拒绝**（越界路径 fail-closed）。`wisp.exe version` rc=0。
  ⚠ AC#1 的**未尽半**：WebView2 窗口本体是票 33 的地界（`internal/panel/doc.go` 明写 host 由票 33 交付，
  go.mod 里至今没有 webview 依赖）。本票交付的是"二进制自带页面、无 node 可取字节"这一段，
  **不是窗口被真的创建出来** ⇒ AC#1 按 A33 记 **PARTIAL**，勾框留着。
  embed 接线：`frontend/embed.go`（`//go:embed all:dist`）+ `internal/panel/assets.go`
  （Resolve/Manifest/Built，纯 Go 无平台标签）+ `cmd/wisp/panel_assets.go`（新文件）。
  `cmd/wisp/main.go` 只加 6 行（usage 一条 + switch 一个 case），逐字理由写在 commit message 里。
  dist 的 gitignore 口径：产物不入库，但 `go:embed` 的模式在干净树上必须可满足 ⇒
  `frontend/.gitignore` 用 `dist/*` + `!dist/.gitkeep` 留一枚锚点，`Assets.Built()==false`
  时主机必须显示"资源未构建"，锚点永远不可能被当成页面渲染。
  `next=` AC#2 四方对账用例 + AC#7 `frontend/VENDORED.md`；之后 AC#4 台账
  （**已知雷**：`frontend/` 一存在，`tools/d22scan` 的 driftedAbsentScope 就 rc=2 要求把 ban #6
  翻成 live:true，而 `tools/d22scan/**` 不是本票能动的地方 —— 见 Progress log 下一条的数字）。

- **15:2x `agent-ticket77` AC#2 + AC#5 + AC#7 落地，AC#3/AC#4 各带回一处硬事实**。
  **AC#2 四方对账**：新用例 `TestC21DesignTokensFourWayAgree`（`internal/panel/tokens_fourway_test.go`）
  一次跑核 P1 `design/assets/tokens.css` / P2 `docs/evidence/s1/c21-native-tokens.md` /
  P3 `internal/ball/tokens.go` / P4 `frontend/src/styles/tokens.generated.css`：
  实测 `135 dark + 71 light tokens reconciled`、`legs 2-4 reconciled for 78 of 78 colour rows`、rc=0。
  ⚠ **口径偏差，先记在这里**：票面写的是"把 `TestC21TableColourRowsMatchTokensCSS` 扩成四方"，
  而那条用例在 `internal/ball/` —— 本票硬规矩明令不许动 `internal/ball`。四方因此落在**我自己的
  作用域里新增的一条用例**上（P2↔P3 那一腿在两边各核一次，语义等价、互不遮蔽），
  ball 的三方用例一字未改。**要不要把它并回 ball 由 ball 的地界决定，不由本票决定。**
  **变异检验（真红）**：只改前端第 42 行 `--accent: #86C2B9` → `#86C2B8` ⇒
  `--- FAIL` rc=1，两条腿同时点名（`tokens.css declares --accent = "#86C2B9" but ... carries "#86C2B8"`
  与 `native Palette.Accent holds "hex(0x86C2B9,1)" but the panel renders --accent = "#86C2B8"`，
  `78 → 77 rows`）。独立第二道 `node scripts/gen-tokens.mjs --check` 同一改动 rc=1。
  **还原证明**：`npm run tokens` 重生成后 `sha256sum` = `93a887626d6ebdb2...`（与变异前同一串），
  `git diff --stat` 对该文件为空，用例回到 rc=0。
  **AC#5**：`TestPanelFrontendIsStateless` 扫 `frontend/src` 20 个文件 × 7 类持久化 API
  （localStorage/sessionStorage/indexedDB/document.cookie/caches/serviceWorker/node:fs）= **0 命中**；
  `TestPanelSnapshotSurvivesWebviewRestart` 用真评估器的决策 → 432 B JSON → 丢弃全部内存态再读回
  `reflect.DeepEqual` 相等（票面允许的"半真机"形状；真 WebView2 归票 33）。
  **AC#7**：`frontend/VENDORED.md` 逐文件台账（Beautiful UI 8 个 + shadcn 4 个，各带来源/组件/许可/
  上游 commit `05dab2d2b5f1`/本地改动/挂载状态），**React Bits 单列一行**：二期、**MIT + Commons Clause**、
  **引入前必须 owner 复核**、并给了可重跑的反查命令 `git grep -in "react-bits\|reactbits" -- frontend/`。
  两条 vendoring 脚本进了仓（`frontend/scripts/vendor.mjs`、`vendor-shadcn.mjs`），台账可被命令复现。
  **AC#4 的两处硬事实**：
  (1) **第一个真命中是我自己的台账** —— `frontend/VENDORED.md` 为了说明 ban #6 而写了那个标识符，
      被 `TestFrontendNeverNamesAnApprovalDecision` 点名 `frontend/VENDORED.md:109`（`--- FAIL` rc=1）。
      处置=**改我们的措辞**，`tools/d22scan/**` 与 `allowlist.txt` **一字未动**（豁免清单只短不长）。
  (2) **planted 真红**：临时 `frontend/src/__ban6_probe.ts` 写入该标识符 ⇒ 扫描器点名
      `frontend/src/__ban6_probe.ts:1: [panel-approval] ...`，但**退出码是 2 不是 1**：
      `verdict()` 的 guard 3（`driftedAbsentScope`）排在 findings 判定之前，而 ban #6 仍登记 `absentOK`。
      同一个 rc=2 在**没有任何命中**的当前树上也一样出现（`scope ban #6 frontend/ examined 37 text files [NOT COVERED]`）。
      ⇒ **`frontend/` 一进树，D22 门禁就是红的，而解除它唯一的一行（`live:true`）在 `tools/d22scan/main.go`，
      不是本票能动的地界**（票面硬规矩 + 编排者指令一致）。扫描器自己的报错文案就是这个意思：
      "Flip it to live:true in the same commit that creates the tree"。**这条要编排者落地**，
      否则 A58 刚转绿的 lint job 会因为我这一票重新红。
      ban #8 一侧：它的声明作用域只有 `design/ internal/ cmd/`，**不含 `frontend/`** ⇒
      票面"选一种并登记"我选的是**自己武装**（`TestFrontendHasNoEmoji`，字符区间照抄 scanner 的 `emojiRe`），
      24 个文件 0 命中；上游 `tool-chips.tsx` 带进来的 2 个 U+2713 在 vendoring 时就 ASCII 化了（不是事后追改）。
  **AC#3 现状（PARTIAL，框不勾）**：`internal/panel/approval.go` 的 `NewApprovalCardView` 直接问
  `risk.NewRiskAssessor()`（与 `internal/tools/bridge.go:169` 同一个构造函数），真评估器实测
  `level=L2 rules=[R1 R8] reason="R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除）"`；
  原因缺失时 `ReasonKnown=false` 且有专门用例钉住"不得显示成无风险"（票 20 第 7 框那条 D22 项未落地的口径）。
  Go↔TS 契约由 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 双向核（10/3/3 个 JSON 键）。
  生产调用者：今天是 `wisp panel-assets -l2 <tool> -- <argv>`（无 node 环境实测 rc=0，输出上面那份 JSON），
  明天是票 35 的桥推送。**没勾的原因**：没有任何浏览器/DOM 证据能证明 React 真把这些值画成了卡
  （一期无 vitest/jsdom 用例，WebView2 本体归票 33）⇒ 按 A33 记 PARTIAL。
  门禁复跑：`go build ./...` rc=0、`go vet ./internal/panel/... ./cmd/wisp/...` rc=0、
  `gofmt -l cmd/wisp internal/panel frontend` 空、`go test ./internal/panel/ -v` **10/10 PASS 0 SKIP**。
  `next=` AC#6 的 `lint-frontend` job 新增行（只加我自己的 job，`ci.yml` 请编排者复核），
  以及编排者那一行 ban #6 `live:true`；AC#3 的 DOM 证据与 AC#1 的真窗口一起排在票 33/35 之后。

- 2026-09-21 16:2x（**编排者记：代理撞 150 轮上限死亡，本文件由我补断点，活没丢**）：
  它交回了 **3 框真绿**（`63ef895`：AC#2 四方对账 135 dark + 71 light token、78/78 配色行四腿全核，
  变异 `--accent #86C2B9 → #86C2B8` 真红并且**独立第二道** `gen-tokens --check` 同改动也红；
  AC#5 无状态 20 文件 × 7 类持久化 API 0 命中 + WebView 重启快照回读相等；AC#7 vendored 台账）。
  死前最后一句是它**自己发现自己的 vendoring 脚本有 bug**：
  `node` 的 `String.replace` 把 `$'` 当特殊替换模式 ⇒ 测试文件被复制成三份，它说"整份重写"。
  ⚠ 我已复核**树里没有留下三份的痕迹**：`internal/panel/tokens_fourway_test.go`(550) 与
  `frontend_hygiene_test.go`(319) 各自 `^package ` 只出现 **1 次**、无重复顶层 `func` 名 ⇒ 污染发生在提交之前，
  它自己修完了才提的（这条判据以后可以复用：**怀疑文件被复制成几份，就数 `^package` 与重名顶层声明**）。
  **未闭**：AC#1（PARTIAL：构建 + `go:embed`）、AC#3（L2 卡渲真数据）、AC#4（`ban #6`/`#8` 作用域台账——
  其中 ban #6 那半边**已被我拆出去成票 88**，因为它是扫描器侧的活、不该由建目录的人顺手改自己的考卷）、
  AC#6（CI 新 job）。
  **工作树里还留着它一枚未提交改动**：`frontend/package.json` 的 `"typecheck": "tsc -b --noEmit false
  --emitDeclarationOnly false"` → `"tsc -b"`（无害，但**归属是它**，接续的人要先决定留还是回）。
  `next=` 交回给接续代理：**先 AC#1 收尾（真构建 + 真 embed + `wisp.exe` 载入），再 AC#3，最后 AC#6**；
  AC#4 只留 ban #8 那半边在你范围内。

- **16:5x `agent-ticket77b`（接续代理）接手 + 那枚未提交改动的处置 = 留（有据）**。
  **处置**：`frontend/package.json` 的 `"typecheck": "tsc -b"` **保留**。判断不是"看着无害"，是**变异检验过它仍有牙**：
  往 `frontend/src/__typecheck_probe.ts` 写一行 `export const probe: number = "not a number";` ⇒
  `npm run typecheck` **rc=2** 并点名 `src/__typecheck_probe.ts(1,14): error TS2322`；删掉探针后 **rc=0**。
  另外两条事实支持"留"：`tsconfig.app.json` 与 `tsconfig.node.json` 都写着 `"noEmit": true`、
  `tsBuildInfoFile` 落在 `node_modules/.tmp/`，所以 `tsc -b` **不向仓内或 dist/ 吐任何产物**
  （跑完 `git status --porcelain frontend/` 只有 package.json 那一行，`find` 无游离 `*.d.ts`），
  而 `"tsc -b && vite build"` 里的 build 脚本本来就是这个形态 ⇒ 两条命令同一语义、少一处自造 flag 组合。
  **本枚 commit 只含票面 + `frontend/package.json`**（checkpoint；`internal/risk/provenance_syncdirs_other_test.go`
  在 `git status` 里是**别人（票 82）在飞的改动**，我没碰也没提）。
  **接续代理的领地重申**：`frontend/**`、`docs/contracts/**`、`internal/panel/**`（只加）、
  `cmd/wisp/panel_assets.go`（本票所有）、`.github/workflows/ci.yml`（只加自己的 job）。
  **不碰**：`internal/ball/**`、`internal/risk/**`、`internal/tools/**`、`tools/d22scan/**`、`allowlist.txt`、`design/**`。
  `next=` AC#1 收尾——给 `panel.Assets` 加一条**产物完整性自检**（entry 引用必须能在同一 embed 树里取到字节），
  配 `internal/panel/assets_test.go` 的"断引用"用例（这条就是票面点名的"能抓到产物没被嵌进去"的用例），
  再跑真构建 + `wisp.exe` 无 node 载入，报产物哈希与退出码。

- 2026-09-21 13:5x（**owner 批了 R19 的取舍：留 / 缓 / 砍 照我给的方案执行**）：接续的人**照这张表做，别再自己发明**：
  - **留（二期第一批就上，5 条）**：`Thinking Dots`（空闲/在听在想）、`Staggered Text`（流式逐字，`delay:30ms`）、
    `Animated List`（消息列表，`maxItems` 兼作渲染上限）、`Blur Highlight`（tool call 参数/代码高亮）、
    `Preloader`（首次加载："在装模型/在起引擎"）。
  - **缓（性能门之后，同屏最多 1 个，4 条互相竞争）**：`Glass Flow`、`Aura Blob`、`Neural Float`、`Fog Sphere`。
    **首选 `Fog Sphere`**（owner 给主页的定位就是"呼吸氛围"，与产品名"一缕"同频）；
    其余 3 个在"同屏≤1"这条规则下**届时判死，不无限期留着**。
    上的前提（**三条都要有用例**）：同屏≤1、**面板隐藏即销毁**（不是暂停渲染循环）、隐藏后 CPU 回落可测。
  - **砍（2 条）**：`Glass Cursor`（桌面常驻工具里是纯噪音，且不表达任何状态）；
    `Agentic Ball` 作为**桌面那颗球的替代品**判死（D32：休眠 CPU ≤0.5% / RSS ≤25MB，WebGL 一定爆；
    原生 Direct2D 球一字不动）。它**若**二期还想当"面板内的次要状态指示"，那是**另案新票**，不在本票范围。
  - ⚠ **一处数字要拧准（我报给 owner 时说成了"12 条"）**：owner 给的表实际是 **11 条组件**，
    所以正确分配是 **留 5 / 缓 4 / 砍 2**（合计 11），**不是 5/4/3**。
    ⇒ 别为了凑够"砍 3"去砍掉一个本没在清单上的东西，也别凭空发明第 12 条。
  - **一期不变**：以上全是 **React Bits = 二期**；一期 **零 React Bits 代码进树**，基座 = Beautiful UI + shadcn + Tailwind。

- **14:5x `agent-ticket77d`（第三任）第一枚 checkpoint：审计收下前任未提交改动 + AC#1 取到真数 + AC#4 台账与真红**。
  **对前任（`agent-ticket77b`）未提交改动的审计结论 = 收，且一字不改**：
  `internal/panel/assets.go`(+52)/`cmd/wisp/panel_assets.go`(+15)/新增 `internal/panel/assets_test.go`(182 行)
  是 **AC#1 的自洽半成品**（不是残件）：`Assets.Check()` 把"entry 在 embed 树里、它点名的哈希资产不在"这个
  真实失效模式（Vite 改哈希 / `all:dist` 模式打错 / 构建中断）变成一条有退出码的判定，`newAssets()` 抽出来
  让同一判定能被**合成 bundle** 驱动，因此"无 node 的机器"上也有证据。它缺的只是**变异检验**与**一次真取数**，
  两者在本条补齐；判据不是"看着能用"：
  (1) **变异真红**：把 `if _, _, err := a.Resolve(ref); err != nil` 改成 `... && false` ⇒
      `go test ./internal/panel/ -run 'TestCheck|TestBuiltin|TestAnchor'` **`--- FAIL: TestCheckRejectsEntryNamingAnUnembeddedAsset`**、包 rc=1；
      还原后 `grep -c "err != nil && false"` = **0**、包重新 `ok`。
      ⚠ 诚实口径：同一变异下 `TestBuiltinBundleIsCompleteOrAbsentNeverHalf` **仍 PASS** —— 因为真 bundle 本来就完整，
      它的牙只在合成 bundle 那条上；这条不是被变异证明的，是被下面 AC#1 的真数证明的。
  (2) **四种假绿逐条查**：`go test -count=2 ./internal/panel/ -v` = **RUN 32 / PASS 32 / FAIL 0 / SKIP 0**，
      32 == 2 × 16 个不同名（对得上，无同名刷分）；无 `--- SKIP`；`-run` 那条的匹配数=5>0；编译 rc=0 非"编译失败当变异"。
  (3) **没有为变绿删断言**：`TestAnchorOnlyBundleIsNotBuiltAndFailsClosed` 钉"只有锚点 ⇒ built=false 且
      Resolve/Check/Manifest 三入口全 fail-closed"，`TestResolveRefusesPathsOutsideTheBundle` 钉 5 条越界路径。
  **AC#1（PARTIAL 收口到"能命令级复现"，框仍不勾，原因见末）**：
  `go build -o build/wisp.exe ./cmd/wisp` **rc=0**（27878718 B）。**剥光 node/npm 的真跑**
  （`PATH=/c/Windows/System32:/c/Windows:/usr/bin`，`command -v node`/`npm` 双双 **ABSENT**）：
  `wisp.exe version` **rc=0**、`wisp.exe panel-assets` **rc=0** = `panel assets embedded: 4 files, entry=index.html built=true`、
  **`wisp.exe panel-assets -check` rc=0** = `entry=index.html built=true 2 asset refs resolve
  [./assets/index-DIrih6dy.js ./assets/index-BQ3J3QdO.css]` —— 即"页面点名的资产二进制自己就有"这条现在是命令级事实。
  `node_modules/` 被 `git check-ignore -v` 证实由
  `frontend/.gitignore:7` 挡住（根 `.gitignore` 的 `node_modules/` 一行是第二层），未进本枚 commit。
  **AC#4（本票范围内的半边闭合）**：`sh scripts/d22scan.sh` 自报**逐作用域文件数**——
  `scope ban #6 frontend/ examined 38 text files`（票 88 已把 ban #6 翻成 `live:true`，与编排者给的 35 相比
  我这次多 3 个文件，因为本票之后 `frontend/` 又落了文件；**取我自己的数 38**）、
  `scope ban #8 design/ 16 text files` / `ban #8 internal/ 322 Go files` / `ban #8 cmd/ 25 Go files`。
  **planted 真红**：临时 `frontend/src/__ban6_probe.ts` 写入 `approval.decide` ⇒
  d22scan **rc=1**（不再是前任遇到的那个 rc=2 guard 抢跑）并点名 `frontend/src/__ban6_probe.ts:1: [panel-approval] ...`，
  探针已删除（`ls` 证不存在）。⇒ AC#4 的两个判据中"逐作用域文件数>0"与"plant⇒rc=1 且点名"都在此闭合，
  **但 ban #8 的声明作用域不含 `frontend/`**（scanner 只报 design/internal/cmd），
  前端零 emoji 这一半靠的是本票自装的 `TestFrontendHasNoEmoji`（24 文件 0 命中）——
  把 ban #8 扩到 `frontend/` 要在 `tools/d22scan/**` 改，那是编排者的地界，本票**未改任何 ban**（台账如实记此口径差）。
  ⚠ **一条先于本票存在的红，不追**：当前工作树 `sh scripts/d22scan.sh` **rc=1**，但**唯一命中的不是我**——
  `internal/winsec/winsec.go:126: [pathresolver-bypass]`（票 89 在飞的未提交改动，我没碰）。
  另 `go test ./cmd/wisp/` 在本机**加载期就 rc=1（`0xc0000135`）**，纯净树 `git archive 63ef895` 同读数 ⇒ 仪器问题，登记不追。
  `next=` AC#3（L2 卡的**渲染证据**：由 `internal/risk` 真决策对象驱动、并给出 DOM/字符串级"这些值真被画出来了"的证明）
  → 之后 AC#6（CI job 的真实 run id，本票不能 push，只能写"欠 push 后复跑"）。

- **15:1x `agent-ticket77d` 第二枚 checkpoint：AC#3 拿到渲染证据 + AC#1 拿到真半嵌红 + AC#6 job 落盘（无 run id）**。
  **AC#3（渲染证据到位，框仍不勾，欠的是 host/pump）**：新文件 `frontend/scripts/render-l2.tsx`
  用 `react-dom/server` 渲**生产根** `src/App.tsx → PanelSkeleton → L2ApprovalCard`，
  输入是 `wisp.exe panel-assets -l2 fs.delete -irreversible delete -- "C:/Users/swq/Documents/notes.md"` 的 JSON。
  为此给 `-l2` 探针补了一条 `-irreversible <ops>`：`risk.Facts.Irreversible` 是**被判断的输入**，
  CLI 不给就永远出不来 R8 ⇒ 真读数 **level=L2 rules=[R1 R8]
  reason="R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除）"**（此前只能到 L1/R1）。
  **判据全部由输入派生，脚本里没有一个被断言的字面量**（否则就等于"卡恰好画了这些字"）：
  rc=0、**1 card / 6885 B HTML**；grep 实到 `>L2<` ×1、`fs.delete` ×2、`R1` ×1、`R8` ×1、
  原因行 ×1、`cli &gt; panel-assets` ×1、`判定来自 原生风险评估` ×1、卡标题 `需要你的确认` ×1。
  **两条负向对照**：(a) 上游 demo 的问卷串（`How many flavors should we launch?` / `Chocolate chips` /
  `Send answers`）在 HTML 里 **0 命中**；(b) 卡数必须等于 pending 数（防止"少画一张卡"读成绿）。
  **Q-23 fail-closed 的渲染级证据**：把同一 JSON 改成 `reasonKnown:false` + `reason:""` 再渲
  ⇒ rc=0、`信息不足` ×1、原原因文本 ×0（真没了）、`>L2<` 仍 ×1（**没有降级成看起来安全的空位**）。
  **变异检验（证明这整套断言有牙）**：把渲染目标临时换成 vendored demo `approval-card.tsx` ⇒
  **rc=1 / 10 条点名失败**（8 条"渲不出真值"+1 条"0 cards for 1 pending"+1 条"带着 demo 串"）。
  还原后 `sha256sum` 与快照同串、`diff` 0 行、`npm run render:l2` 重新 rc=0。
  **同一份真 JSON 已固化为 `frontend/fixtures/l2-card-fs-delete.json`**（与 CLI 输出 `diff` 为空），
  CI 无 Go 也能重放。**为什么放 `scripts/` 不放 `src/`**：它 `import node:fs` 读文件，
  放进 `frontend/src` 会被 AC#5 的 `TestPanelFrontendIsStateless`（扫 7 类持久化 API 含 `node:fs`）判红——
  而给那条用例开例外就是削弱 AC#5 ⇒ 代价如实记在这里：**该文件不被 `tsc -b` 覆盖**（`vite build --ssr` 会做语法转换、`oxlint` 会扫它），
  它渲染的 `src/**` 仍在 typecheck 范围内。**还欠的两件（不在本票地界）**：① 生产装配根今天是
  `cmd/wisp/panel_assets.go` 的 `-l2` 诊断分支，真事件流是票 35 的 bridge 推送 + 票 33 的 WebView2 host；
  ② 真浏览器里的 DOM（本证据是 `renderToStaticMarkup` 的 HTML 串，React 自己画的，但没有 layout/paint）。
  **AC#1 补上"真半嵌红"（不再是合成 bundle）**：`mv frontend/dist/assets/index-UL9kYvYl.js` 出 embed 目录
  → `go build` 得到 3 文件的二进制（entry 在、`Built()=true`）→ **`wisp.exe panel-assets -check` rc=1**
  并点名 `index.html references "./assets/index-UL9kYvYl.js" but the embedded bundle cannot serve it`；
  `npm run build` 重生成 + 重编 → `-check` **rc=0** `2 asset refs resolve [./assets/index-UL9kYvYl.js ./assets/index-CEH-Pz8P.css]`。
  ⚠ **一条诚实的边界**：`Check()` 证的是"embed 内部自洽"，**不是"embed 比磁盘上的 dist 旧"**——
  旧 exe 带着自己那套自洽的旧哈希，`-check` 仍 rc=0（实测）。产物↔二进制同一性由"build 顺序 + CI 同 job"保证，不由它保证。
  无 node 复跑（`node`/`npm` 双 ABSENT，exe 27881312 B）：`panel-assets` rc=0（4 files, built=true）、
  `-render /assets/index-UL9kYvYl.js` → **264929 B / text/javascript**、`-render ../go.mod` **rc=1 拒绝**。
  **AC#4 数字更新**（本枚又往 `frontend/` 落了 2 个文件）：`scope ban #6 frontend/ examined 40 text files`
  （14:5x 那枚是 38，两条新文件各 +1）、ban #8 仍只报 `design/ 16` / `internal/ 322` / `cmd/ 25`，
  **planted 红复现**：`frontend/src/__ban6_probe.ts:1: [panel-approval]` + **rc=1**，探针已删。
  ⚠ **要编排者裁的一条口径**：AC#4 字面要 ban #8 对 `frontend/` 的自报文件数，但它的声明作用域不含 `frontend/`，
  而 `tools/d22scan/**` 不许本票动；What-to-build 第 5 条给的两种合法选择里本票走的是"自装 + ASCII 化"
  （`TestFrontendHasNoEmoji` 现扫 **25 文件 0 命中**、ban #6 孪生检查 **36 文件**）。
  ⇒ 我按字面记 **PARTIAL**，不自己勾；若裁"自装即满足"，请由改 ban #8 作用域的人同批勾。
  **AC#6（job 落盘，run id 明确没有）**：`.github/workflows/ci.yml` 追加 `lint-frontend`（ubuntu-latest，
  `working-directory: frontend`，**8 步 = checkout + setup-node(24,npm cache,lockfile 键) + 6 条命令**：
  `npm ci` / `typecheck` / `lint` / `tokens:check` / `build` / `render:l2 -- fixtures/...`；
  无 `if:`、无 `continue-on-error`、无 path 过滤（D22 mode-6），且**它是独立 job**，别的 job 红遮不住它
  （A44① 那个"步骤在永久红的步骤后面所以从未产出结论"的病，正是靠 job 级隔离而不是靠排序解决的——本票全步只用 node，票 78 的 Go-on-Linux 问题够不到它）。
  **本地逐步 rc**：typecheck 0 / lint 0 / tokens:check 0 / build 0 / render:l2 0；
  ⚠ **`npm ci` 本地故意没跑**——它会先删 `node_modules` 再联网重装，失败会把这台机器上唯一能出这些数字的环境拆掉。
  ⇒ **本票不能 push，所以 AC#6 没有任何 run id**；`gh run list --branch dev` 最新 5 条
  （`35566711794` / `35566346598` / `35566028944` / `35564183459` / `35564090950`）**全部 completed/failure**，
  且都早于本枚 commit ⇒ 里面**不可能有 `lint-frontend`**。**欠：编排者 push 后 `gh run view <id> --json jobs` 取步级结论**，
  我不用"本地跑过了"替代这一格。
  **门禁复跑（本枚碰的东西）**：`gofmt -l cmd/wisp/panel_assets.go internal/panel` 空、`go vet ./internal/panel/` rc=0、
  `go test -count=2 ./internal/panel/` **RUN 32 / PASS 32 / FAIL 0 / SKIP 0**（16 个不同名 × 2）、
  `go build ./cmd/wisp` rc=0、前端 `typecheck`/`lint`(**6 warnings 0 errors**，逐名：`scripts/vendor.mjs:15,53`、
  `src/components/ai-native/tool-chips.tsx:89`、`src/components/panel-skeleton.tsx:16`、
  `src/components/ui/badge.tsx:61`、`src/components/ui/button.tsx:77` —— **6 条全部既有条目，
  新文件 `scripts/render-l2.tsx` 0 warning**；`vendor.mjs` 被扫到也顺带证明 oxlint 覆盖 `scripts/`)/`build` rc=0。
  ⚠ **不追的红**：`sh scripts/d22scan.sh` 现在 rc=1，唯一命中是 **HEAD 上** `internal/winsec/winsec.go:126 [pathresolver-bypass]`
  （票 89 的 `57bdbb2` 引入，编排者已为它立案 **票 94/A64**，`git show HEAD:internal/winsec/winsec.go` 可复核，非本票）；
  `go test ./cmd/wisp/` 本机加载期 rc=1（`0xc0000135`）照旧登记不追。
  `next=` 本票范围内**已无可推的框**：AC#1/AC#3 等票 33/35 的 host+pump，AC#4 等 ban #8 作用域那一行，AC#6 等 push。
  若编排者要我在这些之前再加一层证据，最有价值的是**真浏览器里的 DOM 断言**（需要新测试工具链，一期未定，等裁）。
