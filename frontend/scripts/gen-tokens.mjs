// Wisp design-token generator - SECOND GENERATION (owner ruling 2026-09-25).
//
// Single source of truth: design/doubao/demo/styles.css.
//
// The first generation copied design/assets/tokens.css verbatim. Owner retired
// that table out loud on 2026-09-25 ("全推翻了，前端重新设计 token ... 按照 demo
// 改") and moved the whole directory under design/old/, so this script no longer
// reads it. What replaced it is a two-layer table:
//
//   layer 1  --demo-*   every custom property demo declares, copied verbatim (a
//                       bare shadcn HSL triplet is rewritten as hsl(...) - a
//                       format change, not a value change)
//   layer 2  semantic   the names frontend/src actually speaks, each written over
//                       those primitives, or as a literal this file PROVES the
//                       demo sheet draws
//
// The proof is the point. `cite()` looks the claimed source text up in
// demo/styles.css and stamps the line number it was found on; `verify()` refuses
// an entry that has neither a citation nor an INTERIM marker, so a value demo
// does not draw cannot enter the table without being labelled as somebody's
// guess. That is what keeps "1:1 移植 demo" from degrading into "whoever edited
// this last liked it".
//
// Hand-editing src/styles/tokens.generated.css is still forbidden.
//
// Usage:  node scripts/gen-tokens.mjs [--check]
//   --check  regenerate in memory, exit 1 if the committed file differs, if a
//            citation has gone stale, or if any var(--token) used under src/
//            resolves to nothing.
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join, posix, relative } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = join(here, "..");
const repoRoot = join(frontendRoot, "..");

const DEMO_REL = "design/doubao/demo/styles.css";
const DEMO = join(repoRoot, "design", "doubao", "demo", "styles.css");
const OUT_REL = "src/styles/tokens.generated.css";
const OUT = join(frontendRoot, "src", "styles", "tokens.generated.css");
const SRC = join(frontendRoot, "src");

const demo = readFileSync(DEMO, "utf8");

/** 1-based line where `text` appears in the demo sheet; throws when it does not. */
function lineOf(text) {
  const idx = demo.indexOf(text);
  if (idx < 0) throw new Error(`citation gone stale, not in ${DEMO_REL}: ${JSON.stringify(text)}`);
  return demo.slice(0, idx).split("\n").length;
}

/** "styles.css:120 + styles.css:122" for the claims an entry rests on. */
function cite(claims) {
  return claims.map((t) => `styles.css:${lineOf(t)}`).join(" + ");
}

/** Interim marker: demo draws no such value, so the entry names what it aliases. */
const NO_DEMO = "INTERIM(无 demo 对应值)";

function flat(value) {
  return value.replace(/\s+/g, " ").trim();
}

