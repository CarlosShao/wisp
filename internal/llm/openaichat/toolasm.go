package openaichat

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// toolAssembler reassembles streamed tool_calls fragments into the C6
// ToolCallStart / ArgsDelta / End sequence.
//
// The wire keys fragments by "index"; id and name ride the first fragment of
// each index. Fragments may therefore be buffered until the identity is
// complete, because the normalized contract requires Start before any
// ArgsDelta. All calls close in index order on finish_reason.
type toolAssembler struct {
	loose  bool
	byIdx  map[int]*toolAcc
	order  []int // stable wire-index order for the End events
	closed bool
}

type toolAcc struct {
	id, name string
	started  bool     // Start emitted; args flow through afterwards
	buffered [][]byte // argument fragments seen before Start was possible
	ended    bool
}

func newToolAssembler(loose bool) *toolAssembler {
	return &toolAssembler{loose: loose, byIdx: map[int]*toolAcc{}}
}

// fragment ingests one wire fragment and returns the events it produces.
func (a *toolAssembler) fragment(f wireToolFrag) ([]llm.StreamEvent, error) {
	if a.closed {
		return nil, observe.New(observe.ClassProvider,
			"openai-chat: tool_calls fragment after finish_reason")
	}

	idx := 0
	if f.Index != nil {
		idx = *f.Index
	} else if !a.loose {
		return nil, observe.New(observe.ClassProvider,
			"openai-chat: tool_calls fragment without index (set compat.loose to tolerate)")
	} else {
		idx = len(a.order) // arrival order
	}

	acc, ok := a.byIdx[idx]
	if !ok {
		acc = &toolAcc{}
		a.byIdx[idx] = acc
		a.order = append(a.order, idx)
	}
	if acc.id == "" {
		acc.id = f.ID
	}
	if acc.name == "" {
		acc.name = f.Function.Name
	}

	if !acc.started {
		if len(f.Function.Arguments) > 0 {
			acc.buffered = append(acc.buffered, []byte(f.Function.Arguments))
		}
		if acc.id != "" && acc.name != "" {
			out := []llm.StreamEvent{{
				Type: llm.EvToolCallStart, ToolCallID: acc.id, ToolName: acc.name,
			}}
			for _, b := range acc.buffered {
				out = append(out, llm.StreamEvent{
					Type: llm.EvToolCallArgsDelta, ToolCallID: acc.id, ArgsDelta: string(b),
				})
			}
			acc.buffered = nil
			acc.started = true
			return out, nil
		}
		return nil, nil
	}

	// started: plain argument flow-through.
	if len(f.Function.Arguments) > 0 {
		return []llm.StreamEvent{{
			Type: llm.EvToolCallArgsDelta, ToolCallID: acc.id, ArgsDelta: string(f.Function.Arguments),
		}}, nil
	}
	return nil, nil
}

// finish closes every call in wire-index order (ToolCallEnd). In loose mode
// a never-identified call (no id/name) gets synthesized values so downstream
// code has a stable key; strict mode reports the malformed stream instead.
func (a *toolAssembler) finish() ([]llm.StreamEvent, error) {
	if a.closed {
		return nil, nil
	}
	a.closed = true
	var out []llm.StreamEvent
	for _, idx := range a.order {
		acc := a.byIdx[idx]
		if acc.ended {
			continue
		}
		if !acc.started {
			if acc.id == "" {
				if !a.loose {
					return nil, observe.New(observe.ClassProvider,
						fmt.Sprintf("openai-chat: tool_calls[%d] closed without an id (set compat.loose to tolerate)", idx))
				}
				acc.id = fmt.Sprintf("call_unnamed_%d", idx)
			}
			if acc.name == "" {
				if !a.loose {
					return nil, observe.New(observe.ClassProvider,
						fmt.Sprintf("openai-chat: tool_calls[%d] closed without a function name (set compat.loose to tolerate)", idx))
				}
				acc.name = "unknown_tool"
			}
			out = append(out, llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: acc.id, ToolName: acc.name})
			for _, b := range acc.buffered {
				out = append(out, llm.StreamEvent{
					Type: llm.EvToolCallArgsDelta, ToolCallID: acc.id, ArgsDelta: string(b),
				})
			}
			acc.buffered = nil
			acc.started = true
		}
		acc.ended = true
		out = append(out, llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: acc.id})
	}
	return out, nil
}
