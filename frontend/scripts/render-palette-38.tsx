/* ============================================================================
   Ticket 38's behavioural tooth (render evidence), 2026-11-16.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-l2 / render-approval37 /
   render-history-36 already established for this tree: every check below
   feeds REAL-SHAPED data through the PRODUCTION component and asserts the
   SAME VALUE comes back out in the markup - "喂真数据⇒显示同一份值".

   The fixtures are the wire's own documented shapes:
     - PaletteRosterRow verbatim off internal/panel/commands_roster_213.go:45-51
       (the JSON keys CommandsSection declares in lib/panel.ts, reconciled
       both ways by TestComposerContractTypesMatchFrontend);
     - PaletteNote / PaletteGroup / PaletteItem off the component's own props
       (App builds them from the pushed snapshot at render time).

   Where a dimension has NO wire field, the check asserts its ABSENCE - a
   renderer that invented a value for an undeclared key would be the exact
   假值 the dispatch forbids. The named Go-only dimensions here:
     - the task card's started / cost-to-date / current step (ticket 38 AC#4)
       do not exist on TaskRowView (internal/panel/subagent_roster_197.go:111-142);
       the live tasks surface stays the monitor popover's verbatim rows and
       this file asserts the palette invents none of the three;
     - the palette's send row dispatches panel.message.send, the ONE route
       the composer uses (lib/panel.ts sendMessage) - asserted at source
       level; the page adds no route name of its own (the Go-side roster
       instrument internal/panel/composer_test.go:502 and bridge_test.go's
       page-gap ledger own that rule).

   Placement is checked at SOURCE level where static markup cannot click:
   the palette modal only mounts while local.paletteOpen is true, so its
   feed and its handlers are proven by reading App's source - the same
   provenance render-approval37 used for the strip's every-view placement.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import {
  PaletteScreen,
  type PaletteGroup,
  type PaletteNote,
  type PaletteRosterRow,
} from "@/components/palette-screen";
import type { PanelSnapshot } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
// FE0F rides as its own alternation branch (oxlint: no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** Same escaping React applies to a text node - so a match here is a match on
    screen, not a match on an unescaped substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* ---------------------------------------------------------------------------
   Fixtures: one roster row per documented roster family (a no-argument
   control word and an argument-taking one), the send row's exact label
   shape, and one dormancy note. Transcribed from the wire shapes, runnable
   without Go.
   --------------------------------------------------------------------------- */

const ROSTER: PaletteRosterRow[] = [
  {
    name: "mute",
    summary: "静音或恢复悬浮球的语音应答",
    args: "",
    approval: "本地开关，无审批",
    door: "内置：ball.Mute 门",
  },
  {
    name: "fs.delete",
    summary: "删除工作区内的文件",
    args: "<path>",
    approval: "L2 强确认",
    door: "risk 门控的删除门",
  },
];

const NOTES: PaletteNote[] = [
  { word: "取消任务 未接线", reason: "C17 名册没有取消通路" },
  { word: "模式强制 未接线", reason: "D11③ 规则表未实现" },
];

const QUERY = "把桌面的临时文件清一清";

const GROUPS: PaletteGroup[] = [
  {
    name: "发送任务",
    items: [{ label: `以任务发送：${QUERY}`, icon: "terminal", hint: "Enter", action: "send" }],
  },
  { name: "面板", items: [{ label: "打开设置", icon: "settings", action: "settings" }] },
];

/* ---------------------------------------------------------------------------
   1. Fed values round-trip: the send row, the roster row facts, the notes
   --------------------------------------------------------------------------- */

const html = renderToStaticMarkup(
  <PaletteScreen query={QUERY} onQueryChange={() => undefined} groups={GROUPS} notes={NOTES} roster={ROSTER} />,
);

check("the send row must carry the query verbatim", html.includes(esc(`以任务发送：${QUERY}`)));
check("the send row must advertise the Enter key", html.includes(esc("Enter")));
check("the note's word must reach the screen", html.includes(esc("取消任务 未接线")));
check("the note's reason must reach the screen verbatim", html.includes(esc(NOTES[0]!.reason)));
check("the palette markup must hold zero emoji", !EMOJI_RE.test(html));

// The roster searches with the same query: a query that matches no command
// filters the rows out and says so in its own sentence (never a bare empty).
check(
  "a roster query matching no command must say so, not render rows",
  html.includes(esc("名册里没有匹配的命令")) && !html.includes(esc("mute")),
);

// The roster facts round-trip on their own render (empty query = the full
// menu): name, summary, argument shape, approval and door, all verbatim.
const rosterHtml = renderToStaticMarkup(
  <PaletteScreen query="" onQueryChange={() => undefined} groups={GROUPS} notes={NOTES} roster={ROSTER} />,
);
check("the roster row's name must reach the screen", rosterHtml.includes(esc("mute")));
check("the roster row's summary must reach the screen", rosterHtml.includes(esc(ROSTER[0]!.summary)));
check("the roster row's approval sentence must reach the screen verbatim", rosterHtml.includes(esc(ROSTER[0]!.approval)));
check("the second roster row's args shape must reach the screen", rosterHtml.includes(esc("<path>")));
check("the roster row's door sentence must reach the screen verbatim", rosterHtml.includes(esc(ROSTER[1]!.door)));

