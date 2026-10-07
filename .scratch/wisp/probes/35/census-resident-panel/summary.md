# census-resident-panel —— "机主双击能真点开面板"链路普查（只读，锚 `af68866`，交件腿自量）

> 本腿全程未跑任何 go 命令（cmd/wisp 的 go 窗口归写腿 `35-r4`）；所有读数＝现读文件行号或现跑 grep。

## 0. 先顶回派单一处过期前提（第 107/122 条规矩，具名报回）

派单说"**面板在常驻进程里还没被装配起来**"。**盘上原文不成立**（本腿逐行复量，非转述旧件）：
- `cmd/wisp/resident_windows.go:151`＝`rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)`（注释 `:141-143` 逐字："Building the manager is the first non-test call of NewPanelManager in this repository"）
- `cmd/wisp/resident_windows.go:157`＝`panel = startResidentPanel(rt.Registry, rp)`
- `cmd/wisp/panel_resident_windows.go:253`＝`return NewPanelManager(disp, assets, dataPath, withGeometrySource(...)), nil`（尺＝`*PanelManager`＋`NewPanelManager(` 调用形状，非 test 调用点现量＝本文件 1 枚；`residentPanel.` 无 Show 直调，见下）
- 手势那一跳也在：`cmd/wisp/resident_windows.go:217`＝`startResidentBall(..., withPanelHost(func(via string) bool {...`，执行体 `cmd/wisp/panel_resident_windows.go:409`＝`err := rp.mgr.Show(context.Background())`。
- 入口＝双击那一支：`cmd/wisp/main.go:66`＝`runResident()`（无参数启动）。

⇒ 所以 Q1 的答案不是"装配断在哪一行"，而是**装配之后还断三跳**（见 §1）。引用"生产调用者 0 枚"的旧表（`167-c2`、`wisp-permission-mode-rulings`、票面 AC#1..AC#4 的原文句子）从今天起**全部过期**，引它们的地方要翻。〔票 33 的 `33-r5` 台账句已预告此事，本腿是复量确认，不是转述。〕

## 1. 装配缺口：今天真正断的三跳

链＝双击 → `main.go:66 runResident()` → `resident_windows.go:151/157` 建宿主＋起 STA 线程 → 热键/托盘 → `panel_resident_windows.go:409 Show` → `panel_host_windows.go:480 Show` → `bringUp(:318)` → `firstRoundTripLocked(:425)` → `serveEntry(:427/:453)`。

1. **页面字节**（断在最外层）：`git ls-files frontend/dist`＝只跟踪 `frontend/dist/.gitkeep`；`frontend/embed.go:19` 是 `//go:embed all:dist` ⇒ `newResidentPanelManager` 里 `panel.BuiltinAssets()`（`panel_resident_windows.go:238`）拿到占位 bundle，`serveEntry` 走 `:455` 的"not built"错误分支 → `serveNotBuiltNoticeLocked`（`:441`）显示逐字页面"panel assets unavailable: the embedded bundle is not built."（`:449`）。**机主今天双击能开出的就是这枚告示页。** 票 33 AC#12／票 248 AC#9 都钉这件事，两格都写明"产物由界面侧那枚 agent 产出，本编队⛔ 不写 `frontend/**`"→ **按住，等机主裁页面分支（`D:/wt/fe dsh/feat/frontend-p0-v2`）合入＋构建带页那一格**。
2. **常驻腿没有快照泵**（内容接上后 composer 全维无输入）：尺＝`panel.SnapshotPump`＋`NewSnapshotPump(` 调用形状，非 test 现量**只有 1 处**＝`cmd/wisp/run.go:699 rt.pump = panel.NewSnapshotPump(...)`（`wisp run` 那一支）。`cmd/wisp/resident_windows.go`／`panel_resident_windows.go` 零引用 ⇒ `167-c2` 那句"快照泵无装配连线"对**常驻腿今天仍成立**。归票 145 AC#2／167-a3 的四步表（`12082` 行台账），第 0 跳＝装配根把泵接到常驻。
3. **出向判据的 CI 牙**（票 35 `:75(c)` 新框）：`bringUp→installPanelTransport(:693)→w.Init` 那条真窗绿发**只在 `-tags winlive` 的本机台件里活**（`panel_transport_live_35v2_windows_test.go`）。

