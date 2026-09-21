# 票 91 · 第二会话（`agent-ticket91b`）— OS 级隔离：三条路的实测判定

**性质**：只读评估。未改任何生产文件与配置；未 push；未在仓内建 worktree。
本文件与 `91-isolation-options.md`（第一会话产物）是本票唯一写进仓库的两个文件，
**本文件只 commit 自己**（第一会话那份由编排者以 `2cac098` 入库）。

**实验快照目录（仓外）**：`C:\Users\swq\AppData\Local\Temp\wisp91b-t2\`
（`git archive HEAD` 于 **HEAD=`15c649f`** 解出；spike 源码与原始输出全在里面：
`spike/{main.go,run5.log,run9.log,run10.log}`、`spike3/{main.go,run_ac2.log,run_ac3.log}`）。
会话后缀 `wisp91b-t2`；第一会话用的是 `wisp91-sess91-a7c3`（两目录互不覆盖）。

**完成度（先看这行）**：
- 路 1 受限令牌：**已实测**（§3，独立复测 + 与第一会话读数对照）
- 路 2 AppContainer：**已实测**（§4 —— 第一会话的 §4 从未写出，本节是补齐）
- 路 3 独立低权限账户：**部分实测**（§5：建号/枚举/`LogonUser`/`CreateProcessWithTokenW` 有读数；
  端到端起一个真子进程**未证**，原因写在 §5.3）
- AC#1 备忘录 ✅ · AC#2 ✅（§6–§7）· AC#3 ✅（§3/§4 真读数）· AC#4 ✅（§9 草案，未开工）

---

## 0. 与第一会话那份备忘录的关系（含一条要纠正的话）

第一会话的 `91-isolation-options.md` 我逐字读了。它**不是零产出**：§1–§3 有装置、有读数、
还自己抓出一处引用腐坏（`PLAN.md:1032`）。它**在 §4 第一行断了**——
`91-isolation-options.md:191` 原文是 `（见下文 §4，同一文件续写）`，下文没有写出来。

但它的 §3.3 里有一句**未经实测就写进结论**的话，本节要用真读数纠正它：

> `91-isolation-options.md:174`：「受限令牌不需要任何常驻件（两次调用），**AppContainer 也不需要（一次属性赋值）**」

"一次属性赋值"在**本机非提权进程上是不成立的**：文档给的属性常量直接失败，
只有旧的 `0x00020009` 变体能起得来（§4.1 三条读数）。这句必须纠，因为**后续实现票会照它写**。

---

## 1. 取证环境（每条都是本机命令读数）

| 事实 | 命令 | 读数 |
|---|---|---|
| 身份 | `spike` 打印 `TokenUser` | `S-1-5-21-1228170099-895614386-1166154857-1001`，完整性 `S-1-16-8192`（Medium） |
| 是否提权 | PowerShell `WindowsPrincipal.IsInRole(Administrator)` | **`IsAdmin=False`** |
| 令牌里有什么特权 | `whoami /priv`（用 `$env:WINDIR\System32\whoami.exe` 全路径，绕开 MSYS 路径转换） | **只有 5 条**：`SeShutdownPrivilege`、`SeChangeNotifyPrivilege`(Enabled)、`SeUndockPrivilege`、`SeIncreaseWorkingSetPrivilege`、`SeTimeZonePrivilege`。⇒ **无 `SeIncreaseQuotaPrivilege`、无 `SeImpersonatePrivilege`、无 `SeDebugPrivilege`、无 `SeCreateGlobalPrivilege`** |
| 管理员会话 | `net session` | `System error 5 … Access is denied.` |
| OS / 工具链 | — | Windows 11 Pro `10.0.26100`；`go1.27.1 windows/amd64` |
| ⚠ 补第一会话的缺口 | 它在 `91-isolation-options.md:181-182` 把「本机是否声明 `SeIncreaseQuotaPrivilege`」记为**未证**（它的 `whoami /priv` grep 取到空） | **本节已证**：这台机器的令牌**确实没有**该特权 ⇒ 后面所有"不需要特权"的读数都建立在"特权缺席"之上，不是"特权在场所以当然成功" |

装置：Go，只用 stdlib + `advapi32/kernel32/userenv/crypt32`（**无第三方依赖**，`GOPROXY=off` 可离线构建）。
父进程构造令牌 → `CreateProcessAsUserW` 起子进程 → 子进程把探针结果写进**父进程给的 stdout 文件句柄**，
父进程读回来打印。探针清单与第一会话同源（读 `%USERPROFILE%` 靶文件、列目录、写工作区外、
开主进程 `PROCESS_VM_READ`、写 `HKCU`、看 `PATH`），另加两条它没有的：**显式授权目录的可写性**与 **DPAPI**。

靶子（本会话自建，§10 有清理命令）：
`C:\Users\swq\wisp91b-t2-secret.txt`（`icacls` = `SYSTEM:(I)(F)` `Administrators:(I)(F)` `swq:(I)(F)`，**无 Everyone**）、
`%APPDATA%\wisp91b-t2\config.toml`、`%TEMP%\wisp91b-t2-parentfile.txt`（同用户独占）、
`%TEMP%\wisp91b-t2-public\`（`icacls /grant *S-1-1-0:(OI)(CI)(M)`，后又 `*S-1-15-2-1:(OI)(CI)(M)`）。

---

## 2. 今天的工具执行路径在哪、以谁的身份跑（我自己 grep 的结论，不受人代述）

**工具是跑在主进程里的，不是子进程。** 三处锚点：

1. 唯一收口点：`internal/tools/bridge.go:99`（`// Bridge is the host bridge: the single choke point every capability flows`）、
   `internal/tools/bridge.go:230`（`// Execute implements agent.ToolProvider - the choke point.`）。
