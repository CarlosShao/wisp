# comment-truth-2 / 01 — 15 枚"已过期"产码注释的归口与分档表

锚全部在 **HEAD `7a452d0aca1bf5157a97fb5329ca6eb38b9f0694`** 现取（`git show HEAD:<path> | grep -nF`／`sed -n`），⛔ 不抄表一。
档位三值＝〔只欠文案〕／〔要动产码〕／〔要新判据〕。撞钉＝`git grep -nF "<逐字串>" HEAD -- '*_test.go'`，**15 枚全部 0 命中**（原始输出在 `02` §A）。
本程不写码、不翻框、不立票。

| 号 | HEAD 现量锚 | 过期那句（逐字） | 档位 | 凭据（哪一行代码/哪把尺） | 归口 |
|---|---|---|---|---|---|
| P01 | `HEAD:internal/panel/composer_dispatch.go:50` | `// belongs to the host; and it has NO production caller yet. The honest state of` | 〔只欠文案〕 | `.Handle(` 产码两处＝`cmd/wisp/panel_host_windows.go:821`、`cmd/wisp/panel_inbound.go:163`；入口 `cmd/wisp/main.go:115 case "panel-inbound"` 现存 | **票 33**（未勾 13 枚）——`AC#9（片 B：入向有具名生产听众）` 就是承接格 |
| P02 | `HEAD:internal/agent/approval/doc.go:26` | `// Wiring still owed by ticket 12 (this package has NO production caller yet,` | 〔只欠文案〕 | `approval.New(` 产码两处＝`cmd/wisp/run.go:612`、`cmd/wisp/resident_approval_windows.go:369`；下面五条 bullet 的第 1 条（`compose approval.New(...) as tools.Options.Gate`）已在盘上 | **票 12**（句里点名 owner）——但 12 的 3 枚未勾框＝Idle SLO／人眼验收／3s 回落，**无一承接** |
| P03 | `HEAD:cmd/wisp/panel_inbound.go:11` | `// line and hands the bytes to (*panel.ComposerDispatch).Handle. Nothing else in` | 〔只欠文案〕 | 除本腿外多一处＝`cmd/wisp/panel_host_windows.go:821`；被引仪器 `cmd/wisp/panel_inbound_33_test.go:183 inboundHandleCallSites33` **存在**，但射程只有 `cmd/wisp` 目录（AST 扫非测试 `.Handle(`），⛔ 覆盖不了注释说的 "this repository" | **票 33** `AC#9` 同 P01 |
| P04 | `HEAD:internal/panel/git.go:76` | `// ComposerRequest) has zero production callers in this tree - measured, with the` | 〔只欠文案〕 | 通路已在＝P01 那两处 + `main.go:115`；被引普查件 `docs/evidence/s1/181-186-git-detection-census-c1.md` **存在**；同文件 `:19` 点名的仪器 `TestGitDimensionHasNoModelCallableTool` **存在**（`internal/panel/git_test.go`） | **票 186**（`:79-80` 自指"由票 186 换掉"）——186 未勾 9 枚，最接近的是 `AC#2`（定切换方法＋写进契约），**不是改注释格** |
| P06 | `HEAD:cmd/wisp/config_readers_255.go:43` | `// as its accessor - and TierOf had ZERO production callers (measured at HEAD` | 〔只欠文案〕 | `config.TierOf(` 产码唯一一处＝**同文件** `:235 if tier, ok := config.TierOf(name); ok { // TierOf's first production caller`；定义 `internal/config/tiers.go:93` | **票 255**（未勾 2＝AC#1/AC#4），两格都不是注释更正格 |
| P07 | `HEAD:cmd/wisp/config_reload.go:11` | `//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a` | 〔只欠文案〕 | 产码调用点两处＝`cmd/wisp/config_reload.go:153`、`cmd/balldebug/main.go:243`；定义 `internal/config/manager.go:142` | **票 223**（未勾 2＝AC#2/AC#7），非注释格 |
| P08 | `HEAD:cmd/wisp/config_reload.go:14` | `//   - Manager.ConfirmLocked - zero assignments outside *_test.go, and` | 〔只欠文案〕 | 赋值点＝`cmd/wisp/config_reload.go:114 rt.mgr.ConfirmLocked = rt.confirmLockedLoosening` | **票 223**，同上 |
| P09 | `HEAD:cmd/wisp/config_reload.go:18` | `//   - Manager.OnRestartPending - zero assignments outside *_test.go too, so the` | 〔只欠文案〕 | 赋值点＝`cmd/wisp/config_reload.go:115`；读侧 `internal/config/manager.go:201-202` | **票 223** `AC#7`——⚠ 那格的判据原文正是"`OnRestartPending` … 在生产里有赋值点"，**今天已被产码满足却仍未勾**（勾要非实现者裁，本程不动） |
| P10 | `HEAD:cmd/wisp/firstrun.go:11` | `// no production path ever did was CALL that pair once, at the moment a user` | 〔只欠文案〕 | 那一对已由 `cmd/wisp/firstrun.go:82 config.SaveFile(cfgPath, config.NewDefaults())` 调用；入口 `cmd/wisp/run.go:245 ensureFirstRunConfig(...)` | **票 198 已 `-done`（未勾 0）** ⇒ 在册但**不可承接**；需新立或并入，本程不裁 |
| P11 | `HEAD:cmd/wisp/resident_approval_windows.go:282` | `// wins) but has no production caller; a second assembly in one process belongs` | 〔只欠文案〕 | 调用点＝`cmd/wisp/resident_task_source_windows.go:309 ra.bindResidentGrantLedger(run.session)`；⚠ **同文件 `:322` 已改口**（`... called by startResidentTaskSource`）＝文件内自相矛盾；且这枚事实有仪器钉着＝`cmd/wisp/resident_grant_writer_265_windows_test.go:465` AST 读 bind 位 | **票 246**（未勾仅 1＝AC#8，不承接） |
| P12 | `HEAD:cmd/wisp/approval_reply.go:12` | `// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers` | 〔只欠文案〕 | `.Veto(` 产码两处＝`approval_reply.go:363`、`resident_approval_windows.go:612`；`.DecideFromNative/.DecideFromPanel` 四处＝`internal/agent/approval/replies.go:328/:408/:413/:437` | **票 201**（未勾 5，`AC#1`＝"三枚入口至少两枚有生产调用者"是承接格）。⚠ **已登记**：票 201 票面 `:73` 逐字写着「⛔ 票面 §10 那把尺的读数（"三枚入口生产零调用者"）已过期」⇒ **⛔ 不算本程新发现** |
| P13 | `HEAD:internal/risk/assessor.go:29` | `// not yet wired by tickets 18/19) leave the corresponding rule DORMANT — an` | 〔只欠文案〕 | 三枚 `With*` 全已装配＝`internal/tools/bridge.go:212-213`（canonicalizer+classifier）、`internal/tools/bridge.go:925-926`（`assessorFor`；`:927` 带 TaintDetector）、`cmd/wisp/panel_assets.go:66`；票 18/19 均已 `-done` | 票 18/19 **已 `-done`** ⇒ 无在册承接格；票 259（授权绑定层）不是一件事。**⚠ 已登记＝台账 `A720` §3 ⑤**，那里明写"落哪一票待我先查重再裁"——本程只交读数，裁定归编排者 |
| P14 | `HEAD:cmd/wisp/resident_task_source_windows.go:14` | `// askConfirmation (:205) with ZERO product callers - only this package's own` | 〔要动产码〕 | 两处引用行号**全漂**：`AskOnTaskRoot` 现在 `resident_approval_windows.go:697`（注释写 `:261`）、`askConfirmation` 现在 `:641`（注释写 `:205`）；`askConfirmation` 现有 1 枚包内调用者 `:700`，但它**住在 `AskOnTaskRoot` 体内**；尺 `[A-Za-z0-9_]\.AskOnTaskRoot\(` 在 HEAD 非测试产码 **0 命中（rc=1）** ⇒ 那一跳真没接，不是措辞问题 | **票 246**：`AC#7` **已经是 `[x]`**（完成判据原文＝"常驻那条腿真起一条任务管线 … 四步走通"），只剩 `AC#8` 且不承接 ⇒ **无承接格** |
| P15 | `HEAD:cmd/wisp/resident_windows.go:248` | `// the way out - and AskOnTaskRoot / askConfirmation still had zero product` | 〔要动产码〕 | `:249` 那句 `This call is the caller: it reads` 指的调用是 `:260 src := startResidentTaskSource(rt, ra)`，而它不触 `AskOnTaskRoot`（同上尺 0 命中）⇒ "票已勾、断言未真" | 同 P14 → **票 246**（AC#7 已勾、AC#8 不承接）。⚠ **已登记＝台账 `A720` §3 ①②**（编排者已复认并定档〔建了但没接〕）⇒ 不算本程新发现 |
| P39 | `HEAD:cmd/wisp/run.go:333` | `// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does` | 〔只欠文案〕 | 半枚：结论分句**仍成立**（`modeWrites` 尺＝`run.go:337` 定义、`:674` 赋值、`panel_inbound.go:248/:272` 是**同名局部量**⇒ 字段读侧 0）；理由分句**已过期**（`ParseComposerRequest(` 产码调用＝`internal/panel/composer_dispatch.go:155`，其宿主 `Handle` 有两处产码调用点，`main.go:115` 子命令在册） | **票 114**（未勾 9）＋票 33/35。⚠ 具名附加：票 114 **文件名本身**逐字是 `composer-request-has-no-production-caller…`，那是同一族过期断言的票面级载体（本程不改票名、不立票） |

