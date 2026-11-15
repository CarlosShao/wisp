// Ticket 77 AC#7 - the vendoring step, kept in the repo so "where did this file
// come from" is answerable by reading a command instead of trusting a comment.
//
//   node scripts/vendor.mjs <path-to-ai-native-react-components-checkout>
//
// The upstream project is shipped as a Next.js demo app (one component per file,
// demo data baked in). Vendoring here is deliberately mechanical:
//   - prepend the provenance header (repo / source file / license / commit)
//   - drop the Next.js-only `"use client"` directive (this app is plain Vite)
// and nothing else. No reformatting, no value edits: a hand-tuned copy would
// silently fork the upstream component, which is the thing VENDORED.md exists
// to prevent. Upstream is git-cloned to a scratch dir and never added as an npm
// dependency (Q-20: copy-paste delivery, vendored into the tree).
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { argv, cwd, exit, stdout } from "node:process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { execFileSync } from "node:child_process";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = resolve(here, "..");

const UPSTREAM_REPO = "github.com/TurboKach/ai-native-react-components";
const LICENSE = "MIT (upstream LICENSE: Copyright (c) 2026 Turbo)";

// dest -> { src, component, note }
const JOBS = [
  { src: "components/atoms/StreamText.tsx", dest: "src/components/ai-native/stream-text.tsx", component: "StreamText", note: "NOT mounted: it builds segments with split on a single space, so a Chinese answer reveals as one unit; PLAN.md:3476 逐段 is served by our own src/components/reveal-text.tsx" },
  { src: "components/atoms/Shimmer.tsx", dest: "src/components/ai-native/shimmer.tsx", component: "Shimmer", note: "mounted: the waiting label in src/components/PanelSkeleton.tsx" },
  { src: "components/approval-card.tsx", dest: "src/components/ai-native/approval-card.tsx", component: "ApprovalCard", note: "NOT mounted: upstream demo questionnaire; kept as the visual reference the adapted src/components/l2-approval-card.tsx was cut from" },
  { src: "components/loading-state.tsx", dest: "src/components/ai-native/loading-state.tsx", component: "LoadingState", adapted: true, note: "ADAPTED 2026-09-25 (props-driven, demo timer removed, P9): a re-run must not clobber it - vendor.mjs skips adapted dests without --force" },
  { src: "components/thinking.tsx", dest: "src/components/ai-native/thinking.tsx", component: "Thinking", adapted: true, note: "ADAPTED 2026-09-25 (props-driven, sparkles head replaced by the pixel grid, P9)" },
  { src: "components/tool-chips.tsx", dest: "src/components/ai-native/tool-chips.tsx", component: "ToolChips", adapted: true, note: "ADAPTED 2026-09-25 (props-driven four-state tool chips, P9)" },
  { src: "components/task-rows.tsx", dest: "src/components/ai-native/task-rows.tsx", component: "TaskRows", note: "NOT mounted: demo data intact; history rows are ticket 35+" },
  { src: "components/streaming-text.tsx", dest: "src/components/ai-native/streaming-text.tsx", component: "StreamingText", adapted: true, note: "ADAPTED 2026-09-25 (props-driven, sources downgraded to spans, caret served by reveal-text, P9)" },
];

