# 255-r7 交付件 — AC#6 名册行号引用降级为内容锚 ＋ A878 两处重裁（写面＝cmd/wisp 两枚）

本腿派单原名 255-r2 撞 10-03 旧腿；编排者 10-14 08:5x 更正为 255-r7（存档 `.scratch/wisp/probes/orch/20261013/txt/255r7-correction.txt`）。
本腿 23:58 前后在旧 `r2/` 目录建的两枚件（`inventory-cites.py`＋`logs/{inventory.txt,pre-run.txt}`）已**复制**到 `r7/` 续用，**原件未删**（微残留；`../r2/` 现 20 枚条目）。本机 2026-10-13 23:58 → 10-14 08:53 挂起约 9 小时，本腿此间零落盘。

- 起手 HEAD（派单写）`ad2dadd3`；交付 HEAD `dcaff922`（＝本腿产码笔，父 `7db2fea0`；期间编排者落 58366328..7db2fea0 五笔，最近一笔碰本腿两枚文件的 commit 仍是 `9dade451`＝293-r1 那笔）。
- 写面＝`cmd/wisp/config_readers_255.go`＋`cmd/wisp/config_receipt_255_test.go`（＋本证据目录）。⛔ 零 push。
- 交付态 sha256：roster `6b840071c183b3ed6470935508247037c416ed3787b22c1b50f839a56f2b55f2`／test `2195e669aab19be80643d0040be67c3fdbd716d082b4e04de25d60caf68dfa4e`（gofmt/gofumpt 皆空；gofmt 洗前 roster＝ed6d0143…，备份在仓外）。
- 两笔 commit：产码 `dcaff922`（两枚 Go 文件）＋交付件 `<B>`（本件＋logs 子集＋票尾 bullet）。

## 1. AC#6 落地（判据逐字：token 唯一命中当判据，行号只当读数）

**尺的逻辑改写**（`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`）：
- 旧尺把 cite 读成 `path:LINE [token]` 并要求**那一行**仍含 token ⇒ 插入一行即红。缺陷面现成两枚：313-r2 落地插两行把 llm 行 `:96`→`:98` 打红（`logs/pre-run.txt` 红句逐字）；受控重演＝在 run.go 被引行上方插一行，旧尺 08:59:34→08:59:55 rc=1（`logs/probe-a2-old-ruler-insert-red.txt`）。
- 新尺把 token 当锚：在**被引文件整文件**里找，**恰好命中 1 行**＝绿；**0 命中**级红（token 被改），且逐名点名 row／cite／**该行现读**；**>1 命中**级红（锚不唯一，须削尖）。`:LINE` 降级为**读数**：漂了只 `t.Logf`（不判红）。
- 行号读数**一枚未删**（边界：摘尺不许当修尺）：尺① `grep -o '[A-Za-z0-9_./]*\.go:[0-9]\+' | wc -l` 改前 27／改后 27（去重 20＝20）；尺② `grep -o ':[0-9]\+ \['| wc -l` 23 枚、`grep -c` 22 行，改后同数（`logs/final-measurements.txt`）。

**六枚不唯一 token 削尖**（判据要求「token 唯一命中」，起手实测六枚不唯一⇒不削尖则 ⓑ 在这些 cite 上不会红；起手读数＝`logs/inventory.txt`，终态＝`logs/inventory-after-all-roster.txt`）：

| row | 起手 token（文件内命中数） | 入册 token（现读，命中 1） |
|---|---|---|
| llm | `cfg.LLM.Roles.Chat.Model`（2：98/99） | `ChatModel:      cfg.LLM.Roles.Chat.Model` |
| ball | `SizePx`（13） | `SizePx   int // configured orb size 44..72` |
| audio | `c.Audio.MicMutedDefault`（5） | `ra.mutedAtBoot = c.Audio.MicMutedDefault` |
| privacy | `RedactPaths`（3：50/51/88） | `RedactPaths bool` |
| voice×4 | `c.Voice.Enabled`（3：253/257/262） | `ra.voiceEnabled = c.Voice.Enabled` |

（另：agent 行**注释**里 `cfg.Agent.PerToolTimeoutMS`（2：763/1014）顺手削尖为 `PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS)`——注释不在尺内，为"全名册 token 唯一"齐。）

