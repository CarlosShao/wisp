# D. 页面侧现状与待办＋313 三格 AC 判据草案（只读普查腿 313-a1；frontend=dev 工作树；HEAD 0e5253fe）

每行=结论＋尺＋件行号快照；行号=HEAD 现读，引用前重跑尺。

## D1 现状逐件（现读）
- 设置页 frontend/src/components/harness/settings-page.tsx：rail 12 类（CATEGORIES :37-102）。「模型」waiting 逐字（:51）="还没收到模型清单。C5 抽象了 provider，但面板没有一个字段能列出「有哪些模型可选」。"；「语音」:54-56（waiting 句=:56"还没收到语音链路的状态。…"）；「外观」waiting=""（:46）⇒ 只有外观渲染真屏（:188-193）。
- 真屏 frontend/src/components/config-screen.tsx：**活着的两行**=APPEARANCE_ROWS :87 `panel.opacity`（control alpha，"面板背景不透明度"）、:88 `panel.animations`（control motion，"界面动画（关掉后全部静帧，性能最省）"）；其余 6 行 blocked（:86、:89-93）。两行只写本页样式（:129-137、:207-210）——不出 WebView、不落盘、重启回 token 默认。
- 声明面 frontend/src/lib/panel.ts（唯一出口）：sendRequest=:211-226；活调用点四枚 composer 方法——panel.mode.request(:246)/panel.workspace.request(:255)/panel.attachment.add(:272)/panel.message.send(:300)；另 `panel.approval.request`(:181) 是旧门、Go 名册**拒**（bridge.go:146-152 不含它）；panel.view.request 只在注释(:228)。**config.get/config.set 零调用**（现跑全树 0 命中；node_modules 内命中非页面）。
- 数据通路：出向=postMessage→wispDispatch→bind 闭包**返回回执串**（cmd/wisp/panel_host_windows.go:813-850）⇒ config.get/set 回执今天就走这个 promise 返回；入向（快照推页面）=**新树无接收端**（`grep chrome.webview\|addEventListener("message")` 0 命中；main.tsx:12 注释说"Go pushes snapshots"，实现不在本树）⇒ App 吃 EMPTY 快照（App.tsx:91-92）。快照线上本有 composer.currentModel/modelKnown/credentialState/credentialKnown（internal/panel/composer.go:247-269），前端类型与消费方都没有它们（panel.ts:112-119）。
- 首启引导：产品 App **没挂** FirstRunScreen（grep App.tsx 0 命中；它只在 showcase/harness）⇒ AC#1"（与首次引导，若有）"的落点=设置页。
- 尺：`grep -rn 'config\.get\|config\.set' frontend/src --include='*.ts' --include='*.tsx'`（=0）＋ `grep -n 'sendRequest(\|approval.request' frontend/src/lib/panel.ts`

