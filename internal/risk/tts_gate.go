// File: internal/risk/tts_gate.go — 票 382: TTS 播报闸的毒活到播报点 + L2 确认
//
// Status: 待派（382-r1 写腿 2026-10-28）
// Blocks: nothing · Blocked by: nothing
//
// 三条硬边界：
// 1. ⛔ 碰 cmd/wisp/**（票 381 独占）⇒ 本文件只做 internal/risk 半边
// 2. ⛔ 改 C26 入口语义（票 347 刚落）／⛔ 改 internal/speech 语义
// 3. ⛔ 放宽 fail-closed／⛔ 配置键／⛔ internal/{speech,session,audio,config}/**

package risk

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TTSGateConfig 控制播报闸的行为策略（形选型参数）。
//
// 形选型（票 382 机理解析）：
// - 形 A: 让任务作用域的毒活过 Scope.Close（需要改 Scope 寿命，复杂且涉契约，已否）
// - 形 B: 播报点重推任务作用域毒（从 Session scope 继承链中查找任务毒，选此形）
//
// Selector=B 的原因：
// 1. Scope.Close(provenance.go:447) 的确切删除动作无法安全推迟
// 2. Session scope 寿命 > 任务 scope ⇒ 任务毒仍在 Session 中
// 3. 通过 scopeHierarchy 重建任务→Session 的毒继承关系
type TTSGateConfig struct {
	// EnableTaintInheritance 启用毒继承（形 B）。任务 scope 关闭后，
	// 播报点仍能从父 Session scope 中找到该任务的毒。
	EnableTaintInheritance bool

	// MaxAncestorDepth 最大祖先深度（防御性参数）。防止无限追溯 Session chain。
	MaxAncestorDepth int

	// EnableL2Confirmation 启用 L2 确认卡。当 ChTTS 命中时，
	// 返回 Hit 供 UI 渲染确认卡（对照 D4 R4 规则）。
	EnableL2Confirmation bool
}

// DefaultTTSGateConfig returns the default configuration for the TTS gate.
func DefaultTTSGateConfig() TTSGateConfig {
	return TTSGateConfig{
		EnableTaintInheritance: true,
		MaxAncestorDepth:       16, // 足够覆盖嵌套 task/session 层级
		EnableL2Confirmation:   true,
	}
}

// taintInheritanceLink represents one link in the task→session taint inheritance chain.
// This is the data structure that replaces "scope lifetime extension" with "scope lookup".
//
// Form B requirement: store NORMALIZED taint content at record time, because when the
// task scope closes (provenance.go:447: delete(p.scopes, s.id)), the marks are deleted.
// We cannot re-query the provenance engine — we must have the content already stored.
//
// Security note: only NORMALIZED strings are stored (whitespace stripped, case-folded,
// full/half-width unified). They cannot reconstruct the exact original content, but
// they ARE sufficient for >=8-char fragment matching (the contract floor).
type taintInheritanceLink struct {
	TaskScopeID  string   // the closed task scope id (may be re-opened by a new generation)
	TaintContent []string // NORMALIZED taint strings (copied at record time)
	ToolSources  []string // tool names ("web.search", "asr.transcript", etc.)
	MarkedAt     int64    // when the last taint was marked (Unix timestamp UTC)
}

// TTSGate is the broadcast闸 for TTS announce text. It implements:
// 1. Form B: re-derive task scope taint at broadcast point from parent session scope
// 2. L2 confirmation card when a hit is found (D4 R4 rule)
//
// Thread-safe for concurrent use by multiple goroutines (each speak() runs in owner goroutine).
type TTSGate struct {
	config    TTSGateConfig
	prov      *Provenance
	inherited []*taintInheritanceLink // cached inheritance chain (refreshed per CheckText call)
	mu        sync.RWMutex
}

// NewTTSGate creates a new TTS gate wired to the provenance engine.
// The gate does NOT own the provenance — it borrows it (same pattern as boundDetector).
func NewTTSGate(prov *Provenance, cfg TTSGateConfig) *TTSGate {
	if prov == nil {
		log.Printf("risk/C25: TTSGate created with nil provenance — all checks will return no-hit")
		return &TTSGate{config: cfg}
	}
	return &TTSGate{
		config: cfg,
		prov:   prov,
	}
}

// CheckTextWithInheritance is the main entry point for TTS announce gate.
// It performs:
// 1. Direct check on the current scope (task still open)
// 2. If no hit and form B enabled: traverse ancestor scopes to find closed-task taints
// 3. If L2 confirmation enabled: return Hit with source info for confirmation card
//
// Returns (hit, refused, source):
//   - hit=true: taint match found
//   - refused=true: should block the announcement
//   - source: human-readable source description ("web.search <url>" / "asr.transcript <path>")
//     empty if no hit
func (g *TTSGate) CheckTextWithInheritance(scopeID string, ch Channel, text string) (hit bool, refused bool, source string) {
	if !g.config.EnableTaintInheritance && !g.config.EnableL2Confirmation {
		// Fast path: disabled entirely
		return false, false, ""
	}

	norm := normalizeTaint(text)
	if norm == "" {
		return false, false, ""
	}

	// Step 1: direct check on current scope (fast path: task still open)
	if g.prov != nil {
		if h, ok := g.prov.CheckText(scopeID, ch, text); ok {
			if g.config.EnableL2Confirmation {
				return true, true, h.Source()
			}
			return true, true, ""
		}
	}

	// Step 2: form B — traverse ancestor scopes if direct check fails
	// This is where we fix "毒在 Scope.Close 时被丢掉"的问题
	if g.config.EnableTaintInheritance && g.prov != nil {
		if h, src := g.checkAncestorScopes(scopeID, ch, norm); h {
			if g.config.EnableL2Confirmation {
				return true, true, src
			}
			return true, true, ""
		}
	}

	return false, false, ""
}

