# 05 · `got -1.000` 这个哨兵值：四条出口的现量行号＋可达性

锚＝本腿 HEAD `052b393f`；文件＝`cmd/wisp/panel_host_windows.go`（969 行，`wc -l` rc=0）。

## 现量：`grep -n 'return -1' cmd/wisp/panel_host_windows.go` ⇒ **恰 4 行**（rc=0）

| 出口 | `return -1` 行号 | 触发条件（原文） | 含义 |
|---|---|---|---|
| E1 | **845** | `:842-844 w := m.w` … `if w == nil {` | 拿不到窗口对象 |
| E2 | **860** | `:852 if err := w.Bind("wispProbeRT", func() string { ... }); err != nil {` | 绑定探针失败 |
| E3 | **878** | `:876-877 if ctx != nil { if err := ctx.Err(); err != nil {` | 上下文被取消 |
| E4 | **882** | `:881 if time.Now().After(deadline) {`，`deadline := t0.Add(5 * time.Second)`（`:868`） | 5 秒单调截止到期 |

`firstRoundTripLocked` 定义在 `:840`；`t0` 来自 `bringUp` 的第 **319** 行 `t0 := time.Now()`（⇒ 5 秒是从**开窗之前**起算，不是从探针起算）。
返回值唯一的去处：`coldStartPageHandover`（`:447`，`:448 rtMs := m.firstRoundTripLocked(ctx, t0)`）→ `bringUp:421-425` 把它写进 `m.lastColdMs`（**只有这一处赋值**，`grep -n 'lastColdMs'` rc=0：声明 `:174`、访问器 `:285`、赋值 `:424`）。
测试侧读它：`cmd/wisp/panel_host_windows_test.go:659 coldMs := mgr.LastColdMs()` → `:660 t.Logf(...)` → `:661 if coldMs <= 0 {` → `:662 t.Fatalf("cold bring-up did not produce a browser round trip (got %.3f); ...")`。

## 答复：CI 那发红句**能不能区分是哪一条出口**？

**红句本身不能**——它只带一个 `float64`（`-1.000`），函数签名里**没有 error、没有原因位、没有计数**；`LastColdMs()`（`:285`）也只回 `float64`。缺的那一栏是**"哪条出口"**（E1/E2/E3/E4）。

**但把红句 + 同发别的读数 + 码摆在一起，四条里能排掉三条**（⛔ 这是推理，不是仪器读数，逐条给依据）：
- **E1 排除**：`m.w` 在 `:413-417` 就已赋值（`m.w = w; m.hwnd = hwnd; m.created = true`），且同一发里 `:647 if !mgr.IsCreated()` 与 `:650 if hwnd := mgr.windowHandle(); hwnd == 0` 两道 `t.Fatalf` **都没响**（红句出现在 `:662` 就是证据）⇒ 进入 handover 时 `w != nil`。
- **E3 排除**：测试的 ctx 是 `panel_host_windows_test.go:625 ctx, cancel := context.WithCancel(context.Background())`，`cancel` 只在 `:626` 的 `defer` 里；`hostThreadHarness.bringUp`（`:110-112`）逐字 `return hh.runOnThread(func() error { return mgr.bringUp(ctx) })` —— **没有 `WithTimeout`/`WithDeadline` 包装** ⇒ bringUp 期间 `ctx.Err()` 恒 nil。
- **E2 静态不可达**：本仓用的库是 `webview2 "github.com/jchv/go-webview2"`（`panel_host_windows.go:64`），其 `Bind` 只在两种形状下返回 error（`webview.go:450-457`，模块缓存 `D:\work\base\gopath\pkg\mod\github.com\jchv\go-webview2@v0.0.0-20260205173254-56598839c808\webview.go`，只读）：`v.Kind() != reflect.Func` 或 `NumOut() > 2`。这里的实参是 `func() string`（一个字面函数、1 个返回值）⇒ 那两条永不成立，`Bind` 必回 `nil`。
- ⇒ **只剩 E4**：5 秒单调截止到期。旁证（不是判据）：run 305 该用例 18.94 s，本机 10-08 那发 5.06 s（≈ 恰好撞 5 s 上限）。

⚠ 反过来说：**E4 成立不等于"环境没有运行时"**。同一条 E4 在 8 发里只在 run 305 出现过；另外 7 发走的是"量到了正数、超预算"那一支（`02-*.md`），那 7 发**同时**打出 `our tree webview=7（machine-wide 7）` 的正读数（`07-*.md`）。

## 口径（本项目已定，本腿只是复述并自守）

- `-1` 是**未测量的哨兵**：⛔ 不许读成"时延是 -1 毫秒"，⛔ 不许当"环境没有运行时"的证据。`02-*.md` 里那 7 发的正数读数（3 126 – 3 941 ms）与 `-1` 是**两种不同的形状**，本件逐发分开列。
- `probes/35/v6/live-wave.md:44` 那个 `-1` 是本机已知同类坏读数，⛔ 本件没引它，任何一格都没引。

rc=0
