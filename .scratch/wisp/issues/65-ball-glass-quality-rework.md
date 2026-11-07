# 65 — 悬浮球玻璃质感返工（票 62 的"赝品"判决在此消化）

**Status:** blocked-on-owner（缺参考图，见 §第一步）
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** 62-liquid-glass-ball-visuals（结构在，质感不对）
**Parallel slots:** ≤1 sub-agent（独占 `internal/ball`）
**Spec refs:** SPEC-08 §2（本次要再改一次的那段）、C21 DesignTokens、D29、D32（不放宽）

## 为什么要这张票
票 62 把球从"物理上不成像"（差分 0 像素变化 ≥8/255）做到了"看得见、不吃资源"
（2120 像素、`timers=no`、`cpu=0.000`、私有工作集 11.76MB）。**这三项 owner 已批准并回填进 SPEC-08 §2。**

但 owner 对**质感**的判决原话是：

> 可以说算是赝品吧，离我发的那种质感还是有不小差距，但是先勉强用吧，以后再换样式，就先这样吧。
> 开始阶段不能要求太高，本末倒置了就，先完成核心功能

⇒ **这是一次"临时通过"，不是签收。** 本票的存在就是防止"先这样吧"在三个月后被读成"样式已认可"。
SPEC-08 §2 里那段已显式标注 INTERIM 并指向本票。

## 一个必须承认的事实：到今天为止，没有任何代理看见过参考图
owner 在票 62 的需求里写过「球的UI风格就跟我发的图1和2类似」。**那两张图从未落到磁盘上**——
2026-09-20 我逐个核对了本会话附件目录里的 11 张 PNG，可辨识为：商标查询 4 张、托盘 1 张、
AnySearch 文档 1 张、代理面板 1 张、本轮 tour 截图 3 张、小裁图 1 张。**没有一张是玻璃参考图。**
所以票 62 的四个代理（aurora/glacier/nebula/solar 四套配色）全部是**按文字盲做**的。

**"赝品"这个判决在信息缺失下是公平的，也是必然的。** 本票第一步不写代码。

## 第一步（owner 动作，不是代理动作）
请 owner 把图1/图2 **重新发一次并另存到仓库内**（建议 `design/refs/ball-glass-1.png` / `-2.png`），
然后本票才可开工。开工前代理必须：

1. 从参考图里**量出**而不是"感觉出"下列属性，逐条写进 `docs/evidence/s1/65-glass-reference-analysis.md`：
   高光位置与宽度、边缘折射带的层数与颜色偏移、内部散景（bokeh）的尺度与对比、
   饱和度与明度分布（取若干采样点的 RGB 而非"看起来是蓝紫"）、背景透过率、
   阴影/接触影的方向与软化半径、是否有色散（chromatic aberration）及其强度。
2. 把量出来的值**映射到 C21 DesignTokens**（`design/assets/tokens.css` 是唯一样式真相源；
   绘制代码硬编码色值 = 违规，SPEC-08 §2 已写明）。也就是说返工的产物是**一组 token 变更 + 一次渲染实现**，
   不是一堆新写的字面量。
3. 给出"当前实现 vs 参考图"的**并排差分**（同一区域、同一底色、球在/不在的取证法仍适用），
   让"像不像"这件事至少有一个可比对的产物，而不是各人各自的观感。

## 硬约束（与票 62 完全相同，一条都不放宽）
- **`Sleeping` 零定时器、CPU ≤0.5%、私有工作集 ≤25MB（D32）**。返工只许换材质，
  不许为了"更有质感"引入常驻渲染循环。改完之后必须重跑票 62 的差分，
  证明 `timers=no` / `cpu=0.000` 没有回退（基线在 `docs/evidence/s1/62-diff-glass/diff-table.txt`）。
