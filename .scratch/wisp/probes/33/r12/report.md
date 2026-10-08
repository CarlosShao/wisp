# 33-r12 交付件（票 33 的 `.md` 旧方法名归位＝文本腿）

- 工单＝票 33 `.scratch/wisp/issues/33-panel-host-c27.md`｜腿号＝`33-r12`（引用本腿任何一条活请带这三件：票 33＋`33-r12`＋下面那格的内容锚）
- 起手 HEAD 现量＝`bb871191`（`git log --oneline -1` 真身；派单转述的 `23ade0da` 不是我这把的锚，未采）
- 起手 `probes/33/`＝21 枚（`a1 a2 a3 h1 n1 p1 r1..r9 r8b r10 r11 v2 v4 winc1`），**`r12` 目录起手不存在**＝与"应为空"相符，非"回非空"
- 第 1 笔＝`cb90bdf1`（名册＋取证件）｜第 2 笔＝改字（票面 1 行＋本件与收尾尺）
- 零代码、零测试、零 build、零 vet、零 AC 框、零 push；写面只碰 `.md`

---

## ① 全量名册

尺＝`grep -rn 'firstRoundTripLocked\|serveNotBuiltNoticeLocked' --include='*.md' .`　rc=0
读数＝**102 行／27 枚文件／39,234 字节**，最长行 **1,344 字符**（`awk` 现量）⇒ 逐枚**整行**读（`Read` 全 102 行，没有 `cut -c1-150`）。
原始读数＝`grep-raw.txt`（102 行／39,234 字节）、逐枚文件计数＝`grep-byfile.txt`（27 行）、票面骨架＝`ticket-headings.txt`、证据件上下文＝`evidence-ctx.txt`。

**逐行名册（文件／行／那句在说什么／现在时 or 历史／处置）＝`roster.md`（106 行／15,723 字节）**，本件不重复那 102 行，只给汇总：

| 分类 | 行数 | 说明 |
|---|---|---|
| 写面内（票面 5 ＋ `docs/evidence/**` 6） | 11 | 逐行判语见 `roster.md` §1／§2 |
| 台账 `docs/reports/pending-and-issues.md` | 6 | 禁区：既有行一字不许动（10204／10292／12883／13737／13964／13991），且逐条本身是历史 |
| 别人格 `probes/**`（21 枚文件） | 85 | 不在本程具名写面 ⇒ 未动，逐枚仍给判（`roster.md` §4） |
| —— 按性质 —— | | |
| 判**历史叙述** | **92** | 带日期／带"当时"／带 `file:LINE` 的现量读数／引用源码原文或某枚 commit 的内容／引用他人原话（裁定、票面抄本）／仪器尺读数与突变报句 |
| 判**现在时描述代码** | **10** | 其中：改 1、不可纯换字 3、写面外 6 |

10 行的逐枚下落：
1. `issues/33-panel-host-c27.md:358` —— **改**（唯一一枚）
2. `issues/33-panel-host-c27.md:385` —— 不可纯换字（见 §⑤顶回）
3. `probes/33/v4/verdict.md:57` —— 同上那一形（"…`Locked` 后缀意味着调用点应持锁。`firstRoundTripLocked`（`:751`）…自己取锁"）
4. `probes/33/v4/verdict.md:58` —— 同上（`serveNotBuiltNoticeLocked`（`:460`）…）
5. `probes/33/a3/verdict.md:128` —— 写面外（"`firstRoundTripLocked` 里就有一枚 5 秒 `deadline` 循环＋`pnlPumpOnce()`"）
6. `probes/33/a3/00-anchor.md:24` —— 写面外（"涉及载体：…的 `firstRoundTripLocked` / `serveNotBuiltNoticeLocked`"）
7. `probes/33/v4/verdict.md:50` —— 写面外（"`firstRoundTripLocked` 的四条 `-1` 出口"）
8. `probes/33/r10/impl.md:154` —— 写面外（"探测页…是运行时从 `firstRoundTripLocked` 那里取来的"）
9. `probes/33/p1/probe.md:169` —— 写面外（"那正是今天 `firstRoundTripLocked` 会绿的那一枚"）
10. `probes/35/r1/impl.md:34` —— 写面外（"`installPanelTransport` 在 `bringUp` 里跑于 `firstRoundTripLocked`/`serveEntry` 之前"）

⇒ **本腿动字＝1 枚文件／1 行／1 处。** 写面外那 6 枚具名列出，归编排者处置（派单只具名放开 `issues/**` 与 `docs/evidence/**`，本腿不扩写面）。

