package main

// Ruling A11 / ticket 11 AC#6: `llm.RunProbeSuite` had nine green tests and no
// production caller, so the 「声明 PASS / 实测 FAIL」 event could never fire on a
// real machine. These cases drive the composition root's own provider path
// against the live mockllm and pin both directions with the SERVER's
// capability mode.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
)

// providersFixture is a data dir with a config that names the live mockllm, plus
// the CLI surface to drive.
type providersFixture struct {
	t     *testing.T
	srv   *adaptertest.Mockllm
	dir   string
	out   *bytes.Buffer
	err   *bytes.Buffer
	store *memory.Store
}

func newProvidersFixture(t *testing.T) *providersFixture {
	t.Helper()
	pf := &providersFixture{
		t: t, srv: adaptertest.StartMockllm(t),
		out: &bytes.Buffer{}, err: &bytes.Buffer{},
	}
	pf.dir = sealableTempDir124(t)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/mock-small"]

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q

[llm.providers.acme.models.mock-small]
context_window = 128000

[llm.providers.acme.models.mock-small.capabilities]
text = true
fc = true
vision = true
thinking = true

[llm.providers.acme.models.mock-small.thinking_levels]
`, pf.srv.Base+"/v1")
	// thinking_levels is a []string; leaving it absent keeps the defaults.
	body = strings.Replace(body, "\n[llm.providers.acme.models.mock-small.thinking_levels]\n", "\n", 1)
	if err := os.WriteFile(filepath.Join(pf.dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := memory.Open(pf.dir)
	if err != nil {
		t.Fatal(err)
	}
	pf.store = st
	t.Cleanup(func() { _ = st.Close() })
	return pf
}

func (pf *providersFixture) call(argv ...string) int {
	pf.t.Helper()
	return cmdProviders(argv, providersIO{stdout: pf.out, stderr: pf.err, dataDir: pf.dir})
}

// health reads back the row the probe wrote.
func (pf *providersFixture) health() memory.ProviderHealth {
	pf.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, err := pf.store.ProviderHealthRow(ctx, "acme", "mock-small")
	if err != nil {
		pf.t.Fatalf("provider_health: %v", err)
	}
	return row
}

// probeVerdicts parses the verdict column out of a `wisp providers probe`
// report. A verdict row is printed as "  <capability> <声明|未声明> 实测 <verdict>"
// (providers.go's table), so the rows are exactly the two-space-indented lines
// whose third whitespace field is 实测.
//
// Ticket 67 AC#3 pins this column POSITIVELY: every row's verdict must equal
// the literal "PASS" or "FAIL" - not merely "not a checkmark". The d22scan
// emoji ban covers U+2190-U+2BFF and currently cannot see cmd/ at all; a third
// symbol sneaking into this column must fail HERE, not wait for the gate.
func probeVerdicts(t *testing.T, out string) map[string]string {
	t.Helper()
	verdicts := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "    ") {
			continue // verdict rows carry exactly two leading spaces
		}
		fields := strings.Fields(line)
		if len(fields) != 4 || fields[2] != "实测" {
			t.Errorf("unrecognised two-space row in probe report: %q", line)
			continue
		}
		verdict := fields[3]
		if verdict != "PASS" && verdict != "FAIL" {
			t.Errorf("verdict for %s = %q, want exactly PASS or FAIL", fields[0], verdict)
			continue
		}
		if _, dup := verdicts[fields[0]]; dup {
			t.Errorf("duplicate verdict row for %s", fields[0])
		}
		verdicts[fields[0]] = verdict
	}
	if len(verdicts) == 0 {
		t.Error("the probe report printed no verdict rows")
	}
	return verdicts
}

// TestProvidersProbeRecordsMeasuredThinkingFalse is the event firing on a real
// machine: the catalog declares thinking, the server withholds every reasoning
// delta, and `wisp providers probe` says so out loud AND writes the verdict.
func TestProvidersProbeRecordsMeasuredThinkingFalse(t *testing.T) {
	pf := newProvidersFixture(t)
	pf.srv.Reset(t)
	pf.srv.Control(t, "/__control/capability", `{"thinking":"broken"}`)

	code := pf.call("probe", "acme/mock-small")
	if code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, pf.out.String(), pf.err.String())
	}
	out, errb := pf.out.String(), pf.err.String()
	if !strings.Contains(errb, "能力实测不符") || !strings.Contains(errb, "thinking") {
		t.Errorf("the declared-vs-measured event was not announced:\n%s", errb)
	}
	if !strings.Contains(out, "thinking") {
		t.Errorf("the probe report is missing the thinking line:\n%s", out)
	}
	// Ticket 67 AC#3: the verdict column reads in ASCII, pinned POSITIVELY
	// (every row is exactly PASS or FAIL) and the two rows this case proves
	// are named: thinking withheld => FAIL, fc measured OK => PASS.
	verdicts := probeVerdicts(t, out)
	if verdicts["thinking"] != "FAIL" {
		t.Errorf("thinking verdict = %q, want FAIL", verdicts["thinking"])
	}
	if verdicts["fc"] != "PASS" {
		t.Errorf("fc verdict = %q, want PASS", verdicts["fc"])
	}
	// The server is provably the cause, not the assertion: three probe requests
	// went out (fc, vision, thinking) on the chat route.
	if n := pf.srv.RouteCount(t, "chat"); n < 3 {
		t.Errorf("mockllm served %d chat requests, want >= 3 probes", n)
	}
	h := pf.health()
	if h.Probe.Thinking == nil || *h.Probe.Thinking != false {
		t.Errorf("persisted thinking = %v, want measured false", h.Probe.Thinking)
	}
	if h.Probe.FC == nil || *h.Probe.FC != true {
		t.Errorf("persisted fc = %v, want measured true (the run must not blanket-report broken)", h.Probe.FC)
	}
	if h.LastProbeAt.IsZero() {
		t.Error("the probe record carries no timestamp")
	}
}

// TestProvidersProbeRecordsMeasuredThinkingTrue is the same command with the
// same server in its default mode: the ONLY thing that changed is the server's
// behavior, so a verdict that moved with it is a measurement.
func TestProvidersProbeRecordsMeasuredThinkingTrue(t *testing.T) {
	pf := newProvidersFixture(t)
	pf.srv.Reset(t)
	if code := pf.call("probe", "acme/mock-small"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, pf.out.String(), pf.err.String())
	}
	if strings.Contains(pf.err.String(), "能力实测不符") {
		t.Errorf("a capable server must not produce a mismatch: %s", pf.err.String())
	}
	// The other side of the same positive pin: with a capable server every
	// verdict row reads PASS. The fixture declares fc/vision/thinking all
	// true, so a row that were FAIL here would also have raised the mismatch
	// asserted absent above - this closes the loop on the printed table.
	verdicts := probeVerdicts(t, pf.out.String())
	if verdicts["thinking"] != "PASS" {
		t.Errorf("thinking verdict = %q, want PASS", verdicts["thinking"])
	}
	for name, v := range verdicts {
		if v != "PASS" {
			t.Errorf("capable server measured %s = %q, want PASS", name, v)
		}
	}
	if h := pf.health(); h.Probe.Thinking == nil || *h.Probe.Thinking != true {
		t.Errorf("persisted thinking = %v, want measured true", h.Probe.Thinking)
	}
}

// TestProvidersDiscoverListsWhatTheServerServes covers the other half of the
// save/discovery path: the command asks the provider, not the catalog.
func TestProvidersDiscoverListsWhatTheServerServes(t *testing.T) {
	pf := newProvidersFixture(t)
	pf.srv.Reset(t)
	if code := pf.call("discover", "acme"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, pf.out.String(), pf.err.String())
	}
	if pf.srv.RouteCount(t, "models") == 0 {
		t.Fatal("discovery never touched the provider's /v1/models")
	}
	if !strings.Contains(pf.out.String(), "wisp providers: acme 上报") {
		t.Errorf("no model list printed:\n%s", pf.out.String())
	}
	if code := pf.call("discover", "nope"); code == 0 {
		t.Error("an unknown provider must not exit 0")
	}
}

// TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless: the discovery/probe path
// resolves api_key_ref through the store like the text path, so a ref with no
// blob behind it must fail rather than probe unauthenticated.
func TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless(t *testing.T) {
	pf := newProvidersFixture(t)
	body, err := os.ReadFile(filepath.Join(pf.dir, configFileName))
	if err != nil {
		t.Fatal(err)
	}
	rewritten := strings.Replace(string(body),
		`base_url = "`+pf.srv.Base+`/v1"`,
		`base_url = "`+pf.srv.Base+`/v1"`+"\n"+`api_key_ref = "dpapi:absent"`, 1)
	if err := os.WriteFile(filepath.Join(pf.dir, configFileName), []byte(rewritten), 0o600); err != nil {
		t.Fatal(err)
	}
	pf.srv.Reset(t)
	if code := pf.call("probe", "acme/mock-small"); code == 0 {
		t.Error("a probe with an unresolvable key must not succeed")
	}
	if !strings.Contains(pf.err.String(), "Unconfigured") {
		t.Errorf("the failure must name the Unconfigured path:\n%s", pf.err.String())
	}
	if n := pf.srv.RouteCount(t, "chat"); n != 0 {
		t.Errorf("%d probe requests went out with an unresolvable key", n)
	}
}
