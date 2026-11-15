# 214 — AC#2–AC#6 收件（非实现者／只读＋裁决；逐格两件套＋快照半边＋突变）

- 腿：`214-ac2-ac6-receipt`（**收件腿／非实现者**，只读＋裁决；写面只有本文件与 `.scratch/wisp/probes/214/**`）
- 裁定时刻盘上 HEAD：`cbf88954edb31e053ad9a3fe506560c169a8ac7d`（2026-11-15 15:13:59 +0800）
- 派单存档：`.scratch/wisp/probes/orch/20261115/txt/219-ac1-ac5-and-214-ac2-ac6-dispatch.txt` §「腿 B」（我照它执行）
- 被收的前一格：`docs/evidence/s1/214-ac1-receipt.md`（AC#1 已勾）；实现件：`docs/evidence/s1/214-attachment-wiring.md`；AC#5 取证件：`docs/evidence/s1/214-ac5-stamp-feasibility.md`
- 硬规遵守：**A1875**（①生产点 `file:line` ②会红测试名）／**A1879**（全仓相对路径＋贴命令原文与输出）／**A1880**（分清行为牙／字面钉）／**A1881**（时间条件格在锚点上另跑）
- 禁令遵守：⛔ 改产码／票面／台账／阈值／契约／`frontend/**`；⛔ 写 `cmd/wisp/**`；⛔ push。**零产码改动**（`git status --porcelain --untracked-files=no -- cmd/wisp internal/panel` 全程为空，见 §7）

---

## 0. 起手三核（现跑）

| 核 | 命令 | 读数 |
|---|---|---|
| ① HEAD 与时刻 | `git rev-parse HEAD`；`git log -1 --format='%H\|%ad\|%s' --date=iso` | `cbf88954edb31e053ad9a3fe506560c169a8ac7d`／`2026-11-15 15:13:59 +0800` |
| ② 票面 AC#2–#6 原文还在（未被人改） | `Select-String -Path .scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md -Pattern '^\- \['` | `:38`–`:42` 五格逐字在；`:37` AC#1 已 `[x]`；**本腿未翻任何框** |
| ③ 工作树干不干净（判"我有没有写产码"） | `git status --porcelain --untracked-files=no -- cmd/wisp internal/panel` | **空**（全程复取同值，见 §7 的 residue.log） |

环境注记（复跑要带）：`$env:PATH` 前置 `third_party/sherpa-onnx`、`$env:GOFLAGS=''`。
测试跑法（本腿实际用的三条，全为**单枚/单组**，⛔ 未跑整包）：
`go test ./cmd/wisp -run '^Test214' -count=1 -v`、`go test ./cmd/wisp -run '^Test(214|65)' -count=1 -v`、`go test ./internal/panel -run 'Attachment' -count=1 -v`。

**基线（不摘任何东西）**：`Test214*` **8/8 PASS**，`ok github.com/CarlosShao/wisp/cmd/wisp 0.881s`（载具 `.scratch/wisp/probes/214/ac2-6/green-baseline.log`）；`internal/panel -run Attachment` **ok 0.335s**（`green-baseline-internal-panel.log`）。

---

## 1. 判定汇总（本腿只裁 AC#2–#6，总 5）

| 格 | 判语 | 一句话依据 | 行为牙有无 |
|---|---|---|---|
| **AC#2** | **满足** | 超限拒绝有生产守卫＋行为牙（M-E 1 红）；"拒＋为什么进快照"三跳各有行为牙（M-A 4 红／M-B 2 红）；正控真走门（EndToEnd 8/8 基线绿） | 有（泵层）；**装配两站无牙**（M-F／M-G 各 0 红） |
| **AC#3** | **满足** | 按嗅探判、⛔ 不看扩展名；声明与嗅探不符即拒（M-C 2 红）＋ probe 10 例现测 | 有 |
| **AC#4** | **部分满足** | "越界那一形必红" ✓（M-D 1 红，且拒时零落盘；落点树另有行为断言）；"落盘位置在**授权根**内"按 C26 词义 **✗**（`[fs] allowed_dirs` 默认空，artifacts 根从不入册＝票 174 同一洞） | 有（名字守卫）；"根逃逸"形结构不可达、也无测试 |
| **AC#5** | **不满足** | 盖戳宿主整条缺：`OutgoingMessage` 产码构造 **0**、`ForAgent` 产码调用者 **0**、C25 名册无附件源；只有渲染形状的牙 | **无生产点**（A1875①缺）；②只是形状 |
| **AC#6** | **满足** | 三件逐条具名在票上（`:42`＋实现件头注 `:41-44`）；A1881 在锚点 `fc76d5d7e` 上三件在**产品页**都真；**但今天前两件已具备**（§6 具名） | 文件类判据，无测试牙（性质如此） |

**计数：满足 3／部分 1／不满足 1（总 5）**。
（AC#2 判满足的口径写在 §2「缺什么」：超限那一发只到 broker＋ledger，没走门；"同文件两向"这对不存在。）

---

## 2. AC#2 —— 超限拒＋快照可见＋正控

### 2.1 票面原文（逐字，`.scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md:38`）

> - [ ] **AC#2**：投一枚超限文件 ⇒ 拒、**且快照里能看到"被拒＋为什么"**；正控＝同文件放到限内 ⇒ 收。

