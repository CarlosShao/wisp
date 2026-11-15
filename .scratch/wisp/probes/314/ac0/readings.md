# 票 314 AC#0 读数根（2026-11-15，只读取证腿）

- 腿：314 AC#0（⛔ 非写腿；⛔ 未翻任何 AC 框）
- 证据件：`docs/evidence/s1/314-ac0-forensics.md`
- 原始读数（只建不删）：`D:\tmp\wisp\20261115\w314\`；第 1 发原样副本 `run1-canonical\`
- 发射器（本目录）：`run-direct-canonical.ps1`（canonical 形，`-Suffix` 区分发）、`gate-and-run.ps1`（独占闸＋发车）、`desktop-state.ps1`、`run-direct.ps1`（run A 反面控制）、`run-direct-cwd.ps1`（run B 中间控制）

## 0. 二进制与树

| 项 | 值 |
|---|---|
| 编译 | `go test -c -tags winlive -o D:\tmp\wisp\20261115\w314\wisp-winlive.test.exe ./cmd/wisp` |
| 读数 | rc=0、29.5 s、`63,448,065 B`、SHA256 `1C7F9B324023214601F25B93A24006C7F6FFFBF1A6BFEADAD8896AEEF064029D` |
| 编译时 HEAD | `bf4adaa7b`（跑的中途别的腿推到 `1c19b1953`；读数属 `bf4adaa7b` 的 cmd/wisp） |
| tags | `go version -m` 逐字 `build -tags=winlive` |
| 名册 | `-test.list '.*'` = 669 枚顶层；末枚＝终点锚 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` |
| DLL 暂存 | `onnxruntime.dll`/`sherpa-onnx-c-api.dll`/`sherpa-onnx-cxx-api.dll` 拷到二进制旁，SHA256 与 `third_party/sherpa-onnx/` 同串 |

## 1. 逐发读数

| 发 | 形 | cwd | env/暂存 | rc（十/十六进制） | 耗时 | RUN/PASS/FAIL/SKIP（全行） | 终点锚 | 判 |
|---|---|---|---|---|---|---|---|---|
| A | `-test.v -test.timeout=60m` | **仓根** | 无 | 1 / `0x00000001` | 1219.7 s | 880 / 745 / 121 / 14（顶层 `669/534/121/14`，669 枚全有 verdict） | 命中 | **作废**（不是同一个测试；121 红＝发射形错，见证据件 §1.1） |
| B | 同 A 但 cwd=`cmd/wisp` | `cmd/wisp` | 无 | 无（我 kill，保住"一次一发"） | 13 枚后停 | 13 / 12 / 0 / 0 | — | 作废（只为隔离 cwd 变量，已足够） |
| C | canonical 形 | `cmd/wisp` | 有 | 无（我 kill） | ~1 min | — | — | 作废：发现竞争者在跑（`stage-cli\wisp.test.exe` PID 51336，票 451 腿的载具全量） |
| **1** | **canonical 形**（`-test.count=1 -test.timeout=90m -test.v`） | `cmd/wisp` | 有 | **-1073741819 / `0xC0000005`** | **554.5 s** | **480 / 476 / 3 / 0** | **0 命中** | **★死（native ACCESS_VIOLATION）** |
| **2** | 同 1 | `cmd/wisp` | 有 | **1 / `0x00000001`** | **1389.8 s** | **892 / 861 / 11 / 20** | **命中（PASS）** | **没死，跑完；11 红全属本机噪声族** |

四数口径：RUN＝全部 `=== RUN ` 行；PASS/FAIL/SKIP＝全部 `--- PASS:/FAIL:/SKIP:` 行（含子测试）。
第 1 发 `480 − (476+3+0) = 1` ⇒ 恰 1 枚 unverdicted＝死点那枚。

## 2. 死点（第 1 发）

- 最后一条 `=== RUN`：`TestAC4PriorFocusSurvivesARefusedPanelSample`（`cmd/wisp/panel_resident_windows_test.go:1203`，该文件最后一枚）
- 之后零输出、零 verdict；stdout 最后一行 `wisp: panel window hidden (focus returned to hwnd 0x5701544)`（12:58:20.461）
- stderr 最后一行 `panel thread exited cleanly shows=1 toggles=0 disposals=0`（12:58:20.521）；`disposals=0`
- 进程消失于 12:58:20.5–12:58:21 之间（`Wait-Process` 12:58:21 返回）
- 采样：12:58:16 活着（CPU 70.33 s、WS 115.8 MB、webview2 92）→ 12:58:37 已无（webview2 掉到 37）⇒ 不是"卡住后被外部杀"
- ⛔ 本机没留 WER/dump（Application 日志 0 命中、`CrashDumps` 空、WER 无本时刻新档）⇒ **故障模块点名不了**；盘上唯一 `wisp.test.exe` WER 旧档是 2026-10-22 03:11:47Z（同为 `c0000005`、Fault Module＝`wisp.test.exe` +0x95393），只作历史旁证

## 3. 摘录文件（本目录）

| 文件 | 内容 |
|---|---|
| `run1-rc.txt` | 第 1 发退出码/耗时/字节数原样 |
| `run1-stdout-tail40.txt` | 第 1 发最后 40 行 stdout |
| `run1-stderr-tail25.txt` | 第 1 发最后 25 行 stderr |
| `run1-sampler.csv` | 第 1 发 20 s 一档采样（CPU/WS/字节/webview2/竞争者/桌面） |
| `run1-gate.log` | 第 1 发独占闸日志（谁在占、何时放行） |
| `run1-desktop-before.log` | 第 1 发期间桌面状态（`Default` 未锁屏／`LogonUI=0`／前台 `magpie`） |
| `runA-fail-messages.tsv` | run A 121 枚假红的逐枚首句（用于"cwd/env 会把整包打红"的归类） |

## 4. 环境

- `go1.27.1 windows/amd64`；`GOTMPDIR` 空；`TMP`/`TEMP`＝`C:\Users\swq\AppData\Local\Temp`
- env：`WISP_SPIKE_MODELS`、`WISP_SILERO_MODEL`、`WISP_TTS_DIR`、`WISP_KWS_DIR`、`WISP_PUNC_MODEL`、`WISP_43_REPO_ROOT`（路径与 `scripts/wisp-cli-staged-tests.sh:312-331` 同）
- 驻留球 `wisp.exe` PID **49348** 全程在跑（⛔ 未动）；它占 `Ctrl+Alt+Q/M/P` ⇒ 第 2 发 8 枚热键族红是它的直接后果
- 并发：`gate-and-run.ps1` 等 `wisp.test.exe`/`wisp-winlive.test.exe` 连续 60 s 不存在才发车；第 1 发 12:51:29 有一档竞争者=1（另一腿的单枚定向跑，<20 s），其余全 0；第 2 发全程 0
