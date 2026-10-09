# 票 294 · AC#0 只读普查（腿 294-a1）

- 起手 `date` stdout：`2026-10-09 20:55:18 +0800`
- 身份：只读普查腿。射程＝只新建本目录下的 `.md`＋在票面 Progress log 追加一行；⛔ 零源码改动、⛔ 零翻框、⛔ 不动 `docs/**`。
- ★**这一把我没跑任何 `go` 命令**（`go build/vet/test/list` 全禁）：另一枚写腿 `296-r1` 正独占 Go 编译面并跑 `cmd/wisp` 整包名册。下面所有"能不能跑/读不读得到"的结论一律是**静态读码**判定，需现跑取数的判据具名归 `294-r1`/`296-r1`。
- ⚠ **行号在本把里会漂**：`cmd/wisp/resident_windows.go` 此刻正被 `296-r1` 改（同一枚产码调用在两把 grep 之间从 `:325` 漂到 `:336`）。所以本票一律锚**调用形状**，不锚行号（AC#1 原文也这么要求）。票面正文里的 `:217/:278/:303` 都已过期，见文末「票面与读数不符」。

---

## ⓐ `startResidentBall` 名册：真窗 vs 手搓 `&residentBall{}`

**尺（具名声明射程）**：`grep -rn "startResidentBall" --include=*_test.go cmd/`（只扫 `cmd/` 目录、只带 `--include=*_test.go`、判据锚在调用串 `startResidentBall(`；注释行 / AST 里的 `id.Name != "startResidentBall"` 字符串比较 / `t.Fatalf` 文案**都不算调用**）。子测试不单列（调用点本身即函数级）。手搓那一族的尺＝`grep -rn "residentBall{" --include=*_test.go cmd/`。

**真窗调用点（`rb := startResidentBall(...)`，逐枚）——共 10 枚：**

| 文件:行 | 实参形状 | 判 |
|---|---|---|
| `cmd/wisp/resident_approval_246_windows_test.go:314` | `(observe.NewRegistry(), nil, nil, nil)` | 真窗 |
| `cmd/wisp/resident_approval_live_246_windows_test.go:92` | `(reg, ra.vetoByEsc, nil, nil)` | 真窗 |
| `cmd/wisp/resident_approval_live_246_windows_test.go:202` | `(reg, ra.vetoByEsc, nil, nil)` | 真窗 |
| `cmd/wisp/resident_approval_live_246_windows_test.go:296` | `(reg, ra.vetoByEsc, nil, nil)` | 真窗 |
| `cmd/wisp/resident_hotkey_258_windows_test.go:278` | `(observe.NewRegistry(), nil, src, src)` | 真窗 |
| `cmd/wisp/resident_hotkey_258_windows_test.go:341` | `(observe.NewRegistry(), nil, src, nil)` | 真窗（变异控：hotReload=nil） |
| `cmd/wisp/resident_hotkey_258_windows_test.go:384` | `(observe.NewRegistry(), nil, src, src)` | 真窗 |
| `cmd/wisp/resident_hotkey_296_windows_test.go:158` | `(observe.NewRegistry(), nil, hotCfg, hotReload)` | 真窗（★票 296 在飞的文件） |
| `cmd/wisp/resident_hotkey_296_windows_test.go:203` | `(observe.NewRegistry(), nil, hotCfg, hotReload)` | 真窗（★同上） |
| `cmd/wisp/resident_hotkey_v1probe_test.go:57` | `(observe.NewRegistry(), nil, src, src)` | 真窗 |

`resident_hotkey_258_test.go:9/45`（注释）、`:55`（AST 字符串比较 `id.Name != "startResidentBall"`）、`:65/:68`（`t.Fatalf` 文案）＝**非调用**，已排除。

**手搓 `&residentBall{...}` 结构体字面量（不走装配根）——共 13 枚：**
`resident_approval_246_windows_test.go:69/95/331/364`、`resident_cancel_key_label_260r4_windows_test.go:89/131`、`resident_cancel_key_wording_260r3_windows_test.go:99/141/165/179`、`resident_mute_290_windows_test.go:92/205/220`。

