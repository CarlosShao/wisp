package secret

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"
	"time"
)

// BlobInfo is the metadata of one stored blob: enough to list what exists
// without ever reading (let alone decrypting) its content. This is the only
// shape `wisp secret list` may print - a listing that carried secret material
// would break C28/SPEC-06.
type BlobInfo struct {
	// ID is the blob id (= the file name under secrets\).
	ID string
	// Ref is the full reference ("dpapi:<ID>"), the form a config field takes.
	Ref string
	// Created is the filesystem creation time where the platform exposes one
	// (Windows: FILETIME), else the modification time.
	Created time.Time
	// Modified is the filesystem modification time.
	Modified time.Time
	// Size is the encrypted blob size in bytes (never the plaintext size).
	Size int64
}

// Blobs enumerates the stored blobs by metadata only: the directory is read,
// the files are not. Entries whose name is not a valid blob id are skipped
// (they are not refs; e.g. a stray file left by a tool), and the result is
// sorted by id so listings are stable.
func (s *Store) Blobs() ([]BlobInfo, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil // no secrets directory yet = no blobs
		}
		return nil, fmt.Errorf("secret: list %s: %w", s.dir, err)
	}
	var out []BlobInfo
	for _, e := range entries {
		if e.IsDir() || !ValidBlobID(e.Name()) {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue // raced with a delete
			}
			return nil, fmt.Errorf("secret: list %s: stat %s: %w", s.dir, e.Name(), err)
		}
		out = append(out, BlobInfo{
			ID:       e.Name(),
			Ref:      RefPrefixDPAPI + e.Name(),
			Created:  createdTime(fi),
			Modified: fi.ModTime(),
			Size:     fi.Size(),
		})
	}
	slices.SortFunc(out, func(a, b BlobInfo) int { return strings.Compare(a.ID, b.ID) })
	return out, nil
}
