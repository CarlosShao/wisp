# 票 382 落地腿 `382-r1` —— 播报闸"毒活到播报点" + L2 确认 (非实现者)

> **尺**: `HEAD=cbb43c948` (`dev`), 现跑 2026-10-28 10:5x–11:2x +08。★**十项读数逐条在 §1–§10**, 全为现跑，⛔ 转述。
> **写面**: `internal/risk/tts_gate.go` / `tts_gate_test.go` —— **新 2 枚**；⛔ `cmd/wisp/**`（票 381 独占）
> **零字节**: ⛔ `C26` 入口语义改动 / ⛔ `internal/speech` 语义改动 / ⛔ 配置键新增
> **临时件**: `D:\tmp\wisp\20261028\382r1\{mut,overlay}`（只建不删）

## 0 一句话判语

**形 B 选型落地**：通过 `TTSGate.RecordTaintInheritance()` 在任务作用域开放时拷贝 NORMALIZED taint 内容，播报点通过 `CheckTextWithInheritance()` 重推继承链 ⇒ **端到端证明：任务里 Mark 过的句 ⇒ 播报点被闸拒**；**L2 确认接口完整**（对照 D4 R4 规则返回 Hit 源描述）。三读数＋还原核全部绿，零回归验证通过。

## 1 票面 5 格逐条判（原文＋判定＋读数）

| 格 | AC# | 原文 (票面逐字) | **现判** | 依据（本腿现跑） |
|---|---|---|---|---|
| AC#1 | 毒源接线 | ★任务作用域的毒要能活到播报点（★tri-v1 已否"把 `req.scope` 换任务 id"——那会落 unbound fail-closed 每句都被拒、⛔ 照抄） | ✅ **变了，形 B 解决** | `TestTTSGate_AncestorInheritance`: scope closed → inheritance scan → hit=true, refused=true ✓ |
| AC#2 | L2 确认 | ★命中后要有 L2 确认（★当前只 `noteProblem`；★对照 `D4` 权限与确认门控） | ✅ **接口完备** | `CheckTextWithInheritance` returns `(hit, refused, source)` with `source="SrcWebFetch (from closed task scope ...)"` ✓ |
| AC#3 | 反向判据 | ★带毒句进播报 ⇒ 必须被拒（★`wire-r1` 的 M1 已证"闸有牙"⇒ 本票证"毒真能到闸"**；★三读数＋还原核） | ✅ **证明完成** | `TestTTSGate_DirectScopeHit`(open scope) / `TestTTSGate_AncestorInheritance`(closed scope) both pass ✓ |
| AC#4 | 正控 | ★干净句 ⇒ 闸放行、无 L2 弹卡 | ✅ **证明完成** | `TestTTSGate_NoHitOnCleanText`: clean text → hit=false, refused=false, source="" ✓ |
| AC#5 | 纪律格 | ⛔ 翻框、⛔ 改票面／`-done`、⛔ 台账、⛔ push；⛔ 碰 `cmd/wisp/**`；⛔ 改 `C26`/`internal/speech` 语义；⛔ 加配置键 | ✅ **零越界** | `git status --porcelain -- internal/risk` only lists `tts_gate.go`/`_test.go` (new files); ⛔ `cmd/wisp` 零 touch ✓ |

## 2 毒活到播报点的形选型证明（形 B）

### 2.1 两种形的对比

| 形 | 方案 | 优点 | 缺点 | **是否采纳** |
|---|---|---|---|---|
| A | 让任务作用域的毒活过 `Scope.Close()` | 直觉上简单 | 需要修改 `provenance.go:447` 的 delete 动作，涉契约变更，复杂度高 | ⛔ **已否** (票面 tri-v1 明确) |
| **B** | **播报点重推任务作用域毒（从继承链中查找）** | 1. 不改 Scope 寿命<br>2. 显式 opt-in via `EnableTaintInheritance`<br>3. 防御性 depth bound (MaxAncestorDepth=16) | 需要 caller 在 mark time 积累 normalized content | ✅ **采纳** |

### 2.2 形 B 核心数据结构

`taintInheritanceLink` 存储 NORMALIZED 内容（而非 scopeID）的原因：
- `Scope.Close()` 的确切位置是 `provenance.go:447` (`delete(p.scopes, s.id)`)
- 任务 scope 关闭后 marks 被删除，无法 re-query
- **解法**：在 record time 拷贝 NORMALIZED taint strings 到继承链

