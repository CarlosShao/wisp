# SPEC-01 · 总体架构、仓库骨架与运行时模型

> 追溯：D2/D24/D25/D32/D37/D38/D41、§1/§14.9/§14.10、C11/C30/C31；状态机细节在 SPEC-08。

## 1. 背景与问题

- 轻量级 = 运行时资源约束（D1）：空闲真睡眠（进程树私有内存 ≤25MB/路径Y 或 ≤40MB/路径X，D32）、
  无常驻 WebView、CPU ≤0.5%。
- 开发方是 AI Agent（D22），架构必须消除歧义：模块边界冻结、goroutine 全员有名、关停顺序是契约。
- Go + cgo（sherpa-onnx）的内存行为未实测 → 拓扑 X/Y 由 S0 spike 判定，架构必须两条路径都成立（D25）。

## 2. 进程与内存拓扑

```
常驻主进程（纯 Go，唯一；路径 Y 下含 cgo 语音运行时）
  ├ 悬浮球        Win32 分层窗口（x/sys/windows）+ Direct2D/DirectWrite（STA 线程）
  ├ 托盘 / 全局快捷键 / 单实例锁（per-session）
  ├ BallState 状态机（20 态，C12/D43） + 会话控制层（本地正则，D11）
  ├ 资源看门狗（按态查表，D32 修订③）
  ├ 配置（config.toml 唯一真相源，D6） + SecretStore（C28）
  ├ Job Object 持有者（C30：所有子进程入 Job，度量 TreePrivateBytes()）
  └ [KWS] opt-in（Armed 态，RSS ≤90/110MB）

按需（工作态，峰值 ≤700MB，目标值——S0/S2 实测前不作验收值，D32 16.3.3）
  ├ 语音运行时：路径Y=主进程内延迟初始化；路径X=语音子进程（PCM 走 stdio pipe）
  ├ Agent 循环 / LLM adapter / 工具执行（工具并发 ≤4）
  ├ 面板 WebView（C27 单例单窗口，隐藏而非销毁；WebView2 派生 3–5 子进程全部入 Job）
  └ SQLite（WAL，单写者 db-writer goroutine）

回落：Settling 3s → Dispose（C11 逆序销毁 + debug.FreeOSMemory()）→ 10s 内回到该态上限
异步（不占感知延迟）：Warm 窗口内做 L1 画像提取（D20）
```

- D25 spike 判定（写死，不留裁量）：空闲 private RSS ≤25MB → 路径 Y；>25MB → 路径 X。
  判定输入必须含 **cgo 崩溃概率**（C 侧 segfault 无法 recover）与 **10s 内能否回落**；
  若 X 下回落不达标 → 直接判 Y（D32 16.3.5）。
- 路径 X 的强论据：子进程崩溃只带走语音能力，悬浮球存活并自动重启子进程（D37(b)）。
- 排除项（D2）：stub+门卫双进程、常驻 Tauri、零常驻只留托盘。

## 3. 仓库骨架（【SPEC】由冻结模块表 1:1 映射，不得增删模块）

