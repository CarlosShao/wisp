//go:build windows

package aecspike

// The 389-r1 leg's two-state readings (the product-form gate on the good and
// collapse states) and the reverse criterion, taken OFFLINE over the 59-w-r1
// rig artifacts on disk: no rig, no audio devices, no ASR rerun (every
// episode carries its transcription in gates-probe.json). The product gate
// code (audio.BargeinGate) judges every episode through its real decision
// path, and its criteria run for real: each material's windows are folded
// through ObserveWindow as RAW values (the curve's per-window ERLE and the
// per-window correlation), so the collapse detector's streak logic and the
// correlation floor both execute inside the product gate - nothing is
// pre-ORed. This is the dispatch's material-reuse boundary: the rig is
// saturated (59-v3 section 7), and the 400ms stop-playback chain, the
// product-microphone zero-transcription surface and the -30 margin
// re-measurement belong to ticket 26's integration surface, not here.
//
// Materials (the 59-w gates-probe set, docs/evidence/s1/59wr1-longrun-gates.md
// section 4):
//
//	good     - 9 materials, 222.0s  (long-run steady spans + good cold runs)
//	collapse - 44 materials, 104.4s (long-run event windows + cold collapse
//	           + the archived 59-v2 collapse wav)
//	control  - control-26 / control-36 from gates.json: real TTS speech
//	           mixed over the long-music-60 steady span (the measured
//	           episode levels -26.70 / -36.58 dBFS RMS)
//
// Baseline being re-measured (same table): tier-3 false triggers at the
// loose -45 axis are good 2 events/222s = 5.4/10min vs collapse 2/104.4s =
// 11.5/10min; at -30 both states are zero. The dynamic gate's job: pull the
// collapse state to the good level or below WITHOUT losing the quiet-speech
// reach the loose nominal buys (control-36 passes at -45 where the static
// -30 gate ate it), and survive double-talk (control-26 at -26.70 passes
// even under the tightened cut).
//
// Correlation sources per far content, exactly the rig leg's discriminator
// choice: music runs have a bit-exact far (SynthMusic deterministic) so the
// probe's corr_vs_far is used; TTS runs do not (matcha is not
// bit-deterministic across processes, corr[0]=0.126 on the rebuild), so the
// reference-free self_corr_at_period is used (valid when > 0: the probe
// writes 0 for the early windows and for inactive windows).

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
)

const bargein389r1Out = `D:\tmp\wisp\20261031\389r1\out`

// probeGateMaterial is one material row of gates-probe.json / gates.json.
type probeGateMaterial struct {
	Material    string  `json:"material"`
	State       string  `json:"state"`
	Source      string  `json:"source"`
	DurationSec float64 `json:"duration_sec"`
	Episodes    []struct {
		StartSample int     `json:"start_sample"`
		Samples     int     `json:"samples"`
		RmsDbfs     float64 `json:"rms_dbfs"`
		AsrText     string  `json:"asr_text"`
	} `json:"episodes"`
}

type probeGates map[string]probeGateMaterial

// autocorrRun mirrors the self-correlation probe's output.
type autocorrRun struct {
	Windows []autocorrWindow `json:"windows"`
	Events  []struct {
		StartWin int    `json:"start_window"`
		EndWin   int    `json:"end_window"`
		Class    string `json:"class"`
	} `json:"events"`
}

type autocorrWindow struct {
	Window   int     `json:"window"`
	ERLE     any     `json:"erle_db"`
	SelfCorr float64 `json:"self_corr_at_period"`
	Class    string  `json:"class"`
}

// probeCorrRun is the slice of probe.json the music-arm correlation comes
// from.
type probeCorrRun struct {
	Windows []struct {
		Corr float64 `json:"corr_vs_far"`
	} `json:"windows"`
}

// bargeinCurvePoint is one window of a run's ERLE curve.
type bargeinCurvePoint struct {
	Window int `json:"window"`
	ERLE   any `json:"erle_db"`
}

// bargeinWinVal is one global window's raw criterion values on the re-cut.
type bargeinWinVal struct {
	ERLE any     // the curve's erle_db (a float64, or "inactive")
	Corr float64 // the window correlation; NaN when not measurable
}

// bargeinMatRef locates a material on its source run's window axis.
type bargeinMatRef struct {
	run    string
	kind   string // "steady" | "event" | "whole" | "archived"
	event  int
	offset int // global window of material window 0 (-1: not mappable)
}