**判定（承重）**：`attachMuteGate` 的 3 枚测试直调**全部**挂在手搓 `rb := &residentBall{}`（`resident_mute_290_windows_test.go:92→:93`、`:205→:206`、`:220→...`）与 `var nilBall *residentBall`（`:235`）上；上面 10 枚**真窗**里，**没有一枚在起窗后再调 `rb.attachMuteGate(...)`**。原因＝`attachMuteGate` 是装配根 `runResident` 体内的一步，不在 `startResidentBall` 内部；而这些用例都直呼 `startResidentBall`、从不进 `runResident`（`resident_ball_228_test.go:17` 逐字：`// Why a source walk and not a runtime call: runResident() never returns (it`）。⇒ **"接好了"这一跳在生产形状上今天无人验**：成立。

---

## ⓑ `attachMuteGate` 调用点枚数（锚在调用形状 `.attachMuteGate(`，非符号名）

尺：`grep -rn "attachMuteGate" cmd/wisp/`＝**总 6 行**；再按形状分：

- **产码调用（非测试）：1 枚** —— `cmd/wisp/resident_windows.go:336`（本把两次读数 `:325`→`:336`＝`296-r1` 并发编辑漂移）逐字：`rb.attachMuteGate(raudio.toggleMute)`。
- **测试直调：3 枚** —— `resident_mute_290_windows_test.go:93` `rb.attachMuteGate(ra.toggleMute)`、`:206` 同、`:235` `nilBall.attachMuteGate(func() (string, bool) { return "", false })`。
- **定义：1 枚** —— `cmd/wisp/resident_ball_windows.go:472` 逐字 `func (rb *residentBall) attachMuteGate(fn muteGestureFunc) {`；体内 `:478` 逐字 `rb.muteGate = fn`（带 `rb.muteMux` 锁，读回函数 `currentMuteGate()` 在 `:483`，同一把锁）。定义行以 `) attachMuteGate(` 出现，不含 `.attachMuteGate(`，故与调用形状天然可分。
- **注释：1 枚** —— `resident_ball_windows.go:469` `// attachMuteGate hands this host ...`（⛔ 不作枚数，正是票面"别锚符号名"的理由）。
- **AST 走查断言它在 `runResident` 体内：0 枚** —— 见下 ⓒ；全仓无一处以 `Sel.Name == "attachMuteGate"` 形状钉它。

⇒ **直调有（1 产码 + 3 测试），走查零。** 成立。

---

## 附：装配根里其它"直调有、走查零"的后置 setter / 注入挂钩（⛔ 本票只治静音，名册一次列全）

装配根 `runResident` 体内（`resident_windows.go` 现量约 :250–:336）的挂载/注入调用，逐枚判是否被"走 `runResident` 体的 AST 走查"钉住。现存的 `runResident` 体走查器只有 3 枚（`resident_ball_228_test.go`、`resident_hotkey_258_test.go`、`resident_approval_risk_256_windows_test.go`）；它们断言的符号名分别是：

- `228_test`：`ball.New` 构造、runResident 调用某球宿主（`:247`）、`stop` 的 defer（`:261/:266`）。
- `258_test`：`startResidentBall` 被调用且 `len(call.Args) >= 4`（`:55/:59`）。
- `risk_256`：`newResidentApprovalWithConfig` 恰 1 次（`:561`）、禁 `newResidentApproval()`（`:564`）、实参带 `*.DataDir`（`:550/:569`）。