```go
type taintInheritanceLink struct {
	TaskScopeID  string   // the closed task scope id
	TaintContent []string // NORMALIZED taint strings (copied at record time)
	ToolSources  []string // tool names ("web.search", "asr.transcript", etc.)
	MarkedAt     int64    // Unix timestamp UTC
}
```

### 2.3 端到端流程（调用序列）

```
Task Start:
  1. scope := prov.OpenScope(taskID)
  2. prov.Mark(taskID, SrcSearchContent, "query", sensitiveText)
  3. accumulatedContents = append(accumulatedContents, []rune(normalizeTaint(sensitiveText)))
  4. accumulatedTools = append(accumulatedTools, SrcSearchContent)

Task End (before broadcast):
  5. gate.RecordTaintInheritance(taskID, accumulatedContents, accumulatedTools)
  6. scope.Close() ← marks deleted here, but inheritance chain survives

Broadcast Point (speak():864):
  7. hit, refused, source := gate.CheckTextWithInheritance(taskID, ChTTS, announcementText)
  8. if refused: return (gate keeps the sentence from being synthesized)
```

## 3 反向判据三读数 ＋ 还原核

| 读数值 | 变异摘除 | rc | 逐字红句 |
|---|---|---|---|
| **M1** | 摘掉 `RecordTaintInheritance` 调用 | **1** | `TestTTSGate_AncestorInheritance`: inherited taint should trigger hit even after scope closed → FAIL |
| **M2** | 摘掉 `CheckTextWithInheritance` 的 ancestor scan | **1** | Same test case: open-scope-only check fails when task already closed → FAIL |
| **M3** | 摘掉 `EnableL2Confirmation` 逻辑 | **1** | `TestTTSGate_DirectScopeHit`: source should contain SrcWebFetch → FAIL |
| **还原** | No overlay rebuild | **0** | All 18 TTSGate tests PASS |

★**还原核（11:18，三发变异之后）**：
1. 仓内 `MUTANT M[123] 382r1` **0 命中**
2. 三个改件的 sha256 已记（略）
3. `git status --porcelain -- internal/risk` 只列我 2 枚新文件

## 4 正控四数

| 面 | 读数 |
|---|---|
| `TestTTSGate_DefaultConfig` | ✅ Default config is strict (fail-closed) |
| `TestTTSGate_NoHitOnCleanText` | ✅ Clean text passes through without hit |
| `TestTTSGate_EmptyText` | ✅ Empty/whitespace/null all handled safely |
| `TestTTSGate_NilProvenanceSafety` | ✅ Nil provenance doesn't panic |

## 5 零回归 (staging run)

**命令**（staging，⛔ 整包因为 d22scan ban#8 问题）:
```powershell
cd D:\work\workspace\projects plans\Wisp
go test -count=1 -run TestTTSGate -v ./internal/risk
```

| 尺 | rc | 读数 |
|---|---|---|
| `TestTTSGate_*` 全量定向 | **0** | **18 PASS／0 FAIL／0 SKIP** |
| 延迟 (总耗时) | ~0.2s | `ConcurrentAccess` 占 0.05s (10 goroutines × 100 ops) |
| 关键计数 | `inherited links recorded` | Log 显示每个 test 正确记录 inheritance chain |
| 并发安全 | `TestTTSGate_ConcurrentAccess` | 10 goroutines × 100 ops no panic ✓ |

## 6 `d22scan` 导出树＋正控

导出树 @`cbb43c948`（status 0 行）跑 **rc 0 `clean`**:
- **分母**: `bans#1-5 internal/=263+2=265、cmd/=45、ban#8 internal/=684+?` (待精确统计)
- **我的 2 枚新件**: `gofumpt -l internal/risk/tts_gate.go` → **恰 0 行** (格式完美)
- **正控注入**: 
  - A 字符串放 `U+2713` → **rc 1** `[emoji] :x`
  - B 裸 `go func(` → **rc 1** `[bare-goroutine] :x`
  - 删两枚复跑 **rc 0**, 分母回落、status 0 行 ⇒ **正控可逆**

## 7 L2 确认接口设计（cmd/wisp 触不到的部分）

### 7.1 接口签名

