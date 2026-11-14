//go:build windows

package main

// Ticket 442's observation seam: the one place where a test can see, without a
// window, what each construction point actually handed the panel host as its
// WebView2 user-data folder.
//
// WHY IT EXISTS. The source-shape ruler (panel_profile_leaf_shape_windows_test.go)
// judges how a data path is SPELLED. It is a static heuristic, and by
// construction it can only ever know the shapes it was taught - ticket 442 is the
// third audit round on that one idea (eight holes found, five fixed, twelve more
// found). What the host RECEIVES is a runtime fact no spelling can dodge, so
// NewPanelManager hands the value it was given to this hook and a test asserts on
// that value. Ticket 442's 2026-11-14 ruling (甲＋落①) made this the primary
// criterion and froze the ruler as the secondary heuristic; the two files state
// their division of labour in their own headers.
//
// WHAT IT COSTS IN PRODUCTION: one nil check per construction, and nothing else.
// It is nil in every build, assigned by nothing outside *_test.go files, exported
// to nobody, and it keeps no record of its own - a test installs the recorder,
// the recorder is the only state, and t.Cleanup removes it again. This is the
// standard observability seam a harness gives itself before asserting behaviour;
// it is not product logic, it decides nothing, and with a nil hook this file is
// indistinguishable from the pre-ticket code.
//
// A func VALUE, not a slice, on purpose: the sink order is not the subject, and a
// package-level slice would keep every path alive for the whole run. The caller
// frame is deliberately not passed either - the recorder reads it itself with
// runtime.Caller, so production pays nothing for the diagnostic and the seam
// stays at one argument.
var panelDataPathSeam func(dataPath string)
