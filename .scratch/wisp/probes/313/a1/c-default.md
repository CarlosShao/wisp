# C. "默认=stepfun"的落点候选（只读普查腿 313-a1；HEAD 0e5253fe）

每行=结论＋尺＋件行号快照；行号=HEAD 现读，引用前重跑尺。

## C1 roles 没配时今天的行为（逐字）
- 聊天等四 role（internal/llm/resolver.go:182-190 注释逐字）："Unset roles (both provider and model empty in config) fall back to the FIRST text_chain element with the role's temperature; when no chain exists either, a ClassConfig error names the role."
- 实现 :203-226；"都没有"的错误句逐字 :226 = `"llm: role %q is unset and text_chain is empty (Unconfigured)"`（⇒ run 腿按 Unconfigured 退 2）。
- ASR/TTS 不是 role：`voice.asr.provider` default="local-sherpa"（schema.go:210）、`voice.tts.provider` 同（:216）⇒"没配"=local-sherpa + 空 model/voice。运行时**零消费者**：internal/speech 只有 doc.go（DEFERRED: 票 15/26/41）；voice 面代码里只有 config 自己的校验与 reload 差量（catalog.go:30-38、manager.go:386-387）。
- 尺：`sed -n '182,190p;203,226p' internal/llm/resolver.go` ＋ `ls internal/speech/`

## C2 "默认选中"今天有没有承载物 ⇒ 没有（四证）
- 一：NewDefaults 全来自 `default:"..."` 标签（defaults.go:58-62；D36 rule 3），Role.Provider/Model **无 default 标签**（schema.go:285-294）；map 恒 nil（defaults.go:55-62）。
- 二：首建配置**不带**任何 provider/model/key——firstrun.go:36-38 逐字"the created file carries NO model and NO api_key_ref"；回执句 :94-97"未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺"。
- 三：票 198 的钉**逐字禁**首建文件出现 `[llm.providers` 字样（firstrun.go:149-152；测试 firstrun_198_test.go:215 forbidden 列表）。
- 四：预设表里 stepfun 在（defaults.go:34 `"stepfun": {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.stepfun.com/v1"}`；票面写 :28 ⇒ 现读 :34，行号已漂），但预设只给已声明行补 protocol/base_url（loader.go:199-218）、**不构成"选中"**。
- 尺：`sed -n '34p;39,46p' internal/config/defaults.go` ＋ `grep -n 'carries NO model\|未替你选' cmd/wisp/firstrun.go`

## C3 落点候选＋代价（"用户改过以用户为准"逐候选说）
- 候选① 界面预选（UI-only；本腿建议先落这一枚）：模型类控件在"无任何已配置行"时默认高亮 stepfun，不点保存不写任何东西；用户一旦有配置（roles.chat 或 providers 行非空）即按配置 ⇒"改过以用户为准"天然成立（只读+默认高亮，不落盘）。代价：小。缺件两个读面（预设清单 PresetNames defaults.go:39-46 今天零推送；模型清单见 a-get A3）＋写面仍受"不能建行"限（b-set B6）。判据=页面/DOM；Go 侧今天只能钉预设表（已被 internal/llm/catalog_test.go:151 钉"stepfun 在预设名里"）。
- 候选② 首建文件写入（ensureFirstRunConfig 扩展：写 [llm.providers.stepfun]+models.<id>+roles.chat 对）：代价大且**触钉触契约**——(a) 撞 firstrun_198 的钉（首建⛔ [llm.providers；放宽=违票面"⛔ 为让默认值生效去放宽断言"）；(b) 需一枚**真** stepfun model 名（C4：仓内零枚，⛔ 编名）；(c) 只对全新机器生效，存量配置不动；(d)"用户改过"由文件覆盖 tag 成立。⇒ 要落地须先拿人工裁定。
- 候选③ defaults 标签 / 运行时兜底（两枚都重）：(a) 给 Role.Provider/Model 加 default 标签 ⇒ 每份缺键配置都变 provider=stepfun/model=X，而 catalog.go:108-115 会在缺 `[llm.providers.stepfun]` 行时**拒载**；map 恒 nil（defaults.go:55-57）又不能同时 seed models 行 ⇒ 一改就是"下一次启动直接加载不过"，与 AC#7 精神正相反。(b) resolver 把"role unset+chain 空"的 Unconfigured 改成默认选 stepfun ⇒ 推翻 firstrun.go:36-38 的 R13（"没有配置就当默认跑" stays refused）与 run.go 退 2 的钉 ⇒ 同样要人工批。
- 尺：`grep -n 'stays refused\|carries NO model' cmd/wisp/firstrun.go` ＋ `grep -n 'references unknown provider' internal/config/catalog.go`

## C4 stepfun 名下今天 catalog 有哪些 model 可写（逐名抄；⛔ 编名）
- **零枚**。全仓 stepfun 非测试出现点穷举：defaults.go:20/34（预设行）、schema.go:376（注释）、internal/secret/configrefs.go:15（注释示例"dpapi:stepfun"）；测试里另有 internal/llm/catalog_test.go:151（预设名钉）与 cmd/wisp/secret_test.go 若干样例。**没有任何内置 stepfun 模型名**。
- 今天"可写"的 stepfun 模型名集合=用户 config.toml 已写的 `[llm.providers.stepfun.models.<id>]` 键（校验 catalog.go:112-115）。
- 新名字的合法来源=真端点：`wisp providers discover stepfun`（cmd/wisp/providers.go:115-149；只打印 /v1/models、⛔ 不写文件；且要求 provider 行已存在，:118-122）⇒ 313 要"界面给可选模型"缺的正是"发现→写入"这条口（今天只有手改）。
- 尺：`grep -rn 'stepfun' --include='*.go' internal/ cmd/ | grep -v _test.go`