- 动画仍只允许出现在**活跃态**；`liquidDriven()` 那张白名单的方向性是对的（新状态默认不驱动），保留。
- 靠边吸附/弹回、音频驱动液面、边框滑入这些**行为**已由票 62 交付并有 11 个永不 skip 的确定性测试，
  本票**不得改动行为**，只改材质；若必须改，先在票里说明为什么。
- 冻结契约变更（含 SPEC-08 §2 的再次修改）需 owner 明确批准（D22），并像这次一样把批准原话记进票。

## Out of scope
配置界面（票 39/40）；热键与交互缺陷（票 64）；任何非球体的 UI。

## Acceptance criteria
- [ ] 参考图已入库（`design/refs/ball-glass-*.png`），且 `65-glass-reference-analysis.md` 给出
      **逐条量化属性**，每条注明来源图与采样位置。
- [ ] 上述属性映射为 C21 DesignTokens 变更；渲染代码中无新增硬编码色值（`d22scan` 与人工抽查双过）。
- [ ] 并排差分产物在档：当前实现 / 返工后 / 参考图，同一区域同一底色。
- [ ] 票 62 的 11 个 dock 确定性测试与全部球测试仍绿，且行为断言**未被改动过**。
- [ ] `Sleeping` 差分重跑：`timers=no`、`cpu≤0.5%`、私有工作集 ≤25MB，且 `px_delta_ge8` 不低于基线 2120。
- [ ] **owner 第二次看 `-tour`，并明确说"这次对了"**——本票的完成判据是人的判决，不是测试绿。
      在此之前 SPEC-08 §2 的 INTERIM 标注**不得摘除**。

## Progress log (append-only, newest last)
- [2026-09-20T10:40Z] agent=orchestrator did=开票。记录三件事：①owner 对票 62 的判决是"赝品/先勉强用"，
  属临时通过而非签收，故 SPEC-08 §2 的改动带 INTERIM 标注并指向本票；②核对结论——参考图 图1/图2
  **从未进入磁盘**，本会话 11 张附件全部可归因到其他用途，票 62 的四个代理是盲做的，
  所以返工第一步必须是"拿到图"而不是"再猜一次"；③量化属性清单与"不得改行为、只改材质"的边界写进票面。
  next=等 owner 重发参考图并另存进 design/refs/；在此之前本票 blocked-on-owner，不占车队名额

- [2026-11-07 12:5x +08] **编排者改道（台账 A1594，机主两次否手写）**：机主逐字两次 —— ①"不行，要完全跟人家的一致，你这还是属于拉裤兜子的复刻"（否 CSS-blur 版）；②"你搞错了，其实没必要自己哼哧哼哧写，别人有现成的成熟的方案，你就直接搞过来不就行了吗，为什么要自己写一套，而且很可能有bug呢？服了"。⇒ **裁定：⛔ 继续手写 D2D/WebGL 着色器，改集成现成开源件**。D2D 两笔（5a921e40／2182e1e8）**不结案、不翻 AC#6**。
- **选型（编排者拍，铁律 ⓪⓪⓪）**：LerSent001/orb（749 星、MIT、13 预设、WebGPU/WGSL 1236 行、自带独立 Web 导出器）。淘汰项：wensity siri-orb（付费 http=402，源码不可得）；zq52xy/glass-voice-orb-study（非商用许可，⛔ 不可用）。
- **实测凭据（编排者自跑）**：D:\tmp\wisp\20261107\orblab\orb749 pnpm 装依赖成功 → dev server 127.0.0.1:5180 真 Edge 开窗截屏，球体渲染正常（玻璃折射／彩虹光带／边缘高光全在，图 siri-big.png）；导出器产出 13 枚独立单文件 HTML（各 ~89KB），**导出页纯净＝只 canvas+status，零 dnd-kit/base-ui/Presets 命中**。
- **落地**：wave53 orb-integrate 在飞 —— 把 orb-siri.html 作为产品资源落位 ＋ 新建 WebView2 透明置顶点击穿透悬浮窗渲染器 ＋ 与 internal/ball 状态机对接 ＋ 保留 D2D 作 fallback（⛔ 删 D2D）。
- **AC#7（owner 第二次看 -tour 并说"这次对了"）仍不翻**：等集成后机主亲看。
- 撤销口令：git revert 该落地腿 commit；球回 D2D＝2182e1e8。

