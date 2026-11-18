# 36 AC#5 — 面板隐藏即暂停（rAF 计数器实证）＋顺手 deflake 39 — 写腿收件（agent-140, 2026-11-18 r1）

- 起手锚 `ca30abb85`（＝派单存档笔 `36-ac5-deflake39-dispatch.txt`）；未 push。
- 写面＝`frontend/src/**`＋`frontend/scripts/render-hiddenpause36.tsx`＋`frontend/package.json`（仪器
  script 名注册一行）＋`cmd/wisp/config_loosen_gate_39_test.go`（唯一）＋本两件。`internal/panel/**` 本腿
  零触碰（不改即不涉共享编译面 Go 半）。
- **既有量普查（dispatch §A1）**：`git grep -n "visibilitychange|visibilityState|requestAnimationFrame|
  cancelAnimationFrame" -- frontend/src` 命中只有 `react-bits/line-sidebar.tsx`（指针驱动的近邻 rAF，
  自停型，非面板隐藏语义）。**面板无既有隐藏暂停机制 ⇒ 本腿新建**（dispatch「没有就做」）。

## 1. 机制（`frontend/src/lib/panel-visibility.ts`，产码）

主流形状，逐字对齐 dispatch：「`document.visibilityState` 变 ⇒ 停 rAF 循环、恢复再启动」：

| 件 | 读法 |
|---|---|
| `panelWorkAllowed(doc)` | 唯一事实＝`document.visibilityState === "visible"`；⛔ 任何时间戳推导（本文件零 `Date.now`/`performance.now`）；无 DOM（SSR/静态渲染）＝不允许（停摆是诚实答案） |
| `createPanelAnimationLoop(doc, onFrame)` | rAF 循环工厂：start 挂一帧、tick 尾部自续；stop 取消在飞帧、链就地死。仪器计数器的载具 |
| `installPanelVisibilityPause(doc, loops)` | **每页一次订阅**（panel-root.tsx 挂载点 `useEffect` 挂/StrictMode 卸）：hidden ⇒ 在跑循环逐个 stop 并记入 suspended 集；visible ⇒ **恰 suspended 集重启**——隐藏前已停的循环不被 resume 拉起（暂停不发明工作），隐藏前在跑的循环绝不停在停态（恢复半，MUT-B 的靶） |
| `panelAnimationLoops` | 页级循环名册（本页 JS 循环从这里走，循环不自带 listener——一页一订阅） |
| `usePanelVisibilityTick()` | RevealText 梯子用的恢复信号：visibilitychange 计数；effect 依赖数组带 tick，事实翻转的当刻重新过闸 |

接线（产码消费面，非 shelf-ware）：
- `panel-root.tsx`：挂载点订阅一次（ticket 329 的 snapshot receiver 旁边，同一 mount/unmount 纪律）。
- `reveal-text.tsx`：**面板里真实的 JS 动画源＝RevealText 的 46ms 逐段梯子**（流式中每 46ms 醒一次 timer）。
  现在调度闸在 `panelWorkAllowed()` 上：hidden ⇒ effect 早退，**一个 timer 都不排**（醒源消失，CPU 规则落地）；
  visible ⇒ tick 重开闸，梯子从冻结的 count 续走（不补发、不跳变）。46ms 节奏与「无 reset 当 text 增长」两个
  冻结 ruler 零字节改动。
- `App.tsx`：`PanelHiddenPauseProbe`（可注 document 的只报事实 span，`data-panel-hidden-pause`）——App 自己
  没有 JS 循环（CSS keyframes 是合成器的，浏览器自己停）；它只提供仪器在**整 App 层**钉接线的探针，不发明
  App 没有的动画。

诚实边界（写死在件头）：CSS 半（fade-up/caret/球环）在隐藏页由合成器自己停——那半不需要代码，本件不假装
加了它；本件拥有的是合成器看不见的 JS 半。

## 2. 仪器（`frontend/scripts/render-hiddenpause36.tsx`，npm 名 `render:hiddenpause36`）

解析器级浏览器类引擎（本仓仪器惯例，仓内无 jsdom）：`FakeDocument`（EventTarget＋visibilityState，
`setVisibility`＝改事实＋发 `visibilitychange`——浏览器怎么发就怎么发）＋`FakeView`（rAF 真句柄登记表：
schedule→run→schedule 至 cancel；**cancel 真把在飞帧摘掉**；一 pump＝一帧的活）。全部判读＝**计数器差值**，
本文件零墙钟。