## D2 页面侧待办（15 行，可直接整段转交机主外部 agent）
1. 「模型」类（settings-page.tsx:48-52）：waiting 文案替换为控件区；⛔ 不动其它 11 类与 rail 结构。
2. 读面调用形：`const receipt = await window.wispDispatch(JSON.stringify({method:"config.get", requestId:"pc-…", source:"panel-composer"}))`——回的是**一句人话**（a-get A1），下拉要的结构化数据今天没有 ⇒ 落地腿需在回句之外给结构化读字段（⛔ 页面先别解析 prose）。
3. 服务商下拉数据源=预设清单（internal/config PresetNames()，含 stepfun；defaults.go:39-46）——今天零推送 ⇒ 新增读字段（与票 312 扫描预填共用一口）。
4. 模型下拉数据源=该 provider 的 `llm.providers.<名>.models` 键名——今天只有"目录模型 N 个"计数（config_handlers.go:458）⇒ 新增读字段；无行时给"手填模型 id"入口。
5. 文本写路径：`config.set` + `configField:"role_chat_model"` + `configValue:<模型id>`（settings.go:224-230）；provider 那一半今天**没有写口**（requireChatProvider 拒，:235-243）⇒ 成对写要落地腿补字段或建行口。
6. 语音写路径：新字段（名册/校验见 b-set B2/B3）；实时段两候选形（α 单字段 pair / β 双字段）由落地腿裁定。
7. 密钥录入：`config.set` + `configField:"provider_credential"` + `configProvider:<名>`；值**只能**放 `credentialValue` 键（cmd/wisp/panel_config_store.go:53、:59-66）——⛔ 放 configValue、⛔ 放共享封套其它键（消息路可见）。
8. 密钥永不回显写法：输入框提交后立即清空且不存 state；屏上任何位置不渲染值/尾号；状态只显示「已录入/未录入」+引用名（dpapi:<id> 或 env:NAME）；回执句原文（renderCredentialReceipt，config_handlers.go:401-424）可直接展示。
9. 回执即事实：写入回执带"实际落盘的键路径+生效级别"（renderSettingReceipt :381-396）；现量所有设置写入=**要重启**（panel_config_store.go:228 EffectiveRestart）⇒ 照显示，⛔ 自编"已生效"。
10. 成对写拦截：保存前拦 provider/model 半对（catalog.go:101-106 set both or neither）；拒因三句照原文（settings.go:77-81 三种拒因文案）。
11. 行存在性：页面⛔ 不能建 provider/model 行（settings.go:36-39、:376-380）；今天让用户照首启回执手加三样——除非本票新增建行能力（b-set B6）。
12. 快照读法（可选）：若要用 composer.currentModel/credentialState，先要一条入向承载（本树今天没有监听实现，D1）——那是另一枚票的口，⛔ 本票页面别自建第二通道。
13. 规范：⛔ emoji；用现有 token/控件（Chip/SegmentedControl 等）；⛔ 自创新方法名（C17 白名单，panel.ts:228-235）。
14. 读不到就说"未读到"（a-get A1 两句先例）；⛔ 空壳控件假值（config-screen.tsx:21-27 的 P9 红线同规）。
15. stepfun 默认高亮=仅"零已配置行"时生效；用户有配置就按配置（AC#1"改过以用户为准"）。

## D3 313 AC 判据草案（每格：判据形状＋落点＋今天能否自动判）
- AC#1（默认选中）：
  - Go 侧今天可判的一半=预设表含 stepfun（internal/llm/catalog_test.go:151 已钉）；若落"首建写入"则 firstrun_198 测试族可判，但先触钉（c-default C3 候选②）。
  - 必须页面判（页面归机主 agent）：无配置时模型类默认高亮 stepfun（DOM/harness 或真机）。
  - 结论：Go 侧只能判"stepfun 在预设表"；"默认选中"本身=页面/真机判。
- AC#2（两个模型界面可配）：
  - Go 侧可自测：(a) 新字段进白名单（WritableFields()，internal/panel/config_route_248_test.go 族）；(b) 每枚新 setter 校验（internal/config settings_248_test.go / validate 族：成对、未知 provider/model、链元素拒 local-sherpa、锁定族仍拒）。
  - 必须页面/真机：控件渲染＋真窗里把文本+语音（stepfun audio3）配成并通过回执；还依赖"建行"口（b-set B6）；CI 不判真窗。
- AC#3（界面录 key，不回显/不进日志/不进模型文本）：
  - Go 侧可自测：值只走 credentialValue 键（cmd/wisp/panel_config_248_test.go 族已有）；config.toml 只出现 dpapi:<id>（validateAPIKeyRefs validate.go:189-196）；回执不回显（renderCredentialReceipt config_handlers.go:401-424；票 63 族 argv/日志钉）。
  - 必须页面判：DOM 从不渲染值、提交后清空（页面 agent 的测试）。
  - 必须真机：真窗走一遍（输 key→屏上只出"已录入"+引用名→config.toml 只有 dpapi:<id>）；CI 不做真窗。
- 尺：`grep -rn 'WritableFields\|credentialValue' internal/panel cmd/wisp --include='*_test.go'`