```
wisp/                                  # Go module: github.com/CarlosShao/wisp（P10 定案前为占位）
├── cmd/wisp/main.go                   # 唯一二进制入口；无参=GUI 常驻；`wisp run "..."`=CLI（D12），
│                                      #   CLI 经 AttachConsole 输出；`wisp doctor` 自检（S0 起）
├── internal/
│   ├── buildinfo/                     # 版本、WISP_ENV、C29 minisign 公钥（硬编码进二进制）
│   ├── ball/                          # 悬浮球：分层窗口绘制、鼠标/拖拽、多显示器定位、托盘（无业务逻辑）
│   ├── statemachine/                  # BallState 20 态 + D43 转移表 + 每态超时（只有转移，没有能力）
│   ├── audio/                         # AudioSource(C8)、WASAPI 采集（pinned 线程）、VAD、半双工开关
│   ├── speech/                        # AsrEngine/TtsEngine/WakeWordEngine(C9) sherpa 实现 + 云端接口位
│   ├── agent/                         # ReAct 循环、上下文预算、取消、截断 tool call 处理
│   │   ├── approval/                  #   C18 ApprovalQueue（归属推定：D31/D37「审批 UI 归 agent/panel」）
│   │   └── scheduler/                 #   C20 TaskScheduler + PathLock（并发任务与路径冲突）
│   ├── llm/                           # 三协议 adapter、归一化 StreamEvent(C6)、重试退避、usage
│   ├── tools/                         # ToolProvider(C4)、host bridge（capability/risk 强制收口）
│   ├── plugin/                        # manifest 加载校验(C14/C16)、goja 运行时、DisposalScope(C11)
│   ├── memory/                        # L1/L2/L3 存储(C13)、SQLite 生命周期、RetentionJob
│   ├── panel/                         # WebView 宿主(C27)、PanelBridge(C17)、webassets embed
│   ├── models/                        # 模型下载/校验/续传/本地指向（D26）+ C29 manifest
│   ├── config/                        # config.toml schema(C15/D36)、热加载、迁移
│   ├── secret/                        # SecretStore(C28)：DPAPI/env 引用、明文迁移、日志脱敏 Key
│   ├── risk/                          # RiskAssessor R1–R9(C19)、PathResolver(C26)、Provenance(C25)、A/B 黑名单
│   ├── session/                       # SessionScope(C31)、Warm/Conversation 计时、模型引用计数、D45 授权登记
│   ├── proc/                          # JobScope(C30)、per-session 单实例、关停顺序编排(D38e)、更新(D41b)
│   ├── watchdog/                      # 资源看门狗（按态查表；含配置文件 mtime 轮询复用其循环，SPEC-03 §4.3）
│   └── observe/                       # slog JSONL、脱敏、诊断包、SLO 采样、CostMeter(C23) 落账
├── frontend/                          # React 面板（S5+；技术栈见 SPEC-08 §5）
│   ├── src/…  └── vite.config.ts
├── assets/web/embed.go                # //go:embed frontend 产物（AddWebResourceRequestedFilter 喂给 WebView2，
│                                      #   不起 localhost HTTP 服务，D29）
├── design/                            # 已有静态原型（C21 tokens.css 参考实现）
├── docs/                              # PLAN.md、specs/、contracts/(S1 生成)、slices/(S1 生成)、
│                                      #   STATE_MACHINE.md、TOOLS.md、SLO.md、BUILD.md、DEFERRED.md、…
├── scripts/
│   ├── fetch-deps.ps1                 # 下载 sherpa-onnx/onnxruntime 预编译库（版本+SHA256 钉死在 deps.toml）
│   ├── slo-check.ps1                  # D32 六态采样 + 回落验收（SPEC-10 §3）
│   ├── dev.ps1 / test.ps1             # 本地起 compose mock + 构建 + 测试（SPEC-11）
│   └── spike/                         # S0 五项实测的测量程序（D44 S0 修订）
├── docker/
│   ├── builder.Dockerfile             # Go + mingw-w64 构建镜像（SPEC-11 §3）
│   ├── frontend.Dockerfile            # Node 20 构建镜像（输出 dist）
│   └── compose.{dev,test}.yml         # mock-llm / mock-web / mock-search / model-mirror / ssrf-target
├── tools/mockllm/                     # Go 实现的 OpenAI/Anthropic 兼容 mock（SPEC-11 §4）
├── testdata/
│   ├── asr-baseline/near_clean/       # 近场清洁 wav 集（CER 门禁 ≤6%，必须入库）
│   ├── asr-baseline/noisy_far/        # 带噪远场 wav 集（≤15%，不阻塞发布）
│   ├── golden/                        # 录制的 SSE 流回放文件
│   └── models/                        # 测试用微型 onnx 模型（不下载真模型即可跑 CI 的最小桩）
├── deps.toml                          # 原生依赖版本与哈希钉死【SPEC】
├── .devcontainer/devcontainer.json    # 纯逻辑开发容器（SPEC-11 §5）
├── third_party/                       # fetch-deps 产物（git 忽略）：sherpa c 头/库、onnxruntime.dll
└── go.mod
```

依赖白名单（D22；新增须人批准）：`sherpa-onnx` Go 绑定（cgo）· `goja`（纯 Go）· `jchv/go-webview2`
（纯 Go）· `golang.org/x/sys` · SQLite 驱动 · TOML 解析 · `goreleaser`（打包）。
【SPEC】白名单类别内选型：SQLite 用 `modernc.org/sqlite`（纯 Go，不再叠加 cgo）、TOML 用
`pelletier/go-toml/v2`、结构化日志用标准库 `log/slog`（零新依赖）。前端白名单见 SPEC-08 §5。
明确排除：`wails`、`webview/webview_go`、任何 Node/Python/Deno **运行时**（构建期工具链除外）、
LangGraph 类编排框架。

