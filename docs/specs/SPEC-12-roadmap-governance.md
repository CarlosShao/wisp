# SPEC-12 · 切片路线图与治理

> 追溯：§4/D44(a)(b)/M0/§7/§12/D22/D23；每张切片的完整任务卡按 §12 在 S1 生成到 `docs/slices/`，
> 本 spec 是其权威摘要与治理规则。

## 1. 为什么是垂直切片

草案 Phase 1/2/3 按层横切，Phase 1 结束跑不通一条完整链路。垂直切片保证每片结束都有真能用的
东西——「自用优先」（D23）与 AI Agent 开发（D22，需要频繁可验证里程碑）都依赖这一点。
**每张切片卡必须显式写「本切片明确不做什么」**（M0 要求②，双闸门：防扩范围、防误判完成）。

## 2. 切片路线图（S0–S8；done 判据机器可验收）

| 切片 | 目标 | done 判据（摘要，全文见切片卡） | 明确不做 |
|---|---|---|---|
| **S0** | Spike + 构建链打通（D44 五项输出） | ① 干净环境一键构建 exe，流程冻结 `docs/BUILD.md`（**最高优先级**）② X/Y 判定（输入含 cgo 崩溃概率与 10s 回落）③ goja async/ES module 结论 ④ go-webview2 冷/热拉起实测 ⑤ 五个 RSS 基线 + 模型常驻内存实测表 + ASR↔TTS 串行切换延迟 → 回填 D32 16.3.3 | 任何产品功能、悬浮球动画、Agent 循环 |
| **S1** | 最小通路：快捷键→文字输入→Agent→通知 | 打字→回复→3s 回落；空闲进程树 RSS ≤25/40MB（按 spike 判定路径）实测过；`wisp run` CLI 可用；C21 token 落地（原生侧）；网络/TLS 错误分类、单调时钟、per-session 互斥、Job Object、`notify`+`list_tools` | 语音全部、WebView 面板、门控 UI、插件抽象、记忆 |
| **S2** | 语音输入：麦克风→VAD→流式 ASR→Agent | 说话→出文字→走 S1 通路；CER 双门禁 6%/15%（两 wav 集入库）；D26 模型分发落地（P5 阻塞本片）；音频热插拔；`Downloading` 态；`intra_op_num_threads=1` | TTS、唤醒词、AEC |
| **S3** | 安全层 + 称职助手能力面（判据必须同时含能力项与安全项） | **能力项**：D34 表 S3 工具逐个跑通留证据 + 四场景①②脚本化验收；**安全项**：L0 直通/L1 阻止窗口/L2 强确认、C19 生效（声明 L0 的危险工具仍判 L2）、越界被拒、D30 黑名单硬拒、组合闸门生效、`failToolCallsFromTruncatedMessage` 生效、C26 四连红队、C25 四通道、C29 篡改用例、批量聚合确认（D45-1）、C7 Image 部件可用 | Tier2 goja、命令面板、并发（TaskScheduler 只许单任务——**安全不完整期约束**） |
| **S4** | 语音输出 + L1 画像 + 会话保活 | 短结果播报、长结果落文件；画像异步提取不占感知延迟；TTS 音质 ≥7/10 门禁（P7）；标点（P4）；`Warm`/`Conversation`（C31 + D43 #26/28/29/31）；`Warm` 内二次唤起首字 P50 ≤1.5s；`reminder.*`/`memory.*` | 命令面板、GUI 配置 |
| **S5** | 按需 WebView：命令/结果面板 + L2 确认卡 + 配置编辑器 | 冷 ≤1500ms/热 ≤200ms；用完销毁或隐藏；RSS 达标；TOML↔GUI 双向一致；L2 确认卡完整参数 + **无「允许」按钮**（原生侧批准指引）；CSP+净化+服务端授权三层；无面板模式原生降级卡；CSS 侧复用 C21 token（人工验收） | **常驻 WebView（禁止）** |
| **S6** | 唤醒词 + 看门狗 + 可观测性 + 成本 | KWS 开启空闲 ≤90/110MB、CPU ≤2%；看门狗按态查表自动卸载（Armed 不卸 KWS）并留日志；诊断包不含音频/Key/转写全文；CostMeter 面板可用；D42#1#5#6#8#10 全过 | AEC/barge-in |
| **S7** | 并发 + 审批队列 + Tier2 goja + 生命周期收尾 | C18：并发确认按 correlationId 路由、队头单显+深度计数、超时判拒绝；C20 路径冲突排队面板可见；C22 梯度提醒触发进 Stuck；D45-2 会话授权；C24 契约测试；**端到端真插件 = lark-cli 包装插件**（含 exe_hash 篡改拒执行）；goja 按需加载卸载+加固；崩溃恢复/单实例/自启/更新/DPI/显示器拓扑/失败预演全过 | 插件 SDK 文档、registry |
| **S8** | （开源前才启动）macOS + 签名分发 + i18n + 插件 SDK + registry | 见 `docs/DEFERRED.md` 各条完成判据；更新流程（D41b）+ N-1 回滚；卸载残留提示 | — |

