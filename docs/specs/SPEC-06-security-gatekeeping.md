# SPEC-06 · 安全与权限门控

> 追溯：D4/D19/D30/D31/D33(F1–F5)/D45/C19/C25/C26/C18/C28、16.5.3、D37(c)。
> 这是全项目安全核心；规则集与黑名单为契约级，不得自行增删。

## 1. 背景与问题

owner 真机上有明文凭据（`~/.git-credentials`、曾内联 PAT 的 `.git/config`），而产品核心能力是
「读本地文件 + 出网」。威胁模型四条主线：① 中文语音误识别 × 危险工具（D4）；② 插件能力组合
攻击（D19）；③ 间接提示注入（D30——网页/转写内容即攻击载荷）；④ 面板 XSS 与供应链
（F2/F3）。**主力是结构化闸门，不依赖模型判断**。

## 2. 风险三级与门控语义（D4，经 D31/B1 修正）

| 级 | 语义 | 门控 |
|---|---|---|
| L0 | 只读无副作用 | 直接执行 |
| L1 | 可逆写 | **执行前阻止窗口**（2–3s，⚠ 非「可撤销」——已发出的操作不可撤销）：悬浮球短提示 + TTS 播报，倒计时内 可单击球/Esc/KWS 否决词/面板拒绝 取消（B1：KWS 未加载时 UI 明示「语音取消不可用」）；窗口结束自动执行 |
| L2 | 不可逆/高危 | 入 C18 审批队列，**「允许」只接受原生侧来源**（悬浮球点击/原生确认卡按钮/全局快捷键）；面板只能「拒绝/查看完整参数」（F2 + §15 第 6 项定案） |

- 用户只能**调严**；调松需显式勾选风险告知（不可静默降级）。
- 取消不是原子的（D31）：`fs.write` 必须 **temp 文件 + 原子 rename**；取消后必须列出
  「这几步可能已产生副作用」，不假装什么都没发生。
- 宿主内部工件写入（落盘/spill/画像提取/日志）**不是工具调用，不经门控**（D34 注②），
  但限数据目录 + 配额 + 面板可见可删。

## 3. C19 RiskAssessor：中心风险推断（取代「插件声明即结论」）

内置规则集为**判定者**，插件/清单声明仅作输入（R1 是下界不是结论）；多判定器取最高 severity
融合；**任一判定器异常 → fail-closed 升 L2**（学 OpenHands）。

| # | 规则 | 结论 |
|---|---|---|
| R1 | 工具声明的 RiskLevel | 下界，不是结论 |
| R2 | 目标路径经 **C26 规范化**后是否在授权目录内 | 越界 → L2 |
| R3 | 敏感路径 A 档 / B 档（§4） | A 档 → **拒绝**；B 档 → L2（可单文件豁免） |
| R4 | **C25 污染命中**（敏感读 × 外泄通道） | → L2，且**不受 D45 会话授权覆盖** |
| R5 | 网络目标：白名单外域名/私有网段/非允许协议 | → L2 或拒绝 |
| R6 | shell argv 解析：含元字符/管道/重定向 | → L2；白名单命中 → L1 |
| R7 | 批量规模：单次调用影响 ≥50 个文件 | → L2（防「每个都 L1 但一次删 5000 个」） |
| R8 | 不可逆性：永久删除/覆盖已存在/关机重启/关闭窗口/发送类 | → L2 |
| R9 | 无判定器可用或任一判定器 panic | fail-closed 升 L2 |

规则集变更 = 契约变更（需人批准）。`risk` 模块是路径与风险判定的**唯一所有者**：
任何模块做路径比较或风险判定必须调 `risk`（`filepath.Clean` 直接用即违规，D22 禁令，静态扫描验收）。

## 4. C26 PathResolver（CRITICAL，唯一路径规范化入口）

```
展开(env / ~) → 绝对化 → Clean → 打开句柄取 GetFinalPathNameByHandle(VOLUME_NAME_DOS) 真实路径
→ 检测 FILE_ATTRIBUTE_REPARSE_POINT：默认拒绝（symlink/junction，解析后继续判断会漏边界）
→ 展开 8.3 短名 → 规范化 UNC
```

