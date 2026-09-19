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

// Artifact is one file in the artifacts directory.
type Artifact struct {
	Name      string    `json:"name"`
	SizeBytes int64     `json:"size_bytes"`
	ModTime   time.Time `json:"mod_time"` // wall clock, display only
}

// ListArtifacts returns the artifacts directory contents, oldest-written
// first (the eviction order).
func (s *Store) ListArtifacts(ctx context.Context) ([]Artifact, error) {
	return listArtifactsDir(s.artifactsDir)
}

func listArtifactsDir(dir string) ([]Artifact, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("memory: read artifacts dir: %w", err)
	}
	out := make([]Artifact, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue // raced with a concurrent delete
			}
			return nil, fmt.Errorf("memory: artifact info %s: %w", e.Name(), err)
		}
		out = append(out, Artifact{
			Name:      e.Name(),
			SizeBytes: info.Size(),
			ModTime:   info.ModTime().UTC(),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime.Before(out[j].ModTime) })
	return out, nil
}

// artifactsDirSize sums the file sizes in the artifacts directory.
func artifactsDirSize(arts []Artifact) int64 {
	var total int64
	for _, a := range arts {
		total += a.SizeBytes
	}
	return total
}

// enforceArtifactsQuota applies the LRU quota: while the directory total
// exceeds quota bytes, oldest-written files are deleted. Returns how many
// files were removed and how many bytes were freed. Quota of exactly the
// current size is compliant (strictly over triggers cleanup).
func (s *Store) enforceArtifactsQuota(ctx context.Context, quota int64) (files int, freed int64, err error) {
	arts, err := listArtifactsDir(s.artifactsDir)
	if err != nil {
		return 0, 0, err
	}
	total := artifactsDirSize(arts)
	if total <= quota {
		return 0, 0, nil
	}
	over := total - quota
	for _, a := range arts {
		if over <= 0 {
			break
		}
		if rmErr := os.Remove(filepath.Join(s.artifactsDir, a.Name)); rmErr != nil && !errors.Is(rmErr, fs.ErrNotExist) {
			return files, freed, fmt.Errorf("memory: LRU remove %s: %w", a.Name, rmErr)
		}
		files++
		freed += a.SizeBytes
		over -= a.SizeBytes
	}
	s.logger.Info("artifacts LRU cleanup",
		"component", "memory",
		"files_removed", files, "bytes_freed", freed, "quota_bytes", quota)
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

// PurgeArtifacts removes every artifact file (privacy: one-click clear,
// SPEC-02 §5). Returns the number of files removed.
func (s *Store) PurgeArtifacts(ctx context.Context) (int, error) {
	arts, err := listArtifactsDir(s.artifactsDir)
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
	if removed > 0 {
		s.logger.Info("artifacts purged",
			"component", "memory", "files_removed", removed)
	}
	return removed, nil
}

// validArtifactName guards against path traversal: artifacts live flat in
// the artifacts dir, so only bare names are accepted.
func validArtifactName(name string) error {
	if name == "" {
		return fmt.Errorf("memory: artifact name is required")
	}
	if strings.ContainsRune(name, '/') || strings.ContainsRune(name, '\\') ||
		strings.ContainsRune(name, ':') || name == "." || name == ".." ||
		filepath.Base(name) != name {
		return fmt.Errorf("memory: invalid artifact name %q (bare file names only)", name)
	}
	return nil
}
