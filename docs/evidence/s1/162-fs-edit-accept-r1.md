# 162-v1 验收件（**非实现者**）— 票 162 的 AC#1／AC#2／AC#3／AC#4 四格独立复算

- 派单＝`.scratch/wisp/dispatches/2026-09-27-125x-accept-162-r1r2-v1-four-cells.md`
- 写码位＝**无**（本程一格生产码都没改，禁改清单见 §0 末与 §8）；裁决位＝本程（162-v1），实现程＝162-r1／162-r2
- ⚠ `docs/evidence/s1/162-fs-edit-r1.md` 与 `-r2.md` 是**实现者自己的**证据件 ⇒ 本程把它们的话**全部当待复算断言**，
  两处例外逐条点名在 §9。本程自己跑的台件全部在 `.scratch/wisp/probes/162/v1/**`（含 5 枚变异与全部原始日志）。
- 档别约定（四档，与 r1/r2 同形）：**〔本轮现跑〕**／**〔现读码，未跑用例〕**／**〔盘上有件，我未复算〕**／**〔仅自述〕**。
  本表里**没有一格**用后两档当判据。

## 0. step-0 五件〔本轮现跑〕

```
date                                    → 2026-09-27 15:43:02 +0800
git rev-parse --abbrev-ref HEAD         → dev                （三次逐次跑，均为 dev；见下）
git rev-parse HEAD                      → 19513ccefa45d51206871f00b3efa6ab19dd620a
git status --porcelain -- internal/tools/ → （空，起手与交件前各核一次，见 §8）
票面 AC 四行连行号                       → 抄在下面
```

分支逐次：step-0 现量 `dev`；AC#1 取证后 `dev`；交件前 `dev`（三次同一枚命令，输出逐次贴，见 §8 末）。
⚠ 派单 §4 要求"分支名逐次"＝本程在**每一次提交前**各跑一次 `git rev-parse --abbrev-ref HEAD`，读数都在 §8。

票面＝`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`（**本程一字未动**）：

- `:32` `- [ ] **AC#1 先把"整文件写回会安静丢东西"量成读数**（未修码上：一发"写回时漏抄一段"的变异，量今天有没有任何一处报错）。量不到＝本票前提不成立，停手上报。`
- `:33` `- [ ] **AC#2 落地四枚拒绝**（上面 1–4 条各一枚用例），且**每一枚都要答"这一发在未修码上响不响"**。`
- `:34` `- [ ] **AC#3 换行/BOM 正反两向**（CRLF 文件改完还是 CRLF；带 BOM 的改完 BOM 还在；**不许把 LF 文件的行尾改掉**）。`
- `:35` `- [ ] **AC#4 原子性**：改到一半被取消不许留损坏文件（拿 `PLAN.md:1176` 那条当判据，并答"摘掉临时文件那一味，是否存在一发损坏从此看不见"）。`
- 判据形状 `:22`–`:28`（本表用到 1／2／3／4／5／6 六条）；`PLAN.md:1176` 行号**本程现量**：
  `grep -n "临时文件" docs/PLAN.md` → `1176:- 工程要求：`fs.write` 必须走**临时文件 + 原子 rename**，否则取消会留损坏文件；`
  ⇒ **票面引的那一行今天还在 1176**（r2 自己承认它没去现量、按 r1 转述读的）。

**本程写法约束**：四格全部走 `go test -overlay`（台件在仓外目录、虚拟路径落 `internal/tools/`），
`internal/**` 与任何生产文件**零字节**；每枚变异都做了"编译字节回读证明"（§5 的 `PROOF present-hits=1 absent-hits=0`），
防的正是派单点名的"overlay 路径没匹配上＝静默 no-op 假绿"。⚠ 全程零 `-cover`，**没有把 overlay 与 cover 合跑过**。

## 1. 本程**没**测什么（先写这节，一句不含三件）

1. **票面 `:36` 的 AC#3b 与 `:39` 的 AC#4b 两格本程未判**——它们是编排者 09-27 11:0x 新加、**尚未实现**，
   派单明令"一格都不许替它们判"。本程只做了两件相关的事：
   ① 复算了 AC#3b **前提**那两形今天静不静默（§4，结论＝**两形今天确实都静默 ⇒ 票面 :36/:37 的前提成立，不用改**）；
   ② 复算了 AC#4b 的**分界**（本表 AC#4 判的是 `Hooks.Kill` 这一支＝**优雅取消**这一支的原子性；
   真实 `taskkill` 不跑 Go 清理那一支**本程未验**，`fs.edit` 今天没有那一发的用例，归 AC#4b）。
