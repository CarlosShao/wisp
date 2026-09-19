# SPEC-04 · 语音链路与模型分发

> 追溯：D5/D16/D25/D26/§14.4/C8/C9/C29、D32 16.3.3/16.3.4、P3/P4/P5/P7/P14、D42#2#12。

## 1. 背景与问题

草案语音选型全部不合格（Porcupine 付费墙 / OpenWakeWord NC 许可 / Whisper tiny 中文 CER 15–25% /
Piper GPL / Edge TTS 灰区）→ 定案 sherpa-onnx 本地基座 + provider 抽象（D5）。本 spec 落到：
采集与端点检测参数、模型清单与分发校验、半双工与隐私、脱敏、延迟预算的逐段归属。

## 2. 链路总览与线程归属

```
AudioSource(C8) → 16k/mono PCM → VAD(silero) → ASR(sherpa 流式) → 标点恢复 → 上下文组装
                                                                        ↓
回答：TTS(sherpa) ← LLM 首 token ← Agent 循环
KWS：独立于 ASR 的常驻小模型（Armed 态），命中唤醒词/否决词
```

- 半双工（D16）：TTS 播报期间**关闭麦克风与 KWS**（从根上消除自激）；打断通道 = 全局快捷键 /
  点悬浮球 / 控制词（非播报期）。
- D16 推论（D32 16.3.3）：半双工 ⇒ ASR 与 TTS 时段互斥 ⇒ **串行加载，峰值 = max(ASR, TTS)
  而非 sum**——这是 700MB 预算成立的前提；切换延迟（加载 TTS ~300–800ms）必须实测（S0 输出⑤）。
- 音频缓冲区**永不落盘、不写日志**（D16③）；`keep_audio` 硬编码 false（SPEC-03 §3）。

## 3. AudioSource（C8）

```go
type AudioSource interface {
    // Start 后按 16kHz/mono/int16 推帧；ctx 取消即停
    Start(ctx context.Context, buf chan<- []byte) error
    Stop() error
}
// 实现：WASAPIMicrophone（共享模式，pinned 线程）、WavInjector（测试钩子，C8 定案）、
//       [AecSource] 仅接口位（DEFERRED，D16）
```

- 采集：WASAPI 共享模式；设备原生采样率 ≠16k 时**进程内重采样**【SPEC：线性插值重采样器，
  `audio` 包内实现，不加依赖】；帧长对齐 VAD 输入（512 样本 @16k ≈ 32ms）。
- 有界 channel ≤200ms（D38d），满则丢帧**并计数**（计数进日志/诊断包）。
- 设备热插拔：`IMMNotificationClient::OnDefaultDeviceChanged` → 重枚举一次；失败 →
  `Error(audio_device)` 明示设备名；**不得静默沿用旧句柄**（D42#2）。
- 麦克风被独占/权限被拒：明示 + 指向 Windows 隐私设置，不得静默失败（D42#12）。
- 静音键：一键全局静音 → `Muted` 态；KWS 激活时悬浮球常亮指示（D16②）。

## 4. VAD / ASR / 标点 / TTS / KWS

| 组件 | 模型（D5/D32 16.3.3，估算待实测） | 加载/卸载 | 关键参数 |
|---|---|---|---|
| VAD | silero-vad（MIT，~1M，~8MB） | `Listening` 起 → 会话结束 | 判停静音 ~700ms【SPEC 默认】；最短语音 ≥300ms（D43#11）；【SPEC】最长语句 60s 保护（超长自动判停） |
| KWS | `kws-zipformer-wenetspeech-3.3M`（~25MB 含 session） | Armed 态 → 关 KWS/回落 | `--keywords-file` 文本自定义**无需训练**；多关键词 + **逐词阈值**（P14 实测）；否决词（「取消/停下/别」）与唤醒词同表注册——`Confirming` 态语音否决的唯一可行通道（B1，不走 ASR） |
| ASR | 流式 Paraformer / 离线 SenseVoice（int8 **强制**，fp32 不可行） | 会话起 → Warm 超时 | `intra_op_num_threads = 1`（**强制**，D32）；CER 门禁见 SPEC-10 §4 |
| 标点 | CT-Punc 类（~100–200MB 视量化） | ASR 后 → Warm 超时 | ≤300ms（D32 延迟漏项①）；P4 效果评测落 S4 |
| TTS | matcha / vits 中文（~120–300MB） | 会话起（或首播报前）→ Warm 超时 | 首帧 ≤800ms（Warm 内）；音质主观 ≥7/10 门禁（P7 落 S4） |

- **一个 onnxruntime 实例复用全链路**（D5）；一律 int8 量化。
- 串行化策略：【SPEC】`speech` 模块内建「引擎 slot 互斥」：同一时刻至多一个大模型 family
  （ASR 或 TTS）驻留；TTS 需要而 ASR 驻留时，先 Dispose ASR scope（C11）再加载 TTS；
  若 S0 实测切换 >1.6s 预算 → 退化方案为「TTS 常驻 + ASR 按需」（TTS 更小）。
- 云端 ASR/TTS provider：C9 接口位保留，实现 DEFERRED（D5/D23）；启用必须显式告知「音频将上传」。

## 5. KWS 行为细节（B1/P14）

