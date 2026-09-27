# 票 162 · r3（写码位）＝AC#3b 两形变响亮 ＋ AC#4b 真硬杀那一半 ＋ "改完还是 CRLF" 反向尺半格 ＋ 判据顺序一笔

- 实现程：`162-r3`｜派单＝`.scratch/wisp/dispatches/2026-09-27-131x-impl-162-r3-two-silent-shapes-and-hard-kill.md`
- 锚点＝**本程 step 0 现量**（下面 §0），不是任何旧档基线。
- 本件顺序照派单 §6 的十二节走；**勾不在本件**（派单："所有票面一格都不许勾"）。

## 0. step-0 五件〔本轮现跑，`2026-09-27 16:1x +08`〕

```
$ date                                   → Sun Sep 27 16:13:16 CST 2026
$ git rev-parse --abbrev-ref HEAD        → dev            （逐次见下表）
$ git rev-parse HEAD                     → ff0007840ac9489c6c5ef017ca5092cdd2ed6094
$ git status --porcelain -- internal/tools/   → （空，0 行输出）
```

分支逐次（每次动手前重量，派单第 1 项要求"含分支名逐次"）：

| 时刻 | `git rev-parse --abbrev-ref HEAD` | `git rev-parse --short HEAD` | 本节 |
|---|---|---|---|
| 16:13 | `dev` | `ff000784` | step 0 起手 |

AC#3b／AC#4b 连现量行号（`grep -n "AC#3b\|AC#4b" .scratch/wisp/issues/162-….md` 现量＝票面 **`:36-38`** 与 **`:39`**；行号是本程自己量的，不是抄派单）：

> **`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md:36`**
> - [ ] **AC#3b（09-27 11:0x 编排者追加，来路＝162-r2 登记的"今天静默通过的两形"）**：本票的立票原则是**"响亮失败，绝不安静改坏"**，而 r2 实测有两形**今天不响**（档位〔程读数，我未复算〕，r3 必须自己复现）：
>   **:37** ① **命中段跨过文件头那三个字节（BOM）⇒ 现在会把 BOM 一起换掉**（`14→11 字节`，回执零提示）；② **`new` 用 LF 写法落进一份纯 CRLF 的文件 ⇒ 现在会造出混合行尾**（`3 组 CRLF＋1 条裸 LF`）。
>   **:38** ⇒ 各加一枚拒绝（或各加一条**外部可见**的告警，两形选一形并写明为什么）；⚠ **必须同时保留一发合法用例证明它没被顺手放宽**（"CRLF 文件里 `old` 也写 CRLF ⇒ 命中并落盘"那一发就是现成的对照）。**判据不许写成"以后再说"**。
> **`:39`**
> - [ ] **AC#4b（09-27 11:0x 编排者追加，来路＝162-r2 自己点名的"没测什么"）**：r2 的原子性用的是 `Hooks.Kill`——**它返回 error，会跑 Go 的清理**；**真实进程被杀（`taskkill`）那一半不跑清理，今天没为 `fs.edit` 验过**（`fs.write` 那半由票 112/118 那族钉过）。⇒ 补一发：硬杀之后 (a) 目标文件字节必须＝改前原文，(b) 目录里留下的 `.wisp-tmp-*` 残件**必须被点名**（留下可以，静默不行——它是可回收的垃圾，不是损坏）。⚠ 不许为了"干净"新加删除原语（`AGENTS` 无此禁止，但**删除临时件要与 `fs.write` 同一把尺**，不许两族各造一套）。

### 0a 两形的前提复算（派单："当已知、但仍要复算"）

未修码＝本程 step 0 的锚点 `ff000784`，`go test -count=1 -v ./internal/tools/`（原始件
`.scratch/wisp/probes/162/r3/baseline-gotest-verbose.txt`；命令见 §7 第一行）：

```
fs_edit_ac34_test.go:372: 缺陷读数①：14 字节带 BOM 的文件，old 从偏移 0 起并含 BOM ⇒
        BOM 被当普通字节替换掉，落盘 11 字节、前三字节 41 4c 50，答复无任何提示（IsError=false）
fs_edit_ac34_test.go:413: 缺陷读数②：纯 CRLF 文件里插入一段 LF 写法的 new ⇒ 落盘后 3 组 CRLF + 1 条裸 LF，
        答复（"fs.edit 已改写 …\crlf-mixed.txt：1 枚编辑全部生效，23 字节 / 3 行 → 30 字节 / 4 行
        （临时文件+原子重命名）"）里没有任何行尾提示
```

⇒ **派单/票面写的前提本程复算＝成立**：两形今天确实都静默（①`14→11`、②`3 组 CRLF＋1 条裸 LF`，
两个数与票面 `:37` 逐字对上；162-v1 自己那两发是 `9→6` 与同样 `3+1`，同向）。
本格因此**不需要重新论证问题存在**，直接做修法。

