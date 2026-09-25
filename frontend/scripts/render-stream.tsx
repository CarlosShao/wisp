/* ============================================================================
   Render + rule evidence for the streaming reveal (frontend session, 2026-09-25).

   Why this file exists
   --------------------
   Two defects were reported against the streaming answer:
     1. The vendored atom built its reveal units with text.split(" ")
        (src/components/ai-native/stream-text.tsx:38), so a Chinese answer came
        out as ONE unit -- no 逐段, and the caret vanished 46 ms into a paragraph
        that took seconds to write.
     2. theme.css froze the caret (`.stream-caret.is-streaming { animation: none }`)
        during the only interval it is mounted, and its geometry was upstream's
        2px x 1.05em --ink step blink rather than PLAN.md:3476's 1px x 14px
        --accent breathing curve.
   Neither is visible to `npm run typecheck`, `npm run lint`, `tokens:check` or
   `vite build`: a wrong segmenter and a wrong CSS declaration both compile.

   What the checks are derived from
   --------------------------------
   The caret's expected numbers are NOT written here. They are parsed out of the
   frozen row of PLAN.md's §17.5 table at run time, so this file cannot drift
   into agreeing with a wrong stylesheet, and so changing the spec text breaks
   the check loudly instead of quietly passing it. If PLAN.md or that row moves,
   this harness fails with a message saying so -- that is the intended reading,
   not a bug in the harness.

   What this harness CANNOT see (stated so nobody over-reads the green)
   -------------------------------------------------------------------
   react-dom/server renders a component at its initial state, count = 0, so the
   reveal never advances here and no markup assertion can prove "segments arrive
   one at a time". That half rests on the pure-function checks against
   revealSegments, which is why the segmenter lives in src/lib/reveal.ts rather
   than inline in the component. What nothing here proves is what the caret looks
   like on a real screen: that stays R-92-5 / ticket 114 AC#6 (differential
   screenshot, owner present).
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { RevealText } from "@/components/reveal-text";
import { ResultStream } from "@/components/result-stream";
import { MAX_SEGMENT, revealSegments } from "@/lib/reveal";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// ---------------------------------------------------------------------------
// 1. The frozen spec, parsed rather than copied.
// ---------------------------------------------------------------------------

interface CaretSpec {
  widthPx: string;
  heightPx: string;
  colorVar: string;
  opacityFrom: string;
  opacityMid: string;
  opacityTo: string;
  duration: string;
  source: string;
}

function readCaretSpec(): CaretSpec {
  let plan: string;
  try {
    plan = readFileSync("../docs/PLAN.md", "utf8");
  } catch {
    throw new Error(
      "render-stream: cannot read ../docs/PLAN.md. The caret numbers are derived " +
        "from it, so this check refuses to fall back to hard-coded expectations.",
    );
  }
  const row = plan
    .split("\n")
    .find((line) => line.startsWith("|") && line.includes("SSE") && line.includes("流式输出"));
  if (!row) throw new Error("render-stream: no 'SSE 流式输出' table row found in PLAN.md");
  const cells = row.split("|").map((c) => c.trim());
  const visual = cells[2] ?? "";
  const motion = cells[3] ?? "";
  // The spec writes the size with U+00D7 MULTIPLICATION SIGN, which is Latin-1
  // and sits below the lowest band ban #8 scans (U+2200). The proof of that is
  // running `sh scripts/d22scan.sh`, not this sentence.
  const size = visual.match(/(\d+(?:\.\d+)?)px\s*×\s*(\d+(?:\.\d+)?)px/u);
  const color = visual.match(/`(--[a-z0-9-]+)`/);
  const opacity = motion.match(/opacity\s*([\d.]+)\D+([\d.]+)\D+([\d.]+)/);
  const duration = motion.match(/(\d+(?:\.\d+)?)s/);
  for (const [name, found] of [
    ["a 'Npx x Npx' size", size],
    ["a backticked --token colour", color],
    ["an opacity sequence", opacity],
    ["a duration in s", duration],
  ] as const) {
    if (!found) throw new Error(`render-stream: PLAN.md's SSE row no longer carries ${name}`);
  }
  return {
    widthPx: size![1] as string,
    heightPx: size![2] as string,
    colorVar: color![1] as string,
    opacityFrom: opacity![1] as string,
    opacityMid: opacity![2] as string,
    opacityTo: opacity![3] as string,
    duration: `${duration![1] as string}s`,
    source: "PLAN.md §17.5 'SSE 流式输出'",
  };
}

const spec = readCaretSpec();

// ---------------------------------------------------------------------------
// 2. The stylesheet must carry that spec.
// ---------------------------------------------------------------------------

const css = readFileSync("src/styles/theme.css", "utf8");
function ruleBody(selector: string): string {
  const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const m = css.match(new RegExp(`^${escaped}\\s*\\{([^}]*)\\}`, "m"));
  if (!m) throw new Error(`render-stream: no '${selector}' rule in theme.css`);
  return m[1] as string;
}
function decl(body: string, property: string): string | undefined {
  const m = body.match(new RegExp(`(?:^|;)\\s*${property}:\\s*([^;]+)`));
  return m?.[1]?.trim();
}

const caret = ruleBody(".stream-caret");
const streaming = ruleBody(".stream-caret.is-streaming");

check("caret width must equal the spec", decl(caret, "width") === `${spec.widthPx}px`, `css=${decl(caret, "width")} spec=${spec.widthPx}px`);
check("caret height must equal the spec", decl(caret, "height") === `${spec.heightPx}px`, `css=${decl(caret, "height")} spec=${spec.heightPx}px`);
check("caret colour must be the spec's token", decl(caret, "background") === `var(${spec.colorVar})`, `css=${decl(caret, "background")} spec=${spec.colorVar}`);
check("the caret must not animate from the base rule (the frozen shape lived there)",
  decl(caret, "animation") === undefined, `css=${decl(caret, "animation")}`);

const anim = decl(streaming, "animation") ?? "none";
check("the streaming state must animate -- this is the declaration that read 'none'",
  anim !== "none", anim);
check("the caret animation must loop while text arrives", anim.includes("infinite"), anim);
check("the caret animation duration must equal the spec", anim.includes(` ${spec.duration}`), `css=${anim} spec=${spec.duration}`);

const frameName = anim.match(/caret-[a-z]+/)?.[0];
check("the rule must name a caret keyframes block", frameName !== undefined, anim);
if (frameName) {
  const body = css.match(new RegExp(`@keyframes\\s+${frameName}\\s*\\{([\\s\\S]*?)\\n\\}`))?.[1];
  check(`@keyframes ${frameName} must exist`, body !== undefined, frameName);
  if (body) {
    const stops = [...body.matchAll(/([^{}]+)\{([^}]*)\}/g)].map((m) => [
      (m[1] as string).trim(),
      (m[2] as string).match(/opacity:\s*([\d.]+)/)?.[1],
    ]);
    const at = (label: string) =>
      stops.find(([keys]) => keys.split(/[\s,]+/).includes(label))?.[1];
    check(`@keyframes ${frameName} must breathe to the spec's mid opacity, not blink to zero`,
      at("50%") === spec.opacityMid, `css=${at("50%")} spec=${spec.opacityMid}`);
    check(`@keyframes ${frameName} must return to the spec's endpoints`,
      at("0%") === spec.opacityFrom && (at("to") ?? at("100%")) === spec.opacityTo,
      `css=${at("0%")}/${at("to") ?? at("100%")} spec=${spec.opacityFrom}/${spec.opacityTo}`);
  }
}

// ---------------------------------------------------------------------------
// 3. The segmenter. This is where the Chinese defect is actually caught.
// ---------------------------------------------------------------------------

const CORPUS: [string, string][] = [
  ["chinese answer", "好的，我先看一下这个目录里有什么文件，然后再帮你整理成一份清单。"],
  ["chinese unpunctuated", "这是一段很长的中文回答里面一个标点都没有所以要靠长度上限来切"],
  ["english answer", "Sure, let me look at that directory for you."],
  ["mixed", "好的，I will run `npm test` 先看一下结果。"],
  ["bracketed", "「今天天气不错」他说，然后走了。"],
  ["one char", "好"],
  ["empty", ""],
  ["code line", "const x = 1;\n  return x;"],
];

for (const [label, text] of CORPUS) {
  const segs = revealSegments(text);
  check(`segments must rejoin into the input (${label})`, segs.join("") === text, `n=${segs.length}`);
  check(`no empty segment (${label})`, segs.every((s) => s.length > 0), JSON.stringify(segs));
}

// The bug, re-derived from the input rather than remembered. Stated as a scaling
// property rather than a count: upstream's split(" ") could only ever produce one
// unit for text with no spaces in it, so any such text longer than one segment
// had to arrive all at once. A short answer ("好") legitimately IS one segment,
// which is why length is part of the predicate and not an afterthought.
for (const [label, text] of CORPUS) {
  const hasSpace = /\S \S/.test(text);
  if (hasSpace || Array.from(text).length <= MAX_SEGMENT) continue;
  const oldUnits = text.split(" ").filter((w) => w.length > 0).length;
  const units = revealSegments(text).length;
  check(`text the old rule could only emit as one unit must reveal in many (${label})`,
    oldUnits <= 1 && units > 1, `old=${oldUnits} new=${units} chars=${Array.from(text).length}`);
}

// The no-regression pin: for space-separated LATIN text the rhythm must be what
// it was, or the "fix" would have silently restyled every English answer. CJK
// inputs are excluded on purpose -- there the whole point is that the old rhythm
// did not exist, and a mixed string's first word rides along on the Chinese
// clause before it, which is a different question from English pacing.
const CJK_INPUT = /[\u4E00-\u9FFF\u3040-\u30FF\uAC00-\uD7AF]/;
for (const [label, text] of CORPUS) {
  if (!/\S \S/.test(text) || CJK_INPUT.test(text)) continue;
  const latin = (s: string) => /^[A-Za-z`]/.test(s.trim());
  const ours = revealSegments(text).filter((s) => latin(s)).length;
  const before = text.split(" ").filter((s) => s.length > 0 && latin(s)).length;
  check(`latin word rhythm must be unchanged (${label})`, ours === before, `ours=${ours} split=${before}`);
}

// The bound here is a literal ON PURPOSE. Writing `<= MAX_SEGMENT` made this
// check read the same constant the mutation edits, so raising the cap to 100000
// passed it trivially -- a ruler that cannot be broken by the thing it claims to
// catch. 20 is "any answer must arrive in more than one piece for text this
// long", which is the property PLAN.md:3476's 逐段 actually asks for.
const longRun = "这是没有标点的长句".repeat(8);
check("a long unpunctuated CJK run must be split, not emitted whole",
  revealSegments(longRun).every((s) => Array.from(s).length <= 20) &&
    revealSegments(longRun).length > 1,
  `MAX_SEGMENT=${MAX_SEGMENT} n=${revealSegments(longRun).length} longest=${Math.max(...revealSegments(longRun).map((s) => Array.from(s).length))}`);
check("the exported cap must itself stay small enough to read as 逐段",
  MAX_SEGMENT > 0 && MAX_SEGMENT <= 20, `MAX_SEGMENT=${MAX_SEGMENT}`);

// Built from a code point so this source file carries no astral glyph of its own.
const astral = String.fromCodePoint(0x1f389).repeat(30);
check("an astral code point must never be cut in half",
  revealSegments(astral).every((s) =>
    !/[\uD800-\uDBFF](?![\uDC00-\uDFFF])|(?<![\uD800-\uDBFF])[\uDC00-\uDFFF]/.test(s)),
  `n=${revealSegments(astral).length}`);

// ---------------------------------------------------------------------------
// 4. What is mounted.
// ---------------------------------------------------------------------------

const mountPoint = readFileSync("src/components/result-stream.tsx", "utf8");
check("the result stream must not mount the vendored split-on-space atom",
  !mountPoint.includes("ai-native/stream-text"), "still imports the atom");

const live = renderToStaticMarkup(<RevealText text="好的，我先看一下这个目录里有什么文件。" />);
check("the caret must be mounted while text is arriving",
  live.includes("stream-caret") && live.includes("is-streaming"), live);
check("the caret must be hidden from assistive tech", live.includes("aria-hidden"), live);

const done = renderToStaticMarkup(
  <ResultStream chunks={[{ correlationId: "c1", text: "第一段结果，已经写完。", done: true }]} />,
);
check("a finished chunk must show its text verbatim", done.includes("第一段结果，已经写完。"), done);
check("a finished chunk must not leave an animating caret behind", !done.includes("is-streaming"), done);

if (failures.length > 0) {
  for (const f of failures) console.error("render-stream: FAIL " + f);
  console.error(`render-stream: FAIL, ${failures.length} check(s) failed (spec from ${spec.source})`);
  process.exit(1);
}
console.error(
  `render-stream: OK, spec ${spec.widthPx}x${spec.heightPx} ${spec.colorVar} ` +
    `${spec.opacityFrom}->${spec.opacityMid}->${spec.opacityTo} ${spec.duration}, ` +
    `${CORPUS.length} corpus input(s), caret markup checked at count=0`,
);