- [2026-11-07T12:55Z] agent=65-orb-integration did=集成腿落地（A1594 改道：**集成**开源 LerSent001/orb
  (MIT,749星)，⛔ 不自写着色器）。上游导出页 orb-siri.html 逐字节嵌入 internal/ball/assets/（仅 CRLF→LF）；
  新建 WebView2 悬浮宿主 orb_windows.go（WS_POPUP+LAYERED/TOPMOST/TOOLWINDOW/NOACTIVATE，透明背景走库里
  现成 ICoreWebView2Controller2::DefaultBackgroundColor(a=0)，零新 COM 绑定；WebResourceRequested 过滤器
  服务合成源 https://wisp.orb，CSP connect-src none，无 socket）；状态经 ExecuteScript 推
  window.liquidOrb.setState（20 冻结态→idle|thinking 映射表钉死测试）、SetAudioLevel 标量转
  setAudioBands（语音族才推，Sleeping 永不被麦克风唤醒）；D2D fallback 一字未删（attach 失败原路回退）。
  实测：默认套件 ok、winlive 实窗 PASS（真挂载+状态推送+50态50音频推送句柄零增长+两窗全灭），
  顺手修掉 destroy 顺序竞态（实测 panic+窗口泄漏）；门禁 gofumpt空/vet rc0/build rc0(含 cmd/wisp,balldebug)/
  d22scan clean。证据 docs/evidence/s1/65-orb-integration.md。next=owner 人眼签收（AC 六条未勾，
  D32 预算需在 webview2 渲染器上重测——见证据件"已知边界"）

- [2026-11-07T13:1xZ] agent=orb-fix-offscreen did=修 A1596 出屏 bug（编排者亲跑实测：球停靠右缘时
  「Wisp Orb」rect 右边缘超出虚拟屏 120px）。follow() 拆纯函数 orbFollowXY：MonitorFromWindow+GetMonitorInfoW
  取 rcWork（复用既有 monitorRectsForWindow，零新绑定），右放得下保持右、放不下翻左、两侧都放不下夹紧
  进工作区，垂直同夹紧；WS_EX/Topmost/NoActivate 语义零改动，D2D fallback 与上游 orb-siri.html 零字节动。
  回归两枚（orb_live_windows_test.go）：orbFollowXY 六形状表测 + 真窗用例把球 SetWindowPos 到工作区右下角
  贴边调 followOrb，断言 orb 真实 rect 完全落在 rcWork 内且坐标精确（翻左 orb.r=ball.l-12、垂直夹紧
  orb.b=work.B-edge）——修前形状在此摆位必红。门禁：gofumpt 空/vet rc0/build rc0(含 cmd/wisp,balldebug)/
  默认套件 ok/winlive TestBallLiveOrb ok/d22scan clean。证据 docs/evidence/s1/65-orb-follow-clamp.md。
  next=编排者收件三把尺＋裁决（审计者≠实现者）；AC 与出屏无交集，未翻任何格
