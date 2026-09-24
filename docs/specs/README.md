# Wisp Specs — 实施规格文档集

本目录是 [../PLAN.md](../PLAN.md)（实施级方案，D1–D47 决策 + C1–C32 冻结契约）的**落地规格**。
PLAN.md 定「做什么、为什么、边界在哪」；specs 定「具体怎么做：骨架、接口、数据、参数、流程」。

## 效力关系（冲突裁决顺序）

1. **PLAN.md 定案内容（D1–D47 / C1–C32 / §16 修订 / §15 定案）** —— 最高，spec 不得与之矛盾
2. **本目录 spec 的【SPEC】标记条目** —— 为落地新增的细化（PLAN 未定但实现必需），
   默认按 spec 执行；人工审阅后如不同意，改 spec 即可（成本低）
3. 代码注释 / 实现细节

任何与冻结契约的偏离 = 契约变更，必须人工批准（D22 闸门①），并同步更新 PLAN.md 与受影响切片卡。

## 标记约定

- `（Dxx）` / `（Cxx）` —— 追溯到 PLAN.md 的决策 / 契约编号
- `（Sx）` —— 追溯到切片 S0–S8
- 【PLAN】—— 直接来自 PLAN.md 的定案，本 spec 仅展开
- 【SPEC】—— 本 spec 为落地新增的细化，需人工确认但默认按此执行
- 【OPEN】—— PLAN.md 已登记未决（§15 / §16.11），spec 给出建议但不得视为已定案

## 文档地图

| 文件 | 内容 | 主要服务的切片 |
|---|---|---|
| [SPEC-00-product-overview.md](SPEC-00-product-overview.md) | 产品定义、用户故事全集、范围与场景矩阵、术语表 | 全部 |
| [SPEC-01-architecture.md](SPEC-01-architecture.md) | 进程拓扑、模块边界、**仓库骨架（文件树）**、线程/goroutine 模型、关停顺序、错误模型 | S0–S1 |
| [SPEC-02-data-storage.md](SPEC-02-data-storage.md) | **数据存储**：SQLite 全量 DDL、WAL 策略、保留期清理、数据目录、artifacts 配额、备份导出 | S1, S4, S6 |
| [SPEC-03-config-secrets-envs.md](SPEC-03-config-secrets-envs.md) | 配置模型全量 schema、三档生效级别、SecretStore(DPAPI)、**开发/测试/生产环境隔离** | S1 |
| [SPEC-04-voice-pipeline.md](SPEC-04-voice-pipeline.md) | 语音链路：音频采集/VAD/ASR/标点/TTS/KWS、半双工、模型分发与签名校验、转写脱敏 | S2, S4, S6 |
| [SPEC-05-agent-core.md](SPEC-05-agent-core.md) | Agent 循环、LLM 三协议 adapter、流式事件、上下文预算、记忆系统、会话保活、防死循环、成本 | S1, S3, S4 |
| [SPEC-06-security-gatekeeping.md](SPEC-06-security-gatekeeping.md) | 三级风险门控、R1–R9 规则集、路径规范化、污染追踪、注入五层防御、审批队列、批量授权 | S3, S7 |
| [SPEC-07-tools-and-plugins.md](SPEC-07-tools-and-plugins.md) | 内置工具权威表、ToolProvider、插件 manifest schema、Tier2 goja 约束、CLI 包装插件 | S1, S3, S7 |
| [SPEC-08-ui-ball-panel.md](SPEC-08-ui-ball-panel.md) | 悬浮球（Win32/D2D）、**20 态状态机权威转移表**、WebView 面板、React 技术栈、设计系统 | S1, S5 |
| [SPEC-09-platform-windows.md](SPEC-09-platform-windows.md) | Windows 平台细节：DPI/多显示器、音频设备、电源、时钟、单实例、12 项可靠性事件 | S1–S7 |
| [SPEC-10-testing-acceptance.md](SPEC-10-testing-acceptance.md) | **测试接缝（seams）**、SLO/CER/延迟门禁、安全红队用例、失败预演、对抗验收、场景验收 | 全部 |
| [SPEC-11-build-deploy-containerization.md](SPEC-11-build-deploy-containerization.md) | **构建工具链、容器化（构建镜像/测试依赖/CI）、环境隔离实现、打包分发与更新** | S0, S7, S8 |
| [SPEC-12-roadmap-governance.md](SPEC-12-roadmap-governance.md) | 切片路线图 S0–S8、场景↔切片追踪、DEFERRED/RESERVED/REJECTED 登记、治理流程 | 全部 |

## 关于「容器化」的边界声明（先读）

Wisp 是 **Windows 桌面常驻应用**（Win32 分层窗口 + WASAPI 音频 + WebView2 + sherpa-onnx cgo），
**没有独立的后端 API 服务**——它的「后端」就是常驻主进程，LLM 是第三方云 API。因此：

- ❌ **Wisp 主程序本体不容器化**：Win32 窗口、音频设备、悬浮球必须在真实 Windows 桌面会话中运行，
  Windows 容器无法承载（详见 SPEC-11 §1）
- ✅ **可容器化的全部容器化**：构建工具链镜像、前端构建、测试依赖服务（mock LLM / mock 站点 /
  模型镜像源 / SSRF 靶机）、CI 流水线、纯逻辑开发环境（devcontainer）
- ✅ **开发/测试/生产环境隔离**：`WISP_ENV` 环境模型 + 每环境独立数据目录与单实例锁 + 测试专用
  mock 端点（详见 SPEC-03 §5 与 SPEC-11 §6）

## Agent 使用方式

- 按切片开工前必读：SPEC-00（全局）+ 该切片对应 spec + SPEC-10（测试接缝）+ SPEC-12（该切片 done 判据）
- 契约冻结细节以 `docs/contracts/C1..C31.md`（按 PLAN §12 于 S1 生成）为最终文本；
  本 spec 的接口签名草案与之冲突时**以契约为准**并停下询问（D22 闸门③）
- 仓库骨架以 SPEC-01 §3 为准；模块与包的对应关系 1:1 映射 §1.2 冻结模块表，不得自行增删模块
