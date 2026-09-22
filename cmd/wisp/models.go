package main

// wisp models - the C29 model store's operator surface, and the place ticket 121
// AC#2 put the model hand-off chain into the product.
//
// WHY THIS FILE EXISTS (R-109-4, the sixth case of "built but nobody calls it").
// Before this commit `go list -deps ./cmd/wisp` did not contain
// github.com/CarlosShao/wisp/internal/models on any platform. That made three
// things that are green inside their own package unreachable from the binary at
// once: models.Manager, models.VerifyInstalled (ticket 109's hand-off
// re-verification), and DownloadingBridge.Run - which is the ONLY caller of
// VerifyInstalled anywhere outside tests (internal/models/bridge.go:58). The
// reading is reproducible; it is recorded in the ticket face.
//
// WHAT THE PRODUCTION TRIGGER ACTUALLY IS, stated verbatim from the spec rather
// than invented here. SPEC-04 §7.1: "首次使用某能力时下载对应模型" - the
// download is specced to start when a capability is first used. That capability
// does not exist in this tree: internal/speech is a doc.go-only boundary stub
// (DEFERRED: ticket 15/26/41 per its own doc.go) and there is no
// internal/engines/ directory, so nothing loads model bytes into memory and
// "first use of a capability" has no code behind it yet. `wisp models ensure`
// is therefore the hand-off point that DOES exist - a real process, a real
// signature-verified manifest, a real re-verification - and it is deliberately
// not dressed up as the engine's call site. Whichever ticket lands
// internal/speech owns adding its own handOffModel call; this file cannot
// become that reader by wishing, and AC#5 below says what that leaves open.
//
// AC#5, THE RESIDUAL WINDOW THIS WIRING DOES NOT CLOSE (ticket 109's R-109-1,
// not "fixed", not glossed over). handOffModel returns *after* VerifyInstalled
// passed. Anything that opens those bytes afterwards is still separated from
// this return by a span of exactly the shape ticket 109 measured: the install
// directory stays inherit-wide on purpose (ticket 95's ruling, defended by
// internal/models/no_seal_ruling_windows_test.go), so every principal the
// parent grants write to owns that span, and a swap landing after the return is
// invisible to this process. What is NOT true after this wiring: that the
// window is *reachable*. There is still no reader of model bytes anywhere in
// `cmd/wisp`'s dependency graph - this process never opens a file under
// <store>/<id>/, it only re-hashes it - so the segment has no end point in the
// product today. It gets an end point the moment an engine loader is wired, and
// closing it then is that ticket's work, not a sentence in this comment.
//
// The dependency direction stays what this package has always held: cmd/wisp may
// know every layer at once, layers know nothing of each other. models and
// statemachine import nothing back up into here.

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/models"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

const modelsUsage = `wisp models - the signed model store (C29, SPEC-04 §7.2)

Usage:
  wisp models list            show what the signature-verified manifest carries
                              (id, purpose, license, size, status)
  wisp models verify <id>     re-check the installed files of one model against
                              the signed manifest, without downloading anything
  wisp models ensure <id>     make one model available: download or reuse it,
                              then re-verify it at the hand-off point (ticket
                              109's guard) and announce it through D43 row #37

Store root: [models] dir, relative to the data dir, default <data dir>\models.
Mirrors: [models] mirror serve bytes only; the signed manifest is the authority.
Manifest: WISP_MODELS_MANIFEST, else <exe dir>\models, else <cwd>\models.
`

// modelsIO is `wisp models`' process surface, injectable for the composition
// tests the same way providersIO is.
type modelsIO struct {
	stdout  io.Writer
	stderr  io.Writer
	dataDir string
}

func (m modelsIO) out() io.Writer {
	if m.stdout != nil {
		return m.stdout
	}
	return os.Stdout
}

func (m modelsIO) err() io.Writer {
	if m.stderr != nil {
		return m.stderr
	}
	return os.Stderr
}

