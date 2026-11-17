/* ============================================================================
   Ticket 37's behavioural tooth (render evidence), 2026-11-16.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-l2 / render-task-roster /
   render-history-36 already established for this tree: every check below
   feeds REAL-SHAPED data through the PRODUCTION component and asserts the
   SAME VALUE comes back out in the markup - "喂真数据⇒显示同一份值".

   The fixtures are the wire's own documented shapes:
     - ApprovalCardView verbatim off internal/panel/approval.go:38-59 (the
       JSON keys TestApprovalCardViewJSONKeysMatchFrontendTypes reconciles);
     - TaskRowView / TaskRosterSection verbatim off
       internal/panel/subagent_roster_197.go:111-184 (the keys
       TestComposerContractTypesMatchFrontend reconciles);
     - the CSV rows from the ticket 48 host test's join
       (cmd/wisp/approval_queue_48_windows_test.go:229 Test48PanelRosterJoin-
       DepthSyncUnder100ms) transcribed into those two wire shapes - the join
       is the host's own (indexWaitingKey), so the page reading is the packet
       itself, and the transcript keeps this file runnable without Go.

   Where a dimension has NO wire field, the check asserts its ABSENCE - a
   renderer that invented a value for an undeclared key would be the exact
   假值 the dispatch forbids. The dimensions and where they now stand:
     - the batch display fields (readface wave 1, 2026-11-17): the five flat
       ApprovalCardView keys (batchTotal/batchDirs/batchFirst/batchDeferred/
       batchSummary, internal/panel/approval.go:81-85, produced by
       approval.CardSummary) now RIDE the card wire, so a fed card must draw
       them verbatim and ONE card only, while a plain card (no batch keys -
       omitempty) must draw no batch chrome from zeros;
     - the 270s/300s deadline (gate.go's deadline, EventWarning.Remaining)
       is not on the card wire either - the countdown lives on the native
       card (cardrender.go) and the ball, never on this page;
     - the head's 1-based Depth stamp rides Prompt.Depth (the ball's badge
       channel, approvalDepthBadge), not the packet.

   Placement is checked at SOURCE level where static markup cannot click:
   App renders the strip above the cards on every view (the same rule
   render-task-roster uses for App-level wiring).
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import { ApprovalQueueSection } from "@/components/approval-queue";
import { ApprovalScreen, QUEUE_CAPACITY } from "@/components/approval-screen";
import type { ApprovalCardView, PanelSnapshot, TaskRosterSection, TaskRowView } from "@/lib/panel";

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
   Fixtures: two L2 cards whose correlation ids, tools, args, rules and
   reasons are the wire's own vocabulary (fs.delete / shell.exec L2 cards as
   the assessor writes them - see fixtures/l2-card-fs-delete.json), plus the
   task roster the ticket-48 join needs.
   --------------------------------------------------------------------------- */

const CARD_A: ApprovalCardView = {
  correlationId: "t48-join-a#call_1",
  tool: "fs.delete",
  args: ["--path", "C:/Users/swq/Documents/notes.md", "--permanent"],
  level: "L2",
  rulesHit: ["R1", "R8"],
  reason: "R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除）",
  reasonKnown: true,
  sessionOverrideBlocked: false,
  callChain: [],
  decidedBy: "native",
};

const CARD_B: ApprovalCardView = {
  correlationId: "t48-join-b#call_2",
  tool: "shell.exec",
  args: ["powershell", "-NoProfile", "-File", "clean.ps1"],
  level: "L2",
  rulesHit: ["R6"],
  reason: "脚本执行属不可逆操作，需逐条确认。",
  reasonKnown: true,
  sessionOverrideBlocked: true,
  callChain: [],
  decidedBy: "native",
};

/** The wire's TaskRowView, defaulted the way the host sends a healthy row. */
function row(over: Partial<TaskRowView> & Pick<TaskRowView, "taskId" | "label">): TaskRowView {
  return {
    kind: "subagent",
    parentTaskId: "root-1",
    status: "Acting",
    statusKnown: true,
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
    inFlightSlots: 2,
    poolCap: 8,
    streamTruncated: false,
    streamElidedRunes: 0,
    droppedStreamKeys: [],
    droppedCount: 0,
  };
}

