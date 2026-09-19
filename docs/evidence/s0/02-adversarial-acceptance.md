# T02 对抗验收报告 — 02-s0-spike（五项实测 + X/Y 判定 + SLO 回填）

- 验收角色：T02-adv（对抗验收，独立上下文，只读仓库 + 临时目录实验）
- 验收时间：2026-09-19T09:00Z ~ 09:40Z
- 被验收基线：`e081619`（T02 相关 commit：`dd4be2a` `91eb04d` `c7df446` `e44fa7c`
  `45212b5` `8c5a72c` + 编排者路由 `da81a6a`/`6160bc8`/`e081619`）
- 验收环境：Windows x64（build 26100）· Git Bash · 模型已缓存于 `third_party/spike-models/`
- 复跑方式：全部用仓库内已构建二进制（`scripts/spike/bin/*`）以 stdout 模式复跑
  （不带 `-out`），**未覆盖任何证据 JSON**；临时输出写 `$TEMP/t02adv-*.json`。

## 逐项裁决

### 1. 证据链完整性 — PASS

- `docs/evidence/s0/json/` 恰 13 个 JSON，与报告/`run.ps1` 的工件清单一一对应
  （01×1、02×1、03×2、04×1、05×1、06×2、07×4、08×1）。
- 抽 5 个关键值与原始 JSON 逐字节核对，全部一致：
  ① 空 Go 6.890625MB（报告 6.89）✓；③ KWS session 48.60546875MB（报告 48.61）✓；
  idle-y 16.453125MB（报告 16.45）✓；切换 P50 4726ms（`08-switch.p50AsrToTtsMs`）✓；
  TTS 常驻 174.1953125MB（报告 174.20）✓。webview2 两轮（879.668/1125.71 cold P50、
  25.723/71.359 hot P50、955.745/858.95 recreate P50）亦逐项吻合。
- 模型清单 URL+SHA256 全部登记；抽 3 个实测核对（`sha256sum` 直接算缓存文件）：
  `silero_vad.onnx` = `9e2449e1…cbb1fd6` ✓ · `vocos-22khz-univ.onnx` = `0574a135…4c0081` ✓ ·
  KWS encoder int8 = `dd784973…e447f4d` ✓ —— **与报告 §5 逐字符一致**。
- 慢测量（常驻/切换）未重跑（按验收指示），做内部一致性核对：切换三轮算术自洽
  （56+4370=4426 / 69+4657=4726 / 106+5728=5834）；常驻值 ≥ session 基线
  （KWS 48.64 ≥ 48.61、ASR 281.13 ≈ unload-test 281.41）；"峰值=max 非 sum" 有数据支撑
  （切换期 WS 最大 290.7MB、峰值 commit 314.3MB；若 sum 应 ≈470MB）✓。

### 2. 测量口径真实性 — PASS

- `scripts/spike/common/mem.go` 逐行核实：`NtQuerySystemInformation(SystemProcessInformation=5)`
  真实调用 ntdll，x64 偏移正确（`WorkingSetPrivateSize`@0x08、`UniqueProcessId`@0x50，
  与内核布局一致）；扫描缓冲 `VirtualAlloc`/`VirtualFree` 真实存在（防 Go 堆污染）；
  shared=WorkingSet−private、commit=PrivateUsage、峰值=PeakPagefileUsage 均如实实现。
  "与任务管理器私有工作集同源" 的声称技术上准确（同读 `WorkingSetPrivateSize` 内核字段）。
- `xy-verdict/main.go` 核实：回落判据真实实现（pre+2MB 阈值、250ms 采样、默认 15s 窗口，
  `04-unload-asr.json` 60 样本平坦 22.19→22.38MB、`settled=false`）——15s 不回落是
  10s 不回落的重言证明，FAIL 判定成立。session_cgo.go 是真 sherpa-onnx Go 绑定调用
  （OnlineRecognizer + Paraformer int8 + NumThreads=1 + 1s 音频解码）。
- **cgo 崩溃概率引用独立核实**：sherpa-onnx issue #2694（open，Go 绑定在
  `NewOfflineRecognizer`/session 创建内偶发 SIGSEGV，"signal arrived during cgo execution"）
  与 #3635（open，VAD CircularBuffer int32 溢出，2^31 样本 ≈ 37h16m 后段错误）**真实存在
  且描述准确**，非编造引用。

