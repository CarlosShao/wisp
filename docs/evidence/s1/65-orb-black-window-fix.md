# 65 — orb 窗全黑修复（三层根因，2026-11-07）

症状（编排者亲测）：`build/balldebug-orb.exe` 的「Wisp Orb」窗内容**纯黑**，球不渲染；
stdout 反复出现 `ball: orb response create retry ... err="An attempt was made to reference a token that does not exist."`。

## 一、诊断读数（先诊断后动手，全部本机实测）

**① 失败点诊断（直读 COM vtbl 的探针 + 包装自带读数）**

```
level=WARN msg="ball: orb response create DIAG" attempt=1 errno=1008
  errnoText="An attempt was made to reference a token that does not exist."
  directHRESULT=0x00000000   directResponseNonNil=true
  threadID=55468  envPtr=0x437800004410  controller2Ready=true
```

结论：**COM 调用其实成功了**（HRESULT=S_OK、response 对象就在 out 参数里）。jchv 包装
`ICoreWebView2Environment.CreateWebResourceResponse`（pkg/edge/corewebview2.go:169-203）把
HRESULT 丢掉、拿裸线程 `GetLastError` 当成败判据；orb STA 线程上 1008（ERROR_NO_TOKEN）**持续存在**
（不是面板 2026-10-13 量的"冷启动首报"），包装据此把成功当失败返回 nil（对象还泄漏），`PutResponse`
永远拿不到东西 → 页面 HTML 发不出去 → 球窗全黑。"重试"对这个形状无效：3 次全灭已实测。

**② 页面侧真相（页面经 chrome.webview.postMessage 报话）**

```
msg="diag url=about:blank liquidOrb=object gpu=undefined ready=complete"
```

换 `NavigateToString` 后页面**确实加载了**（上游接口活着），但 `navigator.gpu=undefined`——
页面自报 **"WebGPU is not supported in this environment."**。本机 WebView2 运行时 157.0.4341.42
（WebGPU 完全支持），缺的是 secure context：about:blank 的不透明 origin 拿不到 WebGPU。

**③ 合成层真相（控制器可见性）**

```
Chrome_WidgetWin_0 / Chrome_WidgetWin_1 / Intermediate D3D Window 全部 visible=False
```

jchv `Embed` **从不调 `PutIsVisible(true)`**；orb 的 never-activate show 路径下控制器带着隐藏的
合成子窗上车，窗口什么都不合成（剥 WS_EX_LAYERED、SetLayeredWindowAttributes、ShowWindow+
UpdateWindow 均无效——对照实验逐项排除）。显式 `browser.Show()` 后页面像素立刻上屏。
（面板宿主创建即激活+聚焦，控制器可见性在那里自然成立，故面板无此修复也活着。）

## 二、修法与选型理由

- **选 `SetVirtualHostNameToFolderMapping` 而非 `NavigateToString` 作主路径**：页面的球是 WebGPU
  着色器，运行时只对 secure context 暴露 `navigator.gpu`（② 的实测）。virtual-host 映射给页面真正的
  `https://wisp.orb` origin，WebGPU 起来；`NavigateToString` 只能给出 about:blank（降级路径保留：页面
  界面/状态/音频仍活，页面自报 WebGPU 不支持）。这是"页面需要真 origin"场景的主流做法，不是动作最小。
- **CSP 保住（D29 形状不变）**：两种投递形态都没有响应头通道，CSP 以注入的
  `<meta http-equiv="Content-Security-Policy">` 进 served 副本（orbPageDocument，orbasset.go）——
  `default-src 'none'` + `connect-src 'none'` 全套原样；文档 API 只能收紧不能放宽，页面剥不掉。
  **asset 文件字节零动**：`internal/ball/assets/orb-siri.html` sha256 复测 =
  `d93da649d11bc4bc3ab91c006777bbadce1c2eb7d0ffc9d91608ad8d866e3bc6`（与钉死值逐字节一致）。
- **离线形状不变**：映射只存在于本 WebView2 实例内，host 在无解析的 `.orb` TLD 下，meta CSP 切断一切
  网络动词；无 socket、无端口、无 WebResourceRequested 响应构造缝（含 GetLastError 判错的那类缺陷
  结构性消失；模块 `SetVirtualHostNameToFolderMapping` 包装同款判错缺陷已按"告警+降级 NavigateToString"
  处理，不盲信）。
- **可见性**：newOrbWindow 在 Embed/Resize 后显式 `browser.Show()`（= controller.PutIsVisible(true)）。

## 三、验收（修前黑 / 修后不黑，有牙）

修后独立跑实拍（orb 自然位置，未做任何移动）：液态玻璃球真渲染——彩虹折射带、玻璃质感
（`build/orb-final-natural.png`）；最终跑日志 0 条 WARN/ERROR（修前反复 retry 告警）。

**可复跑命令（编排者自跑）：**

```
go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 15m -v
```

回归用例 `TestBallLiveOrbPageRendersPixels`（orb_live_windows_test.go）：真球+真 WebView2，
把 orb 移进工作区中央，在其 rect 下铺一枚 **topmost 黑底窗**（挡掉壁纸/普通窗的污染），再用屏幕 DC
按 3px 网格采样合成后的真实像素，断言「RGB 任一通道>30 的采样像素占比 >1%」。修前形状（页面发不出/
WebGPU 不起/不合成）在此用例全为全黑采样 → 必红。

本轮读数：**634/1764 = 35.94% 非黑（> 1% 阈值）**，全套 `TestBallLiveOrb` 4/4 PASS。
（诚实注脚：若桌面恰有**别的置顶悬浮件**压进采样区，读数可能被抬——黑底窗已挡掉普通窗这一层，
置顶悬浮件无窗内手段排除；阈值 1% + 常态 30%+ 量级使误报空间很小。）

## 四、门禁读数表

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1 -timeout 20m` | ok |
| `go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 15m` | 4/4 PASS（像素 35.94%） |
| `sh scripts/d22scan.sh clean` | clean（bans #1-8/#10 全量 examined，0 violation） |

## 五、改动面

- `internal/ball/orb_windows.go` — serving 缝重写（virtual-host 映射 + NavigateToString 降级）、
  `browser.Show()` 可见性修复；删除 serve/orbBuildResponse/诊断脚手架。
- `internal/ball/orbasset.go` — origin/entry 常量恢复、orbPageDocument（meta CSP 注入）、orbWriteSite。
- `internal/ball/orb_test.go` — serving 常量钉测更新 + 上游 sha256 钉 + orbWriteSite 钉测。
- `internal/ball/orb_live_windows_test.go` — 新增 TestBallLiveOrbPageRendersPixels（像素回归）。
- 未碰：cmd/wisp/**、internal/panel/**、internal/winsec/**、internal/agent/**、frontend/**、design/**、
  thresholds/golden、票面 AC/Status。

## 六、残留风险（如实）

1. 面板宿主无 `Show()` 调用，但其窗口激活/聚焦路径当前成立——**本腿未动面板**；若面板日后改 never-activate
   路径，会踩同一坑（此处留档）。
2. jchv 模块三处包装（CreateWebResourceResponse / SetVirtualHostNameToFolderMapping / 各 Put*）共用
   裸 GetLastError 判错，同款误报可在任何 STA 线程复现；本腿只在 orb 面内防御，模块级修复属上游票。
3. 站点文件夹落盘 `%LOCALAPPDATA%\wisp\orb\site\orb.html`（0600，每次 attach 幂等重写）——新增一块
   每用户磁盘面，内容= served 文档（含 meta CSP），无可执行面。
