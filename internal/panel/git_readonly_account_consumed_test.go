package panel

// Ticket 189's process-starting leg must consume C26's rewrite account before it
// hands a tree to `-C` (ticket 102's rule, ledger A380).
//
// The static criterion in internal/risk (TestC26RewriteAccountIsConsumedAtEvery
// SecurityLeg) only sees that the account's name appears somewhere in
// git_readonly_exec.go; once Result.Actable() is used there, the file keeps the
// name forever, so that criterion can no longer tell a code path that reads the
// account from one that merely mentions it. These sub-tests are the behavioural
// half: a spelling expansion moved onto another tree must be refused here, while
// an ordinary root must still resolve to exactly what it always did.

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveGitTreeConsumesRewriteAccount(t *testing.T) {
	t.Run("ordinary_root_is_the_ordinary_road", func(t *testing.T) {
		dir := t.TempDir()
		tree, kind, reason := resolveGitTree(dir)
		if kind != "" {
			t.Fatalf("resolveGitTree(%q) refused an ordinary root: kind=%q reason=%q", dir, kind, reason)
		}
		if strings.TrimSpace(tree) == "" {
			t.Fatalf("resolveGitTree(%q) returned no tree and no refusal", dir)
		}
	})

	t.Run("rewritten_spelling_never_reaches_the_subprocess", func(t *testing.T) {
		elsewhere := t.TempDir()
		const name = "WISP185ACC_GITTREE"
		t.Setenv(name, elsewhere)
		// The construct this leg is handed by a caller that did NOT read the
		// account: the spelling expands onto a tree the caller never named.
		spelled := "%" + name + "%" + string(filepath.Separator) + "sub"

		tree, kind, reason := resolveGitTree(spelled)
		if kind == "" {
			t.Fatalf("resolveGitTree(%q) accepted a rewritten spelling (tree=%q): the expanded tree would reach `-C` as if the caller had named it",
				spelled, tree)
		}
		if kind != GitKindUnreadable {
			t.Errorf("kind = %q for a rewritten spelling, want the unreadable third state %q", kind, GitKindUnreadable)
		}
		if !strings.Contains(reason, "rewritten=true") {
			t.Errorf("reason %q does not name the rewrite account: a refusal has to say which account it read", reason)
		}
	})
}
