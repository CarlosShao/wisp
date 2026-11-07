# 票 65 · 设置页「球风格」切换 UI（外观选择器腿）

- 腿：wave77 B 腿（frontend/** 本波独占写面）
- 日期：2026-11-07 22:1x–23:4x +08
- 机主原话（逐字）：「两个都要，但是设置项里面外观可以配置切换风格，默认就是玻璃球，切换 Jarvis 才换球的风格」
- 授权前提：机主 2026-11-07 22:0x 令（A1622 逐字入 7873c8328 commit message）——frontend/** 与全部 UI 归本编队写面；旧「frontend 归外部 agent Q-51 零写面」记载作废。

## 0. 起手三核（现跑读数）

1. **派单引文核**：`frontend/src/lib/panel.ts` 现量确认 `PanelSnapshot` 只声明 `tasks?`，`ballStyle/modelSwitch/commands` 三键缺——引文为真。但基线红名单比派单多：`TestComposerContractTypesMatchFrontend` 还报 `ComposerState` 缺 `files`（票 190），且随后现量长出 `plugins`（票 215 在飞腿）——全部并入本腿声明面（见 §3）。
2. **写入通道核**：派单问「前端有没有 config 写入通道」——**有**：`panel.ts requestConfigWrite`（config.set，configField/configValue）→ `lib/config-actions.ts performConfigWrite` → `cmd/wisp/panel_config_store.go:284 case panel.FieldBallStyle → mgr.SetBallStyle`。派单的「若没有则只落读半」分支不触发。
3. **红/绿归属核**：`git status --porcelain` 现跑——本腿开场时 tracked 面干净（A 腿 internal/ball 零 tracked 改动）；C21 三方红的两件被点名文件对 HEAD 零 diff ⇒ 已提交态红，非在飞件污染。本腿中途 HEAD 前移（644ddc147，票 226 收口）且票 215 腿在 `internal/panel` 落了 `plugins` 键（composer.go/pump.go 已改未提交）——如实记，处理见 §3。

## 1. 落了什么（逐件）

| 件 | 改动 |
| --- | --- |
| `frontend/src/lib/panel.ts` | 契约声明面：`CommandsSection`/`ModelSwitchSection`/`BallStyleSection`/`FilesSection`/`RunSection`/`PluginsSection`（＋行类型 `CommandRowView`/`ModelOptionView`/`TreeEntryView`/`ToolRunRowView`/`RunErrorRowView`/`RunStuckRowView`/`PluginRowView`）全部照 Go json tag 逐字对齐（`ball_style_65.go`/`model_switch_187.go`/`commands_roster_213.go`/`tree_types_190.go`/`run_feed_145.go`/`plugins_inventory_215.go`）；`PanelSnapshot` 补 `commands/modelSwitch/ballStyle/run/plugins` 五键、`ComposerState` 补 `files` 一键 |
| `frontend/src/components/settings-controls.tsx` | 新控件 `BallStylePicker`：照本文件 `Picker` 同款门禁形状（button+popover+useDismiss+Check），零自创范式；值枚举与当前值全部来自快照 section，本文件零风格名字面量；section 为 null（本机没装配球风格读口）时渲染诚实句，不画空壳 |
| `frontend/src/components/config-screen.tsx` | `APPEARANCE_ROWS` 加 `{group:"ball", name:"style", desc:"悬浮球外观风格（玻璃 / Jarvis 粒子）", blocked:"", control:"ballStyle"}`；`Row.control` 联合类型加 `"ballStyle"`；`ConfigScreen` 加 `ballStyle` prop（快照读半透传）；`moveBallStyle` 走 `performConfigWrite({field:"ball_style", value})` 真写（空值本地先拒，收据句落 `FeedbackLine`）；徽章词：ball.style 行标「能改」（生效档位由收据实测句说了算——有球 now、无球 restart，标「即时」会在 restart 档说谎），alpha/motion 两行保持「即时」原文不动；文件头 dormant/ live 分界注释重写如实；页脚过期文案（「缺的是 C17 快照上的一个字段」）改为如实句 |
| `frontend/src/components/harness/settings-page.tsx` | `SettingsPage` 加 `ballStyle` prop 透传给外观类的 `ConfigScreen`；「悬浮球」分类的 waiting 文案如实更新（风格已能切＋去哪切；大小/透明度/穿透仍 D29 原生；唤醒词/热键仍无数据源） |
| `frontend/src/App.tsx` | `SettingsPage` 调用点接 `ballStyle={snapshot.ballStyle ?? null}`（快照是唯一数据源，无本地副本，PLAN.md:1044） |
| `frontend/scripts/render-config-screen.tsx` | 新渲染仪器（照 render-* 同款范式）：6 组断言——ball.style 行上屏（chip 两半拼接）、能改徽章＋picker 挂载＋快照当前值到达按钮、即时徽章计数钉 3、无读口诚实态（无 picker＋那句话）、写路源级钉（`field: "ball_style"`＋只走 performConfigWrite 无 postMessage）、涂出 markup 零 emoji（ban#8 乐器同款区间） |

## 2. 显示值跟手方式（诚实设计，实测后定的）

`config.set` 处理链（composer_dispatch→config_handlers→configStore.ApplySetting）**只回收据句，不推新快照**（现跑核过：publishPanelSnapshot 的触发面＝审批卡/事件流，config.set 不在其列）。因此选择器**不造本地副本**：显示值钉在快照 section 上，写入后的下一帧快照把它带上来；写入事实由收据句当场呈现（「已保存…立即生效／重启后生效」——`renderSettingReceipt`+`tierSentence` 原文，`translateWriteReply` 现有规则能识别「已写入」，无需新增翻译）。这与 model/voice 设置页「write 后 reload 重读」是同一诚实家族的两种形状：它们重读的是 config.get 视图，本行重读的是快照推送——各自的读口是什么就信什么。

## 3. 票 215 的 plugins 键为何在本腿声明（归属写明）

本腿跑对账尺时红名单从 `[commands modelSwitch ballStyle run]` 变为 `[plugins]`：票 215 腿在 `internal/panel` 落了 reader-gated `plugins` 键（composer.go:174 / plugins_inventory_215.go，未提交），其载体注释逐字写明「the two two-way reconciliations against frontend/src/lib/panel.ts name it as undeclared until the UI leg moves them」。本编队 frontend/** 现在就是 UI 腿写面（A1622），故照载体逐字落 `PluginsSection`/`PluginRowView` 声明——两枚对账尺当场收绿。该腿的 Go 件（composer.go/pump.go/plugins_inventory_215*.go）**不在本腿 commit pathspec 里**，其提交归其自己。

## 4. 判据与门禁读数（全部现跑）

| 判据 | 读数 |
| --- | --- |
| ① TS 类型检查 | `npm run typecheck`（＝`tsc -b`，如实报跑的这条）rc0 |
| ② Go 契约红转绿 | `go test ./internal/panel -run 'TestComposerContractTypesMatchFrontend\|TestApprovalCardViewJSONKeysMatchFrontendTypes' -count=1` **ok rc0**（基线红 `Snapshot emits [commands modelSwitch ballStyle run]`＋`ComposerState emits [files]`→中途 `+plugins`→全收绿）；全包 `go test ./internal/panel -count=1` 剩 **1 红＝TestC21DesignTokensFourWayAgree**（红因：c21-native-tokens.md 点名 `Palette.GlassRing/GlassHi/GlassHiSoft` 在 internal/ball/tokens.go 缺失；两件对 HEAD 零 diff、internal/ball 无在飞改动 ⇒ 已提交态红，tokens 面归 A 腿/token 治理，非本腿） |
| ③ 设置页渲染 | 新仪器 `render-config-screen` OK（ball.style 行上屏＋picker 挂载＋glass 到按钮＋无读口诚实句＋写路钉＋零 emoji）；render:nav/composer/stream OK |
| ④ 突变 | 摘 `BallStyleSection` interface ⇒ typecheck 当场红（TS2304 Cannot find name ＋ TS7006 implicit any 一串），还原复绿（sed 范围删除多伤了 interface 本体与注释两处，手工补回后 typecheck/仪器/对账尺全绿——如实记，还原不是 sha256 一发到位） |
| go vet | `go vet ./internal/panel` rc0 |
| d22scan | `sh scripts/d22scan.sh` **clean rc0**（frontend/=116 文件过 ban#6/ban#8；UI 零 emoji） |
| lint | oxlint 0 errors / 35 warnings（其中本腿新仪器 1 枚已修——FE0F 移出字符类改 alternation；App.tsx:147 的 set-state-in-effect 为既存代码，非本腿引入） |
| LF | 六件改动件 `git ls-files --eol` 全 `w/lf`（settings-page.tsx/App.tsx 曾被工具写成 CRLF，已转回） |

## 5. 边界与如实

- **没动的**：`internal/**`、`cmd/wisp/**`、design/** 零字节（A 腿 internal/ball＋票 215 腿 internal/panel 各自写面，本腿零触碰）；thresholds/golden 零触碰；票面 AC 勾与 Status 零触碰；三张票只 append Progress log。
- **写入通道**：真写已通（§0.2），非「只有读」——派单的降级分支未触发。
- **收据档位**：ball_style 写入是实测档（`panel_config_store.go:300-309`：applier 在＝now，无球＝restart），UI 不预设档位词，收据句原文呈现。
- **选择器不做的事**：不造第二写入通道（C17 白名单零新增）；不本地镜像风格值；无读口时不画可点的空壳。
- **临时件**：无新增（dump 用过 /tmp，不在仓库内）。
- 只 commit 不 push；`git add`/`commit` 显式 pathspec 点名具体文件。
