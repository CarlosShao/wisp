package panel

// Embedded panel assets (ticket 77 AC#1).
//
// This is the seam between the Go binary and the WebView bundle: it owns
// nothing but read access to the embedded dist tree. The WebView2 host
// (ticket 33) calls Resolve for every request the page makes, and per D29
// there is no local HTTP server to fall back on - if a byte is not in the
// embed, it does not exist at runtime.
//
// The package is deliberately pure Go and platform-neutral: it compiles and
// is testable on any GOOS, so the "the binary carries its own UI" claim is
// checked by go test rather than by a manual look at a window.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/frontend"
)

// EntryFile is what the host loads for the panel's root URL.
const EntryFile = "index.html"

// errNotBuilt is returned when the embed only carries the tracked anchor file,
// i.e. nobody ran "npm run build" before "go build". Callers must show this
// instead of rendering the anchor.
var errNotBuilt = errors.New("panel: embedded assets are not built (run npm run build in frontend/)")

// Assets is a read view over the embedded dist tree.
type Assets struct {
	tree  fs.FS
	built bool
}

// BuiltinAssets opens the bundle compiled into this binary.
func BuiltinAssets() (*Assets, error) {
	sub, err := fs.Sub(frontend.Dist(), "dist")
	if err != nil {
		return nil, fmt.Errorf("panel: cannot open embedded dist: %w", err)
	}
	a := &Assets{tree: sub}
	if _, err := fs.Stat(sub, EntryFile); err == nil {
		a.built = true
	}
	return a, nil
}

// Built reports whether the binary really carries a built panel, as opposed
// to just the anchor file that keeps the go:embed pattern valid in a clean
// checkout.
func (a *Assets) Built() bool { return a != nil && a.built }

// Resolve maps a request path to embedded bytes plus a Content-Type.
//
// Path handling is string-only on purpose: the embed tree is addressed with
// slash paths, and routing through the OS path layer is what lets ".." and
// drive-relative names look harmless on Windows. Nothing outside dist/ can be
// reached, and an unknown path is an error rather than a fallback - a silent
// index.html for every miss would turn a broken bundle reference into a blank
// panel that still looks alive.
func (a *Assets) Resolve(requestPath string) ([]byte, string, error) {
	if !a.Built() {
		return nil, "", errNotBuilt
	}
	name := strings.TrimPrefix(strings.TrimPrefix(requestPath, "/"), "./")
	if name == "" || name == "." {
		name = EntryFile
	}
	if strings.Contains(name, "..") || strings.Contains(name, "\\") || strings.Contains(name, ":") {
		return nil, "", fmt.Errorf("panel: refused request path %q", requestPath)
	}
	data, err := fs.ReadFile(a.tree, name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, "", fmt.Errorf("panel: %q is not in the embedded bundle: %w", name, err)
	}
	if err != nil {
		return nil, "", fmt.Errorf("panel: read %q: %w", name, err)
	}
	return data, contentTypeOf(name), nil
}

// Manifest lists every embedded file with its size and a fingerprint, sorted.
// It exists for diagnostics: "wisp panel-assets -manifest" prints exactly what
// the binary carries, which is the only way to prove on a node-free machine
// that the page came from the executable and not from disk.
func (a *Assets) Manifest() ([]string, error) {
	if !a.Built() {
		return nil, errNotBuilt
	}
	var out []string
	err := fs.WalkDir(a.tree, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(a.tree, p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		out = append(out, fmt.Sprintf("%s %d %s", p, len(data), hex.EncodeToString(sum[:8])))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("panel: walk embedded bundle: %w", err)
	}
	sort.Strings(out)
	return out, nil
}

// contentTypeOf covers the shapes Vite emits. Anything unlisted is refused by
// the host rather than served as text/html, because WebView2 would then try to
// navigate the panel to a script.
func contentTypeOf(name string) string {
	switch strings.ToLower(name[strings.LastIndex(name, ".")+1:]) {
	case "html":
		return "text/html; charset=utf-8"
	case "js", "mjs":
		return "text/javascript; charset=utf-8"
	case "css":
		return "text/css; charset=utf-8"
	case "json":
		return "application/json; charset=utf-8"
	case "svg":
		return "image/svg+xml"
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "woff2":
		return "font/woff2"
	case "ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
