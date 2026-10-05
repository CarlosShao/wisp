package panel

// Ticket 181 AC#7: the PRODUCER side of the rewrite account.
//
// 181-r2 landed the consumer (ReadGitForWorkspace reads ws.Rewritten) and then
// measured that no producer could ever fill the field: WorkspaceViewFromRoot
// took a bare string and set three of the view's six fields, so the account
// booleans were not "read as false", they were unreachable, and both consumer
// branches in this package (the git dimension here, the project-instruction
// request in instructions_200.go) were decoration. The tests below are the
// ones that break when the producer goes back to dropping the account - which
// is exactly the distinction AC#7 asked for, because "the field exists" and
// "the package compiles" both stay true either way.
//
// accountScope is the scripted side of the seam PathScope already is: ticket 92
// built this interface so a test could hold the handler to scripted verdicts
// without assembling an allowlist (workspace_test.go's scriptedScope is the
// sibling), and internal/tools' real PathCanonicalizer satisfies the same
// interface. It is NOT standing in for the real producer: the real one is
// nailed from its own package in internal/tools/paths_workspace_account_181r3_test.go,
// and the two files together are the evidence that the answer crosses the seam.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// accountScope is a path scope whose answer about the narrowing in force is
// whatever account the test hands it, and nothing more.
type accountScope struct {
	inForce risk.Result
}

func (s *accountScope) WorkspaceRoot() risk.Result { return s.inForce }

func (s *accountScope) ResolveWorkspace(string) (risk.Result, error) {
	return risk.Result{}, risk.ErrRewrittenPath
}

func (s *accountScope) SetWorkspaceRoot(string, risk.Result) error { return nil }

func TestWorkspaceViewCarriesEveryHalfOfTheAccount(t *testing.T) {
	const (
		spelled  = `%WISP181R3ACC%\proj`
		expanded = `C:\elsewhere\proj`
		plain    = `D:\work\Wisp`
	)

	t.Run("clean_answer_is_the_ordinary_road", func(t *testing.T) {
		v := WorkspaceViewFromRoot(risk.Result{
			Spelling: plain, Canonical: plain, Resolved: true,
		})
		if !v.Set || v.Canonical != plain {
			t.Fatalf("view = %+v, want the resolved tree set", v)
		}
		// Spelling is the half that was unreachable by construction before this
		// ticket: composer.go:213 promises "the user's own string, verbatim",
		// and a producer taking a bare string could not promise it at all.
		if v.Spelling != plain {
			t.Errorf("spelling = %q, want the caller's own string %q verbatim", v.Spelling, plain)
		}
		if v.Rewritten || v.Reparse {
			t.Errorf("a clean account rendered as an accounted one: %+v", v)
		}
		if !strings.Contains(v.Reason, "已收窄") || strings.Contains(v.Reason, "改写") {
			t.Errorf("reason %q, want the narrowing sentence and no rewrite claim", v.Reason)
		}
	})

	t.Run("rewritten_answer_must_be_rewritten_in_the_view", func(t *testing.T) {
		// This is AC#7's judgement sub-case. It is not "the field exists": with
		// the pre-181-r3 producer the same input reached this function as a bare
		// string and the view came out Rewritten == false, reason free of any
		// account, which is the lie the ticket names.
		v := WorkspaceViewFromRoot(risk.Result{
			Spelling: spelled, Canonical: expanded, Resolved: true,
			Rewritten: true, Rewrites: []string{"env"},
		})
		if !v.Rewritten {
			t.Errorf("rewritten = false for the account C26 marked rewritten (%+v): "+
				"the producer dropped the book again, and every consumer of this field is decoration", v)
		}
		if v.Spelling != spelled {
			t.Errorf("spelling = %q, want the spelling C26 was given %q", v.Spelling, spelled)
		}
		if v.Canonical != expanded {
			t.Errorf("canonical = %q, want the expanded tree %q", v.Canonical, expanded)
		}
		for _, want := range []string{"改写", "env", spelled, expanded} {
			if !strings.Contains(v.Reason, want) {
				t.Errorf("reason %q does not say %q out loud: a moved coordinate must not read as the operator's own folder", v.Reason, want)
			}
		}
	})

	t.Run("reparse_answer_is_read_not_assumed", func(t *testing.T) {
		v := WorkspaceViewFromRoot(risk.Result{
			Spelling: `D:\work\link`, Canonical: `D:\work\link`, Reparse: true,
		})
		if !v.Reparse {
			t.Errorf("reparse = false although C26 reported an exception-listed traversal: %+v", v)
		}
		if v.Rewritten {
			t.Errorf("a reparse traversal rendered as a rewrite: %+v", v)
		}
		if !strings.Contains(v.Reason, "联接点") {
			t.Errorf("reason %q must name the reparse account it is carrying", v.Reason)
		}
	})

	t.Run("no_answer_still_renders_unset_with_a_reason", func(t *testing.T) {
		// Unset is the one state that must NOT be invented from an account, and
		// the zero Result is what internal/tools answers for it (AC#4 of ticket
		// 92: "no workspace" is a fact the user has to be able to read).
		v := WorkspaceViewFromRoot(risk.Result{})
		if v.Set || v.Reason == "" || v.Rewritten || v.Reparse {
			t.Errorf("unset rendered as %+v", v)
		}
		if v != UnsetWorkspaceView() {
			t.Errorf("unset view = %+v, want the same view UnsetWorkspaceView mints", v)
		}
	})
}

