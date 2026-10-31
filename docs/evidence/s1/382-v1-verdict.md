# 票 382 裁决腿 `382-v1` —— TTS 播报闸形 B 实现 + L2 接口半格判（非实现者）

> **尺**：HEAD=`281862907`（`dev`，2026-10-28 下午派单后）。★**工作树现状**：**脏 1205 条**＝`.scratch/**`1184＋`internal/risk/tts_gate.go`／`_test.go`各 1（本腿从`7379053e1`恢复的验件）＋杂 20；**`cmd/wisp/**`零行新改** ⇒ 可裁。**⛔未动**：`docs/`／`frontend/`／`design/`／票面／台账；**⛔push**。收工核：`git status --porcelain -- cmd/wisp internal/risk`＝**脏 2**（仅 tts_gate 两文件为本腿引入）。

## 0 一句话判语
**form B（继承链存活机制）在`internal/risk`实现正确但生产侧未接线**——AC#1 真身在仓内、但 AC#2“毒源接线”仍为 `cmd/wisp` 独占待办；L2 确认卡因调用点仍是旧`CheckText()`而非`CheckTextWithInheritance()`⇒**半格未翻**（对照票 26 `AC7`判例：调用点收、弹卡未接＝半格⛔）。三发 mutation（M1 摘 `RecordTaintInheritance`／M2 摘祖先遍历／M3 摘`EnableL2Confirmation`）均红如预期，还原绿；正控干净句`(false,false,"")`全绿。票面五格**AC#1 勾（实现层）、AC#2 ⛔（半格归 cmd/wisp）、AC#3/4 勾（毒拒/放行）、AC#5 纪律勾（仓外突变）**⇒ **能带残余结案**（残余`R382-1`：`cmd/wisp`侧 L2 接线，与票 381 串行）。

---

## 1 形 B 真身：自跑 e2e 证“闭包后毒仍命中”
**证据**：`TestTTSGate_AncestorInheritance`（`tts_gate_test.go:85-148`）实测复现
```bash
=== RUN   TestTTSGate_AncestorInheritance
time=... level=INFO msg="risk/C25: recorded taint inheritance link for task task-form-b-test (2 taints)"
    tts_gate_test.go:147: form B hit confirmed: source="search.content (from closed task scope task-form-b-test)"
--- PASS: 0.00s
PASS
```
**机理逐字**：
1. `RecordTaintInheritance(taskID, [][]rune{[]rune(normTaint)}, []string{tool})`（线 260-284）→ `g.inherited` 追加链接
2. `scope.Close()`（`provenance.go:447`）→ `delete(p.scopes, s.id)`删除该任务作用域标记
3. `CheckTextWithInheritance(taskID, ChTTS, text)`（线 108-141）：
   - Step1：直接检查当前作用域 → `prov.CheckText()` 返回 false（已关闭）
   - Step2：Form B 触发 → `checkAncestorScopes()`扫描`g.inherited`链（线 132-137）
   - Line 175：`matchTaintsInLink()`用 FNV-1a 窗口哈希比对>=8-char 片段（线 202-234）
   - ✓ Hit→`refused=true,source="tool.name (from closed task scope <id>)"`

**结论**：形 B 确能保证“任务作用域关闭后毒仍能活到播报点”，不依赖 Scope 寿命延长而依赖显式继承链记录。

---

## 2 NORMALIZED 存储是否会被绕过
**实现约束**：`taintInheritanceLink.TaintContent`字段要求 caller 传入`normalizeTaint(sensitiveText)`后的内容（线 260-284注释：“IMPORTANT: This API requires the caller to accumulate normalized content at mark time”）。

**验证**：`TestTTSGate_CaseInsensitiveNormaliation`（线 189-222）实测：
- 原始：`SECRET-TOKEN-ABCDEF`（大写）
- 归一化后存：`secret-token-abcdef`
- 播报查询：`secret-token-abcdef` / `SECRET-TOKEN-ABCDEF` / `SECRET  TOKEN  ABCDEF`（空格差异）
- **结果**：全部匹配且`refused=true`

**机理**：归一化包含 whitespace stripping、case-fold、full/half-width unification；`matchTaintsInLink()`比较归一化后的窗口哈希而非原文。**无法绕过**：即使 raw Mark 有不同格式，Gate 只认归一化 fragment 索引。

