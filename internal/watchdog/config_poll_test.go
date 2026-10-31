package watchdog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// AC#5, the shared test: ONE planted-edit sequence driven through TWO
// drivers, and the reload events must be identical.
//
//   - driver "05": config.Manager.CheckAndReload called directly, the shape
//     ticket 05's own manager tests froze (the idempotent mtime+size poll,
//     the D36 tier semantics, the OnReload / OnRestartPending events);
//   - driver "42": the same manager handed to the watchdog loop, whose tick
//     carries the config duty (SPEC-03 §4.3: one loop, both duties).
//
// The shared body below runs the same three edits through both drivers and
// applies the same assertions to both, then DeepEqual-compares the two
// records. A watchdog tick that fires different events than the bare poll
// fails here, not in production.

const reloadBase = `
schema_version = 2

[ball]
size = 56

[app]
theme = "dark"
language = "zh-CN"

[voice]
enabled = true

[voice.asr]
provider = "local-sherpa"
model = "parakeet-zh"

[fs]
allowed_dirs = []
`

// reloadRecord is one driver's outcome: the event sequences plus the final
// live values the tiers decide (hot applied, reload applied, restart kept).
type reloadRecord struct {
	Reloads     [][]string
	Restarts    [][]string
	BallSize    int
	VoiceModel  string
	AppLanguage string
}

// reloadRecorder collects the manager's events for one scenario.
type reloadRecorder struct {
	mu      sync.Mutex
	reload  [][]string
	restart [][]string
}

func (r *reloadRecorder) hookOnReload(sections []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reload = append(r.reload, sections)
}

func (r *reloadRecorder) hookOnRestart(sections []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.restart = append(r.restart, sections)
}

func (r *reloadRecorder) reloads() [][]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.reload)
}

func (r *reloadRecorder) restarts() [][]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.restart)
}

// reloadDriver is one side of the shared test: drive applies one planted
// edit (the manager's poll semantics differ per driver) and after runs the
// driver's own audit once the scenario is done. plant writes the step's
// content again - the concurrent poller can lose one write (a transient
// sharing violation during the rename makes CheckAndReload adopt the stat on
// its failure path, the 05 contract), and a real writer would simply write
// again, so the 42 driver replants until the edit lands.
type reloadDriver interface {
	name() string
	drive(t *testing.T, m *config.Manager, plant func(), applied func() bool)
	after(t *testing.T)
}

// barePollDriver05 is the ticket 05 shape: one direct CheckAndReload call
// per edit.
type barePollDriver05 struct{}

func (barePollDriver05) name() string { return "05" }

func (barePollDriver05) drive(t *testing.T, m *config.Manager, plant func(), applied func() bool) {
	t.Helper()
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("[05] reload: %v", err)
	}
	if rep == nil {
		t.Fatalf("[05] the planted edit did not register: CheckAndReload reported no change")
	}
	if !applied() {
		t.Fatalf("[05] the poll ran but the edit did not apply: %+v", rep)
	}
}

func (barePollDriver05) after(t *testing.T) {}

// watchdogTickDriver42 is the ticket 42 shape: the loop carries the manager
// on its tick, so drive only WAITS until the edit applied. Its scenario
// setup also gives the loop resource duty work, so after() can prove one
// loop carried both duties.
type watchdogTickDriver42 struct {
	reg     *observe.Registry
	tree    *fakeTree
	logBuf  *bytes.Buffer
	log     *slog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	started bool
}

func newWatchdogTickDriver42() *watchdogTickDriver42 {
	log, buf := captureLogger()
	ctx, cancel := context.WithCancel(context.Background())
	tree := &fakeTree{}
	tree.set(overCapTree(400 << 20)) // the resource duty always has work here
	return &watchdogTickDriver42{
		reg: observe.NewRegistry(), tree: tree,
		logBuf: buf, log: log, ctx: ctx, cancel: cancel,
	}
}

func (d *watchdogTickDriver42) name() string { return "42" }

