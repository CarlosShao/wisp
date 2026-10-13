# 313-r1 → 接续腿：`cmd/wisp` 精确落点包（文件＋行＋要加什么）

判语：本腿 **⛔ 未动 `cmd/wisp/` 一字**（派单补记 §1：该面此刻归 `294-r1` 的整包门禁）；下面是接续腿照抄即可的两处落点，行号＝HEAD `39f6341d` 现量（`panel_config_store.go` 最后一次被改＝`0d87a681`，票 248）。
尺：`git log --oneline -1 -- cmd/wisp/panel_config_store.go` · `grep -n 'ReadSettings\|switch w.Field' cmd/wisp/panel_config_store.go`

## 落点 1 — `ReadSettings`：把 313 那几枚读字段填上（现量 `:88-133`）

**`cmd/wisp/panel_config_store.go:96`** 这一句（现量逐字）：

```go
view := panel.SettingsView{Readable: true, ChatModel: cfg.LLM.Roles.Chat.Model, ChatModelKnown: strings.TrimSpace(cfg.LLM.Roles.Chat.Model) != ""}
```

换成：

```go
view := panel.SettingsView{
	Readable: true,
	ChatModel: cfg.LLM.Roles.Chat.Model,
	ChatModelKnown: strings.TrimSpace(cfg.LLM.Roles.Chat.Model) != "",
	// 票 313：选择器的数据源 ＋ "用户配过没有" 那一枚旗。
	Presets: config.PresetNames(), PresetsKnown: true,
	ChatProvider: cfg.LLM.Roles.Chat.Provider,
	Configured: len(cfg.LLM.Providers) > 0 || strings.TrimSpace(cfg.LLM.Roles.Chat.Provider) != "",
	ConfiguredKnown: true,
	Voice: panel.VoiceSettingsView{
		CloudASRChain:    cfg.Voice.CloudASRChain,
		CloudTTSChain:    cfg.Voice.CloudTTSChain,
		RealtimeEnabled:  cfg.Voice.Realtime.Enabled,
		RealtimeProvider: cfg.Voice.Realtime.Provider,
		RealtimeModel:    cfg.Voice.Realtime.Model,
	},
}
```

**`cmd/wisp/panel_config_store.go:112-114`** 这一小段（现量是 `for range p.Models { row.ModelCount++ }`），换成（计数与**名字清单**同源，别各算各的）：

```go
	modelNames := make([]string, 0, len(p.Models))
	for id := range p.Models {
		modelNames = append(modelNames, id)
	}
	sort.Strings(modelNames)
	row.ModelCount = len(modelNames)
	row.Models = modelNames
```

- `sort` 该文件已 import（现量 `:39`）；`config` 也已 import（`:43`）。
- 不可读那一支（`:94`）**不用动**：零值即"腿没答"，`presetsKnown:false`，页面按状态渲染。
- `PresetNames()` 含 `stepfun`（`internal/config/defaults.go:34`）——它**只给已声明行补 protocol/base_url，⛔ 不建行**（`loader.go:202-218`）；"默认 stepfun" 由页面高亮实现（`A877` ③），**⛔ 不写首建文件**（票 198 钉）。

## 落点 2 — `ApplySetting` 的 switch：六枚新 case（现量 `:182-223`）

在 `case panel.FieldProviderCredential:`（现量 `:218`）**之前**追加：

```go
	case panel.FieldVoiceCloudASRChain:
		written, err = s.mgr.SetVoiceCloudASRChain(splitChainValue(w.Value))
	case panel.FieldVoiceCloudTTSChain:
		written, err = s.mgr.SetVoiceCloudTTSChain(splitChainValue(w.Value))
	case panel.FieldVoiceRealtimeProvider:
		written, err = s.mgr.SetVoiceRealtimeProvider(w.Value)
	case panel.FieldVoiceRealtimeModel:
		written, err = s.mgr.SetVoiceRealtimeModel(w.Value)
	case panel.FieldProviderRowCreate:
		written, err = s.mgr.CreateProviderRow(w.Provider, w.Model)
	case panel.FieldRoleChatPair:
		written, err = s.mgr.SetRoleChatPair(w.Provider, w.Model)
```

同文件尾（或 `newConfigStore` 旁）加一枚小助手：

```go
// splitChainValue turns the panel's comma-separated chain value into the element
// slice config's setter takes. The panel's own gate refuses an empty value for
// these fields, so this never has to decide what "" means.
func splitChainValue(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
```

- 六枚 setter 本腿已落 `internal/config/settings_313.go`（`CreateProviderRow` / `SetRoleChatPair` / `SetVoiceCloudASRChain` / `SetVoiceCloudTTSChain` / `SetVoiceRealtimeProvider` / `SetVoiceRealtimeModel`）。
- `res.Tier = panel.EffectiveRestart`（现量 `:228`）**不用改**：现口径对每枚设置写入都答"要重启"，理由逐字在文件头（`:25-32`）——语音那几枚同理由成立。
- `credentialValue` 那条路（`:53`、`:59-66`）**一字不动**：key 录入仍只走 `provider_credential` + `credentialValue`。

## 接续腿还要顺手做的一件事

`cmd/wisp/panel_config_248_test.go` 已把甜点铺好（`settingsSeed` 有 deepseek 行），可照它的样式给"空配置 → 建 stepfun 行 → 写云链 → 写 pair"补一发端到端用例；⛔ 一处**既有断言不许改**（本腿已量：`internal/panel/config_route_248_test.go:271-277` 把 `WritableFields()` 钉在票 248 七枚上，故本腿把 313 六枚放在**第二张表**里，`WritableFields()` 一字未动）。
