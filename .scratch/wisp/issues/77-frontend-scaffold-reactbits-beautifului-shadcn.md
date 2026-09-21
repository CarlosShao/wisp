# 77 — Frontend scaffold: React + TS + Tailwind + shadcn, vendored react-bits / beautifului components, sharing C21 tokens

**Status:** **unblocked**（2026-09-21 13:5x，owner 第二次指令 ⇒ 裁定 **R19**）——
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
- [ ] **AC#2** token 四方对账用例落地 + 一次"只改前端一处颜色 ⇒ 用例红"的变异（grep 证落盘 → 跑 → 还原 → 证还原）。
- [ ] **AC#3** L2 确认卡渲染真数据：由 `internal/risk` 的真实决策对象驱动，**不是**组件库 demo 的假 props。
      ⚠ 能力类判据：**说清生产调用者是谁**（哪个装配根、什么事件流），否则按 A33 只能记 PARTIAL。
- [ ] **AC#4** `ban #6` 与 `ban #8` 对真实 `frontend/` 的作用域台账：贴出扫描器自报的**逐作用域文件数**（>0），
      并给出一次故意在前端写 `approval.decide` ⇒ **rc=1 且点名文件** 的真红。
- [ ] **AC#5** 前端无任何本地持久化；重启 WebView 后一切从 Go 侧恢复到同一状态（一条真机或半真机用例证明）。
- [ ] **AC#6** CI 新 job 有真实 run id + 结论（见前提 6）；**不允许用"本地跑过了"替代**。
- [ ] **AC#7** vendored 清单：`frontend/VENDORED.md` 逐文件列 来源仓库 / 组件名 / 许可（含 Commons Clause 标注），
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
