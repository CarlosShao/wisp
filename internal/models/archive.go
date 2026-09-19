package models

import (
	"archive/tar"
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractTarBz2 extracts a sha256-verified .tar.bz2 archive into destDir.
//
// Security posture (C26 spirit applied to archives): only regular files and
// directories are materialized; any link entry rejects the whole archive;
// every member path must be a clean relative path inside destDir; a single
// shared top-level directory (the k2-fsa release convention) is stripped so
// the flattened layout matches ArchiveSpec.Files paths.
func ExtractTarBz2(archivePath, destDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return observeWrap(ClassModel, "open archive", err)
	}
	defer f.Close()

	bz := bzip2.NewReader(f)
	tr := tar.NewReader(bz)

	var topDirs []string
	type entry struct {
		header *tar.Header
		body   []byte
	}
	// Two passes are avoided by buffering small headers first to detect the
	// shared top dir; instead we extract into a temp subdir with structure,
	// then flatten if exactly one top dir exists. Bounded by the staging dir
	// size on disk, not memory: entries are streamed to the temp tree.
	tmp := destDir + ".extracting"
	_ = os.RemoveAll(tmp)
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return observeWrap(ClassModel, "create extraction dir", err)
	}
	defer os.RemoveAll(tmp)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return observeWrap(ClassModel, "read archive", err)
		}
		if hdr.Typeflag != tar.TypeReg && hdr.Typeflag != tar.TypeDir {
			return observeNew(ClassModel, fmt.Sprintf("archive %s: non-file member %q (type %q) - refusing", filepath.Base(archivePath), hdr.Name, string(hdr.Typeflag)))
		}
		name := filepath.ToSlash(hdr.Name)
		// Strip "./" prefixes some tar writers add.
		for strings.HasPrefix(name, "./") {
			name = strings.TrimPrefix(name, "./")
		}
		if name == "" {
			continue
		}
		if err := validRelPath(strings.TrimSuffix(name, "/")); err != nil {
			return observeNew(ClassModel, fmt.Sprintf("archive %s: unsafe member path %q", filepath.Base(archivePath), hdr.Name))
		}
		seg := strings.SplitN(strings.TrimSuffix(name, "/"), "/", 2)
		topDirs = append(topDirs, seg[0])

		target := filepath.Join(tmp, filepath.FromSlash(name))
		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return observeWrap(ClassModel, "archive dir", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return observeWrap(ClassModel, "archive parent", err)
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, normalMode(hdr.FileInfo().Mode()))
		if err != nil {
			return observeWrap(ClassModel, "archive file create", err)
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return observeWrap(ClassModel, "archive file write", err)
		}
		if err := out.Close(); err != nil {
			return observeWrap(ClassModel, "archive file close", err)
		}
	}

	// Flatten: if every member lived under one common top directory, move its
	// children up. Multiple top entries (or none) install as-is.
	unique := uniqueStrings(topDirs)
	if len(unique) == 1 {
		inner := filepath.Join(tmp, unique[0])
		if entries, err := os.ReadDir(inner); err == nil {
			for _, e := range entries {
				if err := os.Rename(filepath.Join(inner, e.Name()), filepath.Join(destDir, e.Name())); err != nil {
					return observeWrap(ClassModel, "flatten extraction", err)
				}
			}
		} else {
			return observeWrap(ClassModel, "read extraction dir", err)
		}
	} else {
		entries, err := os.ReadDir(tmp)
		if err != nil {
			return observeWrap(ClassModel, "read extraction dir", err)
		}
		for _, e := range entries {
			if err := os.Rename(filepath.Join(tmp, e.Name()), filepath.Join(destDir, e.Name())); err != nil {
				return observeWrap(ClassModel, "flatten extraction", err)
			}
		}
	}
	return nil
}

func normalMode(m os.FileMode) os.FileMode {
	if m&0o111 != 0 {
		return 0o755
	}
	return 0o644
}

func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
