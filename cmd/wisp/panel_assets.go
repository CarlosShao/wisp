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
// WISP-LEG-COVERAGE-RULING: panel-assets is dispatched by main and driven by no
// case in this package. It is a read-only diagnostic over bytes the binary already
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
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wisp panel-assets [-manifest] [-check] [-render <path>] [-l2 <tool> <args...>] [-irreversible <ops>]")
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
		view := panel.NewApprovalCardView(risk.NewRiskAssessor(), panel.ApprovalSubject{
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