2. 真做 IO 的地方：`internal/tools/fs.go:144` `f, err := os.Open(canon)`（`fs.read`）、
   `internal/tools/fs.go:208`（`fs.list`）；写侧 `internal/tools/fs_write.go:233/368/439/582`（`fs.write`/`fs.trash`/`fs.move`/`fs.delete`）。
   ⇒ 这些调用**用的就是宿主令牌 = 交互式登录用户 `swq` 的令牌**，没有任何 OS 级降权。
3. 全仓 `os/exec` 生产调用点（`grep -rn "exec\.Command" --include="*.go" internal/ cmd/ | grep -v _test`）只有：
   `internal/llm/adaptertest/mockllm.go:65,72`（**测试装置**）、
   `internal/proc/jobscope_windows.go:110`（C30 的 `StartInJob(cmd *exec.Cmd)` 封装本身）、
   `cmd/wisp/slo_windows.go:457`（SLO 采样自测）、`cmd/wisp/doctor.go:275`（`--version` 探针）。
   `shell.exec` 这个名字只出现在 `internal/risk/rules_shell.go`（**判定**）与 config 的未接线清单里，
   `internal/tools/` 下**没有 shell 工具落地**。

⇒ **对本票的直接含义**：票 90 问"某档要不要真的降权跑工具"——今天**没有可降权的工具子进程**。
降权层唯一能挂的对象是**还没建的** `shell.exec`（票 23）与 D46 Tier-1 命令插件（票 50/52）。
今天真要降权，只有两种做法：① 把工具搬出主进程（新增常驻件 ⇒ 撞 D32，见 §7.2）；
② 把**主进程自己**降权起（`fs.read` 与 DPAPI/SQLite/模型一起被降权，第一会话的配方 5/6 读数显示这条路会把进程废掉）。

现有三层的真实分工（都是**应用层判定**，不是 OS 边界）：
`internal/risk/assessor.go`（冻结，C19 中心推断）、`internal/risk/rules_gateway.go`（冻结，R1/R9 fail-closed）、
`internal/risk/pathresolver.go:19-23`（C26：「"Parse-then-continue" on reparse points is forbidden (SPEC-06 §4) …
the resolver denies any path that traverses a reparse point」）。

**⚠ 纠正一处现场说法（引用腐坏候选）**：编排者与我收到的派单都写「C30 **明写**它不是安全边界」。
我把 C30 的两处原文都读了：`docs/PLAN.md:1380` 的 C30 行只写
「Windows Job Object 封装（`KILL_ON_JOB_CLOSE`）：所有子进程入 Job；提供 `TreePrivateBytes()` …」，
`internal/proc/jobscope_windows.go:15-27` 的注释块只讲生命周期与度量口径——
**两处都没有"不是安全边界"这句话**，全仓 `grep -rn "安全边界"` 在 `internal/proc/` 与 C30 契约文本里**零命中**
（命中的是 `PLAN.md:137`、`:1107` 的别的语境，和票 90/91 票面自己）。
⇒ 实质判断我同意（Job Object 不是安全边界，§3/§4 的实测直接证明：Job 里的进程用宿主令牌读走了整个用户目录），
但**"C30 明写"这个归因是错的**，别让它进任何后续票的判据。

---

## 3. 路 1：受限令牌 / 降权（`CreateRestrictedToken` + `CreateProcessAsUserW`）— 实测

### 3.1 令牌构造与起进程，要不要特权？

| 步骤 | 读数 |
|---|---|
| `CreateRestrictedToken(base, DISABLE_MAX_PRIVILEGE\|LUA_TOKEN\|WRITE_RESTRICTED, RestrictedSids={S-1-1-0})` | **OK**；新令牌 `restrictedSids=1`、`privileges=1`（5→1）、`integrity=S-1-16-8192`（**没降**）、`user=` 同一个 SID |
| 同上但只 `DISABLE_MAX_PRIVILEGE` | OK，`restrictedSids=0`、`privileges=1` |
| `CreateProcessAsUserW` 用降权令牌起子进程 | **成功**，**本机无 `SeIncreaseQuotaPrivilege` 也成功**（§1 已证特权缺席）。启动延迟样本（21 次，含普通/降权/AC）：**中位 ~28ms，区间 12.4–266.7ms**（R4 有一次 266ms 抖动，其余 <74ms）⇒ 与不降权的对照组同一量级 ⇒ **降权本身不增加启动成本** |
| `CreateProcessWithTokenW`（同一个降权令牌） | **FAILED errno=1314 `ERROR_PRIVILEGE_NOT_HELD`** ⇒ 这条 API 要特权，本机不可用；**实现票只能走 `CreateProcessAsUserW`** |
| `CreateProcessAsUserW` + `LOGON_WITH_PROFILE=1` + 受限令牌 | 进程**起来了但 15.5s 没有任何输出**，被我 `TerminateProcess`，`exit=1`；两次复现（`run5.log` 15.556s、`run10.log` 15.560s）⇒ **"给子进程加载用户配置文件"这条配方在本机是挂死，不是慢**（票面 d 项的直接答案） |

### 3.2 子进程能不能读 / 能不能写（票面 AC#3 要的读数）

