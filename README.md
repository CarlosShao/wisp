# Wisp · 一缕

轻量级 Windows 语音 Agent 悬浮助手：本地优先、插件化、默认不碰麦克风、不用时真睡眠
（空闲进程树私有内存 ≤25MB，可测量验收）。

> **当前阶段**：方案与规格已定稿，尚未开始实现（从 S0 spike 开始，见 SPEC-12）。

| 入口 | 说明 |
|---|---|
| [docs/PLAN.md](docs/PLAN.md) | 实施级方案定稿（D1–D47 决策 + C1–C32 冻结契约，四轮打磨 + 用户全部定案） |
| [docs/specs/](docs/specs/README.md) | 落地规格：骨架、数据存储、语音链路、安全门控、UI、构建与容器化、环境隔离等 13 份 |
| [design/](design/index.html) | 前端静态原型（11 屏 + 设计系统总览，`tokens.css` 即 C21 DesignTokens 参考实现） |

## 技术栈（摘要）

Go（单常驻进程，无常驻 WebView）· sherpa-onnx 本地语音全链路（KWS/VAD/ASR/标点/TTS）·
goja（Tier2 插件，默认关闭）· jchv/go-webview2（按需面板）· React + Tailwind + shadcn/ui（面板）·
SQLite（WAL 单写者）。详见 SPEC-00/01/11。