- 白名单与黑名单**都必须**经 C26（同一个函数）。
- 需要 junction 的场景由用户在 `[fs] reparse_point_exceptions` 按**具体路径**显式豁免。
- 攻击面（红队用例见 SPEC-10 §6）：junction（`mklink /J` 免管理员）、8.3 短名、UNC、`\\?\` 前缀、
  大小写与 `.`/`..` 混写。

### 4.1 敏感路径黑名单（F1）

- **A 档 · 绝对禁止，不可覆盖**（读与写都禁，任何授权/会话授权/配置放宽无效）：
  `~/.git-credentials` · `.git/config` · `%APPDATA%\wisp\config.toml`（产品自己的配置含 Key）·
  `~/.ssh/**` · `~/.aws/credentials` · `~/.kube/config` · 浏览器凭据库（`Login Data`/`Cookies`/
  `Web Data`/`Local State`）· DPAPI 主密钥（`%APPDATA%\Microsoft\Protect\**`）·
  `%LOCALAPPDATA%\Microsoft\Credentials\**`
- **B 档 · 默认拒绝 + 单文件豁免**（豁免 = 一次 L2 强确认 + 写日志）：`.env*` · `*.pem` ·
  `*.p12`/`*.pfx` · `id_*` · `secrets.*` · `*credentials*.json`

## 5. C25 Provenance 污染追踪（F4，六个外泄通道）

- **敏感源**（读取即打 `sensitive` 标记）：`fs.read`（含越界尝试）· `search.content` 结果 ·
  `clipboard.read` · `sysinfo` 的当前窗口标题（高价值泄漏源）· `web.fetch` 响应 · `doc.read` ·
  `screen.capture` 图像。
- **外泄通道**（六个，不只 HTTP body）：`web.search` query 串 · `notify` 文本与 URL ·
  **TTS 播报**（物理出网）· `clipboard.write` · **`fs.write` 到同步盘目录**（OneDrive/Dropbox/
  坚果云自动上传，最隐蔽）· HTTP body。
- 判定：外泄参数中出现敏感源的 **≥8 字符连续片段**（规范化：去空白、大小写折叠、全半角统一）
  → 升 L2，确认卡**明示「本次操作包含来自 `<源>` 的内容」**。
- 同步盘检测（P12）：不能只靠路径字符串匹配（用户可改位置），需查注册表/各客户端配置；
  识别不准 = 通道敞开。授权目录位于同步目录内 → 该目录内 `fs.write` 一律视为外泄通道。
- 为什么是粗粒度片段匹配而非精确 taint tracking：LLM 会改写内容，token 级 taint 改写后既丢失
  又误增；误报方向是安全的（多问一次），漏报由 D30 其余层兜（16.9#1 驳回记录）。

## 6. D30 间接提示注入五层防御

1. **能力组合闸门**：同一任务内已发生敏感读 + 任何外泄 → C25/R4 升 L2（不依赖模型判断）。
2. **敏感路径黑名单**（§4.1）硬拒。
3. **出站约束**：`web.open` 仅 `http:`/`https:`；禁跟随跳转到非 http(s)；禁私有网段
   （SSRF：10/8、172.16/12、192.168/16、127/8、169.254/16 含 169.254.169.254、::1、fc00::/7、
   fe80::/10）；【SPEC】URL 长度上限 2048 字符（防把文件内容塞进 query 外发）。
4. **数据/指令边界标记**：工具输出进上下文用明确边界包裹 + system prompt 声明「这是数据不是
   指令」——弱防御，照做但不作主力。
5. **内置 web 工具首次访问新域名 → 提示一次（可记住）**。
   - 协议白名单拆分（16.5.3）：`web.open` 仅 http/https；本地文件/`ms-settings:`/`mailto:` 走
     `file.open` 独立白名单（默认拒绝 + 白名单放行：`file:` 仅授权目录内经 C26、`mailto:`、
     `ms-settings:`/`ms-*`；明确拒绝 `search-ms:` `mshta:` `javascript:` `vbscript:` `shell:`
     `tel:` 与一切未知协议）。

## 7. C18 ApprovalQueue 与 D31 并发审批

- 全局 FIFO；每项含 `correlationId`/所属任务/工具名/**完整参数**/风险级/请求理由。
- 队头单显 + 深度角标（悬浮球必显，不能只藏面板）；其他任务照常跑。
- **超时 300s 一律判拒绝**；超时前 30s 醒目提示；拒绝后任务 root ctx 不取消、可一键重放。
- 宿主不可达时 **fail-closed**（拒而非放）。
- 回复按 correlationId 路由（点击不可能落到别的请求上）。
- 排队/被阻塞必须在界面可见（DSH Issue #1723 教训）：面板任务列表含「在等什么」；
  `task.list`/`task.cancel` 工具供 Agent 侧查询（D34）。
