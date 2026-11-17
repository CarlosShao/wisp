#!/usr/bin/env node
/**
 * lint-seeds.mjs - ticket 77 AC#9 seed scanner (zero new dependencies).
 *
 * AC#9 (merged from ticket 34): the CI job must not merely run green - a
 * hardcoded colour, a UI emoji and a banned import must each be CATCHABLE,
 * proven by seeded violations that ring. This file is both the production
 * scan of frontend/src and its own positive/negative control harness:
 *
 *   node scripts/lint-seeds.mjs              scan frontend/src (CI gate)
 *   node scripts/lint-seeds.mjs --self-test  controls: clean sample must stay
 *                                            green, one seeded violation per
 *                                            family must ring (CI runs BOTH)
 *
 * The three families and their exact rules:
 *
 *   1. hex      - a hardcoded #rgb/#rgba/#rrggbb/#rrggbbaa colour literal in
 *                 any .ts/.tsx/.css under src. Named exclusion, exactly one:
 *                 src/styles/tokens.generated.css, the generated C21 token
 *                 single source (guarded independently by `tokens:check` and
 *                 TestC21DesignTokensFourWayAgree). Every other file must
 *                 reference a token (var(--...)), mirroring the Go-side
 *                 TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme, but
 *                 without that gate's reachability restriction: this scan
 *                 reads every src file, reachable or not.
 *   2. emoji    - the tools/d22scan emojiRe character class (ban #8),
 *                 COPIED VERBATIM, with the comment exemption deliberately
 *                 absent: raw lines only, comments and strings alike. If the
 *                 Go spec wording and the instrument ever disagree, the
 *                 product rule wins: UI zero emoji (dispatch 2026-11-17).
 *   3. imports  - two forbidden shapes, both about demo/fixture data staying
 *                 out of the product face (P9: no fake data behind real UI):
 *                   a. an import whose specifier resolves OUTSIDE src into a
 *                      fixtures/ or demo/ tree (e.g. the render-fixture JSONs
 *                      in frontend/fixtures/ must never be bundled);
 *                   b. an import of src/fixtures/** (the harness demo data)
 *                      from any file that is not the harness face. The
 *                      harness face is a stated scope, not a hidden exemption:
 *                      src/main.tsx, src/components/harness/**, and
 *                      src/components/showcase.tsx - exactly the importers
 *                      that existed when this rule landed.
 *
 * Exit codes (the repo's convention, a gate must not pass by matching
 * nothing): 0 = clean, 1 = findings / a control misbehaved, 2 = cannot
 * conclude (empty denominator or src missing).
 *
 * NOTE: this file sits inside TestFrontendHasNoEmoji's scan set
 * (frontend/src + frontend/scripts), so the emoji sample character below is
 * built with String.fromCodePoint at runtime - the source text itself must
 * stay outside ban #8's bands.
 */

import { readdirSync, readFileSync, existsSync } from "node:fs";
import { dirname, join, relative, resolve, sep } from "node:path";
import { posix } from "node:path";
import { fileURLToPath } from "node:url";

const FRONTEND_ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const SRC_ROOT = join(FRONTEND_ROOT, "src");

// ---- family 1: hardcoded hex colour ---------------------------------------
const HEX_RE = /#[0-9a-fA-F]{3,8}\b|#[0-9a-fA-F]{9,}/g;
const HEX_EXT = new Set([".ts", ".tsx", ".css"]);
// The one named exclusion: the generated token single source.
const HEX_EXCLUDED = new Set(["src/styles/tokens.generated.css"]);

// ---- family 2: UI emoji (d22scan emojiRe, verbatim, no comment exemption) --
// Built through new RegExp from an escaped string: the bands contain FE0F, a
// combining variation selector, and oxlint's no-misleading-character-class
// flags the same character class spelled as a literal (the identical warning
// already sits on vendor*.mjs). Semantics are identical to the Go class.
const EMOJI_BANDS =
  "\\u{1F000}-\\u{1FAFF}\\u{2200}-\\u{22FF}\\u{2600}-\\u{27BF}\\u{2B00}-\\u{2BFF}\\u{FE0F}\\u{1F1E6}-\\u{1F1FF}";
const EMOJI_RE = new RegExp("[" + EMOJI_BANDS + "]", "gu");

