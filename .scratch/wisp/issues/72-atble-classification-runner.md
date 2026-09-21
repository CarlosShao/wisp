# 72 — A 表（受保护路径）在 CI runner 上退化成 B 表：C26 的纵深防御真破了一格

**Status:** in progress（R17 已解除 D22 闸门 ⇒ 直接改实现）
**Claimed by:** implementer（票 72 代理，2026-09-21 10:05 接手；根因已由票 70 的 run 35547905707 诊断输出提供）
**Last update:** 2026-09-21 10:05（implementer 认领；R17 见票尾）
**Blocked by:** —（与票 70 共享同一条测试，但**修的是实现不是 CI**；票 70 只负责"这条 CI 变绿"）
**Parallel slots:** ≤1 sub-agent（碰 `internal/risk/`，那里是**冻结契约**，见"硬约束"）
**Spec refs:** C26 PathResolver、SPEC-06 §4 解析管线、D22（契约变更需人批）、票 18 的红队四连

## 症状（不是占位，不是断言写松）

票面/A27 原本的假设是"`test-windows` 那个步骤名自带 placeholder ⇒ 可能是设计上先红的占位"。
**这个假设被 `a04d3e2` 实测否证**：那条命令跑的是 `internal/risk/pathresolver_junction_windows_test.go:88`
的**票 18 真红队用例**（同文件 Case 1..10 覆盖 junction / 8.3 / UNC / `\?\` / 大小写混排 /
豁免只按精确路径 / A 表不可覆盖 / B 表默认拒 / 幂等 / 展开），**一条都不是占位**；
步骤名只是 18/20 落地后没改的**陈旧命名**。

真实失败（GitHub Actions run **35547905707**，commit `5c8f9e4`，`windows-latest`）：

```
pathresolver_junction_windows_test.go:104: canonical "C:\Users\runneradmin\AppData\Local\Temp\
  TestPathResolverJunctionWindows4168261687\001\.ssh\id_testkey" classified B, want ClassA (defense in depth)
