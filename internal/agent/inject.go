package agent

import (
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/CarlosShao/wisp/internal/llm"
)

// D15(1)+(2) two-tier tool injection, done locally with ZERO LLM round-trips
// (SPEC-05 §4.2): builtin/resident tools are always injected; third-party or
// long-tail tools are selected by a local BM25/keyword score over their
// name+description, top-K with K<=5. When the retrieval misses, the model is
// given the list_tools meta-tool so it can self-serve the full directory from
// inside the same loop turn (D15(2): "检索未命中 → list_tools 元工具自查完整目
// 录（循环内普通工具调用，非额外 LLM 往返）").

// MaxRetrievedTools is the D15(2) top-K ceiling for third-party tools.
const MaxRetrievedTools = 5

// ToolListName is the D34/S1 meta-tool that answers the full directory from
// inside the loop.
const ToolListName = "list_tools"

// Injection is one turn's tool-injection decision.
type Injection struct {
	// Resident are the always-injected builtin tools.
	Resident []ToolInfo
	// Retrieved are the BM25/keyword-selected third-party tools (<=K).
	Retrieved []ToolInfo
	// AdvertiseListTools is true when the retrieval missed, so list_tools is
	// offered even if the provider does not declare it (D15(2) fallback).
	AdvertiseListTools bool
	// IndexText is the (2) suffix section: the selected third-party index.
	IndexText string
	// ToolDefs is the wire declaration set for the turn.
	ToolDefs []llm.ToolDef
}

// Inject runs the two-tier selection for one query against the directory.
// Budgets bound the (2) suffix index (reference 300 tokens at 128k).
func Inject(query string, dir []ToolInfo, b Budgets) Injection {
	var resident, candidates []ToolInfo
	for _, t := range dir {
		if t.Resident {
			resident = append(resident, t)
		} else {
			candidates = append(candidates, t)
		}
	}

	inj := Injection{Resident: resident}
	scored := bm25Select(query, candidates, MaxRetrievedTools, b.SectionToolIndex)
	inj.Retrieved = scored
	inj.AdvertiseListTools = len(scored) == 0 && len(candidates) > 0
	if len(scored) > 0 {
		inj.IndexText = indexText(scored, b.SectionToolIndex)
	}

	defs := make([]llm.ToolDef, 0, len(resident)+len(scored)+1)
	for _, t := range append(append([]ToolInfo{}, resident...), scored...) {
		defs = append(defs, t.toolDef())
	}
	if inj.AdvertiseListTools && !dirHas(dir, ToolListName) {
		defs = append(defs, llm.ToolDef{
			Name:        ToolListName,
			Description: "List the full tool directory (name, description, risk level).",
			Parameters:  listToolsSchema,
		})
	}
	inj.ToolDefs = defs
	return inj
}

func dirHas(dir []ToolInfo, name string) bool {
	for _, t := range dir {
		if t.Name == name {
			return true
		}
	}
	return false
}

func (t ToolInfo) toolDef() llm.ToolDef {
	params := t.Parameters
	if len(params) == 0 {
		params = emptyObjectSchema
	}
	return llm.ToolDef{Name: t.Name, Description: t.Description, Parameters: params}
}

// indexText renders the (2) tool index section for a tool set, clipped to the
// budget (clipping is by line, never mid-entry, so the section stays
// parseable for the model).
func indexText(tools []ToolInfo, budget int) string {
	var kept []string
	used := 0
	for _, t := range tools {
		line := "- " + t.Name + ": " + strings.TrimSpace(t.Description)
		cost := ApproxTokens(line) + 1
		if used+cost > budget {
			break
		}
		kept = append(kept, line)
		used += cost
	}
	if len(kept) == 0 {
		return ""
	}
	return strings.Join(kept, "\n")
}

// ---------------------------------------------------------------------------
// Local BM25 (Okapi, k1=1.5, b=0.75) over name+description. This is a
// retrieval helper, not a model call: no network, no tokens, single-digit ms.

const (
	bm25K1 = 1.5
	bm25B  = 0.75
)

// Tokenize splits text into index terms: ASCII/latin words lowercased, and
// CJK as single characters plus adjacent character bigrams (Chinese tool
// descriptions carry no spaces, so word segmentation would fail on them).
func Tokenize(s string) []string {
	s = strings.ToLower(s)
	var toks []string
	var word []rune
	flush := func() {
		if len(word) > 0 {
			toks = append(toks, string(word))
			word = word[:0]
		}
	}
	var prev rune
	for _, r := range s {
		switch {
		case r == '_' || r == '-' || r == '.':
			// separator inside a tool name: namespace.action stays intact,
			// snake_case words split.
			flush()
			prev = 0
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if isCJK(r) {
				flush()
				toks = append(toks, string(r))
				if prev != 0 && isCJK(prev) {
					toks = append(toks, string([]rune{prev, r}))
				}
			} else {
				word = append(word, r)
			}
			prev = r
		default:
			flush()
			prev = 0
		}
	}
	flush()
	return toks
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r)
}

// bm25Select scores candidates against query and keeps the top k whose score
// is above zero and whose index lines fit budget tokens.
func bm25Select(query string, cands []ToolInfo, k, budget int) []ToolInfo {
	if len(cands) == 0 || strings.TrimSpace(query) == "" {
		return nil
	}
	q := Tokenize(query)
	if len(q) == 0 {
		return nil
	}

	docs := make([][]string, len(cands))
	totalLen := 0.0
	df := map[string]int{}
	for i, c := range cands {
		docs[i] = Tokenize(c.Name + " " + c.Description)
		totalLen += float64(len(docs[i]))
		seen := map[string]bool{}
		for _, t := range docs[i] {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}
	avgLen := totalLen / float64(len(docs))

	n := float64(len(docs))
	type hit struct {
		idx   int
		score float64
	}
	var hits []hit
	for i, doc := range docs {
		tf := map[string]int{}
		for _, t := range doc {
			tf[t]++
		}
		dl := float64(len(doc))
		var score float64
		for _, term := range q {
			f := float64(tf[term])
			if f == 0 {
				continue
			}
			idf := math.Log(1 + (n-float64(df[term])+0.5)/(float64(df[term])+0.5))
			score += idf * (f * (bm25K1 + 1)) /
				(f + bm25K1*(1-bm25B+bm25B*dl/maxF(avgLen, 1)))
		}
		if score > 0 {
			hits = append(hits, hit{idx: i, score: score})
		}
	}
	sort.SliceStable(hits, func(a, b int) bool {
		if hits[a].score != hits[b].score {
			return hits[a].score > hits[b].score
		}
		return cands[hits[a].idx].Name < cands[hits[b].idx].Name
	})

	var out []ToolInfo
	used := 0
	for _, h := range hits {
		if len(out) >= k {
			break
		}
		t := cands[h.idx]
		cost := ApproxTokens("- "+t.Name+": "+strings.TrimSpace(t.Description)) + 1
		if used+cost > budget && len(out) > 0 {
			break
		}
		used += cost
		out = append(out, t)
	}
	return out
}

func maxF(v, min float64) float64 {
	if v < min {
		return min
	}
	return v
}
