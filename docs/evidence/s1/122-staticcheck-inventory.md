# 票 122 · 波 1 · staticcheck 现量清单（122-wave1，2026-11-12）

- 执行腿：122-wave1（写腿，只 commit 不 push）
- HEAD：`095427398`（A1771）
- staticcheck：**2026.2.2 (0.8.2)**，本地二进制 `D:\work\base\gopath\bin\staticcheck.exe`，与 `ci.yml:465` 钉的 `@2026.2.2` 一致
- 跑法：根 module 的 `go list ./...` 全 42 包，**剔除 `github.com/CarlosShao/wisp/build`**（`build/` 目录被 `.gitignore:14` 忽略、CI 上不存在，本地那份是编译坏的孤儿件，11 行 `undefined` 噪音会污染分母；剔除后与 CI 形状对齐）。两个 GOOS 各跑一遍。
- 派单预期「约 54 条：U1000 约 47（internal/ball 约 38）」＝**linux（CI）形状，现量逐位对上：54 条 finding、U1000 47、internal/ball 38**。

## 修前分组计数（改码之前现跑）

| 检查码 | linux（CI 形状） | windows（本机形状） |
|---|---|---|
| U1000 | 47 | 85 |
| S1011 | 2 | 2 |
| S1016 | 1 | 2 |
| SA4004 | 0 | 1 |
| SA4010 | 0 | 1 |
| SA5000 | 1 | 1 |
| SA9009 | 1 | 1 |
| ST1005 | 0 | 2 |
| ST1012 | 0 | 2 |
| ST1018 | 2 | 2 |
| **合计** | **54** | **99** |

linux 另有 2 行非 finding 编译噪音（`sherpa-onnx-go-linux@v1.13.8` build constraints、`cmd/balldebug` build constraints），不计入分母。

## 非-U1000 逐条（linux 形状，7 条）

| # | 位置 | 码 | 一句话 | 处置 |
|---|---|---|---|---|
| 1 | internal/plugin/manifest.go:231 | S1011 | 内层逐元素 append 循环可换成 `append(errs, ...)...` | **修** |
| 2 | internal/plugin/manifestslot/install_50r1f_test.go:48 | S1011 | 同上形状，`out = append(out, m.saved...)` | **修** |
| 3 | internal/tools/web_search.go:252 | S1016 | `webSearchArgs` 与 `SearchQuery` 字段序/型全同，`SearchQuery(a)` 转换替代结构体字面量 | **修** |
| 4 | internal/risk/markseam_394_test.go:24 | ST1018 | 字符串字面量含 U+200B，改 `\u200b` 转义（串值逐字节不变） | **修** |
| 5 | internal/risk/markseam_394_test.go:162 | ST1018 | 同上 | **修** |
| 6 | internal/plugin/panic_isolation_51_test.go:62 | SA5000 | assignment to nil map | **不修（有意为之，具名理由见下）** |
| 7 | frontend/embed.go:4 | SA9009 | 已知假阳性（散文注释恰好以 `// go:embed` 开头，真指令在 `:19`） | **不修（只登记，票面 AC#3 明令）** |

### SA5000 具名理由（不修）

`panic_isolation_51_test.go:62` 的 nil map 写入**就是被测行为本身**：该文件 `:42-43` 注释逐字写着
"panicVal == nil selects the nil-map-write form: the classic nil-dereference host bug, which arrives as a
runtime.Error"。测试要制造的是**隐式 runtime.Error 崩溃**（区别于显式 `panic(val)` 分支），staticcheck 的
"建议"（初始化 map）会直接拆掉这个崩溃，被测的宿主崩溃恢复路径就不复存在——这属于票面铁律
"不许为变绿放宽断言" 的反例形状。票据禁令同时禁 `//lint:ignore` 抑制 ⇒ 本条只能原样保留并具名登记，
是否豁免该夹具的读数归属由编排者裁定。

### windows 形状独有的 7 条（CI 门看不到，本波不动，波 2 候选）

| 位置 | 码 | 一句话 |
|---|---|---|
| cmd/wisp/panel_transport_35r2_test.go:965 | SA4004 | 外层循环被无条件终止 |
| internal/audio/aecspike/longrun_windows_test.go:773 | S1016 | 结构体字面量可换类型转换 |
| internal/ball/hotkey_borrow_refused_260r5_test.go:101 | ST1012 | 错误变量命名非 errFoo 形 |
| internal/ball/hotkey_refused_borrow_release_245_test.go:58 | ST1012 | 同上 |
| internal/secret/delete_test.go:131 | SA4010 | append 结果未被使用 |
| internal/tools/recycle_windows.go:151 | ST1005 | 错误串大写开头 |
| internal/tools/search_apps_windows.go:94 | ST1005 | 同上 |

## 修后分组计数（复跑，同一跑法）

| 检查码 | linux | windows |
|---|---|---|
| U1000 | 47 | 85 |
| SA5000 | 1 | 1 |
| SA9009 | 1 | 1 |
| S1011 / S1016 / SA4004 / SA4010 / ST1005 / ST1012 / ST1018 | **0** | 见上表（windows 独有 7 条不变） |
| **合计** | **49 = 47 U1000 + SA5000 + SA9009** | **94** |

净减 5 条（linux 54→49、windows 99→94），零新增；被修 5 条在两个 GOOS 下都消失。
U1000 逐符号判表见 `docs/evidence/s1/122-u1000-judgment.md`（本波零删）。
