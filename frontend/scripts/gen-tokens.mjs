// Wisp design-token generator - THIRD GENERATION (owner-delegated, 2026-09-25).
//
// Single source of truth: design/doubao/demo/styles.css, plus
// design/doubao/demo/index.html for the two things the stylesheet never
// declares (the font stacks and the window geometry come from the demo's
// Tailwind config block).
//
// Owner's standing ruling (2026-09-25, twice): the high-fidelity demo is the
// style truth source ("全推翻了，前端重新设计 token ... 按照 demo 改"), and the
// frontend owner was told to rebuild the token layer from it directly. The
// second generation's own semantic renaming layer (--bg-base / --fg-primary /
// --accent ...) is retired: every extra name between the demo and a component
// was one more place to drift. This generation speaks the demo's OWN names.
//
//   layer 1   demo's custom properties, verbatim - same names, same values,
//             bare shadcn HSL triplets left as triplets exactly as demo writes
//             them (the point of use spells hsl(...), also exactly as demo does)
//   layer 2   values demo paints WITHOUT declaring as a variable - the success
//             green, the warn amber, the code-block syntax colours, the fog
//             spheres, the ball gradient, the layered window shadow. Each entry
//             cites the exact demo text it was read off, so "1:1 照 demo" stays
//             an instrument and not a slogan.
//   layer 3   the shared scales (type / radius / motion / geometry), each cited
//             to the demo line that uses the value.
//
// `cite()` looks the claimed source text up in the demo and stamps the line
// number; `verify()` refuses an entry with neither a citation nor an INTERIM
// marker. Hand-editing src/styles/tokens.generated.css is still forbidden.
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
const DEMO_HTML_REL = "design/doubao/demo/index.html";
const DEMO_HTML = join(repoRoot, "design", "doubao", "demo", "index.html");
const OUT_REL = "src/styles/tokens.generated.css";
const OUT = join(frontendRoot, "src", "styles", "tokens.generated.css");
const SRC = join(frontendRoot, "src");

const demo = readFileSync(DEMO, "utf8");
const demoHtml = readFileSync(DEMO_HTML, "utf8");

/** 1-based line where `text` appears; throws when it does not. */
function lineIn(haystack, rel, text) {
  const idx = haystack.indexOf(text);
  if (idx < 0) throw new Error(`citation gone stale, not in ${rel}: ${JSON.stringify(text)}`);
  return haystack.slice(0, idx).split("\n").length;
}

/** "styles.css:120" - the default citation surface is the demo stylesheet. */
function lineOf(text) {
  return lineIn(demo, DEMO_REL, text);
}