依赖：`S0 → S1 → {S2 ∥ S3} → S4 → S5 → S6 → S7 → S8`。
**⚠ 安全不完整期**：S3 落地门控但无并发审批队列，S7 才补上——S3–S6 期间
`TaskScheduler` 只允许单任务，此约束写进各切片卡「明确不做」。

## 3. 场景 ↔ 切片追踪矩阵（D44(b)；对抗验收双向核对，任一侧落空 = 缺口回报而非自行补）

见 SPEC-00 §4。核对规则：每场景每项能力能在某切片 done 判据找到；每判据能追到决策或场景。

## 4. 治理规则

### 4.1 契约变更流程

改 C1–C31 或 D1–D46 = **人工批准**；同步更新 PLAN.md、`docs/DECISIONS.md`、受影响切片卡。
agent 单方面改契约 = 跑歪模式 #1，对抗验收判失败。

### 4.2 未定义即停（D22 闸门③）

方案未覆盖的情况 → 停下来问，不得自行假设。已知待定案项：
- `AwaitingApproval` 遇系统挂起（倾向作废判拒绝，S7 切片卡定案）
- `web.search` 实现路径（抓结果页 vs API，S3 前）
- Go module 组织名 + P10 命名核查（S1 建仓前，阻塞）
- 便携模式 × DPAPI（SPEC-02 §6 建议，S1 定案）
- C24 GojaHostAPI 初始集与 C17 方法白名单定稿（S7/S5 切片卡批准）
- WebView2 冷拉起 >2s → 重评 L2 卡是否回原生（P11，S0）

### 4.3 每切片完成后的强制动作

1. 重跑**缺口审计**（审计者 ≠ 实现者，D22 双角色）
2. 失败预演（SPEC-10 §7）
3. 对抗验收报告（SPEC-10 §8 的七项检查）
4. 场景↔切片双向核对
5. DEFERRED/RESERVED/REJECTED 登记表更新（有新增推迟必须五字段齐全）

## 5. 推迟项登记表（M0 要求①：五字段缺一不可）

类型语义：**DEFERRED** = 有计划不做；**RESERVED** = 只留接口位、无实现计划、不得误读为待办；
**REJECTED** = 有意取舍、不是债务、不得被「修复」。

