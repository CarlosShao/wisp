# HANDOVER — Wisp 实施交接文档（给新会话的新 agent）

> 写于 2026-09-19T23:08:11Z（+0800）。**用户平台配额 2026-09-20 09:00 到期**——本文档保证新会话无缝接续。
> 阅读顺序：§1 使命与铁律 → §2 环境 → §3 完成度 → §4 在途状态 → §5 下一步 → §6 技术事实 → §7 事故手册 → §8 待人项。

## 1. 使命与治理铁律

- **目标**：按 `.scratch/wisp/issues/README.md` 的票据顺序完成 Wisp（Windows 语音 Agent 悬浮助手）全部 61 张票。
- **权威文档**：`docs/PLAN.md`（D1–D47 决策 + C1–C32 契约，冻结不得擅改）→ `docs/specs/SPEC-*.md`（13 份落地规格）→ `.scratch/wisp/issues/`（61 张原子票）。
- **票据铁律**（`.scratch/wisp/issues/README.md`）：开工先置 `Status: in-progress` + 填 Claimed by + Progress log 追加（newest last）+ commit push；**每完成一个最小单元立即 commit+push**；完成 → 全验收框打勾 → Status: done → **文件重命名加 `-done` 后缀 + 标题加 (DONE ✅)** → 更新索引。
- **对抗验收**：每票实现完成后由**非实现者**验收（子代理或编排者亲自——独立性=实现者≠验收者）；FAIL → 最小修复集发回原实现代理 → 复核 PASS 才置 done。报告存 `docs/evidence/`。
- **commit 纪律**：**永远显式列路径，禁止 `git add -A`**（两次吞并行 WIP 事故教训）；Conventional Commits；双远程 push（origin 偶发 TLS 失败重试≤5 次，cnb 稳定）。
- **并发纪律**：子代理并发 2 为底线、3 为试验位（3 曾触发平台 captcha/quota 故障，回落 2）；子代理闪断（Captcha timed out / exceed quota）→ 直接重派续传，**不丢弃不降级**。
- **D22 未定义即停**：方案未覆盖 → 停下来问用户；DEFERRED 项不得擅改。

## 2. 环境（全部已就绪并实测）

