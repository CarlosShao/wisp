// Wisp design-token generator - FOURTH GENERATION (owner pivot, 2026-09-25).
//
// Owner's ruling, in person, after rejecting the demo-faithful first batch:
// the demo is dead as a visual source ("不按照demo来了…不管这个demo了"). The
// design system is now MINIMAL, light/dark only, and built ON the beautiful-ui
// component library's own theme: every base component speaks bg-surface /
// text-ink-2 / border-line, and the palette below is that library's theme
// adopted verbatim so a vendored component renders natively.
//
// Provenance: the light and dark palettes, the shadows, the radii and the
// easings are adopted from beautiful-ui (TurboKach/ai-native-react-components,
// MIT, Copyright (c) 2026 Turbo) app/globals.css as of upstream commit
// 05dab2d2. Adoption, not transcription: this table is the truth source from
// here on, and each WISP ADDITION below carries its own reason.
//
//   layer 1   the library's theme, verbatim values, light + dark
//   layer 2   Wisp additions - names the frozen PLAN rows and our own files
//             require, each with a note saying what it aliases and why
//
// Hand-editing src/styles/tokens.generated.css is still forbidden.
//
// Usage:  node scripts/gen-tokens.mjs [--check]
//   --check  regenerate in memory, exit 1 if the committed file differs, or if
//            any var(--token) used under src/ resolves to nothing.
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join, relative } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = join(here, "..");
const repoRoot = join(frontendRoot, "..");
const OUT_REL = "src/styles/tokens.generated.css";
const OUT = join(frontendRoot, "src", "styles", "tokens.generated.css");
const SRC = join(frontendRoot, "src");

/**
 * Entry shape: { name, light, dark, note? }. A missing `dark` means the value
 * is theme-invariant and is emitted into both blocks anyway - a var() inside a
 * custom property resolves where it is declared, so re-declaring in the dark
 * block is what lets color-mix()/aliases see dark hues.
 */
const TOKENS = [
  // ---- layer 1: the beautiful-ui theme (Turbo, MIT, commit 05dab2d2) ------
  { name: "--page", light: "#fafafb", dark: "#17181a" },
  { name: "--canvas", light: "#f1f2f3", dark: "#1c1d1f" },
  { name: "--surface", light: "#ffffff", dark: "#232427" },
  { name: "--inset", light: "#f7f8f9", dark: "#1f2022" },
  { name: "--hover", light: "#f4f5f6", dark: "#2a2b2e" },
  { name: "--hover-2", light: "#e7e9eb", dark: "#313236" },
  { name: "--ink", light: "#1f2124", dark: "#f2f3f4" },
  { name: "--ink-2", light: "#62656b", dark: "#a5a8ad" },
  { name: "--ink-3", light: "#9a9da3", dark: "#6c6f75" },
  { name: "--line", light: "#ecedef", dark: "#2e3033" },
  { name: "--line-strong", light: "#e0e2e5", dark: "#3a3c40" },
  { name: "--field", light: "#f2f2f3", dark: "#2b2c2f" },
  { name: "--stripe", light: "#49494913", dark: "#ffffff0e" },
  { name: "--stripe-bg", light: "#f5f5f5", dark: "#1b1c1e" },
  { name: "--accent", light: "#0285ff", dark: "#3d9aff", note: "PLAN.md:3476 冻结行的 --accent 即此色（SSE 光标）" },
  { name: "--accent-ink", light: "#0170dd", dark: "#7ec0ff" },
  { name: "--accent-tint", light: "#e9f3ff", dark: "#3d9aff29" },
  { name: "--green", light: "#189a4d", dark: "#3dbb72" },
  { name: "--green-tint", light: "#e8f5ed", dark: "#3dbb7224" },
  { name: "--orange", light: "#ef720c", dark: "#f68f3c" },
  { name: "--orange-tint", light: "#fdf1e5", dark: "#f68f3c24" },
  { name: "--red", light: "#e3474c", dark: "#ee5c61" },
  { name: "--red-tint", light: "#fcecec", dark: "#ee5c6124" },
  { name: "--tooltip-bg", light: "#25272b", dark: "#111214" },
  { name: "--tooltip-fg", light: "#f6f7f8", dark: "#f2f3f4" },
  { name: "--tooltip-muted", light: "#a5a8ad", dark: "#a5a8ad" },
  { name: "--tooltip-border", light: "#3a3c40", dark: "#2e3033" },
  { name: "--shadow-hairline", light: "0 0 0 1px var(--line)", dark: "0 0 0 1px var(--line)" },
  { name: "--shadow-btn", light: "0 0 0 1px var(--line-strong), 0 1px 2px #1018280d", dark: "0 0 0 1px var(--line-strong), 0 1px 2px #0000004d" },
  { name: "--shadow-card", light: "0 0 0 1px var(--line), 0 1px 2px #1018280a, 0 2px 6px #10182808", dark: "0 0 0 1px var(--line), 0 1px 2px #0003, 0 2px 6px #0003" },
  { name: "--shadow-raised", light: "0 0 0 1px var(--line), 0 2px 10px #0000000b", dark: "0 0 0 1px var(--line), 0 2px 10px #00000038" },
  { name: "--shadow-overlay", light: "0 0 0 1px var(--line), 0 8px 28px #0001", dark: "0 0 0 1px var(--line-strong), 0 8px 28px #00000057" },
  { name: "--shadow-inset-field", light: "inset 0 1px 2px #0000001f", dark: "inset 0 1px 2px #0006" },

  // ---- layer 2: Wisp additions -------------------------------------------
  // PLAN.md 的十四态表（:3475-3488，硬约束）用的语义名是 --warn / --danger /
  // --success。库的词汇里它们是 orange / red / green——别名接过去，冻结行的
  // 语义不变：琥珀是等待、红是错误（:3480 原文）。
  { name: "--warn", light: "var(--orange)", dark: "var(--orange)", note: "冻结行语义名→库的 orange" },
  { name: "--warn-soft", light: "var(--orange-tint)", dark: "var(--orange-tint)", note: "同上" },
  { name: "--warn-line", light: "color-mix(in srgb, var(--orange) 30%, transparent)", dark: "color-mix(in srgb, var(--orange) 30%, transparent)", note: "冻结行要 1px 描边，30% 是描边惯用浓度" },
  { name: "--danger", light: "var(--red)", dark: "var(--red)", note: "render-l2 仪器经 --stop→--danger 钉色" },
  { name: "--stop", light: "var(--danger)", dark: "var(--danger)", note: "render-l2 钉死的拼法" },
  { name: "--success", light: "var(--green)", dark: "var(--green)", note: "冻结行语义名→库的 green" },
  { name: "--success-soft", light: "var(--green-tint)", dark: "var(--green-tint)" },
  // 设置屏唯一的活控制（owner 2026-09-25：「设置里要能调透明度」）。极简纯面
  // 设计下它调的是窗口本体的不透明度，默认全实。
  { name: "--panel-alpha", light: "100%", dark: "100%", note: "设置屏滑杆的旋钮，作用于 .wisp-shell" },
  // 字体栈：库用 Inter；中文落到 PingFang/雅黑，Windows 首选 Segoe UI。
  { name: "--font-sans-stack", light: "Inter, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei UI', system-ui, sans-serif", dark: "Inter, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei UI', system-ui, sans-serif", note: "库主题用 Inter（next/font）；我们无网打包，同族回退" },
  { name: "--font-mono-stack", light: "ui-monospace, 'SF Mono', 'Cascadia Code', Consolas, monospace", dark: "ui-monospace, 'SF Mono', 'Cascadia Code', Consolas, monospace", note: "库的 --font-mono 回退族" },
];