### 3. 抽复现（快测量三项全做） — PASS

| 项 | 证据声称 | 本次独立复跑 | 判定 |
|---|---|---|---|
| idle-y 基线 | 16.45/16.59MB | **16.49MB**（handles 414/threads 23，与证据 413/23 一致） | ✓ ±10% 内 |
| goja Interrupt | P50 +0.4–0.6ms / max ~12ms | P50 超调 +0.30–1.03ms，本次 max +3.4ms | ✓ 同量级（声称偏保守） |
| goja 能力探针 | async ✓ / modules ✗ / asyncgen ✗ | 逐条一致（`asyncAwaitWorks=true`、`esModulesSupported=false`、`Async generators are not supported yet`） | ✓ |
| webview2 热 show | P50 26–71ms（≤200 PASS） | **P50 20.0ms / P95 45.4ms** | ✓ |
| （参考）recreate | P50 860–956ms | 856.5ms | ✓ |
| （参考）进程内 create | 1163–1809ms | 956ms（暖 profile、负载更低） | ✓ 方向一致 |

- 附带发现：idle-y 首跑被 `Shell_NotifyIconW` shell RPC 挂起 11 分钟（杀进程重跑即正常）
  —— **报告 §4.5 已把该坑如实记录**（"一次实测挂起、重跑正常"），复现即证。
- 见下"轻微问题"第 3 条：idle-y 声称的"两次运行"只提交了一次 JSON。

### 4. SLO.md / PRECHECK.md 回填正确性（对照 PLAN D32/D44） — PASS

- 对照 `docs/PLAN.md` L2245–L2353（16.3.2/16.3.3/16.3.4/16.3.5）与 D44（L2023）：
  - 16.3.2 七行全部处置：仅实测支撑的行带实测值（Sleeping 16.6、Armed 58≈48.6+shell、
    面板/Conversation/CPU 明确标"未测/待票 08"，**未出现无实测却填数**）；
  - **Warm 479MB>350 判 FAIL 如实登记**（"按 D32 原定义 FAIL"，未见静默放宽；
    处置 = D32 16.3.3 预设降级 + Warm 口径修订明示"需编排者拍板"）；
  - **700MB 保留"目标值非验收值"标注**（§2/§3 各一次，满足 D32 16.3.3
    "这一区分必须写进 docs/SLO.md" 的强制要求）；
  - **D47 Path C 注存在**（§6：Conversation = Warm + TTS 增量、S4 回填、下限 479MB、
    "禁止把 Path T 串行数字套用 Path C"——与 PLAN L2304–2308 收 scope 注一致）；
  - D44 的三判定输入（内存/回落/崩溃概率）全部出现在 SLO.md §1 判定表内；
    16.3.5 强制-Y 规则的适用边界处理正确：规则只约束"判 X"，内存已选 Y，
    故回落 FAIL 转为 S1 义务而非翻转判定（SLO §1 与 PRECHECK P1 表述一致且正确）。

### 5. 判定逻辑与路由登记 — PASS

- 串行切换 P50 4726ms > 1600ms → **D32 16.3.3 预设降级引用正确**
  （PLAN L2312："若实测串行加载的切换延迟吃掉 TTS 预算 → 退化为 TTS 常驻 + ASR 按需"；
  SLO §3 强制项 3 + §4 + 编排者裁定 3 一致）。已在票 15、票 26 登记（`da81a6a`/`6160bc8`
  diff 核实，内容与 SLO 数字一致）。
- 回落 FAIL → S1 义务（`debug.SetMemoryLimit`/GOGC、25MB 线复测）已在 SLO §1 义务 1、
  §2 回落行、编排者裁定 2 **及票 12、票 15**（diff 核实）登记 ✓。
- 句柄口径裁定 <600 见 SLO.md 末段 ✓（数字瑕疵见下第 2 条）。
- 票 07（Path Y + 窗口栈句柄）、票 51（goja 结论）路由内容与 SLO/PRECHECK 一致 ✓。

### 6. 越界扫描 — PASS

- `scripts/spike/` 为独立 module（`github.com/CarlosShao/wisp/scripts/spike`，自带
  go.mod/go.sum）；根 go.mod 无 goja/go-webview2 依赖；internal/、cmd/ 无 spike 代码混入。