| 项 | 位置/值 |
|---|---|
| 仓库 | `D:\work\workspace\projects plans\Wisp`（git，dev 分支为工作分支） |
| 远程 | origin=github.com/CarlosShao/wisp，cnb=cnb.cool/CarlosShao/wisp（凭据已存，gh CLI 已登录） |
| Go | `D:\work\base\go\bin\go.exe`（1.27.1，不在 PATH，自行 export）；GOPATH=`D:\work\base\gopath`；GOPROXY=goproxy.cn |
| cgo 工具链 | `E:\work\base\msys64\mingw64\bin`（gcc 16.2.0）；make 在 `E:\work\base\msys64\usr\bin` |
| 构建 | `scripts/build.ps1 [-Env dev|prod]`（一键：fetch-deps→build→DLL 同目录→SHA256SUMS）；`docs/BUILD.md` 已冻结 |
| 自托管 runner | wisp-selfhosted-01（标签 self-hosted/Windows/X64/wisp-slo），登录自启（HKCU Run → E:\work\base\start-wisp-runner.vbs） |
| Docker | Desktop 可用；compose 的 model-mirror（18081）+ fixtures（18082/18083）票 14 已建，mock-llm（18080）票 08 接线 |
| minisign dev 密钥 | `E:\work\base\wisp-minisign\`（盘外，公钥硬编码 internal/buildinfo；生产密钥仪式归 S8） |
| 真麦克风 | Realtek，LIVE 冒烟：`WISP_LIVE_MIC=1 go test ./internal/audio/ -run TestLiveWasapiSmoke` |
| 模型缓存 | `third_party/spike-models/`（git 忽略；URL/SHA256 见 docs/evidence/s0/json/） |

## 3. 完成度（截至本文档：9/61 done）

| 票 | 内容 | 证据 |
|---|---|---|
| 01 ✅ | 构建链 + BUILD.md 冻结 | docs/evidence/s0/01-* |
| 02 ✅ | S0 spike（Path Y 判定、五基线、goja/WebView2/模型常驻实测） | docs/evidence/s0/02-*、docs/SLO.md |
| 03 ✅ | 19 包骨架 + D37 错误模型 + goroutine 注册表 + C11/C30 | docs/evidence/s0/03-* |
| 04 ✅ | SQLite 核心（8 表 DDL、WAL、db-writer、保留期）+ schema v2（provider_health，票 09） | docs/evidence/s0/04-* |
| 05 ✅ | config 全量模型（catalog v2、三档热加载、🔒方向引擎、迁移） | docs/evidence/s1/05-* |
| 06 ✅ | SecretStore(DPAPI) + WISP_ENV 分叉 + 便携/P13 | docs/evidence/s0/06-* |
| 07 ✅ | 悬浮球窗口 + D2D 渲染 + 20 态状态机（D43 全 42 行）+ 热键/托盘/多显示器 | docs/evidence/s1/07-*、ball-states/ 截图 |
| 09 ✅ | LLM seam（C5/C6/C7 全集）+ OpenAI Chat adapter + golden 两跑器 + mockllm + text_chain/令牌桶 + provider_health 原语 | docs/evidence/s1/09-* |
| 13 ✅ | 音频采集（WASAPI 真机验证、热插拔、半双工门、重采样） | docs/evidence/s2/13-* |

**S0 完成；S1 剩 08（在途）/10/12；S2 已开 13✅+14（在途）。**

## 4. 在途状态（截至本文档）

- **票 14（模型分发，T14-impl 运行中）**：已完成 C29 签名清单（6 条目，复用票 02 验证过的 URL/SHA256）、dev minisign 密钥仪式、compose model-mirror（提交 5ab0339）；预计剩余：下载管线测试收尾。已收截止警戒。
- **票 08（可观测性+CI，T08-impl 运行中）**：正在写 observe/logging.go、redact.go、sampler.go（未跟踪 WIP）；剩余：slo-check.ps1、CI workflows、compose mock-llm 接线。已收截止警戒。
- 两代理均已被指示：08:30 前完不成就在干净单元边界停（deadline-pause 日志），**不留半截提交**。

## 5. 下一步（新会话的行动清单，按序）

1. **收在途**：若票 14/08 处于 deadline-pause → 派续传代理完成（提示词模板见历史模式：读票据+进度日志断点+显式路径提交）；完成后编排者（或新子代理）做对抗验收 → 置 done。
2. **票 10（agent 循环核心，2–4h 大票）**：被 05✅+09✅ 解锁。参考：Pi 骨架（SPEC-05 §2）、golden 回放（票 09 的 `internal/llm/golden` + mockllm 控制端点）、failToolCallsFromTruncatedMessage、D15 阈值随 context_window 缩放。
3. **票 12（S1 gate）**：被 04✅+06✅+07✅+08✅+09✅+10 解锁。含 SLO 实测（回落调优义务——T02 发现 Settle 残留 ~31MB>25MB，`debug.SetMemoryLimit` 调优是票 12/15 的明文义务）+ 打字→回复→通知端到端 + 人工视觉项（H1）。
4. **S2 收口**：票 15（TTS 引擎——**TTS 常驻+ASR 按需为默认策略**，spike 实测串行 4.7s>1.6s 预算）、票 16（S2 验收：CER 双门禁 6%/15% + 真机冒烟三条物理场景）。
5. **S3 批量**：17–25 九张票（安全+能力面），两并发按 17∥18 → 19/20 → 21–24 → 25 推进；每片完成跑缺口审计+对抗验收。
6. **Key 类等待项**：真 LLM Key 到位后先跑 `cmd/llmrecord` 补录真 provider 黄金流（票 09 的 H2 DEFER）。

## 6. 关键技术事实（省得新会话重踩）

- **spike 结论（实测，已回填 docs/SLO.md）**：Path **Y**（单进程 cgo，空闲 16.5MB≤25）；**TTS 常驻+ASR 按需**（串行切换 4.7s>1.6s）；WebView2 冷 880–1126ms（P11 过，L2 卡不回退原生）；goja **无 ES modules**（Tier2=ES2020-minus-modules，Interrupt 精度 OK）；句柄口径 <600（D2D 窗口栈实测 ~420）；**Settle 残留义务**（Sleeping ~31MB>25MB → SetMemoryLimit 调优归票 12/15）；matcha 无官方 int8（TTS fp32 174MB）。
- **intra_op_num_threads=1 强制**；音频帧 16k/mono/int16/512 样本；`WISP_ENV` 三环境分叉（prod/dev/test 目录与互斥名全隔离）；config schema_version=**2**（catalog v2）；SQLite schema v2=9 表（provider_health）。
- **golden 格式**：`# wisp golden sse v1` + `# @response/@latency` 指令——单测回放器与 mockllm 子进程两跑器字节一致（契约）。
- **WASAPI 勘误**（真机发现）：IAudioClient vtable Start=10/Stop=11/SetEventHandle=13/GetService=14；共享模式 BUFFER_SIZE_NOT_ALIGNED 对齐重试——`internal/audio/wasapi_windows.go` 有注释留痕。
- **待修技术债**：cmd/wisp 的 sherpa cgo 使 `CGO_ENABLED=0 go build ./...` 在 cmd 失败（internal 全绿——历棒基线，票 08 CI 已按此设计）。

## 7. 事故手册（已验证的应对）

| 症状 | 应对 |
|---|---|
| 子代理派发即死（Captcha timed out / exceed quota） | 直接重派续传（随机性，重试存活率高）；持续失败 → 降并发或编排者亲自 |
| origin push TLS/EOF | 重试≤5 次×12s；cnb 稳定可先推 |
| 子代理死在中途 | 票据 Progress log 是断点：新代理读日志 `next=` 字段接续；盘上未跟踪 WIP 属实且完好 |
| 平台容量（Start Plan busy） | 回落 2 并发；验收可由编排者亲自（独立性=实现者≠验收者） |

## 8. 待人项（docs/reports/pending-and-issues.md 同步维护）

- **H1** 悬浮球 20 态视觉签收（看 docs/evidence/s1/ball-states/*.png）
- **H2** LLM API Key（阻塞真 provider 黄金录制；框架已就绪，`cmd/llmrecord` 一条命令补录）
- **H3** P10 命名残余核查（阻塞票 56）
- **H4** web.search 实现路径拍板（票 22 以接口先行，不返工）
- **H5** 生产 minisign 密钥仪式（S8 前；dev 密钥已可用）
- **H6** 票 16 真机三场景（物理拔插麦克风/隐私开关/独占占用——人工操作）

## 9. 进度节奏参考（新会话据此排程）

实测：2 并发下 ~9 票/24h（含平台故障开销），单票中位 1–2h、大票 2–4h、验收 0.5–1h。
剩余量级：S1 尾（10+12）≈1 天 → S2 收口（15/16）≈半天 → S3 九票 ≈3–4 天 → S4 ≈3 天 → S5 ≈4 天 → S6 ≈2 天 → S7 ≈4 天 → S8 另计。**全项目 ≈3–4 周连续双代理运行。**
