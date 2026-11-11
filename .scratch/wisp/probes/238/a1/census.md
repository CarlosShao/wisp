# 238-a1 · 只读普查 · AC#1「谁负责封数据根」两形代价（⛔ 零产码）

- 日期 2026-11-11 19:4x–20:1x +08 · 腿码 **238-a1** · 只读腿，唯一写点＝本文件
- 起手锚点（现跑 `git log -1`）：`e31937a90` 2026-11-11 19:46:58 +0800（分支 `dev`，`git branch --show-current` 现读）
- 路由 provider=magpie · model=auto/glm-5.3-flash
- ⛔ 本腿零 `.go` 改动、零 `go test`、零 `go build`（同机写腿在飞）；每条结论的尺都可独立复跑
- ⚠ 行号会腐烂：引用前按各节给的尺现跑

---

## §0 先摆一句立场：本腿与 owner 否决的关系

票面 Progress log 记着三笔后置（全部逐字现读于
`.scratch/wisp/issues/238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md`）：

- `:36–:43` **owner 裁定 A476（09-30）＝"本条不做"**，原话「锁个屁，人家都能登录我电脑了，还锁个蛋啊」；重开仅认 **「Q-72 重开」**；
- `:64` **A1558（11-05）**：owner-vetoed、勿再派；**132-r1 已落地** `PrivateDirAll` 创建即封；
- `:66` **wave108 腿 `238-proc`（11-09）**：STOP，零实现；其读数与本腿独立复跑的结果一致（见 §1）。

本腿是编排者 2026-11-11 派的**只读普查**：不实现、不翻框、不把封根再摆给 owner（A476 逐字禁）。产出＝两形代价表＋尺形状草案＋POSIX 半边，全部停在纸面，供编排者在「Q-72 若重开」时一次取用。**本腿不推翻 A476，也不请求重开。**

---

## §1 现量复认（票面与派单的每一条，逐条现跑）

| # | 票面/派单的断言 | 本腿现跑读数 | 尺（可复跑） |
|---|---|---|---|
| 1 | `internal/proc` 对 `internal/winsec` **直接导入＝0** | **属实（直接边＝0）**。生产码内仅注释命中：`internal/proc/envfork.go:113,116,148` | `grep -rn "internal/winsec" --include=*.go internal/proc internal/buildinfo \| grep -v _test`（本腿现跑＝envfork.go 三行注释，零 import 行） |
| 2 | `internal/buildinfo` 对 `internal/winsec` 导入＝0 | **属实**（上同尺，buildinfo 零命中） | 同上 |
| 3 | 5 枚 `PrivateDirAll` 生产点 | **票面过时**。**生产（非 `_test`）调用现有 8 枚**：`internal/observe/logging.go:99`、`internal/secret/store.go:58`、`internal/agent/spill.go:156`、`internal/memory/open.go:270,283,697`、`internal/winsec` 自身实现（`winsec.go:259` 定义＋`winsec_windows.go:313` 注释）。其中 `memory/open.go` 一枚包（3 处）在票 132 立票（09-29）后新增 | `grep -rn "PrivateDirAll" --include=*.go internal cmd \| grep -v _test` |
| 4 | `resident_windows.go:57` 装日志 sink | **过时**。现读 `cmd/wisp/resident_windows.go:173`＝`sink, sinkErr := installLogSink(rt.Layout.DataDir)`；`:57` 现在是 `hotkeyReloadSource296` 注释里的一行（票 296 文字） | `grep -n "installLogSink" cmd/wisp/resident_windows.go` |
| 5 | 票面 `:13` 所指「`observe/logging.go:72` 旧行」造宽树 | **该行已不存在**。现读 `internal/observe/logging.go:99`＝`winsec.PrivateDirAll(cfg.Dir, 0o700)`（132-r1 落地，commit `b22225513`，2026-10-23，现跑 `git log --oneline -3 -- internal/observe/logging.go` 可证）；`logging.go:16`＝`import "github.com/CarlosShao/wisp/internal/winsec"` | `git log --format='%h %ad %s' --date=short -1 b22225513` |
| 6 | `winsec` 不依赖 `proc`（甲形无环的前提） | **属实**：`go list -deps ./internal/winsec` 输出中零 `CarlosShao/wisp` 包（只到 stdlib＋`golang.org/x/sys/windows`） | `go list -deps ./internal/winsec`（⛔ 纯只读，不编译本仓新代码，同机写腿无撞） |
| 7 | proc 对 observe 的依赖现状 | **`internal/proc/boot_windows.go:19` 已 import `internal/observe`**（生产 import 行，非注释）；envfork 等另有多处。`go list -deps ./internal/proc` 输出含 `internal/observe`、`internal/secret`、**`internal/winsec`** | `go list -deps ./internal/proc \| Select-String CarlosShao` |
| 8 | `proc.Boot` 是否造目录 | **不造**。`internal/proc/boot_windows.go:93–156` 全函数自读：env 自检→`DefaultLayout`→Job→mutex→registry 检查，零 `Mkdir*`/`PrivateDir*`。`grep -n "MkdirAll\|os.Mkdir" internal/proc/*.go`（含 `_test`）＝空 | `go doc ./internal/proc`＋读 `boot_windows.go:93-156` |
| 9 | 派单问的「日志 sink 在 `resident_windows.go:57` 装」的现时序 | 现读真身：`resident_windows.go:146`（`proc.Boot`）→`:173`（`installLogSink`）→`logsink.go:157,163`（`observe.InitLog`）→`logging.go:99`（`PrivateDirAll`＝**今日首造者**）。run 腿：`run.go:246`（sink）→`:269`（`ensureFirstRunConfig`）→`:273`（`assembleRuntime`，其内 `secret.go:226` 再装一次） | 本节各 file:line 直接读 |
| 10 | 票 197 段「装配根 cmd/wisp 是唯一接缝」裁定在册 | **属实**：`docs/reports/pending-and-issues.md:8788`（「装配根 `cmd/wisp` 是唯一的接缝」）与 `:8790-8791`（裁定＝注入、不新开依赖边） | `sed -n '8786,8792p' docs/reports/pending-and-issues.md` |

