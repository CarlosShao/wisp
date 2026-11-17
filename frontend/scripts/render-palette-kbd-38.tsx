/* ============================================================================
   Ticket 38 AC#6's behavioural tooth - the keyboard-only palette flow
   (render evidence), 2026-11-17.
   ----------------------------------------------------------------------------
   Companion to render-palette-38.tsx (which owns the wire shapes; this file
   owns ONLY the keyboard flow, so neither instrument grows a second copy of
   the other's job). Same shape that file established: feed REAL data through
   the PRODUCTION component and assert readable values come back out.

   The readables here are the two the dispatch names:
     - the FOCUS ID SEQUENCE: the highlighted-row id (palette-row-N) as the
       arrow keys walk it, driven through the component's own exported
       decision (paletteKeyDecision + paletteCursorStep) - the exact code
       path onKey runs, not a re-derived copy;
     - the aria-ACTIVEDESCENDANT value the markup announces, asserted equal
       to the row the markup itself marks aria-selected.

   Positive and negative controls, one pair per dimension:
     positive : keys walk/select in the open state (id sequence read back);
     negative : the SAME keys with the IME composition flag set return
                "none" for every one of them (AC#2's guard, still armed),
                and the empty result set announces no activedescendant.

   Focus trap / return-focus live in App and only run in a real document
   (keydown listeners, document.activeElement), so they are proven at source
   level in the App section - the same provenance render-palette-38.tsx used
   for the modal's placement, with the mechanism's own strings named.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { PaletteScreen, paletteCursorStep, paletteKeyDecision, type PaletteGroup } from "@/components/palette-screen";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   Fixtures: three rows over two groups (send row first - App puts it first
   when the query has text - then two panel rows), matching the production
   App's own group vocabulary.
   --------------------------------------------------------------------------- */

const QUERY = "";
const ROWS = 3;

const GROUPS: PaletteGroup[] = [
  { name: "发送任务", items: [{ label: "以任务发送：清理桌面临时文件", icon: "terminal", hint: "Enter", action: "send" }] },
  {
    name: "面板",
    items: [
      { label: "打开设置", icon: "settings", action: "settings" },
      { label: "任务监控", icon: "activity", action: "monitor" },
    ],
  },
];

const noop = (): void => undefined;

/* ---------------------------------------------------------------------------
   1. aria face (positive control): listbox / option / activedescendant
   --------------------------------------------------------------------------- */

const html = renderToStaticMarkup(
  <PaletteScreen query={QUERY} onQueryChange={noop} groups={GROUPS} />,
);

check("the row list must carry role=listbox", html.includes(`role="listbox"`));
check("the row list must carry its own aria-label", html.includes(`aria-label="命令列表"`));
check("the search input must carry aria-label=搜索命令", html.includes(`aria-label="搜索命令"`));
check("every row must be role=option", (html.match(/role="option"/g) ?? []).length === ROWS);
check("every row must carry its flat index as the focus id", html.includes(`id="palette-row-0"`) && html.includes(`id="palette-row-1"`) && html.includes(`id="palette-row-2"`));
check(
  "aria-activedescendant must announce the highlighted row (cursor starts at 0)",
  html.includes(`aria-activedescendant="palette-row-0"`),
);
check(
  "the announced row and the aria-selected row must be the SAME row",
  /<button aria-selected="true"[^>]*id="palette-row-0"/.test(html),
);
check(
  "the two unhighlighted rows must be aria-selected=false",
  /<button aria-selected="false"[^>]*id="palette-row-1"/.test(html) &&
    /<button aria-selected="false"[^>]*id="palette-row-2"/.test(html),
);
check("the palette markup must hold zero emoji", !EMOJI_RE.test(html));

/* Tab order, read off the DOM the component actually renders: the search box
   comes first, then the rows in ascending id order - so the browser's native
   Tab walk (which the palette deliberately never intercepts) goes 搜索框 →
   行列表 → 行内操作 in exactly that order. */
const inputAt = html.indexOf("<input");
const rowAt = [0, 1, 2].map((n) => html.indexOf(`id="palette-row-${n}"`));
check(
  "DOM order: the input precedes every row (Tab walks search box first)",
  inputAt >= 0 && rowAt.every((p) => p > inputAt),
);
check(
  "DOM order: rows appear in ascending id order (Tab follows the keyboard walk)",
  rowAt[0]! >= 0 && rowAt[0]! < rowAt[1]! && rowAt[1]! < rowAt[2]!,
);
check(
  "Tab must not be intercepted by the palette (native DOM-order walk stays)",
  paletteKeyDecision("Tab", false, ROWS) === "none",
);

/* ---------------------------------------------------------------------------
   2. Negative control (aria): an empty result set announces NOTHING
   --------------------------------------------------------------------------- */

const emptyHtml = renderToStaticMarkup(
  <PaletteScreen query="查无此命令xyz" onQueryChange={noop} groups={GROUPS} />,
);
check("an empty result set must render no listbox", !emptyHtml.includes(`role="listbox"`));
check("an empty result set must announce no activedescendant", !emptyHtml.includes("aria-activedescendant"));
check("an empty result set must render no row ids", !emptyHtml.includes("palette-row-"));

/* ---------------------------------------------------------------------------
   3. The focus id sequence (positive control): the arrow walk, through the
      component's own exported decision
   --------------------------------------------------------------------------- */