同一次运行里的横向对照（`spike/run9.log`、`run10.log`；`R0` = 不降权对照，`R1` = `WRITE_RESTRICTED`+ restricting SID 只给 `Everyone`，`R3` = 只去特权，`R4` = `R1` 去掉 `LUA_TOKEN`）：

| 探针 | R0 对照 | **R1 降权** | R3 只去特权 | R4 |
|---|---|---|---|---|
| 读 `%USERPROFILE%\wisp91b-t2-secret.txt` | **OK（原文可见）** | **OK（原文可见）** | OK | OK |
| 列 `%USERPROFILE%`（30 项） | OK | **OK** | OK | OK |
| 读 `%APPDATA%\wisp91b-t2\config.toml` | OK | **OK** | OK | OK |
| `stat ~/.git-credentials` | EXISTS | EXISTS | EXISTS | EXISTS |
| 在 `%USERPROFILE%` 新建文件 | OK | **DENIED (5)** | OK | **DENIED (5)** |
| 写 `%TEMP%` 下父进程拥有的文件 | OK | **DENIED (5)** | OK | **DENIED (5)** |
| 写 `D:\`（工作区外） | OK | **DENIED (5)** | OK | **DENIED (5)** |
| 写**显式授权给 Everyone 的目录** | OK | **OK ← 逃逸面** | OK | **OK ← 逃逸面** |
| **`OpenProcess(父, PROCESS_VM_READ)`** | **OK** | **OK** | OK | **OK** |
| 写 `HKCU\Software\wisp91b-t2` | OK | **DENIED (5)** | OK | **DENIED (5)** |
| `PATH` 长度 / `%USERPROFILE%` / `%TEMP%` | 2583 / 有 / 有 | **2583 / 有 / 有（不变）** | 同 | 同 |
| 读 `C:\Windows\System32\cmd.exe` | OK | OK | OK | OK |
| 子进程私有内存（自身 `GetProcessMemoryInfo`） | 1.83 MB | 1.95 MB | — | 1.95 MB |

**判定（这张表就是结论）**：
- `WRITE_RESTRICTED` 只改**写**访问检查，**读完全不受影响**：用户目录、`%APPDATA%`、黑名单文件的 `stat` 全通。
  ⇒ 票 90 若问"某档降权能不能防凭据外泄"：**防不住**——外泄的前半步（读）一步都没挡。
- **降权子进程仍可用 `PROCESS_VM_READ` 打开宿主进程**（R0/R1/R3/R4 全 OK）。
  而 §2 已经证明：今天工具的读、DPAPI 解密、SQLite 页、模型权重**都在宿主内存里**。
  ⇒ 只要秘密还在主进程，降权子进程就是**同房间的另一个进程**，降权形同装饰。
- 唯一能把用户目录挡住的官方配方是给用户 SID 打 deny-only / 禁掉用户 SID，
  **第一会话已实测**：`91-isolation-options.md:135` 配方 5（`SidsToDisable` 含用户 SID）子进程
  **连 `OpenProcessToken` 都 `errno=5`** ⇒ 一个查不了自己令牌的进程跑不了任何 CLI。
  我没重复做这一条（它的装置已在 `wisp91-sess91-a7c3\spike2\`），但它的 §3.2 表格逐项带了 pid 与 errno，可信；
  **标注：这一条是"引第一会话读数"，不是本会话复现**。
- `HKCU` 写被拒（5）但 `PATH` 仍可读（环境变量继承自父进程），OneDrive/`%LOCALAPPDATA%` 下的用户级工具链**没测**（受限令牌下读写分离的行为不干净，测了也不构成推荐依据）。

⇒ **路 1 结论：不采用**。它能给的只有"少一点写权限"，代价是 `shell.exec` 生态半瘫 + 一个"我们降权了"的错误信号。
这与第一会话 §0 第 2 点的结论一致（它从"能用与挡住互斥"的角度得到同一答案）。

---

## 4. 路 2：AppContainer — 实测（补齐第一会话没写出的 §4）

### 4.1 非提权进程能不能走到"起一个 AC 子进程"？能，但只有一条配方

| 步骤 | 读数 |
|---|---|
| `CreateAppContainerProfile("wisp91bt2spike", display, desc, NULL, 0, &sid)` | **HRESULT=0x0，sid=S-1-15-2-1317057317-…** ⇒ **非管理员可创建**（第一会话 `acmd2-run2.txt` 用 `wisp.spike91` 得到同一结论；它那次打印的 `errno=1008` 是**没重置的 `GetLastError` 残值**，函数本身 `HRESULT=0`） |
| 第二次调同名 | `HRESULT=0x800700B7`（`ERROR_ALREADY_EXISTS`）⇒ 配置**持久化在用户档案里**，跨进程/跨重启存在 |
| `DeriveAppContainerSidFromAppContainerName` | HRESULT=0，返回同一个 SID |
| `GetAppContainerFolderPath(sid)` | **`HRESULT=0x80070539`（`ERROR_INVALID_SID`）** ⇒ 手工建的容器**没有落地配置目录**（后果见 §4.3 的 `%TEMP%`） |
| `GetAppContainerToken` | **`advapi32.dll` / `kernel32.dll` / `userenv.dll` / `api-ms-win-security-base-l1-2-2.dll` 全部 NOT EXPORTED** ⇒ 这条 MSDN 文档路线在本机 Windows 11 26100 上取不到符号（"文档说有、导入表没有"） |
| `CreateProcessW` + `PROC_THREAD_ATTRIBUTE_SECURITY_CAPABILITY`（文档值 `0x00020000`，载荷 `SECURITY_CAPABILITY{Length,Sid,Caps,NumCaps}`，Length=32 **和** 28 都试了） | **`UpdateProcThreadAttribute` FAILED errno=24 `ERROR_BAD_LENGTH`** ⇒ 文档给的那条在本机不工作 |
| `CreateProcessW` + `PROC_THREAD_ATTRIBUTE_ALL_ASSIGN_SECURITY_CAPABILITY`（`0x00020009`，载荷**无 `Length`** 字段：`{Sid,Caps,NumCaps,Reserved}`）+ `EXTENDED_STARTUPINFO_PRESENT` | **OK，pid=40576，spawn=18.2ms**（另一轮 24.7ms）⇒ **唯一可用配方**；子进程自报 `token.isAppContainer=1`、`appContainerSid=S-1-15-2-1317…`、`integrity=S-1-16-4096`（**Low**） |

⚠ 给实现票的硬信息：**`0x00020009` 在头文件里被标注为"不要用"的旧常量**。
能用的恰好是被废弃的那个，能跑的那个不被官方文档推荐——这个矛盾**必须**在实现票里正面处理
（钉 Win11 26100 实测可用 + 建 CI 门；不能写成"按 MSDN 用 0x00020000"，那样第一天就红）。

### 4.2 AC 子进程的读写读数（与 §3.2 同格式对照）

`spike3/run_ac3.log`；父进程**事先** `icacls` 了两样东西：sink 目录 `*S-1-15-2-1:(OI)(CI)(M)`、子 exe `*S-1-15-2-1:(R)`。

| 探针 | 受限令牌 R1（§3.2） | **AppContainer** |
|---|---|---|
| 读 `%USERPROFILE%` 靶文件 | OK（挡不住） | **DENIED (Access is denied)** |
| 列 `%USERPROFILE%` | OK | **DENIED** |
| 读 `%APPDATA%\…\config.toml` | OK | **DENIED** |
| `stat ~/.git-credentials` | EXISTS | **"absent"**（见下 caveat） |
| `%USERPROFILE%` 新建 / 写 | DENIED | **DENIED** |
| 写 `D:\`（工作区外） | DENIED | **DENIED** |
| **`OpenProcess(宿主, PROCESS_VM_READ)`** | **OK（挡不住）** | **DENIED（`OpenProcess` 返回 FALSE）** |
| 写 `HKCU\Software\…` | DENIED (5) | **DENIED (5)** |
| **写"显式授权给 `ALL APPLICATION PACKAGES` 的目录"** | OK | **OK + 回读成功** |
| 读 `C:\Windows\System32\cmd.exe` | OK | **OK**（系统目录默认可读可执行 ⇒ shell 本身还能起） |
| `stat %LOCALAPPDATA%\Programs`（用户级工具链） | OK | **DENIED** |
| `PATH` / `%USERPROFILE%` | 2583 / 有 | **2583 / 有**（环境变量继承，没被剥） |
| 子进程私有内存 | 1.95 MB | **1.95 MB**（Go 探针二进制；`cmd.exe` 量级见第一会话 `acmd2-run2.txt` 峰值 3.1–3.2 MB） |

**caveat（我自己的装置缺陷，不许当成收益）**：`STAT_home_dotgitcredentials = absent` 这一行
**不能读成"AC 让黑名单文件消失了"**——我的探针把"访问被拒"与"不存在"合并成同一个输出。
真说法是：**AC 子进程无法判定该文件是否存在**（枚举 `%USERPROFILE%` 已 DENIED）。存在性是**变弱了**的旁道，
但本装置无法区分"拒绝"和"没有" ⇒ **未证**。

### 4.3 代价（票面 b/d 两项，都有读数）

- **成本不在启动**：AC 18–25ms，与不降权同量级；子进程瞬时私有 ~1.9MB。
- **成本在"每一个可访问路径都要显式授权"**：默认**全拒**（§4.2），要给工作区就得对**具体路径**下 ACL。
  好消息：这与 C26 的 `[fs] allowed_dirs` 是**同一份清单**（`internal/tools/paths.go` 的 `PathCanonicalizer` +
  `internal/risk/pathresolver.go` 的解析结果）⇒ 授权面可以从已有配置**推导**，不需要新真相源。
  坏消息：授权是**内核对象 ACL，进程退出不会自己收**（我 `icacls` 加的 ACE 现在还在 `%TEMP%` 上，§10 手工清理），
  所以"给完就撤"必须实现票自己做（`SetFileSecurity` 快照/还原），并处理崩溃残留 ——
  **这条与票 73（孤儿暂存文件清扫）是同一类问题，不要另起炉灶。**
- **UX 会破的具体三处**（实测）：① `%TEMP%` 被重定向到 `…\Local\Packages\wisp91bt2spike\AC\Temp`，
  而该目录**不存在**（`GetAppContainerFolderPath` 已 `ERROR_INVALID_SID`）⇒ 任何用临时文件的 CLI 当场失败；
  ② `%LOCALAPPDATA%\Programs` 不可读 ⇒ 装在用户目录里的 `node`/`npm`/`git` 类工具找不到；
  ③ `HKCU` 写被拒 ⇒ 会写注册表的 CLI 失败。**三处都必须有"明确错误、不许静默"的文案**（票面 d 项要求）。
- **D32 的关系（票面 b 项）**：`docs/PLAN.md:2253` 的 `Sleeping` 行是硬口径——
  进程树私有内存 ≤25MB(Y)/≤40MB(X)、CPU ≤0.5%、句柄 <300、**"无子进程"**。
  ⇒ 走"**宿主内按需起 AC 子进程、用完即退**"**不触碰**休眠预算（休眠时一个子进程都没有）。
  ⇒ 但如果为了实现这条路而**常驻一个降权 executor**，就直接撞 `Sleeping … 无子进程` 这条附加硬约束
  （数字与措辞一字不许动 ⇒ **常驻方案出局**）。
  顺带纠一处引用：票面把 D32 写成 `PLAN.md:527`、`:1032`——`:527` 是 D18 旧表且 `PLAN.md:500` 已声明
  「本节的 SLO 表与延迟表已被 D32 整体取代」，`:1032` 是 D29 的 UI 承载表（**跟内存预算无关**，第一会话也纠过）。
  权威行号 = **`PLAN.md:2253`**。

### 4.4 路 2 判定

**方向和路 1 相反：AppContainer 是本机唯一"读也挡、写也挡、宿主内存也挡"的一层。**
但**它今天没有可挂的对象**（§2：工具在主进程里，没有子进程），并且它有 4 个未解的工程债：
废弃常量、无配置目录、每路径 ACL 与撤销、以及 §4.2 caveat 那类"看起来更安全"的读数陷阱。
⇒ 本票**不推荐在 S1 落地**，推荐按第一会话 §0 第 3 点记为 **RESERVED + 重评触发条件**（触发条件见 §9.3，我把它写成可判定的）。

---

## 5. 路 3：独立低权限本地账户 — 部分实测

### 5.1 有读数的部分

| 动作 | 读数 | 含义 |
|---|---|---|
| `net user wisp91b-test /add`（本会话自己发起） | **`System error 5 … Access is denied.`** | **建号需要一次性提权**（owner 手动做一次，或安装器做一次） |
| `Get-LocalUser` | 非提权可枚举：`Administrator enabled=False`、`CarlosShao True`、`CodexSandboxOffline True`、`CodexSandboxOnline True`、`DefaultAccount False`、`Guest True` | 这台机器上**已经存在**第三方建的沙箱账户；我们没有它们的口令 |
| `LogonUserW("Guest"/"WDAGUtilityAccount"/不存在用户, 错口令, LOGON_INTERACTIVE)` | 三个都 **errno=1326 `ERROR_LOGON_FAILURE`** | **不是 1314 特权不足** ⇒ `LogonUser` 对非提权进程**可用**；挡路的是**凭据**，不是 API |
| `runas` 依赖的 Secondary Logon | `seclogon Status=Running StartType=Manual` | 服务在跑（按需启动）⇒ `runas` 这条路今天可用 |
| `CreateProcessWithTokenW` | **errno=1314 `ERROR_PRIVILEGE_NOT_HELD`**（§3.1，同一个令牌实测） | 用它起"别的令牌"的进程在本机不可用 |
| `CreateProcessAsUserW`（换到别的用户）需要 `SeAssignPrimaryTokenPrivilege`/`SeIncreaseQuotaPrivilege` | 本机令牌里**这两条都没有**（§1 `whoami /priv` 全量 5 条） | 跨用户起进程这一步预期失败；**⚠ 未直接证**——我没有第二账户的令牌可试 |

**安全纪律说明**：我**故意没有**用真实交互账户 `swq` + 错口令去测 `LogonUserW`，因为失败计数可能触发
本地"账户锁定阈值"，那会把 owner 的机器锁在外面。上面三条只用内置/不存在的账户名。

### 5.2 判定（为什么它最重）

- 边界最硬：另一个登录会话，另一个 `%APPDATA%\Microsoft\Protect\<SID>`，DPAPI 主密钥天然不同域
  （这条**目录结构证据**是第一会话取的：`91-isolation-options.md:40`——`Protect` 下只有本用户 SID 一个目录）。
- 但产品要付三样钱：① **一次性提权建号**；② 产品必须**保管这个账户的口令**——于是多出一个
  "能登录本机的凭据"要保护，而它**恰好落在 C28/票 89 要保护的同一类秘密里**（D33 域），
  形成循环：为防泄漏引入一个新泄漏面；③ 每次执行工具都要**跨会话**：文件要能被两个账户同时访问 ⇒
  工作区要么放共享盘 ACL（`Users\Public` 那种，等于对外开放），要么每个路径做双 SID 授权 + 撤销。
- 加上 §5.1 最后两行的特权缺席：**本机连"起进程"这一步都没被证成**。

⇒ **路 3：不推荐**（不是"做不到"，是"收益/新增攻击面/新持久依赖"三条都不划算，且**本机端到端未证**）。

---

## 6. 票面 a–e 五项判定（一次列全，路 1/2/3 各一列）

| 项 | 路 1 受限令牌 | 路 2 AppContainer | 路 3 独立账户 |
|---|---|---|---|
| **a. 与 cgo 主进程的关系**（保护什么/不保护什么） | 只降子进程。保护：**极少**（§3.2 只有写被拒）。不保护：宿主内存（实测 `VM_READ` OK）、所有读、DPAPI | 只降子进程。保护：**该子进程能看到的一切**（读/写/宿主内存/HKCU 实测全 DENIED）。不保护：宿主自身——**§2 已证今天工具就在宿主里** ⇒ 不搬出宿主就等于没有 | 同上；边界最硬（登录会话不同域），但宿主不受保护这一点**完全一样** |
| **b. 与 D32 预算**（`PLAN.md:2253`，数字一字不动） | 无常驻件；按需子进程 12–74ms、瞬时私有 ~1.9MB ⇒ `Sleeping` 不受影响 | 无常驻件；18–25ms、~1.9MB ⇒ 同上。**但"常驻 executor"方案违反 `Sleeping … 无子进程` ⇒ 出局** | 登录会话 + profile 常驻；本机未测出数（建不了号）⇒ **未证**，且它天然需要至少一个跨会话常驻件，与 D32 冲突面最大 |
| **c. 与 C26/C19 的分工，谁兜底**（本票最重要的产出） | **不能当兜底**：读不挡 ⇒ 判错一次照样泄。**分工**：C19/C26 = 判定层（该不该），限令牌只是"少一点写"，不改变判定被绕过的后果 | **唯一能当兜底的候选**：OS 回答"能不能"，与 C19/C26 的"该不该"**不同层**，判定错了也读不到。授权表可从 C26 的 `allowed_dirs` 推导（同一真相源）；撤销与残留归实现票 | 同 b，但每路径要做**双身份**授权 + 登录会话管理；判定层错了也读不到（最硬） |
| **d. 用户体验** | 轻：`PATH`/`%TEMP%`/`%USERPROFILE%` 实测不变；**但 `LOGON_WITH_PROFILE` + 受限令牌 = 15.5s 无输出挂死（两次复现）** | 重：`%TEMP%` 指向不存在的目录、`%LOCALAPPDATA%\Programs` 不可读、HKCU 写被拒 ⇒ 三处会静悄悄失败 | 最重：跨会话、无交互桌面（UIPI/窗口站）、profile 首次加载 |
| **e. 与既有红线** | **安全**（不会成为第三条通道），因为它是 no-op。**但危险在宣传**：说"降权了"会让 owner 以为门控判错也没事，那是对 D4/D19 的**腐蚀**，与 B2 的"确认疲劳"同型 | **必须钉三条**：① 降权**不是**授权——不许任何档以"已降权"为由跳 L2；② 任何授权不得撤销 C26 黑名单的 ACL（`PLAN.md:1629` 原文就是这条，见 §8.3）；③ AC 子进程失败必须显式报错（票面 d）。它**不产生**新的批准通道，`approval.decide` 仍只吃原生侧（`PLAN.md:1588`）；起进程仍受 C18 的 300s 超时约束（`PLAN.md:3143` 驳回无限等待）⇒ **不会成第三条通道** | 同 e2；额外风险：产品保管一个可登录本机的口令，直接进 D33 面 |

---

## 7. AC#2：什么都不加，今天的真实风险是什么（具体到"谁能读到什么"）

### 7.1 直接引用第一会话已证的部分

`91-isolation-options.md:53-95` 的那三条我逐字读了：
① **DPAPI blob 换进程照样解密**（两个 pid 52428/32496 的端到端读数，`:56-60`）；
② `%APPDATA%` 的继承 ACL 只有 `SYSTEM/Administrators/swq` 三条 ⇒ NTFS DACL 不认"哪个进程"；
③ 今天任何以 `swq` 身份跑的东西（计划任务、输入法、同步盘、反作弊、浏览器子进程池）都能读
`secrets\`、`wisp.db`、`config.toml`、`~/.git-credentials`。
⇒ **本会话的增量证明**：§3.2 的 `OpenProcess(宿主, PROCESS_VM_READ)=OK` 说明
**连"进程边界"都拦不住**——不需要读文件，直接读宿主内存就能拿到已解密的凭据与已加载的模型。
这是对 ① 的**加强**：泄漏不需要碰 DPAPI 本身。

### 7.2 我要往里加的一条（不写就是漏报）

**"什么都不加"的准确说法是：什么都不加的代价 = 判定层就是唯一层。**
C19 规则集（`assessor.go`/`rules_gateway.go`，冻结）、C26 路径解析、C25 污染追踪**都在同一个进程、同一个令牌里**，
所以它们共享一个故障域：一次提示注入把 `fs.read` 的参数拼歪，三层同时失效——
**不是三层各失守一次，是一层失守即三层失守**。OS 层的唯一不可替代价值就是"不同故障域"，
而 §3 证明受限令牌提供不了它（读不挡），只有 §4 能。

### 7.3 与"什么都不加"等价的第三个选项，必须说破

如果只做 §3（受限令牌）并对外说"工具有沙箱"，那**比不做更糟**：
它给了一个假的边界信号，而 `PLAN.md:592`/`:2490`/`SPEC-07-tools-and-plugins.md:114` 对"假承诺"的判法
（"文档/注释里出现该说法判失败"）说明**这个仓库已经把'假沙箱'当成红线在管**。
⇒ 本票的措辞遵守同一条：全文**不使用**"内存沙箱"，"沙箱"只用于讨论 OS 隔离并始终写清不保护什么。

---

## 8. 引用核对结果（我自己一行行开的；腐坏逐条列）

| 引用 | 打开后的原文 | 判定 |
|---|---|---|
| `PLAN.md:527` | `| **空闲态**（KWS 关） | ≤ 25MB | ≤ 0.5%（1min 均值）| …` | **内容对但非权威**：`PLAN.md:500` 写「本节的 SLO 表…已被 D32 整体取代」 |
| `PLAN.md:1032` | `\| **L2 强确认卡** \| **WebView** \| …`（D29 UI 承载表） | **腐坏**（票面把它当 D32 预算行）。第一会话已记，本会话复核确认 |
| 权威 D32 行 | `PLAN.md:2253` = `\| Sleeping（空闲，KWS 关，无子进程）\| ≤ 25MB / ≤ 40MB \| ≤ 0.5% \| …句柄 <300 · GDI <200 · goroutine ≤6 \|` | ✅ 后续票引这一行 |
| `PLAN.md:592` | `(b) goja 没有内存上限机制，只有 vm.Interrupt()。声称有内存沙箱是假承诺` | ✅ |
| `PLAN.md:2486` / `:2490` | `### F5（HIGH）goja 沙箱是假承诺 + 同进程崩溃面` / `1. **不得声称有内存沙箱**。` | ✅ |
| `PLAN.md:1380` | C30 `JobScope` 行（`KILL_ON_JOB_CLOSE`、`TreePrivateBytes()`、`D32/D38`） | ✅ 行号对；但**没有"不是安全边界"这句** ⇒ 见 §2 末的纠正 |
| `SPEC-07-tools-plugins.md:114`（票面写法） | 真实文件名是 `docs/specs/SPEC-07-tools-and-plugins.md:114`，内容 `- **不得声称内存沙箱**（goja 无此机制；…）` | **文件名腐坏**（内容/行号对）；本票按真名引 |
| `PLAN.md:1588` | D33/F2 面板 XSS 行，句内含 `approval.decide 的"允许"只接受原生侧` | ✅（票 90 用法准确） |
| `PLAN.md:1629` | `- （新增）A 档黑名单：读 ~/.git-credentials / .git/config / %APPDATA%\wisp\config.toml → 必须拒绝，且任何授权与会话授权都无法覆盖` | ⚠ **释义偏差**：这一行讲的是**黑名单文件读取**，不是"不可逆操作"。票 90 的红线之一挂错锚点；"不可逆操作"的真锚是 `PLAN.md:123` + `internal/risk/rules_irreversible.go:5-14`（R8 类集，SPEC-06 §3） |
| `PLAN.md:3143` | §16.9 第 4 条：`L2 审批超时应设为无限等待 … **驳回**` | ✅ |
| `PLAN.md:2030` | 决策 11：Tier-2 goja 默认关闭（含「不得声称有内存沙箱」） | ✅ |
| `internal/proc/envfork.go:60,69` | `:60` = prod `DataDir = %APPDATA%\wisp`，`:69` = dev `wisp-dev` | ✅（第一会话引用准确） |
| `docs/reports/pending-and-issues.md:1820` | 「…OS 级沙箱没有（只有 Job Object，而 C30 明写它不是安全边界）」 | ⚠ 归因不成立（同 §2）；**账目不改**，但**票 91 的这句结论要落成本文件 §2 的表述** |

