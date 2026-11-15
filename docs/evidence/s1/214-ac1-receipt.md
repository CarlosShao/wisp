# 214 — AC#1 收件（受理器生产调用点 0→N 读数做实）

- 腿：`214-ac1-receipt`（**收件腿／非实现者**，只读＋裁决；写面只有本文件与 `.scratch/wisp/probes/214/**`）
- 裁定时刻盘上 HEAD：`d0639eedfe5539e78bfedfe59a83ce19bd20293c`（2026-11-15 13:27:58 +0800）
- 被收的对象：`docs/evidence/s1/214-attachment-wiring.md`（实现腿 `214-attach-wire`，提交 `fc76d5d7e96664790d4d21904ce72fda6cf18883`）
- 读面的禁令照办：`cmd/wisp/**` **只读**；⛔ 未跑 `go test ./cmd/wisp` 全包（见 §7 的〔未跑包〕说明）；⛔ 未改产码／票面／台账／阈值／契约；⛔ 未 push。
- 台账 **A1875** 两件套：①生产调用点 `file:line`＝§2 四处；②会因摘掉接线而变红的测试名＝`Test214AttachmentAddEndToEndStoresBytesInArtifacts`（§4）＋同链另三枚（§4.4）。**两件都齐** ⇒ AC#1 判**满足**（§5）。

---

## 0. 起手三核（现跑，非转述）

| 核 | 读数 |
|---|---|
| ① 票面 AC#1 原文还在不在 | `grep -n 'AC#1' .scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md` ⇒ `:37`，逐字见 §1 |
| ② 实现腿自报的落点今天还在不在 | `cmd/wisp/panel_inbound.go:386` `Attachment: attachment,` **在**；`cmd/wisp/panel_attachment_214.go:204` `panel.NewAttachmentBroker(...)`、`:243` `broker.Ingest(ctx, src)` **在**（实现腿写 `:204`／`:243`／`panel_inbound.go:351`，**第三处行号已腐烂**：`:351` → 现量 `:386`，事实仍在） |
| ③ 收件前工作树干不干净 | `git status --porcelain -- cmd/wisp` ⇒ **空**（无别的腿的在飞临改）；`Get-FileHash -Algorithm MD5 cmd/wisp/panel_inbound.go` ⇒ `53243FD2654812D927A10313EFAE9E16`（§4 变异跑完复取同值） |

> ⚠ 实现腿 `214-attachment-wiring.md` 写的 `panel_inbound.go:351` 是**当时**的 struct literal 行；今天同一行在 `:386`。本件全部用**现量行号**，并在此具名更正，⛔ 不改票面、⛔ 不改别的证据件。

---

## 1. AC#1 原文（逐字，票面 `:37`）

> - [ ] **AC#1**：受理器有**生产调用者**（现跑 `grep` 点数：改前 0 枚 ⇒ 改后 ≥1 枚，且那枚调用者在 `cmd/wisp` 的装配路径上，不在测试里）。

三个子条件拆开：**(a)** 改前 0 枚、**(b)** 改后 ≥1 枚且逐处点名、**(c)** 那枚调用者在 `cmd/wisp` **装配路径**上（非测试）。

---

## 2. 现跑普查：受理器生产调用点 **4 处**（逐处 `file:line`）

尺与读数（本腿现跑，`git --no-pager grep`，非测试文件）：

```
git --no-pager grep -n -E 'NewAttachmentBroker|\.Ingest\(' HEAD -- internal/ cmd/ | grep -v '_test.go'
```

