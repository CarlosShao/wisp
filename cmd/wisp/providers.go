package main

// wisp providers - the provider save/discovery path, and with it the production
// caller of the capability probe (ruling A11 / ticket 11 AC#6).
//
// Before this file llm.RunProbeSuite had nine green tests and no caller: the
// 「声明 PASS / 实测 FAIL」 event SPEC-05 §3.1 demands could not fire on a real
// machine, because nothing built a provider and asked it what it can actually
// do. Same shape as A8 and A13, and the same fix - the composition root now
// calls it.
//
//	discover <provider>              GET /v1/models through llm.DiscoverModels
//	probe <provider>/<model>         llm.RunProbeSuite -> provider_health + the
//	                                 mismatch notice on stderr

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

const providersUsage = `wisp providers - the provider catalog path (ticket 12, ruling A11)

Usage:
  wisp providers discover <provider>
        ask the provider's /v1/models endpoint what it serves
  wisp providers probe <provider>/<model> [--timeout 30s]
        run the capability probe suite (fc / vision / thinking) against the
        REAL provider, write the verdicts into provider_health and print every
        「声明 PASS / 实测 FAIL」 mismatch

Keys never come from here: both subcommands resolve api_key_ref through the
DPAPI store of the active env, exactly like the text path does.
`

// providersIO is the process surface of one `wisp providers` run.
type providersIO struct {
	stdout  io.Writer
	stderr  io.Writer
	dataDir string
	// now is the clock for the persisted probe record only (a record, never a
	// deadline - D42#9).
	now func() time.Time
}

func (p providersIO) out() io.Writer {
	if p.stdout != nil {
		return p.stdout
	}
	return os.Stdout
}

func (p providersIO) err() io.Writer {
	if p.stderr != nil {
		return p.stderr
	}
	return os.Stderr
}

// cmdProviders runs `wisp providers ...` and returns the exit code.
func cmdProviders(argv []string, io_ providersIO) int {
	if len(argv) == 0 {
		fmt.Fprint(io_.err(), providersUsage)
		return 2
	}
	if io_.now == nil {
		io_.now = func() time.Time { return time.Now().UTC() }
	}
	if io_.dataDir == "" {
		dir, dirErr := resolveDataDir(buildinfo.EnvString())
		if dirErr != nil {
			// Ticket 128 AC#2: no data root, no provider listing. It lands
			// before secret.NewStore on purpose, so a start-up directory is
			// never handed a credential store to create.
			fmt.Fprintf(io_.err(), "wisp providers: %v\n", dirErr)
			return 2
		}
		io_.dataDir = dir
	}
	sub, rest := argv[0], argv[1:]
	st, err := secret.NewStore(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 凭据存储不可用：%v\n", err)
		return 2
	}
	// res=nil: the store stays the only thing that turns a ref into a key.
	cfg, _, err := config.LoadFile(filepath.Join(io_.dataDir, configFileName), nil)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 配置未就绪（Unconfigured）：%v\n", err)
		return 2
	}
	switch sub {
	case "discover":
		return providersDiscover(io_, cfg, st, rest)
	case "probe":
		return providersProbe(io_, cfg, st, rest)
	default:
		fmt.Fprintf(io_.err(), "wisp providers: 未知子命令 %q\n\n", sub)
		fmt.Fprint(io_.err(), providersUsage)
		return 2
	}
}

// providersDiscover lists what the provider's own endpoint reports.
func providersDiscover(io_ providersIO, cfg *config.Config, st *secret.Store, argv []string) int {
	name := firstArg(argv)
	p, ok := cfg.LLM.Providers[name]
	if !ok {
		fmt.Fprintf(io_.err(), "wisp providers: 目录中没有 provider %q\n", name)
		return 2
	}
	key := ""
	if p.APIKeyRef != "" {
		k, err := st.Resolve(p.APIKeyRef)
		if err != nil || k == "" {
			// The Unconfigured path: a ref that does not resolve is never
			// answered with an unauthenticated request.
			fmt.Fprintf(io_.err(), "wisp providers: provider %q 的 api_key_ref 无法解析（Unconfigured）\n", name)
			return 2
		}
		key = k
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	models, err := llm.DiscoverModels(ctx, llm.DiscoverOptions{
		BaseURL: p.BaseURL, APIKey: key,
		HTTPClient: llm.NewDefaultHTTPClient(cfg.Net.Proxy.Mode, cfg.Net.Proxy.URL),
	})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 发现失败：%v\n", err)
		return 1
	}
	fmt.Fprintf(io_.out(), "wisp providers: %s 上报 %d 个模型\n", name, len(models))
	for _, m := range models {
		fmt.Fprintf(io_.out(), "  %s\n", m.ID)
	}
	return 0
}