一处**派单里也过期了的前提**：票 33 AC#13"探测页盖掉真页面"的**时序半**——现读 `bringUp` 已是**先探测后供页**（`panel_host_windows.go:425 firstRoundTripLocked` 在 `:427 serveEntry()` 之前；票面 AC#13 引的旧行号 `:221→:227` 反过来），`33-r2` 题面②那半看起来已落，**但框还开着**（勾归非实现者）。落地腿去勾之前要自己复跑这一把。

## 2. 双角色/双进程形状：ⓐ 已经是既成事实

| 形 | 状态 | 动哪些文件 | D38(e) 十步 | "装配根唯一接缝" | 代价 |
|---|---|---|---|---|---|
| ⓐ 常驻腿自起泵（专用 STA 线程＋库 `Run()`） | **已落地**（票 33；线程＝`panel_resident_windows.go:154 reg.Spawn(panelSTAName,...)`，泵＝`:297` 库 `Run()`，停机＝`resident_windows.go:158 defer panel.stop()` 注释 `:146-150` 逐字"the frozen D38(e) ten-step order ... is not touched: no new step"） | 无新增 | 不撞（defer LIFO 排好，面板线程最后退） | 合形：宿主由装配根 `resident_windows.go:151` 建、以函数值下发球腿（`:217`），与票 246/256 同shape | 0（已完成） |
| ⓑ 面板宿主独立进程＋IPC | 未选，且编排者裁决 P1（`resident_windows.go:134-135` 注释具名）已选ⓐ | cmd/wisp 新进程对＋proc 名册（`observe.ResidentNames` 六枚冻结）＋D38 步序 | **撞**：新进程要有停机源，第 7 枚常驻协程/stop-request 钩子两条都被票 228 AC#2 射程标注为"要人裁形" | 新增第二接缝 | 高；今天无理由重开 |
| ⓒ 维持现状只走 `wisp run` | 并存中（run seam 的 `newComposerDispatchChain`/泵还在 `run.go:648/699`） | 无 | 不撞 | 不撞 | 无；但它救不了"双击"——run 支是控制台一次性进程 |

⇒ **所以呢**：Q2 不该再被当选型题；该问的是"§1.2 的泵接进ⓐ的装配根"这一格何时可派。

## 3. 每形的既有钉与既有红（现量行号）

- **双向键对账两枚红**（`internal/panel`，非 cmd/wisp）：
  - `internal/panel/approval_test.go:129`（`TestApprovalCardViewJSONKeysMatchFrontendTypes` 族）：`subtract(goKeys, tsKeys)` 报「Go Snapshot emits **[instructions tasks]** ... does not declare」。
  - `internal/panel/composer_test.go:74`：同句＋「Go ComposerState emits **[git currentModel modelKnown credentialState credentialKnown]** ... does not declare」。
  - 尺的对面在本分支：`frontend/src/lib/panel.ts:112 interface ComposerState`／`:129 interface PanelSnapshot`，`grep -c "instructions|credentialState|currentModel|tasks"`＝**0**；台账 `12971` 行那把点号尺量过**页面分支 `dsh/feat/frontend-p0-v2` 上这十枚键也全 0**〔旧读数・本腿未跨树复跑，仅行号存在性复认〕⇒ **落地时必须同批处理＝页面声明面先加键**（跨会话带话或具名解冻 `frontend/**`，台账 `13376` 行两条合法路）；⛔ 不许删 Go 字段换绿（派单与台账同令）。
