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

## 1　AC#1 钉成能响的用例（结／不结 + 现量读数）

## 2　AC#2 修法：只分"没写完"与"坏了"，不放宽 fail-closed

## 3　AC#3 变异自证两向（含"摘掉它有没有外部可见读数变过"）

## 4　AC#4 门禁与名册

## 5　AC#5 那句注释改了什么

## 6　本程没测什么（按"漏了它谁先被骗"排序）

## 7　伪授权登记（两个分开的数）