// TestWorkspaceAccountReachesEveryConsumer is the end-to-end half of AC#7: the
// account leaves the path scope as the same answer cmd/wisp's pump line passes
// into WorkspaceViewFromRoot, and both legs that read it react. Before this
// ticket the route existed on paper only - no producer could put a rewritten
// coordinate on it.
func TestWorkspaceAccountReachesEveryConsumer(t *testing.T) {
	tree := mkTree(t, "accounted tree", map[string]string{
		".git/HEAD":                 "ref: refs/heads/release/7\n",
		".git/refs/heads/release/7": gitRefFile,
	})
	spelled := `%WISP181R3ACC%`
	acc := risk.Result{
		Spelling: spelled, Canonical: tree, Resolved: true,
		Rewritten: true, Rewrites: []string{"env"},
	}
	scope := &accountScope{inForce: acc}

	// The one shape the production call site has (cmd/wisp/panel_pump.go:87):
	// the view is built from the scope's answer, byte for byte the same
	// expression, so this is the packet's producer and not a test-local one.
	view := WorkspaceViewFromRoot(scope.WorkspaceRoot())
	if !view.Rewritten {
		t.Fatalf("the view built from a scope that reports a rewritten account is not rewritten: %+v", view)
	}

	t.Run("git_dimension_says_the_account_out_loud", func(t *testing.T) {
		got := ReadGitForWorkspace(view)
		if !strings.Contains(got.Reason, "rewritten=true") {
			t.Errorf("git reason %q does not name the account: the branch it shows would read as if it belonged to the folder the operator named", got.Reason)
		}
		if !strings.Contains(got.Reason, spelled) || !strings.Contains(got.Reason, tree) {
			t.Errorf("git reason %q must carry both the spelling and the expanded coordinate", got.Reason)
		}
		if got.Branch != "release/7" {
			t.Errorf("git branch = %q, want release/7 - the dimension still reports the tree it read, it does not blank it (candidate 2 is a different ruling)", got.Branch)
		}
		if got.CurrentWorktree != view.Canonical {
			t.Errorf("currentWorktree %q and workspace.canonical %q disagree - one packet, two trees",
				got.CurrentWorktree, view.Canonical)
		}
	})

	t.Run("project_instruction_request_refuses_both_trees", func(t *testing.T) {
		req := ProjectInstructionLoadRequestFor(view, `D:\data`)
		if req.Refused == "" || req.WorkspaceDir != "" || req.GlobalDir != "" {
			t.Errorf("a rewritten workspace handed the loader a directory anyway: %+v", req)
		}
	})

	t.Run("the_account_is_on_the_wire", func(t *testing.T) {
		pump := NewSnapshotPump(PumpSources{
			Workspace: func() WorkspaceView { return WorkspaceViewFromRoot(scope.WorkspaceRoot()) },
			Git:       func() GitView { return ReadGitForWorkspace(view) },
		})
		data, err := pump.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var wire struct {
			Composer struct {
				Workspace struct {
					Set       bool   `json:"set"`
					Spelling  string `json:"spelling"`
					Canonical string `json:"canonical"`
					Rewritten bool   `json:"rewritten"`
					Reason    string `json:"reason"`
				} `json:"workspace"`
			} `json:"composer"`
		}
		if err := json.Unmarshal(data, &wire); err != nil {
			t.Fatal(err)
		}
		if !wire.Composer.Workspace.Rewritten {
			t.Errorf("packet workspace = %+v, want rewritten:true - the bytes the panel renders from are the whole judgement here", wire.Composer.Workspace)
		}
		if wire.Composer.Workspace.Spelling != spelled || wire.Composer.Workspace.Canonical != tree {
			t.Errorf("packet workspace coordinates = %+v, want spelling %q and canonical %q",
				wire.Composer.Workspace, spelled, tree)
		}
		if !strings.Contains(wire.Composer.Workspace.Reason, "改写") {
			t.Errorf("packet workspace reason %q lost the account sentence", wire.Composer.Workspace.Reason)
		}
	})
}

// TestScopeThatChangesItsAccountIsSurfaced pins the half of currentAfter that
// only exists once the account travels: a scope whose answer moved in its BOOK
// while its coordinate stood still is a scope that is not sure what it
// authorised, and the handler must say so instead of rendering whichever half
// it likes. With the pre-181-r3 bare-string scope this was structurally
// invisible - a string cannot disagree about an account.
func TestScopeThatChangesItsAccountIsSurfaced(t *testing.T) {
	scope := &movingScope{}
	view, err := RequestWorkspaceSwitch(scope, `D:\work\beta`, nil)
	if err == nil {
		t.Fatal("a scope whose account moved under a standing coordinate reported success")
	}
	if view.Set {
		t.Errorf("the moved scope was rendered as a normal narrowing: %+v", view)
	}
	if !strings.Contains(view.Reason, "账户") || !strings.Contains(view.Reason, "rewritten=true") {
		t.Errorf("reason %q must name the account half of what moved", view.Reason)
	}
}

// movingScope answers the same coordinate twice with two different accounts.
type movingScope struct{ calls int }

func (s *movingScope) WorkspaceRoot() risk.Result {
	s.calls++
	if s.calls == 1 {
		return risk.Result{Spelling: `D:\work\alpha`, Canonical: `D:\work\alpha`}
	}
	return risk.Result{
		Spelling: `%WISP181R3ACC%`, Canonical: `D:\work\alpha`,
		Rewritten: true, Rewrites: []string{"env"},
	}
}

func (s *movingScope) ResolveWorkspace(string) (risk.Result, error) {
	return risk.Result{}, risk.ErrRewrittenPath
}

func (s *movingScope) SetWorkspaceRoot(string, risk.Result) error { return nil }
