# 162-r1 证据件 — 票 162 的两格：AC#1（读数）＋ AC#2（四枚拒绝）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-094x-impl-162-r1-fs-edit-refusals.md`
- 进场锚＝`5e83808`（step-0 现量）｜交件时刻＝09-27 10:2x｜写码位＝本程（162-r1）
- 档别标注约定：**〔本轮现跑〕**＝这一程亲自跑过的命令与输出；**〔盘上有件，我未复算〕**＝读了码/读了件，没重跑；**〔仅自述〕**＝只有我一句话，无凭据。

## 0. step-0 四件〔本轮现跑〕

```
date                       → 2026-09-27 09:47:09 +0800
git rev-parse HEAD         → 5e838081826e4587d27ad579b34db0d8eec92555
git status --porcelain -- internal/tools/ docs/PLAN.md docs/specs/   → （空）
```

票面 AC 两行连行号（`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`）：

- `:32` `- [ ] **AC#1 先把"整文件写回会安静丢东西"量成读数**（未修码上：一发"写回时漏抄一段"的变异，量今天有没有任何一处报错）。量不到＝本票前提不成立，停手上报。`
- `:33` `- [ ] **AC#2 落地四枚拒绝**（上面 1–4 条各一枚用例），且**每一枚都要答"这一发在未修码上响不响"**。`

判据形状 6 条在 `:22`–`:28`（本程用掉 1／2／3／4 四条，第 5 条＝AC#3 留 r2，第 6 条＝红线，本程守住）。

## 1. 本程**没**测什么（先写这节）

一句不许含三件，逐条列清"没测＋为什么是刻意的"：

1. **没测 CRLF／BOM 任何一形**（AC#3）。派单 §3 明令"Unicode 修复那一层本单不做"、§0 明令 AC#3 归 r2。本程只做字节级字面匹配，**行尾还原、剥 BOM、首次出现判定行尾**一次都没跑过。⇒ r2 别把"本程的 `fs.edit` 在 LF 文件上对"当成"它在 CRLF 上对"。
2. **没测原子性的反面**（AC#4：摘掉 temp+rename 那一味会不会有一发损坏从此看不见）。本程**复用了 `fs_write.go` 里已有的 `stageAndRename`**（红线 6），既没新增写路径也没剥掉它，所以 AC#4 那条"摘掉一味"的实验根本不属于本程的形状。
3. **没测风险档路由与越界路径**（AC#5）。本程只把 `Declared` 写成 **L2**（D34 冻结文本那一行），**没有**给 `fs.edit` 挂 `Facts` 钩子、**没有**验证"越界路径必须被 `PathResolver` 拒"、**没有**跑过 C19 对 `fs.edit` 的完整判级（R2/R3/R8 命中组合）。用例里那条 `dec.Level != risk.L2` 只证"审批通道确实按声明档走了一次"，不证 AC#5。
4. **没做契约轴复算**（AC#6）。派单 §5 只要求跑四把尺。`docs/PLAN.md`／`docs/specs/**` 本程零字节（见 §6）。
5. **没测 `internal/risk/provenance.go` 的 exfil 通道对 `fs.edit` 的实际行为**：我只**读码**确认它是按调用形状（`writeGate`，`provenance.go:655-692`）判的、不按工具名硬编码，所以 `fs.edit` 不在 `channelTable`（`:127`）里**不构成缺口**；〔盘上有件，我未复算〕——没有为本程的这句话跑过一发 taint 用例。这句话该由 AC#5/r2 用例钉。
6. **没测"真实模型会不会用对这枚工具"**：`old` 由模型逐字抄，抄错就是 0 命中硬错。本程只证"抄错会响"，不证"抄错率低"。这是产品效果轴，不在票面射程。
7. **`fs.edit` 的 Unicode 修复层（判据形状 1 第二半）没做**，且**不是遗漏**：见 §4 末那一格的理由。

## 2. 第①格＝AC#1 的三件事（全部〔本轮现跑〕）

跑法（唯一读数来源，原始输出已入库）：

