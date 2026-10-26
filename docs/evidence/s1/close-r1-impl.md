# 收口腿 close-r1 交件（写腿·实现者）——票 279／352 的最后一格

★起手锚＝`cbb43c948`（`dev`，19:09）。★★**本腿零产码**：只写本件一枚，⛔ 碰任何 `.go`。
⛔ 翻框／⛔ 改票面／⛔ `Status:`／⛔ `-done`／⛔ 改台账／⛔ push。★腿末锚＝`0855c168f`（测量期间仓库被他腿推进三次，见 §8）。
★读数快照时刻＝**2026-10-26 19:39:37 +08**；该刻 `git status --porcelain -- internal`＝**0 行** ⇒ 本件所有 `internal/**` 读数出自同一份源码态。

## 1. 两张票的最后一格：原文／判定／读数

### 票 279 — 原文（`.scratch/wisp/issues/279-p39-…md:14` 逐字）

> - [ ] **AC#1 同笔联动**：P39 贴＋`config_readers_255.go:109` 那枚引用改到新行号（`[cfg := rt.cfg]` 那句原文不动）＋复跑那枚测试 **rc=0**；⚠ 若测试的**期望**也要改，逐字给"原句/新句"与理由。

★判定＝**不属本腿，具名交回**（该格两处写面全在 `cmd/wisp/**`＝硬边界 1，且派单停手上报②点名此形）；
★★**另加一条：该格前提今天已丧失，⛔ 应按原字面派给任何落地腿**。四读数：

- **R1 P39 的前置块已不在树里**（19:1x 现读）：`git grep -c 'Nothing calls it yet' -- cmd internal`＝**0 命中**；`run.go` 全文无 `ParseComposerRequest`；料文 `:329` 点名的 `run.go:331-336` 现是 `liveCards`／`reply` 字段段（`:331-337`）。
- **R2 真因＝被后来一笔整段重写**：`114-r2`（`80376d58f`，2026-10-24）改写了 `modeWrites` 注释，现文 `run.go:350-363` 已逐字写"The FIELD still has no reader in production…"，并把 114-r2 凭什么关掉那个缺口讲得**比 P39 替换文本更准** ⇒ 照 P39 贴＝**回退 114-r2 的说明**，是毁账不是修账。
- **R3 该格"复跑那枚测试 rc=0"的联动约束亦失效**：`config_receipt_255_test.go:201-227` 现只断言**令牌在该文件里唯一**（`hits` 长度＝1），行号漂移只 `t.Logf`（用例自陈"the line number never is"）；票面 `A` 段"按 `lines[n-1]` 直接取第 N 行"是**旧尺**（`evidenceCite` 仍在 `config_readers_255.go:270`，但用法已改）。
- **R4 但名册的行号账确实过期**（⛔ 静默过期，属"没人报错的错账"）。尺＝逐行现读 + 令牌定位：

| 名册行 | 现写的 | 19:39:37 现量 | 漂 |
|---|---|---|---|
| `:56` | `run.go:426-427`（`cfg := mgr.Config()`） | `run.go:466` | **+40** |
| `:58` `:114` `:117` | `run.go:1023 [cfg := rt.cfg]` | `run.go:1252` | **+229** |
| `:107` | `run.go:438 [res := llm.NewResolver(cfg, st)]` | `run.go:478` | **+40** |
| `:113` | `run.go:427 [rt.cfg = cfg]` | `run.go:467` | **+40** |
| `:115` | `run.go:1046 [PerToolTimeout: …]` | `run.go:1278` | **+232** |

⇒ **交回内容**：⛔ 别重贴 P39（改由接管腿把该格判为"被 114-r2 取代"）；只做上表五处行号刷新。⚠ 该刻 `cmd/wisp/run.go` 正被在飞腿改写（` M`），**执行前必须重取**。

### 票 352 — 原文（`.scratch/wisp/issues/352-wrappath-…md:16` 逐字）

