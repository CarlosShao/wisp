# 92 — 主面板输入框 v2：**档位显示 + 附件（含图片/视频粘贴）+ 工作区选择**，并且明确不做 git 切换

**Status:** open（2026-09-21 17:1x 由编排者建，来源是 owner 对 M5 的定案 ⇒ 见裁定 **R20** 与票 90 末尾"owner 定案"条）
**Type:** 能力（面板交互面）+ 安全边界（"档位/工作区"这两样都是**权限的输入**，不是装饰）
**Blocks:** nothing（它是票 90 的**落点**之一，但不是它的前提：票 90 的语义与持久化先走，本票接它的状态字段）
· **Blocked by:** 票 **90** 的 AC#1（模式作为显式参数进入决策链）——没有那个状态，面板无处显示档位。
**Packages:** `frontend/src/components/`（新建 composer 相关组件）、`internal/panel/`（消息封套与状态快照）、
              `cmd/wisp/`（bridge 的原生侧落点）。附件的**存储**沿用票 79/76 已经钉死的 artifacts 面
              （注入式命名 + 真 O_EXCL + 递归配额），**不许**在本票里新开一条"用户文件"的旁路目录。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、
              `internal/risk/pathresolver*.go`、`rules_gateway.go`、`tools/d22scan/**` 与 `allowlist.txt`（门禁是编排者的）。

## owner 的原话（2026-09-21，M5 定案）

> "不是有主面板吗，主面板的输入框中显示档位就行了啊，主面板输入框内容可以借鉴图片这种来，
> **可以添加附件，可以选择工作区**，但是什么 **git 切换就没有太大必要了，我们产品不要求 coding 能力**。"

他同时给了一张参考截图（左下"添加附件"、输入框内放模式选择器、右侧模型选择 + 发送）。
⇒ **要做**：档位显示、附件、工作区选择。**不要做**：git 分支/仓库切换那一类（他明确砍掉，别"顺手加回来"）。

## 本票最怕的失败模式（写在你动手之前）

输入框里的"档位"和"工作区"**不是外观，是权限的输入口**。仓里已经有一条冻结红线管着这个形状：

- **`PLAN.md:1588`：'允许'只接受原生侧点击**，机器门就是 `ban #6`（`frontend/` 里不得出现 `approval.decide`）。
- 如果面板能自己把档位从"每步都问"改成"全自动"，那**被攻陷的渲染进程就等于给自己签了一张免审通行证**
  —— 比批准单条操作更危险，而且完全绕开 M4（"切到全自动要 L2 确认"）。

⇒ 因此本票的形状**定死**成"显示 + 请求，授权在原生侧"：

| 元素 | 面板可以 | 必须回原生侧 |
|---|---|---|
| 档位 | **显示**当前档（来自 Go 侧状态快照）；点击可**发起切换请求** | 确认与写入：全自动那一次要 L2 强确认（M4），并写审计 |
| 附件 | 采集文件/剪贴板、给出名字与字节数、随消息发**引用** | 落盘与配额（走 artifacts 面），类型白名单不过就**响亮失败** |
| 工作区 | 显示当前工作区、发起"换工作区"请求 | 真实路径解析 + 校验（C26 PathResolver 的 reparse/ junction 拒绝在这里生效），切换写审计 |

## AC（1:1，裁决表 `docs/evidence/s1/92-*.md` 由验收方出，不是自裁）

- [ ] **AC#1** 面板**只读地**显示当前档位（三档之一），并且**面板侧不存在任何能改档位的路径**：
      给出 `ban #6` 家族的正向钉子——在快照里往 composer 组件种一个"直接改 mode"的调用，扫描/用例必须红
      （或者你论证它为什么不该被 `ban #6` 抓住、并把它挂到别的门上）。**不许**新建第二条独立通道。
