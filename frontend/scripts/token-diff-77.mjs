// scripts/token-diff-77.mjs - ticket 77 AC#10 / ticket 34 AC#10 scripted diff.
//
// What it compares
//   A (legacy, first generation): design/assets/tokens.css. That file was
//   deleted from the repo in 72b65c614 (ledger A1010 wave); its authoritative
//   in-repo archive is internal/ball/testdata/tokens_c21.css, whose provenance
//   header pins the deleted file's sha256 and whose body is byte-for-byte the
//   source below that header (verified by diff at this commit: header-only
//   delta). When design/assets/tokens.css exists again it is used directly.
//   B (shipped, fourth generation): src/styles/tokens.generated.css, emitted
//   verbatim from the TOKENS table in scripts/gen-tokens.mjs (owner pivot
//   2026-09-25; sync to the generator is machine-checked by
//   `npm run tokens:check`, run separately - this script does not re-prove it).
//
// What it outputs: EVERY named difference - names only in A, names only in B,
// names in both with per-theme values that differ, and names identical in
// both. A difference count is printed; exit code is 0 whenever the diff ran
// to completion (this is a reporter, not a gate - the gate teeth live in
// gen-tokens --check and internal/panel's four-way test). It never claims
// parity: absence of differences prints as an empty section with a zero.
//
// Theme polarity (documented in internal/panel/tokens_fourway_test.go):
//   first generation: :root = dark, [data-theme="light"] = light
//   fourth generation: :root = light, [data-theme="dark"] = dark
// so same-name values are compared theme-to-theme, not block-to-block.

import { existsSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const frontendRoot = join(dirname(fileURLToPath(import.meta.url)), "..");
const repoRoot = join(frontendRoot, "..");

const DESIGN_TOKENS = join(repoRoot, "design", "assets", "tokens.css");
const ARCHIVE_TOKENS = join(repoRoot, "internal", "ball", "testdata", "tokens_c21.css");
const THEME_TOKENS = join(frontendRoot, "src", "styles", "tokens.generated.css");

function stripComments(text) {
  return text.replace(/\/\*[\s\S]*?\*\//g, "");
}

// One block per selector; declarations are flat (no nesting in either file).
function parseBlocks(text) {
  const blocks = {};
  const re = /(^|\n)\s*(:root|\[data-theme="light"\]|\[data-theme="dark"\])\s*\{([^{}]*)\}/g;
  for (const m of text.matchAll(re)) {
    const decls = {};
    for (const d of m[3].matchAll(/(--[A-Za-z0-9-]+)\s*:\s*([^;]+);/g)) {
      decls[d[1]] = d[2].trim().replace(/\s+/g, " ").toLowerCase();
    }
    blocks[m[2]] = decls;
  }
  return blocks;
}

function loadLegacy() {
  if (existsSync(DESIGN_TOKENS)) {
    return { path: DESIGN_TOKENS, text: readFileSync(DESIGN_TOKENS, "utf8"), direct: true };
  }
  if (!existsSync(ARCHIVE_TOKENS)) {
    console.error(`token-diff: neither ${DESIGN_TOKENS} nor ${ARCHIVE_TOKENS} exists - nothing to diff against`);
    process.exit(1);
  }
  const raw = readFileSync(ARCHIVE_TOKENS, "utf8");
  // Drop the provenance header comment so the body parses as bare CSS.
  const body = raw.replace(/^\s*\/\*[\s\S]*?\*\//, "");
  return { path: ARCHIVE_TOKENS, text: body, direct: false };
}

const legacy = loadLegacy();
const themeText = readFileSync(THEME_TOKENS, "utf8");
const legacyBlocks = parseBlocks(stripComments(legacy.text));
const themeBlocks = parseBlocks(stripComments(themeText));

// Theme-normalised views: dark/light -> { name: value }.
const A = {
  dark: { ...(legacyBlocks[":root"] ?? {}) },
  light: { ...(legacyBlocks['[data-theme="light"]'] ?? {}) },
};
const B = {
  light: { ...(themeBlocks[":root"] ?? {}) },
  dark: { ...(themeBlocks['[data-theme="dark"]'] ?? {}) },
};

const missing = [];
if (Object.keys(A.dark).length === 0) missing.push("A dark(:root)");
if (Object.keys(A.light).length === 0) missing.push("A light([data-theme=light])");
if (Object.keys(B.light).length === 0) missing.push("B light(:root)");
if (Object.keys(B.dark).length === 0) missing.push("B dark([data-theme=dark])");
if (missing.length > 0) {
  console.error(`token-diff: parsed empty theme block(s): ${missing.join(", ")} - parse broke, refusing to diff`);
  process.exit(1);
}

const aNames = [...new Set([...Object.keys(A.dark), ...Object.keys(A.light)])].sort();
const bNames = [...new Set([...Object.keys(B.light), ...Object.keys(B.dark)])].sort();
const onlyA = aNames.filter((n) => !bNames.includes(n));
const onlyB = bNames.filter((n) => !aNames.includes(n));
const both = aNames.filter((n) => bNames.includes(n));

const valueDiff = [];
const identical = [];
for (const n of both) {
  const diffs = [];
  for (const theme of ["dark", "light"]) {
    const av = A[theme][n];
    const bv = B[theme][n];
    if (av !== undefined && bv !== undefined && av !== bv) diffs.push(`${theme}: A=${av} B=${bv}`);
    else if (av === undefined && bv !== undefined) diffs.push(`${theme}: A has no declaration, B=${bv}`);
    else if (bv === undefined && av !== undefined) diffs.push(`${theme}: B has no declaration, A=${av}`);
  }
  if (diffs.length > 0) valueDiff.push({ name: n, diffs });
  else identical.push(n);
}

console.log("token-diff-77 (ticket 77 AC#10 / ticket 34 AC#10)");
console.log(`  A (legacy first generation): ${legacy.path}${legacy.direct ? "" : "  [design/assets/tokens.css is NOT in the repo - deleted in 72b65c614; this archive copy is the named reference, provenance header pins the deleted file's sha256]"}`);
console.log(`  B (shipped fourth generation): ${THEME_TOKENS}`);
console.log(`  names: A=${aNames.length} B=${bNames.length} both=${both.length} A-only=${onlyA.length} B-only=${onlyB.length}`);
console.log("");
console.log(`== names only in A (legacy): ${onlyA.length} ==`);
for (const n of onlyA) {
  console.log(`  ${n}`);
  console.log(`    dark : ${A.dark[n] ?? "(absent)"}`);
  console.log(`    light: ${A.light[n] ?? "(absent)"}`);
}
console.log("");
console.log(`== names only in B (shipped): ${onlyB.length} ==`);
for (const n of onlyB) {
  console.log(`  ${n}`);
  console.log(`    light: ${B.light[n] ?? "(absent)"}`);
  console.log(`    dark : ${B.dark[n] ?? "(absent)"}`);
}
console.log("");
console.log(`== names in both with differing values: ${valueDiff.length} ==`);
for (const { name, diffs } of valueDiff) {
  console.log(`  ${name}`);
  for (const d of diffs) console.log(`    ${d}`);
}
console.log("");
console.log(`== names in both with identical values: ${identical.length} ==`);
for (const n of identical) console.log(`  ${n}`);
console.log("");
console.log("NOTE: the two files are different DESIGN GENERATIONS. Equality is not");
console.log("the acceptance bar; the bar is this named diff plus the machine legs:");
console.log("  - generator <-> emitted theme: npm run tokens:check (run separately)");
console.log("  - C21 table <-> native <-> frontend fourth-gen legs: internal/panel");
console.log("    TestC21DesignTokensFourWayAgree (22 reconciled mapped pairs).");