const COMPOSER = {
  mode: { current: "ask_high_risk", names: ["ask_every_step", "ask_high_risk", "auto_approve"], l2ConfirmNames: ["auto_approve"] },
  workspace: { set: true, spelling: "D:/work/wisp", canonical: "D:/work/wisp", reparse: false, rewritten: false, reason: "" },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

function snapshotWith(pending: ApprovalCardView[], tasks?: TaskRosterSection): PanelSnapshot {
  return { pending, results: [], composer: COMPOSER, generatedAt: "render-approval-37", ...(tasks ? { tasks } : {}) };
}

// The ticket-48 join, transcribed: the host matched card A's corr onto task
// task-a (the corr carries the task id, indexWaitingKey's family), and task-b
// joined through its own corr prefix (the "b#2" stacked shape's task half).
const JOINED_ROSTER = roster([
  row({ taskId: "t48-join-a", label: "整理桌面的子任务", blockedOnApproval: true }),
  row({ taskId: "t48-join-b", label: "清理临时目录的子任务", blockedOnApproval: true }),
  row({ taskId: "t48-free", label: "没被卡的子任务", blockedOnApproval: false }),
]);

/* ---------------------------------------------------------------------------
   1. The card: correlationId on the card (the ticket's forensics row)
   --------------------------------------------------------------------------- */

const singleHtml = renderToStaticMarkup(<App snapshot={snapshotWith([CARD_A, CARD_B], JOINED_ROSTER)} />);
check("the card must carry the correlationId verbatim (forensics)", singleHtml.includes(esc(`corr ${CARD_A.correlationId}`)));
check("the full corr id must sit in a title attribute (hover forensics)", singleHtml.includes(esc(`correlationId：${CARD_A.correlationId}`)));
check("the queue strip must name the head card's corr", singleHtml.includes(esc(CARD_A.correlationId)));
check("the card's rule ids must reach the screen", CARD_A.rulesHit.every((r) => singleHtml.includes(esc(r))));
check("the card's argv must reach the screen", CARD_A.args.every((a) => singleHtml.includes(esc(a))));
check("the card's reason must reach the screen", singleHtml.includes(esc(CARD_A.reason)));
check("the App markup must hold zero emoji", !EMOJI_RE.test(singleHtml));

/* ---------------------------------------------------------------------------
   2. The queue strip: depth + 队头, on every view, from the pushed list
   --------------------------------------------------------------------------- */

check("the strip must show the depth the wire carries (2)", singleHtml.includes(esc("深度 2")));
check("the strip must mark the FIFO head", singleHtml.includes(esc("队头")));
check("the strip must not mark a non-head row as head", !singleHtml.includes(esc("队头") + esc(CARD_B.tool)));
// Strip presence is the App-level wiring: it must sit in the source above the
// cards (the L2 cards' own comment names the order), on every view.
const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "App must mount ApprovalQueueSection from the snapshot's pending list",
  appSrc.includes("snapshot.pending.length > 0 ? (\n                <ApprovalQueueSection snapshot={snapshot} />"),
);
check(
  "the strip must be mounted WITH the cards, not as a separate view",
  appSrc.indexOf("ApprovalQueueSection snapshot={snapshot}") < appSrc.indexOf("snapshot.pending.map((v) => ("),
);

// Empty queue: zero strip markup at all (no fake 深度 0 furniture).
const emptyHtml = renderToStaticMarkup(<App snapshot={snapshotWith([], JOINED_ROSTER)} />);
check("an empty queue must render no strip", !emptyHtml.includes('aria-label="审批队列"'));
check("an empty queue must not fake a depth", !emptyHtml.includes(esc("深度 0")));

/* ---------------------------------------------------------------------------
   3. 等待于: the roster join, verbatim off the packet
   --------------------------------------------------------------------------- */

check("the joined row must name its task label verbatim", singleHtml.includes(esc("等待于：整理桌面的子任务")));
check("the second joined row must name its own task", singleHtml.includes(esc("等待于：清理临时目录的子任务")), singleHtml.slice(singleHtml.indexOf("t48-join-b") - 10, singleHtml.indexOf("t48-join-b") + 240));
// The free row's task is NOT waiting on anything and must not be claimed.
check("a task without a card must not be claimed as waiting", !singleHtml.includes(esc("等待于：没被卡的子任务")));

// No tasks key = no reader for the task dimension. The row says so; it does
// not guess a task name from the corr prefix (the corr-carries-task-id
// fallback is exactly the coincidence ticket 37's 对账 row names - a page
// that "derives" a task from the corr text would be judging a join the host
// owns).
const noRosterHtml = renderToStaticMarkup(
  <ApprovalQueueSection snapshot={snapshotWith([CARD_A])} />,
);
check(
  "without a tasks key the row must say the roster never arrived",
  noRosterHtml.includes(esc("等待于：快照未带任务名册")),
);
check(
  "without a tasks key no task name may render",
  !noRosterHtml.includes(esc("整理桌面的子任务")),
);