### 2.2 五列

| 列 | 内容 |
|---|---|
| **AC 原文** | 见 2.1（三子句：超限⇒拒／快照见拒因／限内⇒收） |
| **判语** | **满足**（带一条具名口径差，见"缺什么"） |
| **依据（生产点 `file:line`）** | 上限常量：`internal/panel/attachments.go:43-47`（`MaxAttachmentBytes int64 = 64 << 20`）；门处传 0 ⇒ 用编译默认：`cmd/wisp/panel_attachment_214.go:204`（`panel.NewAttachmentBroker(store, memory.ValidArtifactName, 0)`）＋`internal/panel/attachments.go:186-188`（`max = MaxAttachmentBytes`）；**声明大小守卫**：`internal/panel/attachments.go:214-217`；**实读超限守卫**：`:226-235`（`LimitReader(rc, b.max+1)`＋`len(data) > b.max`）；**每次尝试都入账（含超限/解码失败/库打不开）**：`cmd/wisp/panel_attachment_214.go:218` `h.drafts.record(ref)`（三条早退路径在 `:237`／`:241`／`:243`）；账本写体：`cmd/wisp/panel_attachment_214.go:88-106`；账本读口：`:108-118`（rows）／`:120-130`（reason）；**快照 hop**：`internal/panel/pump.go:422-434` → 线格式字段 `internal/panel/composer.go:416`（`attachments`）／`:425`（`attachmentError`） |
| **依据（会红测试名＋牙种）** | ① **超限⇒拒**：`Test214OversizeAttachmentIsRefusedBeforeAnyRead`（`cmd/wisp/panel_attachment_214_test.go:191`）——摘 `:214-217` ⇒ **1 红**，红句 `the refusal does not name the ceiling: … 没有可读的内容来源`＋`the wire does not show the oversize refusal`（`red-E-no-ceiling.log`）＝**行为牙**。② **快照行有真写入者**：摘 `panel_attachment_214.go:218` ⇒ **4 红**（`Test214AttachmentAddEndToEndStoresBytesInArtifacts:114`、`Test214AttachmentRefusalReachesTheSnapshotWithItsReason:148`、`Test214DisguisedExecutableIsRefusedByItsBytes:184`、`Test214PathSpellingDisplayNameIsRefusedAndNothingLands:246`，`red-A-no-ledger-record.log`）＝**行为牙**。③ **拒因上线的泵跳**：摘 `internal/panel/pump.go:432-434` ⇒ **2 红**（`Test214AttachmentRefusalReachesTheSnapshotWithItsReason`＋`internal/panel` 的 `TestComposerAttachmentsTravelFromTheIngestDoorToTheWire`，红在 `internal/panel/composer_test.go:1495`；`red-B-*.log`）＝**行为牙**。④ 正控：`Test214AttachmentAddEndToEndStoresBytesInArtifacts`（`:76`）走真链、真 store，断言 `<dir>/artifacts` **恰好 1 枚**且名字 `^attachment-[0-9a-f]{16}\.png$`（`:97-100`），字节逐字节相同（`:106`） |
| **可核命令** | `git --no-pager grep -n -E 'MaxAttachmentBytes\|src\.SizeBytes > b\.max\|h\.drafts\.record' -- internal/panel/attachments.go cmd/wisp/panel_attachment_214.go`；`go test ./cmd/wisp -run '^Test214' -count=1 -v -overlay .scratch/wisp/probes/214/ac2-6/overlay-E.json`（其余见 §7 表） |
| **缺什么** | ⓐ **超限那一发不走门**：测试在 broker 上直接投 `AttachmentSource{SizeBytes: MaxAttachmentBytes+1, Open: nil}`，再**手写** `drafts.record(ref)`（`cmd/wisp/panel_attachment_214_test.go:206-218`）——三个生产函数都是真的，但"一枚超限 pick 经 `panel.attachment.add` 进来"这条组装从未跑过（64 MiB＋base64 ≈ 85 MB，无测试这么做）。ⓑ **"同文件两向"这对不存在**：超限发是**合成 size**，正控是另一枚 11 字节小 png——AC 字面的"同文件"没做到。ⓒ 见 §2.3 的装配无牙 |

### 2.3 ★AC#2 的「快照半边」做实（本腿的专项）

**问题拆三问：快照行有没有真写入者／谁是读口／有没有会红的测试。**

