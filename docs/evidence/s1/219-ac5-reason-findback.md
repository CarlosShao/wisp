# 票 219 AC#5 最后一寸 —— 被截掉的原文留第二份，标记带可找回指针（写腿交件）

- 腿：`219-ac5`（写腿，本波 `cmd/wisp` 唯一写面）；基线锚点 `d943da7bc`（上一腿：标记在、个数在、第二份缺）。
- 票面：`.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md` AC#5（:75）与「AC#5 ⛔ 不翻 —— 只差一件」（:211-219）；
  派单：`.scratch/wisp/probes/orch/20261115/txt/219-ac5-dispatch.txt`。
- ⛔ 本腿未动：`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43／`frontend/**`／`design/**`／`internal/**`（`internal/observe` 只读引用）。

## 0. 只做那一件

上一腿把超长理由的截断做成了「可见＋计数」（`replyTruncationMark`＋`dropped`），但 `dropped` **只留个数**——被截掉的那 40 字没有第二份，标记里的「原文共 200 字」找不回原文。本腿只补这一件：**被截掉的原文留第二份（本地审计档）＋标记带可核指针（档行号）＋今天就能用的读取入口（`decisions` 动词逐字打回原文）**。三条口径与「静默截断仍不许」全部钉在行为牙上。

## 1. 落点（全部在写面内）

| 件 | 内容 |
|---|---|
| `cmd/wisp/approval_reply.go` | `replyReason` 增 `full`（完整洗后原文，**只在真截断时才填**）；`washReply` 成为洗/截/留三值唯一合成点（`washReplyText` 保两返回值形状不动，既有牙零改）；`refuse()` 在**路由之前**mint 指针并把指针句拼进主文本（模型／审计／回看三读者看到同一段）；`record()` 拆出 `book()` 单一写点＋`recordWithOriginal()`（回看条目携带 `original`）；`decisions()` 加档头说明与逐字「原文全文」行；`archivePointer()` 是口径①②的判定门 |
| `cmd/wisp/approval_reason_archive.go`（新） | 本地审计档 `<dataDir>/audit/approval-reply-reasons.jsonl`（0700 目录／0600 文件，追加式 JSONL）：每行 `{seq,when,corr,action,route,dropped,original_runes,original}`；`seq`＝**每次追加现数行数＋1**（跨进程重启后指针仍可核） |
| `cmd/wisp/approval_219ac5_findback_test.go`（新） | 三枚行为牙（§3） |

## 2. 标记现在长什么样（原文，取自 M1 突变日志里真实控制台 dump）

模型文本／审计 `detail=`／回看 `理由/明细=` 三处同一段（截断时）：

```
<前160字>…(理由过长，已截断 40 字，原文共 200 字)；原文另存：audit/approval-reply-reasons.jsonl 第 1 行（corr="t219-findback-corr"，全文 200 字）
```

- 「全文 200 字」是**可核断言**：档第 1 行的 `original` 逐字持有这 200 字、`original_runes=200`、`dropped=40` 与标记同数。
- **512 上限为何不触发**：`internal/observe/redact.go:37 MaxLoggedString=512` 管**到达 slog sink 的字符串**；第二份写的是本包自持文件（`os` 直写，不经 redactor、不在 `logs/`、不被 diagnostics 打包）⇒ 原文完整落档，「全文 N 字」为真，无假称（`approval_reason_archive.go` 头注释逐字钉住这条纪律）。
- 正常长度（≤160）：`full=""`、无指针句、无档行、无回看副本 —— 口径①。
- 空理由：`washReply("")` 早退（`"","",0`）—— 与上一腿逐字同形，口径③。

## 3. decisions 怎么取回原文（命令＋真实输出）

命令：`wisp run` 的答复提示符下敲 **`decisions`**（help 里的既有动词，票 219 第三去向）。真实输出（M1 突变日志的实弹 dump，`<前160字>` 处为 200 个「长」的前 160 个）：

```
wisp run: 本进程最近 1 条答复（新的在前；落盘的证据是日志里 approval: REPLY 那一族行）：
被截断理由的原文另存档：<dataDir>\audit\approval-reply-reasons.jsonl（只有真截断的答复才有行；标记里「原文另存」说的第 N 行就是这份档的第 N 行）
  2026-11-16T10:12:35+08:00 ANSWERED corr="t219-findback-corr" tool=fs.write route="native/reject" 理由/明细="<前160字>…(理由过长，已截断 40 字，原文共 200 字)；原文另存：audit/approval-reply-reasons.jsonl 第 1 行（corr=\"t219-findback-corr\"，全文 200 字）"
      原文全文="<完整200字>"（这条的理由超过了 160 字上限：明细里是截断后的文本，这一行是被截掉之前的完整原文，两者逐字对得上）
```

- **进程内**：`decisions` 的「原文全文」行逐字打回（牙①的取回断言）。
- **进程重启后**：标记/审计行里的「第 N 行」＝档的第 N 行，`original` 字段即全文；跨进程 seq 靠每次追加现数行数自愈。
- 敏感面边界：原文全文只出现在**档**与**操作员主动敲 `decisions` 的控制台**；⛔ 产品日志（`logs/wisp-*.jsonl`）里只有截后 160 字＋标记＋指针（既有审计行为不变，本腿未把全文放进任何 slog 路径）。

## 4. 三枚行为牙（A1880：全部经 `newReplyHost` 真装配＋真控制台管道＋真档文件）

| 牙 | 名 | 拆什么 | 红句（实弹 M 突变各一发出处） |
|---|---|---|---|
| ① | `Test219AC5TruncatedReasonIsFindableThroughTheMark` | 删取回（`recordWithOriginal` 传空）| `:109` "the decisions read entry does not carry the verbatim 200-rune original: the retrieval was deleted or broken" ＋ `:113` "the review line grew no 原文全文 field: the entry lost its second copy"（M1：FAIL 11.89s） |
| ①(档半) | 同上 | 删档写入（`archive.append`）| `:85` "reading the reason档: …"／逐字档行断言红（M1 同时证：`:78` "the review list never printed the archived original"） |
| ② | `Test219AC5ShortReasonLeavesNoSecondCopy` | 删「只在真截断时写」的判断（`archivePointer` 的 `dropped` 门＋`washReply` 的 `full` 门）| `:195` "a reason within the ceiling grew a档 copy (stat err = \<nil\>): the only-on-truncation judgment was deleted" ＋ `:191` "a within-ceiling reason was marked or pointed as truncated" ＋ `:199` "the review list grew an 原文全文 field…"（M2：FAIL 1.73s） |
| ③ | `Test219AC5EmptyReasonStaysDeterminate` | 改空理由语义（`washReply("")` 造句）| `:252` "the model-visible refusal lost its cause: \"（操作员未替你填理由）\"" ＋ `:259` "the audit line did not record the empty reason as empty"（M3：FAIL 12.17s） |

突变方法：`go test -overlay` 经载具（`GOFLAGS=-overlay=<space-free>/M*.json`，突变体放 `%TEMP%\219ac5mut\`，**盘上 `cmd/wisp` 始终只有本腿三件**，每发后 `git status --porcelain -- cmd/wisp` 复核）。基线（未突变）：六枚 AC#5 牙（新三＋旧三）**全绿**，载具 `ok … 12.00s`（2026-11-16 10:09）。

## 5. 整包读数（A1882）

载具：`D:\work\soft\Git\bin\bash.exe scripts/wisp-cli-staged-tests.sh`（整包一发；staged DLL sha256 全 ok，voice fixtures 齐 → 无跳名）。日志：`.scratch/wisp/probes/219/ac5/fullpkg.log`（8037 行，跑到底，无截断点）。

- **RUN=904 / PASS=660 / FAIL=2 / SKIP=17，1730.00s**
- **终点锚命中**：`TestZZ315A1RootLoopFirstBridgeCall`／`TestZZ315A1SiblingChildFirstBridgeCall`／`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` 三枚全 PASS；**最后一条 `=== RUN` ＝ `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（PASS 0.19s，日志行 7986）** ⇒ 本发无截断点。
- **三枚新牙在截断点之前**（日志行 288／312／336，均 PASS：2.54s／3.66s／2.50s）。
- **FAIL=2 均为既有红/环境红，非本腿引入**：
  1. `TestAC4FocusReturnToPriorWindowGap33r5`（5.16s）——**昨天编排者整包基线（`orch/20261115/46-full-log.txt`，head d943da7bc）同一枚同红**（AC#4 focus 桌面态读数）；
  2. `TestAC246DevLegIgnoresTheTestTaskInjection`（37.87s）——红句自陈「**AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session**」，环境拒绝，非断言破坏。

## 6. 诚实边界

1. `decisions` 回看仍是**进程内、上限 64**（既有口径）：窗口滚出后指针靠**档**继续可核（第 N 行），控制台不再显示——档是持久第二份，回看是「今天就能用」的读取入口，两层各有其职。
2. 被拒（REFUSED）的答复：指针在路由前 mint，档行已落（键 corr+route+when 可查）；但审计 REFUSED 行 `detail=` 是错误文本不带理由，回看条目不带 `original`——原文检索走档行（按 corr grep），不在控制台逐字显示。
3. 无管线托盘面（`answerSurfaceWithoutPipeline`）**未接档**（`archive=nil`）：该面今天没有可截断的操作员输入路径（托盘拒绝是固定短句、不经洗），`archivePointer` 对 nil 档明确说「原文另存未接入——装配缺陷」而不是冒指针；若未来该面接理由框，须同时接档。
4. 单进程单写者假设：档的互斥锁覆盖本进程；跨进程并发追加靠 O_APPEND＋每发现数行数尽力对齐（注释写明，best-effort）。
5. 「原文」＝**洗后完整文本**（控制符→空格、空白折叠后）——与标记「原文共 N 字」同一定义，回看/档/标记三处逐字对得上；洗是形状清洗不损语义（票 216 同族）。
