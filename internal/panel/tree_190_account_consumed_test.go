package panel

// Ticket 190's file list must consume C26's rewrite account, the way the git
// dimension already does (git_test.go:TestReadGitForWorkspaceConsumesRewriteAccount).
//
// The static criterion in internal/risk's pathresolver_rewrite_account_test.go
// (TestC26RewriteAccountIsConsumedAtEverySecurityLeg, ledger A380) can only see
// that SOME name from the account appears in a file. It cannot see whether the
// account is READ or merely mentioned. These sub-tests are the behavioural half
// for ticket 190's leg: a rewritten coordinate must reach the page with the
// account in its sentence, and an ordinary one must be byte-for-byte what it was
// before the account-aware entry existed.

import (
	"strings"
	"testing"
)

func TestReadFileTreeForWorkspaceConsumesRewriteAccount(t *testing.T) {
	dir := t.TempDir()

	t.Run("not_rewritten_is_the_ordinary_road", func(t *testing.T) {
		ws := WorkspaceView{Set: true, Spelling: dir, Canonical: dir}
		got := ReadFileTreeForWorkspace190(ws)
		want := ReadFileTree190(dir)
		if got.Kind != want.Kind || got.Reason != want.Reason || got.FilesTotal != want.FilesTotal {
			t.Errorf("account-aware entry disagreed with the plain read: got kind=%q reason=%q total=%d, want kind=%q reason=%q total=%d",
				got.Kind, got.Reason, got.FilesTotal, want.Kind, want.Reason, want.FilesTotal)
		}
		if strings.Contains(got.Reason, "rewritten=true") {
			t.Errorf("a view with rewritten=false was narrated as rewritten: %q", got.Reason)
		}
	})

	t.Run("rewritten_must_say_so_out_loud", func(t *testing.T) {
		ws := WorkspaceView{Set: true, Spelling: `%WISP190ACC%\named`, Canonical: dir, Rewritten: true}
		got := ReadFileTreeForWorkspace190(ws)
		if !strings.Contains(got.Reason, "rewritten=true") {
			t.Errorf("reason %q does not name the rewrite account: a listing of the expanded tree would read as the operator's own folder", got.Reason)
		}
		if !strings.Contains(got.Reason, ws.Spelling) || !strings.Contains(got.Reason, ws.Canonical) {
			t.Errorf("reason %q must carry both the spelling the user named and the tree C26 moved it to", got.Reason)
		}
		// The disposition is the git dimension's (candidate 1 of ticket 181 AC#6):
		// the expanded tree is still read, and it is the SAME reading the plain
		// door gives for that tree - the account changes the sentence, not the data.
		want := ReadFileTree190(dir)
		if got.Kind != want.Kind || got.FilesTotal != want.FilesTotal {
			t.Errorf("rewritten view read %q (total %d), want the same dimension the plain door reads off the expanded tree: %q (total %d)",
				got.Kind, got.FilesTotal, want.Kind, want.FilesTotal)
		}
		if !strings.HasSuffix(got.Reason, want.Reason) {
			t.Errorf("reason %q does not keep the reader's own sentence %q after the account narration", got.Reason, want.Reason)
		}
	})

	t.Run("unset_view_reads_nothing", func(t *testing.T) {
		if got := ReadFileTreeForWorkspace190(UnsetWorkspaceView()); got.Kind != TreeKindNotChosen {
			t.Errorf("kind = %q for an unset workspace, want %q", got.Kind, TreeKindNotChosen)
		}
	})
}
