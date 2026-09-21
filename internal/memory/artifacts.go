package memory

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// artifacts\ management (SPEC-02 §4/§6): host-internal artifacts
// (tool-output-<id>.txt etc. from D10/D15). Flat directory, 500MB quota,
// LRU cleanup by file modification time (Windows often has last-ACCESS time
// disabled, so last-write is the honest proxy for "last used").
//
// "Flat" is enforced, not assumed (ticket 79): an entry that is not a bare
// artifact file is measured recursively, counted against the quota, reclaimed by
// the quota job and by a purge, and reported in the log. It used to be
// `continue`d past, which made a whole subtree invisible to all three at once.

// Artifact is one file in the artifacts directory.
type Artifact struct {
	Name      string    `json:"name"`
	SizeBytes int64     `json:"size_bytes"`
	ModTime   time.Time `json:"mod_time"` // wall clock, display only
}

// strayEntry is a non-conforming artifacts-dir entry: something that is not a
// bare artifact file, in practice a subdirectory that got in there. It is NOT an
// Artifact - Artifact.Name is a bare file name and the public delete route only
// accepts bare names (ticket 76's guard), so a nested path could be listed but
// never deleted. What it must never be again is invisible.
type strayEntry struct {
	name   string    // the top-level offending entry, exactly as listed
	bytes  int64     // every regular file under it, added up
	newest time.Time // newest last-write under it: the LRU proxy for "last used"
}

func strayBytes(strays []strayEntry) int64 {
	var total int64
	for _, st := range strays {
		total += st.bytes
	}
	return total
}

// ErrInvalidArtifactName is the sentinel behind every refusal of the artifact
// name guard (ticket 76). It exists so a caller - or a test - can tell "the
// guard said no before the filesystem was touched" apart from "no such file"
// without matching on message text: the second verdict means a path WAS built,
// which is exactly the thing the artifacts API must never do from a caller name.
var ErrInvalidArtifactName = errors.New("memory: invalid artifact name")

// ListArtifacts returns the artifacts directory contents, oldest-written
// first (the eviction order). Only conforming bare-file artifacts are returned;
// non-conforming entries are logged, counted and reclaimed elsewhere.
func (s *Store) ListArtifacts(ctx context.Context) ([]Artifact, error) {
	arts, strays, err := listArtifactsDir(s.artifactsDir)
	if err != nil {
		return nil, err
	}
	if len(strays) > 0 {
		s.logger.Warn("artifacts dir holds non-conforming entries",
			"component", "memory", "strays", len(strays), "hidden_bytes", strayBytes(strays))
	}
	return arts, nil
}

// listArtifactsDir splits the artifacts directory into conforming artifacts and
// non-conforming entries, the latter measured recursively. Every caller in this
// package feeds it s.artifactsDir and nothing else (ticket 76's audit enforces
// that), so no walk here starts anywhere a caller pointed at.
func listArtifactsDir(dir string) ([]Artifact, []strayEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("memory: read artifacts dir: %w", err)
	}
	out := make([]Artifact, 0, len(entries))
	var strays []strayEntry
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue // raced with a concurrent delete
			}
			return nil, nil, fmt.Errorf("memory: artifact info %s: %w", e.Name(), err)
		}
		if !e.IsDir() {
			out = append(out, Artifact{
				Name:      e.Name(),
				SizeBytes: info.Size(),
				ModTime:   info.ModTime().UTC(),
			})
			continue
		}
		// A subdirectory. WalkDir is Lstat-based and never follows a symlink, so
		// a planted link cannot lead the measurement out of the tree; an entry it
		// cannot read is an error, not a skip (that is the whole point).
		st := strayEntry{name: e.Name(), newest: info.ModTime().UTC()}
		walkErr := filepath.WalkDir(filepath.Join(dir, e.Name()),
			func(_ string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				fi, iErr := d.Info()
				if iErr != nil {
					if errors.Is(iErr, fs.ErrNotExist) {
						return nil // raced with a concurrent delete
					}
					return fmt.Errorf("memory: artifact info %s: %w", d.Name(), iErr)
				}
				st.bytes += fi.Size()
				if m := fi.ModTime().UTC(); m.After(st.newest) {
					st.newest = m
				}
				return nil
			})
		if walkErr != nil {
			return nil, nil, fmt.Errorf("memory: measure artifacts stray %q: %w", e.Name(), walkErr)
		}
		strays = append(strays, st)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime.Before(out[j].ModTime) })
	return out, strays, nil
}

// artifactsDirSize sums the file sizes in the artifacts directory.
func artifactsDirSize(arts []Artifact) int64 {
	var total int64
	for _, a := range arts {
		total += a.SizeBytes
	}
	return total
}

// removeStray reclaims a non-conforming entry: every file under it, then its
// directories bottom-up. Every path is rebuilt as filepath.Join(s.artifactsDir,
// rel) from this store's own walk of its own subtree, and entries are removed
// with os.Remove - never os.RemoveAll - because os.Remove does not recurse: a
// symlink standing where a child was listed gets unlinked, not followed.
func (s *Store) removeStray(name string) error {
	var dirs []string // artifacts-relative, walk (pre-)order; reversed = children first
	err := filepath.WalkDir(filepath.Join(s.artifactsDir, name),
		func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, rErr := filepath.Rel(s.artifactsDir, p)
			if rErr != nil {
				return rErr
			}
			if d.IsDir() {
				if rel != "." {
					dirs = append(dirs, rel)
				}
				return nil
			}
			full := filepath.Join(s.artifactsDir, rel)
			if rErr := os.Remove(full); rErr != nil && !errors.Is(rErr, fs.ErrNotExist) {
				return fmt.Errorf("memory: remove artifacts stray file %q: %w", rel, rErr)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("memory: reclaim artifacts stray %q: %w", name, err)
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		full := filepath.Join(s.artifactsDir, dirs[i])
		if err := os.Remove(full); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("memory: reclaim artifacts stray dir %q: %w", dirs[i], err)
		}
	}
	return nil
}