### 0b 分支与锚点逐次

```
$ git rev-parse --abbrev-ref HEAD      → dev   （step-0 起手 16:13 ＋ 三格提交前各一次 ＝ 四次，全部 dev）
$ git status --porcelain -- internal/tools/  → 0 行（起手；三格提交后各自复量仍 0 行）
```

提交链（**这是提完三枚之后从 `git log` 复原的，不是当时逐次现量**：派单第 1 项要求"含分支名逐次"，
本程逐次量的是**分支名**（四次，全 `dev`），HEAD 只有起手那一次是现读。别程 `171-r1` 与编排者的提交
夹在本程两格之间，故父链如下，读者按它取射程）：

```
ff000784 16:11  ← 本程锚点（step 0 现量 HEAD；此后本程不再改 HEAD 的读法）
3df74960 16:23  别程 171-r1 的仪器程提交
0cb6e8f2 16:26  别程 171-r1 的 Progress log 提交   ← 6823cbc6 的父
6823cbc6 16:33  本程第①格＋反向尺半格（同片文件故一次提）
422c1bc7 16:37  编排者 ledger(A334)＋票 171        ← b1cbdad4 的父
b1cbdad4 16:37  本程第②格 AC#4b
3d98a8ab 16:38  本程判据顺序那一笔（纯测试）
```

⇒ **本程三枚提交之间确实有别人入库**（`3df74960`／`0cb6e8f2`／`422c1bc7`），
所以 §2 的射程用**锚点区间** `ff000784..HEAD` 量，并在同一条命令里把 `internal/risk/` 等禁改面**单列**成空输出；
⚠ 别把别程落在 `docs/evidence/s1/171-…`、`票 171 票面` 上的字节挂到本票。

