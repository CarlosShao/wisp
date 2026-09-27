package main

// `wisp panel-assets` - the diagnostic face of the embedded panel bundle
// (ticket 77 AC#1) and of the L2 card payload (AC#3).
//
// Why a subcommand exists at all: AC#1's judgement is that the panel reaches a
// machine with no node, no npm and no network. That is only provable if
// something in the binary can hand back the page bytes without a browser being
// involved, so this command IS that proof surface, and it is also what an
// operator runs when a panel looks wrong.
//
// It is additive by construction: the resident GUI path, the ball and the
// agent loop are untouched by this file. Ticket 33's WebView2 host is expected
// to call the same panel.Assets API this command prints.
//
// WISP-LEG-COVERAGE-RULING: panel-assets is dispatched by main; the card branch
// below is driven by cmd/wisp's panel_assets_143_test.go and the asset branch is
// still ruled. It is a read-only diagnostic over bytes the binary already
// carries (internal/panel's //go:embed bundle), so deleting this branch costs an
// operator a window on the bundle and books no record anywhere; the bundle itself
// is checked from the other side by internal/panel's asset tests. This sentence is
// here because ticket 133's leg census asks every dispatched command to be nailed,
// driven, or ruled, and this one is ruled.

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

func cmdPanelAssets(args []string) int {
	fs := flag.NewFlagSet("panel-assets", flag.ContinueOnError)
	manifest := fs.Bool("manifest", false, "list every file the binary carries, with size and fingerprint")
	check := fs.Bool("check", false, "verify the entry file and every asset it references are inside this binary")
	render := fs.String("render", "", "write the embedded bytes for this request path to stdout")
	l2 := fs.String("l2", "", "assess this tool name with the remaining args as its argv and print the L2 card JSON")
	l2Irreversible := fs.String("irreversible", "", "with -l2: comma-separated irreversible operations the CALLER declares in its risk.Facts (R8's input)")
	taintSources := &taintSourceFlag{}
	fs.Var(taintSources, "taint-source", "with -l2: repeatable INPUT FACT of the shape <source-tool>|<origin>|<content>, meaning \"this task read <content> from <source-tool> at <origin>\". It feeds the C25 provenance engine the taint rule judges; it declares no verdict, no rule id and no level, and the taint rule stays dormant without it. Every flag has to come before the tool's own arguments: parsing stops at the first argument that does not start with a dash, and everything after it is argv.")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wisp panel-assets [-manifest] [-check] [-render <path>] [-taint-source <source-tool>|<origin>|<content>]... [-irreversible <ops>] -l2 <tool> <args...>")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}

	// The card path is the panel's data source: internal/risk decides, this
	// prints exactly the JSON the WebView2 host pushes. Ticket 35 replaces the
	// printing with a bridge push; nothing else about the payload changes, which
	// is why the frontend's ApprovalCardView keys are pinned against it.
	if *l2 != "" {
		rest := fs.Args()
		// Same assembly the production bridge uses (internal/tools/bridge.go's
		// assessorFor): a bare assessor plus the C25 engine bound to one scope.
		// With no -taint-source there is nothing to bind, so the assessor stays
		// exactly what it was before ticket 143 and the printed card is
		// unchanged, byte for byte.
		assessor := risk.NewRiskAssessor()
		if det := taintSources.detector(); det != nil {
			assessor = assessor.WithTaintDetector(det)
		}
		view := panel.NewApprovalCardView(assessor, panel.ApprovalSubject{
			CorrelationID: "panel-assets-l2",
			Tool:          *l2,
			Args:          rest,
			CallChain:     []string{"cli", "panel-assets", *l2},
			Facts: risk.Facts{
				Declared:     risk.L1,
				Paths:        rest,
				Irreversible: splitFacts(*l2Irreversible),
			},
		})
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(view); err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: encode: %v\n", err)
			return 1
		}
		return 0
	}

	assets, err := panel.BuiltinAssets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp panel-assets: %v\n", err)
		return 1
	}
	switch {
	case *render != "":
		data, contentType, err := assets.Resolve(*render)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "wisp panel-assets: %s %d bytes %s\n", *render, len(data), contentType)
		if _, err := os.Stdout.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: write: %v\n", err)
			return 1
		}
	case *manifest:
		lines, err := assets.Manifest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	case *check:
		// The "the embed carried less than the build produced" detector: an
		// operator, and the AC#1 evidence run, both end here rather than
		// eyeballing a blank panel.
		served, err := assets.Check()
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: %v\n", err)
			return 1
		}
		fmt.Printf("panel assets check: entry=%s built=%t %d asset refs resolve [%s]\n",
			panel.EntryFile, assets.Built(), len(served), strings.Join(served, " "))
	default:
		if !assets.Built() {
			fmt.Fprintln(os.Stderr, "wisp panel-assets: assets NOT BUILT")
			return 1
		}
		lines, err := assets.Manifest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "wisp panel-assets: %v\n", err)
			return 1
		}
		fmt.Printf("panel assets embedded: %d files, entry=%s built=%t\n", len(lines), panel.EntryFile, assets.Built())
	}
	return 0
}