## ② 改动逐枚 diff 摘要＋票面框尺前后读数

改动（唯一一枚，行内换字，未加行未删行）：

```
--- .scratch/wisp/issues/33-panel-host-c27.md:358
- (a) 最后一份文档**逐字节等于**探测步骤自己写的那份吗？——探测页由 `firstRoundTripLocked` 对着一次性接收器**当场跑出来再取**，⛔ 不是抄进台件的字符串；
+ (a) 最后一份文档**逐字节等于**探测步骤自己写的那份吗？——探测页由 `firstRoundTrip`（这枚方法的旧名＝`firstRoundTripLocked`，`33-r11` 的 ⓐ 去掉了后缀）对着一次性接收器**当场跑出来再取**，⛔ 不是抄进台件的字符串；
```

形制说明：这枚句子在"33-r10 进度追加（2026-10-08 10:2x）"节里，但它断的是**取法**（运行时向一次性接收器取，不是抄进台件的字符串），既无日期状语也无行号锚——今天那枚名字已不存在，说的是假话 ⇒ 属"现在时描述代码"。为不让换名把"当时那枚名字"抹平，本腿采**新旧并列**写法（旧名仍在句内可核），而不是纯替换。

框尺（尺＝`grep -cE '^[[:space:]]*- \[ \]'`，带 `[[:space:]]*`，rc 单测不挂管道）：

| 文件 | 改前 | 改后 | 同数 |
|---|---|---|---|
| `.scratch/wisp/issues/33-panel-host-c27.md` | 13 | **13** | 是 |
| `docs/evidence/s1/33-panel-host-c27-r4.md` | 0 | （未动）0 | 是 |
| `docs/evidence/s1/33-panel-host-c27-r5.md` | 0 | （未动）0 | 是 |
| `docs/evidence/s1/33-panel-host-c27-v1.md` | 0 | （未动）0 | 是 |
| `docs/evidence/s1/33-panel-host-c27-v2.md` | 0 | （未动）0 | 是 |

已勾框尺 `grep -cE '^[[:space:]]*- \[x\]'`＝**1 改前／1 改后**（本腿一枚框没翻、一格判据没改、没加 `-done`）。
行数尺：`wc -l` 票面改前 **394**／改后 **394**（＝零增删行）。读数件＝`rulers-before.txt`、`rulers-after.txt`。

## ③ "我判定为历史所以一枚没动"的名册（具名，可核）

写面内 10 枚历史的（一句没动）：

- `issues/33-panel-host-c27.md:39`——AC#13 框正文＝**判据文字**（硬约束禁改），且那句明写"我 11:1x 亲自复认过行号"、同格 10-08 11:1x 的收窄注已自证"那四个行号是 10-01 的现读、已漂"
- `issues/33-panel-host-c27.md:321`——裁定 P2 原文，与台账 `pending-and-issues.md:10292` 逐字同源（台账不可动 ⇒ 单改票面会造出票面↔台账的假分歧）
- `issues/33-panel-host-c27.md:348`——同节上一行就是"现量（尺＝…，本文件**当时** 861 行）"，行号 `:425`/`:757-758` 今天已漂到 `:448`
- `issues/33-panel-host-c27.md:385`——见 §⑤（不可纯换字）
- `docs/evidence/s1/33-panel-host-c27-r4.md:251`——33-r4 本程"现量"（`:221`/`:227`/`:374-375`）
- `docs/evidence/s1/33-panel-host-c27-r5.md:56`——表行"改前（行号现取于 `7a02b121`）"
- `docs/evidence/s1/33-panel-host-c27-r5.md:57`——表行"改后"，记 33-r5 那一笔造了什么
- `docs/evidence/s1/33-panel-host-c27-v1.md:83`——v1 腿顺路核 AGENTS §1.2 的原读数（附当时注释行号 `:168-173`）
- `docs/evidence/s1/33-panel-host-c27-v1.md:89`——§A#28"调用点全仓唯一＝`:227`（grep 现量）"
- `docs/evidence/s1/33-panel-host-c27-v2.md:34`——AC#13 判语的证据半（起手锚 `3216ba6d`、10-01 16:47、`:308`/`:310`）

