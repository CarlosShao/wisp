# 65 — orb 集成腿证据（A1594 改道：集成 LerSent001/orb，不自写着色器）

- 时间：2026-11-07T12:5x+08（写腿单枚，独占 internal/ball）
- 角色：本件是**实现腿自证**；按 D22 双角色，缺口审计/对抗验收须由另一 agent 复核
- 票面：`.scratch/wisp/issues/65-ball-glass-quality-rework.md`（Status/AC 勾**未动**，只追加 Progress log 一行）

## 做了什么（三句）

1. **上游原样上舰**：`internal/ball/assets/orb-siri.html` = LerSent001/orb 的独立导出页
   （MIT，package.json license=MIT；D:\tmp\wisp\20261107\orblab\exported\orb-siri.html 的逐字节拷贝，
   唯一变换 CRLF→LF，89660 字节）。零自写着色器；钉死测试
   `TestOrbAssetShippedUnmodified` 逐 marker 验明"这就是上游页"（liquidOrb 接口/stateSeeds/WGSL @fragment），
   防后来腿悄悄换成仿制品。
2. **WebView2 悬浮宿主**（`orb_windows.go`）：WS_POPUP + WS_EX_LAYERED|TOPMOST|TOOLWINDOW|NOACTIVATE，
   固定方窗（基础边长 108px×目标 DPI），透明背景走库里现成的
   `ICoreWebView2Controller2::PutDefaultBackgroundColor(a=0)`（IID {c979903e-d4ca-4228-92eb-47ee3fa96eab}，
   jchv/go-webview2 已绑定——**未新写任何 COM 绑定**）。服务面无 socket：WebResourceRequested 过滤器
   + 合成源 `https://wisp.orb/orb.html`，每条响应带 CSP 头（connect-src 'none'），404 永远空体。
   嵌入唯一加的字节 = `orbPageBootstrap`（document-created 注入，清掉导出页为浏览器标签准备的深色 body 背景，
   露出宿主透明）。
3. **状态/音频对接**：`window.liquidOrb`（上游 freeze 对象）经 ExecuteScript 推送。
   20 个冻结态（D43/C12）→ 上游仅有的两态 `idle|thinking` 映射表钉死
   （`orbstate.go` + `TestOrbStateMappingCoversAll20States`：工作族→thinking，待机/故障族→idle）；
   语音族（Listening/Speaking/Conversation）把既有 `SetAudioLevel` 标量转 `setAudioBands`（四带等值，
   NaN→0、限幅 [0,1]），Sleeping 永不被麦克风唤醒（与 D2D liquid 同一法则）。

## 回退（A1594 ③，一字未删 D2D）

- `attachOrb` 失败（WebView2 运行时缺失/Embed 拒绝/无 GPU 会话）→ 只打一行 WARN，`b.orb` 保持 nil，
  D2D 路径原样渲染（renderFrame/present 的守卫是三个 early-return；D2D 代码在其下原封未动）。
- 球窗拖拽/DPI 变更 → `followOrb()` 把 orb 窗摆在球右侧 12px（拖动中不跟，松手/DPI 时跟；
  orb 窗自身点击穿透由 NOACTIVATE+层叠窗形状与球窗的 HTTRANSPARENT 法则一致承担）。

## 实测读数（本机，12:46-12:50）

- 默认套件 `go test ./internal/ball -count=1 -timeout 20m` → **ok 0.245s**（orb 8 个确定性测试全绿）。
- 实窗 `go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -v -count=1` → **PASS 1.23s**：
  真挂载（DebugOrbHWND≠0，ex-style 四位全中、方窗）、真推状态（debug 面 idle→thinking），
  50 态+50 音频推送句柄数**零增长**，Close 后两窗全灭（ball 窗+orb 窗 IsWindow=0）。
- 该测试还抓出并修掉一个真竞态：`destroy()` 原顺序先置 browser=nil 再 DestroyWindow，
  销毁途中同步到达的 WM_WINDOWPOSCHANGED 会进库的 NotifyParentWindowPositionChanged 空指针 panic
  （observe 边界接住）并留下活窗——修为**先摘路由再销毁**（orb_windows.go 注释记了实测出处）。
- 门禁：`gofumpt -l`（6 个新 .go）空；`go vet ./internal/ball` rc0；
  `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` rc0；
  `sh scripts/d22scan.sh clean` → **clean**（internal/=938 Go files 扫过）。

## 已知边界（如实，未做假装）

1. **视觉签收未发生**：本腿交付的是"上游页面在 Wisp 宿主里跑起来+状态推得进去"。
   票 65 的完成判据仍是机主人眼（"这次对了"）；AC 六条一个未勾。
2. **D32 预算未在 webview2 渲染器上重测**：Sleeping 零定时器/CPU/内存的差分是对 D2D 路径做的；
   WebView2 渲染器是常驻 GPU 合成面，其空闲功耗需要仪器重测（webview2 冷拉起 >2s 重评 L2 的
   既定待定案项同源）。挂载一次性句柄成本 ~300（同库给面板宿主的形状，库从不释放 environment；
   进程内 CAS 保证只挂一枚）。
3. **D2D 渲染器仍会初始化**：挂上 orb 后 D2D 帧不画（守卫挡住），但其对象仍创建——
   是无害浪费不是泄漏，后续腿可做懒建。
4. **音频只送标量**：上游分析器分低/中/高频带；Wisp 音频捕获腿现在只给一个 RMS 包络，
   四带等值推送。分带是音频腿的缝，不在本腿。
5. **时间窗**：本腿开工晚于 cmd/wisp 写腿的告警窗（派单已注明）；实际碰撞未发生
   （cmd/wisp 全程未动，`go build ./cmd/wisp` rc0 复核）。
6. 挂接碰了两个共享件（`ball_windows.go` +42 行、`liquid_windows.go` +7 行）：
   场点=New 挂载钩子/Close 拆卸/renderFrame+present 守卫/SetState 推送/SetAudioLevel 转发/
   两处 wndproc followOrb 挂点；除 SetState 由单行展开为块（原 `applyStateLocked` 调用原样保留）
   外全部是**追加**，既有逻辑零改写。

## 你自己跑起来看（两条命令）

```
go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -v -count=1
```
→ PASS 且桌面上球+orb 实窗出现过又消失（约 1.2s）。

```
go run ./cmd/balldebug
```
→ 真球+orb 停在桌面，20 态循环每态 2 秒；球右边的玻璃球会随状态在 idle/thinking 两张脸间切换。
（`scripts/dev/ball-cycle.ps1` 是同一循环的打包版。）

## 撤销口令

回滚 = `git revert <本腿 commit>`（单提交，既有文件改动全部随 revert 还原；
新增件随 revert 删除后不残留——临时件纪律只约束 .scratch 证据，不约束 revert 语义）。
