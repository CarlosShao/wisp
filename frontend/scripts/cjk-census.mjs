#!/usr/bin/env node
/**
 * cjk-census.mjs - ticket 58 AC#1 census (zero new dependencies).
 *
 * SURVEY ONLY - this is not a gate. The AC#1 gate ("zero hardcoded Chinese in
 * frontend render paths") lives in internal/panel/frontend_cjk_render_test.go,
 * which scans the same closure with the same bands. This script exists so the
 * census numbers (per-file literal counts, top lists) can be re-derived at any
 * time without a Go toolchain:
 *
 *   node scripts/cjk-census.mjs
 *
 * Semantics (fixed so two runs can be compared):
 *   - scan set   : frontend/src/**\/*.{ts,tsx}; each file is also marked whether
 *                  it sits in the import closure of src/main.tsx (the files the
 *                  shipped bundle actually carries - the same reachability rule
 *                  TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme uses).
 *   - bands      : CJK ideographs U+3400-4DBF, U+4E00-9FFF, U+F900-FAFF
 *                  (compatibility), kana U+3040-30FF, Hangul U+AC00-D7AF.
 *                  Fullwidth punctuation (，。：) is NOT a band member: the
 *                  census counts characters, and any string carrying real copy
 *                  carries ideographs too.
 *   - comments   : stripped first (state machine over //, slash-star, and the
 *                  three string quotes), so documented Chinese never counts.
 *   - literals   : per string-literal occurrence whose stripped content holds a
 *                  band character ("枚数").
 *   - jsx text   : per LINE holding a band character outside any literal, in
 *                  .tsx only (JSX text nodes cannot be quoted - counted by
 *                  line, and a line is the unit the gate reds on too).
 */

import { readdirSync, readFileSync, existsSync } from "node:fs";
import { dirname, join, relative, resolve, sep } from "node:path";
import { posix } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const FRONTEND_ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const SRC_ROOT = join(FRONTEND_ROOT, "src");

const CJK_RE = /[\u3400-\u4DBF\u4E00-\u9FFF\uF900-\uFAFF\u3040-\u30FF\uAC00-\uD7AF]/;

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

/**
 * Strip comments while respecting string literals. Returns { code, inComment }
 * where code has comments replaced by spaces (line structure preserved).
 * Regex literals are not modeled: a `/` that opens a regex can at worst make
 * the stripper mis-read the rest of that one line, and both the census and the
 * gate count line-level CJK afterwards, so no CJK occurrence can be lost by
 * this simplification (it may move a count between the literal and jsx-text
 * columns - the total is what the ticket's number is).
 */
export function stripComments(text) {
  let out = "";
  let state = "code"; // code | line | block | squote | dquote | tick
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    const n = i + 1 < text.length ? text[i + 1] : "";
    if (state === "code") {
      if (c === "/" && n === "/") { state = "line"; out += "  "; i++; continue; }
      if (c === "/" && n === "*") { state = "block"; out += "  "; i++; continue; }
      if (c === '"' || c === "'" || c === "`") { state = c; out += c; continue; }
      out += c;
      continue;
    }
    if (state === "line") {
      if (c === "\n") { state = "code"; out += c; } else out += " ";
      continue;
    }
    if (state === "block") {
      if (c === "*" && n === "/") { state = "code"; out += "  "; i++; }
      else out += c === "\n" ? "\n" : " ";
      continue;
    }
    // inside a string literal
    if (c === "\\") { out += "  "; i++; continue; }
    if (c === state) { state = "code"; out += c; continue; }
    out += c;
  }
  return out;
}

/** Count band characters per stripped-comment text, split literal vs outside. */
export function analyze(code, isTsx) {
  const lines = code.split("\n");
  let literals = 0;
  let jsxTextLines = 0;
  const LIT_RE = /"(?:[^"\\\n]|\\.)*"|'(?:[^'\\\n]|\\.)*'|`(?:[^`\\]|\\.)*`/g;
  for (const line of lines) {
    let outside = line;
    let litHits = 0;
    for (const m of line.matchAll(LIT_RE)) {
      if (CJK_RE.test(m[0])) litHits++;
      outside = outside.replace(m[0], "");
    }
    literals += litHits;
    if (isTsx && CJK_RE.test(outside)) jsxTextLines++;
  }
  return { literals, jsxTextLines };
}

