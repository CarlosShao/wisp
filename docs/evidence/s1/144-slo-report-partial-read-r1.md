# 144 — `wisp slo` 收主体报告：把"还没写完"与"真的坏了"分开 —— 实现程自证 r1

**落点**：票 144（`.scratch/wisp/issues/144-wisp-slo-collectreport-reads-a-still-being-written-file-as-a-corrupt-one-one-empty-read-kills-the-run.md`）
**身份**：实现程（写代码）。本文件是**实现者自证**，不是裁决表 —— 按 `AGENTS.md §0.3` /
`SPEC-12 §4.3` #1，缺口审计与对抗验收必须由**另一个** agent 出，本文件不声称结掉那一格。
**分支 / 起手锚点**：`dev`，起手 `git rev-parse --short HEAD`（现量见 §0.1）。
**时刻**：2026-09-25 19:0x +08。

---

## 0　现量：票面给的行号与读数，我自己重量了一遍

### 0.1　锚点与树形

```
$ git rev-parse --short HEAD          # 起手
40be959
$ git rev-parse --abbrev-ref HEAD
dev
$ wc -l cmd/wisp/slo_windows.go
666 cmd/wisp/slo_windows.go
```

⚠ 起手时刻树里已有**别家程的在飞提交**（我进场时 `git log -1` 是 `598c0c6 ticket(146)`，写这份证据时 HEAD 已推进到
`40be959`）。本文件所有行号都是**在 `40be959` 上现量**的，`slo_windows.go` 由 `git log --oneline -3 -- cmd/wisp/slo_windows.go`
现量最近改动为 `ed18727 / 6f702ae / 7e60d31`（票 133/131），`collectReport` 那一段（`:520-546`）由
`git log --oneline -L 520,546:cmd/wisp/slo_windows.go` 现量＝**只有 1 枚** `00bbb76 slo(66)`。

### 0.2　`collectReport()` 与其邻座（现量行号）

```
$ grep -n "func (s \*sloSubject) collectReport\|func (s \*sloSubject) waitReady\|func fileHas\|func writeSLO\|os.ReadFile(s.outPath)\|json.Unmarshal(data\|subject report: %w\|never races\|in-tree record unavailable" cmd/wisp/slo_windows.go
383:		return 2, fmt.Errorf("in-tree record unavailable (fail-closed, no silent downgrade to the tree basis): %w", err)
503:func (s *sloSubject) waitReady() error {
522:// window, so this never races the observer's last sample.
523:func (s *sloSubject) collectReport() (*sloRun, error) {
526:		if data, err := os.ReadFile(s.outPath); err == nil {
528:			if err := json.Unmarshal(data, &rep); err != nil {
529:				return nil, fmt.Errorf("wisp slo: subject report: %w", err)
587:func fileHas(path, needle string) bool {
628:func writeSLO(run *sloRun, out string) error {
```

`collectReport` 函数体全文（现量 `awk 'NR>=519 && NR<=547'`，逐字节贴）：见 §0.4。

### 0.3　票面 vs 现量：逐格对账

| 票面说 | 我量到 | 判定 |
|---|---|---|
| `collectReport()` 在 `:523-545` | 函数头 `:523`，闭括号 `:546`（`for` 体到 `:545`） | 票面**末行差一枚**（将闭括号写进了循环体）。形状一致，不影响任何一格 |
| 一发即死那一支在 `:528-529` | `:528` 是 `if err := json.Unmarshal(...)`，`:529` 是 `return nil, fmt.Errorf("wisp slo: subject report: %w", err)` | **对上了** |
| `waitReady()` 在 `:503-518`，用 `fileHas(s.readyPath, "ready=1")` 在预算里轮询 | `:503-518` 逐字对上；`fileHas` 调用在 `:506`，`fileHas` 本体在 `:587-593`（读失败即 `false`，不区分"没文件"与"没内容"） | **对上了**：不对称确实成立 |
| `:520-522` 那句 *"never races"* 注释 | 逐字在 `:520-522`，`never races` 落 `:522` | **对上了** |
| `:383` 那一层包住 `:529` | `:383` 是 `runOutOfTree` 里包住 `collectReport()` 错误的那句 `in-tree record unavailable (fail-closed, ...)` | **对上了** |
| 引用的那 7 行代码 | `:526-534` 的**真实体多一段票面没贴的** `if rep.Report == nil { return ... }`（`:531-533`） | **票面引用不完整**（它省略的正好是"能解析出对象但内容不对"那一支 —— 也就是 AC#2 举例的**当场红**腿本来就存在）。不影响判据，影响的是"这一形今天到底缺哪一味"：缺的只有"未写完"那一味 |

⇒ 结论：**票面对形状的描述成立**（`collectReport` 无任何重试、`waitReady` 容忍空读、注释把没验证的时序假设写成了不变式），
只是行号有 1 枚平移、引用少贴了 3 行。**没有一格需要停手报回。**

### 0.4　写入侧（这一形能不能"分得开"的前提，票面没量，我补上）

```
$ awk 'NR>=628 && NR<=638' cmd/wisp/slo_windows.go
```

