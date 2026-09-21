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
	// Ticket 117: the resident process is the leg owner actually uses - double
	// click the icon, no terminal attached, stderr going nowhere. Everything
	// this process logs (the seal notices of internal/winsec, a locked config
	// section loosened under L2, the D33 credential migration, the panic sink)
	// needs a listener on disk, or "you will see it when an authorization is
	// cleared" is a sentence that only holds in a test binary.
	//
	// Ordering: proc.Boot is the only thing that ran before this, and it seals
	// nothing - internal/proc has zero winsec imports - so no event can have
	// been missed. The close defer is registered BEFORE the shutdown defer so
	// LIFO runs the D38(e) sequence first and its own log lines still land.
	sink, sinkErr := installLogSink(rt.Layout.DataDir)
	if sinkErr != nil {
		// Loud, and it does not stop the app: a log directory that will not
		// open must not become a way to keep Wisp from starting. The notices
		// fall back to stderr, which is the state before this ticket.
		fmt.Fprintf(os.Stderr, "wisp: 持久日志未启用（%v）：安全告警只会到 stderr，不会落盘\n", sinkErr)
	} else {
		defer sink.close()
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
