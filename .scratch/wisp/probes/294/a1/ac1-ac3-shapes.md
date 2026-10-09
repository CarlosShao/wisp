# 票 294 · AC#1/AC#2/AC#3 落点形状只读判定（腿 294-a1）

- ★**这一把我没跑任何 `go` 命令**（`go build/vet/test/list` 全禁）：`296-r1` 正独占 Go 编译面跑 `cmd/wisp` 整包名册。下面全是静态读码；"能不能跑成/红绿读数"归 `294-r1`/`296-r1`。
- ⛔ 本件不写用例、不改产码，只判"要动哪几枚文件、要不要新文件、撞不撞别人在改的件、缺什么仪器/世界"。

---

## AC#1 装配根那一跳的形钉（读 AST）

**落点文件**：
- 新建 untagged `cmd/wisp/resident_mute_294_test.go`（照 `resident_hotkey_258_test.go:25-68` 的 `runResidentBody258` 走查 helper 复用/改名，读 `resident_windows.go`，不需 Windows build）。
- ⛔ 不改 `resident_hotkey_258_test.go`（归属 258），⛔ 不改产码。

**撞不撞正在被改的件**：走查**读**的 `resident_windows.go` 此刻被 `296-r1` 在写（同枚产码调用本把从 `:325`→`:336`）。所以 AC#1 **必须锚形状、不锚行号**（票面硬约束已同此），否则 296 一插就把行号锚打死。新测试文件本身无人碰，不与 296 冲突。

**静态可行性（三件事，见 AC#0 件 ⓒ）**：ⓐ 体内 `attachMuteGate` 调用 `== 1`（现量恰好 1）；ⓑ 接收者 `rb`=startResidentBall 返回值、实参点名 `raudio.toggleMute`（`:250`/`:311` 绑定，AST 可判 `SelectorExpr{X:Ident raudio, Sel:toggleMute}`）；ⓒ 位置晚于 `startResidentAudio` 那次调用（现量 :250<:311<:336，`Pos()`/语句序可判）。⚠ 258 模板匹配的是 `call.Fun.(*ast.Ident)`（裸函数名），`rb.attachMuteGate` 是方法调用、`call.Fun` 是 `*ast.SelectorExpr` ⇒ **不能纯照抄，要新增一层 SelectorExpr 判定**。

**缺的仪器**：无新 seam，纯走查即可（`runResident` 体从不返回、不能直调，见 `resident_ball_228_test.go:17`）。⇒ 判据可满足、世界已存在（源码即在树里，`parser.ParseFile` 读得到）。**需现跑者**：走查跑红/跑绿由 `294-r1` 出（我只静态确认形状可满足）。

---

## AC#2 球侧注册那一跳的形钉（真球窗现读 IsLive）

**落点文件**：
- 新建 windows 文件 `cmd/wisp/resident_mute_294_windows_test.go`，**复用**（不是编辑）`resident_hotkey_258_windows_test.go:266-294` 的真窗形状：`startResidentBall(observe.NewRegistry(), nil, src, src)` + `writeConfig258`（其 `:273` 已含 `mute = "Ctrl+Alt+M"`）+ `rb.b.HotkeyReport()`。
- ⛔ 不改 `resident_hotkey_258_windows_test.go` 本体、⛔ 不改 `resident_hotkey_296_windows_test.go`（296-r1 在飞）。

**撞不撞正在被改的件**：新文件不撞；但 AC#2 落真窗 ⇒ 与 `296-r1` 抢**同一枚编译/运行面**（都要 build + 起窗）。票面排程本票在 293 之后、295 之前，且"⛔ 不与任何落地腿混批" ⇒ 实读必须等编译面到手。