- **cmd/wisp 的 winlive 红族**：`-tags winlive` 用例（`panel_host_windows_live_test.go`、`panel_transport_live_35v2_windows_test.go`、`panel_geometry_255_winlive_test.go`、`resident_approval_live_246_windows_test.go` 等，`grep -rln winlive cmd/wisp`＝10 枚文件）默认构建里根本不参编＝不进普通 CI 的红名册；它们的"红"只会出现在本机带 tag 复跑时。任何落格动 `bringUp`/transport 都要**同波复跑**（票 35 `:75(c)` 框的原话）。
- **恒真面两枚**（票 35 新框 (a) 半）：行为裁判 `panel_transport_35r2_test.go` 的夹具对 `native(message)` 丢 receiver 与守卫删半**都绿**（W1/W2 突变，`35-v2` 交件）；落地该框时自造两种错法各带一枚具名红。

## 4. 可派发性表（按写面互斥＋判据有无被测物两把尺过；⛔ 未按票号分桶）

| 格名 | 票面行号（`33-panel-host-c27.md` 等） | 凭据要什么 | 写面 | 今天可派？ | 撞谁 |
|---|---|---|---|---|---|
| 33 AC#14 Go→页投递者 | `:42` | 指名 `installPanelTransport`(`panel_host_windows.go:693`)＋真窗 correlationId 绿发；**时序半疑似已落，勾归非实现者** | 只读验收件（`docs/evidence`/probes） | **可派**（验收腿，非实现者） | 不写码＝不撞 `35-r4` go 窗口；引 winlive 读数需本机带 tag 复跑→**撞 35-r4 窗口，按在其后** |
| 33 AC#13 探测页盖页面 | `:39` | ①现读次序（§1.3 已量=已换向）②"最终文档含真入口"会响断言③换序必红反控 | `cmd/wisp/panel_host_windows.go`＋同包测试 | **可派**（写腿），但 AC#12 产物未到时"真内容"断言只能用受控夹具并明写局限 | cmd/wisp go 窗口（35-r4 在飞） |
| 常驻泵接线（145 AC#2 前置的第 0 跳） | 145 `:26`；台账 `:12082` 四步 | 真装配根里 `NewSnapshotPump(` 出现＋reader 来自真源（AC#6 反向判据） | `cmd/wisp`（装配根+panel_pump.go）＋`internal/panel/pump.go composer.go` | **可派写面，按住时序**：新键落地即触发 §3 两枚对账红 ⇒ 必须与"页面声明"同批，那半今天只有一条合法路（带话/解冻） | 键对账两枚红；cmd/wisp 窗口 |
| 35 `:75`(c) winlive CI 载体 | 35 `:22`（框内） | 要么 CI 可达替身，要么具名〔仅本机可量〕＋逐波复跑台账 | `.github/workflows/ci.yml`＋可能 `scripts/`；自托管 runner 今天在线（`gh api`现量 `wisp-selfhosted-01 online`） | **可派**（写面与 cmd/wisp 互斥，go 命令都不需要） | 无——但⛔ 别把"yaml 里有"读成"跑过"（第 109 条）；新步要带 `if: !cancelled()` 类守卫并自己核 |
| 228 AC#8/AC#10 假自述三处＋godoc | 228 `:38`/`:40` | 改成带条件事实句；⛔ 词面型新仪器 | `cmd/wisp/main.go resident_ball_windows.go` | **可派**（纯注释/slog 句），但仍占 cmd/wisp 文件读窗 | 35-r4 读数窗口（文件级错峰即可，无 go 命令） |
| 33 AC#12／248 AC#9 embed 真产物 | 33 `:37`／248 `:47` | 能力型判据（条目数/入口文件真存在），⛔ 退出码 | `frontend/**`（本编队禁写）＋embed 侧测试 | **按住**：具名理由＝产物侧未落定（票面两格自写"谁把 dist 填上落定前勾不了"），且是同一归口两面 | 机主裁页面分支合入 |
| 248 AC#4 真机整链 | 248 `:29` | 面板点设置→保存→下任务真用上 | 跨 cmd/wisp＋frontend | **按住**：票面自写"票 33 宿主真起来前不许勾"——宿主已起来，但**页面没有**，点击链第一跳就不存在＝**被测物缺失的死格**，现在派＝任何腿只能靠放宽判据交件 | AC#9 归口 |
| 256 AC#1/#2 | 256 `:50-51` | 复用票 224 r2 授权仪器＋AC#0 先落一枚具名 `A##` 批准 | `cmd/wisp resident_approval*.go` | **按住**：AC#1 票面逐字"只在 AC#0 交完并由编排者落具名 A## 之后才许动" | 编排者批准（非技术） |