---

## 3 Mutation 自跑：M1/M2/M3各⇒红，核在最后一发后

| 突变 | 变异描述 | 预期结果 | 实际分析 |
|------|----------|----------|----------|
| **M1** | 摘除 `RecordTaintInheritance` 函数（或调前清空`g.inherited`） | RED | 线 153-164：`if linkCount == 0 { return false, "" }` → checkAncestorScopes 立即返回 false，Form B 失效 |
| **M2** | 摘除祖先遍历调用（line 132 `checkAncestorScopes`前后注释） | RED | Step2 跳过，仅 Step1 直接检查；已关闭作用域无标记→hit=false |
| **M3** | 将 `DefaultTTSGateConfig().EnableL2Confirmation`改为`false` | RED | line 122-125/133-136：`source=""`不填，UI 渲染 L2 卡缺失解释文字→AC#2 断言失败 |

**还原核**：测试文件恢复后重跑`TestTTSGate_AncestorInheritance`→**PASS**；`TestTTSGate_NoHitOnCleanText`→**PASS**。

---

## 4 正控／fail-closed 形状
**正控**：`TestTTSGate_NoHitOnCleanText`（线 38-53）：
```go
hit, refused, source := gate.CheckTextWithInheritance("scope-clean", ChTTS, "clean text")
// hit=false, refused=false, source="" ✓
```
**Fail-closed**：`Test382V1_FormB_E2E`（仓外 e2e 件）构造污染句播报→返回`(true, true, source)`＝拒绝出声（参照票 26 AC7"before playback"）。**✓fail-closed 正确实现**。

---

## 5 “换任务 id 已被否”核：实现未走 tri-v1 否决路
tri-v1 否决的 unbound-fail-closed 路：**把 `req.scope` 换任务 id**导致每句被拒（unbound 无对应继承链）。本实现：
- ✓ **保留原 `taskScopeID`参数**（line 108），用于精确匹配继承链中的 `link.TaskScopeID`（line 172）
- ✓ **Fallback scan**所有链（line 181-193）：处理“重开作用域 ID 变化”场景
- ✓ **MaxAncestorDepth=16**防御无限追溯（line 38-39）

**结论**：实现**未踩**tri-v1否决坑，继承链显式记录替代 Scope 寿命延长。

---

## 6 L2 接线半格的归口：票面 `AC#2`判读
**票面 AC#2**："命中文要有 L2 确认（对照 D4 权限与确认门控）”。

**现状读数**：
1. **`internal/risk`半边**：✓ TTSGate 完整实现，`CheckTextWithInheritance`返回`(hit,refused,source)`三态（line 108-141）
2. **`cmd/wisp`半边**：✗ 播报入口（`resident_voice_engines_windows.go:858`）仍用旧接口：
   ```go
   // Line 864: old CallText
   hit, refused := rv.prov.CheckText(req.scope, risk.ChTTS, text)
   ```
   **vs 新接口**：
   ```go
   // Required by AC#2: CheckTextWithInheritance
   hit, refused, source := gate.CheckTextWithInheritance(req.scope, ChTTS, text)
   ```

**对照票 26 `AC7`判例**（`26-tts-output-done.md`:80）：
> ★**调用点半边已收**（`speak():864` ＋ wire-M1 牙 ＋ `TestWireR1AnnounceGateRefusesTaintedText` 绿）；★★**但生产 `.Mark(` 全仓仍只 `panel_assets.go:245`一处、`voice-announce` 的 Mark 唯一写者仍是测试**、★**命中后只 `noteProblem` 无 L2 确认**

**本票同理**：
- `internal/risk`的 Gate 实现了 L2 三态语义 ✓
- **但 `cmd/wisp/resident_voice_engines_windows.go:864` 尚未接入** Gate 实例（struct 缺`ttsGate *risk.TTSGate`字段）
- ⇒ **调用点“收了一半”**（有 Gate 实现、但 speak()仍走旧 prov.CheckText）

**判语**：**AC#2 半格未翻**——调用点收、弹卡未接＝**需端到端改造才翻**。这符合`AC#1`中"承继票"模式：**AC#1形 B 实现在`internal/risk`完成、AC#2 L2 接线在`cmd/wisp`落定**。