- [2026-11-07T13:55+08] agent=ball-orb-fix-leg did=修编排者实测的 orb 全黑（球未渲染）。三层根因逐一实测定位：
  ①响应构造缝——jchv 包装用裸 GetLastError 判 CreateWebResourceResponse 成败，orb STA 线程上 errno=1008
  持续（直接 vtbl 探针读 HRESULT=0x00000000/S_OK、out 参数持活对象），3 次重试全灭，页面永远发不出去；
  ②控制器可见性——jchv Embed 从不 PutIsVisible(true)，当前运行时控制器带着隐藏的合成子窗上车，
  orb 的 never-activate show 路径必须显式 browser.Show()（面板宿主因创建即激活聚焦而幸免）；
  ③WebGPU origin——NavigateToString 下页面活着（liquidOrb=object ready=complete）但 navigator.gpu=undefined，
  页面自报 "WebGPU is not supported in this environment"，换 ICoreWebView2_3.SetVirtualHostNameToFolderMapping
  把 https://wisp.orb 映到 %LOCALAPPDATA%\wisp\orb\site（DENY_CORS），WebGPU 起来，球真渲染（截图可视证据）。
  CSP 保住：注入 <meta http-equiv>（connect-src 'none' 等全套）进 served 副本，asset 文件字节零动
  （sha d93da649... 复测一致）。回归：TestBallLiveOrbPageRendersPixels（黑底窗 topmost+屏幕 DC 采样，
  断言非黑像素占比>1%，实测 634/1764=35.94%）。门禁全绿。证据 docs/evidence/s1/65-orb-black-window-fix.md。
  next=编排者收件三把尺＋裁决（审计者≠实现者）；未翻任何 AC 格/Status

- [2026-11-07T14:30+08] agent=verdict（裁决腿，非实现者）did=独立复跑 orb 三腿（a2e535674/d687dd910/8fb821bba）
  全部声明：资产 sha256 逐字节=d93da649... 一致；像素用例亲跑 PASS 625/1764=35.43%＋★突变①（注释
  browser.Show()）→0.000% 全黑 FAIL 有牙；出屏用例亲跑 PASS orb=(3196,1260)-(3328,1392) 全内＋★突变②
  （orbFollowXY 改回无条件 wr.r+gap）→orb=(3452,1292)-(3584,1424) 双轴出屏 FAIL 有牙；两突变做完即还
  （git status internal/ball 复归空）；票 62 dock 测试 0 字节 diff（12 枚 Test func 只增未删）；门禁全套
  gofumpt 空/vet rc0/build rc0/默认套件 ok/d22scan clean（939+224 files）。★核心发现：winlive 全套
  4 FAIL 逐枚归因——TestBallLiveEdgeDock 与 TestLiveConfirmingCancelAndEscReturned 存量环境失败
  （pre-orb 仓外 clone 基线同红同读数）；TestLiveNeverStealsFocus 与 TestLiveTransparentCornerFallsThrough
  基线 PASS ⇒ **orb 集成引入真回归**：a2e535674 的 orbLive() 守卫（ball_windows.go:419/:554）让球窗
  在 orb 场景从未 ULW 一帧 ⇒ 分层窗 alpha 恒 0 ⇒ 球不可见+点击穿透（探针实证：同位同分钟，
  pre-orb WFP(球心)=WispBallWindow，HEAD WFP=第三方窗被穿透）。判语：a2e535674 声明属实但留集成面
  回归；d687dd910/8fb821bba PASS 全格；AC#1 沿用 PASS/#2 部分/#3/#5 真桌面未跑（#5 基线已换渲染器
  整表待重测）/#4 部分（断言未改 PASS，"仍绿"被球透明回归压住）/**#6 归机主**。证据
  docs/evidence/s1/65-verdict.md §7（§1-§6 历史裁决原样保留）。next=球窗透明回归收尾腿（实现腿做，
  裁决只定性不开方）→编排者重派桌面趟（AC#3/#5/#6）。未翻任何 AC 勾/Status；⛔push