- `:629` `json.MarshalIndent(run, "", "  ")` —— 先在内存里把**整个**文档产成一块字节；
- `:637` `os.WriteFile(out, data, 0o644)` —— 一次 open(O_TRUNC|CREATE) + 一次写，**没有临时文件、没有 rename**；
- 尾部**没有换行**（`MarshalIndent` 不追加 `\n`）。

⇒ 唯一让"读者看到半个文档"成为可能的机制就是这一枚 `os.WriteFile`：**磁盘上的可见内容必然是最终文档的一个前缀**
（0 字节、1 字节、…、N-1 字节、N 字节）。这条前提是 AC#2 那句"分得开且有凭据"能落地的**全部依据**：
前缀只会缺尾巴、不会改头部，所以**"头部自相矛盾"的字节永远不可能是"还没写完"**，而"尾部缺失"的字节永远可能是。
（`wisp slo` 主体写完报告之后才 `rt.Shutdown` → 退出，见 `:308-326`；写这一次是**唯一一次写**。）

### 0.5　`exited()` 这一路我今天到底能不能靠？（现量）

```
$ grep -n "func (j \*JobScope) StartInJob" -A 12 internal/proc/jobscope_windows.go
$ grep -rn "\.Wait()" cmd/wisp/slo_windows.go
```

读数：`StartInJob` 只做 `cmd.Start()` + `Assign`，**没有任何 goroutine 去 `Wait`**；`slo_windows.go` 里
`cmd.Wait()` 只出现在 `sloSubject.stop()`（`:571`）。而 `exited()` 读的是 `cmd.ProcessState`（`:548-550`），
`ProcessState` 只在 `Wait()` 返回后才非 nil。
⇒ **`collectReport`/`waitReady` 循环里的 `s.exited()` 在真实运行中恒为 false**（既有事实，不是本票新造）。
因此本票**不把"进程已退出"当判别凭据**（那会是一味抓不住的药）；判别只用**字节自身**（§0.4 的前缀论）。
`exited()` 那一支的原样保留、不删不放宽 —— 它是 `stop()` 之前唯一的一道早退，只是它今天不响。
⚠ 这一格登记进 §6"本程没测什么"：**"主体死了没写完"这一形今天仍然只能等预算到点**（30s+3s），
它到点后给的是**点名预算与最后一次字节数**的红，不是 `unexpected end of JSON input` —— 比修前信息量大，但仍不是当场红。

---

## 1　AC#1 钉成能响的用例 —— **结**

落点：`cmd/wisp/slo_report_144_windows_test.go`（新增，348 行，`//go:build windows`）。
落点（现量，改后 `cmd/wisp/slo_windows.go` 840 行）：
三态 `subjectReportState` `:536`（`String()` `:558`）、观测量 `subjectReportRead` `:574`（`summary()` `:589`）、
判别函数 `readSubjectReport` `:623`、生产入口 `collectReport` `:659`、循环 `collectReportWithin` `:682`、
测试用读数接缝字段 `readReportFile` `:140`。
（`git log --format=%h -1 -- cmd/wisp/slo_report_144_windows_test.go` 现量 = `f63f0e3`。）

判据一句一发读数：

| 用例 | 钉的是 AC#1 的哪一支 | 现量读数 |
|---|---|---|
| `TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt` | ⓐ"没写完"这一支的**全部分母** | 真报告 **1978** 字节（夹具由 `json.MarshalIndent` 产的一份完整 `sloRun`，含 2 枚 sample + 2 枚 verdict；用例先验夹具非平凡才继续），`0..1977` 每一枚前缀逐一问判别函数，全 1978 发都是 `unwritten`；整枚是 `complete`。`--- PASS` |
| `TestSLO144ReportsThatContradictThemselvesAreCorruptNow` | ⓑ"内容确实是坏的仍然当场红" | 7 发子格全 `--- PASS`（html 头 / 多余逗号 / 类型不对 / 值后面还有内容 / 两份文档 / 没有 state report / 只有一对花括号），每发都要求错误点名**字节数** |
| `TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport` | ⓐ 的端到端（真文件） | 空文件 → 半截 → 全量，循环带回报告、`Mode="subject-in-tree"`、`Pass=true`。`--- PASS` |
| `TestSLO144LoopGiveUpSentencesOnRealFiles` | ⓐ 的另一形：等不齐时红得**有理由** | 半截 44 字节：错误点名 `within 50ms` + `44 bytes` + `tail had not arrived`，且断言**不得再是** `unexpected end of JSON input`；坏字节 25：点名 `corrupt` + `25 bytes`、且不得含 `within`。`--- PASS` |
| `TestSLO144CorruptReportIsJudgedOnTheFirstRead` | ⓑ 的"当场"凭据（**数读取次数**，不数秒） | `calls == 1`，且脚本读取器一被叫到第 2 发就 `Fatalf`。`--- PASS` |
| `TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes` | ⓐ 的"回预算里重试"凭据 | 缺档 2 发 + 前缀 1 发起持续：`calls > 3`，红句点名 `within 60ms` / `8 bytes`。`--- PASS` |
| `TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected` | ⓐ 的收口：重试真的能收到 | 恰好读满 4 发后拿到 `report.state=Sleeping`。`--- PASS` |

