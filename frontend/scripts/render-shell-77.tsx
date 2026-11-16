/* ============================================================================
   AC#11 render + contrast instrument (ticket 77, batch 2, 2026-11-16).
   ----------------------------------------------------------------------------
   What it does
   ------------
   Renders the PRODUCTION root (src/App.tsx - the panel shell: header,
   sidebar, conversation, composer) plus one sample screen (the L2 strong
   confirmation card, driven by the real risk-decision fixtures) under BOTH
   themes, then measures every text/background pair the rendered tree
   actually paints and judges each pair against WCAG AA (normal text 4.5:1;
   there is no large text on this shell, so the 3:1 tier has no rows).

   The pair list is not hand-picked: the instrument walks the rendered
   markup's tag stack, reads every text/bg colour utility actually
   painted, resolves each to its token value from
   src/styles/tokens.generated.css (light = :root, dark = [data-theme=dark],
   aliases followed), and requires each declared pair below to be OBSERVED in
   the markup - a refactor that silently removes a pair goes red.

   Verdict rules (no adjective verdicts, numbers only)
   ---------------------------------------------------
   - every declared pair is printed with its measured ratio and PASS/FAIL;
   - a failing pair that is NOT registered below is a hard failure (exit 1)
     naming the pair and its ratio - never a blanket "overall pass";
   - KNOWN_BREAKS registers the pairs that are below AA with the NAMED hop
     their fix needs (a contract-surface change outside this wave's write
     face). A registered pair is still measured and printed every run; the
     moment it passes, the run FAILS until the register row is deleted - the
     register cannot hide a fix or accumulate silently.

   Disabled controls: WCAG 1.4.3 exempts inactive UI components. The
   disabled pair is still composited (CSS opacity over its background) and
   printed, but contributes no failure.

   Why scripts/ and not src/: it reads files with node:fs (same reasoning as
   scripts/render-l2.tsx - AC#5's statelessness scanner must not need an
   exception). It is linted by `npm run lint`; everything it renders (src/**)
   stays under `npm run typecheck`.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import type { ApprovalCardView, ComposerState, PanelSnapshot } from "@/lib/panel";

/* --------------------------------------------------------------------------
   inputs: the two REAL decision fixtures (regenerated from
   `wisp panel-assets -l2 ...`, whose -l2 branch is
   panel.NewApprovalCardView(risk.NewRiskAssessor(), ...) - see
   docs/evidence/s1/77-scaffold-batch2-ac3-ac11.md for the chain). The second
   fixture is the wire's documented absent-reason shape
   (internal/panel/approval.go: ReasonKnown=false), i.e. "not fed" data: the
   card must answer it with the explicit fail-closed wording, never a blank.
   -------------------------------------------------------------------------- */

function loadView(path: string): ApprovalCardView {
  const view: unknown = JSON.parse(readFileSync(path, "utf8"));
  if (
    typeof view !== "object" ||
    view === null ||
    typeof (view as ApprovalCardView).correlationId !== "string"
  ) {
    throw new Error(`render-shell-77: ${path} does not hold an ApprovalCardView JSON`);
  }
  return view as ApprovalCardView;
}

const REAL = loadView("fixtures/l2-card-fs-delete.json");
const UNJUDGED = loadView("fixtures/l2-card-unjudged-reason.json");

/** The composer section the cold start actually renders (App's own
    EMPTY_COMPOSER shape): the workspace row - and with it the genuinely
    disabled 切换 button - must be part of the measured shell. */
