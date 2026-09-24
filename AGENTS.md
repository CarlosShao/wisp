# AGENTS.md — Agent 执行版薄索引

> **这份文件是什么**：`docs/specs/SPEC-12-roadmap-governance.md` §6「文档交付物清单」里要求的
> Agent 执行版第一项 —— 根目录薄 `AGENTS.md`（D1–D47 一行摘要 + 禁止清单 + 未定义即停 + docs 索引）。
> **这份文件不是什么**：它**不是权威来源，不新造任何规矩**。下面每一条都逐字抄或具名指回下列文件之一：
> `docs/PLAN.md` · `docs/specs/SPEC-*.md` · `.scratch/wisp/issues/README.md`。
> **任何一处与那些文件不一致 ⇒ 以那些文件为准，并把这一处当作本文件的缺陷上报，不要照它做。**
> 生成时刻 `2026-09-24 09:4x +08`，锚点 `4e66817`；D 表由脚本从 `PLAN.md` 标题现抽（47 枚、无缺号）。

---

## 0. 进场先读这三句

1. **未定义即停**（D22 闸门③）：方案没覆盖的情况 → 停下来问，**不得自行假设**。见 §2。
2. **改契约＝人工批准**：改 `C1–C32` 或 `D1–D47` = 人工批准（`SPEC-12 §4.1` 逐字）。
   agent 单方面改契约 = 跑歪模式 #1，**对抗验收直接判失败**。
3. **裁决者≠实现者**：每片完成的缺口审计与对抗验收必须由**另一个** agent 做
   （`SPEC-12 §4.3` #1/#3、`issues/README` 「Hard global constraints」末条、D22 双角色）。

---

## 1. 禁止清单（逐条抄自 `.scratch/wisp/issues/README.md`「Hard global constraints」与 `SPEC-12 §4.1`）

**1.1 不许改的东西**
- `D1–D47` / `C1–C32` / `R1–R9` / **D43 转移表** —— 契约变更须人工批准。
  ⚠ 按**范围收紧**理解，不是按票号：漏计会让人以为"D47 不存在、可以随便动"。
- **SLO 阈值 / golden / `thresholds.go` 一字节都不许动**；不许为了变绿放宽任何断言。
- `DEFERRED(D-xx)` 代码标记必须与 `SPEC-12 §5` 登记表 **1:1 双向**对得上。

**1.2 禁止的代码形状**（CI 自动扫，`tools/d22scan`）
- 裸 `go func(` 而无 owner/recover；
- 在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策（直接用即判违规，`PLAN.md:1288`）；
- 明文 API 密钥；
- 用"墙钟时间差"实现超时；
- 从镜像站取哈希；
- 由面板侧来源的 L2「允许」；
- 把宿主内部 artifacts 写入实现成受门控的 Tool；
- UI 代码里出现任何 emoji（`U+2190–U+2BFF`、`U+1F300–U+1FAFF`、`U+FE0F`）。

**1.3 测试的注入面**
- 只在接缝注入：`C8 AudioSource`（wav）· `C5 LlmProvider`（golden SSE）· `C17 PanelBridge` · CLI `wisp run`。
- **不许"用 mock 代替真的"来假报完成**。

**1.4 Git 纪律**（权威文本在 `.scratch/wisp/issues/README.md` 规则 1/2 与 `docs/reports/HANDOVER.md`）
- 子代理**只 commit、不 push**；推送由编排者在核过之后做。
- `commit` 必须带**显式 pathspec**；**禁 `git add -A` / `git add .`**。
- **禁 `--amend` / `reset` / `rebase` / `stash` / `checkout .` / `clean`**（共享工作树里会吞掉别人的活）。
- 已推送的历史不改写，要更正就**追加新 commit**。
- 临时件**只建不删**（`issues/README` 规则 8）；**绝不在仓库目录内建 worktree 或 checkout**。

**1.5 交付物落点**
- 工单：`.scratch/wisp/issues/NN-slug.md`；`-done` 后缀是**防重领的唯一键**。
- 裁决表：`docs/evidence/s1/`。真相源台账：`docs/reports/pending-and-issues.md`（`A##`/`R##`/`Q##`，**只追加不删**）。
- 停车点：`docs/reports/HANDOVER.md` §4 **最新时间戳**那一节。

