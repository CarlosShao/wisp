# 证据：fork 级 WebView2 nil-env 守卫（票 408 残余项，编排者裁定）

- 时刻：2026-11-08 晚；腿＝408r3 nil-env 守卫腿；只 commit 不 push。
- 起点 commit：eec09a028（票 408 r2，per-boot UDF 叶＋测试层有界启动等待）。

## 1. 起手三核（现跑，非转述）

1. **回调真位置**：`grep -rn "EnvironmentCompleted" forks/go-webview2/` ⇒ 实现体在
   `forks/go-webview2/pkg/edge/chromium.go` 的 `(*Chromium).EnvironmentCompleted`（COM 完成回调，
   由 `corewebview2.go` 的 vtbl `_...HandlerInvoke` 调入）。
2. **ball 侧现况**：`newOrbWindow`（`internal/ball/orb_windows.go`）判成功＝`browser.Embed(hwnd)` 返 bool；
   `Embed==false ⇒ 具名错误` ⇒ 调用方 `attachOrb` **已有降级路**：`slog.Warn("ball: webview orb not
   attached; D2D fallback stays active")`＋`b.orb` 留 nil，D2D fallback 原样渲染（注释 A1594 ③，
   `ball_windows.go`："the D2D fallback IS the fallback, not a third mode"）。
3. **fork 归属**：`forks/go-webview2` 是本仓 **tracked 的自家 fork**（票 65 步3甲：jchv fork 进
   `forks/`＋主 go.mod `replace`＋已补过 `ICoreWebView2Controller.Close/Release` 包装），**无
   「⛔ 改 forks」规矩**；故按派单改 fork 本体，未停手。

## 2. 根因与守卫形状

- 实测（`eec09a028` 活体栈）：WebView2「出生即死」时 `EnvironmentCompleted` 以 **res=S_OK＋env=0x0**
  同步回调；原码 `if int64(res) < 0 { log.Fatalf }` 之后直接 `env.vtbl.AddRef` ⇒ **nil 解引用崩**，
  启动线程死 ⇒ `ball.New` 无限等（r2 已用测试件有界化兜住，产码仍崩）。
- 主流 harness 形状：COM 完成回调**绝不** `log.Fatal`／解引用空接口指针；失败**上报给等待方**，由
  宿主决定降级。WebView2 起不来 ⇒ 退 D2D 备胎是产品行为，不是测试作弊。

## 3. 改动（零公开 API 签名变化）

- `forks/go-webview2/pkg/edge/chromium.go`：
  1. 新增 unexported 字段 `bootFailed uintptr`（与 `inited` 同款 atomic 用法）。
  2. `EnvironmentCompleted`：`int64(res) < 0 || env == nil` ⇒ **不进成功链**，记具名日志
     `webview2: environment boot failed (hr=%08x, env nil: %v)`＋置 `bootFailed`，返 0。
     （原 `log.Fatalf` 移除——杀进程改成上报。）
  3. `CreateCoreWebView2ControllerCompleted`：同缝同款守卫（`res<0 || controller==nil`）——同一
     boot 失败类的另一半，不补则 controller 失败仍会 `log.Fatalf` 杀进程。
  4. `Embed` 泵循环：见 `bootFailed` 即断并**返 false**；泵耗尽而 `inited==0` 也返 false
     （原码此况会带着 nil webview 调 `e.Init` 崩）。
- `internal/ball/orb_windows.go`（仅接住那半）：`Embed==false` 的既有具名错误文案补准为
  `(runtime missing, blocked, or boot guard: nil environment)`；降级路**零改动**——错误照旧
  流入 `attachOrb` 的 warn＋D2D fallback。
- 新测试件 `forks/go-webview2/pkg/edge/chromium_bootguard_windows_test.go`：直接注入回调形状
  （S_OK+nil env／负 HRESULT+nil env／controller 两款＋标志复位），断言**不崩＋`bootFailed` 置位
  ＋成功链零启动**；不依赖真 WebView2 运行时、不靠偶发崩溃复现。

## 4. 门禁读数（主树；tracked 面本腿改动前 `git status` 干净，无飞件干扰）

- gofumpt（`/d/work/base/gopath/bin/gofumpt.exe -l` 三改动件）：**空**（chromium.go 初检 1 处格式，
  `-w` 后复检空）。