命令（可复制）：`git rev-parse --abbrev-ref HEAD`｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/`。
全程**零切分支、零 switch/checkout/merge/rebase/reset/stash/worktree/clean**；提交一律一步式
`git add -- <显式路径> && git commit -q -F - -- <同批显式路径>`，**没有把 `git add` 当独立一步**、零 `git add -A`、零 `git commit -a`。

## 1. 本程**没**测什么

1. **AC#5（风险档与门控）与 AC#6（契约轴）两格本程一律未做、未判**（派单 §5 的写面里就没有它们）。
   本程只声明一句：`internal/risk/**` 与 `docs/PLAN.md`／`docs/specs/**`／`.github/**`／`frontend/**`／`design/**`／
   `tools/d22scan/**` 在 `ff000784..HEAD` 的 numstat 里**零字节**（命令与输出在 §2），
   这是"没越权"的读数，**不是 AC#6 的判语**。
2. **AC#4b 的非 Windows 那一支＝未经验证**。`fs_edit_ac4b_kill_windows_test.go` 带 `//go:build windows`，
   在 Linux runner 上**根本不进编译**⇒ 那一侧既没有分母也没有读数，本程不许读者把"整包绿"当它响过。
   分母问题（本机 or CI）的答复在 §5c。
3. **行尾归一层（`判据形状 1` 的第二半，票面 :27"剥 BOM／统一成 LF／写完还原"）本程没有实现**，
   也没顺手实现一半：两枚新增检查只**拒绝**形状，不放宽任何匹配（这是 §3 那枚 `…CannotBeFedAPartialBOM…`
   存在的理由——它把"① 的可达形到底有多大"量出来，而不是让注释替它说话）。
4. **已混合行尾的文件（无单一约定）今天仍无任何提示**——本程把它写成一条**显式对照**
   （`…mixed_file_is_not_refused` 断言"不该被拒"）而不是一条静默；要不要给这一形也加告警＝新判据，见 §12 第 3 条。
5. **残件回收那一问本程未复算**：fs.edit 硬杀留下的 `.wisp-tmp-*` 会不会被"下一次写盘"的清扫器带走——
   〔现读码：`fs_write.go:288` 的 `reclaimStaging(parent, se)` 在 `CreateTemp` 之前按目录＋前缀扫，
   fs.edit 走的正是同一条 `stageAndRename` ⇒ 同一把尺、同一目录〕，**本程未跑用例钉它**（A18 那枚只钉了 `fs.write` 触发的清扫）。
6. 没跑全仓 `go test ./...`（派单 §4：逐包单跑，别跑全量）；没跑 `-cover*`（与 `-overlay` 不合跑）。
7. 本程**没有勾任何票面框**，`docs/reports/**` 一字未写（台账与勾归编排者）。

## 2. 写面自证与禁改面

```
$ git diff --numstat ff000784..HEAD -- internal/tools/
132	6	internal/tools/fs_edit.go
54	56	internal/tools/fs_edit_ac34_test.go
384	0	internal/tools/fs_edit_ac3b_test.go
274	0	internal/tools/fs_edit_ac4b_kill_windows_test.go
$ git diff --numstat ff000784..HEAD -- internal/risk/ docs/PLAN.md docs/specs/ .github/ frontend/ \
      design/ tools/d22scan/ internal/config/ internal/agent/
（空输出＝零字节）
```

- 派单 §5 的写面之外本程**没有动过一枚文件**；`probes/162/v1/**`、`probes/154/**`、`probes/171/**` 一字未写
  （`$ git diff --name-only ff000784..HEAD -- .scratch/wisp/probes/162/v1 .scratch/wisp/probes/154 .scratch/wisp/probes/171` → 只出现 `171-r1` 自己的路径，**本程三枚 `--name-only` 里 0 枚**）。
- ⚠ 归因提醒照抄不采信：**`internal/risk` 自 `5e83808` 起那三行 numstat 出处是 `f576cf08`（票 160）**，
  本程锚点区间 `ff000784..HEAD` 里 `internal/risk/` 是**零字节**（上面第二条命令现量），与本票无关。
- **没有触发"必须动 `internal/risk/**` 才能走到 L2"那一支**：本程两格都不需要新的档位判断，
  `FSEditDecl` 一字未动，`allowlist.txt` 一字未动。

## 3. 第①格＝AC#3b：两形从静默变响亮

**两形都取"硬拒"**（票面 `:38` 允许"拒绝 or 外部可见告警"，选一形并写理由）：

1. 告警必须先把损坏落盘才谈得上被读到，而本工具的立票原则是**响亮失败、绝不安静改坏**——
   BOM 没了、行尾混了都在文案被读之前发生了；
2. 两形各有一句模型能直接照做的补救（`old 从 BOM 之后抄`／`new 按文件自己的行尾拼写`），
   够得上 `判据形状 2/3` 的"文案里写补救动作"，不构成一堵墙；
3. 两道检查都落在**任何写入之前**（`fs_edit.go:193` 在定位循环内、`:244` 在 `updated` 之后 `stageAndRename` 之前），
   所以自动继承本工具"整批全有或全无、台账为空"的性质；告警形做不到这一点（它必须已经落盘）。

### 3a 外部可见读数：指名为**两处**，改前／改后两版

派单点名"不许改了判据但没人看得见"。本程指名：**(i) 回执 `agent.ToolOutcome` 的 `IsError`／`ErrorClass`／`Text`；
(ii) 同一枚调用在 bridge 审计行里的 `outcome=` 列**（`bridge.go:505-509` 由 `res.IsError` 推 `ErrorClass=observe.ClassTool`
与 `kind=OutcomeKindError`，`bridge.go:855-858` 把 `kind` 打成 `outcome=` 字段；同一字段也是 `tool_call.outcome` 列的来源）。
第 (ii) 处不是引代码行号充数：`fs_edit_ac3b_test.go` 里 `fsEditBridgeAudited` 造了一枚带 `Logf` 收口的桥，
**把那一行审计原文抓进断言**（`auditLineFor` 还钉"命中的必须是恰好一枚"，防"探针没接到也算绿"）。

| 形 | 改前（锚点 `ff000784`，`probes/162/r3/baseline-gotest-verbose.txt:126/:128`） | 改后（`probes/162/r3/ac3b-after.txt:46-48`/`:251-254`） |
|---|---|---|
| ① BOM | `IsError=false`、文案 `"fs.edit 已改写 …：1 枚编辑全部生效，14 字节 / 2 行 → 11 字节 / 2 行"`（BOM 三字节日光化消失，`前三字节 41 4c 50`） | `IsError=true ErrorClass="tool" AppliedSteps=0`、文案 `"fs.edit 拒绝执行：第 1 枚编辑的 old 从文件第 0 字节起，圈住了文件头那 3 字节 EF BB BF（UTF-8 BOM，文件 14 字节）。BOM 是文件自己的编码标记…请把 old 和 new 都从 BOM 之后的第一个字节开始抄；…"`；审计行 `tool=fs.edit risk=L2 decision=allow outcome=error` |
| ② 行尾 | `IsError=false`、文案 `"fs.edit 已改写 …23 字节 / 3 行 → 30 字节 / 4 行"`，**"行尾/不一致/mixed endings/endings" 四探针全不在文案里**，落盘 `3 组 CRLF + 1 条裸 LF` | `IsError=true ErrorClass="tool" AppliedSteps=0`、文案 `"…第 1 枚编辑的 new 会改写这份文件的行尾约定：文件原本是纯 CRLF，改完会变成 CRLF 3 组 / 裸 CR 0 条 / 裸 LF 1 条。请把 new 里的换行按文件自己的行尾拼写（纯 CRLF 的文件：换行写成 \\r\\n…）；确实要换行尾请改用 fs.write…"`；审计行 `outcome=error` |

r2 留在 `fs_edit_ac34_test.go` 里那两枚"钉住今天行为"的用例（`TestFSEditSilentShapeChangesLandToday`）
**函数名保留**（名册是稳定键，`缺`那一侧必须 0），断言整体翻到响亮那一侧（翻转读数：
`ac3b-after.txt:33/:35`，两枚都 `IsError=true ErrorClass="tool"`），并在注释里写明它已不钉今天行为、
以及"合法方向由 `fs_edit_ac3b_test.go` 的对照钉住"。

### 3b 合法对照（证明门没被焊死）

```
$ go test -count=1 -v -run 'TestFSEdit' ./internal/tools/   （原始件 probes/162/r3/ac3b-after.txt，rc=0）
AC#3b① 合法对照：14 字节 → 14 字节，前三字节仍是 EF BB BF（old 从 BOM 之后抄起 ⇒ 命中并落盘）
AC#3b① 合法对照（无 BOM 文件、old 从第 0 字节起）：11 字节 → 11 字节，未被拒绝
AC#3b① 射程读数：old 只带 BOM 的第 3 个字节 ⇒ 拒因是「0 命中」而不是「吞 BOM」，落盘第 0 字节是 U+FEFF
AC#3b② 合法对照一（同一枚编辑、new 按 CRLF 拼写）：23 字节 → 31 字节，CRLF 3→4 组、裸 CR/裸 LF 均 0
AC#3b② 合法对照二（LF 文件 + LF 写法 new）：20 字节 → 27 字节，未被拒绝
AC#3b② 合法对照三（已混合的文件不被告知该长成什么样）：22 字节 → 29 字节，行尾形 混合(CRLF=2 CR=0 LF=2)
```

**"判据不许写成喂什么都会响"这一条的正面答复**：②的检查只在文件**本来有单一约定**时开火，
且红的是"new 带进异形换行"这一枚——`mixed_file_is_not_refused` 与 `lf_file_lands_an_lf_spelled_new`
两枚对照钉住"没约定就不拒、有约定且照抄就落盘"。①的检查只在文件真有 BOM 且命中段压到那三字节时开火，
`无 BOM 文件从第 0 字节改起`那枚钉住另一半。

### 3c 恒真那一问（逐枚，两向读数）

未修码＝`-overlay` 换成锚点 `ff000784` 那一版 `fs_edit.go`（rig 见 §7 命令，副本反扫 `present-hits=1 absent-hits=0`）。
原始件 `probes/162/r3/logs/base-on-unfixed-code.log`：**顶层 19 枚跑完，PASS=16 / FAIL=3**。

| 用例 | 未修码（base overlay） | 修码后 | 答复 |
|---|---|---|---|
| `TestFSEditRefusesAnOldThatSwallowsTheFileBOM` | **FAIL**（回执是成功、BOM 被吃掉） | PASS | **不响** ⇒ 这一枚是"把今天静默的形变响亮"的尺，不是恒真 |
| `TestFSEditRefusesALFSpelledNewInAPureCRLFFile` | **FAIL**（成功落盘出混合行尾） | PASS | 同上，**不响** |
| `TestFSEditSilentShapeChangesLandToday`（含两枚子例） | **FAIL**（翻转后钉的是响亮，未修码给的是静默） | PASS | 同上，**不响** |
| `TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM` | PASS | PASS | 对照枚**两向都绿＝它不测量修码**，它钉"没被放宽"，不参与恒真 |
| `TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak`（两子例） | PASS | PASS | 同上 |
| `TestFSEditCannotBeFedAPartialBOMThroughJSON` | PASS | PASS | 它钉的是**传输射程**（JSON 送不进残缺 BOM 字节），两向都绿是有意义的：两版都不吞 BOM |
| `TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten` | PASS | PASS | 反向尺：它的红来自 M-3，不来自修码与否（§4） |

## 4. 缺口那半格：`改完还是 CRLF` 的反向尺

162-v1 实测 M-3（落盘前强行改成 CRLF）之下 `TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF`
与五形取证那枚**仍绿** ⇒ 正向有尺、反向无尺。本程先**复算它这两个"仍绿"**（下表 M-3 列里两枚都不在红名册），
再补 `TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten`：夹具是**含 CRLF 但不是纯 CRLF**
的 22 字节文件（2 组 CRLF + 1 条裸 LF），`old/new` 都不带任何换行 ⇒ 唯一能改动它行尾的就是"落盘前被强行改写"这一味。

**为什么必须是混合夹具**：对纯 CRLF 文件强制 CRLF 是恒等变换，任何纯-CRLF 夹具都分不开
"工具保住 CRLF"与"工具盖上 CRLF"；对纯 CR 文件 M-3 也是恒等（`\r` 既不含 `\r\n` 也不含 `\n`）；
纯 LF 文件那一侧 v1 已经量到会被 M-3 抓红。⇒ 判别器只能长在"含 CRLF、又不是只有 CRLF"的形状上。

```
$ bash .scratch/wisp/probes/162/r3/build-rigs.sh     （原始件 logs/m3-on-r3.log，副本反扫 present=1/absent=0）
M-3 叠加在 r3 修码之上：顶层红 11 枚，其中
  M3-RED --- FAIL: TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten   ← 那半格要的正是这一条
  M3-绿（v1 点名的两枚，本程复算同为仍绿）：
    TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF
    TestFSEditLineEndingForensicsRefusesRealSpellings
```

⇒ **反向尺成立**：M-3 之下它红，未修码与修码之后它都绿（§3c 表）＝它不是"喂什么都会响"。
**同一遍里 ② 那枚拒绝用例也红了**（`TestFSEditRefusesALFSpelledNewInAPureCRLFFile` 与 r2 翻转的
`…/an_lf_spelled_new_in_a_crlf_file_is_refused_r3`）——红因正是"强制 CRLF 把要挡的形状抹掉了，于是拒绝不再开火"，
⇒ ② 那一格其实有**两把**反向尺：混合夹具那枚（钉字节）与拒绝那枚（钉"守卫还在开火"）。
⚠ 其余 **8 枚红是 collateral**（`…KeepsABOM…`、`…DoesNotRewriteAnLFFiles…`、`…AppliesABatchOnDisk`、
两枚 AC#4 用例、`…StillLandsAnOldStartingAtByteZero…`、`…RefusesAnOldThatSwallows…`〔它红在自己内部那枚对照上：
BOM+LF 文件被强制成 CRLF ⇒ 撞上②号拒绝〕、`…StillLandsNewEndingsWhereThereIsNoConventionToBreak`〔两子例：
LF 那枚被误拒、混合那枚字节不再等值〕）——红因都是同一件事：强制 CRLF 让 LF 系夹具撞上本程新长的②号拒绝，
**这是"这一味现在有人会挡"的读数，不是回归**；逐枚红因在 `logs/m3-on-r3.log`。
另一处点名：`TestFSEditSilentShapeChangesLandToday/an_old_that_carries_the_bom_is_refused_r3` 在 M-3 下**仍绿**，
因为①的拒绝发生在落盘之前，M-3 那一味根本到不了。

## 5. 第②格＝AC#4b：真 `taskkill` 那一半

新用例 `TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue`
（`internal/tools/fs_edit_ac4b_kill_windows_test.go`，`//go:build windows`，**文件内零 `t.Skip`**）。
形状照 `bridge_a18_kill_windows_test.go`（**那个文件一字未改**，连它的 `stagingFiles` 都是复用不是另造一把尺），
子进程在 `write:8` 这道边界挂死（`WriteChunk:4`），父进程 `taskkill /F /PID` 真杀。

### 5a 字节级读数（`probes/162/r3/ac4b-kill.txt`，rc=0，`--- PASS (0.48s)`）

```
AC#4b(a) 目标 31 字节＝改前原文（"alpha\nbravo\ncharlie\ndelta\necho\n"），
         要落的 74 字节 "ALPHA-KILLED-MID-WITHOUT-RUNNING-ANY-GO-CLEANUP\n\nbravo\n…" 一字未出现在盘上目标里
         断言形状：got 与 original 逐字节比（红时打印 % x）＋ HasPrefix(got,intended) 反向断言
                    ＋ Contains(got, new) 反向断言 ⇒ 停在"没报错"的读数做不到通过这三关
AC#4b(c) 事后 fs.read：IsError=false ErrorClass=""，内容＝原文 31 字节（下一台仪器看见的是干净文件）
```

### 5b 残件点名清单（**只建不删**，本程没清）

| 残件 | 所在目录 | 大小 | 内容 |
|---|---|---|---|
| `.wisp-tmp-2e5a8af8-47540-1693964025` | `…/Temp/TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResi1007168199/001/edit` | 4 字节 | `"ALPH"`＝要落的新内容的**前 4 字节** |

断言是 `len(residue) != 1 即失败`，并把每一枚的**路径、大小、mtime、内容**打进日志——
"半个新内容"确实存在，只是存在于**该在的那个文件**里。0 枚＝taskkill 没打断任何事，多枚＝尺坏了。
⚠ 残件留在用例自己的 temp 目录里，本程代码零删除动作；目录级回收是 Go 的 `t.TempDir()` 收尾做的，
**不是本用例的判语**（也正因如此，断言顺序是"先量残件、后跑 fs.read"，任何后续写盘都会让清扫器把这枚凭据扫走）。

### 5c 分母：本机跑还是有分母，别一支标未经验证

- **本机（Windows＋`System32\taskkill.exe`）有分母**，凭据不是自称：同一遍 `go test ./internal/tools/` 里
  `TestA18RealTaskkill…` 与本格**同遍 PASS**（`final-gates-gotest.txt`：`--- PASS: TestA18RealTaskkill…`、
  `--- PASS: TestFSEditRealTaskkill…`，该遍 `SKIP=0`），⇒ 真杀这一形在这台机器上今天能跑、跑了有数。
- **非 Windows runner 那一支＝未经验证**：`//go:build windows` 让它根本不参与编译，
  本程**没有**把"整包红/整包绿"当它的读数。派单要的答复就这一句：**本机有分母，跨平台那一支无分母。**
- 尺不是死开关（§7 有 rig）：M-4 自造副本（摘掉 `stageAndRename`、同一 `WriteChunk`、同一杀点）下本枚**变红**，
  红因正是字节断言：`AC#4b(a) violated：硬杀之后目标不是原文。got 4 字节 41 4c 50 48`
  （`logs/m4-on-r3-ac4b.log`，rc=1）。同一遍里 `(b)` 那一关还没到就被 `(a)` 挡住 ⇒ 这一点与 §6 是同一族缺陷，
  本程**没有**把它顺手也改成"先看残件"，因为 AC#4b 的判据本体就是 (a)；登记给 r4（§12 第 4 条）。

## 6. 判据顺序那一笔：改前红在哪／改后红在哪

162-v1 点名：`:458` 那枚边界探针会挡住 `:471` 的字节断言。**本程现量**（同一枚探针，行号因上文改动上移）：

```
改前：$ go test -count=1 -overlay=.scratch/wisp/probes/162/r3/overlay-m4.json -v \
        -run '^TestFSEditKilledMidWriteLeavesTheTargetByteIdentical$' ./internal/tools/   → rc=1
  fs_edit_ac34_test.go:448: expected the write to reach the staging boundaries, got ["write:8" "write:16"]
  ⇒ 红在边界探针；字节等值断言一行未执行（logs/order-before-m4.log）

改后：同一命令 → rc=1
  fs_edit_ac34_test.go:456: AC#4 violated: after a mid-write kill the target holds 8 bytes "ALPHA-LO",
                            want the original "alpha\nbravo\ncharlie\ndelta\necho\n"
  ⇒ 红在字节上，与 162-v1 自造那枚 TestV1AC4TargetIsNeverHalfWritten 的读数同向（logs/order-after-m4.log）
```

只动测试：`无 Skip、无删除断言、生产路径零字节`（本枚 diff 只在 `fs_edit_ac34_test.go`，numstat 见 §2）。
未变异仍 PASS，且台账五步读数一字未少（`AC#4(a)(b) 读数 … 台账 ["在 … 创建临时文件 .wisp-tmp-…" …
"删除临时文件 …（目标从头到尾未被改动）"]`）。

## 7. 门禁四数＋名册两向差集（逐包单跑，未跑全量；全部本程现量）

```
$ go test -count=1 -v ./internal/tools/                      （final-gates-gotest.txt）
RUN=150  PASS=150  FAIL=0  SKIP=0  panic=0        顶层 PASS=101（起手基线 94）
  逐格：第①格后 gates-cell1-gotest.txt = RUN=149/PASS=149/0/0，顶层 100
        第②格后 gates-cell2-gotest.txt = RUN=150/PASS=150/0/0，顶层 101

$ git grep -h -E '^func (Test[A-Za-z0-9_]*)' ff000784 -- internal/tools/ | sed -E 's/^func (Test[A-Za-z0-9_]*).*/\1/' | sort -u \
      > probes/162/r3/roster-base-head.txt          # 94 枚