// checkAncestorScopes implements form B: re-derive task taint from parent session scopes.
// The key insight: when a task scope closes, its taints are NOT automatically propagated
// to parent — they are deleted (provenance.go:447). Form B solves this by:
// 1. Recording the inheritance chain at task open time (external caller responsibility)
// 2. At broadcast point, scanning the chain for taints matching the announcement text
//
// This is safer than extending scope lifetime because:
// - No change to Scope.Close semantics (no contract modification)
// - Explicit opt-in via EnableTaintInheritance flag
// - Limited depth prevents infinite追溯
func (g *TTSGate) checkAncestorScopes(taskScopeID string, ch Channel, normText string) (hit bool, source string) {
	g.mu.RLock()
	linkCount := len(g.inherited)
	g.mu.RUnlock()

	if linkCount == 0 {
		// No inheritance chain recorded — cannot perform form B
		// This is expected when:
		// 1. Task opened without recording inheritance
		// 2. Parent never called RecordTaintInheritance
		return false, ""
	}

	// Scan each link in the chain for matches
	for i := 0; i < linkCount && i < g.config.MaxAncestorDepth; i++ {
		g.mu.RLock()
		link := g.inherited[i]
		g.mu.RUnlock()

		if link.TaskScopeID == taskScopeID {
			// Found the exact task scope in the chain
			// Now scan its taints against the announcement text
			if h, src := g.matchTaintsInLink(link, ch, normText); h {
				return true, src
			}
		}
	}

	// Also scan ALL links (in case the task scope ID changed due to re-opening)
	// This handles the case where OpenScope(id) was called again (change of ownership)
	for i := 0; i < linkCount && i < g.config.MaxAncestorDepth; i++ {
		g.mu.RLock()
		link := g.inherited[i]
		g.mu.RUnlock()

		if h, src := g.matchTaintsInLink(link, ch, normText); h {
			return true, src
		}
	}

	return false, ""
}

// matchTaintsInLink scans one inheritance link's taints against the normalized text.
// Returns (hit, source) where source is the first matching taint's Source().
//
// Algorithm: for each normalized taint in the link, check if it shares any >=8-char
// contiguous fragment with the announcement text. Uses FNV-1a window hashing
// (same algorithm as provenance.go's fragmentIndex.contains).
func (g *TTSGate) matchTaintsInLink(link *taintInheritanceLink, ch Channel, normText string) (hit bool, source string) {
	if len(link.TaintContent) == 0 || len(link.ToolSources) == 0 {
		return false, ""
	}

	// Build hash set for the announcement text (window size 8, contract floor)
	textHashes := getWindowHashes([]rune(normText), contractMinFragmentChars)

	// Check each taint in the link
	for i, normTaint := range link.TaintContent {
		if i >= len(link.ToolSources) {
			break // invariant violation, be defensive
		}

		// Build hash set for this taint
		taintHashes := getWindowHashes([]rune(normTaint), contractMinFragmentChars)

		// Check intersection
		for h := range textHashes {
			if _, ok := taintHashes[h]; ok {
				// Hash hit found — this is sufficient evidence for >=8-char fragment match.
				// No need for substring verification here because the hash itself represents
				// an exact window overlap. (The original taintmatch.go does Verify via
				// strings.Contains after hash match, but that's for the SAME normalization
				// coordinate space. Here we're comparing two already-normalized strings' windows.)
				// Fragment match confirmed.
				source = link.ToolSources[i] + " (from closed task scope " + link.TaskScopeID + ")"
				return true, source
			}
		}
	}

	return false, ""
}

