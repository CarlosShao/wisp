# 票 329 AC#5 — 页面侧接收器（`window.wispOnPanelSnapshot`）交付表

- 腿＝页面接收器补票（`frontend/**` 本波独占写面）；日期＝2026-11-08 11:2x +08。
- 合同＝`.scratch/wisp/probes/329/r1/page-contract.md`（65 行，逐字读过；形状未自创一字）。
- Go 半边＝wave81 落地件（票面 Progress log 2026-10-17 bullet；本腿未碰 `cmd/wisp/**` 一字）。

## 1. 交付件（逐枚）

| 件 | 变化 | 一句 |
|---|---|---|
| `frontend/src/lib/panel.ts` | 改 | `Window` 声明新增 `wispOnPanelSnapshot?: PanelSnapshotReceiver`＋导出 `PanelSnapshotReceiver`／`PanelSnapshotPushTarget`／`installPanelSnapshotReceiver(target, cb): () => void`（卸载只删自己的回调，StrictMode 次序安全） |
| `frontend/src/panel-root.tsx` | 新 | `PanelRoot`＝挂载点：`useState<PanelSnapshot | null>(null)`＋`useEffect` 把 `window` 装上接收器；`<App snapshot={pushed ?? undefined} />` |
| `frontend/src/main.tsx` | 改 | 产品根由 `<App />` 换成 `<PanelRoot />`（harness=1/2 两路一字未动） |
| `frontend/scripts/render-snapshot-receiver.tsx` | 新 | 渲染仪器（照 `render-config-screen.tsx` 形状），15 枚 check |
| `frontend/package.json` | 改 | 加 `render:receiver` 脚本一枚 |

**为什么接收器状态放 `panel-root.tsx` 不放 `App`**：`scripts/render-nav.tsx:180` 把 `App.tsx` 钉在**恰好一个 `= useState<`**（视图选择＋伴随偏好）；"宿主推来的最新一包"是文档生命周期级的馈入，不是视图偏好 ⇒ 归挂载点。这也是主流 harness 的形状（root 持宿主馈入，下层是纯渲染）。`App.tsx` 本腿**一字未动**（`git diff` 为空），render-nav 的四条 App 源码钉全部原样。

**合同语义落点**：每次调用＝`setSnapshot(s)` 整包替换（无 diff／无 merge／无节流）；冷启动＝页面只在 React 挂载后定义全局，首包前的推被 WebView2 丢弃、页面**不发明补发**（合同 §4 明文归票 33/35），首包未到时维持 `App` 自己的 EMPTY 帧并如实写明。

## 2. 前端测试（`npm run render:receiver`，rc=0）

15 枚 check 全绿，逐类：

- **接收器门（真实调用，非只渲染）**：安装后具名全局被定义；首推**整对象落到**（恒等，非字符串／非拷贝）；第二推**整包替换**（第二包刻意不带 `tasks` 键 ⇒ 断言合并不存在）；连推 5 次幂等（最后一包赢，不崩）；卸载删全局；陈旧卸载器不误杀重装回调（StrictMode 次序）。
- **未推时**：`<App />` 空帧照常渲染（不崩、无球风格 picker、无待审批徽标、零 emoji）。
- **推到 App**：推送包的正文进对话区、工作区进侧栏；快照带 `view:"config"` 时外观页的 `ball.style` 行显示**推送包里的 family**（`>jarvis<`）而非无读口句。
- **源码级接线**：`panel-root.tsx` 经 `installPanelSnapshotReceiver(window, setSnapshot)` 唯一门安装、路径上无 `Object.assign`／无 spread 合并；`main.tsx` 挂 `<PanelRoot />`；全局名在 `lib/panel.ts` 的 `Window` 声明里（⛔ any）。

## 3. 门禁四数（全部本腿自跑）

| 尺 | 读数 |
|---|---|
| `npm run typecheck`（`tsc -b`） | rc=0 |
| `npm run build` | rc=0（dist 产出 `index-D_cLgpmD.js`，随后经 `go build` 嵌入） |
| `npm run render:receiver` | rc=0（15/15） |
| `npm run render:nav`（App 源码钉回归） | rc=0 |
| `npm run lint`（oxlint） | rc=0（0 errors；35 warnings 全为既有） |
| `sh scripts/d22scan.sh` | rc=0 clean（ban #6/#8 frontend/ 119 text files 全扫） |
| LF | 新件（panel-root.tsx／render-snapshot-receiver.tsx）CRLF=0；main.tsx／package.json 工作树 CRLF 为**检出时既有**（`git ls-files --eol`＝`i/lf`＋`text=auto`，提交 blob 为 LF，本腿未引入新 CRLF） |

