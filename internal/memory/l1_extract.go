package memory

// Ticket 29 AC#1: the real L1 profile extraction chain.
//
// WHAT WAS MISSING BEFORE THIS FILE. golden_extract.go pinned the AC#2 target
// and left a swap point, but the production side had the shape 29-v1 named:
// the host never assigned session.WarmWindow, and no C5 llm.LlmProvider
// call ever fed profile upserts. The rule-based stub in golden_extract_test.go
// proved the fixture self-consistent and nothing else ever ran. This file is
// the real chain; l1_warm.go is the seam the host hooks it to.
//
// THE SHAPE, and why it is this shape:
//
//	L1Extractor   ONE llm.LlmProvider.Stream call per task over that task's
//	              conversation turns. One call, not one per fact: the request
//	              prefix is identical for every slot a task yields, so a task
//	              costs at most ONE prompt-cache invalidation. That is what
//	              "batch-commit" means on the C5 side. (The AC#2 deliverable -
//	              one assistant-text golden-SSE recording per scripted
//	              conversation - lives under internal/memory/testdata/golden_sse
//	              and replays through this extractor. Those bytes are pinned to
//	              the reviewed fixture, not captured from a live endpoint: a
//	              live capture needs a real key and would pin whatever the model
//	              happened to answer that day.)
//	L1Pipeline    runs the extractor, then commits the WHOLE outcome through
//	              ONE store transaction (Store.ApplyExtraction): upserts first,
//	              forgets after, net per task. A failed extraction commits
//	              nothing at all.
//	L1WarmWindow  (l1_warm.go) the host seam: Record() on the response path
//	              does no LLM work and no DB work, and Fire() is the body
//	              assigned to session.WarmWindow, which hands the batch
//	              to the "memory-extract" roster goroutine.
//
// It implements the Extractor interface golden_extract.go defines, and the AC#2
// swap point (newGoldenExtractor) returns an *L1Extractor over the replayed C5
// chain without touching the fixture, the comparator or the DAO assertions. The
// golden suite therefore asserts the reviewed contract, not a live endpoint's
// mood.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The prompt is a package constant so the cache-stable prefix is byte-identical
// on every task (SPEC-05 section 4.1's discipline). The slot names and the value
// formats are the golden contract's, verbatim: the real chain has to emit the
// same normalized values the fixture pins (golden_extract.go, "Value formats
// are part of the contract").
const l1ExtractSystemPrompt = `You are Wisp's memory extractor. Read the conversation transcript and return the durable personal profile facts the USER stated about themselves.

Return exactly ONE JSON object and nothing else - no prose and no code fences:
{"upserts":[{"slot":"pref.tone","content":"concise"}],"forgets":[]}

Rules:
1. Only user turns are evidence. Never derive a slot from an assistant turn.
2. Slot names must stay in one of these namespaces, with these normalized values:
   pref.language     "zh" | "en"  (the language the ASSISTANT replies in, not the language the user typed)
   pref.tone         "casual" | "concise" | "direct"
   habit.work_hours  "weekday HH:MM-HH:MM" (24h, zero padded)
   habit.editor      "vscode" | "goland"
   fact.family       "mother birthday MM-DD" | "father birthday MM-DD"
3. A turn that asks to stop remembering a slot ("别记了" / "forget this") contributes that slot to "forgets" and never to "upserts".
4. Later turns override earlier turns for the same slot. Emit each slot at most once.
5. Do not invent slots, and do not record a one-off event or a single meeting time as a habit.
6. If nothing is worth remembering, return {"upserts":[],"forgets":[]}.`

// Defaults for the extraction call. maxTurns keeps the prompt bounded on a long
// task (the most recent turns are the ones that carry a preference change);
// maxTurnRunes keeps one pasted document from eating the whole request.
const (
	l1DefaultMaxTurns    = 24
	l1DefaultMaxTurnRune = 2000
	l1DefaultMaxTokens   = 512
)

// ErrL1NoProvider is returned when the chain is assembled without a C5 seam.
// It is a composition error, never a runtime one: a host that wires extraction
// without a provider would otherwise spend a Warm window discovering it.
var ErrL1NoProvider = errors.New("memory: L1 extraction has no C5 provider")

// L1Extractor is the REAL Extractor: one C5 call per task over the whole
// conversation. It is stateless per call and safe for concurrent use.
type L1Extractor struct {
	provider llm.LlmProvider
	model    string
	maxTurns int
	logf     func(format string, args ...any)
}

// L1ExtractorOption adjusts the extractor.
type L1ExtractorOption func(*L1Extractor)