$ grep -h -E '^func (Test[A-Za-z0-9_]*)' internal/tools/*.go | sed -E 's/^func (Test[A-Za-z0-9_]*).*/\1/' | sort -u \
      > probes/162/r3/roster-now.txt                # 101 枚
$ comm -23 roster-base-head.txt roster-now.txt > roster-missing.txt   # 缺 = 0 枚 ✅
$ comm -13 roster-base-head.txt roster-now.txt > roster-added.txt     # 多 = 7 枚，逐枚全是本程新件：
  TestFSEditCannotBeFedAPartialBOMThroughJSON
  TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten
  TestFSEditRefusesALFSpelledNewInAPureCRLFFile
  TestFSEditRefusesAnOldThatSwallowsTheFileBOM
  TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM
  TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak
  TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue
  （更名 0 枚：r2 的 14 枚名册全部原位，只翻断言不翻名字）

$ sh scripts/d22scan.sh                              rc=0  clean（ban #7 internal/tools/ 19 枚生产文件在册）
$ bash tools/d22scan/runtests.sh -C tools/d22scan ./...   rc=0  PASS=34 FAIL=0 SKIP=0 RUN=76（基线现量，与 v1 同数）
$ export PATH="$PATH:$(go env GOPATH)/bin" && gofumpt --version   → v0.12.0 (go1.27.1)
$ git ls-files '*.go' > tracked-go-final.txt && wc -l < tracked-go-final.txt   → 561 枚（v1 现量 554，差 7 枚＝本程与别程入库件）
$ gofumpt -l $(tr '\n' ' ' < tracked-go-final.txt)                       → 0 行 ✅（甲形）
  （未跟踪新件另测：gofumpt -l 一并带上 fs_edit_ac3b_test.go / fs_edit_ac4b_kill_windows_test.go → 仍 0 行）