**这把尺今天存不存在**（见 AC#0 件 ⓓ）：
- `HotkeyReport.IsLive(id)` 方法**存在**（`258_windows_test.go:290` summon 在用；`internal/ball/hotkey_live_test.go:404` 已对 `hkMute` 读）。
- 但 `hkMute` 是 `internal/ball/hotkey_windows.go:37` 的**非导出常量**（hkSummon=1/hkMute=2/hkCancel=3/hkPanel=4）；`cmd/wisp` 只镜像了 `hkSummon258=1`、`hkPanel258=4`（`258_windows_test.go:74-75`），**没有 mute 的本地镜像常量、也没有任何 cmd/wisp 用例对 mute 的 id 调过 IsLive**。
- ⇒ **缺的仪器（具名）**：需在 cmd/wisp 测试里补一枚 `hkMute` 的本地镜像常量（值须对齐 internal/ball 的 `2` 并注释来路），才能把 `boot.IsLive(hkMute258)` 这把尺架起来。球窗形状＋`mute=Ctrl+Alt+M` 夹具都已存在，唯"为 mute 读 IsLive"这一动作缺失。**世界造得出来**（真窗可起、热键可注册），只缺那枚 id 的名字与一次现读。

**skip 纪律（本机具名）**：`258_windows_test.go:280` 的退路逐字＝`t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the live rebind has no Win32 to land in: %q", ...)`。票面 ⚠ 已定式"skip 不算读数"。编排者 20:5x 现量本机起得来真窗（296 夹具跑出改前红/改后绿）⇒ **SKIP-LOUD 在本机不是退路、不许当交付**；AC#2 必须在本机出实读，红/绿具名。

---

## AC#3 `executed` 的另一半（M2 那一格）

**happy path 上有没有人断 `executed == true`？判：没有。**

**`executed` 断言名册**（尺＝`grep -rn "executed" cmd/wisp/ internal/audio/`；逐枚分类，仅列与静音那枚 bool 相关的，其余标"无关"）：

| 文件:行 | 逐字/内容 | 类 |
|---|---|---|
| `cmd/wisp/resident_ball_windows.go:89` | `type muteGestureFunc func() (outcome string, executed bool)` | 定义 |
| `cmd/wisp/resident_ball_windows.go:86` | 注释（executed=false 含义） | 注释 |
| `cmd/wisp/resident_ball_windows.go:527` | `outcome, executed := fn()` | 赋值（产码，muteGesture 体内） |
| `cmd/wisp/resident_ball_windows.go:528` | `if !executed {` | 分支（产码） |
| `cmd/wisp/resident_audio_windows.go:152` | 注释 `// executed=false is the shape...` | 注释 |
| `cmd/wisp/resident_mute_290_windows_test.go:196` | `outcome, executed := ra.toggleMute()` | 赋值（测试，无门用例） |
| `cmd/wisp/resident_mute_290_windows_test.go:197-198` | `if executed { t.Fatalf(...) }` | **断言，但只断 executed==false（无门那一支）** |
| approval_reply_201:474 / instructions_200r2:7 / panel_resident_windows.go:15 / resident_approval_246:128 / live_246:192,281 / resident_task_source_246:120 / live_246:90,161 / resident_windows.go:190 / run_mode101:311 | 均为"卡片/任务 executed"字样的文案或注释 | 无关 |

- `internal/audio/` 对 `executed`＝**零命中**：这枚 bool 的命名与返回全在 `cmd/wisp`（`residentAudio.toggleMute` 于 `resident_audio_windows.go:156`）。
- **toggleMute 的 happy 返回**（M2 要翻的那两枚 `true`）逐字：`:172` `return "已静音：采集已关闭，设备未打开（再按一次取消静音）", true`、`:174` `return "已取消静音：设备已交接，采集线程在跑；...", true`。
- ⇒ **唯一断言 `executed` 的用例（`:197`）断的是"无门 ⇒ false"；无任何用例在门真被拧时断 `executed == true`。** 与票面 M2（把 :172/:174 改 false ⇒ 六枚全绿、`rc_m2=0`）一致。

