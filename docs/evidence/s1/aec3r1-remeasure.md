# AEC3 接入重测（腿 aec3-r1）— 构建读数与 PRECHECK P15 重裁材料

日期 2026-10-31 · 腿 `aec3-r1`（AEC3 接入重测腿，写）· 起手锚 `0bb87859c`（dev；开工三核过）。
接续件：license/fetch 判定在 `aec3r1-license.md`；基线＝`59r1-aec-spike.md`（ERLE 16.8dB、
自回声 corr 0.978、ASR 逐字转写、corr 判别器阈值 0.5）。
**PASS/NO-GO 判语留给编排者/59-v2——本件只给读数与建议。**

## 0. 一句话结论

**webrtc AEC3（2.1）在本腿钉定工具链上构建成功并已插进 59-r1 同一 rig**：同一把 ERLE 尺、
同一 fixture 上，基线 16.8dB（单增益消参考替身）→ **AEC3 24.5dB（59-r1 原 fixture）/
30.7dB（同日新采）**；**ASR 对回声的逐字转写消失**（「你好我是威斯普今天天气」→「哦」→
新采臂→空串）；corr 判别器读数上 raw 0.978/0.979 → AEC 输出 −0.012/+0.023（远低于 0.5 阈）。
判据 ≥30dB **两臂跨线**（一臂达标、一臂 24.5dB），且判别器的保护角色在真 AEC 之后移
（见 §4.3）。重裁所需的全部数字在 §3。

## 1. 环境与工具链（与 59-r1 钉定一致）

- Go 1.27.1（CGO_ENABLED=1）；C/C++ 工具链＝**MSYS2 mingw-w64 GCC 16.2.0 (Rev3, MSYS2
  project)**（`E:\work\base\msys64\mingw64\bin\{gcc,g++}.exe`——59-r1 已证 cgo 面可用）。
- **meson/ninja 不在场**（mingw python 3.14.7 无 meson 模块）且下载边界不许补装工具 ⇒ 构建
  走**手搓驱动**：源单逐字取自 release 自带 meson.build（webrtc 侧）与官方 wrap 钉死的 wrapdb
  abseil meson.build（abseil 侧），编译参数逐字照抄根 meson.build 的 windows/x86_64 分支
  （`-DWEBRTC_WIN -D_WIN32 -D__STDC_FORMAT_MACROS=1 -DNOMINMAX -D_USE_MATH_DEFINES
  -DWEBRTC_LIBRARY_IMPL -DWEBRTC_ENABLE_SYMBOL_EXPORT -D_WINSOCKAPI_ -DNDEBUG
  -DWEBRTC_ENABLE_AVX2 -DWEBRTC_APM_DEBUG_DUMP=0`，C++17，SSE2 文件 `-msse2`、AVX2 文件
  `-mavx2 -mfma`，链接 `winmm`）。驱动脚本＝仓外 `D:\tmp\wisp\20261031\aec3r1\build.py`。
- 音频错峰：每次 rig 前现核进程表（go/aecspike/wisp 测试进程零在场；桌面常驻 wisp-* 非测试
  进程，rig 自带 runaway/idle 守卫与 ratio_ppm 体检可暴露污染）。

## 2. 构建读数（成功＝产物路径）

**产物**：`D:\tmp\wisp\20261031\aec3r1\build\libwebrtc_audio_processing.a`（3,747,614 字节；
**380 个目标文件**归并＝abseil 12 源集＋webrtc 全部源单；`ar rcs`）。全程 3 轮：

1. 第 1 轮 **FAIL 59 文件**——卡点：我的驱动脚本把 base 目录比较写错，webrtc 文件被喂了
   abseil 参数表（缺 `-DWEBRTC_WIN`），全部失败都是这一个的级联。修复在驱动（仓外），
   源树零改动。
2. 第 2 轮 **FAIL 2 文件**——卡点＝**GCC 16 libstdc++ 不再传递包含 `<cstdint>`**：
   `rtc_base/trace_event.h`（`uint8_t`）与 `aec3/multi_channel_content_detector.h`（`int64_t`）。
   修复＝编译旗标 `-include cstdint`（只加旗标、**源树零改动**，provenance 干净）。
3. 第 3 轮 **rc=0**：381 jobs → 380 objects → 归并 .a，**177s**（12 并行）。
   （380 vs 381 的 1 差额＝重复源名归并，链接面由 cgo 链接全绿背书。）

