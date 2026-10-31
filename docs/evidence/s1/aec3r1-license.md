# AEC3 fetch + license 判定（腿 aec3-r1）

日期 2026-10-31 · 腿 `aec3-r1`（AEC3 接入重测腿，写）· 起手锚 `0bb87859c`（dev，非 detached；
开工三核过：HEAD=锚点、索引无本腿写面脏、`internal/audio/**`＋`deps.toml`＋evidence 零他人脏——
cmd/wisp 的脏件属 384-r1 在飞、.scratch 残件不属写面）。

机主已批下载（A1311③ 逐字「批准下载 AEC3」）＝⛔下载边界对 AEC3/webrtc-audio-processing 解除
（仅此一项）。本腿据批实取、哈希逐件记账。

## 1. 判定结论（接 59-r1 license 件的 NO-GO 缺口）

**webrtc-audio-processing 2.1（AEC3）license＝BSD-3-Clause，逐字现核成立**——59-r1
`59r1-license-verdict.md` 里「上游声明 BSD-3、本机未现核」的待验项，本腿以真实取回的 license
文件全文关闭。**未定义即停①（license 不是 BSD 类）不触发。**

| 文件 | 来源（release 内原位） | license 文本 | evidence 落点（逐字拷贝） |
|---|---|---|---|
| 主 COPYING | `webrtc-audio-processing-2.1/COPYING` | BSD-3-Clause（"Copyright (c) 2011, Google Inc. All rights reserved."，三条款：保留版权声明＋保留 disclaimer＋不得用名义背书） | `aec3r1-webrtc-audio-processing-2.1-COPYING.txt` |
| webrtc 树 LICENSE | `webrtc-audio-processing-2.1/webrtc/LICENSE` | 同上 BSD-3（WebRTC 项目文本） | `aec3r1-webrtc-audio-processing-2.1-webrtc-LICENSE.txt` |
| abseil LICENSE | `abseil-cpp-20240722.0/LICENSE` | Apache-2.0（全文 11,361 字节） | `aec3r1-abseil-cpp-20240722.0-LICENSE.txt` |
| vendored ooura fft | `webrtc/common_audio/third_party/ooura/LICENSE` | ooura fft 原作者页头自述授权文本 | `aec3r1-vendored-ooura-LICENSE.txt` |
| vendored spl_sqrt_floor | `webrtc/common_audio/third_party/spl_sqrt_floor/LICENSE` | Wilco Dijkstra 1996 邮件授权文本（允许用于 any purpose） | `aec3r1-vendored-spl-sqrt-floor-LICENSE.txt` |
| vendored webrtc fft | `webrtc/modules/third_party/fft/LICENSE` | Mark Olesen 1995/97（FFTPACK 派生，自由使用声明） | `aec3r1-vendored-webrtc-fft-LICENSE.txt` |
| vendored pffft | `webrtc/third_party/pffft/LICENSE` | BSD 类（Julien Pommier 2013，四条款） | `aec3r1-vendored-pffft-LICENSE.txt` |
| vendored rnnoise | `webrtc/third_party/rnnoise/COPYING` | BSD 类（Mozilla/Jean-Marc Valin 2017，三条款） | `aec3r1-vendored-rnnoise-COPYING.txt` |

**再分发义务（本产品的）**：BSD-3 条款只要求分发时保留版权＋许可声明＋disclaimer——
与 D17（代码签名与分发）不冲突；abseil Apache-2.0 同理（保留 LICENSE 与 NOTICE）。以上 8 个
license 文件已逐字入 `docs/evidence/s1/`（`⛔ 记忆抄＝伪造`的对照格：本腿全部来自真实 tarball
原位拷贝，非凭记忆转写）。

## 2. fetch 读数（版本钉死＋sha256 记账）

| 件 | 版本 | 来源 URL | sha256 |
|---|---|---|---|
| webrtc-audio-processing-2.1.tar.xz | **2.1**（freedesktop 官方 release，NEWS：WebRTC **M131** 基线；"forward port some missing patches to fix Windows builds"） | `https://freedesktop.org/software/pulseaudio/webrtc-audio-processing/webrtc-audio-processing-2.1.tar.xz` | `ae9302824b2038d394f10213cab05312c564a038434269f11dbf68f511f9f9fe` |
| abseil-cpp-20240722.0.tar.gz | **20240722.0** | `https://github.com/abseil/abseil-cpp/releases/download/20240722.0/abseil-cpp-20240722.0.tar.gz` | `f50e5ac311a81382da7fa75b97310e4b9006474f9560ac46f54a9967f07d4ae3` |
| abseil wrapdb patch | 20240722.0-3 | `https://wrapdb.mesonbuild.com/v2/abseil-cpp_20240722.0-3/get_patch` | `12dd8df1488a314c53e3751abd2750cf233b830651d168b6a9f15e7d0cf71f7b`（＝wrap 内 `patch_hash` 逐字） |

三个哈希都是**下载后现算**、与官方清单对照：前两个记入 `deps.toml`
（`[webrtc-audio-processing]`＋`[webrtc-audio-processing.abseil]`）；wrapdb 补丁哈希与 release
自带 `subprojects/abseil-cpp.wrap` 的 `patch_hash` 逐字吻合（该补丁只提供 abseil 的 meson 构建
文件，本腿源单即从它提取）。

### 2.1 版本选择理由（铁律：主流＋最有价值）

- **主流**：freedesktop/pulseaudio 打包是 webrtc AEC3 的主流 Linux 打封装（发行版即用它），
  2.1 是其**最新稳定 release**（官方索引：2.1→2.0→1.3→…）。
- **最有价值**：2.0 起 WebRTC 基线升到 M131（"Various improvements to the AEC implementation"），
  2.1 的 release notes 逐字写明修 Windows 构建——对本腿（Windows/mingw）是最有利组合。
- **候选 1.3 落选理由**：基线旧（Chromium 88 世代）、AEC3 实现旧；非首选。
- **⛔下载边界的边界解读（显式申报，供机主否决）**：2.1 与 1.3 的 release tarball 都**不**捆绑
  abseil——构建依赖由 release 自带 wrap 清单钉死（版本＋sha256 逐字在案）。本腿把 abseil 视为
  webrtc-audio-processing fetch 的组成部分（deps.toml onnxruntime 先例：主依赖的钉死组成件，
  随主条目记账），**未下载任何 wrap 清单之外的包/工具**；meson/ninja 等构建工具一个都没装
  （见 §构建）。

## 3. fetch 落点

`D:\tmp\wisp\20261031\aec3r1\`（仓外，派单允许的落点；`third_party/` 已 gitignore 但本腿取仓外
以避免 mingw 构建撞仓库路径里的空格）。目录形状：tarball×2＋wrapdb patch＋两棵解包源树＋
`build\`（目标文件与 `libwebrtc-audio-processing.a`）＋`build.py`（构建驱动，见 remeasure 件）＋
各轮构建/重测日志。仓内只落：deps.toml 条目＋license 拷贝＋Go 绑定＋gated 测试＋本件。

## 4. 未定义即停三门的现值

1. license 非 BSD 类 ⇒ **不触发**（§1，BSD-3 逐字现核）。
2. mingw 构建诚实努力后仍不可行 ⇒ **不触发**（构建成功，读数见 `aec3r1-remeasure.md`）。
3. rig 在 AEC3 链下读数发散 ⇒ 重测件里如实记（判读留给编排者/59-v2）。