func bargeinParseMaterial(name string) bargeinMatRef {
	if name == "archived-59v2-collapse" {
		return bargeinMatRef{kind: "archived", offset: -1}
	}
	if rest, ok := strings.CutSuffix(name, "-steady"); ok {
		return bargeinMatRef{run: rest, kind: "steady", offset: 10}
	}
	if i := strings.Index(name, "-event"); i >= 0 {
		k, err := strconv.Atoi(name[i+len("-event"):])
		if err != nil {
			return bargeinMatRef{kind: "archived", offset: -1}
		}
		return bargeinMatRef{run: name[:i], kind: "event", event: k, offset: -1}
	}
	return bargeinMatRef{run: name, kind: "whole", offset: 0} // the cold runs
}

// bargeinWindowStat rebuilds a closed window's stat from its curve ERLE (the
// re-cut's ERLE unit is dB; powers are reconstructed at res=1).
func bargeinWindowStat(idx int, erle any) audio.BargeinWindowStat {
	w := audio.BargeinWindowStat{Index: idx}
	if f, ok := erle.(float64); ok && !math.IsNaN(f) {
		w.ResPow, w.NearPow, w.ActiveFrames = 1, math.Pow(10, f/10), 1
	}
	return w
}

// bargeinRunWindows reads one run's per-window raw values: the ERLE curve
// plus the correlation source (music: the probe's corr_vs_far against the
// bit-exact far; TTS: the reference-free self_corr_at_period, 0 meaning not
// measured). Windows past the curve's end get empty values.
func bargeinRunWindows(t *testing.T, out59, run string) []bargeinWinVal {
	t.Helper()
	var rep probeRunJSON
	b, err := os.ReadFile(filepath.Join(out59, run+".json"))
	if err != nil {
		t.Fatalf("run report %s: %v", run, err)
	}
	if err := json.Unmarshal(b, &rep); err != nil {
		t.Fatalf("parse %s.json: %v", run, err)
	}
	var curve []bargeinCurvePoint
	cb, err := os.ReadFile(filepath.Join(out59, run+"-curve.json"))
	if err != nil {
		t.Fatalf("curve %s: %v", run, err)
	}
	if err := json.Unmarshal(cb, &curve); err != nil {
		t.Fatalf("parse %s-curve.json: %v", run, err)
	}
	isMusic := rep.FarContent == "music"

	var corr []float64
	if isMusic {
		pb, err := os.ReadFile(filepath.Join(out59, "probe.json"))
		if err != nil {
			t.Fatalf("probe.json: %v", err)
		}
		all := map[string]probeCorrRun{}
		if err := json.Unmarshal(pb, &all); err != nil {
			t.Fatalf("parse probe.json: %v", err)
		}
		for _, w := range all[run].Windows {
			corr = append(corr, w.Corr)
		}
	} else {
		ab, err := os.ReadFile(filepath.Join(out59, "autocorr.json"))
		if err != nil {
			t.Fatalf("autocorr.json: %v", err)
		}
		all := map[string]autocorrRun{}
		if err := json.Unmarshal(ab, &all); err != nil {
			t.Fatalf("parse autocorr.json: %v", err)
		}
		for _, w := range all[run].Windows {
			v := math.NaN()
			if w.SelfCorr > 0 {
				v = w.SelfCorr
			}
			corr = append(corr, v)
		}
	}

	out := make([]bargeinWinVal, len(curve))
	for i, cp := range curve {
		v := bargeinWinVal{ERLE: cp.ERLE, Corr: math.NaN()}
		if i < len(corr) {
			v.Corr = corr[i]
		}
		out[i] = v
	}
	return out
}

// bargeinMaterialWindows slices a run's global-window values onto a
// material's own window axis: material window i reads global window
// offset+i; where nothing was measured (past the replayed span, or no axis
// at all) the window is empty - ERLE inactive, correlation NaN - and the
// gate's channels skip it (the nominal cut applies).
func bargeinMaterialWindows(global []bargeinWinVal, offset, durSec int) []bargeinWinVal {
	n := durSec + 1 // a material of D seconds spans windows 0..floor(D)
	out := make([]bargeinWinVal, n)
	if offset < 0 {
		return out
	}
	for i := range out {
		g := offset + i
		if g < len(global) {
			out[i] = global[g]
		}
	}
	return out
}

// probeEvents is the slice of probe.json the event segmentation comes from:
// the gate materials were cut on the probe's event spans, so the material
// axis must use them. The autocorr probe's own event list segments the same
// curve differently (merge and classification rules differ) - its per-window
// self-correlations are used as values, never its event list.
type probeEvents struct {
	Events [][2]int `json:"events"`
}

