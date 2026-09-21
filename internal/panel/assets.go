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
	"regexp"
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
	return newAssets(sub), nil
}

// newAssets decides the built / not-built state for any tree the same way, so
// the tests can hand it a synthetic bundle - including the "embed carried the
// entry but not the assets" case that only happens by accident in the real one.
func newAssets(tree fs.FS) *Assets {
	a := &Assets{tree: tree}
	if _, err := fs.Stat(tree, EntryFile); err == nil {
		a.built = true
	}
	return a
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

// entryRefRe matches the asset references Vite writes into index.html
// (src="./assets/index-<hash>.js", href="./assets/index-<hash>.css").
var entryRefRe = regexp.MustCompile(`(?i)\b(?:src|href)\s*=\s*["']([^"']+)["']`)

// Check proves the embedded bundle is internally complete, which is the
// failure mode "go:embed shipped less than npm run build produced" actually
// looks like: the entry file is there, the page loads, and then every hashed
// asset it names 404s, so the panel renders blank while the binary still
// reports built=true.
//
// It returns the bundle-relative names it resolved so a caller can print how
// many bytes the page actually depends on. Anything the entry references that
// cannot be resolved from the same tree is an error naming that reference.
func (a *Assets) Check() ([]string, error) {
	entry, _, err := a.Resolve(EntryFile)
	if err != nil {
		return nil, err
	}
	var served []string
	seen := map[string]bool{}
	for _, match := range entryRefRe.FindAllStringSubmatch(string(entry), -1) {
		ref := match[1]
		if cut := strings.IndexAny(ref, "?#"); cut >= 0 {
			ref = ref[:cut]
		}
		if ref == "" || strings.HasPrefix(ref, "#") || strings.HasPrefix(ref, "//") ||
			strings.HasPrefix(ref, "data:") || strings.Contains(ref, "://") || seen[ref] {
			continue
		}
		seen[ref] = true
		if _, _, err := a.Resolve(ref); err != nil {
			return served, fmt.Errorf("panel: %s references %q but the embedded bundle cannot serve it", EntryFile, ref)
		}
		served = append(served, ref)
	}
	return served, nil
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