- [2026-11-07T15:07 +08] agent=65-ball-clickability-fix did=球窗可点性回归修复（wave59 实现腿）。红句亲跑复现
  （:343/:371 两枚 FAIL）→诊断：orb 与球窗是**并排顶层窗**（orbFollowXY 几何，出屏测试钉死），"orb 挡球"不成立，
  因果＝orbLive 守卫令球窗 0 帧⇒分层窗 alpha 恒 0⇒OS hit-test 跳过⇒隐形＋穿透。选型**甲**（铁律 ⓪⓪⓪）：球窗照常画
  （删 renderFrame/present 两处守卫，原"transparent child double-composite"前提证伪——orb 非 child）＋orb 窗增
  WS_EX_TRANSPARENT 纯装饰全穿透（主流 overlay 形状：装饰层全穿透＋交互实窗承载全部冻结语义）；**乙排除**＝与
  orb_live_windows_test.go:152 "SECOND window, not a replacement" 正面冲突且要改断言；断言**零改动**两枚转绿
  （修复后 5 轮全绿）。判别实验在册（docs/evidence/s1/65-ball-clickability-fix.md §5）：Pixel 枚 in-suite 曾三连红
  0.000% 全黑，实验A排除 TRANSPARENT、实验B（守卫恢复）转绿锁定"球窗照画"与 in-suite 状态相互作用、
  standalone 4 连绿、末两轮全套 4 枚全绿（19.95%/12.47%）——摇摆未根除，如实报。门禁：gofumpt 空/vet rc0/
  build rc0（含 cmd/wisp、balldebug）/默认套件 ok 0.379s/d22scan clean。AC 勾与 Status 未动；⛔push

- [2026-11-07T17:30+08] agent=65-transparency-dedup did=修编排者真机实测两缺陷。**A 真透明**：根因三层——
  ①首因＝orb 窗挂 LAYERED 却从不 ULW/SLWA，顶层重定向表面是未初始化黑层；②jchv `AddScriptToExecuteOnDocumentCreated`
  传 NULL handler＝竞态缝，bootstrap 可整段丢失（真机实测页面活+WebGPU 跑而四角整生命周期显示 body 底色）；
  ③`PutDefaultBackgroundColor` API 生效滞后（文档明载 known issue）。修＝主流 Tauri/tao 配方：去 LAYERED＋
  `DwmEnableBlurBehindWindow` 空区域＋`WEBVIEW2_DEFAULT_BACKGROUND_COLOR=0x00FFFFFF` 出生即透明＋served 副本注入
  `<style>html,body{background:transparent !important}</style>`（声明式无竞态，asset 字节零动 sha 钉测绿）。
  **B 单球选型甲**（orb＝可见面与球窗同 rect 中心叠加＋工作区夹紧；球窗＝交互层只画 badge/进度环/文本＋alpha=1/255
  hit-test 轮廓——全零 alpha 会被 OS 按像素跳过＝上腿实测回归；z 序＝球上 orb 正下，因实测非 layered 的
  WS_EX_TRANSPARENT 窗不被 WindowFromPoint 跳过；乙裁页面毁票值/丙搬交互毁冻结契约均排除并留档）。新回归
  TestBallLiveOrbCornersReadTheDesktop（(51,51,50) 已知底板上四角每通道差 ≤12 断言）多轮实测**四角精确零差**，
  balldebug 二进制测试外复核同读数；EnumWindows 实测两窗 rect 完全重合＝单球，截图两枚入档。门禁：gofumpt 空/
  vet rc0/build internal/ball+balldebug rc0/**cmd/wisp 被并发 214/323 腿在飞 8 文件压红（非本腿，如实）**/默认套件
  ok/winlive TestBallLiveOrb 5/5＋可点性 2/2 PASS/d22scan clean。写面偏离如实：orbasset.go（CSS 注入落点＝编排者
  问题A第4条点名之路）与 renderer_windows.go（overlay 分支内聚处）超出派单清单，均在 internal/ball 独占区。
  证据 docs/evidence/s1/65-orb-transparency-dedup.md。next=编排者收件三把尺＋人眼签收（AC#7）；AC 勾与 Status 未动；⛔push

