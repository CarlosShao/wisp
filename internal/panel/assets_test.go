package panel

// Ticket 77 AC#1 - the bundle inside the binary is the bundle the build
// produced.
//
// The failure this file exists to catch is not "dist/ is missing": that case
// is loud. It is the half-embedded case - the go:embed pattern carries
// index.html but drops the hashed assets next to it (a pattern typo, a stale
// dist/, a build interrupted). The page then loads, names an asset that is
// not in the binary, and the panel renders blank while Built() still reports
// true. Check() is the assertion that catches it, and the tests below pin
// Check() against synthetic bundles so it is verified even where no frontend
// build has ever run.

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestAnchorOnlyBundleIsNotBuiltAndFailsClosed(t *testing.T) {
	// Exactly what a clean checkout carries: the tracked anchor, nothing else.
	assets := newAssets(bundleFS(map[string]string{
		"dist/.gitkeep": "",
	}))
	if assets.Built() {
		t.Fatal("a bundle holding only the anchor file must report built=false")
	}
	if _, _, err := assets.Resolve("/"); err == nil || !strings.Contains(err.Error(), "not built") {
		t.Fatalf("Resolve on an unbuilt bundle must fail closed, got %v", err)
	}
	if _, err := assets.Check(); err == nil || !strings.Contains(err.Error(), "not built") {
		t.Fatalf("Check on an unbuilt bundle must fail closed, got %v", err)
	}
	if _, err := assets.Manifest(); err == nil {
		t.Fatal("Manifest on an unbuilt bundle must be an error, not an empty list")
	}
}

func TestCheckRejectsEntryNamingAnUnembeddedAsset(t *testing.T) {
	// index.html made it into the embed; the script it loads did not.
	assets := newAssets(bundleFS(map[string]string{
		"dist/index.html":           `<!doctype html><link href="./assets/index-AAA.css"><script src="./assets/index-BBB.js"></script>`,
		"dist/assets/index-AAA.css": "body{}",
	}))
	if !assets.Built() {
		t.Fatal("entry file is present, so this bundle is 'built' - that is the trap")
	}
	served, err := assets.Check()
	if err == nil {
		t.Fatalf("Check accepted a bundle missing %q; served=%v", "assets/index-BBB.js", served)
	}
	if !strings.Contains(err.Error(), "index-BBB.js") {
		t.Fatalf("Check must name the dangling reference, got: %v", err)
	}
	if len(served) != 1 || served[0] != "./assets/index-AAA.css" {
		t.Fatalf("expected the one resolvable ref to be reported before the failure, got %v", served)
	}
}

func TestCheckAcceptsAConsistentBundle(t *testing.T) {
	assets := newAssets(bundleFS(map[string]string{
		"dist/index.html":           `<link rel="stylesheet" href="./assets/index-AAA.css"/><script type="module" src="./assets/index-BBB.js"></script>`,
		"dist/assets/index-AAA.css": "body{}",
		"dist/assets/index-BBB.js":  "console.log(1)",
	}))
	served, err := assets.Check()
	if err != nil {
		t.Fatalf("consistent bundle rejected: %v", err)
	}
	want := []string{"./assets/index-AAA.css", "./assets/index-BBB.js"}
	if strings.Join(served, ",") != strings.Join(want, ",") {
		t.Fatalf("served refs = %v, want %v", served, want)
	}
	manifest, err := assets.Manifest()
	if err != nil {
		t.Fatalf("Manifest: %v", err)
	}
	if len(manifest) != 3 {
		t.Fatalf("expected 3 embedded files, got %d: %v", len(manifest), manifest)
	}
}

func TestCheckIgnoresOutsideAndFragmentReferences(t *testing.T) {
	// No runtime CDN (D29) means such a reference should not appear at all;
	// if one ever does, Check must not count it as an embedded asset, and must
	// not report it as a dangling one either.
	assets := newAssets(bundleFS(map[string]string{
		"dist/index.html":           `<a href="#top">x</a><script src="https://cdn.example/x.js"></script><link href="./assets/index-AAA.css">`,
		"dist/assets/index-AAA.css": "body{}",
	}))
	served, err := assets.Check()
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(served) != 1 || served[0] != "./assets/index-AAA.css" {
		t.Fatalf("only the bundle-relative ref should be checked, got %v", served)
	}
}

func TestBuiltinBundleIsCompleteOrAbsentNeverHalf(t *testing.T) {
	assets, err := BuiltinAssets()
	if err != nil {
		t.Fatalf("BuiltinAssets: %v", err)
	}
	if !assets.Built() {
		// The clean-checkout branch: legal, loud, and never a skip - a
		// checkout that has not run npm run build must report not-built
		// through every entry point instead of serving the anchor.
		if _, _, err := assets.Resolve("/"); err == nil || !strings.Contains(err.Error(), "not built") {
			t.Fatalf("unbuilt binary must refuse to serve the entry, got %v", err)
		}
		if _, err := assets.Check(); err == nil || !strings.Contains(err.Error(), "not built") {
			t.Fatalf("unbuilt binary must fail Check loudly, got %v", err)
		}
		t.Log("verdict: no frontend build in this checkout - bundle absent, not half-embedded")
		return
	}
	served, err := assets.Check()
	if err != nil {
		t.Fatalf("the binary carries an incomplete bundle: %v", err)
	}
	manifest, err := assets.Manifest()
	if err != nil {
		t.Fatalf("Manifest: %v", err)
	}
	entry, contentType, err := assets.Resolve("/")
	if err != nil {
		t.Fatalf("Resolve(/): %v", err)
	}
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		t.Fatalf("entry content type = %q, want html", contentType)
	}
	if !strings.Contains(strings.ToLower(string(entry)), "<!doctype html") {
		t.Fatalf("embedded entry is not an HTML document (%d bytes)", len(entry))
	}
	if len(served) == 0 {
		t.Fatal("built bundle references no assets at all - the page would be inert")
	}
	if len(manifest) < len(served)+1 {
		t.Fatalf("manifest carries %d files but the entry needs %d assets plus itself", len(manifest), len(served))
	}
	t.Logf("verdict: bundle complete - entry %d bytes, %d asset refs, %d embedded files", len(entry), len(served), len(manifest))
}

func TestResolveRefusesPathsOutsideTheBundle(t *testing.T) {
	assets := newAssets(bundleFS(map[string]string{
		"dist/index.html": "<!doctype html>",
		"dist/inside.js":  "x",
	}))
	for _, path := range []string{"../go.mod", "/..\\go.mod", "C:/Windows/win.ini", "/assets/../../go.mod", "./"} {
		if path == "./" {
			// "." normalises to the entry file and must NOT be refused.
			if _, _, err := assets.Resolve(path); err != nil {
				t.Fatalf("%q should resolve to the entry, got %v", path, err)
			}
			continue
		}
		if _, _, err := assets.Resolve(path); err == nil {
			t.Fatalf("Resolve(%q) returned bytes for a path outside dist/", path)
		}
	}
	if _, _, err := assets.Resolve("/missing.js"); err == nil {
		t.Fatal("unknown path must be an error, not a fallback to index.html")
	}
}

// bundleFS turns "dist/..." keyed content into the fs.FS shape BuiltinAssets
// gets after fs.Sub, so the synthetic bundles below are addressed exactly like
// the real embedded one.
func bundleFS(files map[string]string) fs.FS {
	tree := fstest.MapFS{}
	for name, body := range files {
		tree[name] = &fstest.MapFile{Data: []byte(body)}
	}
	sub, err := fs.Sub(tree, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
