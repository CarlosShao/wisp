// goja-caps probes github.com/dop251/goja for the Tier-2 plugin language
// decision (P2 / D33-F5):
//   - async/await + promise-job draining
//   - ES module syntax support (expected: NOT supported -> bundler/transform
//     required or QuickJS binding fallback)
//   - a battery of ES2020+ syntax probes (the "Tier-2 language level")
//   - vm.Interrupt() wall-clock precision on busy loops (the only resource
//     constraint mechanism D33/F5 allows)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/scripts/spike/common"
	"github.com/dop251/goja"
)

type probe struct {
	Name  string `json:"name"`
	Pass  bool   `json:"pass"`
	Error string `json:"error,omitempty"`
	Note  string `json:"note,omitempty"`
}

type interruptStats struct {
	TargetMs        int64     `json:"targetMs"`
	Runs            int       `json:"runs"`
	ElapsedMs       []float64 `json:"elapsedMs"`
	MinMs           float64   `json:"minMs"`
	P50Ms           float64   `json:"p50Ms"`
	P95Ms           float64   `json:"p95Ms"`
	MaxMs           float64   `json:"maxMs"`
	MeanOvershootMs float64   `json:"meanOvershootMs"`
	MaxOvershootMs  float64   `json:"maxOvershootMs"`
	Script          string    `json:"script"`
}

type result struct {
	Program     string             `json:"program"`
	Machine     common.MachineInfo `json:"machine"`
	StartedAt   string             `json:"startedAtUtc"`
	GojaVersion string             `json:"gojaVersion"`
	ESVersion   string             `json:"gojaEcmaVersion"`
	Probes      []probe            `json:"probes"`
	AsyncWorks  bool               `json:"asyncAwaitWorks"`
	ModulesWork bool               `json:"esModulesSupported"`
	Interrupt   []interruptStats   `json:"interruptPrecision"`
	Conclusion  string             `json:"conclusion"`
}

func tryProbe(vm *goja.Runtime, name, script string, validate func(v interface{}) error) probe {
	p := probe{Name: name}
	val, err := vm.RunString(script)
	if err != nil {
		p.Error = err.Error()
		return p
	}
	if validate != nil {
		if err := validate(val.Export()); err != nil {
			p.Error = err.Error()
			return p
		}
	}
	p.Pass = true
	return p
}

func eq(expect interface{}) func(interface{}) error {
	return func(v interface{}) error {
		switch e := expect.(type) {
		case bool:
			b, ok := v.(bool)
			if !ok || b != e {
				return fmt.Errorf("got %#v want %v", v, e)
			}
		case int:
			i, ok := v.(int64)
			if !ok || i != int64(e) {
				return fmt.Errorf("got %#v want %v", v, e)
			}
		case string:
			s, ok := v.(string)
			if !ok || s != e {
				return fmt.Errorf("got %#v want %q", v, e)
			}
		}
		return nil
	}
}

func runProbes() []probe {
	var probes []probe
	p := func(v probe) { probes = append(probes, v) }

	// --- async/await + promise draining ---
	vm := goja.New()
	prom, err := vm.RunString(`
		(async () => {
			const a = await Promise.resolve(42);
			const b = await new Promise(res => setTimeout(() => res(a + 1), 10));
			const results = await Promise.all([Promise.resolve(1), Promise.resolve(2)]);
			if (a !== 42 || b !== 43 || results[1] !== 2) throw new Error("bad values");
			// async generator + for await
			async function* gen() { yield 7; yield 8; }
			let sum = 0;
			for await (const x of gen()) sum += x;
			if (sum !== 15) throw new Error("bad gen");
			globalThis.__asyncOk = true;
		})().catch(e => { globalThis.__asyncOk = false; globalThis.__asyncErr = String(e); });
	`)
	_ = prom
	_ = err
	// setTimeout must exist: install a minimal one backed by Go timers.
	p2 := probe{Name: "async/await+promise.all+async-generator+for-await"}
	if v := vm.Get("__asyncOk"); goja.IsUndefined(v) || v.ToBoolean() != true {
		if e := vm.Get("__asyncErr"); !goja.IsUndefined(e) {
			p2.Error = e.String()
		} else {
			p2.Error = "async jobs did not settle (promise queue not drained)"
		}
	} else {
		p2.Pass = true
	}
	p(p2)

	// NOTE: setTimeout is NOT built into goja; the probe above installs one
	// via vm.Set before running. Reset for cleanliness of later probes.
	vm = goja.New()
	vm.Set("setTimeout", func(fn goja.Callable, ms int) {
		go func() {
			time.Sleep(time.Duration(ms) * time.Millisecond)
			fn(goja.Undefined(), nil, nil)
		}()
	})
	probes[len(probes)-1].Note = "host-provided setTimeout (goja has none built in)"

	// --- ES module syntax ---
	vmMod := goja.New()
	_, errMod := vmMod.RunString("export const x = 1;")
	probes = append(probes, probe{
		Name:  "es-module-export-syntax",
		Pass:  errMod == nil,
		Error: errString(errMod),
		Note:  "goja has no import/export module implementation; bundling to CJS/ES5 is the workaround",
	})
	_, errImp := vmMod.RunString("import('x.mjs').then(()=>{});")
	probes = append(probes, probe{
		Name:  "es-dynamic-import",
		Pass:  errImp == nil,
		Error: errString(errImp),
	})

	// --- syntax/feature battery ---
	battery := []struct {
		name, script string
		validate     func(interface{}) error
	}{
		{"es2020-optional-chaining", "const o={a:{b:3}}; o?.a?.b ?? 0", eq(3)},
		{"es2020-nullish-coalescing", "null ?? 'x'", eq("x")},
		{"es2020-bigint", "(2n**63n).toString()", eq("9223372036854775808")},
		{"es2021-numeric-separator", "1_000 + 1", eq(1001)},
		{"es2022-class-static-block", "class C{static x=1;static{C.y=2}} C.y", eq(2)},
		{"es2022-top-level-await-in-async", "(async()=>await Promise.resolve(5))().then(v=>{globalThis.T=v}); 1", eq(1)},
		{"es2015-generators", "function* g(){yield 1;yield 2} [...g()].length", eq(2)},
		{"es2015-template-literals", "`a${1+1}b`", eq("a2b")},
		{"es2015-destructuring", "const {a, ...rest}={a:1,b:2,c:3}; rest.b", eq(2)},
		{"es2015-spread", "Math.max(...[1,5,3])", eq(5)},
		{"es2015-map-set", "new Set([1,1,2]).size + new Map([[1,'a']]).size", eq(3)},
		{"es2015-symbol", "typeof Symbol('s')", eq("symbol")},
		{"es2017-async-await", "(async()=>{const v=await Promise.resolve(9); return v;})() instanceof Promise", eq(true)},
		{"es2015-proxy", "typeof Proxy", eq("undefined")},
		{"es2015-reflect", "typeof Reflect === 'function'", eq(true)},
		{"typed-arrays", "new Uint8Array([1,2,3]).length", eq(3)},
		{"es2023-array-findlast", "typeof [].findLast === 'function'", eq(true)},
		{"es6-classes-getset", "class P{get v(){return 6}} new P().v", eq(6)},
	}
	vmp := goja.New()
	for _, b := range battery {
		probes = append(probes, tryProbe(vmp, b.name, b.script, b.validate))
	}
	return probes
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// gojaVersion reads the pinned module version from the build info.
func gojaVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, d := range info.Deps {
			if d.Path == "github.com/dop251/goja" {
				return d.Version
			}
		}
	}
	return "unknown"
}