- [ ] **AC#2** 附件通路：粘贴/选择**图片与视频**能把字节送到 Go 侧并出现在消息里。
      ⚠ **不支持的类型必须响亮失败并告知用户**，不许静默丢弃（票 83 的"说谎的配置键"同族）：
      给出你测过的**每一类**输入与各自的结论（png / jpg / mp4 / 一个 0 字节文件 / 一个伪装成 .png 的 .exe /
      一个 2GB 的文件 / 一个带 `\` 和 `..` 的文件名）。**字节进去了 ≠ agent 能理解** ⇒
      视频**语义理解**不在本票（挂 **Q-28**），本票只保证"不骗人 + 不越权 + 不吃掉用户的意图"。
- [ ] **AC#3** 工作区选择：一次真实切换后，**(i)** 后续操作的风险判定用的是新工作区，
      **(ii)** reparse/junction 指向外部的那次切换**被拒且给出原因**（C26），**(iii)** 切换写进审计。
- [ ] **AC#4** 反向：面板**关闭**（球唤醒但不展开面板）时，档位与已选工作区都**不丢**，
      重开面板从 Go 侧恢复到同一状态（AC#5 的既有判据形状，票 77 做过一次，照做别削弱）。
- [ ] **AC#5** 变异三向：(i) 把"面板只能显示"改成"面板能写 mode" ⇒ 必须有用例红；
      (ii) 把不支持类型的报错改成 `return nil` ⇒ AC#2 红；(iii) 把 reparse 那次的拒绝改成放行 ⇒
      **既有**安全用例红（不是本票新写的）。锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净；
      **编译失败不算变异**。
- [ ] **AC#6** 台账与门禁（只跑自己碰的范围）：`sh scripts/d22scan.sh` 纯净树 rc=0 且贴出**逐作用域文件数**
      （`ban #6 frontend/` 的 N 必须**因为本票而变大**，这就是覆盖面证明）；`ban #8` 零 emoji；
      `gofmt -l`/`gofumpt -l` 空、`go vet ./internal/panel/` rc=0、`go test -count=2 ./internal/panel/` rc=0
      且逐条点名 SKIP/FAIL。⚠ `go test ./cmd/wisp/` 在本机是**加载期 `0xc0000135` 的既有红**（票 87 已在纯净树复现），
      不要去追，如实登记即可。
- [ ] **AC#7** **负判据**：把"不做 git 切换"变成可检查的东西——在 `frontend/` 与 `internal/panel/` 里
      `grep -rn` 证明没有任何分支切换/仓库选择的能力入口（owner 明令砍掉，防止后人"顺手加回来"）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc，A31）；提交前 `git diff --cached --name-only`；
**禁** `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**；
不在仓内建 worktree/checkout（A38④，快照用 `git archive HEAD | tar -x -C /tmp/<带会话后缀>`）；
不跑整仓门禁；票面 append-only，**要改的那行先读再替换**，在标题前插段落必须把标题原文重抄进 `new_string`
并让 `git diff --numstat` 的删除列为 0；四种假绿逐条点名；数字不达标就 FAIL 附数字，不许调阈值。
`node_modules/` 不许进 commit。15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`。

## Progress log（append-only）

- 2026-09-21 17:1x（编排者）：建票。owner 把 M1–M5 全部拍完（登记为 **R20**），M5 选了"档位显示在主面板
  输入框里、加附件与工作区选择、砍掉 git 切换"。⚠ **我把这条 UI 决策的风险面先钉在这里**：档位与工作区都是
  **权限的输入**，所以本票的形状必须是"面板显示 + 发起请求、原生侧确认 + 落账"，
  判据是 `PLAN.md:1588`（'允许'只在原生侧）与 M4（全自动要 L2 确认）。
  **owner 补的两条不在本票**：视频语义理解（**Q-28**）与悬浮球/看门狗改造（**Q-29**），他都说了"不急、放最后"。
  next= 等票 90 的 AC#1 落地后派单；在此之前可以先把附件通路的 AC#2 做起来（它不依赖档位）。

## Progress log（agent-ticket92，append-only）

- 2026-09-21 18:45（agent-ticket92）：**checkpoint 1 / 基线读数**。`date` 实测 `Mon Sep 21 18:45:15 CST 2026`。
  开工前 `sh scripts/d22scan.sh` 真实读数（树 = 本工作树，含邻居 `internal/winsec/winsec_windows.go` 的 ` M`）：
  **rc=0**，逐作用域 `bans #1-5 internal/=197`、`bans #1-5 cmd/=20`、**`ban #6 frontend/=40`**、
  `ban #7 internal/tools/=17`、`ban #8 design/=16`、**`ban #8 frontend/=40`**、`ban #8 internal/=347`、`ban #8 cmd/=26`。
  ⇒ 收尾时 `ban #6/#8 frontend/` 的 40 必须**只增不减**（本票新建 composer 组件会把分母推高，这就是 AC#6 要的覆盖面证明）。
  上面那串 `frontend/src/app.js:1: [panel-approval] ...` 与 `scan_test.go:1157` 是 d22scan **自己包的 fixture 输出**
  （它跑自己的 `go test ./...`），末行对本仓的判定是 `clean - no D22 ban violations`。
  已读码定位的现状：档位由 `internal/perm.Store`（票 90）持有，`cmd/wisp/run.go:377 confirmModeSwitch` 是唯一的
  L2 强确认落点（R20/M4）；面板契约面 = `frontend/src/lib/panel.ts` ↔ `internal/panel/approval*.go`
  （`Snapshot`/`ResultChunk` 现在**只声明在 `approval_test.go` 里**，生产侧无处构造——票 35 的泵还没落地留下的形状）；
  附件存储复用面 = `internal/memory/artifacts.go`（`validArtifactName` 守卫 + `winsec.PrivateFileExclusive` 真 O_EXCL +
  `enforceArtifactsQuota` 递归配额，票 76/79 成果）；工作区校验复用面 = `risk.Resolve` + `Result.Actable`
  （票 102/107 的"改写要记账"，`res.Resolved` 在 POSIX 恒 false，所以本票不拿它当唯一确认腿）。
  计划落点：`internal/panel/composer.go`（档位/工作区**只读视图** + 请求类型）、`internal/panel/attachments.go`
  （白名单 + 魔数嗅探 + 响亮失败）、`internal/memory` 导出一条复用既有守卫的 artifact 写入腿（**不新建旁路目录**）、
  `internal/tools/paths.go` 加"工作区**收窄**"（默认空 = 行为逐字节不变）、
  `frontend/src/components/composer.tsx` + `frontend/src/lib/panel.ts`（显示 + 只发起请求）。
  **不做**：git 分支/仓库切换（AC#7 用 grep 证）、视频语义理解（Q-28）、改扫描器/豁免。
  next= 按 AC#1 先在快照里种"直接改 mode"的正向钉子并跑出红，再动实现。