// ---- import closure of main.tsx (the files the bundle ships) ---------------
const importFromRe = /(?:from|import)\s*\(?\s*["']([^"']+)["']/g;
const cssImportRe = /@import\s+["']([^"']+)["']/g;

function resolveSpec(srcRoot, from, spec) {
  let base;
  if (spec.startsWith("@/")) base = join(srcRoot, filepathFromSlash(spec.slice(2)));
  else if (spec.startsWith("./") || spec.startsWith("../")) base = join(dirname(from), filepathFromSlash(spec));
  else return [];
  const cands = [base];
  for (const e of [".ts", ".tsx", ".css", ".js", ".jsx"]) cands.push(base + e);
  if (!spec.endsWith(".css")) {
    cands.push(join(base, "index.tsx"), join(base, "index.ts"));
  }
  return cands;
}

function filepathFromSlash(p) {
  return p.split("/").join(sep);
}

export function closureOfMain(srcRoot = SRC_ROOT) {
  const seen = new Set();
  const visit = (path) => {
    if (seen.has(path)) return;
    seen.add(path);
    let text;
    try { text = readFileSync(path, "utf8"); } catch { return; }
    const specs = [];
    for (const m of text.matchAll(importFromRe)) specs.push(m[1]);
    for (const m of text.matchAll(cssImportRe)) specs.push(m[1]);
    for (const spec of specs) {
      for (const cand of resolveSpec(SRC_ROOT, path, spec)) {
        if (existsSync(cand)) visit(cand);
      }
    }
  };
  visit(join(srcRoot, "main.tsx"));
  return seen;
}

function run() {
  if (!existsSync(SRC_ROOT)) {
    console.error("cjk-census: frontend/src does not exist");
    return 2;
  }
  const closure = closureOfMain();
  const files = listFiles(SRC_ROOT)
    .map((f) => ({ full: f, rel: toPosix(relative(FRONTEND_ROOT, f)) }))
    .filter((f) => /\.tsx?$/.test(f.rel))
    .sort((a, b) => a.rel.localeCompare(b.rel));
  const rows = [];
  for (const f of files) {
    const text = readFileSync(f.full, "utf8");
    const code = stripComments(text);
    const { literals, jsxTextLines } = analyze(code, f.rel.endsWith(".tsx"));
    rows.push({
      rel: f.rel,
      inClosure: closure.has(f.full),
      literals,
      jsxTextLines,
      total: literals + jsxTextLines,
    });
  }
  const inC = rows.filter((r) => r.inClosure);
  const outC = rows.filter((r) => !r.inClosure);
  const sum = (rs) => rs.reduce((a, r) => a + r.total, 0);
  console.log("cjk-census: frontend/src/**/*.{ts,tsx}  bands=ideographs+kana+hangul, comments stripped");
  console.log(
    "closure of src/main.tsx: files=" + inC.length + " total=" + sum(inC) +
    " (literals=" + inC.reduce((a, r) => a + r.literals, 0) +
    " jsxTextLines=" + inC.reduce((a, r) => a + r.jsxTextLines, 0) + ")",
  );
  console.log(
    "outside closure        : files=" + outC.length + " total=" + sum(outC) +
    " (unmounted demos/fixtures - not the render path)",
  );
  console.log("-- top 15 by total (closure first) --");
  for (const r of [...rows].sort((a, b) => b.total - a.total || a.rel.localeCompare(b.rel)).slice(0, 15)) {
    console.log(
      "  " + String(r.total).padStart(4) + "  lit=" + String(r.literals).padStart(4) +
      "  jsx=" + String(r.jsxTextLines).padStart(3) +
      "  " + (r.inClosure ? "SHIP " : "unmnt") + "  " + r.rel,
    );
  }
  console.log("-- closure files with total>0, descending --");
  for (const r of inC.filter((r) => r.total > 0).sort((a, b) => b.total - a.total)) {
    console.log(
      "  " + String(r.total).padStart(4) + "  lit=" + String(r.literals).padStart(4) +
      "  jsx=" + String(r.jsxTextLines).padStart(3) + "  " + r.rel,
    );
  }
  return 0;
}

const invokedDirectly =
  process.argv[1] !== undefined && import.meta.url === pathToFileURL(process.argv[1]).href;
if (invokedDirectly) process.exit(run());