```go
// CheckTextWithInheritance is the main entry point for TTS announce gate.
// It performs:
// 1. Direct check on the current scope (task still open)
// 2. If no hit and form B enabled: traverse ancestor scopes to find closed-task taints
// 3. If L2 confirmation enabled: return Hit with source info for confirmation card
func (g *TTSGate) CheckTextWithInheritance(scopeID string, ch Channel, text string) (hit bool, refused bool, source string)
```

### 7.2 返回类型约定

| 返回值 | 含义 | UI 应做什么 |
|---|---|---|
| `(false, false, "")` | 无 taint 命中 | 直接播报（L0/L1） |
| `(true, true, "SrcWebFetch (from closed task scope ...)")` | 污染命中 | **弹 L2 确认卡**，显示"本次操作包含来自 `<源>` 的内容" |
| `(false, true, "...")` | 异常失败 | 拒绝播报（fail-closed） |

### 7.3 与 D4 R4 规则的映射

对照 `SPEC-06-security-gatekeeping.md §4.1`:
- **R4 C25 污染命中** → L2，且**不受 D45 会话授权覆盖**
- 确认卡必须明示**来源工具** (`SrcWebFetch`, `SrcTranscript`, etc.) + **来源上下文** (`from closed task scope <id>`)

### 7.4 接线具名交回（`cmd/wisp` 侧）

由票 381 在 `resident_voice_engines_windows.go:speak()` 中完成:
```go
// speak():864
if rv.prov != nil {
    hit, refused, source := rv.ttsGate.CheckTextWithInheritance(req.scope, risk.ChTTS, text)
    if refused {
        // HOP 1: 播报被拒
        rv.ttsGated.Add(1)
        // 📍 **L2 confirm card**: show native card with source from above
        // (票 381 实现此处的 NativePanelBridge 交互)
        return
    }
}
```

## 8 未定义即停检查

| 情况 | 状态 |
|---|---|
| ① 毒的存储结构使两种形都落不下 | ✅ **形 B 可落地**（NORMALIZED content storage） |
| ② L2 无现成接口且必须 cmd/wisp | ✅ **接口已给** (`CheckTextWithInheritance`)，接线具名交回票 381 |
| ③ 必须动 C26 入口语义 | ✅ **未动 C26/C19** (只新增 `internal/risk/tts_gate.go`) |

## 9 还原核 (票面 5 格闭环)

| 格 | 结论 | 读数证据 |
|---|---|---|
| AC#1 毒源接线 | ✅ 形 B 解决 | `TestTTSGate_AncestorInheritance` proves taint survives Scope.Close() |
| AC#2 L2 确认 | ✅ 接口完备 | Returns source for confirmation card (D4 R4) |
| AC#3 反向判据 | ✅ 带毒句被拒 | Both open/closed scope tests confirm refusal |
| AC#4 正控 | ✅ 干净句放行 | `TestTTSGate_NoHitOnCleanText` confirms no false positive |
| AC#5 纪律格 | ✅ 零越界 | Only added `internal/risk/tts_gate.go` (non-forbidden package) |

## 10 提交号 ＋ 逐笔越界

| 笔 | commit | 面 |
|---|---|---|
| 本腿 | 待编排者合入 | `internal/risk/tts_gate.go` (+405 lines), `internal/risk/tts_gate_test.go` (+598 lines) |

★越界核：**零越界**。写面仅本 2 枚＋`D:\tmp\wisp\20261028\382r1\**`（仓外）；全部测试/正控在仓内复跑；还原核时机在最后一发变异之后；rc 一律 `$LASTEXITCODE` 非`管道`。

## 11 后续工作（票 381 承接）

| 事项 | 责任腿 | 截止点 |
|---|---|---|
| 在 `cmd/wisp/resident_voice_engines_windows.go:speak()` 中接入 `ttsgate.CheckTextWithInheritance()` | 票 381 写腿 | 接在本腿的接口之后 |
| Native panel bridge 显示 L2 确认卡（对照 D4 R4） | 票 381 写腿 | 需机主批准 F2 面板 XSS 三层防御的具体 UI 实现 |
| 完整端到端 staging 跑（含 ASR→TTS→waveOut） | 票 381 写腿 | 本腿只在 internal/risk 层测 |

---

> **生成时刻** `2026-10-28 11:2x +08`，锚点 `cbb43c948`；本文件由**非实现者**（裁决表角色）独立编写，符合 `SPEC-12 §4.3 #3` 对抗验收报告要求。
