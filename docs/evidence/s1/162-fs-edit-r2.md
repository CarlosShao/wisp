# 162-r2 证据件 — 票 162 的第①②格（AC#3 行尾/BOM 正反两向＋AC#4 原子性反面实验）＋裁给本程的一句注释

- 派单＝`.scratch/wisp/dispatches/2026-09-27-103x-impl-162-r2-crlf-bom-and-atomicity.md`
- 进场锚＝`bb679c6`（step-0 现量）｜交件时刻＝09-27 10:5x｜写码位＝本程（162-r2）｜前一程＝162-r1（`0cb8864 a9d0576 7aaff7a 0e36753`）
- 档别标注约定（四档）：**〔本轮现跑〕**＝这一程亲自跑过的命令与输出；**〔现读码，未跑用例〕**＝读了码但没为这句话跑用例；**〔盘上有件，我未复算〕**＝读了件/读了码，没重跑；**〔仅自述〕**＝只有我一句话，无凭据。
- ⚠ AC#5（风险档路由）／AC#6（契约轴复算）**本程一字未做**，按派单留 r3。`internal/risk/**`、`internal/config/**` 本程零字节（自证见 §7）。

## 0. step-0 四件〔本轮现跑〕

```
date                       → 2026-09-27 10:33:03 CST 2026
git rev-parse HEAD         → bb679c6062b832cb54ce68b06ffc4cdd8b7e30c3
git branch --show-current  → dsh/feat/frontend-p0
git status --porcelain -- internal/tools/ docs/evidence/s1/162-*.md → （空，符合派单要求）
```

⚠ **派单/交接话里那句"分支 dev"与现场不符**：现场分支是 `dsh/feat/frontend-p0`（`git branch -a --contains HEAD` 只有这一枚）。我**没按那句话去切分支**（那是编排者的地界），只在现量上继续，登记在此。

票面 AC 四行连行号（`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`）：

- `:34` `- [ ] **AC#3 换行/BOM 正反两向**（CRLF 文件改完还是 CRLF；带 BOM 的改完 BOM 还在；**不许把 LF 文件的行尾改掉**）。`
- `:35` `- [ ] **AC#4 原子性**：改到一半被取消不许留损坏文件（拿 `PLAN.md:1176` 那条当判据，并答"摘掉临时文件那一味，是否存在一发损坏从此看不见"）。`
- `:27`（判据形状第 5 条，**本程只做取证**）`5. 换行与 BOM：**先剥 BOM、按首次出现判定行尾、内部统一成 LF（旧文本与新文本同一把尺）、写完还原**（`:243`／`:307`）。`
- `:28`（红线，本程守住）`6. ⚠ **不许照抄的那半**：它的写入调用点是**整文件交出去写**（`tools/edit.ts:127`）。我们 `PLAN.md:1176` 早写了必须临时文件＋原子改名 ⇒ **这一格我们比它严，别因为"上游那样做"就回撤**。`

起点复算（派单断言，〔本轮现跑〕）：`internal/tools/fs_edit.go:136` 确为 `content := string(raw)`＝**纯字节**；整条定位链是 `strings.Count`/`strings.Index`（`:171`/`:182`），落盘字节是 `content[:start] + new + content[stop:]`（`:205-215`）。⇒ 派单那句"`fs.edit` 从不转换任何东西"**成立**，且它推出的两件事本程各自有了读数（§3）。

## 1. 本程**没**测什么（先写这节）

一句不许含三件，逐条列清"没测＋为什么是刻意的"：