| 类型 | 项 | 为什么现在不做（依据） | 完成判据 | 前置 | 当前残缺表现 |
|---|---|---|---|---|---|
| DEFERRED | macOS 平台层 | D23 自用期仅 Windows | 两套悬浮球实现均过交互测试；CI 出 macOS 产物 | macOS 实机/可靠 CI | macOS 完全无法运行 |
| DEFERRED | 代码签名/包管理器分发 | D17+D23 | SignPath 签名 SmartScreen 不拦；winget/Scoop 可装 | SignPath 审批 | SmartScreen 警告需「仍要运行」 |
| DEFERRED | 插件 SDK+文档+示例 | D23 | 第三方照文档独立写出 Tier1/Tier2 各一过契约测试 | S7 稳定 | 只有作者能写插件 |
| DEFERRED | 社区插件 registry | D23 | 可搜索索引+安装命令+版本兼容校验 | 插件 SDK | 只能本地路径/Git URL 安装 |
| DEFERRED | i18n/非中文 | D23（⚠ 语音链路是中文模型，非 locale 问题；非中文用户语音功能不可用只能走文字通道） | UI 全外置+英文 locale+英文 ASR/TTS/KWS 模型过基线 | S5+英文语音模型 | 界面中文；语音仅中文 |
| DEFERRED | 无障碍 | 未讨论遗漏后补；自用期非阻塞 | 面板过 axe 检查；悬浮球状态有非视觉等价反馈；全流程纯键盘 | S5 | 屏幕阅读器不可用面板 |
| DEFERRED | AEC/真 barge-in | D16 | 播报期可语音打断不自激 | C8 已留 Aec 接口位 | 播报时无法语音打断 |
| DEFERRED | 快捷键路径语音否决（B1） | KWS 未加载+ASR 加载 1–3s > 2–3s 窗口，物理不可能 | 快捷键会话中说「取消」能否决 L1 | AEC | Confirming 明示「语音取消不可用」 |
| DEFERRED | 剪贴板历史（D34） | 无差别记录隐私风险过高 | 敏感内容识别+自动跳过+逐条删 | 敏感识别器 | 只有当次读写 |
| DEFERRED | `doc.read` xlsx/OCR（D34） | 工量大/额外模型冲击内存 | OCR 模型常驻 ≤80MB 不破 700MB；xlsx 读出表格结构 | S3 的 PDF/docx | 读不了 Excel 与图片文字 |
| DEFERRED | `system.eject`（D34） | 低频，与 system.power 共路径 | 能弹出设备并反馈失败 | S3 的 system.* | 无法语音弹 U 盘 |
| DEFERRED | 完整错误文案体系 | D23 | 所有用户可见错误有人话文案+建议（含 D37 全 17 类） | S6 | 部分错误暴露技术细节 |
| DEFERRED | 云端 ASR/TTS provider | D5 抽象已留，本地够用 | 至少一个云端 provider 过 C9 契约测试 | 用户自备 Key | 嘈杂环境无高精度逃生通道（带噪 CER 门禁靠这条兜） |
| DEFERRED | 竞品对比文档 | 自用期非阻塞 | README 第一段能说清「为什么不用现成的」（四条可验证差异：真轻量/默认不碰麦克风/本地优先/插件不常驻） | S8 前 | README 差异化论证缺位 |
| RESERVED | MCP client | D13 与 host bridge 收口冲突（16.9#7 驳回引入） | — | ToolProvider 接口位已留 | 无法接 MCP server |
| RESERVED | Codex guardian 双轴评分 | D31：+300–800ms 在响应路径上；**重评触发条件：L1 误执行率不可接受时优先启用此项而非放宽门控** | — | C19 可插拔判定器 | 「风险不高但用户没要求过」抓不到 |
| RESERVED | 跨会话持久授权档（workspace/user） | D45 后仅剩跨会话档；持久授权×提示注入风险窗口大 | — | C18 决策类型可扩展 | 每会话重新授权 |
| RESERVED | `available_decisions`（服务端下发按钮） | 确认动作已固定三种 | — | C18 | 按钮硬编码 |
| RESERVED | embedding 语义记忆 | 16.9#3：常驻 embedding 击穿 90MB；D20 关键词已够 | — | MemoryStore 接口位 | 模糊描述查不到 |
| RESERVED | Linux 支持 | D7 GNOME/Wayland 无法实现悬浮球 | — | 平台抽象仅 Windows impl | Linux 上无法运行 |
| RESERVED | 精确 taint tracking | 16.9#1：LLM 改写使 token 级失效；C25 粗粒度兜 | — | — | 污染检测是片段级：改写/翻译后外发可能漏检 |
| REJECTED | 周期性调度器（cron） | D12 有意取舍；一次性 `reminder` 已移出为窄例外（16.5.4） | — | 改 CLI+OS 调度 | 需自配任务计划；文档给配方 |
| REJECTED | 邮件/日历/IM **核心内置** | §15 第 10 项定案：核心永不内置；**能力走 D46 插件路径已采纳**（lark-cli 包装） | — | D46 已定 | 核心不内置，但插件路径完整可用（S7 验收用例） |
| REJECTED | Coding/IDE 能力 | 草案 1.1 | — | — | 不能写代码 |
| REJECTED | 多 Agent 协作 | 草案 6.6（与 D31 多任务并发区分：并发是同一 Agent 跑多任务） | — | — | 单 Agent |
| REJECTED | 意图分类四级流水线 | D11 过度设计已废 | — | D11 两层替代 | 无前置分类器（特性非缺失） |
| REJECTED | `fs.delete` 默认提供 | D34：有 `fs.trash` 够用；须 `[fs] delete_enabled=true` | — | fs.trash | 默认无法语音永久删除（有意） |

