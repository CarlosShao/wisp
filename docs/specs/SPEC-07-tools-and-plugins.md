# SPEC-07 · 内置工具权威表与插件系统

> 追溯：D3/D14/D19/D34/D46/16.5、C1/C3/C4/C10/C14/C16/C24；风险判定见 SPEC-06 §3。

## 1. 背景与问题

D14 六件套已被 D34 权威表取代（原表散在三处且互相矛盾；按「除 coding 外必须称职」标准严重缺项）。
本 spec 将 D34 表落成实现规格：接口、capability 面、manifest schema、D46 command 插件约束、
goja 宿主 API、插件生命周期。

## 2. 工具接口与能力面

```go
// C1（Pi AgentTool 语义对齐，D21#1）
type Tool interface {
    Name() string          // 全局唯一，'namespace.action'
    Description() string
    Parameters() JSONSchema
    Execute(ctx context.Context, params json.RawMessage,
            onUpdate func(delta string)) (ToolResult, error) // onUpdate 流式进度
}
// C2：L0 只读 / L1 可逆写 / L2 不可逆——判定由 C19 中心推断，声明只是 R1 下界
// C3 capability（11 项）：fs.read fs.write net clipboard shell sysinfo
//                        screen input window notify memory
//                        未声明即拒绝调用（不是报错，是拒绝）
```

- **host bridge 是唯一收口点**：所有能力（内置/插件）的执行、capability 校验、风险门控、审计
  落账都经 `tools` 模块；`plugin` 不得自行判权限（§1.2 边界原则）。
- `ToolProvider`（C4）：`Builtin` / `Manifest`(Tier1) / `Goja`(Tier2) / `[Mcp]` 仅接口位不实现。

## 3. D34 内置工具权威表（唯一来源；落 `docs/TOOLS.md`）

RiskLevel 列是 C19 融合后的**默认结论**；实际判定在调用时算出（越界路径会把 L0 升 L2）。
「落切片」列 = 实现归属。

