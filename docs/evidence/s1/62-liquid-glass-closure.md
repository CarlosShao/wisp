# 62 · wave-3 玻璃牙收件 — 状态与音频必须抵达可见像素（2026-11-07）

**腿**：62-liquid（写面 `internal/ball/**` 本波独占）　**执行**：写腿（GLM 5.3 Flash 路由组）
**本件是什么**：62-v1 对抗验收（`62-adversarial-acceptance.md`）点名的两条「判据今天不可满足／缺用例」里，
**今天仍成立且在本票写面内可闭的那两条**的收件记录。AC 勾与 Status 零触碰（纪律），本件是证据不是裁决表。

---

## 1. 起手三核（2026-11-07 21:2x 现跑）

| 核 | 命令/读数 | 结论 |
|---|---|---|
| ① 票面锚还在不在 | `liquidMotion`/`dockRampFrame`/`SetAudioLevel`/`TestBallLiveAudioLiquidGate` 全部 grep 命中在位；票据 Progress log 62-C/D/E 各段逐字在文件里 | 票面锚**未过期** |
| ② 引文出处 | `62-adversarial-acceptance.md` 存在（55,777B）；AC#3 判语「没有任何一条用例会因两个态长得一样而红」、AC#4 判语「单调律两道门真删过、都咬住，但 **wav→C8 注入口那条路全仓没有生产者**」逐字核过 | 引文属实 |
| ③ 红/绿是谁的 | 第一次全组秒挂的元凶 = **foreign 球窗**（22:0x `balldebug-orb.exe` pid 40144／更早 pid 58792，标题 "Wisp"，非本腿进程）＝**别人的在飞 live 探针占桌**；`git status --porcelain` 被追踪面净 | 秒挂非本腿产物，净桌后全绿 |

## 2. 今天仍成立的判据缺口（现核，非照抄票面）

- **AC#3 侧**：「两个态长得一样」至今**没有任何用例会红**。65 波把渲染搬到 WebView2 页面后，
  状态→页面唯一通路是 `SetState → pushState → orbStateScript`（`ball_windows.go:330`、`orb_windows.go:423`），
  该通路只有**字符串形状**钉子（`TestOrbStateScript`），**没有像素判据**。
- **AC#4 侧**：音频→页面通路 `applyLevel → pushAudio → orbAudioScript`（`liquid_windows.go:60`、`orb_windows.go:440`）
  同样只有脚本形状钉子；62-v1 咬过的单调律牙全在 **D2D 面**（`liquid.go` 模型），**orb 页像素零判据**。
- ⚠ AC#4 的「wav→C8 生产者」缺口**今天仍未闭**（生产者属 `internal/audio` 面，不在本腿写面）；本牙
  用 `Ball.SetAudioLevel`（项目唯一渲染侧音频入口，C25 合规）从**渲染面**钉像素单调律，**不冒充 C8 闭环**。

## 3. 新牙：`TestBallLiveOrbStateAndAudioReachPixels`（`orb_glass62_live_windows_test.go`，winlive）

真 WebView2 球（默认 StyleGlass＝siri 页）挂在已知黑板底上，**量合成屏像素**（`GetPixel` 屏 DC，
961 采样 px，3px 栅格，4px 内缩），五读数固定序：Sleeping(idle) → Thinking → Listening(静音) →
喂 0.5（喂中采样）→ 喂 0.9（喂中采样）。

### 读数表（两次独立全绿跑 + 探针，同一台机 2026-11-07）

| 读数 | 探针(21:41) | 绿跑#1(21:52) | 绿跑#2 收件(22:0x) |
|---|---|---|---|
| idle (Sleeping) | 0.1873 | 0.1863 | 0.1904→0.1883（M2 面） |
| thinking (Thinking) | 0.4204 | 0.4277 | 0.4246→0.4225（M2 面） |
| **T1 态差** | **0.2331** | **0.2414** | 0.0052（M1 面塌掉） |
| silence (Listening) | 0.4225 | 0.4256 | 0.4277 |
| level 0.5 | 0.4329 | 0.4311 | 0.4308 |
| level 0.9 | 0.4339 | 0.4357 | 0.4370 |
| **T2 静音→0.9 抬升** | **+0.0156** | **+0.0101** | −0.0017（M2 面塌掉） |
| T3 排序 silence<0.5<0.9 | ✓ | ✓ | ✓（M2 面 0.5=0.9=0.4221 平） |