```
bash .scratch/wisp/probes/162/run-ac1.sh
  → .scratch/wisp/probes/162/ac1-reading.txt（17 行，本程未编辑）
用例：internal/tools/fswrite_silentloss_ac1_test.go:78 TestFSWriteSilentLossIsNotReported
形状：走 internal/tools 的桥（fsBridgeWith 同一套组成：真 C26＋真 C19＋真 gate＋真 SQLite journal＋审计行捕获），
      发一发 fs.write，content＝原文件少抄一段。fs.write 一字未改（本格在 AC#2 落码之前先提交，0cb8864）。
```

### (a) 文件真的少了那一段

```
AC#1(a) 落盘前 672 字节 / 24 行；落盘后 450 字节 / 18 行；少了 6 行 222 字节
```

用例断言：六行 `beta-*` 逐行 `strings.Contains(after, l)` 必须为 false（少了），且 `after == mutated`（少得和发出去的一字不差）。**⇒ 前提成立：这六行确实没了。**

### (b) 返回里没有任何报错／警告／计数

```
AC#1(b) IsError=false Truncated=false RiskLevel=L2 ErrorClass=""
AC#1(b) Result.Text="已写入 C:\...\notes.txt（450 字节，临时文件+原子重命名）"
AC#1(b) AppliedSteps[0]="在 C:\...\001 创建临时文件 .wisp-tmp-e3fd961f-61988-4182457482"
AC#1(b) AppliedSteps[1]="向临时文件写入 450 字节"
AC#1(b) AppliedSteps[2]="临时文件已落盘并关闭"
AC#1(b) AppliedSteps[3]="原子重命名 .wisp-tmp-… → …\notes.txt（目标此刻起为新内容）"
```

断言两向：`IsError==false && ErrorClass==""`（没有失败）**且** 文案与台账里不得出现
`少|丢|缺|删除|行差|警告|warn|missing|dropped` 任一（一旦有人补了信号，这条断言会红）。
唯一的数字 **450** 是"发出去多少字节"，不是"少了多少"——它和落盘前那 672 无任何关系。

### (c) 事后没有任何仪器能看出"少了 N 行"

```
AC#1(c) tool_call 行：tool=fs.write risk=L2 decision=allow outcome=success error_class=""
AC#1(c) tool_call 行的全部列名：[ID TaskID Seq Tool ArgsJSON RiskLevel Decision DecidedAt
                                 StartedAt EndedAt Outcome ErrorClass CorrelationID GrantID]
AC#1(c) audit: "tools: call kind=success task=task-1 corr=corr-1 tool=fs.write risk=L2 decision=allow
                outcome=success rules_hit=[R1 R8] in_allowlist_scope=true
                reason=\"R1: 工具声明为下界（L1）; R8: 不可逆操作（覆盖已有内容）\""
AC#1(c) audit: "tools: PATH-ACCOUNT task=task-1 tool=fs.write roots=1 rewritten=[] unusable=[]"
AC#1(c) fs.read 事后复读：450 字节，Truncated=false
```

三条命令化的判据，不是"我看没有"：

- **tool_call 表上没有能装下这个数的那一列**：用例反射 `memory.ToolCall` 的字段名，
  含 `size|bytes|line|diff|sha|hash` 者必须为 0 枚 ⇒ 记录**结构上**不可能有差值。
- **审计行里没有落盘前的尺寸 672**：用例在捕获的全部 audit 行里找字符串 `672`，必须找不到
  ⇒ 读者**拿不到被减数**，也就算不出差。
- **事后 `fs.read` 只给内容**：450 字节的原文，`Truncated=false`，没有行数/大小字段
  ⇒ 唯一的对照物是"你自己另存过旧字节"。

**AC#1 结论（票面前提成立）**：丢 6 行 222 字节，五处信号（Result.Text / AppliedSteps /
tool_call 行 / audit 行 / 事后 fs.read）全为 0。今天的通路确实是"我们在赌模型抄得全"。

## 3. 第②格＝AC#2 四枚拒绝（〔本轮现跑〕，`internal/tools/fs_edit.go`＋`fs_edit_test.go`）

