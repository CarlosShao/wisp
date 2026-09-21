# 票 91 — OS 级隔离选型备忘录：受限令牌 / AppContainer / 独立低权限账户

**性质**：只读评估。未改任何生产文件与配置；未 commit / push / checkout / stash / 建 worktree。
本文与票面是本票唯一写进仓库的两个文件。
**实验快照目录（仓外）**：`C:\Users\swq\AppData\Local\Temp\wisp91-sess91-a7c3\`
（`spike2\` = 受限令牌装置，`acmd\`、`acmd2\` = AppContainer 装置；原始输出 `battery6.txt`、
`battery7.txt`、`acmd3.txt`、`acmd2-run1.txt`、`acmd2-run2.txt` 都在里面）。
**交回日期**：2026-09-21。

---

## 0. 先给答案（三句话）

1. **我们确实没有 OS 级沙箱**。现有的是 C19 风险分级 + C26 路径解析 + C30 Job Object（治理与度量），
   三者都**不是安全边界**：Job 里的子进程用的是主进程的令牌，照样能读整个用户目录。
2. **票面点名的路 1（受限令牌）在实测中不解决问题**：Windows 上"挡得住用户目录"的唯一可用配方
   会把子进程废掉；所有"子进程还能正常干活"的配方都**挡不住读**，也**挡不住 DPAPI 解密**。
   ⇒ 推荐**不采用**受限令牌（作为"沙箱"来宣传或作为门控的替代品）。
3. **推荐**：本切片只补一条**便宜的 OS 事实层**——把 `shell.exec` / D46 子进程的 DACL 边界交给
   **票 89 已有的 `winsec` 授权方向**（"只有当前用户 SID 能读"），并把 AppContainer 记为
   **RESERVED + 重评触发条件**（不进 S1 实现）。理由见 §6。**代价与不确定项也写在 §6，不粉饰。**

> ⚠ 本节任何"我们有什么"的表述都不使用"内存沙箱"字样（`PLAN.md:592`、`PLAN.md:2490`、
> `SPEC-07-tools-and-plugins.md:114` 明令：出现该说法判失败）。本文把"沙箱"一词只用于
> 讨论**OS 级隔离**，并始终写清它保护什么、不保护什么。

---

## 1. 取证环境（读数的可复现性）

| 事实 | 取证命令 | 读数 |
|---|---|---|
| 当前身份 | `whoami /user` | `desktop-lvs7839\swq` = `S-1-5-21-1228170099-895614386-1166154857-1001` |
| 是否管理员 | `net localgroup Administrators` | `swq` **在** Administrators 里 |
| 实际有效性 | `whoami /groups` | `BUILTIN\Administrators … Group used for deny only`、`Mandatory Label\Medium Mandatory Level S-1-16-8192` ⇒ **UAC 过滤后的中完整性令牌**（不是提权 shell） |
| 完整性级别开关 | PowerShell `Get-CimInstance Win32_DeviceGuard` | `CodeIntegrityPolicyEnforcementStatus=2`、`UsermodeCodeIntegrityPolicyEnforcementStatus=0`（无用户模式 CI 强制） |
| Smart App Control | 注册表 `VerifiedAndReputablePolicyState` | `0`（关闭） |
| Secondary Logon | `sc query seclogon` | `RUNNING` |
| Go 工具链 | `go version` | `go1.27.1 windows/amd64` |
| DPAPI 主密钥位置 | `ls %APPDATA%\Microsoft\Protect` | 只有 `S-1-5-21-…-1001\`（**每个用户 SID 一个目录**）+ `CREDHIST` |

产品数据根（`internal/proc/envfork.go:60,69`）= `%APPDATA%\wisp`（prod）/ `%APPDATA%\wisp-dev`（dev）。
**这台机器上这两个目录当前都是空目录**（`ls -laR` 无文件）⇒ 没有活的 `secrets\` blob 与 `wisp.db` 可取 ACL。
因此 §2 用**同构样本**（同一个用户、同一父目录继承规则下的文件与 DPAPI blob）取证，并在每处标明"这是样本、不是生产数据"。
票 89 的 `winsec` 落地后，生产目录的真实 ACL 应由那张票自己出证据，我不替它写。

---

## 2. AC#2：什么都不加，今天的真实风险是什么（有读数）

### 2.1 同用户的任何进程都能读到我们的秘密 —— 实测

**(a) DPAPI blob 可被"另一个进程"解密。** 用两个互不相干的进程做端到端：

```
$ ./spike2.exe -role=dpapi -dpapi-mode=protect -dpapi-file=...\secrets-test\key.bin
DPAPI protect pid=52428 => 262-byte blob written to C:\Users\swq\AppData\Local\Temp\wisp91-sess91-a7c3\secrets-test\key.bin