/** "index.html:56" - for the claims that live in the demo's HTML. */
function citeHtml(text) {
  return `index.html:${lineIn(demoHtml, DEMO_HTML_REL, text)}`;
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
function block(source, sel) {
  const body = stripComments(source);
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
 * demo's `:root` is LIGHT and `.dark` overrides it; demo/app.js:32 ships
 * `theme: 'light'`, owner confirmed the default stays light on 2026-09-25.
 * The emitted file follows demo: `:root` is light, dark is the opt-in through
 * [data-theme="dark"] (demo toggles a `.dark` class on <html>; our index.html
 * carries the same switch as a data-attribute - same element, either selector).
 */
const lightBlock = block(demo, ":root");
const darkBlock = block(demo, ".dark");

// ---------------------------------------------------------------------------
// layer 2 - values demo paints inline, named so components can reach them
// ---------------------------------------------------------------------------

/**
 * Entry shape:
 *   name    the custom property
 *   value   theme-invariant expression over layer-1 names (declared in BOTH
 *           blocks: a var() inside a custom property resolves where it is
 *           declared, so a dark redefinition of --primary would not reach a
 *           --x-soft declared only in :root - re-declaring it in the dark
 *           block is what makes the mix pick up the dark hue)
 *   light / dark   per-theme literal, for the handful demo itself changes
 *   from    the demo text this rests on - each claim must be findable
 *   interim what it stands in for, when demo draws no such thing
 */
const EXTENDED = [
  // ---- the frozen caret row's colour name ---------------------------------
  // PLAN.md:3476 freezes "1px × 14px 竖线光标（--accent）" - the brand teal.
  // demo spells the same colour --primary (same three numbers in :root and
  // .dark), so --accent aliases it and the pale tint demo calls --accent is
  // emitted one block above as --accent-soft.
  { name: "--accent", value: "hsl(var(--primary))", from: ["--primary: 173 34% 55%"], note: "C21 冻结行的 --accent＝品牌青，即 demo 的 --primary" },

  // ---- the frost material itself ------------------------------------------
  { name: "--glass-blur", value: "blur(28px) saturate(1.25)", from: ["backdrop-filter: blur(28px) saturate(1.25)"] },
  { name: "--glass-blur-sm", value: "blur(4px)", from: ["backdrop-filter: blur(4px)"] },

  // ---- success green: demo draws it at three stops, never names it --------
  { name: "--success", light: "hsl(142 71% 45%)", dark: "#7DB896", from: ["color: hsl(142 71% 45%)", ".dark .badge-success { color: #7DB896; }"] },
  { name: "--success-ink", light: "hsl(142 71% 38%)", dark: "#7DB896", from: ["color: hsl(142 71% 38%)"], note: "文字/图标用的深一档绿" },
  { name: "--success-soft", value: "color-mix(in srgb, var(--success) 10%, transparent)", from: ["hsl(142 71% 45% / 0.1)"] },
  { name: "--success-line", value: "color-mix(in srgb, var(--success) 30%, transparent)", from: ["hsl(142 71% 45% / 0.3)"] },

  // ---- warn amber: the L2 等待色, red is for errors (PLAN.md:3480) --------
  { name: "--warn", light: "#b8860b", dark: "#D9B26A", from: ["color: #b8860b", ".dark .badge-warn { color: #D9B26A; }"] },
  { name: "--warn-soft", value: "rgba(217, 178, 106, 0.1)", from: ["rgba(217, 178, 106, 0.1)"] },
  { name: "--warn-line", value: "rgba(217, 178, 106, 0.3)", from: ["rgba(217, 178, 106, 0.3)"] },
  { name: "--warn-chip", value: "rgba(217, 178, 106, 0.14)", from: ["rgba(217, 178, 106, 0.14)"], note: "bu-chip-blocked 的底" },

  // ---- code-block syntax colours ------------------------------------------
  { name: "--tok-key", value: "#c792ea", from: [".code-block .tok-key { color: #c792ea; }"] },
  { name: "--tok-str", value: "#c3e88d", from: [".code-block .tok-str { color: #c3e88d; }"] },
  { name: "--tok-num", value: "#f78c6c", from: [".code-block .tok-num { color: #f78c6c; }"], note: "demo 唯一的暖色，暖语义沿用它" },

  // ---- desktop backdrop behind the frost (harness/desktop preview) --------
  { name: "--ambient-a", light: "rgba(134, 194, 185, 0.08)", dark: "rgba(134, 194, 185, 0.06)", from: ["rgba(134, 194, 185, 0.08)", "rgba(134, 194, 185, 0.06)"] },
  { name: "--ambient-b", light: "rgba(138, 170, 210, 0.07)", dark: "rgba(138, 170, 210, 0.05)", from: ["rgba(138, 170, 210, 0.07)", "rgba(138, 170, 210, 0.05)"] },
  { name: "--ambient-a-geo", value: "ellipse 80% 60% at 20% 30%", from: ["radial-gradient(ellipse 80% 60% at 20% 30%"] },
  { name: "--ambient-a-stop", value: "60%", from: ["radial-gradient(ellipse 80% 60% at 20% 30%"] },
  { name: "--ambient-b-geo", value: "ellipse 60% 50% at 80% 70%", from: ["radial-gradient(ellipse 60% 50% at 80% 70%"] },
  { name: "--ambient-b-stop", value: "55%", from: ["radial-gradient(ellipse 60% 50% at 80% 70%"] },
  { name: "--desktop-1", light: "#f0f2f5", dark: "#14161b", from: ["#f0f2f5", "#14161b"] },
  { name: "--desktop-2", light: "#e8eaed", dark: "#0f1115", from: ["#e8eaed", "#0f1115"] },
  { name: "--desktop-3", light: "#eef0f3", dark: "#121419", from: ["#eef0f3", "#121419"] },
  { name: "--desktop-geo", value: "160deg", from: ["linear-gradient(160deg, #f0f2f5"] },

  // ---- fog spheres (demo's ambience layer, at most one in production) -----
  { name: "--fog-1", value: "radial-gradient(circle, rgba(134,194,185,0.5), transparent 70%)", from: ["radial-gradient(circle, rgba(134,194,185,0.5), transparent 70%)"] },
  { name: "--fog-2", value: "radial-gradient(circle, rgba(90,106,180,0.35), transparent 70%)", from: ["radial-gradient(circle, rgba(90,106,180,0.35), transparent 70%)"] },

  // ---- the ball glyph (title-bar dot, preloader, L2 card 图示) ------------
  { name: "--ball-grad", value: "linear-gradient(135deg, rgba(134, 194, 185, 0.9), rgba(108, 168, 160, 0.9))", from: ["linear-gradient(135deg, rgba(134, 194, 185, 0.9), rgba(108, 168, 160, 0.9))"] },
  { name: "--ball-dot-grad", value: "linear-gradient(135deg, #86C2B9, #6CA8A0)", from: ["linear-gradient(135deg, #86C2B9, #6CA8A0)"] },
  { name: "--ball-dot-shadow", value: "0 1px 3px rgba(134,194,185,0.4)", from: ["box-shadow: 0 1px 3px rgba(134,194,185,0.4)"] },
  { name: "--ball-glow", value: "0 4px 16px rgba(134, 194, 185, 0.35), inset 0 1px 2px rgba(255,255,255,0.3)", from: ["0 4px 16px rgba(134, 194, 185, 0.35)"] },
  { name: "--ball-glow-hover", value: "0 6px 20px rgba(134, 194, 185, 0.45), inset 0 1px 2px rgba(255,255,255,0.4)", from: ["0 6px 20px rgba(134, 194, 185, 0.45)"] },
  { name: "--preloader-grad", value: "radial-gradient(circle at 32% 28%, rgba(255,255,255,.75), rgba(134,194,185,.9) 46%, rgba(108,168,160,.95))", from: ["radial-gradient(circle at 32% 28%, rgba(255,255,255,.75)"] },
  { name: "--preloader-glow", value: "0 0 32px rgba(134,194,185,.4), inset 0 2px 6px rgba(255,255,255,.5)", from: ["0 0 32px rgba(134,194,185,.4)"] },

  // ---- elevation beyond --window-shadow -----------------------------------
  { name: "--shadow-window-v4", value: "0 1px 2px rgba(0,0,0,0.04), 0 8px 24px rgba(0,0,0,0.06), 0 24px 64px rgba(0,0,0,0.10)", from: ["0 1px 2px rgba(0,0,0,0.04)"], note: "v4 的三层窗口投影，压过 :root 里那枚单层 --window-shadow" },
  { name: "--shadow-menu", value: "0 8px 24px rgba(0,0,0,.12)", from: ["box-shadow: 0 8px 24px rgba(0,0,0,.12)"] },
  { name: "--shadow-toast", value: "0 4px 16px rgba(0,0,0,0.15)", from: ["box-shadow: 0 4px 16px rgba(0,0,0,0.15)"] },
  { name: "--shadow-knob", value: "0 1px 3px rgba(0,0,0,0.2)", from: ["box-shadow: 0 1px 3px rgba(0,0,0,0.2)"] },
  { name: "--shadow-subtab", value: "0 1px 2px rgba(0, 0, 0, 0.06)", from: ["box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06)"] },
  { name: "--shadow-btn", value: "0 1px 3px rgba(0,0,0,0.08)", from: ["box-shadow: 0 1px 3px rgba(0,0,0,0.08)"], note: "sub-tab.active 的微投影，也作库组件的按钮投影词" },
  { name: "--shadow-ctx", value: "0 2px 8px rgba(0,0,0,0.04)", from: ["box-shadow: 0 2px 8px rgba(0,0,0,0.04)"] },

  // ---- the opacity knob (the settings screen's one live control) ----------
  {
    name: "--panel-alpha",
    value: "72%",
    from: ["rgba(255, 255, 255, 0.72)"],
    note: "demo 窗口 rgba alpha 拎出来的旋钮；设置屏唯一能动的就是它",
  },
  {
    name: "--panel-surface",
    value: "rgb(255 255 255 / var(--panel-alpha))",
    light: "rgb(255 255 255 / var(--panel-alpha))",
    dark: "rgb(22 24 29 / var(--panel-alpha))",
    from: ["rgba(255, 255, 255, 0.72)", "rgba(22, 24, 29, 0.72)"],
    note: "窗口磨砂面：同 demo 三个通道，alpha 换成旋钮",
  },
];

// ---------------------------------------------------------------------------
// layer 3 - the shared scales, each cited to the demo line using the value
// ---------------------------------------------------------------------------

const SCALES = [
  // radius: demo paints 4 / 6 / 8 / 10 / 12 / 999 (kbd 4 · title-bar-btn 6 ·
  // btn 8 · nav-item 10 · card 12 · pill 999). Mapped onto Tailwind's
  // --radius-* names one stop apart from the defaults on purpose.
  { name: "--r-xs", value: "4px", from: ["border-radius: 4px"] },
  { name: "--r-sm", value: "6px", from: ["border-radius: 6px"] },
  { name: "--r-md", value: "8px", from: ["border-radius: 8px"] },
  { name: "--r-lg", value: "10px", from: ["border-radius: 10px"] },
  { name: "--r-xl", value: "12px", from: ["border-radius: 12px"] },
  { name: "--r-full", value: "999px", from: ["border-radius: 999px"] },
  { name: "--r-nav", value: "60px", from: ["width: 60px"], note: "左竖条宽" },

  // type: demo's steps are 11 / 12 / 13 / 20 / 28. Owner ruled 2026-09-25 that
  // the demo wins over PLAN.md §17.3.2's 14px floor where the two collide.
  { name: "--t-micro", value: "11px", from: ["font-size: 11px"] },
  { name: "--t-small", value: "12px", from: ["font-size: 12px"] },
  { name: "--t-body", value: "13px", from: ["font-size: 13px"] },
  { name: "--t-title", value: "20px", from: ["font-size: 20px"] },
  { name: "--t-display", value: "28px", from: ["font-size: 28px"] },
  { name: "--lh-body", value: "1.6", from: ["line-height:1.6"] },
  { name: "--lh-mono", value: "1.4", from: ["line-height: 1.4"] },
  { name: "--lh-code", value: "1.7", from: ["line-height: 1.7"] },
  { name: "--fw-medium", value: "500", from: ["font-weight: 500"] },
  { name: "--fw-semibold", value: "600", from: ["font-weight: 600"] },

  // motion
  { name: "--dur-fast", value: "120ms", from: ["transition: all 120ms ease"] },
  { name: "--dur-base", value: "150ms", from: ["transition: all 150ms ease"] },
  { name: "--dur-slow", value: "200ms", from: ["transition: all 200ms ease"] },
  { name: "--dur-screen", value: "220ms", from: ["animation: rbScreenIn 220ms"], note: "屏切换" },
  { name: "--dur-glide", value: "280ms", from: ["transition: transform 280ms"], note: "nav-glide 滑翔条" },
  { name: "--ease-out", value: "cubic-bezier(0.22, 1, 0.36, 1)", from: ["cubic-bezier(0.22, 1, 0.36, 1)"] },
  { name: "--ease-stream", value: "cubic-bezier(0.22,0.61,0.25,1)", from: ["cubic-bezier(0.22,0.61,0.25,1)"], note: "stream-word 逐词" },
  { name: "--ease-fadeup", value: "cubic-bezier(0.23, 1, 0.32, 1)", from: ["cubic-bezier(0.23, 1, 0.32, 1)"], note: "animated-list-item" },

  // geometry of the window itself
  { name: "--panel-w", value: "720px", from: ["width: 720px"] },
  { name: "--panel-h", value: "780px", from: ["height: 780px"] },
  { name: "--titlebar-h", value: "40px", from: ["height: 40px;\n  min-height: 40px"] },

  // fonts: the stylesheet never declares a stack; the demo's Tailwind config does.
  { name: "--font-sans-stack", value: "Inter, 'PingFang SC', 'Microsoft YaHei UI', system-ui, sans-serif", htmlFrom: ["sans: ['Inter', 'PingFang SC', 'Microsoft YaHei UI', 'system-ui', 'sans-serif']"] },
  { name: "--font-mono-stack", value: "'JetBrains Mono', 'Cascadia Code', 'Consolas', monospace", htmlFrom: ["mono: ['JetBrains Mono', 'Cascadia Code', 'Consolas', 'monospace']"] },
];

// ---------------------------------------------------------------------------
// the two rules that make this a gate and not a comment
// ---------------------------------------------------------------------------

/** An entry either cites the demo or admits it is an interim alias. Never neither. */
function verify(entry) {
  if (entry.interim) {
    if (!entry.interim.startsWith(NO_DEMO)) throw new Error(`${entry.name}: interim must start with the marker`);
    return;
  }
  const claims = entry.from ?? [];
  const htmlClaims = entry.htmlFrom ?? [];
  if (claims.length === 0 && htmlClaims.length === 0) {
    throw new Error(`${entry.name}: neither a demo citation nor an interim marker`);
  }
  for (const claim of claims) lineOf(claim);
  for (const claim of htmlClaims) citeHtml(claim);
}

for (const entry of [...EXTENDED, ...SCALES]) verify(entry);

function entryTag(entry) {
  if (entry.interim) return entry.interim;
  const parts = [];
  if (entry.from?.length) parts.push(cite(entry.from));
  if (entry.htmlFrom?.length) parts.push(citeHtml(entry.htmlFrom));
  return parts.join(" + ");
}

const DUPLEX = [...EXTENDED, ...SCALES].filter((e) => e.dark !== undefined).map((e) => e.name);

function extendedDecls(theme) {
  const lines = [];
  for (const e of [...EXTENDED, ...SCALES]) {
    const perTheme = theme === "dark" ? e.dark : e.light;
    const value = perTheme ?? e.value;
    if (value === undefined) continue;
    const tag = entryTag(e);
    const note = e.note ? ` ${e.note}` : "";
    lines.push(`  ${e.name}: ${value};  /* ${tag}${note} */`);
  }
  return lines;
}

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

function checkReferences(out) {
  const defined = new Set(lightBlock.map((d) => d.name));
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
    for (const m of text.matchAll(/var\((--[a-z0-9-]+)(\s*,[^)]*)?\)/g)) {
      // A var() with a fallback (var(--mx, 50%)) survives never being declared;
      // the spotlight's --mx/--my are written at runtime from JS, never in CSS.
      if (m[2] !== undefined) continue;
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

const HEADER = `/* ============================================================================
   GENERATED FILE - DO NOT EDIT BY HAND  (third generation, owner-delegated)
   ----------------------------------------------------------------------------
   Source:  design/doubao/demo/styles.css - the high-fidelity demo, which owner
            ruled the style truth source on 2026-09-25 ("全推翻了，前端重新设计
            token ... 按照 demo 改") and confirmed again when handing the
            frontend to a new owner-agent the same day ("以原型demo为主").
            Font stacks and window geometry cite design/doubao/demo/index.html.

   Shape (third generation): the table speaks the demo's OWN names.
     layer 1  demo's custom properties, verbatim - same names, same values,
              HSL triplets left as triplets (the point of use spells hsl(),
              also exactly as demo does)
     layer 2  values demo paints without declaring (success green, warn amber,
              syntax colours, fog, ball, layered shadow, the alpha knob) - each
              cited to the demo text it was read off
     layer 3  shared scales (type / radius / motion / geometry), each cited

   The second generation's renaming layer (--bg-base / --fg-primary) is
   retired with owner's delegation; scripts/gen-tokens.mjs refuses an entry
   with neither a demo citation nor an INTERIM marker, and --check re-runs
   every citation, so a value that stops existing in the demo fails the build
   instead of quietly staying in the table.

   :root is LIGHT (demo/app.js:32 theme:'light'; owner confirmed the default).
   Dark is the opt-in through [data-theme="dark"] (demo toggles .dark on the
   same element). Theme-invariant entries are declared in BOTH blocks on
   purpose: a var() inside a custom property resolves where it is declared, so
   the dark block's re-declaration is what lets a color-mix() see dark hues.

   A colour, radius, font-size or duration literal written anywhere else under
   frontend/src is a second style truth source. That is signed by
   TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme in
   internal/panel/frontend_hygiene_test.go, not by this comment.
   ============================================================================ */

`;

const primitiveLines = (decls) =>
  decls.map((d) => {
    // PLAN.md:3476 (frozen SSE row) names the caret colour token `--accent` and
    // means the brand teal. Demo spells that colour `--primary` and uses
    // `--accent` for the pale tint. The frozen row wins, so demo's pale tint is
    // emitted as --accent-soft and --accent is reserved (extended layer below).
    const name = d.name === "--accent" ? "--accent-soft" : d.name;
    const note = d.name === "--accent" ? " demo 叫 --accent；--accent 留给 C21 冻结行的品牌青" : "";
    return `  ${name}: ${d.value};  /* styles.css:${lineOf(`${d.name}: ${d.value}`)}${note} */`;
  });

const out =
  HEADER +
  ":root {\n" +
  [...primitiveLines(lightBlock), ...extendedDecls("light")].join("\n") +
  "\n}\n\n" +
  '[data-theme="dark"] {\n' +
  [...primitiveLines(darkBlock), ...extendedDecls("dark")].join("\n") +
  "\n}\n";

const keyCount = checkReferences(out);
const interimCount = [...EXTENDED, ...SCALES].filter((e) => e.interim).length;
const summary =
  `${EXTENDED.length} extended + ${SCALES.length} scale keys (${interimCount} interim) over ` +
  `${lightBlock.length} light + ${darkBlock.length} dark demo primitives, ` +
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
