# 197-a1 尺与 rc 逐条（现跑；每条尺跑完**下一句立刻**取 rc，⛔ 中间不插 echo/管道）

---- R1 mergeOverflow 全仓 .go
$ grep -rn mergeOverflow --include=*.go internal cmd tools
rc=1
lines=0


---- R2 Truncated( 调用形状
$ grep -rn --include=*.go -F Truncated( internal cmd
rc=0
lines=12
internal/observe/logging_test.go:92:func TestRedactLongArgTruncated(t *testing.T) {
internal/panel/pump.go:632:func (s *StreamLog) Truncated() bool {
internal/panel/pump_test.go:196:	if !sl.Truncated() {
internal/panel/subagent_roster_197_test.go:340:	if !log.Truncated() {
internal/panel/subagent_roster_197_test.go:383:				StreamTruncated:   log.Truncated(),
internal/panel/subagent_roster_197_test.go:434:			StreamTruncated:   log.Truncated(),
internal/panel/subagent_stream_197_test.go:106:	if !sl.Truncated() {
internal/panel/subagent_stream_197_test.go:204:	if sl.Truncated() {
internal/panel/subagent_stream_197_test.go:225:	if !sl.Truncated() {
cmd/wisp/panel_pump.go:203:	state.StreamTruncated = rt.stream.Truncated()
cmd/wisp/subagent_carrier_197_test.go:553:	if !sect.StreamTruncated || !driven.stream.Truncated() {
cmd/wisp/subagent_carrier_197_test.go:555:			sect.StreamTruncated, driven.stream.Truncated(), panel.DefaultStreamKeys+1)

---- R3 ElidedRunes( 调用形状
$ grep -rn --include=*.go -F ElidedRunes( internal cmd
rc=0
lines=13
internal/panel/pump.go:644:func (s *StreamLog) ElidedRunes() int {
internal/panel/pump_test.go:199:	if sl.ElidedRunes() != 0 {
internal/panel/pump_test.go:200:		t.Errorf("elided = %d, want 0: short streams keep full fidelity past the bound", sl.ElidedRunes())
internal/panel/subagent_roster_197_test.go:384:				StreamElidedRunes: log.ElidedRunes(),
internal/panel/subagent_roster_197_test.go:396:	if sect.StreamElidedRunes != log.ElidedRunes() {
internal/panel/subagent_roster_197_test.go:397:		t.Errorf("section total = %d, want the log's own %d", sect.StreamElidedRunes, log.ElidedRunes())
internal/panel/subagent_roster_197_test.go:435:			StreamElidedRunes: log.ElidedRunes(),
internal/panel/subagent_stream_197_test.go:109:	if got := sl.ElidedRunes(); got != 0 {
internal/panel/subagent_stream_197_test.go:266:	if sl.ElidedRunes() != totalElided {
internal/panel/subagent_stream_197_test.go:268:			sl.ElidedRunes(), totalElided)
cmd/wisp/panel_pump.go:204:	state.StreamElidedRunes = rt.stream.ElidedRunes()
cmd/wisp/subagent_carrier_197_test.go:557:	if sect.StreamElidedRunes != driven.stream.ElidedRunes() {
cmd/wisp/subagent_carrier_197_test.go:559:			sect.StreamElidedRunes, driven.stream.ElidedRunes())

---- R4 DroppedKeys( 调用形状
$ grep -rn --include=*.go -F DroppedKeys( internal cmd
rc=0
lines=12
internal/panel/pump.go:408:// NAMED in DroppedKeys() - "nobody is showing this one" is a different sentence from
internal/panel/pump.go:657:func (s *StreamLog) DroppedKeys() []string {
internal/panel/pump_test.go:219:	if len(big.DroppedKeys()) != 50-len(got) {
internal/panel/pump_test.go:220:		t.Errorf("dropped = %v, want the %d keys no row covers", big.DroppedKeys(), 50-len(got))
internal/panel/subagent_roster_197_test.go:421:	dropped := log.DroppedKeys()
internal/panel/subagent_roster_197_test.go:433:			DroppedStreamKeys: log.DroppedKeys(),
internal/panel/subagent_stream_197_test.go:112:	if got := sl.DroppedKeys(); len(got) != 0 {
internal/panel/subagent_stream_197_test.go:270:	if got := sl.DroppedKeys(); len(got) != 0 {
internal/panel/subagent_stream_197_test.go:292:	dropped := sl.DroppedKeys()
cmd/wisp/panel_pump.go:205:	state.DroppedStreamKeys = rt.stream.DroppedKeys()
cmd/wisp/subagent_carrier_197_test.go:561:	if strings.Join(sect.DroppedStreamKeys, ",") != strings.Join(driven.stream.DroppedKeys(), ",") {
cmd/wisp/subagent_carrier_197_test.go:563:			driven.stream.DroppedKeys())

---- R5 pump.go 三处锚
$ grep -n -e Overflow TRUNCATES and never merges -e no key is ever folded into another -e never merged into a row internal/panel/pump.go
rc=0
lines=3
397:// Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule, which
573:// no key is ever folded into another: past maxKeys the log declares itself
655:// and the host can say WHICH ones and say they were never merged into a row that is

---- R6 panel_pump.go 锚
$ grep -n merge/backpressure rule cmd/wisp/panel_pump.go
rc=0
lines=1
397:// for one more word of a reply, and the merge/backpressure rule that decides how

---- R7 票35 里 merge 的每一枚命中行
$ grep -n -i merg .scratch/wisp/issues/35-panel-bridge-c17.md
rc=0
lines=7
14:approval requests, ball state, cost ticks) with bounded-queue merge, and `panel.resync` full-
27:- Event push: bounded queue per panel; overflow MERGES increments (never drops content);
49:- [ ] Backpressure: flood events under blocked consumer → merges, no unbounded memory (heap cap
51:  〔**10-08 11:3x 编排者就地裁冲突，来源＝只读普查腿 `35-a6`（件 `.scratch/wisp/probes/35/a6/verdict.md`，盘上逐字我未复跑，⛔ 本框保持未勾）**：这一格三句里 **"merges" 与 "content integrity kept" 两句已被票 197 leg B 的重裁取代**——产码 `internal/panel/pump.go:397` 逐字写着 `Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule`，且该规则有命名测试钉死（三枚既有尺反咬 merge）。⇒ 本框今天**只剩"no unbounded memory（heap cap asserted）"那半句算账**，而那一半按普查是"每面板推队列没写"＝一块没写。⛔ **任何后续程不许把 pump 改成"会合并"来迎合上面那句旧话**（那是把一枚已裁的契约规则倒回去），要改必须先落 `A##` 并经人工批准。〕
94:⇒ 票面 §"What to build" 里那句 "Go→frontend event push (task deltas, tool chips, approval requests, ball state, cost ticks) with bounded-queue merge" **今天兑现枚数＝0**，而它正是 `panel.resync` 那一格的前置。
392:| `:49` Backpressure | **部分在、且钉的方向与 AC 相反**：`internal/panel/pump.go:397` 逐字 "**Overflow TRUNCATES and never merges** (ticket 197 leg B re-cut this rule…)"，三枚计数出口 `:632/:644/:657`；面板线程队列有界 `cmd/wisp/panel_resident_windows.go:148 make(chan func(), 16)`＋背压计数 `:130/:369 failedPost`（**零测试引用**）。票面 `:27` 那枚"每面板 bounded queue"**一块没写**（出向零产码） | **先裁冲突才谈落点**；按"merge 只用于新事件推、不碰 StreamLog"走＝2–3 枚 250–450 行、**前置＝`:47`** | **不要**（同上例外：事件名进 `bridge.go` 才要） | 三枚既有尺 `pump_test.go:165 TestTheStreamLogTruncatesInsteadOfMerging`／`subagent_stream_197_test.go:209`/`:275` 会因改回 merge 而红＝**唯一有牙的尺，咬的方向与 AC 相反** |
406:**给编排者的排程后果（本腿不替他裁，只把可派性摊清）**：**今天真正可派的只有 `:42` 前半格与 `:51` 的残差那一格**（两枚都是纯测试、零产码、零新名、且有牙）。`:44`/`:47`/`:49` 三枚共用同一枚前置＝**出向那一跳**，而它被三枚外部条件挡着（票 `:97` 的 ⓐ/ⓑ/ⓒ 归属未定＋`frontend/**` 只读＋票 194 名册未补）；`:49` 额外要先裁"merge vs truncate"那处与票 197 的正面冲突。`:45`＝归口、零落点。**⛔ 本腿未翻任何框、未加 `-done`、未 push。**

---- R8 票36 里 merge 命中
$ grep -n -i merg .scratch/wisp/issues/36-result-history-panel.md
rc=0
lines=1
44:- [ ] Long-task stream: 200-chunk golden replay renders with merge (no jank), token counter

---- R9 PLAN.md 里合并增量/背压与合并
$ grep -n -e 合并增量 -e 背压与合并 docs/PLAN.md
rc=0
lines=2
2848:- LLM stream → UI：有界 channel，满则**合并增量**（不丢内容）。
2971:③ correlationId 路由 + 事件推送的背压与合并（D38(d)）

---- R10 SPEC-01 合并增量
$ grep -n 合并增量 docs/specs/SPEC-01-architecture.md
rc=0
lines=1
133:- LLM 流→UI：有界 channel，满则**合并增量**（不丢内容）。

---- R11 doc.go 背压/合并
$ grep -n backpressure/merge internal/panel/doc.go
rc=0
lines=1
10://     routing + push backpressure/merge (D38d)

---- R12 票面框数尺 票36 未勾
$ grep -c ^- \[ \] .scratch/wisp/issues/36-result-history-panel.md
rc=0
lines=1
6

---- R13 票面框数尺 票36 已勾
$ grep -c ^- \[x\] .scratch/wisp/issues/36-result-history-panel.md
rc=1
lines=1
0

