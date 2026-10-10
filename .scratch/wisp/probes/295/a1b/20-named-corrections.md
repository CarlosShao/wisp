# 票 295 · 补位腿 `295-a1b` · 第 3 笔：具名更正（死腿原话 ↔ HEAD 现量 ↔ 差在哪）

HEAD `29081a13…`；死腿件 `.scratch/wisp/probes/295/a1/10-ac0-field-roster.md`（提交 `3f9bedf7…`）。
⛔ 本件**不改死腿那份件一字**，按本仓规矩另写一节具名更正。凡枚数一律同句写清：尺逐字／射程目录／含不含 `_test.go` 与注释行。

> **前提读数**：`git diff --stat 448f5a57 HEAD -- cmd/wisp` 输出为空 ⇒ 死腿锚 `448f5a57` 与本腿锚 `29081a13` 之间
> `cmd/wisp/**` 一字未动。⇒ 下面 6 枚对不上里 **没有一枚是行号漂移**，全是**枚数／归属／尺的口径**；
> 落地腿照死腿的行号改，行号本身不会错，但**照它的枚数做缺口审计会漏**。

## D-1 · `.started` 命中枚数：死腿 12 / HEAD 11

- 死腿原话（名册 A #6 行末）：尺＝``git grep -n -o "\.started\b" HEAD -- cmd/wisp``＝**12 命中**（声明外：写 2 枚、读 10 枚，其中生产读 1 枚）。
- HEAD 现量：同一条尺（含 `_test.go`、含注释行，目录 `cmd/wisp`）＝**11 命中**——写 2（`resident_audio_windows.go:168`、`:323`）、读 9（生产 1＝`:379`，测试 8＝`resident_audio_247_windows_test.go:113`/`:155`/`:256`、`resident_audio_247_live_windows_test.go:148`、`resident_mute_290_windows_test.go:124`/`:141`/`:142`/`:174`）。
- 同尺在死腿自己的锚 `448f5a57` 上也是 **11** ⇒ 不是漂移，是它数错；它自己那一行逐枚列的点合计也只有 11（2+1+8），属**自相矛盾**。
- ⇒ 派单口径应写：**11 枚点／3 枚生产（2 写 + 1 读）／8 枚测试**。

## D-2 · 测试面读点枚数：死腿"7 处"与"＝8 处"自相矛盾

- 死腿原话（ⓑ 末节"给 `295-r1` 的三条具名前提"第 1 条）：测试面另有 **7 处读点**要跟着换形状（列了 `247_windows:113`/`:155`/`:256`、`247_live:148`、`mute_290:124`/`:141`/`:142`/`:174`＝8 处含 `:142` 的 `t.Fatalf` 参数位）。
- HEAD 现量：它自己列的清单就是 **8 枚**，逐枚原文本腿已验；"7 处"那个数没有任何一把尺能得出来（尺：同 D-1，去掉生产那一枚＝8）。
- ⇒ 派单只许写 **8 枚测试读点**，并保留它那枚最有价值的点名：`:142` 是 `t.Fatalf` **参数位**（见 `30-landing-and-anti-form.md` ⓓ）。

## D-3 · 类型对齐尺的文件归属：`resident_approval_windows.go:167` 不是那句

- 死腿原话（名册 A 前的"锚点陷阱"段）：对齐依据＝`resident_approval_windows.go` 里 `ra` 的声明是 `:167` 逐字 `	ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`。
- HEAD 现量：`git show HEAD:cmd/wisp/resident_approval_windows.go` 第 **167** 行逐字是注释 `// one, otherwise the live ball the UI is holding.`；
  那句原文真身在 **`cmd/wisp/resident_windows.go:167`**（行号恰好同为 167，**文件错**）。
  相关真锚补三枚：`resident_approval_windows.go:301` `	return newResidentApprovalWithConfig("")`、`:363` `func newResidentApprovalWithConfig(dataDir string) *residentApproval {`、`:365` `	ra := &residentApproval{root: root, cancel: cancel}`。
- 差在哪：**结论对、尺指错文件**。它据以判"grant_writer_265:585 的 `ra.cancel` 不是 `residentAudio` 的点"这件事**成立**（真凭据＝同件 `:584` `	ra := newResidentApproval()`），但落地腿⛔ 不许照它去 `resident_approval_windows.go:167` 找那句话。

## D-4 · `resident_windows.go:260` 的接收者不是本票那两枚 struct

- 死腿原话（线程词汇表 `boot` 行）：`:260` 逐字 `	defer ra.detachBall()`（作为 boot 线程锚点，ⓐ 与名册 B #1 又拿它算 LIFO defer 序的一环）。
- HEAD 现量：行原文 ✓ 对上；但 `detachBall` 定义在 `resident_approval_windows.go:781` `func (ra *residentApproval) detachBall() {` ⇒ 那一行的 `ra` 是 **`*residentApproval`**，⛔ 不是 `residentAudio`/`residentBall` 的字段点。这正是死腿自己"锚点陷阱"段警告过的 `ra` 复用。
- 差在哪：线程归属（`boot`）与 defer 序结论**不受影响**，但落地腿若把 `:260` 当作"本票两枚 struct 的读/写点"会多算一枚。测试面同名 defer 另有 3 枚（`resident_approval_live_246_windows_test.go:101`/`:210`/`:302`），同样属 `residentApproval`。