| 问 | 读数 | 证据 |
|---|---|---|
| **① 有真写入者吗** | **有，唯一一处**：`cmd/wisp/panel_attachment_214.go:218` `h.drafts.record(ref)`。它写在**成功与拒绝两条路之后**（`ingest` 的三种早退都先造好 `AttachmentRef`，见 `:231-244`），所以"被拒"也有行、也带人话原因（`ref.Reason`）。账本自身写体 `:88-106`：`items` 追行、`last` 只保留**最近一发**的拒因、成功后清空 | 摘它 ⇒ 4 红（M-A，§7）＝这不是装饰 |
| **② 谁是读口** | **两个生产读口，都通同一台泵**：<br>**(a) pipeline 形**：`cmd/wisp/run.go:1197-1198`（`Attachments: rt.attachmentDraftRows,`／`AttachmentReason: rt.attachmentDraftReason,`，**method value**，因为驻留 boot 在字面量之后才绑定账本）→ `internal/panel/pump.go:422-434` → `internal/panel/composer.go:416/:425`；账本绑定：`cmd/wisp/resident_windows.go:596` `src.run.bindAttachmentDrafts(attachmentDrafts)`（绑体 `cmd/wisp/panel_attachment_214.go:251-258`）<br>**(b) 无 pipeline 的驻留形（`-H=windowsgui` 双击，出厂形）**：`cmd/wisp/panel_resident_pump.go:101-102`（`Attachments: drafts.rows,`／`AttachmentReason: drafts.reason,`）→ 同一台泵；账本创建 `cmd/wisp/panel_resident_windows.go:306` `drafts := newAttachmentDraftLedger()`，装配 `cmd/wisp/resident_windows.go:634`<br>**CLI 缝**：`cmd/wisp/panel_inbound.go:343` 传 `drafts`（`wisp panel-inbound` 那条路拿 nil，只有审计没有快照行——实现件注释 `:340-342` 明写这是诚实缺省） | 两读口 grep 见 §9 census.log[1][2] |
| **③ 有没有会红的测试** | **泵层：有（行为牙，2 枚）**：`Test214AttachmentRefusalReachesTheSnapshotWithItsReason`（`cmd/wisp/panel_attachment_214_test.go:128`，真链→账本→`Snapshot()`＋`Marshal()` 双断言，`:157`／`:167`）与 `internal/panel` 的 `TestComposerAttachmentsTravelFromTheIngestDoorToTheWire`（`internal/panel/composer_test.go:1495`）——摘泵跳 M-B ⇒ 两枚同时红。<br>**装配两站：没有行为牙 ——★这是本格唯一的洞**：<br>· 摘 `cmd/wisp/run.go:1197-1198`（M-F）⇒ `-run '^Test(214\|65)'` **0 红**（`red-F-run-readers.log`：8/8 + 5/5 全 PASS，`ok`）。唯一相关的 `Test214PumpAssemblyCarriesTheAttachmentReaders`（`:284-312`）断言的是**盘上 run.go 的源码字符串**——它连 overlay 改的编译产物都看不见（overlay 不改盘），所以连"字面钉"都没响；**改盘**才会红，那是字符串审判（A1880 判"装饰性"）。<br>· 摘 `cmd/wisp/panel_resident_pump.go:101-102`（M-G）⇒ 同样 **0 红**（`red-G-resident-readers.log`）。resident 这条连字面钉都没有（`Test65ResidentShowTriggersResyncPins:237` 只钉构造函数的**实参串**，不钉两个 reader）。 |

**这一格的判语**：快照半边**是实的**——"写入者→账本→泵→线"每一跳都有生产点、都有一枚会因摘掉而红的行为测试；**但它的两个注射口（run.go 的 pipeline 形、panel_resident_pump.go 的出厂形）今天零行为牙**，删掉任何一处都没有一条行为测试会响。⇒ 能力活着（泵层证明），**接线可被无声删除**（A1880 那个形状），登记在 §10 交编排者。

---

## 3. AC#3 —— 伪装文件按字节判，⛔ 不按扩展名

### 3.1 票面原文（`:39`）

> - [ ] **AC#3**：投一枚伪装文件（扩展名说自己是图片、字节不是）⇒ 按嗅探结果判，**不按扩展名判**（这条尺要现读 `matchesISOBaseMedia` 的语义再写死）。

### 3.2 五列

| 列 | 内容 |
|---|---|
| **AC 原文** | 见 3.1（两个子句：按嗅探判／不按扩展名判） |
| **判语** | **满足**（并附一条**新测出的语义洞**：真 QuickTime 声明 `video/quicktime` 会被误拒——见 3.3，须另立票，本腿不改码） |
| **依据（生产点）** | 白名单表：`internal/panel/attachments.go:115-134`（6 枚：png/jpeg/gif/webp/mp4/quicktime）；嗅探循环：`:240-247`；不受支持拒绝（带 `describeHead`）：`:248-251`＋`:327-343`；**声明 vs 嗅探一致性**：`:252-255`＋`mimeAgrees` `:346-...`；判定结果写进 ref：`:256`；**产物名**是内容哈希＋**嗅探出的**扩展名：`:258-263`（`"attachment-" + shortHash(data) + "." + accepted.ext`）。**扩展名从不参与判定**：`src.DisplayName` 只过 `checkDisplayName`（`:207`，裸名守卫，交给 `memory.ValidArtifactName`）与 `sanitizeDisplayName`（`:310`，控制字符清洗），然后只用作显示名 `ref.Name` |
| **依据（会红测试名＋牙种）** | 摘嗅探、改成"信声明"（M-C：`internal/panel/attachments.go:248` 前插一段"声明非空即认"）⇒ **2 红**：`Test214DisguisedExecutableIsRefusedByItsBytes`（`cmd/wisp/panel_attachment_214_test.go:172`，红句 `an executable named .png was not refused: <nil>`）＋`Test214AttachmentRefusalReachesTheSnapshotWithItsReason:137`（`red-C-sniff-by-declared.log`）＝**行为牙**。另：同族的 92 家判据 `internal/panel/attachments_test.go:129` `TestRefusesUnsupportedAndMasqueradingInputsLoudly` 也在盘上（本腿未摘它做突变） |
| **可核命令** | `go test ./cmd/wisp -run '^Test214' -count=1 -v -overlay .scratch/wisp/probes/214/ac2-6/overlay-C.json`；`go run ./.scratch/wisp/probes/214/ac2-6/probe`（10 例语义表，日志 `sniff-semantics-probe.log`） |
| **缺什么** | ⓐ 见 3.3 的 `.mov` 误拒（新缺陷，须另立票）。ⓑ AC 说"这条尺要现读 `matchesISOBaseMedia` 的语义再写死"——本腿已现读并写死在 3.3，但**票面上还没有这段文字**（我的写面不含票面） |