| # | 生产调用点（`file:line`） | 逐字原文 | 它在这条链上的位置 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_inbound.go:343` | `attachment := newAttachmentIngestHandler214(dir, drafts, auditf, actor)` | **装配路径**：`newComposerDispatchChain` 体内构造受理器 handler（AC#1 那句"在 `cmd/wisp` 的装配路径上"指的就是这一站） |
| 2 | `cmd/wisp/panel_inbound.go:386` | `Attachment: attachment,` | **装配路径**：把 handler 接进 `panel.ComposerDispatch` 的 `Attachment` 槽（改前这一行是 `Attachment: nil`，见 §3-C3） |
| 3 | `cmd/wisp/panel_attachment_214.go:204` | `h.broker, h.openErr = panel.NewAttachmentBroker(store, memory.ValidArtifactName, 0)` | **受理器本体第一次被生产构造**（真 `memory.Store` 的 artifacts 树、票 76/79 守卫） |
| 4 | `cmd/wisp/panel_attachment_214.go:243` | `return broker.Ingest(ctx, src)` | **受理器本体第一次被生产调用**（`Ingest`＝票面 §现量点名的"受理器"那一枚） |

另有一枚**派发点**（不在上表计数内，且**不是本票引入**）：

| — | `internal/panel/composer_dispatch.go:235` | `return "", d.Attachment.HandleAttachmentRequest(ctx, req)` | 路由器把 `panel.attachment.add` 交给 `Attachment` 槽；改前这一行就在（`25eb216…:214` 同形），差别只在于槽里是 `nil` ⇒ 走到 `composer_dispatch.go:232-233` 的 `d.unattached(req)` |

**计数口径（⛔ 不混着说）**

- **口径 A（票面 §现量原尺＝"谁调用受理器"）**：`0 → 2`（上表 #3、#4）。这正是实现腿自报的 `0→2`，本腿复跑同值。
- **口径 B（本件的判据口径＝受理器整条链的生产站点，含装配与路由器接线）**：`0 → 4`（#1–#4）。
- 两个口径同向、不矛盾；AC#1 要的是"≥1 且在装配路径上"，**两个口径都成立**。本件标题里的 **N＝4**。

**生产可达性（"在装配路径上"这一条的正面凭据，逐跳 `file:line`）**

| 入口 | 逐跳 |
|---|---|
| CLI 缝 | `cmd/wisp/main.go:187 case "panel-inbound":` → `cmd/wisp/panel_inbound.go:157 disp, err := newPanelInboundDispatch(dir, auditf)` → `:246 newComposerDispatchChain(...)` → `:343`＋`:386` → 循环里 `cmd/wisp/panel_inbound.go:174 reply, err := disp.Handle(ctx, raw)` |
| 常驻（`wisp run`）缝 | `cmd/wisp/panel_resident_windows.go:360 disp, cfgMgr, workspaceScope, drafts, err := newResidentComposerDispatch(...)` → `:296` 函数体 → `:306 drafts := newAttachmentDraftLedger()` → `:312 newComposerDispatchChain(dataDir, auditf, residentPanelActor, scope, confirm, drafts, preview)` → 同一条 `:343`＋`:386` → `cmd/wisp/panel_host_windows.go:1247 return m.disp.Handle(ctx, raw)` |

⇒ 两条产线入口**共用同一枚** `newComposerDispatchChain`，受理器挂在它里面 ⇒ **一处接线、两条缝都活**。两处 `disp.Handle` 都在**非测试**文件里。

---

## 3. "改前 0"这一半的依据（`git log -S` ＋ 改前树现跑 grep）

### 3.1 `git log -S`：引入调用点的那一枚提交

```
$ git --no-pager log --oneline -S 'panel.NewAttachmentBroker' -- '*.go'
fc76d5d7e 214: 附件接线——panel.attachment.add 挂上真受理器，附件真传到落盘与快照

$ git --no-pager log --oneline -S 'Attachment: attachment' -- '*.go'
fc76d5d7e 214: 附件接线——panel.attachment.add 挂上真受理器，附件真传到落盘与快照

$ git --no-pager log --oneline --diff-filter=A -- cmd/wisp/panel_attachment_214.go
fc76d5d7e 214: 附件接线——panel.attachment.add 挂上真受理器，附件真传到落盘与快照
```

- 引入提交：`fc76d5d7e96664790d4d21904ce72fda6cf18883`（author `CarlosShao`，`Sat Nov 7 17:52:15 2026 +0800`）
- 父提交（＝"改前"那一棵树）：`25eb216aa8a9b71c216bb19f5111bd80459fe328`
- 该提交 `--name-status`：**新增** `cmd/wisp/panel_attachment_214.go`（＋其测试件、证据件），**改** `cmd/wisp/panel_inbound.go` 等 ⇒ 受理器生产调用点**随这一枚提交从 0 变有**。