**双向交叉核对**：代码内每个 `// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点` 必须在表有对应条目，
反向亦查（对抗验收逐片执行）。

## 6. 文档交付物清单（§12；S1 时按此生成骨架）

- 人读版：`PLAN.md`（定稿）· `DECISIONS.md`（D1–D46，agent 只读）· `DEFERRED.md`（§5 表）·
  `RISKS.md`（§9 + D42）· `PRECHECK.md`（P1–P14 及结论回填）· `ARCH.md`（D35–D38/D41）·
  `SEQUENCES.md`（D40 五时序）
- Agent 执行版：**薄** `AGENTS.md`（根：D1–D46 一行摘要 + 禁止清单 + 未定义即停 + docs 索引）·
  `docs/contracts/C1..C31.md` · `docs/STATE_MACHINE.md`（D43）· `docs/TOOLS.md`（D34）·
  `docs/slices/S0..S8.md` · `docs/SLO.md`（D32 + 进程树口径 + 可重跑命令，阈值禁改）·
  `docs/BUILD.md`（S0 冻结）
- 本目录 `docs/specs/`：与上列互补——spec 讲「怎么做」，契约/状态机/工具表讲「是什么」。
- 时序视图五条（D40，落 SEQUENCES.md）：冷唤起单轮回落 · Warm 多轮 · 并发审批冲突 ·
  失败与恢复（LLM 5xx + cgo 崩溃）· 注入检出。

## 7. 前置核实任务（P1–P14，阻塞关系）

| 任务 | 阻塞 |
|---|---|
| P1 spike（五项输出）/ P10 命名核查 / P11 WebView2 冷拉起 / P13 DPAPI 便携 | **S1** |
| P5 模型镜像（升级为阻塞）/ P3 模型权重许可 | **S2** |
| P12 同步盘识别 / P14 KWS 否决词 | **S3** |
| P4 标点评测 / P7 TTS 音质（门禁） | **S4** |
| P6 WebView2 Runtime 普及 | S5 |
| P2 goja async/ES 实测 | S7 |
| P8 SignPath | S8 |

## 8. 下一步（PLAN §15 建议第 3–5 步的执行顺序）

1. ~~B 段 8 项定案~~（已全部定案，2026-09-18）
2. P1–P3、P5、P10–P14 前置核实（阻塞关系见 §7）
3. 按 §6 生成两版文档骨架（contracts/slices/STATE_MACHINE/TOOLS/SLO/BUILD）
4. **从 S0 spike 开始**——五项结论会反确认/修正 D25/D24/D29/D32；结论出来前不开 S1
5. 每切片后重跑缺口审计（§4.3）——**四轮打磨未穷尽是方案自己写下的结论，不得省略**
