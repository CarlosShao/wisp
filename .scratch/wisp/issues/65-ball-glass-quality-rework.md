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

- [2026-11-07 19:2x +08] agent=65-jarvis-style-bg-fix did=修 Jarvis 粒子球页底色/配色＋winlive 判据加严（机主实测凭据：上游 demo 深黑底亮青环 vs 我们产物页浏览器白底灰白环）。步1 现跑取证：demo 深底＝上游 examples/basic/index.html body #050913（App.tsx 零背景，组件默认 palette=cyan 不依赖宿主底）；我们的入口 palette 已显式传 cyan、three renderer 与 demo 同段零差异（alpha:true/清屏 a=0/SRGB/无 toneMapping），白底唯一根源＝导出模板 body 写 transparent（浏览器透白＋AdditiveBlending 洗白）；产品透明架构五证据核清＝orb_windows.go DefaultBackgroundColor a=0/无 LAYERED＋orbasset.go 注入 transparent !important（压过页面自身底色）＋siri 页 #030409 同构先例＋四角 ±12 测试钉。步2 裁定（⓪⓪⓪）：导出页自带 demo 原值 #050913——产品里被注入层清成透明（透明契约不靠页面），浏览器直开呈现 demo 原貌；⛔ 绝对暗阈值黑底（那才是破坏透明架构）。步3 重建 jarvis-orb.html（759,031B 纯 LF CR=0，sha 740c4868…同 commit 换钉 orb_test.go，外链 0 命中）。步4 判据加严＝废 non-dark>1% 白底下恒真的弱牙，换三牙：①四角内缩 3px 在已知 (51,51,50) 底板 ±12 内（桌面一致形态）；②青色主导像素（g-r≥12 且 b-r≥12 且 b≥80，上游 cyan 色板 0x38f4ff 系）占比>2%（实测 7.91-8.32%，4 倍余量；洗白灰环/重上色归零）；③活性由②蕴含。步5 门禁：gofumpt 空/vet rc0(含 winlive)/build 三目标 rc0/默认套件 ok 0.403s/★Jarvis 新牙 PASS 8.65s（四角精确零差＋cyan 8.12%）/★TestBallLiveOrb 6/6 PASS 63.47s（四角 (51,51,50) 零差未破）/d22scan clean（首跑 FAIL=panel-187 在飞件幻影引用瞬态，其证据件落盘后复跑 clean，非本腿写面）。写面超派单清单如实：orb_test.go 仅换 sha 钉一行（asset 变钉必须跟，README 纪律）。commits=fe4fec342/d8da4e83e/5762be0d9/本 commit；AC 勾与 Status 未动；⛔push；证据 docs/evidence/s1/65-jarvis-style-bg-fix.md