## 汇总

- 档位：**〔只欠文案〕13 枚**（P01 P02 P03 P04 P06 P07 P08 P09 P10 P11 P12 P13 P39）／**〔要动产码〕2 枚**（P14 P15）／**〔要新判据〕0 枚**。
- 撞钉：**15 枚逐字串在 `*_test.go` 里 0 命中**；但有 **4 枚"改文案有形状约束"**（不是逐字钉，是仪器扫同一文件）：
  P01/P03/P39 → `internal/panel/composer_dispatch_test.go:375`（非注释行的 `"panel.*"` 字面量）与 `:708`（**整份文件含注释**扫 `checkoutBranch/changeRepo/repoPicker/branchSelect/vcs.switch` 五个禁词）；
  P10 → `cmd/wisp/firstrun_257_test.go:430-438`（钉 `\nfunc ` 计数＝1 与 `ensureFirstRunConfig` 签名逐字）；
  P15 → `cmd/wisp/resident_hotkey_258_test.go:73-79`（钉 `resident_windows.go` 里含 `config.LoadFile`）；
  P11 → `cmd/wisp/resident_grant_writer_265_windows_test.go:465`（AST 读 bind 位，注释不进 AST）。
- **〔引用假凭据〕：0 枚**。15 枚里凡点名测试件／证据件／commit 的，被引对象**全部存在**（清单见 `02` §C）。两处**引用缺陷**但**不是假凭据**，具名：
  ① P10 的 `loader.go:238` 指 `SaveFile` ⇒ 现量在 **`internal/config/loader.go:252`（漂 14 行）**，符号在、行号假；
  ② P12 的 `:16` 把票 201 现量表转写成 `"三种人能答复的入口 - 生产零调用者"` ⇒ 票 201 `:10` 原文是 `| 三种"人能答复"的入口 | **生产零调用者** |`，**对象在、逐字不符**（近似转写）。
- 名册分类**无一枚不成立**：15 枚在 HEAD `7a452d0a` 复跑后仍判"已过期"，⛔ 不需报回降档。两处**档位内的限定**具名（不改档）：
  P06 与 P12 的"过期"只到**现在式读法**——两句都自带锚点/日期（`measured at HEAD 8a3790f0`／`The 09-29 census …`，且 commit 与证据件均存在），按历史读法仍成立，这与表一把 `config_reload.go:7` 判"仍成立（历史）"的口径同形；
  P39 只过期**理由分句**、结论分句今天仍对（表一原文如此，本程复跑一致）。