**四处漂读更新**（读数字段更新为现读数，§A878-2 在内）：
- llm：`:96`→`:98`（313-r2 两行插入所致）；
- privacy：`:50`→`:51`（削尖后 token 真身在字段行 51，50 是注释行）；
- panel 行注释：`panel_resident_windows.go:253`→`:273`（本腿现量抓出的既有漂读：`NewPanelManager(` 唯一非测试调用点在 `:273`，`:253` 是 `newResidentComposerDispatch`）；
- audio 行注释：`internal/audio/gate.go:57`→`:55-56`（那句 doc 注释真身在 55-56，57 是 func 行）。

**测试侧 4 行 `.go:数字`→2 行**：L440 历史指针（`:304`，AC#4 后该处已非几何字面量）转内容锚＋`HEAD 67ab595d` 限定；L445 断言 `contains "panel_resident_windows.go:207"` 改为 file＋`[cfg.Panel.Width]`（「断言不许再依赖行号」）。余下 2 枚（L81/L529 的 `panel_config_store.go:92`）是**活的、准确**读数，照边界保留为读数。

## 2. A878 两处重裁

1. **voice 读者集**：`sectionReadSites["Voice"]` 补 `cmd/wisp/panel_config_store.go`＋`internal/panel/config_handlers.go`（313-r1/r2 落点所致；红句逐字在 `logs/pre-run.txt`）；四枚 voice 键行去 "the section's only production reader since ticket 247" 表述，改为"其余读者答别的键"（摄取腿 boot 读 voice.enabled＋313 设置腿/面板视图读云链、realtime 键；三枚文件由 sectionReadSites 钉）。
2. **llm 证据行**：`:96 [cfg.LLM.Roles.Chat.Model]`→`:98 [ChatModel:      cfg.LLM.Roles.Chat.Model]`（读数字段更新＋token 削尖，见 §1）。

## 3. 两形读数（能力形；全部在**交付态 sha** 上复跑；仓外备份＋sha256 逐字节证）

ⓐ **在同一枚被引文件里插一行、token 原样保住 ⇒ 不许变红**：

| # | 被引文件（插在 cite 上方） | 窗口（+08） | rc | 读数 |
|---|---|---|---|---|
| a1 | `cmd/wisp/run.go`（agent 行） | 09:13:56 → 09:14:09 | 0 | 绿＋`t.Logf`「recorded reading cmd/wisp/run.go:991 drifted, the token "cfg := rt.cfg" now uniquely sits at line 992」 |
| a2 | `cmd/wisp/resident_audio_windows.go`（voice 四行） | 09:14:32 → 09:14:44 | 0 | 绿＋4 行 `t.Logf` 漂读（`:257`→`:258`） |

ⓑ **把被引那一行的 token 改掉 ⇒ 必须红且逐名点到那一行**：

| # | 被引行（改法可编译） | 窗口（+08） | rc | 红句要点 |
|---|---|---|---|---|
| b1 | `run.go:991` `cfg := rt.cfg`→双空格 | 09:14:19 → 09:14:32 | 1 | `row "agent" cites cmd/wisp/run.go:991 for "cfg := rt.cfg" and the token is now nowhere in that file (the recorded line reads "cfg :=  rt.cfg")` |
| b2 | `resident_audio_windows.go:257` 同法 | 09:14:54 → 09:15:05 | 1 | 红 4 条，逐名 `voice.tts.speed`／`voice.punctuation`／`voice.wake_word.thresholds`／`voice.wake_word.veto_words` |

ⓒ（支撑格，非判据要求）把 ball token 暂回 `SizePx`（13 命中）⇒ 红「not unique … hits lines [64 144 145 147 148 150 151 229 233 346 377 560 698]」：09:15:06 → 09:15:19 rc=1（`logs/probe-v2-c1-ambiguous-token-red.txt`）。

**还原证**（`sha256sum` 改前=改后，`git status` 空）：`run.go` `d60c50e3e6af98dc5c6f6842574c3cf47c0093956063b5bd8813dc171eb4c68d`；`resident_audio_windows.go` `a841f5220ea68f2158a8ddb34cd378a26b672daf1e20b0f62a9cfed8d2609a84`；roster ⓒ 后回装=交付态 `6b840071…`。备份＝`C:/tmp/wisp255r7backup/{run.go.orig, resident_audio_windows.go.orig, config_readers_255.go.deliverable-final, config_receipt_255_test.go.deliverable}`。