2. **没判 AC#5（风险档路由／越界路径必须被 `PathResolver` 拒）**、**没判 AC#6（契约轴）**。
   §8 里那三行 `internal/risk/**` 的 numstat **不归 162**：出处＝`f576cf08`（票 160 第②格），现量在 §8。
3. **没复算 `fs.edit` 的 Unicode 修复层**（判据形状 1 第二半）——r1 明记"本程不做"，票面也没要求今天做。
4. **没测"真实模型会不会用对这枚工具"**：本程所有调用都是手写字节走桥，没有走 `C5 LlmProvider` 的 golden SSE。
5. **没做 `stageAndRename` 的其余失败分支**（`Sync` 失败／`Close` 失败／`os.Rename` 失败）——那是票 112/118 那族的射程，本程未复算。
6. **没测大小上限与行尾的交互**（`writeCap()` 两枚拒绝本程没为它们补用例，也没测"一个大文件里同时有 BOM 和 CRLF"）。
7. **AC#2 第一枚（非字面 old）没做单点回退**——本程的 4 枚回退只覆盖 判据 3／判据 4 与 AC#3／AC#4；
   撤掉 `count == 0` 那一支会掉到哪一发上，本程未量（登记在此，不当结论用）。
8. **没跑 `staticcheck`**（r1/r2 都登记它本机读不了 go1.27 的 export data），也不是派单门禁项。
9. **没跑全仓测试**（派单 §2 明令"逐包单跑，别跑全量"）；本表所有 `go test` 读数只覆盖 `./internal/tools/` 与 `tools/d22scan`。

## 2. 第①格＝AC#1〔**成立**〕——前提由验收位自己另造一发复算出来，且那把尺**会响**

派单要求："实现者那发不算，你自己另造一发漏抄一段"。本程自己的台件＝
`.scratch/wisp/probes/162/v1/probe162v1_ac13_test.go` 的 `TestV1AC1WholeFileRewriteDropsMiddleBlockUnreported`，
**与 r1 不同形**：13 行 / 375 字节，漏的是**中段四行**（r1 是 24 行 / 672 字节漏 6 行），
跑法（唯一读数来源，原始输出已入库 `logs/mine.txt`）：

```
go test -count=1 -v -overlay=.scratch/wisp/probes/162/v1/overlay-mine.json -run 'TestV1' ./internal/tools/   → rc=0
V1-AC#1(a) 落盘前 375 字节 / 13 行 → 落盘后 219 字节 / 9 行（少了 4 行 156 字节，中段四行）
V1-AC#1(b) IsError=false Truncated=false RiskLevel=L2 ErrorClass=""
           Text="已写入 …\v1notes.txt（219 字节，临时文件+原子重命名）"
V1-AC#1(c) row tool=fs.write risk=L2 decision=allow outcome=success error_class=""
V1-AC#1(c) memory.ToolCall 字段名=[ID TaskID Seq Tool ArgsJSON RiskLevel Decision DecidedAt
                                    StartedAt EndedAt Outcome ErrorClass CorrelationID GrantID]
V1-AC#1(c) agent.ToolOutcome 字段名=[Text IsError RiskLevel ErrorClass Truncated AppliedSteps] size-bearing=[]
V1-AC#1(c) Result 字段名=[Text IsError Truncated Origin AppliedSteps] size-bearing=[]
V1-AC#1(c) audit: "tools: call kind=success task=task-1 … tool=fs.write risk=L2 decision=allow
                   outcome=success rules_hit=[R1 R8] …" / "tools: PATH-ACCOUNT … roots=1 rewritten=[] unusable=[]"
V1-AC#1(c) 事后 fs.read：219 字节 Truncated=false ErrorClass=""
```

⇒ 票面 `:32` 那句话（"漏抄一段就是安静丢东西"）**由验收位自己复算成立**，且比 r1 的读数多一处：
**回执里唯一的数字是"发出去多少字节"（219），落盘前的 375 在 Result／台账／tool_call 行／audit 行／事后 fs.read 五处全不出现**
⇒ 读者手上没有被减数。