// A roster present but nobody claiming this card: the honest third sentence.
const CLAIMLESS_ROSTER = roster([row({ taskId: "other-task", label: "别人的任务", blockedOnApproval: false })]);
const claimlessHtml = renderToStaticMarkup(
  <ApprovalQueueSection snapshot={snapshotWith([CARD_A], CLAIMLESS_ROSTER)} />,
);
check(
  "an unclaimed card must say no task claims it, never a guessed name",
  claimlessHtml.includes(esc("等待于：名册里没有认领这张卡的任务")),
);

/* ---------------------------------------------------------------------------
   4. The approval screen: FIFO numbering, head note, capacity == constant
   --------------------------------------------------------------------------- */

const screenHtml = renderToStaticMarkup(<ApprovalScreen snapshot={snapshotWith([CARD_A, CARD_B], JOINED_ROSTER)} />);
check("the screen must number rows in FIFO order", screenHtml.includes(esc("审批队列第 1 项")) && screenHtml.includes(esc("审批队列第 2 项")));
check("the first row must be card A's tool (wire order preserved)", screenHtml.indexOf(esc(CARD_A.tool)) < screenHtml.indexOf(esc(CARD_B.tool)));
check("the screen must say the head is what the ball answers", screenHtml.includes(esc("队头在先（FIFO），悬浮球只答队头")));
check(
  "the capacity sentence must carry the queue's own bound (DefaultMaxPending = 8)",
  screenHtml.includes(esc(`队列容量上限 ${QUEUE_CAPACITY} 条`)) && QUEUE_CAPACITY === 8,
);
check("the screen markup must hold zero emoji", !EMOJI_RE.test(screenHtml));

const emptyScreenHtml = renderToStaticMarkup(<ApprovalScreen snapshot={snapshotWith([])} />);
check("the empty queue must be its own sentence, not a blank list", emptyScreenHtml.includes(esc("队列当前为空")));

/* ---------------------------------------------------------------------------
   4b. The aggregate card (ticket 37 AC#5, readface wave 1): a many-target
       call renders ONE card wearing the batch's display fields, verbatim
   --------------------------------------------------------------------------- */

// The wire's own batch shape, transcribed off readface_wave1_test.go's
// NativeVerdict.Batch fixture (10 targets, 2 dirs, preview rows, deferred
// remainder) and Go's CardSummary sentence format (cardsummary_37.go:67-76).
const BATCH_CARD: ApprovalCardView = {
  correlationId: "t48-join-batch#call_1",
  tool: "fs.write",
  args: ["--batch", "C:/data/cleanup.txt"],
  level: "L2",
  rulesHit: ["R7"],
  reason: "本批 10 项同类写入，超过 R7 的批量阈值，升级逐项确认。",
  reasonKnown: true,
  sessionOverrideBlocked: false,
  callChain: [],
  decidedBy: "native",
  batchTotal: 10,
  batchDirs: ["C:/data/x", "C:/data/y"],
  batchFirst: ["C:/data/x/a.txt", "C:/data/x/b.txt", "C:/data/x/c.txt"],
  batchDeferred: 7,
  batchSummary: "本次调用包含 10 项同类 L2 操作，涉及 2 个目录：C:/data/x、C:/data/y；前 3 项：C:/data/x/a.txt、C:/data/x/b.txt、C:/data/x/c.txt（另有 7 项可展开查看）",
};