`-S 'broker.Ingest'` 另命中的三枚（`0eaf5a6a2`／`dd2506f06`／`92e8cc575`）本腿逐枚 `git show --name-only` 核过：全部只动 `*_test.go` 与 `.scratch/wisp/probes/**`，**零产码** ⇒ 不构成新的生产调用点（HEAD 非测试 grep 也只有 `:243` 一处，同向）。

### 3.2 改前树现跑 grep（不是从 `-S` 推的，是直接在父提交那棵树上量的）

| 尺 | 父提交 `25eb216…` 读数 | HEAD 读数 |
|---|---|---|
| C1 **票面原尺**：`git grep -n -E 'panel\.Attachment\|Attachments\(' <rev> -- internal/ cmd/ \| grep -v internal/panel/attachments` | **0 行**（空结果） | **15 行**：生产 12 行（`cmd/wisp/panel_attachment_214.go` 11 行 ＋ `internal/panel/pump.go:427` 快照读者那一半）＋测试 3 行（`panel_attachment_214_test.go:206/:263/:317`） |
| C2 `git grep -n -E 'NewAttachmentBroker\|\.Ingest\(' <rev> -- internal/ cmd/ \| grep -v _test.go` | **2 行，全是它自己的定义**：`internal/panel/attachments.go:175` 注、`:179` `func NewAttachmentBroker(...)` ⇒ **调用点 0** | **4 行**：2 行定义 ＋ **2 行调用**（`panel_attachment_214.go:204`、`:243`） |
| C3 `git grep -n 'Attachment:' <rev> -- internal/ cmd/ \| grep -v _test.go` | `cmd/wisp/panel_inbound.go:331  Attachment: nil,` ⇒ **槽是明写的 nil** | `cmd/wisp/panel_inbound.go:386  Attachment: attachment,` |
| C4 `newAttachmentIngestHandler214` | 0 行（`git cat-file -e` ⇒ `fatal: path 'cmd/wisp/panel_attachment_214.go' exists on disk, but not in '25eb216…'`，rc=128） | 3 行（`:157` 注、`:166` 定义、`panel_inbound.go:343` 调用） |

⇒ **改前 0**：三个独立尺（票面原尺／受理器本体调用尺／接线行尺）在**同一棵父提交树**上同时读到 0／nil，互相对得上；改后 ≥1（且每一处都已点名）。

---

## 4. 会红的测试：摘掉那处接线 ⇒ 哪枚红

### 4.1 变异怎么做的（⛔ 改动没留树上）

- 载体＝**`go build -overlay`**：只把 `cmd/wisp/panel_inbound.go` 的**一份副本**挂进构建，盘上原文件一字节不动。
- 副本：`.scratch/wisp/probes/214/ac1/mut-attachment-nil-panel_inbound.go`
- overlay：`.scratch/wisp/probes/214/ac1/overlay-attachment-nil.json`
- 变异**只有一行**（`git diff --no-index` 逐字）：

```
-		Attachment: attachment,
+		Attachment: func() panel.AttachmentRequestHandler { _ = attachment; return nil }(), // MUT-214-AC1 overlay: the one wiring line pulled
```

  ⇒ 结构性改动只有"槽从真 handler 变回 nil"；`_ = attachment` 只为编译得过（不然 `attachment` 未使用），**不改变行为**。
- 盘上无残留：跑完 `Get-FileHash -Algorithm MD5 cmd/wisp/panel_inbound.go` 仍是 `53243FD2654812D927A10313EFAE9E16`、`git status --porcelain -- cmd/wisp` 仍为**空**。

### 4.2 正控（不摘接线 ⇒ 同一枚测试绿）

```
$ go test ./cmd/wisp -run '^Test214AttachmentAddEndToEndStoresBytesInArtifacts$' -count=1 -v
--- PASS: Test214AttachmentAddEndToEndStoresBytesInArtifacts (0.26s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.417s        （日志：.scratch/wisp/probes/214/ac1/green-baseline.log）
```