### 3.3 `matchesISOBaseMedia` 的现读语义（真跑，不是转述）

源码 `internal/panel/attachments.go:136-149`：`len(h) >= 12` 且 `h[4:8] == "ftyp"` 且 12–15 字节的品牌 ∈ `{isom, avc1, mp42, "M4V ", "qt  ", mmp4}`。**不校验前 4 字节的 box size**（我给 8 字节短文件它也不 panic，落"未知类型"）。表序 `:132` mp4 在 `:133` quicktime 之前 ⇒ 同一谓词下**先命中 mp4**。

真跑 10 例（`.scratch/wisp/probes/214/ac2-6/sniff-semantics-probe.log`，走真 `panel.AttachmentBroker.Ingest`，sink 是丢弃式探针，⛔ 不落盘）：

| # | 输入 | 结果 |
|---|---|---|
| 01 | png 字节／`photo.png`／声明 `image/png` | 收，`mime=image/png`，`artifact=attachment-1c8927715e5b55c7.png` |
| 02 | `MZ…` 字节／`photo.png`／声明 `image/png` | **拒**：`"photo.png" 的类型（可执行文件（MS-DOS/PE 头 "MZ"））不受支持…` |
| 03 | png 字节／`photo.exe`／声明空（**扩展名反向骗**） | **收**，`mime=image/png`，产物名 `.png` ⇒ 扩展名**根本没被看** |
| 04 | png 字节／`photo.jpg`／声明 `image/jpeg` | 拒：`声明为 image/jpeg，但文件头表明它是 image/png` |
| 05 | jpeg 字节／`photo.png`／声明 `image/png` | 拒：`声明为 image/png，但文件头表明它是 image/jpeg` |
| 06 | ISO `qt  ` 品牌／`clip.mov`／声明 `video/quicktime` | **拒**：`声明为 video/quicktime，但文件头表明它是 video/mp4` ← ★误拒 |
| 07 | 同 06 字节／声明 `video/mp4` | 收，`mime=video/mp4`，`.mp4` |
| 08 | ISO `mmp4`／`clip.mp4`／声明空 | 收，`video/mp4` |
| 09 | ISO 未知品牌 `heic`／声明 `video/mp4` | 拒：`未知类型，文件头 0x0000001866747970` |
| 10 | ISO 头但只有 8 字节 | 拒（谓词要求 ≥12） |

★**第 06 例是真缺陷**（不是测试口径）：一枚**诚实的** QuickTime 文件、声明也是 `video/quicktime`，因为两个白名单类型共用一个谓词、表序先命中 mp4，被按"声明与文件头不符"拒掉。⇒ 嗅探对 ISO 容器**不可分辨**；要么按品牌细分（`qt  `⇒quicktime），要么给这两枚合并成一条。**本腿不改码**（产码禁令），具名登记，请编排者立票或并入 92 家。

---

## 4. AC#4 —— 落盘位置在授权根内／越界那一形必红

### 4.1 票面原文（`:40`）

> - [ ] **AC#4**：落盘位置在授权根内 ⇒ 越界那一形**必红**（票 174/175 那族判据同族）。

### 4.2 五列