- 唤醒词命中 → `Listening`，**KWS 暂停**（半双工内不需要同时跑）；`Confirming` 态例外：
  KWS 保持运行以识别**否决词**（~100ms 量级推理，与 2–3s 阻止窗口相容）。
- **KWS 未加载（快捷键唤起路径）时**：`Confirming` UI 必须明示「语音取消不可用」，不得假装支持；
  否决通道只剩 单击球 / Esc / 面板拒绝（B1 定案，DEFERRED：依赖 AEC）。
- 否决词误触发率实测（P14）：误触发高 → 提高否决词阈值或去掉该通道（退化为「两路径都不支持
  语音否决」），S3 切片卡定案。

## 6. 隐私与转写脱敏（D16 三件套 + §14.4）

1. 默认不打开麦克风（KWS opt-in）——可写进 README 第一行的强论证。
2. KWS 激活常亮指示 + 一键静音快捷键。
3. 音频缓冲永不落盘、不写日志。
4. **转写文本**（§14.4）：
   - 日志：默认不进日志；诊断模式只记长度与置信度。
   - L3 任务日志/SQLite 历史：产品功能必须存，但 GUI 全量可见、逐条删、一键清空。
   - 诊断包：默认不含转写；显式勾选才含（导出前展示将包含哪些条目）。
   - 模式匹配脱敏（密码/卡号/token 样式掩码）：尽力而为，UI 必须标注「非完全脱敏」，
     不得给虚假安全感。

## 7. 模型分发（D26 + C29，阻塞 S2 的 P5）

### 7.1 触发与体验

- 首次使用某能力时下载对应模型；KWS 模型仅 3.3MB【SPEC：随安装包附带，免首跑下载】。
- 下载：断点续传（HTTP Range）+ 失败重试 + 进度（`Downloading` 态进度环）+ 可取消。
- 下载源：`[models] mirror`（国内可达镜像优先，HuggingFace 国内访问困难）+ 官方源兜底；
  **镜像只提供字节，认证完全来自仓库内签名清单**（D33/F3）。
- 逃生门：`[models] local_override` 指向本地已有模型目录（离线安装唯一路径）；GUI 显示已下载
  模型与体积、可单独删除。

### 7.2 C29 ModelManifest（CRITICAL）

- 仓库内置 `models/manifest.json`：每模型 `{id, purpose, urls[], sha256, size_bytes, license, quant}`；
  用 **minisign 私钥签名**，**公钥硬编码进 `buildinfo`**。
- 下载流程：先验 manifest 签名（离线可验）→ 下载 → 验文件 sha256 → 落 `models\<id>\`。
- `[models] verify_signature` **不可关闭**：写 `false` 直接报错（配置项只为显示状态）。
- **不得从镜像站取哈希**（哈希与模型同源 = 零认证价值，D33/F3）。
- 许可证：P3 逐模型核实权重许可，结果回填 manifest 的 `license` 字段（OpenWakeWord 的坑）。

## 8. 延迟预算逐段归属（D32 16.3.4，验收见 SPEC-10 §5）

| 段 | 预算 |
|---|---|
| 唤醒词说完 → 球转监听中 | ≤300ms |
| VAD 判停 → ASR 全文 | ≤800ms 流式 / ≤1500ms 离线 |
| ASR 全文 → 标点完成 | ≤300ms |
| 标点 → 上下文组装（BM25 + SQLite 读） | ≤50ms |
| 组装 → LLM 首 token | ≤1500ms（含网络） |
| 首 token → TTS 首帧 | ≤800ms（Warm，已加载）/ ≤1600ms（需加载） |
| 冷会话全程 | P50 ≤4s / P95 ≤6s；Warm 内 P50 ≤2.5s / P95 ≤4s |
| 冷启动（温启动，非首次运行） | ≤1s；首次运行含模型下载，不计入 |

## 9. 测试决策

- **全部走 C8 WavInjector**，不依赖真麦克风（整个测试策略的地基）。
- wav 注入集：`testdata/asr-baseline/{near_clean,noisy_far}`（CER 双门禁 6%/15%，随 CI 跑，
  带噪不阻塞发布但必须验证 UI 提示与 DEFERRED 登记）。
- 端点检测：构造「说话-停顿-说话」wav 断言判停时机；最长语句保护触发。
- 半双工：注入「TTS 播报音频」到输入侧，断言播报期间采集被关（无 ASR 结果产生）。
- KWS：否决词 wav 命中断言 `Confirming` 取消；KWS 未加载时断言 UI 提示「语音取消不可用」。
- 模型分发：compose model-mirror 提供正常/损坏/缺清单三类文件 → 断言 sha256 失败拒绝、
  manifest 签名失败拒绝、`verify_signature=false` 报错、断点续传续跑、取消清理半成品。
- 热插拔：枚举器注入假设备变更事件（接口层 mock 设备列表，**不是 mock 音频数据**）。

## 10. 不做什么

- 不做 AEC/barge-in（DEFERRED，`AecSource` 仅接口位）；不做全双工。
- 不做唤醒词训练功能（keywords 文本文件即自定义，D5）。
- 不做英文/多语言语音模型（i18n DEFERRED；注意语音链路是中文模型，非 locale 问题）。
- 不把 ASR/TTS 模型打进安装包（体积不可能，D24/D26）。