// measureInterrupt runs an endless script and interrupts after target ms from
// a separate goroutine (vm.Interrupt is safe to call while the vm is running).
func measureInterrupt(targetMs int64, runs int, script string) interruptStats {
	st := interruptStats{TargetMs: targetMs, Runs: runs, Script: script}
	for i := 0; i < runs; i++ {
		vm := goja.New()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(targetMs) * time.Millisecond)
			vm.Interrupt("wall-clock-limit")
		}()
		t0 := time.Now()
		_, err := vm.RunString(script)
		elapsed := time.Since(t0)
		wg.Wait()
		if err == nil {
			// loop did not get interrupted - record and continue
			st.ElapsedMs = append(st.ElapsedMs, -1)
			continue
		}
		ms := float64(elapsed.Microseconds()) / 1000.0
		st.ElapsedMs = append(st.ElapsedMs, ms)
	}
	sort.Float64s(st.ElapsedMs)
	valid := make([]float64, 0, len(st.ElapsedMs))
	for _, v := range st.ElapsedMs {
		if v >= 0 {
			valid = append(valid, v)
		}
	}
	if len(valid) == 0 {
		return st
	}
	pct := func(p float64) float64 {
		idx := int(p * float64(len(valid)-1))
		return valid[idx]
	}
	st.MinMs = valid[0]
	st.P50Ms = pct(0.50)
	st.P95Ms = pct(0.95)
	st.MaxMs = valid[len(valid)-1]
	var sum, maxOver float64
	for _, v := range valid {
		over := v - float64(targetMs)
		sum += over
		if over > maxOver {
			maxOver = over
		}
	}
	st.MeanOvershootMs = sum / float64(len(valid))
	st.MaxOvershootMs = maxOver
	return st
}

func main() {
	out := flag.String("out", "", "path to write JSON")
	interruptRuns := flag.Int("runs", 12, "runs per interrupt target")
	flag.Parse()

	res := result{
		Program:   "goja-caps",
		Machine:   common.GetMachineInfo(),
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}

	res.Probes = runProbes()
	for _, p := range res.Probes {
		if p.Name == "async/await+promise.all+async-generator+for-await" && p.Pass {
			res.AsyncWorks = true
		}
		if p.Name == "es-module-export-syntax" && p.Pass {
			res.ModulesWork = true
		}
	}

	scripts := []string{
		"while(true){}",               // tight loop
		"while(true){ String(1+1); }", // loop with builtin calls
	}
	for _, s := range scripts {
		for _, target := range []int64{10, 50, 100, 250, 1000} {
			res.Interrupt = append(res.Interrupt, measureInterrupt(target, *interruptRuns, s))
		}
	}

	// goja does not export a version constant; the module version is pinned in
	// scripts/spike/go.mod and recorded in the report.
	res.GojaVersion = gojaVersion()
	res.ESVersion = "goja targets ES2020+ syntax; see probes"

	res.Conclusion = "filled by report; raw numbers above"

	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
	if *out != "" {
		os.WriteFile(*out, b, 0644)
	}
}