| 工具 | 动作 | RiskLevel | capability | 落切片 | 备注 |
|---|---|---|---|---|---|
| `fs.read` | 读文本/二进制摘要 | L0（授权内）/ L2（越界） | fs.read | S3 | 经 C26；长输出宿主内部 spill；结果打 taint |
| `fs.list` | 列目录 | L0 / L2（越界） | fs.read | S3 | — |
| `fs.write` | 新建文件 | L1 | fs.write | S3 | temp + 原子 rename（D31） |
| `fs.write` | 覆盖已存在 | L2 | fs.write | S3 | C19/R8 |
| `fs.move` | 移动/重命名 | L1 | fs.write | S3 | 跨盘 = 复制+删 → L2 |
| `fs.trash` | 移入回收站 | **L1** | fs.write | S3 | ⭐ 可用性收益最大：场景①大部分操作 L2→L1，缓解 B2 |
| `fs.delete` | 永久删除 | L2 | fs.write | S3 | **默认不注册**，`[fs] delete_enabled=true` 才有 |
| `file.open` | 系统默认程序打开本地文件/`ms-settings:` | L1 | shell | S3 | 协议白名单见 SPEC-06 §6 |
| `app.launch` | 启动应用（名称/路径/UWP AUMID） | L1 | shell | S3 | ⭐ 一等公民工具，统一 L1（消除原 L0/L1 矛盾） |
| `search.files` | 按名找文件 | L0 | fs.read | S3 | 白名单内 |
| `search.content` | 全文检索 | L0 | fs.read | S3 | 结果打 taint |
| `search.apps` | 检索已安装应用 | L0 | sysinfo | S3 | 只检索不启动 |
| `web.search` | 网页搜索 | L0（query 含 taint → L2） | net | S3 | 实现路径【OPEN：抓结果页 vs API，S3 前】 |
| `web.fetch` | 取网页/JSON→Markdown | L0（私有网段一律拒绝） | net | S3 | 响应打 taint；结果永不以 HTML 渲染 |
| `web.open` | 浏览器打开 URL | L1（仅 http/https） | net | S3 | 非 http(s) 走 `file.open` |
| `doc.read` | PDF/docx/pptx → 文本 | L0 | fs.read | S3 | xlsx 与 OCR DEFERRED |
| `screen.capture` | 截屏（全屏/窗口/区域）→送 LLM | L1 | screen | S3 | 前提 C7 Image 部件；截图只存 `bytes_ref` 不落盘不进日志 DB |
| `input.type` | 向当前焦点窗口注入文本/按键 | L2（首次）→ 可按 D45 会话授权降 L1 | input | S3 | 绑定目标进程名；UIPI 限制见 SPEC-09 |
| `clipboard.read` | 读剪贴板 | L0（打 taint） | clipboard | S3 | — |
| `clipboard.write` | 写剪贴板 | L1 | clipboard | S3 | 算外泄通道（F4） |
| `system.get` | 电量/网络/音量/亮度/时间/焦点窗口标题/输入法 | L0（窗口标题打 taint） | sysinfo | S3 | — |
| `system.set` | 音量/亮度/静音/媒体键/输入法切换 | L1 | sysinfo | S3 | 可逆 |
| `system.power` | 锁屏/睡眠/重启/关机/注销 | L2 | sysinfo | S3 | 关机/重启须二次输入确认词 |
| `window.list` | 列窗口 | L0 | window | S3 | — |
| `window.focus` | 切换/前置窗口 | L1 | window | S3 | — |
| `window.manage` | 最小化/最大化/移动/缩放/分屏 L1；关闭 L2 | L1 / L2 | window | S3 | 关闭可能丢未保存内容 |
| `notify` | 系统通知 | L0（内容含 URL → URL 打 taint） | notify | **S1** | D10 分流依赖；Focus Assist 退化 D42#8 |
| `media.control` | 播放/暂停/上下曲 | L1 | sysinfo | S3 | — |
| `reminder.set/list/cancel` | 一次性提醒（只发通知） | L1 / L0 / L1 | notify | **S4** | D12 窄例外：不做周期性/触发执行/跨重启补发/OS 级调度 |
| `memory.save` / `memory.recall` | 显式记忆 | L1 / L0 | memory | S4 | D20 |
| `task.list` / `task.cancel` | 查看/取消任务 | L0 / L1 | — | **S7** | 「被阻塞的东西必须可见」的 Agent 侧手段 |
| `list_tools` | 元工具 | L0 | — | **S1** | D15②兜底 |
| `shell.exec` | 任意命令 | L2（白名单命中 → L1） | shell | S3 | **默认禁用**；argv 向量强制 |

## 4. 插件 manifest schema（C14/C16，TOML）

```toml
schema_version = 1
id = "com.example.lark"
host_api = "^1.0"                # C16 semver；主版本不兼容 → 拒绝加载并给明确错误，不静默降级
name = "Lark CLI"
version = "0.1.0"

[[tools]]                        # Tier1 声明式工具：参数 schema → host 能力映射
name = "lark.calendar_agenda"
description = "查询飞书日程"
risk = "L0"                      # 仅 R1 下界；实际由 C19 判定
capabilities = ["net"]           # 未声明的 capability → 拒绝调用
net_allowlist = ["open.feishu.cn"]   # net 能力必须带目标域名白名单；白名单外 → 升 L2
parameters = { … JSON Schema … }

[[tools]]                        # D46 command 型：包本机 CLI（S7；能力面扩展主路径）
name = "lark.send_message"
kind = "command"
exe = "C:/Users/me/bin/lark-cli.exe"   # 绝对路径；可变路径（%TEMP% 等）→ 拒绝安装
exe_hash = "sha256:…"            # 安装时钉死；每次执行前校验，不符 → 拒绝执行并告警
argv_template = ["im", "send", "--to", "{{to}}", "--text", "{{text}}"]  # 直接 CreateProcess
                                 # 传 argv 向量，绝不拼接字符串、绝不经 cmd.exe/sh
parameters = { … 每个占位符的 JSON Schema（类型/枚举/范围/最大长度），不符 schema 拒绝进 argv … }
extract = { kind = "json", path = ".data.items[].summary" }  # 声明式提取；原始 stdout 不进上下文
timeout_ms = 15000               # 默认 15000；stdout/stderr 各 ≤1MB 超出 kill 并标 truncated
env_allowlist = ["PATH", "HOME", "USERPROFILE", "LANG"]      # 子进程环境白名单制；
                                 # 绝不传含 API Key 的环境
```

