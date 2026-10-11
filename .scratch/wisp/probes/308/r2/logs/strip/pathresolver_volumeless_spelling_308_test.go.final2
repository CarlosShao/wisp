package risk

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// driveShare takes a path's volume name and returns the admin-share token
// ("c$") the UNC spellings need -- WITHOUT the `vol[:1]` slice that panicked
// ticket 308. On a volume-less temp filepath.VolumeName returns "", and the
// old inline `strings.ToLower(vol[:1])` then panicked
// ("slice bounds out of range [:1] with length 0"), which did not fail one
// case: it killed the whole test binary so every OTHER case in internal/risk
// got zero evaluation (the shape ticket 308 is about -- looks like red, is
// really "this run never finished").
//
// Within this test surface this is the single place the drive letter is taken.
// The repo takes it in one other place, guarded and outside this ticket:
// production internal/risk/pathresolver.go (normalizeLocalUNC) checks
// len(drive) == 2 && drive[1] == '$' before its drive[:1], and
// internal/tools/volume_windows.go reads filepath.VolumeName (volumeID,
// volumeRootOf) without slicing it at all. Census 2026-10-11, scope
// `internal cmd tools` (ruler: git grep -e 'vol[:1]' -e 'drive[:1]' -- '*.go'):
// drive-letter slice evaluation sites = 2 -- that production line and this
// helper -- and every other line the pattern hits is a comment. 编排者裁语 ⓑ/ⓒ: the
// guard's teeth are on the INPUT layer (a caller asks driveShare and reacts to
// ok), NOT on a machine that happens to hand out a volume-less temp -- this
// ticket may not depend on being able to build such a surface (ⓑ: 本机造不出).
func driveShare(vol string) (share string, ok bool) {
	if len(vol) == 0 {
		return "", false
	}
	return strings.ToLower(vol[:1]) + "$", true
}

// volumeInvariantSpellings is the spelling-selection half of
// TestClassifySpellingInvariance (anchor_spelling_windows_test.go), extracted
// so the volume-less shape is unit-testable by feeding it a literal. The seven
// volume-free spellings are ALWAYS present; the two volume-derived UNC spellings
// (unc-admin, unc-extended) are dropped -- returned by name in `dropped` --
// ONLY when driveShare reports there is no drive letter to build them from.
// Dropping those two named items is exactly what ticket AC#1 sanctions
// (⛔ 去掉整枚, ⛔ t.Skip); the rest keep being measured.
func volumeInvariantSpellings(home, long string) (keep map[string]string, dropped []string) {
	keep = map[string]string{
		"long":           long,
		"upper":          strings.ToUpper(long),
		"lower":          strings.ToLower(long),
		"extended":       `\\?\` + long,
		"slashes":        strings.ReplaceAll(long, `\`, "/"),
		"dot-and-dotdot": filepath.Join(home, ".", ".ssh", "..", ".ssh", "id_testkey"),
		"not-yet-there":  filepath.Join(home, ".ssh", "brand-new-key"),
	}
	vol := filepath.VolumeName(long)
	share, ok := driveShare(vol)
	if !ok {
		return keep, []string{"unc-admin", "unc-extended"}
	}
	rest := strings.TrimPrefix(long, vol)
	keep["unc-admin"] = `\\localhost\` + share + rest
	keep["unc-extended"] = `\\?\UNC\localhost\` + share + rest
	return keep, nil
}

// TestVolumeInvariantSpellingsVolumeLess308 is ticket 308 AC#1's direct proof.
// It feeds volumeInvariantSpellings / driveShare a volume-less LITERAL
// ("/tmp/...") as INPUT. This literal is the helper's own input; per the ticket
// hard constraint it never replaces an existing case's t.TempDir() (票 72 :201
// 那条禁线同形). On every OS filepath.VolumeName("/tmp/...") == "", the exact
// shape that used to panic. Assertions, in AC#1 order:
//
//	(1) no panic -- this case reaches its own PASS/FAIL verdict instead of
//	    taking the whole binary down;
//	(2) the two unc-* spellings are dropped BY NAME (driveShare reports ok=false);
//	(3) long/upper/lower/extended -- the four spellings AC#1 names, plus the
//	    other volume-free ones -- are still present and will still be evaluated.
func TestVolumeInvariantSpellingsVolumeLess308(t *testing.T) {
	home := "/tmp/wisp-308-no-volume"
	long := filepath.Join(home, ".ssh", "id_testkey")

	// Premise: the literal really is the volume-less shape (else this proves
	// nothing about the guard -- fail the premise rather than report a green
	// that measured the wrong thing).
	if vol := filepath.VolumeName(long); vol != "" {
		t.Fatalf("premise broken: filepath.VolumeName(%q)=%q, want the volume-less shape \"\"", long, vol)
	}

	// (1)+(2): driveShare must return, not panic, and report the missing volume.
	share, ok := driveShare(filepath.VolumeName(long))
	if ok {
		t.Fatalf("driveShare(\"\") returned ok=true share=%q: a volume-less path must report not-ok, never fabricate a share", share)
	}
	if share != "" {
		t.Errorf("driveShare(\"\") share=%q, want empty string", share)
	}

	// (2)+(3): the selection drops exactly the two UNC spellings by name and
	// keeps every volume-free spelling.
	keep, dropped := volumeInvariantSpellings(home, long)
	sort.Strings(dropped)
	if len(dropped) != 2 || dropped[0] != "unc-admin" || dropped[1] != "unc-extended" {
		t.Fatalf("dropped=%v, want exactly [unc-admin unc-extended]", dropped)
	}
	for _, name := range []string{"unc-admin", "unc-extended"} {
		if _, present := keep[name]; present {
			t.Errorf("keep still holds %q with no volume: the guard failed to drop it", name)
		}
	}
	volumeFree := []string{"long", "upper", "lower", "extended", "slashes", "dot-and-dotdot", "not-yet-there"}
	for _, name := range volumeFree {
		if _, present := keep[name]; !present {
			t.Errorf("volume-free spelling %q went missing: only the two volume-derived UNC spellings may be dropped (AC#1)", name)
		}
	}

	// Keep the denominator visible for the next reader (ticket: 跑了 7 项、去掉 2
	// 项, ⛔ 写成"全跑"): the volume-qualified surface evaluates 9 spellings,
	// this volume-less input still evaluates 7 of them.
	got := make([]string, 0, len(keep))
	for name := range keep {
		got = append(got, name)
	}
	sort.Strings(got)
	t.Logf("ticket 308 guard on volume-less input %q: %d volume-free spellings still evaluated %v; %d UNC spellings dropped by name %v",
		long, len(keep), got, len(dropped), dropped)
}
