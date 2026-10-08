# 167-a2b 格一：面板「草稿」（用户打了字还没发出去那段文本）在 Go 侧的承载面

- 现量时间 2026-10-08 18:2x +0800；根锚 HEAD `8dd239f1`；全部 `git show HEAD:` / `git grep HEAD` 现量，零 Go 命令。
- 三问：①快照里有没有承载它的字段（档位三选一）②有没有读者 ③有没有写者（入向有没有一条路能把页面那段未发的文本递给 Go）。

## 0. 判语（先说结论）

- 档位＝**〔没字段〕**：出向快照（`Snapshot`／`ComposerState`）全字段无草稿位。
- 读者＝**0**（没有出向字段可读；入向 `Text` 在产码里也没有任何读者）。
- 写者＝**插座有、接线断**：名册/白名单/封套/插座四层全在，生产装配那行逐字 `Message:    nil`（`cmd/wisp/panel_inbound.go:277`）⇒ 页面真发按名拒绝。
- 判语同 167-a2 §9 第 3 条"两层：门有、接法无"（复认成立；行号按当前 HEAD 现量更新，且**未抄旧行号**）。

## 1. 主尺与读数

| 尺 | 命令 | 读数 |
|---|---|---|
| draft 字面 | `git grep -niE 'draft' HEAD -- 'internal/panel/*.go' 'cmd/wisp/*.go' ':(exclude)*_test.go'` | **rc=1，零命中** |
| 字段名册 | `git show HEAD:internal/panel/composer.go \| sed -n '40,90p'` ＋ `sed -n '240,346p'`（文件共 346 行） | 逐枚见 §2／§3 |

## 2. 出向面（Go→页面）——没有承载草稿的字段

- `ComposerState`（`internal/panel/composer.go:235`，11 枚 JSON 键）：`mode` / `workspace` / `attachments` / `acceptedAttachmentMimes` / `maxAttachmentBytes` / `attachmentError` / `git` / `currentModel` / `modelKnown` / `credentialState` / `credentialKnown`。**无草稿键**。
- `Snapshot`（`composer.go:57`，6 枚键）：`pending` / `results` / `composer` / `generatedAt` / `instructions`（omitempty） / `tasks`（omitempty）。**无草稿键**。
- 构造器 `func NewComposerState(mode risk.Mode, ws WorkspaceView, atts []AttachmentRef, maxBytes int64)`（`composer.go:273`）：四参无文本、无失败原文。
- ⇒ 档位落 **〔没字段〕**。

## 3. 有没有读者——零

- 出向无字段可读（§2 全字段已枚举）。
- 入向封套的 `Text`（`internal/panel/bridge.go:98`，逐字 `Text string \`json:"text,omitempty"\``）在产码里的消费接口 `MessageRequestHandler`：**全仓命中仅 3 处、全在定义自身**——`internal/panel/composer_dispatch.go:87`（注释）/ `:88`（接口）/ `:130`（插座字段）⇒ **产码实现体 0 枚**。
- `OutgoingMessage`（`composer.go:295`，"一条发送"= Text+Attachments 的载体）：全仓产码命中仅定义自身（`:292` 注释 / `:295` type / `:308` `ForAgent`），唯一的构造在测试（`attachments_test.go:284`）⇒ **产码零构造**。
- ⇒ 即便封套解析收下 Text，产码里没有任何东西读它。读者＝0。

## 4. 有没有写者——插座有、接线断（逐点现量）

1. 名册：`internal/panel/bridge.go:45` 逐字 `MethodMessageSend = "panel.message.send"`（const 块 `:41` 起第 4 枚）；
2. 白名单：`bridge.go:148` `case MethodModeRequest, ..., MethodMessageSend, ...:` ⇒ `return true`；
3. 封套：`ComposerRequest.Text`（`bridge.go:98`）；
4. 插座：`ComposerDispatch.Message MessageRequestHandler`（`composer_dispatch.go:130`）；
5. 派发：`composer_dispatch.go:192` `case MethodMessageSend:` → `:193-195` `if d == nil || d.Message == nil { return "", d.unattached(req) }`（否则 `d.Message.HandleMessageRequest(ctx, req)` `:196`）；
6. **生产装配（断点）**：`cmd/wisp/panel_inbound.go:277` 逐字 `Message:    nil,`；其上注释逐字（`:274-276` 区间，sed 现量）："Untouched by this slice, and refused by name when they arrive: workspace = ticket 186, attachment = ticket 92, message = ticket 35."。

⇒ 页面若真发 `panel.message.send`，落 `unattached` 拒绝（`composer_dispatch.go:230` 函数），拒绝句逐字（`:231-232`）："方法 %q 的处理器未接入（requestId=%q），档位/工作区/附件/消息/设置均未变化"。⇒ 写者路＝**门有、接法无**。

## 5. 判不动（缺哪把尺）

- 票 AC#5 下半段"失败原文单独存住、不许覆盖用户后来新写的字"落在页面侧 `frontend/**`：本程禁读（派单硬禁区），页面今天怎么处理草稿**判不动**。
- "草稿该不该由 Go 出向承载"＝契约判定（撞票 145 局部解冻面／`Q-51` 谁写 `panel.ts`）：本程只交读数，⛔ 不裁。
- ⛔ 未提议动 `frontend/**` 任何一寸。
