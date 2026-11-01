# Path-C 重 key 验收小腿 `26p-k1` —— 实现腿申报件（写腿）

日期 2026-11-01 11:0x-11:2x +08 · 腿 `26p-k1`（写、cmd/wisp 测试面）· 基准 HEAD `24f45b78f`
（dev，开工三核：HEAD≥锚点 ✓、staged 区干净 ✓、写面无他人脏 ✓）· 派单存档
`.scratch/wisp/probes/orch/20261101/txt/26p-k1-dispatch.txt`。

**裁定溯源**：26p-v2 §3/§4 终判（有条件 PASS 条件①）＋59-v3 重开条件——(a) 验收改
**episode keyed**（K=5/N=1 起步）：**5 窗至少 1 火＋凡火必 ≤400ms＋每窗深压读数入档**；
「加稳」归 28-PathC/393 Reset 轨；「放宽」驳回。本腿只把该裁定落成测试形状，零产品码、
零契约文字。

## 0. 一句话

**(a) 案 keyed 改写落成并三发全绿（火数 2/5、4/5、3/5，凡火 28..133ms 全带内）；反向判据
闭合——overlay 改回静态五轮断言三发全红（红形 `:354` 超时与 26p-v2 引的逐字同形）＋还原核
哈希与复跑双核；零回归（26p 其余 4/4＋六族 33/2/0 两红单跑全归因）；四数/两尺/d22scan
全对上。**

## 1. 改了什么（只改断言形状；窗材料/时序零动）

`cmd/wisp/resident_voice_bargein_26p_windows_test.go` 唯一写面。diff 两个区：
文件头 (a) 从句（:6-8）与 `Test26PBargeinStopsWithin400ms` 全案（:293-490）。
`Test26PFeedOnlineMatchesWindowReplay`（:499 起）与 (b)(c) 件（echo 测试文件）零触碰。

| 旧（静态，已退） | 新（keyed，26p-v2 §3） |
|---|---|
| 每轮 `Triggers != 1` 即红（「每轮必火」） | 火/不火都合法落窗：`fired = ttsAborted 移动`（:390，392 `voice392Capture` 同款 either-or 等待 :387） |
| 每轮断 vadToStop ≤400ms | **凡火必 ≤400ms**（火窗逐窗断言 :432，原字面原样）；不火窗不断时延 |
| 无档（只有 summary 行） | **每窗读数入档**：`26p-a-keyed.json`（:483）逐窗 fired/vad_ms/episode rms/阈值/tighten＋**不火窗逐 episode 深压读数**（rms/阈值/tighten/energy_pass/suppressed）＋每窗 HELD 行 t.Logf（:463）——FINDING MOVED 牙式「测试即档案」 |
| —（隐含「全火」） | **K=5/N=1 地板**（:479）：`FINDING MOVED: 0/5 windows fired - the keyed floor (K=5/N=1) is red`——五窗全不火＝红（N 地板防空转死），消息里把「静态五轮全火断言已退役、不许回摆」钉死 |
| `Triggers != 1` 打包两义 | **零自激字面独立保住**：`Triggers > 1` 即红（:400，S4「不自激」逐 episode 原样）；火窗 `Triggers==1`＋passing episode 必在（原两断言原样）；被打断播报不计 spoken（原字面，改相对形）；不火窗播报恰计 +1 |

产品码（bargein.go/gate.go/align.go/monitor）**零字节**——`git diff --stat` 全程只有本测试件
（`.scratch/wisp/issues/30-result-routing-d10-done.md` 的 M 是开工前既有他人脏件，本腿不碰
不提交）。

## 2. 读数（staged 全链；全部本腿现跑）

staged＝A1326 三 DLL 贴测试二进制＋WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_TTS_DIR/
WISP_KWS_DIR/WISP_15E_TMP/WISP_43_REPO_ROOT＋cwd=包目录＋unset GOTMPDIR；webrtcap 二进制
`go test -c -tags webrtcap`（CGO 宏集照 aec3r1 build.py root-meson windows 形全量、`D:/`
盘符形 393-v1 坑）；Git Bash runner（390-v1 坑，`D:\tmp\wisp\20261101\26pk1\run.sh`）；
每跑前进程表现核（在飞测试零在场全过窗）。

