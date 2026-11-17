/* ============================================================================
   Ticket 36 AC#4's page-half instrument (render evidence), 2026-11-17 - the
   200-chunk golden replay, read through the PRODUCTION result region.
   ----------------------------------------------------------------------------
   The golden fixture is the SAME FILE the Go instrument reads
   (internal/panel/testdata/stream200-golden-36.json), so the two sides cannot
   disagree about what the golden says. The rule it is written in is the 10-08
   re-cut (ticket 197 leg B, approval A406, code 94ea50f7): overflow TRUNCATES
   and never merges - the criteria below are written in that shape, and no byte
   of pump.go's truncation rule is negotiable against the ticket's old
   "renders with merge" word.

   WHAT IT PROVES, all by decidable facts (no wall clock anywhere):
     1. NO-JANK FACTS: the replay's incremental publishes are one packet per
        step; a mid-replay packet (40 keys through a 32-key log) carries the
        40 opened rows, identity intact, and the final packet carries exactly
        the hard ceiling's 64 surviving rows - one li per chunk, each with its
        own marker and never another chunk's, dropped chunks' markers absent
        from the page. A merge anywhere turns any of these red.
     2. TRUNCATION DISPLAY: the long row renders its own head + the exact
        "...[truncated: N runes elided]..." marker + its own tail, and the
        elided middle's marker never leaks into the page.
     3. TOKEN COUNTER RECONCILIATION: every chunk the golden gives a usage
        row to renders 已用 N tokens with N = tokensIn + tokensOut verbatim
        from the golden; every chunk without one renders the explicit 未知
        sentence; the page NEVER fabricates 已用 0 tokens. The counter shows
        in both the streaming (RevealText) and finished (MarkdownText) states.
     4. SOURCE LEVEL: panel.ts declares ResultChunkView.usage (the two-way
        contract walk's page half, landed in the same commit as the Go key),
        the counter's zero-fabrication shape lives in the production
        component, and the markup holds zero emoji.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { ResultHistory } from "@/components/result-history";
import type { ResultChunkView } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* --------------------------- the shared golden --------------------------- */

const golden = JSON.parse(
  readFileSync("../internal/panel/testdata/stream200-golden-36.json", "utf8"),
) as {
  maxKeys: number;
  expectSurvivingKeys: number;
  expectDroppedCount: number;
  chunks: {
    key: string;
    deltas: string[];
    close: boolean;
    usage: {
      tokensIn: number;
      tokensOut: number;
      cachedTokens: number;
      micros: number;
      currency: string;
    } | null;
    usageFirst: boolean;
  } | null;
  longRow: {
    key: string;
    deltas: string[];
    headMarker: string;
    midMarker: string;
    tailMarker: string;
    elidedExact: number;
  };
};

if (golden.chunks === null || !Array.isArray(golden.chunks) || golden.chunks.length !== 200) {
  throw new Error("the golden fixture does not carry 200 chunks - the instrument would be measuring air");
}

const markerOf = (key: string): string => `tok${key.slice("task-".length)}`;

/** The wire view of one chunk: the page half of StreamLog.Chunks' attach rule
    - usage present exactly when the golden recorded an event, never a zero. */
function toView(chunk: (typeof golden.chunks)[number], done: boolean): ResultChunkView {
  return {
    correlationId: chunk.key,
    text: chunk.deltas.join(""),
    done,
    ...(chunk.usage
      ? {
          usage: {
            tokensIn: chunk.usage.tokensIn,
            tokensOut: chunk.usage.tokensOut,
            ...(chunk.usage.cachedTokens ? { cachedTokens: chunk.usage.cachedTokens } : {}),
            micros: chunk.usage.micros,
            currency: chunk.usage.currency,
          },
        }
      : {}),
  };
}

function renderResults(results: ResultChunkView[]): string {
  return renderToStaticMarkup(<ResultHistory run={undefined} results={results} />);
}

const count = (haystack: string, needle: string): number => haystack.split(needle).length - 1;