- MinGW 官方 patch（`patches/0001-Some-fixes-for-MinGW.patch`）**已在 2.1 树内**（NEWS 逐字
  "forward port some missing patches to fix Windows builds"；两处修复——SEH 去 mingw 化、
  `windows.h` 小写——现核在源内），无需再打。
- **cgo 绑定链接全绿**：`go build -tags webrtcap ./internal/audio/aec/` rc=0（shim 编译＋
  静态库符号全解析）；无 tag 默认面 rc=0（stub）；`GOOS=linux go vet` rc=0。

## 3. 重测读数（同 rig 同 fixture；判据 ≥30dB）

rig＝59-r1 同一 `RunLoopback`（共享模式 render＋端点 loopback 采集、音量 0.5、同一 TTS
fixture「你好，我是威斯普，今天天气不错。」）；AEC3 链＝far→`ProcessReverseStream`、
near→`ProcessStream`，10ms 帧锁步，`set_stream_delay_ms(30)`＝实测 lag 493 样本@16k 取整；
float32 接缝；**AGC1/AGC2/HPF/NS 全关**（官方离线示例开着它们，但增益器会重标定残差、污染
ERLE——偏离示例的原因在此如实记账）；**每臂冷启动全新 engine**（首跑曾两臂共池，第二臂继承
已收敛滤波器＝方法学污染，已修并重跑，本节数字全部来自冷启动终版 run4）。

| 读数 | fixture-59r1（=基线那条采捕） | fresh-rig（同日新采） |
|---|---|---|
| 对齐 lag / strength | 493 样本 / 0.9895（＝59-r1 原值复现） | 493 / 0.9895 |
| **ERLE：AEC 关（raw）** | 0dB（定义） | 0dB（定义） |
| **ERLE：单增益消参考（59-r1 替身）** | **16.79dB（＝基线 16.8dB 逐位复现⇒仪器链验证）** | 16.82dB |
| **ERLE：AEC3 开** | **24.49dB** | **30.69dB** |
| AEC3 实时 CPU | 0.0138 CPU-秒/音频-秒＝**1.38% 单核** | 0.01387＝1.39%（批处理 burst 利用率 104.8%/137.5% 是离线跑的瞬时值，实时成本以 CPU-秒/音频-秒为准） |
| ASR（产品同源 streaming ASR）raw | 「你好我是威斯普今天天气」（逐字＝59-r1 基线复现） | 同左 |
| **ASR 对 AEC 输出** | **「哦」**（逐字转写消失，残一音节伪迹） | **「」空串**（完全转写不出） |
| VAD 起点数 raw / AEC 输出 | 1 / 1 | 1 / 1 |
| corr(near, far_ref) raw / AEC 输出 | 0.978 / **−0.012** | 0.979 / **+0.023** |

读数判读（材料，非判语）：

1. **仪器链自证**：59-r1 原 fixture 上 lag=493/strength=0.989/单增益 ERLE=16.79dB 与基线记录
   三值全对上 ⇒ 本腿测的 AEC3 数字与 16.8dB 基线在同一把尺上可比。
2. **AEC3 对基线的净提升＝+7.7dB（原 fixture）～ +13.9dB（新采）**；判据 ≥30dB **两臂跨线**：
   新采臂达标（30.7dB），原 fixture 臂 24.5dB 未达。两臂唯一差异是采捕数据本身（同协议同
   参数同引擎），提示 AEC3 冷启动 3.4s 短窗内的收敛速度对信号细节敏感；产品是连续流会话
   （分钟级），其稳态 ERLE 只会高于冷启动短窗——但这句是推断，**本 rig 没有长窗读数**。
3. **ASR 逐字转写消失（本腿问句的直接答案：消失了）**：59-r1 的「回声被逐字转写＝Path C 无
   AEC 时必自触发」的实锤，在 AEC3 后变成「哦」（原 fixture）→ 空串（新采）——残差已低于
   产品同源 ASR 的可转写面（新采臂完全不可转写）。
4. **corr 判别器（判据 corr<0.5）**：AEC 输出残差与参考的相关 |corr|≤0.023，深低于 0.5 阈。
   但要如实指出其语义变化：**该判别器的原始负载是「把自回声从外部源里分出来」；AEC 把残差
   打到 |corr|≈0 之后，残差在 corr 轴上与外部源不可分**——保护职责移到「残差能量低于
   VAD/ASR 敏感面」。而本 rig 的 VAD（silero，与 59-r1 同仪器同复用形态）在 AEC 输出上仍
   起 1 次/3.4s ⇒ **「barge-in 自触发是否消失」不能只靠 corr 回答**，需要组合判据（VAD∧
   能量门∧corr）的产品级读数；本 rig 无该组合判据的现场读数。如实记入重裁材料。
