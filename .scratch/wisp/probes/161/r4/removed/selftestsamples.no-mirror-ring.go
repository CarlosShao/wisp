package main

// The self-test sample table (ticket 161 AC#2). Machinery and provenance notes
// live in selftest.go; this file is the samples, so that the one edit AC#4 asks
// about - "delete the violating sample" - is visible in exactly one place.
//
// Every tag listed here is one the scanner can emit, and every tag it can emit is
// required to appear here in BOTH directions by auditSelfCases(). Shapes are taken
// from the 27 measurements in docs/evidence/s1/161-gate-blindspot-r1.md §3 (the
// expect-ring / expect-silent pairs that program produced on a throwaway bench CI
// never ran) plus the six registered blind spots and three temperaments of §4/§4.1,
// so this table is a transcription of readings, not a restatement of hopes.
//
// Where a case pins something uncomfortable, its `note` says so and says which
// registered hole it closes the moment someone fixes that hole. A pin that flips
// because a ban was WIDENED is a contract change (Q-46 / ticket 141 territory) and
// must arrive with an owner's approval, never as a side effect of test work.

// Glyphs are built from their code points so this file's bytes stay pure ASCII (a
// literal U+2713 in a Go source file is hard to tell apart from a violation of the
// very ban this table samples). The scanner reads the bytes of the seeded FIXTURE,
// and that is where the rune has to actually appear. U+2713 and U+2264 are inside
// emojiRe's bands; U+2192 and U+2460 are the two bands Q-46 deliberately left out
// (main.go:161-163), and their cases below pin that gap rather than wish it away.
var (
	glyphMath    = string(rune(0x2264))  // U+2264 - inside the math band ticket 141 added
	glyphCheck   = string(rune(0x2713))  // U+2713 - inside U+2600-U+27BF
	glyphArrow   = string(rune(0x2192))  // U+2192 - deliberately NOT scanned
	glyphCircled = string(rune(0x2460))  // U+2460 - deliberately NOT scanned
	glyphFace    = string(rune(0x1F600)) // U+1F600 - inside U+1F000-U+1FAFF
)