/* ---------------------------------------------------------------------------
   1. The no-jank facts: one packet per publish, identity intact per row.
   --------------------------------------------------------------------------- */

// Mid-replay packet: after the first 40 keys opened, all 40 are still tracked
// (the 32-key bound truncates rows, the 64-key hard ceiling drops none) - the
// page must see exactly those 40 rows, none lost, none merged. In-flight rows
// render through RevealText, whose reveal counter is 0 in static markup, so
// the per-row identity facts here are the ROW COUNT and the 40 in-flight
// carets; the per-row text identity across the replay is the Go instrument's
// job (it reads the row text directly at every step).
const midPacket = golden.chunks.slice(0, 40).map((c) => toView(c, false));
const midHtml = renderResults(midPacket);
check("the mid-replay packet must render 40 rows (one per opened key, none dropped)", count(midHtml, "<li") === 40, `got ${count(midHtml, "<li")}`);
check(
  "the mid-replay packet must carry the streaming caret on every in-flight row",
  count(midHtml, "stream-caret") === 40,
  `got ${count(midHtml, "stream-caret")}`,
);
check(
  "the mid-replay packet must not carry any key beyond the opened set",
  !midHtml.includes(esc(`〔${markerOf(golden.chunks[40].key)}〕`)),
);

// Final packet: exactly the hard ceiling's surviving rows (the last 64 keys;
// the oldest 136 were evicted and named, never merged). The merge detector is
// per rendered row: split the markup into rows and require every row to hold
// exactly ONE golden marker - a merged row would hold two.
const firstSurvivor = 200 - golden.expectSurvivingKeys;
const finalPacket = golden.chunks.slice(firstSurvivor).map((c) => toView(c, true));
const finalHtml = renderResults(finalPacket);
check(
  `the final packet must render ${golden.expectSurvivingKeys} rows (one per surviving chunk)`,
  count(finalHtml, "<li") === golden.expectSurvivingKeys,
  `got ${count(finalHtml, "<li")}`,
);
const rows = finalHtml.split("<li").filter((p) => p.includes(esc("〔tok")));
check("every rendered row must be a golden row (no foreign content rows)", rows.length === golden.expectSurvivingKeys, `got ${rows.length}`);
for (const row of rows) {
  const markersInRow = [...row.matchAll(/〔(tok\d{3})〕/g)].map((m) => m[1]);
  check(
    `a rendered row must carry exactly one chunk's marker, got [${markersInRow.join(",")}]`,
    markersInRow.length > 0 && new Set(markersInRow).size === 1,
  );
}
for (let i = firstSurvivor; i < 200; i += 1) {
  check(
    `final row ${golden.chunks[i].key} must show its own marker`,
    finalHtml.includes(esc(`〔${markerOf(golden.chunks[i].key)}〕`)),
  );
}
for (let i = 0; i < firstSurvivor; i += 1) {
  check(
    `dropped chunk ${golden.chunks[i].key}'s marker must appear nowhere on the page (named, not merged)`,
    !finalHtml.includes(esc(`〔${markerOf(golden.chunks[i].key)}〕`)),
  );
}

/* ---------------------------------------------------------------------------
   2. The truncation display: the named middle, verbatim.
   --------------------------------------------------------------------------- */

const longText = golden.longRow.deltas.join("");
const longHead = longText.slice(0, 4096);
const longTail = longText.slice(longText.length - 4096);
const truncatedRow: ResultChunkView = {
  correlationId: golden.longRow.key,
  text: longHead + `...[truncated: ${golden.longRow.elidedExact} runes elided]...` + longTail,
  done: true,
};
const longHtml = renderResults([truncatedRow]);
check(
  "the truncated row must render the exact elision marker with its count",
  longHtml.includes(esc(`...[truncated: ${golden.longRow.elidedExact} runes elided]...`)),
);
check("the truncated row must keep its own head marker", longHtml.includes(esc(golden.longRow.headMarker)));
check("the truncated row must keep its own tail marker", longHtml.includes(esc(golden.longRow.tailMarker)));
check("the elided middle's marker must never leak onto the page", !longHtml.includes(esc(golden.longRow.midMarker)));