// enforceArtifactsQuota applies the LRU quota: while the directory total
// exceeds quota bytes, oldest-written entries are deleted. Returns how many
// entries were removed and how many bytes were freed - an entry is one
// conforming file, or one non-conforming subtree (whose bytes were all counted
// against the cap before any of it was reclaimed). Quota of exactly the current
// size is compliant (strictly over triggers cleanup).
func (s *Store) enforceArtifactsQuota(ctx context.Context, quota int64) (files int, freed int64, err error) {
	arts, strays, err := listArtifactsDir(s.artifactsDir)
	if err != nil {
		return 0, 0, err
	}
	// Strays are in the total: that addition is the fix. A subtree the cap could
	// not see could hold any number of bytes.
	total := artifactsDirSize(arts) + strayBytes(strays)
	if total <= quota {
		return 0, 0, nil
	}
	type reclaim struct {
		name    string
		bytes   int64
		at      time.Time
		subtree bool
	}
	items := make([]reclaim, 0, len(arts)+len(strays))
	for _, a := range arts {
		items = append(items, reclaim{name: a.Name, bytes: a.SizeBytes, at: a.ModTime})
	}
	for _, st := range strays {
		items = append(items, reclaim{name: st.name, bytes: st.bytes, at: st.newest, subtree: true})
	}
	// One LRU queue across both kinds, so a stray dir cannot buy itself immunity
	// by being a directory.
	sort.Slice(items, func(i, j int) bool { return items[i].at.Before(items[j].at) })
	over := total - quota
	for _, it := range items {
		if over <= 0 {
			break
		}
		if it.subtree {
			if rmErr := s.removeStray(it.name); rmErr != nil {
				return files, freed, rmErr
			}
		} else if rmErr := os.Remove(filepath.Join(s.artifactsDir, it.name)); rmErr != nil &&
			!errors.Is(rmErr, fs.ErrNotExist) {
			return files, freed, fmt.Errorf("memory: LRU remove %s: %w", it.name, rmErr)
		}
		files++
		freed += it.bytes
		over -= it.bytes
	}
	s.logger.Info("artifacts LRU cleanup",
		"component", "memory",
		"files_removed", files, "bytes_freed", freed, "quota_bytes", quota,
		"stray_entries", len(strays))
	return files, freed, nil
}

// DeleteArtifact removes one artifact file by name. Path traversal is
// rejected: only bare file names inside the artifacts directory are valid.
func (s *Store) DeleteArtifact(ctx context.Context, name string) error {
	if err := validArtifactName(name); err != nil {
		return err
	}
	full := filepath.Join(s.artifactsDir, name)
	if err := os.Remove(full); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("artifact %q: %w", name, ErrNotFound)
		}
		return fmt.Errorf("memory: delete artifact %q: %w", name, err)
	}
	return nil
}

// PurgeArtifacts removes every artifact file and every non-conforming entry
// (privacy: one-click clear, SPEC-02 §5). Ticket 79 is the second half of that
// sentence: a stray subdirectory used to survive a purge, which for a privacy
// clear means it kept bytes the user had just asked to clear. Returns the number
// of entries removed - one per top-level thing reclaimed, a stray subtree
// counting once with all of its bytes reported in the log.
func (s *Store) PurgeArtifacts(ctx context.Context) (int, error) {
	arts, strays, err := listArtifactsDir(s.artifactsDir)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, a := range arts {
		if err := os.Remove(filepath.Join(s.artifactsDir, a.Name)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return removed, fmt.Errorf("memory: purge artifact %q: %w", a.Name, err)
		}
		removed++
	}
	var purgedStrayBytes int64
	for _, st := range strays {
		if err := s.removeStray(st.name); err != nil {
			return removed, err
		}
		purgedStrayBytes += st.bytes
		removed++
	}
	if removed > 0 {
		s.logger.Info("artifacts purged",
			"component", "memory", "files_removed", removed,
			"stray_entries", len(strays), "stray_bytes", purgedStrayBytes)
	}
	return removed, nil
}

// validArtifactName guards against path traversal: artifacts live flat in
// the artifacts dir, so only bare names are accepted. Every refusal carries
// ErrInvalidArtifactName (ticket 76) so the reject is machine-distinguishable
// from a filesystem error.
func validArtifactName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: required (bare file names only)", ErrInvalidArtifactName)
	}
	if strings.ContainsRune(name, '/') || strings.ContainsRune(name, '\\') ||
		strings.ContainsRune(name, ':') || name == "." || name == ".." ||
		filepath.Base(name) != name {
		return fmt.Errorf("%w %q (bare file names only)", ErrInvalidArtifactName, name)
	}
	// Windows strips TRAILING dots and spaces from a path component before it
	// resolves it, so "...." is a spelling of "." and "report.txt." is a spelling
	// of "report.txt": the guard would say "bare name, fine" and os.Remove would
	// then act on a path the caller did not name (measured: DeleteArtifact("....")
	// reached os.Remove on <artifactsDir>\.... and only failed because the
	// directory was not empty). Ticket 76 closes that alias: a name must survive
	// Win32 trailing-dot/space stripping unchanged to be its own name.
	if trimmed := strings.TrimRight(name, ". "); trimmed != name {
		return fmt.Errorf("%w %q (a trailing dot or space is a different path on Windows)",
			ErrInvalidArtifactName, name)
	}
	return nil
}
