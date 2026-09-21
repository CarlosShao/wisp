package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// A18 / ticket 73: reclaiming staging files a real kill left behind
// ---------------------------------------------------------------------------
//
// D31's writer stages beside the target and touches the target exactly once
// (os.Rename). That is why a kill mid-write leaves the destination byte-exact -
// and it is also why the kill leaves the STAGING FILE behind: `taskkill /F`
// runs no Go cleanup, so nothing removes it. 761447f measured it with a real
// child process and a real taskkill: exactly one `.wisp-tmp-*` per kill, and
// neither the next successful write nor a freshly built bridge looks at it
// (docs/evidence/s1/20-a18-taskkill-residue-characterization.md §2).
//
// This file is the sweeper. Two decisions are written down because they are
// behavior, not mechanics:
//
//   - WHEN: on the NEXT WRITE into the directory (called from stageAndRename
//     and from fs.move's crossVolume half), not at bridge start-up. There is no
//     start-up seam in internal/tools worth hooking - New() is a constructor,
//     and cmd/wisp is not this ticket's to touch - whereas "the next write into
//     that directory" is exactly the moment A18's characterization already
//     reaches for, so the fix and the proof sit on the same seam.
//   - WHOSE: a staging file is attributable to the program that created it by
//     NAME, because attribution has to survive the death of the process that
//     wrote it (an in-process ledger can never reclaim what a killed session
//     left). The name is
//
//	.wisp-tmp-<owner>-<pid>-<random>
//
//     with <owner> a stable 8-hex token derived from (executable path, user
//     name) and <pid> the creator's process id. Stable across runs, so a new
//     process can reclaim a dead one's leftovers; different between programs
//     and between users, so a foreign file that merely shares the
//     `.wisp-tmp-` prefix is never touched. <pid> is what keeps a CONCURRENT
//     live writer safe: an in-flight staging file belongs to a process that is
//     still running, so it is not an orphan and the sweep leaves it alone
//     (measured, not assumed: stagingCreatorAlive is the guard for it).
//
// The trade-off is stated on purpose: an orphan whose <owner> no longer matches
// (Wisp was reinstalled to a different path, or runs under another user) becomes
// permanent litter. That is the conservative side of the mistake - this ticket
// forbids guessing at deletion, because a sweeper that deletes what it cannot
// attribute is just a delete primitive pointed at whatever the model typed.
//
// What the sweep NEVER does: follow a reparse point (that is how a cleanup
// becomes a delete outside the allowed dirs - C26's whole objection, see
// bridge_junction_windows_test.go), touch a non-regular entry, delete a file it
// cannot attribute, or fail the write it was invited by. A file another process
// holds open is retried a couple of times and then left in place.

const (
	// stagingOwnerLen is how much of the owner digest goes into a name. 4 bytes
	// is 1.8e9 values: collision with some other program's hand-made file name
	// is not a realistic risk, and a shorter token would be.
	stagingOwnerLen = 8
	// stagingSeparator splits the three attributable fields.
	stagingSeparator = "-"
	// stagingScanCap bounds how many directory entries one sweep examines, so a
	// pathological directory cannot stall a write.
	stagingScanCap = 4096
	// stagingRetries is how many extra times an unlink is attempted before the
	// file is judged busy and left alone. The first attempt is not counted.
	stagingRetries = 2
	// stagingRetryWait is the pause between attempts - long enough for a
	// short-lived handle to close, short enough not to be a write's cost.
	stagingRetryWait = 15 * time.Millisecond
)

// stagingOwner identifies the files THIS program creates for THIS user. It is
// computed once; a failure to read either input degrades to a fixed label
// rather than to "unknown", so the scheme never collapses into "match anything".
var stagingOwner = sync.OnceValue(func() string {
	exe, err := os.Executable()
	if err != nil || exe == "" {
		exe = "<unreadable-executable>"
	}
	user := stagingUserName()
	sum := sha256.Sum256([]byte(strings.Join([]string{exe, user, tempPrefix}, "\x00")))
	return hex.EncodeToString(sum[:])[:stagingOwnerLen]
})

