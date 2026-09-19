# SPEC-09 · Windows 平台行为与可靠性

> 追溯：§6/D38/D41/D42（12 项）/§16.11#2#4、D43；平台异常态不得静默失败（§2 末强制）。

## 1. 背景与问题

常驻 Win32 + WASAPI + WebView2 应用死于平台细节：电源事件、设备热插拔、显示器拓扑、时钟跳变、
Focus Assist。D42 登记了 12 项可靠性事件，本 spec 给出每项的检测与应对实现约定。

## 2. 生命周期（§6）

| 项 | 决策 |
|---|---|
| 单实例 | 命名互斥体 **per-session**：`Local\wisp-single-instance`（env 分叉见 SPEC-03 §5.2）；二次启动 → 激活已有实例的悬浮球并退出。global 互斥会让第二个登录用户无法启动——真实 bug 源 |
| 开机自启 | 默认关，opt-in（注册表 Run 键，D1「不用不占资源」） |
| 崩溃恢复 | 重启回 `Sleeping`；上次未完成任务标记「中断」面板可见；**不自动重跑**（可能已产生副作用）；无自动重启守护者（§14.10 倾向定案） |
| 自动更新 | 默认手动检查；minisign 校验（D17）；更新运行中 exe 的完整流程见 SPEC-11 §7 |
| WebView2 Runtime | Win11 预装、Win10 多数已装；缺失 → 无面板模式 + 引导安装，不得崩溃（D42#11） |
| 多显示器 DPI | Per-Monitor V2；位置按显示器保存；拔插自动回可见区 |
| 全屏抑制 | 前台全屏应用自动隐藏悬浮球（可关） |
| 麦克风被占/被拒 | `Error(audio_device)` 明确文案，不得静默失败 |
| 语言 | 中文优先；UI 字符串从 S1 起外置（i18n 本体 DEFERRED） |

## 3. 平台硬限制（必须明示，不得尝试绕过）

- **UIPI**：`input.type` 无法向管理员权限窗口注入，且无解。错误文案必须明确：
  「目标程序以管理员身份运行，Wisp 无法向其输入」；文档写清并引导用户以普通权限运行目标应用。
- **DRM 窗口截屏黑屏**（Netflix 等）：`screen.capture` 平台行为，明示即可不修。
- **Focus Assist / Win11 勿扰无稳定公开 API**（D42#8）：拿不到状态就**无条件走退化路径**——
  长结果总是同时更新悬浮球角标 + 面板任务列表，不依赖系统通知可达。

## 4. 12 项可靠性事件（D42，每项须真实触发验证——拔设备、断网、填磁盘、改系统时间）

| # | 事件 | 检测 | 应对 | 落切片 |
|---|---|---|---|---|
| 1 | 睡眠/唤醒 | `WM_POWERBROADCAST`（PBT_APMSUSPEND/RESUME） | 挂起前 Dispose 会话 scope（释放模型）+ 停采集；恢复后重开音频设备 + 重进 Armed/Sleeping + **重校系统时钟** | S6 |
| 2 | 音频设备热插拔 | `IMMNotificationClient::OnDefaultDeviceChanged` | 重枚举一次；失败 → `Error(audio_device)` 明示设备名；不得静默沿用旧句柄 | S2 |
| 3 | 显示器拓扑/DPI 变化 | `WM_DISPLAYCHANGE`/`WM_DPICHANGED` | 重算悬浮球位置到可见区；按新 DPI 重建 D2D 设备相关资源 | S7 |
| 4 | 网络/代理/PAC/企业 TLS 拦截 | 请求失败分类 | 尊重系统代理（WinHTTP 默认）+ `HTTP(S)_PROXY`；TLS 拦截时区分「证书不受信」与「连不上」 | S1 |
| 5 | 磁盘满/WAL 无界 | 写前检查可用空间（阈值 200MB） | 不足 → 停写 L3 与日志并告警但不得影响主响应；WAL 启动 TRUNCATE + autocheckpoint=1000；artifacts 500MB LRU | S6 |
| 6 | 省电模式/热节流 | `SYSTEM_POWER_STATUS.BatterySaver` · Power Throttling | 关 KWS 常驻（改纯快捷键唤起）+ Warm 缩 30s、Conversation 缩 15s；热节流下降 ASR 线程优先级 | S6 |
| 7 | Windows Update 重启/会话结束 | `WM_QUERYENDSESSION`/`WM_ENDSESSION` | ≤5s 走 D38 快速版关停（跳过非关键 flush，第 4/5/9 步不可跳）；单实例互斥必须 per-session | S1 |
| 8 | Focus Assist 抑制通知 | 无稳定 API | 无条件退化：总是同时更新悬浮球角标 + 面板任务列表 | S6 |
| 9 | 系统时钟跳变（NTP/时区） | — | **所有超时与保留期用单调时钟**（Go `time.Since` 天然单调）；只有落盘时间戳用 wall clock；禁 wall clock 差值做超时（D22 禁令） | S1 |
| 10 | GDI/User 对象/句柄泄漏 | SLO 采样含 GDI/User/句柄/goroutine/线程 | 阈值：句柄 >2000 或 GDI >5000 或 goroutine 超上限 → 告警；只量 RSS 看不到这类泄漏（Win32 自绘 + WebView2 的典型泄漏形态） | S6 |
| 11 | WebView2 Runtime 缺失/过旧 | go-webview2 创建失败 | 无面板模式（C27 原生 L2 降级卡）+ 引导 Evergreen 安装；不自动下载安装器 | S5 |
| 12 | 麦克风被独占/权限被拒 | 打开设备失败 | 明示 + 指向系统设置（Windows 隐私设置） | S2 |

## 5. 必须显式处理的平台异常态（§2 强制清单，对照 SPEC-08 §3 状态表）

- 麦克风被其他应用占用 → `Error(audio_device)`
- 麦克风权限被系统拒绝 → `Error(audio_device)` + 引导设置页
- WebView2 Runtime 缺失 → 无面板模式提示
- 显示器拔插悬浮球落入不可见区域 → 自动回主屏可见位置

## 6. 测试决策

- 失败预演要求**真实触发**（D22/§10 第 7 条）：填满磁盘（预留虚拟盘）、改系统时间、拔显示器
  （或改分辨率）、断网、独占麦克风（用测试进程占住设备）、模拟挂起（`psshutdown -d` 或手动）。
- 单实例：同会话第二实例断言激活已有实例并退出；不同会话语义以文档测试为准（CI 无多会话环境，
  人工验收项）。
- 时钟：注入时钟跳变（改系统时间 +5h）后断言超时行为不变（单调时钟）。
- 会话结束：发 `WM_QUERYENDSESSION` 断言 5s 内退出且音频/Job 清理完成。

## 7. 不做什么

- 不做电源管理之外的系统服务集成（无 Windows Service 形态——常驻语义是 per-user 桌面应用）。
- 不做 UIPI 绕过/提权注入；不做 DRM 窗口截屏对抗。
- 不做 Focus Assist 状态的私有 API 探测（拿不到就退化，不做脆弱 hack）。