// splitFacts turns a comma-separated flag value into the string slice a
// risk.Facts field expects, dropping empties so "-irreversible ," is the same
// as not passing it at all.
//
// This exists because risk.Facts is judged INPUT, not something the card or
// this command may invent: R8 fires on what the caller declared irreversible,
// so a probe that wants the real L2 verdict has to be able to say so. Nothing
// here derives a level - internal/risk still decides that.
func splitFacts(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// taintSourceScopeID is the C25 scope this command opens for its one card. In
// the product the scope id is a task id owned by the agent loop (ticket 19's
// DEFERRED(C25-loop-wiring) item 1); a CLI probe has no task, so it names one
// and closes nothing - the engine lives and dies inside this process.
const taintSourceScopeID = "panel-assets-l2"

// taintSource is one declared sensitive-source read: the operator says which
// tool produced content, where it came from, and what the content was. Those
// three fields are exactly C25's Mark(scope, tool, origin, content) input, and
// nothing else - no level, no rule id, no explanation text. The taint rule's
// sentence and its hit are produced by internal/risk, which is the only thing
// in this binary allowed to decide them (SPEC-06 §3).
type taintSource struct {
	tool    string
	origin  string
	content string
}

// taintSourceFlag collects repeatable -taint-source values as a flag.Value so
// one probe can declare several sources, which is what the engine's
// oldest-mark-first hit order is for.
type taintSourceFlag struct {
	sources []taintSource
}

// String keeps flag's default-printing honest: the default is "no source
// declared", and printing the collected values here would echo declared
// content back at whoever ran the command.
func (f *taintSourceFlag) String() string { return "" }

// Set parses <tool>|<origin>|<content>. SplitN(3) means the content may itself
// contain a pipe; only the first two separators are structure. The origin may
// be empty (a source with no location); the tool and the content may not,
// because an unnamed source prints an empty slot on the card and empty content
// is a fact nobody declared. Content shorter than C25's fragment floor cannot
// match anything (internal/risk/taintmatch.go documents that residual); this
// probe does not pretend to re-judge it and passes it to the engine as given.
func (f *taintSourceFlag) Set(raw string) error {
	parts := strings.SplitN(raw, "|", 3)
	tool := strings.TrimSpace(parts[0])
	if tool == "" {
		return fmt.Errorf("-taint-source needs <source-tool>|<origin>|<content>, the source tool name is empty")
	}
	if len(parts) < 3 {
		return fmt.Errorf("-taint-source %q needs three |-separated parts: <source-tool>|<origin>|<content>", tool)
	}
	content := parts[2]
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("-taint-source %q has empty content: nothing can match against it", tool)
	}
	f.sources = append(f.sources, taintSource{tool: tool, origin: strings.TrimSpace(parts[1]), content: content})
	return nil
}

// detector builds the C25 engine over the declared sources and returns it bound
// to this command's scope - the same call shape the production bridge uses
// (risk.NewProvenance(...).Detector(taskID), see internal/tools/bridge.go's
// assessorFor). Returns nil when no source was declared, which is the caller's
// signal to wire nothing at all.
//
// ProvOptions.NoProbe is the one difference from the bridge, and it cannot
// change a verdict on this path: the panel builds its parameters from the
// argv (internal/panel's paramsFromArgs yields the "command" and "argv" keys
// only), neither of which is a write target, so the sync-directory gate this
// option would feed is never consulted for a call of that shape.
func (f *taintSourceFlag) detector() risk.TaintDetector {
	if len(f.sources) == 0 {
		return nil
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	prov.OpenScope(taintSourceScopeID)
	for _, s := range f.sources {
		prov.Mark(taintSourceScopeID, s.tool, s.origin, s.content)
	}
	return prov.Detector(taintSourceScopeID)
}
