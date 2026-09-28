/* ============================================================================
   Nav + shell evidence (frontend; last rewritten 2026-09-28).
   ----------------------------------------------------------------------------
   Why this file was rewritten twice, and the lesson recorded so it does not
   happen a third time:

   The 2026-09-27 rewrite made the gate PASS by widening it to whatever the
   new sidebar happened to render. The sidebar then shipped nine rows -
   including a duplicate 设置, because the gate's own loop was iterating
   `PANEL_VIEWS` and demanded every one of them. **A gate derived from the
   thing it is checking cannot catch that thing being wrong.** It only
   restates it.

   So the shape assertions below are written by hand, from the shape a panel
   is supposed to have, and several of them are the exact bugs that shipped:

     - 设置 must appear in the sidebar EXACTLY ONCE. It shipped twice (once as
       a row in the section list, once at the bottom) and every gate was
       green, because "every view is a row" was the assertion.
     - The four rows that are not navigable views (命令 / 安全 / 隐私 / 成本)
       must NOT be sidebar rows. 命令 is the Ctrl+K palette; the other three
       are settings categories. Shipping them in the sidebar is how a
       settings list gets mistaken for navigation.
     - The sidebar must carry a search field, because the approved shape has
       one. Its absence shipped unnoticed for the same reason.

   What is still checked from the old version, unchanged and still load-bearing:

     A. Every icon name in panel-views.ts is either in PLAN.md §17.4's frozen
        list or in the owner-approved extension. Parsed out of PLAN.md at run
        time, with positive AND negative controls, so a parser that silently
        stopped finding anything fails loudly instead of passing on zero hits.
     B. The navigation layer holds no view state (Q2). Nothing survives a
        WebView restart; a snapshot that names a view outranks any local pick.
     C. A row whose icon had to be substituted records why, and the retired
        `interim` marker never comes back.

   Why design/doubao is NOT read here: `design/**` is 16 uncommitted deletions
   deep (owner moved it; nobody may commit or restore that) and CI runs on a
   clean checkout where those files exist. A gate that depends on which of us
   committed what is not a gate.

   Colours: none - this file paints nothing.
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
// A. The frozen icon list, parsed out of PLAN.md, plus the owner-approved
//    extension from the 2026-09-25 demo ruling. That ruling is branch 乙's
//    human approval for exactly the seven names below, nothing more.
// ---------------------------------------------------------------------------

const APPROVED_EXTENSION: ReadonlySet<string> = new Set([
  "message-square-text",
  "shield-check",
  "list-checks",
  "orbit",
  "settings",
  "chart-column",
  "life-buoy",
]);

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
  // name, so a whitespace-only split silently drops one name per group. Peel
  // the label; the positive controls below are what make that load-bearing
  // rather than hopeful.
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
check("the panel must have at least one screen to name", PANEL_VIEWS.length > 0, String(PANEL_VIEWS.length));

for (const v of PANEL_VIEWS) {
  check(
    `icon ${JSON.stringify(v.icon)} (screen ${v.id}) must be a frozen §17.4 name or an owner-approved extension`,
    frozen.has(v.icon) || APPROVED_EXTENSION.has(v.icon),
  );
}

// C. The substitution ledger. The 2026-09-25 ruling drew every row with the
// demo's own icon, retiring the Q1 = 甲 substitutions and their interim
// markers. The shape stays so a future substitution cannot be silent:
//   substituted - arithmetic: demoIcon !== icon. Must carry a `note`.
//   interim     - RETIRED. No row may carry one again.
for (const v of PANEL_VIEWS) {
  const substituted = v.demoIcon !== v.icon;
  if (substituted) {
    check(`screen ${v.id}: drew ${JSON.stringify(v.icon)} where the demo drew ${JSON.stringify(v.demoIcon)}, so it must record why`,
      Boolean(v.note), v.note ?? "no note");
  } else {
    check(`screen ${v.id}: ${JSON.stringify(v.demoIcon)} is what the demo drew and what is drawn, so it must not claim a substitution`,
      !v.note && !v.interim, `${v.note ?? ""} ${v.interim ?? ""}`);
  }
  check(`screen ${v.id}: the Q1 = 甲 interim marker is retired (owner 2026-09-25)`,
    !v.interim, v.interim ?? "");
}

const glyphs = PANEL_VIEWS.map((v) => v.icon);
check("screen icons must be distinct", new Set(glyphs).size === glyphs.length, glyphs.join(" "));

// ---------------------------------------------------------------------------
// B. Purity: nothing in the navigation layer remembers the screen.
// ---------------------------------------------------------------------------

for (const f of ["src/components/panel-sidebar.tsx", "src/lib/panel-views.ts"]) {
  const src = readFileSync(f, "utf8");
  check(`${f} must not hold view state (Q2)`, !/\buseState<PanelViewId|\buseReducer\b/.test(src),
    `${(src.match(/useReducer\b/g) ?? []).length} useReducer`);
}
// The sidebar may hold a search box and a collapse set - both are view
// preferences that die with the document and feed nothing back to the host,
// the same class as the settings screen's alpha knob. What it may NOT hold is
// the screen: that arrives from currentView(snapshot) in App.
const sidebarSrc = readFileSync("src/components/panel-sidebar.tsx", "utf8");
check("the sidebar must not decide which screen is showing",
  !/\bactive\s*=\s*useState\b/.test(sidebarSrc));
check("the sidebar must take its active section as a prop, not compute one",
  /active:\s*PanelSectionId/.test(sidebarSrc));

// App is the one place allowed to know which of its own screens it is showing.
// Exactly ONE state object, so the count cannot creep: the view pick plus the
// four view preferences that travel with it. Q2's real subject - nothing
// survives a WebView restart, and a snapshot naming a view outranks the local
// pick (PLAN.md:1043-1044) - is checked in currentView's own three cases and
// at the markup loop below. Undo: 撤「rail 可点」.
const appSrc = readFileSync("src/App.tsx", "utf8");
check("App must hold exactly one state object (the view pick plus its companions)",
  (appSrc.match(/= useState</g) ?? []).length === 1,
  `useState calls=${(appSrc.match(/= useState</g) ?? []).length}`);
check("App's state must still name the view pick as a field",
  /view:\s*PanelViewId \| null/.test(appSrc));
check("App must keep the ban on persisting anything (no storage API may appear)",
  !/\b(localStorage|sessionStorage|indexedDB|document\.cookie|caches|navigator\.serviceWorker)\b/.test(appSrc));
check("a view named by the snapshot must outrank the local pick",
  appSrc.includes("hostView !== undefined ? currentView(snapshot)"));

// currentView is the single door. Exercised directly because a wrong fallback
// here is invisible on screen: it just shows a different row than the host asked for.
check("an unnamed view falls back to the documented default",
  currentView({} as PanelSnapshot) === DEFAULT_VIEW, DEFAULT_VIEW);
check("a view named by the snapshot is honoured",
  currentView({ view: "cost" } as PanelSnapshot & { view?: string }) === "cost");
check("an unknown name must not be guessed at",
  currentView({ view: "settings-page" } as PanelSnapshot & { view?: string }) === DEFAULT_VIEW);

// Q-50 = 甲 (owner, 2026-09-25): the panel sends NO fifth route. A frontend
// belt for internal/panel/composer_test.go:522, which reads the same literals
// from the Go side - if someone re-adds a view request, both go red.
const libSource = readFileSync("src/lib/panel.ts", "utf8");
const routes = [...libSource.matchAll(/"(panel\.[a-z.]+)"/g)].map((m) => m[1] as string);
check("panel.ts must keep naming the routes Go already answers", routes.length >= 4, routes.join(" "));
check("no view-change route may come back without owner's ruling (Q-50 = 甲, undo: 撤 Q-50 甲)",
  !routes.some((r) => r.includes("view")), routes.join(" "));

// ---------------------------------------------------------------------------
// The shape a panel is supposed to have. Written by hand, NOT derived from
// PANEL_VIEWS - see the header for why that distinction is the whole point.
// ---------------------------------------------------------------------------

/** The four places the main column can go. */
const SECTION_ROWS = ["对话", "审批", "任务", "球状态"] as const;