// WithL1MaxTurns bounds how many of the task's most recent turns are sent.
func WithL1MaxTurns(n int) L1ExtractorOption {
	return func(e *L1Extractor) {
		if n > 0 {
			e.maxTurns = n
		}
	}
}

// WithL1Logf is the extraction chain's own log sink (debug lines only: a
// failed extraction must never disturb the user).
func WithL1Logf(f func(format string, args ...any)) L1ExtractorOption {
	return func(e *L1Extractor) {
		if f != nil {
			e.logf = f
		}
	}
}

// NewL1Extractor composes the chain. provider and model are both required: a
// provider with no model id would fail Request.Validate on every Warm window,
// which is a composition mistake worth refusing at assembly.
func NewL1Extractor(provider llm.LlmProvider, model string, opts ...L1ExtractorOption) (*L1Extractor, error) {
	if provider == nil {
		return nil, ErrL1NoProvider
	}
	if strings.TrimSpace(model) == "" {
		return nil, errors.New("memory: L1 extraction needs a model id")
	}
	e := &L1Extractor{
		provider: provider,
		model:    model,
		maxTurns: l1DefaultMaxTurns,
		logf:     func(string, ...any) {},
	}
	for _, o := range opts {
		o(e)
	}
	return e, nil
}

// Extractor is the seam golden_extract.go defines; this line is the compile-time
// proof that the real chain is what newGoldenExtractor() returns.
var _ Extractor = (*L1Extractor)(nil)

// Extract runs the one C5 call and parses its single JSON answer into the
// batched outcome for the task.
func (e *L1Extractor) Extract(ctx context.Context, conv GoldenConversation) (ExtractionOutcome, error) {
	if e == nil || e.provider == nil {
		return ExtractionOutcome{}, ErrL1NoProvider
	}
	req, err := e.buildRequest(conv)
	if err != nil {
		return ExtractionOutcome{}, err
	}
	col := llm.NewTurnCollector()
	streamErr := e.provider.Stream(ctx, req, func(ev llm.StreamEvent) error {
		col.Observe(ev)
		return nil
	})
	res := col.Result()
	if streamErr != nil {
		return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction call: %w", streamErr)
	}
	if res.Err != nil {
		return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction call: %w", res.Err)
	}
	switch res.Stop {
	case llm.StopEndTurn, llm.StopStopSequence, llm.StopMaxTokens:
		// max_tokens keeps the partial text: the parser below either reads a
		// complete JSON object or refuses the whole batch.
	case llm.StopCancelled:
		return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction cancelled: %w", context.Canceled)
	default:
		return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction stopped with %q", string(res.Stop))
	}
	out, perr := ParseExtractionReply(res.Text)
	if perr != nil {
		return ExtractionOutcome{}, perr
	}
	e.logf("memory: L1 extraction task=%s upserts=%d forgets=%d", conv.ID, len(out.Upserts), len(out.Forgets))
	return out, nil
}

// buildRequest renders the transcript into ONE user message behind the constant
// system prompt. The empty-content check is Request.Validate's, so this only
// refuses the shape that has no extraction evidence at all.
func (e *L1Extractor) buildRequest(conv GoldenConversation) (*llm.Request, error) {
	turns := e.window(conv.Turns)
	if len(turns) == 0 {
		return nil, errors.New("memory: L1 extraction: conversation has no user turn")
	}
	req := &llm.Request{
		Model:  e.model,
		System: []llm.Content{llm.TextPart{Text: l1ExtractSystemPrompt}},
		Messages: []llm.Message{{
			Role:    llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: renderTranscript(turns)}},
		}},
		MaxOutputTokens: l1DefaultMaxTokens,
		// Temperature 0: extraction is a classification, not a composition.
		Temperature: l1ZeroTemperature(),
	}
	// No CacheBreakpoints: the prefix that matters is the constant system
	// prompt, and a one-shot call per task has nothing to reuse. The
	// invalidation budget is held by the CALL COUNT instead (one per task).
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("memory: L1 extraction request: %w", err)
	}
	return req, nil
}

func l1ZeroTemperature() *float64 {
	v := 0.0
	return &v
}

// window keeps the most recent maxTurns turns and drops assistant-only tails
// that carry no extraction evidence. Order is preserved (oldest first), because
// the "later turn overrides earlier" rule is a contract rule.
func (e *L1Extractor) window(turns []GoldenTurn) []GoldenTurn {
	out := make([]GoldenTurn, 0, len(turns))
	for _, t := range turns {
		role := strings.ToLower(strings.TrimSpace(t.Role))
		if role != "user" && role != "assistant" {
			continue
		}
		text := truncateRunes(t.Text, l1DefaultMaxTurnRune)
		if strings.TrimSpace(text) == "" {
			continue
		}
		out = append(out, GoldenTurn{Role: role, Text: text})
	}
	if e.maxTurns > 0 && len(out) > e.maxTurns {
		out = out[len(out)-e.maxTurns:]
	}
	return out
}