---

## 9. AC#4：如果落地，要改哪些契约 + 后续实现票判据草稿

**⚠ 这一节是草案，属于 D22，我没有动任何冻结文件；owner 批准前不许开工。**

### 9.1 契约草案

- **新增 `C33 ProcessIsolation`**（建议进 `PLAN.md` C 表，紧跟 C30）：
  > 职责：**只**为"会执行外部代码"的子进程（`shell.exec`、D46 Tier-1 命令插件）提供 OS 能力边界；
  > 接口 = `Start(cmd, allowedPaths) -> proc`，内部 = 一个**按需创建、用完即退**的 AppContainer +
  > 对该次 `allowedPaths` 的**临时 ACL 授予与撤销**。
  > 边界：① **C33 不是授权**——C19/C26/C18 的判定与门控**一字不让**，C33 只在判定通过后收窄能力；
  > ② 宿主自身与主进程内工具（今天的 `fs.*`）**不在 C33 保护范围**，须明写；
  > ③ 失败**必须**返回可判别的错误（`WISP_ERR_ISOLATION_GRANT_DENIED` 等），**禁止**静默降级为"不降权跑"；
  > ④ 授权表**只能从 C26 的 `allowed_dirs` 推导**，不得新增真相源；
  > ⑤ 崩溃残留的 ACL 由**票 73 的清扫器**同一机制回收（不另造）。