> - [ ] **AC#2 消费者面不成对⛔ 结案**：`internal/secret`⇄`config`⇄`memory`⇄`agent`⇄`panel` 五包各自跑色（终裁腿那 38 枚判定点的归属包）＋★`cmd/wisp` 整包⛔ 本趟跑（归编排者那窗，防互抢）；落地⇄基线成对逐名作差＝新增红 0 枚。★既知连带具名：`private_fail_test.go:37/50/61/72/75` 五枚 `ErrNotSealable` 判定须仍绿；读 `ErrUnresolvedPath` 的 `openleaf_322`／`placement_symlink_113` ⛔ 受影响。

★判定＝**成立**。★补的正是终裁腿 `352-v1-verdict.md:35` 具名缺的那一味「落地后**现势跑色**」：本腿在现势树上现跑成对，⛔ 只搬底片。★`cmd/wisp` 整包**一枚没跑**（遵票面）。落地那一行已在盘上并逐字复核：`internal/winsec/winsec.go:452`＝`return fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)`（落地笔 `ecffb93b`，`git show -U0` 净变动 1 增 1 删）。

## 2. 票 352 的「成对」怎么成对：基线是什么、三条边读数

★**基线定义（仓外，零工作树突变）**：取工作树 `internal/winsec/winsec.go`（blob `fc4e0d69` ＝ `HEAD:` 同值），**只**把已批准那半 `%w` 改回 `%v`，落 `D:\tmp\wisp\20261026\closer1\ovl\winsec-baseline.go`，用 `go test -overlay=` 挂进编译。★判定尺＝`git diff --no-index -U0` 两侧**只差 1 行**（`%w`⇄`%v`），其余字节同。⇒ 对照物与落地件**只差机主批的那一枚字符**，比拿历史提交当基线更干净。

**三条边（读 errno 分类的消费者，`errors.Is(err, fs.ErrNotExist)` 那侧）**——三条边两侧同数、同红集：

| 边 | 消费点（现读） | 落地 run/pass/fail/skip | 基线 run/pass/fail/skip | 新增红 |
|---|---|---|---|---|
| `internal/secret` | `store.go:140`／`list.go:38`／`migrate.go:88`／`configrefs.go:25` | 29/29/0/0 | 29/29/0/0 | **0** |
| `internal/config` | `parse.go:234`（`settings_257_test.go:286` 断言它） | 205/205/0/0 | 205/205/0/0 | **0** |
| `internal/memory` | `artifacts.go:115/125/152/204/214/264/286/307`（`winsec.RemoveUnlinked` 的 errno） | 134/132/0/2 | 134/132/0/2 | **0** |

★另两包（票面也点名）：`internal/agent` **138/138/0/0** ⇄ 同；`internal/panel` **269/265/4/0** ⇄ **268/263/5/0**、新增红 **0**。
★panel 那 4 枚＝**存量红**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`），两侧逐名同；基线多出的那 1 枚 `TestTheBoundDispatchDoorIsReachableFromThePage` **不是本树之红**，见 §8 具名。
★**两发独立成对**（pair-1 全六包、pair-3 五包）⇒ 五包**新增红 0／消失红 0**（pair-3 那 1 枚"消失"＝上句那枚瞬时假红）。

★**第六包 `internal/winsec`（改动本包，票面五包之外，照实报）**：127/125/2/0 ⇄ 127/125/2/0，新增红 0；但那 2 枚是**竞态探针**，须具名：`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（每发 400 次尝试，命中 1 次即红）。
**交错量（landed⇄baseline 各 5 发，逐发交替）＝ landed 3/5 红（witness 1~3 次/400）⇄ baseline 5/5 红（witness 1~2 次/400）** ⇒ **两侧都咬、与那一枚字符无关**，⛔ 归本票。剔掉该族后 winsec 两侧**确定性红＝0／0**。

★**既知连带逐名（票面点名）**：`private_fail_test.go` 五枚 `ErrNotSealable`（行 37/50/61/72/75，宿主 `TestAC5FailedSealRefusesTheWrite`＋四子测试）两侧**逐名 pass**；`openleaf_322`（`openleaf_guard_322_windows_test.go`）两侧 pass；★`placement_symlink_113_other_test.go` 是 `//go:build !windows`（8 枚）⇒ **本平台未跑**，照实报"未触及"而不是"绿"。★读 `ErrUnresolvedPath` 的 **9 枚 windows 测试文件共 30 个用例**：两侧 **30/30 pass**、非 pass **0**。

