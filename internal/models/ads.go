package models

// Ticket 357's report face: which installed files carry a named data stream the
// signed manifest never named.
//
// WHAT THIS IS, IN ONE SENTENCE. It is a read-only enumeration of the install
// directory VerifyInstalled just verified, answering one question - "is there a
// second stream on any of these files" - and it decides nothing: it never refuses
// a hand-off, never changes an exit code, and never runs inside the hash face.
// The ticket's ruling (form 丁) is that a silent pass must become a pass that
// says one sentence, because ticket 121's "nothing unnamed" walk only sees
// directory entries: an ADS is invisible to os.ReadDir and fully readable as
// model.onnx:attacker (measured by 357-a1, D1, on real NTFS), so a tree can match
// every hash the manifest pins and still carry bytes nobody signed.
//
// WHY IT TAKES AN ID AND HANDS THE DIRECTORY BACK. Ticket 357's first face is
// exactly the trap this signature closes: the guard verifies <DataDir>/<id> OR the
// local_override directory for that id, while `wisp models verify` prints
// store=<store root>. A caller that re-derived the directory from the printed
// store would enumerate the WRONG tree whenever a local_override is configured, and
// would then report "no second streams" about a directory nobody verified. So the
// id goes in, and the answer carries the directory that was actually walked - the
// same spelling VerifyInstalled uses, from the same resolveInstalledDir.

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// UnnamedStreams enumerates the named data streams standing on the files under
// the install directory of id, returning the directory it walked and the streams
// it found spelled the way the OS does on the command line
// ("model.onnx:attacker:$DATA"), sorted so two runs print the same sentence.
//
// It walks the files rather than asking a directory handle for its streams: a
// directory handle answers with an empty list AND no error (measured by 357-r1),
// which is the one reading this must never produce by accident.
//
// An empty stream list means "walked every installed file and none carries a
// second stream". A failure to ask is returned as an error and must NOT be read as
// an empty answer - the caller reports it in its own words.
func (m *Manager) UnnamedStreams(id string) (dir string, streams []string, err error) {
	dir = m.resolveInstalledDir(id)
	st, statErr := os.Stat(dir)
	if statErr != nil || !st.IsDir() {
		return dir, nil, fmt.Errorf("model dir %s missing", dir)
	}

	walkErr := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Symlinks are skipped rather than followed: ticket 121's walk already
		// refuses a link standing where a file or a directory was expected, so a
		// link reaching this enumeration means the tree was refused before the
		// report ran. Following it here would read streams outside the verified
		// tree and report them as if they were inside it.
		if d.IsDir() || d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		found, serr := listNamedStreams(p)
		if serr != nil {
			return fmt.Errorf("enumerate named streams of %s: %w", p, serr)
		}
		if len(found) == 0 {
			return nil
		}
		rel, rerr := filepath.Rel(dir, p)
		if rerr != nil {
			return rerr
		}
		rel = filepath.ToSlash(rel)
		for _, s := range found {
			// The OS spells a stream ":name:$DATA"; the concatenation below is the
			// same spelling `dir /r` prints, which is the reading an operator can
			// reproduce without this program.
			streams = append(streams, rel+s)
		}
		return nil
	})
	if walkErr != nil {
		return dir, nil, fmt.Errorf("enumerate named streams under %s: %w", dir, walkErr)
	}
	sort.Strings(streams)
	return dir, streams, nil
}