- **修正 C30 的描述一句**（现在 `PLAN.md:1380` 只讲生命周期与度量）：加
  「**Job Object 不是安全边界**：Job 内进程使用宿主令牌，C30 不得被引用为隔离手段」。
  ⇒ 这是把票面上已经在传的话**写进契约**，同时把 §2 的"明写"错误就地终结。
- **SPEC-06 加一节"两层分工"**：判定层（C19/C26/C25，同故障域）/ 能力层（C33，不同故障域）；
  **兜底者是能力层**，且**能力层永远不成为批准通道**。
- **D32 不改**（数字与 `Sleeping … 无子进程` 是验收判据）。C33 的存在性判据必须引用 `PLAN.md:2253`，
  并**显式禁止**为隔离引入常驻子进程。

### 9.2 票 90 的直接答案（M1 要不要含"降权跑工具"这一维）

**不要。** 三档（每步都问 / 只问高危 / 全自动，owner 已定 R20）改的是**问不问**，
不引入"降权"第四维；理由：§3 证明降权提供不了任何可宣传的保证，而 §4 证明真正的那一层**今天没有可挂的对象**（§2）。
**唯一要落进票 90 的一条**：**明写"权限模式不改变 OS 能力"**，防止后续有人把"全自动"实现成"顺手降个权"
（那是 §7.3 的假承诺）。这条我建议在票 90 的 AC 里加半句，**不需要**新决策。