## 3. 反向判据（票 352 一发）：把改动摘掉 ⇒ 必须红

★器械＝仓外探针 `closer1\probe\zz_closer1_probe_test.go`（`package winsec`，只取**真 errno**：`os.Open` 一个父目录不存在的路径，再喂 `wrapPath`），`-overlay` 挂入，⛔ 工作树零突变。判据＝机主批的那句"两分类并存"：`errors.Is(err, ErrNotSealable)` **且** `errors.Is(err, fs.ErrNotExist)`，另加反向控制"不得冒出 `ErrUnresolvedPath`"。

| # | 读数 | 命令（`-overlay` 三种） | rc | 逐字 |
|---|---|---|---|---|
| ① | **变异红**（摘掉改动：`%w`→`%v`） | probe＋baseline | **1** | `the errno no longer survives wrapPath, so a caller asking whether the path is missing still cannot get a true answer: winsec: cannot apply a private security descriptor: leaf.bin: open …\no-such-parent\leaf.bin: The system cannot find the path specified.` |
| ② | **还原绿** | probe（改动在） | **0** | `--- PASS: TestCloseR1WrapPathCarriesBothSentinels` |
| ③ | **原样绿** | 无 overlay（仓自带尺） | — | ★仓里**没有**一枚用例读这个判据（票面 `:5` 已具名"全仓无真依赖它的调用者"）⇒ 本腿**不假装仓里会红**：如实报"仓自带尺对此无牙，牙在探针"；替代读数是 §2 的"两侧成对新增红 0"。 |
| ④ | **还原核**（时机＝**最后一发变异之后**） | `git hash-object` ⇄ `git rev-parse HEAD:` | — | worktree `fc4e0d69…` ＝ HEAD `fc4e0d69…`；`git status --porcelain -- internal/winsec`＝**0 行**（在 19:39:37 快照复验，且在最后一次变异之后） |

★**票 279 的反向判据＝本腿发不出**（其写面全在 `cmd/wisp/**`，⛔ 碰＝越界）；对 279 的交付是 §1 那四条**可复跑的现量**，不是产物。★照实登记本腿自抓一发尺错：首轮 `-overlay=$ov` 在 PowerShell 里未展开（`go: reading overlay: open $ov`），那 1 发**作废、⛔ 计入任何数**，已用字符串拼接重跑（§2 的 5 发交错即重跑后的读数）。

## 4. 零回归 + 四数

| 包 | 改前（基线） | 改后（落地） | 判 |
|---|---|---|---|
| `internal/secret` | 29/29/0/0 | 29/29/0/0 | 零回归 |
| `internal/config` | 205/205/0/0 | 205/205/0/0 | 零回归 |
| `internal/memory` | 134/132/0/2 | 134/132/0/2 | 零回归 |
| `internal/agent` | 138/138/0/0 | 138/138/0/0 | 零回归 |
| `internal/panel` | 268/263/5/0（含 1 枚瞬时假红） | 269/265/4/0 | 零回归（存量红两侧同） |
| `internal/winsec` | 127/125/2/0 | 127/125/2/0 | 零回归（342 竞态族见 §2） |

★**本腿改的包＝0 个**（零产码）⇒ "改前/改后"即上表两侧，四数＝**run/pass/fail/skip**。
★`internal/speech`＝**零字节**（`git status --porcelain -- internal/speech`＝0 行）；`internal/session`＝**零字节**（同尺 0 行）。
★`cmd/wisp`＝**本腿零字节**（本腿 ⛔ 一次 `write`/`edit` 指向该树）。⚠ **照实登记异动**：该树在腿内被他腿写了 5 项——`M cmd/wisp/run.go`／`M cmd/wisp/resident_windows.go`／`M cmd/wisp/resident_voice_windows.go`／`?? cmd/wisp/resident_voice_engines_windows.go`（19:21:17，28062B）／`?? cmd/wisp/voice_out_windows.go`（19:18:50，13869B）。⇒ 归在飞装配腿，⛔ 本腿（本腿对它只有读）。