// providersProbe measures one provider/model pair and records the verdicts.
func providersProbe(io_ providersIO, cfg *config.Config, st *secret.Store, argv []string) int {
	fs := flag.NewFlagSet("wisp providers probe", flag.ContinueOnError)
	fs.SetOutput(io_.err())
	timeout := fs.Duration("timeout", 60*time.Second, "whole-suite budget")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	target := fs.Arg(0)
	provider, model, ok := splitRef(target)
	if !ok {
		fmt.Fprintf(io_.err(), "wisp providers: 参数需为 provider/model，收到 %q\n", target)
		return 2
	}
	// The same resolver the text path uses, pointed at the pair under test:
	// ResolveChain turns "provider/model" into an Endpoint whose key came out
	// of the store.
	res := llm.NewResolver(cfg, st)
	res.TextChain = []string{target}
	eps, err := res.ResolveChain()
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 端点未就绪（Unconfigured）：%v\n", err)
		return 2
	}
	prov, err := llm.BuildEndpointProvider(eps[0], chainOptions(cfg))
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 无法构造 provider：%v\n", err)
		return 2
	}
	store, err := memory.Open(io_.dataDir)
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 存储不可用：%v\n", err)
		return 2
	}
	defer func() { _ = store.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var mismatches []llm.ProbeMismatch
	rep, err := llm.RunProbeSuite(ctx, prov, llm.ProbeSuiteOptions{
		Provider: provider,
		Model:    model,
		Declared: cfg.LLM.Providers[provider].Models[model].Capabilities,
		Sink:     storeHealthSink{store: store},
		Mismatch: func(m llm.ProbeMismatch) {
			mismatches = append(mismatches, m)
			// SPEC-05 §3.1: the contradiction is announced where the user is.
			fmt.Fprintf(io_.err(), "wisp providers: %s（%s）\n", m.Label(), m.Detail)
		},
		Now: io_.now,
	})
	if err != nil {
		fmt.Fprintf(io_.err(), "wisp providers: 探针未完成：%v\n", err)
		return 1
	}
	fmt.Fprintf(io_.out(), "wisp providers: %s/%s 实测（p50 %dms，写入 %v）\n",
		rep.Provider, rep.Model, rep.LatencyMS, rep.Written)
	for _, o := range rep.Outcomes {
		// AC#3 (ticket 67 / A23): the verdict column is ASCII PASS/FAIL, not
		// U+2713/U+2717 - the console font does not guarantee those codepoints,
		// and both fall inside the D22 ban #8 range (U+2190-U+2BFF).
		verdict := "FAIL"
		if o.Result.OK {
			verdict = "PASS"
		}
		declared := "未声明"
		if o.Declared {
			declared = "声明"
		}
		fmt.Fprintf(io_.out(), "  %s %s 实测 %s\n", o.Capability, declared, verdict)
		if o.Result.Detail != "" && !o.Result.OK {
			fmt.Fprintf(io_.out(), "    细节：%s\n", oneLine(o.Result.Detail))
		}
	}
	if len(mismatches) > 0 {
		fmt.Fprintf(io_.out(), "wisp providers: %d 项声明与实测不符\n", len(mismatches))
	}
	return 0
}

// splitRef splits "provider/model" at the first slash (model ids keep their
// slashes, so this mirrors llm.Resolver's own chain-element rule).
func splitRef(s string) (provider, model string, ok bool) {
	i := strings.Index(s, "/")
	if i <= 0 || i == len(s)-1 {
		return "", "", false
	}
	return s[:i], s[i+1:], true
}

// oneLine keeps a probe detail on a single output line without hiding it.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