### 9.3 后续实现票（建议编号 95）判据草稿 —— 含**重评触发条件**

- **AC#0（触发门，先判这条）**：`shell.exec` 或任一 D46 Tier-1 插件**在生产路径真的会起子进程**
  （判据 = `grep -rn "exec.Command" internal/tools/ internal/plugin/` 出现**非测试**命中且被 bridge 调到）。
  不满足 ⇒ 票 95 保持 blocked，**不许用"先在主进程里做 ACL"绕过**。
- **AC#1 能力实测可复现**：把本文件 §4.1 的**两条否定 + 一条肯定**做成 CI 用例
  （0x00020000 → 必须证成 `ERROR_BAD_LENGTH` 或更好；0x00020009 → 必须起得来且子进程自报 `isAppContainer=1`）。
  **不许**只写"按 MSDN"。
- **AC#2 读边界**：AC 子进程读 `%USERPROFILE%` 靶文件、读 `%APPDATA%`、`OpenProcess(宿主, VM_READ)` **三条都必须 DENIED**，
  探针输出用本文件的 `PROBE k=v` 格式（可被日志判定，不靠人眼）。
- **AC#3 授权与撤销**：给定 `allowed_dirs`，子进程**只**能读写清单内路径；
  进程退出（含被 kill）后 ACL **必须**回到快照态——用例要包含 `TerminateProcess` 与主进程被 Job 连带清理
  （C30 `KILL_ON_JOB_CLOSE`）**两条**崩溃路径。
