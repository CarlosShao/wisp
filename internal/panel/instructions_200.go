package panel

import "github.com/CarlosShao/wisp/internal/projctx"

// Ticket 200 AC#7: the panel carrier for "which project instruction files this
// run actually followed". The panel shows it; how it shows it is the UI leg's
// decision, not this file's. The view type is write-free on purpose, exactly
// like ModeView further up in composer.go: a view model with no setter cannot
// be written through from the panel side, and nothing here gives the panel a
// route back to a permission mode or an allowed-dir list.
//
// The producer fills it from the loop's loader (Loop.ProjectInstructionManifest
// / projctx.Bundle). NewSnapshot's signature is unchanged on purpose: this is
// an additive carrier, applied by the composition root through WithInstructions
// so no existing caller - and no existing wire-key nail - has to move.

// ProjectInstructionFile is one loaded instruction file, as the panel sees it.
type ProjectInstructionFile struct {
	// Path is the file the loader actually read.
	Path string `json:"path"`
	// Tier is "project" (found by walking up from the workspace) or "global"
	// (Wisp's own data dir, next to config.toml).
	Tier string `json:"tier"`
	// Depth is 0 for the workspace directory itself, larger further out, -1 for
	// the global tier.
	Depth int `json:"depth"`
	// Bytes is what this file contributed before any truncation.
	Bytes int `json:"bytes"`
	// TruncatedBytes is how much the context budget cut (omitted when zero).
	TruncatedBytes int `json:"truncatedBytes,omitempty"`
	// Dropped is true when the budget removed the whole file - the panel must
	// be able to say "this one did not take effect", not pretend it did.
	Dropped bool `json:"dropped,omitempty"`
	// DuplicateOf names the file this path collapsed into: one physical file
	// seen through two directory levels is injected only once.
	DuplicateOf string `json:"duplicateOf,omitempty"`
	// Source is the C25 provenance name stamped on the content ("fs.read").
	Source string `json:"source,omitempty"`
}

// ProjectInstructionsFromBundle turns one turn's load result into the carrier.
// A nil bundle yields an empty slice, which marshals as no key at all.
func ProjectInstructionsFromBundle(b *projctx.Bundle) []ProjectInstructionFile {
	if b == nil || len(b.Files) == 0 {
		return nil
	}
	out := make([]ProjectInstructionFile, 0, len(b.Files))
	for _, f := range b.Files {
		out = append(out, ProjectInstructionFile{
			Path:           f.Path,
			Tier:           f.Tier,
			Depth:          f.Depth,
			Bytes:          f.Bytes,
			TruncatedBytes: f.TruncatedBytes,
			Dropped:        f.Dropped,
			DuplicateOf:    f.DuplicateOf,
			Source:         f.Source,
		})
	}
	return out
}

// WithInstructions returns s carrying the loaded-instruction list. Snapshots
// are value types handed to the pump, so this is a copy, not a mutation.
func WithInstructions(s Snapshot, files []ProjectInstructionFile) Snapshot {
	if len(files) == 0 {
		s.Instructions = nil
		return s
	}
	s.Instructions = append([]ProjectInstructionFile(nil), files...)
	return s
}