---

## 2. 未定义即停：已知待定案项（逐字抄自 `SPEC-12 §4.2`）

碰到下面任何一条，**停手上报**，不要按自己的判断填：

- `AwaitingApproval` 遇系统挂起（倾向作废判拒绝，S7 切片卡定案）
- `web.search` 实现路径（抓结果页 vs API，S3 前）
- Go module 组织名 + P10 命名核查（S1 建仓前，**阻塞**）
- 便携模式 × DPAPI（`SPEC-02 §6` 建议，S1 定案）
- `C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿（S7/S5 切片卡批准）
- WebView2 冷拉起 >2s → 重评 L2 卡是否回原生（P11，S0）

---

## 3. 每片完成后必须做的五件事（逐字抄自 `SPEC-12 §4.3`）

1. 重跑**缺口审计**（审计者 ≠ 实现者，D22 双角色）
2. 失败预演（`SPEC-10 §7`）
3. 对抗验收报告（`SPEC-10 §8` 的七项检查）
4. 场景↔切片双向核对（`SPEC-00 §4` 矩阵，`SPEC-12 §3`）
5. 推迟项登记更新（有新增推迟必须**五字段齐全**，`SPEC-12 §5` + M0 要求①）

---

## 4. docs 索引

| 路径 | 是什么 |
|---|---|
| `docs/PLAN.md` | 实施级方案定稿（47 条决策 + 32 条冻结契约，四轮打磨 + 用户全部拍板） |
| `docs/specs/SPEC-00 … SPEC-12` | 13 份切片规格：`00` 产品总览与场景矩阵 · `01` 架构 · `02` 数据存储 · `03` 配置与密钥 · `04` 语音链路 · `05` Agent 内核 · `06` 安全门控 · `07` 工具与插件 · `08` 球与面板 · `09` Windows 平台 · `10` 测试与验收 · `11` 构建部署 · `12` 路线图与治理 |
| `docs/specs/README.md` | spec 目录的优先级声明（**PLAN.md 定案内容最高，spec 不得与之矛盾**） |
| `.scratch/wisp/issues/` | 工单池（含 `README.md` 的派单规矩与「Hard global constraints」） |
| `docs/evidence/s1/` | 每格的裁决表（**必须出自非实现者**） |
| `docs/reports/pending-and-issues.md` | 真相源台账：`A##` 账 · `R##` 退回原因 · `Q##` 待人拍板 |
| `docs/reports/HANDOVER.md` | 停车点，§4 最新时间戳那一节＝从哪儿接 |
| `docs/reports/injection-timeline.md` | 伪授权形状的世代记录（**判"是不是授权"看这份的判据**） |

⚠ `SPEC-12 §6` 还列了 `DECISIONS.md` / `DEFERRED.md` / `RISKS.md` / `PRECHECK.md` / `ARCH.md` /
`SEQUENCES.md` / `docs/contracts/C1..C32.md` / `docs/STATE_MACHINE.md` / `docs/TOOLS.md` /
`docs/slices/S0..S8.md` / `docs/SLO.md` / `docs/BUILD.md` 这些交付物 ——
**截至锚点 `4e66817`，它们在仓里不存在**（含 `docs/` 目录本身没有 `DECISIONS.md`）。
其内容目前仍在 `PLAN.md` 正文与 `docs/specs/**` 里。**不要按上面那张索引去找这些文件然后以为丢了。**

---

## 5. D1–D47 一行摘要

> 每行**逐字抄自 `PLAN.md` 的决策标题**（脚本抽取，非本文作者改写）。要看内容请点最后一列的行号。
> ⚠ 抽取时已剥掉加粗标记，所以**标题本身**是原文、周围格式不是。已知的取代关系有几枚直接写在标题里
> （`D32` 取代 `D18` 的两张表、`D34` 取代 `D14` 的六件套），其余修订关系一律以
> `PLAN.md §16` 的差异总览表（`PLAN.md:3162` 起）为准。**本表只负责让你知道那一枚存在、不要当它不存在。**

