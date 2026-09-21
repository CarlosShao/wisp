package panel

// Ticket 92 AC#3(ii)/AC#4 and AC#5(i): the native handler for a workspace
// request. The scope is scripted, so each case pins one decision the handler is
// responsible for: read the account, refuse loudly, audit both outcomes, and
// never leave the panel showing a scope the runtime is not using.

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

type scriptedScope struct {
	root    string
	res     risk.Result
	err     error
	setErr  error
	setArgs []string
}

func (s *scriptedScope) WorkspaceRoot() string { return s.root }

func (s *scriptedScope) ResolveWorkspace(string) (risk.Result, error) {
	return s.res, s.err
}

func (s *scriptedScope) SetWorkspaceRoot(root string) error {
	s.setArgs = append(s.setArgs, root)
	if s.setErr != nil {
		return s.setErr
	}
	s.root = root
	return nil
}

type auditLog []string

func (a *auditLog) f(format string, args ...any) {
	*a = append(*a, fmt.Sprintf(format, args...))
}

func (a auditLog) joined() string { return strings.Join(a, "\n") }

func TestWorkspaceSwitchAuditsAndAppliesACleanPath(t *testing.T) {
	scope := &scriptedScope{res: risk.Result{
		Spelling: `D:\work\Wisp\notes`, Canonical: `D:\work\Wisp\notes`, Resolved: true,
	}}
	var log auditLog
	view, err := RequestWorkspaceSwitch(scope, `D:\work\Wisp\notes`, log.f)
	if err != nil {
		t.Fatalf("a clean in-scope switch failed: %v", err)
	}
	if !view.Set || view.Canonical != `D:\work\Wisp\notes` {
		t.Errorf("view = %+v, want the resolved tree", view)
	}
	if len(scope.setArgs) != 1 || scope.setArgs[0] != `D:\work\Wisp\notes` {
		t.Errorf("scope.SetWorkspaceRoot calls = %v", scope.setArgs)
	}
	line := log.joined()
	for _, want := range []string{"workspace: SWITCH", "from=\"\"", "to=", "result=ok", "spelling="} {
		if !strings.Contains(line, want) {
			t.Errorf("audit line is missing %q:\n%s", want, line)
		}
	}
	t.Logf("audit: %s", line)
}

func TestWorkspaceSwitchPropagatesC26sReparseRefusal(t *testing.T) {
	scope := &scriptedScope{err: risk.ErrReparseDenied}
	var log auditLog
	view, err := RequestWorkspaceSwitch(scope, `D:\work\link`, log.f)
	if err == nil {
		t.Fatal("a junction-crossing workspace switch returned no error")
	}
	if !errors.Is(err, risk.ErrReparseDenied) {
		t.Errorf("the handler re-labelled C26's refusal: %v", err)
	}
	if view.Set {
		t.Errorf("the panel was told a workspace is set after a refusal: %+v", view)
	}
	if len(scope.setArgs) != 0 {
		t.Errorf("SetWorkspaceRoot was called despite the refusal: %v", scope.setArgs)
	}
	if !strings.Contains(log.joined(), "SWITCH-REFUSED") {
		t.Errorf("a refused switch was not audited:\n%s", log.joined())
	}
}

// TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolverSaysOK is AC#5(i)
// in its permission form: the handler reads ticket 102's book itself, so a
// resolver that stopped objecting cannot turn an expanded spelling into
// authority by silence.
func TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolverSaysOK(t *testing.T) {
	scope := &scriptedScope{res: risk.Result{
		Spelling: `%WORKON_HOME%\proj`, Canonical: `C:\other\tree`, Resolved: true,
		Rewritten: true, Rewrites: []string{"env"},
	}}
	var log auditLog
	view, err := RequestWorkspaceSwitch(scope, `%WORKON_HOME%\proj`, log.f)
	if err == nil {
		t.Fatal("a rewritten workspace was authorized - AC#5(i): the account was not read")
	}
	if !errors.Is(err, risk.ErrRewrittenPath) {
		t.Errorf("error %q does not carry risk.ErrRewrittenPath", err)
	}
	if len(scope.setArgs) != 0 {
		t.Errorf("the rewritten tree was handed to SetWorkspaceRoot anyway: %v", scope.setArgs)
	}
	if view.Rewritten {
		t.Errorf("the refusal view claims a rewritten workspace is in force: %+v", view)
	}
	if !strings.Contains(log.joined(), "SWITCH-REFUSED") {
		t.Errorf("the refusal was not audited:\n%s", log.joined())
	}
}

func TestWorkspaceSwitchKeepsThePreviousScopeOnFailure(t *testing.T) {
	scope := &scriptedScope{root: `d:\work\alpha`, err: errors.New("boom")}
	view, err := RequestWorkspaceSwitch(scope, `d:\work\beta`, nil)
	if err == nil {
		t.Fatal("expected the refusal to carry an error")
	}
	if view.Canonical != `d:\work\alpha` {
		t.Errorf("after a refusal the panel was shown %q instead of the scope still in force %q",
			view.Canonical, scope.root)
	}
	if !view.Set {
		t.Errorf("the still-active workspace rendered as unset: %+v", view)
	}
}

func TestWorkspaceSwitchDetectsAScopeThatMovedAnyway(t *testing.T) {
	// A scope whose SetWorkspaceRoot failed after ResolveWorkspace said yes, or
	// whose root changed behind the handler, must be reported as inconsistent -
	// not rendered as the folder the user asked for.
	scope := &scriptedScope{root: `d:\work\alpha`, setErr: errors.New("denied by acl")}
	view, err := RequestWorkspaceSwitch(scope, `d:\work\beta`, nil)
	if err == nil {
		t.Fatal("SetWorkspaceRoot failed but the switch reported success")
	}
	if strings.Contains(view.Reason, "已收窄") && view.Canonical == `d:\work\beta` {
		t.Errorf("the panel was told beta is in force while the scope refused it: %+v", view)
	}
	if view.Canonical != `d:\work\alpha` {
		t.Errorf("view = %+v, want the unchanged alpha", view)
	}
}

func TestWorkspaceViewRenderedForSnapshots(t *testing.T) {
	if v := WorkspaceViewFromRoot(""); v.Set || v.Reason == "" {
		t.Errorf("an unset workspace rendered as %+v: it must say so", v)
	}
	v := WorkspaceViewFromRoot(`d:\work\a`)
	if !v.Set || v.Canonical != `d:\work\a` {
		t.Errorf("a set workspace rendered as %+v", v)
	}
	if _, err := RequestWorkspaceSwitch(nil, "x", nil); err == nil {
		t.Error("a handler with no path scope reported success")
	}
}
