# 393 — AEC3 绑定层 flush/reset API（392 残余：保温真解的前置）

**Status:** **DONE-带残余**（2026-11-01 编排者收口：五格全过（393-v1 逐面现核）、半停手裁法对；★主新发现被反证修正=stale 窗最迟首帧前关闭、Reset 即统一器⇒保温修复无需动生产路径；残余四件见末节）
- **写面**：`internal/audio/aec/**`（★**绑定层＋shim；⛔ AEC3 源树**）
- **机理**：★**保温姿势判死的根因链**＝APM 跨窗携带失准＋无 flush/reset ⇒ **绑定层新增 flush/reset 面**（shim wap_flush/wap_reset＋Go 侧 Flush/Reset 方法）⇒ **后续保温修复票的前置**

## 判据
- [x] **AC#1 绑定层 API 落地**：★**shim＋Go 方法（flush/reset）**（★**对照 wap_new_tuned 形状**）
- [x] **AC#2 单测**：★**flush/reset 后状态清洁读数**（★**对照 392-r1 的 K=4 测量形状**）
- [x] **AC#3 反向判据**（★三读数＋还原核）
- [x] **AC#4 零回归**（★internal/audio 全套 94/83/0/11 基线＋392 的调参面不退化）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ AEC3 源树／`cmd/wisp`（26p 集成已收）／thresholds
## ★★★终裁落地（2026-11-01 09:5x 编排者，凭非实现者 `393-v1`，件 `docs/evidence/s1/393-v1-verdict.md` 131 行、`9b2805d50`）

★**五格全过 ⇒ 可收口**。
★**半停手裁法核＝对**（钉定树逐行现读：EchoControl 无 reset/flush；audio_processing.h:497-507 Initialize() 逐字＝真 reset 公开面；InitializeLocked :516-588 全重建＋:1887-1889 经 factory 重造；kBufferFlush/delay_controller_->Reset 全内部、kPlayoutVolumeChange 只置 gain_change＝加速重适配非清态；★**精度注记：reset 面在 AudioProcessing 非 EchoControl、shim 持 APM 句柄、落点正确**）。
★**AC#2 逐字节复现**（本腿自建两发 rc=0、工件两枚 JSON 与 393-r2 在盘件 diff rc=0；reset==rebuild 4/4 帧差 0）。
★**AC#3 carry 退化 7.87dB 复现且 Reset 移除**。
★**AC#4 口径差如实＋判据落点核过**（-v 口径 RUN=98 与 392-r1 自记逐位同；派单基线 94/83/0/11 溯源＝A1326 时代 389-v1 行旧测试总体口径——**我派单的基线数引的是旧口径**，教训入册）。
★**★★主新发现被反证修正（393-v1 新反证）**：★393-r2 的「stale submodule_states_ ⇒ 持续吃全带」**机制归属未证**——★**393-v1 读出两条实现未计入的懒同步路径：MaybeInitializeCapture :878 每次 float ProcessStream 都先跑 UpdateActiveSubmoduleStates、stale 即触发全量 InitializeLocked；格式变更 reinit 同样同步 ⇒ stale 窗最迟在首个处理帧前关闭**。
★**★★Reset 即统一器**（default+Reset ≡ factory 形逐位、dTuned=0/dKnobs=0 坐实）⇒ ★**保温修复票无需动生产路径：窗间调 Reset 后第二窗起全落 Initialize 形、fresh 形只在首窗、量级 0.0411/≤1.7dB＝启动暂态尺寸**。
★**处置四条（393-v1 建议、编排者采纳）**：(a) 维持不动生产路径（对）(b) 后续票把 §4 归属降级为「两形实存、全带/分带标签未证」(c) 动 wap_new→factory 形之前＝先跑 392 K=4 对照＋仓外仪器化探针坐实归属 (d) Reset 即统一器⇒保温修复票无需动生产路径。
★**残余四件**：flush 半边公开面不存在（未来升级带来时照 wap_reset 形加一法）／两引擎形机制归属未证（探针待需要时做）／旋钮同形效应 0＝合成局限（权威读数仍 392-r1 cmd/wisp K=4）／392 K=4 重跑不在已跑写面（非收口阻塞）。
★**393-r1 死亡核＝工件/证据件 0 命中**（死亡判定维持）。★**踩坑入件：cgo 旗标 MSYS /d/ 形会被当相对路径、必须 D:/ 盘符形**。