### 4.3 反控（摘掉那处接线 ⇒ 同一枚测试红）

```
$ go test ./cmd/wisp -run '^Test214AttachmentAddEndToEndStoresBytesInArtifacts$' -count=1 -v -overlay .scratch/wisp/probes/214/ac1/overlay-attachment-nil.json
--- FAIL: Test214AttachmentAddEndToEndStoresBytesInArtifacts (0.02s)
    panel_attachment_214_test.go:88: an accepted attachment came back refused:
      panel: 该方法在名册内，但本机未接入处理器: 方法 "panel.attachment.add" 的处理器未接入（requestId="r-214-accept"），档位/工作区/附件/消息/设置均未变化
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.130s
```

（日志：`.scratch/wisp/probes/214/ac1/red-mutant-attachment-nil.log`）

**⇒ 会红的测试名：`Test214AttachmentAddEndToEndStoresBytesInArtifacts`**（`cmd/wisp/panel_attachment_214_test.go:76`，红在 `:88`）。红因逐字是"该门未接入处理器"，即**摘掉的正是 AC#1 那处接线**，不是别的面抖红。

### 4.4 同链另三枚也红（变异有鉴别力，不是一片红）

`-run '^Test214'` 同一条件跑（日志 `red-mutant-test214-group.log`）：

| 测试 | 变异后 | 红句（`panel_attachment_214_test.go`） |
|---|---|---|
| `Test214AttachmentAddEndToEndStoresBytesInArtifacts` | **FAIL** | `:88` an accepted attachment came back refused |
| `Test214AttachmentRefusalReachesTheSnapshotWithItsReason` | **FAIL** | `:137` a non-image payload was not refused with the sentinel |
| `Test214DisguisedExecutableIsRefusedByItsBytes` | **FAIL** | `:181` an executable named .png was not refused |
| `Test214PathSpellingDisplayNameIsRefusedAndNothingLands` | **FAIL** | `:243` a path-spelling name was not refused |
| `Test214OversizeAttachmentIsRefusedBeforeAnyRead` | PASS | （它判的是"拒＋零读"，nil 门也是拒 ⇒ 本变异对它无鉴别力，**不算它的正控**） |
| `Test214IngestRefRenderingCarriesTheArtifactPathAndRefusalsToTheAgent` | PASS | （直测渲染，不经门） |
| `Test214PumpAssemblyCarriesTheAttachmentReaders` | PASS | （钉的是 `run.go` 快照字面量，不是这道门） |
| `Test214DraftLedgerDropsOldestPastItsCap` | PASS | （ledger 自己的界） |

⇒ 4 红／4 绿：变异**只打穿"经这道门"的那一组**，说明红的因果链对得上，不是整包塌。

### 4.5 ⚠ 与本票另一枚既有钉的分界（防误认）

`TestAC9ComposerDispatchHasAProductionCaller`（`cmd/wisp/panel_inbound_33_test.go:182`）钉的是 **`ComposerDispatch.Handle` 有生产调用者**（票 33 AC#9），**不是**附件受理器的调用者；本变异摘的是 `Attachment` 槽而留 `disp.Handle`，它**按设计不红**。AC#1 的牙在本节 4.3 那一枚，别拿 33 的那枚冒充。

---

## 5. 逐格判语（五列）

### AC#1