**结论（现量层）**：票 238 标题所写的缺陷（首启那棵树无封可继承）**已被 132-r1 在 `logging.go:99` 收口**——今天新机器先跑常驻腿，`installLogSink`→`InitLog` 的第一次造根动作本身就是 `PrivateDirAll`（逐层先封后下钻，契约见 `winsec.go:238-250`）。这正与 `238-proc` 腿 11-09 的 STOP 读数一致（票面 `:66`）。**AC#2「谁造根」的消费者今天不存在**；以下两形摆开按派单要求照做，落在纸面上。

---

## §2 AC#1 两形摆开（本腿核心）

**共同前提（两种形都建立在这一句上）**：今日的造根者不是 `proc.Boot`（§1#8），而是**各腿装配代码里第一次把数据根/子目录交给 `PrivateDirAll` 的那一刻**。常驻腿＝`resident_windows.go:173`（sink）；run 腿＝`run.go:246`（sink，早于 `:269` 首建配置与 `:273` 装配）。所以两形改的都不是"补一道工序"，而是"把封根动作从 sink 的随行者提到所有腿之前"。

### 2.1 甲形＝常驻腿（`internal/proc`）自己调封条

**形状**：在 `proc.Boot`（`boot_windows.go:93`）内、`DefaultLayout` 定下 `rt.Layout`（`:120–127`）之后，加一步"对 `rt.Layout.DataDir` 调 `winsec.PrivateDirAll(DataDir, 0o700)`"。

**依赖边与循环判定（派单指定判据，现跑）**：

- 新增直接边：`internal/proc` → `internal/winsec`。
- **不造循环**：`go list -deps ./internal/winsec` 现跑输出**零 `CarlosShao/wisp` 包**（§1#6），即 winsec 今天不依赖 proc，反向加边无环。
- ⚠ 但这条**直接边在成文上正是禁区**：`internal/winsec/resolve.go:35–40` 逐字「the graph already runs risk → observe → secret → winsec … so the reverse edge closes a cycle」；且 `internal/risk/winsec_c26.go:9` 同形。那句针对的是"winsec→risk"反向边；proc→winsec 是**顺依赖方向的向下边**，不属于那两句禁的形状。
- ⚠ 然而**它冗余**：proc 经 observe→secret **已经传递依赖 winsec**（`go list -deps ./internal/proc` 含 `internal/winsec`，§1#7）。把传递边升级成直接边，proc 的依赖闭包一个包都不变。