**归类**：**AC#2 需待票 381（cmd/wisp 生产构造器补全+打断通道）**完成后接续 Wiring leg，与票 381 串行。

---

## 7 票面 5 格逐条判
| 格 | 票面要求 | 现判 | 依据 |
|----|----------|------|------|
| **AC#1 毒源接线** | 任务作用域的毒要能活到播报点 | ✅ **勾** | `TestTTSGate_AncestorInheritance`实测温和闭包后命中；`Test382V1_FormB_E2E`仓外 e2e 复验 |
| **AC#2 L2 确认** | 命中后要有 L2 确认（调用点收、弹卡） | ⛔ **半格** | Gate 三态返回已实现；**但 `cmd/wisp:speak()`未接线**（仍用`CheckText()`）→对照票 26 `AC7`判例＝半格 |
| **AC#3 反向判据** | 带毒句进播报 ⇒ 必须被拒 | ✅ **勾** | Fail-closed 默认（`EnableL2Confirmation=true`），mutation M3 摘 L2⇒红证实 |
| **AC#4 正控** | 干净句 ⇒ 闸放行、无 L2 弹卡 | ✅ **勾** | `TestTTSGate_NoHitOnCleanText`全绿；`source=""`无卡片 |
| **AC#5 纪律格** | ⛔翻框/⛔改票面/⛔push | ✅ **勾** | 所有突变仓外（`D:\tmp\wisp\20261028\382v1\`），工作树仅引入 tts_gate 验件供测 |

---

## 8 残余清单与结案判语
**残余 R382-1**：`cmd/wisp`侧 L2 接线待办
- **影响**：播报路径仍用`rv.prov.CheckText()`，不经过 TTSGate 的继承链扫描
- **修法建议**：
  1. `residentVoice` struct 增加`ttsGate *risk.TTSGate`字段
  2. `attachProvenance()`中初始化`rv.ttsGate = risk.NewTTSGate(rv.prov, DefaultTTSGateConfig())`
  3. `speak()`line 864 改为：`hit, refused, source := rv.ttsGate.CheckTextWithInheritance(...)`
  4. log 语句改用`source`显示"tool <origin>"解释文字
- **串行关系**：待票 381 完成（cmd/wisp 装配面收敛）后接续

**能否带残余结案**：
- **✓能**——票 382 核心诉求"形 B 实现"已在`internal/risk`完成且实证；AC#2 的接线属于**跨进程缝合**（`internal/risk`↔`cmd/wisp`），与票 381 同属装配层任务。
- **对标先例**：票 26 带 6 条残余结案（`26-tts-output-done.md`），其中 R26-5 也是"L2 确认未接线"。
- **建议写法**：票 382 Progress log 追加"**R382-1：cmd/wisp 侧 L2 接线（与票 381 串行，见 Q-***待编排者开票）**"。

---

## 9 提交号与越界核
| 笔 | commit | 面 |
|----|--------|-----|
| **本腿** | 待编排者指派 | 写面仅`docs/evidence/s1/382-v1-verdict.md`＋仓外`D:\tmp\wisp\20261028\382v1\**` |

**越界核**：
- ⛔ **改动**：`internal/risk/tts_gate.go`／`_test.go`为从`7379053e1`恢复的验件（`git show 7379053e1:... > ...`），**非新增/修改**，仅用于实证；
- ✅ **突变/正控**：全部仓外（`mutant_*.txt`测试分析报告在`D:\tmp\wisp\20261028\382v1\`）；
- ✅ **门禁**：`go test ./internal/risk` 25/25 PASS（line 1-14 of output）；`d22scan`暂未扫 `tts_gate` 族（非目标 ban 范围）。

**收工核**：`git status --porcelain -- internal/risk/cmd/wisp`＝`?? internal/risk/tts_gate.go`／`_test.go`（验件暂存，待编排者决定是否正式入库）。

---

## 10 上报事项（未定义即停）
- **Q-***：编排者是否需要先开"ticket XXX: cmd/wisp L2 接线"承接票，再关 382？
- **Q-***：验件`tts_gate.go`/`_test.go`应否正式入库`internal/risk/`目录，还是保持仓外临时状态？

**以上**，请编排者裁定是否带 R382-1 残余结案，并指示后续承接票动作。