let cursor = 0;
const walk: string[] = [`palette-row-${cursor}`];
for (const key of ["ArrowDown", "ArrowDown", "ArrowUp", "ArrowUp"]) {
  const action = paletteKeyDecision(key, false, ROWS);
  check(`an open palette answers ${key} (not none)`, action !== "none");
  cursor = paletteCursorStep(cursor, action, ROWS);
  walk.push(`palette-row-${cursor}`);
}
check(
  "the focus id sequence for Down Down Up Up must read 0 1 2 1 0",
  walk.join(" ") === "palette-row-0 palette-row-1 palette-row-2 palette-row-1 palette-row-0",
  walk.join(" "),
);
check("ArrowUp from the first row must wrap to the last (down never dead-ends)", paletteCursorStep(0, "up", ROWS) === ROWS - 1);
check("ArrowDown from the last row must wrap to the first (up never dead-ends)", paletteCursorStep(ROWS - 1, "down", ROWS) === 0);

/* The component must run THIS decision on the input (not a private copy):
   the blocked flag is exactly the AC#2 composition guard. (Sources are
   normalized to LF so the check is line-ending agnostic.) */
const paletteSrc = readFileSync("src/components/palette-screen.tsx", "utf8").replace(/\r\n/g, "\n");
check(
  "the input's key handler must drive the pure decision with the composition guard",
  paletteSrc.includes("paletteKeyDecision(\n      e.key,\n      composing || e.nativeEvent.isComposing,\n      flat.length,\n    )"),
);
check(
  "the arrow path must step the cursor through paletteCursorStep",
  paletteSrc.includes("setActive(paletteCursorStep(cursor, action, flat.length))"),
);
check(
  "Enter must select the highlighted row (flat[cursor], the announced one)",
  paletteSrc.includes("onSelect?.(flat[cursor])"),
);

/* ---------------------------------------------------------------------------
   4. The negative control: composing period blocks the WHOLE vocabulary
      (ticket 38 AC#2, still armed after the AC#6 work)
   --------------------------------------------------------------------------- */

for (const key of ["ArrowDown", "ArrowUp", "Enter", "Escape", "Tab", "k"]) {
  check(
    `composing period must return none for ${key}`,
    paletteKeyDecision(key, true, ROWS) === "none",
  );
}
check(
  "an empty palette must return none even for the arrow keys (nothing to walk)",
  paletteKeyDecision("ArrowDown", false, 0) === "none" && paletteKeyDecision("Enter", false, 0) === "none",
);

/* ---------------------------------------------------------------------------
   5. ESC / Ctrl+K / trap / return-focus (App, source level - these run in
      keydown listeners and document.activeElement, which static markup
      cannot exercise)
   --------------------------------------------------------------------------- */

const appSrc = readFileSync("src/App.tsx", "utf8").replace(/\r\n/g, "\n");

// ESC: the palette itself declines it ("none" above); App's dialog handler
// owns the close.
check("the palette must not swallow Escape (App owns the close)", paletteKeyDecision("Escape", false, ROWS) === "none");
check("App must close the palette on Escape", appSrc.includes(`event.key === "Escape"`) && appSrc.includes("paletteOpen: false"));

// Ctrl+K: the toggle, with the composition guard BEFORE it.
check(
  "App must keep the Ctrl/Cmd+K toggle",
  appSrc.includes('(event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k"'),
);
check(
  "App's key handler must return while an IME composition is open (guard precedes the toggle)",
  appSrc.indexOf("if (event.isComposing) return;") >= 0 &&
    appSrc.indexOf("if (event.isComposing) return;") < appSrc.indexOf("(event.ctrlKey || event.metaKey)"),
);

// Tab trap: the dialog is queried, the walk wraps, stray focus is pulled in.
check(
  "App must trap Tab inside the palette dialog while it is open",
  appSrc.includes(`event.key === "Tab" && local.paletteOpen`) &&
    appSrc.includes(`'[role="dialog"][aria-label="命令面板"]'`) &&
    appSrc.includes("focusables[focusables.length - 1]"),
);
check(
  "the trap must wrap both ends and reclaim stray focus",
  appSrc.includes("!dialog.contains(current)") &&
    appSrc.includes("event.shiftKey && current === first") &&
    appSrc.includes("!event.shiftKey && current === last"),
);

// Return focus: captured at open (both paths), restored on close.
check(
  "App must remember the opener on both open paths and restore focus on close",
  appSrc.includes("paletteReturnFocusRef") &&
    (appSrc.match(/rememberPaletteReturnFocus\(\)/g) ?? []).length >= 2 &&
    appSrc.includes("back.focus()"),
);
check(
  "the remembered opener must be document.activeElement (an HTMLElement)",
  appSrc.includes("document.activeElement instanceof HTMLElement"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-palette38kbd: FAIL " + f);
  console.error(`render-palette38kbd: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-palette38kbd: OK, ticket 38 AC#6's keyboard-only flow: listbox/option aria " +
    "with activedescendant announcing the aria-selected row, DOM-order Tab walk " +
    "(input first, rows in id order, Tab never intercepted), arrow walk reads " +
    "palette-row-0 1 2 1 0 with wrap both ways, Enter selects the announced row, " +
    "Escape/Ctrl+K owned by App with trap + return-focus + composition guard, " +
    "and the composing period returns none for every key (negative control)",
);