/** Rows that must NOT be sidebar navigation, each with where it belongs
    instead. Shipping these in the sidebar is the 2026-09-27 bug. */
const NOT_NAVIGATION: readonly { label: string; belongs: string }[] = [
  { label: "命令", belongs: "the Ctrl+K palette (header)" },
  { label: "设置", belongs: "the settings page (sidebar bottom)" },
  { label: "安全", belongs: "a settings category" },
  { label: "隐私", belongs: "a settings category" },
  { label: "成本", belongs: "a settings category" },
];

const base: PanelSnapshot = {
  pending: [],
  results: [],
  generatedAt: "render-evidence",
} as PanelSnapshot;

const defaultHtml = renderToStaticMarkup(<App snapshot={{ ...base } as PanelSnapshot} />);
const nav = defaultHtml.match(/aria-label="面板导航"[\s\S]*?<\/nav>/)?.[0] ?? "";
// The 设置 uniqueness test spans the WHOLE sidebar, not just the nav: the
// bottom block that holds it is a sibling of the nav landmark, and the
// 2026-09-27 duplicate was one copy inside and one outside. Narrowing this to
// the nav is how it would have passed again.
const sidebar = defaultHtml.match(/<aside[\s\S]*?<\/aside>/)?.[0] ?? "";
check("the sidebar must render a navigation landmark", nav.length > 0, `${nav.length} chars`);
check("the sidebar's bottom block must be reachable for the uniqueness test", sidebar.length > 0);