`$ go test -count=1 -v -run 'TestSLO144' ./cmd/wisp/`（PATH 挂 DLL，见 §4）→
`ok github.com/CarlosShao/wisp/cmd/wisp 0.239s`，顶格 7 枚 + 子格 7 枚，`--- FAIL` 0、`--- SKIP` 0。

**这一形造得出可重放用例** ⇒ 不走票面"AC#1 自证作废"那一支。可重放的依据不是运气：判别函数的输入是**一枚字节切片**，
所以"文件存在但字节没齐"这一态是被**全体前缀**直接问出来的（用例 1），循环那一侧用真文件 + 读数计数器复现，
不需要真等一次真实调度竞态。

⚠ 明说这一格**没有**证的东西：**没有在真进程上复现那一发竞态**（`wisp slo` 真跑时读半个文件）。
§3 的外证那一问会给出为什么，以及它落在"本程没测什么"第 1 条。


## 2　AC#2 修法只有一个方向 —— **结**

一句话：**"没写完"回预算里重读，"坏了"当场红，凭据是"这次读到的字节是不是那份文档的一个合法前缀"** ——
不是"重试几次都失败就算坏"（那一支没有凭据，票面禁止，本票也不采用）。

判别（`readSubjectReport`，`cmd/wisp/slo_windows.go:623`）分成三态，两态的边界是**类型化错误**、不是字符串比对：

| 观察到的字节 | 判别 | 循环动作 |
|---|---|---|
| 0 字节 / 只有空白 | `io.EOF` ⇒ `reportUnwritten` | 预算内重读 |
| 停在某个值内部（`{`、`{"mo`、`"tru`、`5.`…） | `io.ErrUnexpectedEOF` ⇒ `reportUnwritten` | 预算内重读 |
| 头部就不可能是这份文档（`<`、多余的 `,`）、类型不匹配 | 其它 `*json.SyntaxError` / `*json.UnmarshalTypeError` ⇒ `reportCorrupt` | **当场红** |
| 值已闭合后面还有内容（尾巴、第二份文档） | `dec.Token()` 不再 `io.EOF` ⇒ `reportCorrupt` | **当场红** |
| 解析干净、完整、但没有 state report（含 `{}`） | `rep.Report == nil` ⇒ `reportCorrupt` | **当场红**（这一腿修前就有，本票没放宽） |
| 完整一份文档 | `reportComplete` | 收 |

**为什么"分得开"这句在这里站得住**（凭据不在"我猜它没写完"，在写入侧的形状）：
`writeSLO` 先 `json.MarshalIndent` 把整份文档产进内存，再**一次** `os.WriteFile(…, 0o644)` 从 offset 0 推下去，
没有临时文件、没有 rename、尾部不追加换行（§0.4 现量）。⇒ 读者在写入期间能看到的**只有那份文档的前缀**。
前缀只会缺尾巴、不会改头部，也不会凭空多出一份第二文档 ——
所以"头部/尾部自相矛盾"这三类**按构造不可能是没写完**，而"字节用尽"这两类**按构造可能是**。
边界因此是双向可证伪的：§1 的用例 1 把**全部 1978 枚前缀**问了一遍（真报告没有一枚前缀被判坏），
用例 2 把 7 发反面形状问了一遍（没有一枚被判成"还没写完"）。

fail-closed 没有动：
- 报告真坏 ⇒ 仍然 `return nil, error`，且 `runOutOfTree` 那句 `in-tree record unavailable (fail-closed, no silent
  downgrade to the tree basis)`（`:383`，本票未改）继续包在最外层；
- 前缀等不齐 ⇒ 预算到点仍然 `return nil, error`，**没有任何一支变成 continue／静默放过／降级到树内基准**；
- 到点那一句现在点名预算与最后一发的读数：
  `subject %d never wrote a complete report within %s (last read: %d bytes read, document still open at offset %d: the tail had not arrived)`
  —— 这正是 AC#1ⓐ 括号里给的另一支（不假装读到过完整报告）。

预算／阈值／节奏三个数一个都没动（`:80-82` 三枚常量原文未改；`collectReport()` 仍传
`subjectReportBudget+subjectGrace` 与 `subjectPollInterval`）。`collectReportWithin` 只是把这两枚数收进参数，
为的是让用例能不开 33 秒地驱动**同一个**循环 —— 这是形状接缝，不是新的宽限期。

## 3　AC#3 变异自证两向（仓外副本 `D:/tmp/wisp-144-mut`，仓内一字未改）

副本怎么来的（现量命令）：

```
$ cd "<repo>" && tar -cf - --exclude=node_modules --exclude=.git \
    cmd internal frontend go.mod go.sum third_party/sherpa-onnx | (cd /d/tmp/wisp-144-mut && tar -xf -)   # 29M
$ cd /d/tmp/wisp-144-mut && PATH=$PWD/third_party/sherpa-onnx:$PATH go test -count=1 -run TestSLO144 ./cmd/wisp/
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.295s     ← 副本先证明"没变异时是绿的"，否则下面的红不算数
```

### 3.1　(i) "主体写的报告确实坏了（非空、语义错）"⇒ 红，且红因点名

用**语义错**那一形（解析干净、文档完整、里面没有 state report —— 不是截断、不是空）：