## 4. 线程与 goroutine 模型（D38，契约级）

### 4.1 专属线程（仅两条，其余交给 Go 调度器）

| 线程 | 约束 | 原因 |
|---|---|---|
| `ui-sta`（STA/UI） | 拥有悬浮球窗口 + 托盘 + **WebView2 窗口** + 全部 Win32 消息循环 | WebView2 要求 STA + 消息泵；D2D factory 有 thread affinity → 悬浮球与面板**共用同一条 STA 线程**，不建第二个 D2D factory |
| `audio-capture`（pinned） | `runtime.LockOSThread()`；**该线程不得跑任何其他 Go 代码** | WASAPI 回调需要固定线程；goroutine 会被 M 迁移，不 pin 则爆音/丢帧 |

### 4.2 goroutine 花名册（必须有名字/owner/生命周期/退出条件；无名 goroutine = 泄漏）

| 类别 | 名字 | 说明 |
|---|---|---|
| 常驻（6） | `ui-sta` · `audio-capture`(pinned) · `hotkey-listener` · `db-writer` · `watchdog` · `log-flusher` | 上限 6，SLO 采样比对 |
| 按需（+1） | `kws-infer` | 仅 Armed 态 |
| 每任务（3） | `agent-task-<id>`（持 root ctx）· `tool-exec-<id>` · `approval-waiter` | 派生自任务 root ctx |
| 临时 | `asr-infer` · `tts-infer` · `panel-host` · `retention-job` · `model-downloader` · `memory-extract` | 均在对应 DisposalScope 内 |

总数上限 = 常驻 6 + 每任务 3；超出即泄漏征兆 → 看门狗记录并进 SLO/诊断包。
**禁止裸 `go func()`**：必须有名字、owner、退出条件，入口 `defer recover()`（D22 禁令，AST 扫描验收）。

### 4.3 并发与背压

- 每任务 root context：任务内所有 goroutine 派生自它；**任务完成必须等全部派生 goroutine 退出
  （WaitGroup）才算完成**——否则 RSS 不回落（D32 的 10s 目标直接依赖此条）。
- 音频→ASR：有界 channel ≤200ms 音频，满则丢帧**并计数**（计数进日志与诊断包，静默丢帧=识别率
  莫名下降却查不出）。
- LLM 流→UI：有界 channel，满则**合并增量**（不丢内容）。
- 工具并发上限 **4**（防 LLM 一次 50 个 tool call 打爆机器；与 D45 批量聚合配合）。

### 4.4 关停顺序（D38(e)，契约级，顺序错误即事故）

```
1. TaskScheduler 关门（不再收新任务）
2. 停热键与唤醒词监听
3. 取消所有任务 root ctx → 等 WaitGroup ≤3s（超时放弃等待并记录未收尾任务）
4. 停音频采集线程 → 关设备        ← 必须在释放 ASR 之前（否则 ASR 拿到半截 buffer）
5. 释放 ASR/TTS/KWS session（DisposalScope 逆序销毁）
6. 销毁面板 WebView（若存活）→ 等子进程退出 ≤2s
7. flush 日志 → db-writer 排空 → wal_checkpoint(TRUNCATE) → 关 SQLite
8. debug.FreeOSMemory()
9. 关闭 Job Object（子进程由 OS 连带清理）
10. 退出
```

Windows 会话结束（`WM_QUERYENDSESSION`）走快速版：跳过第 7 步非关键 flush，但 **4/5/9 不可跳**。

## 5. 错误模型（D37，契约级）

### 5.1 error_class 有限枚举（17 类，落 `task_log.error_class` 与日志字段；不得自由字符串）