5. VAD 复用形态＝59-r1 原样（同一 Detector 实例先后喂 raw 与处理流），其状态跨喂的边界
   与 59-r1 相同，未引入新偏差。
6. 本 rig 全程数字回线（无房间混响）⇒ 全部 ERLE 是真实现场**上界方向**的读数（59-r1 同注）。

## 4. 四数／两尺／d22scan 导出树＋正控（第 21 处）

- **四数（`./internal/audio/...` 全集，含 PATH 前缀的文档化跑法）**：`-count=1 -v` ⇒
  **RUN=71 PASS=53 FAIL=0 SKIP=6**（SKIP 逐枚＝TestLiveWasapiSmoke live 硬件门＋本腿
  TestSpikeAec3ErleRemeasure＋59-r1 四枚 rig 门控测试；rc=0）。⚠ bare `go test`（无 PATH
  前缀）在 aecspike 测试二进制**加载期** 0xc0000135（speech 的 sherpa DLL 动态依赖，
  59-r1 时代既有形状，非本腿引入——aec 包 stub 面纯 Go 零 DLL；文档化跑法带 PATH 前缀即绿）。
  另记同口径单包数（对照 59-r1 的 52/39/0/1）：`./internal/audio/` ⇒ RUN=57 PASS=44 FAIL=0
  SKIP=1（+5 全是两腿之间其它腿提交的 internal/audio 测试，本腿零新增该包测试）。
