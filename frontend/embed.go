// Package frontend holds the built panel bundle and nothing else.
//
// Ticket 77 AC#1: `npm run build` in this directory produces dist/, and the
// go:embed below is the ONLY way those bytes reach wisp.exe. There is no
// runtime CDN and no local HTTP server (D29 - the WebView2 host answers
// AddWebResourceRequested from these bytes), so the binary must be
// self-sufficient on a machine with no node, no npm and no network.
//
// dist/ is gitignored (see the note in the root .gitignore) except for the
// .gitkeep anchor, which exists so that `go build ./...` still resolves the
// embed pattern in a clean checkout that has not run the frontend build. With
// only the anchor present, panel.Assets.Built() reports false and the host
// shows its own "assets not built" state - it can never ship a page that
// quietly renders the placeholder.
package frontend

import "embed"

//go:embed all:dist
var distFS embed.FS

// Dist returns the embedded build output rooted at dist/.
func Dist() embed.FS {
	return distFS
}

// AnchorName is the tracked placeholder that keeps the embed pattern
// satisfiable before a build has run. Its presence as the ONLY entry is how
// panel detects "not built yet".
const AnchorName = ".gitkeep"