- 判别层：`TestSLO144ReportsThatContradictThemselvesAreCorruptNow/complete-but-no-state-report`
  ⇒ `corrupt`，错误句子是 `57 bytes parsed as a subject run but carry no state report`（字节数在句里）；
- 循环层（真文件、真 `os.ReadFile`）：`TestSLO144CorruptReportIsJudgedOnTheFirstRead`
  ⇒ `wisp slo: subject report is corrupt (failing closed on the first read, not a partial write): 38 bytes parsed as a subject run but carry no state report`，
  并且 `readReportFile` 计数器 **calls == 1**。

⇒ 坏字节不会被"再给它一次机会"洗白，红句点名它坏在哪。
⚠ 这一形的**进程外**版本（真 `wisp slo` 产出一份坏报告）没做到，见 §6 第 3 条。

### 3.2　(ii) 摘掉 AC#2 新加的那一味 ⇒ AC#1 那批从绿变红

**变异 A**：删掉 `readSubjectReport` 里 `io.EOF`/`io.ErrUnexpectedEOF` ⇒ `reportUnwritten` 那一支，
即"解码报错一律当判决"＝**修前的形状**。`go test -count=1 -v -run TestSLO144 ./cmd/wisp/`，rc=1：

| 用例 | 变异 A 下 | 落点红句（原文摘） |
|---|---|---|
| `…EveryPrefixOfARealReportIsUnwrittenNotCorrupt` | **FAIL** | `slo_report_144_windows_test.go:108: prefix of 0/1978 bytes classified as corrupt, want unwritten (0 bytes read, report corrupt: … offset 0: EOF)` —— 1978 枚前缀逐枚红 |
| `…ReportsThatContradictThemselvesAreCorruptNow` | PASS | ← 该绿的保持绿 |
| `…LoopRetriesAnUnfinishedFileAndReadsTheWholeReport` | **FAIL** | `:200: collectReport on a file that finishes writing: wisp slo: subject report is corrupt (failing closed on the first read, not a partial write): 0 bytes contradict a subject report at offset 0: EOF` |
| `…LoopGiveUpSentencesOnRealFiles` | **FAIL** | `:230: unfinished give-up "…subject report is corrupt…44 bytes contradict…unexpected EOF" does not name "within 50ms"` |
| `…CorruptReportIsJudgedOnTheFirstRead` | PASS | ← 同上 |
| `…UnfinishedReportKeepsPollingThenNamesBudgetAndBytes` | **FAIL** | `:320: the loop took 3 reads; a missing file and an unfinished file must both be re-read inside the budget` |
| `…ReportThatArrivesAfterMissingReadingsIsCollected` | **FAIL** | `:340: collectReportWithin: wisp slo: subject report is corrupt …: 0 bytes contradict a subject report at offset 0: EOF` |

⇒ **摘掉它，AC#1 的 5 枚用例从绿变红**，红的正是"没写完"那几条腿；"当场红"那两条不受影响 ⇒ 这一味不是装饰。

**变异 B**（反方向，AC#2 明禁的那形："解析失败就 `continue` 到预算尽头"）：把循环里
`case reportCorrupt: return nil, …` 吞掉。rc=1，红的正好是另外两枚：

```
--- FAIL: TestSLO144LoopGiveUpSentencesOnRealFiles (20.06s)
    slo_report_144_windows_test.go:250: corrupt bytes were treated as a wait, not a verdict:
      "wisp slo: subject 4244 never wrote a complete report within 20s (last read: 25 bytes read, report corrupt: 25 bytes contradict a subject report at offset 0: invalid character '<' looking for beginning of value)"
--- FAIL: TestSLO144CorruptReportIsJudgedOnTheFirstRead (0.00s)
    slo_report_144_windows_test.go:270: scripted.json: read #2 of a 1-reading script - the loop is
      retrying something it should have judged on the first look
```

⇒ **A 红"没写完"那 5 枚、B 红"当场红"那 2 枚**：两向非对称都量过，不是单向自证，也没有一枚恒绿。

### 3.3　另问那句：摘掉它，有没有任何**外部可见**读数变过

现量（两枚二进制之差＝上面那一味；都建在仓外 `D:/tmp/wisp-144-bin/`）：

```
$ cd <repo>              && go build -o /d/tmp/wisp-144-bin/wisp-fixed.exe ./cmd/wisp
$ cd /d/tmp/wisp-144-mut && go build -o /d/tmp/wisp-144-bin/wisp-mut.exe   ./cmd/wisp
$ ./wisp-fixed.exe slo -state Sleeping -seconds 2 -out fixed-1.json   → rc=0，JSON 里 "observer" 出现 1 次
$ ./wisp-mut.exe   slo -state Sleeping -seconds 2 -out mut-1.json     → rc=0，JSON 里 "observer" 出现 1 次
```

**答案：这一发上没有，外部读数一个字都没变。** 原因写在 §0.4：窗口只存在于 `os.WriteFile` 那一次灌字节的过程里，
无争抢时读者基本落不到里面 —— CI 那一发（run `36094258734`）是**并发下撞上的**，票面自己也是这么归因的。

