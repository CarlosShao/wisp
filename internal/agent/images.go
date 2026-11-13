// Package agent - file-level note: images.go is the memory half of the C7
// image seam (ticket 23 AC#2). The package's own charter - the ReAct loop,
// budgeting, cancellation - is documented in doc.go; this file adds the
// process-local image store the loop hands screenshot bytes to.
//
// Ticket 23 AC#2 - the memory half of the C7 image seam.
//
// screen.capture hands the loop the raw pixels of a window it just captured.
// C7 (SPEC-05 sec 3.3) fixes what may happen to them: the model may see the
// image, and nothing else may. ImagePart.BytesRef is "a reference to memory or
// a temp file" that is NEVER persisted to the DB, NEVER written to a log and
// NEVER included in a diagnostics bundle.
//
// This file implements the strictest legal choice of those two, the
// memory-only one: bytes live in a process-local store keyed by an opaque
// mem:<n> handle. Nothing here opens a file, so there is no path to leak and no
// artifact a diagnostics bundle could pick up. Two consequences are worth
// stating out loud:
//
//   - The bytes die with the session. Loop.Reset clears the store at the same
//     moment it drops the history that referenced it, so a resolvable handle
//     and a live history cannot drift apart. Clearing at TASK end would be
//     wrong: the conversation keeps every earlier screenshot and the next turn
//     would fail to resolve one of them.
//   - A tool cannot express "keep this for later" - by design. A screenshot
//     that outlives the session is the user's explicit save, not a tool
//     default (D34 note on screen.capture: bytes_ref is memory-only unless the
//     user saves it).
package agent

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

const (
	// maxImageBytes is the ceiling for ONE image. A screenshot of a 4K desktop
	// as PNG is well below it; anything above is a sign the caller handed us
	// the wrong buffer (a raw BGRA frame, a whole file) and must not be copied
	// into a request that every later turn pays for.
	maxImageBytes = 4 << 20

	// imageSessionBudget bounds the WHOLE session's in-memory images, because
	// the store's lifetime is the session's. The C7 history accumulates one
	// image per capture, so an unbounded store is an unbounded leak with a
	// UI-shaped trigger.
	imageSessionBudget = 32 << 20
)

var (
	// errImageEmpty and errImageTooLarge are structural refusals: the bytes
	// cannot be what the caller says they are.
	errImageEmpty    = errors.New("图像没有任何字节")
	errImageTooLarge = errors.New("图像超过单图上限")
	// errImageBudget is the session bound: the image is real and well-formed,
	// it just no longer fits. Both refusals are model-facing, because the only
	// correct behaviour is to TELL the model it did not receive pixels.
	errImageBudget = errors.New("会话内存里的图像已装满")
	// errImageStoreMissing is a wiring fault, not a policy refusal: a Loop built
	// by hand (no New) has nowhere to keep pixels. It is stated to the model in
	// the same breath as a budget refusal rather than dropped silently.
	errImageStoreMissing = errors.New("本循环没有可用的内存图像仓库")
)

// storedImage is one in-memory image. data is a private copy: the tool's slice
// is never retained, so a provider that reuses its buffer cannot corrupt a
// screenshot the history still points at.
type storedImage struct {
	mime string
	data []byte
	alt  string
}

// imageStore is the session's memory-only image store. It is safe for concurrent
// use (tool calls run concurrently) even though the loop currently attaches
// images from one goroutine.
type imageStore struct {
	mu     sync.Mutex
	budget int
	used   int
	seq    int
	items  map[string]storedImage
}

// newImageStore builds a store with the given total budget; a non-positive
// budget falls back to imageSessionBudget (never to "unlimited").
func newImageStore(budget int) *imageStore {
	if budget <= 0 {
		budget = imageSessionBudget
	}
	return &imageStore{budget: budget, items: map[string]storedImage{}}
}

// put stores one image and returns its opaque handle. The handle is the only
// thing that ever leaves this package's memory: it is not a path, it has no
// separator, and it means nothing outside this store.
func (s *imageStore) put(mime string, data []byte, alt string) (string, error) {
	if len(data) == 0 {
		return "", errImageEmpty
	}
	if len(data) > maxImageBytes {
		return "", fmt.Errorf("%w（%d 字节 > %d）", errImageTooLarge, len(data), maxImageBytes)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used+len(data) > s.budget {
		return "", fmt.Errorf("%w（已用 %d / %d 字节）", errImageBudget, s.used, s.budget)
	}
	s.seq++
	ref := "mem:" + strconv.Itoa(s.seq)
	s.items[ref] = storedImage{mime: mime, data: append([]byte(nil), data...), alt: alt}
	s.used += len(data)
	return ref, nil
}

// resolve hands the adapter the bytes of one handle (Request.ResolveBytes).
// The copy keeps a provider that mutates its input from editing the session's
// image. The error deliberately does not quote the handle: a ref is a
// conversation-level secret and this error can reach a log.
func (s *imageStore) resolve(ref string) (string, []byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[ref]
	if !ok {
		return "", nil, errors.New("agent: 图像句柄在本会话内已不存在（会话已结束或从未写入）")
	}
	return it.mime, append([]byte(nil), it.data...), nil
}

// clear drops every image and every byte of budget (session end).
func (s *imageStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = map[string]storedImage{}
	s.used = 0
}

// imageDropNotice is what the model reads when an image could not be attached.
// Silence is not an option: a model that believes it saw a screen it never
// received will describe pixels out of thin air (D37/D22 honesty rule).
func imageDropNotice(alt string, err error) string {
	if alt == "" {
		return fmt.Sprintf("（图像未送入模型：%s）", err)
	}
	return fmt.Sprintf("（图像 %q 未送入模型：%s）", alt, err)
}