- **两尺**：gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`）`-l internal/audio/`＝空；
  `go vet` windows（无 tag＋`-tags webrtcap` 两面）＝rc=0；GOOS=linux vet＝rc=0。
- **d22scan（Git Bash `scripts/d22scan.sh`）**：实扫 **clean**；导出树＝bans#1-5 internal/=283、
  cmd/=49、ban#6 frontend/=116、ban#7 internal/tools/=25、ban#8 design/=39、frontend/=116、
  internal/=719、cmd/=171（与 59-r1 的 280/48/…/715/169 同形，小额增长属两腿间其它腿提交）。
- **正控（咬合记录，如实报）**：selftest 抓到本腿 `aec.go:10` 的 **[phantom-citation]
  （ticket 212 ban #9）**——注释引用本证据件而文件尚未落盘；证据件落盘后复跑转绿
  （`TestScannerSelfScanOfRealRepoIsGreen`/`TestRealRepoLedgerIsHonest` 两枚由红转绿）。
  自测总分＝**PASS=35 FAIL=0 SKIP=0（RUN=77）**。
- ⚠ **新 cgo 件 ban 射程如实报**：本腿新增 `wap_shim.cpp/.h`（C++）不在 d22scan ban#8 的
  Go 文件扫描射程内（仪器只扫 Go 文件）；两文件内容为 ASCII 注释＋代码，无 emoji/数学带
  字符，如实备注而非依赖仪器。
- 路径长度门：本腿写面最长相对路径 `docs/evidence/s1/aec3r1-webrtc-audio-processing-2.1-webrtc-LICENSE.txt`
  （60 字符＜100 帽），不涉名册。

## 5. PRECHECK P15 重裁材料汇总（判语摆编排者/59-v2）

59-r1 NO-GO 的三条支柱在 AEC3 落地后的现值：

| 59-r1 支柱 | 本腿后现值 |
|---|---|
| ① AEC 引擎三候选全灭（乙不可 fetch） | **翻**：机主批准 fetch（A1311③）⇒ 乙（webrtc-audio-processing 2.1）已 fetch/构建/接入，license BSD-3 现核（`aec3r1-license.md`） |
| ② 无 AEC 时 Path C 必自触发（ERLE 16.8、ASR 逐字转写） | **翻**：AEC3 后 ASR 逐字转写消失（「哦」→空串）、ERLE 升到 24.5/30.7dB、残差 corr≈0 |
| ③ 判据 ERLE≥30dB | **跨线**：新采臂 30.7dB 达标、59-r1 原 fixture 臂 24.5dB 未达；冷启动短窗 vs 产品连续流的收敛差是未测缺口 |

悬而未决（重裁前值得知道的缺口，本腿射程内不补）：
- **VAD 在 AEC 输出上仍起 1 次/3.4s** ⇒ barge-in 组合判据（VAD∧能量∧corr）的产品级读数缺失；
- 长窗稳态 ERLE（产品连续会话形状）无 rig 读数；
- 真跨端点（蓝牙音箱＋有线麦）依旧本机不可测（59-r1 §5 原样）；
- SPEC-10 Path C 行的 CPU 回填：AEC3 引擎本体 **1.38% 单核**（远低于 5% 暂定值），但 VAD
  26%/遍（59-r1 §7）仍是该行的大头——回填材料齐了，动作仍摆票 13/机主。

## 6. 写面申报（逐枚）

仓内（全部 commit pathspec 逐名列出）：
1. `deps.toml` —— 追加 `[webrtc-audio-processing]`＋`[webrtc-audio-processing.abseil]` 两节
   （票 13 形状；其它条目一字未动）。
2. `internal/audio/aec/aec.go` —— 绑定 API 面（ErrUnavailable/Apm/FrameSamples）。
3. `internal/audio/aec/aec_stub.go` —— 非原生构建形状（`!windows || !cgo || !webrtcap`）。
4. `internal/audio/aec/aec_webrtc_windows.go` —— 原生绑定（`windows && cgo && webrtcap`）。
5. `internal/audio/aec/wap_shim.h` ／ 6. `wap_shim.cpp` —— C ABI shim（同 tag；webrtc 类型
   全部封在 C++ 侧）。
7. `internal/audio/aecspike/aec3_windows_test.go` —— 门控重测测试（WISP_AEC_SPIKE=1）。
8. `docs/evidence/s1/aec3r1-license.md` ＋ 9. `aec3r1-remeasure.md`（本件）。
10-17. `docs/evidence/s1/aec3r1-*.txt` 八枚 license 逐字拷贝（清单见 license 件 §1）。

仓外（`D:\tmp\wisp\20261031\aec3r1\**`，只建不删）：两 tarball＋wrapdb patch＋两棵源树＋
`build.py`＋`build\libwebrtc_audio_processing.a`＋`run\`（测试 exe＋三 DLL）＋`out\`
（两臂 JSON＋AEC 输出 wav）＋全部轮次日志（build-run1/2/3.log、remeasure-run1/2/3/4.log、
four-numbers*.log、d22scan.log）。音频 artifact＝合成 fixture＋自渲染 TTS＋其回声，无用户
语音（D16）。

## 7. 复跑处方（逐字）

fetch+构建：

```
curl -L -o webrtc-audio-processing-2.1.tar.xz https://freedesktop.org/software/pulseaudio/webrtc-audio-processing/webrtc-audio-processing-2.1.tar.xz
# sha256 = ae9302824b2038d394f10213cab05312c564a038434269f11dbf68f511f9f9fe
curl -L -o abseil-cpp-20240722.0.tar.gz https://github.com/abseil/abseil-cpp/releases/download/20240722.0/abseil-cpp-20240722.0.tar.gz
# sha256 = f50e5ac311a81382da7fa75b97310e4b9006474f9560ac46f54a9967f07d4ae3
tar -xf <两个包>；E:\work\base\msys64\mingw64\bin\python.exe build.py   # 产物 build\libwebrtc_audio_processing.a
```

重测（59-r1 处方＋webrtcap 两处扩展）：

```
cd <repo>
set CGO_CXXFLAGS=-I<D:\tmp\wisp\20261031\aec3r1\webrtc-audio-processing-2.1\webrtc> -I<D:\tmp\wisp\20261031\aec3r1\abseil-cpp-20240722.0>
set CGO_LDFLAGS=-L<D:\tmp\wisp\20261031\aec3r1\build>
go test -c -o D:\tmp\wisp\20261031\aec3r1\run\aecspike.test.exe -tags webrtcap ./internal/audio/aecspike/
copy third_party\sherpa-onnx\{onnxruntime.dll,sherpa-onnx-c-api.dll,sherpa-onnx-cxx-api.dll} D:\tmp\wisp\20261031\aec3r1\run\
cd D:\tmp\wisp\20261031\aec3r1\run
set WISP_AEC_SPIKE=1
set WISP_AEC_SPIKE_ROOT=<repo>
set WISP_AEC_SPIKE_OUT=D:\tmp\wisp\20261031\aec3r1\out
aecspike.test.exe -test.run TestSpikeAec3 -test.v -test.timeout 12m
```

（无 tag 四数：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./internal/audio/... -count=1 -v`。）