// D23 / ban #8 zero-emoji, applied to OUR copy: upstream writes U+2713 into two
// demo strings of tool-chips.tsx and U+2212 into the diff rows of both
// tool-chips.tsx and thinking.tsx. The ticket demands picking a side - ASCII-ize
// or register the tree as out of scope - so we ASCII-ize. An unmapped glyph is
// fatal rather than silently shipped, because a re-vendor must not be able to
// reopen the hole.
// Written with escapes rather than the glyphs themselves: internal/panel's
// TestFrontendHasNoEmoji scans this directory with the scanner's own ranges, and
// the table that removes the glyphs must not be the thing that reintroduces
// them. U+2713/U+2714 check, U+2717/U+2715 cross, U+2212 minus.
//
// EMOJI_RE's bands are a third copy of tools/d22scan/main.go's emojiRe and must
// travel with it. It previously omitted U+2200-U+22FF, which made this guard
// structurally blind to the two U+2212 that upstream ships - the hand-fix in
// commit 3b59512 would have been reverted silently by the next re-vendor.
// Measured before widening: on the pre-fix tree the old class reported 0 hits
// where the widened one reports 2. On today's tree both report 0, so widening
// costs nothing here and only buys the forward catch.
const GLYPH_MAP = new Map([
  ["\u2713", "[ok]"],
  ["\u2714", "[ok]"],
  ["\u2717", "[x]"],
  ["\u2715", "[x]"],
  ["\u2212", "-"],
]);
const EMOJI_RE = /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{FE0F}\u{1F1E6}-\u{1F1FF}]/gu;

function asciiize(text, dest) {
  let n = 0;
  const out = text.replace(EMOJI_RE, (g) => {
    const mapped = GLYPH_MAP.get(g);
    if (!mapped) {
      console.error(`vendor.mjs: ${dest} carries U+${g.codePointAt(0).toString(16).toUpperCase()}, which GLYPH_MAP does not cover - add a mapping or ASCII-ize by hand and record it in VENDORED.md`);
      exit(2);
    }
    n += 1;
    return mapped;
  });
  return { text: out, replaced: n };
}

function header(job, commit, glyphCount) {
  return `/* ============================================================================
   Vendored file - ticket 77 AC#7 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : ${UPSTREAM_REPO}  (registry name "ai-native", beautifului.dev)
   Source file    : ${job.src}
   Component      : ${job.component}
   License        : ${LICENSE}
   Upstream commit: ${commit}
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; line endings normalized to LF;
                    ${glyphCount} dingbat glyph(s) ASCII-ized for D23/ban #8.${glyphCount ? "" : " (none found in this file)"}
                    Nothing else - see scripts/vendor.mjs.
   Panel status   : ${job.note}
   ============================================================================ */

`;
}


const upstream = argv[2];
if (!upstream || !existsSync(join(upstream, "LICENSE"))) {
  console.error("usage: node scripts/vendor.mjs <checkout-of-ai-native-react-components>");
  exit(2);
}
const commit = execFileSync("git", ["-C", upstream, "rev-parse", "HEAD"], { encoding: "utf8" }).trim();

mkdirSync(join(frontendRoot, "src", "components", "ai-native"), { recursive: true });
for (const job of JOBS) {
  // 2026-09-25 pivot: four of these dests are no longer verbatim copies - they
  // were ADAPTED into props-driven components (demo data removed, P9). A
  // mechanical re-vendor would silently restore upstream demo pages over the
  // adaptations, so adapted dests are skipped unless --force.
  if (job.adapted && !argv.includes("--force")) {
    stdout.write(`skipped ${job.dest}  (adapted by hand - see VENDORED.md; pass --force to overwrite)\n`);
    continue;
  }
  const raw = readFileSync(join(upstream, job.src), "utf8");
  // \r\n tolerance: upstream is a Windows checkout, so the directive can end
  // with either line ending (a plain \n regex silently matches nothing). Line
  // endings are then normalized to LF so the whole tree stays uniform.
  const body = raw.replace(/^\s*"use client";\r?\n\r?\n?/, "").replace(/\r\n/g, "\n");
  if (/"use client"/.test(body)) {
    console.error(`vendor.mjs: ${job.dest} still contains a use-client directive`);
    exit(2);
  }
  const cleaned = asciiize(body, job.dest);
  writeFileSync(join(frontendRoot, job.dest), header(job, commit, cleaned.replaced) + cleaned.text, "utf8");
  stdout.write(`vendored ${job.dest}  (${job.component}, upstream ${commit.slice(0, 12)}, glyphs ASCII-ized: ${cleaned.replaced})\n`);
}