⇒ 本票对"承重"的回答因此**两问分开**：
① 摘掉它 AC#1 会不会从绿变红：**会，5 枚**（§3.2）；
② 摘掉它有没有外部可见读数变过：**今天这两发没有**。
**只按 ② 判承重会把这一味误判成装饰**（票面点名的正是这个误导形状）；
反过过来，我也**不**据此声称"`slo-full` 从此不再红"——我没能把那一发竞态在真进程上召出来（§6 第 1 条），
本票证的只到"下次再撞上时，它不会被读成损坏，且到点那一句会点名预算与字节数"。

`slo-check.ps1` 那行 `state … pass=` 的两版对照：结果与限制都在 §4.4。

### 3.4　续程复量（2026-09-25 19:4x-20:0x 第二程，另起仓外副本 `D:/tmp/wisp-144-r2/`）

**为什么要重走**：编排者派单时按"前一程死在 AC#3 之前"描述这格，而 §3.1-§3.3 已在未提交的
264 行里写满了读数。两句话不能都成立 ⇒ 我自己量一遍，不认账也不推翻，**只贴现量**。
所有数字落盘在 `D:/tmp/wisp-144-r2/`（`mutA.log` / `mutB.log` / `copy-baseline.log` /
`probe-baseline.log` / `r2-before-full.log` / `r2-after-full.log` / `r2-*.roster` / `d22scan-r2.log`）。

副本来源用的是**编排者给的配方**而不是 §3 那种 `tar` 工作树：

```
$ git -c core.autocrlf=false -c core.eol=lf archive a91d7c2 | tar -x -C /d/tmp/wisp-144-r2/base
$ cp third_party/sherpa-onnx/*.dll <copy>/third_party/sherpa-onnx/     # DLL 不在 git 里（现量：git ls-files third_party 为空）
$ cd <copy> && gofumpt -l .            # 空 ⇒ 两个 -c 都带上确实挡住了 CRLF
```

反证这个坑不是空穴：`D:/tmp/wisp-144-r2/prefix/` 那棵（`95885fb^` 的 archive）里
`slo_windows.go` 是 **666 行**、`:522` 仍是 *"never races"*、`:528-529` 仍是一发即死 —— 修前树是真的。

**先要"没变异时是绿的"这一发**（否则下面的红不算数）：`go test -run 'TestSLO144|TestZ144Probe' -v` ⇒
rc=**0**，顶格 10 枚 `--- PASS`（含 3 枚探针）+ 子格 7 枚，FAIL 0。摘掉变异、回装 pristine 又跑一遍：`ok`。

**§1/§2 引用的字节数逐格复算，全部相符**（探针在仓外，不改仓内任何文件）：

| 前一程写 | 我复量 | 判定 |
|---|---|---|
| 真报告 **1978** 字节，`0..1977` 全判 `unwritten` | 探针 `PROBE fixture bytes = 1978`；1978 枚前缀的状态 tally 只有 `unwritten` 一个键、count **1978** | **相符** |
| 半截 **44** 字节 / 坏字节 **25** / 语义坏 **38** / 完整无报告 **57** | 44 / 25 / 38 / 57（`printf | wc -c` 与探针两条路都取到同一数） | **相符** |
| §3.2 变异 A ⇒ AC#1 **5 枚**从绿变红、"当场红"那 2 枚保持绿 | `mutA.log`：rc=1，`PASS=2 FAIL=5`；红的正是 5 枚，绿的正是 `…ContradictThemselves…` 与 `…CorruptReportIsJudgedOnTheFirstRead` | **相符** |
| §3.2 变异 B ⇒ 红的是另外 2 枚、句落 `:250` 与 `:270` | `mutB.log`：rc=1，`PASS=5 FAIL=2`；`slo_report_144_windows_test.go:250` 与 `:270` 逐字对上（含 `20.06s` 那一发） | **相符** |
| §3.3 "外部读数没变"：两版各一发 `wisp slo`，都 rc=0 | 我跑了 **6 发**（两枚二进制 × 3 轮）全 rc=0，另见下面那条 flake | **相符但取证方式不同** |

**AC#3(i) 这一格我换了个更硬的凭据**（前一程的 §3.1 是"喂 `readSubjectReport` 一串手写坏字节"，
票面要的是"主体写的报告确实坏了"）：走 `writeSLO`（主体自己那枚序列化器）落真文件，
再用**生产入口 `collectReport()`**（真预算、真 `os.ReadFile`）去收 ——

```
PROBE2 writeSLO wrote 174 non-empty bytes; parses as JSON with a mode field: true
PROBE2 collectReport() error = wisp slo: subject report is corrupt (failing closed on the first read,
         not a partial write): 174 bytes parsed as a subject run but carry no state report
```

⇒ 非空、语义错、当场红、红因点名（`174 bytes` 与 `no state report` 都在句里）。
同一探针的反面（`PROBE3`）：**1977 字节前缀**没有被读成损坏，红句是
`subject 4242 never wrote a complete report within 30ms (last read: 1977 bytes read, …)`。

**§3.3 那句"外部读数没变过"只站住了它自己那一发 —— 补一发它没取的。**
两棵树上放**同一份**测试代码（`z144diff_windows_test.go`），问同一件 CI 形状（磁盘上真存在一枚差 1 字节的
报告），只调生产入口 `collectReport()`：