// bargeinEventSpan reads one event's [start, end) window span from
// probe.json.
func bargeinEventSpan(t *testing.T, out59, run string, k int) [2]int {
	t.Helper()
	pb, err := os.ReadFile(filepath.Join(out59, "probe.json"))
	if err != nil {
		t.Fatalf("probe.json: %v", err)
	}
	all := map[string]probeEvents{}
	if err := json.Unmarshal(pb, &all); err != nil {
		t.Fatalf("parse probe.json: %v", err)
	}
	evs := all[run].Events
	if k >= len(evs) {
		t.Fatalf("%s event %d: only %d events on record", run, k, len(evs))
	}
	return evs[k]
}

// bargeinReadings is the readings report (the JSON artifact).
type bargeinReadings struct {
	Matrix    []bargeinCell     `json:"matrix"`
	Events    []bargeinEventRow `json:"tier3_events"`
	Controls  map[string]any    `json:"controls"`
	Untrusted map[string]int    `json:"untrusted_windows_per_run"`
}

type bargeinCell struct {
	Axis           float64 `json:"axis_dbfs"`
	Dynamic        bool    `json:"dynamic"`
	GoodEvents     int     `json:"good_events"`
	CollapseEvents int     `json:"collapse_events"`
	GoodDur        float64 `json:"good_dur_sec"`
	CollapseDur    float64 `json:"collapse_dur_sec"`
	GoodRate       float64 `json:"good_per_10min"`
	CollapseRate   float64 `json:"collapse_per_10min"`
}

type bargeinEventRow struct {
	Material string  `json:"material"`
	State    string  `json:"state"`
	Axis     float64 `json:"axis_dbfs"`
	Dynamic  bool    `json:"dynamic"`
	Window   int     `json:"material_window"`
	Untrust  bool    `json:"untrusted"`
	RmsDbfs  float64 `json:"rms_dbfs"`
	AsrText  string  `json:"asr_text"`
}