// selfCases is the table runSelfTest executes, in the order it prints.
var selfCases = []selfCase{
	// ---- ban #1 bare-goroutine ------------------------------------------------
	{
		tag: "bare-goroutine", want: wantRing,
		file:    "internal/probe/leak.go",
		src:     "package probe\n\nfunc leak() {\n\tgo func() { println(\"unnamed\") }()\n}\n",
		summary: "closure form `go func(){}` (D22/D38b)",
	},
	{
		tag: "bare-goroutine", want: wantRing,
		file:    "internal/probe/leaknamed.go",
		src:     "package probe\n\nfunc worker() {}\n\nfunc leakNamed() {\n\tgo worker()\n}\n",
		summary: "named form `go worker()` (R16#1 widened the ban to this)",
		note:    "161-r1 §3 cell b1-named-bad. Before R16#1 the matcher saw only `go func(`, so a clean report meant \"no bare closures\", not \"nothing bypasses Spawn\".",
	},
	{
		tag: "bare-goroutine", want: wantSilent,
		file: "internal/probe/spawned.go",
		src: "package probe\n" +
			"\n" +
			"// The banned shape, written down so a reviewer can see it was considered\n" +
			"// and rejected: `go func() { work() }()`. Comments are invisible to this\n" +
			"// ban because it is an AST check (GoStmt), which is the exemption the\n" +
			"// D22 wording intends.\n" +
			"func schedule(reg interface{ Spawn(func()) }) {\n\treg.Spawn(work)\n}\n" +
			"\n" +
			"func work() {}\n",
		summary: "the words `go func()` only in a comment, real call goes through a Spawn-shaped seam",
	},
	{
		tag: "bare-goroutine", want: wantRing,
		file:    "cmd/probe/leak.go",
		src:     "package main\n\nfunc leak() {\n\tgo func() { println(\"cmd scope\") }()\n}\n\nfunc main() {}\n",
		cover:   "cmdGo",
		summary: "same shape in the cmd/ half of the Go scope",
		note:    "pins that cmd/ is a live walk of its own, not a shadow of internal/ (declaredScopes lists both; ticket 71 AC#4 generalised the empty-scope guard to each).",
	},

	// ---- ban #2 pathresolver-bypass ------------------------------------------
	{
		tag: "pathresolver-bypass", want: wantRing,
		file: "internal/probe/paths.go",
		src: "package probe\n" +
			"\n" +
			"import \"path/filepath\"\n" +
			"\n" +
			"func Normalize(p string) string { return filepath.Clean(p) }\n" +
			"\n" +
			"func Absolute(p string) string { return filepath.Abs(p) }\n",
		summary: "filepath.Clean and filepath.Abs, both arms in one file",
	},
	{
		tag: "pathresolver-bypass", want: wantSilent,
		file: "internal/probe/pathsok.go",
		src: "package probe\n" +
			"\n" +
			"import \"path/filepath\"\n" +
			"\n" +
			"type lexical struct{ raw string }\n" +
			"\n" +
			"func (l lexical) Clean() string { return filepath.Base(l.raw) }\n" +
			"\n" +
			"func Join(a, b string) string { return filepath.Join(a, b) }\n" +
			"\n" +
			"func Dir(p string) string { return filepath.Dir(p) }\n",
		summary: "Join/Dir/Base and a method named Clean on another receiver",
		note:    "161-r1 §3 cell b2-clean: the ban is a C26 PathResolver bypass, not a ban on the filepath package.",
	},
	{
		tag: "pathresolver-bypass", want: wantSilent,
		file: "internal/probe/aliased.go",
		src: "package probe\n" +
			"\n" +
			"import fp \"path/filepath\"\n" +
			"\n" +
			"func Normalize(p string) string { return fp.Clean(p) }\n",
		summary: "KNOWN BLIND SPOT B1: the same call under an aliased import",
		note:    "161-r1 §4 B1 - the matcher compares the receiver's LOCAL NAME (main.go:715), not the package's identity, so `fp \"path/filepath\"` + `fp.Clean` decides outside the PathResolver and stays silent. Pinned as expect-silent to record the ban's REAL range today: if someone resolves the check by import path, this line has to be flipped to expect-ring, and that flip is a scope change belonging to Q-46 / ticket 141 (owner approval), not a test fix someone may make in passing.",
	},

	// ---- ban #3 plaintext-key -------------------------------------------------
	{
		tag: "plaintext-key", want: wantRing,
		file:    "internal/probe/keyconst.go",
		src:     "package probe\n\nconst probeAPIKey = \"PROBE161NOTACREDENTIAL0000\"\n",
		summary: "const with a key-named identifier holding a 24-char literal",
		note:    "The literal is a synthetic string invented for this cell (no service, no account, no environment). 161-r1 §8 states the same rule for its samples; nothing here is or was ever a real credential.",
	},
	{
		tag: "plaintext-key", want: wantRing,
		file:    "internal/probe/keyassign.go",
		src:     "package probe\n\nfunc read() string {\n\tvar probeSecret string\n\tprobeSecret = \"PROBE161NOTACREDENTIAL0001\"\n\treturn probeSecret\n}\n",
		summary: "assignment arm: identifier on the left, literal on the right",
	},
	{
		tag: "plaintext-key", want: wantSilent,
		file:    "internal/probe/keyref.go",
		src:     "package probe\n\nconst probeAPIKeyRef = \"ref:wisp/llm-api-key\"\n\nconst shortToken = \"abc123\"\n",
		summary: "key-named identifiers whose values are SecretStore refs or too short",
		note:    "looksLikePlaintextKey() (main.go:803) exempts the ref: / env: / dpapi: / placeholder prefixes and requires >=16 chars of [A-Za-z0-9._~-] - C28 says the DB may store a REFERENCE, so a reference must not read as a violation.",
	},
	{
		tag: "plaintext-key", want: wantSilent,
		file: "internal/probe/keystruct.go",
		src: "package probe\n" +
			"\n" +
			"type cfg struct{ APIKey string }\n" +
			"\n" +
			"var loaded = cfg{APIKey: \"PROBE161NOTACREDENTIAL0002\"}\n",
		summary: "KNOWN BLIND SPOT B2: the same literal in a composite-literal field",
		note:    "161-r1 §4 B2 - the AST walks ValueSpec and Ident-lvalue AssignStmt only (main.go:723-755), so a plaintext key inside a struct/map literal or an argument position is invisible. Not sampled: map literals, call-argument position, pointer fields (161-r1 §2 N3).",
	},

	// ---- ban #4 wallclock-timeout --------------------------------------------
	{
		tag: "wallclock-timeout", want: wantRing,
		file: "internal/probe/deadline.go",
		src: "package probe\n" +
			"\n" +
			"import \"time\"\n" +
			"\n" +
			"func remaining(deadline time.Time) time.Duration {\n" +
			"\treturn deadline.Sub(time.Now())\n" +
			"}\n",
		summary: "deadline.Sub(time.Now()) - the wall-clock delta arm",
	},
	{
		tag: "wallclock-timeout", want: wantRing,
		file: "internal/probe/unixdeadline.go",
		src: "package probe\n" +
			"\n" +
			"import \"time\"\n" +
			"\n" +
			"func nowAsUnix() int64 {\n" +
			"\ttimeoutUnix := time.Now().Unix() + 30\n" +
			"\treturn timeoutUnix\n" +
			"}\n",
		summary: "time.Now().Unix() on a line that also names a timeout - the second arm",
		note:    "The timeout WORD must be on the SAME line as the unix call: main.go:767 conjoins unixTimeRe and timeoutWordRe per line. The first draft of this case had the word only in the enclosing function name, and the self-test printed FAIL for it - the case was wrong, the gate was not. Recorded because 161-r1 §3 cell b4-unix-bad used a one-line assignment for exactly this reason, and because a self-test that never contradicts its author is not measuring anything.",
	},
	{
		tag: "wallclock-timeout", want: wantSilent,
		file: "internal/probe/monotonic.go",
		src: "package probe\n" +
			"\n" +
			"import \"time\"\n" +
			"\n" +
			"func monotonic(start, deadline time.Time) time.Duration {\n" +
			"\t// time.Until / time.Since read the monotonic clock; so does a\n" +
			"\t// Duration subtraction between two Times taken from it.\n" +
			"\treturn time.Since(start) + time.Until(deadline) + deadline.Sub(start)\n" +
			"}\n",
		summary: "time.Until / time.Since / a.Sub(b) - the sanctioned clock shapes",
	},

	// ---- ban #5 mirror-hash ---------------------------------------------------
	{
		tag: "mirror-hash", want: wantSilent,
		file: "internal/probe/mirrorsplit.go",
		src: "package probe\n" +
			"\n" +
			"const upstreamURL = \"https://mirror.example/file.bin\"\n" +
			"\n" +
			"const digestAlgorithm = \"sha256\"\n",
		summary: "the same two vocabularies on separate lines",
		note:    "161-r1 §4.1 A2: the check is same-line co-occurrence and reads no data flow, so it neither catches a differently-spelled mirror fetch nor ignores a benign mention. Its second ring case (`b5-bad` line 7) fires on a pure reference line for the same reason.",
	},

	// ---- ban #6 panel-approval ------------------------------------------------
	{
		tag: "panel-approval", want: wantRing,
		file:    "frontend/src/app.js",
		src:     "export function decide(o) { return approval.decide(o); }\n",
		summary: "approval.decide called from the panel",
		note:    "D33/F2: allow decisions are native-side only, so an L2 \"yes\" must never originate in frontend/.",
	},
	{
		tag: "panel-approval", want: wantSilent,
		file:    "frontend/src/ask.js",
		src:     "export function ask(o) { return approval.request(o); }\n",
		summary: "approval.request - the panel may ask, it may not decide",
	},
	{
		tag: "panel-approval", want: wantRing,
		file: "frontend/src/notes.js",
		src: "// Documenting the banned call for reviewers: approval.decide(...) must never\n" +
			"// be reached from this tree; the decision belongs to the native side.\n" +
			"export const note = 1;\n",
		summary: "TEMPERAMENT A1: the call name appears only in a comment - and rings anyway",
		note:    "161-r1 §4.1 A1: bans #6 and #7 go through walkText, which strips no comments at all (unlike #8's Q-46(c) exemption and #4/#5's whole-line skip). Consequence in production today: one explanatory comment naming approval.decide turns CI red. Pinned as expect-ring because that IS the shipped range; narrowing or exempting it is a scope change (owner's call), not something a test may quietly do.",
	},

	// ---- ban #7 internal-artifact-tool ---------------------------------------
	{
		tag: "internal-artifact-tool", want: wantRing,
		file:    "internal/tools/gated.go",
		src:     "package tools\n\nvar gated = []string{\"spill\", \"internal.logwrite\"}\n",
		summary: "artifact-ish tool names as string literals in the tool registry tree",
	},
	{
		tag: "internal-artifact-tool", want: wantSilent,
		file:    "internal/tools/oklist.go",
		src:     "package tools\n\nvar publicTools = []string{\"fs.read\", \"fs.list\", \"web.search\"}\n",
		summary: "the real D34 tool names, none of them an artifact write",
	},
	{
		tag: "internal-artifact-tool", want: wantRing,
		file: "internal/tools/notes.go",
		src: "// Why \"internal.logwrite\" is not a tool: host-internal artifact writes are\n" +
			"// not gated tools (D34 note 2). The quoted name here is prose.\n" +
			"package tools\n",
		summary: "TEMPERAMENT A1 again: quoted name inside a comment - rings anyway",
		note:    "Same walkText shape as ban #6. This is the second reason 161-r1 called #6/#7 \"not blind, but temperamental\": prose in these two trees is judged as code.",
	},

	// ---- ban #8 emoji ---------------------------------------------------------
	{
		tag: "emoji", want: wantRing,
		file:    "internal/probe/glyph.go",
		src:     "package probe\n\nconst banner = \"ready " + glyphCheck + "\"\n",
		cover:   "emojiInternal",
		summary: "U+2713 inside a string literal in internal/",
		note:    "Q-46(c): strings are NOT exempt, comments are. 161-r1 §3 cell b8-probe-bands measured this exact glyph ringing.",
	},
	{
		tag: "emoji", want: wantRing,
		file:    "design/screens/states.html",
		src:     "<html><body>idle " + glyphFace + " glow</body></html>\n",
		cover:   "emojiDesign",
		summary: "U+1F600 in the design/ scope",
	},
	{
		tag: "emoji", want: wantRing,
		file:    "frontend/src/glyph.ts",
		src:     "export const label = \"built " + glyphMath + " 20 rows\";\n",
		cover:   "emojiFrontend",
		summary: "U+2264 in the frontend/ scope (the band ticket 141 added)",
		note:    "The frontend/ U+2713 that made three repo-wide rulers vanish for hours was ticket 169's subject; this sample is the same SCOPE with a different glyph, so a scope that stopped walking is visible here (161-r1 §4.1 A3: frontend/ is scanned with no suffix filter).",
	},
	{
		tag: "emoji", want: wantRing,
		file:    "cmd/probe/glyph.go",
		src:     "package main\n\nconst tail = \"done " + glyphCheck + "\"\n\nfunc main() {}\n",
		cover:   "emojiCmd",
		summary: "U+2713 in the cmd/ scope",
	},
	{
		tag: "emoji", want: wantSilent,
		file: "internal/probe/glyphcomment.go",
		src: "package probe\n" +
			"\n" +
			"// A glyph in prose is invisible to a renderer: " + glyphCheck + " and " + glyphFace + ".\n" +
			"func noop() {}\n",
		cover:   "emojiInternal",
		summary: "the same glyphs only inside a comment",
		note:    "Q-46(c) exemption, signed by the owner (ticket 141). 161-r1 §3 cell b8-comment-clean measured all four scopes silent on this shape.",
	},
	{
		tag: "emoji", want: wantSilent,
		file:    "internal/probe/arrowgap.go",
		src:     "package probe\n\nconst flow = \"read " + glyphArrow + " decide\"\n",
		cover:   "emojiInternal",
		summary: "KNOWN GAP: U+2192 in a string literal does not ring",
		note:    "main.go:161-163 states the arrow band U+2190-U+21FF is deliberately absent, and AGENTS.md §1.2 draws the practical consequence (\"a UI may show an arrow; the instrument cannot see it\"). Pinned expect-silent so a future widening is a loud flip on this line, which belongs to Q-46 / ticket 141 and needs an owner's approval.",
	},
	{
		tag: "emoji", want: wantSilent,
		file:    "internal/probe/circledgap.go",
		src:     "package probe\n\nconst step = \"" + glyphCircled + " first\"\n",
		cover:   "emojiInternal",
		summary: "KNOWN GAP: U+2460 (circled number) in a string literal does not ring",
		note:    "The second deliberate gap, U+2460-U+24FF. Same rule as the arrow pin above: TestBan8MathBandAndRemainingGaps already pins both directions in the test surface; this adds the entry-point reading a person running -self-test sees without opening a test file.",
	},

	// ---- unparseable (a finding type, not a numbered ban) ---------------------
	{
		tag: "unparseable", want: wantRing,
		file:    "internal/probe/broken.go",
		src:     "package probe\n\nfunc broken() {\n\tgo func() {\n}\n",
		summary: "a file the Go parser refuses",
	},
	{
		tag: "unparseable", want: wantSilent,
		file:    "internal/probe/fixed.go",
		src:     "package probe\n\nfunc fine() {}\n",
		summary: "a well-formed file in the same position",
	},
	{
		tag: "bare-goroutine", want: wantSilent,
		file: "internal/probe/brokenhidden.go",
		src: "package probe\n" +
			"\n" +
			"func broken() {\n" +
			"\tgo func() {\n" +
			"}\n",
		summary: "KNOWN BLIND SPOT B5: a bare goroutine inside a file that will not parse",
		note:    "161-r1 §4 B5 - scanGoFile() returns at the parse error (main.go:684-689), so bans #1-#3 (AST) AND #4-#5 (the line loop further down the same function) are all skipped for that file; the run prints one [unparseable] line and never says \"five bans did not look at this file\". This is ticket 161's founding shape - the gate still reports, it just stopped seeing. Fixing it means either a regex fallback or printing the skip; both change what the tool says, so the pin is expect-silent until someone owns that change.",
	},

	// ---- allowlist suppression (blind spot B6) --------------------------------
	{
		tag: "pathresolver-bypass", want: wantSilent,
		file: "internal/probe/allowed.go",
		src: "package probe\n" +
			"\n" +
			"import \"path/filepath\"\n" +
			"\n" +
			"func Normalize(p string) string { return filepath.Clean(p) }\n",
		allow:   "pathresolver-bypass\tinternal/\tfixture: a directory-prefix entry on purpose (161-r1 B6)\n",
		cover:   "internalGo",
		summary: "KNOWN BLIND SPOT B6: one directory-prefix entry silences the whole tree",
		note:    "161-r1 §4 B6 - suppression is per-ban (correct) but counted nowhere and printed nowhere, so \"the allowlist got one line wider\" is invisible in the self-report, including for the 5 real entries in tools/d22scan/allowlist.txt. allowlist.txt itself is FROZEN by ticket 161 AC#5; this case touches only a throwaway copy in a temp fixture root. The recomputable fix is a `suppressed=N (by entry ...)` line in verdict(), which is 161-r1 §2 N2's minimal closure action and is not this program's call.",
	},
	{
		tag: "bare-goroutine", want: wantRing,
		file: "internal/probe/allowedalso.go",
		src: "package probe\n" +
			"\n" +
			"import \"path/filepath\"\n" +
			"\n" +
			"func Normalize(p string) string { return filepath.Clean(p) }\n" +
			"\n" +
			"func leak() {\n\tgo func() { println(\"still visible\") }()\n}\n",
		allow:   "pathresolver-bypass\tinternal/\tfixture: same entry as the B6 pin\n",
		cover:   "internalGo",
		summary: "same file, same entry: ban #1 still rings while #2 is suppressed",
		note:    "The pair to the pin above: suppression is per-ban-id, so the B6 hole is \"invisible\", not \"global\". Stated because a reader who only saw the line above would conclude the allowlist can mute everything, which 161-r1 disproved.",
	},
}
