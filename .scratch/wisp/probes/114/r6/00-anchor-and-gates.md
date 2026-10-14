# 114-r6 探针锚 — 票 114 AC#10「可达性钉认第三装法」（写腿；实现者）

- 身份／时刻：`114-r6`，2026-10-14 16:3x–17:0x +08；起手 HEAD `6c5515e8`，产码 commit 锚 `14e4961d`，产码 sha `3b6f170f`。
- 窗口：他腿 `201-a1` 在飞并落件（`64bf9b29`/`14e4961d`）；`frontend/**` 被界面侧持续写（`frontend/fixtures/composer-states.html` 16:17 重生、`tokens.generated.css` 未提交脏）；本波 `internal/panel` 只有本腿一枚写腿。
- 面：`internal/panel` 测试面（仪器）。⛔ 未触 `cmd/wisp/**` 产码、⛔ 未触 `frontend/**`、⛔ 未动契约／冻结件、零 push。
- 未定义即停核：`w.Bind` 现读与 `114-v1` 报的**逐字对上**（`:80` const／`:842` Bind），换名突变做得圆（见 §2）⇒ 未触发停手。

## 1. 事实链（全部现读，逐字）

- 真门＝`cmd/wisp/panel_host_windows.go:80` `panelDispatchBinding = "wispDispatch"` ＋ `:842` `w.Bind(panelDispatchBinding, func(raw string) string {`；`:892` 另有测试门 `w.Bind("wispProbeRT", ...)`（literal，**不入 reviewed**）。
- 页真调用＝`frontend/src/lib/panel.ts:350` `const door = window.wispDispatch;`（声明在 `:155`；`$(window.wispDispatch)` 在 `:329` 注释里）。
- 旧可达性钉＝`internal/panel/composer_test.go:441` `hostBridgeCallRe = regexp.MustCompile(`\.\s*postMessage\s*\(`)` ⇒ 第三装法不在其射程（`.postMessage(` 只数页侧调用点）。
- 仓内他尺：`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（253-r1，windows 整包内、cmd/wisp 自读自，**不读页**且不在本队可跑面）——本腿不重造它，只补它射程外的"页—宿主同名"能力钉。

## 2. 落了什么（产码 `3b6f170f`；单枚新件 530 行）

`internal/panel/bound_dispatch_door_114_test.go`／`TestTheBoundDispatchDoorIsReachableFromThePage`：把第三装法当**能力**读盘定钉，两侧同名：

- **host 侧**：走遍 `cmd/wisp` 产码（AST 解析，跳过 `_test.go`/`testdata`/`logs`），取每个 `w.Bind(<名>, 回调)` 调用；`<名>`＝字符串字面量或**包级字符串常量**（经包常量表解析，含 offset 供种形用）。①即此：产码里那个"常量绑的门"存在性被断言；`ws.Bind(run.paths)`（`:311`，单参）记为 skipped 不判。
- **page 侧**：该常量值必须是 `frontend/src/lib/panel.ts` 以 `window.<名>` 读到的属性，且**整棵 `frontend/` 里该门的读只许在 panel.ts**（注释按 `//`／`/* */`／`<!-- -->` 剥掉，`panel.ts:329` 的 prose 不能顶数）。
- **floor**：reviewed 集为空＝红（空读不等于绿）；绑定名解析不出＝红；两枚域名各点名。
- **自证种形**（同尺、仓内零突变）：用刚读到的真字节按本尺自己找到的 offset 造三种突变（见 §3）。

判读对照（派单三向）：①＝host 侧 exists 断言；②＝三突变红且点名；③＝既有 `.postMessage(` 牙不动（log 01 "red as required" 行逐字）。

## 3. 三向读数（逐字；全在 `logs/`）

- **基线绿**（live，log 01；rc=0）：
  `bound door held: 36 cmd/wisp production file(s), 2 w.Bind door site(s) [cmd/wisp/panel_host_windows.go:842 binds "wispDispatch" via panelDispatchBinding; cmd/wisp/panel_host_windows.go:892 binds "wispProbeRT" via "wispProbeRT"], 1 skipped non-door .Bind call(s); 27 window.* read site(s) under frontend/, door-relevant: [frontend/src/lib/panel.ts:161 window.wispBridge; frontend/src/lib/panel.ts:161 window.chrome; frontend/src/lib/panel.ts:350 window.wispDispatch]`
  既有牙：`renderer door held: 113 files scanned, 2 host call sites ...` ＋ `red as required: 3 outside call sites, ...`（不动）。
- **突变红**（仓外真种盘 `D:\tmp\wisp114r6\tree`＝`git archive HEAD`；A713：overlay 换不到 `os.ReadFile`，故真盘；log 02）：
  - m1 把 `w.Bind` 摘掉〔span `[42508,42619)` 换 `error(nil)`〕⇒
    `cmd/wisp's production sources bind no door through a w.Bind(<package-level string constant>, callback) call, so the page's door has no host-side receiver: ... w.Bind sites read: cmd/wisp/panel_host_windows.go:889 binds "wispProbeRT" via "wispProbeRT"; page window.* reads: ... panel.ts:350 window.wispDispatch`（FAIL rc=1）
  - m2 把绑定名换掉〔`:80` `"wispDispatch"`→`"wispDispatchX"`〕⇒
    `cmd/wisp/panel_host_windows.go:842 binds "wispDispatchX" through w.Bind, and frontend/src/lib/panel.ts never reads window.wispDispatchX ...`（FAIL rc=1）
  - m3 把页读换掉〔panel.ts 两处 `window.wispDispatch`→`window.wispDispatchX`，宿主不动〕⇒
    `... binds "wispDispatch" through w.Bind, and frontend/src/lib/panel.ts never reads window.wispDispatch - ... The page's window.* reads: ... panel.ts:350 window.wispDispatchX`（FAIL rc=1）
  - **对照**：同一棵突变树上旧 `.postMessage(` 尺（`TestTheRendererHoldsExactlyOneDoorToTheHost|TestPlantedRendererDoorShapesGoRed`）三发全 **PASS**（`113 files scanned, 2 host call sites`）——m3 更是页面侧改动、仍在其射程外＝"旧射程盲、新钉有牙"现场。
- **抽核**（log 03）：`internal/panel` 整包同树 A/B＝A（有新钉）`225/142/7` vs B（移走新钉）`223/141/7`，**红集 diff SAME SET**（树 7 红＝live 既有 5＋2 枚树工件：`TestReadGit…` 树无 `.git`、`TestComposerRenderFixture…` 导出件 CRLF 而该尺分隔符写字面 `" -->\n"`）；live 整包＝`225/144/5/0 rc=1`，5 红与 v1/a5 记录**同集**（新钉在此绿）。

## 4. 门禁（现跑）

build `./internal/panel ./cmd/wisp` rc=0；`go vet ./internal/panel/` rc=0；`gofumpt -l` 本件零行 rc=0；`sh scripts/d22scan.sh` clean rc=0（ban #8 internal/=556 含本件）；长度门 `VERDICT GREEN`（tracked 12855，over-budget 60=roster 60，本件不在超预算内）。仓内零突变核：`cmd/wisp/panel_host_windows.go` sha256 `45c1059b…`、`frontend/src/lib/panel.ts` `2b7a9fd0…`、`composer_test.go` `d2b2e848…`（起手＝收手同值）。

## 5. 本波⛔没做（具名）

- ⛔ 未动 `cmd/wisp/**` 产码（`:842` 的 `w.Bind` 是既有的，一个字节未改）；⛔ 未动 `frontend/**`（`panel.ts:146/329` 注释自陈两种装法那半按派单归前端）。
- 未跑：`go test ./cmd/wisp`（314 死区，派单禁）；POSIX 三合包；`bash -n` 两脚本（本波未触脚本）。
- 未跑真窗／未做 AC#10 之外任何格的翻转（⛔ 框未翻、`Status:`/票名未动）。
- 待裁：`v1` 口径的**非实现者**复核本钉（AC#10 判语+翻勾仍归编排者／验收腿）。
