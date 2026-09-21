package panel

// Frontend hygiene gates that exist because the repo's own gates do not reach
// this tree yet (ticket 77 AC#4 and AC#5).
//
//  1. STATELESSNESS (AC#5, PLAN.md:1044). The panel may hold nothing across a
//     WebView restart: no Web Storage, no IndexedDB, no cookies, no service
//     worker cache. Everything on screen comes from a Go-side push, so a
//     frontend that remembered things would be a second state holder and the
//     first place a divergence shows up is the history list.
//
//  2. ONE STYLE SOURCE (AC#2's other half). A colour literal may only exist in
//     the file generated from C21. This walks the import closure of
//     src/main.tsx - the files that actually reach the shipped bundle - and
//     rejects a hex/rgb(a) written anywhere else. Vendored demo components that
//     nothing imports are allowed their upstream colours precisely BECAUSE the
//     same walk proves they are unreachable; that exemption is earned, not
//     assumed (see TestVendoredDemoComponentsAreNotMounted).
//
//  3. ZERO EMOJI IN frontend/ (AC#4's ban #8 half). tools/d22scan's ban #8
//     walks design/, internal/ and cmd/ - frontend/ is not in its declared
//     scope, and arming it is a scanner change this ticket may not make. The
//     ticket demanded picking a side and registering the reason instead of
//     letting the glyphs through unnoticed, so the character ranges below are
//     copied from the scanner's own emojiRe and enforced here. D23's icon rule
//     (lucide-react SVGs, never emoji) is what this pins.
//
//  4. NO APPROVAL DECISION IN THE PANEL (AC#4's ban #6 half). d22scan's ban #6
//     already scans frontend/ text files; this is the same rule asserted from
//     the inside, so a violation is caught by two independent instruments and
//     the naming convention the panel uses ('requestApprovalResolution',
//     'panel.approval.request') is written down where a reader can see why the
//     callback is not named after a decision.

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// bannedStorageAPIs are the persistence surfaces PLAN.md:1044 forbids, written
// as the calls a code writer would actually type.
var bannedStorageAPIs = map[string]*regexp.Regexp{
	"localStorage":   regexp.MustCompile(`\blocalStorage\b`),
	"sessionStorage": regexp.MustCompile(`\bsessionStorage\b`),
	"indexedDB":      regexp.MustCompile(`\bindexedDB\b| indexedDB\.open|\bIDBDatabase\b`),
	"cookie":         regexp.MustCompile(`\bdocument\.cookie\b`),
	"caches":         regexp.MustCompile(`\bcaches\.(open|match|keys)\b`),
	"serviceWorker":  regexp.MustCompile(`\bnavigator\.serviceWorker\b`),
	"fs-access":      regexp.MustCompile(`\brequire\(["']fs["']\)|from ["']node:fs["']`),
}

var (
	// colourLitRe matches the two spellings a colour can be written in, in TS
	// or CSS. It deliberately does not match var(--x).
	colourLitRe  = regexp.MustCompile(`#[0-9A-Fa-f]{3,8}\b|\brgba?\(`)
	importFromRe = regexp.MustCompile(`(?:from|import)\s*\(?\s*["']([^"']+)["']`)
	cssImportRe  = regexp.MustCompile(`@import\s+["']([^"']+)["']`)
	// emojiRangesRe is tools/d22scan/main.go's emojiRe, copied verbatim so the
	// two instruments cannot disagree about what a dingbat is.
	emojiRangesRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
	// panelDecisionIdentifierRe is tools/d22scan's ban #6 pattern (approval.decide).
	panelDecisionIdentifierRe = regexp.MustCompile(`approval\.decide`)
)

func frontendSrcFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(filepath.Join(root, "frontend", "src"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".ts") || strings.HasSuffix(p, ".tsx") || strings.HasSuffix(p, ".css") {
			out = append(out, p)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk frontend/src: %v", err)
	}
	sort.Strings(out)
	return out
}

// reachableFrom computes the module closure of frontend/src/main.tsx across
// both JS imports ("@/x", "./y") and CSS @import, resolving through the "@"
// alias the frontend's tsconfig and vite config pin.
func reachableFrom(t *testing.T, root string) map[string]bool {
	t.Helper()
	srcRoot := filepath.Join(root, "frontend", "src")
	seen := map[string]bool{}
	var visit func(path string)
	visit = func(path string) {
		if seen[path] {
			return
		}
		seen[path] = true
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("cannot read %s while walking the import closure: %v", relToRoot(root, path), err)
			return
		}
		text := string(data)
		var specs []string
		for _, m := range importFromRe.FindAllStringSubmatch(text, -1) {
			specs = append(specs, m[1])
		}
		for _, m := range cssImportRe.FindAllStringSubmatch(text, -1) {
			specs = append(specs, m[1])
		}
		for _, spec := range specs {
			for _, cand := range resolveSpec(srcRoot, path, spec) {
				if st, err := os.Stat(cand); err == nil && !st.IsDir() {
					visit(cand)
				}
			}
		}
	}
	visit(filepath.Join(srcRoot, "main.tsx"))
	if len(seen) < 5 {
		t.Fatalf("import closure of main.tsx holds only %d files - the walk broke and the checks below would be vacuous", len(seen))
	}
	return seen
}

