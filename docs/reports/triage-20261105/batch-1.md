# 批 1 只读分档（21 23 25 32 33 34 35 36 37 38 39 40 46 47 48 49 51 53 54 55）

> 只读分档，2026-11-05。票面与代码一字未改。依据列取自票面 Status/Blocked by/AC/进度记录＋盘上 `-done` 现量（17/20/22/24/27/29/30/31/07/12/10/04/41–45/52 均已 done）。

| 票 | 档 | 依据(一行) | 阻断者/首动作 |
|---|---|---|---|
| 21 | T1 | 段1 决策层已落且自测绿（L1 四否决、C18 单任务队列、D45-1 聚合、native-source nonce），`Blocked by: 17` 已 done，剩段2；AC#1 是纯 internal/tools＋internal/agent 的 Go 判据，不碰 frontend/** | 首动作＝AC#1（非实现者复跑 L1 超时执行＋四否决通道，当前零生产可达故勾不了）；Go 半＝删 `loop.go` 的 RiskL1/RiskL2 前置门＋以 `approval.New` 作 `Options.Gate`/`AdmitTextTask`（internal/agent＋internal/tools，约 2 包＋1 处 cmd/wisp 装配点）；native 卡与球视觉属真机 GUI，不在首动作内 |
| 23 | T1 | `Blocked by: 20` 已 done；AC1 风险矩阵/AC5 关机确认词/taint→R4 已满足，残余 AC2 为纯 Go 接缝（C7 ImagePart→LLM、bytes_ref 仅内存），不碰 frontend/** | 首动作＝AC#2：`screen.capture` 的 C7 `ImagePart` 经 internal/agent loop＋`llm.ImagePart` 直达 mock LLM，并断言 bytes_ref 不落 disk/db/log；范围 llm/agent/tools 约 3 文件＋测；AC3 确认卡三选项＝frontend、UIPI/剪贴板 round-trip＝真机，另计 |
| 25 | T2 | S3 场景①②＋全套红队；`Blocked by` 中 22/24 已 done，但 21/23 均未 done，且证据本身要真机场景脚本＋windows CI 的 junction 真测 | 阻断：21, 23 |
| 32 | T3 | S4 ③听写进目标 app 光标＋④多轮陪聊，需真机麦克风/扬声器、UIPI 真窗口、barge-in ≤400ms 实测，且主观分与 D29 人工签字是验收条 | 阻断：真机 desktop-audio/GUI＋用户人工签字（human acceptance） |
| 33 | T3 | WebView2 宿主窗口/焦点回还/冷热延迟实测全在真机，写面在 `cmd/wisp`（宿主 `panel_host_windows.go`）；AC#3 的“无监听端口”判据今天恒真（常量错位），AC#12/13/14 亦需真窗 | 阻断：真机 WebView2 桌面 GUI＋cmd/wisp 宿主面（非本编队纯 Go 写面） |
| 34 | T3 | `frontend/**` 脚手架、tokens 端口、docker 构建与 CI 前端 job，本编队对 frontend 零写面 | 阻断：frontend owner（机主外部 agent）；Go 侧无首动作 |
| 35 | T2 | C17 桥；出向 Go→页今天零调用者、页面无接收器，`Blocked by` 的 33/34 均未 done；核心传输仍落 cmd/wisp＋frontend 客户端 SDK | 阻断：33, 34 |
| 36 | T3 | 结果流/历史/transcript 面板页＝frontend/**，14 态视觉需人工签字；另依赖未 done 的 35（票面自记“真数据 e2e 等 145/35/33”） | 阻断：frontend owner（且 35 未 done） |
| 37 | T3 | 审批队列/L2 卡/native fallback 视觉＝frontend/**＋design 对照人工签字；`Blocked by` 21/35 均未 done | 阻断：frontend owner（另 21/35 未 done） |
| 38 | T3 | cmdk 命令面板/tasks 视图＝frontend/**，并需 Windows IME 真机 e2e；依赖未 done 的 35 | 阻断：frontend owner（＋IME 真机） |
| 39 | T3 | 配置编辑器 GUI＝frontend/**（schema 表单、三档徽章、api_key 只写），依赖未 done 的 35 | 阻断：frontend owner（＋35 未 done） |
| 40 | T3 | 安全/隐私/成本/模型页＝frontend/**＋三张 design 稿签字，依赖未 done 的 35；票内唯一 Go 落点（`internal/memory/privacy.go` 的 decision 人话文案）可作独立小活切出 | 阻断：frontend owner（＋35）；可切出 sidecar＝privacy.go 文案 AC（纯 Go，1 文件＋测） |
| 46 | T3 | S6 门：休眠/Armed/Warm SLO、5 类 settle、7 枚 power pre-mortem 全需真机采样（CPU/句柄/GDI/RSS 读数），`Blocked by` 41–45 已全 done，唯一缺口是真机与采样 | 阻断：真机桌面采样（非代码缺口） |
| 47 | T2 | TaskScheduler＋PathLock 本体是纯 Go，但 `Blocked by: 21` 未 done：审批-等待独立性/取消语义要 21 段2 把 Gate 接进生产 Loop 才有被测物 | 阻断：21 |
| 48 | T2 | C18 全量队列（FIFO/head-only/300s→reject/fail-closed）是 Go 活，但 `Blocked by` 的 37（前端卡）、47（调度器）均未 done | 阻断：37, 47 |
| 49 | T2 | D45-2 授权存储与风险路径评估纯 Go（04 已 done），但 `Blocked by: 48` 未 done：决策持久化与 grant_id 链接在 48，卡上第三选项在 37 | 阻断：48 |
| 51 | T1 | Tier2 goja（VM 生命周期/interrupt/结果帽/默认关闸门）纯 Go，`Blocked by` 17/20 均已 done；仅 AC#3 的 C24 初始集条目要冻结（契约面），其余 4 格可先做 | 首动作＝AC#1：死循环 JS 约 5s wall-clock 中断→tool 错误、宿主存活（命名 watchdog goroutine，禁裸 `go func`）；范围 internal/plugin 新包＋测；C24 面先留注册表接口，条目等冻结 |
| 53 | T2 | lark-cli 真插件 e2e；`Blocked by` 的 48/51 未 done（52 已 done），另需本机已认证 lark-cli 与真语音/文本输入 | 阻断：48, 51 |
| 54 | T2 | S7 总门；`Blocked by` 47/48/49/53 全未 done；另含真机 pre-mortem（拔显示器/DPI/kill -9 回滚）与人工裁决 | 阻断：47, 48, 49, 53 |
| 55 | T3 | macOS 移植票面即 DEFERRED 到开源准备期，`Blocked by: 54`＋自用一周＋机主明确放行（D23），且需 macOS 真机与签名策略 | 阻断：机主放行＋macOS 真机（DEFERRED，不可排期） |

T1-count=3
T1: 21, 23, 51
