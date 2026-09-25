# Ticket 149 mutation harness (OUT OF REPO; the committed copy of this file is
# .scratch/wisp/probes/149/harness.py, byte-identical).
#
# One literal replacement per mutation, each asserted to hit EXACTLY ONE place in
# the file it targets (else FATAL) - the same rule ticket 147's acceptance harness
# used, so a re-run that points at the wrong line cannot return a number that
# looks like green.
#
# Mutations are applied with `go test -overlay`, so nothing in any snapshot tree is
# rewritten or deleted: each mutation writes a new file under mut/<name>/ and an
# overlay.json, and every run's raw log is kept under logs/.
import json
import os
import subprocess
import sys

ROOT = os.path.dirname(os.path.abspath(__file__))
REPO = os.path.abspath(os.path.join(ROOT, "..", "..", ".."))  # unused, kept for clarity

GO_FILE = "cmd/wisp/slo_windows.go"
TEST_FILE = "cmd/wisp/slo_report_144_windows_test.go"

# --- literals (shared between the pre and post trees where they are identical) --
CORRUPT_TO_ZERO_PRE = (GO_FILE, """\t\tobs.offset = dec.InputOffset()
\t\tobs.state = reportCorrupt""", """\t\tobs.offset = 0
\t\tobs.state = reportCorrupt""")

COMPLETE_TO_ZERO_PRE = (GO_FILE, """\tobs.offset = dec.InputOffset()
\tif tok, terr := dec.Token(); !errors.Is(terr, io.EOF) {""", """\tobs.offset = 0
\tif tok, terr := dec.Token(); !errors.Is(terr, io.EOF) {""")

EXITED_DROP_SUMMARY = (GO_FILE, """\t\t\treturn nil, fmt.Errorf("wisp slo: subject %d exited (code %d) without writing its report (%s)",
\t\t\t\ts.pid, s.exitCode(), last.summary())""", """\t\t\treturn nil, fmt.Errorf("wisp slo: subject %d exited (code %d) without writing its report",
\t\t\t\ts.pid, s.exitCode())""")

TIMEOUT_DROP_SUMMARY = (GO_FILE, """\t\t\treturn nil, fmt.Errorf("wisp slo: subject %d never wrote a complete report within %s (last read: %s)",
\t\t\t\ts.pid, timeout.Budget(), last.summary())""", """\t\t\treturn nil, fmt.Errorf("wisp slo: subject %d never wrote a complete report within %s",
\t\t\t\ts.pid, timeout.Budget())""")

EXITED_SUMMARY_TO_CONSTANT = (GO_FILE, """\t\t\t\ts.pid, s.exitCode(), last.summary())""",
                              """\t\t\t\ts.pid, s.exitCode(), "last reading unavailable")""")

CORRUPT_FILL_TO_ZERO = (GO_FILE, """\t\tobs.offset = contradictionOffset(err, dec.InputOffset())""",
                        """\t\tobs.offset = 0""")
CORRUPT_FILL_TO_BYTES = (GO_FILE, """\t\tobs.offset = contradictionOffset(err, dec.InputOffset())""",
                         """\t\tobs.offset = int64(obs.bytes)""")
CORRUPT_FILL_REVERTED = (GO_FILE, """\t\tobs.offset = contradictionOffset(err, dec.InputOffset())""",
                         """\t\tobs.offset = dec.InputOffset()""")
SYNTAX_BRANCH_OFF = (GO_FILE, """\t\treturn syn.Offset""", """\t\treturn inputOffset""")
TYPE_BRANCH_OFF = (GO_FILE, """\t\treturn typ.Offset""", """\t\treturn inputOffset""")
CONTRADICTION_ALWAYS_BYTES = (GO_FILE, """\treturn inputOffset
}""", """\treturn inputOffset
}

func p149Unused() {}""")  # no-op placeholder, never selected
CORRUPT_SENTENCE_DROPS_OFFSET = (
    GO_FILE,
    """\t\tobs.err = fmt.Errorf("%d bytes contradict a subject report at offset %d: %w", obs.bytes, obs.offset, err)""",
    """\t\tobs.err = fmt.Errorf("%d bytes contradict a subject report: %w", obs.bytes, err)""")
CORRUPT_SENTENCE_FIRST_COUNT_TO_OFFSET = (
    GO_FILE,
    """\t\tobs.err = fmt.Errorf("%d bytes contradict a subject report at offset %d: %w", obs.bytes, obs.offset, err)""",
    """\t\tobs.err = fmt.Errorf("%d bytes contradict a subject report at offset %d: %w", obs.offset, obs.offset, err)""")
MULTIDOC_TAKES_CONTRADICTION_OFFSET = (
    GO_FILE,
    """\tif tok, terr := dec.Token(); !errors.Is(terr, io.EOF) {
\t\tobs.state = reportCorrupt""",
    """\tif tok, terr := dec.Token(); !errors.Is(terr, io.EOF) {
\t\tobs.state = reportCorrupt
\t\tobs.offset = contradictionOffset(terr, obs.offset)""")
NOSTATERPORT_SENTENCE_GAINS_OFFSET = (
    GO_FILE,
    """\t\tobs.err = fmt.Errorf("%d bytes parsed as a subject run but carry no state report", obs.bytes)""",
    """\t\tobs.err = fmt.Errorf("%d bytes parsed as a subject run but carry no state report at offset %d", obs.bytes, obs.offset)""")