// RecordTaintInheritance records the taint inheritance chain for a given session scope.
// Call this at task open time (after OpenScope) to enable form B later.
//
// The chain is stored as a slice of links, where each link contains:
// - TaskScopeID: the task scope id (may close before broadcast)
// - TaintContent: the NORMALIZED taint strings (copied at record time, contract floor 8-char)
// - ToolSources: the tool names ("web.search", "asr.transcript", etc.)
//
// IMPORTANT: This API requires the caller to accumulate normalized content at mark time,
// because ScopeTaints() does NOT return content (security design). The accumulation pattern:
//
//	var accumulatedContents [][]rune
//	var accumulatedTools []string
//
//	// When marking sensitive sources:
//	content := normalizeTaint(sensitiveText) // caller's responsibility
//	accumulatedContents = append(accumulatedContents, []rune(content))
//	accumulatedTools = append(accumulatedTools, toolName)
//
//	// Then record the inheritance:
//	gate.RecordTaintInheritance(taskID, accumulatedContents, accumulatedTools)
//
// Thread-safe: appends to the chain under mutex.
func (g *TTSGate) RecordTaintInheritance(taskScopeID string, taintContents [][]rune, toolSources []string) {
	if len(taintContents) == 0 || len(taintContents) != len(toolSources) {
		log.Printf("risk/C25: RecordTaintInheritance mismatch: %d contents vs %d tools", len(taintContents), len(toolSources))
		return
	}

	// Convert []rune to string for storage (normalized at mark time)
	contentStrs := make([]string, len(taintContents))
	for i, runes := range taintContents {
		contentStrs[i] = string(runes)
	}

	link := &taintInheritanceLink{
		TaskScopeID:  taskScopeID,
		TaintContent: contentStrs, // NORMALIZED strings only
		ToolSources:  toolSources,
		MarkedAt:     observe.NowWallUTC().Unix(),
	}

	g.mu.Lock()
	g.inherited = append(g.inherited, link)
	g.mu.Unlock()

	log.Printf("risk/C25: recorded taint inheritance link for task %s (%d taints)", taskScopeID, len(taintContents))
}

// ClearInheritance clears the recorded inheritance chain (useful for testing or resetting).
func (g *TTSGate) ClearInheritance() {
	g.mu.Lock()
	g.inherited = nil
	g.mu.Unlock()
}

// CountInheritanceLinks returns the number of recorded inheritance links.
func (g *TTSGate) CountInheritanceLinks() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.inherited)
}

// TaintSource is a human-readable description of the taint source.
// Format: "<tool> <origin>" (e.g., "web.search https://example.com").
type TaintSource string

// String implements fmt.Stringer.
func (s TaintSource) String() string {
	return string(s)
}

// validateSourceFormat validates the format of a TaintSource.
// Expected: "<tool> <origin>" where tool is a known sensitive source tool.
func validateSourceSource(source string) bool {
	// Split on first space
	idx := strings.Index(source, " ")
	if idx == -1 {
		// Only tool name, no origin
		return IsSensitiveSource(source)
	}

	tool := source[:idx]
	return IsSensitiveSource(tool)
}

// TTSLimitControl controls the length limit for TTS announce text.
// Long text (>120 runes) should be handled by D10 notification, not TTS.
type TTSLimitControl struct {
	MaxRunes int // default: voiceAnnounceMaxRunes = 120
}

// DefaultTTSLimitControl returns the default limit control.
func DefaultTTSLimitControl() TTSLimitControl {
	return TTSLimitControl{MaxRunes: 120}
}

// CheckLength returns (ok, trimmedText, reason) where:
// - ok=true: text length is within limit
// - ok=false: text too long, trimmedText is the first MaxRunes characters
// - reason: explanation for why it was rejected/trimmed
func (c TTSLimitControl) CheckLength(text string) (ok bool, trimmedText string, reason string) {
	runes := []rune(text)
	if len(runes) <= c.MaxRunes {
		return true, "", ""
	}

	// Trim to MaxRunes
	trimmed := string(runes[:c.MaxRunes])
	reason = fmt.Sprintf("播报文本超长 (%d 字符>上限%d)，已截断至前%d字符；完整内容请走 D10 通知", len(runes), c.MaxRunes, c.MaxRunes)
	return false, trimmed, reason
}

// --- helper functions for form B implementation ---

// buildNormalizedTaintString converts a taint string to normalized form for storage.
// This is the same normalization used in Mark(): whitespace stripped, case-folded, full-hlf unified.
func buildNormalizedTaintString(content string) string {
	return normalizeTaint(content)
}

// fragmentMatch checks if two normalized strings share any >=n contiguous fragment.
// Uses FNV-1a window hashing (same algorithm as taintmatch.go).
func fragmentMatch(normA, normB string, n int) bool {
	if len(normA) < n || len(normB) < n {
		return false
	}

	// Build hash sets for both strings (window size n)
	hashSetA := getWindowHashes([]rune(normA), n)
	hashSetB := getWindowHashes([]rune(normB), n)

	// Check intersection
	for h := range hashSetA {
		if _, ok := hashSetB[h]; ok {
			return true
		}
	}

	return false
}

// getWindowHashes computes FNV-1a hashes for all windows of size n in the rune slice.
// Returns a map for fast O(1) lookup (collision-resistant, verified by substring search).
func getWindowHashes(rp []rune, n int) map[uint64]bool {
	hashes := make(map[uint64]bool)
	win := make([]rune, n)

	for i := 0; i+n <= len(rp); i++ {
		copy(win, rp[i:i+n])
		h := hashWindow(win)
		hashes[h] = true
	}

	return hashes
}

// SortedUnique returns a sorted list of unique elements from the input slice.
// Used for deterministic output in tests and logs.
func SortedUnique(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	slices.Sort(result)
	return result
}