### 2.1 keyed 三发（pristine 二进制 `wisp26pk1.test.exe`，SHA 前缀 a217b79cbb54a983）

| 发 | 火数 | 凡火 vad→stop | 深压窗（不火） | 判 |
|---|---|---|---|---|
| 1 (11:04) | **2/5** | R1 66ms @-24.64、R2 123ms @-20.39 | R3 -40.23、R4 -42.81、R5 -39.49 | PASS 37.4s |
| 2 (11:06) | **4/5** | R1 124ms @-27.18、R3 42ms @-20.47、R4 92ms @-23.39、R5 47ms @-27.44 | R2 -40.75 | PASS 27.3s |
| 3＝还原核 (11:10) | **3/5** | R1 28ms @-20.44、R3 133ms @-19.89、R4 55ms @-20.52 | R2 -44.43、R5 -43.42 | PASS 32.5s |

- **凡火 9/9 全部 28..133ms**——对照 26p-v2 唯一火 24ms、26p-r1 40-99ms：同带内，无越 400ms 反例。
- **火数 2→4→3 摆动**＝被重 key 显式化的受测维度，与 392-v2「摆动脆弱」主发现同形；两发非全塌
  （非派单停手条件②的 NO-GO 材料）。深压带本腿 -39.49..-44.43，落 392-v2/26p-v2 记的
  -39..-45 带（见 §2.2 静态发还测得 -50.39，摆幅再向下扩——如实记档）。
- 对齐面每窗健康（lag 29-40ms、strength 0.27-0.83、ready=9728）——与 26p-v2「失败全落引擎
  双讲抑制段」同源读数。