| 号 | 一行摘要（原文标题） | 出处 |
|---|---|---|
| D1 | 「轻量级」的定义边界 | `PLAN.md:64` |
| D2 | 深度睡眠时的常驻进程拓扑 | `PLAN.md:73` |
| D3 | 插件运行时 | `PLAN.md:108` |
| D4 | 权限与确认门控 | `PLAN.md:122` |
| D5 | 语音链路技术基座 | `PLAN.md:139` |
| D6 | 目标用户与配置界面形态 | `PLAN.md:172` |
| D7 | 平台范围 | `PLAN.md:184` |
| D8 | LLM 后端抽象 | `PLAN.md:205` |
| D9 | 非语音输入通道 | `PLAN.md:226` |
| D10 | 结果呈现层 | `PLAN.md:242` |
| D11 | 意图分类 | `PLAN.md:274` |
| D12 | 主动触发与定时任务 | `PLAN.md:306` |
| D13 | MCP 取舍 | `PLAN.md:320` |
| D14 | 内置工具集与能力边界 | `PLAN.md:348` |
| D15 | 上下文预算管理 | `PLAN.md:391` |
| D16 | 音频双工与隐私边界 | `PLAN.md:444` |
| D17 | 代码签名与分发 | `PLAN.md:473` |
| D18 | 资源 SLO 与延迟预算 | `PLAN.md:497` |
| D19 | 插件信任模型 | `PLAN.md:558` |
| D20 | 记忆系统 | `PLAN.md:610` |
| D21 | 核心参考与复用策略 | `PLAN.md:653` |
| D22 | AI Agent 开发的约束模型 | `PLAN.md:716` |
| D23 | 自用期范围 | `PLAN.md:776` |
| D24 | 实现语言 | `PLAN.md:815` |
| D25 | 语音进程拓扑（onnxruntime vs D18 空闲上限） | `PLAN.md:865` |
| D26 | 模型分发（草案零提及，自用第一天就会撞上） | `PLAN.md:901` |
| D27 | 两版交付物落位 | `PLAN.md:968` |
| D28 | 项目命名 | `PLAN.md:977` |
| D29 | 前端与视觉架构 | `PLAN.md:980` |
| D30 | 间接提示注入防护 | `PLAN.md:1057` |
| D31 | 并发与审批模型 | `PLAN.md:1113` |
| D32 | 资源 SLO 与延迟预算（整体取代 D18 的两张表） | `PLAN.md:2219` |
| D33 | 安全必修五项（实现前必须全部落地，不得推迟到「以后加固」） | `PLAN.md:2357` |
| D34 | 内置工具权威表（取代 D14 的六件套；`docs/TOOLS.md` 的唯一来源） | `PLAN.md:2514` |
| D35 | 数据模型（SQLite，`%APPDATA%\wisp\wisp.db`） | `PLAN.md:2686` |
| D36 | 配置模型（`config.toml` 全量 section 树 + 三档生效级别） | `PLAN.md:2715` |
| D37 | 错误模型（分类 + 跨边界传播 + 映射到状态） | `PLAN.md:2764` |
| D38 | 并发与线程模型（原文只有一行，而 D32 的可达性完全依赖它） | `PLAN.md:2812` |
| D39 | 契约深化（把"只有名字的契约"变成规格） | `PLAN.md:2926` |
| D40 | 关键时序视图（原文完全没有） | `PLAN.md:3010` |
| D41 | 部署与打包（原文完全没有） | `PLAN.md:2876` |
| D42 | 可靠性缺口登记（原文完全没有；每条注明落哪个模块/切片） | `PLAN.md:3029` |
| D43 | 状态机权威转移表（C12 冻结的就是这张表；`docs/STATE_MACHINE.md`） | `PLAN.md:3049` |
| D44 | 切片路线图修订 + 场景↔切片追踪矩阵 | `PLAN.md:3096` |
| D45 | 批量授权与授权梯度（解 B2；§7 RESERVED 条目须据此改写） | `PLAN.md:2137` |
| D46 | Tier-1 外部命令插件（CLI 包装）【由 §15 第 10 项定案牵出】 | `PLAN.md:2633` |
| D47 | 双语音路径：半双工只约束干活路径，陪聊路径全双工 + 可选云端 realtime 大脑 | `PLAN.md:3235` |