| 列 | 内容 |
|---|---|
| **AC 原文** | 见 4.1（两半：落点在授权根内／越界形必红） |
| **判语** | **部分满足**。下半（越界形必红）**成立**且是行为牙；上半按 **C26 的"授权根"词义不成立**（`[fs] allowed_dirs` 默认空、artifacts 根从不入册——这是票 174 已记的同一洞），按 "artifacts 树＝它该落的那棵树" 词义成立。**两个词义都写下来，判定交编排者**（改票面文字也轮不到我） |
| **依据（生产点）** | 落盘**唯一路径**：`internal/panel/attachments.go:260-264`（名字＝内容哈希 → `b.sink.PutArtifact(ctx, name, data)`）；sink 是**真** `memory.Store`：`cmd/wisp/panel_attachment_214.go:199`（`memory.Open(h.dir)`）／`:204`；store 的 artifacts 根＝`filepath.Join(abs(dataRoot), "artifacts")`：`internal/memory/open.go:277`（常量 `:42`）；写入处再走一次名字守卫：`internal/memory/artifacts.go:71-77`（`validArtifactName`＋`winsec.PrivateFileExclusive(filepath.Join(s.artifactsDir, name), data)`）；**裸名守卫**：`internal/panel/attachments.go:207`（`checkDisplayName`→`internal/memory/ValidArtifactName`）。**整条路没有 `risk.PathResolver` 参与**（`internal/panel/attachments.go` 的 import 块里没有 `internal/risk`，也没有 `filepath` 决策） |
| **依据（会红测试名＋牙种）** | 摘裸名守卫（M-D：`internal/panel/attachments.go:207-209`）⇒ **1 红**：`Test214PathSpellingDisplayNameIsRefusedAndNothingLands`（`cmd/wisp/panel_attachment_214_test.go:234`，红句 `a path-spelling name was not refused: <nil>`；`red-D-no-name-guard.log`）＝**行为牙（拒的那一半）**；该测试还断言拒后 `<dir>/artifacts` 为空（`:248` 一带）。落点"树内"另有行为断言：`Test214AttachmentAddEndToEndStoresBytesInArtifacts:97-100` 要求 `<dir>/artifacts` **恰好 1 枚**、名字是内容哈希形——若 sink 换到别的树，这枚立刻红 |
| **可核命令** | `go test ./cmd/wisp -run '^Test214' -count=1 -v -overlay .scratch/wisp/probes/214/ac2-6/overlay-D.json`；`git --no-pager grep -n -E 'AllowedDirs\|NewPathCanonicalizer' -- cmd/wisp/run.go`；`git --no-pager grep -n 'AddAllowedDir(' -- '*.go'` |
| **缺什么** | ⓐ **"授权根"按 C26 词义今天不成立**：授权根只从配置来——`cmd/wisp/run.go:772-774` `allowed = append(allowed, cfg.FS.AllowedDirs...)` → `tools.NewPathCanonicalizer(allowed, …)`；生产里唯一的"加根"是用户手批 `cmd/wisp/approval_always.go:134`（`s.rt.mgr.AddAllowedDir(dir)`）。artifacts 根**从不入册** ⇒ 附件落点在授权根**之外**（票 174 的同一洞；214 的票面 §禁区/09-29 裁定明确不许动配置键，本腿也不动）。ⓑ **"根逃逸"形不可达也无测试**：产物名是内容哈希、display name 从不拼路径 ⇒ 今天唯一可能的越界是"名字拼路径"，已在 M-D 测到；真正的"根逃逸"没有、也不可能有测试（结构上不存在这条路径）。ⓒ 若编排者按 C26 词义要这一格过，唯一合法路径是**改 174 那条裁定**（配置/契约轴），⛔ 不是 214 能顺手做的 |

---

## 5. AC#5 —— C25 盖戳／不可信外部文本

### 5.1 票面原文（`:41`）

> - [ ] **AC#5**：来源戳（C25 现成名）盖上 ⇒ 附件内容进上下文时被当作**不可信外部文本**处理（票 200 那节"说明文件≠授权"同一条纪律）。

### 5.2 五列

| 列 | 内容 |
|---|---|
| **AC 原文** | 见 5.1（两个子句：C25 现成名盖戳／进上下文时按不可信外部文本处理） |
| **判语** | **不满足**（A1875 的"②有①没有"形状：渲染函数在盘上、有一枚形状测试，但**没有任何生产路径**把附件内容送进上下文，也没有为附件存在的 C25 源名） |
| **依据（生产点）** | **无。** 逐条现跑：<br>· `OutgoingMessage` 产码构造 **0**：`git --no-pager grep -n 'OutgoingMessage{' -- '*.go'`（剥 `_test.go`）⇒ **空**（结构体定义在 `internal/panel/composer.go:482-485`）<br>· `ForAgent` 产码调用者 **0**：`git --no-pager grep -n 'ForAgent' -- '*.go'`（剥测试/探针）⇒ 只有定义 `internal/panel/composer.go:495`（渲染体 `:495-520`）<br>· 唯一生产 `MessageRequestHandler` 实现是 **192 的预览引用门**：`cmd/wisp/panel_inbound.go:374-377` 装配 → `cmd/wisp/panel_preview_192.go:269`；它只认 `panel.PreviewPathPrefix` 的 path（`:272-277` 其余按名拒绝），**完全不碰 `req.AttachmentPayload` 或草稿行**，最后一行还明写"消息派发通路属票 35，本机未接入"（`:303-304`）<br>· **C25 名册没有附件源**：`internal/risk/provenance.go:84-94` 九枚（fs.read／search.content／clipboard.read／system.get／web.fetch／doc.read／screen.capture／asr.transcript／task.output），无 `panel.attachment` 一类 ⇒ "C25 **现成名**"里今天**不存在**可给附件盖的戳（这正是取证件 §2 末行的结论，本腿复跑同值） |
| **依据（会红测试名＋牙种）** | 只有**形状牙**：`Test214IngestRefRenderingCarriesTheArtifactPathAndRefusalsToTheAgent`（`cmd/wisp/panel_attachment_214_test.go:253`）直接构造 `panel.OutgoingMessage` 调 `ForAgent(artDir)` 断言文本形状——删掉 `ForAgent` 它会红（所以渲染器本身有牙），但它**不证明任何生产调用**。⇒ 对 AC#5 而言：**无生产点、无可红的行为测试**（A1875 两件套齐缺①） |
| **可核命令** | `git --no-pager grep -n 'ForAgent\|OutgoingMessage{' -- '*.go' \| Select-String -NotMatch '_test.go\|\.scratch'`；`git --no-pager grep -n 'Src[A-Za-z]* +=' -- internal/risk/provenance.go`；`git --no-pager grep -n 'HandleMessageRequest' -- '*.go'` |
| **缺什么** | 整条宿主：消息门 handler 本体（`ComposerDispatch.Message` 的**非预览**实现）＋盖戳一跳（scope 有值的那一跳）＋C25 名册新增附件源＋回执路径。落点唯一＝`cmd/wisp/panel_inbound.go:387` `Message: message`（今天只在 `preview != nil` 时非 nil） |

