# 187 面板要画哪几枚控件、点下去发哪条请求（交编排者转前端会话；票 187 AC#9）

本表只描述 Go 侧已备好的能力；`frontend/**` 归界面会话（Q-51），本票零触碰。

## 1. 数据源（读，快照新键）

`PanelSnapshot.modelSwitch`（`internal/panel/model_switch_187.go`）：

| 字段 | 内容 | 页面用法 |
|---|---|---|
| `current.provider` / `current.model` | 活配置里 chat 角色此刻的名对 | 下拉默认高亮；与 `composer.currentModel`（boot 解析值）可能不同，重启前如实分别显示 |
| `currentThinking` | `llm.roles.chat.thinking_intensity` 原文；`""`＝文件没写这枚键 | 档位控件默认值 |
| `models[]` | `{provider, model}`，只含 `enabled=true` 的目录条目 | 模型下拉选项 |
| `thinkingLevels[]` | **当前模型自己声明**的档（可能少于四档、可能空） | 档位控件只画这些；空＝不画档位控件 |
| `reason` | 「没有可用模型」等显式句（目录空/全 disabled/配置不可读） | reason 非空且 models 空时显示原因，不渲染成空白菜单 |

键缺席规则（213 同款）：泵装配没接 `ModelSwitch` 读口 ⇒ 整个键不上线；页面按
`"modelSwitch" in snapshot` 判可画，勿按 null 判。

## 2. 写操作（点下去发什么）

全部走**既有** `config.set`（`internal/panel/bridge.go` 已在册的 settings 通路），
`source: "panel-composer"`，带 `requestId`。

### 换模型

```
{ method: "config.set", configField: "role_chat_model", configValue: "<model id>" }
```

- 该字段已在 248 名册；换 provider+model 一起换则用 313 名册的
  `role_chat_pair`（`configProvider`＋`configModel`）。
- 收据档位＝restart：模型通路本进程不重建，写完要重启进程才用上。

### 换思考档位

```
{ method: "config.set", configField: "role_chat_thinking", configValue: "<off|low|medium|high 中当前模型声明的一枚>" }
```

- 收据档位＝next_task：下一任务即用，不用重启；当前正在跑的一轮不受影响。
- 值若不是当前模型声明档（或不在枚举内），宿主拒收并在收据里说明；页面请以
  `thinkingLevels[]` 为可选项来源，勿自造四档。

## 3. 状态对齐

- 发送成功后无需本地改值：下一次快照推送的 `modelSwitch.currentThinking` /
  `current` 即新值（宿主写盘＋内存同步后由泵读出）。
- 拒收（空值/非法档/无写腿）会返回明确错误句；页面显示句子，勿静默重试。

## 4. 交编排者的接续点（不在本票写面）

1. `cmd/wisp/panel_config_store.go` 的 `ApplySetting` switch 需要一行
   `case panel.FieldRoleChatThinking: written, err = s.mgr.SetRoleChatThinking(w.Value)`，
   否则该字段在活机上按 default 分支被 leg 拒绝（不静默，但功能不通）。归 188 或下一枚
   cmd/wisp 腿。
2. 装配根（run.go 装配区）把 `mgr.ChatSwitchOptions()` 接上
   `PumpSources.ModelSwitch`，键才会在活机快照上线。同上。
3. 前端 panel.ts 的 `PanelSnapshot.modelSwitch` 类型声明＋控件，归界面腿（Q-51）；
   在那之前 Go 侧两枚对账尺保持预定红（红因已在票面与裁决表登记）。