$ ./spike2.exe -role=dpapi -dpapi-mode=unprotect -dpapi-file=...\secrets-test\key.bin
DPAPI unprotect pid=32496 => SUCCESS 32 bytes, value="wisp91-dummy-value-not-a-real-key"
```

pid 52428 与 pid 32496 是两个**不同进程**，唯一共同点是同一个用户。C28 用的是 DPAPI CurrentUser
（`PLAN.md:1378` C28 条目），所以这不是"实现有 bug"，而是**该机制的定义域**：
`CryptProtectData` 的密钥派生自**登录用户**，官方 `Requirements` 把它的作用域写为 CurrentUser；
它防的是**别的用户/别的机器**，不防**同一用户里的别的进程**（同用户进程本来就能拿到主密钥
`%APPDATA%\Microsoft\Protect\<SID>\`，见 §1 取证）。

**(b) 文件层的"DACL 只授到用户"= 同用户全部进程可用。** 生产数据根的继承 ACL 原文：

```
C:\Users\swq\AppData\Roaming  NT AUTHORITY\SYSTEM:(I)(OI)(CI)(F)
                              BUILTIN\Administrators:(I)(OI)(CI)(F)
                              DESKTOP-LVS7839\swq:(I)(OI)(CI)(F)
```

（`icacls C:\Users\swq` 同形；票 89 的 `winsec` 未落地前，目录里新建的文件继承的就是这三条。）
NTFS 的 DACL **不按进程区分**：它只问"你的令牌里有没有这个 SID"。所以
**"把 ACL 收成只有我"（票 89 的目标）与"所有以我身份跑的进程都能读"是同一条读数**——
票 89 值得做（它能把"Users/Everyone/其他本地账户"这类**多出来的主体**请出去），
但它**在定义上**挡不住同用户进程。本票与它的分工见 §5。

**(c) 今天谁能读到什么（具体化）：**

| 读到的东西 | 谁 | 证据 |
|---|---|---|
| `secrets\*.bin` 的明文 API Key | **任何以 `swq` 身份运行的进程**：一次被 LLM 拼歪的 `shell.exec`/D46 CLI 子进程、任何计划任务、任何用户级常驻程序（输入法、同步盘、游戏反作弊、浏览器自身的子进程池） | 2.1(a) 实测：换一个 pid 照样解密 |
| `wisp.db`（会话历史、L1 用户画像、任务日志、审计） | 同上（同一个 `%APPDATA%` 继承域；SQLite 只是文件，不认进程） | 1 节 ACL 原文 + 2.1(b) 语义 |
| `config.toml` | 同上 | 同上 |
| 用户整个 home：`~/.git-credentials`、`.git/config`、浏览器 profile、OneDrive 同步盘 | 同上 | `icacls C:\Users\swq` 原文；且 §3 实测：连 `HKCU\Environment` 与 `%APPDATA%` 目录枚举（109 项）都通 |
| **只有提权进程**能读的东西 | System/keylogger 级 | 我们**没有**任何一层能挡（Job Object 只是治理，见 §5） |

⇒ **一句话**：门控（C19/C26/C18）判对了，秘密不外泄；门控被判错一次（提示注入、ASR 听错、
`fs.read` 参数被模型拼歪），泄漏面是**整个用户目录 + 全部凭据 + 全部历史**，
且事后无法归因给"某条边界"，因为那条边界不存在。
今天唯一的刹车是 C25 污染追踪（≥8 字符片段命中升 L2）与 C26 越界拒绝——两者都在**同一层**
（判断该不该允许），都在同一个进程里，**任何一个被判错，另一个不会兜住**。

---

## 3. AC#3：真实读数（受限令牌）

### 3.1 装置

`spike2\main.go`（Go 1.27，只用 stdlib + `advapi32/kernel32/userenv/crypt32`，无第三方依赖）。
它做四件事：

1. 用 `CreateRestrictedToken` 构造 9 种"降权配方"（`priv-only` … `untrusted-il-noreadup`），
   打印每种配方的令牌画像（用户 SID / 完整性标签 / restricting SID 是否存在 / 特权数 / 组数）。
2. 对每种能构造出来的令牌，调 **`CreateProcessAsUserW` 真的起一个子进程**，
   子进程 `-role=probe -recipe=none` 打印自己的令牌 + 执行探针（不做任何二次限制）。
3. 另一条起子进程的路：普通 `CreateProcessW`（继承令牌）起子进程，子进程**自己**构造同一配方、
   `SetTokenInformation(TokenType=TokenImpersonation)` + `SetThreadToken` 把它挂到当前线程
   （`runtime.LockOSThread()`），然后跑探针；跑完 detach 再跑一遍作对照。
   ⇒ 这一条量的是**访问检查的真实结果**（`CreateFile`/`RegOpenKeyEx`/`CryptUnprotectData` 的成败），
   不是"读了文档觉得应该不行"。
4. 探针清单：读 `%USERPROFILE%` 下测试文件、列 `%APPDATA%`、在工作区外建文件、
   读并解密 DPAPI blob（样本，见 §2.1）、开 `HKCU\Environment`、看 `PATH`/`OneDrive`。

测试靶子（实验自建、事后已清理）：
`C:\Users\swq\wisp91-probe-read-sess91a7c3.txt`（内容 `WISP91-SENTINEL-CONTENT-DO-NOT-CAREFULLY`，
`icacls` = SYSTEM/Administrators/`swq` 三条 `(I)(F)`）、
`C:\Users\swq\wisp91-probe-write-sess91a7c3.txt`（工作区外写入靶）、
`C:\Users\Public\wisp91-probe-write-sess91a7c3.txt`（共享目录写入靶）、
`secrets-test\key.bin`（DPAPI 样本 blob）。

### 3.2 关键读数（`battery6.txt` / `battery7.txt`，原始文件在同目录）

| 配方 | 令牌变化 | `CreateProcessAsUserW` | 读 `%USERPROFILE%` 文件 | 列 `%APPDATA%` | 写工作区外 | 写 `C:\Users\Public` | 读+解密 DPAPI | `HKCU\Environment` | 子进程可用？ |
|---|---|---|---|---|---|---|---|---|---|
| `1-priv-only`（`DISABLE_MAX_PRIVILEGE`） | privs 5→1 | **OK** pid=18516 wall=42.9ms | **ALLOWED 41 bytes 原文可见** | ALLOWED 109 项 | ALLOWED | ALLOWED | **ALLOWED + 解出明文** | ALLOWED | 可用（`PATH` 2583 字符、`OneDrive` 变量都在） |
| `2-drop-builtin-groups`（禁 Administrators/Users/Everyone/Authenticated Users） | privs 5→1，`groups=15` 仍列出（TokenGroups 含禁用项） | **OK** pid=8468 | **ALLOWED** | ALLOWED 109 项 | ALLOWED | — | **探针当场炸**：`panic: Failed to load crypt32.dll: winapi error #5` | — | **不可用**：连映射系统 DLL 都被拒 |
| `3-restrict-user-sid`（用户 SID 进 restricting 列表，`Attributes=0`） | `restrictedSids=true` | **OK** pid=41604 | **ALLOWED** | ALLOWED 109 项 | ALLOWED | — | **ALLOWED + 解出明文**（同配方在自受限路线下也 ALLOWED） | ALLOWED | 可用 |
| `4-denyonly-user-sid`（用户 SID 进 restricting 列表且带 `SE_GROUP_USE_FOR_DENY_ONLY`） | — | 未达 | — | — | — | — | — | — | **`CreateRestrictedToken` 直接拒绝：`errno=87 ERROR_INVALID_PARAMETER`**（battery6 行 130、battery7 同位置） |
| `5-disable-user-sid`（用户 SID 进 `SidsToDisable` = 官方"deny-only"路线） | 令牌可建 | **OK** pid=56120，**但子进程立刻废**：`PROBE OpenProcessToken: FAIL errno=5` | **DENIED** | **DENIED** | **DENIED** | ALLOWED | **DENIED（连文件都读不到）** | **DENIED** | **不可用** |
| `6-disable-user+nologon`（5 + 再禁 Logon SID） | 令牌可建 | OK pid=53800（同样废） | **DENIED** | **DENIED** | **DENIED** | ALLOWED | **DENIED** | **DENIED** | 不可用 |
| `7-low-il`（`SetTokenInformation(TokenIntegrityLevel=Low S-1-16-4096)`，默认策略） | IL 成功降到 `S-1-16-4096` | **OK** pid=34368 | **ALLOWED（读照样通）** | ALLOWED | **DENIED**（写被 NoWriteUp 挡） | **DENIED** | **ALLOWED + 解出明文** | ALLOWED | 可用（但读全通） |
| `8/9`（同上再叠 `TokenMandatoryPolicy=NoReadUp|NoWriteUp|NoExecuteUp`） | **`SetTokenInformation` 报"A required privilege is not held by the client"** ⇒ 策略位设不上，IL 仍变 Untrusted(`S-1-16-0`) | OK | **ALLOWED**（策略位没生效） | ALLOWED | **DENIED** | **DENIED** | **ALLOWED + 解出明文** | ALLOWED | 可用（读全通） |