### 5.3 ★AC#5「指向新票」具名：**没有具名新票**

派单要求"具名是哪张新票、它现在什么状态（若已 `-done`，这一格要不要翻由你给判语）"。现读结论：

| 问 | 读数 | 尺 |
|---|---|---|
| 票面/证据指向哪张 | 只写"**须编排者立票**"，没有票号：`.scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md:54`（"工单池今天没有 owner，须编排者立票"）；`docs/evidence/s1/214-ac5-stamp-feasibility.md:43`（"须编排者立新票（或按规矩给 35 具名增格）"）；台账 `docs/reports/pending-and-issues.md:31704`（"需**立新票**（消息门 handler 本体）"） | 三处逐字 |
| 代码注释点名的票 35 是不是它 | **不是**：`.scratch/wisp/issues/35-panel-bridge-c17.md:3` `Status: ready-for-agent`（**未 -done**），八枚 AC 全传输层（白名单 fuzz／correlationId 路由／伪造 allow 拒／resync 无状态／背压／secret 扫描／传输协议），**无一格**认领 handler 本体 | 票面 §AC + 取证件 §3 |
| 工单池有没有新票认领 | **没有**：`.scratch/wisp/issues/` 452 张逐名扫过，`MessageRequestHandler`／`panel.message.send`／`消息门`／`附件`＋`Provenance` 交叉 grep 后，无一张新票认领"消息门 handler 本体"（命令与输出见 `.scratch/wisp/probes/214/ac2-6/census.log` §[7][8]） | `Get-ChildItem .scratch/wisp/issues -Filter '*.md' \| Select-String -Pattern 'MessageRequestHandler\|消息门' -List` |
| 它现在什么状态 | **不存在**（既非 open 也非 -done）⇒ 派单里"若已 `-done`"那一支**不适用**，⛔ 我没有默默跳过：我在此明确写"指向的是一张**尚未立**的票" | 同上 |

**这一格要不要翻：不翻（AC#5 判不满足）。** 判语理由（给编排者）：
1. A1875 的规矩正是为这个形状写的——渲染器在、测试在、**生产调用者 0** ＝ "②有①没有"。拿"渲染形状已交付"翻 AC#5，等于把"内容真的进了上下文并被当不可信文本处理"换成"有个函数能渲染它"。
2. 若要收 214，唯一诚实做法是**先立那张票**（票名 ≤100 字符），AC 里写死三件：①`ComposerDispatch.Message` 的非预览实现；②盖戳点＝`ForAgent` 渲染后、喂 agent 之前，对每枚 `Stored` 附件 `c25.Mark(scope, 新源名, 内容, origin)`；③`internal/risk/provenance.go` 名册增一枚附件源并进聚合表（这是 C25 的**代码名册**扩张，按纪律要具名登记）。然后 214 的 AC#5 才能以"移交该票"结线——**不是以"已具备"结线**。

---

## 6. AC#6 —— 三件"没接之前不许说已具备"，逐条具名

### 6.1 票面原文（`:42`）

> - [ ] **AC#6**：三件"没接之前不许说已具备"——界面点不了、没有选择器、没有 @引用，**逐条具名留在票上**。

### 6.2 五列

| 列 | 内容 |
|---|---|
| **AC 原文** | 见 6.1（这是一格**文件类判据**：三件必须逐条具名在票上） |
| **判语** | **满足**（字面：三件逐条具名在票面 `:42`，实现件头注 `cmd/wisp/panel_attachment_214.go:41-44` 又逐条写了一遍）；**A1881**：在锚点 `fc76d5d7e` 上三件在**产品页**上都真；**★但今天前两件已经具备**（下表的现读），⛔ 不能拿它当"今天界面点不了"的凭据 |
| **依据（生产点）** | 具名处：`.scratch/wisp/issues/214-composer-plus-menu-attachments-and-context-wired.md:42`（票面）＋`cmd/wisp/panel_attachment_214.go:41-44`（"no page-side picker, no @reference parser and no message-send door exist yet"）。**无生产调用点**（这格要的就是"没有"） |
| **依据（会红测试名＋牙种）** | **无测试牙**——文件类判据的性质；本腿不以"有测试"来判它，而是以"字在不在＋当时真不真"判 |
| **可核命令** | 见 6.3 的 anchor 日志 `.scratch/wisp/probes/214/ac2-6/anchor-2026-11-07.log` 与 §9 的 ten 条 grep |
| **缺什么** | 票面措辞已过期两件（下节）。⛔ 我没改票面（写面禁令）；改判据文字须编排者/人工，不归我 |

### 6.3 A1881：在锚点上复跑（时间条件格）

锚点＝本票交付提交 `fc76d5d7e96664790d4d21904ce72fda6cf18883`（2026-11-07 17:52:15 +0800，`214: 附件接线…`）。