禁区不动（另具名）：台账 6 行（10204／10292／12883／13737／13964／13991，**含 `A704`/`A705`/`A708`–`A713` 那些旧名，连追加都没做**）；别人格 `probes/**` 21 枚 85 行；`cmd/wisp/panel_locked_naming_33r11_windows_test.go:360` 的正控夹具 `roundTripHost.firstRoundTripLocked`（那枚名字是 ⓒ 那枚会响的钉的输入，改了钉就哑——产码是禁区，本腿没碰也不打算碰）。

## ④ 收尾自证：`git diff --numstat bb871191..HEAD -- ':(exclude).scratch/wisp/probes/33/r12'`

```
1	1	.scratch/wisp/issues/33-panel-host-c27.md
```

逐枚：只有这一枚文件被本腿改过（第 1 笔只写 `probes/33/r12/**`，被 exclude 剔掉）。
**删除行＝1，不是 0——这不是违令，是"改字"这件事在 numstat 上的固有形状**：git 的 numstat 以**行**为粒度，任何行内换字都记成 `1/1`（一行旧文本消失、一行新文本出现）。本腿真正满足的判据是**行数不变**（`wc -l` 394→394，`git diff --stat` 里那枚文件的 `+1 -1` 就是同一行），也就是仓里那条实测结论"脆的是增删行，不是符号名"。若要 `删除行=0` 的字面形状，只能**追加新行**（加一行注记），而那正是本程明令禁止的"只加字不改行"之外的动作，且会挪动票面下游靠行号认东西的引用（`:385` 段与台账都按票面行号引）。这一条按第 ⑤ 顶回。

预测尺（提交前 `git diff --numstat HEAD -- <票面>` 现量）＝`1 1`，与上表一致；取证件＝`numstat-beforecommit.txt`。

## ⑤ 我没量到的／要顶回与告知的

1. **顶回派单（判据形状）**：§3 要求的"每枚文件删除行数必须为 0"与"只改字"在数学上不能同时成立（见 §④）。本腿按"行数为零增删"执行并落了两把尺（`wc -l`、numstat 的 `1/1`）。若编排者要的是"票面零改动面"，那这格的正确答案是**由他追加一行带时刻的更正**，不是由我扫字。
2. **顶回派单（示例句其实是历史）**：派单举的"该改"例子——票面那句"当前 `bringUp` 会调 `firstRoundTripLocked`，它自己取锁"——现量落在 `issues/33-panel-host-c27.md:385`，**我没改**。那枚句子的谓词就是"`Locked` 后缀却自己取锁"，换成如实名后句子自相矛盾（"带 `Locked` 后缀"当场变假），比留旧名更坏；同一形还有 `probes/33/v4/verdict.md:57`、`:58` 两枚。这一格的如实更正只能靠追加注记（"ⓐ 已于 `33-r11` 落地，两枚现名 …"），而写票面那段措辞台账 `A713` 里编排者已自留（"这枚票面措辞要我改与①同批"）。同段 `:391` 那句"现存那 2 枚进具名豁免名册直到 ⓐ 落地"不含旧方法名、本腿的尺抓不到，也一并**没碰**，具名告知。
3. **约 20 枚的差**：33-r11 报"约 20 枚 `.md`"，我这把全量尺＝**27 枚文件**（含 33-r11 自己没数到的 `probes/111/c2/logs/*`、`probes/panel-resident/1`、`probes/e2e-panel-1`、台账与 4 枚证据件）。枚数差属正常（那句是"约"），但**具名事实是：写面内只有 11 行、其中只有 1 行该由我改**——绝大多数不是"该扫的旧字"，是别人的读数与裁定原文。
4. 我没量的：整包 `go test`／`-race`／`d22scan` 复跑（产码面禁区、且派单禁跑）；票面 `:39` 那枚 AC 框的翻勾与 `ⓑ`/`ⓒ` 落地性；台账 `A201②` 那条"规格比仪器宽"的缺口（与本格无关）；`probes/**` 里那 6 枚现在时句子要不要改——**要改，但写面不在我这程，需编排者具名放开或他自己落**。
5. 本件与取证件（本腿目录共 10 枚：`roster.md`／`report.md`／`grep-raw.txt`／`grep-byfile.txt`／`ticket-headings.txt`／`evidence-ctx.txt`／`rulers-before.txt`／`rulers-after.txt`／`numstat-beforecommit.txt`／`oldnames-go-files.txt`）全部非 0 字节、无一枚 `.out` 命名（根 `.gitignore` 有全仓 `*.out`）；`git check-ignore` 对本腿件 rc=1（没被忽略）。临时件只建不删。
