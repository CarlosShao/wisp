# Wisp 外部对标 · 第五版块 B —— VS Code 扩展那一屏（Cline ＋ OpenChamber 扩展形态）

> 本腿只做一件事：**装了扩展之后 VS Code 里出现了哪些界面件**，逐块逐控件。
> ⛔ 不越界：移动端＝块 A（`block-A.md` 已交）、组件分包＝块 C；本文件一处不提它们的屏。
> ⛔ 零读零引本仓 `frontend/**` 与 `design/**`；"我方有没有对应物"这一列只引 `docs/**` 行号或本腿现读的 Go 文件行号。
> ⛔ 本腿不跑任何编译／测试／门禁。
> 三档标法：**〔已证〕**＝读到实现／**〔建了但没接〕**＝有代码无调用者／**〔仅文档〕**＝只有 README 或文档里有。

---

## §0 本轮取数的家与出处

| 家 | 仓库 | 取数方式 | commit sha | 日期 | 扩展本体在哪 | 版本 |
|---|---|---|---|---|---|---|
| **Cline** | `github.com/cline/cline` | `curl` 拉 codeload tarball（`main`）解到 `/tmp/v5b/cline-fresh`，40,530,055 字节 | `c604735fe1ce16d745f9dbf799ade79a4e3dad86`（`api.github.com` 具名） | 提交时间 `2026-09-30T02:57:24Z`（"chore(desktop): release v0.0.38"） | `apps/vscode/`（扩展 id 内部名 `claude-dev`、`displayName` ＝ `Cline`、publisher `saoudrizwan`）`apps/vscode/package.json:2-14` | **4.1.21**（`package.json:5`） |
| **OpenChamber** | `github.com/openchamber/openchamber`（10,942 star，`api.github.com` 查得） | 同一台机 `D:\work\AI\open source\openchamber` 那份 **2026-09-28 tarball 快照**为主读面（前几轮已用同一份，行号可对），另拉 `main` 新鲜 tarball 对版本 | 新鲜快照对到的 `main` HEAD＝`d78dac542d796f20536a3bb94971aa94fbb9e48e`（`2026-09-30T00:43:13Z`） | 快照落地 `2026-09-28` | `packages/vscode/`（`displayName` `OpenChamber`、publisher `fedaykindev`）`packages/vscode/package.json:2-6` | 快照内 **2.0.3**（`packages/vscode/package.json:5`）；新鲜版另记，见 §8 |

**取数代价（必须写明白）**：tarball 快照**没有 `.git`、没有 `node_modules`** ⇒
①**读不到提交历史**，本文件**每一条结论都带 `文件:行`**，一律不引 commit／PR／历史；
②**读不到依赖源码**，凡是"依赖 VS Code 宿主能力"的地方只能读调用点、读不到实现；
③Cline 那份是**monorepo**（`apps/vscode` ＋ `apps/cli` ＋ `sdk/` ＋ `docs/`），本腿只读 `apps/vscode` 这一支与其 webview 子目录 `apps/vscode/webview-ui`，**`sdk/` 与 `apps/cli` 只在需要确认"扩展形态有没有对应界面"时点名，不做 CLI 调研**（那是块 C／别的腿）。

**Cline 与 OpenChamber 的"扩展"不是同一种东西**（先定口径，否则表一会读歪）：
- Cline：扩展宿主进程＋**一枚常驻 webview 视图**，所有对话、设置、审批都画在那枚 webview 里；VS Code 原生设置页 **一项都没有**（`contributes.configuration.properties` 空，实测读数见 §2）。
- OpenChamber：扩展宿主进程里跑一台 web 服务器＋**webview 装的是它自家 web 应用**，所以"扩展形态的界面"＝web 界面＋一串只有扩展才有的宿主件（评论线程、标题栏按钮、状态栏、命令）。本腿只登记**扩展独有的那一层**，web 那一层是块 A／第四版的地界（重复即缺陷，见 §0.1）。

### §0.1 与前几轮的分工（防重）

- 第四版 `docs/reports/missing-features-2026-09-29-v4.md` 里 OpenChamber 扩展只被点了**三枚具名件**：`packages/vscode/src/quotaProviders.ts:15-52`（用量窗口）、`bridge-permission-auto-accept-runtime.ts:1-103`（每会话自动接受）、`bridge-git-*.ts` 一族（转引 `docs/reports/survey-2026-09-28-oc-mobile-vscode-extensions.md:38`）。
- `survey-2026-09-28-oc-mobile-vscode-extensions.md` §1-Q1 给过一份"IDE 多出来的"**散文式清单**（评论线程／右键菜单／会话进标签页／配额盘／bridge 32 种消息）。
- **本块的增量＝把它换成三张表**：`package.json` 的**每一枚贡献点逐枚列**（视图／命令 19 枚／键位／七组菜单／getStarted 走查）、**设置页逐条标签原文**、**状态文案逐句**，并**第一次读 Cline**（前四版六大素材家里没有 Cline，`missing-features-2026-09-29-v4.md:13-16` 逐名列的是 Step-Code／deepseek-harness／minimax-code／openchamber／pi／pi-upstream）。
- 已点过的三枚具名件本块**不重述内容**，只在需要"它长在哪一屏"时行号指回。

---

## §1 表一 · VS Code 扩展逐块逐控件表

（填写中：位置 ｜ 控件名（界面原文标签）｜ 它干什么 ｜ 我方有没有对应物 ｜ 出处 文件:行）

---

## §2 表二 · 设置页逐项表

（填写中：设置项标签原文 ｜ 改了会怎样 ｜ 我方对应物或"没有" ｜ 出处）

---

## §3 表三 · 状态与文案表

（填写中）

---

## §4 必答一：点进某个子任务/子代理，能不能看到它各自的流式工作页面？

（填写中）

---

## §5 必答二：审批卡有没有"本次／本次会话内／长期"三档＋拒绝/允许的理由输入框？

（填写中，逐家答）

---

## §6 我写错的条目（拿第四版清单逐条对抗本腿自己）

（填写中）

---

## §7 不建议抄（只给形状级理由）

（填写中）

---

## §8 没调研完／留给编排者

（填写中）