- 工件：`D:\tmp\wisp\20261031\26pr1\out\26p-a-keyed.json`（尾发读数；前三发逐窗在
  `26pk1\logs\26p-a-keyed-run{1,2}.log`）。26pr1\out 原 26p-v2 四件先备份
  `26pk1\out-backup-26pr1\` 再按既定处置让本腿读数落同根。

### 2.2 反向判据：改回静态五轮断言 ⇒ 摆动敏感复现（三读数＋还原核，仓外 overlay）

突变件＝HEAD 原文（`git show` 出的静态断言形状）仓外 overlay
（`26pk1\mut\overlay-static.json` → `-overlay` 构建独立二进制
`wisp26pk1-static.test.exe`，SHA 前缀 8fe6032b54e1253d ≠ pristine ≠ 实现腿件）：

| 发 | 火形 | 深压 | 判 |
|---|---|---|---|
| 1 (11:07) | R1 火 -21.37 | R2 深压 **-39.50** 不触发 | **红**：`timed out waiting for round 2: the barge-in stop to land`（:354——与 26p-v2 引的红形逐字同形） |
| 2 (11:08) | R1 火 -23.15 | R2 深压 **-50.39** 不触发 | **红**（同 :354） |
| 3 (11:09) | R1-R2-R3 火 -22.28/-22.82/-23 | R4 深压 **-40.22** 不触发 | **红**：`timed out waiting for round 4`（:354） |

**静态 3/3 全红、keyed 3/3 全绿**——「每轮必火」在本面不可复现（392-v2/26p-v2 主发现）由本腿
第三次独立复证；keyed 把摆动从掷币变显式受测维度且不丢任何 S4 硬界。**还原核**：overlay 从未
触碰工作树（pristine 文件 SHA256 8A546FC8…F8EF72 前后逐位同）＋pristine 二进制复跑 PASS
（§2.1 发 3）。

## 3. 零回归

- **26p 套件其余 4/4 PASS**（pristine 二进制）：(b) 19.0s、(b)-M2 36.5s、(c)+M1 44.7s、
  首格等价 15.7s（6/6 窗在线 Feed==离线 ObserveWindow 零漂移，lag=628/0.74）。
- **六族抽案 staged 合跑 33 PASS / 2 FAIL / 0 SKIP（35 RUN）**——总数与 26p-v2 基线逐位同；
  两红 `Test26R3AnnounceClosesTheCaptureGate`（10.5s）与 `Test384ResidentWatchdogEndToEnd`
  （3.3s）**单跑全 PASS**＝合跑负载时窗归因（A1329 先例形状；红名集合零新增）。
- 四数（staged 面累计）：keyed (a) 三发 3 RUN 3 PASS；26p 其余 4 RUN 4 PASS；六族 35 RUN
  33/2/0（两红单跑 2 RUN 2 PASS）；反向静态发 3 RUN 3 FAIL（**预期红**＝反向判据本体）。
  预期红与意外红零混淆：全部 FAIL 都来自突变二进制、全部 pristine FAIL 零枚。

## 4. 两尺＋d22scan（正控）

- **gofumpt** `-l`（`D:\work\base\gopath\bin\gofumpt.exe`，cmd/wisp＋internal/audio 全树）
  ＝0（rc=0 listed=0）。
- **vet** 双形 rc=0：no-tag＋`-tags webrtcap`（CGO env 同构建）。
- **d22scan**（`tools/d22scan` 独立模块，门形 `runtests.sh -C tools/d22scan ./...`＋
  `go run . -root <repo>`）：**正控 selftest 35/0/0（RUN=77）**与 26p-v2 逐位同；扫描
  **clean rc=0**，四数 bans#1-5 internal/=287、cmd/=52、ban#8 internal/=730、cmd/=183
  ——对 26p-v2 基线**逐位对上、零增量**（本腿只改既有测试件，无新文件；design/=39 为机主
  外部面文本数漂移，非本腿面、非违规项）。

## 5. 复跑处方（逐字）

```
# 构建（Git Bash；CGO 宏集=aec3r1 build.py root-meson windows 形全量；D:/ 盘符形）
cd <repo> && unset GOTMPDIR
export CGO_CXXFLAGS="-ID:/tmp/wisp/20261031/aec3r1/webrtc-audio-processing-2.1/webrtc -ID:/tmp/wisp/20261031/aec3r1/abseil-cpp-20240722.0 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_ENABLE_SYMBOL_EXPORT -D_WINSOCKAPI_ -DNDEBUG -DWEBRTC_WIN -D_WIN32 -D__STDC_FORMAT_MACROS=1 -DNOMINMAX -D_USE_MATH_DEFINES -DWEBRTC_ENABLE_AVX2 -DWEBRTC_APM_DEBUG_DUMP=0 -include cstdint"
export CGO_LDFLAGS="-LD:/tmp/wisp/20261031/aec3r1/build"
go test -c -o D:/tmp/wisp/20261101/26pk1/bin/wisp26pk1.test.exe -tags webrtcap ./cmd/wisp/
# 三 DLL（third_party/sherpa-onnx/）copy 进 bin\
# 跑（run.sh：WISP_* env＋cwd=cmd/wisp＋unset GOTMPDIR＋-test.v -test.timeout 30m；每跑前进程表现核）
sh D:/tmp/wisp/20261101/26pk1/run.sh D:/tmp/wisp/20261101/26pk1/bin/wisp26pk1.test.exe \
  Test26PBargeinStopsWithin400ms <log-name>
```

⚠ 摆动脆弱如实申明：火数按发摆（本腿 2/5、4/5、3/5），凡火带内是稳定保证；单发全塌
（0/5）才会撞 keyed 地板（那是要上报的 NO-GO 材料）。

## 6. 写面申报与仓外件

- 仓内写面＝`cmd/wisp/resident_voice_bargein_26p_windows_test.go`（本件为唯一码面改动）＋
  `docs/evidence/s1/26pk1-impl.md`（本件）。禁 push/amend/reset/rebase/stash/checkout；
  commit 显式 pathspec、主题逐字 `Path-C 重key验收小腿 26p-k1`、无 `# ` 前缀。
- 仓外 `D:\tmp\wisp\20261101\26pk1\**`（bin＝pristine 二进制＋三 DLL；mut＝静态突变件＋
  overlay json＋突变二进制＋SHA 记录；logs 全量 12 件；out-backup-26pr1＝26p-v2 四件原读数
  备份；run.sh；cases），只建不删。
- 测试工件落 `D:\tmp\wisp\20261031\26pr1\out\`（测试硬码根，26p-v2 同款处置；旧件已备份）。
- ⛔ 未碰：产品码零字节、(b)(c)/等价案零字节、fixture/时序/窗材料零字节、thresholds/golden/
  testdata 仓根/票面/台账/specs/PLAN/frontend/design。
