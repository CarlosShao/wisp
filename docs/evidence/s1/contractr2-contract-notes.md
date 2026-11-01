# 契约文字执行腿 `contract-r2` —— 两项已批契约文字落账（写腿，非裁决者）

日期 2026-11-01 · 腿 `contract-r2`（写腿/契约执行者；写面＝本件＋恰两枚 .md）·
机主批复（A1362 逐字）：「**D43 批准、DEFERRED 批准更新**」。
派单存档 `.scratch/wisp/probes/orch/20261101/txt/contract-r2-dispatch.txt` 在案。

## 0. 开工三核（本腿现跑，2026-11-01 16:4x +08）

| # | 项 | 读数 |
|---|---|---|
| ① | HEAD 在 dev ≥ `608693adc` | `dev` @ `608693adca2bbf2ef9ea184bbc493445ba6b360a`（＝派单下限本身，编排者 A1362 落账提交） |
| ② | 索引干净 | `git diff --cached --stat` 空 |
| ③ | docs/specs 无他人脏 | `git status --porcelain -- docs/PLAN.md docs/specs/` 空（工作树其余脏面＝frontend/internal 在飞腿，零接触） |

## 1. 坐标现核（两处均无漂移）

- **PLAN.md :3244-3247**：现读与派单坐标逐字一致——§16.11#10 具体待办段，原文倾向「作废并判拒绝」、
  「本轮未定案，须在 S7 的切片卡里明确」。43-v1 终裁（`43-v1-verdict.md` §2 裁 7）引的同一坐标同段原文。
- **SPEC-12-roadmap-governance.md :73-75**：三行 DEFERRED（26-PathC barge-in／28-PathC 全双工／60-C32 gated）
  与 contract-r1 2026-10-31 落的格式逐行对上（A1307 P15 NO-GO＋A1311③ 批推迟）。
- **依据具名核实**：`59-v3-verdict.md`（§0 P15＝PASS 无条件、「Path C 可以开工（票 26 重开进实现）」）、
  `26p-v2-verdict.md`（§0 终判＝有条件 PASS、(a) 语义重 key、摆动残余具名移交 28-PathC/393 轨）、
  `26pr1-impl.md`、`389r1-bargein-gate.md`、`26pk1-impl.md`、`392r1-session-integration.md`、
  `392r2-keepalive-fix.md`、`397r2-monitor-revival.md` 全部在盘。

## 2. 第一处：PLAN.md D43 段 note（机主「D43 批准」）

改法＝**纯追加**：原文四行一字不动，note 追加在原文段之后、`---` 分隔线之前；措辞照
`43-v1-verdict.md` §4 草案；格式照该待办段现有样式（粗体引导行＋硬换行续行）。

**Before（:3244-3247 逐字）**：

```
**关于第 10 条的具体待办**（记下来，因为它是一个已识别但未定义的行为）：
`AwaitingApproval` 态遇到系统挂起 → 审批请求应当**保留还是作废**？
倾向**作废并判拒绝**（挂起可能长达数小时，恢复后用户早已忘记自己在批准什么，
"恢复一个陈旧的批准请求"是安全事故），但**本轮未定案**，须在 S7 的切片卡里明确。
```

**After（原文保留＋追加 8 行；追加部分逐字）**：

```
**★ 定案（2026-11-01 机主批准，台账 A1362：「D43 批准、DEFERRED 批准更新」）**：
`AwaitingApproval` 遇系统挂起 ⇒ **一律作废**全部待确认项——L2 走 D43 #24 拒绝漏斗
（队头被拒绝 → `Acting`、取消该调用）、L1 倒计时窗口作废并留审计行；**按拒绝处理**，
并以 console 计数告知「错过 N 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）」；
恢复后**不自动重放**（重跑/重放只能由操作者重新发起，与崩溃恢复同一姿态）。
上面"倾向作废并判拒绝"的原判断即最终方向，43-r1 实现给出其 lean 可执行化
（判语与草案措辞见 `docs/evidence/s1/43-v1-verdict.md` §2 裁 7、§4）。
```

`git diff --stat`：`docs/PLAN.md | 8 ++++++++`（+8/-0，纯追加）。

## 3. 第二处：SPEC-12 §5 三行状态更新（机主「DEFERRED 批准更新」）

改法＝**改状态列（第一列）、依据列尾追加、其余四列一字不动、行不删**；状态列样式照本表
既有先例（:81 `~~DEFERRED~~ → **已采纳（…用户批准…）**`）；原 DEFERRED 历史（NO-GO＋A1311③）
全部保留在依据列。

**Before（:73-75 逐字，三行）**：

```
| DEFERRED | 票 26（TTS 播报）的 Path-C barge-in 部分——播报中说话停播转听 | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
| DEFERRED | 票 28（SessionScope C31）的 Path C 全双工 | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
| DEFERRED | C32 云端 realtime 引擎门（票 60，gated） | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
```