## 4. 端到端实机（如实，两半分开报）

**Go 半边（本腿现跑复证，非转抄）**：`wisp.exe`（`go build -o wisp.exe ./cmd/wisp` rc=0，`panel-assets -manifest` 确认嵌入新 bundle `assets/index-D_cLgpmD.js`）；托盘"打开面板"（`EvTrayPanel`）唤出面板，每次 show 一发：

```
panel: SNAPSHOT at=2026-11-08T03:12:27Z depth=0 pending= mode=ask_every_step ws=unset results=0/0 done bytes=3656 sha256=…
wisp: panel window is up (tray-open-panel, cold …)
```

与编排者 wave81 读数同形（bytes=3656）。

**页面那半在真 Chromium 里活体取证（已证）**：把**同一个嵌入 bundle** 经 `vite preview` 供给真浏览器（Tabbit/Chromium），页面挂载后 `typeof window.wispOnPanelSnapshot === "function"`；按合同调用形状执行 `window.wispOnPanelSnapshot(<整包对象字面量，ballStyle.style="jarvis", view:"config">)` ⇒ 外观页 `ball.style` 行右侧显示 **jarvis**（能改徽标在），无读口句不出现；console 0 error。截图＝`D:\tmp\wisp\20261108\orblab\PANEL-BALLSTYLE-WIRED.png`（球风格行）＋`PANEL-PUSH-REACHES-SETTINGS.png`（推送包整体落屏）。**读数：该行右侧文字＝`jarvis`，不是「缺省／快照里还没有球风格这一节」。**

**wisp.exe 的 WebView2 窗口本机今天画不出（缺口，如实报）**：面板窗（`wisp-panel-host`）打开后内容区**纯白**（至少 5 次独立复现；新 profile 也一样），且 **wisp.exe 进程在面板打开后数十秒～两分钟内无退出日志地消失**（stderr 止于正常 wndproc 行，无 panic／无 WER 记录）。**A/B 对照＝`D:\tmp\wisp\20261106\170ac1\wisp.exe`（2026-11-06 旧二进制，先于 wave81 与本腿）同机同日同样白窗＋同法死亡** ⇒ 缺口在本机环境／面板宿主侧，**与本接收器无关**（本腿只改页面 JS；JS 不能杀死 Go 进程）。另一佐证：同一 bundle 在真浏览器全绿。该格（wisp 窗内 ball.style 行的像素级截图）**本机取不到**，请编排者裁（换机／或先裁面板宿主白窗＋进程死亡两缺陷）。

附带环境读数（可能相关，供宿主腿排查）：本机常年跑着 NBA2K Online（反作弊＋全局钩子＋独占热键——首启四枚热键全报 `occupied by another program`，`[hotkey]` 改成空闲组合后 bridge 1s 内重绑成功 `panel=Ctrl+Shift+F9 live=3`，验后 config 已还原 `panel = ''`）；magpie（DWM 放大）在跑；对 `wisp-panel-host` 做 UIA 无障碍查询（WM_GETOBJECT 洪泛）与进程死亡强相关（两发死亡前 stderr 尾均是它）——宿主腿可循这条查。

## 5. 本腿没做／没碰（如实）

- ⛔ 未碰 `cmd/wisp/**`、`internal/ball/**`、`internal/panel/**`（只读）、`internal/config/**`、`design/**`。
- ⛔ 未重造写路（config.set 写半边＝票 65 界面腿既有，`a3fdb50dd`）；未发明冷包补发／节流。
- `PanelSnapshot` 类型**不缺键**：`ballStyle`／`plugins`／`run` 等九节全部在 `lib/panel.ts:398-434` 声明（与 Go json tag 双向对账由 Go 侧测试钉着），本腿零补声明。
- 临时件（探针脚本、截图、日志）在 `D:\tmp\wisp\20261108\`，只建不删。
