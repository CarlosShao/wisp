/* ============================================================================
   Nav-layer evidence (Q1 = 甲, frontend session 2026-09-25).

   Three claims this harness is responsible for, chosen because each is the kind
   a later round can silently break:

     A. Every icon the rail draws carries a name from PLAN.md §17.4's frozen
        list. Adding a name to that list is a human-approval change (that is
        branch 乙), so the check reads the list out of PLAN.md at run time
        instead of copying it - a copied list would drift into agreeing with
        whatever the rail happened to import.

     B. Nothing in the navigation layer remembers the screen (Q2). Checked as
        text: no useState in nav-rail or panel-views, and the view reaches the
        tree only through currentView(snapshot). App is the one named exception -
        a single transient pick, counted and pinned by identifier - because
        owner asked for a settings screen a click can reach (2026-09-25), and a
        snapshot naming a view still outranks it.

     C. A row whose icon had to be substituted says so. `interim` is required
        exactly when the demo's own name is not in the frozen list, and is
        forbidden otherwise, so the marker cannot be left behind after the
        substitution is fixed, nor skipped when one is made.

   Why design/doubao is NOT read here: the demo's rail is the source of the nine
   rows and of each `demoIcon`, but `design/**` is currently 16 uncommitted
   deletions deep (owner moved it, nobody may commit or restore that - brief §0),
   and CI runs on a clean checkout where those files exist while a local tree may
   not. A gate that depends on which of us committed what is not a gate. The
   values copied out of it are checked against PLAN instead, which is stable.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import { DEFAULT_VIEW, PANEL_VIEWS, currentView } from "@/lib/panel-views";
import type { PanelSnapshot } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// ---------------------------------------------------------------------------
// A. The frozen icon list, parsed out of PLAN.md.
// ---------------------------------------------------------------------------

function frozenIcons(): Set<string> {
  let plan: string;
  try {
    plan = readFileSync("../docs/PLAN.md", "utf8");
  } catch {
    throw new Error("render-nav: cannot read ../docs/PLAN.md, which the icon list is derived from");
  }
  const lines = plan.split("\n");
  const head = lines.findIndex((l) => l.includes("需要的图标清单"));
  if (head === -1) throw new Error("render-nav: PLAN.md no longer has the '需要的图标清单' heading");
  const body: string[] = [];
  let fence = 0;
  for (let i = head; i < lines.length; i++) {
    const t = lines[i]!.trim();
    if (t === "```") {
      fence += 1;
      if (fence === 2) break;
      continue;
    }
    if (fence === 1) body.push(t);
  }
  // Each group in that block starts with a Chinese label glued to its first
  // name ("状态：circle-dot circle-dashed ..."), and splitting on whitespace
  // alone silently drops one name per group - which is how this file first
  // reported that `circle-dot` and `list` were not in the frozen list when they
  // plainly are. The label is peeled here, and the POSITIVE CONTROLS below are
  // what make that fix load-bearing rather than hopeful.
  const out = new Set(
    body
      .map((l) => l.replace(/^[^：:]*[：:]/, ""))
      .join(" ")
      .split(/\s+/)
      .filter((w) => /^[a-z0-9][a-z0-9-]*$/.test(w)),
  );
  if (out.size < 40) {
    throw new Error(`render-nav: only ${out.size} icon names parsed - the fence moved, fix this check`);
  }
  // A "0 hits" reading from a parser that has never been shown a 1 is not a
  // reading. Names that must be found, and names that must not.
  for (const present of ["circle-dot", "list", "shield", "settings-2", "command", "lock", "shield-alert", "hash", "timer", "cpu", "layers", "corner-down-left"]) {
    if (!out.has(present)) {
      throw new Error(`render-nav: parser lost ${JSON.stringify(present)}, which PLAN.md §17.4 does list - the parse is broken, do not trust any other answer here`);
    }
  }
  for (const absent of ["message-square-text", "orbit", "chart-column", "coin", "wallet", "sparkles"]) {
    if (out.has(absent)) {
      throw new Error(`render-nav: parser invented ${JSON.stringify(absent)}, which §17.4 does not list`);
    }
  }
  return out;
}

const frozen = frozenIcons();
check("the rail must draw something", PANEL_VIEWS.length > 0, String(PANEL_VIEWS.length));

for (const v of PANEL_VIEWS) {
  check(`icon ${JSON.stringify(v.icon)} (row ${v.id}) must be a frozen §17.4 name`, frozen.has(v.icon));
}

// C. The substitution ledger.
//
// Two different facts, kept in two different fields, because owner approved a
// count ("那两枚") and a count needs a rule that cannot drift:
//   note    - required whenever the row draws a different name than the demo
//             did. This is arithmetic: demoIcon !== icon.
//   interim - required on top of that only where the substitute is NOT apt,
//             i.e. the frozen list has nothing chat-shaped or money-shaped at
//             all. Owner's ruling named exactly two such rows. Whether a
//             near-relative counts as apt is a judgement, so it is recorded in
//             the data and the harness checks the SHAPE of the record, not the
//             judgement: interim must be a subset of substituted, and a row
//             that drew the demo's own name may carry neither.
for (const v of PANEL_VIEWS) {
  const substituted = v.demoIcon !== v.icon;
  if (substituted) {
    check(`row ${v.id}: drew ${JSON.stringify(v.icon)} where the demo drew ${JSON.stringify(v.demoIcon)}, so it must record why`,
      Boolean(v.note), v.note ?? "no note");
  } else {
    check(`row ${v.id}: ${JSON.stringify(v.demoIcon)} is frozen and is what is drawn, so it must not claim a substitution`,
      !v.note && !v.interim, `${v.note ?? ""} ${v.interim ?? ""}`);
  }
  if (v.interim) {
    check(`row ${v.id}: interim without a substitution`, substituted, v.interim);
    check(`row ${v.id}: the interim marker must name the ruling and the date`,
      /INTERIM\(图标不贴切，Q1=甲 2026-09-25\)/.test(v.interim ?? ""), v.interim);
  }
}
const interimRows = PANEL_VIEWS.filter((v) => v.interim).map((v) => v.id);
const substitutedRows = PANEL_VIEWS.filter((v) => v.demoIcon !== v.icon).map((v) => v.id);
check("the interim set must be exactly the two rows owner's ruling covers (对话 and 成本)",
  interimRows.slice().sort().join(",") === "chat,cost",
  `interim=${interimRows.join(" ") || "none"} substituted=${substitutedRows.join(" ")}`);

// No two rows may light up with the same glyph: the rail is icons-only, so a
// duplicate is not a smaller label, it is the same picture meaning two places.
const glyphs = PANEL_VIEWS.map((v) => v.icon);
check("rail icons must be distinct", new Set(glyphs).size === glyphs.length,
  glyphs.join(" "));

// ---------------------------------------------------------------------------
// B. Purity: nothing in the navigation layer remembers the screen.
// ---------------------------------------------------------------------------

for (const f of ["src/components/nav-rail.tsx", "src/lib/panel-views.ts"]) {
  const src = readFileSync(f, "utf8");
  check(`${f} must not hold view state (Q2)`, !/\buseState\b|\buseReducer\b/.test(src),
    `${(src.match(/\buseState\b/g) ?? []).length} useState`);
}

// App is the one place allowed to know which of its own screens it is showing,
// and only as a single transient pick (owner, 2026-09-25: 设置里要能调透明度 -
// a settings screen no click reaches is not a setting). The blanket ban above
// stays for everything else; this one is narrowed to an exact count and an
// exact name so the exception cannot grow sideways. Q2's real subject - nothing
// survives a WebView restart, and a snapshot that names a view outranks the
// pick (PLAN.md:1043-1044) - is checked at the markup loop below and in
// currentView's own three cases. Undo: 撤「rail 可点」.
const appSrc = readFileSync("src/App.tsx", "utf8");
check("App may hold exactly one piece of state, and it must be the view pick",
  (appSrc.match(/= useState</g) ?? []).length === 1 &&
    /const \[picked, setPicked\] = useState<PanelViewId \| null>\(null\)/.test(appSrc),
  `useState calls=${(appSrc.match(/= useState</g) ?? []).length}`);
check("App must keep the ban on persisting anything (no storage API may appear)",
  !/\b(localStorage|sessionStorage|indexedDB|document\.cookie|caches|navigator\.serviceWorker)\b/.test(appSrc));
check("a view named by the snapshot must outrank the local pick",
  appSrc.includes("hostView !== undefined ? currentView(snapshot)"));
const skeleton = readFileSync("src/components/panel-skeleton.tsx", "utf8");
// Matched on the declaration, not the word: this file's own header comment
// names PANEL_TABS to explain what it replaced, and a plain substring test read
// that explanation as the tabs still being there.
check("the old 乙 header tabs must be gone, not left alongside the rail",
  !/export const PANEL_TABS/.test(skeleton) && !/type PanelTab\b/.test(skeleton));
check("the rail must be the only thing offering a screen change",
  (skeleton.match(/NavRail/g) ?? []).length >= 2);

// currentView is the single door. Exercised directly because a wrong fallback
// here is invisible on screen: it just shows a different row than the host asked for.
check("an unnamed view falls back to the documented default",
  currentView({} as PanelSnapshot) === DEFAULT_VIEW, DEFAULT_VIEW);
check("a view named by the snapshot is honoured",
  currentView({ view: "cost" } as PanelSnapshot & { view?: string }) === "cost");
check("an unknown name must not be guessed at",
  currentView({ view: "settings-page" } as PanelSnapshot & { view?: string }) === DEFAULT_VIEW);

// ---------------------------------------------------------------------------
// Markup, once per row.
// ---------------------------------------------------------------------------

const base: PanelSnapshot = {
  pending: [],
  results: [],
  // composer intentionally absent: the fallback path is part of what is checked
  generatedAt: "render-evidence",
} as PanelSnapshot;

for (const v of PANEL_VIEWS) {
  const html = renderToStaticMarkup(<App snapshot={{ ...base, view: v.id } as PanelSnapshot} />);
  const rows = html.match(/aria-label="面板视图"[\s\S]*?<\/nav>/)?.[0] ?? "";
  check(`row ${v.id}: the rail must render all ${PANEL_VIEWS.length} labels`,
    PANEL_VIEWS.every((r) => rows.includes(`aria-label="${r.label}"`)),
    `${(rows.match(/<button/g) ?? []).length} buttons`);
  check(`row ${v.id}: exactly one row may claim aria-current`,
    (rows.match(/aria-current="page"/g) ?? []).length === 1,
    String((rows.match(/aria-current="page"/g) ?? []).length));
  check(`row ${v.id}: the lit row must be ${v.label}`,
    rows.includes(`aria-current="page"`) &&
      new RegExp(`aria-current="page"[^>]*aria-label="${v.label}"|aria-label="${v.label}"[^>]*aria-current="page"`).test(rows));
  check(`row ${v.id}: its hidden name must still be in the markup`,
    rows.includes(`>${v.label}</span>`));
  if (!v.fed && !v.selfFed) {
    check(`row ${v.id}: an unfed screen must say it is unfed rather than show nothing`,
      html.includes("还没有接到数据"));
  }
  if (v.selfFed) {
    check(`row ${v.id}: a self-fed screen must name the one thing it can actually do`,
      html.includes("即时") && html.includes("未接线"));
    // The config keys are stored in two halves so they do not read as host routes
    // (see src/components/config-screen.tsx). That split is only safe while the
    // JOINED text still appears, so it is asserted from the markup side - after
    // dropping the `<!-- -->` React puts between adjacent text nodes, which is
    // what makes the joined string invisible otherwise.
    const joined = html.replace(/<!--\s*-->/g, "");
    check(`row ${v.id}: the config keys must still render joined, not as halves`,
      joined.includes("panel.opacity") && joined.includes("ball.opacity_idle"),
      joined.slice(joined.indexOf("font-mono"), joined.indexOf("font-mono") + 120));
  }
}

// The hidden labels are hidden by CSS. If that ever becomes display:none the
// names leave the accessibility tree and the rail becomes icons-only for a
// screen reader too - which is the part of 甲 nobody would notice missing.
const css = readFileSync("src/styles/theme.css", "utf8");
const labelRule = css.match(/^\.nav-rail-label\s*\{([^}]*)\}/m)?.[1] ?? "";
check(".nav-rail-label must exist", labelRule.length > 0);
check("the rail labels must be hidden by opacity/visibility, never display:none",
  !/display:\s*none/.test(labelRule) && /visibility:\s*hidden/.test(labelRule), labelRule.trim());
check("hover AND focus must reveal the label (keyboard users get the name too)",
  /\.nav-rail-item:hover > \.nav-rail-label/.test(css) &&
    /\.nav-rail-item:focus(-visible)? > \.nav-rail-label/.test(css));

// Q-50 = 甲 (owner, 2026-09-25): the panel sends NO fifth route. This is a
// frontend-side belt for internal/panel/composer_test.go:522, which reads the
// same literals from the Go side - if someone re-adds a view request, both
// should go red, and the one in this tree fails before a push does.
const libSource = readFileSync("src/lib/panel.ts", "utf8");
const routes = [...libSource.matchAll(/"(panel\.[a-z.]+)"/g)].map((m) => m[1] as string);
check("panel.ts must keep naming the routes Go already answers", routes.length >= 4, routes.join(" "));
check("no view-change route may come back without owner's ruling (Q-50 = 甲, undo: 撤 Q-50 甲)",
  !routes.some((r) => r.includes("view")), routes.join(" "));

// The rail is a live control now (owner, 2026-09-25: 设置里要能调透明度 - a
// settings screen nobody can reach is not a setting). It must be a real button,
// and it must hand the choice UP rather than keep it: the nav layer holding no
// view state is what keeps Q2 true while the click works. Undo the click half:
// 撤「rail 可点」, and the two checks below flip back with it.
const railSource = readFileSync("src/components/nav-rail.tsx", "utf8");
check("the rail must carry a click handler", railSource.includes("onClick={() =>"));
check("the nav layer must hold no view state of its own", !/\buseState\b/.test(railSource));
const firstRail = renderToStaticMarkup(<App snapshot={{ ...base, view: "chat" } as PanelSnapshot} />);
const railMarkup = firstRail.match(/aria-label="面板视图"[\s\S]*?<\/nav>/)?.[0] ?? "";
check("no rail row may claim disabled now that clicking works",
  (railMarkup.match(/ disabled=""/g) ?? []).length === 0 &&
    (railMarkup.match(/aria-disabled="true"/g) ?? []).length === 0,
  `aria-disabled=${(railMarkup.match(/aria-disabled="true"/g) ?? []).length} disabled=${(railMarkup.match(/ disabled=""/g) ?? []).length}`);
check("exactly one row must be lit, and lit must be the row the snapshot named",
  (railMarkup.match(/aria-current="page"/g) ?? []).length === 1 &&
    railMarkup.includes("aria-label=\"对话\""));

if (failures.length > 0) {
  for (const f of failures) console.error("render-nav: FAIL " + f);
  console.error(`render-nav: FAIL, ${failures.length} check(s) failed (${frozen.size} frozen icon names)`);
  process.exit(1);
}
console.error(
  `render-nav: OK, ${PANEL_VIEWS.length} rail rows against ${frozen.size} frozen §17.4 names, ` +
    `interim=${interimRows.length} (${interimRows.join(",")}), no view state in the nav layer`,
);