func (d *watchdogTickDriver42) drive(t *testing.T, m *config.Manager, plant func(), applied func() bool) {
	t.Helper()
	if !d.started {
		d.started = true
		l := New(Options{
			TreeReader:   d.tree,
			Registry:     d.reg,
			Config:       m,
			Log:          d.log,
			TickEvery:    2 * time.Millisecond,
			SustainTicks: 1,
		})
		l.Register(ComponentBM25IndexCache, &fakeComponent{dropped: true, detail: "evicted"})
		l.Start(d.ctx)
		// One named goroutine, while both of its duties are live.
		deadline := observe.NewTimeout(10 * time.Second)
		for d.reg.CountByName("watchdog") != 1 {
			if deadline.Expired() {
				t.Fatalf("[42] the watchdog goroutine never appeared in the registry")
			}
			time.Sleep(time.Millisecond)
		}
	}
	// Await the edit, replanting periodically: one write can be swallowed by
	// the manager's failed-load adoption (a transient sharing violation
	// during the rename, read by the concurrent poller), and a real writer
	// would write again. The observable guards the replant, so an edit that
	// already fired its events is never replanted (no double events).
	deadline := observe.NewTimeout(10 * time.Second) // test-side ceiling, not product logic
	lastPlant := time.Now()
	for !applied() {
		if deadline.Expired() {
			t.Fatalf("[42] the watchdog tick never applied the edit; log:\n%s", d.logBuf.String())
		}
		if time.Since(lastPlant) >= 500*time.Millisecond {
			plant()
			lastPlant = time.Now()
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func (d *watchdogTickDriver42) after(t *testing.T) {
	t.Helper()
	d.cancel()
	select {
	case <-d.ctx.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("[42] the loop did not exit after cancel")
	}
	if !strings.Contains(d.logBuf.String(), "config poll applied changes") {
		t.Errorf("[42] the config duty never logged an applied change; log:\n%s", d.logBuf.String())
	}
	if !strings.Contains(d.logBuf.String(), "unload action") {
		t.Errorf("[42] the resource duty never ran on the loop; log:\n%s", d.logBuf.String())
	}
}

// driveReloadScenario is the shared body. Three edits, one per D36 tier:
// [voice] model swap (reload: apply + OnReload), [ball] size (hot: apply,
// no event), [app] language (restart: keep old value + OnRestartPending).
func driveReloadScenario(t *testing.T, d reloadDriver) reloadRecord {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(reloadBase), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("initial load: %v", err)
	}
	rec := &reloadRecorder{}
	m.OnReload = rec.hookOnReload
	m.OnRestartPending = rec.hookOnRestart

	// The 10ms sleep mirrors the config package's own mutator: it keeps the
	// mtime+size fingerprint unambiguous across writes. The write itself is
	// temp+rename (the C19 shape): the concurrent 2ms poller must never read
	// a torn file - a torn read makes CheckAndReload adopt the stat on its
	// failure path, and that write would then be invisible forever.
	step := 0
	mutate := func(next string) {
		step++
		time.Sleep(10 * time.Millisecond)
		tmp := filepath.Join(dir, fmt.Sprintf("config.toml.next-%d", step))
		if err := os.WriteFile(tmp, []byte(next), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(tmp, path); err != nil {
			t.Fatal(err)
		}
	}

	// Chained edits: each step builds on the previous content, so every
	// write changes exactly one thing (an edit built from the original base
	// would silently revert the earlier ones and double-fire their tiers).
	voiceSwap := strings.Replace(reloadBase, `"parakeet-zh"`, `"paraformer-zh"`, 1)
	ballEdit := strings.Replace(voiceSwap, "size = 56", "size = 60", 1)
	appEdit := strings.Replace(ballEdit, `"zh-CN"`, `"en-US"`, 1)

	steps := []struct {
		what    string
		next    string
		applied func(rec *reloadRecorder, cfg *config.Config) bool
	}{
		{
			what: "reload-tier [voice] model swap",
			next: voiceSwap,
			applied: func(rec *reloadRecorder, cfg *config.Config) bool {
				return len(rec.reloads()) == 1
			},
		},
		{
			what: "hot-tier [ball] size",
			next: ballEdit,
			applied: func(rec *reloadRecorder, cfg *config.Config) bool {
				return cfg.Ball.Size == 60
			},
		},
		{
			what: "restart-tier [app] language",
			next: appEdit,
			applied: func(rec *reloadRecorder, cfg *config.Config) bool {
				return len(rec.restarts()) == 1
			},
		},
	}
	for _, s := range steps {
		mutate(s.next)
		d.drive(t, m, func() { mutate(s.next) }, func() bool { return s.applied(rec, m.Config()) })
	}
	d.after(t)

	cfg := m.Config()
	return reloadRecord{
		Reloads:     rec.reloads(),
		Restarts:    rec.restarts(),
		BallSize:    cfg.Ball.Size,
		VoiceModel:  cfg.Voice.ASR.Model,
		AppLanguage: cfg.App.Language,
	}
}

// assertReloadContract is the shared assertion set: both drivers must pass
// the same D36 tier semantics with the same events.
func assertReloadContract(t *testing.T, d reloadDriver, got reloadRecord) {
	t.Helper()
	if len(got.Reloads) != 1 || !reflect.DeepEqual(got.Reloads[0], []string{"voice"}) {
		t.Errorf("[%s] OnReload = %v, want exactly [[voice]]", d.name(), got.Reloads)
	}
	if len(got.Restarts) != 1 || !reflect.DeepEqual(got.Restarts[0], []string{"app"}) {
		t.Errorf("[%s] OnRestartPending = %v, want exactly [[app]]", d.name(), got.Restarts)
	}
	if got.BallSize != 60 {
		t.Errorf("[%s] live ball.size = %d, want 60 (hot tier applies immediately)", d.name(), got.BallSize)
	}
	if got.VoiceModel != "paraformer-zh" {
		t.Errorf("[%s] live voice model = %q, want the planted value (reload tier applies + fires)", d.name(), got.VoiceModel)
	}
	if got.AppLanguage != "zh-CN" {
		t.Errorf("[%s] live app.language = %q, want the old value (restart tier defers to restart)", d.name(), got.AppLanguage)
	}
}

// TestConfigMtimePollFiresReloadEventsIdenticalTo05Contract is AC#5: the
// same edits through both drivers, the same assertions, then one DeepEqual
// over the two records.
func TestConfigMtimePollFiresReloadEventsIdenticalTo05Contract(t *testing.T) {
	d05 := barePollDriver05{}
	got05 := driveReloadScenario(t, d05)
	assertReloadContract(t, d05, got05)

	d42 := newWatchdogTickDriver42()
	got42 := driveReloadScenario(t, d42)
	assertReloadContract(t, d42, got42)

	if !reflect.DeepEqual(got05, got42) {
		t.Fatalf("watchdog-tick reload events differ from the 05 bare-poll contract:\n05: %+v\n42: %+v", got05, got42)
	}
}