**这张表是本票最重要的实测结论**：

* 能"建令牌 + 起子进程 + 子进程还能正常执行外部 CLI"的三种配方（1/3/7/9），
  **对 `%USERPROFILE%` 的读一律 ALLOWED，对 DPAPI blob 一律"读到 + 解出明文"**。
  也就是说：**受限令牌挡不住 §2 里那三类真实泄漏**。
* 唯一真把用户目录挡住的配方（5/6，官方 `SidsToDisable` 语义 = 给用户 SID 打 deny-only）
  让子进程**连 `OpenProcessToken` 都失败**——一个不能查询自己令牌的进程，谈不上能跑 `git status`。
  配方 2 同理（`crypt32.dll` 映射被拒，进程当场 panic）。
* 中间态（7/9，Low/Untrusted 完整性）**只挡写、不挡读**：`WRITE => DENIED` 与 `READ => ALLOWED` 同框出现。
  想把"读"也挡掉必须设 `TOKEN_MANDATORY_POLICY_NO_READ_UP`，而本机实测
  **`SetTokenInformation` 返回"需要特权未持有"**（`ERROR_PRIVILEGE_NOT_HELD`）⇒
  **非提权进程连这个开关都拧不动**。这条是"受限令牌做不到"的硬理由，不是观点。