**要动的文件（甲形清单）**：
1. `internal/proc/boot_windows.go`（唯一：`:127` 之后插封根步＋import 行）。
2. 无测试硬撞：`cmd/wisp` 里直接 `proc.Boot` 的生产调用仅 `resident_windows.go:146` 与 `slo_windows.go:261`（`grep -n "proc.Boot(" cmd/wisp/*.go` 现跑＝18 处全是 `_test`）；测试全走 `WithLayout` 指临时根，封临时根的父目录不碰（`PrivateDirAll` 契约 `winsec.go:246–248`："Ancestors above path are never touched"）。
3. 失败方向要新写：proc 内被封根失败怎么办——`slo` 腿的先例是 `return 2`（票面 `:27` 引 `slo_windows.go:272–278`），`resident` 腿的先例是降级照跑（`resident_windows.go:174–178`）。同一动作两种方向＝票面 `:27` 第 3 刀已判过的缺陷形状。

**代价（甲形）**：
- 新增直接依赖边 proc→winsec（虽然无环、虽然闭包不变，但它让 proc 从"进程/生命周期包"变成"生命周期＋存储安全策略包"，包的单一职责被打破）；
- **跨平台不对称内嵌进 proc**：POSIX 半边要么跟着做递归 chmod（见 §5，被 `winsec_other.go:427–432` 的存在理由顶回），要么 proc 里写 `_test`/`other` 分叉；
- proc 内要为"封根失败"再造一条失败方向裁定（两腿现行为不一致，见上）；
- **时序收益≈0**：`proc.Boot`（`:146`）只比 `installLogSink`（`:173`）早 27 行、早零次造盘动作；把封根提前到 Boot 里，改变的只是"封的动作离造根的动作近了 0 个文件系统事件"——因为 Boot 自己不造根，造根仍发生在 :173 的 `PrivateDirAll`，而那一步**本来就带封**（132-r1）。甲形在今天的代码上**无可保护的时间窗**。

### 2.2 乙形＝装配根（`cmd/wisp`）在起任何腿之前把根封一次

**形状**：`cmd/wisp` 的两个入口（`resident_windows.go` 的 `runResident`、`run.go` 的 `runTextTask`）在拿到 `rt.Layout.DataDir`/`s.dataDir` 之后、装 sink 之前，先 `winsec.PrivateDirAll(dataDir, 0o700)` 封根一层。

**依赖边与循环判定**：

- 新增直接边：`cmd/wisp` → `internal/winsec`。**今天已存在（测试文件里）**：`cmd/wisp/secret_dataroot_119b_test.go:49` 已 import winsec；生产码还没有 import 行，加一行即是首次生产直连。
- 无环：winsec 是叶子包（§1#6），cmd/wisp 在依赖图的根，向下加边天然无环。
- **与在册裁定同形**：台账 `pending-and-issues.md:8788–8791`（票 197 段）逐字「装配根 `cmd/wisp` 是唯一的接缝」＋「不新开 tools→panel 的依赖边…注入」——乙形正是"装配根自己动手"，不引入任何新包间边。

**要动的文件（乙形清单）**：
1. `cmd/wisp/resident_windows.go`：`:146`（Boot 返回）与 `:173`（installLogSink）之间插一次封根调用（一处）。
2. `cmd/wisp/run.go`：`:237`（dataDir 定下）与 `:246`（installLogSink）之间插一次（一处）。
3. `cmd/wisp/logsink.go`：把"谁先封根"写进 `:27–37` 那段 Ordering 注释（不改行为，只改文字）——否则下一位读码者会以为 sink 仍是第一个碰盘的人。
4. （可选，倾向不做）`cmd/wisp/slo_windows.go:261` 同插——诊断腿与主腿不同根场景罕见，且 `wisp slo` 的 sink 也走 `InitLog` 已封 `logs` 一层。

