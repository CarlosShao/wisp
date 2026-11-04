package panel

// The states of TreeView.kind (ticket 190). An enum, not a boolean, and for the
// same reason GitView's is one (git.go:84-89): "nothing was chosen", "I am not
// allowed to look", "I could not read it" and "it is not a repository" are four
// different facts. A panel that renders any of them as another is exactly the lie
// ticket 92/145 refuse in every other section of this packet, and ticket 190's
// AC#4 forbids it in this one: the refusal shapes must not show up as a silently
// shorter list.
const (
	// TreeKindNotRepo: the root exists and is a directory, but it is not a git
	// work tree, so there is no tracked/untracked reading to report.
	TreeKindNotRepo = "not_repo"
	// TreeKindNoGit: no usable git binary on this machine. A machine problem,
	// stated as one instead of being reported as an empty tree.
	TreeKindNoGit = "no_git"
	// TreeKindTooLarge: kept from the first leg of this ticket. The reader now
	// answers a failed listing with unreadable (a failure is not a size); the
	// value stays declared because kind values travel to the page.
	TreeKindTooLarge = "too_large"
	// TreeKindCounted: the list was read. A counted answer may still carry
	// FilesTruncated and link/crossing accounting in Reason.
	TreeKindCounted = "counted"

	// TreeKindNotChosen is AC#2's answer when no workspace was narrowed (root
	// empty or unset). The whole disk is NOT a fallback: an unset root is an
	// unread dimension, and it says so instead of listing the process's cwd or
	// the drive it happens to sit on.
	TreeKindNotChosen = "not_chosen"
	// TreeKindUnreadable: the root is not there, is not a directory, or could
	// not be read for a reason that is not a permission refusal.
	TreeKindUnreadable = "unreadable"
	// TreeKindPermissionDenied: the filesystem refused this process. Kept apart
	// from unreadable because "you may not look" and "it is not there" call for
	// different answers on screen (AC#4(b)).
	TreeKindPermissionDenied = "permission_denied"
	// TreeKindLinkOutside: a symlink/junction crossing, either the root itself
	// or a component the listing would have to pass through. C26's default-deny
	// (SPEC-06 §4) refuses to follow it, and this state is how the refusal is
	// said out loud instead of being rendered as a short list (AC#4(c)).
	TreeKindLinkOutside = "link_outside"
)

type TreeEntry struct {
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type TreeView struct {
	Kind           string      `json:"kind"`
	Reason         string      `json:"reason"`
	Entries        []TreeEntry `json:"entries"`
	FilesTotal     int         `json:"filesTotal"`
	FilesTruncated bool        `json:"filesTruncated"`
	BytesTotal     int64       `json:"bytesTotal"`
	RecoverPointer string      `json:"recoverPointer"`
}

type TreeLimits struct {
	MaxEntries int   `json:"maxEntries"`
	MaxBytes   int64 `json:"maxBytes"`
	MaxDepth   int   `json:"maxDepth"`
}