* `PATH`/`OneDrive`/`HKCU\Environment`（票面 d 项）：1/3/7/9 全 ALLOWED（体验不破），
  5/6 全 DENIED（体验碎）——**"能用"与"挡住"在受限令牌这条路上是互斥的**。

### 3.3 附带读数：起子进程到底要不要特权、要多久

| 项 | 读数 | 说明 |
|---|---|---|
| `CreateRestrictedToken` 调用耗时 | 0 – 1.4 ms（9 个配方逐个） | 忽略不计 |
| `CreateProcessAsUserW` 用降权令牌起子进程 | **不需要提权、不需要 `SeAssignPrimaryTokenPrivilege`/`SeIncreaseQuotaPrivilege`**，本机全部返回 OK；wall 5.8 – 96.7 ms（同一台机器三次量级抖动，取区间不取单点） | `whoami /priv` 在本环境输出不可靠（见 §3.4），所以这条以"调用返回 OK"为准 |
| 对照：普通 `CreateProcessW` 起同样的 Go 子进程 | 同一区间（16 – 50 ms） | ⇒ **降权本身几乎不增加启动成本**；成本不在这里 |
| 空闲 Go 子进程足迹（`battery7.txt` [4]，3 次） | **私有内存 12.1 – 12.5 MB**、工作集 6.6 – 6.7 MB、**子进程句柄 112 – 118** | Go 运行时基线，不是隔离机制的成本 |
| `cmd.exe` 作 AC 子进程（`acmd2-run1/2.txt`） | 峰值私有 3.1 – 3.2 MB | 系统 shell 的真实量级 |
| 父进程句柄增量 | 一整轮 9×2 子进程 + 9 个令牌跑完后 `delta=+80`（`battery7.txt`）；未 close 的才计入 ⇒ 逐子进程稳态成本 = **2 个句柄 + 1 个令牌句柄** | 对 `D32` 的 `Sleeping` 句柄 <300 无影响（`Sleeping` 定义本就"无子进程"） |

