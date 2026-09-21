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
  { src: "components/atoms/StreamText.tsx", dest: "src/components/ai-native/stream-text.tsx", component: "StreamText", note: "mounted: the result stream reuses it (src/components/ResultStream.tsx)" },
  { src: "components/atoms/Shimmer.tsx", dest: "src/components/ai-native/shimmer.tsx", component: "Shimmer", note: "mounted: the waiting label in src/components/PanelSkeleton.tsx" },
  { src: "components/approval-card.tsx", dest: "src/components/ai-native/approval-card.tsx", component: "ApprovalCard", note: "NOT mounted: upstream demo questionnaire; kept as the visual reference the adapted src/components/l2-approval-card.tsx was cut from" },
  { src: "components/loading-state.tsx", dest: "src/components/ai-native/loading-state.tsx", component: "LoadingState", note: "NOT mounted: ships with upstream demo strings; the panel's own loading state is PanelSkeleton.tsx until ticket 35 feeds it" },
  { src: "components/thinking.tsx", dest: "src/components/ai-native/thinking.tsx", component: "Thinking", note: "NOT mounted: demo data intact, see the note above the default export" },
  { src: "components/tool-chips.tsx", dest: "src/components/ai-native/tool-chips.tsx", component: "ToolChips", note: "NOT mounted: demo data intact; tool traces arrive through the C17 bridge (ticket 35)" },
  { src: "components/task-rows.tsx", dest: "src/components/ai-native/task-rows.tsx", component: "TaskRows", note: "NOT mounted: demo data intact; history rows are ticket 35+" },
  { src: "components/streaming-text.tsx", dest: "src/components/ai-native/streaming-text.tsx", component: "StreamingText", note: "NOT mounted: demo answer text + external source links; the panel must not render upstream URLs" },
];

// D23 / ban #8 zero-emoji, applied to OUR copy: upstream writes U+2713 into two
// demo strings of tool-chips.tsx. The ticket demands picking a side - ASCII-ize
// or register the tree as out of scope - so we ASCII-ize. An unmapped glyph is
// fatal rather than silently shipped, because a re-vendor must not be able to
// reopen the hole.
const GLYPH_MAP = new Map([
  ["✓", "[ok]"],
  ["✔", "[ok]"],
  ["✗", "[x]"],
  ["✕", "[x]"],
]);
const EMOJI_RE = /[\u{1F000}-\u{1FAFF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{FE0F}\u{1F1E6}-\u{1F1FF}]/gu;

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