# --- weakenings of ticket 149's own assertions -------------------------------
D11A_FIELD_LEG_OFF = (TEST_FILE, """\t\tif obs.offset != int64(tc.want) {""", """\t\tif false {""")
D11B_SENTENCE_LEG_OFF = (TEST_FILE, """\t\tcase offset != tc.want:""", """\t\tcase false:""")
D11C_READ_COUNT_LEG_OFF = (TEST_FILE, """\t\tcase read != len(tc.body):""", """\t\tcase false:""")
D11D_NAMED_COUNT_LEG_OFF = (TEST_FILE, """\t\tcase named != read:""", """\t\tcase false:""")
D11E_CASE_NOT_A_TEST = (TEST_FILE, """func TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt(""",
                        """func testSLO149CorruptSentenceNotRun(""")
D12_ASSERT_OFF = (TEST_FILE, """\t\tif obs.offset != int64(len(doc)) {""", """\t\tif false {""")
D12_NOTE_LEG_OFF = (TEST_FILE, """\t\tif strings.Contains(obs.summary(), "offset") {""", """\t\tif false {""")
D13A_LITERAL_WEAKENED = (TEST_FILE,
                         '''\tif want := "8 bytes read, document still open at offset 8: the tail had not arrived"; !strings.Contains(withPosition, want) {''',
                         '''\tif want := "code 7"; !strings.Contains(withPosition, want) {''')
D13B_CASE_NOT_A_TEST = (TEST_FILE, "func TestSLO149ExitedGiveUpSentenceCarriesTheLastReading(",
                        "func testSLO149ExitedCaseNotRun(")

MUT = {
    # pre tree = HEAD 64858d6, code and tests both as shipped: AC#1's readings
    "pre-asis": [],
    "pre-B14-corrupt-fill-to-zero": [CORRUPT_TO_ZERO_PRE],
    "pre-B15-complete-fill-to-zero": [COMPLETE_TO_ZERO_PRE],
    "pre-B13-exited-drops-summary": [EXITED_DROP_SUMMARY],
    "pre-B5-timeout-drops-summary": [TIMEOUT_DROP_SUMMARY],
    # post tree = ticket 149's code + tests: the load-bearing matrix
    "post-asis": [],
    "p1-corrupt-fill-to-zero": [CORRUPT_FILL_TO_ZERO],
    "p2-corrupt-fill-to-bytes": [CORRUPT_FILL_TO_BYTES],
    "p3-corrupt-fill-reverted": [CORRUPT_FILL_REVERTED],
    "p4-syntax-branch-off": [SYNTAX_BRANCH_OFF],
    "p5-type-branch-off": [TYPE_BRANCH_OFF],
    "p6-corrupt-sentence-drops-offset": [CORRUPT_SENTENCE_DROPS_OFFSET],
    "p7-exited-drops-summary": [EXITED_DROP_SUMMARY],
    "p8-exited-summary-to-constant": [EXITED_SUMMARY_TO_CONSTANT],
    "p9-timeout-drops-summary": [TIMEOUT_DROP_SUMMARY],
    "p10-multidoc-takes-error-offset": [MULTIDOC_TAKES_CONTRADICTION_OFFSET],
    "p11-nostate-report-sentence-gains-offset": [NOSTATERPORT_SENTENCE_GAINS_OFFSET],
    "p12-first-count-to-offset": [CORRUPT_SENTENCE_FIRST_COUNT_TO_OFFSET],
    "d11a-field-leg-off": [D11A_FIELD_LEG_OFF],
    "d11b-sentence-leg-off": [D11B_SENTENCE_LEG_OFF],
    "d11c-read-count-leg-off": [D11C_READ_COUNT_LEG_OFF],
    "d11d-named-count-leg-off": [D11D_NAMED_COUNT_LEG_OFF],
    "d11e-case-not-a-test": [D11E_CASE_NOT_A_TEST],
    "d12-offset-assert-off": [D12_ASSERT_OFF],
    "d12-note-leg-off": [D12_NOTE_LEG_OFF],
    "d13a-exited-literal-weakened": [D13A_LITERAL_WEAKENED],
    "d13b-case-not-a-test": [D13B_CASE_NOT_A_TEST],
    # combinations: which leg closes which mutation
    "x-p1-d11a": [CORRUPT_FILL_TO_ZERO, D11A_FIELD_LEG_OFF],
    "x-p1-d11b": [CORRUPT_FILL_TO_ZERO, D11B_SENTENCE_LEG_OFF],
    "x-p1-d11a-d11b": [CORRUPT_FILL_TO_ZERO, D11A_FIELD_LEG_OFF, D11B_SENTENCE_LEG_OFF],
    "x-p2-d11a-d11b": [CORRUPT_FILL_TO_BYTES, D11A_FIELD_LEG_OFF, D11B_SENTENCE_LEG_OFF],
    "x-p3-d11a-d11b": [CORRUPT_FILL_REVERTED, D11A_FIELD_LEG_OFF, D11B_SENTENCE_LEG_OFF],
    "x-p1-d11e": [CORRUPT_FILL_TO_ZERO, D11E_CASE_NOT_A_TEST],
    "x-p7-d13a": [EXITED_DROP_SUMMARY, D13A_LITERAL_WEAKENED],
    "x-p7-d13b": [EXITED_DROP_SUMMARY, D13B_CASE_NOT_A_TEST],
    "x-p12-d11c": [CORRUPT_SENTENCE_FIRST_COUNT_TO_OFFSET, D11C_READ_COUNT_LEG_OFF],
    "x-p12-d11d": [CORRUPT_SENTENCE_FIRST_COUNT_TO_OFFSET, D11D_NAMED_COUNT_LEG_OFF],
    "x-p10-d12": [MULTIDOC_TAKES_CONTRADICTION_OFFSET, D12_ASSERT_OFF],
    "x-preB14-newtests": [CORRUPT_TO_ZERO_PRE],
}