| AC 原文 | 判语 | 依据 | 可核命令 | 缺什么 |
|---|---|---|---|---|
| 受理器有**生产调用者**（现跑 `grep` 点数：改前 0 枚 ⇒ 改后 ≥1 枚，且那枚调用者在 `cmd/wisp` 的装配路径上，不在测试里） | **满足**（A1875 两件套齐：①四处 `file:line` 见 §2；②会红的测试 `Test214AttachmentAddEndToEndStoresBytesInArtifacts` 见 §4.3） | ① 生产调用点 4 处：`cmd/wisp/panel_inbound.go:343`／`:386`／`cmd/wisp/panel_attachment_214.go:204`／`:243`，全在非测试文件、全在 `newComposerDispatchChain` 这条装配路径上，两条产线入口（`main.go:187` CLI 缝、`panel_resident_windows.go:360` 常驻缝）都从 `main` 可达（§2 末表）；② 改前 0：父提交 `25eb216…` 三把尺同时读到 0／`Attachment: nil`／文件不存在（§3.2）；引入提交 `fc76d5d7e…`（§3.1）；③ 摘线变异 ⇒ `Test214AttachmentAddEndToEndStoresBytesInArtifacts` 红在 `:88`，还原后绿（§4.2／§4.3） | `git --no-pager grep -n -E 'NewAttachmentBroker\|\.Ingest\(' HEAD -- internal/ cmd/ \| grep -v _test` ／ `git --no-pager grep -n 'Attachment:' HEAD -- cmd/wisp/panel_inbound.go` ／ `git --no-pager log --oneline -S 'panel.NewAttachmentBroker' -- '*.go'` ／ `git --no-pager grep -n -E 'panel\.Attachment\|Attachments\(' 25eb216aa8a9b71c216bb19f5111bd80459fe328 -- internal/ cmd/` ／ `go test ./cmd/wisp -run '^Test214AttachmentAddEndToEndStoresBytesInArtifacts$' -count=1 -overlay .scratch/wisp/probes/214/ac1/overlay-attachment-nil.json` | 无（就 AC#1 这三个子条件而言）。**边界两条**：① 本件只裁 AC#1，别的格未裁（§6）；② "界面能点"不在 AC#1 的文字里（它在 AC#6），⛔ 本件的满足不许被读成"附件今天能点" |

### 其余五格（本腿未裁，具名交回）

本件是 AC#1 的收件腿，射程只到 `:37` 一格。其余五格**今天没有本腿的裁决**，逐格写明归属，⛔ 不许把本件的"满足"外溢过去：

| 格 | 票面 | 本腿状态 | 归属／理由 |
|---|---|---|---|
| AC#2 超限拒＋快照可见、限内收 | `:38` | **未裁** | 实现件 §3 给了两枚测试名（`Test214OversizeAttachmentIsRefusedBeforeAnyRead`／`Test214AttachmentAddEndToEndStoresBytesInArtifacts`），但**快照那一半要非实现者另裁**；且本腿变异发现超限那枚**对 nil 门无鉴别力**（§4.4），按 A1875 它单独不足以充当"会红"的正控 ⇒ 该格要另立收件腿，⛔ 本件不代裁 |
| AC#3 伪装文件按嗅探判 | `:39` | **未裁** | 实现件 §3 只给了 `Test214DisguisedExecutableIsRefusedByItsBytes`；本腿变异证明它**确实是这道门的牙**（§4.4 第三行），但"按嗅探不按扩展名"的语义还要现读 `matchesISOBaseMedia`（`internal/panel/attachments.go:140-149`）另裁 |
| AC#4 落盘位置在授权根内 | `:40` | **未裁** | 票 174/175 同族判据，需另腿 |
| AC#5 来源戳＋不可信外部文本 | `:41` | **未裁（已知未完成）** | 实现件 §3 自陈"C25 盖戳半边未完成"，且 `docs/evidence/s1/214-ac5-stamp-feasibility.md` 已裁"乙"并指向**须新立一枚票**（工单池今天无 owner）⇒ 这格由编排者处理，不属 AC#1 |
| AC#6 三件"没接之前不许说已具备" | `:42` | **未裁** | 实现件 §3 写了三件（界面点不了／无选择器／无 @引用）都在；本腿零前端写面，不可判 |

---

## 6. 诚实边界（本件**不是**什么）

