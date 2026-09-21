// C21 design-token generator (ticket 77 AC#2).
//
// Single source of truth: design/assets/tokens.css (contract C21 DesignTokens,
// "唯一样式真相源"). This script is the ONLY sanctioned way for a colour,
// radius, font-size or easing value to reach the frontend: it copies the CSS
// custom properties verbatim into src/styles/tokens.generated.css, which is
// committed so that `go build` and the four-way reconciliation test need no
// node toolchain.
//
// Hand-editing the generated file is forbidden - internal/panel's
// TestC21DesignTokensFourWayAgree re-parses both sides and fails on drift, and
// CI's lint-frontend job re-runs this generator and diffs the tree.
//
// Usage:  node scripts/gen-tokens.mjs [--check]
//   --check  regenerate in memory and exit 1 if the committed file differs.
import { readFileSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join, posix } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = join(here, "..");
const repoRoot = join(frontendRoot, "..");

const TOKENS_CSS = join(repoRoot, "design", "assets", "tokens.css");
const OUT_REL = "src/styles/tokens.generated.css";
const OUT = join(frontendRoot, "src", "styles", "tokens.generated.css");

/** Remove slash-star block comments, keeping declaration boundaries intact. */
function stripComments(text) {
  return text.replace(/\/\*[\s\S]*?\*\//g, "");
}

/**
 * Parse `:root { ... }` and `[data-theme="light"] { ... }` into ordered
 * { name, value } declaration lists. Values are kept verbatim (minus the
 * trailing `;`), including `var(...)` references and multi-layer backgrounds.
 */
function parseBlocks(text) {
  const body = stripComments(text);
  const blocks = new Map();
  const ruleRe = /(:root|\[data-theme="light"\])\s*\{([^}]*)\}/g;
  let m;
  while ((m = ruleRe.exec(body)) !== null) {
    const key = m[1] === ":root" ? "dark" : "light";
    const decls = [];
    for (const chunk of m[2].split(";")) {
      const decl = chunk.trim();
      if (!decl) continue;
      const idx = decl.indexOf(":");
      if (idx < 0) continue;
      const name = decl.slice(0, idx).trim();
      const value = decl.slice(idx + 1).trim();
      if (!name.startsWith("--")) continue;
      decls.push({ name, value });
    }
    blocks.set(key, decls);
  }
  if (!blocks.has("dark")) throw new Error("no :root block in " + posix.join("design", "assets", "tokens.css"));
  return blocks;
}

/** A value is a plain colour when it is a hex or an rgb()/rgba() call. */
function isPlainColour(value) {
  return /^#[0-9A-Fa-f]{3,8}$/.test(value) || /^rgba?\(/.test(value);
}

/**
 * C21 layers some backgrounds as `var(--sheen), rgba(...)`: a sheen gradient on
 * top of a colour. The native side (internal/ball/tokens.go) copies only the
 * colour component - the C21 table states this out loud for --bg-overlay. The
 * frontend needs the same slice, because a CSS utility that sets
 * `background-color` cannot carry a two-layer background stack. So for every
 * composite value we derive a `<name>-color` token holding the LAST colour
 * literal found in the parent value. Nothing is invented here: the extraction
 * is re-computed by TestC21DesignTokensFourWayAgree in Go, which fails if this
 * file and tokens.css disagree.
 */
const COLOUR_IN_VALUE_RE = /rgba?\([^)]*\)|#[0-9A-Fa-f]{3,8}/g;

function colourSlice(value) {
  const found = flat(value).match(COLOUR_IN_VALUE_RE);
  if (!found || found.length === 0) return null;
  return found[found.length - 1];
}

/** Flatten whitespace inside a value so the emitted file stays one line per token. */
function flat(value) {
  return value.replace(/\s+/g, " ").trim();
}

function emitSelector(sel, decls) {
  const lines = [`${sel} {`];
  for (const d of expand(decls)) lines.push(`  ${d.name}: ${d.value};`);
  lines.push("}", "");
  return lines.join("\n");
}

/**
 * The declared tokens plus the derived `<name>-color` slices, in source order.
 * emitSelector writes this list verbatim and emitTheme maps it, so both sides
 * read exactly one list.
 */
function expand(decls) {
  const out = [];
  for (const d of decls) {
    const value = flat(d.value);
    out.push({ name: d.name, value });
    if (!isPlainColour(value)) {
      const slice = colourSlice(value);
      if (slice) out.push({ name: `${d.name}-color`, value: slice, derived: true });
    }
  }
  return out;
}

/**
 * Tailwind v4 `@theme` mapping. Only namespaces Tailwind actually owns, and
 * only from tokens declared above: plain colours, the --r-* radius scale and the
 * two font stacks. Composite values (sheen background stacks, shadows, blur)
 * stay raw custom properties and are referenced through var(...) by the
 * component layer. Inventing a value here is forbidden.
 */
function emitTheme(dark) {
  const lines = ["@theme inline {"];
  for (const d of expand(dark)) {
    const key = d.name.slice(2);
    if (isPlainColour(d.value)) lines.push(`  --color-${key}: var(${d.name});`);
  }
  for (const d of expand(dark)) {
    const key = d.name.slice(2);
    if (key.startsWith("r-") && key !== "r-full") lines.push(`  --radius-${key.slice(2)}: var(${d.name});`);
    if (key === "font-sans" || key === "font-mono") lines.push(`  --${key}: var(${d.name});`);
  }
  lines.push("}", "");
  return lines.join("\n");
}



const HEADER = `/* ============================================================================
   GENERATED FILE - DO NOT EDIT BY HAND  (ticket 77 AC#2)
   ----------------------------------------------------------------------------
   Source:  design/assets/tokens.css  (contract C21 DesignTokens, the ONLY
            style truth source of this repo)
   Command: node scripts/gen-tokens.mjs        (npm run tokens)
   Checked: internal/panel's TestC21DesignTokensFourWayAgree reconciles this
            file against tokens.css, docs/evidence/s1/c21-native-tokens.md and
            internal/ball/tokens.go - four parties, one value per token.

   A colour, radius, font-size or duration literal written anywhere else under
   frontend/src is a second style truth source and fails the same test.
   ============================================================================ */

`;

const blocks = parseBlocks(readFileSync(TOKENS_CSS, "utf8"));
const dark = blocks.get("dark");
const light = blocks.get("light") ?? [];

const out =
  HEADER +
  emitSelector(":root", dark) +
  (light.length ? emitSelector('[data-theme="light"]', light) : "") +
  emitTheme(dark);

if (process.argv.includes("--check")) {
  let current = "";
  try {
    current = readFileSync(OUT, "utf8");
  } catch {
    current = "";
  }
  if (current !== out) {
    console.error(
      "gen-tokens: " +
        OUT_REL +
        " is not what design/assets/tokens.css generates - run `npm run tokens` and commit the result",
    );
    process.exit(1);
  }
  console.log(`gen-tokens: ${OUT_REL} matches design/assets/tokens.css (${dark.length} dark + ${light.length} light declarations)`);
} else {
  writeFileSync(OUT, out, "utf8");
  console.log(`gen-tokens: wrote ${OUT_REL} (${dark.length} dark + ${light.length} light declarations)`);
}