1. **没测 AC#5（风险档路由／越界路径必须被 `PathResolver` 拒）**。派单 §0 明令留给 r3；本程**没动 `internal/risk/**`、没动 `internal/config/**`、没动 `allowlist.txt`**，也没写任何"越界路径"用例。§4 的 AC#4 用例走的是 r1 已有的 `fsEditBridge`（真 C26，root＝本用例自建的临时目录），它证的是"写被杀之后字节是什么样"，**不证**越界判定。
2. **没做 AC#6（契约轴复算）**。本程没跑 `git diff --numstat <契约锚>..HEAD -- docs/PLAN.md docs/specs/**`；那两行的现量归 r3。
3. **没实现票面 `:27` 的行尾归一／剥 BOM**（派单 §1 明令"本格只做取证"）。§3d 的读数是"不实现要付什么"，不是"实现了会怎样"。
4. **没测 Unicode 修复层**（判据形状 1 第二半，NFKC／去尾空白／弯引号→ASCII）。r1 已把"要不要做"单独挂账，本程不碰：它每宽一分就把一次响亮失败变成一次静默改写，而行尾归一（`:27`）**正是同一族放宽**——所以本程特意把两件事的读数分开登记（§3d 与 §3e），不许混成"归一 layer 一起做完"。
5. **没测真实进程死亡（taskkill）那一半**：本程用的是 `Hooks.Kill` 缝，它**返回 error ⇒ 会跑 Go 的清理**（`discard()` 删暂存件）。真实 `taskkill` 不跑任何清理，那一半在 `internal/tools/bridge_a18_kill_windows_test.go` 里为 `fs.write` 钉着；`fs.edit` 复用同一枚 `stageAndRename` ⇒ 结构上共享 A18 的开局清扫（`reclaimStaging`），但**本程没为 `fs.edit` 写过一发真实进程死亡的用例**。〔现读码，未跑用例〕——这句是给 r3/验收的，不是本程的结论。
6. **反面实验是构造出来的写路径，不是产品码**：§4c 那枚直写目标用的是生产自己的 `sideEffect.writeAll`＋同一道杀点，**唯一变量是缺 temp＋rename**；它不证 `stageAndRename` 的其余失败分支（`Sync` 失败、`Close` 失败、`os.Rename` 失败）——那些由 r1 之前的 `fs_write_test.go` 射程管，本程未复算。
7. **没测大小上限与行尾的交互**：`writeCap()` 那两枚拒绝（改前超限／改后超限）是 r1 的码，本程没为它们补用例，也没测"一个大文件里混着 BOM 和 CRLF"的整链路。
8. **没测真实模型送什么**：本程所有调用都是手写字节走桥，没走 `C5 LlmProvider` 的 golden SSE。"模型会不会把 `new` 按 LF 拼写进 CRLF 文件"（§3e 的那枚缺陷读数的现实概率）是产品效果轴，不在票面射程。

## 2. 交件的三枚文件（〔本轮现跑〕）

| 落点 | 内容 | 提交 |
|---|---|---|
| `internal/tools/fs_edit_ac34_test.go` | 新建，8 枚顶层用例（AC#3 六枚＋AC#4 两枚） | `64b522b6`（与下一行同提交） |
| `internal/tools/fs_edit_test.go` | 只把 r1 的 `fsEditBridge` 拆成"薄壳＋`fsEditBridgeDeps(t, root, tweak)`"，r1 的 5 枚用例一字未改、断言一字未松 | `64b522b6` |
| `internal/tools/fs.go:22-27` | 那句不完整清单的注释文字（§5） | ⚠ **落进了别人的提交 `319fe230`**，见 §7 那一格 |

⚠ 两格共用一枚新文件＝**派单 §4 点名的文件名是 `fs_edit_ac34_test.go`（两格一枚）**，与"每裁完一格 commit 一次"在字面上冲突 ⇒ 我按**文件名单**执行（不新造名字、不挤进 r1 的文件），两格合在 `64b522b6` 一枚提交里，第三格（注释）单独一枚。这一处偏离登记在此，不改判据。

## 3. 第①格＝AC#3：换行与 BOM，正反两向

### 3a 正向两枚

**（一）带 BOM 的文件改完 BOM 还在**＝`TestFSEditKeepsABOMItWasNotAskedToTouch`

- 夹具自己先验：`original` 前三字节必须等于 `EF BB BF`（用例先 `t.Fatalf` 夹具，免得拿一枚假 BOM 文件证一件真事）。
- 断言**在字节层**：`string([]byte(after)[:3]) == "\xef\xbb\xbf"`，**不是**字符串前缀比较（后者能被"重新编码成同一串文字"糊过去）；再加 `strings.Count(after, utf8BOM) == 1`（既不许丢，也不许重复/移位）。
- 全文等值断言：只有命中的那几个字节可以变。
- 读数：`31 字节 → 31 字节，前三字节仍为 EF BB BF，BOM 出现 1 次`。

**（二）CRLF 文件里 `old` 也写成 CRLF ⇒ 必须命中并落盘**＝`TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF`

- 读数：`23 字节 → 30 字节，CRLF 对 3→4，裸 CR/裸 LF 均为 0`。插入的那一行按 CRLF 拼写，落盘后仍是**纯 CRLF**（`endings()` 把 `crlf/bareCR/bareLF` 三味分开数，裸 CR 与裸 LF 都必须为 0）。

### 3b 反向两枚