| 枚 | 判据（票面 `:22-26`） | 用例名 | 未修码上响不响 | 落盘前后字节证据 |
|---|---|---|---|---|
| 1 | `:22` 定位＝字面子串＋`:24` 匹不上硬错带补救 | `TestFSEditRefusesANonLiteralOld` | **不响**：fs.edit 不存在 ⇒ "缩进记错了"这一发今天走不到定位；今天唯一通路是整文件重抄，缩进错了也算成功 | `命中前 50 字节 / 命中失败后 50 字节（一字未变）`；断言＝`readString` 与原文全等＋目录内无 `.wisp-tmp-*` 残留 |
| 2 | `:24` 文案逐字带补救动作 | 同上（`assertRefusal` 查 `"must match exactly including all whitespace and newlines"` 与 `"fs.read"` 两段补救文字） | 同上 | 同上 |
| 3 | `:25` 多处命中＝**另一种**硬错，先数再拒，不许取第一处 | `TestFSEditRefusesAnAmbiguousOldIsADistinctError` | **不响**：今天的 fs.write 没有"命中几处"的概念；歧义在旧通路上表现为"改了第一处、第二处原样留着"，而且成功 | `old 在 40 字节文件里命中 2 处，先数后拒，落盘前后一字未变`；额外断言`strings.Count(after,"item: apricot")==0`（没吃掉第一处）＋`not 0 命中`（两枚错必须可分辨） |
| 4 | `:26` 空 old 拒；无变化／重叠整批拒＝全有或全无 | `TestFSEditRefusesAnEmptyOldIsAWholeFileInsert`（空 old 两形：`old:""` 与缺 `old` 字段）＋`TestFSEditIsAllOrNothingAcrossABatch`（四批：no-op／重叠／第二批找不到／空批） | **不响**，且方向相反：今天要做"往文件里插一段"只能整文件重抄交回来（fs.write 接受），插没插对没有读数 | `两形都拒，文件保持 29 字节`；四批各 `本批 N 枚…一个字也不落盘，文件仍是 26 字节`（good edit 也没落＝全有或全无） |
| — | 红线 `:28`（第 6 条不许照抄上游） | `TestFSEditAppliesABatchOnDisk` | 对照格（不然四枚拒绝可被"永远说不"满足） | `26 字节 / 4 行 → 32 字节 / 5 行（临时文件+原子重命名）`；台账四步含 `创建临时文件 .wisp-tmp-8c8b7c13-…`／`原子重命名`；用例断言 AppliedSteps 必含"原子重命名"、且目录内无残留 |

四枚拒绝共同的"响亮"形状（`assertRefusal`）：`IsError==true`、`ErrorClass=="tool"`（D37 让模型能自纠）、
文案含补救动作、`AppliedSteps` 必须为**空**（一个副作用都没发生）。

读到的两条实际文案（证明补救是**逐字在错里**，不是我转述）：

```
fs.edit 拒绝执行：第 1 枚编辑的 old 在文件里找不到（0 命中，文件 50 字节）。old 必须与文件内容逐字一致，
包含全部空白与换行（the old string must match exactly including all whitespace and newlines）；
先用 fs.read 读回原文再照抄那一段。本次调用一个字也没有写，目标文件保持原状。

fs.edit 拒绝执行：第 1 枚编辑的 old 在文件里命中 2 处（多于 1 处即为歧义，不取第一处）。
请在 old 里带上更多上下文，让它只命中一处（provide more surrounding context to make it unique）。
本次调用一个字也没有写，目标文件保持原状。
```

### 为什么"多处命中"和"匹不上"是两枚不同的错

票面 `:25` 要求"另一种硬错"。用例里两枚的断言词互斥（`0 命中` vs `命中 N 处`），
且互相验否（歧义那枚断言文案里**不得**出现 `0 命中`）。模型对这两件事的修法相反：
前者要抄得更准，后者要抄得更长。合并成一枚就是没读判据。

### Unicode 修复层（判据形状 1 第二半）本程**没做**，理由不是省事

Pi 在精确匹配失败后再走一层 NFKC／逐行去尾空白／弯引号→ASCII 的修复。本程不做的实质理由：
**一枚"要修复才能命中"的 old 正是要避免的东西**——修复层每宽一分，就把一次响亮失败变成一次
静默改写；本票存在的全部理由就是消灭"静默改写"。要做就得单独定案（哪些形可修、修完要不要在
台账里报"我是修复命中的"），不该塞进一枚"落地四枚拒绝"的程里。已写进 `fs_edit.go` 的
文件头注释，r2／后续若要做得单独一格。