// The four sections, and only the four.
for (const label of SECTION_ROWS) {
  check(`the sidebar must offer ${label} as a section`,
    nav.includes(`aria-label="${label}"`), "no such row");
}
for (const { label, belongs } of NOT_NAVIGATION) {
  check(`${label} must not be a sidebar navigation row - it belongs in ${belongs}`,
    !nav.includes(`aria-label="${label}"`), "found as a nav row");
}

// The bug that shipped, as a standing regression guard.
check("设置 must appear in the sidebar exactly once (it shipped twice on 2026-09-27)",
  (sidebar.match(/aria-label="设置"/g) ?? []).length === 1,
  `${(sidebar.match(/aria-label="设置"/g) ?? []).length} occurrences in the sidebar`);

// The search field the approved shape has, and whose absence also shipped.
check("the sidebar must carry a search field (the approved shape has one)",
  nav.includes('aria-label="搜索会话"'));
check("the search field must actually be an input, not a div that looks like one",
  /<input[^>]*aria-label="搜索会话"/.test(nav));

// Exactly one section lit, and it is the one the snapshot named.
for (const [view, label] of [["chat", "对话"], ["approval", "审批"], ["tasks", "任务"], ["ball", "球状态"]] as const) {
  const html = renderToStaticMarkup(<App snapshot={{ ...base, view } as PanelSnapshot} />);
  const rows = html.match(/aria-label="面板导航"[\s\S]*?<\/nav>/)?.[0] ?? "";
  const lit = [...rows.matchAll(/aria-current="page"[^>]*aria-label="([^"]+)"|<button[^>]*aria-label="([^"]+)"[^>]*aria-current="page"/g)]
    .map((m) => m[1] ?? m[2]);
  check(`snapshot naming ${view}: exactly one sidebar row may claim aria-current`,
    lit.length === 1, `lit=${JSON.stringify(lit)}`);
  check(`snapshot naming ${view}: the lit row must be ${label}`,
    lit[0] === label, `lit=${JSON.stringify(lit)}`);
  check(`snapshot naming ${view}: no sidebar row may be disabled now that clicking works`,
    (rows.match(/ disabled=""/g) ?? []).length === 0 &&
      (rows.match(/aria-disabled="true"/g) ?? []).length === 0);
}