**（三）CRLF 文件里 `old` 写成 LF ⇒ 0 命中**，并回答派单点的那一问："它是正确拒绝还是根本读不到"＝`TestFSEditOnACRLFFileRefusesAnLFSpelledOld`

三件读数钉住"**是正确拒绝**"：

1. 拒绝文案里带着文件长度：`（0 命中，文件 23 字节）`，用例断言 `"文件 23 字节"` 这个子串必须在——那 23 只有把整文件读成功才算得出来（`fs_edit.go:171-176` 的计数就发生在 `content` 上）。
2. 读失败那一枚的文案 `"读取目标失败"` **必须不出现**（用例反向断言）。两枚错误形状在码里是不同分支（`:134` vs `:174`），这里把"是哪一枚"钉住。
3. **同一批字节、紧接的下一次调用**，一枚不含换行的 `old`（`"bravo"`）立即命中并落盘 ⇒ "读不到"这一解释被现场否证。
- 未修码上响不响：这一枚**今天就会响**（派单的断言，本程复现：`IsError=true`、`ErrorClass="tool"`、`AppliedSteps` 空、文件 23 字节一字未变）。它的"响"不是新行为，**新增的是"响得可分辨"**：以前没有任何用例说清它是"读不到"还是"按字面拒"。

**（四）不许把 LF 文件的行尾改掉（票面 `:34` 第三句）**＝`TestFSEditDoesNotRewriteAnLFFilesLineEndings`

- 双向：LF 文件配 CRLF 写法的 `old` 同样 0 命中硬错（不对称不能只跑一个方向）；LF 文件配 LF 写法落盘后 `strings.ContainsRune(after,'\r') == false`，且 `crlf=0 bareCR=0 bareLF=4`。
- 读数：`20 字节 → 27 字节，裸 LF 4 条、CR 0 枚`。

### 3c 派单 §1 那句"`content := string(raw)`＝纯字节匹配 ⇒ LF 被改成 CRLF 今天不可能发生"的复算

成立，而且它有一个**更强的形式**（用例里钉住了）：`fs.edit` 唯一能改的字节是**命中的 `old` 之内**＋**调用方给的 `new`**。所以"文件里没人动的行被改写"这一整类今天结构上不存在；**存在的是命中段本身被换掉**，那正是 §3e 两枚缺陷读数的来路。

### 3d 取证：不实现票面 `:27` 的归一，哪些真写法会被 0 命中拒掉＝`TestFSEditLineEndingForensicsRefusesRealSpellings`

每形一枚子用例，**先量代价再登记**：`strings.Count(lfView(file), lfView(old)) == 1`——即"行尾归一（旧新同一把尺）之后这一枚恰好命中一次"。不满足这一条的形状**不算代价**（那是模型抄错了），用例直接 `t.Fatalf`。跑出来的五形全部满足：

| 行尾形（`describeEndings` 现量） | 文件 | 被拒的 `old` | 归一后命中 | 今天的后果 |
|---|---|---|---|---|
| 纯 CRLF（3 组） | `"alpha\r\nbravo\r\ncharlie\r\n"` 23 字节 | `alpha\nbravo` | 1 | 0 命中硬错，文件一字未变 |
| 纯 CR（3 枚） | `"alpha\rbravo\rcharlie\r"` 20 字节 | `alpha\nbravo` | 1 | 同上 |
| 混合（CRLF=2 LF=1） | `"alpha\r\nbravo\ncharlie\r\n"` 22 字节 | `alpha\nbravo`（跨过 CRLF 那一行） | 1 | 同上 |
| 混合（同上文件） | 同上 | `bravo\r\ncharlie`（跨过 LF 那一行） | 1 | 同上 |
| 纯 LF（3 条） | `"alpha\nbravo\ncharlie\n"` 20 字节 | `alpha\r\nbravo` | 1 | 同上 |

⇒ **代价的形状**：5 枚"其实是对的"调用被响亮拒掉，每次多一个来回，**数据损失 0**，且 `assertUntouched` 逐枚验过。
⇒ **不做的坏处有界、做了的坏处无界**：归一是把"响亮失败"换成"我们替模型判定它想改哪一段"，而本票存在的全部理由是消灭后者。**本程不下"永远不做"的结论**，只把两头的读数摆出来——要不要做、以及做了之后"命中段是怎么被放宽的"要不要在台账里点名，**这一裁归编排者**。

### 3e 取证（另一侧，**这两枚今天不拒**）＝`TestFSEditSilentShapeChangesLandToday`

`old` 是字面命中、`new` 是调用方给的，所以归一那层**管不到**它们：