## D-5 · `panel_resident_windows.go` 的 atomic.Int 枚数：死腿"四枚 `:126-129`" / HEAD 五枚 `:126-130`

- 死腿原话（ⓑ 第 2 行）：`cmd/wisp/panel_resident_windows.go:126-129` 四枚 `atomic.Int*`。
- HEAD 现量（尺＝`git show HEAD:cmd/wisp/panel_resident_windows.go | cat -n | sed -n '124,138p'`，含注释行）：`:126` `	toggles    atomic.Int64`、`:127` `	shows      atomic.Int64`、`:128` `	hides      atomic.Int64`、`:129` `	disposals  atomic.Int64`、`:130` `	failedPost atomic.Int64` ＝ **5 枚**，区间该写 `:126-:130`。
- 差在哪：少算一枚、行区间少一行。⚠ 这条不影响"生产 `atomic.Bool`＝2 枚"（那一句复跑通过，见 `10-anchor-recheck.md` §3 第 3 行）。

## D-6 · `cancelHosted` 的测试面：把构造当读点、真读点漏列、枚数少一

- 死腿原话（名册 B #2 读点列）：测试（`resident_approval_246_windows_test.go:331`、`resident_cancel_key_*` 两枚字面量）。
- HEAD 现量（尺＝`git grep -n -E "cancelHosted" HEAD -- cmd/wisp`，含 `_test.go`、含注释行）：
  - 真**读**点＝`resident_approval_246_windows_test.go:315` `	if rb.cancelHosted {`（＋`:316` 的 `t.Fatal` 文本），死腿**未列**；
  - 它列的 `:331` 逐字 `	if bound := ra2.bindBallHost(&residentBall{cancelHosted: true}); bound {`＝**构造字面量（写形）**，⛔ 不是读点；
  - `resident_cancel_key_*` 里的字面量 HEAD 现量＝**3 枚**（`label_260r4_windows_test.go:89`、`wording_260r3_windows_test.go:99`、`wording_260r3_windows_test.go:141`），死腿写"两枚"。
- 差在哪：读/写类别搞混 ＋ 枚数少一枚。⇒ 对本票无直接后果（`cancelHosted` 不在 `AC#0` 射程），但它是"名册尺必写类别"这一课的第二次样本，报编排者记账。

## D-7 · `voiceEnabled`/`mutedAtBoot` 那把尺的射程口径（结论对、尺不严）

- 死腿原话（名册 A #4）：尺＝``git grep -n -E "voiceEnabled\|mutedAtBoot" HEAD`` ⇒ 源码命中只有 `:100`/`:101` 声明与 `:256`/`:257` 两处写，读＝0。
- HEAD 现量：限定 `-- cmd/wisp`（含 `_test.go`、含注释行）＝**4 命中**（上面那四行），零读者 ✓ 成立；
  但它写的那把尺**没带目录**，HEAD 全仓现量另有 **12 枚** `.scratch/**` 与 `docs/**` 里 `.md` 文件的命中（票 247/290/293/295 的探针件与台账、`HANDOVER.md`、`pending-and-issues.md`）⇒ 照原尺抄数会得 16，"源码/文档"混在一个数里。
- 差在哪：⛔ 不是结论错，是**尺没写射程**。本票第 4 次口径病。

## D-8 · 死腿件里 ⓓⓔ 两节的实际状态（编排者转述"从未产出"＝与原文一致，另有 ⓒ）

- 原文核：`.scratch/wisp/probes/295/a1/10-ac0-field-roster.md` 通篇只有 线程词汇表 ＋ ⓐ 名册 A ＋ ⓐ 名册 B ＋ ⓑ 形状表 ＋ "给 `295-r1` 的三条具名前提"。
  ** ⓒ／ⓓ／ⓔ 三节标题在原文里根本不存在**，`00-anchor-and-rulers.md` 之外该腿也没有第二枚交件（`ls .scratch/wisp/probes/295/a1/` 现跑＝**只有 1 个文件**）⇒ 编排者转述"ⓒⓓⓔ 从未产出"**属实**，本腿不具名更正它。
- ⚠ 但那条"给 `295-r1` 的三条具名前提"实际上**已经把 ⓓ 做了半节**（点到 `:142` 是格式串参数位）。本腿 `30-landing-and-anti-form.md` 的 ⓓ 是在它之上**逐枚补全并改正枚数**，⛔ 不是重复劳动，也⛔ 不改它原句。