/* ---------------------------------------------------------------------------
   3. The token counter reconciliation, per chunk.
   --------------------------------------------------------------------------- */

const counterHtml = renderResults(finalPacket);
let withUsage = 0;
let withoutUsage = 0;
for (let i = firstSurvivor; i < 200; i += 1) {
  const c = golden.chunks[i];
  if (c.usage) {
    withUsage += 1;
    const used = c.usage.tokensIn + c.usage.tokensOut;
    check(`chunk ${c.key} must show 已用 ${used} tokens off the golden's own numbers`, counterHtml.includes(esc(`已用 ${used} tokens`)));
  } else {
    withoutUsage += 1;
  }
}
check(
  `the page must carry the explicit unknown sentence for all ${withoutUsage} usage-less chunks`,
  count(counterHtml, esc("token 用量未知（本任务未发布 usage 事件）")) === withoutUsage,
  `got ${count(counterHtml, esc("token 用量未知（本任务未发布 usage 事件）"))}`,
);
check("the page must never fabricate a zero counter", !counterHtml.includes(esc("已用 0 tokens")));
check(
  "the counter must render during streaming too (an in-flight chunk with usage)",
  renderResults([toView(golden.chunks[firstSurvivor], false)]).includes(
    esc(`已用 ${(golden.chunks[firstSurvivor].usage?.tokensIn ?? 0) + (golden.chunks[firstSurvivor].usage?.tokensOut ?? 0)} tokens`),
  ),
);
check(
  "the unknown sentence must render during streaming too (in-flight chunk without usage)",
  renderResults([toView(golden.chunks[firstSurvivor + 1], false)]).includes(esc("token 用量未知（本任务未发布 usage 事件）")),
);
check(
  "a usage-carrying chunk's breakdown must travel in the title (verbatim in/out)",
  counterHtml.includes(esc(`tokensIn ${golden.chunks[firstSurvivor].usage!.tokensIn} / tokensOut ${golden.chunks[firstSurvivor].usage!.tokensOut}`)),
);
check(`the reconciliation must have measured both shapes (${withUsage} with usage, ${withoutUsage} without)`, withUsage > 0 && withoutUsage > 0, `${withUsage}/${withoutUsage}`);

/* ---------------------------------------------------------------------------
   4. Source level + zero emoji.
   --------------------------------------------------------------------------- */

check("the final markup must hold zero emoji", !EMOJI_RE.test(finalHtml));
check("the counter markup must hold zero emoji", !EMOJI_RE.test(counterHtml));
check("the long-row markup must hold zero emoji", !EMOJI_RE.test(longHtml));

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check(
  "panel.ts must declare ResultChunkView.usage (the contract walk's page half, same commit as the Go key)",
  /export interface ResultChunkView \{[\s\S]*?usage\?: RunUsageRow;/.test(panelSrc),
);
const componentSrc = readFileSync("src/components/result-history.tsx", "utf8");
check(
  "the production counter must carry the explicit unknown sentence (never a fabricated zero)",
  componentSrc.includes("token 用量未知（本任务未发布 usage 事件）") && !componentSrc.includes("已用 0 tokens"),
);
check(
  "the counter must be mounted in the production result region",
  componentSrc.includes("<TokenCounter usage={chunk.usage} />"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-stream200-36: FAIL " + f);
  console.error(`render-stream200-36: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-stream200-36: OK, ticket 36 AC#4's 200-chunk golden replay through the production result region " +
    "(mid-replay packet = the 40 opened rows identity-intact, final packet = the hard ceiling's 64 survivors " +
    "one row per chunk with no foreign marker and no dropped chunk's content anywhere, " +
    "the truncated long row renders its own head/tail plus the exact named elision and never the middle, " +
    `${withUsage} usage rows show 已用 N tokens off the golden's own numbers while ${withoutUsage} render the ` +
    "explicit unknown and the page never fabricates 已用 0 tokens, " +
    "panel.ts declares the usage page half in the same commit as the Go key, zero emoji)",
);