| 类 | 典型 | 用户可见 | 映射态 | 可重试 |
|---|---|---|---|---|
| `config` | schema 非法/缺 Key/目录未授权 | 指引 + 开配置面板 | Unconfigured | 否 |
| `auth` | 401/403、Key 失效 | 提示重填 Key | Unconfigured/Error | 否 |
| `network` | DNS/超时/TLS/代理 | 「网络不可用，可切文字模式」 | NoNetwork | 指数退避 ≤3 |
| `rate_limit` | 429/配额 | retry-after + 本次成本 | Error | 尊重 retry-after |
| `provider` | 5xx/协议不合规 | 「服务异常」+ 切 fallback_provider | Error | 是 |
| `model` | 缺模型/哈希签名不符/加载失败 | 指向下载或校验详情 | Downloading/Error | 否 |
| `audio_device` | 被独占/权限拒绝/热拔出 | 明示设备名 + 指向系统设置 | Error | 重枚举一次 |
| `asr` | 识别失败/纯静音/过短 | 「没听清，请再说一次」 | Listening(重试1次)/Error | 1 次 |
| `tool` | 工具内部失败/参数不合法 | 回给 LLM 自纠，不打断用户 | Acting | Agent 决定 |
| `permission_denied` | 越界/未声明 capability/黑名单/reparse | 回给 LLM + 面板可见 | Acting | 否 |
| `user_rejected` | 审批拒绝/超时 | 「已取消」 | Settling | 可一键重放 |
| `cancelled` | 用户取消/关停 | 静默 | Settling | 否 |
| `budget` | token/轮数/时间/费用超限 | 明示哪一项 | Stuck | 否 |
| `loop` | 梯度刹车触发 | 明示重复了什么 | Stuck | 否 |
| `injection` | C25/R4 检出 | 明示来源与命中 | Error | 否 |
| `internal` | panic 被 recover/意外 nil | 诊断包入口 | Error | 否 |
| `resource` | 磁盘满/内存超限/看门狗/GDI 句柄超限 | 一键重启 | WatchdogAlert | 否 |

### 5.2 跨边界传播（三条最容易出事的边界）

- **cgo 边界**：C 侧返回码 + `char*` 错误串 → Go 侧**立即拷贝为 Go string 再释放**（不得持有 C 指针，
  否则 GC 移动后 = 随机崩溃不可复现）；C 侧 segfault 无法 recover → 路径 X 的最强论据。
- **goroutine 边界**：panic 不跨 goroutine；每个 goroutine 入口 `defer recover()` → 转 `internal`
  + 结构化日志（goroutine 名 + stack）→ 取消所属任务 root ctx。
- **goja 边界**：JS `throw` → 转 `tool` 错误回给 LLM（让模型自纠）；不得穿透到 Go 调用栈。

### 5.3 取消边界（架构级承认）

`context.Context` 是唯一取消通道，但 **cgo 推理调用不可取消**（onnxruntime `Run` 无法中断）→
推理一律在「可丢弃」的 goroutine/子进程中执行，取消 = 放弃等待结果；路径 X 下取消 = kill 子进程。

## 6. DisposalScope（C11，从「氛围」变规格）

- 语义：`scope.Defer(fn)` 注册；`Dispose()` **幂等、逆序、每个 fn 独立 recover、总超时 3s**；
  超时项记日志并标记 `disposal_incomplete`（必须可见）。
- 强制收尾顺序：① 取消派生 ctx → ② 等 WaitGroup → ③ 释放 native session →
  ④ **`debug.FreeOSMemory()`** → ⑤ 记录本 scope 峰值/末尾内存到 SLO 采样。
- 层级：会话 scope（C31）⊃ 任务 scope ⊃ 工具 scope ⊃ 插件 scope。

## 7. 测试决策

- goroutine 花名册用测试断言：跑完一个任务后 `runtime.NumGoroutine()` 回到常驻基线（+容差 1）。
- 关停顺序写一条**顺序审计测试**：各步骤埋点，断言实际发生序与 §4.4 一致（含会话结束快速版）。
- 错误分类：17 类各一条用例（注入真实错误源：坏配置/断网 mock/429 mock/坏模型文件…）断言
  `error_class`、映射态与重试行为（mock 端点由 SPEC-11 compose 提供）。
- cgo 边界：`tools/mockllm` 之外另提供**崩溃注入的语音子进程桩**（路径 X 下 kill -9 场景），
  断言主进程存活、悬浮球不受影响、任务判 `internal`。

## 8. 不做什么

- 不做自动重启守护者（§14.10 倾向已定：崩溃后用户手动重启 + 上次任务标「中断」）。
- 不为「未来跨平台」建平台抽象层之外的东西（D23：预留 = 接口边界画对，禁止空实现）。
- 不引入任何编排框架/事件总线库；模块间通信 = 显式接口 + channel。

## 9. 未决（沿用 PLAN 登记）

- X/Y 拓扑、五项 RSS 基线、goja ES 能力、WebView2 冷/热拉起、模型常驻内存表 → S0 spike 输出（P1/P11）。
- Go module 组织名（P10）→ `github.com/CarlosShao/wisp` 为占位。
- C18/C20 的包归属（§3 中 `agent/approval`、`agent/scheduler`）为【SPEC】推定，S7 切片卡定稿时确认。