// cmdModels runs `wisp models ...` and returns the exit code. The codes are the
// ones `wisp run` established (SPEC-05 §3.4's classified failure, ticket 12):
// 0 the requested thing happened, 2 Unconfigured (no manifest, no readable
// config, unknown model id - a setup problem, not an outcome), 1 anything else,
// including a refused hand-off, which is the security verdict this command
// exists to deliver and must not be softened into a warning.
func cmdModels(argv []string, io_ modelsIO) int {
	if len(argv) == 0 {
		fmt.Fprint(io_.err(), modelsUsage)
		return 2
	}
	if io_.dataDir == "" {
		io_.dataDir = resolveDataDir(buildinfo.EnvString())
	}
	sub, rest := argv[0], argv[1:]
	switch sub {
	case "list":
		return modelsList(io_, rest)
	case "verify":
		return modelsVerify(io_, rest)
	case "ensure":
		return modelsEnsure(io_, rest)
	default:
		fmt.Fprintf(io_.err(), "wisp models: 未知子命令 %q\n\n", sub)
		fmt.Fprint(io_.err(), modelsUsage)
		return 2
	}
}

// modelStore is the assembled hand-off chain for one process: the signature-
// verified manifest, the Manager over the store root this build's config names,
// and the manifest path it came from (so a refusal can name its own source).
type modelStore struct {
	mgr      *models.Manager
	manifest *models.Manifest
	root     string
	src      string
}

// openModelStore is the assembly half of this command, and the reason
// internal/models is in this binary's dependency graph at all. C29 has no off
// switch here on purpose: VerifySignature is passed as the literal true, because
// models.NewManager rejects a false (internal/models/downloader.go:104) and
// config.validateModels rejects it even earlier, so a field that can only be
// true must not be threaded through as if it could be false.
func openModelStore(dataDir string, cfg *config.Config) (*modelStore, error) {
	manifestPath, sigPath, err := models.ResolveManifestPath("")
	if err != nil {
		return nil, err
	}
	manifest, err := models.LoadSignedManifest(manifestPath, sigPath, buildinfo.MinisignPublicKey)
	if err != nil {
		return nil, err
	}
	var mirrors []string
	var overrides map[string]string
	if cfg != nil {
		mirrors = cfg.Models.Mirror
		overrides = cfg.Models.LocalOverride
	}
	root := modelStoreDir(dataDir, cfg)
	mgr, err := models.NewManager(models.Options{
		DataDir:         root,
		Manifest:        manifest,
		VerifySignature: true,
		Mirrors:         mirrors,
		LocalOverride:   overrides,
	})
	if err != nil {
		return nil, err
	}
	return &modelStore{mgr: mgr, manifest: manifest, root: root, src: manifestPath}, nil
}

// loadModelConfig reads [models] the same way every other command reads its
// section: one LoadFile, res=nil (the model store resolves no api_key_ref, so
// there is nothing for a SecretResolver to turn into a key here).
func loadModelConfig(dataDir string) (*config.Config, error) {
	cfg, _, err := config.LoadFile(filepath.Join(dataDir, configFileName), nil)
	return cfg, err
}

// modelStoreDir implements SPEC-04 §7.2's install layout ("落 models\<id>\") on
// top of SPEC-03 §5.2's data dir. An empty [models] dir is not a bug to guess
// around, it is the documented default: <data dir>\models.
func modelStoreDir(dataDir string, cfg *config.Config) string {
	if cfg != nil && cfg.Models.Dir != "" {
		if filepath.IsAbs(cfg.Models.Dir) {
			return cfg.Models.Dir
		}
		return filepath.Join(dataDir, cfg.Models.Dir)
	}
	return filepath.Join(dataDir, "models")
}

// modelsList prints the manifest. It is the cheapest leg and still not nothing:
// a listing that succeeds means LoadSignedManifest verified the minisign
// signature against the key embedded in buildinfo BEFORE any URL in the file
// was trusted, which is C29's whole claim.
func modelsList(io_ modelsIO, _ []string) int {
	// cfg=nil: listing is a question about the manifest, not about this boot's
	// [models] section, so it must not be answered by a config file that may not
	// exist yet. The store path printed is then the documented default.
	store, err := openModelStore(io_.dataDir, nil)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp models list: %v\n", err)
		return 2
	}
	fmt.Fprintf(io_.out(), "wisp models: 已验签清单 %s（%d 个模型，store=%s）\n",
		store.src, len(store.manifest.Models), store.root)
	for _, e := range store.manifest.Models {
		status := e.Status
		if status == "" {
			status = "ok"
		}
		fmt.Fprintf(io_.out(), "  %-40s purpose=%-6s %-14s %10d B license=%s status=%s\n",
			e.ID, e.Purpose, e.Quant, e.SizeBytes, e.License, status)
	}
	return 0
}