| 三件 | 锚点 `fc76d5d7e` 读数 | 今天 HEAD `cbf88954` 读数 |
|---|---|---|
| **界面点不了** | **真**：产品页 `frontend/src/App.tsx` 挂 `PromptBar`，四个回调是桩——anchor `frontend/src/App.tsx:332` `onPermissionChange={() => undefined}`、`:333` `onSend={() => undefined}`（wave128 后来才接的线） | **已不真**：`frontend/src/App.tsx:496` 挂 `ComposerRow`；`frontend/src/components/composer-row.tsx:291` `onAttach={composer ? () => fileInput.current?.click() : undefined}`；`:216` `await submitAttachment(f)` → `frontend/src/lib/panel.ts:718` → `:709` `sendRequest("panel.attachment.add", …)` |
| **没有选择器** | **真（产品页）**：当时唯一带 `type="file"` 的是 `frontend/src/components/composer.tsx:210`，而 App.tsx **不 import 它**（只被 `frontend/scripts/render-composer.tsx:31` 这个渲染脚本引）⇒ 产品页上没有挂载的选择器 | **已不真**：`frontend/src/components/composer-row.tsx:299-306` `<input type="file" multiple accept={…acceptedAttachmentMimes…}>` 挂在产品页上；构建产物 `frontend/dist/assets/index-QUtCIbHm.js`（2026-11-11 16:59，640318 B）里能搜到 `panel.attachment.add` |
| **没有 @引用** | **真** | **仍真**：`frontend/src/components/ai-native/prompt-bar.tsx:101-106`／`:213` 有 @ 菜单机制，但 `ComposerRow` **不传 `sources`**（`git --no-pager grep -n 'sources=' -- frontend/src/components/composer-row.tsx` ⇒ **空**）；Go 侧四个入向方法（`internal/panel/bridge.go:42-45`：mode.request／workspace.request／attachment.add／message.send）里**没有**"引用→受理器"的路 ⇒ 工作区文件/目录的 @ 引用仍无 |

**判语展开**：AC#6 当时（11-07）**是真话**，所以这格按 A1881 判**满足**（锚点上成立）。但**今天**票面这句已经落后于前端那条腿（11-11 的 `6c0cbd244` 把四个 composer 回调接上、并把 `composer-row.tsx` 挂进产品页）——我的读数是"前两件已具备、只第三件仍无"。⇒ 建议编排者把票面那句改成"11-15 起：选择器与附件按钮已上线（前端腿），@引用仍无"，⛔ **别**继续拿"界面点不了"当今天的现状（会误判成"这票还有前端半没做"）。

---

## 7. 突变（七发，全部 `-overlay`，盘上原件一字节未动）

载体目录：`.scratch/wisp/probes/214/ac2-6/`（副本＋overlay JSON＋红日志）。**每一发的 diff 只有一处**（`mutation-diffs.log`）；`-overlay` 只改编译输入，不改盘。

| # | 摘掉/改掉什么（盘上坐标） | 跑法 | 红了几枚／红在哪 | 牙种 |
|---|---|---|---|---|
| **M-A** | `cmd/wisp/panel_attachment_214.go:218` `h.drafts.record(ref)` → 注释＋`_ = ref`（快照行的**唯一写入者**） | `go test ./cmd/wisp -run '^Test214' -overlay overlay-A.json` | **4 红**：EndToEnd`:114`、RefusalReachesTheSnapshot`:148`、Disguised`:184`、PathSpelling`:246` | **行为牙** |
| **M-B** | `internal/panel/pump.go:432-434`（拒因上线的泵跳） | 同上＋`go test ./internal/panel -run 'Attachment'` | **2 红**：`cmd/wisp` 的 RefusalReachesTheSnapshot；`internal/panel` 的 `TestComposerAttachmentsTravelFromTheIngestDoorToTheWire`（`composer_test.go:1495`） | **行为牙** |
| **M-E** | `internal/panel/attachments.go:214-217`（声明大小上限守卫） | `-run '^Test214'` | **1 红**：Oversize`:211`＋`:227` | **行为牙** |
| **M-C** | `internal/panel/attachments.go:248` 前插"声明非空即认"（**信扩展名/声明、不信字节**） | `-run '^Test214'` | **2 红**：Disguised`:181`、RefusalReachesTheSnapshot`:137` | **行为牙** |
| **M-D** | `internal/panel/attachments.go:207-209`（裸名守卫） | `-run '^Test214'` | **1 红**：PathSpelling`:243` | **行为牙** |
| **M-F** | `cmd/wisp/run.go:1197-1198`（pipeline 装配的两枚 reader） | `-run '^Test(214\|65)'` | **0 红**（8/8＋5/5 全 PASS） | **无牙**（只剩"改盘才红"的源码字面钉；overlay 下连它都不响） |
| **M-G** | `cmd/wisp/panel_resident_pump.go:101-102`（驻留/出厂形的两枚 reader） | `-run '^Test(214\|65)'` | **0 红** | **无牙**（resident 连字面钉都没有） |

**留痕（盘上读数，`residue.log`）**：