| 装配根挂钩（现量行号仅示意） | 测试直调 | runResident 体走查 | 结论 |
|---|---|---|---|
| `startResidentBall(...)`（:250） | 有 | **有**（258/228） | 已钉 |
| `newResidentApprovalWithConfig(...)`（约 :200 前段） | 有 | **有**（risk_256） | 已钉 |
| `ra.bindBallHost(rb)`（:261） | 有（246/260r3/260r4 多枚手搓） | **零** | 欠账 |
| `rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)`（:273） | 有（`resident_approval_246:141`、`live_246:212/303`） | **零** | 欠账 |
| `startResidentTaskSource(rt, ra)` 的调用（:293） | 函数体被 `grant_writer_265` 走查，但**调用点**在 runResident 无人钉 | 调用点**零** | 欠账（半个） |
| `startResidentAudio(rt, rb.setAudioLevel)`（:311） | 手搓 `assembleCapture` 有（290/247） | runResident 体**零** | 欠账 |
| `rb.attachMuteGate(raudio.toggleMute)`（:336） | 有（290×3 手搓） | **零** | **本票目标** |

⇒ 除 `attachMuteGate` 外，装配根里同样"直调有、走查零"的后置 setter/注入挂钩至少 **4 枚**：`bindBallHost`、`RegisterShutdownHook(StepCancelTasks,...)`、`startResidentTaskSource`(调用点)、`startResidentAudio`。票面"欠账名册一次列全"这一问＝可答且非空。

---

## ⓒ AC#1 要的那枚"读 AST 的形钉"——仓里现成的走查长什么样（静态读，⛔ 未跑）

**模板＝`cmd/wisp/resident_hotkey_258_test.go`（untagged 文件，读源码不需 Windows build）。逐字摘关键段：**

- 解析：`:27-28` `fset := token.NewFileSet()` + `f, err := parser.ParseFile(fset, filepath.Join(".", "resident_windows.go"), nil, parser.ParseComments)`。
- 定位函数体：`:32-35` 遍历 `f.Decls`，`if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == "runResident" { return fn.Body, ... }`。
- 遍历调用：`:50-63` `ast.Inspect(body, func(n ast.Node) bool { call, ok := n.(*ast.CallExpr); if !ok { return true }; if id, ok := call.Fun.(*ast.Ident); !ok || id.Name != "startResidentBall" { return true }; calls++; if len(call.Args) >= 4 { hotArgs++ }; return true })`。
- 断言：`:64-68` `calls == 0` → `t.Fatalf`；`hotArgs == 0` → `t.Fatalf`。

**用的 `go/ast` API**：`token.NewFileSet`、`parser.ParseFile`（`parser.ParseComments`）、`ast.Inspect`、`*ast.FuncDecl`/`*ast.BlockStmt`、`*ast.CallExpr`、`*ast.Ident`。

**能不能直接照抄成 `attachMuteGate` 版本？判：不能纯照抄，差一枚形状。** 258 模板匹配的是 `call.Fun.(*ast.Ident)`（裸函数名 `startResidentBall`），而 `rb.attachMuteGate(...)` 是**方法调用**，其 `call.Fun` 是 `*ast.SelectorExpr{X: *ast.Ident{rb}, Sel: *ast.Ident{attachMuteGate}}`。⇒ AC#1 必须新增一层 `call.Fun.(*ast.SelectorExpr)` 判 `sel.Sel.Name == "attachMuteGate"`（这正是票面 ⛔"别锚符号名字面、锚形状"的落点）。三件事逐条静态可行性：
- ⓐ `== 1` 次：现量该调用在 `runResident` 体内**恰好 1 次**（`:336`，唯一产码枚），形状可数、可满足。
- ⓑ 接收者是 `startResidentBall` 返回值 `rb`、实参点名 `raudio.toggleMute`：静态成立——`rb` 于 `:250` `rb := startResidentBall(...)` 绑定；实参 `raudio.toggleMute` 是 `*ast.SelectorExpr{X:Ident{raudio}, Sel:Ident{toggleMute}}`，`raudio` 于 `:311` 绑定。⛔"不接受任意函数值"＝要判 `X` 那个 `Ident` 的字面名为 `raudio`、`Sel` 名为 `toggleMute`，AST 可表达。
- ⓒ 位置晚于 `startResidentAudio` 那次调用：可判——`ast.Inspect` 按源序访问，`sel.Pos()` 单调；或按 `runResident` 体的 statement 序列比 `attachMuteGate` 与 `startResidentAudio` 两 call 的先后。现量顺序 `startResidentBall(:250) < startResidentAudio(:311) < attachMuteGate(:336)` 满足时序前提。
- ⚠ 硬约束照票面：走查锚落形状，⛔ 不写行号锚（`resident_windows.go` 此刻正在漂，行号锚会被 `296-r1` 一插打死）。
- **落点文件建议**：新增 untagged `cmd/wisp/resident_mute_294_test.go`（照 258_test 的 `runResidentBody258` helper 复用/改名），⛔ 不改 258_test（其归属非本票）。