## 5. 格式甲乙两尺

★两尺在 `cmd internal tools` 上**逐行同**：甲 `gofmt -l`＝4 行、乙 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`）＝4 行，**rc 均 0**。
★归属判读用 `[\\/]`（8/8 行带分隔符，`NO_SEPARATOR`＝0）：

- **既存 3 条（⛔ 本腿责任）**：`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go` —— 三枚均 TRACKED 且 ` M`＝0（改前就在）。
- **1 条外来**：`cmd/wisp/resident_voice_engines_windows.go` ＝ **UNTRACKED**、19:21:17 由他腿新建 ⇒ 非既存 3 条之列，也非本腿。

★另注：`gofmt -l .`（全树）rc＝2、49 行，全部落在 `.scratch/**` 探针件（含 4 行语法坏件）⇒ **不是仓门禁**，仓门禁是目录/跟踪件那两尺。

## 6. `d22scan`

★**工作树（含他腿在飞件）**：`sh scripts/d22scan.sh` **rc＝1**，且红在第 1 步正控——33 个顶层用例 **PASS=33／FAIL=2**，两枚红同因一处、**均非本腿**：

- `TestScannerSelfScanOfRealRepoIsGreen`／`TestRealRepoLedgerIsHonest` 都在跑时重扫真实工作树，命中：
  `cmd/wisp/resident_voice_engines_windows.go:695` 与 `:817` ＝ `[emoji] ban #8 glyph in scope cmd/ is banned (D23)`（他腿 19:21 新建的那枚 UNTRACKED 件）。

★**HEAD 导出树（仓外 `git archive HEAD` ⇒ `headtree\`，1225 枚 `.go`，无在飞件、无 `.git`）**：`go run . -root <headtree>` **rc＝0 CLEAN**——`examined 306 production Go files`，ban#8 覆盖 design/14＋frontend/116＋internal/679＋cmd/155，逐 scope 全 clean。⇒ **提交态干净；工作树之红 100% 来自他腿在飞件**（导出件数为 306 ⇄ 工作树 308）。
★遵令**未**用 ban#8 扫 `docs/`。

## 7. 写面自证

★本腿写面＝本件一枚（`docs/evidence/s1/close-r1-impl.md`）＋仓外 `D:\tmp\wisp\20261026\closer1\**`（只建不删：`sh\` 3 枚驱动／`ovl\` 基线件＋3 枚 overlay／`probe\` 探针／`logs\` 50 枚原样日志／`headtree\` 导出树）。
⛔ 未碰 `cmd/**`、`internal/speech/**`、`internal/session/**`、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/specs/**`、`docs/PLAN.md`、`docs/audit/**`。
★`git add`／`commit` 两处**显式逐枚 pathspec**，`-F` 写在 `--` 之前；⛔ `add -A`／amend／reset／rebase／stash／checkout ./clean；⛔ push。

## 8. 两条须交编排者处置的现场事实（⛔ 本腿可修）

1. **`cmd/wisp` 在腿内被他腿持续改写**（三次 HEAD 推进 `cbb43c948`→`385afead9`→`0855c168f`；`run.go` 的 `cfg := rt.cfg` 从 `:1246` 漂到 `:1252`；19:38 曾读到**半保存态**）⇒ 任何 `cmd/wisp` 行号读数（含 §1 那张表）**执行前必须重取**。
2. **两条门禁因他腿在飞件而红、⛔ 因本票**：D22 的 2 枚 ban#8（`:695`／`:817`）；`internal/panel` 的 `TestTheBoundDispatchDoorIsReachableFromThePage` 在两发里红了一次，逐字＝`cmd/wisp/resident_voice_engines_windows.go does not parse (…:323:10: rune literal not terminated (and 2 more errors)), so the door binding in it cannot be read and the host half of the tie is blind`（该用例 `WalkDir` + `os.ReadFile` 读 `cmd/wisp/**`）⇒ **他腿收线后两处自消**；独跑该用例 3 发 × 两侧＝6/6 绿。