// A screen Go has not fed says so. Checked for the three unfed sections the
// sidebar can reach, and for the ones only the host can name.
for (const view of ["tasks", "ball"] as const) {
  const html = renderToStaticMarkup(<App snapshot={{ ...base, view } as PanelSnapshot} />);
  check(`${view}: an unfed section must say it is unfed rather than show nothing`,
    html.includes("还没有收到"));
}
const hostOnly = renderToStaticMarkup(<App snapshot={{ ...base, view: "security" } as PanelSnapshot} />);
check("a view only the host can name must still land somewhere real, not blank",
  hostOnly.includes("还没有") || hostOnly.includes("不可逆"),
  hostOnly.slice(hostOnly.indexOf("还没有") - 40, hostOnly.indexOf("还没有") + 80));

// One right-rail control, not two. The 2026-09-27 build had a toggle in the
// header AND another in a content-area header, which is one action with two
// buttons.
const railToggles = (defaultHtml.match(/收起右侧栏|展开右侧栏/g) ?? []).length;
check("the right rail must have exactly one toggle in the chrome",
  railToggles === 2, // one control, two states
  `${railToggles} labels`);

// The context/usage line, and its honest absence.
check("the composer must sit under a context/usage readout",
  defaultHtml.includes("上下文与用量"));
check("with no usage fields the readout must say so, not print zeros",
  !/0 tok|缓存命中 0%|占用 0%/.test(defaultHtml));

// The empty conversation column must state the shape of what arrives. It used
// to draw nothing, so a panel with no data looked like one that failed.
check("the empty conversation column must name what will arrive there",
  defaultHtml.includes("流式回复") && defaultHtml.includes("工具调用") && defaultHtml.includes("思考过程"));
check("the empty conversation column must say nothing has arrived yet",
  defaultHtml.includes("还没收到任何一条"));

// The settings page: eleven categories, each saying its own thing. The first
// pass gave every unfed category the SAME sentence, so the rail could not tell
// them apart - the same failure as nine identical sidebar rows.
const settingsSrc = readFileSync("src/components/harness/settings-page.tsx", "utf8");
const categoryNames = ["通用", "外观", "模型", "语音", "记忆", "技能", "MCP", "安全", "隐私", "成本", "引导"];
for (const name of categoryNames) {
  check(`the settings page must offer the ${name} category`,
    settingsSrc.includes(`name: "${name}"`), "missing");
}
check("the settings page must not render a duplicate close affordance (返回 handles it)",
  !/aria-label="关闭设置"/.test(settingsSrc));
// Each unfed category carries its own sentence, so the rail is a status board.
const waitingBlocks = (settingsSrc.match(/waiting:\s*\n?\s*"/g) ?? []).length;
check("each settings category must carry its own waiting-for sentence",
  waitingBlocks >= 9, `${waitingBlocks} found; 外观 is the one with live controls`);

// The settings screen must not stack two 设置 headings on itself. ConfigScreen
// takes the heading from the page for exactly this reason.
check("ConfigScreen must take its heading from the caller",
  /heading\?:\s*string/.test(readFileSync("src/components/config-screen.tsx", "utf8")));

if (failures.length > 0) {
  for (const f of failures) console.error("render-nav: FAIL " + f);
  console.error(`render-nav: FAIL, ${failures.length} check(s) failed (${frozen.size} frozen + ${APPROVED_EXTENSION.size} approved icon names)`);
  process.exit(1);
}
console.error(
  `render-nav: OK, ${PANEL_VIEWS.length} screens against ${frozen.size} frozen §17.4 names + ${APPROVED_EXTENSION.size} owner-approved (2026-09-25), ` +
    `${SECTION_ROWS.length} sidebar sections, ${NOT_NAVIGATION.length} rows correctly NOT navigation, ` +
    `${categoryNames.length} settings categories, no substitutions, no view state in the nav layer`,
);