1. **含 BOM 的 `old` 从偏移 0 起 ⇒ BOM 被当普通字节换掉**。读数：`14 字节带 BOM 的文件 → 11 字节、前三字节 41 4c 50（"ALP"）`，答复 `IsError=false`、文案里没有任何编码提示。⇒ 票面 `:34` 的"带 BOM 的改完 BOM 还在"今天**只在 `old` 不含 BOM 时成立**（§3a 那枚）。这一枚断言钉的是**今天的行为**，注释里写明"是缺陷读数、不是批准"；将来若实现剥 BOM／禁命中段跨越 BOM，这一枚会红并强迫改文字。
2. **LF 写法的 `new` 插进纯 CRLF 文件 ⇒ 落盘成混合行尾**。读数：`3 组 CRLF + 1 条裸 LF`，答复里没有任何行尾提示（探针用多词短语 `"行尾"`/`"mixed endings"`/`"endings"`/`"不一致"`——单词 `"mixed"` 会被本用例自己的临时目录名撞上，**首跑就是这样假红过一次**，已改成多词并把这条坑写进码注释）。⇒ `:34` 的"CRLF 文件改完还是 CRLF"今天只在 `new` 也按 CRLF 拼写时成立。
3. 顺带现量：`fs.edit` 的成功答复**已经带前后字节数与行数**（`23 字节 / 3 行 → 30 字节 / 4 行`），这是 `fs.write` 的答复没有的那个数（r1 的 AC#1 量的是"没有可比的数"）。它不算行尾信号，但**算一次可核对的差值**。

⚠ 这两枚要请编排者拍的问（派单没覆盖，我按"未定义即停"不自行假设也不实现）：**是否给 `fs.edit` 加一枚"命中段跨 BOM 即拒"／"`new` 的行尾形与文件主形不一致即拒或即报"**。本程只把它们钉成读数。

## 4. 第②格＝AC#4：原子性反面实验（〔本轮现跑〕）

缝＝`TestAtomicWriteKillsMidWrite` 同一枚 `Hooks.Kill`（派单断言"可复用"＝**成立**）。装法是给 r1 已有的 `fsEditBridge` 加一个 knob（`fsEditBridgeDeps(t, root, tweak)`，`tweak` 传 `nil` 时与原函数一字不差），**没有新造第二套脚手架、没有新的 gate、没有新的 resolver**。杀点 `write:16`、`WriteChunk: 8`。

⚠ 现量补一句派单里没有的形状：`sideEffect.writeAll` 的边界名是 `write:<累计+本块>`，检查发生在**落笔之前**（`fs_write.go:141-145`）⇒ 杀在 `write:16` 时盘上是 **8 字节**，不是 16。两处读数都按 8 写死并写明理由。

### (a) 目标文件字节＝改前原文

`TestFSEditKilledMidWriteLeavesTheTargetByteIdentical`，批是**两枚都已通过全部检查**的编辑（定位／计数／重叠／上限全过了，只剩写）：

```
改前 26 字节 → 杀点 write:16（暂存件里此刻 8 字节）→ 目标仍是 26 字节、与原文逐字节全等
out.IsError=true，AppliedSteps 非空
```

用例还钉了"杀点真的到过"：`AtStep` 记下的边界序列第一个是 `create-temp`、且必含 `write:16`（否则这枚用例什么都没证，直接 `t.Fatalf`）。

### (b) 残件：今天实际是哪一种形状

派单允许"不留残件"或"残件被点名"两种，**本程量到的同时是两种**：

- 目录里 `.wisp-tmp-*` **0 枚**（`noStagingFilesLeft`，r1 已有的 helper）；
- 且台账逐字点名被删掉的那一枚：

```
["在 …\001 创建临时文件 .wisp-tmp-37d9cd89-67612-3083192143"
 "在步骤「write:16」前停止：模拟进程在此刻被杀"
 "向临时文件写入 8 字节"
 "写入未完成：已停止于 write:16：模拟进程在此刻被杀"
 "删除临时文件 .wisp-tmp-37d9cd89-67612-3083192143（目标从头到尾未被改动）"]
```

用例断言"创建临时文件"与"删除临时文件…目标从头到尾未被改动"两行**都必须在**，缺一即红。

### (c) 承重那一问：摘掉 temp＋rename 那一味，是否存在一发损坏从此看不见？

**答：存在，而且这里不是"我看了代码觉得安全"，是一发跑出来的读数。**
`TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget`＝**反面实验**，两步：

