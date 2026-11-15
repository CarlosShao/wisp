/* ============================================================================
   Snapshot-receiver evidence (ticket 329 AC#5's page half, 2026-11-08).
   ----------------------------------------------------------------------------
   What nothing else can see: before this ticket the page never defined
   window.wispOnPanelSnapshot, so every whole snapshot the Go pump pushed was
   dropped by WebView2 and the panel sat on EMPTY forever. This instrument
   checks the page half of the contract in .scratch/wisp/probes/329/r1/
   page-contract.md:

     1. the receiver door (lib/panel.ts installPanelSnapshotReceiver): the
        named global is defined, the pushed argument lands as the WHOLE
        object (identity, never a string to parse), a second push REPLACES
        the package (no merge: a key the first package had and the second
        does not is gone), arbitrary repeated pushes are idempotent
        (latest wins), and the uninstaller removes the global without
        clobbering a reinstalled callback (StrictMode's order);
     2. App with NO push: renders EMPTY honestly - no picker, no invented
        family, no pending count, no crash;
     3. App with a pushed snapshot: the package's own content reaches the
        screen - the result text in the conversation, the workspace in the
        sidebar, and (via the snapshot-named 外观 category) the ball.style
        picker showing the pushed family, not the no-reader sentence.

   What this harness CANNOT see (stated so nobody over-reads the green):
   react-dom/server renders initial state only, so "host calls the global,
   React re-renders" is exercised here at the door level (tooth 1) plus as a
   source-level wiring check, not as a live re-render - the live loop is the
   real-machine acceptance (wisp.exe serving the rebuilt dist, push on show,
   ball.style row showing the live config's family).

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import {
  installPanelSnapshotReceiver,
  type PanelSnapshot,
  type PanelSnapshotReceiver,
} from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
// FE0F rides as its own alternation branch: inside a character class it is a
// combining sequence and oxlint warns (no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   The two whole packages. SECOND deliberately omits the tasks key FIRST
   carries - that omission is the whole-replacement probe.
   --------------------------------------------------------------------------- */

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
  acceptedAttachmentMimes: ["image/png"],
  maxAttachmentBytes: 10485760,
  attachmentError: "",
};

const FIRST: PanelSnapshot = {
  pending: [],
  results: [],
  composer: COMPOSER,
  generatedAt: "2026-11-08T10:00:00Z",
  ballStyle: { style: "glass", vocabulary: ["glass", "jarvis"] },
  tasks: {
    rows: [],
    inFlightSlots: 0,
    poolCap: 0,
    streamTruncated: false,
    streamElidedRunes: 0,
    droppedStreamKeys: [],
    droppedCount: 0,
  },
};

const SECOND: PanelSnapshot = {
  pending: [],
  results: [{ correlationId: "c1", text: "第二包推上来的正文。", done: true }],
  composer: COMPOSER,
  generatedAt: "2026-11-08T10:00:05Z",
  ballStyle: { style: "jarvis", vocabulary: ["glass", "jarvis"] },
};

/* ---------------------------------------------------------------------------
   1. The receiver door, exercised for real (a plain object stands in for
      Window; the door is the exact function App's root installs through).
   --------------------------------------------------------------------------- */

const target: { wispOnPanelSnapshot?: PanelSnapshotReceiver } = {};
let latest: PanelSnapshot | null = null;
const uninstall = installPanelSnapshotReceiver(target, (s) => {
  latest = s;
});
check(
  "installing the receiver must define the named global with the callback",
  target.wispOnPanelSnapshot !== undefined,
);

target.wispOnPanelSnapshot?.(FIRST);
check(
  "the first push must land as the whole object (identity, not a copy or a string)",
  latest === FIRST,
);
check(
  "the pushed argument must be the snapshot object, not a JSON string",
  latest !== null && typeof latest === "object" && latest.ballStyle?.style === "glass",
  `style=${latest?.ballStyle?.style}`,
);

target.wispOnPanelSnapshot?.(SECOND);
check("a second push must replace the whole snapshot (latest wins)", latest === SECOND);
check(
  "the replacement must not merge the previous package (SECOND carries no tasks key)",
  FIRST.tasks !== undefined && latest?.tasks === undefined,
  `tasks=${JSON.stringify(latest?.tasks)}`,
);

for (let i = 0; i < 5; i += 1) {
  target.wispOnPanelSnapshot?.(FIRST);
}
check(
  "arbitrary repeated pushes must be idempotent (the last whole package wins, no crash)",
  latest === FIRST,
);