```

**同一条命令本机 windows 跑是 PASS**（`:88` 的 reparse 门与豁免门都过，只有**分类**不同结果）。

## 为什么它是安全洞而不是"环境差异导致的假红"

A 表（`blacklist.go:132` 的 `~/.ssh/**` 之类）语义是**不可放行的禁区**；B 表是
"确认可放行"。分类从 A 掉到 B ⇒ **一次 L2 确认就能读到本该永不可达的路径**。
所以判据不是"让 CI 绿"，是"**A 表在任何 home 落点下都必须赢**"。

## 已经被排除的（不要重新排除一遍）

死前那位代理逐条查过、并写进 `a04d3e2` 的 message：
- `os.UserHomeDir()` 在 windows 就是 `os.Getenv("USERPROFILE")`（读的是 `$(go env GOROOT)/src/os/file.go:605`，
  **不是** known-folder），`t.Setenv` 一定生效；
- `normPath` 两侧都小写、`isUnder` 是 `dir+\` 前缀匹配；
- A 表那条 `~/.ssh/**` 在 `home == temp 根` 时**必然**命中；B 表命中的是 `id_*`（`:155`）。

⇒ 症状**精确等价于**："Resolve 交回的 canonical 与 USERPROFILE 锚点**不在同一棵子树下**"。
它留下的未决问题（也就是它的死因）：**canonicalizer 是 Windows-by-design 会产出这种锚点，还是这里有一个真 bug**——
最可疑的两条路径是 ①runner 上 temp/home 走了 **8.3 短名或大小写混排**（`RUNNER~1`）而锚点没有，
②canonical 化把 `\\?\` 前缀或卷大小写规范化成与 `USERPROFILE` 不同的形状。

## 诊断增量已经有了，别重做

`a04d3e2` **一字未改断言**，只把失败信息改成能一次跑定位的：多打三行 `USERPROFILE` / `HOME` /
`userHomeDir()` 实值、A 表实际去比的前缀锚点、以及 `isUnder(...)` 的布尔结果。
本机复跑仍 PASS、gofumpt 空。⇒ **下一次真 run 会直接把 runner 上的锚点打印出来**，先读那段再动手。

## Acceptance criteria

- [ ] **AC#1 根因定位**：读新 run 的诊断输出（**不要**在本地猜完就改码），给出"A 表锚点为何不比 B 表更近"的
      **机制级**解释，并指明具体那一行代码。允许结论是"Windows-by-design 需要归一化两侧"，但必须给出证据。
- [ ] **AC#2 修复方向：归一化两侧，不是加特例**。修法必须让"A 表优先"成为**不变式**
      （同一条路径既命中 A 又命中 B 时永远判 A），**禁止**用"把 `.ssh` 的特例写进 canonicalizer"或
      "调整 case 顺序"来过关。
- [ ] **AC#3 双向变异检验**：①把修复退回旧实现 ⇒ 分类用例必须转红；②**构造一条 A 与 B 同时命中的形状**
      （比如把禁目录做成 temp 下的子目录、文件名同时匹配 B 表模式），证明修复后的不变式判 A。
- [ ] **AC#4 本机 + runner 两侧都过**：贴出两侧**逐字同命令**的输出与真实 exit code；
      runner 侧必须来自一次**真 run**，不是本地等价环境。
- [ ] **AC#5 不放宽任何断言**：`ClassA` 的期望值、A/B 表内容、豁免语义（只按精确路径）
      一律不许改；若证据表明**必须**改契约，**停下来上报编排者**（D22）。
- [ ] **AC#6 顺带把票 70 的 AC#3 闭掉**：本票的修复就是那条 CI 步骤转绿的因；
      修好后在票 70 的 log 里补一行指回本票，并把"步骤名陈旧"这件事改成非误导性命名（改名不算放宽门）。

## 硬约束

- `internal/risk/assessor.go`、`internal/risk/pathresolver*.go`（**除测试外**）、`rules_gateway.go`、
  `docs/PLAN.md`、`docs/specs/**` 都是**冻结契约**。本票大概率要动 `pathresolver*.go` 的实现体——
  **动之前先把改法报给编排者判**，不要先改了再说。测试文件（`*_test.go`）可以直接动，但**不许改断言来迁就实现**。
- 共树共索引：`git add` 只加显式路径；commit 前必须 `git diff --cached --name-only` 确认只有你的文件；
  **禁止 `git commit --amend` / `git reset` / `git rebase` / `git stash` / `git checkout .` / `git add -A`**
  （registry A31/A34：这两条今天各真出过一次事故）。**绝不 push**。
- commit message 用 `git commit -F -` + **带引号的 heredoc**（`<<'EOF'`）；中文长段落不要用 `printf`。
- **每个 commit 同步票面**（勾框 + 一行 log + `next=`）：票 70 那位做了 5 个 commit、174 次调用，
  票面 Progress log 里**一条都没写**——工作全在、账目全空，下一个代理会把这些活重做一遍。
- 别碰 `cmd/wisp/providers*`、`internal/tools/bridge_*_test.go`、`internal/ball/tokens*`、`tools/d22scan/**`
  （四个在途代理的包）。

## log
- [2026-09-21T10:10Z] agent=agent-ticket72 did=**认领 + 根因定位（AC#1 的证据链已闭合，先记 checkpoint）**。
  起点证据是票 70 从 run **35549859581**（commit `6c1b5e9`，job `test-windows`）里抠出的 `a04d3e2` 诊断实值：
  `USERPROFILE="C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestPathResolverJunctionWindows2581009681\\001"`
  / `A-tier anchor tried="c:\\users\\runner~1\\...\\001\\.ssh" isUnder=false`
  / `canonical "C:\\Users\\runneradmin\\...\\.ssh\\id_testkey" classified B`。
  **机制级结论**：`Resolve` 的产物两侧不对称——候选路径过了 `resolveHandle`（`GetFinalPathNameByHandle`，已展开 8.3），
  而 A 锚点走 `internal/risk/blacklist.go:120` 的 `normDir(userHomeDir())`，也就是 `os.Getenv("USERPROFILE")` 的
  **字面拼写**，从没进过 C26 管线 ⇒ `blacklist.go:132` 的 `isUnder(p, home+`\.ssh`)` 拿长名去比短名，静默 miss，
  掉到 B 表 `id_*`（`blacklist.go:155`）。修法按 R17 不变式做：**两侧都进句柄真实路径**，锚点无法证明已展开时 fail-closed。
  next=落 `pathresolver*.go` 的锚点展开实现 + `blacklist.go` 改为多形比较（不动表内容），然后跑 AC#3 双向变异。
- 2026-09-21 09:10 编排者建票。**为什么单开一张**：票 70 的代理在 `a04d3e2` 里明确写了
  "这条按真缺陷登记、够格单开一张票"，而它自己撞 turn 上限死了；把一个安全分类失效留在
  "让 CI 变绿"那张票里，很容易被下一个代理用"改断言/加 skip"的最短路径解决掉——那正是本项目最贵的一类错。

## 裁定 R17（2026-09-21 09:52，编排者）：**本票不需要 D22 批准**——契约早就写死了

我建票时在票面写了"⚠ 本票要动 `internal/risk/`＝冻结区 ⇒ 改法先报编排者判"。
**那句话我现在自己更正**，理由是我去读了契约原文而不是继续凭印象设闸：

- `docs/specs/SPEC-06-security-gatekeeping.md:52` 把解析管线明写成
  `→ 展开 8.3 短名 → 规范化 UNC`；
- `docs/PLAN.md:1376`（C26 本体）要求"**取句柄真实路径（Win: `GetFinalPathNameByHandle`）**"；
- `docs/PLAN.md:1796` 更是**直接推翻**"字面比较"的黑名单，原话理由就是
  "**junction/8.3/UNC/`\?\` 可全部绕过**"。

⇒ **8.3 短名必须在解析层被展开，这是既有契约要求，不是新政策。** 所以：
runner 上 `RUNNER~1` 与 `runneradmin` 两种拼写导致 A 表不命中，是**实现缺陷（bug）**，
修它**不触发 D22**。**D22 仍然管住的是**：若有人想把契约文本改成"分类容忍拼写差异"，那才是契约变更，须人批。

**不变式（本票真正的验收对象，写死）**：
> **安全分类的结果不得依赖路径的拼写形式。**
> 任何"A 锚点与候选路径比较"的代码，两侧都必须是**已展开的句柄真实路径**；
> 若某一侧无法证明已展开，**只能 fail-closed（判 A / 拒绝），不允许 fail-open**。

**三种"看起来能变绿"的解法，本票明令禁止**（这才是我设闸的原意）：
1. 把测试里的锚点也换成短名、或比较前对两侧做大小写/短名"对齐"的**局部补丁**——那是**在测试里复刻 bug**；
2. 给 runner 加特例（`if isCI` / 按环境变量放宽）——分类逻辑的输入里出现**调用方可控的选择器**，
   与 registry **A19/M-7/C-3 同一族**（见编排者记忆第 6 条）；
3. 用 `t.Skip` 或 `//go:build` 把它挡出 `test-windows`——**票 70 在 `secret` 上用 build tag 是合法的
   （DPAPI 按 C28 真是 Windows-only），在本票不合法**：路径拼写不是平台 API 限制，是我们的 bug（A40③）。

**优先级说明**：本票高于普通票，因为它让 **A 表（不可放行的敏感路径）在一种真实环境下静默降级成 B 表（一次确认可放行）**。
本机不复现不代表不存在：`GetShortPathNameW` 在启用 8.3 的卷上同样会产短名，
票 20 已经证明"**桥会拒 8.3**"，但那是**工具层第二次解析**在挡（A38②），**判定层本身仍会被拼写骗过**。