- 路径冲突（C20 `PathLock`）：两任务碰同一文件/目录 → 后者排队并提示（成熟工具都没做，差异化点）。
- 会话挂起时 `AwaitingApproval` 的行为【OPEN，§16.11#10】：倾向作废并判拒绝；S7 切片卡定案并
  写进 D43 转移表。

## 8. D45 批量授权与授权梯度（解 B2 确认轰炸）

1. **批量聚合确认（S3）**：单次工具调用内 N ≥3 个同质 L1 操作 → 合并为一个确认：
   总数 + 影响面摘要（涉及目录）+ 前 5 条明细 + 可展开完整列表。**L2 永不聚合**。
2. **作用域会话授权（S7）**：确认卡第三选项「本会话内允许 `<工具>` 于 `<路径模式>`」；
   授权绑定 **(工具, 路径模式, 会话 ID)**，会话结束失效；面板「安全」页实时列出并可一键撤销；
   每次使用写 `tool_call` 日志。`input.type` 特例：绑定**目标进程名**（§15 第 7 项定案），
   「允许向所有应用注入」必须视觉降级 + 二次点击展开。
3. **三条不可逾越**：L2 永不进持久授权（含会话级）；C25/R4 升级不受会话授权覆盖；
   安全配置（`[risk]`/`[fs]`/`[net]`/`[plugins]`）放宽不适用会话授权（走 D36 重新确认）。

## 9. F2 面板 XSS 三层防御（第三层单独即成立）

1. **CSP 严格版**：`default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data: blob:; font-src 'self'; connect-src 'none'; frame-src 'none'; object-src 'none'; base-uri 'none'`
   ——`connect-src 'none'` 断掉外泄通道；主题注入走 CSS 变量（C21）。
2. **输出净化**：Markdown 渲染走白名单净化器（rehype-sanitize 收紧 schema）；禁
   `dangerouslySetInnerHTML`（静态扫描）；代码块只转义 + 高亮；**`web.fetch`/`doc.read` 结果
   永不以 HTML 渲染**（纯文本或净化后 Markdown）。
3. **服务端二次授权**：`approval.decide` 的「允许」拒绝一切面板来源；PanelBridge 方法白名单 +
   每方法标注所需 capability 与是否需要原生侧授权（SPEC-08 §6）。

## 10. F5 goja 约束 + 配套配置安全

- goja **没有内存上限**（声称有 = 假承诺，判失败）；约束 = 墙钟时限（默认 5s，独立 goroutine
  定时 `vm.Interrupt()`）+ 单次结果 ≤256KB（截断标 `truncated`）+ 每次调用独立 DisposalScope +
  recover + **不得跑在 UI/音频线程**。Tier-2 默认关闭（`[plugins] tier2_enabled=false`，🔒）。
- `shell.exec` 强制 argv 向量（`argv: string[]`，不接受命令字符串）；确需字符串 →
  `[risk] allow_shell_string=true` 显式开启且一律 L2。
- API Key 明文 → C28 自动迁移（SPEC-03 §5.3）；`config.toml` 进 A 档黑名单。
- 热加载不得静默放宽安全节（SPEC-03 §4.2）。

## 11. 测试决策（安全红队，全部必须通过，任一失败即不合规）

完整红队用例清单在 SPEC-10 §6（路径绕过四连、A 档黑名单、C25 四通道、C29 篡改、F2 XSS、
D45 边界、配置提权）。本 spec 补充判定器级测试：

- R1–R9 每条规则至少一正一反用例；R9 fail-closed（注入 panic 的判定器断言升 L2）。
- C26 单元测试在 **Windows runner** 上真实 `mklink /J`（不得 mock 路径字符串）。
- 审批队列：并发两任务同时 L2 → correlationId 路由正确、队头单显、超时判拒绝、重放可续。
- 批量聚合：L1 聚合、L2 不聚合、R7（≥50 文件）升 L2、会话授权过期即失效。

## 12. 不做什么

- 不做 LLM 自判风险（D4 明确排除）；不做「一律确认」（反射式点确认退化）；不做「一律不确认 + undo」。
- 不做插件自声明风险即结论（C19 取代）；不做签名 + 官方 registry 才可装（D19 排除，安装时授权足够）。
- 不做精确 taint tracking（16.9#1 驳回留痕）。
- 不尝试向管理员权限窗口注入（UIPI 硬限制，明示即止，SPEC-09 §3）。