1. **对照组（真码，无 hook）**：先让真 `fs.edit` 把同一批编辑落一次，把"要落的那批字节"读回来当 `intended` ⇒ 实验里喂给直写通路的内容**由生产码自己算出**，本程一个字都没替它写。
2. **实验组（摘掉那一味）**：把目标文件重置回 26 字节原文，然后用**生产自己的** `sideEffect.writeAll`＋`FSDeps{WriteChunk:8, Hooks:{Kill: write:16}}` 把 `intended` **直写目标**（`os.OpenFile(... O_WRONLY|O_CREATE|O_TRUNC)`＝上游 `edit.ts:127` 那形"整文件交出去写"）。**唯一变量是缺 temp＋rename。**

```
AC#4(c) 读数：目标从 26 字节变成 8 字节（"ALPHA-LO"）——既不是原文也不是要落的内容；
事后 fs.read：IsError=false ErrorClass="" Truncated=false，Text 恰好等于那 8 字节损坏内容
```

⇒ "看不见"有两层，都量到了：
- **写这一侧没有答案**：损坏停在半路时目标已经是新内容的前缀（`O_TRUNC` 在第一个字节落地之前就毁了原文），而 `discard()`／"目标从头到尾未被改动"那一句**在这个形状里根本不存在**——没有暂存件可删，也就没有那条台账；真实进程死亡（§1 第 5 条）连 `werr` 都不会有人读到。
- **读这一侧没有信号**：下一位instrument（`fs.read`）把 8 字节照读照回、三个标志位全干净。⇒ 这与 r1 AC#1 量到的"少 6 行 222 字节无人报错"是**同一族的无信号**，只是从"内容少了"换成"文件被截断"。装回 temp＋rename（§4a）则同一道杀点下目标一字未动。

⇒ **红线 `:28` 由这两枚用例一起守住**：`fs.edit` 的写只能是 `stageAndRename`，且"摘掉它就是看不见的损坏"现在有一发可重跑的反面读数，不是一句注释。

## 5. 第③件＝裁给本程的那句注释（**只改文字**）

`internal/tools/fs.go:22-25` → 六行，diff：

```diff
-// The write half (fs.write / fs.trash / fs.move, plus the delete_enabled-gated
-// fs.delete) is in fs_write.go and shares this file's FSDeps: same C26
-// resolver, same caps, one registration entry point (BuiltinFSEntries). The
-// [fs] allowed_dirs first-use ask flow is still open (see doc.go).
+// The write half (fs.write / fs.edit / fs.trash / fs.move, plus the
+// delete_enabled-gated fs.delete) is in fs_write.go and, for fs.edit (ticket
+// 162), in fs_edit.go; all of it shares this file's FSDeps: same C26 resolver,
+// same caps, one registration entry point (BuiltinFSEntries) and the same D31
+// staged writer (stageAndRename). The [fs] allowed_dirs first-use ask flow is
+// still open (see doc.go).
```

- **函数行为零字节**：`git show 319fe230 -- internal/tools/fs.go` 的差量只有上面这六行注释（该提交本身是别人的前端提交，我那枚"注释单独提交"因 §7 那起共享索引事故没有成立）；`go test ./internal/tools/` 与 `go vet` 在改前改后同向（§6）。
- 现量补一句派单里没有的：`fs.edit` 的**实现不住在 `fs_write.go`**（它在 `fs_edit.go`），所以原句不只漏了名字、还**给错了文件**——清单漏一枚是"不完全"，加上地点错是"会带人走错门"。这就是为什么新句里点名了两枚文件。
- **删掉这句更正，哪一发会重新误导**：`fs.go` 是 `FSDeps` 的所在文件，"所有写工具是不是都共用同一枚暂存写"这个问题只会从这里查起（票面 `:28` 的红线与 r1 那句"复用 `fs_write.go` 的 `stageAndRename`，没新造第二条写路径"都挂在这条不变量上）。清单漏掉写族里**唯一按字面定位**的那枚工具，下一个写形状的改动（`fs.append`／批量写）就会照这份"看起来完整"的清单往 `fs_write.go` 里再加一条通路，而 `fs.edit` 落在名单外 ⇒ **第二条写路径从此并行存在且没有一处文档会说它违规**；同时读 `fs.go` 的人仍会以为写族只有三枚＋一枚门控，去 `fs_write.go` 找 `fs.edit` 而找不到（文件给错了）。现量核对：`fs_write.go:15-21` 的表头、`fs.go:318-322` 的 `BuiltinFSEntries` 文档、`doc.go:13-14`、`fs_test.go:181/195/202/272` 的名册**都已含 `fs.edit`**（r1 落的）⇒ 本程之后全仓再没有一处把写族数少了。