**代价（乙形）**：
- 生产码多两次直接 `winsec.PrivateDirAll` 调用（同一个包已提供的导出名，零新 API、零新钩子、零 `Q-72` 那种契约级新导出名）；
- 与甲形同样的"失败方向"问题（封根失败降级还是拒绝），但它可以整体复用 sink 失败的既有形态（`resident_windows.go:174–178` 的大声降级），不用新裁；
- 命令腿（`wisp secret`、`wisp models` 等）不在乙形射程内——它们今天同样由自己的首个 `PrivateDirAll`（`store.go:58` 等）带封；乙形若要全覆盖，每条腿入口都插（但这与甲形同病：`proc.Boot` 只盖两条腿，票面 `:25` 第 1 刀与 `238-c1` §1 都记过"造根的腿不止两枚"）。
- **时序收益≈0（与甲形同因）**：乙形插在 sink 之前，但 sink 的第一步 `InitLog→PrivateDirAll` 本来就"封着造"；乙形提前封的只是"数据根目录本身"这一层，而这层在 `PrivateDirAll` 造 `logs` 子目录时作为"missing parent"已被逐层封过（`winsec.go:271–284`：climb to first existing ancestor, sealed parent first——**父档先封再建子**）。也就是说：新机器上 `PrivateDirAll("<data>\logs")` 会把 `<data>` 与 `<data>\logs` **都以 0700 descriptor 造出**。乙形比现状真正多买的只有一件事：**`<data>` 这层的封条不依赖"恰好由 logs 触发"**，从此根的封是显式意图而不是随行者。

### 2.3 两形对照表

| 维度 | 甲形（proc 自己封） | 乙形（装配根先封） |
|---|---|---|
| 新增包级依赖边 | proc→winsec **直接边**（现仅传递） | cmd/wisp→winsec（cmd/wisp 已是根，测试已 import） |
| 循环风险 | 无环（`go list -deps ./internal/winsec` 零 wisp 包） | 无环 |
| 与 `resolve.go:35–40` 成文禁令 | 不冲突（那是反向边），但**与"低层不感知存储安全"的分层意图相逆** | 不冲突；与台账 `:8791`「装配根是唯一接缝」**同形** |
| 要动文件数 | 1 枚（boot_windows.go） | 2–3 枚（resident_windows.go、run.go、注释） |
| 新增契约级导出名 | 0 | 0 |
| 新失败方向裁定 | 要（proc 内两腿行为不一致） | 可复用 sink 降级形态 |
| 跨平台 | 不对称进 proc（§5） | 不对称留在 winsec（每平台实现自洽） |
| 对今天缺陷的边际收益 | **0**（Boot 不造根，造根处已带封） | ≈0 但**显式化**：`<data>` 根的封条从"随行者"变"声明意图" |
| 造根腿覆盖 | 仅走 Boot 的两条腿 | 仅插了的腿（同病） |

### 2.4 判定：谁更符合主流 harness

**判据（机主铁律：按主流 harness 怎么做；判据＝最有价值，⛔ 不是动作最小）**。

主流做法＝**"能力/封装边界按依赖方向定"**，三条可以逐字对上：

1. **低层包不反向依赖上层**（dependency inversion 的"依赖朝着抽象/稳定方向"）。`internal/winsec` 是叶子（§1#6），proc 是中间层，cmd/wisp 是装配根。甲形给 proc 加直接 winsec 边虽不造环，但把"存储安全策略"塞进"进程生命周期"包——proc 从此既管 Job/mutex/registry，又管 ACL，这是**把两个变化原因装进一个包**，正是主流 layered 架构要避免的。
2. **装配根（composition root）是唯一知道"整套东西怎么拼起来"的地方**。本仓自己已有同一裁定的在册先例：`pending-and-issues.md:8788–8791`（票 197 段，逐字「装配根 cmd/wisp 是唯一的接缝」，裁定＝注入、不新开边）。乙形与它同形；甲形要靠"给 BootOption 递钩子"来绕，而**钩子绕法已被票面 `:25` 第 1 刀裁过一轮**（"两枚直接边都不开，封根做成注入"）——那轮裁定的历史语境是"反向边成文禁"，今天边是顺向的，该刀的理由半失效，但**它的落点（注入＞新边）恰好与主流分层一致**。
3. **安全动作离它保护的对象越具体的层越好，而"根是谁的"是装配层的知识**。"这棵树的根在哪里"由 `buildinfo.ResolveEnv`＋`proc.Layout` 决定——proc 知道路径，但"这棵树该以什么安全姿态存在"是整个产品的决定，属于装配根。

**本腿判定：乙形更符合主流 harness。** 一句理由：乙形零新依赖边、与仓内在册的"装配根是唯一接缝"裁定同形、把封根从 sink 的随行者变成装配根的显式声明；甲形用一条新的跨层直接边换 27 行的提前量，且那 27 行内没有一次真实的盘上动作可保护。**但同时如实记录：在 132-r1 已经落地的今天，两形的边际安全收益都≈0——AC#1 的消费者（"根是谁造的"这个问题）已被 `logging.go:99` 的创建即封消解**，这也是 owner A476 否决此票的客观背景（否决发生在 132-r1 之前三天，但 132-r1 的落地让否决的实操后果更小）。若「Q-72 重开」，推荐乙形，3 枚文件，一天内可落。