- **AC#4 双向变异**：① 去掉撤销 ⇒ 必须有用例红；② 把"授权失败"改成"当作没降权继续跑" ⇒ 必须有用例红；
  ③ 把 §4.2 的 DENIED 断言改成"跳过不存在的路径" ⇒ 必须有用例红（防 §4.2 caveat 那类假绿）。
- **AC#5 体验面**：§4.3 三处（`%TEMP%`、`%LOCALAPPDATA%`、HKCU）各有**明确错误文案**用例，
  **禁止**"看起来成功了"。
- **AC#6 D32 无回归**：`Sleeping` 句柄/进程数采样与 `PLAN.md:2253` 两档数字一字不动地跑；
  另测"按需拉起 + 用完退出"下的 `Conversation` 峰值不越过同表工作态上限。

---

## 10. 未证清单 + 现场清理

**未证（明写，不当结论用）**：
1. **AC 子进程能否解密本用户的 DPAPI blob**——本会话装置有缺陷：`CryptProtectData` 产出的 blob
   不以 `DPAPI` 魔数开头（base64 头是 `QEmIGsEB…`），**连控制组 R0 自解都返回 errno=13 `ERROR_INVALID_DATA`**
   ⇒ 我这组 `DPAPI_unprotect` 读数**全部作废，不采信**（AC 那行 DENIED 同样是垃圾进垃圾出）。
   只能给推断：AC 读 `%APPDATA%` 已实测 DENIED，而主密钥在 `%APPDATA%\Microsoft\Protect\<SID>\`
   ⇒ **推断为解不开，未证**。要证就用第一会话 `wisp91-sess91-a7c3\spike2\` 那套已跑通的 DPAPI 装置重跑。
2. 跨用户 `CreateProcessAsUserW`（路 3 的起进程步骤）：本机无第二账户令牌 ⇒ 未证（依据见 §5.1 末行）。
3. `TOKEN_MANDATORY_POLICY_NO_READ_UP`（把读也挡掉的 Low IL 配方）：第一会话 `:149-151` 拿到的是
   「`SetTokenInformation` 报需要特权未持有」——**非提权环境的负向读数**；提权宿主上的行为未证。
4. Windows 10 / 其他 Windows 11 build 上 `0x00020009` 是否仍可用：未证 ⇒ 票 95 AC#1 必须是**跑出来的门**，
   不能是本文件这种一次性手工读数。
5. `wisp` 产品数据根（`%APPDATA%\wisp`、`wisp-dev`）的真实 ACL：本机这两目录当前为空 ⇒ 归票 89。

**现场清理：本会话收尾已执行并验证**（留档如下，便于后来人复现同样的取证靶）：

```
icacls %TEMP%\wisp91b-t2-public  /remove *S-1-1-0 *S-1-15-2-1     # 撤掉 Everyone / AC 授权
icacls …\wisp91b-t2\spike\spike91b.exe /remove *S-1-15-2-1
del %TEMP%\wisp91b-t2-public\* ; del %TEMP%\wisp91b-t2-parentfile.txt
del C:\Users\swq\wisp91b-t2-secret.txt ; rmdir %APPDATA%\wisp91b-t2
reg delete HKCU\Software\wisp91b-t2 /f
userenv!DeleteAppContainerProfile("wisp91bt2spike")
```

验证读数：四个靶 `Test-Path` 全 `False`（含 `D:\wisp91b-t2-outside.txt`——控制组自己写了又自己删了）、
`DeleteAppContainerProfile` 返回 `0`（S_OK）。**保留**的只有仓外快照目录
`%TEMP%\wisp91b-t2\{spike,spike2,spike3}\`（源码 + `run*.log`/`run_ac*.log` 原始输出），
因为它是本文所有读数的出处；快照树 `git archive` 出来的副本与仓无关，可整目录删。

