# 票 264 AC#0 现量复算（读/分析腿，2026-11-06）

> 隔离（规则 21）：全部读数出自共享克隆 `/d/tmp/wisp/20261106/264ac0/repo`（`git clone --shared` 自主树，branch `dev`，起手 `git status --porcelain` = 0 行，主树在飞未跟踪件一律不在克隆里）。克隆 HEAD = `86fb145ee`（2026-11-06 现跑；派单锚 `bca0990f8` 为其祖先，读数取更新的现 HEAD）。本腿**零产码、零票改、零 push**；未跑 `cmd/wisp` 全量、未开 GUI。

## 0. 两枚读数（先清点、不引票面旧数）

| # | 读数 | 值 |
|---|---|---|
| **①** | `internal/agent` + `internal/tools`（非测试、非 testdata）里 `.Error()` 装配点，**能进模型可见文本**（经已核实的生产消费链）的枚数 | **103 ／ 112**（tools 95／100；agent 8／12） |
| **②** | 同两包出站路径上**既有脱敏点**（`redact|Redact` 词面） | **0 枚功能性脱敏**（全 grep 唯一命中是一句注释，注释自己写着"这里没有 redaction"） |

①的子切分：**39／103 构造上就带机器路径**——12 枚错误源是 C26 `Canonicalize`（`internal/risk/pathresolver.go:95` 的 `%w: %s expands to %s (%s)` 同时带原始拼写与规范化绝对路径），25 枚错误源是 `os.*` 文件调用（`*fs.PathError`/`*os.LinkError` 逐字带绝对路径），另 2 枚（`web_open.go:367/371` app.launch）大概率带。

## 1. 复算命令（Git Bash，克隆内逐条可重跑）

```bash
git clone --shared "/d/work/workspace/projects plans/Wisp" /d/tmp/wisp/20261106/264ac0/repo
cd /d/tmp/wisp/20261106/264ac0/repo && git rev-parse HEAD   # 86fb145eed4ee8ce0912c5468d3208fc9fe3bf8d

# [C1] 原始普查（112 枚，逐枚留档 /d/tmp/wisp/20261106/264ac0/assembly-raw.txt）
git grep -nE "\.Error\(\)" -- internal/agent internal/tools ":(exclude)*_test.go" ":(exclude)*/testdata/*" | wc -l   # = 112

# [C2] 读数②（唯一命中 = internal/agent/approval/cardrender.go:114 注释）
git grep -nE "redact|Redact" -- internal/agent internal/tools ":(exclude)*_test.go" ":(exclude)*/testdata/*"
# [C3] 宽网（mask/sanitize/scrub）：仅 sanitizeArgs（bridge.go:1396，司法取证列轴，非出站文本）
git grep -niE "mask|sanitiz|scrub" -- internal/agent internal/tools ":(exclude)*_test.go" ":(exclude)*/testdata/*"

# [C4] 分母：全仓非测试（含 .scratch 在册探针）= 667；产码域（internal+cmd+tools）= 544
git grep -n "%w" -- "*.go" ":(exclude)*_test.go" | wc -l
git grep -n "%w" -- internal cmd tools ":(exclude)*_test.go" | wc -l
# [C5] %w 且自带另一插值动词：产码域 = 335（全仓含探针 = 401）
git grep -nE "Errorf\(" -- internal cmd tools ":(exclude)*_test.go" | grep "%w" | grep -E "%[sdvqxfe]" | wc -l
```

分类方法：112 枚逐枚读码，判据 = **有已核实的生产消费者把该串放进 `Result.Text` 或 `log.Text`**（`internal/agent/loop.go:819/815`）→ D15 spill（:830-838）→ `toolResultMessageParts` 入史（:843 → `l.append` :1171-1174，`toolResultParts` :1251-1253 逐字放行）→ `buildRequest` `asm.Build(in, l.History())`（:664）→ `provider.Stream`（:639）。五跳全部现读核过，非转述。

## 2. 112 → 103：排除的 9 枚（逐一写明去向）