NEW_NAMES = [
    "TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt",
    "TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd",
    "TestSLO149ExitedGiveUpSentenceCarriesTheLastReading",
]


def apply(base, name, specs):
    """Write mutated copies + overlay.json for one mutation; return overlay path.

    ⚠ Bug found while reading the first pass of the combination runs: when two
    specs target the SAME file, the replacements must chain (each hits=1 check is
    against the text as it stands after the previous one), not be written
    independently - two separate writes to one overlay path silently keep only the
    last spec, which is how x-p1-d11a-d11b first ran with only d11b applied. That
    reading is re-taken as logs/combos2/x-p1-d11a-d11b.log.
    """
    out = os.path.join(ROOT, "mut", name)
    os.makedirs(out, exist_ok=True)
    texts = {}
    replace = {}
    for rel, old, new in specs:
        if rel not in texts:
            src = os.path.join(base, rel.replace("/", os.sep))
            with open(src, "r", encoding="utf-8", newline="") as fh:
                texts[rel] = fh.read()
        hits = texts[rel].count(old)
        if hits != 1:
            sys.exit("FATAL %s: literal hits=%d in %s, want exactly 1" % (name, hits, rel))
        texts[rel] = texts[rel].replace(old, new)
    for rel in texts:
        dst = os.path.join(out, os.path.basename(rel))
        with open(dst, "w", encoding="utf-8", newline="") as fh:
            fh.write(texts[rel])
        replace[os.path.join(base, rel).replace("\\", "/")] = dst.replace("\\", "/")
    overlay = os.path.join(out, "overlay.json")
    with open(overlay, "w", encoding="utf-8") as fh:
        json.dump({"Replace": replace}, fh)
    return overlay


def run(base, name, specs, run_pat, log_dir):
    overlay = apply(base, name, specs)
    cmd = ["go", "test", "-count=1", "-overlay", overlay, "-v", "-run", run_pat, "./cmd/wisp/"]
    env = dict(os.environ)
    env["PATH"] = os.path.join(base, "third_party", "sherpa-onnx") + os.pathsep + env["PATH"]
    proc = subprocess.run(cmd, cwd=base, env=env, capture_output=True, text=True)
    out = proc.stdout + proc.stderr
    logpath = os.path.join(log_dir, name + ".log")
    with open(logpath, "w", encoding="utf-8") as fh:
        fh.write("$ %s\n\n%s" % (" ".join(cmd), out))
    top = {}
    for line in out.splitlines():
        for tag in ("--- PASS: ", "--- FAIL: ", "--- SKIP: "):
            if line.startswith(tag):
                top[line[len(tag):].split(" ")[0]] = tag.strip("-: ").split()[0]
    ran = sum(1 for line in out.splitlines() if line.startswith("=== RUN"))
    reds = sorted([n for n, v in top.items() if v == "FAIL"])
    greens = sorted([n for n, v in top.items() if v == "PASS"])
    verdicts = {n: top.get(n, "NORUN") for n in NEW_NAMES}
    print("%-34s rc=%d RUN=%d topPASS=%d topFAIL=%d 144red=%d 147red=%d new=%s" % (
        name, proc.returncode, ran, len(greens), len(reds),
        sum(1 for n in reds if n.startswith("TestSLO144")),
        sum(1 for n in reds if n.startswith("TestSLO147")),
        " ".join("%s:%s" % (n[8:26], v) for n, v in verdicts.items())))
    if reds:
        print("      RED: %s" % ", ".join(reds))
    return reds


def main():
    base = sys.argv[1]
    log_dir = sys.argv[2]
    names = sys.argv[3:]
    os.makedirs(log_dir, exist_ok=True)
    if not names:
        names = list(MUT)
    for name in names:
        if name not in MUT:
            sys.exit("FATAL unknown mutation %s" % name)
        run(base, name, MUT[name], "TestSLO144|TestSLO147|TestSLO149", log_dir)


if __name__ == "__main__":
    main()