**After（三行，逐字）**：

```
| ~~DEFERRED~~ → **重开·已落地（2026-11-01 机主批 A1362「DEFERRED 批准更新」）** | 票 26（TTS 播报）的 Path-C barge-in 部分——播报中说话停播转听 | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）**；**重开（2026-11-01）依据＝59-v3 终判 P15 无条件 PASS（`docs/evidence/s1/59-v3-verdict.md` §0「Path C 可以开工」）＋26p-v2 终判有条件 PASS（`docs/evidence/s1/26p-v2-verdict.md`）＋机主批 A1362；重开已落地＝26p-r1 集成接线（`docs/evidence/s1/26pr1-impl.md`）＋389 组合门（`docs/evidence/s1/389r1-bargein-gate.md`）＋26p-k1 重 key（`docs/evidence/s1/26pk1-impl.md`）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
| ~~DEFERRED~~ → **重开·集成进行中（2026-11-01 机主批 A1362「DEFERRED 批准更新」）** | 票 28（SessionScope C31）的 Path C 全双工 | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）**；**重开（2026-11-01）依据＝59-v3 终判 P15 无条件 PASS（`docs/evidence/s1/59-v3-verdict.md` §0）＋26p-v2 终判（摆动残余具名移交本轨，`docs/evidence/s1/26p-v2-verdict.md` §0）＋机主批 A1362；集成进行中＝392-r1 会话集成（`docs/evidence/s1/392r1-session-integration.md`）＋392-r2 保温修复（`docs/evidence/s1/392r2-keepalive-fix.md`）＋397-r2 复活在飞（`docs/evidence/s1/397r2-monitor-revival.md`）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
| ~~DEFERRED~~ → **重开·门条件计通过、按 S6 可并行（2026-11-01 机主批 A1362「DEFERRED 批准更新」）** | C32 云端 realtime 引擎门（票 60，gated） | **PRECHECK P15 NO-GO**（台账 A1307，回环 rig 实测无 AEC 必自触发；`docs/evidence/s1/59r1-aec-spike.md` §8）；**机主已批推迟（A1311③，2026-10-31：「批准推迟」）**；**重开（2026-11-01）依据＝59-v3 终判 P15 无条件 PASS（D47 C32 门条件「S4 的 AEC spike 通过」已计通过，`docs/evidence/s1/59-v3-verdict.md` §0）＋机主批 A1362；按 S6 可并行（「用户有 Key」半边仍归用户自备）** | AEC3 重测 PASS 或 OS-AEC 在场机器 ERLE≥30dB（aec3-r1 在飞） | aec3-r1 重测（在飞） | 全双工暂缓、半双工不受影响 |
```

`git diff --stat`：`docs/specs/SPEC-12-roadmap-governance.md | 6 +++---`（+3/-3，三行各改状态列＋依据列尾追加）。

## 4. 写面清单（恰三件）

1. `docs/PLAN.md` —— 仅 §16.11#10 待办段后追加定案 note（+8/-0）
2. `docs/specs/SPEC-12-roadmap-governance.md` —— 仅 §5 :73-75 三行（+3/-3）
3. `docs/evidence/s1/contractr2-contract-notes.md` —— 本件

⛔ 未触碰：PLAN 其它行／SPEC 其它段／任何 .go／thresholds/golden/testdata/frontend/design／
票面／台账（`pending-and-issues.md` 零写入，A1362 本就是编排者已落账项）／docs/SLO.md／其它 specs。

## 5. 具名观察（⛔ 顺手改、只具名上报）

**PLAN.md:1536**（§15 批量授权段的 DEFERRED 登记）行「**`AwaitingApproval` 遇系统挂起的行为定案**（§16.11 第 10 条）」
现仍写「本轮识别但未定案……S7 切片卡明确定义，并写进 D43 转移表」。43-v1-verdict §4 的建议落点是
**三处**（D43 note／§16.11 第 10 条／DEFERRED 行），本腿派单明示只批两处、「PLAN 其它行一字不动」，
故 :1536 未动。⇒ **PLAN.md 内部现在一处说「未定案」（:1536）、一处说「已定案」（:3249 起）**，
时序差真实存在；该行的更新（改写为已定案指向）建议编排者另派或并入下一波契约文字腿，非本腿写面。

## 6. 边界自陈

- 本腿零代码改动、零测试运行需求（纯 .md 文字）；d22scan/名册等仪器射程不含 .md，无需跑。
- commit 只带显式 pathspec 三件、只 commit 不 push、无 amend/reset。
- 临时件（commit message 文件）在 `D:\tmp\wisp\20261101\contractr2\`，只建不删。
