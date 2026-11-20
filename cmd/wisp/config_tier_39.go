package main

// Ticket 39 AC#1's tier-events leg: the registry-to-receipt tier mapping the
// settings leg stamps from.
//
// THE TIER JUDGMENT IS THE REGISTRY'S, NEVER A HAND STAMP. tiers.go's
// TierRegistry is the one machine-readable answer to "which key lands on which
// D36 tier", and TierOf is its accessor. This helper maps the registry's own
// word for a written key path onto the receipt vocabulary the panel carries
// (panel.EffectiveTier), so a receipt's tier sentence and a reload event's
// tier field are derived from the same table the write path obeys:
//
//	hot     -> EffectiveNow     (applied to this process's memory on write)
//	reload  -> EffectiveNow     (applied AND the owning subsystem got
//	                             Manager.OnReload - sectionwrite.go fires it)
//	restart -> EffectiveRestart (file written, memory keeps the boot value)
//	locked  -> EffectiveRestart (the locked marker is not a reload tier; the
//	                             conservative word, same as the old blanket)
//
// SCOPE, STATED SO NOBODY WIDENS IT BY ACCIDENT: this mapping is consulted for
// the ticket-39 field (app_language, the per-key "app.language" row) and,
// since ticket 193's wiring leg, for panel_animation - whose tier comes off
// the same registry's whole-section "panel" row (tiers.go: "panel" = hot,
// D36's [panel] 行), keyed by the key path the guarded write itself reports.
// The older fields keep their measured tier stamps (187's next_task
// correction, 65's applier-measured now, 422's static restart for ball.size,
// the appearance wave's [panel] restart-with-Note shape) - those are
// per-field measurements with their own evidence pins, not registry rows, and
// folding them into this table would re-litigate tickets 187/255/422 with no
// new fact.
import (
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

func registryEffectiveTier39(key string) (panel.EffectiveTier, bool) {
	tier, ok := config.TierOf(key)
	if !ok {
		return "", false
	}
	// TierOf answers the string view of the registry's typed value (tiers.go's
	// own contract); the comparison is against the constants themselves.
	switch config.Tier(tier) {
	case config.TierHot, config.TierReload:
		return panel.EffectiveNow, true
	case config.TierRestart:
		return panel.EffectiveRestart, true
	default:
		// "locked" and any future marker: the conservative word, never a
		// made-up fourth tier.
		return panel.EffectiveRestart, true
	}
}