// stagingUserName reads the account the sweeper runs as, degrading to the
// environment and finally to a fixed label. The label is never empty: an empty
// field would make two different programs hash alike.
func stagingUserName() string {
	if u, err := user.Current(); err == nil && u != nil && u.Username != "" {
		return u.Username
	}
	for _, k := range []string{"USERNAME", "USER", "LNAME", "LOGNAME"} {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return "<unreadable-user>"
}

// stagingPattern is the os.CreateTemp pattern for a Wisp staging file. The
// trailing '*' is where CreateTemp puts its random token, and it must stay the
// last byte of the pattern: the matcher below assumes exactly three fields.
func stagingPattern() string {
	return tempPrefix + stagingOwner() + stagingSeparator +
		strconv.Itoa(os.Getpid()) + stagingSeparator + "*"
}

// stagingNameOf renders one concrete staging name without creating it. Tests
// plant orphans with it (a killed process is not here to name its own files),
// so the naming scheme has exactly ONE spelling in this package.
func stagingNameOf(owner string, pid int, random string) string {
	return tempPrefix + owner + stagingSeparator + strconv.Itoa(pid) +
		stagingSeparator + random
}

// stagingAttributable reports whether one directory entry name is a Wisp
// staging file THIS program created. It is a pure string test on the NAME (no
// stat, no open), and it errs toward "not mine".
func stagingAttributable(name string) bool {
	return stagingAttribution(name) != ""
}

// stagingAttribution parses a staging name and returns the creator's pid as a
// string, or "" when the name is not attributable. Splitting is done on fixed
// field counts rather than with a wildcard prefix match, because "starts with
// .wisp-tmp-" is precisely the rule this ticket forbids.
func stagingAttribution(name string) string {
	if !strings.HasPrefix(name, tempPrefix) {
		return ""
	}
	owner, pid, _, ok := splitStagingName(strings.TrimPrefix(name, tempPrefix))
	if !ok || owner != stagingOwner() {
		return ""
	}
	return pid
}

// splitStagingName cuts "<owner>-<pid>-<random>" into its three fields. The
// random field comes from os.CreateTemp, whose alphabet this package does not
// control, so it is accepted as "non-empty and contains no separator and no
// path character" rather than as a specific radix.
func splitStagingName(rest string) (owner, pid, random string, ok bool) {
	parts := strings.Split(rest, stagingSeparator)
	if len(parts) != 3 {
		return "", "", "", false
	}
	owner, pid, random = parts[0], parts[1], parts[2]
	if len(owner) != stagingOwnerLen || random == "" || !allDigits(pid) {
		return "", "", "", false
	}
	if strings.ContainsAny(random, `\/*:`) {
		return "", "", "", false
	}
	return owner, pid, random, true
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > 0
}

// stagingSweep is what one sweep did, in facts the ledger can render.
type stagingSweep struct {
	reclaimed int
	skipped   int
}

func (s stagingSweep) String() string {
	return fmt.Sprintf("清扫 %d 个，跳过 %d 个", s.reclaimed, s.skipped)
}

// reclaimStaging is the write paths' entry into the sweeper, and the only place
// that speaks about it in the D31 ledger. Deletions are a side effect, so they
// are reported as facts rather than happening silently under a write the user
// approved - and "the target file is a few bytes smaller than it should be" is
// not a thing a reader should have to wonder about.
func reclaimStaging(dir string, se *sideEffect) {
	swept := sweepStagingOrphans(dir)
	if swept.reclaimed+swept.skipped == 0 {
		return
	}
	se.record("清扫 %s 里本程序留下的历史暂存文件（上次写盘被中断的残口）：%s；"+
		"跳过的都不是孤儿：创建者进程仍在运行、正被占用、或不是可归因的普通文件",
		dir, swept)
}

// sweepStagingOrphans reclaims the staging files this program left in dir under
// a killed write. It is called BEFORE the new temp file is created, so the
// sweep can never see - let alone delete - the file the current write owns.
//
// It returns counts, not an error: an interrupted-write janitor that can fail
// the write it was invited by has turned a 4-byte leak into a lost write.
func sweepStagingOrphans(dir string) stagingSweep {
	var got stagingSweep
	des, err := os.ReadDir(dir)
	if err != nil {
		// Cannot even list the directory: nothing is deleted. A write into a
		// directory we cannot read will fail on its own terms and say so.
		return got
	}
	for i, e := range des {
		if i >= stagingScanCap {
			break
		}
		pid := stagingAttribution(e.Name())
		if pid == "" {
			continue // not ours, or not a staging file at all
		}
		path := filepath.Join(dir, e.Name())
		if !stagingDeletable(path) {
			got.skipped++
			continue
		}
		if alive, err := stagingCreatorAlive(pid); err != nil || alive {
			// Undecidable counts as alive: an orphan is a file whose creator is
			// demonstrably gone, and a concurrent writer's in-flight staging
			// file must survive this sweep.
			got.skipped++
			continue
		}
		if removeOrReclaim(path) {
			got.reclaimed++
		} else {
			got.skipped++
		}
	}
	return got
}

// stagingDeletable is the on-disk half of attribution: a plain regular file of
// ours, judged with Lstat so a reparse point is seen as ITSELF and not through
// whatever it reaches. Directories, junctions, symlinks and pipes are all left
// alone no matter how well their name matches.
func stagingDeletable(path string) bool {
	st, err := os.Lstat(path)
	if err != nil || !st.Mode().IsRegular() {
		return false
	}
	return !stagingIsReparse(st)
}

// removeOrReclaim unlinks with a short retry budget and reports whether the file
// is gone. A file another process still holds open is LEFT IN PLACE (ticket 73:
// retry-then-skip, never a hard failure), which is also what makes the sweep
// idempotent - the next write tries again.
func removeOrReclaim(path string) bool {
	for attempt := 0; ; attempt++ {
		err := os.Remove(path)
		if err == nil {
			return true
		}
		if attempt >= stagingRetries {
			return false
		}
		time.Sleep(stagingRetryWait)
	}
}
