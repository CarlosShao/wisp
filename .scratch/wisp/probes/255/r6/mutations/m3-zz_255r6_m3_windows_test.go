//go:build windows

package main

// zz_255r6_m3 is the r6 M3 plant: a NewWithOptions site inside a TEST file.
// The 255r6 census judges production sources only (shipped code owns the promise),
// so this is the named, proven limit of the widened range, not an accident.
import webview2 "github.com/jchv/go-webview2"

func probe255r6TestFileSite() {
	w := webview2.NewWithOptions(webview2.WebViewOptions{})
	_ = w
}