function verify() {
  for (const t of TOKENS) {
    if (!t.light) throw new Error(`${t.name}: missing light value`);
    if (!t.dark) throw new Error(`${t.name}: missing dark value`);
  }
}

function decls(theme) {
  return TOKENS.map((t) => {
    const value = theme === "dark" ? t.dark : t.light;
    const note = t.note ? `  /* ${t.note} */` : "";
    return `  ${t.name}: ${value};${note}`;
  });
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
  const defined = new Set(TOKENS.map((t) => t.name));
  const files = walk(SRC).filter((f) => f !== OUT);
  const texts = new Map();
  for (const file of files) {
    const text = readFileSync(file, "utf8");
    texts.set(file, text);
    // theme.css may alias vocabulary; a name declared anywhere in the tree is
    // not dangling just because this table lacks it.
    for (const m of text.matchAll(/^\s*(--[a-z0-9-]+):/gm)) defined.add(m[1]);
  }
  const missing = new Map();
  for (const [file, text] of texts) {
    const rel = relative(repoRoot, file).replace(/\\/g, "/");
    for (const m of text.matchAll(/var\((--[a-z0-9-]+)(\s*,[^)]*)?\)/g)) {
      // var() with a fallback (var(--mx, 50%)) survives never being declared;
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
   GENERATED FILE - DO NOT EDIT BY HAND  (fourth generation, owner pivot)
   ----------------------------------------------------------------------------
   Design system: MINIMAL, light/dark only (owner 2026-09-25: "必须是极简风的，
   就两种风格，暗色/浅色…不搞花里胡哨的"). The demo-faithful generations are
   retired with that ruling; the high-fidelity demo is no longer a style
   source.

   Provenance: layer 1 is the beautiful-ui component library's own theme,
   adopted verbatim from TurboKach/ai-native-react-components (MIT, Copyright
   (c) 2026 Turbo) app/globals.css at upstream commit 05dab2d2 - adopting the
   library's theme is what makes every vendored beautiful-ui component render
   natively, with no per-component remapping. Layer 2 carries the names the
   frozen PLAN.md rows and our own files require; each carries a note.

   Command: node scripts/gen-tokens.mjs          (npm run tokens)
   Check:   this file must equal what the script emits, and every var() under
   src/ must resolve to a key declared here or aliased in theme.css. The
   colour-literal rule itself is signed by
   TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme in
   internal/panel/frontend_hygiene_test.go - this file is the one place a
   colour literal may live.
   ============================================================================ */

`;

verify();
const out =
  HEADER +
  ":root {\n" + decls("light").join("\n") + "\n}\n\n" +
  '[data-theme="dark"] {\n' + decls("dark").join("\n") + "\n}\n";

const keyCount = checkReferences(out);

if (process.argv.includes("--check")) {
  let current = "";
  try {
    current = readFileSync(OUT, "utf8");
  } catch {
    current = "";
  }
  if (current !== out) {
    console.error(`gen-tokens: ${OUT_REL} is not what the generator emits - run \`npm run tokens\` and commit the result`);
    process.exit(1);
  }
  console.log(`gen-tokens: ${OUT_REL} matches the generator (${keyCount} keys, ${TOKENS.length} tokens, light+dark)`);
} else {
  writeFileSync(OUT, out, "utf8");
  console.log(`gen-tokens: wrote ${OUT_REL} (${keyCount} keys, ${TOKENS.length} tokens, light+dark)`);
}