绿控读数（正控两次：真仓一次、导出树一次，读数同）：
- visible：pump 后 1→2（增长）；
- hidden：3 连 pump 采样 **2/2/2 相等**（dispatch「同一 tick 内多次采样相等」的原形）＋在飞帧被摘下（scheduler 空）；
- 切回：2→4（**恢复后计数必须再增**）；
- 恢复成员规则：恰 suspended 集重启；从未启动的循环不被 resume 发明；
- uninstall 后订阅不再驱动循环；
- 探针（App 层）hidden⇒`paused` / visible⇒`running`；
- 源级钉（静态渲染跑不了 effect，钉接线）：reveal-text 的闸行＋`[count, done, visTick]` 依赖、panel-root 的
  订阅行＋cleanup、**零自持 listener**（一页一订阅）、两枚变异锚在盘。

## 3. 牙（dispatch §A3，变异全在仓外导出树，真仓零字节）

台面＝`git archive ca30abb85` 导出 `D:\tmp\wisp36mut2`＋工作树改件覆入＋`node_modules` junction（只建不删）。
tar 解包对 `.scratch` 16 枚异形字符历史件解包失败＝与腿 2 的记录同形（前端树与 1906 枚 .go 全在）。

| 发 | 变异（导出树逐字） | rc | 红名/红句 |
|---|---|---|---|
| 正控 | 未变异导出树 | **0** | `OK, rAF counters judge the pause by deltas (visible 1->2 grows, hidden frozen 2/2/2 equal, resume 2->4 grows, ...)` |
| MUT-A | 删 hidden 支的 stop 循环（MUTATION-36A 锚 5 行＋stop 循环 4 行整删） | **1** | `FAIL hidden: the counter is frozen (pumps 1/2/3 equal: 3/4/5) [2/3/4/5]`＋`the pause stopped the loop`／`pending frame off scheduler` 红 |
| 还原 | 工作树原件覆回 | **0** | 复绿（同读数） |
| MUT-B | 删 visible 支的 restart 行（`for (const loop of suspended) loop.start();`） | **1** | `FAIL resume: the counter grows again (2 -> 2)`＋`the growth continues (2 -> 2)` 红——只停不恢复即红 |
| 还原 | 覆回 | **0** | 复绿收台 |

⚠ 诚实注记：牙的靶＝`installPanelVisibilityPause` 的两条循环（机制体本身）。MUT-A/B 各自的红句是计数器差值
检查（§2 的 hidden/resume 判读），变异锚检查只是伴红。**RevealText 的闸行没有动态牙**（SSR 跑不了 effect，
本树惯例是源级钉＋共享机制的动态证明）——这句写在件头，不假报。

## 4. 门禁读数（全部现跑，2026-11-18 21:5x–22:3x +08）

| 门 | rc | 备注 |
|---|---|---|
| `npm run typecheck` / `build` / `lint` | 0 / 0 / 0 | lint 0 errors / 50 warnings（存量）；本腿文件零新增命中（App.tsx 两条 = HEAD 既有的 monitorAgents 导出与 set-state-in-effect，现核于 `git show HEAD`） |
| 新仪器 `render:hiddenpause36` | 0 | 读数 §2 |
| 既有仪器 `render:history36 / stream20036 / stream / receiver / historypage36 / shell77` | 全 0 | stream20036/stream 首跑红＝`usePanelVisibilityTick` 默认参在无 DOM SSR 下 `document is not defined`，已修（默认参守卫，同 panelWorkAllowed 形）；复跑全绿，markup 逐字不变 |
| `dist` | 未入库 | `git ls-files frontend/dist`＝1（.gitkeep）；build 产物零提交（票 248 AC#9 的既有边界） |
| `sh scripts/d22scan.sh` | **0** | `OK - control rc=0, scan rc=0`；产码注释改动后复跑仍 0 |
| `go build ./...` | 0 | |
| 前台 M 件 | 未动 | `frontend/fixtures/composer-states.html` 为他人在飞件（M）：只报不动，⛔ restore/checkout |
| 载具整包（deflake 39 半） | 见 `deflake39-20261118.md` | rc=1＝唯一红 33r5（名册环境红）；本例整包 PASS；读数与作废的第 1 发记在同件 |

## 5. 诚实边界（A 半）

- `PanelHiddenPauseProbe` 的接探针（App 层）是**静态渲染**读数（探针不挂进页面树，只由仪器渲染）；
  RevealText 的闸行是**源级钉**（SSR 跑不了 effect）——rAF 循环半（机制本体）才是动态计数器证明；
  这两行写在仪器件头，不假报。
- CSS 半（合成器自停）不加代码、不冒领，件头明写。