const COLD_COMPOSER: ComposerState = {
  mode: { current: "unknown", names: [], l2ConfirmNames: [] },
  workspace: {
    set: false,
    spelling: "",
    canonical: "",
    reparse: false,
    rewritten: false,
    reason: "尚未收到原生侧的状态快照",
  },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

function snapshotOf(view: ApprovalCardView): PanelSnapshot {
  return { pending: [view], results: [], composer: COLD_COMPOSER, generatedAt: "contrast-evidence" };
}

/* --------------------------------------------------------------------------
   theme blocks: parse the generated token table (the one legal colour-literal
   home) into per-theme value maps, and resolve var() aliases.
   -------------------------------------------------------------------------- */

const THEME_CSS = readFileSync("src/styles/tokens.generated.css", "utf8");

function parseBlock(header: string): Map<string, string> {
  const start = THEME_CSS.indexOf(header);
  if (start < 0) throw new Error(`render-shell-77: no ${header} block in tokens.generated.css`);
  const open = THEME_CSS.indexOf("{", start);
  const close = THEME_CSS.indexOf("}", open);
  const body = THEME_CSS.slice(open + 1, close).replace(/\/\*[\s\S]*?\*\//g, "");
  const out = new Map<string, string>();
  for (const line of body.split(";")) {
    const m = line.match(/^\s*(--[A-Za-z0-9-]+)\s*:\s*(.+?)\s*$/);
    if (m) out.set(m[1]!, m[2]!);
  }
  return out;
}

const LIGHT = parseBlock(":root");
const DARK = parseBlock('[data-theme="dark"]');
if (LIGHT.size === 0 || DARK.size === 0) {
  throw new Error("render-shell-77: empty theme block(s) - parse broke, every ratio below would be vacuous");
}

for (const [theme, block] of [["light", LIGHT], ["dark", DARK]] as const) {
  const alpha = block.get("--panel-alpha");
  if (alpha !== "100%") {
    throw new Error(
      `render-shell-77: --panel-alpha is ${JSON.stringify(alpha)} in ${theme}; ` +
        "the translucency knob moved, so surface contrast needs a real composite, not this shortcut",
    );
  }
}

/** Resolve one token name to an sRGB colour with its alpha, following var()
    aliases. Dark-theme tints are 8-digit hex (#f68f3c24): the alpha rides
    along and the caller composites over the backdrop the element sits on
    (the card/header surface). bg-page/bg-surface utilities ride
    color-mix(... var(--panel-alpha) ...) with the factory 100% (asserted
    above), which is an identity, so the base vars are the painted colours. */
function resolveRGBA(name: string, block: Map<string, string>, depth = 0): { rgb: [number, number, number]; a: number } {
  if (depth > 8) throw new Error(`render-shell-77: alias cycle at ${name}`);
  const value = block.get(name);
  if (value === undefined) throw new Error(`render-shell-77: token ${name} is not declared`);
  const varRef = value.match(/^var\((--[A-Za-z0-9-]+)\)$/);
  if (varRef) return resolveRGBA(varRef[1]!, block, depth + 1);
  const hex = value.match(/^#([0-9A-Fa-f]{6})([0-9A-Fa-f]{2})?$/);
  if (!hex) throw new Error(`render-shell-77: ${name} = ${JSON.stringify(value)} is not a plain #rrggbb[aa]`);
  const n = parseInt(hex[1]!, 16);
  return {
    rgb: [(n >> 16) & 0xff, (n >> 8) & 0xff, n & 0xff],
    a: hex[2] ? parseInt(hex[2], 16) / 255 : 1,
  };
}

function resolveColor(name: string, block: Map<string, string>): [number, number, number] {
  const c = resolveRGBA(name, block);
  if (c.a < 1) throw new Error(`render-shell-77: ${name} is translucent; composite it over a backdrop first`);
  return c.rgb;
}

/** A translucent token paints as itself composited over the surface behind
    it (the card/header plane the pill and badges sit on). */
function resolveOverSurface(name: string, block: Map<string, string>): [number, number, number] {
  const c = resolveRGBA(name, block);
  return c.a < 1 ? composite(c.rgb, resolveColor("--surface", block), c.a) : c.rgb;
}

/* --------------------------------------------------------------------------
   WCAG 2 relative luminance / contrast ratio, plus the sRGB compositing a
   CSS opacity blending needs for the disabled pair.
   -------------------------------------------------------------------------- */

function linearize(channel: number): number {
  const c = channel / 255;
  return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
}

function luminance(rgb: [number, number, number]): number {
  const [r, g, b] = rgb.map(linearize) as [number, number, number];
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrast(fg: [number, number, number], bg: [number, number, number]): number {
  const l1 = luminance(fg);
  const l2 = luminance(bg);
  const [hi, lo] = l1 >= l2 ? [l1, l2] : [l2, l1];
  return (hi + 0.05) / (lo + 0.05);
}

function composite(fg: [number, number, number], bg: [number, number, number], alpha: number): [number, number, number] {
  return [0, 1, 2].map((i) => Math.round(fg[i]! * alpha + bg[i]! * (1 - alpha))) as [number, number, number];
}

function hexOf(rgb: [number, number, number]): string {
  return "#" + rgb.map((c) => c.toString(16).padStart(2, "0")).join("");
}

/* --------------------------------------------------------------------------
   markup scan: walk the rendered tag stack, collect (text colour, nearest
   background) pairs and disabled controls.
   -------------------------------------------------------------------------- */

/** Tailwind colour names this theme actually declares (theme.css @theme
    inline maps text-<name>/bg-<name> to --color-<name> -> var(--<name>)). */
const COLOUR_NAMES = new Set([
  "page", "canvas", "surface", "inset", "hover", "hover-2",
  "ink", "ink-2", "ink-3", "line", "line-strong", "field",
  "stripe", "stripe-bg",
  "accent", "accent-ink", "accent-tint",
  "green", "green-tint", "orange", "orange-tint",
  "red", "red-tint",
  "tooltip-bg", "tooltip-fg", "tooltip-muted", "tooltip-border",
  "warn", "warn-soft", "warn-line", "danger", "stop", "success", "success-soft",
]);

const VOID_TAGS = new Set([
  "br", "hr", "img", "input", "meta", "link", "area", "base", "col", "embed", "source", "track", "wbr",
]);

interface Observed {
  pairs: Set<string>;
  disabled: { fg: string; alpha: number; bg: string }[];
}

function colourUtility(token: string, prefix: "text" | "bg"): { name: string; alpha: number } | null {
  if (token.includes(":")) return null; // hover:/disabled:/... variants are not painted at rest
  const m = token.match(new RegExp(`^${prefix}-([a-z0-9-]+?)(?:/(\\d+))?$`));
  if (!m) return null;
  const name = m[1]!;
  if (prefix === "text" ? name !== "white" && !COLOUR_NAMES.has(name) : !COLOUR_NAMES.has(name)) return null;
  return { name, alpha: m[2] ? parseInt(m[2]!, 10) / 100 : 1 };
}

function scanMarkup(html: string): Observed {
  const pairs = new Set<string>();
  const disabled: Observed["disabled"] = [];
  const stack: { tag: string; classes: string[] }[] = [];
  const tagRe = /<(\/?)([a-zA-Z][a-zA-Z0-9-]*)((?:"[^"]*"|'[^']*'|[^>"'])*)>/g;
  for (const m of html.matchAll(tagRe)) {
    const closing = m[1] === "/";
    const tag = m[2]!.toLowerCase();
    let attrs = m[3] ?? "";
    const selfClosing = attrs.trimEnd().endsWith("/");
    if (selfClosing) attrs = attrs.trimEnd().slice(0, -1);
    const classAttr = attrs.match(/class="([^"]*)"/);
    const classes = classAttr ? classAttr[1]!.split(/\s+/) : [];

    if (closing) {
      for (let i = stack.length - 1; i >= 0; i--) {
        if (stack[i]!.tag === tag) {
          stack.length = i;
          break;
        }
      }
      continue;
    }

    const ownBg = classes.map((c) => colourUtility(c, "bg")).find((b) => b !== null) ?? null;
    const ownText = classes.map((c) => colourUtility(c, "text")).find((t) => t !== null) ?? null;
    if (ownText) {
      let bgName: string | null = ownBg?.name ?? null;
      if (!bgName) {
        for (let i = stack.length - 1; i >= 0 && bgName === null; i--) {
          const up = stack[i]!.classes.map((c) => colourUtility(c, "bg")).find((b) => b !== null);
          if (up) bgName = up.name;
        }
      }
      if (bgName) pairs.add(`${ownText.name}/${bgName}`);
    }
    if (tag === "button" && /\bdisabled(?:=|\s|$)/.test(attrs) && ownText) {
      let bgName: string | null = ownBg?.name ?? null;
      for (let i = stack.length - 1; i >= 0 && bgName === null; i--) {
        const up = stack[i]!.classes.map((c) => colourUtility(c, "bg")).find((b) => b !== null);
        if (up) bgName = up.name;
      }
      const op = classes.map((c) => c.match(/^disabled:opacity-(\d+)$/)).find((x) => x !== null);
      if (bgName) {
        disabled.push({ fg: ownText.name, alpha: op ? parseInt(op[1]!, 10) / 100 : 0.5, bg: bgName });
      }
    }

    if (!selfClosing && !VOID_TAGS.has(tag)) stack.push({ tag, classes });
  }
  return { pairs, disabled };
}

/* --------------------------------------------------------------------------
   the declared pair list. "where" names the producer element; each pair is
   REQUIRED to be observed in the listed trees, so the measurement rides the
   render instead of a colour table nobody can prove is on screen.
   -------------------------------------------------------------------------- */

interface PairSpec {
  id: string;
  category: string;
  fg: string; // "white" or a token name like "ink-2"; "@<alpha>" suffix composites
  bg: string;
  where: string;
  trees: ("real" | "unjudged")[];
}

const PAIRS: PairSpec[] = [
  { id: "body-on-surface", category: "正文", fg: "ink", bg: "surface", where: "卡片正文/按钮文字（拒绝卡 bg-surface 卡面）", trees: ["real", "unjudged"] },
  { id: "body-on-page", category: "正文", fg: "ink", bg: "page", where: "会话占位标题（一缕在这儿）", trees: ["real", "unjudged"] },
  { id: "body-on-inset", category: "正文", fg: "ink", bg: "inset", where: "完整参数块（bg-inset 里的命令文本）", trees: ["real", "unjudged"] },
  { id: "secondary-on-surface", category: "次要", fg: "ink-2", bg: "surface", where: "理由行/调用链/常驻提示/顶栏图标钮", trees: ["real", "unjudged"] },
  { id: "label-on-surface", category: "标签", fg: "ink-3", bg: "surface", where: "小节标题（完整参数/为什么需要确认）与占位 chips", trees: ["real", "unjudged"] },
  { id: "label-on-page", category: "标签", fg: "ink-3", bg: "page", where: "占位副文案与 composer 工作区行", trees: ["real", "unjudged"] },
  { id: "kbd-on-field", category: "标签", fg: "ink-2", bg: "field", where: "顶栏 Ctrl K 键帽（.kbd 自定义类，非 utility）", trees: ["real", "unjudged"] },
  { id: "button-white-on-red", category: "按钮", fg: "white", bg: "red", where: "拒绝按钮（bg-red text-white）", trees: ["real", "unjudged"] },
  { id: "button-ink-on-surface", category: "按钮", fg: "ink", bg: "surface", where: "查看完整参数（outline 按钮）", trees: ["real", "unjudged"] },
  { id: "pill-warn-on-soft", category: "标签", fg: "warn", bg: "warn-soft", where: "顶栏 待审批 计数胶囊", trees: ["real", "unjudged"] },
  { id: "unjudged-warn-on-surface", category: "标签(未判)", fg: "warn", bg: "surface", where: "信息不足 前缀（reasonKnown=false 渲染）", trees: ["unjudged"] },
  { id: "badge-red-on-surface", category: "标签", fg: "red", bg: "surface", where: "L2 强确认 徽标文字", trees: ["real", "unjudged"] },
  { id: "disabled-on-shell", category: "禁用态", fg: "ink-3@0.5", bg: "*", where: "composer 切换 按钮 disabled:opacity-50（WCAG 1.4.3 豁免，仍实测）", trees: ["real", "unjudged"] },
];

/** Pairs measured below AA whose fix needs a contract-surface change outside
    this wave's write face (frontend/** + docs/evidence/s1/77-*). Adding a
    colour token to the generator makes the token uncovered by the C21 four-way
    reconciliation (internal/panel/tokens_fourway_test.go check 2) until the
    hand-maintained table and the pinned list grow a waiver row - both of those
    files are other territories, so this wave reports the break and its hop
    instead of hand-waving it. A registered row is still measured every run:
    the run FAILS the moment the pair passes, so the register cannot outlive
    its reason. */
const KNOWN_BREAKS: Record<string, string> = {
  "dark/button-white-on-red":
    "断跳: 深色 --red 同时充当徽标文字（text-red on surface 实测 4.68 过 AA），压深按钮红必然击穿文字红——单 token 无解，修复需契约面新增 token（如 --red-strong）；新增彩色 token 会让四方对账 check 2 未覆盖，须同步 docs/evidence/s1/c21-native-tokens.md 第四代未覆盖 + internal/panel/tokens_fourway_test.go 钉名表，两者均在本波写面之外（亮色半边已由 --red 亮色值上调修复，4.71:1）",
  "light/pill-warn-on-soft":
    "断跳: 橙字在橙底浅 tint 上无法到 4.5:1，需文本专用 token（如 --warn-ink，库已有 --accent-ink 先例）；新增彩色 token 触发四方对账 check 2 未覆盖，须同步 c21-native-tokens.md 与钉名表，本波写面之外",
  "light/unjudged-warn-on-surface":
    "断跳: 同 pill-warn-on-soft，需 --warn-ink 文本 token；契约面改动本波不可写",
};

/* --------------------------------------------------------------------------
   run: four trees = {light, dark} x {real, unjudged}.
   -------------------------------------------------------------------------- */

function renderTree(theme: "light" | "dark", view: ApprovalCardView): string {
  const app = <App snapshot={snapshotOf(view)} />;
  return theme === "dark"
    ? renderToStaticMarkup(<div data-theme="dark">{app}</div>)
    : renderToStaticMarkup(<div>{app}</div>);
}

const trees = {
  real: { light: renderTree("light", REAL), dark: renderTree("dark", REAL) },
  unjudged: { light: renderTree("light", UNJUDGED), dark: renderTree("dark", UNJUDGED) },
} as const;

// The fail-closed half belongs to AC#3; asserted here too because this
// instrument already holds the unjudged tree.
for (const theme of ["light", "dark"] as const) {
  const html = trees.unjudged[theme];
  if (!html.includes("信息不足")) {
    console.error(`render-shell-77: ${theme} unjudged tree lost the fail-closed wording`);
    process.exit(1);
  }
  if (html.includes(REAL.reason)) {
    console.error(`render-shell-77: ${theme} unjudged tree leaks the real reason it was not fed`);
    process.exit(1);
  }
}

// The kbd chip is a custom class (colour in theme.css), not a utility pair:
// require the class to be painted, then measure its declared colours.
for (const kind of ["real", "unjudged"] as const) {
  for (const theme of ["light", "dark"] as const) {
    if (!trees[kind][theme].includes('class="kbd')) {
      console.error(`render-shell-77: ${theme} ${kind} tree does not paint the kbd chip`);
      process.exit(1);
    }
  }
}

const observed = {
  real: { light: scanMarkup(trees.real.light), dark: scanMarkup(trees.real.dark) },
  unjudged: { light: scanMarkup(trees.unjudged.light), dark: scanMarkup(trees.unjudged.dark) },
} as const;

if (process.env.SHELL77_DUMP) {
  for (const kind of ["real", "unjudged"] as const) {
    for (const theme of ["light", "dark"] as const) {
      console.error(`render-shell-77: observed pairs ${theme}/${kind}: ${[...observed[kind][theme].pairs].sort().join(", ")}`);
      console.error(`render-shell-77: observed disabled ${theme}/${kind}: ${JSON.stringify(observed[kind][theme].disabled)}`);
    }
  }
}

function fgColor(spec: string, block: Map<string, string>): [number, number, number] {
  const [name, alphaStr] = spec.split("@");
  const base = name === "white" ? ([255, 255, 255] as [number, number, number]) : resolveColor(`--${name}`, block);
  return alphaStr ? composite(base, [255, 255, 255], parseFloat(alphaStr!)) : base;
  // note: a plain white composite is only correct over white; the disabled
  // pair composites over its real bg below - this branch exists for a future
  // alpha fg declared with a fixed bg, and is not exercised by PAIRS today.
}

const lines: string[] = [];
const failures: string[] = [];
let measured = 0;

for (const theme of ["light", "dark"] as const) {
  const block = theme === "light" ? LIGHT : DARK;
  for (const spec of PAIRS) {
    for (const kind of spec.trees) {
      const obs = observed[kind][theme];
      const key = `${theme}/${spec.id}`;
      const registryKey = `${theme}/${spec.id}`;

      // observation: the render must actually paint the pair
      if (spec.id === "kbd-on-field") {
        // custom class pair, colours from theme.css, presence asserted above
      } else if (spec.bg === "*") {
        if (obs.disabled.length === 0) {
          failures.push(`${key}: the render paints no disabled control for ${spec.where}`);
          continue;
        }
      } else if (!obs.pairs.has(`${spec.fg}/${spec.bg}`)) {
        failures.push(`${key}: the rendered tree does not paint ${spec.fg} on ${spec.bg} (${spec.where})`);
        continue;
      }

      // measurement
      let fg: [number, number, number];
      let bg: [number, number, number];
      let alphaNote = "";
      if (spec.bg === "*") {
        const d = obs.disabled[0]!;
        const base = resolveColor(`--${d.fg}`, block);
        bg = resolveOverSurface(`--${d.bg}`, block);
        fg = composite(base, bg, d.alpha);
        alphaNote = ` (${d.fg} @${d.alpha} over ${d.bg})`;
      } else {
        fg = fgColor(spec.fg, block);
        bg = resolveOverSurface(`--${spec.bg}`, block);
      }
      const ratio = contrast(fg, bg);
      measured++;
      const pass = ratio >= 4.5;
      const registered = KNOWN_BREAKS[registryKey];
      const row =
        `[${theme}] ${spec.category} ${spec.id}: ${hexOf(fg)} on ${hexOf(bg)}${alphaNote} = ${ratio.toFixed(2)}:1 vs 4.5:1 ` +
        (pass ? "PASS" : "FAIL(<AA)") + (registered ? " [registered break]" : "");
      lines.push(row);
      if (!pass && !registered && spec.id !== "disabled-on-shell") {
        failures.push(`${key}: ${ratio.toFixed(2)}:1 is below AA 4.5:1 (${spec.where})`);
      }
      if (pass && registered) {
        failures.push(`${key}: registered break now passes (${ratio.toFixed(2)}:1) - delete its KNOWN_BREAKS row`);
      }
    }
  }
}

// the disabled pair is exempt (WCAG 1.4.3 inactive components) but must still
// be observed exactly once per tree
for (const theme of ["light", "dark"] as const) {
  for (const kind of ["real", "unjudged"] as const) {
    const d = observed[kind][theme].disabled;
    if (d.length > 1) {
      failures.push(`${theme}/${kind}: ${d.length} disabled controls matched; the scanner over-collects`);
    }
  }
}

process.stdout.write(`render-shell-77: contrast pairs measured = ${measured}\n`);
for (const line of lines) process.stdout.write(line + "\n");

if (failures.length > 0) {
  for (const f of failures) console.error("render-shell-77: FAIL " + f);
  console.error(`render-shell-77: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
const registered = Object.keys(KNOWN_BREAKS).length;
console.error(
  `render-shell-77: OK, ${measured} pair measurements across light+dark x real+unjudged; ` +
    `${measured - registered} pairs pass AA 4.5:1, ${registered} pairs registered below AA with a named fix hop ` +
    `(KNOWN_BREAKS; measured and printed every run, auto-fails when fixed)`,
);