## 4. 注册链与 L2：编排者留的两问，实测答案（〔本轮现跑〕）

- **问①（`fs.edit` 要不要进"配置允许面"的名单？）＝只进 `BuiltinFSWriteEntries`，不用动任何别处。**
  现读注册链只有一条：`cmd/wisp/run.go:341` `for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: rt.paths, DeleteEnabled: cfg.FS.DeleteEnabled})`
  → `fs.go:323 BuiltinFSEntries` → `fs_write.go:772 BuiltinFSWriteEntries`。
  本程把 `{Tool: fsEdit{d: d}, Decl: FSEditDecl(d)}` 插进后者一枚，模型侧即可见；
  **`internal/config/**`、`internal/risk/**` 零字节**（§6）。`Registry.Register` 的 C1 名字校验
  （`registry.go:247 validToolName`）现读通过（`fs.edit` 是合法 `namespace.action`）。
- **问②（"覆盖已存在如何走到 L2"的机理）＝它不按工具名走，所以本程不需要表态。**
  现读：`fs_write.go:638 writeFacts` 把 `risk.Facts.OverwriteExisting` 喂给 C19，`R8` 升 L2；
  这条链对 `fs.edit` 的必要性是 0——**本工具唯一能做的事就是覆盖一个已存在的文件**
  （目标不存在时它先拒，见 `fs_edit.go:127`），所以 D34 那一行直接把声明档写成 **L2**，
  不挂 `Facts`。⇒ **AC#5 的路由断言本程一条没写，`internal/risk` 本程一字未动，不需要停手上报。**
  顺带现读：`internal/agent/loop.go:773 decideRisk` 对**声明 L1/L2** 的条目在 `AdmitTask==nil` 时拒；
  今天 `fs.write/fs.trash/fs.move`（L1）与 `fs.delete`（L2）已在同一条款下，`fs.edit` 不比它们更受限。

## 5. 门禁（〔本轮现跑〕，两格全部落盘之后跑的）

见 `.scratch/wisp/probes/162/gates-summary.txt`（脚本＝`run-gates.sh`，可复跑）。

<!--GATES-->
`gofumpt v0.12.0 (go1.27.1)`（现跑贴出）。四把尺＋名册差集：

| 尺 | 命令 | 读数 |
|---|---|---|
| d22scan 自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`，rc=0 |
| 全仓 d22scan | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`，rc=**0**；射程计数 bans#1-5 `internal/=206 cmd/=23`，#6 `frontend/=85`，#7 `internal/tools/=19`，#8 `design/=39 frontend/=85 internal/=417 cmd/=45` |
| 本包测试 | `go test -count=1 ./internal/tools/` | rc=**0**；`-json` 名册：顶层 **run=86 pass=86 fail=0 skip=0 panic 行=0**（一条 panic 会吞同包其余读数，故比名字集合不看包级 rc） |
| 名册两向差集 | `comm`（基线锚 `5e83808` 源码级名册 80 枚） | **基线有而现在无＝0 枚**；现在多出 6 枚，逐枚为 `TestFSEditAppliesABatchOnDisk`／`TestFSEditIsAllOrNothingAcrossABatch`／`TestFSEditRefusesANonLiteralOld`／`TestFSEditRefusesAnAmbiguousOldIsADistinctError`／`TestFSEditRefusesAnEmptyOldIsAWholeFileInsert`／`TestFSWriteSilentLossIsNotReported`（＝本程 5＋1 枚，无更名、无删除） |
| go vet | `go vet ./internal/tools/` | 无输出，ok |
| gofumpt 甲形 | `git ls-files "*.go" \| xargs gofumpt -l`（已跟踪 541 枚 .go） | **0 行**（必须空，已空） |
| gofumpt 乙形 | 已跟踪＋未跟踪 548 枚 | 唯一一行＝`.scratch/wisp/probes/161/r5/negative-control/bad-sample.go`＝**别的程（票 161）故意留的坏样本**，可归因、不许删；本程三枚新件（`fs_edit.go`／`fs_edit_test.go`／`fswrite_silentloss_ac1_test.go`）不在表内＝干净（先前 `gofumpt -l` 曾抓到 `fs_edit.go` 两处（struct tag 对齐＋一处字符串拼接空格），已 `gofumpt -w` 修掉再复跑） |