**那句反射字段名清单＝本程自己跑的**（派单点名"它是那条结论的唯一凭据"）。本程自己实现 `v1Fields`/`v1SizeBearing`
（不复用 r1 的 `structFieldsOf`），对六枚 token `size|bytes|line|diff|sha|hash` 逐枚比：

- `memory.ToolCall` 命中 **0 枚** ⇒ r1 那句"台账结构上装不下被减数"**复算一致**；
- 并且**多做一枚控制实验**证明这把尺不是死开关：同一枚函数喂一枚含 `SizeBytes/LineCount/DiffText/SHA` 的控制结构体
  ⇒ 命中 5 枚（`SizeBytes` 同时命中 size 与 bytes），喂一枚无名结构体 ⇒ 0 枚。读数行：
  `V1-AC#1 反射尺自证：控制结构体命中 5 枚 [SizeBytes/size SizeBytes/bytes LineCount/line DiffText/diff SHA/sha]（非死开关），无名结构体命中 0 枚`
- 本程**比 r1 多读两枚结构**（`agent.ToolOutcome`、`Result`）：size-bearing 同样 0 枚 ⇒ 这条结论比 r1 写的更宽。

**恒真那一问（AC#1）**：**响。** 本程做了 M-5 变异＝**给 fs.write 装上那枚缺的仪器**
（`probes/162/v1/mut/m5/fs_write.go`：改名前先 `os.Stat` 目标拿到改前尺寸，回执文案变成
"…（450 字节，比原文件少 222 字节，临时文件+原子重命名）"），跑 `-overlay=overlay-m5mine.json`：

```
--- FAIL: TestFSWriteSilentLossIsNotReported   fswrite_silentloss_ac1_test.go:159:
      (b) would be falsified: the answer mentions "少": "已写入 …（450 字节，比原文件少 222 字节 …）"
--- FAIL: TestV1AC1WholeFileRewriteDropsMiddleBlockUnreported   probe162v1_ac13_test.go:146:
      (b) would be falsified: the answer mentions "少": "已写入 …（219 字节，比原文件少 156 字节 …）"
```

⇒ 这一格**不是**"我看了没有所以没有"：仪器一旦出现，实现者那枚与本程这枚**同时变红**。
AC#1 的"没有信号"是**被断言钉住的读数**，不是叙述。本程另钉一条更硬的方向盘：
`if strings.Contains(out.Text, "375")` ⇒ 只要有人把**被减数**放进回执，本程这发也红。
**没有一支是恒真的**（本格的"未修码"就是今天这块码，M-5 是它的反向对照）。

## 3. 第②格＝AC#2〔**成立**〕——四枚拒绝各有一枚用例，且**两问都过**

### 3a 本程自己跑的正读数（未变异，`-overlay` 之外的真实工作树）

```
go test -count=1 -v -run 'TestFSEditRefuses|TestFSEditIsAllOrNothing|TestFSEditAppliesABatchOnDisk' ./internal/tools/  → rc=0
判据 1+2 读数：命中前 50 字节 / 命中失败后 50 字节（一字未变），文案含 "0 命中，文件 50 字节" +
              "must match exactly including all whitespace and newlines" + "先用 fs.read 读回原文"
判据 3 读数：  old 在 40 字节文件里命中 2 处，先数后拒，落盘前后一字未变
判据 4 读数：  两形（old:"" 与缺 old 字段）都拒，文件保持 29 字节
全有或全无读数：四批各 "本批 N 枚…一个字也不落盘，文件仍是 26 字节"（四腿 no_op/overlapping/second_misses/empty_batch 全 PASS）
对照读数：    26 字节 / 4 行 → 32 字节 / 5 行，台账四步含 "创建临时文件 .wisp-tmp-…" 与 "原子重命名 …（目标此刻起为新内容）"
```

⇒ 与 r1 §3 表里的四组数**逐枚一致**（50→50／40 字节 2 处→0／29→29／26→26／对照 26→32、4 行→5 行）。**复算一致，不是照抄。**

### 3b "这一发在未修码上响不响"——本程**实测**，不是叙述