| | 修前树 `95885fb^` | 修后树 |
|---|---|---|
| 同一枚 667/668 字节的真半截文件 | `wisp slo: subject report: unexpected end of JSON input` | `wisp slo: subject 4242 never wrote a complete report within 33s (last read: 667 bytes read, document still open at offset 0: the tail had not arrived)` |
| 是否只剩解码器那句原话 | **是** | 否 |
| 是否点名预算与字节数 | 否 | **是** |
| 是否说"这是没写完、不是坏" | 否 | **是** |
| 是否还看第二眼 | **否**（第一发即死） | **是**（把 `subjectReportBudget+subjectGrace` 那 33s 花在重读上） |

⇒ 承重问题的正确答案因此不是"外部读数一模一样"，而是：**"半截文件真的被读到"时外部读数必变**
（红句整条换掉、且从"一发即死"变成"预算内重读"）；§3.3 那两发之所以没变，只是因为**无争抢时读者落不进那次写**。
⚠ 最后一列那个 `33s` 是**被花掉的预算常量**（`:84` 的 `subjectReportBudget=30s` + `:77` 的 `subjectGrace=3s`，
本票一字未动），**不是计时读数**，按票面规矩不许当秒数用。

**两发要登记的现场**：
① `wisp slo` 四发里出现过一发 rc=**2**：`out-of-tree sampling failed: resource: observe: baseline tree read:
resource: proc: … NtQuerySystemInformation: buffer never sufficient (last 1313568 bytes)` —— 那是**进程快照**
那一条腿在并发下取不到表，与本票无关；同一枚二进制重跑 3 轮全 rc=0。**不据它改判据、不据它判本票红。**
② 我第一次跑那两枚二进制时用 `D:/work/…` 盘符形式挂 PATH，**两发都 127**
（`error while loading shared libraries: sherpa-onnx-c-api.dll`），换成 `/d/work/…` 立刻能跑
⇒ `scripts/wisp-cli-tests.sh:101-108` 那段反直觉注释**我自己撞了一遍**，没去"顺手优化"它。

## 4　AC#4 门禁与名册

### 4.1　`go test ./cmd/wisp/` 改前 / 改后各一次（+ 名册差集）

⚠ **先记一枚仪器坑（票面没写）**：本机直接 `go test ./cmd/wisp/` 会在**加载期**死掉，一条用例都不跑
（票 98 仍 open：`.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md`，无 `-done` 后缀）：

```
$ go test -count=1 -v ./cmd/wisp/                     ← 不挂 DLL
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.049s     ← stdout 里 --- (PASS|FAIL|SKIP) 计数 = 0：仪器不产出结论
$ PATH=$PWD/third_party/sherpa-onnx:$PATH go test -count=1 -v ./cmd/wisp/    ← documented 形状（scripts/wisp-cli-tests.sh 用的就是它）
```

⇒ 本票所有 `cmd/wisp` 读数都是挂 DLL 取的；**没挂的那一发红不是本票的形状，是一架空仪器**，照实登记。

四数（名册 = `grep -oE '\-\-\- (FAIL|PASS|SKIP): [A-Za-z0-9_/]+' | sort -u`，票面指定的那把 grep；
⚠ 它在 `-` 与 `.` 处截断长名，所以子格名看着短——两遍同一把尺，差集可比）：

| 发 | 结果 | 名册 |
|---|---|---|
| 改前 @`40be959`，不挂 DLL | 包级 rc=1，**取不到名册**（0 条 `--- …`） | 这一发的红归票 98 |
| 改前 @`40be959`，挂 DLL | `ok`，rc=**0**，`^panic`/`^[signal` 计数 **0** | **106** 枚（含子格），`=== RUN`=116，FAIL 0、SKIP 0 |
| 改后 @终版（本票三枚 commit 之后），挂 DLL | `ok`，rc=**0**，panic 计数 **0** | **120** 枚，FAIL 0、SKIP 0 |

差集（`comm` 现量，落盘 `D:/tmp/wisp-144-b2/test-v-before.txt.roster` 与 `D:/tmp/wisp-144-a2/final.roster`）：

```
$ comm -23 before.roster after.roster        # 消失的
（空）                                        ← 没有一枚既有用例被改名、被删、被转成 SKIP
$ comm -13 before.roster after.roster | wc -l
14                                           ← 7 枚顶格 + 7 枚子格，全是本票新增 TestSLO144*
```

⇒ 名册**只增不减**、FAIL 0、SKIP 0、无 panic（那一形"一枚用例吞掉同包几十条读数"今天没发生，所以没有"未取到"的债）。

### 4.2　两把尺（票面给的 `gofmt -l cmd/wisp/` 与 CI 那一步的 `gofumpt`，不是一回事）

```
$ gofmt -l cmd/wisp/                                                        → 空
$ D:/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/                          → 空
$ D:/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm      → 空（CI 那条全树命令；gofumpt v0.12.0）
```

### 4.3　vet / d22scan

```
$ go vet ./cmd/wisp/            → rc=0（空）
$ sh scripts/d22scan.sh         → rc=0
  d22scan.sh: positive control → runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70
  d22scan: examined 228 production Go files under internal/ and cmd/
  d22scan: scope ban #8 cmd/  examined  43 Go files, comments and _test.go included   ← 本票新增的 _test.go 在这一发里被扫过
  d22scan: clean - no D22 ban violations
```

