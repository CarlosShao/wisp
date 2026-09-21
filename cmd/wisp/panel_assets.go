package main

// `wisp panel-assets` - the diagnostic face of the embedded panel bundle
// (ticket 77 AC#1).
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

import (
	"flag"
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/panel"
)

func cmdPanelAssets(args []string) int {
	fs := flag.NewFlagSet("panel-assets", flag.ContinueOnError)
	manifest := fs.Bool("manifest", false, "list every file the binary carries, with size and fingerprint")
	render := fs.String("render", "", "write the embedded bytes for this request path to stdout")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wisp panel-assets [-manifest] [-render <path>]")
	}
	if err := fs.Parse(args); err != nil {
		return 2
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
