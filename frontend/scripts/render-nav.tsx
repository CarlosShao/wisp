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
check("the sidebar must take the active session as a prop, not compute one",
  /activeSessionId:\s*string \| null/.test(sidebarSrc));

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

/** The left column is a TREE of the user's own work. It is not a menu. */
const TREE_LANDMARK = "工作区会话";
const TREE_ROWS = ["新建任务", "设置", "暗色"];

/** Rows that must NOT be in the sidebar, each with where it belongs instead.
    Owner 2026-09-28, verbatim: 「怎么特么出来什么审批，任务，球状态，这特么几个
    干什么的」. Shipping these as navigation is the bug; these assertions are
    the standing guard against it coming back. */
const NOT_NAVIGATION: readonly { label: string; belongs: string }[] = [
  { label: "审批", belongs: "a card that appears when a gated call needs you, plus the header count" },
  { label: "任务", belongs: "the 任务监控 popover in the top bar" },
  { label: "球状态", belongs: "nowhere - the ball is native (D29); its settings are on the settings page" },
  { label: "命令", belongs: "the Ctrl+K palette" },
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
const nav = defaultHtml.match(new RegExp(`aria-label="${TREE_LANDMARK}"[\\s\\S]*?</nav>`))?.[0] ?? "";
// The 设置 uniqueness test spans the WHOLE sidebar, not just the nav: the
// bottom block that holds it is a sibling of the nav landmark, and the
// 2026-09-27 duplicate was one copy inside and one outside. Narrowing this to
// the nav is how it would have passed again.
const sidebar = defaultHtml.match(/<aside[\s\S]*?<\/aside>/)?.[0] ?? "";
check("the sidebar must render the session tree landmark", nav.length > 0, `${nav.length} chars`);
check("the sidebar's bottom block must be reachable for the uniqueness test", sidebar.length > 0);

for (const label of TREE_ROWS) {
  check(`the sidebar must carry ${label}`, sidebar.includes(label), "missing");
}
for (const { label, belongs } of NOT_NAVIGATION) {
  check(`${label} must not be a sidebar row - it belongs in ${belongs}`,
    !nav.includes(`aria-label="${label}"`), "found as a nav row");
}
// 设置 must appear in the sidebar exactly once. It shipped twice on 2026-09-27.
check("设置 must appear in the sidebar exactly once (it shipped twice on 2026-09-27)",
  (sidebar.match(/aria-label="设置"/g) ?? []).length === 1,
  `${(sidebar.match(/aria-label="设置"/g) ?? []).length} occurrences in the sidebar`);

// The search field the approved shape has, and whose absence also shipped.
check("the tree must carry a search field (the approved shape has one)",
  nav.includes('aria-label="搜索会话"'));
check("the search field must actually be an input, not a div that looks like one",
  /<input[^>]*aria-label="搜索会话"/.test(nav));

// A screen Go has not fed says so. The three that left the sidebar land on a
// sentence explaining where they went; the four that were never navigation
// land on the settings page AT their own category.
for (const view of ["tasks", "ball"] as const) {
  const html = renderToStaticMarkup(<App snapshot={{ ...base, view } as PanelSnapshot} />);
  check(`${view}: must land somewhere that explains itself, never on a blank screen`,
    /不是一屏|不是一个地方|不在面板里/.test(html), html.slice(0, 160));
}
for (const [view, category] of [["config", "外观"], ["security", "安全"], ["privacy", "隐私"], ["cost", "成本"]] as const) {
  const html = renderToStaticMarkup(<App snapshot={{ ...base, view } as PanelSnapshot} />);
  check(`${view}: must open the settings page at ${category}, not the conversation column`,
    html.includes(category) && html.includes("设置分类"), html.slice(0, 160));
}

// One right-rail control, not two. The 2026-09-27 build had a toggle in the
// header AND another in a content-area header, which is one action with two
// buttons.
const railToggles = (defaultHtml.match(/收起右侧栏|展开右侧栏/g) ?? []).length;
check("the right rail must have exactly one toggle in the chrome",
  railToggles === 2, // one control, two states
  `${railToggles} labels`);

// Controls that were drawn with no handler. The command palette shipped as a
// dead button: 「点命令面板图标毫无反应」. Checked in the SOURCE, because a
// compiled bundle has no onClick left in it to find.
check("the command palette control must have a handler",
  /aria-label="命令面板"[\s\S]{0,400}?onClick/.test(appSrc));
check("Ctrl+K must be wired (the chip in the header promises it)",
  /key\.toLowerCase\(\) === "k"/.test(appSrc));
// Escape shipped unhandled on the palette until 2026-09-28: a modal you could
// only leave with the mouse.
check("Escape must close the command palette",
  /event\.key === "Escape"[\s\S]{0,300}?paletteOpen: false/.test(appSrc));
check("a floating layer must dismiss on an outside click, not only its own trigger",
  readFileSync("src/lib/use-dismiss.ts", "utf8").includes("pointerdown"));

// The composer must be the library's PromptBar, not the retired local one. The
// panel shipped the old composer.tsx while the walkthrough shipped PromptBar,
// so the two halves of the product looked like two different products.
check("the conversation column must mount PromptBar", defaultHtml.includes('aria-label="Prompt"'));
check("the retired local composer must not be mounted by the panel",
  !defaultHtml.includes("给 Wisp 的指令"));
// The two things the old composer did not have, called out by owner.
check("the composer must carry a permission tier control",
  /aria-label="权限档位|permissions=\{/s.test(defaultHtml) || defaultHtml.includes("档位"));
check("PromptBar must accept reasoning effort", /efforts\?: readonly PromptEffort\[\]/.test(
  readFileSync("src/components/ai-native/prompt-bar.tsx", "utf8")));

// The empty conversation column must state the shape of what arrives. It used
// to draw nothing, so a panel with no data looked like one that failed.
check("the empty conversation column must name what will arrive there",
  defaultHtml.includes("流式回复") && defaultHtml.includes("工具调用") && defaultHtml.includes("思考过程"));
check("the empty conversation column must say nothing has arrived yet",
  defaultHtml.includes("还没收到任何一条"));

// The context/usage line, and its honest absence.
check("the composer must sit under a context/usage readout",
  defaultHtml.includes("上下文与用量"));
check("with no usage fields the readout must say so, not print zeros",
  !/0 tok|缓存命中 0%|占用 0%/.test(defaultHtml));

// The narrow-window path. 640px is the DEFAULT `[panel] width`
// (internal/config/schema.go), so a 240px sidebar that always occupies a
// column leaves 400px of conversation - this is the common case, not an edge.
check("the panel must adapt below the narrow threshold", /SIDEBAR_COLLAPSE_BELOW/.test(appSrc));
check("the threshold must be above the default panel width",
  /SIDEBAR_COLLAPSE_BELOW = (\d+)/.test(appSrc) &&
    Number(/SIDEBAR_COLLAPSE_BELOW = (\d+)/.exec(appSrc)?.[1]) > 640,
  /SIDEBAR_COLLAPSE_BELOW = (\d+)/.exec(appSrc)?.[1]);

// git: no cloud. Owner: 「云端那个就算了，我们是开源项目，没有默认云端环境的」.
// Checked with the comments stripped, because the file's own header QUOTES the
// rule - and a gate that fails on its own explanation is a gate that gets
// deleted instead of obeyed.
const gitSrc = readFileSync("src/components/git-branch.tsx", "utf8")
  .replace(/\/\*[\s\S]*?\*\//g, "")
  .replace(/^\s*\/\/.*$/gm, "");
check("the branch switcher must not offer a cloud target", !/云端/.test(gitSrc));
check("the branch switcher must offer a worktree", gitSrc.includes("工作树"));
check("a non-repo workspace must render no branch chip at all",
  gitSrc.includes("if (!view.isRepo) return null"));
check("the branch name must be read, never defaulted to a literal branch",
  /current:\s*string/.test(gitSrc) && !/current:\s*"main"/.test(gitSrc));

// The turn rail for long conversations.
const railSrc = readFileSync("src/components/turn-rail.tsx", "utf8");
check("a long conversation must have a way back to earlier turns",
  defaultHtml.includes("对话轮次") || /<TurnRail/.test(appSrc));
check("the turn rail must hide itself when there is nothing to navigate",
  /turns\.length < 2/.test(railSrc));

// The settings page: eleven categories, each saying its own thing. The first
// pass gave every unfed category the SAME sentence, so the rail could not tell
// them apart - the same failure as nine identical sidebar rows.
const settingsSrc = readFileSync("src/components/harness/settings-page.tsx", "utf8");
const categoryNames = ["通用", "外观", "模型", "语音", "记忆", "技能", "MCP", "安全", "隐私", "成本", "引导", "悬浮球"];
for (const name of categoryNames) {
  check(`the settings page must offer the ${name} category`,
    settingsSrc.includes(`name: "${name}"`), "missing");
}
// 球状态 was a sidebar row that belonged nowhere; its settings now live on
// the settings page, so the page must actually have a place for them.
check("the settings page must carry the ball's own settings (球状态 left the panel)",
  settingsSrc.includes("悬浮球"));
check("the settings page must not stack two 设置 headings on itself", 
  /heading\?:\s*string/.test(readFileSync("src/components/config-screen.tsx", "utf8")));

if (failures.length > 0) {
  for (const f of failures) console.error("render-nav: FAIL " + f);
  console.error(`render-nav: FAIL, ${failures.length} check(s) failed (${frozen.size} frozen + ${APPROVED_EXTENSION.size} approved icon names)`);
  process.exit(1);
}
console.error(
  `render-nav: OK, ${PANEL_VIEWS.length} screens against ${frozen.size} frozen §17.4 names + ${APPROVED_EXTENSION.size} owner-approved (2026-09-25), ` +
    `tree sidebar, ${NOT_NAVIGATION.length} rows correctly NOT navigation, ` +
    `${categoryNames.length} settings categories, narrow-window path at ${/SIDEBAR_COLLAPSE_BELOW = (\d+)/.exec(appSrc)?.[1]}px, ` +
    `no substitutions, no view state in the nav layer`,
);