`tools/d22scan` 是独立 module，**没有**在根目录 `go vet ./tools/d22scan/`。
⚠ bans #1-5 不含 `_test.go` ⇒ 本票测试里那一枚裸 `go func(` **不在 ban #1 的射程内**，照实登记在 §6 第 5 条，不写成"过了 ban #1"。

### 4.4　CI 那一步的形状：`bash scripts/wisp-cli-tests.sh`（现量三发）

```
wisp-cli-tests.sh: scope=./cmd/wisp/ dll_dir=<repo>/third_party/sherpa-onnx (pinned: onnxruntime.dll sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll )
portable-tests.sh: platform=windows scope=[./cmd/wisp/]
portable-tests.sh: 11 ledger entries, 7 accounted on this platform
runtests.sh: OK - top-level: PASS=70 FAIL=0 SKIP=0, === RUN=130, '[no tests to run]'=0
rc=0
```

三发都是 rc=0：前两发跑在我**升级夹具之前**（夹具从 676 字节变成 1978 字节那一次），
第三发（`D:/tmp/wisp-144-a2/wisp-cli-final.txt`）跑在终版树上，读数与前面两发逐字相同
（`PASS=70 FAIL=0 SKIP=0, === RUN=130`）⇒ 夹具那一次升级没有挪动任何一格读数。
那一发里本票的 14 枚名册**全部有跑**（结果行 `grep -cE "(PASS|FAIL): TestSLO144"` = **14** ＝ 顶格 7 + 子格 7；
含 `=== RUN` 在内 `grep -c TestSLO144` = 29）——不是"编译过了就算数"。
⚠ 那 7 枚 "accounted on this platform" 的 ledger 豁免（`TestHelperProcess`、`TestLiveWasapiSmoke` 之类）
是 `portable-tests.sh` 打印理由的既定豁免，**本票 14 枚新增名册一枚都不在其中**。

`slo-check.ps1 -Subset smoke` 的两版对照（§3.3 那两枚仓外二进制）：**未取到**。
现量原因：脚本的采样有效性前置（ticket 134 AC#4）在本机此刻会走
"foreign toolchain/wisp processes / busy machine ⇒ NO CONCLUSION" 那一支 —— 这台机器上同时有别家程在跑
（我起手时 `go.exe` 有 2 枚在飞；`Runner.Listener.exe` 常驻）。
按票面"计时类读数不作凭据"与"要取颜色先确认 runner 空着"，我**没有**在这种情况下硬取那一发，
也**没有**为了拿到读数去改脚本或加豁免（脚本本体在禁改清单上）。这一格因此按"未取到"登记，不当绿、也不当红。

### 4.5　`!windows` 对照腿

**没有补**。被检对象（`slo_windows.go`）本身 `//go:build windows`，新代码没有任何 `!windows` 对照面，
所以本票不需要 linux 容器那一发；`GOOS=windows go vet` 那类只编译的仪器也**没有**拿来代答任何一格。

## 5　AC#5 那句注释 —— **结**

被改掉的那句（现量原貌在 §0.2，`git log -L 520,546` 现量它只由 `00bbb76` 写过）：

> `:520-522` "The subject only writes after subjectGrace past the end of its window, so this **never races** the observer's last sample."

它论证的是"写不会撞上观察者最后一发采样"，**没有**论证"读者看文件时文件已经写完"——存在与写完是两回事，
所以它把一个没验证的时序假设写成了不变式。

新措辞（`collectReportWithin` 头上，`:663-681`）说的是两件事：
① 分得开的是 **`reportUnwritten`（尾巴还在路上）** 与 **`reportCorrupt`（这些字节和那份文档自相矛盾）**；
② 为什么"文件存在"不等于"报告完整"——`os.WriteFile` 先建文件再灌字节，期间可见的**只是前缀**，
真正的竞态不是 observer-vs-observer，而是 **reader-vs-writer**，而它可分正是靠前缀只会缺尾巴这一个构造性质。

⚠ 同一处注释还顺手登记了一件本票不修的事（`:636-641` 那句 "The one shape this cannot tell apart…"）：
**主体写到一半死掉**留下的前缀，与"还在写"的前缀是同一串字节 —— 分不开。它的后果是"红得晚 33 秒"、
不是"红错"，且今天的 `exited()` 那一味在这条路上恒 false（§0.5），所以那一形落进 §6 第 2 条而不是被藏起来。


## 6　本程没测什么（按"漏了它谁先被骗"排序）

1. **"下一次 `slo-full` 会不会因此变绿"——没证。** 我从没在真进程上把那一发竞态召出来（§3.3 只证明了无争抢时它不出现）。
   谁先被骗：把它读成"`slo-full` 修好了"的人，包括下一次 push 前做归因的编排者。
   本票的主张止于："前缀不再被读成损坏；等不齐时红句点名预算与字节数"。