```
$ git status --porcelain --untracked-files=no -- cmd/wisp internal/panel
（空）
$ git diff --stat -- cmd/wisp internal/panel
（空）
$ Get-FileHash -Algorithm MD5 cmd/wisp/panel_attachment_214.go, internal/panel/pump.go,
    internal/panel/attachments.go, cmd/wisp/run.go, cmd/wisp/panel_resident_pump.go
6521BAFB709E636F9C8DB70EAE99BBB5  …\cmd\wisp\panel_attachment_214.go
35BB0F6E16480908BF815F359A106AD0  …\internal\panel\pump.go
F235522DE855917BDF9E76B3D31F6B98  …\internal\panel\attachments.go
B7D7EBC0534528EFAF3BA1B3B287B327  …\cmd\wisp\run.go
01DB09B43D7BF5551C422C5513F4B700  …\cmd\wisp\panel_resident_pump.go
HEAD = cbf88954edb31e053ad9a3fe506560c169a8ac7d
```

**七发的合并读数**：AC#2 的三跳（上限／入账／理由 hop）、AC#3 的嗅探、AC#4 的裸名守卫 —— **全部是行为牙**；**唯一零牙的是两处装配读口**（M-F／M-G），这正是 §2.3 那个洞的实测。

---

## 8. 〔未跑包〕与门禁（如实）

- **`go test ./cmd/wisp` 全包未跑**（该包归写腿，且本机整包有已知的 native 死法 `0xC0000005`，台账 `A1882`）⇒ 本件所有 `cmd/wisp` 读数只来自**单枚/单组**：`-run '^Test214'`、`-run '^Test(214|65)'`。
- **`go test ./internal/panel` 全包未跑**：只跑了 `-run 'Attachment'`（基线 ok 0.335s）与 M-B 那一发（红在具名那一枚；`red-B-internal-panel.log`）。
- **跑前一发 `git status --porcelain --untracked-files=no -- cmd/wisp` 空**（§0 三核③；§7 residue.log 复取同值）。
- **`go vet`／`gofumpt`／`d22scan` 未跑：本腿零 Go 改动，无对象可扫**（探针 `go run` 只编译 `.scratch` 下的探针 main，不动产码）。
- 环境：`$env:PATH` 前置 `third_party/sherpa-onnx`（否则 `0xc0000135` 假失败，A1877 规矩一）；`$env:GOFLAGS=''`。
- 写面自检：本件＋`.scratch/wisp/probes/214/ac2-6/**` 之外**零字节**；`frontend/**`、`design/**`、`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt` 全未动；⛔ 未 push。

---

## 9. 载具清单（`.scratch/wisp/probes/214/ac2-6/`）

| 文件 | 是什么 |
|---|---|
| `green-baseline.log`／`green-baseline-internal-panel.log` | 正控：8/8 PASS＋internal/panel ok |
| `mut-A…go`／`mut-B…go`／`mut-C…go`／`mut-D…go`／`mut-E…go`／`mut-F…go`／`mut-G…go` | 七个变异副本（盘上原件未动） |
| `overlay-A.json` … `overlay-G.json` | 七张 `go build -overlay` 替换表（绝对路径） |
| `red-A-no-ledger-record.log`／`red-B-no-attachment-reason.log`／`red-B-internal-panel.log`／`red-C-sniff-by-declared.log`／`red-D-no-name-guard.log`／`red-E-no-ceiling.log`／`red-F-run-readers.log`／`red-G-resident-readers.log` | 每发的逐枚红绿原始输出 |
| `mutation-diffs.log` | 七发 diff（每发只有一处改动，自证不是"顺手多改"） |
| `probe/main.go`＋`sniff-semantics-probe.log` | AC#3 的 10 例语义表（走真 `AttachmentBroker`，丢弃式 sink，不落盘） |
| `census.log` | §[1]–[10] 十组 grep 的命令原文＋输出（生产点普查、账本绑定、快照字段、嗅探、守卫、授权根、AC#5 三查、前端两件） |
| `anchor-2026-11-07.log` | A1881：AC#6 在锚点 `fc76d5d7e` 上的复跑读数（含 `composer-row.tsx` 在锚点不存在的 `git cat-file` 读数） |
| `residue.log` | 七发之后的 `git status`／`git diff --stat`／五枚 md5／HEAD |

---

## 10. 一句话给编排者（含要你拍的）

**AC#2–#6 的盘上读数**：AC#2 满足（泵层行为牙齐；**装配两站零牙**）、AC#3 满足（**新发现**：诚实 `.mov` 会被误拒）、AC#4 部分满足（**"授权根"按 C26 读不成立，是 174 的同一洞**）、AC#5 不满足（盖戳宿主整条缺，指向的是一张**尚未立**的票——票 35 不是它、也未 -done）、AC#6 满足（锚点上真；**今天前两件已上线，票面措辞过期**）。

**要你拍的两件（我不改票面/不再往下做）**：
1. **AC#4 的词义**：按 "artifacts 树＝该落的那棵树" 读就翻；按 C26 "授权根" 读，就得先动 174 的裁定 —— 这是配置/契约轴，⛔ 不是 214 能顺手补的。
2. **AC#5 的结线方式**：要么先立"消息门 handler 本体＋C25 附件源"那张票（我可以给判据草案），要么明确 214 以"移交该票"结线。**⛔ 别拿"渲染形状已有"翻框**（A1875 的原形就是这句）。