// ---- family 3: forbidden imports ------------------------------------------
const SPEC_RES = [
  /\bfrom\s*["']([^"']+)["']/g,
  /\bimport\s*\(\s*["']([^"']+)["']\s*[),]/g,
  /(?:^|[;\s(])import\s*["']([^"']+)["']/g,
];
const HARNESS_FACE = [
  (rel) => rel === "src/main.tsx",
  (rel) => rel.startsWith("src/components/harness/"),
  (rel) => rel === "src/components/showcase.tsx",
];

function toPosix(p) {
  return p.split(sep).join("/");
}

function listFiles(dir) {
  const out = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) out.push(...listFiles(full));
    else out.push(full);
  }
  return out;
}

function lineOf(text, index) {
  let line = 1;
  for (let i = 0; i < index; i++) {
    if (text.charCodeAt(i) === 10) line++;
  }
  return line;
}

function resolveSpecifier(spec, importerRel) {
  if (spec.startsWith("@/")) {
    return posix.normalize("src/" + spec.slice(2));
  }
  if (spec.startsWith("./") || spec.startsWith("../")) {
    return posix.normalize(posix.join(posix.dirname(importerRel), spec));
  }
  if (spec.startsWith("fixtures/") || spec.startsWith("demo/")) {
    return posix.normalize(spec);
  }
  return null; // bare package specifier: not this gate's business
}

function importFinding(importerRel, spec) {
  const target = resolveSpecifier(spec, importerRel);
  if (!target) return null;
  const segs = target.split("/");
  const namesAFixtureTree = segs.includes("fixtures") || segs.includes("demo");
  const insideSrc = target === "src" || target.startsWith("src/");
  if (namesAFixtureTree && !insideSrc) {
    return (
      'forbidden-import: "' + spec + '" resolves outside src into a ' +
      "fixtures/demo tree (" + target + ") - render-fixture and demo data must " +
      "never be bundled into the product face"
    );
  }
  if (
    insideSrc &&
    target.startsWith("src/fixtures/") &&
    !HARNESS_FACE.some((allowed) => allowed(importerRel))
  ) {
    return (
      'forbidden-import: "' + importerRel + '" imports the harness fixture "' +
      spec + '" (' + target + ") but is not the harness face (src/main.tsx, " +
      "src/components/harness/**, src/components/showcase.tsx)"
    );
  }
  return null;
}

/**
 * Scan one in-memory file. The real tree and the self-test samples both go
 * through this exact function, so the controls exercise the same rules CI
 * consumes.
 */
function scanFile(rel, text, ext) {
  const findings = [];
  const lines = text.split("\n");
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    // family 2 first: every file, every line, comments included.
    EMOJI_RE.lastIndex = 0;
    const emoji = [...line.matchAll(EMOJI_RE)];
    if (emoji.length > 0) {
      const cps = emoji
        .map((m) => "U+" + m[0].codePointAt(0).toString(16).toUpperCase().padStart(4, "0"))
        .join(" ");
      findings.push({
        family: "emoji",
        rel,
        line: i + 1,
        detail: "UI emoji character present (d22scan emojiRe bands, no comment exemption): " + cps,
      });
    }
    if (HEX_EXT.has(ext) && !HEX_EXCLUDED.has(rel)) {
      HEX_RE.lastIndex = 0;
      const hex = [...line.matchAll(HEX_RE)];
      for (const m of hex) {
        findings.push({
          family: "hex",
          rel,
          line: i + 1,
          detail: 'hardcoded colour literal "' + m[0] + '" - reference a C21 token (var(--...)) instead',
        });
      }
    }
  }
  if (ext === ".ts" || ext === ".tsx") {
    const seen = new Set();
    for (const re of SPEC_RES) {
      re.lastIndex = 0;
      for (const m of text.matchAll(re)) {
        const spec = m[1];
        const line = lineOf(text, m.index);
        const key = spec + "@" + line;
        if (seen.has(key)) continue;
        seen.add(key);
        const detail = importFinding(rel, spec);
        if (detail) findings.push({ family: "forbidden-import", rel, line, detail });
      }
    }
  }
  return findings;
}

function scanTree() {
  if (!existsSync(SRC_ROOT)) {
    return { error: "frontend/src does not exist - the scan cannot conclude" };
  }
  const files = listFiles(SRC_ROOT).sort();
  const findings = [];
  let emojiFiles = 0;
  let hexFiles = 0;
  let importFiles = 0;
  for (const full of files) {
    const rel = toPosix(relative(FRONTEND_ROOT, full));
    const ext = posix.extname(rel);
    const text = readFileSync(full, "utf8");
    emojiFiles++;
    if (HEX_EXT.has(ext) && !HEX_EXCLUDED.has(rel)) hexFiles++;
    if (ext === ".ts" || ext === ".tsx") importFiles++;
    findings.push(...scanFile(rel, text, ext));
  }
  if (emojiFiles === 0 || hexFiles === 0 || importFiles === 0) {
    return { error: "empty denominator (emoji=" + emojiFiles + " hex=" + hexFiles + " imports=" + importFiles + ") - refusing to conclude" };
  }
  return { findings, emojiFiles, hexFiles, importFiles };
}