- 2026-09-21 19:2x（agent-ticket92）：**checkpoint 2 / 新增用例的"修前红"原文**（AC#1/AC#4 的契约腿）。
  `go test -v -count=1 -run 'TestModeViewHasNoWriteSurface|TestComposerContract|TestPlantedComposer' ./internal/panel/` ⇒ **rc=1**：
  红名 `TestComposerContractTypesMatchFrontend`，断言原文两行：
  `composer_test.go:74: Go Snapshot emits [composer] that interface PanelSnapshot does not declare`
  `composer_test.go:71: frontend/src/lib/panel.ts declares no interface ComposerState - the composer contract was renamed on one side only`
  同一批里既有守卫 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 也被这条新字段拖红（`approval_test.go:129` 同一句 `Go Snapshot emits [composer] ...`）
  ⇒ 证明 `Snapshot ↔ PanelSnapshot` 那把双向锁**是活的**：Go 侧加一个字段而前端不声明，两枚独立用例同时红。
  同一跑里已绿的三枚（不是"修前红"，是 AC#1 的正向钉子本身）：
  `TestPlantedComposerModeWriteGoesRed`（往临时 `composer.tsx` 里种两种"页面自己改档位"的写法，
  门一 = d22scan ban #6 的同款 `approval.decide`（票 77 已把正则逐字复制进本包 `panelDecisionIdentifierRe`）、
  门二 = 本票新增 `composerModeWriteRe`（`setMode`/`panel.mode.set`/`PermissionMode =`/`perm.Store.Set`），
  日志原文 `red as required: 1 line(s) matched, first=composer.tsx:4: ... method: "approval.decide", target: "mode", to: "auto_approve"`，
  以及 `first=composer.tsx:1: export function setMode(next: string): void {`）；
  `TestModeViewHasNoWriteSurface`（反射证明 `ModeView` 零方法 + 本包不导出 `Set|Write|Apply|ChangeMode`）；
  `AC#7` 的 `TestNoGitSwitchCapabilityInThePanelSurface` 绿（`frontend/` 与 `internal/panel/` 生产文件里
  没有任何 checkout/switchBranch/repoPicker/worktree 形状的能力入口）。
  ⇒ 本票的形状按 AC#1 交了两道门而不是一道：**ban #6 抓不到 `panel.mode.set`**，
  所以"面板能写档位"这一族不能只挂在 `approval.decide` 上；正向钉子已种进 `composer_test.go`（临时 fixture 目录，
  **没有**往真 `frontend/` 里种违规文本——那会把 AC#6 的 `rc=0` 变成红的，且需要动扫描器豁免，本票禁做）。
  附件（AC#2）本轮已在 Go 侧全绿，逐类结论见下一条 commit 后的读数表。

- 2026-09-21 19:3x（**编排者插入，关于你工作树里那半份 `internal/tools/paths.go`**）：
  你在这枚文件里的 `workspace` 收窄 hunk 与 **票 107b 的放行侧修复**（`treeResolvedAsNamed`/`resolvedForm`/`rootsContain`，
  判定点 `:87`）此刻**混在同一份未提交改动里**。107b 已交件、它的会话已结束 ⇒ 这份修复现在只有工作树里这一份。
  ⚠ 三条要求：① **不要**把它当成你自己的功能改写归因；② 你 commit `internal/tools/paths.go` 时
  message 里**分署**（例：标题 `feat(92,AC#3)+fix(107b,AC#3)`，正文写明哪几行是 107b 的）；
  ③ 如果你判断自己还要改这个文件很久，就在 `next=` 里说清，我另派一次"只提 107b 那几行"的落地。
  兜底：我已把整份工作树状态存成 `docs/evidence/s1/107b-pending-paths-go.patch`（189 行，进 git 了）⇒ **丢了也长得回来**。
  ⚠ 顺带：你交件里若出现"票 102 那条在 POSIX 继续绿"之类的读数，请注意它来自 107b 的快照、不是你亲测的，标清来源。