## 5. `:75(c)`：真窗读数进 CI 了吗——没有，历史上也没有

- `.github/workflows/ci.yml`：`grep -n winlive` **零命中**（RC=1）；全文件 `grep -n -- -tags` 也零命中（go test 步骤全走 `scripts/runtests.sh` 那形，无 tag 组合带 winlive）。
- `scripts/*.ps1`：`grep -rn winlive scripts`＝零命中；`scripts/slo-check.ps1` 存在但无 `tags` 行。
- 历史尺：`git log --oneline -S winlive -- .github/workflows/`＝**空** ⇒ ci yaml 历史上从未有过 winlive 步骤。
- 机器事实（只读）：`gh api .../actions/runners`＝`wisp-selfhosted-01 online`，labels `self-hosted Windows X64 wisp-slo`；`ci.yml:798 runs-on: [self-hosted, wisp-slo]` 的 SLO 作业有载体。**结论只到"CI 从未出现过 winlive 步骤的 success/failure"**（yaml 无＋历史无），载体本身在线。

## 6. 机主视角（零术语；可直接给他看）

| 今天的样子 | 缺的那块落地后会怎样 | 不做的后果 |
|---|---|---|
| 双击图标，程序在后台跑起来了，按面板快捷键也确实会弹出一个窗口——但窗口里只有一行英文告示："页面内容还没打包进来"。 | 窗口里出现真正的面板：输入框、状态、结果，都是看得见能点的东西。 | 他永远只能看到那行告示；功能全都"在"，但没有一扇能用的门。 |
| 面板就算打开了内容，上面各项状态（用的哪个模型、凭据录没录、任务清单）今天在那种"双击打开"的进程里是空的，因为喂数据的管道只接在了另一种"命令行跑任务"的用法上。 | 打开面板就能看见当前模型、任务、凭据状态，并随运行刷新。 | 界面是死的：显示不出来或永远空白，他会以为程序坏了。 |
| 鼠标点了没反应算不算坏了——今天算，而且不算意外：负责"从界面点一下→后台动一下"的验收，有一条只能在开发这台电脑上手工验，没进自动检查。 | 每次提交代码，自动检查会把"点开、点了有反应"这条链自己跑一遍。 | 只在这台电脑上手工验过一次＝换台机器、下次改动都可能悄悄坏掉，没人知道。 |
| 他眼睛能签收的唯一办法是"真窗口里看见真页面"；那个页面源码在另一棵没合的分支上，机主已裁"先不合，等构建带页面那格落地"。 | 构建产物里真带上页面字节，出货的程序双击就有效果。 | 现在出货 exe 有闸门检查"带没带页面"，但闸门的另一头还没人把页面交进来。 |

## 7. 给编排者的"所以呢"汇总

1. **最可派**＝35 `:75(c)`（写面 `.github/**`，与 cmd/wisp go 窗口互斥、无被测物缺口）与 33 AC#13（写腿，夹具判据可先立牙）；33 AC#14 派**验收**腿翻勾（时序半已落）。
2. **泵接线**写面可行但**必须与页面声明同批**，否则交件即撞 §3 两枚对账红——今天合法路只有带话/解冻两选一，先问机主。
3. **248 AC#4 是死格**（被测物缺失），⛔ 别派；256 按 AC#0 批准；33 AC#12/248 AC#9 等页面分支裁合。
4. 派单第一句的前提要改：**"面板没装配"已过期**，缺的是页面字节＋常驻泵＋CI 牙三跳。