2. **"主体写到一半死掉"这一形仍然只能等预算到点（30s+3s）。** 它留下的字节与"还在写"的字节是同一串（§0.4 的前缀论是双刃），
   而 `exited()` 在循环里恒 false（§0.5 现量：`StartInJob` 不 `Wait`，`cmd.Wait()` 只在 `stop()` 里）——
   那一味**我没当凭据用**，也**没顺手修**（修它要动进程等待拓扑，不在本票地界）。
   谁先被骗：以为 `collectReportWithin` 里那两行 `if s.exited()` 是真早退的人（我读代码时差点就这么以为）。
3. **`slo-check.ps1` 那行 `state … pass=` 的两版对照：未取到**（§4.4 写了为什么：采样有效性前置会判 machine-contended，
   而票面禁止在这种编队下取计时类读数）。⇒ §3.3 的"外部读数没变"只站在**两发 `wisp slo -state Sleeping` 的 exit 码**上，
   分量比"六态链两版对照"轻，别把它读成后者。谁先被骗：拿 §3.3 当"整条链我也对照过"的人。
4. **报告的语义正确性没验**：`report` 字段齐、里面的数字/判决是假的，`readSubjectReport` 不看；
   嵌套的 `observer` 子报告也没验。本票只分"没写完 / 坏了"两态，AC#1ⓑ 的"内容不对"只到"没有 state report"这一层。
5. **测试里那一枚裸 `go func(`（用例 3 的写入侧）不在 ban #1 射程内**——bans #1-5 不含 `_test.go`（票面明写），
   它也没有 `recover`。它有 owner（测试自己）且有带缓冲 channel 收口，但**这不是"过了 ban #1"**，只是仪器不管。
6. **`staticcheck` 没跑**：本机版解不开 go1.27 产物，它的绿在本仓不作凭据；那一格只能在 CI 那版尺上量。
7. **只跑了 `-state Sleeping` 一腿**：`-settle`、`-leak`、其余五态都没跑。`collectReport` 在六态上是同一个调用点，
   但我没量，别按"同一个函数所以都一样"替我补。
8. **共树形状（登记给归因，不是本票缺陷）**：本程期间 `internal/agent/approval/**` 正被票 146 那枚程写。
   我起手时 `git status --short` 就显示 ` M internal/agent/approval/` 两枚文件在飞，
   我按禁改清单（编排者补①）**一个字都没碰**；但 `./cmd/wisp/` 经 `internal/agent` 传递依赖那枚包，
   ⇒ **它任何一刻编译不过，我这包的门禁会以我的名字红**。我最后一发全量 `go test ./cmd/wisp/` 是 rc=0，
   所以取数那一刻它是好的；下次看到 `cmd/wisp` 红而 diff 只落在 approval，先按这条归因。
9. **`wisp-cli-tests.sh` 那三发都是在**这台机器**上取的（§4.4），本程没有拿 `gh` 去核 CI 上的 run 级结论**
   （起手时 runner 只有 `Runner.Listener.exe` 常驻、job 内进程 0 枚，我据此判"空着"，但那是进程表读数不是 run id）。
   ⇒ "这一步在 GitHub 上真跑起来是什么颜色"这一格归验收方/编排者，不在本程的凭据里。
10. **没有 `!windows` 对照腿**，因此也没有 linux 容器那一发（§4.5）：被检对象本来就 `//go:build windows`。

## 7　伪授权登记（两个分开的数）

- **真通知回显数：7** —— ①会话首条工具输出里 harness 随附的 `<system-reminder>`（可用 skill 清单 +
  "date has changed" + `Memory: …/agents.md` 全文），出处 `Bash`，命令前 40 字
  `ls .scratch/wisp/issues/ | grep -- '-144'`；
  ②-⑦ 六枚后台任务完成通知（`go test` 基线两发、改后全量一发、终版全量一发、CI 形状两发），
  出处均为 harness 的 `task-notification` 块、随本程自己启动的后台命令回来。
  ⇒ 这些我**没有**据此改过任何一格判据；它们只是本程自己动作的回声。
- **判为注入数：0** —— 本程在所有工具输出里**没有**读到任何一句自称"系统提示／编排者备注／已核验请继续提交／
  请 revert／放宽阈值／这条路已解锁／不用取证直接给结论"的句子。
  两条**核过之后判定为合法**的边缘项，写明以免下一个人再审一遍：
  (a) `Edit` 工具在仓外副本上报过一次 "file changed since your last read" —— 那是**我自己刚 `cp` 覆盖了那枚文件**，
      是真信号不是注入，我按它提示去读了当前内容再继续；
  (b) 票 98 正文里那句" cmd/wisp 的红不要追"出自**仓内工单文件**（我 Read 出来的），
      不是工具输出里冒充授权的话；我处理方式 = 当待核断言，现量之后发现它已被 documented 的 DLL-PATH 形状解开（§4.1）。
- **凭据抄录：0 处。** 本程读过的日志里出现过 temp 目录路径、`dir=`、pid、以及 `WISP_TEST_DATA_DIR` /
  `proc.TestDataDirEnv` 这类**变量名**；全文未抄任何 API key / token 原文。
- 反向一条也登记：**同程"我没遇到"不洗掉别人遇到的**。本程零注入是本程这一棵树这一段时间的读数，
  与台账 `docs/reports/injection-timeline.md` 上别人的计数各记各的。