---

## §3 AC#2 那把常驻尺怎么写才不是「测空气」

**票面禁的形状**（`issues/238….md:18` 逐字）：⛔「至少一枚带非 `(I)` 的 ACE」——`(I)`（inherited 标记）只说明"这条 ACE 是继承来的"，不说明"宽不宽"（`internal/winsec/acl_windows_test.go:383–387` 原话，`238-c1` §2 已引）；且 132-r1 后创建即封，那种尺今天天生绿＝测空气。

**前提**：AC#2 的射程已被裁定收窄为「本进程创建的那串目录（`R`、`R\logs`）」（票面 `:28` 第 4 刀），每枚 `*.jsonl` 的收窄属票 132——判据注释必须逐字带这句话。

**逐枚判据的形状草案**（落在 `cmd/wisp/*_test.go` 或 `internal/observe/*_test.go`，票面地界 `:52`）：

- **枚 1（根的存在性＋封态）**：fixture 造一个**不存在的**临时根路径（`t.TempDir()` 下再拼一层，断言 `os.Stat` 返回 `ErrNotExist`，同 `logprivate_132_test.go:33–36` 的形状），走生产装配链（常驻尺＝`installLogSink`；若判据射程含根，则含乙形新插的封根调用）。读数：`<R>` 与 `<R>\logs` 各 `os.Stat` 存在，且 descriptor 断言**逐枚**：
  - 目录无 `Everyone`／`BUILTIN\Users` 的**允许型** ACE（不是"数 (I)"，是"点名这两principal的 allow 必须不存在"——判宽窄不判继承）；
  - 存在当前用户 SID 的 FullControl allow。
  - 断言仪器照 `internal/observe/logprivate_132_windows_test.go` 现成的 `assertSealedDir132` 形状复制（票面 `:29` 第 5 刀：复制、逐字注明出处，⛔ 不复用 winsec 测试内符号）。
- **枚 2（正控 X1＝「注释掉新封条必红」的造法）**：关键难点是常驻尺判的是**生产链的行为**，而"注释掉"是源码变更，测试自己动不了源码。可行的造法有两档：
  - **可注入缝优先**：`internal/observe/logging.go:223` 现成 `var sealLogFile = winsec.SealFile`——132-r1 当年就是为这一发留的缝（`logprivate_132_test.go:114–123` 的 `TestSealFailureRefusesTheWrite132r1` 用它证明失败方向）。若乙形封根走 `cmd/wisp` 直调 `winsec.PrivateDirAll`，cmd 包没有等价缝 ⇒ **判据注释要写明**：正控 X1 落在 observe 层（换 `sealLogFile` 为恒错函数 ⇒ `InitLog` 必错、文件零字节，即 `TestSealFailureRefusesTheWrite132r1` 本身），常驻尺只补"根这一层"的存在性判据；若乙形落成 `BootOption` 注入钩子（第 1 刀原裁定形状），则钩子是显式参数，测试传 nil／恒错函数即得 X1，**必红路径由参数显式性保证**——这是注入形对常驻尺的真实红利。
  - **不可注入时的替代**：枚 1 的 fixture 断言"封条跑过"的唯一痕迹（目录存在＋descriptor 已窄）。X＝注释掉封条后枚 1 的 **`os.Stat` 存在性断言**仍绿（sink 还是会造目录）——所以存在性断言**不构成正控**，正控只能是 descriptor 断言＋"fixture 造的根原本无封"的前置反证（fixture 先 `icacls`/`SetNamedSecurityInfo` 放宽，再断言生产链把它收窄了＝`TestExistingWideLogTreeIsRepaired132r1`（`logprivate_132_test.go:96–107`）的形状）。**这个替代只测"有收窄行为"，测不出"谁收窄的"**——要在判据注释里如实写这句，否则就是把自己骗成 X1。