1. **本件不翻勾**：票面 `:37` 的 `- [ ]` 一个字没动（裁决归非实现者，翻勾归编排者；本腿的写面只有本文件与 `.scratch/wisp/probes/214/**`）。
2. **本件不证明"能点"**：AC#1 只要求"有生产调用者"。附件今天**没有页面选择器、没有 @引用**（AC#6），`frontend/**` 零写面。
3. **`TestAC9ComposerDispatchHasAProductionCaller` 不构成 AC#1 的证据**：它钉 `Handle`，不钉 `Attachment` 槽（§4.5）。
4. **⛔ 不读注释当证据**：`internal/panel/composer_dispatch.go:165-166` 今天仍写着 `// Attachment answers the attachment request. Untouched by this slice.` ⇒ **这句已过期**（`Attachment` 自 `fc76d5d7e` 起由 `cmd/wisp` 装配真 handler）。本件**没有**拿它当依据，而是拿 §2 的调用点与 §4 的红。该注释属 `internal/panel` 的下一枚写腿，本腿只登记、⛔ 不改。
5. **口径差异已具名**：实现件写 `0→2`（票面原尺，受理器本体），本件标题写 **N=4**（整条链的生产站点）。两个数都对，只是尺不同（§2"计数口径"）。⛔ 别拿 4 去说实现件报错了。
6. **行号会腐烂**：本件所有行号是 HEAD `d0639eed…` 的现量；引用前先现跑 §5 那把尺。

---

## 7. 门禁与〔未跑包〕

| 项 | 读数 |
|---|---|
| `go test ./cmd/wisp -run '^Test214AttachmentAddEndToEndStoresBytesInArtifacts$'`（无 overlay） | `PASS`，`ok github.com/CarlosShao/wisp/cmd/wisp 0.417s`，rc=0 |
| 同上 ＋ `-overlay` 变异 | `FAIL`（指名 `Test214AttachmentAddEndToEndStoresBytesInArtifacts`），rc=1 |
| `go test ./cmd/wisp -run '^Test214'` ＋ `-overlay` 变异 | 4 FAIL／4 PASS（§4.4），rc=1 |
| **〔未跑包〕`go test ./cmd/wisp`（全包）** | **未跑**。理由：`cmd/wisp` 现归别的写腿（452）在用，派单明令全包只读、⛔ 跑全包；本腿只跑 `-run` 过滤的**单枚／单组**，且跑前 `git status --porcelain -- cmd/wisp` 为空（§0③）、跑后 md5 同值 ⇒ 与在飞写腿无写冲突。**全包红绿不在本件射程**，要全包读数由持有该包写权的腿出。 |
| `go vet` / `gofumpt` / `d22scan` | **未跑**，理由同上＋本件**零 Go 改动**（写面只有 md + `.scratch/wisp/probes/**`），无对象可扫。 |
| 盘上残留 | 无：`cmd/wisp/panel_inbound.go` md5 前后同为 `53243FD2654812D927A10313EFAE9E16`，`git status --porcelain -- cmd/wisp` 空。 |

环境注记（复跑要带）：`PATH` 需含 `third_party/sherpa-onnx`（否则 `0xc0000135`＝假失败），`GOFLAGS=` 置空。

---

## 8. 载具清单（`.scratch/wisp/probes/214/ac1/`）

| 文件 | 是什么 |
|---|---|
| `mut-attachment-nil-panel_inbound.go` | 变异副本（盘上原件未动） |
| `overlay-attachment-nil.json` | `go build -overlay` 的替换表（绝对路径） |
| `green-baseline.log` | 正控：不摘接线的 PASS 全量日志 |
| `red-mutant-attachment-nil.log` | 反控：摘线后单枚 FAIL 日志 |
| `red-mutant-test214-group.log` | `-run '^Test214'` 变异组日志（4 红／4 绿） |
| `census.log` | 改前树 vs HEAD 的四把尺原始读数 |

---

## 9. 一句话给编排者

AC#1 的 0→N **做实了**：改前树（`25eb216…`）三把尺同时读到 0／`Attachment: nil`／受理器文件不存在；引入提交 `fc76d5d7e…`；今天生产调用点 **4 处**（`panel_inbound.go:343`／`:386`、`panel_attachment_214.go:204`／`:243`）；摘掉 `:386` 那一行 ⇒ **`Test214AttachmentAddEndToEndStoresBytesInArtifacts` 红在 `:88`**，还原即绿。按 A1875 两件套 ⇒ **AC#1 判满足**（勾由编排者翻）。