**缺的仪器/世界**：需要一枚 happy-path 用例——voice.enabled=true 下 `assembleCapture` 造出**真门**（`:157` 的 `ra.gate == nil` 分支证明两态可造），调 `ra.toggleMute()`（或经 `rb.attachMuteGate`+`rb.muteGesture`），断 **`executed == true`** ＋ 门自身读数真翻了（`gate.Muted()`/`gate.Open()`，票面 ⛔"不许读返回值字样"、M1 已证读返回值会被"句子里写了成功"骗）。判据反形正控＝改 `:172/:174` 两 true→false 指名用例必须红。
- **世界造得出来吗**（本仓夹具）：造真门的世界**已在票 290 夹具里**（`bootAudioRuntime`+`writeAudioConfig`+`assembleCapture`+`newRealCaptureSource`；`resident_mute_290_windows_test.go:189-191` 就是这套），只是当前没人对 happy 支断 `executed`。⇒ 属"仪器躺着、没指向这一跳"，不需新造 seam。

**落点文件**：新建 `cmd/wisp/resident_mute_294_windows_test.go`（与 AC#2 同文件，happy 用例）；⛔ 不改 290 既有断言（AC#4 禁放宽）。

**撞不撞正在被改的件（重要）**：AC#3 的**反形正控**要临时把 `cmd/wisp/resident_audio_windows.go:172/:174` 的 `true` 翻 `false`。该文件＝**票 295 的同写面**（票面"排程与串行：与票 295（同写 cmd/wisp/resident_audio_windows.go）不同批"）。⇒ AC#3 落地腿与 295 **必须串行**，且翻 `true→false` 只在红/绿验证时临时改、验后必须复原（⛔ 不留改动、⛔ 不为此放宽断言）。

---

## 「我没跑的尺归谁」（本腿禁跑 Go 命令）

| 判据 | 谁必须现跑 | 为什么我不能 |
|---|---|---|
| AC#1 走查跑红/绿（attachMuteGate 计数=1、时序、实参点名） | `294-r1`（落地腿，持 Go 面时） | 需 `go test` 执行；`296-r1` 正占编译面 |
| AC#2 真球窗 `boot.IsLive(hkMute)` 实读 | `294-r1`/`296-r1`（本机可起窗） | 需 build+起窗；且 hkMute 镜像常量尚不存在 |
| AC#3 `executed==true` 用例 + M2 反形（:172/:174 true→false ⇒ 指名红） | `294-r1`（与 295 串行后） | 需 `go test`，且涉及产码临时改 |
| AC#5 门禁四数（build/gofumpt/test 双发改前后/d22scan/tasklist） | 落地腿 `294-r1` | 我全程禁跑 go，且不得翻框 |

## 「票面与我读数不符之处」（⛔ 我未自改票面正文）

1. **`:303` 已漂**：票面现量逐字"`cmd/wisp/resident_windows.go:303` 生产挂载点"。本把静态读数＝该调用现落 **`:336`**（且我两把 grep 之间它从 `:325` 变 `:336`，＝`296-r1` 正在并发写该文件）。票面"建球 :217 / 建门 :278"亦同步过期（现量约 :250 / :311）。**次序关系（球<门<挂载）仍成立**，AC#1 时序断言不受影响，但**所有绝对行号引用都应作废、改锚形状**（票面本身已如此要求）。
2. `resident_mute_290_windows_test.go:93/:206/:235` 三枚直调、`resident_ball_windows.go` 的 `attachMuteGate` 定义体——本把读数与票面**逐字对得上**，无不符。
3. AC#2 `258_windows_test.go` 的 `SKIP-LOUD` 行：票面现量指 `:278/:289-291`；本把读数 `startResidentBall` 在 `:278`、`boot.IsLive` 在 `:290`、`SKIP-LOUD` 在 `:280`——行号小漂、内容对得上。
