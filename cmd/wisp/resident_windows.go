//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
)

// runResident is the no-args path: boot (env layout, Job Object, single
// instance, goroutine registry), run the empty event loop, then exit through
// the D38(e) shutdown order. A second launch in the same session signals the
// running instance's activation event and exits (D42#7).
//
// This body is verbatim from main.go (ticket 78); it moved here because every
// call it makes is Windows-only - proc.Boot, proc.ErrAlreadyRunning and
// proc.SignalExistingInstance are declared in internal/proc's
// //go:build windows files, so an untagged main.go could not type-check under
// GOOS=linux. resident_other.go carries the refusal for every other platform.
func runResident() {
	printVersions("")

	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: %v\n", err)
		os.Exit(2)
	}

	rt, err := proc.Boot(env)
	if errors.Is(err, proc.ErrAlreadyRunning) {
		layout, lerr := proc.DefaultLayout(env)
		if lerr == nil && layout.MutexEnabled {
			_ = proc.SignalExistingInstance(layout.ActivateEventName)
		}
		fmt.Println("wisp: another instance is running in this session; activated it; exiting")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: boot failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		records := rt.Shutdown(false)
		failed := 0
		for _, rec := range records {
			if rec.Err != nil {
				failed++
				fmt.Printf("wisp: shutdown step %d (%s) FAILED: %v\n", rec.Step, rec.Name, rec.Err)
			}
		}
		fmt.Printf("wisp: exited through the D38(e) shutdown order (10 steps, %d failed)\n", failed)
	}()

	sum := rt.Layout.Summary()
	fmt.Printf("wisp: resident runtime booted (%s, data dir = %s, portable = %v, job object = on, single instance = %v)\n",
		sum.EnvBadge(), sum.DataDir, sum.Portable, rt.Instance != nil)
	fmt.Printf("wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)\n")

	reason := rt.RunEventLoop()
	fmt.Printf("wisp: event loop ending (%s); running the D38(e) shutdown order\n", reason)
}