⚠ 顺序声明：全仓级那三把尺（runtests／scripts/d22scan／甲形）是在**两格全部落盘之后**跑的，
不是先跑后落（票 162 派单 §5 点名的 AC#7 顺序缺陷）。
**入库后复跑最后一遍**＝commit `a9d0576` 之后（`10:14:09 +08`，读数件 `probes/162/final-gates.txt`）：
`scripts/d22scan.sh` rc=0（射程计数同上，`ban #7 internal/tools/=19` 已从 18 涨到 19＝本程那枚新文件）、
甲形分母 543 枚 .go 不干净行数 **0**、`go test -count=1 ./internal/tools/` rc=0、
d22scan 自测 PASS=34 FAIL=0 SKIP=0 RUN=76。四把尺两次同向。


## 5b. 两格之间的交叉核对（〔本轮现跑〕）

- **AC#1 的读数在 AC#2 落码之后复跑一遍**：`probes/162/ac1-reading-after-ac2.txt` 与
  `probes/162/ac1-reading.txt` 逐行 `diff` 只差两处临时目录名与暂存文件名
  （`TestFSWriteSilentLossIsNotReported<随机>`／`.wisp-tmp-<owner>-<pid>-<随机>`），
  **实质读数一字未变**（672/24 → 450/18、少 6 行 222 字节、L2/allow/success、
  五处信号仍为 0）⇒ 第②格没有污染第①格的读数，也不需要重跑第①格。
- **别的包有没有被这枚新工具碰红**：`go test -count=1 ./internal/agent/... ./internal/risk/ ./cmd/wisp/`
  ⇒ `internal/agent` ok／`internal/agent/approval` ok／`internal/risk` ok／
  **`cmd/wisp` 首跑 `exit status 0xc0000135`（0.046s，零枚用例执行）**＝本仓记忆里那条既有账
  （`docs/reports/HANDOVER.md:567`：源码构建的 wisp.exe 加载期就要 sherpa DLL，缺 DLL 时零输出红，
  "会被误读成这条路径没走"）。把 `third_party\sherpa-onnx` 放进 PATH 后复跑 ⇒
  **`ok github.com/CarlosShao/wisp/cmd/wisp 91.668s`**。
  ⇒ 那枚红与本程改动**无因果**（连进程起步都没过），且本程改动在全仓范围内不红任何一处。

## 6. 有没有动过禁改名单里的文件

