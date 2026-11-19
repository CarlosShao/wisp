package panel

// Ticket 193's animation-switch carrier, panel half: the snapshot's
// panelAnimation key - the [panel] row the config screen does not have yet
// ("动画会有一个选项让用户选，切换最低就是关闭动画了，那就不会卡了", owner
// 09-28, ledger A373), now with the read half this page's row needs. The
// schema half (panel.animation, default "reduce" - the owner-pending
// placeholder) and the config-side setter/reader landed with the same
// ticket's config leg; THIS file is the wire half the page consumes with the
// frontend wiring leg.
//
// READER-GATED, THE 213/187/65/413 PATTERN, NO SECOND PARADIGM. A pump
// assembled without the reader sends no key at all; a pump that HAS one
// always sends it. The cmd/wisp assembly root wires the reader from
// internal/config's PanelAnimationOptions exactly like it already wires
// PanelWidth/PanelFontSize (run.go's pump sources + the resident pump) - that
// wiring is NOT this leg's to land (cmd/wisp/** is outside this dispatch's
// write face), and until it lands the production snapshot sends no
// panelAnimation key at all: absence, never a made-up tier.
//
// NO C17 method is added by this carrier, and NO TIMER arrives with it. The
// Go side schedules nothing for this key - the tier is a wire fact the page
// consumes; whether the page's rAF handles and CSS transitions really stop
// under "off" is the frontend half's proof, and the checklist for that leg
// (rAF not registered, transitions not attached, getAnimations() empty -
// "动画真停", not "跑但看不见") rides docs/evidence/s1/193-animation-opt-r1.md.
//
// THE KEY IS A REGISTERED GAP (composerTSGapLedger's PanelSnapshot entry):
// frontend/src/lib/panel.ts does not declare panelAnimation yet, and
// frontend/** is outside this leg's write face, so the two-way
// reconciliation judges the key through the ledger the same way ticket 167's
// queue numbers were judged. When the UI leg declares it, the ledger entry
// goes stale-red and that leg deletes it - the same procedure the appearance
// wave's four keys followed (4017962b3).

// PanelAnimationSection is the snapshot's panelAnimation key.
type PanelAnimationSection struct {
	// Tier is panel.animation as the live config holds it RIGHT NOW - one of
	// the schema's own three words ("full" | "reduce" | "off"; validate.go's
	// panelAnimationVocabulary is the single source). "" only when no Manager
	// backs the read - the same honest-empty rule PanelFontSizeSection's 0
	// follows: the page must treat "" as unset, never as "off".
	Tier string `json:"tier"`
}