本程造了 M-0＝**摘掉 `fs_write.go:775` 那一行注册**（`{Tool: fsEdit{d: d}, Decl: FSEditDecl(d)},`），
即回到"162 没落码时 fs.edit 不存在"的形状，再跑同样五枚用例（原始件 `logs/m0.log`）：

```
COUNT=1 / PROOF FSEditDecl-hits-in-mutated-copy=0
--- FAIL: TestFSEditRefusesANonLiteralOld                    文案="未知工具 fs.edit，可用工具见 list_tools"
--- FAIL: TestFSEditRefusesAnAmbiguousOldIsADistinctError    （同形）
--- FAIL: TestFSEditRefusesAnEmptyOldIsAWholeFileInsert      （同形）
--- FAIL: TestFSEditIsAllOrNothingAcrossABatch               （同形）
--- FAIL: TestFSEditAppliesABatchOnDisk                      （同形）
```

⇒ 逐枚答复：**五枚用例在"未修码"上全部不响**——它们**红**，而且红因是"未知工具"，
不是"这一段被拒了"。这条实测有两层用处：
① 它证明这四枚拒绝**不是靠"工具不存在"混过去的**（那样它们会**恒绿**，实测是**恒红**）；
② 它同时点名**"不响的那一支"**（派单要求单独点名）：这四枚意图今天真正走的通路是
**`fs.write` 整文件重抄**，那一支**今天确实不响**——由本程 §2 自己量到（375→219、五处信号全安静）。
⇒ **"未修码响不响"这一问的答案不能引用"未知工具"那一支当证据**（那本来就会响）；
能引用的只有 §2 那发。本表把这两支分开写，就是防这一形。

### 3c 判"算不算放水"只看两条——本程现量

- **断言有没有被动过？** 两向都查了：
  - `git diff -U1 5e83808..HEAD -- internal/tools/fs_test.go` ⇒ 名册三处计数 **5→6／6→7／5→6**（变大＝多一枚工具必须被登记），
    且两枚 `map[string]risk.Level` / `declared` 表**各加一行 `"fs.edit": risk.L2` / `"fs.edit": "L2"`**
    ⇒ 这两枚表是**逐名遍历**的（漏一名即红），**方向是变严**，没有一处删除断言或放宽阈值。
  - `git diff a9d05761..HEAD -- internal/tools/fs_edit_test.go`（＝r1 交件后 → 现在）⇒
    **删除行只有 1 行**：`for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {` 被换成用 `d` 的同形；
    其余全是新增的 builder 函数。**r1 那 5 枚用例的断言文本一字未动**（r2 §2 那句"一字未改"复算＝**成立**）。
- **helper 是不是改前就有的？** `git log --oneline -S <名字> --reverse -- internal/tools/` 逐枚现量：

  | helper | 最早出现的提交 | 判定 |
  |---|---|---|
  | `sealableTempDir124` | `dfa3dc44`（票 124/119 那批） | **改前就有** |
  | `mustCanonical` | `64c5fea3`（票 20） | **改前就有** |
  | `noStagingFilesLeft` | `94827e45`（票 20） | **改前就有** |
  | `readString` | `94827e45`（票 20） | **改前就有** |
  | `gateSpy` | `64c5fea3`（票 20） | **改前就有** |
  | `structFieldsOf`／`countLines` | `0cb8864b`（＝162-r1 自己） | 本票新建，**但只长在 162 自己的新测试文件里**，没替换任何既有断言 |

- **四枚拒绝的"响亮"形状本身有牙**：`assertRefusal` 同时要求 `IsError`、`ErrorClass=="tool"`（D37 让模型能自纠）、
  文案含补救动作、**`AppliedSteps` 必须为空**。⚠ 一票否决项：**判据 3 那枚反过来断言文案里不得出现 `0 命中`**
  （两枚错必须可分辨）——本程 M-1 变异（见 §6）把它打到红，说明这条可分辨性不是装饰。

⇒ **AC#2 判成立**：四枚各有一枚用例（判据 1／2／3／4 全覆盖，且 4 拆成"空 old 两形＋四批全有或全无"）、
对照组 `TestFSEditAppliesABatchOnDisk` 让"永远说不"无法满足本格、放水两问都过。

