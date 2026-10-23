package perm

import "github.com/CarlosShao/wisp/internal/risk"

// IsWidening reports whether the switch from -> to removes questions instead of
// adding them.
//
// This is the ONE definition of the direction. It lives here, next to the switch
// it guards, because "loosening the safety posture" is permission policy rather
// than panel vocabulary: the store is the lower layer, so both the store's own
// gate (Set) and the panel's caller-side gate (panel.ModeWriteHandler) can call
// it without either package owning a private copy. Ticket 367 is why that
// matters - the two transition points had two different judges (one asked "is
// the target auto_approve?", the other "is the direction widening?"), and a
// widening that was not auto_approve fell through the pair of them.
//
// The comparison is a raw `to > from` over risk.Mode's constants, which is only
// sound because internal/risk/mode.go orders them strict -> loose and says "Do
// not reorder these constants". Two tests pin that premise from both sides: the
// ladder test in internal/panel (which drives this function over
// risk.ModeNames()) and TestIsWideningOnTheDocumentedLadder here.
//
// Narrowing is deliberately NOT a refusal anywhere in this package: getting
// stricter never needs a confirmation, and blocking it would guard nothing
// (R20's direction-asymmetry).
func IsWidening(from, to risk.Mode) bool { return to > from }
