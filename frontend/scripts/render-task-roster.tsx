/* ============================================================================
   Task-roster evidence (ticket 219 AC#6's page half, 2026-11-08).
   ----------------------------------------------------------------------------
   What nothing else could see before this wave: the host has pushed a real
   TaskRosterSection since ticket 197 (and blockedOnApproval:true has fired in
   live packets since ticket 220, including with an EMPTY pending list - the L1
   window join), but App.tsx fed the task monitor `agents={[]}`, so every
   pushed row died unseen and 等批准 had no tier on the page at all.

   This instrument pins the page half of AC#6:

     1. the derivation (App's monitorAgents) copies the wire fields verbatim:
        an absent tasks key yields zero rows WITHOUT inventing one, and a
        present key yields one row per wire row - blockedOnApproval riding
        exactly as the host joined it;
     2. 等批准 is its own tier: a row with blockedOnApproval:true shows
        「等批准」 and NOT the D43 status word, whatever the row's status says
        (blocked outranks it - a task whose own card is on the queue is waiting
        on the user), and blockedOnApproval:false shows the D43 word;
     3. a row whose status was never filed (statusKnown:false, the wire's
        fail-closed shape) shows the HOST's own sentence, never a guessed word;
     4. no tasks key at all is an HONEST empty state: the popover says the
        roster never arrived (a different sentence from "nothing spawned"),
        and neither shape crashes - with or without the key, App renders.

   What this harness CANNOT see (stated so nobody over-reads the green):
   the popover mounts only on a click, so react-dom/server cannot reach it
   through App; the three teeth above render the real production popover
   component directly over the real derivation's output, and the App-level
   wiring is checked at source level (the `agents={monitorAgents(...)}` call
   and the rosterPresent flag), not by clicking.

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App, { monitorAgents } from "@/App";
import { MonitorPopover } from "@/components/harness/monitor-popover";
import type { PanelSnapshot, TaskRowView, TaskRosterSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
// FE0F rides as its own alternation branch: inside a character class it is a
// combining sequence and oxlint warns (no-misleading-character-class).
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
   Wire fixtures. Field names are lib/panel.ts TaskRowView / TaskRosterSection
   verbatim - the same names internal/panel/subagent_roster_197.go puts on the
   wire (the Go test reconciles the key sets both directions).
   --------------------------------------------------------------------------- */

const UNFILED_REASON = "宿主没有给出这一维的读数（fail-closed：不替任务编一个状态）";

function row(over: Partial<TaskRowView> & Pick<TaskRowView, "taskId" | "label">): TaskRowView {
  return {
    kind: "subagent",
    parentTaskId: "root-1",
    status: "",
    statusKnown: false,
    streamKey: over.taskId,
    blockedOnApproval: false,
    streamTruncated: false,
    streamElidedRunes: 0,
    streamDropped: false,
    ...over,
  };
}

function roster(rows: TaskRowView[]): TaskRosterSection {
  return {
    rows,
    inFlightSlots: 1,
    poolCap: 8,
    streamTruncated: false,
    streamElidedRunes: 0,
    droppedStreamKeys: [],
    droppedCount: 0,
  };
}

const COMPOSER = {
  mode: {
    current: "ask_high_risk",
    names: ["ask_every_step", "ask_high_risk", "auto_approve"],
    l2ConfirmNames: ["auto_approve"],
  },
  workspace: {
    set: true,
    spelling: "D:/work/wisp",
    canonical: "D:/work/wisp",
    reparse: false,
    rewritten: false,
    reason: "",
  },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

function snapshotWith(tasks: TaskRosterSection | undefined): PanelSnapshot {
  return {
    pending: [],
    results: [],
    composer: COMPOSER,
    generatedAt: "2026-11-08T14:00:00Z",
    tasks,
  };
}

const ENV = [
  { label: "工作区", value: "D:/work/wisp" },
  { label: "分支", value: "dev" },
] as const;

/* ---------------------------------------------------------------------------
   1. The derivation: wire fields copied verbatim, absent key stays absent.
   --------------------------------------------------------------------------- */

const fromAbsent = monitorAgents(undefined);
check("an absent tasks key must yield zero rows without inventing any", fromAbsent.length === 0);

const BLOCKED_ROW = row({
  taskId: "child-9",
  label: "等批准的子任务",
  status: "Acting",
  statusKnown: true,
  blockedOnApproval: true,
});
const fromPresent = monitorAgents(roster([BLOCKED_ROW]));
check(
  "a present tasks key must yield one derived row per wire row",
  fromPresent.length === 1,
  `rows=${fromPresent.length}`,
);
check(
  "the derived row must carry blockedOnApproval verbatim from the wire",
  fromPresent[0]?.blockedOnApproval === true,
);
check(
  "the derived row must carry the D43 status verbatim (this layer never judges it)",
  fromPresent[0]?.status === "Acting" && fromPresent[0]?.statusKnown === true,
);
check(
  "the derived row must carry the host's id and label verbatim",
  fromPresent[0]?.id === "child-9" && fromPresent[0]?.name === "等批准的子任务",
);

/* ---------------------------------------------------------------------------
   2. Tooth one: blockedOnApproval:true shows 等批准, not the status word.
   --------------------------------------------------------------------------- */

const blockedHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={monitorAgents(roster([BLOCKED_ROW]))} rosterPresent />,
);
check("the blocked row's own label must reach the screen", blockedHtml.includes(esc("等批准的子任务")));
check("the blocked row must show 等批准", blockedHtml.includes(esc("等批准")));
check(
  "the blocked row must NOT show the D43 status word (等批准 outranks it)",
  !blockedHtml.includes(esc("执行中")),
);
check("the blocked row's dot must be the warn tier", blockedHtml.includes("bg-warn"));
check("the blocked markup must hold zero emoji", !EMOJI_RE.test(blockedHtml));

