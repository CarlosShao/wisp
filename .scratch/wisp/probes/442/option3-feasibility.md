# 票 442 的「候选③ 运行时判据」可行性（读码所得，未实现）

> 触发：票 442 尾部 `## 编排者改向裁定（2026-11-14 15:2x）` 把主判据定为候选③，
> 并要求「源码形状尺降为次要启发式、**本票不再往它里面加规则**」。本件只回答一件事：
> **③ 在本票写面内能不能做成有牙的判据**。⛔ 不含实现、⛔ 改任何文件（本件是 scratch 读数）。

## 1 运行时能被看见的只有一处

| 事实 | 读数（行号＝今天盘上） |
|---|---|
| `NewPanelManager` 只做赋值、⛔ 建窗 | `cmd/wisp/panel_host_windows.go:331-332`（`m := &PanelManager{… dataPath: dataPath}`） |
| 收下的值存在 `PanelManager.dataPath`（非导出字段，测试同包可读） | `panel_host_windows.go:197` |
| 真正把值交给 WebView2 的一跳 | `panel_host_windows.go:625` → `panel_control_windows.go:355`（`browser.DataPath = dataPath`） |
| 产码自己的 dataPath 是**固定**的（产品行为：一个应用一个 profile） | `panel_resident_windows.go:362`（`filepath.Join(dataDir, "panel-webview2")`） |
| 测试侧每个构造点都是**自己算**一个值再直接调 `NewPanelManager` | 例：`panel_unbound_253r6_windows_test.go:71-72`（`panelProfileLeaf("wisp-253r6-panel-profile")`）、`panel_holder_319_windows_test.go:53+62/121/246/326`、`panel_assets_299_windows_test.go:167` |
| 今天**没有任何**测试读 `.dataPath`；`cmd/wisp` **没有** `TestMain` | 现跑 `Select-String '\.dataPath' cmd/wisp/*_test.go` ＝ 2 命中，全是形状尺里的**源文本常量**（`:1513`/`:1517`，尺子还原成 HEAD 那枚后的行号）；`Select-String 'func TestMain'` ＝ 0 |

⇒ 结论一：**运行时不存在"收集这一发所有构造点拿到的 dataPath"的接缝**。值只在
`PanelManager` 实例里，而实例只有构造它的那段代码手里有；`cmd/wisp` 里既没有包级记录器，也没有
`TestMain` 收口。

## 2 三种"③"的牙口（逐种写清代价）

### ③-a 「既有 fixture 能观测到 DataPath 的地方」——按票面字面做 ⇒ **没牙**

可做的写法只有：调一次 `panelProfileLeaf("…")`／`panel319DataPath(t)` 两次，断言两次不同。
但那断言的是**helper 自己**（`panelProfileLeaf` 本来就是 `prefix-pid-纳秒`），**没人会去弄坏它**；
真出事的形状是「某个构造点不再用 helper、改成固定名」——③-a 看不见（除非那个点自己再跑一遍）。
⇒ 这形是「看着像判据、其实是同义反复」，按机主「宁缺毋造」⛔ 交。

### ③-b 逐点双跑：把某个构造点的整段驱动两遍，比 `mgr.dataPath` ⇒ **有牙，但只在一两个点**

- 可行：`panel_unbound_253r6_windows_test.go:69` 的 `startPanelFor253r6(t, disp)` 会返回 `*PanelManager`，
  调两次读 `mgr.dataPath` 就能断言「两次不同」。这正好盖住审计 A3 用来证"这形没牙"的那条回退路径
  （把 `:71` 改回固定名 ⇒ 两次相同 ⇒ 红）。
- 代价：① 每次要 `startResidentPanel` 起一条常驻线程（该文件头注释 `:34-38` 自己写了"这机器上第二个窗
  紧接第一个会间歇性 mid-create 死"）⇒ 两发同进程有 flake 面；② 覆盖**一个点**：新加的固定名构造点
  不会被任何 ③-b 用到；③ 319 holder 与 spike 两个**豁免点是故意固定的**，③-b 必须把它们的名字/身份再抄一遍
  才不误红——等于把静态豁免表在运行时重抄一份。
- 本腿**没有实现**它：它买到的覆盖只有一点，却带来 2 条常驻线程 + 一份重抄的豁免语义。

### ③-c 产码装一个观察点（真正"主流 harness"的那形）⇒ **有牙，且每点都盖，但要动产码**

`NewPanelManager` 里加 3-5 行（包级 `var panelDataPathSink func(string)`，非 nil 时回调），
测试侧在 `TestMain`/`t.Cleanup` 里把收到的每个 dataPath 记下来，收口断言「全部互不相同、
且不等于源码里任何固定名字」⇒ 每个构造点（含下一发新加的）自动被覆盖，跨过程/跨文件/常量折叠全不用管。
**这正是裁定理由①说的「主流 harness 断言行为」，也是唯一能同时满足"不追写法"和"每点都盖"的形。**

## 3 价值排序（按铁律⓪⓪⓪：最有价值，⛔ 按风险/动作）

1. **③-c（产码观察点 + 收口断言）**：价值最高——判据变成"这一发面板宿主真收到的每个 dataPath 都不同"，
   与写法无关、与点无关。代价：动 `cmd/wisp` 产码 3-5 行（本票写面明写 ⛔ 产码）⇒ **需要开写面或另立新票**。
2. **①（自写单层摘要，已现成可落）**：静态、每点都盖、确定性、零新依赖；代价是"追写法"永远差一形
   （二轮审计的 12 形就是证据）。已在仓外树验过 A1b/A1c 三形真红并点名，读数见
   `option1-readings.md`。
3. **③-b（逐点双跑）**：只盖一点、带常驻线程 flake 面、要把豁免语义重抄一遍。
4. **③-a（照票面字面）**：同义反复，⛔ 交。

## 4 我卡在哪（要一个决定，不是要我猜）

票面 15:2x 裁定「必须做运行时判据」＋「⛔ 不再往尺子里加规则」＋票面边界「⛔ 产码」三条同时成立时，
**有牙的运行时判据（③-c）落不下去**。三条里必须松一条：

- **甲：松「⛔ 产码」**（推荐，最主流）＝准我在 `cmd/wisp/panel_host_windows.go` 加 3-5 行观察点 +
  新测试件收口断言；约 30-45 分钟，两侧正控（把 `:71` 改回固定名 ⇒ 真红；还原 ⇒ SHA 同串再绿）。
- **乙：松「不再加规则」**＝按「442 改向复议」回①：`option1-ruler.patch`（57796 字节，post-image SHA `9F97…A4E2`）
  一条 `git apply` 就能落，读数已齐（A1b/A1c 真红并点名；35 形＋1 汇总全 PASS；0.68s → 0.72s；四门禁 GREEN）。
- **丙：两条都不松**＝本票收成「静态尺冻结，运行时观察点另立新票」，并把 442 的 A1b/A1c 缺口如实留在票上。