- [2026-11-07 19:3x +08] agent=65-style-switch-config did=机主点名的原始要求落地（此前漏了：设置项外观可配切换球风格，默认玻璃球，切Jarvis才换球的风格）。步1取证：ApplyStyle签名/机制现跑核清（tokens.go:391 SetStyle进程级+orb_windows.go:457 STA换页导航，orbWindow:280创建种子ActiveStyle⇒New前SetStyle无闪；**产品调用方0仅测试**=缺的一环确证）；config键[ball].style早已在（30d8f2e38：schema词表+validate枚举+默认glass）但setter/reader零命中；面板形状=213 reader-gated+248/313/187写路径三roster。步2裁定**甲**完整四段（乙=重启才换背叛机主原话+闲置已建成的ApplyStyle机器；丙=重启丢失违反D36三档模型），按187冲突拆两段：187未提交hunks压composer/pump/config_handlers⇒快照键/并集行/pump reader/run.go注入四小段留下一腿（如实报，⛔吞同文件他腿hunks）。步3落地=①config setter SetBallStyle(writeOneKey守卫写+枚举config侧把关)+BallStyleOptions reader（schema词表零自造，**新件零玻璃球字面量**，默认三重成立）；②cmd/wisp启动ballStyleSeed65照hotCfg258按次重读[ball].style+ball.New前SetStyle+装配后ApplyStyle幂等复同步=**产品调用方0→实**；③写路径预接线=config_write_65.go roster(FieldBallStyle)+ApplySetting加ball_style case+ballStyleApplier65注入缝（球装配installBallStyleApplier65挂Ball.ApplyStyle），档位实测=applier在答now不在答restart（187静态标法不适用，档取决于本进程有无球）。步4四牙全绿（Test65*改名对齐派单门禁）：默认glass/写jarvis ActiveStyle真变/写回glass真变/突变摘config读取默认用例红（FAIL 2枚逐字复跑+还原全绿）。步5门禁：gofumpt空/vet rc0/build rc0/config ok 2.197s/ball ok 0.242s/panel rc1三红逐枚归属他腿（2预定红[commands modelSwitch]+C21 tokens非我面）/cmd/wisp Test65 4/4 ok 0.226s/**winlive TestBallLiveOrb 7/7 PASS 50.797s**（四角精确零差未破）/d22scan clean rc0；LF七件逐一核过。Bug自报两处：config断言初版双引号形状不符canonical单引号写盘（红后按真形状改钉）；初版测试名不匹配门禁pattern改Test65*。收件边界如实：面板设置页面上**还看不到**外观项（门拒ball_style+快照无键=四小段留下一腿），但启动换风格+写入路径全线打通有测试钉，ApplyStyle零产品调用缺陷已消除。commits=110e3bdcf/bc48ae6c7/本 commit；AC勾与Status未动；⛔push；证据 docs/evidence/s1/65-style-switch-config.md

- [2026-11-07 19:4x +08] agent=65-style-switch-config did=追加更正（append-only，证据件§6）：并发收口腿（fe0f740d6/1b25d931c，187已落007de34f0）已把本腿§4留的四小段落掉三件——knownWritableField四表并集+两附加并集缝+composer快照键BallStyle+pump reader+ball_style_65.go，并如实改写本腿config_write_65.go过时注释；仍差run.go PumpSources注入（该件脏在飞，收口腿未提交）。组合态本腿现跑复核：vet三包rc0/Test65 4/4 ok 0.212s/config ok 2.446s/panel契约红现读数=[commands modelSwitch ballStyle]（Q-51既定程序，非缺陷）。§4为写作时点属实，以§6为最新状态；本腿三commits=110e3bdcf/bc48ae6c7/635922c76+本更正；AC勾与Status未动；⛔push

- [2026-11-07 19:5x +08] agent=65-style-switch-closeout did=收口腿终态（§6「仍差run.go」已了结：run.go注入落在7858108a0，是d86e591ba的祖先，四小段+cmd/wisp半全在库）。步2a/2b=knownWritableField四表并集+acceptedWritableFields/takesValueFor两附加缝（步1新发现，187三处全接同款）+composer快照键BallStyle（wire名ballStyle，§2占位OrbStyle未采用，对齐FieldBallStyle线名）+pump reader+fill+ball_style_65.go载体+run.go装配根注入BallStyle reader（读rt.settings同一个mgr）。步3=panel门级三测（roster四表且三旧名册零污染/门级config.set ball_style直达leg含正负控/ballStyle键reader-gated双向）+run.go源钉；★突变牙④摘writableField65项：门级测试红rc=1红行逐字入证据件§3.4（拒绝文案中ball_style仍在列=文案缝独立起作用的活证据），sha256 2048692e…1616逐字节还原复绿。步4组合态终态门禁：gofumpt九件点名空/vet四包rc0/build四包rc0/config ok 3.136s/cmd/wisp筛选5/5 PASS rc0（4牙+源钉）/★winlive TestBallLiveOrb 7/7 PASS rc0 56.819s（Jarvis ApplyStyle真窗往返在列）/d22scan clean rc0（internal 956+cmd 234）/panel rc1三红逐枚归属（契约红读数=[commands modelSwitch ballStyle]＝Q-51既定程序+C21 tokens非本票）。终态链：快照键ballStyle→门放行→ApplySetting→SetBallStyle→applyBallStyle65→Ball.ApplyStyle即时生效（有球now无球restart）；ApplyStyle产品调用者=2非0；默认玻璃球三重成立。commits=fe0f740d6/1b25d931c/7858108a0/432d5ab7b/本笔；AC勾与Status未动；⛔push；证据§3（节序§1§2§4§5§6§3系双腿append-only交错，未重排）

- [2026-11-07 20:0x~20:1x +08] agent=65-jarvis-transparent-in-product did=机主真机实测两缺陷收口。步1取证：注入层siri/jarvis同一函数同一串（orbasset.go orbServeDocument对两style一视同仁,注入html,body{background:transparent !important}于<head>后,无<head>退化doctype注入）＋页面底色仅一处body #050913无!important（注入必胜）＋canvas本身alpha:true/setClearColor(0,0)排除不透明canvas假设＋D2D badge=renderer step11不分style（overlay模式保留无数字红点=早前裁定,jarvis粒子上即机主红点）。★真机复现：同机同码20:02现跑winlive jarvis PASS四角(51,51,50)逐通道零偏差青粒子7.08%；机主19:46 exe与HEAD代码零差异（19:56收口commit只动2文档件）。★裁定：四角白=透明球后面的白色NBA2K桌面（漏底只会漏暗色#050913不可能白;机主自记连远处桌面都读255;球心170,203,204=青粒子合成于白面读数）——jarvis页产品里此前已是真透明,真实残余缺陷仅badge圆盘。步2落地（⓪⓪⓪选构造即透明非最小动作）：①模板body改background:transparent重建jarvis-orb.html 759648B纯LF外链0命中,sha 66eb4b80…同commit换钉orb_test.go——注入层从承重链退役为belt,「注入优先级被顶掉」整类风险消除（代价:浏览器直开artifact粒子偏淡,demo观感留jarvis-src）;②renderer step11条件改v.BadgeCount>0 && !(r.overlayOnly && activeStyle==StyleJarvis)——jarvis overlay不画badge盘,glass与D2D fallback面不动。步3判据加严：①白角之谜=测试底板(51,51,50)已知色窗非特设页面背景,真机白=NBA2K盖桌面,「四角==环境桌面色」牙在活动游戏面上不可证伪（两次采样间像素自变）,替代=served-page牙（读产品WebView2实载的%LOCALAPPDATA%\wisp\orb\site\jarvis.html,断言belt在+body transparent构造+无#050913）+四角已知地板牙保持;②无badge牙：SetBadge(3)后STA走ULW源DIB数红主导px(α≥100且R≥1.5G且R≥1.5B,Danger双主题满足,预乘缩放不变;BorderSoft/Accent/Info/剪影全不进色带)断言0,glass阳性对照245px(floor50)证审计不瞎。步4门禁：gofumpt空/vet rc0/build三目标rc0/默认套件ok 0.300s/★winlive TestBallLiveOrb 8/8 PASS 58.605s（旧7+新1,四角精确零差未破,新产物上cyan 7.39%）/焦点+穿透2/2齐绿（NeverStealsFocus第1/2次SKIP-LOUD=环境拒SetCursorPos,第3次PASS 1.21s,如实记非首跑绿）/d22scan clean。写面=tools/jarvis-export/index.html+assets/jarvis-orb.html+orb_test.go钉+renderer_windows.go+orb_jarvis_live_windows_test.go+证据件+本log；assets/orb-siri.html钉死sha未动；cmd/wisp与panel与config与tools与frontend与design零触碰；AC勾与Status未动；⛔push；证据 docs/evidence/s1/65-jarvis-transparent-in-product.md；commits=3104ac37a/步2/步3/本笔

- [2026-11-07 20:3x +08] agent=65-ball-style-cli did=机主今天就能用的CLI切球路径落地（`wisp config`子命令，主流git config/kubectl set形状）。步1取证：main.go分发=switch args[0]非冻结契约（133门禁只要求usage对账+legCovers钉+minLegs12下限现13）；SetBallStyle(settings_65.go:46)守卫在写前可直接CLI调；写盘真链=writeOneKey→WriteSection（task提示的panelwrite.go不存在，如实记）；★档位实测=[ball]热档但换球只挂两钩（面板applyBallStyle65进程内+开机种子resident_windows.go:419/432）⇒CLI另一进程答不了now，收据诚实答重启档；另证neutral方向零SECTION-WRITE audit行⇒直用SetBallStyle零假账零新增internal/config面。步2落地=cmd/wisp/config_cmd_65.go（get/set/list+封闭键名册ball.style+默认值走NewDefaults零字面量+缺文件set走ensureFirstRunConfig唯一首建口）+main.go usage行与case+legCovers133注册TestConfigLegGetSetGuarded65→cmdConfig。步3五牙：get默认glass（stderr明说是默认值）/set jarvis读回/set glass切回/非法值紫色rc1守卫具名句+config.toml逐字节相等/★突变=摘枚举守卫牙④红（validate第二层仍拒写但CLI不回前置守卫具名句，红行逐字+sha256 3d87f824…5b06f逐字节还原复绿），附牙：键名册拒绝+list双态。步4门禁：gofumpt四件空/vet rc0/build rc0/cmd/wisp Test65|TestConfig ok 0.394s/133三测PASS（新config腿双向对账核过）/internal/config ok 4.457s/d22scan clean（956+236）/LF四件核过。边界如实：设置页那行仍归Q-51（frontend零写面）；CLI档位=重启档（收据已明说）；cmd/wisp两枚323未跟踪件系他人，未碰未收。AC勾与Status未动；⛔push；证据docs/evidence/s1/65-ball-style-cli.md

- [2026-11-07 23:2x +08] agent=65-jarvis-dark-disc did=机主裁定「给球自己的深色圆盘底」落地（机主原话「把Jarvis风格搞出来给我看」＋HEAD 7873c8328 裁定：盘内不透明给粒子对比度、盘外仍透明）。步1三核：HEAD即授权commit、四件引文现读核齐、起跑时跟踪面零改动（收尾出现他腿frontend 5件+panel 2件+票215在飞件，A/B腿写面零重叠，commit pathspec严格点名零吞并）。步2裁定**甲=页面侧圆盘**（判据非感觉）：z律（orb_windows.go follow文档+:398-403 ball窗顶层、orb窗其下）⇒乙（D2D画盘）画在**粒子上面**结构性不成立；页面是唯一粒子下方的层。几何：盘直径96%窗边（角探针对角66.5%边长在48%半径外⇒四角仍读桌面）、粒子最大投影37.1%边长（相机35°/z=4.1/sceneScale0.82/最宽环半径1.17推导74.2%半边）全被盘盖住⇒AdditiveBlending青粒子不再与桌面合成；盘色#101a2e→#0a1322→#050913深底族非纯黑。步3落地=模板#jarvis-disc+重建jarvis-orb.html（761,765B纯LF外链0）+sha钉换2cb15fc2…（基线重建先证逐字节=已提交66eb4b80…⇒delta恰为盘）；siri钉d93da649…未动。步4判据（白地板=机主缺陷形状的极端）：静态2枚（盘元素/圆形/#050913盘内/无平底负牙/belt与盘共存serving后）+winlive 3牙=①盘环8点42%边长半径均值亮度17.9<60（盘真的在）②四角(255,255,255)逐通道零偏差（透明未破）③cyan主导7.08/7.49%>2%（与既有7.08-8.32%同族）；突变④摘盘→盘环255.0三发全红FAIL 94.41s→备份逐字节cmp还原（0f497d04/2cb15fc2）→复绿PASS 13.08s。步5门禁：gofumpt全目录空/vet rc0/build三目标rc0/默认套件ok 0.223s/★winlive TestBallLiveOrb **10/10 PASS 160.176s**（旧8+新1，四角(51,51,50)钉未破）/d22scan clean（963+239）/LF五件全零。球放大如实报：44–72冻结契约（PLAN.md:3557+SPEC-08+config validate.go:62+tokens.go:524-528）且config/cmd/balldebug不在写面⇒未硬改，待机主批契约文字。AC勾与Status未动；⛔push；证据docs/evidence/s1/65-jarvis-dark-disc.md

- [2026-11-07 23:4x +08] agent=65-appearance-picker-ui (wave77 B leg) did=设置页「球风格」切换 UI 落地（机主原话入口：外观可配置切换风格，默认玻璃球，切换 Jarvis 才换）。前提=A1622（frontend 归本编队写面）。落地六件：panel.ts 契约声明面（BallStyleSection 照 ball_style_65.go json tag 逐字：style/vocabulary；并收全 ModelSwitch/Commands/Run/Files/Plugins 五组同款键——两枚对账尺当场收绿）；settings-controls.tsx 新 BallStylePicker（照本文件 Picker 同款门禁形状，值枚举全来自快照零字面量，无读口渲染诚实句）；config-screen.tsx APPEARANCE_ROWS 加 ball.style 行（blocked=""，control:"ballStyle"，徽章「能改」——档位由收据实测句说了算，有球 now 无球 restart，不预设）；App.tsx+settings-page.tsx 透传 snapshot.ballStyle；新渲染仪器 render-config-screen（行上屏/picker 挂载/当前值到按钮/无读口诚实态/写路钉 field:"ball_style"/涂出 markup 零 emoji 六组断言）。写路=performConfigWrite({field:"ball_style",value})→panel_config_store.go:284→SetBallStyle，真写（写入通道现跑核实存在）；显示值钉快照不造本地副本（config.set 不推快照，现跑核过）。门禁：typecheck(tsc -b) rc0/对账尺两枚 ok（基线红[commands modelSwitch ballStyle run]+[files]→中途+plugins→全绿）/go test ./internal/panel 全包剩 1 红=TestC21DesignTokensFourWayAgree（已提交态 tokens 面，红因 GlassRing/GlassHi/GlassHiSoft 在 ball/tokens.go 缺，非本腿）/vet rc0/d22scan clean（frontend 116 文件）/lint 0 errors/六件 LF 全核。突变牙：摘 BallStyleSection 声明 typecheck 当场红（TS2304+TS7006 一串）还原复绿（sed 多伤两处手工补回，如实记非 sha256 一发到位）。票 215 plugins 键声明并入本腿（其载体注释点名的 UI leg=本写面，其 Go 件不在本腿 pathspec）。AC勾与Status未动；⛔push；证据 docs/evidence/s1/65-appearance-picker-ui.md

- [2026-11-08 10:0x +08] agent=65-jarvis-pure-particles did=机主原话推翻深盘裁定（「不要有这种球体……直接就是粒子态」逐字最高优先）⇒删盘+换混合模式。①基线先证：改前重建逐字节=已提交资产（sha 2cb15fc2…，cmp 过）⇒delta 可信。②删 #jarvis-disc 整段（CSS+元素+注释），body transparent 构造不动。③混合模式正解=不加任何背景，tools/jarvis-export/vite.config.mjs 新增 wisp-jarvis-normal-blending 构建期插件把上游 10 处 blending: AdditiveBlending→NormalBlending（硬护栏=恰好10处否则fail build；上游克隆零改动，干净机器复跑路径不变，patch-package 同思路零新依赖）。④重建资产 759,964B 双跑确定性 sha 754037b9…，同commit换钉 orb_test.go（siri 钉 d93da649…未动）；静态尺 jarvis-disc 0/#050913 0/外链 0/CR=0。⑤判据翻转+测试件重命名如实报：orb_jarvis_disc_test.go→orb_jarvis_pure_particles_test.go、orb_jarvis_disc_live_windows_test.go→orb_jarvis_pure_particles_live_windows_test.go（git rm+新建）；winlive 新牙 TestBallLiveOrbJarvisPureParticlesOverWhiteAndBlackFloors 四牙=①42%环非粒子均值亮度 255.0>200（无盘；盘在时同环 17.9）②四角白/黑双底 ±12 零偏差 ③白底青 12.59%>2% ④黑底青 3.33%>2%（双底=像 Jarvis 任何背景都是粒子态）。★突变④如实：additive 重建重跑牙不红（白 13.63%/黑 4.89% 与 normal 12.28%/3.54% 同区间）——非重叠粒子两种模式缓冲逐位相同，平底色相尺原理上分不开混合模式；四牙咬的是盘回归/整窗背景/粒子丢失，不咬混合常量，NormalBlending 落位依据=机主裁定+three.js 默认实践，实测数已记测试件头注与 vite.config 注释。门禁：gofumpt 空/vet rc0/build 三目标 rc0/默认套件 ok/★winlive TestBallLiveOrb 10/10 PASS 165.313s/d22scan clean（963+243）/LF 六件 CR=0。残留 wisp.exe pid 41396 挡测量按牙自带指令 taskkill 后重跑。AC勾与Status未动；⛔push；证据 docs/evidence/s1/65-jarvis-pure-particles.md
