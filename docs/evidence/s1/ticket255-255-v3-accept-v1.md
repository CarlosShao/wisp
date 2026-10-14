# 255-v3 非实现者判语 — 票 255 AC#6（名册行号引用降级为内容锚）＋ A878 两处重裁

**验收腿**：`255-v3`（非实现者；实现者＝落地腿 `255-r7`）。本件只出判语：⛔ 不翻票框、⛔ 不动 `Status:`、⛔ 不 push。
**被验收对象**：产码 `dcaff922`（`cmd/wisp/config_readers_255.go`＋`cmd/wisp/config_receipt_255_test.go`）／交付 `9bd5ff21`（锚＋logs 子集＋票尾 1 bullet）／更正 `94398165`（追加不 amend，锚 §7 计数字）。
**环境**：win32 + Git Bash，go1.27.1 windows/amd64；`PATH` 带 MSYS 形 `third_party/sherpa-onnx`，`GOTMPDIR=C:\tmp\wisp255v3gotmp`；与在飞 `237-r1`（`internal/tools` 面）零文件重叠。
**起手锚**：HEAD＝`94398165`；交付态 sha ＝ roster `6b840071c183b3ed6470935508247037c416ed3787b22c1b50f839a56f2b55f2`／test `2195e669aab19be80643d0040be67c3fdbd716d082b4e04de25d60caf68dfa4e`；`git status --porcelain cmd/wisp`＝空。跑期编排者顺路落 `edf643c9`→`9e9964f8`→`2bb6394a`→`07a4d5ba`（frontend/**＋docs/reports/**＋probes/orch/**，与判面零重叠）；判面两枚文件每发窗口前后复核，全程停在交付态 sha。
**证物**：逐发原文＝`.scratch/wisp/probes/255/v3/logs/probe-P{1..6}-*.txt`；汇总＝`logs/probes-summary-v3.txt`；独立复量＝`verify-roster-255v3.py` ＋ `logs/static-verify.json`。

## §0 门禁（本腿自跑，⛔ 未引实现者日志）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| 本票全名册 | `go test -count=1 -v -run 'TestTicket255' ./cmd/wisp/` | 0 | 20 枚顶层全 PASS／0 FAIL；`ok … 23.186s`；窗口 09:27:09–09:27:39 +08 |
| 格式 | `gofmt -l` 两枚 | 0 | 空 |
| 独立静态复量 | `python verify-roster-255v3.py`（自写：复刻 evidenceCite＋`readsFieldAccess`，⛔ 不 import 被测件） | 0 | 14 枚 verdict cite 全 unique 且行号读数==现读；14 段全 match；尺① 27/20、尺② 23 枚/22 行、测试侧 2 行 |

## §1 四件判语（结论先行）

| # | 被裁件 | 判语 | 凭据一句话 |
|---|---|---|---|
| 1 | AC#6 两形判据（token 唯一命中当判据，行号只当读数） | **成立** | 自造 5 发：ⓐ 插行绿＋Logf（P1/P5a）、ⓑ 改 token 红逐名（P2）、二义红（P3）、错读数仍绿（P4）；旧尺缺陷面父状态实红（P6） |
| 2 | 六枚不唯一 token 削尖（唯一自裁点） | **可以接受** | 逐枚＝被引行真实子串（byte 级）、含原 token、同指一行；裸形今仍不唯一（削尖才是唯一性来源）；零放宽方向 |
| 3 | A878-① `sectionReadSites["Voice"]` 补两名读者＋四行去"唯一读者" | **成立** | 独立重写扫描逻辑跑全部 14 段：生产读者集==名册集合；两枚新读者读的是 cloud/realtime 键（`cfg` 取自 live Manager） |
| 3 | A878-② llm 证据行 `:96→:98` 并削尖 | **成立** | `panel_config_store.go:98` 现读＝`ChatModel:      cfg.LLM.Roles.Chat.Model,`，token 是其真实子串、文件内 unique、行号读数==现读 |
| 4 | 零放宽＋行号读数保留 | **成立（一处具名降级）** | 断言点 57→56（三处行号判据错→两处 token 判据错）；`citedRowsFloor=8` 未动；⛔ 无 `t.Skip`；尺①/② 改前=改后；测试侧 4→2，全部具名（见 §5） |

## §2 AC#6 两形 — 自造反证逐发（全部在交付态 sha 上；piggyback ⛔ 未用）

| 发 | 突变（⛔ 与实现者组合不重合） | 窗口（+08） | rc | 句 |
|---|---|---|---|---|
| P1 | `internal/observe/logging.go` 被引行 51 上方插一行注释，token 原样 | 09:29:14→09:29:30 | **0** | Logf：`row "privacy": recorded reading internal/observe/logging.go:51 drifted, the token "RedactPaths bool" now uniquely sits at line 52` ⇒ ⓐ 成立 |
| P2 | `cmd/wisp/logsink.go:149` 的 token 改掉（`Level: `→双空格，可编译） ⇒ 文件内 0 命中 | 09:29:41→09:29:54 | **1** | 红：`row "observe" cites cmd/wisp/logsink.go:149 for "Level: logSinkLevel" and the token is now nowhere in that file (the recorded line reads "…Level:  logSinkLevel})"` ⇒ ⓑ 成立且逐名点到该行 |
| P3 | 第三形：`logging.go` 加一行含 `RedactPaths bool` 的注释 ⇒ 2 命中 | 09:30:04→09:30:14 | **1** | 红：`…that token is not unique in the file (it hits lines [51 52])…sharpen the token` ⇒ >1 命中必红 |
| P4 | 反证「行号只当读数」：把名册 llm cite 的 **:98 写成 :96**（错读数），token 未动 | 09:30:30→09:30:44 | **0** | 仅 Logf 漂读（token 唯一命中 98），不红 ⇒ `:LINE` 确实只当读数 |
| P5a | 同一窗口：`internal/ball/ball_windows.go` 被引行 64 上方插一行，**新尺** | 09:30:55→09:31:19 | **0** | Logf 漂读 64→65 ⇒ 插入不打红（另选一枚被引文件，独立于实现者的 run.go/audio 组合） |
| P5b | 同一窗口换**旧尺**（`dcaff922^` 的 test 全文装入） | 同上 | **1** | 红：`row "ball" … that line now reads "// probe 255-v3 P5: inserted line above the cited ball option" - the evidence drifted` ⇒ 旧形缺陷面受控重演，且本发证明新形不是靠删尺换绿 |
| P6 | 父状态（roster＋test 都装 `dcaff922^`）跑旧尺 | 09:33:43→09:33:52 | **1** | 红：`row "llm" cites cmd/wisp/panel_config_store.go:96 for "cfg.LLM.Roles.Chat.Model", that line now reads "view := panel.SettingsView{"` ⇒ 交付前确有实红缺陷面（`96` 是 `view :=` 那一行） |

**还原证**（`sha256sum` 逐发 before==after、`git status` 空）：`logging.go ae6ddfe2…0a189e`（P1/P3 两发同值往返）；`logsink.go f06e2a5b…ae62c11`；`ball_windows.go 89c54360…5984d1`；roster `6b840071…b55f2`（P4）；test `2195e669…dfa4e`（P5 换回后）。备份＝`C:/tmp/wisp255v3backup/`；P6 装入的父对＝roster `367d4a43…278cb884`／test `5701ce4c…67fd275`。

## §3 削尖逐枚核定（起手读数＝`logs/inventory.txt`，本腿现量在 `static-verify.json`）

| row | 起手 token（命中数，实现者量） | 入册 token（削尖后） | 削尖后命中 | token 是被引行真实子串 | 方向 |
|---|---|---|---|---|---|
| llm | `cfg.LLM.Roles.Chat.Model`（2：98/99） | `ChatModel:      cfg.LLM.Roles.Chat.Model` | 1（98） | 是（含 6 空格对齐，逐字节） | 原 :96 已是 stale 读数（P6 实红）；削尖指回真身 98 |
| ball | `SizePx`（13） | `SizePx   int // configured orb size 44..72` | 1（64） | 是（截断于 `(0 = default 56)` 前，仍是子串） | 同指 64，无位移 |
| audio | `c.Audio.MicMutedDefault`（5） | `ra.mutedAtBoot = c.Audio.MicMutedDefault` | 1（256） | 是 | 同指 256，无位移 |
| privacy | `RedactPaths`（3：50/51/88） | `RedactPaths bool` | 1（51） | 是 | 读数 50→51：旧 50 是**注释行**，新指**声明行**——同主题、更实 |
| voice×4 | `c.Voice.Enabled`（3：253/257/262） | `ra.voiceEnabled = c.Voice.Enabled` | 1（257） | 是 | 同指 257，无位移 |
| agent（注释，不在尺内） | `cfg.Agent.PerToolTimeoutMS`（2：763/1014） | `PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS)` | 1（1014） | 是 | 注释性对齐，实现者已具名披露"注释不在尺内" |

**判语：可以接受。** 理由① 六枚裸 token **本腿现量今仍不唯一**（2/13/5/3/3/2 命中）⇒ 唯一性是削尖买来的，不是顺手改别处买来的；② 每枚削尖都**含**原 token 且锚在同一主题行（上表），不存在"改写成意思是那行"；③ 削尖只会**加严**（不唯一即红，P3 实证）；④ 唯一读数位移（privacy 50→51）是从注释行移到声明行；⑤ "未停手"可接受——AC#6 判据本身命令"不唯一⇒红、须削尖"，不削尖则该判据无一条 cite 可满足；实现者已把起手读数与候选留在件内、具名声明可退回（锚 §5.3），不构成未定义即停的违反。
**具名残留**：llm/ball 两枚 token 含 gofmt 对齐空格，纯重排格式会打红（fail-loud，逼重裁，不是漏洞；方向安全）。

## §4 A878 两处重裁

- **① voice 读者集**：本腿**自写**扫描（复刻 `scanReadFiles`/`readsFieldAccess`/排除项，⛔ 不调被测件）跑全部 14 段 ⇒ 段段 `scan==sectionReadSites`；`Voice == {cmd/wisp/panel_config_store.go, cmd/wisp/resident_audio_windows.go, internal/panel/config_handlers.go}`。两枚新读者的读点复核：`panel_config_store.go:105-110` 从 `cfg.Voice.CloudASRChain/CloudTTSChain/Realtime.*` 填 `SettingsView`（`cfg := s.mgr.Config()` 在 :92，每次调用现读）；`config_handlers.go:587-591` 读 `v.Voice.*`（`SettingsView`，由前者填充）。四行"唯一读者"表述已消失（`grep 'only production reader'`＝0）。四枚 voice 键行的 `扫描零命中` 是**键级**主张：本腿直查生产码里 `Voice.{TTS,Punctuation,WakeWord}` 唯一命中在 `internal/config/manager.go`（扫描刻意排除的自家包）⇒ 键级 no-reader 仍真。**判语：成立。**
- **② llm 证据行**：`:98` 行现读＝token 真身；token 在文件内 unique；`:LINE`==唯一命中行 ⇒ 读数与现读一致。**判语：成立。**

## §5 零放宽 ＋ 行号读数保留（摘尺不许当修尺的边界）

- 逐块读 `config_receipt_255_test.go` diff（4 个 hunk）：**移除的断言只有三处，全部是"行号判据"族**——① `unparseable line number`；② `n<1||n>len(lines)` 越界；③ 旧核心 `strings.Contains(lines[n-1], token)`。**换上两处 token 判据**（0 命中红／>1 命中红）。断言点 57→56（净 −1＝3 换 2）。`citedRowsFloor=8` 原样保留；`len(fields)<14` count guard 原样保留；⛔ 无 `t.Skip`、⛔ 无新跳过。
- 一处**开关外断言**具名：`TestTicket255ReceiptOmitsPanelFromTheImmediateSentence` 里 `contains "panel_resident_windows.go:207"` → `file 名 + [cfg.Panel.Width]`。它改前是**绿的**（不是换绿）；行号那一维确被放宽，但方向＝AC#6 本旨"断言不许再依赖行号"；token 的唯一性由名册尺独立钉住，`:207` 仍以"读数"形式在名册里（尺②计数未减）。判语：成立（具名降级，非粉饰）。
- 行号读数**一枚未删**：尺① 27 枚 `path:LINE`（唯一 20）与尺② 23 枚/22 行在父态与交付态**同数**（本腿两态各跑）；测试侧 4→2：删的二是 L440 注释里的历史指针（`:304`→内容锚＋`HEAD 67ab595d` 限定）与 L445 上述断言；留的二是 L81/L529 的 `panel_config_store.go:92`——本腿现量该行逐字 `cfg := s.mgr.Config()`，**活的、准确**。
- 旧尺红形＝P5b/P6 独立两发（⛔ 未引实现者 probe-a2/pre-run）。

## §6 具名没做 / 不开的门

- ⛔ 未跑整包 `go test ./cmd/wisp/`（314 死区，按令不拿整包当判据）；⛔ 未跑 `go vet`、d22scan、长度门（不属本腿门禁表；本腿零 Go 产码新增）；⛔ 未动任何在飞/别人的脏件；⛔ 未动三冻结件/`thresholds.go`/golden/`docs/PLAN.md`/`docs/specs/**`；⛔ 未 push。
- 突变窗口逐发已报（见 §2）；窗口内 `cmd/wisp` 与 `internal/{observe,ball}` 只本腿在写；`237-r1` 的 `internal/tools` 面零重叠。
- 未验的格：`panel cargo` 面真窗不属本票本格；`scanReadFiles` 按**字段名**取读（`v.Voice.` 计入）这一设计口径是既有件自述，本腿只核"名册集合==它量到的集合"，未替它扩口径。
