// Generator for internal/panel/testdata/stream200-golden-36.json - the 200-chunk
// golden fixture ticket 36 AC#4's replay instruments share (the Go test and the
// frontend render:stream200-36 instrument read the same file, so the two sides
// cannot disagree about what the golden says).
//
// Deterministic by construction: chunk i carries the unique marker tok<i> in
// its text (no other chunk's marker ever appears in it, which is what makes a
// cross-key merge detectable), deltas are three pieces whose concatenation is
// the golden text, even-numbered chunks carry the usage row the EvUsage event
// would have delivered (odd ones carry none - the page's "未知" shape), and
// the longRow part is one 10,000-rune stream with head/middle/tail markers at
// fixed offsets, fed in 100 deltas, to pin the per-row truncation below the
// key bound (head 4096 + marker + tail 4096, middle elided and NAMED).
//
// Run: node scripts/gen-stream200-golden-36.mjs   (from frontend/)
import { writeFileSync } from "node:fs";

const N = 200;

function chunkText(i) {
  const id = String(i).padStart(3, "0");
  return (
    `块 ${id}：本块只属于 task-${id}，标记 tok${id} 全场独有；` +
    `golden 回放逐块对账用（200 块之一，第 ${i + 1} 块）。`
  );
}

function usageRow(i) {
  return {
    tokensIn: 100 + i,
    tokensOut: 20 + i,
    cachedTokens: i % 7,
    micros: 1000 + i * 10,
    currency: "CNY",
  };
}

function split3(text, marker) {
  // Each delta carries the chunk's own marker so EVERY intermediate packet
  // (one publish per delta in production) can be identity-checked, not just
  // the final one; the golden text is still the deltas' concatenation.
  const third = Math.ceil(text.length / 3);
  return [
    text.slice(0, third) + marker,
    text.slice(third, 2 * third) + marker,
    text.slice(2 * third) + marker,
  ];
}

const chunks = [];
for (let i = 0; i < N; i += 1) {
  const text = chunkText(i);
  const marker = `tok${String(i).padStart(3, "0")}`;
  const entry = {
    key: `task-${String(i).padStart(3, "0")}`,
    deltas: split3(text, `〔${marker}〕`),
    close: true,
    // Usage arrives where the production sink delivers it: after the deltas,
    // before the close - except chunk 0, which exercises the other order
    // (usage before any delta) to pin that the read attaches either way.
    usage: i % 2 === 0 ? usageRow(i) : null,
    usageFirst: i === 0,
  };
  chunks.push(entry);
}

// One long stream, below the key bound: 10,000 runes total through 100 deltas.
// head marker at rune 0, mid marker at rune 5000 (inside the elided middle,
// which spans [4096, 10000-4096)), tail marker at the very end. The rendered
// row must keep head 4096 + the named marker + tail 4096 and must NOT carry
// the mid marker - that is the per-row truncation shape, never a merge.
const LONG = 10000;
const longBody = [];
for (let p = 0; longBody.join("").length < LONG; p += 1) {
  let piece = `第 ${p} 段：长流回放内容，`;
  while (piece.length < 100) piece += "x";
  longBody.push(piece.slice(0, 100));
}
const longText = longBody.join("").slice(0, LONG);
if (longText.length !== LONG) {
  throw new Error(`longText length ${longText.length} != ${LONG}`);
}

const headMarker = "HEAD0x长流头标记";
const midMarker = "MID0x长流中段标记(应被截断点名)";
const tailMarker = "TAIL0x长流尾标记";
const MID_POS = 5000;
const filler = LONG - headMarker.length - midMarker.length - tailMarker.length;
const before = MID_POS - headMarker.length; // filler runes ahead of the mid marker
const after = filler - before;
if (before <= 0 || after <= 0) throw new Error("marker layout does not fit the window shape");
const withMarkers =
  headMarker +
  longText.slice(0, before) +
  midMarker +
  longText.slice(longText.length - after) +
  tailMarker;
if (withMarkers.length !== LONG) {
  throw new Error(`withMarkers length ${withMarkers.length} != ${LONG}`);
}
const longDeltas = [];
for (let off = 0; off < withMarkers.length; off += 100) {
  longDeltas.push(withMarkers.slice(off, off + 100));
}

const fixture = {
  schema: "wisp/36-stream-token-golden@1",
  // The dispatch's names: 200 chunks across a DefaultStreamKeys=32 log, so the
  // replay crosses the key bound (truncate-and-name, never merge) and the hard
  // ceiling (64 tracked keys; 136 keys dropped and counted).
  maxKeys: 32,
  expectSurvivingKeys: 64,
  expectDroppedCount: 136,
  chunks,
  longRow: {
    key: "task-long",
    deltas: longDeltas,
    totalRunes: withMarkers.length,
    headMarker: "HEAD0x长流头标记",
    midMarker: "MID0x长流中段标记(应被截断点名)",
    tailMarker: "TAIL0x长流尾标记",
    // Must be computed from the retained windows, not stored: the golden says
    // the SHAPE, the log's own constants say the numbers.
    elidedExact: withMarkers.length - 2 * 4096,
  },
};

writeFileSync(
  "../internal/panel/testdata/stream200-golden-36.json",
  JSON.stringify(fixture, null, 1),
  "utf8",
);
console.log(`stream200-golden-36.json written: ${chunks.length} chunks + 1 long row`);
