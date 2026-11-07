# 票 187 收尾腿：`ApplySetting` 的 `role_chat_thinking` 末行 case

- 腿：Wisp 编队 187 欠账接续腿（上一腿 `007de34f0` 自报的欠账）
- 日期：2026-11-07 21:1x +08
- 锚：接续时 HEAD `6e16b287c`
- 写面：`cmd/wisp/panel_config_store.go`（唯一产码件）＋ `cmd/wisp/panel_config_187_test.go`（新）＋ 本件 ＋ 票 187 Progress log 一行

## 1. 交付内容

`panel_config_store.go` 的 `ApplySetting` switch 补一条 case，照既有 `FieldRoleChatModel` 同款形状（值＝枚举词本身，词表由 config 侧 setter 把关，320 裁定的同一条形状）：

```go
case panel.FieldRoleChatThinking:
    written, err = s.mgr.SetRoleChatThinking(w.Value)
```

落位在 `FieldRoleChatModel` 与 `FieldVoiceCloudASRChain` 两 case 之间。此前该字段落 default 分支被明确拒（「没有对应的写入腿」）——不静默，但页面在册的字段永远写不进去。

## 2. 起手三核读数（现跑）

| 核 | 命令 | 读数 |
|---|---|---|
| ① switch 形状 | `grep -n ApplySetting cmd/wisp/panel_config_store.go`＋`sed -n 195,310p` | 15 个既有 case 照抄形状成立；default 分支拒文原样 |
| ② setter 签名 | `grep -n SetRoleChatThinking internal/config/settings_187.go` | `func (m *Manager) SetRoleChatThinking(level string) ([]string, error)`；空值/出词表/模型声明档三道拒，`writeOneKey` 守卫路径 |
| ③ ModelSwitch 接线 | `grep -rn ModelSwitch cmd/wisp/*.go` | **cmd/wisp 零命中**；读口在 `internal/panel/pump.go:231/:432`。⇒ **快照 `modelSwitch` 键的装配根注入未接，归 145 腿的 `cmd/wisp/run.go` 写面，本腿未碰**（如实报） |

`git status --porcelain cmd/wisp`：本腿起手时仅两枚 323 未跟踪测试件（非本腿写面，未动）；`panel_config_store.go` 干净无在飞 hunk。

## 3. 牙（三发）

1. **写进去**：`Test187ApplySettingThinkingLandsTheKey` —— 真 Manager＋真守卫 setter，收据 `Written == [llm.roles.chat.thinking_intensity]`、store 层 tier＝restart（面板门按 `fieldTier187` 改 next_task，这是面板侧既有语义，本腿只钉 store 戳）、内存 `LLM.Roles.Chat.ThinkingIntensity == "low"`、盘上文件按行找到 `thinking_intensity ... low`（写入器规范化整文件引号，故按行断言不按字节）。
2. **突变**：摘掉 case → `go test -run Test187ApplySetting` 当场红：
   `ApplySetting(role_chat_thinking=low): 字段 "role_chat_thinking" 没有对应的写入腿，未落盘` —— case 承重证实。还原后 `git diff --stat`＝6 行插入，逐行同款。
3. **拒绝路径**：`Test187ApplySettingThinkingRefusesOutsideVocabulary` —— 值 `maximum` 出四值词表（off|low|medium|high），setter 落盘前拒；收据 `Tier == EffectiveNotApplied`，内存与盘上文件零改动。拒因句带键路径 `llm.roles.chat.thinking_intensity`（setter 自有词表句），断言钉键路径。

## 4. 门禁读数（2026-11-07 21:2x，本机现跑）

| 门禁 | 读数 |
|---|---|
| `gofumpt -l`（点名两件） | 空，rc0 |
| `go vet ./cmd/wisp` | rc0 |
| `go build ./cmd/wisp` | rc0 |
| `go test ./cmd/wisp -run 'Test187\|TestConfig' -count=1 -timeout 30m` | `ok github.com/CarlosShao/wisp/cmd/wisp 0.156s`（全绿；含既有 TestConfig* 65/313/320 族） |
| `tools/d22scan/d22scan.exe` | clean，449 生产 Go 文件，rc0 |
| 行尾 | 两件均 LF（`file` 无 CRLF 标记） |

## 5. 没做到的（如实）

- **快照 `modelSwitch` 键在活机不出**：`PumpSources.ModelSwitch` 的注入点在 `cmd/wisp/run.go`，那是 145 腿的写面，本腿未接（三核③读数见上）。⇒ 面板快照仍无 `modelSwitch` 键（reader-gated 缺席，213 同款），直到 145 接线。
- 未跑 `internal/config`、`internal/panel` 任何测试（非本腿写面，187 腿已各报）。