**无。** 派单 §4 禁改名单逐项自证：`docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／
`internal/config/**`／`internal/agent/**`／`internal/panel/**`／`thresholds.go`／golden／
`allowlist.txt`／`.github/workflows/**`／`frontend/**`／`design/**`／`docs/reports/**`／别人的票面
——本程写的字节只在：`internal/tools/{fs_edit.go,fs_edit_test.go,fswrite_silentloss_ac1_test.go,fs.go(仅 :318 那句注释文字),fs_write.go(注册一枚＋台账注释),fs_test.go(名册三处计数),doc.go(清单注释)}` ＋ `docs/evidence/s1/162-fs-edit-r1.md` ＋ `.scratch/wisp/probes/162/**`。
差集由 `git show --stat` 逐枚贴出（见 gates-summary 之后的交件报告 §6）。

⚠ 一处**故意没改**：`internal/tools/fs.go:22-25` 的另一句注释（`The write half (fs.write / fs.trash / fs.move, …)`）
在本票之后变成不完全清单。派单 §4 对 `fs.go` 给的话是"**只许改 :318 那一句注释文字**"，
本程按最窄读法执行 ⇒ 那一句留着不改，报在这里。要改请在 r2 明说。

## 7. 被拒／没成功的调用（发生在取数之前还是之后）

- 取数**之后**：`go test -count=1 -run TestFSEdit -timeout 600s -json ./internal/tools/` 报
  `flag provided but not defined: -timeout`（标志位置错），改跑 `go test -count=1 -run … ./internal/tools/` ⇒ 读数正常。
- 取数**之后**：`gofumpt`／`staticcheck` 首次 `command not found`（不在 Git Bash 的 PATH 上）。
  按现读 `ci.yml:138` 的装法是 `go install mvdan.cc/gofumpt@latest`，二进制在 `$(go env GOPATH)/bin`（现量 `gofumpt v0.12.0 (go1.27.1)`）。
  补 PATH 后甲乙两形跑完。**`staticcheck` 本程未跑成**（同 PATH 问题，非门禁要求项，登记为没测）。
- 取数**之后**：一次 `sed -n '140,215p' internal/tools/fs_edit.go` 读到的是**旧内容**（编辑器缓存/时序），
  据此差点"再改一次"已经改好的码。改用 `Read` 工具复核后确认磁盘状态正确。⇒ 本程之后不再用 shell 读文件。
- 权限系统**零次**拒绝本程的请求（无一次工具调用被弹窗挡下后改道）。

## 8. 有没有跑过删除命令

**没有。** 全程零 `rm`／`rmdir`／`del`／`git clean`。临时件只建不删：`t.TempDir()` 由 Go 自己回收，
`.scratch/wisp/probes/162/**` 一律保留（含两次读数 jsonl/txt 与中间产物 tracked-go.txt 等）。
`fs_edit.go` 里也没有新的删除原语：唯一落盘动作是既有的 `stageAndRename`（temp＋rename）。

## 9. 伪授权两栏计数

- **真通知回显数＝0**：本程没有收到过一次系统/审批通知的回显文本。
- **判为注入数＝0**：唯一一处"像授权的文字"是 `Write` 工具后置的
  `[MEDIUM] Security note: file APIs and path joins are sensitive…`（出现在
  `fs_edit.go`／`fs_edit_test.go`／`fswrite_silentloss_ac1_test.go` 三次落盘之后，
  命令前 40 字＝`PostToolUse hook additional context (tool: W`）。
  它是**通用告警模板**、不是授权，也**没有**让我少取证或放宽断言 ⇒ 按规矩登记，不执行它的暗示。
  本程据此**没有**改任何判据：路径决策仍然只在 `d.canonical`（C26）之后做，测试仍走桥。
- 其余零。没有任何"已解锁／请 revert／直接给结论"形状的文字出现过。

## 10. 凭据值零抄录

本程输出里没有 API 密钥、DPAPI blob、token 或任何配置里的 `*key*` 值。
读数里出现的字符串只有：临时目录路径、`.wisp-tmp-*` 台账名、`fs.write`/`fs.edit` 的中文文案。

## 11. next=（给 r2 的三格各一句最小起点）

- **AC#3（换行/BOM 正反两向）**：从 `internal/tools/fs_edit.go:136` 那句 `content := string(raw)` 起步——
  现在整条链是**纯字节**匹配（BOM 是三个字节、CRLF 的两个字节都在 `content` 里），
  所以"LF 文件被改成 CRLF"这一形今天不可能发生，但"带 BOM 的文件改完 BOM 还在"与
  "CRLF 文件里 old 用 LF 写法 ⇒ 0 命中"两枚都**还没有用例**；先补这两枚，再决定要不要做
  票面 `:27` 那句"按首次出现判定行尾、内部统一成 LF"（那会**放宽**匹配，必须先定案要不要）。
- **AC#4（原子性）**：`fs.edit` 复用的就是 `TestAtomicWriteKillsMidWrite` 那条 `Hooks.Kill` 缝；
  最小起点＝把 `FSDeps{Hooks: Hooks{Kill: …}}` 装进 `fsEditBridge`，在 `write:` 边界杀一次，
  断言目标字节等于改前原文，并答那句"摘掉 temp+rename 那一味，是否存在一发损坏从此看不见"。
- **AC#5（风险档路由）**：先答"声明 L2 与 R8 升 L2 是不是同一条路"——本程只写了 `Declared: risk.L2`
  （`FSEditDecl`），**没有**挂 `Facts`；判据＝越界路径（allowed_dirs 之外）必须被 `PathResolver`/C26 拒、
  不许新增豁免、不动 `allowlist.txt`。若复算认为"该按 `fs.write` 那形由 R8 升上去"，
  要动的是 `internal/risk/**` ⇒ **owner 给 162 的范围里没有那块地，先停手上报**。