## 4. 门禁 rc 表（每件自落 rc，无管道吞退码）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | 0 | —（log 0B 在 `logs/gate-build.txt`） |
| 格式 | `gofumpt -l <两枚>` | 0 | 空；`gofmt -l` 同空 |
| 本票全名册 | `go test -count=1 -run 'TestTicket255' ./cmd/wisp/` | 0 | `ok … 25.371s`（20 枚顶层全 PASS；log 49B 在 `logs/gate-test255.txt`） |
| D22 | `sh scripts/d22scan.sh` | 0 | clean，no D22 ban violations（tail 在 `logs/gate-d22scan.txt`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN（longest=184、wall interval=0；`logs/gate-pathlength.txt`） |

## 5. 具名裁定（判语归非实现者；本腿自己裁的形，摆出来受审）

1. **六枚 token 削尖**＝「token 唯一命中当判据」的机械实现；起手六枚不唯一（§1 表）。不削尖则 ⓑ 在这些 cite 上不会红（同 token 另有命中）＝两形缺腿。每枚削尖前后命中数两件 inventory 可复算。
2. **四处漂读更新＋测试侧 L440/L445**：依据＝「行号只当读数」「断言不许再依赖行号」；⛔ 未删任何行号读数（尺①/②同数＝27/20、23/22）。
3. **「token 在文件里不唯一」本腿按判据处理为"削尖后入册"，未停手**：若判语认为该形应先停手上报，§1 表的起手读数与削尖候选都在件内；退回则按判语重裁（这是本件唯一需要非实现者重裁的形）。

## 6. 具名没做（⛔ 不许读成已做）

- ⛔ 零 push；⛔ 未动 313 面在飞/刚落件（`panel_config_store.go`／`config_handlers.go` 等只许读——本腿只读）；⛔ 未动三冻结件／`thresholds.go`／golden／allowlist／`docs/PLAN.md`／`docs/specs/**`；⛔ 未扩到别的按行号引产码的仪器；⛔ 未放宽任何断言；⛔ 未 `t.Skip`；⛔ 未跑整包 `go test ./cmd/wisp/`（314 死区，按令不拿整包当判据）；⛔ 未跑真窗（本格无此项，`-run 'TestTicket255'` 覆盖）；⛔ 未动共享树里别人的脏件（design/**、go.mod、.gitignore、`?? -` 等一概不碰）。
- 未跑的格：`go vet ./cmd/wisp/`（不属本腿派单门禁表，未跑）。
- 大原始日志留盘**不入库**（点名＋字节数见 §7）。

## 7. 件与残留（logs/ 两档）

**入库（本 commit）**：`00-anchor-and-gates.md`／`inventory-cites.py`（r7 副本，已清一行死赋值；原始件在 `../r2/`）／`logs/{inventory.txt(起手,4333B), inventory-after-all-roster.txt(终态,4095B), final-measurements.txt(343B), sha-before-probes.txt(382B), probe-a1-old-ruler-green.txt(48B), probe-a2-old-ruler-insert-red.txt(572B), probe-v2-a1-insert-run-green.txt(474B), probe-v2-b1-token-change-run-red.txt(493B), probe-v2-a2-insert-voice-green.txt(1147B), probe-v2-b2-token-change-voice-red.txt(1480B), probe-v2-c1-ambiguous-token-red.txt(567B)}`。

**留盘不入库**：`logs/pre-run.txt`(14673B，改前两枚红原始日志)／`gate-d22scan.txt`(23734B)／`gate-test255.txt`(49B)／`gate-build.txt`(0B)／`gate-pathlength.txt`(1457B)／`inventory-after-roster.txt`(4319B)／`roster-tests-after.txt`(565B)／v1 四枚 probe 日志 `probe-b1-…`(474B)／`probe-b2-…`(1147B)／`probe-c1-…`(493B)／`probe-c2-…`(1480B)／`probe-d1-…`(567B)（与 v2 同形，roster gofmt 洗格式前取；v2 为交付态正本）。

## 8. 下一位要点

- 两形可原样复跑（§3）；被引文件挑非 313 稳定件（run.go／resident_audio_windows.go 本腿用过）；突变窗口前后 `sha256sum` 对拉。
- AC#6 框未勾、判语归非实现者；票尾只追加 1 条 bullet（⛔ 不翻 `- [ ]`、⛔ 不动 `Status:`）。
- §5.3 是唯一"自行裁了形"的点（token 不唯一→削尖入册），若要退回首选手续在 `logs/inventory.txt`（起手六枚读数）＋§1 表（候选）。