// modelsVerify answers the owner's question directly: "is what is installed
// right now still what the signed manifest pinned?" It never touches the
// network and never installs anything, so it is safe to run against a store an
// engine is using, and it is the leg that stays useful after ticket 109's
// guard - a swap that landed between two runs is caught by this one.
func modelsVerify(io_ modelsIO, argv []string) int {
	id := firstArg(argv)
	if id == "" {
		fmt.Fprint(io_.err(), "wisp models verify: 需要模型 id（wisp models list 看清单）\n")
		return 2
	}
	cfg, err := loadModelConfig(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp models verify: 配置未就绪（Unconfigured）：%v\n", err)
		return 2
	}
	store, err := openModelStore(io_.dataDir, cfg)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp models verify: %v\n", err)
		return 2
	}
	started := time.Now()
	if err := store.mgr.VerifyInstalled(id); err != nil {
		fmt.Fprintf(io_.err(), "wisp models verify: %v\n", err)
		return 1
	}
	fmt.Fprintf(io_.out(), "wisp models verify: %s 与已验签清单一致（store=%s，%.2fs）\n",
		id, store.root, time.Since(started).Seconds())
	return 0
}

// modelsEnsure performs the full walk: download-or-reuse, then the hand-off
// re-verification, then the D43 announcement - in that order, in one call to
// DownloadingBridge.Run, which is the only production path to
// Manager.VerifyInstalled.
func modelsEnsure(io_ modelsIO, argv []string) int {
	id := firstArg(argv)
	if id == "" {
		fmt.Fprint(io_.err(), "wisp models ensure: 需要模型 id（wisp models list 看清单）\n")
		return 2
	}
	cfg, err := loadModelConfig(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp models ensure: 配置未就绪（Unconfigured）：%v\n", err)
		return 2
	}
	store, err := openModelStore(io_.dataDir, cfg)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp models ensure: %v\n", err)
		return 2
	}
	// The persistent sink is installed here for the same reason run.go installs
	// it before assembling (ticket 117): the resident path has no terminal, and
	// a refused hand-off nobody records is a security verdict that leaves no
	// trace. A store that cannot open a log dir degrades loudly, it does not
	// stop the verification.
	var logf func(string, ...any)
	if sink, sinkErr := installLogSink(io_.dataDir); sinkErr != nil {
		fmt.Fprintf(io_.err(), "wisp models: 持久日志未启用（%v）：本轮告警只会到 stderr\n", sinkErr)
	} else {
		defer sink.close()
		sinkLogger := sink.logger()
		logf = func(format string, args ...any) {
			sinkLogger.Info("models: " + fmt.Sprintf(format, args...))
		}
	}
	return store.handOffModel(io_, id, logf)
}

// handOffModel runs one Downloading walk (D43 row #2 in, row #37 out) for id.
//
// This function is the grep target of ticket 121 AC#2: it is the non-test
// caller of models.WireDownloading and therefore of DownloadingBridge.Run,
// whose body holds the only non-test call of Manager.VerifyInstalled
// (internal/models/bridge.go:58). logf receives progress ticks and may be nil.
func (ms *modelStore) handOffModel(io_ modelsIO, id string, logf func(string, ...any)) int {
	machine := statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
	defer machine.Close()
	bridge := models.WireDownloading(ms.mgr, machine)

	ctx, cancel := context.WithTimeout(context.Background(), modelHandoffTimeout)
	defer cancel()

	started := time.Now()
	state, err := bridge.Run(ctx, id)
	if err != nil {
		// Loud and classified: a refused hand-off is reported as the verdict it
		// is, with the state the walk ended in, because "Error" and "still
		// Downloading" are two different failures (the second means the walk
		// could not even be exited, which is a table problem, not a model one).
		fmt.Fprintf(io_.err(), "wisp models ensure: 交还被拒绝（state=%s）：%v\n", state, err)
		if logf != nil {
			logf("hand-off REFUSED id=%s state=%s err=%v", id, state, err)
		}
		return 1
	}
	fmt.Fprintf(io_.out(), "wisp models ensure: %s 已交还（state=%s，ticks=%d，末次进度 %d%%，%.2fs，store=%s）\n",
		id, state, bridge.Ticks(), int(bridge.LastPercent()), time.Since(started).Seconds(), ms.root)
	if logf != nil {
		logf("hand-off OK id=%s state=%s root=%s", id, state, ms.root)
	}
	return 0
}

// modelHandoffTimeout is a wall-clock bound on one walk, not a latency target:
// the largest S1 model in the manifest is ~300 MB, so an hour is generous and
// the number that matters is "it cannot hang forever".
const modelHandoffTimeout = time.Hour