- `go vet ./internal/ball ./cmd/wisp`：rc0；fork 模块 `go vet ./pkg/edge`：rc0。
- `go build ./internal/ball ./cmd/wisp`：rc0。
- fork 守卫牙：`go test ./pkg/edge -run "TestEnvironmentCompleted|TestControllerCompleted|TestBootGuard" -count=1`：
  **5/5 PASS**（具名日志逐条可见）。
- `go test ./internal/ball -count=1`：`ok ... 1.739s`。
- `scripts/d22scan.sh`：**clean**（rc0）。
- LF 核：全部改动件 CRLF-lines=0。
- ★winlive `go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 30m`：
  **终跑 rc=0，10/10 全 PASS，267.260s**（20:05-20:10；读数矩阵与抖动判定见 §5）。

## 6. 暗页抖动族的读数矩阵（同代码翻绿翻红＝环境，守卫清白的决定性证据）

| 时刻(2026-11-08) | 树 | 两枚暗页牙结果 |
|---|---|---|
| 19:37 run1 | 守卫树 | PageRendersPixels 红（ServedPage 也红＝r2 债） |
| 19:44 run2 | 守卫树 | **两枚绿**（唯一红＝ServedPage r2 债） |
| 19:51 run3 | 守卫树 | 两枚红 |
| 19:58 隔离 | 守卫树 | 两枚红 |
| 19:59 对照 | 基线导出树(无守卫) | 两枚绿 |
| 20:02 复采 | 守卫树（与 19:51/19:58 零代码差） | **两枚绿（57.4s）** |
| 20:05 终跑 | 守卫树 | **全量 10/10 PASS（267.3s）** |

- 守卫树同代码先红(19:51/19:58)后绿(20:02/20:05)⇒ 与代码无关；基线只采到一个绿窗不构成
  "基线更稳"。所有红轮 attach 全部成功（无 "did not attach"），守卫零触发——红的是"附着后
  WebGPU 页暗"，属机器 GPU/环境态抖动，**另立票跟踪**（本腿不碰渲染基建）。

## 5. winlive 首跑红：基线对照实锤＝r2 欠账，非本守卫（已就地补一行修）

- 首跑（守护已落地）单枚红：`TestBallLiveOrbJarvisServedPageAndOverlay`
  `orb_jarvis_live_windows_test.go:380 read served jarvis page: open
  ...\wisp\orb\<pid>-<hex>\site\jarvis.html: The system cannot find the path specified`。
- 根因（读码定案）：该牙用 `filepath.Join(orbDataPath(), "site", ...)` 取"被服务的页面"——**r1 形状**；
  `eec09a028`（r2）已把 site 挪进 **`boot-<n>` 叶**（`orbBootDataDir`）。r2 只跑了 default-tag 的
  `go test ./internal/ball`（winlive 牙全在 `-tags winlive` 后面），故这枚红是 **r2 落地即欠、
  首次被本腿的 winlive 跑数暴露**。
- 基线对照（`git -c core.autocrlf=false clone --shared` 导出树
  `/d/tmp/wisp/20261108/408r3-baseline`，HEAD `34af3374e`＝本腿基座、不含本腿任何改动）：
  同尺 `-run 'TestBallLiveOrbJarvisServedPageAndOverlay'` **同红同文案**（读同形 r1 老路径）。
  守卫清白。
- 修补（一行＋一枚取值器）：`orb_windows.go` 加 `DebugOrbUserDataDir()`（读附着 orb 实际的
  per-boot 叶，D2D 回退时空串）；测试牙改用它取页面路径。**注：`orb_jarvis_live_windows_test.go`
  一行越出本腿写面**（派单写面＝`orb_windows.go` 仅守卫半），具名上报编排者——internal/ball/**
  本波空闲无冲突；撤销口令＝把该行改回 `orbDataPath()` 即复原红。
- 另记：winlive 家族本身带抖动——首跑 `TestBallLiveOrbPageRendersPixels` 报"三试全黑"红，
  复跑该枚绿（环境/冷启动抖动类，非本改动路径；守卫若触发只会表现为"did not attach"，
  不会表现为"attach 后全黑"）。
