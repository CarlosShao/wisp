# 390 — 球单击打断装配（381 残余①：381-r1 自认项、381-v3 判"建议独立票"）

**Status:** **DONE-带残余**（2026-10-31 编排者收口：五格全翻（390-v1 独立突变件复刻）、ball API 无需扩展推翻坐实；残余三条见末节）
- **写面**：`cmd/wisp/**`（球单击事件 → `rv.TriggerInterrupt()` 的挂接；★**Esc/veto 通路已落、只缺球单击这一跳**）
- **机理**：★381-r1 的 Esc→`vetoByEsc()`→`TriggerInterrupt` 通路已在（381-r2 增强＝speaker.Abort 真停 ≤400ms 实测 MET）；★**球单击＝第二个触发源**（381-r1 申报"球单击→预留钩子但未装配执行器（需 internal/ball API 扩展，不在本片范围）"）

## 判据
- [x] **AC#1 球单击→TriggerInterrupt 真挂**（★**对照 Esc 通路的挂接形状**；★**internal/ball API 若需扩展 ⇒ 具名交回**——ball 是冻结面）
- [x] **AC#2 端到端**：★**播报中球单击 ⇒ TTS 真停 ≤400ms ＋ 门重开**（★**复用 381-r2 的 staged 测量形状**）
- [x] **AC#3 反向判据**：★**摘球单击挂接 ⇒ 球单击不再打断 ⇒ 红**（三读数＋还原核）
- [x] **AC#4 零回归**（★**六族＋Esc 通路不动**）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `internal/{ball,risk,session,speech}/**`（★**ball API 扩展须具名交回**）
## ★★★终裁落地（2026-10-31 19:5x 编排者，凭非实现者 `390-v1`，件 `docs/evidence/s1/390-v1-verdict.md` 83 行、`b6548ed56`）

★**五格全翻 ⇒ 可带残余结案**。
★**"ball API 无需扩展"推翻坐实**：★**`ball_windows.go:50`＝`OnClickBall func()`、`:632`＝`b.fire(...)`（WM_LBUTTONUP 非拖拽非 dock）逐字命中**；pre-390 旧闭包本就挂此面、e12200317 只换闭包体、**`internal/ball` 零字节（最后改动 10-04 票 260）**——381-r1 申报（`381r1-impl.md:49`）不成立＝"API 已够"非"停手交回"。
★**复跑读数**：★e2e **33/38/28/37/29ms**（合跑发再 5 轮 35-50ms）＋控制轮照播完 2.108s＋门 6/6；★**M1 独立造突变件**（SHA BD37DD5F≠原件≠实现腿件）主案红＋**同突变下 Test381r2 两案＋Unattached 案 3/3 绿＝Esc 照常自证定向**；★正控 idle 惰性＋UNATTACHED 反形；★**纯净对照树 f62b0fc3b 同枚同句同点（:226 3s Armed cap）＝Test384 合跑红 HEAD 既有、红名集合零新增**；★线程形状＝两触发源同枚 Events 字面量同 fire() ui-sta 内联（381-v2 要求同形满足）。
★**两处口径差具名（不动摇实质）**：A1326 卡"四 DLL"实为三枚即足；`WISP_PUNC_MODEL` 未设（无案需要）。★**踩坑供复跑腿**：pwsh 直调测试二进制拆 `-test.v` 旗标（381-v3 同款）⇒ **用 Git Bash、runner 在 `D:\tmp\wisp\20261031\390v1\run.sh`**。
★**残余三条**：①**单击 summon 半仍无执行者**（ballGestureWhy 已随真；★**归口＝GUI 手势面挂账、按 390-v1 建议记 228 AC#2 家族、未立票＝待需要时开**）②Test384 3s 时窗＝HEAD 既有（归 384 侧）③Drain 早回＝381-v3 §10③ 既有、归 26 侧。