// ---- positive/negative controls (seeded violations, one per family) -------
const CHECK_MARK = String.fromCodePoint(0x2713); // U+2713, inside emojiRe's 2600-27BF band

function sample(name, rel, text, expect) {
  return { name, rel, text, expect };
}

const SAMPLES = [
  sample(
    "positive control (clean file must stay green)",
    "src/components/seed-clean-sample.tsx",
    [
      'import { type PanelSnapshot } from "@/lib/panel";',
      "",
      "export function SeedCleanRow() {",
      '  return <div className="bg-surface text-ink" style={{ color: "var(--accent)" }}>token-fed</div>;',
      "}",
    ].join("\n"),
    { green: true },
  ),
  sample(
    "negative control 1 - hardcoded hex colour must ring",
    "src/components/seed-bad-hex.tsx",
    [
      "export function SeedBadHex() {",
      '  const brand = "#0285ff";',
      '  return <div style={{ color: brand }}>brand</div>;',
      "}",
    ].join("\n"),
    { family: "hex" },
  ),
  sample(
    "negative control 2 - UI emoji must ring, even inside a COMMENT (the d22scan comment exemption is deliberately absent)",
    "src/components/seed-bad-emoji.tsx",
    [
      "// the done marker on this line is U+2713: " + CHECK_MARK,
      'export const SeedBadEmoji = "not a real component";',
    ].join("\n"),
    { family: "emoji" },
  ),
  sample(
    "negative control 3 - forbidden import must ring (product-face file importing src/fixtures)",
    "src/components/seed-bad-import.tsx",
    [
      'import { HARNESS_SNAPSHOT } from "@/fixtures/harness";',
      "",
      "export const SeedBadImport = HARNESS_SNAPSHOT;",
    ].join("\n"),
    { family: "forbidden-import" },
  ),
];

function runSelfTest() {
  let failed = 0;
  for (const s of SAMPLES) {
    const findings = scanFile(s.rel, s.text, posix.extname(s.rel));
    if (s.expect.green) {
      if (findings.length === 0) {
        console.log("PASS (green) " + s.name + " -> 0 findings");
      } else {
        failed++;
        console.log("FAIL " + s.name + " -> the clean sample rang " + findings.length + " finding(s):");
        for (const f of findings) console.log("      " + f.rel + ":" + f.line + " [" + f.family + "] " + f.detail);
      }
      continue;
    }
    const hit = findings.find((f) => f.family === s.expect.family && f.rel === s.rel);
    if (hit) {
      console.log("RED  (caught) " + s.name + " -> " + hit.rel + ":" + hit.line + " [" + hit.family + "] " + hit.detail);
    } else {
      failed++;
      console.log("FAIL " + s.name + " -> no " + s.expect.family + " finding named " + s.rel + " (scanner went soft)");
    }
  }
  console.log(
    "self-test verdict: " + (failed === 0 ? "positive control green + 3 seeded violations red - the scanner has teeth" : failed + " control(s) FAILED - the scanner must not be trusted"),
  );
  return failed === 0 ? 0 : 1;
}

function runScan() {
  const result = scanTree();
  if (result.error) {
    console.error("lint-seeds: " + result.error);
    return 2;
  }
  console.log("lint-seeds: scan frontend/src (ticket 77 AC#9)");
  console.log(
    "denominators: emoji files=" + result.emojiFiles +
    " | hex files=" + result.hexFiles + " (excluded by name: src/styles/tokens.generated.css)" +
    " | import files=" + result.importFiles,
  );
  if (result.findings.length === 0) {
    console.log("findings: 0");
    return 0;
  }
  console.log("findings: " + result.findings.length);
  for (const f of result.findings) {
    console.log("  " + f.rel + ":" + f.line + " [" + f.family + "] " + f.detail);
  }
  return 1;
}

const mode = process.argv[2];
if (mode === "--self-test") {
  process.exit(runSelfTest());
} else if (mode === undefined) {
  process.exit(runScan());
} else {
  console.error("usage: node scripts/lint-seeds.mjs [--self-test]");
  process.exit(2);
}