// The roster search also matches a summary: a query against the 删除 family
// keeps fs.delete and filters mute out.
const searchedHtml = renderToStaticMarkup(
  <PaletteScreen query="删除" onQueryChange={() => undefined} groups={[]} notes={[]} roster={ROSTER} />,
);
check("a summary-matching query keeps its row", searchedHtml.includes(esc("fs.delete")));
check("a summary-matching query filters the other rows", !searchedHtml.includes(esc("静音或恢复悬浮球的语音应答")));

/* ---------------------------------------------------------------------------
   3. The roster is read-only: no button may carry a roster row
   --------------------------------------------------------------------------- */

check(
  "roster rows must not render as buttons (the roster names doors, it opens none)",
  !/<button[^>]*data-roster-row/.test(rosterHtml),
);

/* ---------------------------------------------------------------------------
   4. The IME guard and the keyboard vocabulary (AC#2 / AC#6)
   --------------------------------------------------------------------------- */

const paletteSrc = readFileSync("src/components/palette-screen.tsx", "utf8");
check("the input must track composition start", paletteSrc.includes("onCompositionStart"));
check("the input must clear the guard on composition end", paletteSrc.includes("onCompositionEnd"));
check("the key handler must skip while composing (isComposing)", paletteSrc.includes("e.nativeEvent.isComposing"));
check("the palette must focus its input on mount (keyboard-first)", paletteSrc.includes("inputRef.current?.focus()"));
check("the palette must keep its arrow/enter vocabulary", paletteSrc.includes('"ArrowDown"') && paletteSrc.includes('"ArrowUp"') && paletteSrc.includes('"Enter"'));
check(
  "the palette must not reference task dims the wire does not carry (started/cost/step)",
  !paletteSrc.includes("costToDate") && !paletteSrc.includes("startedAt") && !paletteSrc.includes("currentStep"),
);

/* ---------------------------------------------------------------------------
   5. App wiring, at source level: the send route, the roster feed, the keys
   --------------------------------------------------------------------------- */

const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "the palette send must dispatch the composer's own route (lib/panel sendMessage, no second channel)",
  appSrc.includes("sendMessage(body, [])"),
);
check("App must feed the palette the snapshot's commands roster", appSrc.includes("snapshot.commands"));
check("App must hold the palette's ephemeral query as view preference", appSrc.includes("paletteQuery"));
check("App must render the send failure on the palette's own line", appSrc.includes("paletteError"));
check(
  "the palette must stay reachable by hotkey (Ctrl/Cmd+K)",
  appSrc.includes('(event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k"'),
);
check("Escape must close the palette", appSrc.includes("paletteOpen: false"));
// Re-adjudicated 2026-11-17 (consume wave 2): roster.cancel landed on the
// C17 roster (internal/panel/bridge.go's MethodRosterCancel, ticket 38
// AC#4's Go half, commit 3a62bed97), so the old dormancy pin ("取消任务
// 未接线" naming the missing engine) pins a fact that is no longer true.
// The note must now say the door LANDED, name it, and the App must bind the
// real door into the task monitor's cancel flow - a stronger pin, not a
// loosened one (the gate itself is pinned tooth-by-tooth by
// render:taskcard38).
check(
  "the cancel note must say the door landed and name it (roster.cancel; the monitor carries the real gate)",
  appSrc.includes("取消任务 已接线") &&
    appSrc.includes("roster.cancel") &&
    appSrc.includes("cancelDoors={{ cancelTask: requestRosterCancel, historyLoad: requestHistoryPage }}"),
);
check(
  "the force-rule dormancy note must name DEFERRED(D11-3) and its file",
  appSrc.includes("DEFERRED(D11-3)") && appSrc.includes("internal/agent/control.go"),
);
check(
  "the commands no-reader note must name the missing reader (PumpSources.Commands)",
  appSrc.includes("PumpSources.Commands"),
);

/* ---------------------------------------------------------------------------
   6. App still renders with the roster present and absent (both wire shapes)
   --------------------------------------------------------------------------- */

const COMPOSER = {
  mode: { current: "ask_high_risk", names: ["ask_every_step", "ask_high_risk"], l2ConfirmNames: ["auto_approve"] },
  workspace: { set: true, spelling: "D:/work/wisp", canonical: "D:/work/wisp", reparse: false, rewritten: false, reason: "" },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

function snapshotWith(commands?: PanelSnapshot["commands"]): PanelSnapshot {
  return {
    pending: [],
    results: [],
    composer: COMPOSER,
    generatedAt: "render-palette-38",
    ...(commands ? { commands } : {}),
  };
}

const withRosterHtml = renderToStaticMarkup(
  <App snapshot={snapshotWith({ rows: [{ name: "mute", summary: ROSTER[0]!.summary, args: "", approval: ROSTER[0]!.approval, door: ROSTER[0]!.door }] })} />,
);
check("App must still render the conversation with a commands key present", withRosterHtml.includes("一缕在这儿"));
check("the closed palette must not leak roster rows into the page", !withRosterHtml.includes(esc(ROSTER[0]!.approval)));

const noRosterHtml = renderToStaticMarkup(<App snapshot={snapshotWith()} />);
check("App must still render without a commands key (no-reader shape)", noRosterHtml.includes("一缕在这儿"));
check("the closed palette must not leak the dormancy notes into the page", !noRosterHtml.includes(esc("命令名册 未接线")));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-palette38: FAIL " + f);
  console.error(`render-palette38: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-palette38: OK, ticket 38's palette face consumes the pushed commands roster " +
    "(send row verbatim onto panel.message.send, roster facts round-trip read-only, " +
    "IME guard + autofocus + arrow/enter keys, cancel/force/roster dormancy notes name " +
    "their missing engines, no invented started/cost/step dims)",
);