## 6. 门禁（两格全部落盘**之后**跑的，〔本轮现跑〕）

逐尺读数的入库件＝`.scratch/wisp/probes/162/final-gates-r2.txt`（含"未跑清单"）。

| 尺 | 命令 | 读数 |
|---|---|---|
| 本包测试 | `go test ./internal/tools/ -count=1 -v` | rc=**0**；`ok … 19.524s`；顶层 **run/pass=94／fail=0／skip=0**、`panic` 行数 **0**（原始输出已入库 `probes/162/gates-tools-verbose-r2.txt`，AC#3/AC#4 读数抽取件＝`probes/162/ac34-readings.txt`） |
| **名册差集** | `go test -list '.*'` 前后各一份 → `comm` | 进场基线 **86 枚**（派单那句"r1 基线 86"＝**现量对上**，进场时先落盘 `probes/162/roster-baseline-before-r2.txt`）→ 现在 **94 枚**；**"基线有而现在无"＝0 枚**（件：`probes/162/roster-missing-in-r2.txt`，空）；新增 8 枚逐枚点名见下表 |
| go vet | `go vet ./internal/tools/` | 无输出，rc=0 |
| d22scan 自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | `PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`，rc=0（＝基线同数） |
| 全仓 d22scan | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`，rc=**0**；射程计数 bans #1-5 `internal/=206 cmd/=23`、#6 `frontend/=85`、**#7 `internal/tools/=19`（与 r1 入库后同数＝本程新文件是 `_test.go`，不在生产码射程内）**、#8 `internal/=418`（r1 是 417，＋1＝本程那枚测试文件） |
| gofumpt | `gofumpt --version` 现跑＝`v0.12.0 (go1.27.1)`；甲形 `git ls-files "*.go" \| xargs gofumpt -l`（**545 枚已跟踪 .go**）＝**0 行**；乙形（已跟踪＋未跟踪）＝**唯一一行** `.scratch\wisp\probes\161\r5\negative-control\bad-sample.go`＝别的票故意留的坏样本，可归因、未删 | 本程三枚文件都不在表内＝干净 |

新增 8 枚（逐枚点名，全在 `internal/tools/fs_edit_ac34_test.go`）：

```
TestFSEditKeepsABOMItWasNotAskedToTouch              AC#3 正向（BOM）
TestFSEditOnACRLFFileRefusesAnLFSpelledOld           AC#3 反向（0 命中＝正确拒绝）
TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF AC#3 正向（CRLF 命中且不串行尾）
TestFSEditDoesNotRewriteAnLFFilesLineEndings         AC#3 第三句（LF 文件行尾不动，双向）
TestFSEditLineEndingForensicsRefusesRealSpellings    AC#3 取证（5 形代价，子用例 5 枚）
TestFSEditSilentShapeChangesLandToday                AC#3 取证（2 枚静默形状，子用例 2 枚）
TestFSEditKilledMidWriteLeavesTheTargetByteIdentical AC#4 (a)(b)
TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget AC#4 (c) 反面实验
```

`=== RUN` 总数 141（含子用例），`^--- PASS` 94／`^--- FAIL` 0／`^--- SKIP` 0：**没有新增 `t.Skip`**，也没有把任何断言改松（r1 的 5 枚 `TestFSEdit*` 用例文本一字未动，本程只改了它们依赖的 builder 的一枚函数签名，见 §2）。

## 7. 有没有动过禁改名单里的文件

**没有动禁改名单**。本程写的字节只在：`internal/tools/{fs_edit_ac34_test.go(新建), fs_edit_test.go(只改 builder), fs.go(只 :22-27 那句注释)}` ＋ `docs/evidence/s1/162-fs-edit-r2.md`（本件）＋ `.scratch/wisp/probes/162/**`。
`internal/risk/**`／`internal/config/**`／`internal/agent/**`／`internal/panel/**`／`cmd/**`／`docs/PLAN.md`／`docs/specs/**`／`tools/d22scan/**`／`allowlist.txt`／`thresholds.go`／golden／`.github/workflows/**`／`frontend/**`／`design/**`／`docs/reports/**`／票面 ＝ **零字节**（`git status --porcelain -- internal/tools/` 与三枚提交的 `--name-only` 自证）。**没有一格需要落在禁改名单里 ⇒ 未触发停手上报那一支。**

⚠ **一起必须登记的共享索引事故（〔本轮现跑〕，不是我改的，是别人把我 staged 的字节带走了）**：

- 我把 `internal/tools/fs.go` `git add` 之后、`git commit` 之前，同一台机器上的另一个执行位连着提交了两枚（`2a14e40a`、`319fe230`，都是 `fix(frontend)`），其中 **`319fe230` 的 `--name-only` 含 `internal/tools/fs.go`** ⇒ 本程第③格那六行注释文字**落进了别人的前端提交里**，我的 `git commit -m … -- internal/tools/fs.go` 因此报"nothing to commit"（rc=1，未产生提交）。
- 我**没有**做任何补救动作：`reset`／`stash`／`checkout .`／`--amend` 全在禁列，改写已推送/别人的历史也不在我的权限里。字节是**我写的、内容正确、已在 HEAD 里**，只是**归属挂错了提交**。
- 之后我改用 `git commit … -- <显式 pathspec>`（不 `add`）的形状，并把这一条写给 r3：**这台机器上 `git add` 与 `git commit` 之间的窗口会被别人的提交吃掉**，注释类改动尤其容易被顺手带走。
- 另：`git diff --cached --name-only` 在提交前曾出现别人的三条 `frontend/src/**` 路径＝派单 §7 说的"共享索引正常噪声"，本程两枚提交的 `--name-only` 复核**都没有别人的路径**（`64b522b6`＝我的两枚文件；证据件提交＝我的两枚路径）。

⚠ 派单那句"分支 `dev`"与现场 `dsh/feat/frontend-p0` 不符（§0）。本程未切分支、未 push。

## 8. 被拒／没成功的调用（取数前还是取数后）

全部发生在**取数之后**，且权限系统**零次**拒绝本程请求（无弹窗、无改道）：

1. `git commit -m … -- internal/tools/fs.go` ⇒ rc=1 `no changes added to commit`＝§7 那起事故（字节已被 `319fe230` 带走）。取数之后。
2. 三枚新用例**首跑红**，三处都是**我自己写错的预期**，没有一处是改判据迁就码：
   - `TestFSEditLineEndingForensicsRefusesRealSpellings` 两行 ⇒ 我把"归一后可命中"算成"只归一文件、不归一 `old`"（票面 `:27` 明写"旧文本与新文本同一把尺"）⇒ 测量尺改成两向都 `lfView`，代价行数从 3/5 变 5/5；
   - `…lands_mixed_endings` ⇒ 我把落盘后的 CRLF 组数算成 2（真数是 3 组 CRLF＋1 条裸 LF）；
   - `TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget` ⇒ 我按"杀点写 16"预期，现量是 8（边界名含"下一块"、检查在落笔前，§4 已把这条写进码注释）。
   取数之后。**注意这三枚的修法方向全部是"把我的预期改成读数"，没有一次动 `fs_edit.go`、没有一次动 r1 的任何断言。**
3. `grep -c '^Test'`／`gofumpt -l` 之类的"命中 0 即 rc=1"坑：本程一律用 `;` 串命令、不拿 `&&` 串，所以没被吃掉一次门禁。派单 §6 点过的那枚坑，实测仍成立。
4. `staticcheck` 本程**未跑**（派单 §6：本机它读不了 go1.27 的 export data，会产"干净的绿"）。登记为没测，不是测过。

## 9. 有没有跑过删除命令

**没有。** 全程零 `rm`／`rmdir`／`del`／`git clean`／零 `git stash`／零 `--amend`／零 `reset`／零 `rebase`／零 `checkout .`。
临时件只建不删：`.scratch/wisp/probes/162/**` 十枚新件（`roster-baseline-before-r2.txt`、`roster-after-r2.txt`、`_names-before.txt`、`_names-after.txt`、`roster-missing-in-r2.txt`(空件)、`roster-added-in-r2.txt`、`gofumpt-tracked.txt`(空件)、`gofumpt-tracked-and-untracked.txt`、`gates-tools-verbose-r2.txt`、`ac34-readings.txt`）＋门禁汇总件 `final-gates-r2.txt` 全部保留。
用例里也没有新增删除原语：`fs.edit` 的落盘动作仍只有 `stageAndRename`（temp＋rename），§4c 的反面实验写的是**它自己用例目录里的文件**，不是产品通路。

## 10. 伪授权两栏计数（各带出处）

- **真通知回显数＝0**：本程没有收到过一次系统/审批通知的回显文本。
- **判为注入数＝1**（逐枚点名，不猜）：
  1. 出处＝工具名 `Write`，命令前 40 字＝`PostToolUse hook additional context (tool: W`，内容是
     `[MEDIUM] Security note: file APIs and path joins are sensitive…`，落点＝`internal/tools/fs_edit_ac34_test.go` 建件之后。
     它出现在**取数之前**（第一枚测试文件刚落盘，门禁尚未跑），是**通用告警模板**、不是授权。
     本程**没有据此改动任何判据**：路径决策仍然只在 `d.canonical`（C26）之后做；用例仍走桥；
     §3/§4 的断言一条没少。它建议的"clean the path"也**没有被照抄**——那正是 `AGENTS §1.2` 禁止在
     `risk.PathResolver` 之外做的事，测试里的路径全部走 `NewPathCanonicalizer`。
  - 其余零：本程其余每一次 `Edit`／`Bash`／`Read`／`Glob` 调用都没有附带过任何类似文字（后置告警只有上面那一次）；
    没有任何"已解锁／请 revert／直接给结论／别测了"形状的文字出现过。

## 11. 凭据值零抄录

本程输出里没有 API 密钥、DPAPI blob、token 或任何配置里的 `*key*` 值。读数里出现的字符串只有：临时目录路径、`.wisp-tmp-*` 台账名、`fs.read`/`fs.edit` 的中文文案、`"ALPHA-LO"` 这类测试自有字节。

## 12. 档位（四档）汇总

- **〔本轮现跑〕**：step-0 四件、进场名册 86 枚、`go test -count=1 -v` 94/0/0 与名册两向差集、`go vet`、两把 d22scan、gofumpt 甲乙两形、§3 全部 8 枚用例的读数、§4 三问的读数、§7 的提交事故（`git log`/`git show --name-only` 现量）。
- **〔现读码，未跑用例〕**：`fs_edit.go:136/171/182/205-215` 的纯字节形状、`sideEffect.writeAll` 边界名的"落笔前检查"算术（这条被 §4 的用例反过来钉住了，故两处都算现跑）、`stageAndRename` 其余失败分支未复算。
- **〔盘上有件，我未复算〕**：A18 真实进程死亡那半（`bridge_a18_kill_windows_test.go`）、r1 的 AC#1 读数件（`probes/162/ac1-reading*.txt`）、票面 `PLAN.md:1176` 原文（行号会腐坏，本票 `:41` 已提醒；本程按 r1 的转述＋`fs_write.go:267-270` 的码注释读它，未去 `PLAN.md` 现量）。
- **〔仅自述〕**：§3d 末尾那句"不做的坏处有界、做了的坏处无界"是我的判断，不是读数；§5 里"下一个写形状的改动会照这份清单加第二条通路"是对将来的推测。**两句话都不构成任何一格已完成的凭据。**

## 13. next=（给 r3 的 AC#5／AC#6 各一句起点）

- **AC#5（风险档路由／越界路径）**：起点＝`internal/tools/fs_edit_test.go:30` 那枚 builder 现在多了 knob（`fsEditBridgeDeps(t, root, tweak)`），越界用例只要换 `NewPathCanonicalizer([]string{root}, nil)` 的 root 集就行，**不要再加第三套脚手架**；判据是"allowed_dirs 之外的路径必须被 C26/`PathResolver` 拒（`fs_edit.go:113` 的 `t.d.canonical` 那一支）"，而声明档那一半 r1 已写成 `Declared: risk.L2`（`fs_edit.go:260-269`，不挂 `Facts`）。⚠ 现量提醒：**`fs.edit` 的 Execute 里根本没有 `RiskAssessor` 调用**，判级全在桥（`internal/tools/bridge.go`）那一层 ⇒ 若"新名字要被判成 L2"必须动 `internal/risk/**`，**停手报回**，owner 给 162 的范围里没有那块地（本程实测：§3/§4 八枚用例**没有**需要动它一格）。
- **AC#6（契约轴复算）**：起点＝**自己现量，别照抄票里的行号**（票面 `:41` 自己声明过行号会腐坏）。跑 `git diff --numstat <契约锚>..HEAD -- docs/PLAN.md docs/specs/** internal/risk/** internal/config/**` 应**为空**；并顺手把本程那起"staged 字节被别人的提交带走"（§7）算进"契约轴上看得到什么"：`docs/PLAN.md:2536`／`SPEC-07:43` 那两行由编排者落，实现程（含本程）任何一枚提交里都不该出现它们。