- **枚 3（时序半）**：票面 `:28` 第 4 刀把时序半留作 `238-v1` 可选进攻、**不作判据**——判据注释照抄这句，别顺手加。
- **枚 4（逐字注释草案，落判据文件头）**：
  ```
  // Scope: this gauge judges ONLY the directory string this process creates
  // (<data-root> and <data-root>\logs). Per-file narrowing of every *.jsonl is
  // ticket 132's AC#2 and is pinned by internal/observe/logprivate_132_*_test.go;
  // this gauge must not duplicate or weaken those. The (I) inheritance marker is
  // NOT what this gauge reads: it says where an ACE came from, not whether it is
  // wide (acl_windows_test.go:383-387). X-control: with the seal call disabled
  // (injected hook = nil / always-error; direct call = the observe seam
  // sealLogFile), the descriptor assertions MUST go red via
  // TestSealFailureRefusesTheWrite132r1's route or the hook's error path; if the
  // only red available is the widened-fixture repair route, say so in this
  // comment - that is a narrowing gauge, not a provenance gauge.
  ```

---

## §4 AC#3 POSIX 半边

- 现状（现读）：POSIX 上 `winsec` 的 seal＝`chmod`：`sealDir`＝`applyDescriptor(path, true)`→0o700、`sealFile`→0o600（`internal/winsec/winsec_other.go:427–432` 与 `logprivate_132_other_test.go:6–10` 的注解）；**不递归**，且 `winsec_other.go:427–431` 逐字写明理由：「on POSIX a child never inherited its parent's mode in the first place, so "seal the tree" here would be a chmod over files whose permissions somebody else set on purpose, and **the hole being closed does not exist on this platform**」。
- 票面 `:30` 第 6 刀已裁：**甲＝不递归**，注释逐字写明「POSIX 上封根不保护任何东西」引用 `winsec_other.go` 那段；**乙（补递归）拒＝推翻该段存在理由，收益零（这台不验）**。本腿复核该裁定在 132-r1 之后仍然成立：
  - 真正没被管住的 POSIX 面＝每枚 jsonl 的 mode：`logging.go:297` 现读＝`os.OpenFile(..., 0o600)`——**已不是票面 `:30` 说的 `0o644`**（`grep -n "0o644" internal/observe/logging.go` 现跑＝零命中，仅 `logging.go:95–98` 注释里提过"used to pass 0o755"）。即票 132 落地时连 POSIX 文件 mode 一起收了。⇒ **AC#3 在今天的代码上已无事可做**。
  - 若「Q-72 重开」且仍要补判据注释，文案草案：`// POSIX: the root seal here is a single chmod (winsec_other.go:427-432 - "the hole being closed does not exist on this platform"); this platform is verified for placement and mode, not for ACL inheritance, and "it passed in a container" is NOT accepted as the verification.`（照票面 `:19` ⛔"不许用容器里跑过了代替"）
- 递归方案（若未来有人要推翻第 6 刀）：改动＝`winsec_other.go` 新增递归走查＋`winsec.go` 加导出名（契约级，`Q-72` 同款）＋`private_other_test.go`/`placement_symlink_113_other_test.go` 两族尺补断言；代价＝推翻 `:427–431` 的成文存在理由（违反"改契约须人工批准"的同一级动作）。**本腿不推荐**（复杂买不到真价值：POSIX 的洞在继承语义，不在遍历）。

---

## §5 交付自证（Git 纪律遵守）

- 本腿全绿清单：零 `.go` 改动、零删文件、零 `go test`、零 `go build`（`go list -deps` 为纯只读包图查询，不触发编译缓存写入被撞包）。
- 工作树入口自证：`git status --porcelain -- cmd internal` 起手读数＝仅 `?? cmd/wisp/panel_resident_323_windows_test.go`、`?? cmd/wisp/panel_resident_assembly_323_test.go` 两枚他人未跟踪件（现跑记录于 §1，本腿未触碰）。
- 本次 commit：仅本文件，显式 pathspec（`git commit -- .scratch/wisp/probes/238/a1/census.md` 形状），⛔ push、⛔ add -A / amend / reset / rebase / stash / checkout . / clean。

## §6 next（给编排者）

1. 本腿结论与 A476/A1558/`238-proc` STOP 三笔在册记录同向：**缺陷主体已被 132-r1 收口，AC#1..3 的消费者不存在**。除非 owner 说「Q-72 重开」，本普查件与 `238-c1`、`238-proc` 同等，**保留当账、不当待办**。
2. 若重开：推荐乙形（§2.4），3 枚文件；AC#2 尺按 §3 注入缝形状落，正控 X1 的诚实边界（narrowing gauge ≠ provenance gauge）写进判据注释；AC#3 按 §4 文案落注释即可，零递归。