$ bash .scratch/wisp/probes/162/r3/build-rigs.sh     # 三枚副本各自反扫：present-hits=1 / absent-hits=0
  base：'Deliberately NOT here' 在、'leadingBOMLen'/'lineEndingConvention' 均 0 命中（＝该符号在未修码里不存在）
  m3  ：'MUT-R3 same shape as 162-v1 MUT-3' 在、'	updated := b.String()' 0 命中
  m4  ：'MUT-R3-M4 no staging file' 在、't.d.stageAndRename' 0 命中
  ⚠ `-overlay` 与 `-cover*` 全程未合跑。
```

**入库后复跑**（改完最后一行注释——`fs_edit_ac3b_test.go` 里那枚夹具的字节数从"7 bytes"改成现量的 9——之后重测）：
`go test -count=1 ./internal/tools/` → `ok 15.702s`；`sh scripts/d22scan.sh` → rc=0（`final-d22scan-after-comment.txt`）；
`gofumpt -l <561 枚已跟踪 .go> fs_edit_ac4b_kill_windows_test.go` → **0 行**（`gofumpt-A-final2.txt`）。
四数与名册不再重跑（注释改动不参与名册；那一遍 `final-gates-gotest.txt` 是收口数）。

**变异副本先自证那一道救回了一次真事故**：rig 第一遍跑出来的形状是 `rc=1 且 红名册 0 枚`——
不是"全绿"，是我那枚新用例里 `Logf` 用了 `%d` 装了 string（vet 编译失败）。如果按"红名册为空＝没抓到"下判语，
这一格会被记成"M-3 抓不到反向尺"。**它没进任何判语**：修完重跑才是 §4/§5 的读数（详见 §8 第 4 条）。

## 8. 被拒／没成功的调用（各标"取数之前／之后"）

| # | 调用 | 结果 | 发生在 |
|---|---|---|---|
| 1 | `diff internal/tools/… .scratch/…/mut/m3/fs_edit.go` 相对路径写错 | exit 2，文件找不到 | **取数之前**（step 0 之后第一发；改用绝对路径重跑，读出 M-3 形状） |
| 2 | `gofumpt -l … .scratch/wisp/probes/162/r3/*.go` | `CreateFile … 语法不正确`（glob 未展开） | 取数之后（基线四数与 AC#3b 改后读数已在手）；去掉 glob 重跑＝§7 的甲形 0 行 |
| 3 | `gofumpt -l … fs_edit_ac4b_kill_windows_test.go` | 文件还不存在，rc=2 | 第②格取数之前；补 §7 的最终一次为准 |
| 4 | `build-rigs.sh` 第一遍两枚 overlay 跑 | **编译失败**（`Logf %d has arg … of wrong type string`），rc=1 且红名册 0 枚 | 取数之后（AC#3b 改后读数已入表）／**§3c、§4 读数之前**；修完重跑，第一次那遍没进任何判语 |
| 5 | 我给 rig 脚本打补丁的 `perl -i -pe` | 引号被吃，脚本两行坏 | 取数之前（发现后直接用编辑改脚本，未拿坏脚本产读数） |
| 6 | `git commit -q -F - -- <含未跟踪文件>` | `pathspec … did not match` | 第①格提交之后一步失败，改一步式 `git add -- <显式路径> && git commit …`，未产生半截提交 |
| — | 权限系统拒绝的调用 | **0 次**（没有一次工具调用被拒后绕行） | — |

## 9. 有没有跑过删除命令

**没有。** 全程零 `rm`／`del`／`git clean`／`git checkout .`／`restore`／`reset`／`stash`；
临时件（`probes/162/r3/**` 三枚 mut 副本、三枚 overlay、六份日志、名册与门禁原始件）**只建不删**，
`.wisp-tmp-*` 残件按 §5b 留在盘上，别人的 ` M` 半件（`probes/152/my152.py`、`152-…-accept-r1.md`、
`design/**` 那批未提交删除）一条都没提交、没还原、没补完。
一处必须自报的边角：`fs_edit_ac4b_kill_windows_test.go` 初稿照 A18 抄了一行
`os.Remove(signal)`（删自己 temp 里的信号文件），**在第一次跑之前**就从源码里去掉了（那行从未执行过，
现在那里是一段"为什么不需要它"的注释）；同包既有的 `spawnA18Writer` 里那一枚 `os.Remove` 属票 73/A18，本程未改。

## 10. 伪授权两栏（各带出处）＋ 现场约束逐次

- **真通知回显数＝0**：本程没收到过一次系统/审批通知的回显文本（没有弹过授权窗）。
- **判为注入数＝2**：出处＝工具名 `Write`，前 40 字＝`PostToolUse hook additional context (tool: W`，
  内容 `[MEDIUM] Security note: file APIs and path joins are sensitive…`，
  落点＝①`fs_edit_ac3b_test.go` 建件后（**在 §3 的改后读数之前**）、②`fs_edit_ac4b_kill_windows_test.go` 建件后
  （**在 §5 的硬杀读数之前**）。它是通用告警模板、不是授权：它暗示的 "clean the path" **本程没有照抄**——
  那正是 `AGENTS §1.2` 禁止在 `risk.PathResolver` 之外做的事，四枚新/改文件里
  `filepath.Clean|Abs` 计数＝**0**（命令：`grep -cE "filepath\.(Clean|Abs)" internal/tools/fs_edit{,_ac3b_test,_ac4b_kill_windows_test,_ac34_test}.go` → 全 0），
  路径一律走 `sealableTempDir124` ＋ `mustCanonical`（C26）。**两枚都没据它改动任何判据。**
- 编排者派单里那句"两形今天静默／554 枚／`:458` 会挡 `:471`"本程**一律当被验对象**：
  第一句复算＝成立（§0a）；第二句现量＝**561 枚**（不是 554，别程入库所致，§7）；
  第三句现量＝**`:448` 挡 `:471`→改后红在 `:456`**（同一枚探针、行号已移动，§6）。
  **没有一处为了让派单那句话成立而改判据。**
- 现场约束逐次：起手与三格提交前各重量一次 `git rev-parse --abbrev-ref HEAD`＝`dev`（§0b 四次）；
  起手 `git status --porcelain -- internal/tools/`＝**0 行**（无兄弟程在半截里），三格提交后各自复量仍为 0 行；
  同树另一程 `171-r1`（`probes/154`＋`probes/171`＋它自己的表）与本程写面**全程不相交**；
  本程三枚提交的 `--name-only` 里不含别人的路径。

## 11. 凭据值零抄录

本件与本程全部输出里没有 API 密钥、DPAPI blob、token 或任何配置里的 `*key*` 值。
读数里出现的字符串只有：临时目录路径、`.wisp-tmp-*` 名字、`fs.edit`／`fs.read`／审计行的中文与英文文案、
本程自造的字节样本（`"ALPHA-KILLED…"`、`"ALPH"`、`"ALPHA-LO"`、`"alpha\r\nbravo…"`）。

## 12. next=

1. **AC#5（风险档与门控，票面 `:40`）＝本票剩下最硬的一格**：`fs.edit` 与 `fs.write` 同级（L2 已定案）＋
   **越界路径必须被 `PathResolver` 拒**那一枚需要真门控读数；派单那条硬话原样留给 r4——
   **若判"必须动 `internal/risk/**` 才能走到 L2"＝停手上报**，本程未触这一支（§2）。
2. **AC#6（契约轴，票面 `:41`）＝只核零字节＋归因**：`PLAN.md:2536`／`SPEC-07:43` 两行已由编排者落完，
   实现程再动＝越权；复算时记得 `internal/risk` 那三行出自 `f576cf08`（票 160），不归本票。
3. **混合行尾文件要不要也提示**＝本程显式留下的射程边界（§1 第 4 条、§3b 对照三）。
   要放开＝新判据（"没有约定也要挡"或"没有约定就告警"），归 owner 裁，不归实现程假设。
4. **AC#4b 的两关顺序**：本程量到 M-4 之下 `(a)` 先红、`(b)`（残件点名）没被执行
   ——与 §6 是同一族形状。本格**故意没改**（AC#4b 的判据本体就是 (a)，且改了就动刚提交那一格）；
   r4 若要闭合，最小形状＝把 (b) 的残件量测提到 (a) 之前一次、两处都断言，别删任何一关。
5. **残件回收那一问**（§1 第 5 条）：现在只有〔现读码〕没有用例。要钉它＝在 A18 那枚的第③子例里
   把触发改成 `fs.edit`，或照本程 `fs_edit_ac4b…` 的夹具在杀后再发一次合法写盘并断言残件归零——
   ⚠ 后者与本程"残件必须留在盘上被点名"是同一目录里的对立断言，**必须拆成两枚用例**，别混在一枚里。
6. **勾与台账 `A##` 由编排者落**（本程票面一格未勾、`docs/reports/**` 一字未写）。


