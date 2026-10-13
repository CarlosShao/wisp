# A. config.get 今天真读到什么（只读普查腿 313-a1；HEAD 0e5253fe，2026-10-13）

每行=结论＋尺＋件行号快照（行号=HEAD 0e5253fe 现读；引用前重跑尺）。

## A0 逐跳链（现读点名，读通）
- 入向门：页面 `window.chrome.webview.postMessage` → host Init 钩子折进 `window.wispDispatch` → bind 闭包，**返回回执串**。（cmd/wisp/panel_host_windows.go:813-850）
- 方法名：`config.get` / `config.set`，名册在 bridge.go:146-152（internal/panel/bridge.go:66-67 定义；名称出处=SPEC-08 已有拼法，bridge.go:46-65）。
- 路由：internal/panel/composer_dispatch.go:199-208（两方法同进 ConfigRequestHandler.HandleConfigRequest）→ 读 handler = internal/panel/config_handlers.go:269-303。
- 腿：cmd/wisp/panel_config_store.go:88-133（configStore.ReadSettings）；装配点=cmd/wisp/panel_inbound.go:267-271 与 cmd/wisp/panel_resident_windows.go:215（同一 newComposerDispatchChain）。
- 真源：internal/config/manager.go:115-119（Manager.Config() 深拷贝快照）→ `cfg.LLM.Providers`（internal/config/schema.go:379-401）。
- 尺：`grep -n 'ReadSettings\|cfg.LLM.Providers' cmd/wisp/panel_config_store.go`

## A1 逐字段：页面今天真能拿到什么
页面拿到的是**一句散文**（renderSettingsView，internal/panel/config_handlers.go:445-462），不是 JSON：
- 可读性两态：可读→"配置可读。服务商 %d 家；凭据状态 %s。"；不可读→"配置当前不可读，所以无法回答是否已录入密钥（这不是未配置）。"（:446-451）
- 聚合凭据态 CredentialState 五值：unknown / config_unreadable / no_ref_declared / all_recorded / partly_missing（:180-191）
- 聊天模型：`cfg.LLM.Roles.Chat.Model` 逐字（cmd/wisp/panel_config_store.go:96；known=false 时句="当前聊天模型未读到。"）
- 每 provider 一行："- 名：引用 <ref>，目录模型 N 个，密钥 <词>"（config_handlers.go:457-460）；词=不适用/无法核对:<error>/已录入/未录入（recordedWord :471-482）
- 凭据值：**绝无**——本文件注释逐字"它不携带任何凭据值、也不携带半脱敏值——连后四字符都没有"（:138-141）；引用名是 dpapi:<id>/env:NAME（随机 id，不带 vendor 信息，:158-162）
- 注：SettingsView/ProviderSettings 结构字段（:142-172）=Readable / Credential / Providers{Name,APIKeyRef,HasBaseURL,ModelCount,CredentialSet,CredentialError} / ChatModel / ChatModelKnown；其中 **HasBaseURL 填了但渲染句不带它**（panel_config_store.go:110 填充；:457-459 只印引用/计数/密钥）⇒ 页面连"这行有没有 base_url"都读不到
- 尺：`sed -n '445,462p' internal/panel/config_handlers.go`

## A2 必答①：provider 行来自哪张表 ⇒ stepfun 在不在今天的行里
- 结论：来自**用户 config.toml 已声明的行**（`for name := range cfg.LLM.Providers`，panel_config_store.go:98-102）。预设表不是行来源：applyPresets 只给已存在的行补 protocol/base_url、**不建行**（internal/config/loader.go:199-218）；NewDefaults 对 map 恒留 nil（internal/config/defaults.go:55-62）。
- ⇒ 全新首建配置（cmd/wisp/firstrun.go:74-97 写 NewDefaults）里**零行**；stepfun 今天出现在行里的唯一前提=用户手写了 `[llm.providers.stepfun]`。
- 尺：`grep -n 'applyPresets\|LookupPreset' internal/config/loader.go` ＋ `sed -n '98,102p' cmd/wisp/panel_config_store.go`

## A3 必答②：页面那句 waiting"面板没有一个字段能列出「有哪些模型可选」"判真伪
- 判：**真**（就 config.get 这条口）。今天"模型清单"=每 provider 的 `llm.providers.<名>.models` 键集合（schema.go:399-400）；出行只有计数：`row.ModelCount`（panel_config_store.go:112-114）→ 句"目录模型 N 个"。模型**名字一枚都到不了页面**。
- 另有无清单口：没有。快照仅有 `composer.currentModel/modelKnown`——“这一个正在用的 id”（internal/panel/composer.go:247-256；pump 填充 pump.go:280-286），不是可选清单；且新前端树未消费（见 D）。
- 尺：`grep -n 'ModelCount' cmd/wisp/panel_config_store.go internal/panel/config_handlers.go`

## A4 附：今天这条口的四条硬性质（供落地腿引用）
1. 白名单先于腿：未列名字在 panel 内被拒（config_handlers.go:321-326），腿再拒一次（panel_config_store.go:221-222）。
2. 读不出=有状态的句子，不是空值（:299-301；:447-449）。
3. 写入的"生效级别"必须说（tierSentence :428-441）；现量全=重启（panel_config_store.go:228）。
4. 全程无一字节凭据值过此路由（:18-25 头注 + :339-347 credential 分支只回状态）。
- 尺：`sed -n '307,365p' internal/panel/config_handlers.go`