// TestSpike389r1BargeinReadings is the leg's readings driver. Requires the
// 59-w artifacts in their temp root; skips with the reason recorded when
// they are absent (a four-numbers run on a host without them stays green).
func TestSpike389r1BargeinReadings(t *testing.T) {
	spikeGate(t)
	out59 := longWindowOutDefault
	if _, err := os.Stat(filepath.Join(out59, "gates-probe.json")); err != nil {
		t.Skipf("59-w artifacts not present (%v); the readings need them on disk", err)
	}

	pb, err := os.ReadFile(filepath.Join(out59, "gates-probe.json"))
	if err != nil {
		t.Fatalf("gates-probe.json: %v", err)
	}
	probes := probeGates{}
	if err := json.Unmarshal(pb, &probes); err != nil {
		t.Fatalf("parse gates-probe.json: %v", err)
	}
	gb, err := os.ReadFile(filepath.Join(out59, "gates.json"))
	if err != nil {
		t.Fatalf("gates.json: %v", err)
	}
	gates := probeGates{}
	if err := json.Unmarshal(gb, &gates); err != nil {
		t.Fatalf("parse gates.json: %v", err)
	}

	// Raw per-run window values, read once.
	runVals := map[string][]bargeinWinVal{}
	vals := func(run string) []bargeinWinVal {
		if v, ok := runVals[run]; ok {
			return v
		}
		v := bargeinRunWindows(t, out59, run)
		runVals[run] = v
		return v
	}

	names := make([]string, 0, len(probes))
	for n := range probes {
		names = append(names, n)
	}
	sort.Strings(names)

	type matRow struct {
		name string
		mat  probeGateMaterial
		wins []bargeinWinVal
	}
	var rows []matRow
	goodDur, collapseDur := 0.0, 0.0
	goodN, collapseN := 0, 0
	for _, n := range names {
		m := probes[n]
		ref := bargeinParseMaterial(n)
		row := matRow{name: n, mat: m}
		switch ref.kind {
		case "steady", "whole":
			row.wins = bargeinMaterialWindows(vals(ref.run), ref.offset, int(m.DurationSec))
		case "event":
			ev := bargeinEventSpan(t, out59, ref.run, ref.event)
			row.wins = bargeinMaterialWindows(vals(ref.run), ev[0], int(m.DurationSec))
		case "archived":
			// No window axis exists for the archived wav; the material has
			// no episodes, so every window reads empty (nominal cut).
		}
		// The state split follows the 59-w table's numbers exactly: good =
		// 9 materials/222.0s, collapse = 44/104.4s (43 event windows + the
		// archived wav). The state literally named "collapsed" is the
		// cold-tts-2 run - the table's own arithmetic (44 = 43+1) leaves it
		// out of the collapse denominator, and the 59-v3 recheck uses the
		// same numbers, so the re-cut does too (it carries no episodes;
		// recorded here so the exclusion is deliberate, not lost).
		switch m.State {
		case "good":
			goodDur += m.DurationSec
			goodN++
		case "collapse":
			collapseDur += m.DurationSec
			collapseN++
		}
		rows = append(rows, row)
	}
	// Duration and material-count checks against the 59-w table.
	if math.Abs(goodDur-222.0) > 0.5 || math.Abs(collapseDur-104.4) > 0.5 {
		t.Fatalf("material durations drifted from the 59-w set: good=%.1f collapse=%.1f", goodDur, collapseDur)
	}
	if goodN != 9 || collapseN != 44 {
		t.Fatalf("material counts drifted from the 59-w set: good=%d collapse=%d, want 9/44", goodN, collapseN)
	}

	// The controls: the base material is named in the source string the rig
	// leg wrote ("long-music-60-steady + TTS speech at -26 dBFS RMS").
	type controlRow struct {
		name string
		mat  probeGateMaterial
		base string
		wins []bargeinWinVal
	}
	var controls []controlRow
	for _, n := range []string{"control-26", "control-36"} {
		m, ok := gates[n]
		if !ok {
			t.Fatalf("control material %s missing from gates.json", n)
		}
		base, _, _ := strings.Cut(m.Source, " + TTS speech")
		if _, ok := probes[base]; !ok {
			t.Fatalf("control %s base material %q not in the probe set", n, base)
		}
		ref := bargeinParseMaterial(base)
		var wv []bargeinWinVal
		if ref.kind == "steady" {
			wv = bargeinMaterialWindows(vals(ref.run), ref.offset, int(m.DurationSec))
		}
		controls = append(controls, controlRow{name: n, mat: m, base: base, wins: wv})
	}

	readings := bargeinReadings{
		Controls:  map[string]any{},
		Untrusted: map[string]int{},
	}
	// The per-run untrusted counts, replayed once through the product gate
	// with the loose axis (the flag surface the gate would see live).
	for run, wv := range runVals {
		g := audio.NewBargeinGate(audio.BargeinOptions{DynamicGating: true})
		for i, v := range wv {
			g.ObserveWindow(bargeinWindowStat(i, v.ERLE), v.Corr)
		}
		n := 0
		for _, u := range g.UntrustedWindows() {
			if u {
				n++
			}
		}
		readings.Untrusted[run] = n
	}

	// The four cells of the readings matrix: axis x dynamic.
	for _, axis := range []float64{-45, -30} {
		for _, dyn := range []bool{true, false} {
			cell := bargeinCell{Axis: axis, Dynamic: dyn, GoodDur: goodDur, CollapseDur: collapseDur}
			for _, row := range rows {
				g := audio.NewBargeinGate(audio.BargeinOptions{
					EnergyThresholdDbfs: axis,
					TightenedEnergyDbfs: audio.BargeinTightenedEnergyDbfs,
					DynamicGating:       dyn,
				})
				for i, v := range row.wins {
					g.ObserveWindow(bargeinWindowStat(i, v.ERLE), v.Corr)
				}
				for _, ep := range row.mat.Episodes {
					d := g.Evaluate(audio.BargeinEpisode{
						StartSample: ep.StartSample, Samples: ep.Samples,
						RmsDbfs: ep.RmsDbfs, AsrText: ep.AsrText,
					})
					if !d.Pass {
						continue
					}
					switch row.mat.State {
					case "collapse":
						cell.CollapseEvents++
					case "good":
						cell.GoodEvents++
					}
					readings.Events = append(readings.Events, bargeinEventRow{
						Material: row.mat.Material, State: row.mat.State, Axis: axis, Dynamic: dyn,
						Window:  ep.StartSample / audio.BargeinWindowSamples,
						Untrust: d.TightenActive, RmsDbfs: ep.RmsDbfs, AsrText: ep.AsrText,
					})
				}
			}
			cell.GoodRate = Rate10Min(cell.GoodEvents, cell.GoodDur)
			cell.CollapseRate = Rate10Min(cell.CollapseEvents, cell.CollapseDur)
			readings.Matrix = append(readings.Matrix, cell)
			t.Logf("[389r1 axis=%.0f dynamic=%-5v] good %d/%.1fs = %.2f/10min | collapse %d/%.1fs = %.2f/10min",
				axis, dyn, cell.GoodEvents, cell.GoodDur, cell.GoodRate, cell.CollapseEvents, cell.CollapseDur, cell.CollapseRate)
		}
	}

	// The controls through the same decision path (base material windows).
	for _, c := range controls {
		for _, axis := range []float64{-45, -30} {
			g := audio.NewBargeinGate(audio.BargeinOptions{
				EnergyThresholdDbfs: axis,
				TightenedEnergyDbfs: audio.BargeinTightenedEnergyDbfs,
				DynamicGating:       true,
			})
			for i, v := range c.wins {
				g.ObserveWindow(bargeinWindowStat(i, v.ERLE), v.Corr)
			}
			for _, ep := range c.mat.Episodes {
				d := g.Evaluate(audio.BargeinEpisode{
					StartSample: ep.StartSample, Samples: ep.Samples,
					RmsDbfs: ep.RmsDbfs, AsrText: ep.AsrText,
				})
				key := fmt.Sprintf("%s@%.0f", c.name, axis)
				readings.Controls[key] = map[string]any{
					"rms_dbfs": ep.RmsDbfs, "pass": d.Pass, "tighten_active": d.TightenActive,
					"threshold": d.ThresholdDbfs, "suppressed_by_dynamic": d.SuppressedByDynamic,
				}
				t.Logf("[389r1 control %s axis=%.0f] rms=%.2f pass=%v tighten=%v threshold=%.0f",
					c.name, axis, ep.RmsDbfs, d.Pass, d.TightenActive, d.ThresholdDbfs)
				// A real-speech control the gate eats is a product-fatal
				// shape: the -26.70dBFS control must pass on both axes,
				// tight or not.
				if c.name == "control-26" && !d.Pass {
					t.Fatalf("control-26 (%.2f dBFS) rejected at axis %.0f (tighten=%v, threshold=%.0f): the gate eats real speech",
						ep.RmsDbfs, axis, d.TightenActive, d.ThresholdDbfs)
				}
			}
		}
	}

	// --- The assertions: relationship cuts, no golden numbers. ---
	// (The readings are written first: a failed assertion must leave the
	// full readings on disk for the re-cut to inspect.)
	if err := os.MkdirAll(bargein389r1Out, 0o755); err != nil {
		t.Fatalf("out dir: %v", err)
	}
	writeJSON(t, filepath.Join(bargein389r1Out, "bargein-readings.json"), readings)

	byCell := map[string]bargeinCell{}
	for _, c := range readings.Matrix {
		byCell[fmt.Sprintf("%.0f/%v", c.Axis, c.Dynamic)] = c
	}
	off45 := byCell["-45/false"]
	on45 := byCell["-45/true"]
	off30 := byCell["-30/false"]
	on30 := byCell["-30/true"]

	// Recompute alignment with the 59-w static baseline: at -45 with the
	// dynamic gating off, the gate IS the static tier chain the rig leg
	// measured - good 2 events (5.4/10min), collapse 2 (11.5/10min).
	if off45.GoodEvents != 2 || off45.CollapseEvents != 2 {
		t.Fatalf("static -45 recompute drifted from the 59-w baseline: good=%d collapse=%d, want 2/2", off45.GoodEvents, off45.CollapseEvents)
	}
	// AC#2 (the product-form reading): with the dynamic gating on, the
	// collapse state drops to the good-state level or below.
	if on45.CollapseEvents > on45.GoodEvents {
		t.Fatalf("dynamic gate did not pull the collapse state to the good level: collapse=%d > good=%d", on45.CollapseEvents, on45.GoodEvents)
	}
	// The dynamic gate actually eats collapse-state triggers (it is not a
	// no-op): strictly fewer than the static chain on the same material.
	if on45.CollapseEvents >= off45.CollapseEvents {
		t.Fatalf("dynamic gate is a no-op on the collapse state: on=%d >= off=%d", on45.CollapseEvents, off45.CollapseEvents)
	}
	// The good state does not get worse (the quiet-speech reach the loose
	// nominal buys is not paid for with more good-state triggers).
	if on45.GoodEvents > off45.GoodEvents {
		t.Fatalf("dynamic gate worsened the good state: on=%d > off=%d", on45.GoodEvents, off45.GoodEvents)
	}
	// The -30 axis keeps both states at zero with the gating on or off
	// (zero regression against the static -30 readings).
	if off30.GoodEvents != 0 || off30.CollapseEvents != 0 || on30.GoodEvents != 0 || on30.CollapseEvents != 0 {
		t.Fatalf("-30 axis regressed: off=%+v on=%+v", off30, on30)
	}
}