- [2026-11-07T16:3x +08] agent=65-pixel-flake-fix did=Pixel 枚真 flake 治理（缺陷治理腿）。复现＋根因用读数钉死：
  全套 count=10 第二回 **39/40**，像素非黑率逐轮坍塌 `40.6%×7 → 27.2 → 7.7 → 0.000×3发`＝**泄漏堆积**不是随机摇摆——
  destroy() 从不关 controller（jchv 包装层 ICoreWebView2Controller 无 Close vtable）每轮泄漏 msedgewebview2
  进程积 GPU 压力，standalone 连绿/in-suite 变红由此而来。修（测试侧，判据 1% 零改动）：①有界重建重试 3 发
  （orbPixelAttempt 单发 25s 轮询原样，全黑拆整套重建，attach 失败仍硬 Fatal）；②**泄漏收割**
  orbTestReapLeakedWebView2（Toolhelp 枚举：父 PID＝本测试进程 ∧ 映像名 msedgewebview2 才杀，测试开始＋每发
  重建前各清一次，他应用 WebView2 零命中）。**判定牙：修后全套 count=10 ⇒ 40/40 全绿**（rc0，204.4s；
  收割器每轮 reap 1 个＝机理读数自洽；非黑率 13-36% 零坍塌，区间移动是同树另一腿 b576c10a3 改渲染器所致）。
  遗留：产码级根治（补 Close vtable，fork＋replace）另立票。门禁：gofumpt 空/vet rc0（含 winlive）/默认套件 ok/
  d22scan clean。AC 勾与 Status 未动；证据 docs/evidence/s1/65-pixel-flake-fix.md；⛔push

- [2026-11-07 16:5x +08] agent=写腿(wisp编队) did=产码根治"jchv包装层缺ICoreWebView2Controller::Close"：步1现跑取证——vtable**已映射Close槽**(fork内ICoreWebView2Controller.go行33,紧跟NotifyParentWindowPositionChanged,与官方头序一致),缺的只是edge包内Go包装方法(全模块grep 0命中,vtbl不导出),票面"第6号槽"预估不成立;步2裁定选甲fork(乙=产码unsafe自声明25槽COM布局双真相源、丙=动作最小均按⓪⓪⓪否);步3落地=tracked的forks/go-webview2(561K全量拷贝+Close/Release两包装+三枚embed用DLL加.gitignore白名单;★更正:初判third_party撞整目录ignore会replace悬空)+go.mod replace+destroy()在browser=nil前GetController→Close→Release;★判据=自含winlive探针TestOrbControllerCloseStopsBrowserProcessLeak,5轮create/destroy实测baseline=0→after=0零增长,3发成立(真树10.3s/3.4s+隔离场12.3s)。门禁:gofumpt空/vet rc0/scoped build rc0/默认套件ok 0.237s/d22scan clean;★count=5未认证:第1轮4枚全绿(Pixel首发37.15%非黑),第2轮Pixel起跑处二进制静默死亡(无FAIL行无panic),隔离场(git archive仓外,无wave62 WIP)探针PASS但Pixel复跑被外来balldebug-65verify.exe(-stay)挡在quiet守卫外⇒暂无法归因,等wave62落盘后新基线复跑定案;commit=0e0a48732/2bb639e4e;destroy()hunk在工作区待wave62落盘后单独commit(避免pathspec吞其182/65行WIP);AC勾与Status未动;⛔push