const batchHtml = renderToStaticMarkup(<App snapshot={snapshotWith([BATCH_CARD])} />);
check(
  "10 操作 must stay ONE card (one 需要你的确认 head, one card body)",
  (batchHtml.match(/需要你的确认/g) ?? []).length === 1,
  `${(batchHtml.match(/需要你的确认/g) ?? []).length} card heads`,
);
check("the card must carry batchSummary verbatim (Go's own words, unparaphrased)", batchHtml.includes(esc(BATCH_CARD.batchSummary)));
check("the card must carry the batch total verbatim", batchHtml.includes(esc("本批共 10 项")));
check("the card must carry the deferred remainder verbatim", batchHtml.includes(esc("7 项未在卡上展开")));
check("the card must carry every affected dir verbatim", BATCH_CARD.batchDirs.every((d) => batchHtml.includes(esc(d))));
check("the card must carry every preview path verbatim", BATCH_CARD.batchFirst.every((p) => batchHtml.includes(esc(p))));
check("the batch block must hold zero emoji", !EMOJI_RE.test(batchHtml));
check(
  "the batch block must add no approval-shaped control (the native confirm covers the whole batch)",
  (() => {
    const names: string[] = [];
    for (const m of batchHtml.matchAll(/<button\b([^>]*)>([\s\S]*?)<\/button>/g)) {
      const aria = m[1]?.match(/aria-label="([^"]*)"/)?.[1] ?? "";
      const text = (m[2] ?? "").replace(/<[^>]*>/g, "").replace(/&quot;/g, '"').trim();
      names.push(`${text} ${aria}`.trim());
    }
    return names.length > 0 && names.every((n) => !/允许|同意|批准|通过|approve|allow|grant/i.test(n));
  })(),
);
check(
  "the plain cards above must render NO batch chrome (no batch keys on their wire, no zeros dressed up)",
  !singleHtml.includes(esc("本批共")) && !singleHtml.includes(esc("涉及目录")) && !screenHtml.includes(esc("本批共")),
);
check(
  "panel.ts must declare the five batch keys (the Go judged-gap ledger goes stale otherwise)",
  readFileSync("src/lib/panel.ts", "utf8").includes("batchTotal") &&
    readFileSync("src/lib/panel.ts", "utf8").includes("batchSummary"),
);
check(
  "the Go ledger's ApprovalCardView batch-key registration must be gone the same commit (the both-ways rule; the pairs-table name stays)",
  !/ApprovalCardView[^\n]*batchTotal/.test(readFileSync("../internal/panel/composer_test.go", "utf8")),
);

/* ---------------------------------------------------------------------------
   5. The Go-only dimensions: named absences, never invented values
   --------------------------------------------------------------------------- */

// The 270s/300s ladder: the wire carries no deadline, the countdown is the
// native card's. A countdown on this page would be an invented number.
check(
  "no invented deadline may render (the wire carries no timeout field)",
  !screenHtml.includes(esc("270")) && !screenHtml.includes(esc("300")) && !singleHtml.includes(esc("剩余")),
);
// The head's depth badge (Prompt.Depth -> ball badge) is the native channel;
// the page derives depth from the list length only.
check(
  "the page must not draw a per-card depth stamp (the badge is the ball's channel)",
  !singleHtml.includes(esc("深度 1")),
);

/* ---------------------------------------------------------------------------
   6. The no-allow invariant, re-asserted for the new surfaces
   --------------------------------------------------------------------------- */

// render-l2 asserts the card; the strip and the screen are new buttons-free
// surfaces, so their rendered markup must offer no allow-shaped control.
const stripHtml = renderToStaticMarkup(<ApprovalQueueSection snapshot={snapshotWith([CARD_A], JOINED_ROSTER)} />);
for (const html of [stripHtml, screenHtml]) {
  const names: string[] = [];
  for (const m of html.matchAll(/<button\b([^>]*)>([\s\S]*?)<\/button>/g)) {
    const aria = m[1]?.match(/aria-label="([^"]*)"/)?.[1] ?? "";
    const text = (m[2] ?? "").replace(/<[^>]*>/g, "").replace(/&quot;/g, '"').trim();
    names.push(`${text} ${aria}`.trim());
  }
  check(
    "the queue surfaces must offer no approval-shaped control (panel rejects; the ball allows)",
    names.every((n) => !/允许|同意|批准|通过|approve|allow|grant/i.test(n)),
    names.join(" | "),
  );
  check("the queue surface must hold zero emoji", !EMOJI_RE.test(html));
}
check(
  "the approval screen must not wire any decision route (the card is the only sender)",
  !readFileSync("src/components/approval-screen.tsx", "utf8").includes("requestApprovalResolution") &&
    !readFileSync("src/components/approval-queue.tsx", "utf8").includes("requestApprovalResolution"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-approval37: FAIL " + f);
  console.error(`render-approval37: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-approval37: OK, ticket 37's panel faces consume the pushed queue " +
    "(corr on the card, depth + FIFO head strip on every view, 等待于 off the roster join, " +
    "capacity = the queue's own 8, the aggregate card: 10 ops = ONE card wearing the five batch fields verbatim with no allow-shaped control, plain cards render no batch chrome, " +
    "deadline + depth stamps stay named absences, no allow-shaped control)",
);