- Tier2（goja JS）插件 manifest 增加 `entry = "main.js"`；其余字段同构。
- 元字符拒绝：参数值含 `` | & ; $ ` < > \n \r `` → 拒绝（即使 argv 向量安全，也防下游 CLI
  自己解释）。
- 输出提取 ≤32KB 进上下文（截断上限）。

## 5. Tier2 goja 运行时（D3/D19④/F5，默认关闭）

- `[plugins] tier2_enabled = false` 为默认；开启属 🔒 安全节变更（L2 重新确认）。
- 约束：墙钟时限默认 5s（独立 goroutine `vm.Interrupt()`）+ 结果 ≤256KB + 每次调用独立
  DisposalScope + recover + 不跑在 UI/音频线程 + 禁 eval/Function 构造器逃逸。
- **不得声称内存沙箱**（goja 无此机制；文档/注释出现该说法 = 假承诺，判失败）。
- 所有能力经 host 注入（goja 无 `fetch`/定时器/IO）；**C24 `GojaHostAPI` 是 JS 侧唯一宿主入口
  全集**，未列出的入口即契约违规（字符串扫描可验）。
  【SPEC 提案，S7 定稿走契约批准】初始入口集：
  `wisp.http.request`（cap: net）· `wisp.fs.read|write`（cap: fs.read|fs.write）·
  `wisp.clipboard.read|write` · `wisp.notify.send` · `wisp.sysinfo.get` ·
  `wisp.state.get|set`（插件自身 KV，无 capability——插件自己的数据）。
- 需要**外部能力**的场景一律优先 Tier1 + D46 包 CLI（零常驻开销、崩溃隔离、无 goja 攻击面）；
  Tier2 只留给真需要进程内计算逻辑的场景（§15 第 11 项定案）。

## 6. 插件生命周期（§14.7）

- **安装来源**：本地路径 / Git URL /（S8 后）registry；授权界面必须显示来源。
- **安装流程**：下载 → manifest schema 校验（C14）→ `host_api` 版本校验（C16）→
  **展示完整能力清单 + 来源 → 用户显式授权** → 落盘 `plugins\<id>\` → 注册 `plugin_state`
  （含 manifest `hash` 与 command 插件 `exe_hash`）。
- **加载时校验**：manifest 哈希不符 → 拒绝加载并告警（防文件被就地替换）；`host_api` 主版本
  不兼容 → 拒绝并明确报错。
- **卸载**：移除注册 + 删文件 + 询问是否删数据（默认保留）；运行中实例经 C11 DisposalScope
  确定性清理。
- **更新**：主版本不兼容 → 拒绝更新并明确报错；更新 = 重走安装授权流程。
- DisposalScope 层级：插件 scope ⊂ 工具 scope ⊂ 任务 scope ⊂ 会话 scope（C11）。

## 7. 测试决策

- 契约测试（§5.2 测试策略）：manifest schema 校验、capability 拒绝行为、**插件自降级 risk 必须被拒**
  （C19/R1：声明 L0 的发送类动作仍判 L2）。
- C14 每字段一组正反例；command 插件红队：`exe` 可变路径拒装、篡改 exe 一字节拒执行、argv 元字符
  注入被拒、环境白名单生效（子进程看不到 API Key——用回显环境变量的测试 CLI 验证）。
- C24 契约测试：JS 侧调用白名单外宿主入口 → 必须拒绝；扫描 goja 绑定表与 C24 文档一致。
- 端到端「真插件跑起来」（S7 判据）：**lark-cli 包装插件**——查日程 L0 直通、发消息 R8 自动 L2
  （插件声明无效）、篡改 exe 拒执行；不另造玩具插件。
- 生命周期：安装→授权→使用→卸载全流程；卸载时运行中实例被确定性清理（无残留 goroutine/子进程）。

## 8. 不做什么

- 不做 WASM 插件运行时（D3 排除）；不做 MCP client（D13/16.9#7 驳回留痕）。
- 不做中心 registry 与强制签名（D19 排除；S8 才重评）。
- 不做插件后台常驻/定时触发（D12：所有插件不常驻；调度交给 OS + `wisp run`）。
- 示例插件可进仓库，但**凭据/令牌绝不进仓库**（README 注明需用户自行安装认证 lark-cli）。
