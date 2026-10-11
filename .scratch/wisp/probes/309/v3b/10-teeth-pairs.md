# 309-v3b - AC#2 teeth pairs, label census and false-positive shot (my own readings)

- Leg: 309-v3b (non-implementer). Anchor of this leg: 3b8066a2. HEAD at rig time: afd954ed (sibling commits only; the two files under judgement are byte-identical to a1f84d39, see 20-cells).
- Rig trees: out-of-repo only, one fresh mktemp -d per tree, export = `git archive 3b8066a2 go.mod go.sum internal cmd tools scripts | tar -x` (six pieces; missing scripts/ would fake a walk red). go1.27.1 windows/amd64.
- Runner: `GOFLAGS= go test ./internal/tools/ -run '^Test309ToolRequestProductionConstructionNail$' -count=1`.
- Every red below carries the exact line that produced it: the go-test frame `toolrequest_nail_309_test.go:390:` (the t.Errorf call site inside Test309ToolRequestProductionConstructionNail). Void-token ruler per run (`walk |GetFileAttributesEx|build failed|no test files`): count 0 for all three injected runs, stderr landed (0 bytes) at logs/void-<tag>.stderr. Restore proofs: injected file moved aside (`mv`, never deleted), `sha256sum` of injected == sha256sum of aside file, `cmp` of aside file vs staged copy rc=0, and whole-tree fingerprint (all files hashed, then hash of the sorted hash list) pre == post with cmp rc=0.

## Pair 1 - `new`

- Inject: `internal/tools/z309b_new.go`, body `func z309bRigNew() *agent.ToolRequest { return new(agent.ToolRequest) }`.
- rc triple (pristine / injected / restored) = 0 / 1 / 0.
- Red, exact lines (logs/run-new-injected.txt):
  - `--- FAIL: Test309ToolRequestProductionConstructionNail (0.17s)`
  - `    toolrequest_nail_309_test.go:390: 枚数钉响：... 里数到 2 枚 agent.ToolRequest 生产构造点，票 309 钉的是 1 枚且那枚落在 internal/agent/loop.go ...`
  - roster line: `名册=internal/agent/loop.go:675(same-package-bare-literal) | internal/tools/z309b_new.go:5(new-of-tool-request)。`
  - `分母=296 枚非测试源` (295 baseline + the injected file itself; the numerator is the red key).
- sha injected = sha aside = 0894926883d1788e66e7aec48e3db4b09355ac5d649d98687554dbc6f15fb395; aside cmp rc=0; tree fingerprint a8d3bc60... pre == post, cmp rc=0.
- Pristine/restored both `ok github.com/CarlosShao/wisp/internal/tools`.

## Pair 2 - `make([]T, n)`

- Inject: `internal/tools/z309b_make.go`, `func z309bRigMake(n int) []agent.ToolRequest { return make([]agent.ToolRequest, n) }`.
- rc = 0 / 1 / 0. Red lines: same `:390` frame, `数到 2 枚`, roster `internal/tools/z309b_make.go:5(make-slice-of-tool-request)`; 分母 296.
- sha injected = sha aside = d57182001be4ec2fae9d5a1beacd7466fa9a1d3e9e3ad7b6f652d71f0e3e875b; aside cmp rc=0; fingerprint pre/post a8d3bc60... cmp rc=0.

## Pair 3 - conversion from a defined type

- Inject: `internal/tools/z309b_conv.go`, `type z309bRigSame agent.ToolRequest` + `func z309bRigConv(x z309bRigSame) agent.ToolRequest { return agent.ToolRequest(x) }`.
- rc = 0 / 1 / 0. Red lines: same `:390` frame, `数到 2 枚`, roster `internal/tools/z309b_conv.go:7(conversion-to-tool-request)`; 分母 296.
- sha injected = sha aside = f31e22ffa32c23a2c29fd1f217fcfe93116bca8a220c3e48404e84e9991d34f8; aside cmp rc=0; fingerprint pre/post a8d3bc60... cmp rc=0.

## Census rig - all seven roster labels seen on disk, plus two false-positive vectors

- Phase 1 inject `internal/tools/z309b_census.go` (qualified literal, in-file alias literal, var-decl with no initializer). rc pristine=0, phase1=1: `数到 4 枚`, roster:
  `internal/agent/loop.go:675(same-package-bare-literal) | internal/tools/z309b_census.go:11(alias-literal) | internal/tools/z309b_census.go:7(var-decl-no-initializer) | internal/tools/z309b_census.go:9(qualified-literal)。`
- Phase 2 additionally injects `internal/tools/fakeagent/fake.go` (package clause `agent`, its own `type ToolRequest struct`, `var Sample = ToolRequest{...}`) and `internal/tools/z309b_fp.go` (imports the fakeagent package bound as `agent`, calls `agent.ToolRequest{...}`). rc=1: `数到 6 枚`, roster adds:
  `internal/tools/fakeagent/fake.go:5(same-package-bare-literal) | internal/tools/z309b_fp.go:5(qualified-literal)。`
  - Reading: the matcher resolves by spelling, not by import path or type identity. A package-clause-only `agent` name and a foreign import bound as `agent` both ring (over-report direction). Live risk today = 0 (see 20-cells: no non-wisp `/agent` import under roots; the only 2 non-wisp hits of the path regex are comment lines inside the nail test itself).
- Phase 3: all three injected artefacts moved aside; rc=0; fingerprint pre/post a8d3bc60... cmp rc=0; aside cmp vs staged copies rc=0 (census, fp) and sha of fake.go unchanged (8529b969...).
- Seven reachable labels observed across rigs: same-package-bare-literal, new-of-tool-request, make-slice-of-tool-request, conversion-to-tool-request, var-decl-no-initializer, alias-literal, qualified-literal. This equals the bridge.go comment's "seven named shape labels" list (1:1).

## Files

- Script (rerunnable verbatim): logs/rigs-a.txt. Summary: logs/rigs-a-summary.txt. Per-rig logs: logs/rig-{new,make,conv,census}.txt, raw run logs logs/run-*.txt (pristine/injected/restored phases), fingerprints logs/fp-*.txt, void stderr files logs/void-*.stderr.