/* ---------------------------------------------------------------------------
   3. Tooth two: blockedOnApproval:false shows the normal status word. Plus the
      fail-closed shape: a status the host never filed shows ITS sentence.
   --------------------------------------------------------------------------- */

const freeRow = row({
  taskId: "child-4",
  label: "没被卡住的子任务",
  status: "Acting",
  statusKnown: true,
  blockedOnApproval: false,
});
const freeHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={monitorAgents(roster([freeRow]))} rosterPresent />,
);
check("the free row must show the D43 status word", freeHtml.includes(esc("执行中")));
check("the free row must NOT show 等批准", !freeHtml.includes(esc("等批准")));
check("the free row's dot must be the live-work tier", freeHtml.includes("bg-accent"));
check("the free markup must hold zero emoji", !EMOJI_RE.test(freeHtml));

const unknownRow = row({
  taskId: "child-7",
  label: "状态没登记的子任务",
  status: "",
  statusKnown: false,
  blockedOnApproval: false,
  statusReason: UNFILED_REASON,
});
const unknownHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={monitorAgents(roster([unknownRow]))} rosterPresent />,
);
check(
  "an unfiled status must show the HOST's own sentence, never a guessed word",
  unknownHtml.includes(esc(UNFILED_REASON)),
  unknownHtml.slice(unknownHtml.indexOf("状态没登记") - 10, unknownHtml.indexOf("状态没登记") + 220),
);
check(
  "an unfiled status must not render any D43 word",
  !["执行中", "休眠", "收尾中", "出错"].some((w) => unknownHtml.includes(esc(w))),
);
check("the unfiled-status dot must be the idle tier", unknownHtml.includes("bg-ink-3"));

/* ---------------------------------------------------------------------------
   4. Tooth three: no tasks key at all = the honest empty state, no crash.
   --------------------------------------------------------------------------- */

const emptyRosterHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={[]} rosterPresent={false} />,
);
check(
  "with no tasks key the popover must say the roster never arrived",
  emptyRosterHtml.includes(esc("宿主没带任务这一节")),
);
check("the no-roster markup must hold zero emoji", !EMOJI_RE.test(emptyRosterHtml));

const spawnedEmptyHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={[]} rosterPresent />,
);
check(
  "with a present-but-empty roster the popover must say nothing is spawned (a different sentence)",
  spawnedEmptyHtml.includes(esc("还没有子任务在跑")) &&
    !spawnedEmptyHtml.includes(esc("宿主没带任务这一节")),
);

/* ---------------------------------------------------------------------------
   5. App with a pushed roster renders without crashing, and the wiring is at
      source level what static markup cannot click through.
   --------------------------------------------------------------------------- */

const appHtml = renderToStaticMarkup(<App snapshot={snapshotWith(roster([BLOCKED_ROW]))} />);
check("App with a pushed tasks section must render without crashing", appHtml.includes("一缕"));
check("the App markup must hold zero emoji", !EMOJI_RE.test(appHtml));

const appNoTasksHtml = renderToStaticMarkup(<App snapshot={snapshotWith(undefined)} />);
check("App without a tasks section must render without crashing", appNoTasksHtml.includes("一缕"));

const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "the monitor must be fed the derivation over snapshot.tasks, never a bare empty array",
  appSrc.includes("agents={monitorAgents(snapshot.tasks)}") && !appSrc.includes("agents={[]}"),
);
check(
  "the popover must be told whether the snapshot carried a tasks section at all",
  appSrc.includes("rosterPresent={snapshot.tasks !== undefined}"),
);
const popoverSrc = readFileSync("src/components/harness/monitor-popover.tsx", "utf8");
check(
  "等批准 must be derived from the wire's blockedOnApproval, not from any status word",
  popoverSrc.includes("if (row.blockedOnApproval) return { label: \"等批准\""),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-roster: FAIL " + f);
  console.error(`render-roster: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-roster: OK, the task monitor consumes snapshot.tasks " +
    "(等批准 is its own tier from blockedOnApproval, D43 words verbatim, " +
    "unfiled status shows the host's sentence, absent key is an honest empty state)",
);