function stripComments(text) {
  return text.replace(/\/\*[\s\S]*?\*\//g, "");
}

/** The declarations of the first rule whose selector is exactly `sel`. */
function block(sel) {
  const body = stripComments(demo);
  const esc = sel.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const m = body.match(new RegExp(`(?:^|\\})\\s*${esc}\\s*\\{([^}]*)\\}`, "m"));
  if (!m) throw new Error(`no ${sel} block in ${DEMO_REL}`);
  const decls = [];
  for (const chunk of m[1].split(";")) {
    const decl = chunk.trim();
    if (!decl) continue;
    const idx = decl.indexOf(":");
    if (idx < 0) continue;
    const name = decl.slice(0, idx).trim();
    if (!name.startsWith("--")) continue;
    decls.push({ name, value: flat(decl.slice(idx + 1)) });
  }
  return decls;
}

/**
 * demo stores colours the shadcn way: a bare HSL triplet in the variable, the
 * call spelled out at the point of use. A table consumed through Tailwind needs
 * a complete colour, so the triplet is wrapped here - the same three numbers,
 * one function call earlier. Nothing is recomputed and nothing becomes hex.
 */
const TRIPLET = /^-?[\d.]+(deg)? [\d.]+% [\d.]+%$/;
function toColour(value) {
  return TRIPLET.test(value) ? `hsl(${value})` : value;
}

/**
 * demo's `:root` is LIGHT and `.dark` overrides it; demo/app.js:32 ships
 * `theme: 'light'`. That is the opposite way round from the first generation,
 * so the emitted file follows demo: `:root` is light, dark is the opt-in.
 */
const lightBlock = block(":root");
const darkBlock = block(".dark");

function primitives(decls) {
  return decls.map((d) => ({
    name: `--demo-${d.name.slice(2)}`,
    value: toColour(d.value),
    cite: cite([`${d.name}: ${d.value}`]),
  }));
}

const lightPrims = primitives(lightBlock);
const darkPrims = primitives(darkBlock);
const primNames = new Set([...lightPrims, ...darkPrims].map((p) => p.name));

/**
 * Take demo's alpha off one of its rgba() literals. demo paints the window
 * rgba(255,255,255,0.72) light / rgba(22,24,29,0.72) dark; an opacity control
 * needs the hue without the number welded to it, so the alpha can be a variable.
 * Same three channels, one argument removed.
 */
function hueOf(rgbaLiteral) {
  const m = rgbaLiteral.match(/^rgba\(([^,]+), *([^,]+), *([^,]+), *([0-9.]+)\)$/);
  if (!m) throw new Error(`hueOf: not an rgba() literal: ${JSON.stringify(rgbaLiteral)}`);
  lineOf(rgbaLiteral);
  return { hue: `rgb(${m[1]},${m[2]},${m[3]})`, alpha: `${Number(m[4]) * 100}%` };
}

const WINDOW_LIGHT = hueOf("rgba(255, 255, 255, 0.72)");
const WINDOW_DARK = hueOf("rgba(22, 24, 29, 0.72)");

// ---------------------------------------------------------------------------
// layer 2 - the vocabulary frontend/src speaks
// ---------------------------------------------------------------------------

/**
 * Entry shape:
 *   name    the custom property
 *   value   theme-invariant expression over --demo-* (var() indirection means
 *           it picks up whichever theme the primitives hold)
 *   light / dark   per-theme literal, for the handful demo itself changes
 *   from    the demo text this rests on - each claim must be findable
 *   interim what it stands in for, when demo draws no such thing
 *   kind    "colour" (default) or something else; only colours get a --color-*
 *           entry, because that is the namespace Tailwind reads
 */
const SEMANTIC = [
  // ---- surfaces ----------------------------------------------------------
  { name: "--bg-base", value: "var(--demo-background)", from: ["--background:"] },
  { name: "--bg-subtle", value: "var(--demo-secondary)", from: ["--secondary:"] },
  { name: "--bg-inset", value: "var(--demo-muted)", from: ["--muted:"] },
  { name: "--bg-raised", value: "var(--demo-card)", from: ["--card:"] },
  { name: "--bg-overlay", value: "var(--demo-popover)", from: ["--popover:"] },
  {
    name: "--panel-hue",
    light: WINDOW_LIGHT.hue,
    dark: WINDOW_DARK.hue,
    from: ["rgba(255, 255, 255, 0.72)", "rgba(22, 24, 29, 0.72)"],
    note: "demo's window colour with its own alpha taken off, so the alpha can move",
  },
  {
    name: "--panel-alpha",
    value: WINDOW_LIGHT.alpha,
    from: ["rgba(255, 255, 255, 0.72)"],
    kind: "other",
    note: "0.72 written as a percentage; the settings row is the only thing allowed to move it",
  },
  {
    name: "--panel-surface",
    value: "color-mix(in srgb, var(--panel-hue) var(--panel-alpha), transparent)",
    from: ["--window-bg:"],
  },

  // ---- ink ---------------------------------------------------------------
  { name: "--fg-primary", value: "var(--demo-foreground)", from: ["--foreground:"] },
  { name: "--fg-secondary", value: "var(--demo-muted-foreground)", from: ["--muted-foreground:"] },
  {
    name: "--fg-tertiary",
    value: "color-mix(in srgb, var(--demo-muted-foreground) 70%, transparent)",
    interim: NO_DEMO + ": demo 只有一级灰 --muted-foreground，这一级是它再淡三成",
  },
  {
    name: "--fg-disabled",
    value: "color-mix(in srgb, var(--demo-muted-foreground) 40%, transparent)",
    interim: NO_DEMO + ": demo 表达禁用写的是 style=opacity 0.4，不是另一枚颜色",
  },
  {
    name: "--on-solid",
    value: "var(--demo-primary-foreground)",
    from: ["--primary-foreground:"],
    note: "同一枚白两个名字：面板词表叫 on-solid，shadcn 词表叫 primary-foreground",
  },

  // ---- accent ------------------------------------------------------------
  { name: "--accent", value: "var(--demo-primary)", from: ["--primary:"] },
  { name: "--accent-fg", value: "var(--demo-primary-foreground)", from: ["--primary-foreground:"] },
  { name: "--accent-soft", value: "var(--demo-accent)", from: ["--accent:"] },
  { name: "--accent-line", value: "color-mix(in srgb, var(--demo-primary) 30%, transparent)", from: ["hsl(var(--primary) / 0.3)"] },
  {
    name: "--accent-hover",
    value: "color-mix(in srgb, var(--demo-primary) 88%, var(--demo-foreground))",
    interim: NO_DEMO + ": demo 的按钮悬停改的是阴影，不是色阶",
  },
  {
    name: "--accent-pressed",
    value: "color-mix(in srgb, var(--demo-primary) 80%, var(--demo-foreground))",
    interim: NO_DEMO + ": 同上，只按『压得比悬停重一档』排",
  },

  // ---- status ------------------------------------------------------------
  { name: "--danger", value: "var(--demo-destructive)", from: ["--destructive:"] },
  { name: "--danger-soft", value: "color-mix(in srgb, var(--demo-destructive) 10%, transparent)", from: ["hsl(var(--destructive) / 0.1)"] },
  { name: "--danger-line", value: "color-mix(in srgb, var(--demo-destructive) 30%, transparent)", from: ["hsl(var(--destructive) / 0.3)"] },
  { name: "--success", light: "hsl(142 71% 45%)", dark: "#7DB896", from: ["hsl(142 71% 45%)", "color: #7DB896"] },
  { name: "--success-soft", value: "color-mix(in srgb, var(--success) 10%, transparent)", from: ["hsl(142 71% 45% / 0.1)"] },
  { name: "--success-line", value: "color-mix(in srgb, var(--success) 30%, transparent)", from: ["hsl(142 71% 45% / 0.3)"] },
  { name: "--warn", light: "#b8860b", dark: "#D9B26A", from: ["color: #b8860b", "color: #D9B26A"] },
  { name: "--warn-soft", value: "color-mix(in srgb, var(--warn) 10%, transparent)", from: ["rgba(217, 178, 106, 0.1)"] },
  { name: "--warn-line", value: "color-mix(in srgb, var(--warn) 30%, transparent)", from: ["rgba(217, 178, 106, 0.3)"] },
  { name: "--warm", value: "#f78c6c", from: ["#f78c6c"], note: "demo 唯一的暖色在代码块高亮里，沿用它的色相" },
  { name: "--warm-soft", value: "color-mix(in srgb, var(--warm) 10%, transparent)", from: ["#f78c6c"] },
  { name: "--warm-line", value: "color-mix(in srgb, var(--warm) 30%, transparent)", from: ["#f78c6c"] },
  { name: "--info", value: "rgb(138,170,210)", from: ["rgba(138, 170, 210, 0.07)"], note: "demo 没有 info，这是它第二枚环境光的色相" },
  { name: "--info-soft", value: "color-mix(in srgb, var(--info) 10%, transparent)", from: ["rgba(138, 170, 210, 0.07)"] },
  { name: "--info-line", value: "color-mix(in srgb, var(--info) 30%, transparent)", from: ["rgba(138, 170, 210, 0.07)"] },

  // ---- lines -------------------------------------------------------------
  { name: "--border-hair", value: "var(--demo-window-border)", from: ["--window-border:"] },
  { name: "--border-soft", value: "var(--demo-border)", from: ["--border:"] },
  {
    name: "--border-strong",
    value: "color-mix(in srgb, var(--demo-border) 55%, var(--demo-foreground))",
    interim: NO_DEMO + ": demo 只画了一级 --border",
  },
  { name: "--ring", value: "var(--demo-ring)", from: ["--ring:"] },

  // ---- material ----------------------------------------------------------
  { name: "--glass-blur", value: "blur(28px) saturate(1.25)", from: ["backdrop-filter: blur(28px) saturate(1.25)"], kind: "material" },
  { name: "--glass-blur-sm", value: "blur(4px)", from: ["backdrop-filter: blur(4px)"], kind: "material" },
  { name: "--ambient-a", light: "rgba(134, 194, 185, 0.08)", dark: "rgba(134, 194, 185, 0.06)", from: ["rgba(134, 194, 185, 0.08)", "rgba(134, 194, 185, 0.06)"] },
  { name: "--ambient-b", light: "rgba(138, 170, 210, 0.07)", dark: "rgba(138, 170, 210, 0.05)", from: ["rgba(138, 170, 210, 0.07)", "rgba(138, 170, 210, 0.05)"] },
  { name: "--ambient-a-geo", value: "ellipse 80% 60% at 20% 30%", from: ["radial-gradient(ellipse 80% 60% at 20% 30%"], kind: "other" },
  { name: "--ambient-a-stop", value: "60%", from: ["radial-gradient(ellipse 80% 60% at 20% 30%"], kind: "other" },
  { name: "--ambient-b-geo", value: "ellipse 60% 50% at 80% 70%", from: ["radial-gradient(ellipse 60% 50% at 80% 70%"], kind: "other" },
  { name: "--ambient-b-stop", value: "55%", from: ["radial-gradient(ellipse 60% 50% at 80% 70%"], kind: "other" },
  { name: "--desktop-1", light: "#f0f2f5", dark: "#14161b", from: ["#f0f2f5", "#14161b"] },
  { name: "--desktop-2", light: "#e8eaed", dark: "#0f1115", from: ["#e8eaed", "#0f1115"] },
  { name: "--desktop-3", light: "#eef0f3", dark: "#121419", from: ["#eef0f3", "#121419"] },
  { name: "--desktop-geo", value: "160deg", from: ["linear-gradient(160deg, #f0f2f5"], kind: "other" },
  { name: "--nav-bg", value: "var(--demo-nav-bg)", from: ["--nav-bg:"] },
  { name: "--nav-active", value: "var(--demo-nav-active)", from: ["--nav-active:"] },
  { name: "--titlebar-bg", value: "var(--demo-titlebar-bg)", from: ["--titlebar-bg:"] },

  // ---- elevation ---------------------------------------------------------
  { name: "--shadow-card", value: "var(--demo-window-shadow)", from: ["--window-shadow:"], kind: "material" },
  { name: "--shadow-pop", value: "0 8px 24px rgba(0,0,0,.12)", from: ["box-shadow: 0 8px 24px rgba(0,0,0,.12)"], kind: "material" },
  { name: "--shadow-inset", value: "inset 0 1px 2px rgba(255,255,255,0.3)", from: ["inset 0 1px 2px rgba(255,255,255,0.3)"], kind: "material" },
  { name: "--shadow-btn", value: "0 1px 3px rgba(0,0,0,0.08)", from: ["box-shadow: 0 1px 3px rgba(0,0,0,0.08)"], kind: "material" },

  // ---- geometry ----------------------------------------------------------
  { name: "--r-xs", value: "3px", from: ["border-radius: 3px"], kind: "geom" },
  { name: "--r-sm", value: "4px", from: ["border-radius: 4px"], kind: "geom" },
  { name: "--r-md", value: "6px", from: ["border-radius: 6px"], kind: "geom" },
  { name: "--r-lg", value: "8px", from: ["border-radius: 8px"], kind: "geom" },
  { name: "--r-xl", value: "12px", from: ["border-radius: 12px"], kind: "geom" },
  { name: "--r-full", value: "999px", from: ["border-radius: 999px"], kind: "geom" },
  { name: "--panel-w", value: "720px", from: ["width: 720px"], kind: "geom" },

  // ---- type --------------------------------------------------------------
  { name: "--t-body", value: "13px", from: ["font-size: 13px"], kind: "type" },
  { name: "--t-small", value: "12px", from: ["font-size: 12px"], kind: "type" },
  { name: "--t-micro", value: "11px", from: ["font-size: 11px"], kind: "type" },
  { name: "--t-title", value: "20px", from: ["font-size: 20px"], kind: "type" },
  { name: "--t-display", value: "28px", from: ["font-size: 28px"], kind: "type" },
  { name: "--lh-body", value: "1.6", from: ["line-height:1.6"], kind: "type" },
  { name: "--lh-mono", value: "1.4", from: ["line-height: 1.4"], kind: "type" },
  { name: "--fw-body", value: "500", from: ["font-weight: 500"], kind: "type" },
  { name: "--fw-title", value: "600", from: ["font-weight: 600"], kind: "type" },
  {
    name: "--font-sans",
    value: "ui-sans-serif, system-ui, sans-serif",
    interim: NO_DEMO + "：demo 在 body 上写 class=font-sans 且从不覆盖，这就是它的解析值",
    kind: "type",
  },
  {
    name: "--font-mono",
    value: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace",
    from: ["ui-monospace, SFMono-Regular, Menlo, Consolas, monospace"],
    kind: "type",
  },

  // ---- motion ------------------------------------------------------------
  { name: "--dur-fast", value: "120ms", from: ["transition: all 120ms ease"], kind: "motion" },
  { name: "--dur-base", value: "150ms", from: ["transition: all 150ms ease"], kind: "motion" },
  { name: "--dur-slow", value: "200ms", from: ["transition: all 200ms ease"], kind: "motion" },
  { name: "--ease-out", value: "cubic-bezier(0.22, 1, 0.36, 1)", from: ["cubic-bezier(0.22, 1, 0.36, 1)"], kind: "motion" },
];

// ---------------------------------------------------------------------------
// the two rules that make this a gate and not a comment
// ---------------------------------------------------------------------------

/** An entry either cites demo or admits it is an interim alias. Never neither. */
function verify(entry) {
  if (entry.interim) {
    if (!entry.interim.startsWith(NO_DEMO)) throw new Error(`${entry.name}: interim must start with the marker`);
    return;
  }
  if (!entry.from || entry.from.length === 0) {
    throw new Error(`${entry.name}: neither a demo citation nor an interim marker`);
  }
  for (const claim of entry.from) lineOf(claim);
}

for (const entry of SEMANTIC) verify(entry);

const DUPLEX = SEMANTIC.filter((e) => e.dark !== undefined).map((e) => e.name);

function semanticDecls(theme) {
  const lines = [];
  for (const e of SEMANTIC) {
    const perTheme = theme === "dark" ? e.dark : e.light;
    if (perTheme === undefined) {
      if (theme === "dark" && e.light !== undefined) continue;
      if (e.value === undefined) continue;
    }
    const value = perTheme ?? e.value;
    const tag = e.interim ? e.interim : cite(e.from);
    const note = e.note ? ` ${e.note}` : "";
    lines.push(`  ${e.name}: ${value};  /* ${tag}${note} */`);
  }
  return lines;
}

/**
 * Tailwind reads colours out of the --color-* namespace, radii out of
 * --radius-*, fonts out of --font-*. A semantic key joins those only when its
 * own kind says it is one of those things; shadows and blurs stay raw custom
 * properties, because no single-value utility can carry them.
 */
function emitTheme() {
  const lines = ["@theme inline {"];
  for (const e of SEMANTIC) {
    if ((e.kind ?? "colour") !== "colour") continue;
    lines.push(`  --color-${e.name.slice(2)}: var(${e.name});`);
  }
  for (const e of SEMANTIC) {
    const key = e.name.slice(2);
    if (e.kind === "geom" && key.startsWith("r-") && key !== "r-full") lines.push(`  --radius-${key.slice(2)}: var(${e.name});`);
    if (key === "font-sans" || key === "font-mono") lines.push(`  --${key}: var(${e.name});`);
  }
  lines.push("}", "");
  return lines.join("\n");
}

const HEADER = `/* ============================================================================
   GENERATED FILE - DO NOT EDIT BY HAND  (second generation, owner 2026-09-25)
   ----------------------------------------------------------------------------
   Source:  design/doubao/demo/styles.css - the high-fidelity demo, which owner
            ruled the style truth source on 2026-09-25 ("全推翻了，前端重新设计
            token ... 按照 demo 改"). The first generation's source,
            design/assets/tokens.css, was moved to design/old/ by that same
            ruling and this script does not read it.
   Command: node scripts/gen-tokens.mjs          (npm run tokens)

   Layer 1 (--demo-*) is demo's own custom properties, verbatim; a bare shadcn
   HSL triplet is wrapped in hsl() and nothing else is done to it. Every other
   line is a name frontend/src speaks, written over layer 1, and carries either
   the demo line it was read off or an INTERIM marker saying demo draws no such
   thing. scripts/gen-tokens.mjs refuses to emit an entry with neither, and
   --check re-runs every citation, so a value that stops existing in the demo
   sheet fails the build instead of quietly staying in the table.

   :root is LIGHT, because demo's :root is light (demo/app.js:32 theme:'light').
   Dark is the opt-in through [data-theme="dark"], which demo also puts on <html>
   - the same element :root matches, which is why the theme-invariant lines below
   need no dark twin: they say var(--demo-x), and --demo-x is what changes.

   A colour, radius, font-size or duration literal written anywhere else under
   frontend/src is a second style truth source. That is signed by
   TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme in
   internal/panel/frontend_hygiene_test.go, not by this comment.
   ============================================================================ */

`;

const out =
  HEADER +
  ":root {\n" +
  [...lightPrims.map((p) => `  ${p.name}: ${p.value};  /* ${p.cite} */`), ...semanticDecls("light")].join("\n") +
  "\n}\n\n" +
  '[data-theme="dark"] {\n' +
  [...darkPrims.map((p) => `  ${p.name}: ${p.value};  /* ${p.cite} */`), ...semanticDecls("dark")].join("\n") +
  "\n}\n\n" +
  emitTheme();

// ---------------------------------------------------------------------------
// nothing under src/ may reference a key this table does not define
// ---------------------------------------------------------------------------

function walk(dir, acc = []) {
  for (const e of readdirSync(dir)) {
    const p = join(dir, e);
    if (statSync(p).isDirectory()) walk(p, acc);
    else if (/\.(css|ts|tsx)$/.test(e)) acc.push(p);
  }
  return acc;
}

function checkReferences() {
  const defined = new Set(primNames);
  for (const line of out.split("\n")) {
    const m = line.match(/^  (--[a-z0-9-]+):/);
    if (m) defined.add(m[1]);
  }
  const files = walk(SRC).filter((f) => f !== OUT);
  const texts = new Map();
  for (const file of files) {
    const text = readFileSync(file, "utf8");
    texts.set(file, text);
    // theme.css is allowed to define its own alias vocabulary; a name declared
    // anywhere in the tree is not dangling just because this table lacks it.
    for (const m of text.matchAll(/^\s*(--[a-z0-9-]+):/gm)) defined.add(m[1]);
  }
  const missing = new Map();
  for (const [file, text] of texts) {
    const rel = relative(repoRoot, file).replace(/\\/g, "/");
    for (const m of text.matchAll(/var\((--[a-z0-9-]+)/g)) {
      if (defined.has(m[1])) continue;
      if (!missing.has(m[1])) missing.set(m[1], rel);
    }
  }
  if (missing.size) {
    throw new Error(
      "dangling var() references, nothing under src/ defines them: " +
        [...missing].map(([k, f]) => `${k} (${f})`).join(", "),
    );
  }
  return defined.size;
}

const keyCount = checkReferences();
const interimCount = SEMANTIC.filter((e) => e.interim).length;
const summary =
  `${SEMANTIC.length} semantic keys (${interimCount} interim) over ` +
  `${lightPrims.length} light + ${darkPrims.length} dark demo primitives, ` +
  `${DUPLEX.length} of them per-theme, no dangling var()`;

if (process.argv.includes("--check")) {
  let current = "";
  try {
    current = readFileSync(OUT, "utf8");
  } catch {
    current = "";
  }
  if (current !== out) {
    console.error(`gen-tokens: ${OUT_REL} is not what ${DEMO_REL} generates - run \`npm run tokens\` and commit the result`);
    process.exit(1);
  }
  console.log(`gen-tokens: ${OUT_REL} matches ${posix.join("design", "doubao", "demo", "styles.css")} (${keyCount} keys, ${summary})`);
} else {
  writeFileSync(OUT, out, "utf8");
  console.log(`gen-tokens: wrote ${OUT_REL} (${keyCount} keys, ${summary})`);
}
