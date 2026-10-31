# 票 59 AC#1 — AEC 候选 license 判定（腿 59-r1）

日期 2026-10-31 · 腿 `59-r1`（写腿 AEC spike）· 起手锚 `dcf6e910`（dev，非 detached）
写面＝`internal/audio/aecspike/**` ＋ 本件；⛔ cmd/wisp（41a2-r1 在飞）、⛔ internal/speech 语义、⛔ 下载。

## 判定总表

| 候选 | 是什么 | license 结论 | 「license 文件入 evidence」 | 状态 |
|---|---|---|---|---|
| 乙：webrtc-audio-processing（AEC3，票面点名） | 需 fetch 源码/绑定后自建 | 上游声明 BSD-3，**本机未现核**（见下） | **无**（⛔ 下载边界挡住 fetch ⇒ 无从提交真实 license 文件） | **本腿 NO-GO（评估不可达）** |
| 甲：Windows 内置语音 AEC | 系统组件（Voice Capture DMO / APO），零下载 | Windows OS license 附带组件，**无独立 license 文件**（结构性不存在，非缺失） | 无文件可交；以「组件在场性探针＋OS build＋APO 清单」代之 | **本机不存在**（探针 0x80040154） |
| 丙：不上 AEC（Path C barge-in 退化为热键/单击打断，D47 白名单条目撤回） | 无依赖 | 无 license 面 | 不适用 | 可用（代价＝砍功能） |

## 乙（票面点名候选）——为什么本腿只能判 NO-GO

- 盘面现核：`deps.toml` 无 webrtc 条目、`go.mod` 无任何 webrtc/speex 绑定、Go module 缓存
  （`D:\work\base\gopath\pkg\mod`）无 webrtc/speex/opus 目录、MSYS2 mingw64 未装该包 ⇒
  **本机不存在任何可评估的 AEC3 源/构建/绑定**。
- 票面 AC#1 的判据是「verify the exact build/distribution we fetch（deps.toml entry +
  fetch-deps, SPEC-11 §2.1）」——**前提是 fetch**。派单硬闸＝⛔ 下载任何东西；三条未定义即停
  之首＝「候选方案全部要下载/花钱 ⇒ 停手」。故本腿**不 fetch、不代填 license 结论**：
  凭记忆抄 BSD-3 文本进 evidence 是伪造，⛔ 做。
- ⇒ **乙的 license 格＝「未完成，且按『无 license ⇒ NO-GO 候选』规则记 NO-GO（本腿射程内）」**。
  解锁条件＝机主批准 fetch（属「花钱/改依赖」两类人工批准项），批准后由票 13 的 module 工作落
  deps.toml＋license 文件入 evidence，再重跑本 spike 的质量/CPU/误触发测量。
- 跨编译状态（票面 AC#2 附带项）：同样因不可 fetch **无法现核**；票面预期「mingw Linux-cross
  likely NOT viable」照抄为待验假设。本仓 `GOOS=linux go vet ./internal/audio/...` rc=0（纯
  Go 面），不构成 AEC3 的交叉结论。

## 甲（Windows 内置）——本机不存在，license 前提照实记录

- Voice Capture DMO（CLSID_CWMAudioAEC）探针＝`CoCreateInstance` → **hr=0x80040154
  （REGDB_E_CLASSNOTREG）**；`HKCR\CLSID\{745057C7-...}\InprocServer32` 不存在。
- 捕获端点 APO 清单（MMDevices\Audio\Capture\*\FXProperties 值级转储＋CLSID→DLL 解析）＝
  **RtkRecVolMFX / RtkRecMFX / RtkRecInvVolMFX / RtkRecEFX（RltkAPOU64.dll）＋Creative SFX
  （MBAPO264.dll）**，无 Microsoft 语音 AEC APO；`System32` 无 AEC 类 APO DLL 文件名。
- ⇒ 甲的 license 格**在本机没有可判对象**。license 前提照实记录：OS 附带组件、随 Windows
  license 走（本产品 D7＝Windows-only）、**不可独立再分发** ⇒ 若在有的机器上可用，也只能
  运行时探测＋优雅降级（`wisp doctor` 形），不能进 deps.toml。
- 环境＝Windows 10 Pro build 26100（24H2），Realtek Audio（默认渲染/采集）＋NVIDIA HDMI。

## 丙——合法裁决出口

D47 原文即给出退化路径：license 惊喜 ⇒ barge-in 退化为 hotkey/click、白名单条目撤回、
⛔ 复活 Path T 全双工。本腿测量（见 `59r1-aec-spike.md`）补上了丙形的量化代价：无 AEC 时
自回声在采集通路里被 ASR 逐字转写（self-echo corr=0.978，ERLE 仅 16.8dB 的单增益扣除不足
以消音）⇒ 丙形不是「体验差一点」，是 Path C 全双工**不成立**。

## 结论（接 PRECHECK P15）

三候选在本机全部不可用：乙要 fetch（⛔）、甲不存在（探针实锤）、丙＝砍功能。
⇒ **PRECHECK P15 = NO-GO（本机证据范围内）**，完整判语与逐格读数在 `59r1-aec-spike.md`。