// renderTranscript is the user message body: a labeled, oldest-first
// transcript. It carries no fixture metadata (id/title): the model must answer
// from the script, never from the case name.
func renderTranscript(turns []GoldenTurn) string {
	var b strings.Builder
	b.WriteString("Conversation transcript (oldest first):\n")
	for _, t := range turns {
		b.WriteString(t.Role)
		b.WriteString(": ")
		b.WriteString(t.Text)
		b.WriteByte('\n')
	}
	return b.String()
}

// truncateRunes is a local copy of the llm package's log truncation helper:
// internal/memory does not reach into that package's unexported surface.
func truncateRunes(s string, n int) string {
	if n <= 0 || s == "" {
		return s
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// extractionReply is the wire shape the prompt asks for. It is deliberately
// separate from ExtractionOutcome: a model reply is untrusted input and goes
// through normalizeOutcome before it can become an outcome.
type extractionReply struct {
	Upserts []SlotUpsert `json:"upserts"`
	Forgets []string     `json:"forgets"`
}

// ParseExtractionReply turns the model's single JSON object into the batched
// outcome. It tolerates the two shapes real models add anyway - a fenced code
// block and prose around the object - and refuses everything else: an invalid
// slot, an empty content, a duplicate or a reply with no JSON object at all all
// fail the WHOLE batch (a half-applied extraction is worse than none).
func ParseExtractionReply(text string) (ExtractionOutcome, error) {
	body, ok := firstJSONObject(text)
	if !ok {
		return ExtractionOutcome{}, errors.New("memory: L1 extraction: reply holds no JSON object")
	}
	var raw extractionReply
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&raw); err != nil {
		return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction: reply is not the expected JSON: %w", err)
	}
	return normalizeOutcome(raw.Upserts, raw.Forgets)
}

// firstJSONObject extracts the outermost object from a model reply: fenced
// blocks are stripped, then everything before the first '{' and after the last
// '}' is dropped.
func firstJSONObject(text string) (string, bool) {
	s := strings.TrimSpace(text)
	if i := strings.Index(s, "```"); i >= 0 {
		s = s[i+3:]
		if nl := strings.IndexByte(s, '\n'); nl >= 0 {
			s = s[nl+1:]
		}
		if j := strings.LastIndex(s, "```"); j >= 0 {
			s = s[:j]
		}
	}
	start := strings.IndexByte(s, '{')
	end := strings.LastIndexByte(s, '}')
	if start < 0 || end <= start {
		return "", false
	}
	return s[start : end+1], true
}

// normalizeOutcome applies the net semantics the golden contract pins:
// slot namespaces are checked, content must be non-empty, a later value for the
// same slot wins, a forget cancels an upsert of the same slot, and both lists
// come back sorted and duplicate-free so the batch is deterministic.
func normalizeOutcome(upserts []SlotUpsert, forgets []string) (ExtractionOutcome, error) {
	out := ExtractionOutcome{}
	bySlot := make(map[string]string, len(upserts))
	for _, u := range upserts {
		slot := strings.TrimSpace(u.Slot)
		content := strings.TrimSpace(u.Content)
		if err := validateGoldenSlot(slot, content); err != nil {
			return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction: %w", err)
		}
		bySlot[slot] = content
	}
	forgetSet := make(map[string]bool, len(forgets))
	for _, f := range forgets {
		slot := strings.TrimSpace(f)
		if slot == "" {
			return ExtractionOutcome{}, errors.New("memory: L1 extraction: empty forget slot")
		}
		if err := validateGoldenSlot(slot, "x"); err != nil {
			return ExtractionOutcome{}, fmt.Errorf("memory: L1 extraction: forget %q: %w", slot, err)
		}
		forgetSet[slot] = true
		delete(bySlot, slot) // net per task: a forget cancels a same-slot upsert
	}
	for slot, content := range bySlot {
		out.Upserts = append(out.Upserts, SlotUpsert{Slot: slot, Content: content})
	}
	for slot := range forgetSet {
		out.Forgets = append(out.Forgets, slot)
	}
	sort.Slice(out.Upserts, func(i, j int) bool { return out.Upserts[i].Slot < out.Upserts[j].Slot })
	sort.Strings(out.Forgets)
	return out, nil
}

// ExtractionCommit is what one batch wrote. Evicted is the number of profile
// rows the LRU cap removed while committing (the DAO logs each eviction with
// its slot; this is the count the caller can book).
type ExtractionCommit struct {
	Upserts int
	Forgets int
	Evicted int
}

// ApplyExtraction commits one task's whole outcome in ONE write transaction.
// This is the store half of "batch-commit": every slot of a task lands together
// or none of them does, the LRU cap is enforced once after the inserts (the
// same eviction the single-row UpsertProfile performs), and a forget targets a
// slot rather than a row id, so a forget for a slot the user no longer has is
// an idempotent no-op instead of an error that would abort the batch.
func (s *Store) ApplyExtraction(ctx context.Context, out ExtractionOutcome) (ExtractionCommit, error) {
	norm, err := normalizeOutcome(out.Upserts, out.Forgets)
	if err != nil {
		return ExtractionCommit{}, err
	}
	if len(norm.Upserts) == 0 && len(norm.Forgets) == 0 {
		return ExtractionCommit{}, nil
	}
	now := observe.NowWallUTC().Unix()
	commit := ExtractionCommit{}
	err = s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		for _, u := range norm.Upserts {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO profile(slot, content, source, updated_at) VALUES(?,?,?,?)
				 ON CONFLICT(slot) DO UPDATE SET content=excluded.content, source=excluded.source, updated_at=excluded.updated_at`,
				u.Slot, u.Content, string(ProfileExtracted), now); err != nil {
				return err
			}
			commit.Upserts++
		}
		if len(norm.Upserts) > 0 {
			// keepSlot "" matches no row (slots are non-empty), so every row is
			// an eviction candidate and the oldest wins - the same ordering
			// UpsertProfile uses. The rows written above share the newest
			// updated_at and the highest ids, so they are last in that order.
			var before int
			if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM profile`).Scan(&before); err != nil {
				return err
			}
			if err := evictProfileLRU(ctx, tx, s.logger, ""); err != nil {
				return err
			}
			var after int
			if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM profile`).Scan(&after); err != nil {
				return err
			}
			commit.Evicted = before - after
		}
		for _, slot := range norm.Forgets {
			res, err := tx.ExecContext(ctx, `DELETE FROM profile WHERE slot=?`, slot)
			if err != nil {
				return err
			}
			if n, _ := res.RowsAffected(); n > 0 {
				commit.Forgets++
			}
		}
		return nil
	})
	if err != nil {
		return ExtractionCommit{}, fmt.Errorf("memory: commit L1 extraction: %w", err)
	}
	return commit, nil
}

// L1Pipeline is one task's extraction exec: one C5 call, one store transaction.
// It is synchronous by contract; the Warm window owns the goroutine it runs on
// (l1_warm.go).
type L1Pipeline struct {
	store *Store
	ex    Extractor
	logf  func(format string, args ...any)
}

// L1PipelineOption adjusts the pipeline.
type L1PipelineOption func(*L1Pipeline)

// WithL1PipelineLogf is the pipeline's log sink (commit lines + failures).
func WithL1PipelineLogf(f func(format string, args ...any)) L1PipelineOption {
	return func(p *L1Pipeline) {
		if f != nil {
			p.logf = f
		}
	}
}

// NewL1Pipeline composes the exec. Both halves are required: a pipeline with no
// extractor would claim work it cannot do, and one with no store would call C5
// and throw the answer away.
func NewL1Pipeline(store *Store, ex Extractor, opts ...L1PipelineOption) (*L1Pipeline, error) {
	if store == nil {
		return nil, errors.New("memory: L1 pipeline needs a store")
	}
	if ex == nil {
		return nil, errors.New("memory: L1 pipeline needs an extractor")
	}
	p := &L1Pipeline{store: store, ex: ex, logf: func(string, ...any) {}}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// RunTask runs one task end to end: ONE extraction call over the whole
// conversation, then ONE batched commit. An extraction failure returns before
// the store is touched, so a broken reply can never half-write the profile.
func (p *L1Pipeline) RunTask(ctx context.Context, conv GoldenConversation) (ExtractionOutcome, error) {
	if p == nil || p.ex == nil || p.store == nil {
		return ExtractionOutcome{}, errors.New("memory: L1 pipeline is not assembled")
	}
	out, err := p.ex.Extract(ctx, conv)
	if err != nil {
		return ExtractionOutcome{}, err
	}
	commit, err := p.store.ApplyExtraction(ctx, out)
	if err != nil {
		return out, err
	}
	p.logf("memory: L1 extraction committed task=%s upserts=%d forgets=%d evicted=%d",
		conv.ID, commit.Upserts, commit.Forgets, commit.Evicted)
	return out, nil
}
