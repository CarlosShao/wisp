package main

// modeRequest114r2 is one panel.mode.request envelope, spelled the way the page
// spells it: bridge.go accepts exactly one sender identity, and anything else is
// refused by a door in front of the one under test.
//
// Cross-platform test helper: it is a pure string builder over the
// cross-platform test const inboundPanelSource, so both the windows-only
// panel_mode_confirm_114_windows_test.go and the cross-platform
// r9_mode_card_test.go can use it. It lived in the windows file until the
// Linux build of cmd/wisp needed r9_mode_card_test.go to type-check (ticket
// 424: a windows-only test helper referenced from an untagged test file is
// undefined on GOOS=linux, which reddens `go vet ./...`).
func modeRequest114r2(rid, to string) string {
	return `{"method":"panel.mode.request","requestId":"` + rid + `","source":"` +
		inboundPanelSource + `","to":"` + to + `"}`
}
