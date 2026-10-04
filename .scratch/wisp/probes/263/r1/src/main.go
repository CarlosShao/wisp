// Command probe263 is a throwaway stand-in for wisp.exe used by ticket 263's
// positive controls. It mimics the one contract scripts/slo-check.ps1 depends
// on: `wisp slo ...` samples for -seconds N, writes a StateReport/SettleReport
// JSON to -out, and exits 0 when that report passes (exit 1 when it does not,
// which is also what the forced leak fixture does).
//
// It is built twice into bin/: once as a Windows console-subsystem binary and
// once with -ldflags -H=windowsgui - the flag scripts/build.ps1 has passed
// since ticket 244, which is what made wisp.exe a process PowerShell does not
// wait for. Same program, two subsystems, so the gate can be observed against
// both without touching build/ or anyone else's artifact.
//
// Env knobs used by the evidence runs (they let a run stay byte-identical while
// the SAMPLE changes, i.e. "what if the machine really is over budget"):
//
//	PROBE263_VERDICT=fail   write pass:false, i.e. a real failed sample
//	PROBE263_LEAK_EXIT=<n>  exit code used for a -leak run (default 1)
package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

func flagValue(args []string, name, def string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == name {
			return args[i+1]
		}
	}
	return def
}

func hasFlag(args []string, name string) bool {
	for _, a := range args {
		if a == name {
			return true
		}
	}
	return false
}

func writeOut(path string, doc any) {
	if path == "" {
		return
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

func main() {
	args := os.Args[1:]
	secs, _ := strconv.Atoi(flagValue(args, "-seconds", "2"))
	out := flagValue(args, "-out", "")
	state := flagValue(args, "-state", "")
	pass := os.Getenv("PROBE263_VERDICT") != "fail"
	code := 0

	if hasFlag(args, "-settle") {
		writeOut(out, map[string]any{
			"pass":    pass,
			"settle":  map[string]any{"free_os_memory_count": 1, "seconds": secs},
			"sampler": "probe263",
		})
	} else {
		writeOut(out, map[string]any{
			"pass":    pass,
			"state":   state,
			"cpu_pct": 0.4,
			"sampler": "probe263",
		})
	}
	if hasFlag(args, "-leak") {
		// The forced 100MB leak is supposed to FAIL the state gate; if it ever
		// exits 0 the sampler is broken, and slo-check.ps1 reddens on purpose.
		code = 1
		if v, err := strconv.Atoi(os.Getenv("PROBE263_LEAK_EXIT")); err == nil {
			code = v
		}
	} else if !pass {
		// Contract-faithful: `wisp slo` exits 0 pass / 1 fail (see the header
		// comment of slo-check.ps1), so a sample over budget carries BOTH the
		// pass:false report and a non-zero exit code.
		code = 1
	}

	os.Stderr.WriteString("probe263: sampling started pid=" + strconv.Itoa(os.Getpid()) +
		" state=" + state + " seconds=" + strconv.Itoa(secs) + "\n")
	time.Sleep(time.Duration(secs) * time.Second)
	os.Stderr.WriteString("probe263: exiting with code " + strconv.Itoa(code) + "\n")
	os.Exit(code)
}