**与 D32 预算的关系（票面 b 项）**：现行口径在 `PLAN.md:2253`
（`Sleeping`：进程树私有内存 ≤25MB(Y)/≤40MB(X)、CPU ≤0.5%、句柄 <300、**无子进程**）。
⚠ 顺带纠一处引用腐坏：**票面引的 `PLAN.md:1032` 不是预算行**——那一行是 D29 的 UI 承载表
（`| **L2 强确认卡** | **WebView** | …`）；预算的旧表在 `PLAN.md:527`（D18 原表，
`PLAN.md:500-516` 已声明被 D32 取代）。后续票别再引 `:1032`。

结论：**隔离机制的开销不在"起进程"上（毫秒级、句柄个位数），而在"要不要为它常驻一个东西"**——
受限令牌不需要任何常驻件（`OpenProcessToken`+`CreateRestrictedToken` 两次调用），
AppContainer 也不需要（一次属性赋值），独立账户需要（登录会话 + profile）。
D32 的 `Sleeping` 态不受影响（工具子进程只在 `Conversation`/工作态存在）。

### 3.4 未证清单（明写，不猜）

* `SetTokenInformation(TokenIntegrityLevel=Low)` 与 `TokenMandatoryPolicy` 在**提权/带 `SeTcbPrivilege`
  的宿主**上的行为——本机非提权，只拿到"需要特权未持有"这一条负向读数。
* **`whoami /priv` 的取值**：这台机器返回的是中文本地化 + 带列宽的输出，我的 `grep` 取证得到空结果，
  所以"该账户是否**声明**了 `SeIncreaseQuotaPrivilege`"这一项**未证**；我以"API 实际返回 OK"为准。
* `ShellExecuteEx` + `runas.exe /trustlevel` 路线未测（票面没要求，且它不是产品可用的 API）。
* 真实生产目录（`%APPDATA%\wisp\secrets`、`wisp.db`）的 ACL 读数：本机数据目录是空的 ⇒ **未取到**；
  这属于票 89 的验收面。

---

## 4. 路 2：AppContainer

（见下文 §4，同一文件续写）