| 枚 | 去向（现读核实） |
|---|---|
| `agent/spill.go:215`、`tools/task.go:840` | 注释（后者明记原接缝已被 A587 拆除——票面引的那枚接缝今天已不存在） |
| `tools/system_backend_windows.go:309,453` | `e.Error() != "…"` 比较，非装配 |
| `tools/bridge.go:1122`（displayPaths） | 审计行 + 审批卡（用户轴） |
| `agent/approval/gate.go:658` | `q.abandon` 队列/卡记录（模型侧返回 :660 是静态串） |
| `agent/loop.go:1046,1047` | 任务终态 `res.Message` + `EvError` 事件（用户面板轴） |
| `tools/task_backfill.go:136` | 返回串经 `cmd/wisp/run.go:1782` 只进 `rt.auditf`（审计轴） |

## 3. 带路径子集（39 枚清单出处）

- **C26 源 12 枚**：`fs.go:144,208`；`fs_edit.go:119`；`fs_write.go:260,389,460,464,611`；`document_text.go:553`；`search_content.go:197`；`search_files.go:90`；`web_open.go:264`（消息形全部是「路径无法解析（按 fail-closed 拒绝）：」/「拒绝检索/搜索（fail-closed）：」+ `err.Error()`）。
- **os 调用源 25 枚**：`fs.go:148,158,212,218`；`fs_edit.go:128,138`；`fs_write.go:292,325,343,408,477,505,514,540,559,620`（含 :343 原子重命名，:341 台账同记、:342 `baseOf(tmpName)`）；`search_content.go:204,217`；`search_files.go:97,110`；`web_open.go:253,270,276,279,286`。
- **大概率 2 枚**：`web_open.go:367,371`（app.launch 进程启动错误，可能含 exe 路径）。
- 其余 64 枚错误源为 JSON 解码/上下文取消/断言校验等，构造上不带路径（ctx.Err 为静态英文串）。

## 4. 仪器侧现状

- 出站侧脱敏 = 0（读数②）；`redact` 词面在产码里只以一句否认注释存在（`cardrender.go:114`），**不许留下"我们有脱敏开关"式的句子**——本仓产码里连开关都没有。
- `cmd/wisp/task_output_pointer_notice_174r8_test.go`（原票面引的文件名已被 174-r8 重构取代）的 T1/T2 钉的是 task.output 回执授权文案，且 T1 **要求回执含绝对路径**（`全文见 <artifact>`）——任何"出站掩路径"的改动必须保住这枚指针钉，否则 174 轴先红。
- 今天没有任何仪器量"失败句里的绝对路径进了模型正文"这一格（这正是 AC#1 要造的能力尺）。

## 5. AC#1.. 决策/建议表（取形归机主——票面明文"未定义即停"格；技术落点本腿拍）

| 形 | 需要的生产改动与落点 | 规模 | 买到 / 买不到 | 撤销口令 |
|---|---|---|---|---|
| **甲＝出站前脱敏** | **咽喉一处**：`internal/agent/loop.go` `toolResultParts`（:1251）或 spill 后、append 前加一步路径掩码，**一处覆盖全部 103 枚**；只掩模型侧副本，用户侧（:847 `EvToolEnd` 事件、`res.Message`）不掩 ⇒ "用户看得见的提示仍说得出是哪个文件"天然成立。⛔ 不许做成 103 处逐点改（不可维护、必漏）。须保住 174r8 的 task.output 指针钉 | 1 函数 + 尺；写面跨 `internal/agent` ⇒ 需具名解冻 | 服务商正文不再带机器目录结构；代价 = 模型看不见真路径可能多试一轮 | 票面只发了乙/不做的口令；取甲须编排者补一枚 |
| **乙＝只装仪器** | 零行为改动；能力尺落 `internal/agent`（loop 侧）或 `cmd/wisp` 端到端（174r8 同形）：种一发带绝对路径的失败句 ⇒ 指名断言到"历史里那条 RoleTool 文本含/不含绝对路径"那一跳；反向去掉必不响（AC#2 两向） | 0 行产码 + 测试 | 只买"以后谁改这块会被量到"；**不减少已外发内容** | 「264 走乙」 |
| **不做＝登记取舍** | 零；票面与台账写死"工具失败句带真路径发给服务商＝自用期既有取舍"，并写明读数②=0（连开关都不存在） | 0 | 什么都不买，但把实情说到明处 | 「264 不做」 |

**本腿技术裁定（取形本身仍归机主）**：若取甲，落点必须是咽喉一处而非逐点——理由：103 枚装配点散在 24 个文件，逐点改任何新工具一落地就漏；咽喉一处同时天然满足 AC#1 的"模型侧无路径、用户侧仍说得出"两向。分母与落点全部出自 §0-§4 现量。