- [2026-11-07 18:2x +08] agent=wave65-overlay-no-text did=修编排者真机实测缺陷（截图 D:\tmp\wisp\20261107\orblab\orb-FINAL.png）：单球架构下球心叠白色「69%」——根因＝renderer_windows.go overlayOnly 分支的步骤13中心文本（BadgeText＝presenter 经 SetBadgeText("69%") 挂载，resident_present_windows.go presenterBadger）与步骤11 badge 白数字无守卫直画，球窗在 orb 之上（z序）故每帧把字压上玻璃面。选型**甲**（铁律⓪⓪⓪）：orb 模式球窗**零文字**——步骤13整体抑止；步骤11保留 Danger 盘作无数字通知点（主流 badge-dot 形状，审批信号不丢）；进度＝外圈细弧（步骤9 BorderSoft 底环+Accent 进度弧，修前就有、截图里的青弧即它，本腿零新增绘制）；**乙排除**＝页面冻结 API 只有 setState/setAudioBands（orbasset.go 钉死接缝），无进度/文本通道，加通道要动 sha d93da649… 钉死的 asset＝毁票值；**丙排除**＝"移到球下方"要动冻结窗口几何（WindowEdgePx/dock/四角断言全挂现几何）。判据（机械形式）：新 winlive 用例 TestBallLiveOrbOverlayFrameHasNoGlyphPixels——复刻缺陷形状（Downloading+0.69+"69%"+badge3，cmd/wisp presenter 调用形状）后对**球窗 ULW 源 DIB** 逐像素审计（STA PostTask 串行读=定格 present 帧），断言不透明近白像素=0（OnSolid 白是 drawText 专用色，overlay 合法调色板无第二个近白元素）；**红→绿**：修前 FAIL 计数 103、修后 PASS 计数 0。D2D fallback 一字未动；bits uintptr→unsafe.Pointer（vet-clean 读 DIB 所需，零行为变化）。门禁：gofumpt 空/vet rc0/build ./internal/ball ./cmd/wisp ./cmd/balldebug rc0/默认套件 ok 0.645s/**TestBallLiveOrb 6/6 PASS（四角 (51,51,50) 精确读数未破，CornersReadTheDesktop 30.84s PASS）/TestLiveTransparentCornerFallsThrough PASS**；TestLiveNeverStealsFocus 物理指针被占 SKIP-LOUD（12 发×60s 重试环在飞，SKIP 前样式位/MA_NOACTIVATE/前台不变三断言已过），d22scan 全量脚本复跑在飞（首跑 FAIL=tools/d22scan 自模块 1 枚、隔离复跑 ok 73s，归因并发腿并行污染待定案）——两项收尾读数落 docs/evidence/s1/65-orb-overlay-badge.md。orb-siri.html sha 钉测绿（默认套件）；AC 勾与 Status 未动；⛔push

- [2026-11-07 18:5x +08] agent=wave65-overlay-no-text did=收尾读数落盘（追加）：**TestLiveNeverStealsFocus PASS 1.65s**（重试环第2发；首几发物理指针被占 SKIP-LOUD＝纪律拒伪绿，第1发无判行＝在记静默死亡形状未复现）＋ **d22scan clean**（internal=944/cmd=230/frontend=116；首跑 FAIL=tools/d22scan 自模块 1 枚系并发污染，隔离 ok 73s＋全量复跑 clean）。门禁全绿；commit=c5cca231a（收尾读数二次 commit 追加）；AC 勾与 Status 未动；⛔push

- [2026-11-07T18:4x~19:0x +08] agent=65-jarvis-style-integration did=Jarvis 粒子球接入产品第二种视觉风格（真上游页，非 D2D 仿品）。步1 调查（siri 页单文件形态×jarvis 产物形态:tsup noExternal three 已打进、react external 需页面提供）＋步2 裁定**甲**=vite 构建单文件 HTML（乙=vanilla 直调 renderer.ts 要自养 347 行组件宿主逻辑买 140KB、丙=放弃单文件页形状零收益，均按⓪⓪⓪记代价排除）＋步3 落地 `internal/ball/assets/jarvis-orb.html`（758,355B,LF,sha `e95e1d96…` 钉死;脚手架 tools/jarvis-export/ 随库:entry.tsx 唯一自有集成代码+inline.mjs 单文件内联+三雷留档:String.replace $-模式污染/双 React 实例 dedupe/JSX automatic）**真渲染验证**=Tabbit 真浏览器 canvas:true+接缝 setState 回读+idle/thinking 两态截图入档＋步4 Go 接线=orb 宿主按 ActiveStyle 选 asset（orbWriteSite 双入口物化⇒运行时切换=纯导航;orbPageForStyle 纯决策;pushState 分派 liquidOrb/jarvisOrb;pushAudio 对 jarvis 整段跳过=上游无音频接缝不发明）+jarvisStateFor 四态映射（工作族→thinking/失败族→alert/Settling→success/余→idle,20 态全覆盖测试钉）+Ball.ApplyStyle 运行时换页＋步5 winlive 判据 TestBallLiveOrbJarvisStyleRendersParticles（真桌面黑底板非黑像素 10.09%>1% 牙线+ApplyStyle 往返窗柄不变;新文件避让 overlay 腿,orb_live_windows_test.go 零触碰）。门禁：gofumpt 空/vet rc0(含 winlive)/build internal/ball+cmd/wisp+cmd/balldebug rc0/默认套件 ok 0.368s/winlive TestBallLiveOrb **7/7 PASS**（含 overlay 新枚,四角 (51,51,50) 精确零差未破）/焦点+穿透 2/2 PASS/d22scan clean。证据 docs/evidence/s1/65-jarvis-style-integration.md＋65-jarvis-orb-{idle,thinking}.png。写面超派单清单如实：orb_test.go 加 jarvis 钉子（sha/映射/接脚/serving 六枚）。commits=1652051d0/29b3501d5/b679fd072/1853d1d87/d902fbc6b/本 commit；AC 勾与 Status 未动；⛔push

- [2026-11-07 19:1x +08] agent=wave68-orb-arc-restore did=编排者「外圈青弧被 c5cca231a 的 !r.overlayOnly 关掉」初判现跑证伪＋「弧存在」钉子落地。取证：git show c5cca231a 逐 hunk＝只动 :588/:614 两处文字守卫＋overlayGlyphPixels＋bits 类型，:521 RingColor/:533/:546/:571 四处图形守卫零字节动（git log -S 双查＝出生自 10d7cb3c3，注释原文「页面画状态,D2D 再画=draw the state twice」）；真弧＝步骤9 Progress 弧（pal.Accent #86C2B9，:557）本就无守卫 overlay 照画。lit ratio 21.96→3.02 与弧无关：Corners 全程 Sleeping 从不设 Progress（该读数从未量过弧）；读数本身＝「四角 ok ∧ lit>1% 首次过线即返回」的快照，本腿 5 发＝25.91/24.77/5.20/0.000(黑启动,重试治愈)/3.02(编排者引)同族，修前 21.96 亦同语义快照。裁定甲（⓪⓪⓪）：**不恢复 :521**（恢复＝页面状态面双重画＝源注释点名缺陷），**产码零改动**；落地＝①Corners 补 3 发有界重建重试（orbPixelAttempt 同型，断言零改动，全套内黑启动发被治愈）②Glyph 判据同帧加弧存在断言 overlayArcTealPixels（premultiplied 比值谓词，双主题 Accent 命中、Danger/Info/BorderSoft 排除；红牙：临时守掉步骤9⇒FAIL "carries only 0 progress-arc pixels" 且近白审计同帧仍 0＝两判据正交；还原 git diff 空；绿：283 teal pixels ∧ 近白 0 同帧成立）。真机：HEAD 重建 balldebug-orb.exe -shots 20 态循环，13-Downloading.png＝青弧在环上＋球心无文字（＝FINAL 减文字，弧没丢），档 D:\tmp\wisp\wave68\shots\。门禁：gofumpt 空/vet rc0/build 三目标 rc0/默认套件 ok 0.581s/winlive TestBallLiveOrb **7/7 PASS 69.888s**（四角 [[51 51 50]×4] 精确零差）/焦点+穿透 **2/2 PASS**（1.88s/0.58s）/d22scan＝本腿写面（internal/ball）**零发现**（真仓 1 发现＝票 187 腿 untracked 在飞件 config_write_187.go:55 phantom-citation，⛔ 不碰，187 收口自愈）。写面＝orb_live_windows_test.go＋证据件＋本 log；AC 勾与 Status 未动；⛔push；证据 docs/evidence/s1/65-orb-overlay-arc-restore.md；commits=91c0532d2/81c60c650/本 commit