### 阈值（冻结在测试注释里）

- **T1**：`|thinking − idle| ≥ 0.10`（实测 0.2331/0.2414，≥2 倍余量）＋ 双侧 ≥0.02 活性地板（黑页 0=0 不算态差）。
- **T2**：`level0.9 − silence ≥ 0.005`（实测 +0.0156/+0.0101，≥2 倍余量）＋ silence ≥0.02 活性地板。
- **T3**：`silence < 0.5` 且 `0.5 < 0.9`（两轮排序必立，纯序断言）。
- **如实记录**：0.5→0.9 档步实测 **+0.0010**，在单读噪声内 ⇒ **不钉**（页面对高电平饱和；排序牙已把方向钉死）。
- **采样方法纪律**：音频读数必须**喂中采样**（`sampleWhile` 30Hz 持续喂 + 3 帧 400ms 均值）——
  收件前绿跑#0 曾在**停喂 1.5s 后**采样，读到的是页面「zero them on stop」衰减尾（+0.003），
  不是驱动态；方法差会把真信号量丢。

## 4. 红牙（突变面，62-v1 `go test -overlay` 先例，工作树零改动；台件 `.scratch/wisp/probes/62/w3-mut/`）

| 牙 | 突变 | 红读数 |
|---|---|---|
| **M1** | `ball_windows.go` SetState 内 `o.pushState(s)` 摘除（映射名不再随态走） | T1 塌到 **0.0052 < 0.10** → FAIL「the pushed state does not reach pixels」 |
| **M2** | `liquid_windows.go` applyLevel 内 `o.pushAudio(level)` 摘除（包络不再进页） | T2 塌到 **−0.0017 < 0.005** → FAIL「the pushed audio bands do not reach pixels」；T3 平序 |

M1 同时牵连既有 `TestBallLiveOrbAttachStateAndTeardown`（debug face 停在 idle）＝第二锚。

## 5. 门禁名册（2026-11-07 22:0x）

| 门 | 读数 |
|---|---|
| `gofumpt -l internal/ball` | 空（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1` | ok 0.524s |
| `go test -tags winlive -run 'TestBallLiveOrb|TestBallLiveEdgeDock' -count=1 -timeout 25m` | **11/11 PASS**（ok 173.1s；含新牙 115.9s） |
| `d22scan` | clean（emoji 覆盖 internal/ 959 Go 文件含本件） |

测试环境：`unset GOTMPDIR TMPDIR TEMP TMP`＋`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，CGO 默认（⛔ 未置 0）。

## 6. 改动文件（全部 LF；两枚资产 sha 钉零触碰）

- `internal/ball/orb_glass62_live_windows_test.go`（新，winlive 测试件）
- `docs/evidence/s1/62-liquid-glass-closure.md`（本件）
- `.scratch/wisp/issues/62-liquid-glass-ball-visuals.md`（Progress log 追加一行，append-only）

⛔ 未触碰：`cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`、`internal/config/**`、`frontend/**`、`design/**`、
`assets/orb-siri.html`（sha `d93da649…`）、`assets/jarvis-orb.html`（sha `66eb4b80…`）、`thresholds/golden`。

## 7. 没做到的（如实）

- **AC#4 的 C8 生产者**（`internal/audio` 真采集帧 → `SetAudioLevel`）不在本腿写面，未闭；本牙从渲染面钉像素律。
- **AC#1 深色桌布半、AC#5 弹出率框面顶反、AC#6 SLO 回填**：仍开放，非本腿判据（62-v1 判语原文）。
- **SPEC-08 回填**：签字流程未动（票面契约变更条款），本腿零触碰。
- winlive 门禁期间两轮被**别人在飞探针占桌**秒挂（pid 58792/40144），按三核③如实记，未杀任何别家进程，净桌后复跑全绿。