- T02 六个 commit `--stat` 逐一核对：只触碰 `scripts/spike/**`、
  `docs/evidence/s0/json/**`、`docs/SLO.md`/`docs/PRECHECK.md`/报告、以及
  `.scratch/wisp/issues/02-s0-spike.md`（本票进度日志）——**未改动任何其他票据**
  （票据 07/12/15/26/51 的改动来自编排者 commit，非 T02-impl）。
- emoji：`scripts/spike/**`、报告、JSON、PRECHECK 零命中；SLO.md 含 `⚠ ✓ ✗` 文本符号
  （U+26A0/U+2713/U+2717，非 emoji 码位），与 PLAN.md/specs 既有文风一致，不判违例。
- 密钥：无命中（`tokens.txt` 为模型词表文件名，非凭据）。模型目录在 `third_party/`（gitignore ✓）。

### 7. 票据对照 — PASS

- 五项输出逐条：
  ① 五 RSS 基线（含 ④⑤）✓ ② X/Y 判定 = 路径 Y + 三输入齐备 ✓ ③ goja 能力 + Interrupt ✓
  ④ webview2 冷/热 ✓ ⑤ 常驻表 + 切换延迟 + 峰值=max ✓；机器配置（CPU/RAM/OS）全部入 JSON。
- 验收标准"verdict cited in tickets 07/12/15"：由编排者完成（含 26/51）✓。
- Progress log：6 条 append-only 覆盖全部输出 + 移交 ✓（格式瑕疵见下第 6 条）。
- Status 字段仍为 `ready-for-agent`（见下第 4 条）。

## 轻微问题（均不阻塞，无 BLOCKER/MAJOR）

1. **[MINOR] goja Interrupt 报告表与提交 JSON 非同一轮**：报告 §3.3 表（及 PRECHECK
   "max ~12.2ms"、票 51 "max ~12ms"）来自未提交的一轮；提交的 `05-goja.json` P50 相差
   ≤0.2ms 但极值更小（max 超调 8.25ms；本次复跑 3.4ms）。声称方向偏保守、结论
   （D33/F5 时限 ≥250ms）在任一轮数据下均成立，且报告 §4.8 已声明"绝对值以 JSON 为准"。
2. **[MINOR] SLO.md 句柄数 421/窗口栈 +279 无 JSON 出处**：提交证据为 407/413
   （增量 265/271），本次复跑 414。编排者 <600 裁定不受影响（余量巨大）。
3. **[MINOR] idle-y "两次运行 16.5/16.6" 只有 16.45 一轮入 JSON**；本次复跑 16.49
   构成独立第三点，稳定性疑虑消除。
4. **[MINOR] 票 02 Status 未流转**：仍 `ready-for-agent`/`Claimed by: —`。仓库惯例是
   编排者拥有 Status（8c5a72c "ticket Status untouched (orchestrator owns it)"），
   属流程簿记缺口，非交付缺口。
5. **[MINOR] `01-shell-baselines.json` 的 note 字段陈旧**：写 "QueryWorkingSetEx"，
   实际实现是 NtQSI（mem.go 与报告均正确说明）——纯注释残留。
6. **[MINOR] 进度日志 09:20Z 条目被截断**（结尾 "…D32 fallback TTS-resident" 无收尾/
   无 next=），且 09:40Z 条目紧随其后无空行分隔；内容完整性不受影响。
7. **[MINOR] goja 探针 `es2015-reflect-availability` 期望方向与其名字相反**：
   断言 `typeof Reflect === "undefined"`，实测 got "object"（即 Reflect 存在）→ pass:false
   实为"存在"的证据；报告/PRECHECK 按原始数据正确解读（"Proxy/Reflect 均存在"）。
   throwaway harness 质量已被票据豁免，但后续票引用 JSON 汇总字段时需注意。

## 结论

五项实测、X/Y 判定（路径 Y）、SLO/PRECHECK 回填的真实性、口径与判定逻辑全部独立核实
通过；三快测量复现全部命中证据区间；慢测量数据内部自洽；越界扫描干净。7 个轻微问题
均为簿记/出处标注级别，不改变任何结论。

VERDICT: PASS