// resolveSpec turns an import specifier into the candidate files it could name.
func resolveSpec(srcRoot, from, spec string) []string {
	var base string
	switch {
	case strings.HasPrefix(spec, "@/"):
		base = filepath.Join(srcRoot, filepath.FromSlash(spec[2:]))
	case strings.HasPrefix(spec, "./"), strings.HasPrefix(spec, "../"):
		base = filepath.Join(filepath.Dir(from), filepath.FromSlash(spec))
	default:
		return nil // a bare package specifier: node_modules, not this tree
	}
	cands := []string{base}
	exts := []string{"", ".ts", ".tsx", ".css", ".js", ".jsx"}
	for _, e := range exts[1:] {
		cands = append(cands, base+e)
	}
	if !strings.HasSuffix(base, ".css") {
		cands = append(cands, filepath.Join(base, "index.tsx"), filepath.Join(base, "index.ts"))
	}
	return cands
}

func relToRoot(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(rel)
}

func TestPanelFrontendIsStateless(t *testing.T) {
	root := panelRepoRoot(t)
	files := frontendSrcFiles(t, root)
	if len(files) == 0 {
		t.Fatal("no files under frontend/src - this check would pass by seeing nothing")
	}
	var hits []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", relToRoot(root, f), err)
		}
		for name, re := range bannedStorageAPIs {
			for i, line := range strings.Split(string(data), "\n") {
				if re.MatchString(line) {
					hits = append(hits, filepath.ToSlash(relToRoot(root, f))+":"+strconv.Itoa(i+1)+": "+name)
				}
			}
		}
	}
	if len(hits) > 0 {
		t.Errorf("frontend holds local state, which PLAN.md:1044 forbids (everything must come back from Go):\n  %s",
			strings.Join(hits, "\n  "))
	}
	t.Logf("scanned %d frontend/src files for %d persistence APIs: 0 hits", len(files), len(bannedStorageAPIs))
}

func TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme(t *testing.T) {
	root := panelRepoRoot(t)
	reachable := reachableFrom(t, root)
	var offenders []string
	for path := range reachable {
		rel := relToRoot(root, path)
		if rel == "frontend/src/styles/tokens.generated.css" {
			continue // the generated copy of C21 is where colours are allowed to exist
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if colourLitRe.MatchString(line) {
				offenders = append(offenders, rel+":"+strconv.Itoa(i+1)+": "+strings.TrimSpace(line))
			}
		}
	}
	if len(offenders) > 0 {
		t.Errorf("a second style source appeared in code the bundle actually ships:\n  %s\n"+
			"Add or reference a C21 token in design/assets/tokens.css instead.", strings.Join(offenders, "\n  "))
	}
	t.Logf("%d files reachable from main.tsx carry colour literals only in the generated theme", len(reachable))
}

func TestVendoredDemoComponentsAreNotMounted(t *testing.T) {
	root := panelRepoRoot(t)
	reachable := reachableFrom(t, root)
	dir := filepath.Join(root, "frontend", "src", "components", "ai-native")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", relToRoot(root, dir), err)
	}
	var mounted, unmounted []string
	for _, e := range entries {
		rel := relToRoot(root, filepath.Join(dir, e.Name()))
		if reachable[filepath.Join(dir, e.Name())] {
			mounted = append(mounted, rel)
			continue
		}
		unmounted = append(unmounted, rel)
	}
	if len(mounted) == 0 {
		t.Errorf("not one vendored Beautiful UI component is mounted - the ticket's premise is that this library is the base, so an unreachable copy of every file means the base is not really in use")
	}
	t.Logf("vendored ai-native: %d mounted (%s), %d unmounted upstream demos (%s)",
		len(mounted), strings.Join(basenames(mounted), " "), len(unmounted), strings.Join(basenames(unmounted), " "))
}

func TestFrontendHasNoEmoji(t *testing.T) {
	root := panelRepoRoot(t)
	var files []string
	for _, dir := range []string{"src", "scripts"} {
		err := filepath.WalkDir(filepath.Join(root, "frontend", dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			files = append(files, p)
			return nil
		})
		if err != nil {
			t.Fatalf("walk frontend/%s: %v", dir, err)
		}
	}
	files = append(files, filepath.Join(root, "frontend", "index.html"))
	var hits []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", relToRoot(root, f), err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if m := emojiRangesRe.FindAllString(line, -1); m != nil {
				hits = append(hits, relToRoot(root, f)+":"+strconv.Itoa(i+1)+" "+strings.Join(m, " "))
			}
		}
	}
	if len(hits) > 0 {
		t.Errorf("D23/ban #8 zero-emoji violated in frontend/ (ranges copied from tools/d22scan):\n  %s",
			strings.Join(hits, "\n  "))
	}
	t.Logf("ban #8 self-armed: %d frontend files scanned, 0 emoji-range characters", len(files))
}

func TestFrontendNeverNamesAnApprovalDecision(t *testing.T) {
	root := panelRepoRoot(t)
	var checked int
	var hits []string
	err := filepath.WalkDir(filepath.Join(root, "frontend"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if strings.Contains(relToRoot(root, p), "/node_modules/") || strings.Contains(relToRoot(root, p), "/dist/") {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		checked++
		for i, line := range strings.Split(string(data), "\n") {
			if panelDecisionIdentifierRe.MatchString(line) {
				hits = append(hits, relToRoot(root, p)+":"+strconv.Itoa(i+1))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk frontend: %v", err)
	}
	if checked == 0 {
		t.Fatal("walked 0 frontend files - the ban #6 twin check would be vacuous")
	}
	if len(hits) > 0 {
		t.Errorf("ban #6 (D33/F2): allow decisions are native-side only, but the panel tree contains the forbidden identifier:\n  %s",
			strings.Join(hits, "\n  "))
	}
	t.Logf("ban #6 twin check: %d frontend files (node_modules and dist excluded) carry no approval-decision identifier", checked)
}

func basenames(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		out = append(out, filepath.Base(p))
	}
	sort.Strings(out)
	return out
}