uninstall();
check("uninstalling must remove the named global", target.wispOnPanelSnapshot === undefined);

// StrictMode order: mount installs A, cleanup A, mount installs B. A stale
// cleanup running after B's install must not kill B's receiver.
const cbA = (): void => undefined;
const cbB = (): void => undefined;
const strict: { wispOnPanelSnapshot?: PanelSnapshotReceiver } = {};
const unmountA = installPanelSnapshotReceiver(strict, cbA);
const unmountB = installPanelSnapshotReceiver(strict, cbB);
unmountA();
check(
  "a stale uninstaller must not clobber a reinstalled callback (StrictMode order)",
  strict.wispOnPanelSnapshot === cbB,
);
unmountB();
check("the live uninstaller must remove its own callback", strict.wispOnPanelSnapshot === undefined);

/* ---------------------------------------------------------------------------
   2. App with no push: EMPTY, honestly, without crashing.
   --------------------------------------------------------------------------- */

const emptyHtml = renderToStaticMarkup(<App />);
if (process.env.RCV_DUMP) {
  console.error(emptyHtml);
}
check("with no push the panel must render without crashing", emptyHtml.includes("一缕"));
check(
  "with no push the conversation must show the honest placeholder",
  emptyHtml.includes("一缕在这儿"),
);
check(
  "with no push the ball.style picker must not mount (no invented family)",
  !emptyHtml.includes('aria-label="球风格"'),
);
check("with no push there must be no pending count badge", !emptyHtml.includes("待审批"));
check("the empty markup must hold zero emoji", !EMOJI_RE.test(emptyHtml));

/* ---------------------------------------------------------------------------
   3. App with a pushed package. Two renderings: the chat screen (the result
      text arrives) and the snapshot-named 外观 category (the pushed
      ballStyle family reaches the picker - the row AC#5 is about).
   --------------------------------------------------------------------------- */

const chatHtml = renderToStaticMarkup(<App snapshot={SECOND} />);
check(
  "the pushed result text must reach the conversation verbatim",
  chatHtml.includes("第二包推上来的正文。"),
);
check(
  "the pushed workspace must reach the sidebar (whole package, not the EMPTY frame)",
  chatHtml.includes("D:/work/wisp") && !chatHtml.includes("尚未收到工作区"),
);
check("the pushed markup must hold zero emoji", !EMOJI_RE.test(chatHtml));

const configHtml = renderToStaticMarkup(
  <App snapshot={{ ...SECOND, view: "config" } as PanelSnapshot} />,
);
check(
  "the 外观 category must mount the ball.style picker for a pushed package",
  configHtml.includes('aria-label="球风格"'),
);
check(
  "the ball.style row must show the pushed family, not the no-reader sentence",
  configHtml.includes(">jarvis<") && !configHtml.includes("快照里还没有球风格这一节"),
  configHtml.slice(configHtml.indexOf("球风格") - 40, configHtml.indexOf("球风格") + 80),
);
check("the config markup must hold zero emoji", !EMOJI_RE.test(configHtml));

/* ---------------------------------------------------------------------------
   4. Source-level wiring (what static markup cannot reach): the mount point
      must install the receiver through the one shared door and pass the
      pushed object to state verbatim - no merge anywhere on the path.
   --------------------------------------------------------------------------- */

const rootSrc = readFileSync("src/panel-root.tsx", "utf8");
check(
  "the panel root must install the receiver on window through the shared door",
  rootSrc.includes("installPanelSnapshotReceiver(window, setSnapshot)"),
);
check(
  "the pushed object must reach state verbatim (whole replacement, no merge on the path)",
  !rootSrc.includes("Object.assign") && !/\{\s*\.\.\.\w+\s*,\s*\.\.\.\s*\w+\s*\}/.test(rootSrc),
);
check(
  "the mount point must render the panel through PanelRoot, not a bare <App />",
  readFileSync("src/main.tsx", "utf8").includes("<PanelRoot />"),
);
const libSrc = readFileSync("src/lib/panel.ts", "utf8");
check(
  "the named global must be declared on Window in the contract layer",
  libSrc.includes("wispOnPanelSnapshot?: PanelSnapshotReceiver"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-receiver: FAIL " + f);
  console.error(`render-receiver: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-receiver: OK, window.wispOnPanelSnapshot defined via the shared door " +
    "(whole-object pushes, latest wins, no merge, strict-safe uninstall), " +
    "EMPTY honest without a push, pushed package reaches chat + ball.style row",
);