---

## ⓓ AC#2 真球窗形状 + "mute 那枚 id 的 IsLive"这把尺今天存不存在

**逐字摘 `cmd/wisp/resident_hotkey_258_windows_test.go` 用 `HotkeyReport` 现读的那枚用例：**
- `:266` `func Test258BridgeRebindsLiveKeysFromConfigEdit(t *testing.T) {`
- `:272-275` `writeConfig258(t, dir, `summon = "Ctrl+Alt+Z"` / `mute = "Ctrl+Alt+M"` / `cancel = "Esc"` / `panel = "Ctrl+Alt+P"`)` ⇒ **`mute = "Ctrl+Alt+M"` 确实已在夹具里写了**（票面现量对）。
- `:278` `rb := startResidentBall(observe.NewRegistry(), nil, src, src)` 起真球窗。
- **SKIP-LOUD 退路那一行（`:280`）逐字**：`t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the live rebind has no Win32 to land in: %q", rb.verdictStatusLineFor258())`。
- `:289` `boot := rb.b.HotkeyReport()`；`:290` 逐字 `if !boot.IsLive(hkSummon258) {`；`:293` `if got := boot.Live()[hkSummon258].VK; got != 'Z' {`。

**判"mute 那枚 id 的 IsLive"这把尺今天存不存在：半存在——方法在、id 在 cmd/wisp 不可达。**
- `HotkeyReport.IsLive(id)` 方法**存在**：`cmd/wisp/resident_hotkey_258_windows_test.go:290` 现读 summon；`internal/ball/hotkey_live_test.go:404` `if !b.HotkeyReport().IsLive(hkMute) {`——**internal/ball 自己那族用例已对 `hkMute` 读过 IsLive**。
- 但在 **`cmd/wisp` 包**里，`hkMute` 是 `internal/ball/hotkey_windows.go:37` 的**非导出常量**（`hkSummon=1, hkMute=2, hkCancel=3, hkPanel=4`），`cmd/wisp` 不能直接引用。现量 `cmd/wisp` 只镜像了 `hkSummon258 uint32 = 1`、`hkPanel258 uint32 = 4`（`:74/:75`），**没有 `hkMute` 的本地镜像常量**，也就**没有任何一枚 cmd/wisp 用例对 mute 的 id 调过 `IsLive`**。
- ⇒ **缺的仪器（具名）**：AC#2 需 (1) 在 cmd/wisp 新增一枚本地镜像常量（如 `hkMute258 uint32 = 2`，须与 internal/ball 的 `hkMute=2` 对齐并具名注释来路），(2) 在真球窗形状上把 `boot.IsLive(hkMute258)` 读出来。球窗形状＋`mute=Ctrl+Alt+M` 夹具都在，唯独"为 mute 读 IsLive"这一动作与那枚 id 的可达性缺失。成立。

**skip 路纪律（本台机器具名）**：票面 ⚠ 已定式"skip 不算读数"。编排者 20:5x 现量**本机起得来真窗**（票 296 夹具在本机跑出改前红/改后绿），⇒ `SKIP-LOUD` 在本机**不是退路、不许当交付**。AC#2 若最终落在真窗上，必须在本机出**实读**（红/绿具名），不能停在 skip。

---

**（本 AC#0 件未跑任何 `go` 命令；以上全是静态读码。需现跑者见 AC#1–3 件末尾「我没跑的尺归谁」。）**
